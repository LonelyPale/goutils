package msgpack

import (
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func Test(t *testing.T) {
	type User struct {
		Id   string
		Name string `msgpack:"-"`
	}

	v := User{
		Id:   "123",
		Name: "tom",
	}

	bs, err := msgpack.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	var user User
	if err := msgpack.Unmarshal(bs, &user); err != nil {
		t.Fatal(err)
	}

	t.Log(user)
}

func Test2(t *testing.T) {
	s := "123"
	bs, err := msgpack.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bs)

	ss := ""
	if err := msgpack.Unmarshal(bs, &ss); err != nil {
		t.Fatal(err)
	}
	t.Log(ss)
}
