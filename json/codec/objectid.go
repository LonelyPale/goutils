// Created by LonelyPale at 2018-11-30

package codec

import (
	"unsafe"

	"github.com/json-iterator/go"

	"github.com/LonelyPale/goutils/database/mongodb/types"
)

const RegisterTypeObjectID = "primitive.ObjectID"

func RegisterObjectIDCodec() {
	jsoniter.RegisterTypeEncoder(RegisterTypeObjectID, &objectIDCodec{})
	jsoniter.RegisterTypeDecoder(RegisterTypeObjectID, &objectIDCodec{})
}

type objectIDCodec struct {
}

func (codec *objectIDCodec) IsEmpty(ptr unsafe.Pointer) bool {
	obj := *((*types.ObjectID)(ptr))
	return obj.IsZero()
}

func (codec *objectIDCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	obj := *((*types.ObjectID)(ptr))
	str := "\"" + obj.Hex() + "\""
	if _, err := stream.Write([]byte(str)); err != nil {
		panic(err)
	}
}

func (codec *objectIDCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	str := iter.ReadString()
	id, err := types.ObjectIDFromHex(str)
	if err != nil {
		panic(err)
	}
	*((*types.ObjectID)(ptr)) = id
}
