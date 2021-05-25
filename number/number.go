package number

import "math"

// 将 float64 转成精确 retain 位的 int64
func Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

// 将 int64 恢复成精确 retain 位的正常 float64
func Unwrap(num int64, retain int) float64 {
	return float64(num) / math.Pow10(retain)
}

// 精准 retain 位的 float64
func WrapToFloat64(num float64, retain int) float64 {
	return num * math.Pow10(retain)
}

// 精准 retain 位的 int64
func UnwrapToInt64(num int64, retain int) int64 {
	return int64(Unwrap(num, retain))
}
