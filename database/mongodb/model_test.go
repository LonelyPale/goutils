package mongodb

import (
	"encoding/json"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/lonelypale/goutils"
	"github.com/lonelypale/goutils/types"
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
	Test1    int        `bson:"test1,omitempty"`
	Test2    string     `bson:"test2,omitempty"`
	TestTime types.Time `bson:"test_time,omitempty"`
	//TestTime types.Time `bson:"-"`
}

func TestModel_Create(t *testing.T) {
	RegisterModel((*User)(nil), "test")
	coll := client.Database("test").Collection("test")

	user := &User{
		Name: "james",
		Sub: Sub{
			Test1:    123,
			Test2:    "abc",
			TestTime: types.Now(),
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
	RegisterModel((*User)(nil), "test")
	coll := client.Database("test").Collection("test")

	user := &User{
		Name: "tom",
	}
	user.Model = NewModel(user, coll)

	transaction := NewTransaction(coll.client)
	err := transaction.Run(nil, func(sctx mongo.SessionContext) error {
		if err := user.Create(sctx); err != nil {
			return err
		}
		t.Log(user.ID)

		//id, err := coll.InsertOne(sctx, user)
		//if err != nil {
		//	return err
		//}
		//t.Log(id)

		n, err := coll.Delete(sctx, &map[string]string{})
		if err != nil {
			return err
		}
		t.Logf("clear record: %v", n)

		//startTime := time.Now()
		//if err := user.Create(sctx); err != nil {
		//	return err
		//}
		//t.Log("time:", time.Since(startTime))
		//t.Log(user.ID)
		//t.Log(user)

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestModel_Update(t *testing.T) {
	RegisterModel((*User)(nil), "test")
	coll := client.Database("test").Collection("test")

	id, err := types.ObjectIDFromHex("6058376fd7889c8703058df4")
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
	//var tTime = reflect.TypeOf(time.Time{})
	//builder := bson.NewRegistryBuilder()
	//builder.RegisterCodec(tTime, bsoncodec.NewTimeCodec(&bsonoptions.TimeCodecOptions{UseLocalTimeZone: types.Bool(true)}))
	//bson.DefaultRegistry = builder.Build()
	//b, err := bson.DefaultRegistry.LookupDecoder(tTime)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(b)

	RegisterModel((*User)(nil), "test")
	coll := client.Database("test").Collection("test")

	id, err := types.ObjectIDFromHex("6058481014eee1d530389413")
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

func TestModel_Json(t *testing.T) {
	str := `{"id":"6058481014eee1d530389413","createTime":"2021-02-12 13:22:52","updateTime":"2021-02-17 21:47:11", "Name":"abc123"}`
	var user User
	if err := json.Unmarshal([]byte(str), &user); err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func TestTemp(t *testing.T) {
	user := &User{}
	if err := goutils.Inject(user, "NameP", types.String("testing...")); err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}
