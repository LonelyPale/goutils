package types

import "testing"

func TestNewHash(t *testing.T) {
	bytes := NewBytes()

	if err := bytes.FromHex("0010FF"); err != nil {
		t.Fatal(err)
	}

	t.Log(bytes.Hex())
	t.Log(bytes.Base64())
	t.Log(bytes.Bytes())

	if err := bytes.FromBase64("ABD_"); err != nil {
		t.Fatal(err)
	}

	t.Log(bytes.Hex())
	t.Log(bytes.Base64())
	t.Log(bytes.Bytes())
}
