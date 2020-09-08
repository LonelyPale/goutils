package main

import "fmt"

var (
	slice = []string{"", "", "a", "", "b", "", "c"}
)

func test1() {
	elems := make([]string, len(slice))
	copy(elems, slice)
	for i, a := range elems {
		if len(a) == 0 {
			elems = append(elems[:i], elems[i+1:]...)
		}
	}

	fmt.Println(len(elems), elems)
}

func test2() {
	elems := make([]string, 0)
	for _, s := range slice {
		if len(s) > 0 {
			elems = append(elems, s)
		}
	}

	fmt.Println(len(elems), elems)
}

func main() {
	test1()
	test2()
}
