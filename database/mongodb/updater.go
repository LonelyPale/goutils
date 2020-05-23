package mongodb

const (
	SetKey = "$set"
)

type Updater map[string]interface{}

func NewUpdater() Updater {
	return make(Updater)
}

func (u Updater) Set(values ...interface{}) Updater {
	number := len(values)

	if number == 0 {
		return u
	} else if number == 1 {
		u[SetKey] = values[0]
		return u
	} else if number == 2 {
		// 只支持map
		set, ok := u[SetKey].(Updater)
		if !ok {
			set = NewUpdater()
			u[SetKey] = set
		}
		key := values[0].(string)
		val := values[1]
		set[key] = val
		return u
	} else {
		set, ok := u[SetKey].(Updater)
		if !ok {
			set = NewUpdater()
			u[SetKey] = set
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

func (u Updater) GetSet(key string) interface{} {
	if set, ok := u[SetKey].(Updater); ok {
		return set[key]
	}
	return nil
}
