package encoding

import (
	"testing"

	"github.com/lonelypale/goutils/encoding/gob"
	"github.com/lonelypale/goutils/encoding/msgpack"
)

type userStruct struct {
	Name string
	Age  int
}

func Test1(t *testing.T) {
	user := &userStruct{"tom", 10}

	for i := 0; i < 10; i++ {
		_, err := gob.Serialize(user)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 10; i++ {
		_, err := msgpack.Marshal(user)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Test2(t *testing.T) {
	user := &userStruct{"tom", 10}

	bs, err := gob.Serialize(user)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		_, err = gob.Deserialize(bs)
		if err != nil {
			t.Fatal(err)
		}
	}

	bs, err = msgpack.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		u := new(userStruct)
		if err := msgpack.Unmarshal(bs, u); err != nil {
			t.Fatal(err)
		}
	}
}
