package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LonelyPale/goutils/pointer"
)

const (
	SequenceCollectionName = "sequence"
)

// 注册模型类型
func init() {
	RegisterType((*Sequence)(nil), SequenceCollectionName)
}

type Sequence struct {
	ID    string `bson:"_id,omitempty" json:"id,omitempty" validate:"omitempty,min=1" vCreate:"required" vModify:"required" vDelete:"required"`
	Value int64  `bson:"value,omitempty" json:"value,omitempty" validate:"omitempty,gte=1"`
}

func GetSequence(name string) (int64, error) {
	filter := ID(name)
	update := Inc("value", 1)

	var result *Sequence
	coll := DB().Collection(SequenceCollectionName)
	opts := &options.FindOneAndUpdateOptions{
		Upsert: pointer.Bool(true),
	}

	if err := coll.FindOneAndUpdate(nil, &result, filter, update, opts); err != nil {
		if err == mongo.ErrNoDocuments {
			if err = coll.FindOneAndUpdate(nil, &result, filter, update, opts); err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return result.Value, nil
}
