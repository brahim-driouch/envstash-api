// internal/utils/email.go
package utils

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"

	"github.com/resend/resend-go/v3"
)

// EmailData holds the data to inject into email templates
type EmailData struct {
	Fullname         string
	Email            string
	VerificationLink string
	ResetLink        string
}

// SendEmail sends an email using Resend
func SendEmail(ctx context.Context, emailParams *resend.SendEmailRequest) error {
	resendAPIKey := os.Getenv("RESEND_API_KEY")
	if resendAPIKey == "" {
		return fmt.Errorf("RESEND_API_KEY not set")
	}

	client := resend.NewClient(resendAPIKey)

	sent, err := client.Emails.SendWithContext(ctx, emailParams)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if sent != nil && sent.Id != "" {
		fmt.Printf("✓ Email sent successfully (ID: %s)\n", sent.Id)
	}

	return nil
}

// RenderEmailTemplate renders an HTML template with data
func RenderEmailTemplate(templatePath string, data EmailData) (string, error) {
	// Parse the template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// BuildVerificationEmail creates a verification email with rendered template
func BuildVerificationEmail(toEmail, fullname, verificationToken string) (*resend.SendEmailRequest, error) {
	// Build verification link
	frontendURL := os.Getenv("FRONTEND_URL")
	verificationLink := fmt.Sprintf("%s/auth/verify?token=%s", frontendURL, verificationToken)

	// Prepare template data
	data := EmailData{
		Fullname:         fullname,
		Email:            toEmail,
		VerificationLink: verificationLink,
	}

	// Render HTML template
	htmlBody, err := RenderEmailTemplate("templates/verification.html", data)
	if err != nil {
		return nil, err
	}

	// Build email request
	emailParams := &resend.SendEmailRequest{
		From:    "Envify <noreply@dataentryjobs.io>",
		To:      []string{toEmail},
		Subject: "Verify Your Email - Envify",
		Html:    htmlBody,
		// Text:    fmt.Sprintf("Hey %s! Click here to verify your email: %s", fullname, verificationLink),
	}

	return emailParams, nil
}

// BuildPasswordResetEmail creates a password reset email with rendered template
func BuildPasswordResetEmail(toEmail, fullname, resetToken string) (*resend.SendEmailRequest, error) {
	// Build reset link
	frontendURL := os.Getenv("FRONTEND_URL")
	resetLink := fmt.Sprintf("%s/auth/reset-password?token=%s", frontendURL, resetToken)

	// Prepare template data
	data := EmailData{
		Fullname:  fullname,
		Email:     toEmail,
		ResetLink: resetLink,
	}

	// Render HTML template
	htmlBody, err := RenderEmailTemplate("templates/emails/password-reset.html", data)
	if err != nil {
		return nil, err
	}

	// Build email request
	emailParams := &resend.SendEmailRequest{
		From:    "Envify <noreply@envify.dev>",
		To:      []string{toEmail},
		Subject: "Reset Your Password - Envify",
		Html:    htmlBody,
		Text:    fmt.Sprintf("Hey %s! Click here to reset your password: %s", fullname, resetLink),
	}

	return emailParams, nil
}
