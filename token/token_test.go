package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	jwt.StandardClaims
	UserID string `json:"userID"`
}

const signkey = "abc123"

func TestToken(t *testing.T) {
	now := time.Now()
	claim := claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(time.Duration(3) * time.Second).Unix(),
		},
		UserID: "asdf",
	}

	tk, err := GenerateToken(claim, signkey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tk)

	//time.Sleep(time.Second * 4)
	claim2 := new(claims)
	if err := ParseToken(tk, signkey, claim2); err != nil {
		t.Fatal(err)
	}
	t.Log(claim2)
}
