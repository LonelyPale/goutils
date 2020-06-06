package crypto

import (
	"encoding/base64"
	"encoding/hex"
)

type Hash []byte

// bss 可以是 nil
func NewHash(bss ...[]byte) Hash {
	var hash []byte
	for _, bs := range bss {
		hash = append(hash, bs...)
	}
	return hash
}

func (h Hash) Bytes() []byte {
	return h
}

func (h Hash) String() string {
	return h.Hex()
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h)
}

func (h Hash) Base64() string {
	return base64.RawURLEncoding.EncodeToString(h)
}

func (h *Hash) FromHex(s string) error {
	bs, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	*h = bs
	return nil
}

func (h *Hash) FromBase64(s string) error {
	bs, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	*h = bs
	return nil
}
