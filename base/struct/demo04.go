package teststruct

import "fmt"

func main() {
	type person1 struct {
		name, city string
		age        int8
	}

	//5.结构体默认初始化
	var p3 person1
	fmt.Printf("p3=%#v\n", p3)

	//6.键值初始化
	p4 := person1{
		name: "frc",
		city: "wuhan",
		age:  17,
	}
	fmt.Printf("p4=%#v\n", p4)

	//指针形式
	p5 := &person1{
		name: "sw",
		city: "anhui",
		age:  18,
	}
	fmt.Printf("p5=%#v\n", p5)

	//值列表形式
	/**
	这种初始化方式注意：
		1.必须初始化所有字段
		2.填充顺序必须和结构体定义顺序一致
		3.不能和键值初始化方式混用
	*/
	p6 := person1{
		"sw",
		"china",
		14,
	}
	fmt.Printf("p6=%#v\n", p6)
}
