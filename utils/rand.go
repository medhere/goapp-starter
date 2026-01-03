package utils

import (
	"crypto/rand"
	"fmt"
	"io"
)

func generateRandomCode(length int, charset string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than zero")
	}

	// Convert the string charset to a slice of bytes
	charsetBytes := []byte(charset)

	// Get the length of the character set
	charsetLength := len(charsetBytes)

	// Create a byte slice to hold the random index values
	randomBytes := make([]byte, length)

	// Read cryptographically secure random bytes
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	// Create the final code string
	code := make([]byte, length)

	// Map the random bytes to indices within the charset range
	for i := 0; i < length; i++ {
		// Use the modulo operator to ensure the index is within the charset boundaries.
		// This distributes the random bytes evenly across the charset.
		code[i] = charsetBytes[int(randomBytes[i])%charsetLength]
	}

	return string(code), nil
}

func GenerateAlphaCode(length int) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return generateRandomCode(length, chars)
}

func GenerateNumberCode(length int) (string, error) {
	const chars = "0123456789"
	return generateRandomCode(length, chars)
}

func GenerateAlphaNumCode(length int) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return generateRandomCode(length, chars)
}
