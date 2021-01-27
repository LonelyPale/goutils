package mongodb

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/LonelyPale/goutils/pagination"
	"github.com/LonelyPale/goutils/types"
)

func TestCollection_FindOrInsert(t *testing.T) {
	collection := client.Database("test").Collection("test")
	user := &User{Num: 1, Name: "mongo-0"}
	filter := Filter{"num": 1}
	id, err := collection.FindOrInsert(nil, filter, user)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}

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

func TestCollection_FindPagination(t *testing.T) {
	collection := client.Database("TestDB").Collection("test")
	var result []struct {
		ID    types.ObjectID `bson:"_id,omitempty"`
		Title string         `bson:"title,omitempty"`
		Name  string         `bson:"name,omitempty"`
	}
	pager := pagination.Pagination{Current: 2, PageSize: 2, Data: &result}
	filter := map[string]string{}

	err := collection.Find(nil, &pager, filter)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pager)
}
