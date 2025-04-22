package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(hash, password string) error {
	// slog.Info("check", "hash", hash, "passwd", password)
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
			Subject:   userID.String(),
		},
	)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}
	if !parsedToken.Valid {
		return uuid.Nil, fmt.Errorf("error parsing token")
	}
	subject, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	retuid, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, err
	}
	return retuid, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	key := strings.TrimPrefix(headers.Get("Authorization"), "Bearer ")

	if key == "" {
		return "", errors.New("missing authorization header")
	}
	return key, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	key := strings.TrimPrefix(headers.Get("Authorization"), "ApiKey ")

	if key == "" {
		return "", errors.New("missing authorization header")
	}
	return key, nil
}

func MakeRefreshToken() (string, error) {
	tokenkey := make([]byte, 32)
	rand.Read(tokenkey)
	refreshToken := hex.EncodeToString(tokenkey)

	return refreshToken, nil
}
