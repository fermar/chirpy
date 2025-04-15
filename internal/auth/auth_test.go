package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Test: %v\nCheckPasswordHash() error = %v, wantErr %v",
					tt.name,
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestJWT(t *testing.T) {
	unuuid := uuid.New()
	tests := []struct {
		name       string
		iniuuid    uuid.UUID
		durstr     string
		dursleep   string
		claveMake  string
		claveParse string
		wantErr    bool
	}{
		{
			name:       "Correct",
			iniuuid:    unuuid,
			durstr:     "10m",
			dursleep:   "0s",
			claveMake:  "unaclaveparamake",
			claveParse: "unaclaveparamake",
			wantErr:    false,
		},
		{
			name:       "clave incorrecta",
			iniuuid:    uuid.Nil,
			durstr:     "10m",
			dursleep:   "0s",
			claveMake:  "unaclaveparamake",
			claveParse: "unaclaveparamake222",
			wantErr:    true,
		},
		{
			name:       "expirado",
			iniuuid:    uuid.Nil,
			durstr:     "2s",
			dursleep:   "3s",
			claveMake:  "unaclaveparamake",
			claveParse: "unaclaveparamake",
			wantErr:    true,
		},
	}
	// otrouuid := uuid.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration, _ := time.ParseDuration(tt.durstr)
			tokenString, _ := MakeJWT(tt.iniuuid, tt.claveMake, duration)
			dsleep, _ := time.ParseDuration(tt.dursleep)
			time.Sleep(dsleep)
			resuuid, err := ValidateJWT(tokenString, tt.claveParse)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Test: %v\nValidateJWT() error = %v, wantErr %v",
					tt.name,
					err,
					tt.wantErr,
				)
			}
			// if (err == nil) && (resuuid != tt.iniuuid) {
			if resuuid != tt.iniuuid {
				t.Errorf(
					"Test: %v\nValidateJWT() error = %v, wantErr %v",
					tt.name,
					"uuid missmatch",
					tt.wantErr,
				)
			}
		})
	}
}
