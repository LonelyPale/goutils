package redis

import (
	"testing"
	"time"
)

func TestNewRedisDB(t *testing.T) {
	conf := DefaultConfig()
	db, err := NewRedisDB(conf)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Set("keyString", "abc!123,杭州。", time.Minute); err != nil {
		t.Error(err)
	}
	var keyString string
	if err := db.Get("keyString", &keyString); err != nil {
		t.Error(err)
	}
	t.Log(keyString)

	type name struct {
		One string
		Two string
		Num int
	}
	n1 := &name{
		One: "abc123",
		Two: "def中国！",
		Num: 888,
	}
	if err := db.Set("name", n1, time.Minute); err != nil {
		t.Error(err)
	}
	var n2 *name
	if err := db.Get("name", &n2); err != nil {
		t.Error(err)
	}
	t.Log(n2)
}
