package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TokenCreate(ttl time.Duration, email string) (string, error) {
	c := jwt.MapClaims{
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"sub": email,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	r, err := t.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return r, nil
}

func TokenValidate(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return nil, err
	}

	return claims, nil
}
