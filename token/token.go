package token

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"github.com/LonelyPale/goutils/errors"
)

var (
	ErrToken        = errors.New("invalid token")
	ErrTokenTimeout = errors.New("token timeout")
)

func GenerateToken(claims jwt.Claims, signingKey string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(signingKey))
}

func ParseToken(tokenStr string, signingKey string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return ErrTokenTimeout
			}
		}

		return err
	}

	if token != nil {
		if _, ok := token.Claims.(jwt.Claims); ok && token.Valid {
			return nil
		}
	}

	return ErrToken
}
