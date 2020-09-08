package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"testing"
)

// ssh-keygen -t ed25519 -f my_github_ed25519  -C "me@github"

func Test(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("PublicKey:", hex.EncodeToString(publicKey))
	t.Log("PrivateKey:", hex.EncodeToString(privateKey))
}
