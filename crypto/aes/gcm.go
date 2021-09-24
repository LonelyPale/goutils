package aes

import (
	"crypto/aes"
	"crypto/cipher"
)

func GCMEncrypt(secretKey, nonce, plaintext, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//// 向量
	//nonce := make([]byte, aesGcm.NonceSize())
	//if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	//	return nil, err
	//}

	cipherText := aesGcm.Seal(nil, nonce, plaintext, additionalData)

	return cipherText, nil
}

func GCMDecrypt(secretKey, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	c, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	dataBytes, err := gcm.Open(nil, nonce, ciphertext, additionalData)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}
