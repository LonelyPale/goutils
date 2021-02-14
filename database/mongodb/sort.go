package mongodb

type Sort map[string]int

// 默认升序
func NewSort(keys ...string) Sort {
	return make(Sort).Asc(keys...)
}

func Asc(keys ...string) Sort {
	return Sort{}.Asc(keys...)
}

func Desc(keys ...string) Sort {
	return Sort{}.Desc(keys...)
}

// 升序
func (s Sort) Asc(keys ...string) Sort {
	for _, key := range keys {
		if len(key) == 0 {
			continue
		}
		s[key] = 1
	}
	return s
}

// 降序
func (s Sort) Desc(keys ...string) Sort {
	for _, key := range keys {
		if len(key) == 0 {
			continue
		}
		s[key] = -1
	}
	return s
}
