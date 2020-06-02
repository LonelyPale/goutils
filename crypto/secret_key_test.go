package crypto

import (
	"testing"
	"time"
)

func TestGenerateSecretKey(t *testing.T) {
	for i := 0; i < 10; i++ {
		startTime := time.Now()

		key, err := GenerateSecretKey(16)
		if err != nil {
			t.Fatal(err)
		}

		if len(key) != 16 {
			t.Fatal("gen length error")
		}

		t.Log("duration", time.Since(startTime), key)
	}
}
