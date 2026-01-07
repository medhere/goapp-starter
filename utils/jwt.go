package utils

import (
	"errors"
	"goapp/config"
	"goapp/database"
	"goapp/model"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
)

var CONN = database.DB

var accessSecret = config.Env("JWT_SIGNING_KEY", "secret")
var refreshSecret = config.Env("REFRESH_TOKEN_KEY", "refresh-secret")
var accessTokenDuration = time.Hour * 12
var refreshTokenDuration = time.Hour * 24 * 7

func hashToken(t string) string {
	return "hash_" + t
}

// --- JWT Middleware Setup ---
func JWTAuthMiddlewareConfig() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Env("JWT_SIGNING_KEY", "secret"))},
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

// Retrieve in Fiber handler
// user := c.Locals("user").(*jwt.Token)
// claims := user.Claims.(jwt.MapClaims)
// userID := int(claims["uid"].(float64))

// CreateAccessAndRefreshToken generates a new access token and either
// inserts or updates a refresh token row for a given user/device.
// func CreateAccessAndRefreshToken(uid int, deviceID string) (accessToken string, refreshToken string, err error) {
// 	now := time.Now()
// 	accessExp := now.Add(accessTokenDuration)
// 	refreshExp := now.Add(refreshTokenDuration)

// 	// Create access token
// 	accessClaims := jwt.MapClaims{
// 		"uid": uid,
// 		"exp": accessExp.Unix(),
// 	}
// 	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessSecret)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	// Create refresh token
// 	refreshClaims := jwt.MapClaims{
// 		"uid": uid,
// 		"exp": refreshExp.Unix(),
// 	}
// 	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecret)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	// Check if a refresh token already exists for this user/device
// 	var authToken model.JWTAuthToken
// 	if err := CONN.Where("user_id = ? AND device_id = ?", uid, deviceID).First(&authToken).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			// Insert new row
// 			CONN.Create(&model.JWTAuthToken{
// 				UserID:    uid,
// 				DeviceID:  deviceID,
// 				TokenHash: hashToken(refreshToken),
// 				ExpiresAt: refreshExp,
// 			})
// 		} else {
// 			return "", "", err
// 		}
// 	} else {
// 		// Update existing row
// 		authToken.TokenHash = hashToken(refreshToken)
// 		authToken.ExpiresAt = refreshExp
// 		authToken.RevokedAt = nil
// 		CONN.Save(&authToken)
// 	}

// 	return accessToken, refreshToken, nil
// }

// CreateAccessAndRefreshToken encapsulates the JWT creation and storage logic
func CreateAccessAndRefreshToken(uid uint, deviceID string, deviceName string) (accessToken string, refreshToken string, err error) {
	now := time.Now()
	accessExp := now.Add(accessTokenDuration)
	refreshExp := now.Add(refreshTokenDuration)

	// Create access token
	accessClaims := jwt.MapClaims{
		"uid": uid,
		"exp": accessExp.Unix(),
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessSecret)
	if err != nil {
		return "", "", err
	}

	// Create refresh token
	refreshClaims := jwt.MapClaims{
		"uid": uid,
		"exp": refreshExp.Unix(),
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	// Hash the refresh token for secure storage
	refreshTokenHash, hashErr := HashString(refreshToken, 4)
	if hashErr != nil {
		return "", "", hashErr
	}

	// Check if a refresh token already exists for this user/device
	var authToken model.JWTAuthToken
	if err := CONN.Where("user_id = ? AND device_id = ?", uid, deviceID).First(&authToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Insert new row
			if err := CONN.Create(&model.JWTAuthToken{
				UserID:     int(uid),
				DeviceID:   deviceID,
				DeviceName: &deviceName,
				TokenHash:  refreshTokenHash,
				ExpiresAt:  refreshExp,
			}).Error; err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	} else {
		// if avaialble, Update existing row
		authToken.TokenHash = refreshTokenHash
		authToken.ExpiresAt = refreshExp
		authToken.RevokedAt = nil
		if err := CONN.Save(&authToken).Error; err != nil {
			return "", "", err
		}
	}

	return accessToken, refreshToken, nil
}

// VerifyAccessToken validates an access token string.
func VerifyAccessToken(bearer string) (*jwt.Token, error) {
	tokenStr := strings.TrimPrefix(bearer, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("TOKEN_EXPIRED")
		}
		return nil, errors.New("TOKEN_INVALID")
	}

	return token, nil
}

// VerifyRefreshToken validates a refresh token, updates its hash (rotation),
// and returns a new access token while keeping one row per device.
func VerifyRefreshToken(refreshToken string, deviceID string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return "", "", errors.New("INVALID_REFRESH")
	}

	claims := token.Claims.(jwt.MapClaims)
	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	userID := int(claims["uid"].(float64))

	// Find the token in CONN
	var authToken model.JWTAuthToken
	if err := CONN.Where("token_hash = ? AND device_id = ?", hashToken(refreshToken), deviceID).First(&authToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("REFRESH_NOT_FOUND")
		}
		return "", "", err
	}

	// Check if revoked or expired
	if authToken.RevokedAt != nil || time.Now().After(exp) {
		CONN.Delete(&authToken)
		return "", "", errors.New("REFRESH_EXPIRED")
	}

	// Create new access token
	accessClaims := jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(accessTokenDuration).Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessSecret)
	if err != nil {
		return "", "", err
	}

	// Rotate refresh token by updating the same row
	refreshClaims := jwt.MapClaims{
		"uid": userID,
		"exp": exp.Unix(), // preserve original expiration
	}
	newRefresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	authToken.TokenHash = hashToken(newRefresh)
	CONN.Save(&authToken)

	return accessToken, newRefresh, nil
}

// GetAccessID extracts the user ID from a verified access token.
func GetAccessID(token *jwt.Token) int {
	claims := token.Claims.(jwt.MapClaims)
	return int(claims["uid"].(float64))
}
