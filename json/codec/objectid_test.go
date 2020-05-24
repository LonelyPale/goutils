// Created by LonelyPale at 2018-12-11

package codec

import (
	"testing"

	"github.com/json-iterator/go"

	"github.com/LonelyPale/goutils/database/mongodb/types"
)

func TestObjectIDCodec_IsEmpty(t *testing.T) {
	RegisterObjectIDCodec()
	objId := types.NilObjectID
	output, err := jsoniter.Marshal(objId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(output))
}

func TestObjectIDCodec_Encode(t *testing.T) {
	RegisterObjectIDCodec()
	objId := types.NewObjectID()
	t.Log(1, objId)
	output, err := jsoniter.Marshal(objId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(2, string(output))
}

func TestObjectIDCodec_Decode(t *testing.T) {
	RegisterObjectIDCodec()
	objId := "\"5c0ef2732f20398f4b9c5f5f\""
	t.Log(1, objId)
	var val types.ObjectID
	err := jsoniter.Unmarshal([]byte(objId), &val)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(2, val)
}
