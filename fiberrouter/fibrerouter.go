package fiberrouter

import "github.com/gofiber/fiber/v2"

func FiberRoutes(app *fiber.App) {
	AuthRoutes(app)
}
