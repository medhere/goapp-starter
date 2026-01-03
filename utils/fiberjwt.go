package utils

import (
	"net/http"
	"strconv"
	"time"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// --- JWT Middleware Setup ---
func FiberJWTAuthMiddlewareConfig() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(JwTSigningKey)},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
			}
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"status": "error", "message": "Invalid or expired access token", "data": nil})
		},
	})
}

// --- JWT Helper Functions ---
func FiberCreateAccessToken(userID string) (string, error) {
	// jwt.MapClaims is used instead of a custom struct
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(AccessTokenDuration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JwTSigningKey))
}

// CreateRefreshToken generates a signed refresh JWT string using jwt.MapClaims.
func FiberCreateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(RefreshTokenDuration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(RefreshTokenKey))
}

// GetAccessClaims extracts and returns the jwt.MapClaims from a validated JWT token within the Fiber context.
func FiberGetAccessClaims(c fiber.Ctx) jwt.MapClaims {
	// Fiber's JWT middleware stores the token in c.Locals("user")
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	return claims
}

// ValidateAndExtractRefreshClaims validates the refresh token and returns its jwt.MapClaims.
func FiberValidateAndExtractRefreshClaims(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(http.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(RefreshTokenKey), nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fiber.NewError(http.StatusUnauthorized, "Invalid refresh token")
	}

	return claims, nil
}

// ValidToken checks if the user ID in the token claims matches the requested ID.
// This helper is often used inside protected handlers like UpdateUser/DeleteUser.
func FiberValidToken(c fiber.Ctx, requestedID string) bool {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	tokenUserIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return false
	}
	tokenUserID := strconv.Itoa(int(tokenUserIDFloat))

	return tokenUserID == requestedID

	// Example usage in a protected handler:
	// 1. Authorization check using the Fiber helper
	// id := c.Params("id")
	// if !utils.FiberValidToken(c, id) {
	//     return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	//         "message": "Unauthorized: Token ID does not match resource ID",
	//     })
	// }
}

// // Login (Simplified stub)
// func Login(c fiber.Ctx) error {
// 	// ... authentication logic ...
// 	userID := "user456"

// 	accessToken, _ := utils.CreateAccessToken(userID)
// 	refreshToken, _ := utils.CreateRefreshToken(userID)

// 	return c.JSON(fiber.Map{
// 		"access_token":  accessToken,
// 		"refresh_token": refreshToken,
// 	})
// }

// // RefreshToken (Simplified stub)
// func RefreshToken(c fiber.Ctx) error {
// 	// ... token validation logic using utils.ValidateAndExtractRefreshClaims ...
// 	return c.SendString("New access token issued")
// }

// // Restricted (Protected handler stub)
// func Restricted(c fiber.Ctx) error {
// 	claims := utils.GetAccessClaims(c)

// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"message": "Access granted!",
// 		"user_id": claims["user_id"],
// 	})
// }
// app.Post("/login", handler.Login)
// app.Post("/refresh", handler.RefreshToken)
// r := app.Group("/restricted")
// r.Use(utils.Protected())
// r.Get("", handler.Restricted)
