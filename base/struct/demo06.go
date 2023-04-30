package main

import "fmt"

/**
面试题
*/
func main() {
	type student struct {
		name string
		age  int
	}

	m := make(map[string]*student)

	stus := []student{
		{name: "A", age: 10},
		{name: "B", age: 20},
		{name: "C", age: 30},
	}

	for _, stu := range stus {
		m[stu.name] = &stu
	}

	for k, v := range m {
		fmt.Println(k, "=>", v.name)
	}
}
