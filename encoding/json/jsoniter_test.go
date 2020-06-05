package json

import (
	"testing"
	"time"

	"github.com/json-iterator/go"
)

func Test(t *testing.T) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	type carinfo struct {
		Id      string
		License string
		Color   int
		Device  string // "<设备类型>.<设备id>"
	}

	type carlist struct {
		Result  int
		Message string
		Cars    []carinfo
	}

	startTime := time.Now()

	var msg carlist
	jsonstr := `{"result":0, "message":"ok", "cars":[{"id":"311111", "license":"豫A1111", "color":2, "device":"VA3K.10001"}, {"id":"311112", "license":"豫A1112", "color":2, "device":"VA3K.10002"}]}`
	err := json.Unmarshal([]byte(jsonstr), &msg)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("jsoniter Unmarshal duration:", time.Since(startTime))

	t.Log(msg)
}

func Test0(t *testing.T) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	type User struct {
		Name string
		Age  int
	}

	startTime := time.Now()

	var msg User
	jsonstr := `{"name":"hello beijing.","age":18}`
	err := json.Unmarshal([]byte(jsonstr), &msg)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("jsoniter Unmarshal duration:", time.Since(startTime))
	t.Log(msg)
}

func Test1(t *testing.T) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	type User struct {
		Name string
		Age  int
	}

	startTime := time.Now()

	var msg User
	msg.Name = "hello"
	msg.Age = 18
	bs, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("jsoniter Marshal duration:", time.Since(startTime))
	t.Log(string(bs))
}
