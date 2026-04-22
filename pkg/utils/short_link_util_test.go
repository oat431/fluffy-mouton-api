package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateName(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"generates random name 1"},
		{"generates random name 2"},
		{"generates random name 3"},
	}

	stringPool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateName()
			assert.Len(t, got, 5, "GenerateName() should return a 5 character string")
			
			for _, char := range got {
				assert.True(t, strings.ContainsRune(stringPool, char), "GenerateName() contains invalid character: %c", char)
			}
		})
	}
}
