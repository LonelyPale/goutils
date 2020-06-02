package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

//对明文进行填充
func Padding(plainText []byte, blockSize int) []byte {
	n := blockSize - len(plainText)%blockSize //计算要填充的长度
	temp := bytes.Repeat([]byte{byte(n)}, n)  //对原来的明文填充n个n
	plainText = append(plainText, temp...)
	return plainText
}

//对密文删除填充
func UnPadding(cipherText []byte) []byte {
	end := cipherText[len(cipherText)-1]               //取出密文最后一个字节end
	cipherText = cipherText[:len(cipherText)-int(end)] //删除填充
	return cipherText
}

//AES加密（CBC模式）
//指定初始向量vi,长度和block的块尺寸一致
func AESCBCEncrypt(plainText []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) //指定加密算法，返回一个AES算法的Block接口对象
	if err != nil {
		return nil, err
	}

	plainText = Padding(plainText, block.BlockSize()) //进行填充
	blockMode := cipher.NewCBCEncrypter(block, iv)    //指定分组模式，返回一个BlockMode接口对象
	cipherText := make([]byte, len(plainText))        //加密连续数据库
	blockMode.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}

//AES解密（CBC模式）
//指定初始化向量IV,和加密的一致
func AESCBCDecrypt(cipherText []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) //指定解密算法，返回一个AES算法的Block接口对象
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv) //指定分组模式，返回一个BlockMode接口对象
	plainText := make([]byte, len(cipherText))     //解密
	blockMode.CryptBlocks(plainText, cipherText)
	plainText = UnPadding(plainText) //删除填充

	return plainText, nil
}
