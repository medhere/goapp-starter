package utils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type AccessClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// --- JWT Middleware Setup ---
func EchoJWTAuthMiddlewareConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(AccessClaims)
		},
		SigningKey: []byte(JwTSigningKey),
		// TokenLookup: "header:Authorization:Bearer ",
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired access token")
		},
	}
}

// --- JWT Helper Functions ---
func EchoCreateAccessToken(userID string) (string, error) {
	claims := &AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JwTSigningKey))
}

// CreateRefreshToken generates a signed refresh JWT string.
func EchoCreateRefreshToken(userID string) (string, error) {
	claims := &RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(RefreshTokenKey))
}

// GetAccessClaims extracts and returns the AccessClaims from a validated JWT token within the Echo context.
func EchoGetAccessClaims(c echo.Context) *AccessClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*AccessClaims)
	return claims
}

// ValidateAndExtractRefreshClaims validates the refresh token and returns its claims.
func EchoValidateAndExtractRefreshClaims(tokenString string) (*RefreshClaims, error) {
	claims := new(RefreshClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(RefreshTokenKey), nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}

	return claims, nil
}

// EchoValidToken checks if the user ID in the token claims matches the requested ID.
// This helper is used inside protected handlers like UpdateUser/DeleteUser in an Echo application.
func EchoValidToken(c echo.Context, requestedID string) bool {
	token := c.Get("user").(*jwt.Token)
	claims, ok := token.Claims.(*AccessClaims)
	if !ok {
		return false
	}

	return claims.UserID == requestedID

	// Example usage in a protected handler:
	// 1. Authorization check using the Echo helper
	// id := c.Param("id")

	// if !utils.EchoValidToken(c, id) {
	//     return c.JSON(http.StatusUnauthorized, echo.Map{
	//         "message": "Unauthorized: Token ID does not match resource ID",
	//     })
	// }
}

// Login simulates a successful login and returns both tokens.
// func Login(c echo.Context) error {
// 	userID := "user456"

// 	accessToken, err := CreateAccessToken(userID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create access token"})
// 	}

// 	refreshToken, err := CreateRefreshToken(userID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create refresh token"})
// 	}

// 	return c.JSON(http.StatusOK, echo.Map{
// 		"access_token":  accessToken,
// 		"refresh_token": refreshToken,
// 		"token_type":    "Bearer",
// 		"expires_in":    AccessTokenDuration.Seconds(),
// 	})
// }

// RefreshToken endpoint uses a valid refresh token to issue a new access token.
// func RefreshToken(c echo.Context) error {
// 	type refreshRequest struct {
// 		RefreshToken string `json:"refresh_token"`
// 	}
// 	req := new(refreshRequest)
// 	if err := c.Bind(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
// 	}

// 	claims, err := ValidateAndExtractRefreshClaims(req.RefreshToken)
// 	if err != nil {
// 		return err
// 	}

// 	// Create a new access token
// 	newAccessToken, err := CreateAccessToken(claims.UserID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create new access token"})
// 	}

// 	// Optionally, you might also rotate the refresh token here (issue a new one and invalidate the old one)
// 	// For simplicity, we only issue a new access token.

// 	return c.JSON(http.StatusOK, echo.Map{
// 		"access_token": newAccessToken,
// 		"token_type":   "Bearer",
// 		"expires_in":   AccessTokenDuration.Seconds(),
// 	})
// }

// Restricted is a handler protected by the access token middleware.
// func Restricted(c echo.Context) error {
// 	claims := GetAccessClaims(c)

// 	return c.JSON(http.StatusOK, echo.Map{
// 		"message": "Access granted!",
// 		"user_id": claims.UserID,
// 	})
// }

// Public routes
// e.POST("/login", Login)
// e.POST("/refresh", RefreshToken)

// // Restricted group: Apply the access token middleware using e.Use syntax as requested
// r := e.Group("/restricted")
// r.Use(echojwt.WithConfig(JWTAuthMiddlewareConfig()))
// r.GET("", Restricted)
