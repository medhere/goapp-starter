package fiberrouter

import "github.com/gofiber/fiber/v3"

func AuthRoutes(app *fiber.App) {
	authGroup := app.Group("/auth")

	authGroup.Post("/login", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Login endpoint")
	})

	authGroup.Post("/register", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Register endpoint")
	})

	authGroup.Post("/logout", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Logout endpoint")
	})
}
