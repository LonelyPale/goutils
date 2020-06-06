package crypto

import "testing"

func TestNewHash(t *testing.T) {
	hash := NewHash()

	if err := hash.FromHex("0010FF"); err != nil {
		t.Fatal(err)
	}

	t.Log(hash.Hex())
	t.Log(hash.Base64())
	t.Log(hash.Bytes())

	if err := hash.FromBase64("ABD_"); err != nil {
		t.Fatal(err)
	}

	t.Log(hash.Hex())
	t.Log(hash.Base64())
	t.Log(hash.Bytes())
}
