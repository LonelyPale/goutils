package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name   string
	Data   interface{}
	Number int
}

func main() {
	//test1()
	test2()
}

func test1() {
	user := User{
		Name: "",
		Data: "123",
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))

	u := new(User)
	str := `{"Name":"","Data":123}`
	if err := json.Unmarshal([]byte(str), &u); err != nil {
		fmt.Println(err)
	}
	fmt.Println(u)
}

func test2() {
	str := `{"number":12.3}`
	var obj User

	if err := json.Unmarshal([]byte(str), &obj); err != nil {
		fmt.Println(err)
	}

	fmt.Println(obj)

}
