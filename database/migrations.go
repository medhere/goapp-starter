package database

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
		&model.JWTAuthToken{},
	)
	fmt.Println("Database Migrated")

}
