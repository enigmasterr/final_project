package jwt_tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const hmacSampleSecret = "super_secret_signature"

func GenerateSignedToken(userID string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": userID,
		"nbf":  now.Unix(),
		"exp":  now.Add(5 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})
	signedToken, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ParseToken(signedToken string) (*jwt.Token, error) {
	tokenFromString, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(hmacSampleSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return tokenFromString, nil
}
