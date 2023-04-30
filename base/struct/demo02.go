package teststruct

import (
	"fmt"
)

/*
*
结构体定义:

	字段名 	字段类型
	字段名	字段类型
	...

注意点:

	1.类型名：标识自定义结构体的名称，在同一个包内不能重复
	2.字段名：
	3.字段类型：
*/
func main() {

	//2.结构体实例化
	// var 实例名 结构体类型
	var p1 person0
	p1.name = "sw"
	p1.city = "wuhu"
	p1.age = 23
	fmt.Printf("p1=%v\n", p1)
	fmt.Printf("p1=%#v\n", p1)

	//通过实例名.字段名 的方式访问实例字段
}
