// Created by LonelyPale at 2018-11-30

package json

import "github.com/LonelyPale/goutils/encoding/json/codec"

// register jsoniter custom plugins
func init() {
	codec.RegisterTimeAsFormartCodec()
	codec.RegisterTypesTimeAsFormartCodec()
	//codec.RegisterObjectIDCodec() //mongo-driver已原生实现
	//codec.RegisterDecimal128Codec() //暂未使用
}
