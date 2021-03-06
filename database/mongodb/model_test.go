package mongodb

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LonelyPale/goutils"
	"github.com/LonelyPale/goutils/types"
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
	Age        int   `bson:"age,omitempty"`
	AgeP       *int  `bson:",omitempty"`
	Is         bool  `bson:",omitempty"`
	IsP        *bool `bson:",omitempty"`
	prv        string
	prvp       *string
}

func init() {
	RegisterType((*User)(nil), "test")
}

func TestModel_Create(t *testing.T) {
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
	if err := user.Create(); err != nil {
		t.Fatal(err)
	}
	t.Log("time:", time.Since(startTime))
	t.Log(user.ID)
	t.Log(user)
}

func TestModel_Create2(t *testing.T) {
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
		if err := user.Create(sctx); err != nil {
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

func TestModel_Update(t *testing.T) {
	coll := client.Database("test").Collection("test")

	id, err := types.ObjectIDFromHex("5ec81672973997c00c8ba6c8")
	if err != nil {
		t.Fatal(err)
	}

	user := &User{ID: id, Name: "Jerry"}
	user.Model = NewModel(user, coll)

	//startTime := time.Now()
	//if err := user.Get(); err != nil {
	//	t.Fatal(err)
	//}
	//t.Log("time-get:", time.Since(startTime), user)

	user.Name = "Jerry-00"
	//updater := types.M{"$set": types.M{"modifyTime": time.Now()}, "$unset": types.M{"test": "", "temp": ""}}
	//updater := types.M{"$set": user}
	//updater := NewUpdater().Set(user)
	//updater := NewUpdater().Set("name", "Jerry-10", "age", 18)
	//updater := NewUpdater("name", "Jerry-10", "age", 18).Unset("test","temp","namep")
	startTime1 := time.Now()
	if err := user.Update(); err != nil {
		t.Fatal(err)
	}
	t.Log("time-set:", time.Since(startTime1), user)
}

func TestModel_Find(t *testing.T) {
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
	if err := user.Find(); err != nil {
		t.Fatal(err)
	}
	t.Log("time:", time.Since(startTime))
	t.Log(user.ID)
	t.Log(user)
}

func TestTemp(t *testing.T) {
	user := &User{}
	if err := goutils.Inject(user, "NameP", types.String("testing...")); err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}
