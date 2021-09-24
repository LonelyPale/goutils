package random

import (
	"crypto/rand"
	"math"
	"math/big"

	"github.com/lonelypale/goutils/errors"
	"github.com/lonelypale/goutils/types"
)

// 生成区间[-m, n]的安全随机数, 真随机, 性能比较慢
func RangeRandom(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

func Random(length int) (types.Bytes, error) {
	if length <= 0 {
		return nil, errors.New("length must be greater than 0")
	}

	bs := make([]byte, length)
	if _, err := rand.Read(bs); err != nil {
		return nil, err
	}

	return bs, nil
}
