package mongodb

import "github.com/lonelypale/goutils"

func (f Filter) TimeCondition(keys ...string) Filter {
	for _, key := range keys {
		val, ok := f[key]
		if !ok {
			return f
		}

		switch v := val.(type) {
		case string:
			if len(v) == 0 {
				delete(f, key)
				return f
			}
		case []interface{}:
			if len(v) < 2 {
				delete(f, key)
				return f
			}
			start, err := goutils.TimeParseToUTC(v[0].(string))
			if err != nil {
				panic(err)
			}
			end, err := goutils.TimeParseToUTC(v[1].(string))
			if err != nil {
				panic(err)
			}
			delete(f, key)
			f.Gte(key, start).Lte(key, end)
		}
	}
	return f
}

func (f Filter) RegexCondition(keys ...string) Filter {
	for _, key := range keys {
		val, ok := f[key]
		if !ok {
			continue
		}

		if len(val.(string)) > 0 {
			f.Regex(key, val)
		} else {
			delete(f, key)
		}
	}
	return f
}
