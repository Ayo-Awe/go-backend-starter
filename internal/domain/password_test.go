package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	tests := []struct {
		password        string
		isValidPassword bool
	}{
		{"BackroomAuth2", true},
		{"authFlow 2", true},

		{"", false},
		{"short", false},
		{"nouppercase1", false},
		{"NOLOWERCASE2", false},
		{"no Digit Here", false},
	}

	for _, test := range tests {
		t.Run(test.password, func(t *testing.T) {
			pwd, err := NewPassword(test.password)

			if test.isValidPassword {
				assert.Equal(t, test.password, pwd.Value())
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
