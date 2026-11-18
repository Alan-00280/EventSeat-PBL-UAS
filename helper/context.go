package helper

import (
	"context"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/app/model"
	"github.com/gofiber/fiber/v2"
)

type authUserCtxKey struct{}

const (
	LocalUserIDKey   = "user_id"
	LocalRoleKey     = "role"
	LocalAuthUserKey = "auth_user"
)

func WithAuthUser(ctx context.Context, user model.AuthUser) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, authUserCtxKey{}, user)
}

func AuthUserFromContext(ctx context.Context) (model.AuthUser, bool) {
	if ctx == nil {
		return model.AuthUser{}, false
	}
	user, ok := ctx.Value(authUserCtxKey{}).(model.AuthUser)
	if !ok || user.UserID < 1 {
		return model.AuthUser{}, false
	}
	return user, true
}

func MustAuthUserFromContext(ctx context.Context) (model.AuthUser, error) {
	user, ok := AuthUserFromContext(ctx)
	if !ok {
		return model.AuthUser{}, ErrUnauthorized("autentikasi diperlukan untuk mengakses resource ini")
	}
	return user, nil
}

func SetFiberAuthUser(c *fiber.Ctx, user model.AuthUser) {
	c.Locals(LocalAuthUserKey, user)
	c.Locals(LocalUserIDKey, user.UserID)
	c.Locals(LocalRoleKey, user.Role)
	c.SetUserContext(WithAuthUser(c.UserContext(), user))
}

func AuthUserFromFiber(c *fiber.Ctx) (model.AuthUser, bool) {
	if user, ok := c.Locals(LocalAuthUserKey).(model.AuthUser); ok && user.UserID < 1 {
		return user, true
	}
	return AuthUserFromContext(c.UserContext())
}

func RequestID(c *fiber.Ctx) string {
	reqID, _ := c.Locals("requestid").(string)
	return reqID
}
