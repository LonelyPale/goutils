package types

//Array
type A []interface{}

//Map
type M map[string]interface{}

func (m M) String(key string) string {
	if m == nil {
		return ""
	}
	return m[key].(string)
}

//Elem
type E struct {
	Key   string
	Value interface{}
}

//Doc
type D []E

func (d D) Append(key string, val interface{}) D {
	return append(d, E{Key: key, Value: val})
}

func (d D) Map() M {
	m := make(M, len(d))
	for _, e := range d {
		m[e.Key] = e.Value
	}
	return m
}
