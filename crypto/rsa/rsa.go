package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"io"

	"github.com/LonelyPale/goutils/errors"
)

const (
	CharSet          = "UTF-8"
	Base64Format     = "UrlSafeNoPadding"
	AlgorithmKeyType = "PKCS1"
	AlgorithmSign    = crypto.SHA256
	Bits             = 2048
)

var Base64Encoding = base64.StdEncoding

type XRsa struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func CreateKeysToBase64(keyLengths ...int) (string, string, error) {
	pubKey := bytes.NewBuffer(make([]byte, 0))
	priKey := bytes.NewBuffer(make([]byte, 0))
	if err := CreateKeys(pubKey, priKey, keyLengths...); err != nil {
		return "", "", err
	}

	return pubKey.String(), priKey.String(), nil
}

func CreateKeysToBytes(keyLengths ...int) ([]byte, []byte, error) {
	pubKey := bytes.NewBuffer(make([]byte, 0))
	priKey := bytes.NewBuffer(make([]byte, 0))
	if err := CreateKeys(pubKey, priKey, keyLengths...); err != nil {
		return nil, nil, err
	}

	return pubKey.Bytes(), priKey.Bytes(), nil
}

// 生成密钥对
func CreateKeys(publicKeyWriter, privateKeyWriter io.Writer, keyLengths ...int) error {
	var bits int
	if len(keyLengths) > 0 && keyLengths[0] > 0 {
		bits = keyLengths[0]
	} else {
		bits = Bits
	}

	// 生成私钥文件, bits位数长度
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(privateKeyWriter, block)
	if err != nil {
		return err
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	err = pem.Encode(publicKeyWriter, block)
	if err != nil {
		return err
	}

	return nil
}

func NewXRsa(publicKey []byte, privateKey []byte) (*XRsa, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub := pubInterface.(*rsa.PublicKey)
	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}

	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &XRsa{
		publicKey:  pub,
		privateKey: pri,
	}, nil
}

// 公钥加密
func (r *XRsa) Encrypt(data []byte) ([]byte, error) {
	partLen := r.publicKey.N.BitLen()/8 - 11
	chunks := split(data, partLen)
	buffer := bytes.NewBuffer([]byte{})
	for _, chunk := range chunks {
		encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, chunk)
		if err != nil {
			return nil, err
		}

		buffer.Write(encrypted)
	}
	return buffer.Bytes(), nil
}

// 私钥解密
func (r *XRsa) Decrypt(encrypted []byte) ([]byte, error) {
	partLen := r.publicKey.N.BitLen() / 8
	chunks := split(encrypted, partLen)
	buffer := bytes.NewBuffer([]byte{})
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, chunk)
		if err != nil {
			return nil, err
		}

		buffer.Write(decrypted)
	}
	return buffer.Bytes(), nil
}

// 私钥签名
func (r *XRsa) Sign(data []byte) ([]byte, error) {
	h := AlgorithmSign.New()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}

	hashed := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, AlgorithmSign, hashed)
	if err != nil {
		return nil, err
	}

	return sign, nil
}

// 公钥验签
func (r *XRsa) Verify(data []byte, sign []byte) error {
	h := AlgorithmSign.New()
	if _, err := h.Write(data); err != nil {
		return err
	}

	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(r.publicKey, AlgorithmSign, hashed, sign)
}

func (r *XRsa) EncryptToBase64(data []byte) (string, error) {
	bs, err := r.Encrypt(data)
	if err != nil {
		return "", err
	}

	return Base64Encoding.EncodeToString(bs), nil
}

func (r *XRsa) DecryptFromBase64(encrypted string) ([]byte, error) {
	raw, err := Base64Encoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	return r.Decrypt(raw)
}

// 数据签名
func (r *XRsa) SignToBase64(data []byte) (string, error) {
	sign, err := r.Sign(data)
	if err != nil {
		return "", err
	}

	return Base64Encoding.EncodeToString(sign), nil
}

func (r *XRsa) VerifyFromBase64(data []byte, sign string) error {
	decodedSign, err := Base64Encoding.DecodeString(sign)
	if err != nil {
		return err
	}

	return r.Verify(data, decodedSign)
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}

func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) ([]byte, error) {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}

	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)

	k, err := asn1.Marshal(info)
	if err != nil {
		return nil, err
	}

	return k, nil
}
