package mongodb

import "github.com/LonelyPale/goutils/errors"

var ErrNilObjectID = errors.New("ObjectID is nil")
var ErrNilFilter = errors.New("mongodb: filter cannot be nil")
var ErrNilResult = errors.New("mongodb: the result point cannot be nil")
var ErrNilCollection = errors.New("mongodb: Collection nil")
var ErrResultSlice = errors.New("mongodb: result slice type conversion failure")
var ErrDocumentExists = errors.New("mongodb: document already exists")
var ErrMustPointer = errors.New("Must be a pointer")
