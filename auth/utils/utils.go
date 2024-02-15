package utils

import (
	"auth/shared"
	"fmt"

	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT generates a JWT token with the provided claims and secret key
func GenerateJWT(claims jwt.Claims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateUserToken(username string, secretKey string) (shared.UserToken, error) {
	claims := jwt.MapClaims{
		"username": "user123",
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (24 hours from now)
	}

	// Generate JWT token
	token, err := GenerateJWT(claims, secretKey)

	userToken := shared.UserToken{
		Token:  token,
		Expire: claims["exp"].(string),
	}
	
	if err != nil {
		return userToken, err
	}

	return userToken, nil
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// In this example, we use a simple string as the key. In a real-world scenario, you should use a secure key management solution.
		return []byte("your_secret_key"), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
