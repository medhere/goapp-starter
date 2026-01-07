package utils

import "github.com/gofiber/fiber/v3"

func GetBody(c fiber.Ctx, req interface{}) error {
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body format or type mismatch.",
			"error":   err.Error(),
		})
	}
	return nil
}
