package main

import "fmt"

func main() {
	type person1 struct {
		name, city string
		age        int8
	}
	//3.匿名结构体
	var user struct {
		Name string
		Age  int
	}
	user.Name = "frc"
	user.Age = 22
	fmt.Printf("user=%#v\n", user)

	//4.创建指针类型结构体
	var p2 = new(person1)
	//go语言 指针类型结构体 支持 使用 .字段名的方式 直接访问字段 【go语言的语法糖;实际底层是 *p.property】
	p2.name = "somebody"
	p2.age = 100
	p2.city = "somewhere"
	fmt.Printf("p2 type of %T\n", p2)
	fmt.Printf("p2 = %#v\n", p2)
}
