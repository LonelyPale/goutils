package springweb

type ResultHandler interface {
	Invoke(result []interface{}) []interface{}
}

type ResultHandlerFunc func(result []interface{}) []interface{}

func (f ResultHandlerFunc) Invoke(result []interface{}) []interface{} {
	return f(result)
}

var resultHandlers = make([]ResultHandler, 0)

func RegisterResultHandler(handler ResultHandler) {
	resultHandlers = append(resultHandlers, handler)
}

func RegisterResultHandlerFunc(handler ResultHandlerFunc) {
	resultHandlers = append(resultHandlers, handler)
}
