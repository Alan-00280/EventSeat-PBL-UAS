package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/config"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/database"
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

	// APP
	deps := routes.Dependencies{
		Pool: pool,
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
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("FAIL TO SHUTDOWN SERVER", slog.String("ERROR", err.Error()))
	}

	logger.Info("SERVER SHUTTED DOWN...")

}
