package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/app/repository"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/app/service"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/config"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/database"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/helper"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/routes"
)

func main() {
	// LOAD ENV
	env, err := config.LoadEnv()
	if err != nil {
		fmt.Printf("Something Wrong When Loading Environment Variable: %s\n", err.Error())
		os.Exit(1)
	}

	// LOGGER CONFIG
	logger := config.NewLogger(env.AppEnv)

	// DB POOL
	pool, err := database.NewPostgresPool(env)
	if err != nil {
		fmt.Printf("CAN'T CREATE DB POOL: %v", err)
		os.Exit(1)
		return
	}
	defer pool.Close()

	mongoConn, err := database.NewMongoDatabase(env)
	if err != nil {
		logger.Error("gagal terhubung ke MongoDB", slog.String("[ERROR]", err.Error()))
		return
	}
	mongoDB := mongoConn.Database
	mongoClient := mongoConn.Client
	logger.Info("berhasil koneksi ke database MongoDB", slog.Any("MONGO_DB", mongoDB), slog.Any("MONGO_CLIENT", mongoClient))

	// JWT MANAGER
	jwtManager := helper.NewJWTManager(
		env.JWTSecret,
		env.JWTIssuer,
		env.JWTAccessTTL,
	)

	// COMMON PASSWORD
	passwordCommonPath, err := filepath.Abs("./files/common_password.txt")
	if err != nil {
		logger.Error("gagal memuat password umum", slog.String("error", err.Error()))
		os.Exit(1)
	}

	passwordCommonSet, err := helper.NewPasswordCommonSet(passwordCommonPath)
	if err != nil {
		logger.Error("gagal memuat password umum", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("Loaded Password Common Set", slog.Any("Loaded From", passwordCommonPath))

	// APP VALIDATOR
	appValidator := helper.NewValidator(passwordCommonSet)

	// REPO & SERVICES
	userRepo := repository.NewUserRepository(pool)

	authRepo := repository.NewAuthRepo(pool)
	authService := service.NewAuthService(
		userRepo,
		authRepo,
		jwtManager,
		env.JWTRefreshTTL,
		appValidator,
	)

	// APP
	deps := routes.Dependencies{
		Pool:        pool,
		JWT:         jwtManager,
		AuthService: authService,
	}

	app := config.NewFiberApp(env, logger, deps)
	port := env.AppPort

	// RUNNING THE APP
	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("SERVER STOPPED", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()
	logger.Info("SERVER RUNNING", slog.String("PORT", port))

	// GRACEFULL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("SHUTTING DOWN SERVER...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongoConn.Client.Disconnect(ctx)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("FAIL TO SHUTDOWN SERVER", slog.String("ERROR", err.Error()))
	}

	logger.Info("SERVER SHUTTED DOWN...")

}
