package main

import (
	"goapp/config"
	"goapp/fiberrouter"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// config.ConnectDB()
	// config.MigrateDB()

	//Fiber
	c := fiber.New(fiber.Config{
		Prefork:           false,
		ServerHeader:      "goappServer",
		AppName:           "goapp Application",
		ReduceMemoryUsage: true,
		EnablePrintRoutes: true,
		DisableKeepalive:  true,
	})
	c.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Login endpoint")
	})

	fiberrouter.FiberRoutes(c)

	// Start server
	c.Listen(config.Env("APP_HOST") + ":" + config.Env("APP_FIBER_PORT"))

}
