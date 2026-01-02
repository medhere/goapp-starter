package model

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model

	Name      string `gorm:"uniqueIndex;not null"`
	Hierarchy int    `gorm:"uniqueIndex;not null"` // 1 (Highest) -> N (Lowest)

}
