package mongodb

import "testing"

func TestModelRegister(t *testing.T) {
	RegisterModel((*Model)(nil))

	for k := range modelRegistry {
		t.Log(k)
	}

	model := MakeInstance("github.com/lonelypale/goutils/database/mongodb.Model")
	t.Log(model)
}
