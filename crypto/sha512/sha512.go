package sha512

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

func Hash(data []byte) ([]byte, error) {
	hash := sha512.New()
	if _, err := hash.Write(data); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func Hex(data []byte) (string, error) {
	bs, err := Hash(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bs), nil
}

func Base64(data []byte) (string, error) {
	bs, err := Hash(data)
	if err != nil {
		return "", nil
	}
	return base64.RawURLEncoding.EncodeToString(bs), nil
}
