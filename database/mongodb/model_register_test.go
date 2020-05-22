package mongodb

import "testing"

func TestModelRegister(t *testing.T) {
	RegisterType((*Model)(nil))

	for k := range typeRegistry {
		t.Log(k)
	}

	model := MakeInstance("github.com/LonelyPale/goutils/database/mongodb.Model")
	t.Log(model)
}
