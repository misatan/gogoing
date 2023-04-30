package main

import "fmt"

/**
什么时候使用指针类型接受者?
	1. 需要改变接受者的值(外部需要感知变化)
	2. 接受者是拷贝代价比较大的对象
	3. 保证数据一致性
*/
func main() {
	p1 := Person{
		"sw",
		"wuhu",
		18,
	}
	fmt.Printf("p1=%#v\n", p1)
	p1.setAge(19)
	fmt.Printf("p1=%#v\n", p1)
	p1.setAge2(20)
	fmt.Printf("p1=%#v\n", p1)

	var m1 MyInt
	m1.SayHello()
	m1 = 100
	fmt.Printf("%#v %T\n", m1, m1)

}

type Person struct {
	name, city string
	age        int8
}

//引用类型接受者
func (p *Person) setAge(newAge int8) {
	p.age = newAge
}

//值类型接受者[接受者为副本，不改变原有变量]
func (p Person) setAge2(newAge int8) {
	p.age = newAge
}

//接受者类型可以是任何类型;不仅仅是结构体
type MyInt int

func (num MyInt) SayHello() {
	fmt.Println("hello 我是一个int")
}
