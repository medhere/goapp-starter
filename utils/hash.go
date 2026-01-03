package utils

import "golang.org/x/crypto/bcrypt"

// HashString hashes a plain text string using bcrypt with a specified cost.
func HashString(toHash string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		// Use default cost if the specified cost is out of the valid range [4, 31]
		cost = bcrypt.DefaultCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(toHash), cost)
	return string(bytes), err
}

// CheckHash compares a plain text string with a hashed string.
// Returns nil on success, or an error on failure (e.g., mismatch or invalid hash).
func CheckHash(toCheck, hashString string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashString), []byte(toCheck))
}
