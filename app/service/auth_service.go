package service

// TODO
// 1. Create JWTManager Helper
// 2. Make a Process for: Load and Set Common Password to Memory
// 3. Create App Validator and Implement the common password

// TODO
// *Create Auth Service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/app/model"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/app/repository"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/helper"
	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	users         repository.UserRepository
	tokens        repository.TokenRepository
	jwt           *helper.JWTManager
	refreshTTL    time.Duration
	appValidator  *helper.AppValidator
	permissionSet *helper.PermissionSet
}

func NewAuthService(
	users repository.UserRepository,
	tokens repository.TokenRepository,
	jwt *helper.JWTManager,
	refreshTTL time.Duration,
	appValidator *helper.AppValidator,
) *AuthService {
	return &AuthService{
		users:        users,
		tokens:       tokens,
		jwt:          jwt,
		refreshTTL:   refreshTTL,
		appValidator: appValidator,
	}
}

const refreshTokenByte = 32

// --- METHODS ---
// func (s *AuthService) Register() error :
//
//	 POST /auth/register
//		ambil ctx, check body JSON pakai c.BodyParser, trimspace username email, check validasi, hash password, panggil s.users.Create(), kembalikan helper.Created
func (s *AuthService) Register(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RegisterReq
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrBadRequest("BAD_REQUEST_AUTHENTICATION", "JSON invalid!")
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if errs := helper.ValidateStruct(req, *s.appValidator); len(errs) > 0 {
		vals := helper.ToValidationFields(errs)
		return helper.ErrValidation("Validasi Register", vals)
	}

	hash_password, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.ErrInternal(err)
	}

	user, err := s.users.Create(ctx, model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hash_password,
		Role:     "user",
		IsActive: true,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.ErrConflict("username telah terdaftar")
		}

		return helper.ErrInternal(err)
	}

	return helper.Created(c, "user berhasil didaftarkan", user, "/api/v1/users/"+strconv.Itoa(user.ID))
}

// func (s *AuthService) Login() error :
//
//	 POST /auth/loign
//		ambil ctx, cek body JSON, validasi, panggil s.users.FindByUsername(), helper.VerifyPassword, check is active, buat token pair, kembalikan helper.Success bersama token pair
func (s *AuthService) Login(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.LoginReq
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrBadRequest("BAD_REQUEST_AUTHENTICATION", "JSON invalid!")
	}

	if errs := helper.ValidateStruct(req, *s.appValidator); len(errs) > 0 {
		vals := helper.ToValidationFields(errs)
		return helper.ErrValidation("Validasi Register", vals)
	}

	user, err := s.users.FindByUsername(ctx, req.Username)
	if err != nil {
		helper.VerifyDummyPassword(req.Password)
		return helper.ErrUnauthorized("password / username salah")
	}

	if !helper.VerifyPassword(user.Password, req.Password) {
		return helper.ErrUnauthorized("password / username salah")
	}

	if !user.IsActive {
		return helper.ErrForbidden("akun dinonaktifkan")
	}

	// create pair
	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return helper.ErrInternal(err)
	}

	return helper.Success(c, "login berhasil", pair)
}

// func (s *AuthService) Refresh() error :
//
//	 POST /auth/refresh
//		ambil ctx, cek body JSON, cek string token, panggil findActive() pakai token, ambil user, Revoke(), issueTokenPair(), return helper.Success
func (s *AuthService) Refresh(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RefreshReq
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrBadRequest("BAD_REQUEST_AUTHENTICATION", "JSON invalid!")
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		return helper.ErrBadRequest("BAD_REQUEST_AUTHENTICATION", "refresh token kosong")
	}

	hash := helper.SHA256Hex(req.RefreshToken)
	refresh_token, err := s.tokens.FindActive(ctx, hash)
	if err != nil {
		return helper.ErrUnauthorized("token refresh tidak valid / kedaluarwsa")
	}

	user, err := s.users.FindByID(ctx, refresh_token.UserID)
	if err != nil || !user.IsActive {
		return helper.ErrUnauthorized("akun tidak bisa dipakai")
	}

	if err := s.tokens.Revoke(ctx, refresh_token.TokenHash); err != nil {
		return helper.ErrInternal(err)
	}

	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return helper.ErrInternal(err)
	}

	return helper.Success(c, "berhasil membuat token baru", pair)
}

// func (s *AuthService) Logout() error :
//
//	 POST /auth/logout
//		ambil ctx, cek body JSON, cek refreshToken, Revoke Token, return helper.Success
func (s *AuthService) Logout(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RefreshReq
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrBadRequest("BAD_REQUEST_AUTHENTICATION", "JSON invalid!")
	}

	if strings.TrimSpace(req.RefreshToken) != "" {
		_ = s.tokens.Revoke(ctx, helper.SHA256Hex(req.RefreshToken))
	}

	return helper.Success(c, "berhasil logout", nil)
}

// func (s *AuthService)  Me() error :
//
//	 GET /auth/me
//		ambil ctx, ambil current user pakai helper function, return helper.success dengan user
func (s *AuthService) Me(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	authUser, ok := helper.AuthUserFromFiber(c)
	if !ok {
		return helper.ErrUnauthorized("belum terautentikasi")
	}

	user, err := s.users.FindByID(ctx, authUser.UserID)
	if err != nil {
		return helper.ErrUnauthorized("user tidak ditemukan")
	}

	return helper.Success(c, "user berhasil ditemukan", fiber.Map{
		"user": user,
		// "permissions": ISI PERMISSION YANG DIMILIKI,
	})
}

// func (s *AuthService) issueTokenPair(*fiber.Ctx, model.User) (model.TokenPair, error) :
//
//	generate accesstoken jwt.generateAccess(), generate refresh token use helper function, token.Save ini ngesave refresh token, return model.TokenPair
func (s *AuthService) issueTokenPair(ctx context.Context, user model.User) (model.TokenPair, error) {
	acces_token, err := s.jwt.GenerateAccessToken(user)
	if err != nil {
		return model.TokenPair{}, err
	}

	refresh_token, err := helper.RandomToken(refreshTokenByte)
	if err != nil {
		return model.TokenPair{}, err
	}

	if err := s.tokens.Save(ctx, model.RefreshToken{
		UserID:    user.ID,
		TokenHash: helper.SHA256Hex(refresh_token),
		ExpiredAt: time.Now().Add(s.refreshTTL),
	}); err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{
		AccessToken:  acces_token,
		RefreshToken: refresh_token,
		TokenType:    "Bearer",
		ExpiredIn:    int(s.refreshTTL.Seconds()),
	}, nil
}
