package main

import (
	"goapp/config"
	"goapp/echorouter"
	"net/http"

	"github.com/labstack/echo/v4"
	// "golang.org/x/time/rate"
)

func main() {

	//Echo
	e := echo.New()
	// config.ConnectDB()
	// config.MigrateDB()

	// e.Use(middleware.Recover())
	// e.Use(middleware.BodyLimit("2M"))
	// e.Use(middleware.Decompress())
	// e.Use(middleware.Gzip())
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "method=${method}, uri=${uri}, status=${status}\n error=${error}\n query=${query}\n form=${form}\n",
	// }))
	// e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
	// 	ContentSecurityPolicy: "default-src 'self'",
	// }))
	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// 	AllowCredentials: true,
	// 	AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	// }))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	echorouter.EchoRoutes(e)

	// Start server
	// echorouter.ListRoutesOnConsole(e)
	e.Logger.Fatal(e.Start(config.Env("APP_HOST", "0.0.0.0") + ":" + config.Env("APP_ECHO_PORT", "8080")))

}
