package mongodb

import (
	"reflect"
	"strings"

	"github.com/LonelyPale/goutils/validator"
)

// 用来判断 type T 是否实现了接口 I, 用作类型断言, 如果 T 没有实现接口 I, 则编译错误.
var _ ModelType = new(modelType)
var _ FieldType = new(fieldType)
var _ ValidateTag = new(validateTag)

var modelRegistry = make(map[string]ModelType)

// mongodb.RegisterModel((*User)(nil), UserCollectionName)
func RegisterModel(typedNil interface{}, names ...string) {
	typed := reflect.TypeOf(typedNil).Elem()
	if len(names) > 0 {
		modelRegistry[names[0]] = newModelType(typed)
	} else {
		modelRegistry[typed.PkgPath()+"."+typed.Name()] = newModelType(typed)
	}
}

func MakeInstance(name string) interface{} {
	if val, ok := modelRegistry[name]; ok {
		return reflect.New(val.Type()).Elem().Interface()
	}
	return nil
}

func GetModelType(name string) ModelType {
	if val, ok := modelRegistry[name]; ok {
		return val
	}
	return nil
}

type ModelType interface {
	Field(jsonTagName string) FieldType
	Kind() reflect.Kind
	Type() reflect.Type
}

type modelType struct {
	fields map[string]*fieldType
	kind   reflect.Kind
	typ    reflect.Type
}

func newModelType(mtype reflect.Type) *modelType {
	if mtype.Kind() != reflect.Struct {
		panic("expect struct")
	}

	mt := &modelType{
		kind:   mtype.Kind(),
		typ:    mtype,
		fields: make(map[string]*fieldType),
	}

	num := mtype.NumField()
	for i := 0; i < num; i++ {
		field := mtype.Field(i)
		ft := newFieldType(field)
		if len(ft.json) > 0 && ft.json != "-" {
			mt.fields[ft.json] = ft
		} else {
			mt.fields[ft.field] = ft
		}
	}

	return mt
}

func (mt *modelType) Field(jsonTagName string) FieldType {
	if val, ok := mt.fields[jsonTagName]; ok {
		return val
	}
	return nil
}

func (mt *modelType) Kind() reflect.Kind {
	return mt.kind
}

func (mt *modelType) Type() reflect.Type {
	return mt.typ
}

type FieldType interface {
	Field() string
	Bson() string
	Json() string
	ValidateTag() ValidateTag
	Kind() reflect.Kind
	Type() reflect.Type
}

type fieldType struct {
	field string
	bson  string
	json  string
	valid *validateTag
	kind  reflect.Kind
	typ   reflect.Type
}

func newFieldType(field reflect.StructField) *fieldType {
	bsonTag := field.Tag.Get("bson")
	jsonTag := field.Tag.Get("json")

	return &fieldType{
		field: field.Name,
		bson:  strings.Split(bsonTag, ",")[0],
		json:  strings.Split(jsonTag, ",")[0],
		valid: newValidateTag(field.Tag),
		kind:  field.Type.Kind(),
		typ:   field.Type,
	}
}

func (ft *fieldType) Field() string {
	return ft.field
}

func (ft *fieldType) Bson() string {
	return ft.bson
}

func (ft *fieldType) Json() string {
	return ft.json
}

func (ft *fieldType) ValidateTag() ValidateTag {
	return ft.valid
}

func (ft *fieldType) Kind() reflect.Kind {
	return ft.kind
}

func (ft *fieldType) Type() reflect.Type {
	return ft.typ
}

type ValidateTag interface {
	Validate() string
	Create() string
	Update() string
	Delete() string
	Find() string
	Label() string
}

type validateTag struct {
	validate string
	create   string
	update   string
	delete   string
	find     string
	label    string
}

func newValidateTag(tag reflect.StructTag) *validateTag {
	return &validateTag{
		validate: tag.Get(validator.DefaultTagName),
		create:   tag.Get(CreateTagName),
		update:   tag.Get(UpdateTagName),
		delete:   tag.Get(DeleteTagName),
		find:     tag.Get(FindTagName),
		label:    tag.Get(validator.DefaultLabelTagName),
	}
}

func (vt *validateTag) Validate() string {
	return vt.validate
}

func (vt *validateTag) Create() string {
	return vt.create
}

func (vt *validateTag) Update() string {
	return vt.update
}

func (vt *validateTag) Delete() string {
	return vt.delete
}

func (vt *validateTag) Find() string {
	return vt.find
}

func (vt *validateTag) Label() string {
	return vt.label
}

func Var(field interface{}, name string, mname string, tags ...string) error {
	mt := GetModelType(mname)
	ft := mt.Field(name)
	if ft == nil {
		return ErrNilResult
	}

	valTag := ft.ValidateTag().Validate()
	valTag = strings.ReplaceAll(valTag, "omitempty", "")

	err := validator.Var(field, valTag)

	return err
}
