// Created by LonelyPale at 2018-11-30

package codec

import (
	"unsafe"

	"github.com/json-iterator/go"

	"github.com/LonelyPale/goutils/types"
)

// 旧版本的 go.mongodb.org/mongo-driver 不支持完善的 json 序列化和反序列化，所以需要用 jsoniter plugins 自己处理。
// 新版本的 go.mongodb.org/mongo-driver 已支持完善的 json 序列化和反序列化，不再需要用 jsoniter plugins 自己处理。

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
