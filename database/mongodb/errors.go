package mongodb

import "github.com/LonelyPale/goutils/errors"

var ErrNilObjectID = errors.New("ObjectID is nil")
var ErrNilFilter = errors.New("filter cannot be nil")
var ErrNilResult = errors.New("the result point cannot be nil")
var ErrNilCollection = errors.New("Collection nil")
var ErrResultSlice = errors.New("result slice type conversion failure")
var ErrDocumentExists = errors.New("document already exists")
