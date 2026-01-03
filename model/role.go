package model

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model

	Name      string
	Hierarchy int // 1 (Highest) -> N (Lowest)

}
