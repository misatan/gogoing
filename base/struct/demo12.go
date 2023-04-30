package main

import "fmt"

func main() {
	d1 := Dog{
		Feet: 4,
		Animal: &Animal{
			name: "炫神",
		},
	}
	d1.move()
	d1.wang()
}

type Animal struct {
	name string
}

func (a *Animal) move() {
	fmt.Printf("%s会移动\n", a.name)
}

// 结构体的继承
type Dog struct {
	Feet    int8
	*Animal //通过嵌套匿名结构体实现继承
}

func (d *Dog) wang() {
	fmt.Printf("%s会汪汪汪\n", d.name)
}

// 可以对父类定义方法进行覆盖
// func (d *Dog) move() {
// 	fmt.Printf("%s不止会移动\n", d.name)
// }
