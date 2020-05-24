// Created by LonelyPale at 2018-11-30

package json

import . "github.com/LonelyPale/goutils/json/codec"

func init() {
	RegisterTimeAsFormartCodec()
	RegisterObjectIDCodec()
	//RegisterDecimal128Codec()
}
