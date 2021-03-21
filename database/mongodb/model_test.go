package mongodb

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LonelyPale/goutils"
	"github.com/LonelyPale/goutils/types"
)

type User struct {
	Model `bson:",inline"`
	Num   int
	Test  []string
	Temp  *[]int
	Name  string
	NameP *string
	Age   int   `bson:"age,omitempty"`
	AgeP  *int  `bson:",omitempty"`
	Is    bool  `bson:",omitempty"`
	IsP   *bool `bson:",omitempty"`
	prv   string
	prvp  *string
	Sub   `bson:",inline"`
}

type Sub struct {
	Test1 int    `bson:"test1,omitempty"`
	Test2 string `bson:"test2,omitempty"`
}

func init() {
	RegisterModel((*User)(nil), "test")
}

func TestModel_Create(t *testing.T) {
	coll := client.Database("test").Collection("test")

	user := &User{
		Name: "james",
		Sub: Sub{
			Test1: 123,
			Test2: "abc",
		},
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

	id, err := types.ObjectIDFromHex("6056b38fb0eb31305e24c08a")
	if err != nil {
		t.Fatal(err)
	}

	user := &User{Name: "Jerry"}
	user.Model = NewModel(user, coll)
	user.ID = id

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
	coll := client.Database("test").Collection("test")

	id, err := types.ObjectIDFromHex("6056b38fb0eb31305e24c08a")
	if err != nil {
		t.Fatal(err)
	}

	user := &User{}
	user.Model = NewModel(user, coll)
	user.ID = id

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
