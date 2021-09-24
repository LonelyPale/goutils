package sha256

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/lonelypale/goutils/types"
)

func HashReader(src io.Reader) (types.Bytes, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func Hash(data []byte) ([]byte, error) {
	hash := sha256.New()
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
