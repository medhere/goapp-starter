package utils

import (
	"github.com/go-playground/validator/v10" // Assuming this library
	"github.com/gofiber/fiber/v3"
)

var VALID = validator.New()

func ValidateThis(c fiber.Ctx, req interface{}, msg string) error {

	if err := VALID.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": msg,
			"errors":  errors,
		})
	}

	return nil
}
