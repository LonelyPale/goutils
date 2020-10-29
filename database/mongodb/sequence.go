package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LonelyPale/goutils/pointer"
)

type Sequence struct {
	Model `bson:"-" json:"-"`
	ID    string `bson:"_id,omitempty" json:"id,omitempty" validate:"omitempty,min=1" vCreate:"required" vModify:"required" vDelete:"required"`
	Value int64  `bson:"value,omitempty" json:"value,omitempty" validate:"omitempty,gte=1"`
}

func NewSequence(coll *Collection) *Sequence {
	sequence := new(Sequence)
	sequence.Model = NewModel(sequence, coll)
	return sequence
}

func (s *Sequence) Inc(name string) (int64, error) {
	filter := ID(name)
	update := Inc("value", 1)
	opts := &options.FindOneAndUpdateOptions{
		Upsert: pointer.Bool(true),
	}

	if err := s.coll.FindOneAndUpdate(nil, s, filter, update, opts); err != nil {
		if err == mongo.ErrNoDocuments {
			if err = s.coll.FindOneAndUpdate(nil, s, filter, update, opts); err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return s.Value, nil
}
