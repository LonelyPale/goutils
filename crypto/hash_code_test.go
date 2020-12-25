package crypto

import "testing"

func TestHashCode(t *testing.T) {
	address := "bm1q5u8u4eldhjf3lvnkmyl78jj8a75neuryzlknk0" + "bm1q5u8u4eldhjf3lvnkmyl78jj8a75neuryzlknk0bm1q5u8u4eldhjf3lvnkmyl78jj8a75neuryzlknk0"
	t.Log(HashCode(address + "0"))
	t.Log(HashCode(address + "1"))
	t.Log(HashCode(address + "2"))

	t.Log(HashCode("address" + "0"))
	t.Log(HashCode("address" + "1"))
	t.Log(HashCode("address" + "2"))

}
