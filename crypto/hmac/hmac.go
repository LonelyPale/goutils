package hmac

import (
	"crypto/hmac"
	"hash"
)

func Hash(h func() hash.Hash, data []byte, key []byte) ([]byte, error) {
	hasher := hmac.New(h, key)
	if _, err := hasher.Write(data); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}
