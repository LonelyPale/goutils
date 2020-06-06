package crypto

import "github.com/LonelyPale/goutils/crypto/sha512"

// 慢哈希加盐
func Bcrypt(data, salt []byte) ([]byte, error) {
	var hash []byte

	hash, err := sha512.Hash(data)
	if err != nil {
		return nil, err
	}

	hash, err = sha512.Hash(append(hash, salt...))
	if err != nil {
		return nil, err
	}

	for i := 0; i < 16; i++ {
		hash, err = sha512.Hash(hash)
		if err != nil {
			return nil, err
		}
	}

	return hash, nil
}
