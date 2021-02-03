package crypto

import (
	"github.com/LonelyPale/goutils"
	"github.com/LonelyPale/goutils/crypto/sha512"
)

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

	//1 正向哈希
	hash, err := sha512.Hash(data)
	if err != nil {
		return nil, err
	}

	//逆序盐
	saltReverse := goutils.ReverseByte(append([]byte{}, salt...))

	//2 前置盐
	salt, err = sha512.Hash(goutils.MergeSliceByte(salt, hash))
	if err != nil {
		return nil, err
	}

	//3 逆向哈希
	hashReverse, err := sha512.Hash(goutils.ReverseByte(data))
	if err != nil {
		return nil, err
	}
	goutils.ReverseByte(hashReverse)

	//4 后置盐
	saltReverse, err = sha512.Hash(goutils.MergeSliceByte(hashReverse, saltReverse))
	if err != nil {
		return nil, err
	}
	goutils.ReverseByte(saltReverse)

	//5 合并哈希
	newhash := make([]byte, 0)
	newhash = append(newhash, salt...)
	for i := range hash {
		newhash = append(newhash, hash[i])
		newhash = append(newhash, hashReverse[i])
	}
	newhash = append(newhash, saltReverse...)
	newhash, err = sha512.Hash(newhash)
	if err != nil {
		return nil, err
	}

	//6 重复哈希
	for i := 0; i < loop; i++ {
		if i%2 != 0 {
			goutils.ReverseByte(newhash)
		}

		newhash, err = sha512.Hash(newhash)
		if err != nil {
			return nil, err
		}
	}

	//7 重组哈希 4231是1、2、3、4能组成的最大素数
	b1 := make([]byte, 16)
	b2 := make([]byte, 16)
	b3 := make([]byte, 16)
	b4 := make([]byte, 16)
	for i := 0; i < 16; i++ {
		b1[i] = newhash[i*4]
		b2[i] = newhash[i*4+1]
		b3[i] = newhash[i*4+2]
		b4[i] = newhash[i*4+3]
	}
	a1 := b4
	a2 := goutils.ReverseByte(b2)
	a3 := goutils.ReverseByte(b3)
	a4 := b1
	newhash = goutils.MergeSliceByte(a1, a2, a3, a4)

	newhash, err = sha512.Hash(append(newhash, []byte("Bcrypt")...))
	if err != nil {
		return nil, err
	}

	return newhash, nil
}

func BcryptSimple(data, salt []byte, loops ...int) ([]byte, error) {
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

	salt, err = sha512.Hash(salt)
	if err != nil {
		return nil, err
	}

	hash, err = sha512.Hash(append(hash, salt...))
	if err != nil {
		return nil, err
	}

	tempHash := hash
	tempSalt := salt
	hashs := make([]byte, 0)
	for i := 0; i < 10; i++ {
		tempHash, err = sha512.Hash(append(tempHash, byte(i)))
		if err != nil {
			return nil, err
		}

		tempSalt, err = sha512.Hash(append(tempSalt, byte(i)))
		if err != nil {
			return nil, err
		}

		tempHash, err = sha512.Hash(append(tempHash, tempSalt...))
		if err != nil {
			return nil, err
		}

		hashs = append(hashs, tempHash...)
	}

	hash, err = sha512.Hash(hashs)
	if err != nil {
		return nil, err
	}

	idxHash := len(hash) / 2
	hash1 := append([]byte{}, hash[:idxHash]...)
	hash2 := append([]byte{}, hash[idxHash:]...)

	idxSalt := len(salt) / 2
	salt1 := salt[:idxSalt]
	salt2 := salt[idxSalt:]

	for i := 0; i < loop; i++ {
		hash1, err = sha512.Hash(append(hash1, salt1...))
		if err != nil {
			return nil, err
		}

		hash2, err = sha512.Hash(append(hash2, salt2...))
		if err != nil {
			return nil, err
		}
	}

	hash, err = sha512.Hash(append(hash1, hash2...))
	if err != nil {
		return nil, err
	}

	hash, err = sha512.Hash(append(hash, []byte("BcryptSimple")...))
	if err != nil {
		return nil, err
	}

	return hash, nil
}
