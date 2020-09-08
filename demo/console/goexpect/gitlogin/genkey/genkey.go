package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/LonelyPale/goutils/crypto"
	"github.com/LonelyPale/goutils/crypto/aes"
)

func main() {
	data := []byte("github2140264")

	key, err := crypto.GenerateSecretKey(crypto.DefaultSecretKeyLength) //256bits
	if err != nil {
		log.Fatal(err)
	}

	bs, err := aes.Encrypt(data, key)
	if err != nil {
		log.Fatal(err)
	}

	bs2, err := aes.Decrypt(bs, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("data:=", string(bs2))
	fmt.Println("key:=", format(key))
	fmt.Println("ciphertext:=", format(bs))
}

func format(bs []byte) string {
	s := ""
	for _, b := range bs {
		s += strconv.Itoa(int(b)) + ","
	}
	return "[]byte{" + s[:len(s)-1] + "}"
}
