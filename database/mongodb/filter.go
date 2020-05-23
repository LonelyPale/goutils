// Created by LonelyPale at 2020-04-12

package mongodb

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LonelyPale/goutils/database/mongodb/types"
)

const (
	OrKey         = "$or"
	AndKey        = "$and"
	InKey         = "$in"
	NinKey        = "$nin"
	AllKey        = "$all"
	NotKey        = "$not"
	NeKey         = "$ne"  // !=
	GtKey         = "$gt"  // >
	GteKey        = "$gte" // >=
	LtKey         = "$lt"  // <
	LteKey        = "$lte" // <=
	ModKey        = "$mod" // mod 取模运算
	ExistsKey     = "$exists"
	TypeKey       = "$type"
	SizeKey       = "$size"
	MatchKey      = "$elemMatch"
	IDKey         = "_id" // ObjectID Field Name
	CreateTimeKey = "createTime"
	ModifyTimeKey = "modifyTime"
)

const (
	IDField         = "ID"
	CreateTimeField = "CreateTime"
	ModifyTimeField = "ModifyTime"
)

type Filter map[string]interface{}

func NewFilter(values ...interface{}) Filter {
	return make(Filter).Set(values...)
}

func Or(values ...interface{}) Filter {
	return Filter{}.Or(values...)
}

func And(values ...interface{}) Filter {
	return Filter{}.And(values...)
}

func ID(value interface{}) Filter {
	return Filter{}.ID(value)
}

// { name: "go" }
func (f Filter) Set(values ...interface{}) Filter {
	number := len(values)

	if number == 0 {
		return f
	} else if number == 1 {
		return f
	} else {
		for i, n := 0, number/2; i < n; i++ {
			key := values[i*2].(string)
			val := values[i*2+1]
			f[key] = val
		}
	}

	return f
}

func (f Filter) Get(key string) interface{} {
	return f[key]
}

func (f Filter) Delete(key string) {
	delete(f, key)
}

// { $or: [ {name: "c"}, {name: "go"} ] }
// { $or: [ { quantity: { $lt: 20 } }, { price: 10 } ] }
// { $or: [ { <expression1> }, { <expression2> }, ... , { <expressionN> } ] }
func (f Filter) Or(values ...interface{}) Filter {
	number := len(values)

	if number == 0 {
		return f
	} else if number == 1 {
		f[OrKey] = values[0]
	} else {
		orArr, ok := f[OrKey]
		if !ok {
			orArr = types.A{}
		}

		or, ok := orArr.(types.A)
		if !ok {
			log.Warn("$or object not valid types.A")
			return f
		}

		for i, n := 0, number/2; i < n; i++ {
			key := values[i*2].(string)
			val := values[i*2+1]
			or = append(or, Filter{key: val})
		}

		f[OrKey] = or
	}

	return f
}

// { $and: [ { name: { $ne: "c"} }, { name: { $ne: "go"} } ] }
// { $and: [ { price: { $ne: 1.99 } }, { price: { $exists: true } } ] }
// { $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
func (f Filter) And(values ...interface{}) Filter {
	number := len(values)

	if number == 0 {
		return f
	} else if number == 1 {
		f[AndKey] = values[0]
	} else {
		andArr, ok := f[AndKey]
		if !ok {
			andArr = types.A{}
		}

		and, ok := andArr.(types.A)
		if !ok {
			log.Warn("$and object not valid types.A")
			return f
		}

		for i, n := 0, number/2; i < n; i++ {
			key := values[i*2].(string)
			val := values[i*2+1]
			and = append(and, Filter{key: val})
		}

		f[AndKey] = and
	}

	return f
}

func (f Filter) In(key string, value interface{}) Filter {
	f[key] = Filter{InKey: value}
	return f
}

func (f Filter) All(key string, value interface{}) Filter {
	f[key] = Filter{AllKey: value}
	return f
}

func (f Filter) Ne(key string, value interface{}) Filter {
	f[key] = Filter{NeKey: value}
	return f
}

func (f Filter) Gt(key string, value interface{}) Filter {
	f[key] = Filter{GtKey: value}
	return f
}

func (f Filter) Gte(key string, value interface{}) Filter {
	f[key] = Filter{GteKey: value}
	return f
}

func (f Filter) Lt(key string, value interface{}) Filter {
	f[key] = Filter{LtKey: value}
	return f
}

func (f Filter) Lte(key string, value interface{}) Filter {
	f[key] = Filter{LteKey: value}
	return f
}

func (f Filter) Exists(key string, value bool) Filter {
	f[key] = Filter{ExistsKey: value}
	return f
}

func (f Filter) Type(key string, value interface{}) Filter {
	f[key] = Filter{TypeKey: value}
	return f
}

func (f Filter) Size(key string, value interface{}) Filter {
	f[key] = Filter{SizeKey: value}
	return f
}

// todo: 待整理
func (f Filter) Regex(key string, value primitive.Regex) Filter {
	f[key] = value
	return f
}

func (f Filter) ID(value interface{}) Filter {
	f[IDKey] = value
	return f
}
