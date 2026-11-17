package config

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/helper"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/middleware"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hibiken/asynq"
)

const (
	MaxRequestBodySize = 1 * 1024 * 1024
)

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success   bool          `json:"success"`
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Fields    []ErrorDetail `json:"fields,omitempty"`
	RequestID string        `json:"request_id"`
}

func NewFiberApp(env *Env, logger *slog.Logger, deps routes.Dependencies) *fiber.App {
	if logger == nil {
		logger = slog.Default()
	}

	app := fiber.New(fiber.Config{
		AppName:               env.AppName,
		BodyLimit:             MaxRequestBodySize,
		DisableStartupMessage: env.AppEnv == "production",
		ReadTimeout:           15 * time.Second,
		WriteTimeout:          15 * time.Second,
		IdleTimeout:           60 * time.Second,
		ErrorHandler:          newErrorHandler(logger),
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: env.AppEnv != "production",
	}))

	// HELLO WORLD!
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Middleware
	middleware.Register(app, logger, env.AllowedOrigins)

	// App Routes
	routes.Register(app, deps)

	// 404
	app.Use(func(c *fiber.Ctx) error {
		return helper.ErrNotFound("404 - routes not found")
	})

	return app
}

func accessLoggerMiddleware(baseLogger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		reqID, _ := c.Locals("requestid").(string)

		reqLogger := baseLogger.With(
			slog.String("request_id", reqID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("ip", c.IP()),
		)
		c.SetUserContext(WithLogger(c.UserContext(), reqLogger))

		err := c.Next()

		userID, _ := c.Locals("user_id").(string)
		role, _ := c.Locals("role").(string)

		status := resolveHTTPStatus(c, err)

		reqLogger.Info("http_request",
			slog.Int("status", status),
			slog.Int64("latency_ms", time.Since(start).Milliseconds()),
			slog.String("user_id", userID),
			slog.String("role", role),
		)

		return err
	}
}

func resolveHTTPStatus(c *fiber.Ctx, err error) int {
	if err == nil {
		return c.Response().StatusCode()
	}

	var appErr *helper.AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return fiberErr.Code
	}

	return http.StatusInternalServerError
}

func newErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		reqID, _ := c.Locals("requestid").(string)

		var appErr *helper.AppError
		if errors.As(err, &appErr) {
			var fields []ErrorDetail
			if len(appErr.Fields) > 0 {
				fields = make([]ErrorDetail, 0, len(appErr.Fields))
				for _, f := range appErr.Fields {
					fields = append(fields, ErrorDetail{
						Field:   f.Field,
						Message: f.Message,
					})
				}
			}

			return c.Status(appErr.StatusCode).JSON(ErrorResponse{
				Success:   false,
				Code:      appErr.Code,
				Message:   appErr.Message,
				Fields:    fields,
				RequestID: reqID,
			})
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code := httpStatusToCode(fiberErr.Code)
			msg := fiberErr.Message
			if fiberErr.Code == http.StatusRequestEntityTooLarge {
				code = "PAYLOAD_TOO_LARGE"
				msg = "ukuran payload melebihi batas maksimal 1 MB"
			}

			return c.Status(fiberErr.Code).JSON(ErrorResponse{
				Success:   false,
				Code:      code,
				Message:   msg,
				RequestID: reqID,
			})
		}

		logger.Error("unhandled_internal_error",
			slog.String("request_id", reqID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Any("error", err),
		)

		return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
			Success:   false,
			Code:      "INTERNAL_SERVER_ERROR",
			Message:   "terjadi kesalahan internal pada server",
			RequestID: reqID,
		})
	}
}

func httpStatusToCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	case http.StatusNotAcceptable:
		return "NOT_ACCEPTABLE"
	case http.StatusUnsupportedMediaType:
		return "UNSUPPORTED_MEDIA_TYPE"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusUnprocessableEntity:
		return "VALIDATION_ERROR"
	case http.StatusTooManyRequests:
		return "RATE_LIMIT_EXCEEDED"
	default:
		return "HTTP_ERROR"
	}
}

func NewAsynqRedisOpt(env *Env) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     env.RedisAddr(),
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	}
}

func NewAsynqClient(env *Env) *asynq.Client {
	return asynq.NewClient(NewAsynqRedisOpt(env))
}

func NewAsynqInspector(env *Env) *asynq.Inspector {
	return asynq.NewInspector(NewAsynqRedisOpt(env))
}
