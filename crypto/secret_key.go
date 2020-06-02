package crypto

import (
	"crypto/rand"
	"math/big"

	"github.com/LonelyPale/goutils/errors"
)

// 生成指定长度的密钥
func GenerateSecretKey(length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("digits must be greater than 0")
	}

	bs := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(256)) // 返回一个在[0, max)区间服从均匀分布的随机值，如果max<=0则会panic
		if err != nil {
			return nil, err
		}

		bs[i] = uint8(num.Uint64())
	}

	return bs, nil
}
