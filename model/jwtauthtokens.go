package model

import (
	"time"

	"gorm.io/gorm"
)

type JWTAuthToken struct {
	gorm.Model
	UserID    int
	TokenHash string
	DeviceID  string
	ExpiresAt time.Time
	RevokedAt *time.Time
}
