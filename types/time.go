package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	"github.com/lonelypale/goutils/errors"
)

var (
	UseLocalTimeZone  = true                  //是否使用本地时区
	DefaultTimeFormat = "2006-01-02 15:04:05" //默认时间格式化字符串
)

// Time extension time
type Time struct {
	time.Time
}

// Now returns the current local time.
func Now() Time {
	return Time{time.Now()}
}

// MarshalJSON marshal json
func (t Time) MarshalJSON() ([]byte, error) {
	bs := make([]byte, 0, len(DefaultTimeFormat)+2)
	bs = append(bs, '"')
	bs = t.Time.AppendFormat(bs, DefaultTimeFormat)
	bs = append(bs, '"')
	return bs, nil
}

// UnmarshalJSON unmarshal json
func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == `""` || string(data) == `null` {
		return nil
	}

	var now time.Time
	var err error
	switch UseLocalTimeZone {
	case true:
		now, err = time.ParseInLocation(`"`+DefaultTimeFormat+`"`, string(data), time.Local)
	case false:
		now, err = time.Parse(`"`+DefaultTimeFormat+`"`, string(data))
	}
	if err != nil {
		return err
	}

	t.Time = now
	return nil
}

// MarshalBSON marshal bson
func (t Time) MarshalBSON() ([]byte, error) {
	dt := t.UnixNano() / 1000000
	return writei64(dt), nil
}

// UnmarshalBSON unmarshal bson
func (t *Time) UnmarshalBSON(data []byte) error {
	dt := readi64(data)
	timeVal := time.Unix(dt/1000, dt%1000*1000000) //转换后是本地时区，解决 mongodb 保存(默认无时区)后再读出来时没有本地时区的问题。

	if !UseLocalTimeZone {
		timeVal = timeVal.UTC()
	}

	t.Time = timeVal
	return nil
}

// MarshalBSONValue marshal bson value
func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	bytes, err := t.MarshalBSON()
	return bson.TypeDateTime, bytes, err
	//return bson.MarshalValue(t.Time) //原生实现
}

func (t *Time) UnmarshalBSONValue(btype bsontype.Type, data []byte) error {
	if btype == bson.TypeDateTime {
		return t.UnmarshalBSON(data)
	}
	return errors.Errorf("unsupported data type %v into a types.Time", btype)
}

func readi64(src []byte) int64 {
	_ = src[7] // bounds check hint to compiler; see golang.org/issue/14808
	return int64(src[0]) | int64(src[1])<<8 | int64(src[2])<<16 | int64(src[3])<<24 |
		int64(src[4])<<32 | int64(src[5])<<40 | int64(src[6])<<48 | int64(src[7])<<56

}

func writei64(i64 int64) []byte {
	var dst []byte
	return append(dst, byte(i64), byte(i64>>8), byte(i64>>16), byte(i64>>24),
		byte(i64>>32), byte(i64>>40), byte(i64>>48), byte(i64>>56))
}
