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

func TestGetModelType(t *testing.T) {
	RegisterModel((*Model)(nil))

	mt := GetModelType("github.com/lonelypale/goutils/database/mongodb.Model")
	t.Log("mt", mt)
	t.Log("mt.Kind", mt.Kind())
	t.Log("mt.Type", mt.Type())

	ft := mt.Field("id")
	t.Log("ft", ft)
	t.Log("ft.Kind", ft.Kind())
	t.Log("ft.Type", ft.Type())

	typ := ft.Type()
	t.Log(typ.Kind())
	t.Log(typ.String())
	t.Log(typ.PkgPath())
	t.Log(typ.Name())
	t.Log(typ.PkgPath() + "." + typ.Name())
}
