package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LonelyPale/goutils/errors"
)

var ErrNoDocuments = mongo.ErrNoDocuments

var ErrNilObjectID = errors.New("ObjectID is nil")
var ErrNilFilter = errors.New("filter cannot be nil")
var ErrNilResult = errors.New("the result point cannot be nil")
var ErrNilCollection = errors.New("collection is nil")
var ErrNilDocument = errors.New("document is nil")
var ErrDocumentExists = errors.New("document already exists")
var ErrDocumentNotExist = errors.New("document does not exist")
