package rsa

import (
	"bytes"
	"testing"
)

const (
	bits = 2048
)

var data = []byte("123qwe你好杭州")

func TestXRsa(t *testing.T) {
	pubKey := bytes.NewBuffer(make([]byte, 0))
	priKey := bytes.NewBuffer(make([]byte, 0))

	if err := CreateKeys(pubKey, priKey, bits); err != nil {
		t.Fatal(err)
	}

	xrsa, err := NewXRsa(pubKey.Bytes(), priKey.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	encode, err := xrsa.EncryptToBase64(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("encode:", encode)

	decode, err := xrsa.DecryptFromBase64(encode)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("decode:", string(decode))

	signcode, err := xrsa.SignToBase64(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Sign:", signcode)

	if err := xrsa.VerifyFromBase64(data, signcode); err != nil {
		t.Fatal(err)
	}
	t.Log("Verify:", err)
}

func TestXRsaKeys(t *testing.T) {
	pub, pri, err := CreateKeysToBase64()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", pub)
	t.Log("\n", pri)
}
