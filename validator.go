package goutils

import (
	"github.com/lonelypale/goutils/validator"
)

//validation

var (
	Validate                   = validator.Validate
	RegisterCustomValidateType = validator.RegisterCustomValidateType
)
