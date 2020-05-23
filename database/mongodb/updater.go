package mongodb

import log "github.com/sirupsen/logrus"

const (
	SetKey   = "$set"
	UnsetKey = "$unset"
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
		// 只支持map
		setMap, ok := u[SetKey]
		if !ok {
			setMap = NewUpdater()
			u[SetKey] = setMap
		}

		set, ok := setMap.(Updater)
		if !ok {
			log.Warn("$set object not Updater")
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
	unsetMap, ok := u[UnsetKey]
	if !ok {
		unsetMap = NewUpdater()
		u[UnsetKey] = unsetMap
	}

	unset, ok := unsetMap.(Updater)
	if !ok {
		log.Warn("$unset object not Updater")
		return u
	}

	for _, key := range keys {
		unset[key] = ""
	}

	return u
}

func (u Updater) Get(key string) interface{} {
	return u[key]
}

func (u Updater) Delete(key string) {
	delete(u, key)
}
