package mongodb

import (
	"testing"
	"time"

	"github.com/LonelyPale/goutils/database/mongodb/types"
)

type User struct {
	Model `bson:"-"`
	ID    types.ObjectID `bson:"_id,omitempty"`
	Num   int
	Test  []string
	Temp  *[]int
	Name  string
	NameP *string
	Age   int   `bson:",omitempty"`
	AgeP  *int  `bson:",omitempty"`
	Is    bool  `bson:",omitempty"`
	IsP   *bool `bson:",omitempty"`
}

func TestModel_Put(t *testing.T) {
	opts := NewClientOptions("")
	client, err := Connect(opts)
	if err != nil {
		t.Fatal(err)
	}
	coll := client.Database("test").Collection("test")

	user := &User{
		Name: "james",
	}
	user.Model = NewModel(user, coll)

	n, err := coll.Delete(nil, &map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("clear record: %v", n)

	startTime := time.Now()
	if err := user.Put(); err != nil {
		t.Fatal(err)
	}
	t.Log("time:", time.Since(startTime))
	t.Log(user.ID)
	t.Log(user)
}

func TestModel_Get(t *testing.T) {
	opts := NewClientOptions("")
	client, err := Connect(opts)
	if err != nil {
		t.Fatal(err)
	}
	coll := client.Database("test").Collection("test")

	id, err := types.ObjectIDFromHex("5ec3f79bb4a13447d28955e0")
	if err != nil {
		t.Fatal(err)
	}

	user := &User{ID: id}
	user.Model = NewModel(user, coll)

	startTime := time.Now()
	if err := user.Get(); err != nil {
		t.Fatal(err)
	}
	t.Log("time:", time.Since(startTime))
	t.Log(user.ID)
	t.Log(user)
}
