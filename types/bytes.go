package types

import (
	"encoding/base64"
	"encoding/hex"
)

type Bytes []byte

// bss 可以是 nil
// 等同于 Bytes(bs)
func NewBytes(bss ...[]byte) Bytes {
	var bytes []byte
	for _, bs := range bss {
		bytes = append(bytes, bs...)
	}
	return bytes
}

// 等同于 []byte(h)
func (b Bytes) Bytes() []byte {
	return b
}

func (b Bytes) String() string {
	return string(b)
}

func (b Bytes) Hex() string {
	return hex.EncodeToString(b)
}

func (b Bytes) Base64() string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func (b *Bytes) FromString(s string) {
	*b = []byte(s)
}

func (b *Bytes) FromHex(s string) error {
	bs, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	*b = bs
	return nil
}

func (b *Bytes) FromBase64(s string) error {
	bs, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	*b = bs
	return nil
}
