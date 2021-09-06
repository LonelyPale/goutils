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

func TestClaims(t *testing.T) {
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

type token struct {
	StandardToken
	UserID string `json:"userID"`
}

func TestToken(t *testing.T) {
	opt := &Options{
		SecretKey: signkey,
		Expire:    3,
		Cache:     nil,
	}
	tkn := &token{
		StandardToken: New(opt),
		UserID:        "zxcv",
	}

	tk, err := GenerateToken(tkn, opt.SecretKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tk)

	//time.Sleep(time.Second * 4)
	tkn2 := new(token)
	if err := ParseToken(tk, opt.SecretKey, tkn2); err != nil {
		t.Fatal(err)
	}
	t.Log(tkn2)
}
