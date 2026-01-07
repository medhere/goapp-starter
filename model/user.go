package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model // Includes ID, CreatedAt, and UpdatedAt (which map to TimestampsTz and ID in migration)

	// Personal Details
	Nickname  string `gorm:"not null"` // FullText index is applied in migration
	Firstname *string
	Lastname  *string
	Gender    *string
	Age       *int    // Use int or *int for nullable integer
	Avatar    *string // URL or path to avatar image
	Cover     *string // URL or path to cover image
	Bio       *string

	// Location (Using float64 for better geolocation data storage)
	Country *string
	State   *string
	City    *string
	Address *string

	// Contact Details
	Email string  `gorm:"uniqueIndex;not null"`
	Phone *string `gorm:"uniqueIndex"`

	// Authentication Credentials
	Roles    *datatypes.JSON
	Password *string
	Active   bool `gorm:"default:false"`

	// Signin Metadata
	SigninFlow        *datatypes.JSON //can select one or more signin flows (email/phone and (password, magic_link, otp, passkey, security_question))
	LastSigninAt      *time.Time
	SigninExpiresAt   *time.Time
	SigninCheckCount  *int `gorm:"default:0"`
	SigninFailedCount *int `gorm:"default:0"`
	SigninLockedUntil *time.Time

	// Failed Signin Metadata
	// FailedSigninCount *int `gorm:"default:0"` //TODO: fix duplication issue on codebase with SigninFailedCount
	// LockedUntil       *time.Time				//TODO: fix with SigninLockedUntil

	// Password Reset Fields
	PasswordResetTokenHash      *string
	PasswordResetExpiresAt      *time.Time
	PasswordResetFailedAttempts *int `gorm:"default:0"`
	PasswordResetLockedUntil    *time.Time

	// Verification Status (Email)
	EmailVerified             bool `gorm:"default:false"`
	EmailVerifyCodeHash       *string
	EmailVerifyCodeExpiresAt  *time.Time
	EmailVerifyFailedAttempts *int `gorm:"default:0"`
	EmailVerifyLockedUntil    *time.Time

	// Verification Status (Phone)
	PhoneVerified             *bool `gorm:"default:false"`
	PhoneVerifyCodeHash       *string
	PhoneVerifyCodeExpiresAt  *time.Time
	PhoneVerifyFailedAttempts *int `gorm:"default:0"`
	PhoneVerifyLockedUntil    *time.Time

	// OTP Fields
	OTPCodeHash *string
	OTPChannel  *string // (email | sms | whatsapp)

	// Magic Link Fields
	MagicLinkTokenHash *string

	// Passkeys/WebAuthn Credentials
	Passkeys *datatypes.JSON // Stores public_key and sign_count (JSON format)

	//Security Questions
	SecurityQuestion *string
	SecurityAnswer   *string

	// Soft Delete
	gorm.DeletedAt
}
