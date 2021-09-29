package aead

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/lonelypale/goutils/debug"
	"github.com/lonelypale/goutils/random"
)

func TestAEAD(t *testing.T) {
	startTime := time.Now()
	debug.TraceMemStats()

	masterkey, err := random.Random(32)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("masterkey:", time.Since(startTime), masterkey.Hex())

	nonce, err := random.Random(32)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nonce:", time.Since(startTime), nonce.Hex())

	// 1,215,168,512
	inFile, err := os.Open("/Users/wyb/backup/software/os/ubuntu-20.04.2-live-server-amd64.iso")
	if err != nil {
		t.Fatal(err)
	}
	defer inFile.Close()
	t.Log("Open duration:", time.Since(startTime))

	hash := sha256.New()
	if _, err := io.Copy(hash, inFile); err != nil {
		t.Fatal(err)
	}
	t.Log("sha256 duration:", time.Since(startTime), fmt.Sprintf("%x", hash.Sum(nil)))

	ret, err := inFile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Seek duration:", time.Since(startTime), ret)

	cfg := Config{
		Secret: masterkey,
		Salt:   nonce,
	}
	encrypted, err := EncryptReader(inFile, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("EncryptReader duration:", time.Since(startTime))

	decrypted, err := DecryptReader(encrypted, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("DecryptReader duration:", time.Since(startTime))

	// 1,215,168,512
	outFile, err := os.Create("/Users/wyb/backup/software/os/test.tmp")
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()
	t.Log("Create duration:", time.Since(startTime))

	num, err := io.Copy(outFile, decrypted)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Copy duration:", time.Since(startTime), num)

	time.Sleep(time.Second * 10)
}

func TestDeriveSecretKey(t *testing.T) {
	keyStr := "67a1b2010bf968076bb5f5fd01d150b5fe544aba54199df0e80576048a86e1e9"
	secret := []byte{1, 2, 3}
	salt := []byte{3, 2, 1}
	key, err := deriveSecretKey(secret, salt)
	if err != nil {
		t.Fatal(err)
	}

	hexStr := hex.EncodeToString(key)
	t.Log(hexStr, hexStr == keyStr)
}
