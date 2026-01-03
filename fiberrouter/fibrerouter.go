package fiberrouter

import "github.com/gofiber/fiber/v3"

func FiberRoutes(app *fiber.App) {
	AuthRoutes(app)
}
