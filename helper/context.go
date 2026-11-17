package helper

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type authUserCtxKey struct{}

const (
	LocalUserIDKey = "user_id"
	LocalRoleKey = "role"
	LocalAuthUserKey = "auth_user"
)

type AuthUser struct {
	UserID string
	Role   string
	Email  string
}

func WithAuthUser(ctx context.Context, user AuthUser) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, authUserCtxKey{}, user)
}

func AuthUserFromContext(ctx context.Context) (AuthUser, bool) {
	if ctx == nil {
		return AuthUser{}, false
	}
	user, ok := ctx.Value(authUserCtxKey{}).(AuthUser)
	if !ok || user.UserID == "" {
		return AuthUser{}, false
	}
	return user, true
}

func MustAuthUserFromContext(ctx context.Context) (AuthUser, error) {
	user, ok := AuthUserFromContext(ctx)
	if !ok {
		return AuthUser{}, ErrUnauthorized("autentikasi diperlukan untuk mengakses resource ini")
	}
	return user, nil
}

func SetFiberAuthUser(c *fiber.Ctx, user AuthUser) {
	c.Locals(LocalAuthUserKey, user)
	c.Locals(LocalUserIDKey, user.UserID)
	c.Locals(LocalRoleKey, user.Role)
	c.SetUserContext(WithAuthUser(c.UserContext(), user))
}

func AuthUserFromFiber(c *fiber.Ctx) (AuthUser, bool) {
	if user, ok := c.Locals(LocalAuthUserKey).(AuthUser); ok && user.UserID != "" {
		return user, true
	}
	return AuthUserFromContext(c.UserContext())
}