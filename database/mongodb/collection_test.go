package mongodb

import (
	"github.com/LonelyPale/goutils/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestCollection_Find(t *testing.T) {
	collection := client.Database("TestDB").Collection("test")
	opts := options.Find().SetSkip(2 * 1).SetLimit(2)

	filter := map[string]string{}
	result := make([]bson.M, 0)
	err := collection.Find(nil, &result, filter, opts)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestCollection_Find1(t *testing.T) {
	collection := client.Database("TestDB").Collection("test")
	opts := options.Find().SetSkip(2 * 1).SetLimit(2)

	filter := map[string]string{}
	var result []struct {
		ID    types.ObjectID `bson:"_id,omitempty"`
		Title string         `bson:"title,omitempty"`
		Name  string         `bson:"name,omitempty"`
	}
	err := collection.Find(nil, &result, filter, opts)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}
