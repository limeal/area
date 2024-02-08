package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// It takes a map of keys and values, and an expiration time in seconds, and returns a JWT token
func GenerateJWT(keys map[string]interface{}, exp int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	for key := range keys {
		claims[key] = keys[key]
	}

	if exp > 0 {
		claims["exp"] = time.Now().Add(time.Second * time.Duration(exp)).Unix()
	}

	secret, present := os.LookupEnv("JWT_SECRET")

	if !present {
		return "", errors.New("You must precise a secret key")
	}

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, nil
}

// It takes a token string, and returns a map of claims and an error
func ValidateJWT(tokenString string) (map[string]interface{}, error) {

	secret, present := os.LookupEnv("JWT_SECRET")

	if !present {
		return nil, errors.New("You must precise a secret key")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("Invalid token")
	}
}
