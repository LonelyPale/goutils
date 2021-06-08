package disk

import (
	"fmt"
	"testing"
)

func TestDiskSize(t *testing.T) {
	fmt.Println("KB", KB)
	fmt.Println("MB", MB)
	fmt.Println("GB", GB)
	fmt.Println("TB", TB)
	fmt.Println("PB", PB)
	fmt.Println("EB", EB)
	//fmt.Println("ZB", ZB) //constant 1180591620717411303424 overflows int
	//fmt.Println("YB", YB) //constant 1208925819614629174706176 overflows int

	fmt.Println()
	fmt.Println(YB / ZB)
	fmt.Println("1180591620717411303424 B = ", 1180591620717411303424/ZB, "ZB")
}

func TestFree(t *testing.T) {
	free, err := Free("/")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(free/GB, "GB")
}
