package validator

import (
	"testing"

	"github.com/LonelyPale/goutils"
	"github.com/LonelyPale/goutils/database/mongodb/types"
)

type testObjectID struct {
	ID types.ObjectID `validate:"isdefault"`
	//ID types.ObjectID `validate:"required"`
}

func TestValidateObjectID(t *testing.T) {
	//goutils.RegisterCustomValidateType(ValidateObjectID, types.ObjectID{})

	objid := testObjectID{}
	//objid := testObjectID{ID: types.NewObjectID()}
	if err := goutils.Validate(objid); err != nil {
		t.Fatal(err)
	}

	t.Log(objid)
}
