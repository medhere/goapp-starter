package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	switch Env("DB_CONNECTION", "sqlite") {
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			Env("DB_USERNAME", "root"),
			Env("DB_PASSWORD", ""),
			Env("DB_HOST", "localhost"),
			Env("DB_PORT", "3306"),
			Env("DB_DATABASE", "goapp"),
		)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	case "pgsql":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			Env("DB_HOST", "localhost"),
			Env("DB_PORT", "5432"),
			Env("DB_USERNAME", "postgres"),
			Env("DB_PASSWORD", ""),
			Env("DB_DATABASE", "goapp"),
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(Env("DB_DATABASE", "goapp")+".db"), &gorm.Config{})
	default:
		panic("Unsupported DB_CONNECTION type")
	}

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
}
