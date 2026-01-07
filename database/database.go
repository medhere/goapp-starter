package database

import (
	"fmt"
	"goapp/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	switch config.Env("DB_CONNECTION", "sqlite") {
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Env("DB_USERNAME", "root"),
			config.Env("DB_PASSWORD", ""),
			config.Env("DB_HOST", "localhost"),
			config.Env("DB_PORT", "3306"),
			config.Env("DB_DATABASE", "aroundme"),
		)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	case "pgsql":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Env("DB_HOST", "localhost"),
			config.Env("DB_PORT", "5432"),
			config.Env("DB_USERNAME", "postgres"),
			config.Env("DB_PASSWORD", ""),
			config.Env("DB_DATABASE", "aroundme"),
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(config.Env("DB_DATABASE", "aroundme")+".db"), &gorm.Config{})
	default:
		panic("Unsupported DB_CONNECTION type")
	}

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
}
