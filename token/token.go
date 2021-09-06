package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/LonelyPale/goutils/database/redis"
	"github.com/LonelyPale/goutils/errors"
	"github.com/LonelyPale/goutils/types"
)

const (
	CachePrefix = "token:"
)

var (
	ErrToken        = errors.New("invalid token")
	ErrTokenTimeout = errors.New("token timeout")
)

type Token interface {
	jwt.Claims
	redis.CacheAble
	ID() string
}

type StandardToken struct {
	jwt.StandardClaims
}

func (t *StandardToken) ID() string {
	return t.Id
}

func (t *StandardToken) CacheKey() string {
	return CachePrefix + t.Id
}

type Options struct {
	SecretKey string
	Expire    int //Unit: Second
	Cache     redis.Cache
}

func New(opts ...*Options) StandardToken {
	now := time.Now()
	token := StandardToken{
		jwt.StandardClaims{
			Id:       types.NewObjectID().Hex(),
			IssuedAt: now.Unix(),
		},
	}

	if len(opts) > 0 && opts[0] != nil && opts[0].Expire > 0 {
		expire := opts[0].Expire
		token.ExpiresAt = now.Add(time.Second * time.Duration(expire)).Unix()
	}

	return token
}

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
