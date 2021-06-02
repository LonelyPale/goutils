package mongodb

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LonelyPale/goutils/errors"
	"github.com/LonelyPale/goutils/types"
)

type Sequence struct {
	Model `bson:"-" json:"-"`
	Key   string `bson:"key,omitempty" json:"key,omitempty" validate:"omitempty,min=1" vCreate:"required"`
	Value int64  `bson:"value,omitempty" json:"value,omitempty" validate:"omitempty,gte=1" vCreate:"required"`
	mutex sync.Mutex
}

func NewSequence(coll *Collection) *Sequence {
	sequence := new(Sequence)
	sequence.Model = NewModel(sequence, coll)
	return sequence
}

func (s *Sequence) IncByID(id interface{}) (int64, error) {
	objid, err := types.ObjectIDFrom(id)
	if err != nil {
		return 0, err
	}

	s.ID = objid
	return s.Inc()
}

func (s *Sequence) IncByKey(key string) (int64, error) {
	if len(key) == 0 {
		return 0, errors.New("key is empty")
	}

	s.Key = key
	return s.Inc()
}

func (s *Sequence) Inc() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var filter Filter
	if !s.ID.IsZero() {
		filter = ID(s.ID)
	} else if len(s.Key) > 0 {
		filter = NewFilter("key", s.Key)
	} else {
		return 0, errors.New("filter is nil")
	}

	update := Inc("value", 1)
	opts := &options.FindOneAndUpdateOptions{
		Upsert: types.True(),
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
