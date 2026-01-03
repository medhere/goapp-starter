package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Scope     string
	TargetApp string

	Category    string
	SubCategory string

	PermissionName string

	AssignableRoles datatypes.JSON // JSON array of role names
}
