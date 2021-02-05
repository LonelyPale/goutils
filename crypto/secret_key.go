package crypto

import (
	"crypto/rand"
	"math/big"

	"github.com/LonelyPale/goutils/types"
)

const (
	DefaultSecretKeyLength = 32 // 32byte = 256bit
	DefaultSaltLength      = 32
)

// 生成指定字节长度的密钥
func GenerateSecretKey(lengths ...int) (types.Bytes, error) {
	var length int
	if len(lengths) > 0 && lengths[0] > 0 {
		//digits must be greater than 0
		length = lengths[0]
	} else {
		length = DefaultSecretKeyLength
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

func GenerateSalt(lengths ...int) (types.Bytes, error) {
	var length int
	if len(lengths) > 0 && lengths[0] > 0 {
		//digits must be greater than 0
		length = lengths[0]
	} else {
		length = DefaultSaltLength
	}
	return GenerateSecretKey(length)
}
