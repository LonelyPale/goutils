package mongodb

// 验证tag: 增删查改 CRUD
const (
	CreateTagName = "vCreate" //增
	UpdateTagName = "vUpdate" //改
	DeleteTagName = "vDelete" //删
	ReadTagName   = "vRead"   //查
)

type ModelValidator interface {
	Validate(tags ...string) error
}

func validate(doc interface{}, tags ...string) error {
	if doc == nil {
		return ErrNilDocument
	}

	mod, ok := doc.(ModelValidator)
	if ok {
		return mod.Validate(tags...)
	}

	return nil
}
