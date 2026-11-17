package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/config"
	"github.com/Alan-00280/EventSeat-PBL-UAS.git/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// LOAD ENV
	env, err := config.LoadEnv()
	if err != nil {
		fmt.Printf("Something Wrong When Loading Environment Variable: %w", err)
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
	app := config.NewFiberApp(env, logger)

	// RUNNING THE APP

	// GRACEFULL SHUTDOWN

}
