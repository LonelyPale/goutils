package types

import (
	"encoding/base64"
	"encoding/hex"
)

// default base64 encoder
var Base64Encoding = base64.StdEncoding

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

func (b Bytes) Base64(encodings ...*base64.Encoding) string {
	var encoding *base64.Encoding
	if len(encodings) == 0 || encodings[0] == nil {
		encoding = Base64Encoding
	} else {
		encoding = encodings[0]
	}

	return encoding.EncodeToString(b)
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

func (b *Bytes) FromBase64(s string, encodings ...*base64.Encoding) error {
	var encoding *base64.Encoding
	if len(encodings) == 0 || encodings[0] == nil {
		encoding = Base64Encoding
	} else {
		encoding = encodings[0]
	}

	bs, err := encoding.DecodeString(s)
	if err != nil {
		return err
	}

	*b = bs
	return nil
}
