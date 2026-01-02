package echorouter

import (
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all API routes
func EchoRoutes(e *echo.Echo) {

	// Define your routes here
	AuthRoutes(e)
}

func ListRoutesOnConsole(e *echo.Echo) {
	routes := e.Routes()
	for _, route := range routes {
		println(route.Method, route.Path, route.Name)
	}

	// routesJSON, _ := json.MarshalIndent(routes, "", "  ")
	// println(string(routesJSON))

}
