package validator

func init() {
	//SpringWeb.Validator = SpringWeb.NewDefaultValidator()
	if err := DefaultValidator.SetLanguage(ZH); err != nil {
		panic(err)
	}
}

// Validator 参数校验器接口
type Validator interface {
	Engine() interface{}
	Validate(obj interface{}, tags ...string) error
	Var(field interface{}, tag string) error
	SetLanguage(language Language) error
}

// Validator 全局参数校验器
var DefaultValidator Validator = NewDefaultValidator()

// Validate 验证 struct 结构体
func Validate(obj interface{}, tags ...string) error {
	if DefaultValidator != nil {
		return DefaultValidator.Validate(obj, tags...)
	}
	return nil
}

// Var 验证 var 变量
func Var(field interface{}, tag string) error {
	if DefaultValidator != nil {
		return DefaultValidator.Var(field, tag)
	}
	return nil
}
