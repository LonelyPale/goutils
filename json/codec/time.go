// Created by LonelyPale at 2018-11-30

package codec

import (
	"time"
	"unsafe"

	"github.com/json-iterator/go"
)

const RegisterTypeTime = "time.Time"

// 2006-01-02 15:04:05.999
const DefaultTimeFormart = "2006-01-02 15:04:05"

func RegisterTimeAsFormartCodec(formarts ...string) {
	//jsoniter.RegisterTypeEncoder("time.Time", &timeAsInt64Codec{precision})
	//jsoniter.RegisterTypeDecoder("time.Time", &timeAsInt64Codec{precision})

	if len(formarts) > 0 && len(formarts[0]) > 0 {
		jsoniter.RegisterTypeEncoder(RegisterTypeTime, &timeAsFormartCodec{formarts[0]})
		jsoniter.RegisterTypeDecoder(RegisterTypeTime, &timeAsFormartCodec{formarts[0]})
	} else {
		jsoniter.RegisterTypeEncoder(RegisterTypeTime, &timeAsFormartCodec{DefaultTimeFormart})
		jsoniter.RegisterTypeDecoder(RegisterTypeTime, &timeAsFormartCodec{DefaultTimeFormart})
	}
}

type timeAsFormartCodec struct {
	timeFormart string
}

func (codec *timeAsFormartCodec) IsEmpty(ptr unsafe.Pointer) bool {
	ts := *((*time.Time)(ptr))
	return ts.IsZero()
}

func (codec *timeAsFormartCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := *((*time.Time)(ptr))
	str := "\"" + ts.Format(codec.timeFormart) + "\""
	if _, err := stream.Write([]byte(str)); err != nil {
		panic(err)
	}
}

func (codec *timeAsFormartCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	str := iter.ReadString()
	t, err := time.Parse(codec.timeFormart, str)
	if err != nil {
		panic(err)
	}
	*((*time.Time)(ptr)) = t
}
