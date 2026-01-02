package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	switch Env("DB_CONNECTION") {
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			Env("DB_USERNAME"),
			Env("DB_PASSWORD"),
			Env("DB_HOST"),
			Env("DB_PORT"),
			Env("DB_DATABASE"),
		)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			Env("DB_HOST"),
			Env("DB_PORT"),
			Env("DB_USERNAME"),
			Env("DB_PASSWORD"),
			Env("DB_DATABASE"),
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
}
