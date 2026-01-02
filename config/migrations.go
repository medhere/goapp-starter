package config

import (
	"fmt"

	"goapp/model"
)

// MigrateDB migrate db
func MigrateDB() {
	DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserPermission{},
	)
	fmt.Println("Database Migrated")

}
