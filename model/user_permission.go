package model

import (
	"gorm.io/gorm"
)

type UserPermission struct {
	gorm.Model

	UserID       uint `gorm:"not null;index:idx_user_perm"`
	PermissionID uint `gorm:"not null;index:idx_user_perm"`

	// Granular Revocation Logic
	AssignerID uint `gorm:"not null"` // The ID of the user who assigned this permission
}
