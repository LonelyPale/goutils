package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func main() {
	vmess()
	vless()
	test()
}

func vmess() {
	filename := "/Users/wyb/Downloads/V2RayN_1641874422.txt"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	//res, err := base64.URLEncoding.DecodeString(string(data))
	res, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
}

func vless() {
	filename := "/Users/wyb/Downloads/VLESS_1641875121.txt"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	res, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
}

func test() {
	str := "eyJ2IjoiMiIsInBzIjoi8J+MjyB2MXzkuprmtLJ8Q3zmt7flkIjotJ/ovb186KeG6aKR5Li75YqbIiwiaWQiOiIwYTc4MWZhYS0wOTRlLWFlNjktZjg1NC0zMzVlYzlkZWJhMmEiLCJhZGQiOiJrMy43aDcueHl6IiwicG9ydCI6IjQxMzIxIiwiYWlkIjoiMCIsIm5ldCI6InRjcCIsInR5cGUiOiJodHRwIiwiaG9zdCI6IjYxNzY3YTcyNjM3MTYzNzczYTZjNzc3NS5zaW5hLmNuIiwicGF0aCI6Ii8iLCJ0bHMiOiIifQ=="
	res, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
}
