package crypto

import "github.com/LonelyPale/goutils/crypto/sha512"

const (
	DefaultLoopNumber = 2048
)

// 慢哈希加盐
func Bcrypt(data, salt []byte, loops ...int) ([]byte, error) {
	var loop int
	if len(loops) > 0 && loops[0] > 0 {
		loop = loops[0]
	} else {
		loop = DefaultLoopNumber
	}

	hash, err := sha512.Hash(data)
	if err != nil {
		return nil, err
	}

	hash, err = sha512.Hash(append(hash, salt...))
	if err != nil {
		return nil, err
	}

	for i := 0; i < loop; i++ {
		hash, err = sha512.Hash(hash)
		if err != nil {
			return nil, err
		}
	}

	return hash, nil
}
