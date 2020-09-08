package validator

import (
	"testing"

	"github.com/LonelyPale/goutils"
	"github.com/LonelyPale/goutils/types"
)

type testObjectID struct {
	ID types.ObjectID `validate:"isdefault"`
	//ID types.ObjectID `validate:"required"`
	Number int `validate:"required,numeric,gt=0"`
}

func TestValidateObjectID(t *testing.T) {
	//goutils.RegisterCustomValidateType(ValidateObjectID, types.ObjectID{})

	objid := testObjectID{Number: -1}
	//objid := testObjectID{ID: types.NewObjectID()}
	if err := goutils.Validate(objid); err != nil {
		t.Fatal(err)
	}

	t.Log(objid)
}
