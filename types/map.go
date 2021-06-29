package types

import "strconv"

type Map map[string]interface{}

func (m Map) Map(key string) Map {
	if val, ok := m[key]; ok {
		return val.(map[string]interface{})
	}
	return nil
}

func (m Map) List(key string) []interface{} {
	if val, ok := m[key]; ok {
		return val.([]interface{})
	}
	return nil
}

func (m Map) String(key string) string {
	if val, ok := m[key]; ok {
		return val.(string)
	}
	return ""
}

func (m Map) Bool(key string) bool {
	if val, ok := m[key]; ok {
		return val.(bool)
	}
	return false
}

func (m Map) Int(key string) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float32:
			return int(v)
		case float64:
			return int(v)
		case string:
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			return i
		default:
			return val.(int)
		}
	}
	return 0
}
