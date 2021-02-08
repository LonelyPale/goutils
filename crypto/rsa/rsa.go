package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"

	"github.com/LonelyPale/goutils/errors"
	"github.com/LonelyPale/goutils/types"
)

type Type uint

//加密解密流程、签名验签流程
const (
	PKCS1v15 Type = 1 + iota // RSAES-PKCS1-v1_5、RSASSA-PKCS1-v1_5
	OAEP                     // RSAES-OAEP
	PSS                      // RSASSA-PSS
	maxRsaType
)

type KeyType uint

//公私钥格式
const (
	PKCS1 KeyType = 1 + iota //PKCS #1
	PKCS8                    //PKCS #8
	PKIX                     //公钥通用
	maxRsaKeyType
)

const (
	DefaultType    = PKCS1v15
	DefaultKeyType = PKCS1
	DefaultBits    = 2048
	DefaultHash    = crypto.SHA256
)

type Options struct {
	Type       Type
	KeyType    KeyType
	Bits       int
	Hash       crypto.Hash //RSA-OAEP RSA-PSS
	Label      []byte      //RSA-OAEP
	SaltLength int         //RSA-PSS
}

func DefaultOptions() *Options {
	return &Options{
		Type:       DefaultType,
		KeyType:    DefaultKeyType,
		Bits:       DefaultBits,
		Hash:       DefaultHash,
		Label:      nil,
		SaltLength: rsa.PSSSaltLengthAuto,
	}
}

type XRsa struct {
	opts       *Options
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// 生成密钥对(公私钥)
func CreateKeys(opts ...*Options) (types.Bytes, types.Bytes, error) {
	var opt *Options
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	} else {
		opt = DefaultOptions()
	}

	publicKeyWriter := bytes.NewBuffer(make([]byte, 0))
	privateKeyWriter := bytes.NewBuffer(make([]byte, 0))

	// 生成私钥, bits位数长度
	privateKey, err := rsa.GenerateKey(rand.Reader, opt.Bits)
	if err != nil {
		return nil, nil, err
	}

	var derStream []byte
	switch opt.KeyType {
	case PKCS1:
		derStream = x509.MarshalPKCS1PrivateKey(privateKey)
	case PKCS8:
		derStream, err = x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, errors.New("crypto/rsa: invalid KeyType for options")
	}

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(privateKeyWriter, block)
	if err != nil {
		return nil, nil, err
	}

	// 生成公钥
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}

	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	err = pem.Encode(publicKeyWriter, block)
	if err != nil {
		return nil, nil, err
	}

	return publicKeyWriter.Bytes(), privateKeyWriter.Bytes(), nil
}

// 解析密钥对(公私钥)
func NewXRsa(publicKey []byte, privateKey []byte, opts ...*Options) (*XRsa, error) {
	var opt *Options
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	} else {
		opt = DefaultOptions()
	}

	var pubKey *rsa.PublicKey
	var priKey *rsa.PrivateKey

	if len(publicKey) > 0 {
		block, _ := pem.Decode(publicKey)
		if block == nil {
			return nil, errors.New("public key error")
		}

		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		pubKey = pubInterface.(*rsa.PublicKey)
	}

	if len(privateKey) > 0 {
		block, _ := pem.Decode(privateKey)
		if block == nil {
			return nil, errors.New("private key error!")
		}

		var err error
		var priInterface interface{}
		switch opt.KeyType {
		case PKCS1:
			priKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		case PKCS8:
			priInterface, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			priKey = priInterface.(*rsa.PrivateKey)
		default:
			return nil, errors.New("crypto/rsa: invalid KeyType for options")
		}
	}

	return &XRsa{
		opts:       opt,
		publicKey:  pubKey,
		privateKey: priKey,
	}, nil
}

// 公钥加密 encrypted
func (x *XRsa) Encrypt(plaintext []byte, opts ...*Options) (types.Bytes, error) {
	var opt *Options
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	} else {
		opt = x.opts
	}

	//明文过长时分段
	var partLen int
	switch opt.Type {
	case PKCS1v15:
		partLen = x.publicKey.Size() - 11
	case OAEP:
		partLen = x.publicKey.Size() - 2*opt.Hash.Size() - 2
	default:
		return nil, errors.New("crypto/rsa: invalid Type for options")
	}
	chunks := split(plaintext, partLen)

	buffer := bytes.NewBuffer([]byte{})
	for _, chunk := range chunks {
		var ciphertext []byte
		var err error

		switch opt.Type {
		case PKCS1v15:
			ciphertext, err = rsa.EncryptPKCS1v15(rand.Reader, x.publicKey, chunk)
		case OAEP:
			ciphertext, err = rsa.EncryptOAEP(opt.Hash.New(), rand.Reader, x.publicKey, chunk, opt.Label)
		default:
			return nil, errors.New("crypto/rsa: invalid Type for options")
		}
		if err != nil {
			return nil, err
		}

		buffer.Write(ciphertext)
	}

	return buffer.Bytes(), nil
}

// 私钥解密  decrypted
func (x *XRsa) Decrypt(ciphertext []byte, opts ...*Options) (types.Bytes, error) {
	var opt *Options
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	} else {
		opt = x.opts
	}

	//密文分段，rsa(PKCS1v15、OAEP)密文段的长度和公钥的长度相等
	partLen := x.publicKey.Size()
	chunks := split(ciphertext, partLen)
	buffer := bytes.NewBuffer([]byte{})
	for _, chunk := range chunks {
		var plaintext []byte
		var err error

		switch opt.Type {
		case PKCS1v15:
			plaintext, err = rsa.DecryptPKCS1v15(rand.Reader, x.privateKey, chunk)
		case OAEP:
			plaintext, err = rsa.DecryptOAEP(opt.Hash.New(), rand.Reader, x.privateKey, chunk, opt.Label)
		default:
			return nil, errors.New("crypto/rsa: invalid Type for options")
		}
		if err != nil {
			return nil, err
		}

		buffer.Write(plaintext)
	}

	return buffer.Bytes(), nil
}

// 私钥签名 signature
func (x *XRsa) Sign(data []byte, opts ...*Options) (types.Bytes, error) {
	var opt *Options
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	} else {
		opt = x.opts
	}

	h := opt.Hash.New()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	hashed := h.Sum(nil)

	var signature []byte
	var err error
	switch opt.Type {
	case PKCS1v15:
		signature, err = rsa.SignPKCS1v15(rand.Reader, x.privateKey, opt.Hash, hashed)
	case PSS:
		signature, err = rsa.SignPSS(rand.Reader, x.privateKey, opt.Hash, hashed, &rsa.PSSOptions{
			SaltLength: opt.SaltLength,
			Hash:       opt.Hash,
		})
	default:
		return nil, errors.New("crypto/rsa: invalid Type for options")
	}
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// 公钥验签 verifies
func (x *XRsa) Verify(data []byte, sign []byte, opts ...*Options) error {
	var opt *Options
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	} else {
		opt = x.opts
	}

	h := opt.Hash.New()
	if _, err := h.Write(data); err != nil {
		return err
	}
	hashed := h.Sum(nil)

	switch opt.Type {
	case PKCS1v15:
		return rsa.VerifyPKCS1v15(x.publicKey, opt.Hash, hashed, sign)
	case PSS:
		return rsa.VerifyPSS(x.publicKey, opt.Hash, hashed, sign, &rsa.PSSOptions{
			SaltLength: opt.SaltLength,
			Hash:       opt.Hash,
		})
	default:
		return errors.New("crypto/rsa: invalid Type for options")
	}
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

//MarshalPKCS8PrivateKey
//Pkcs1ToPkcs8 将 pkcs1 PrivateKey 转到 pkcs8 PrivateKey 自定义
func Pkcs1ToPkcs8(key []byte) (types.Bytes, error) {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}

	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = key
	//info.PrivateKey = x509.MarshalPKCS1PrivateKey(key *rsa.PrivateKey)

	return asn1.Marshal(info)
}
