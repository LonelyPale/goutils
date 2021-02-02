package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/LonelyPale/goutils/errors"
)

//todo: 优化base64等输出格式
//todo: 现在是前置iv，增加后置iv
//todo: 加入salt，把不标准的key生成统一的key

// AES的区块长度固定为128比特，向量iv长度固定为128比特，密钥长度固定为128、192、256比特。
// 标准的AES是没有salt概念的，自定义时可以扩展salt用来生成长度不标准的key的hash，用于统一key的长度。
//          key  iv   block
// AES-128  128  128  128
// AES-192  192  128  128
// AES-256  256  128  128

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
func CBCEncrypt(plainText []byte, key []byte, iv []byte) ([]byte, error) {
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
func CBCDecrypt(cipherText []byte, key []byte, iv []byte) ([]byte, error) {
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

func Encrypt(plainText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) //指定加密算法，返回一个AES算法的Block接口对象
	if err != nil {
		return nil, err
	}

	plainText = Padding(plainText, block.BlockSize()) //进行填充

	//对IV有随机性要求，但没有保密性要求，所以常见的做法是将IV包含在加密文本当中
	cipherText := make([]byte, aes.BlockSize+len(plainText)) //加密连续数据库
	//随机一个block大小作为IV
	//采用不同的IV时相同的秘钥将会产生不同的密文，可以理解为一次加密的session
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCEncrypter(block, iv) //指定分组模式，返回一个BlockMode接口对象
	blockMode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

func Decrypt(cipherText []byte, key []byte) ([]byte, error) {
	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("cipherText too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	return CBCDecrypt(cipherText, key, iv)
}
