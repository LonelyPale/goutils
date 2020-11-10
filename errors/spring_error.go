package errors

import SpringUtils "github.com/go-spring/spring-utils"

// Panic 抛出一个异常值
func Panic(err error) *SpringUtils.PanicCond {
	return SpringUtils.NewPanicCond(func() interface{} {
		return err
	})
}
