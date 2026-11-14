package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New(fiber.Config{
		AppName: "api",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})
}
