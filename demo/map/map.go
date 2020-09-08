package main

import "fmt"

type Map map[string]interface{}

func main() {
	m := Map{"name": "bob", "age": 10}
	sex, ok := m["sex"]
	fmt.Println(sex, ok)

	if sex, ok := m["sex"]; !ok || len(sex.(string)) == 0 {
		fmt.Println("Sex does not exist")
	}

	switch val := m["sex"].(type) {
	case nil:
		fmt.Println("sex is nil")
	case string:
		if len(val) == 0 {
			fmt.Println("sex len is 0")
		}
	default:
		fmt.Println("Unknown type")
	}
}
