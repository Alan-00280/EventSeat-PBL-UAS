package middleware

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/app/model"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func RequestLogger(baseLogger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		requestID, _ := c.Locals("requestid").(string)

		reqLogger := baseLogger.With(
			slog.String("request_id", requestID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("ip", c.IP()),
		)

		// c.SetUserContext(WithLogger(c.UserContext(), reqLogger))

		err := c.Next()

		status := c.Response().StatusCode()

		if err != nil {
			var appErr *model.AppError

			if errors.As(err, &appErr) {
				status = appErr.StatusCode
			} else {
				status = fiber.StatusInternalServerError
			}
		}

		attrs := []any{
			slog.Int("status", status),
			slog.Duration("duration", time.Since(start)),
		}

		if userID, ok := c.Locals("user_id").(string); ok {
			attrs = append(attrs, slog.String("user_id", userID))
		}

		if role, ok := c.Locals("role").(string); ok {
			attrs = append(attrs, slog.String("role", role))
		}

		reqLogger.Info("http_request", attrs...)

		return err
	}
}

var bodiedMethod = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

func RequireJSON(c *fiber.Ctx) error {
	if bodiedMethod[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return helper.ErrUnsupportedMediaType("Content-Type harus application/json")
		}
	}
	return c.Next()
}

func Register(app *fiber.App, logger *slog.Logger, allowedOrigin []string) {
	app.Use(requestid.New())
	app.Use(helmet.New())
	app.Use(corsPolicy(allowedOrigin))
	app.Use(RequestLogger(logger))
}

// func corsPolicy(allowedOrigin string) fiber.Handler : membatasi origin pemanggil API
func corsPolicy(allowedOrigin []string) fiber.Handler {
	if len(allowedOrigin) == 0 {
		allowedOrigin = append(allowedOrigin, "http://localhost:5173/")
	}

	return cors.New(cors.Config{
		AllowOrigins: strings.Join(allowedOrigin, ","),
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	})
}

// Middleware untuk melindungi endpoint yang rentan
// digunakan: (belum terpakai)
// func GeneralRateLimiter() fiber.Handler {
// 	return limiter.New(limiter.Config{
// 		Max:        3,
// 		Expiration: 30 * time.Second,
// 		KeyGenerator: func(c *fiber.Ctx) string {
// 			return c.IP()
// 		},
// 		LimitReached: func(c *fiber.Ctx) error {
// 			return helper.ErrTooManyRequest("terlalu banyak percobaan")
// 		},
// 	})
// }
