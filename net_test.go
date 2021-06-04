package goutils

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestMac(t *testing.T) {
	//bytes := []byte{0x78, 0x56, 0x34, 0x12}
	bytes := []byte{0x0, 0x0, 0x0, 0x12}
	fmt.Printf("0x%x\n", binary.LittleEndian.Uint32(bytes))
	fmt.Printf("0x%x\n", binary.BigEndian.Uint32(bytes))

	t.Log(Mac())
}
