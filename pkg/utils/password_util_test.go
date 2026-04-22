package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptPassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid password", input: "my_secure_password", wantErr: false},
		{name: "empty password", input: "", wantErr: false}, // bcrypt allows empty password
		{name: "password too long", input: string(make([]byte, 100)), wantErr: true}, // bcrypt limits to 72 bytes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptPassword(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEqual(t, tt.input, got, "hashed password should not equal plaintext")
		})
	}
}

func TestComparePassword(t *testing.T) {
	validHash, err := EncryptPassword("password123")
	require.NoError(t, err)

	tests := []struct {
		name     string
		hash     string
		password string
		wantErr  bool
	}{
		{name: "correct password", hash: validHash, password: "password123", wantErr: false},
		{name: "incorrect password", hash: validHash, password: "wrong", wantErr: true},
		{name: "invalid hash format", hash: "invalid_hash_string", password: "password123", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ComparePassword(tt.hash, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
