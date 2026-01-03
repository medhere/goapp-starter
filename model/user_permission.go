package model

import (
	"gorm.io/gorm"
)

type UserPermission struct {
	gorm.Model

	UserID       uint
	PermissionID uint

	// Granular Revocation Logic
	AssignerID uint // The ID of the user who assigned this permission
}
