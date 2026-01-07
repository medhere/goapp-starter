package main

import (
	"goapp/config"
	"goapp/database"
	"goapp/fiberrouter"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	// "github.com/gofiber/fiber/v3/middleware/compress"
	// "github.com/gofiber/fiber/v3/middleware/cors"
	// "github.com/gofiber/fiber/v3/middleware/etag"
	// "github.com/gofiber/fiber/v3/middleware/limiter"
)

func main() {
	database.ConnectDB()
	// database.MigrateDB()

	//Fiber
	isDebugEnabled := strings.ToLower(config.Env("APP_DEBUG", "false")) == "true"
	app := fiber.New(
		fiber.Config{
			ReduceMemoryUsage: true,
			BodyLimit:         10 * 1024 * 1024, // 10 MB
			DisableKeepalive:  false,
			// CaseSensitive:     true,
			ServerHeader: config.Env("APP_NAME", "GoFiber") + "Server",
			AppName:      config.Env("APP_NAME", "GoFiber") + " Application",
		})

	if isDebugEnabled {
		app.Use(logger.New())
		log.Println("✅ Debug Mode: Fiber Logger Middleware ENABLED.")
	}
	app.Use(recover.New())
	// app.Use(compress.New())
	// app.Use(etag.New())
	// app.Use(limiter.New())
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"*"},
	// 	AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	// }))

	fiberrouter.FiberRoutes(app)

	listenAddr := config.Env("APP_HOST", "0.0.0.0") + ":" + config.Env("APP_FIBER_PORT", "3000")
	listenConfig := fiber.ListenConfig{}

	if isDebugEnabled {
		listenConfig.EnablePrintRoutes = true
		log.Println("✅ Debug Mode: Route Printing ENABLED.")
	}
	//Set Prefork
	listenConfig.EnablePrefork = false

	// Start server
	err := app.Listen(listenAddr, listenConfig)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
