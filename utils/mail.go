package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/mail.v2"
)

// SendEmail sends a comprehensive email with optional attachments.
// It returns true on success, false on failure, with errors logged internally.
func SendEmail(
	fromEmail string,
	toEmail string,
	subject string,
	htmlBody string,
	attachmentPath string, // Optional path to a file
) bool {
	// --- 1. Load Configuration ---
	// In a production app, load this once at startup, not on every call.
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpPortStr == "" || smtpUser == "" || smtpPass == "" {
		log.Println("ERROR: SMTP environment variables are not fully set. Skipping email send.")
		return false
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Printf("ERROR: Invalid SMTP_PORT value: %s", smtpPortStr)
		return false
	}

	// --- 2. Compose the Message ---
	m := mail.NewMessage()

	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)

	// Set the HTML body and a basic plain-text fallback
	m.SetBody("text/html", htmlBody)
	// Attempt to strip HTML tags for plain-text fallback (simple version)
	m.AddAlternative("text/plain", fmt.Sprintf("Email Content:\n%s", stripHtmlTags(htmlBody)))

	// --- 3. Handle Optional Attachment ---
	if attachmentPath != "" {
		// Check if the file exists before adding
		if _, err := os.Stat(attachmentPath); os.IsNotExist(err) {
			log.Printf("WARNING: Attachment file not found at path: %s", attachmentPath)
			// Continue without attachment, but log the warning
		} else {
			// Include attachment with the original filename
			m.Attach(attachmentPath, mail.Rename(filepath.Base(attachmentPath)))
			log.Printf("Attaching file: %s", filepath.Base(attachmentPath))
		}
	}

	// --- 4. Create a Dialer and Send ---
	d := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("ERROR: Failed to send email to %s (Subject: %s): %v", toEmail, subject, err)
		return false
	}

	log.Printf("Successfully sent email to %s (Subject: %s)", toEmail, subject)
	return true
}

// stripHtmlTags is a basic utility to remove HTML tags for plain text fallback.
// Note: A robust implementation would use a proper HTML parsing library.
func stripHtmlTags(s string) string {
	// Simple replacement, replacing <...> with an empty string
	// This is not foolproof but serves for a basic text fallback
	inTag := false
	var result []rune
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result = append(result, r)
		}
	}
	return string(result)
}
