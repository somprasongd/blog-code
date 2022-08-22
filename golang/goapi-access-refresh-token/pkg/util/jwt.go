package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(uid string, payload map[string]any, secretKey string, expires time.Duration) (string, time.Time, error) {
	claims := jwt.MapClaims{}
	claims["sub"] = uid
	claims["iat"] = jwt.NewNumericDate(time.Now())

	expiresAt := time.Now().Add(expires)
	claims["exp"] = jwt.NewNumericDate(expiresAt)

	for k, v := range payload {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	encodedToken, err := token.SignedString([]byte(secretKey))
	return encodedToken, expiresAt, err
}

func ValidateToken(encodedToken string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token %v", token.Header["alg"])

		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}
