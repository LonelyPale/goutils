package qrcode

import "testing"

func TestWriteFile(t *testing.T) {
	if err := WriteFile("test.png", "123"); err != nil {
		t.Fatal(err)
	}

	if err := WriteFile("test.jpg", "qwe"); err != nil {
		t.Fatal(err)
	}
}
