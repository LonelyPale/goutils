package hash

import (
	"encoding/base64"
	"encoding/hex"
	"hash"
)

func Hash(h func() hash.Hash, data []byte) ([]byte, error) {
	hasher := h()
	if _, err := hasher.Write(data); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func Hex(h func() hash.Hash, data []byte) (string, error) {
	bs, err := Hash(h, data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bs), nil
}

func Base64(h func() hash.Hash, data []byte) (string, error) {
	bs, err := Hash(h, data)
	if err != nil {
		return "", nil
	}
	return base64.RawURLEncoding.EncodeToString(bs), nil
}
