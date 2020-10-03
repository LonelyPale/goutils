package qrcode

import (
	"testing"
)

func TestWriteFile(t *testing.T) {
	opts := DefaultOptions()
	opts.Size = 200
	//content := "1234567890"
	content := "abcdefghijklmnopqrstuvwxyz"

	if err := WriteFile("test.png", content, opts); err != nil {
		t.Fatal(err)
	}

	if err := WriteFile("test.jpg", content, opts); err != nil {
		t.Fatal(err)
	}
}
