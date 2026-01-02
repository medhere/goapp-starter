package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Scope     string `gorm:"type:enum('backend','frontend');not null"`
	TargetApp string `gorm:"type:enum('dashboard','client');not null"`

	Category    string `gorm:"size:100;not null"`
	SubCategory string `gorm:"size:100"`

	PermissionName string `gorm:"size:100;uniqueIndex;not null"`

	AssignableRoles datatypes.JSON `gorm:"type:json;not null"` // JSON array of role names
}
