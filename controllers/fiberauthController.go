package controllers

import (
	"fmt"
	"goapp/database"
	"goapp/model"
	"goapp/utils"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

var fibrevalidate = validator.New()

func FiberSignup(c fiber.Ctx) error {
	var req struct {
		Nickname        string `json:"nickname" validate:"required,min=3,max=50"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required,min=8"`
		PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
	}

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if err := fibrevalidate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	var existingUser model.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already registered",
		})
	} else if err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error during check",
		})
	}

	hashedPassword, err := utils.HashString(req.Password, 4)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	verificationCode, err := utils.GenerateNumberCode(6)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate verification code",
		})
	}

	codeHash, _ := utils.HashString(verificationCode, 4)
	expiresAt := time.Now().Add(24 * time.Hour)

	newUser := model.User{
		Nickname:                 req.Nickname,
		Email:                    req.Email,
		Password:                 &hashedPassword,
		Active:                   false,
		EmailVerified:            false,
		EmailVerifyCodeHash:      &codeHash,
		EmailVerifyCodeExpiresAt: &expiresAt,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	go func() {
		htmlContent := fmt.Sprintf(`
			<p>Your verification code is: <strong>%s</strong></p>
			<p>This code is valid for 24 hours.</p>
		`, verificationCode)

		if !utils.SendEmail(
			"Your App <no-reply@yourapp.com>",
			newUser.Email,
			"Verify Your Email Address",
			htmlContent,
			"",
		) {
			log.Printf("CRITICAL: Mail sending failed for user %d (%s).", newUser.ID, newUser.Email)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully. Check your email for verification.",
		"user_id": newUser.ID,
	})
}
