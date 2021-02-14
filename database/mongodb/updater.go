package mongodb

import (
	log "github.com/sirupsen/logrus"

	"github.com/LonelyPale/goutils/errors"
)

// Field Update Operators
const (
	SetKey   = "$set"
	UnsetKey = "$unset"
	IncKey   = "$inc"
)

type Updater map[string]interface{}

func NewUpdater(values ...interface{}) Updater {
	return make(Updater).Set(values...)
}

func Set(values ...interface{}) Updater {
	return Updater{}.Set(values...)
}

func Unset(keys ...string) Updater {
	return Updater{}.Unset(keys...)
}

func Inc(values ...interface{}) Updater {
	return Updater{}.Inc(values...)
}

// 只支持map
func (u Updater) checkout(key string) (Updater, error) {
	val, ok := u[key]
	if !ok {
		val = Updater{}
		u[key] = val
	}

	updater, ok := val.(Updater)
	if !ok {
		err := errors.Errorf("%v object not Updater", key)
		log.Warn(err)
		return updater, err
	}

	return updater, nil
}

// { $set: { "details.make": "zzz" } }
// { $set: { "tags.1": "rain gear", "ratings.0.rating": 2 } }
// { $set: { <field1>: <value1>, ... } }
func (u Updater) Set(values ...interface{}) Updater {
	number := len(values)

	if number == 0 {
		return u
	} else if number == 1 {
		u[SetKey] = values[0]
	} else {
		set, err := u.checkout(SetKey)
		if err != nil {
			return u
		}

		for i, n := 0, number/2; i < n; i++ {
			key := values[i*2].(string)
			val := values[i*2+1]
			set[key] = val
		}
	}

	return u
}

// { $unset: { quantity: "", instock: "" } }
// { $unset: { <field1>: "", ... } }
func (u Updater) Unset(keys ...string) Updater {
	unset, err := u.checkout(UnsetKey)
	if err != nil {
		return u
	}

	for _, key := range keys {
		unset[key] = ""
	}

	return u
}

// { $inc: { quantity: -2, "metrics.orders": 1 } }
// { $inc: { <field1>: <amount1>, <field2>: <amount2>, ... } }
func (u Updater) Inc(values ...interface{}) Updater {
	number := len(values)

	if number == 0 {
		return u
	} else if number == 1 {
		u[IncKey] = values[0]
	} else {
		set, err := u.checkout(IncKey)
		if err != nil {
			return u
		}

		for i, n := 0, number/2; i < n; i++ {
			key := values[i*2].(string)
			val := values[i*2+1]
			set[key] = val
		}
	}

	return u
}

func (u Updater) Get(key string) interface{} {
	return u[key]
}

func (u Updater) Delete(key string) {
	delete(u, key)
}

func (u Updater) UpdateTime(value interface{}) Updater {
	set, err := u.checkout(SetKey)
	if err != nil {
		return u
	}

	set[UpdateTimeBson] = value
	return u
}
