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
	Roles    *datatypes.JSON // Stores multiroles
	Password *string         // Stores the hashed password
	Active   bool            // Stores the active status

	// Verification Status (Email)
	EmailVerified            bool
	EmailVerifyCodeHash      *string
	EmailVerifyCodeExpiresAt *time.Time

	// Verification Status (Phone)
	PhoneVerified            *bool
	PhoneVerifyCodeHash      *string
	PhoneVerifyCodeExpiresAt *time.Time

	// OTP Fields
	OTPCodeHash  *string
	OTPChannel   *string // (email | sms | whatsapp)
	OTPExpiresAt *time.Time

	// Magic Link Fields
	MagicLinkTokenHash *string
	MagicLinkExpiresAt *time.Time

	// Passkeys/WebAuthn Credentials
	Passkeys *datatypes.JSON // Stores public_key and sign_count (JSON format)

	//Security Questions
	SecurityQuestion *string
	SecurityAnswer   *string

	// Login Metadata
	SigninFlow   *datatypes.JSON //can select one or more signin flows (email/phone and (password, magic_link, otp, passkey, security_question))
	LastLoginAt  *time.Time
	LoginCount   *int `gorm:"default:0"` // Use int for login_count
	FailedLogins *int `gorm:"default:0"`
	LockedUntil  *time.Time

	// Soft Delete
	gorm.DeletedAt
}
