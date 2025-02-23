package domain

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 64
)

var passwordValidationPatterns = map[*regexp.Regexp]error{
	regexp.MustCompile(`[A-Z]`): errors.New("password must contain an uppercase character"),
	regexp.MustCompile(`[a-z]`): errors.New("password must contain an lowercase character"),
	regexp.MustCompile(`\d`):    errors.New("password must contain a digit"),
}

// Password is a domain type encapsulating the validation rules for a password
// It can be used during signups, password resets and password changes but not login!
type Password struct {
	value string
}

func (p Password) Value() string { return p.value }
func (p Password) IsZero() bool  { return p.value == "" }

func NewPassword(password string) (*Password, error) {
	if len(password) < MinPasswordLength {
		return nil, fmt.Errorf("password must be at least %d characters", MinPasswordLength)
	}

	if len(password) > MaxPasswordLength {
		return nil, fmt.Errorf("password must be at most %d characters", MaxPasswordLength)
	}

	// validate password against regexp
	for pattern, err := range passwordValidationPatterns {
		if !pattern.Match([]byte(password)) {
			return nil, err
		}
	}

	return &Password{value: password}, nil
}
