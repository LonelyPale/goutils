package mongodb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LonelyPale/goutils/database/mongodb/types"
)

type User struct {
	Model      `bson:"-"`
	ID         types.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreateTime time.Time      `bson:"createTime,omitempty" json:"createTime,omitempty"`
	ModifyTime time.Time      `bson:"modifyTime,omitempty" json:"modifyTime,omitempty"`
	Num        int
	Test       []string
	Temp       *[]int
	Name       string
	NameP      *string
	Age        int   `bson:",omitempty"`
	AgeP       *int  `bson:",omitempty"`
	Is         bool  `bson:",omitempty"`
	IsP        *bool `bson:",omitempty"`
}

func (User) BeforeUpdate(ctx context.Context, coll *Collection) error {
	fmt.Println("BeforeUpdate")
	return nil
}

func (User) AfterUpdate(ctx context.Context, coll *Collection) error {
	fmt.Println("AfterUpdate")
	return nil
}

func init() {
	RegisterType((*User)(nil), "test")
}

func TestModel_Put(t *testing.T) {
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

func TestModel_Put2(t *testing.T) {
	coll := client.Database("test").Collection("test")

	user := &User{
		Name: "tom",
	}
	user.Model = NewModel(user, coll)

	transaction := NewTransaction(coll.client)
	err := transaction.Run(nil, func(sctx mongo.SessionContext) error {
		n, err := coll.Delete(sctx, &map[string]string{})
		if err != nil {
			return err
		}
		t.Logf("clear record: %v", n)

		startTime := time.Now()
		if err := user.Put(sctx); err != nil {
			return err
		}
		t.Log("time:", time.Since(startTime))
		t.Log(user.ID)
		t.Log(user)

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestModel_Set(t *testing.T) {
	coll := client.Database("test").Collection("test")

	id, err := types.ObjectIDFromHex("5ec815d1b732edbfd90159f6")
	if err != nil {
		t.Fatal(err)
	}

	user := &User{ID: id}
	user.Model = NewModel(user, coll)

	startTime := time.Now()
	if err := user.Get(); err != nil {
		t.Fatal(err)
	}
	t.Log("time-get:", time.Since(startTime), user)

	updater := types.M{"$set": types.M{"modifyTime": time.Now()}}
	startTime1 := time.Now()
	if err := user.Set(updater); err != nil {
		t.Fatal(err)
	}
	t.Log("time-set:", time.Since(startTime1), user)
}

func TestModel_Get(t *testing.T) {
	opts := NewClientOptions("")
	client, err := Connect(opts)
	if err != nil {
		t.Fatal(err)
	}
	coll := client.Database("test").Collection("test")

	id, err := types.ObjectIDFromHex("5ec5d7637ebc566c3b8f3b46")
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
