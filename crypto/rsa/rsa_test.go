package rsa

import (
	"bytes"
	"testing"
)

const (
	bits = 2048
	data = "123qwe你好杭州"
)

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

	encode, err := xrsa.PublicEncrypt(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("encode:", encode)

	decode, err := xrsa.PrivateDecrypt(encode)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("decode:", decode)

	signcode, err := xrsa.Sign(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Sign:", signcode)

	if err := xrsa.Verify(data, signcode); err != nil {
		t.Fatal(err)
	}
	t.Log("Verify:", err)
}
