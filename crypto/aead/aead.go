package aead

import (
	"crypto/sha256"
	"io"

	"github.com/minio/sio"
	"golang.org/x/crypto/hkdf"
)

const (
	keySize = 32
)

type Config struct {
	Secret []byte
	Salt   []byte
}

func EncryptReader(src io.Reader, cfg Config) (io.Reader, error) {
	key, err := deriveSecretKey(cfg.Secret, cfg.Salt)
	if err != nil {
		return nil, err
	}
	return sio.EncryptReader(src, sio.Config{Key: key})
}

func DecryptReader(src io.Reader, cfg Config) (io.Reader, error) {
	key, err := deriveSecretKey(cfg.Secret, cfg.Salt)
	if err != nil {
		return nil, err
	}
	return sio.DecryptReader(src, sio.Config{Key: key})
}

func deriveSecretKey(secret, salt []byte) ([]byte, error) {
	// derive an encryption key from the master key and the nonce
	var key [keySize]byte
	kdf := hkdf.New(sha256.New, secret, salt, nil)
	if _, err := io.ReadFull(kdf, key[:]); err != nil {
		return nil, err
	}
	return key[:], nil
}
