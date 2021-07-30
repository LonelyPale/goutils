package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Name    string `form:"name" json:"name" validate:"required" label:"姓名"`
	Age     uint8  `form:"age" json:"age" validate:"required,gt=18" label:"年龄"`
	Passwd  string `form:"passwd" json:"passwd" validate:"required,max=20,min=6" label:"密码"`
	Code    string `form:"code" json:"code" validate:"required,len=6" label:"代码"`
	Address string `form:"address" json:"address" validate:"omitempty" vCreate:"required" label:"地址"`
}

func TestValidate(t *testing.T) {
	users := &User{
		Name:   "admin",
		Age:    12,
		Passwd: "123",
		Code:   "123456",
	}

	validate := validator.New()

	trans, err := SetLanguageZH(validate)
	if err != nil {
		t.Fatal(err)
	}

	err = validate.Struct(users)
	t.Log(err)
	t.Log(Translate(err, trans))
}

func TestValidate2(t *testing.T) {
	users := &User{
		Name:   "admin",
		Age:    12,
		Passwd: "123",
		Code:   "123456",
	}

	validate := NewDefaultValidator()
	err := validate.SetLanguage(ZH)
	if err != nil {
		t.Fatal(err)
	}

	err = validate.Validate(users)
	t.Log(err)
}

func TestValidateTag(t *testing.T) {
	validate := NewDefaultValidator()
	err := validate.SetLanguage(ZH)
	if err != nil {
		t.Fatal(err)
	}

	user := User{
		Name:   "abc",
		Passwd: "123456",
		Age:    123,
		Code:   "qwe123",
	}

	err = validate.Validate(user)
	t.Log(err)

	err = validate.Validate(user, "vCreate")
	t.Log(err)
}
