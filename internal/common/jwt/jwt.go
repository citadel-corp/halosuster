package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	key = []byte(os.Getenv("JWT_SECRET"))

	ErrUnknownClaims = errors.New("unknown claims type")
	ErrTokenInvalid  = errors.New("invalid token")
)

type CustomClaims struct {
	UserType string `json:"userType"`
	jwt.RegisteredClaims
}

func Sign(ttl time.Duration, subject string, userType string) (string, error) {
	now := time.Now()
	expiry := now.Add(ttl)
	claims := CustomClaims{
		userType,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   subject,
		},
	}
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	return t.SignedString(key)
}

func VerifyAndGetSubject(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return "", "", err
	}

	// Checking token validity
	if !token.Valid {
		return "", "", ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims.Subject, claims.UserType, nil
	} else {
		return "", "", ErrUnknownClaims
	}
}
