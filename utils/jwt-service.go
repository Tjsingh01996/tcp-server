package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key")

// GenerateJWT creates a JWT token with a given payload and expiration time
func GenerateJWT(secretKey []byte, payload map[string]any, expiry time.Duration) (string, error) {
	// Create token claims
	claims := jwt.MapClaims{}

	// Add user-defined payload to claims
	for key, value := range payload {
		claims[key] = value
	}
	if expiry != 0 {
		// Set expiration time
		claims["exp"] = time.Now().Add(expiry).Unix()
	}

	// Create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(secretKey []byte, tokenString string) (map[string]interface{}, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	// Handle parsing errors
	if err != nil {
		return nil, err
	}

	// Extract claims if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return nil, fmt.Errorf("token has expired")
			}
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
