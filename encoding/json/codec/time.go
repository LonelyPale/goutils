// Created by LonelyPale at 2018-11-30

package codec

import (
	"time"
	"unsafe"

	"github.com/json-iterator/go"
)

// RFC3339 = "2006-01-02T15:04:05Z07:00"
// golang 自带的 time package，默认时间格式是RFC3339，不支持 json 序列化和反序列化时自定义时间格式。
// type 别名和 struct 嵌套继承的方法有侵入性，需要修改原本的类型。
// struct 转换的方法又太繁琐，每种类型的 struct 都需要处理。
// 用 jsoniter 自定义类型插件序列化和反序列化能比较好的处理以上问题，只是需要引入第三方包。

// 1: type Time time.Time
// 2: type Time struct { time.Time }
// 3:
//type Alias Person
//return json.Marshal(&struct {
//	Alias
//	CreateTime string `json:"create_time"`
//}{
//	Alias:      (Alias)(d),
//	CreateTime: d.CreateTime.Format("2006/01/02 15:04:05"),
//})

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
