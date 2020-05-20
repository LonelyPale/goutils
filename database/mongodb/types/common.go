package types

type A []interface{}

type M map[string]interface{}

type E struct {
	Key   string
	Value interface{}
}

type D []E

func (d D) Map() M {
	m := make(M, len(d))
	for _, e := range d {
		m[e.Key] = e.Value
	}
	return m
}
