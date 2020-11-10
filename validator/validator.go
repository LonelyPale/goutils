package validator

// Validator 参数校验器接口
type Validator interface {
	Engine() interface{}
	Validate(obj interface{}, tags ...string) error
	SetLanguage(language Language) error
}

// Validator 全局参数校验器
var DefaultValidator Validator = NewDefaultValidator()

// Validate 参数校验
func Validate(obj interface{}, tags ...string) error {
	if DefaultValidator != nil {
		return DefaultValidator.Validate(obj, tags...)
	}
	return nil
}
