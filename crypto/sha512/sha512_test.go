package sha512

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	startTime := time.Now()
	for i := 0; i < 1000; i++ {
		_, err := Hash([]byte("123qweasdzxc"))
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Log("duration:", time.Since(startTime))
}

func TestHash(t *testing.T) {
	hashed, err := Hash([]byte("123qweasdzxc"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashed)
	t.Log(string(hashed))
}

func TestHex(t *testing.T) {
	hashed, err := Hex([]byte("123qweasdzxc"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashed)
}

func TestBase64(t *testing.T) {
	hashed, err := Base64([]byte("123qweasdzxc"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashed)
}
