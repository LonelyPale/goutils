// Created by LonelyPale at 2020-04-12

package mongodb

import (
	"github.com/LonelyPale/goutils/database/mongodb/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func NewFilter() Filter {
	return make(Filter)
}

func In(value interface{}) Filter {
	return Filter{InKey: value}
}

func Lt(value interface{}) Filter {
	return Filter{LtKey: value}
}

// { name:"go" }
func (f Filter) Set(key string, value interface{}) Filter {
	f[key] = value
	return f
}

func (f Filter) Get(key string) interface{} {
	return f[key]
}

func (f Filter) Delete(key string) {
	delete(f, key)
}

// { $or:[ {name:"c"}, {name:"go"} ] }
func (f Filter) Or(key string, value interface{}) Filter {
	or, ok := f[OrKey]
	if ok {
		orArray, ok := or.(types.A)
		if ok {
			orArray = append(orArray, types.M{key: value})
		} else {
			panic("type error: Not valid types.A")
		}
	} else {
		f[OrKey] = types.A{types.M{key: value}}
	}
	return f
}

// { $and:[ {name:{$ne:"c"}}, {name:{$ne:"go"}} ] }
func (f Filter) And(key string, value interface{}) Filter {
	f[key] = Filter{AndKey: value}
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

//todo
func (f Filter) Regex(key string, value primitive.Regex) Filter {
	f[key] = value
	return f
}

func (f Filter) ID(value interface{}) Filter {
	f[IDKey] = value
	return f
}
