package mail

import (
	"embed"
)

//go:embed "templates"
var templateFS embed.FS

const (
	EmailVerificationTemplate Template = "email_verification.tmpl"
)

type EmailVerificationEmail struct {
	OTP              string
	OTPExpiryMinutes int
}

func (l EmailVerificationEmail) Template() Template { return EmailVerificationTemplate }
