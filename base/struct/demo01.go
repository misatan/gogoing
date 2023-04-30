package main

import "fmt"

// 类型定义
type newInt int

// 类型别名
type myInt = int

func main() {
	var a newInt
	var b myInt

	fmt.Printf("type of a:%T\n", a) // tyep of a:main.newInt
	fmt.Printf("type of b:%T\n", b) // type of b:int
}
