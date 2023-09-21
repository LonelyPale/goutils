package types

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestTime(t *testing.T) {
	fmt.Println("Local:", time.Local.String())

	testTimeJson(t)
	testTimeBson(t)
	//testTimeBsonMarshal(t) //有错误，未支持
	//testTimeBsonUnmarshal(t) //有错误，未支持
	//testTimeBsonMarshalOther(t)
}

func testTimeJson(t *testing.T) {
	dt1 := Now()
	bs1, err := json.Marshal(dt1)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, bs1)

	var dt2 Time
	if err := json.Unmarshal(bs1, &dt2); err != nil {
		t.Fatal(err)
	}

	bs2, err := json.Marshal(dt2)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, bs2)
	assert.Equal(t, bs1, bs2)

	var dt3 Time
	if err := json.Unmarshal(bs2, &dt3); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dt2, dt3)

	t.Log(string(bs1), bs1)
	t.Log(string(bs2), bs2)
	t.Log(dt1)
	t.Log(dt2)
	t.Log(dt3)
}

func testTimeBson(t *testing.T) {
	dt1 := Now()
	bs1, err := dt1.MarshalBSON()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, bs1)

	var dt2 Time
	if err := dt2.UnmarshalBSON(bs1); err != nil {
		t.Fatal(err)
	}

	bs2, err := dt2.MarshalBSON()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, bs2)
	assert.Equal(t, bs1, bs2)

	var dt3 Time
	if err := dt3.UnmarshalBSON(bs2); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dt2, dt3)

	t.Log(string(bs1), bs1)
	t.Log(string(bs2), bs2)
	t.Log(dt1)
	t.Log(dt2)
	t.Log(dt3)
}

func testTimeBsonMarshal(t *testing.T) {
	dt := Now()
	bs, err := bson.Marshal(dt)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(bs)
}

func testTimeBsonUnmarshal(t *testing.T) {
	data := []byte{98, 68, 173, 88, 120, 1, 0, 0}
	var dt Time

	if err := bson.Unmarshal(data, &dt); err != nil {
		t.Fatal(err)
	}

	t.Log(dt)
}

func testTimeBsonMarshalOther(t *testing.T) {
	dt := time.Now()
	_, bs, err := bson.MarshalValue(dt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bs)
}
