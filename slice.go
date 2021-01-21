package goutils

func MergeSliceByte(ss ...[]byte) []byte {
	switch len(ss) {
	case 0:
		return nil
	case 1:
		return ss[0]
	default:
		length := 0
		for _, s := range ss {
			length += len(s)
		}

		slice := make([]byte, length)
		index := 0
		for _, s := range ss {
			copy(slice[index:], s)
			index += len(s)
		}
		return slice
	}
}

func ReverseByte(arr []byte) []byte {
	length := len(arr)
	for i := 0; i < length/2; i++ {
		idx := length - 1 - i
		temp := arr[i]
		arr[i] = arr[idx]
		arr[idx] = temp
	}
	return arr
}
