package cache

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lonelypale/goutils/types"
)

type user struct {
	Name string
	Age  int
	Data map[string]interface{}
}

func TestCache(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Fatal(err)
	}

	u := user{
		Name: "Tom",
		Age:  10,
		Data: types.M{"phone": 1234567, "email": "tom@mail.com"},
	}
	//test(&u)

	if err := cache.Set("tom", u); err != nil {
		t.Fatal(err)
	}

	uu := new(user)
	if err := cache.Get("tom", uu); err != nil {
		t.Fatal(err)
	}
	t.Log(uu)

	m := make(types.M)
	m["321"] = 123
	if err := cache.Get("tom", &m); err != nil {
		t.Fatal(err)
	}
	t.Log(m)

	tom, err := cache.gettest("tom")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tom)
}

func test(v interface{}) {
	typ := reflect.TypeOf(v)
	fmt.Println(typ.String(), typ.Name(), typ.PkgPath())
	fmt.Println(typ.Elem().String(), typ.Elem().Name(), typ.Elem().PkgPath())

	//elem := reflect.Indirect(reflect.New(typ)).Addr()
}
