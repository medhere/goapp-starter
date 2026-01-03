package main

import (
	"goapp/config"
	"goapp/fiberrouter"

	"github.com/gofiber/fiber/v3"
)

func main() {

	//Fiber
	app := fiber.New(
		fiber.Config{
			ReduceMemoryUsage: true,
			BodyLimit:         10 * 1024 * 1024, // 10 MB
			DisableKeepalive:  false,
			CaseSensitive:     true,
			ServerHeader:      config.Env("APP_NAME", "Fiber") + "Server",
			AppName:           config.Env("APP_NAME", "Fiber") + " Application",
		})
	// config.ConnectDB()
	// config.MigrateDB()

	// app.Use(recoverer.New())
	// app.Use(compress.New())
	// app.Use(logger.New())
	// app.Use(etag.New())
	// app.Use(limiter.New())
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     "*",
	// 	AllowHeaders:     "Origin, Content-Type, Accept",
	// 	AllowCredentials: true,
	// 	AllowMethods:     "GET, POST, PUT, DELETE",
	// }))

	fiberrouter.FiberRoutes(app)

	// Start server
	// ReduceMemoryUsage: true,
	// EnablePrintRoutes: true,

	app.Listen(config.Env("APP_HOST", "0.0.0.0")+":"+config.Env("APP_FIBER_PORT", "3000"),
		fiber.ListenConfig{
			EnablePrefork: true,
		})
}
