package ipfs

import (
	"bytes"
	"io"
	"io/ioutil"

	shell "github.com/ipfs/go-ipfs-api"

	"github.com/lonelypale/goutils/crypto/aead"
	"github.com/lonelypale/goutils/crypto/aes"
)

type Client struct {
	*shell.Shell
	Cfg *Config
}

func NewClient(cfg *Config) *Client {
	if cfg != nil && cfg.URI != "" {
		return &Client{
			Shell: shell.NewShell(cfg.URI),
			Cfg:   cfg,
		}
	} else {
		return &Client{
			Shell: shell.NewLocalShell(),
			Cfg:   cfg,
		}
	}
}

func (c *Client) AddEncrypt(src io.Reader, cfg aead.Config) (string, error) {
	encrypted, err := aead.EncryptReader(src, cfg)
	if err != nil {
		return "", err
	}

	return c.Add(encrypted)
}

func (c *Client) CatDecrypt(hash string, cfg aead.Config) (io.Reader, error) {
	reader, err := c.Cat(hash)
	if err != nil {
		return nil, err
	}

	return aead.DecryptReader(reader, cfg)
}

// 用 aes256 cbc 加密后存储到 ipfs 上，返回 hash
func (c *Client) AddEncryptBuffer(plain []byte, key []byte) (string, error) {
	bs, err := aes.Encrypt(plain, key)
	if err != nil {
		return "", err
	}

	return c.Add(bytes.NewBuffer(bs))
}

// 根据 hash，从 ipfs 上取出后用 aes256 cbc 解密
func (c *Client) CatDecryptBuffer(hash string, key []byte) ([]byte, error) {
	reader, err := c.Cat(hash)
	if err != nil {
		return nil, err
	}

	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return aes.Decrypt(bs, key)
}
