package mongodb

// 验证tag: 增删查改 CRUD
const (
	CreateTagName = "vCreate"
	UpdateTagName = "vUpdate"
	DeleteTagName = "vDelete"
	ReadTagName   = "vRead"
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
