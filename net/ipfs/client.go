package ipfs

import (
	"bytes"
	"io/ioutil"

	shell "github.com/ipfs/go-ipfs-api"

	"github.com/lonelypale/goutils/crypto/aes"
)

type Client struct {
	*shell.Shell
}

func NewClient(baseURLs ...string) *Client {
	if len(baseURLs) > 0 && len(baseURLs[0]) > 0 {
		return &Client{shell.NewShell(baseURLs[0])}
	} else {
		return &Client{shell.NewLocalShell()}
	}
}

// 用 aes256 cbc 加密后存储到 ipfs 上，返回 hash
func (c *Client) AddEncrypt(plain []byte, key []byte) (string, error) {
	bs, err := aes.Encrypt(plain, key)
	if err != nil {
		return "", err
	}

	return c.Add(bytes.NewBuffer(bs))
}

// 根据 hash，从 ipfs 上取出后用 aes256 cbc 解密
func (c *Client) CatDecrypt(hash string, key []byte) ([]byte, error) {
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
