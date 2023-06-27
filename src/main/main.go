package main

import "strconv"

type A struct {
	name string
	age  int
}

func (a A) String() string {
	return "Name: " + a.name + ", Age: " + strconv.Itoa(a.age)
}

func main() {
	obj := A{name: "Alice", age: 25}
	print(obj)
}
