package ipfs

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/LonelyPale/goutils/crypto"
)

var (
	client = NewClient("localhost:15001")
)

func TestNewClient(t *testing.T) {
	id, err := client.ID()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

	version, commit, err := client.Version()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(version)
	t.Log(commit)
}

func TestAdd(t *testing.T) {
	//QmeV1kwh3333bsnT6YRfdCRrSgUPngKmAhhTa4RrqYPbKT
	res, err := client.Add(strings.NewReader("hello world!\n"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)

	//QmXcQrx1CPanVhXGgBvQwPr9qCTmknDueQqeNyaMKsRr9r
	res1, err := client.Add(strings.NewReader("/Users/wyb/project/github/goutils/net/ipfs/client.go\n"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res1)

	//QmeSqn7sfnEKkLqcPNVxbsNvbxup5tFGNf1FkHeNoF6scD
	hash, err := client.Add(bytes.NewBufferString("你好，中国！"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}

func TestCat(t *testing.T) {
	reader, err := client.Cat("QmeSqn7sfnEKkLqcPNVxbsNvbxup5tFGNf1FkHeNoF6scD")
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(body)
	t.Log(string(body))
}

func TestEncryptDecrypt(t *testing.T) {
	key, err := crypto.GenerateSecretKey()
	if err != nil {
		t.Fatal(err)
	}

	hash, err := client.AddEncrypt([]byte("你好，中国！"), key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)

	bs, err := client.CatDecrypt(hash, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func TestBigFile(t *testing.T) {
	t1 := time.Now()
	bs, err := ioutil.ReadFile("/Users/wyb/backup/software/os/ubuntu-20.04.2-live-server-amd64.iso")
	if err != nil {
		t.Fatal(err)
	}
	elapsed := time.Since(t1)
	t.Log(len(bs), elapsed)

	key, err := crypto.GenerateSecretKey()
	if err != nil {
		t.Fatal(err)
	}

	t1 = time.Now()
	hash, err := client.AddEncrypt(bs, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash, time.Since(t1))

	t1 = time.Now()
	bs2, err := client.CatDecrypt(hash, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(bs2), time.Since(t1))

	t1 = time.Now()
	if err := ioutil.WriteFile("/Users/wyb/backup/software/os/test.tmp", bs2, 0755); err != nil {
		t.Fatal(err)
	}
	t.Log(time.Since(t1))
}
