package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type authClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(uid string, email string, role string, secretKey string) (string, error) {
	claims := &authClaims{
		email,
		role,
		jwt.RegisteredClaims{
			Subject:  uid,
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(encodedToken string, secretKey string) (bool, jwt.MapClaims, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token %v", token.Header["alg"])

		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return false, nil, err
	}

	if !token.Valid {
		return false, nil, nil
	}
	claims := token.Claims.(jwt.MapClaims)
	return true, claims, nil
}
