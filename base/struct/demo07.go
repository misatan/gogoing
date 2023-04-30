package teststruct

import "fmt"

func main() {
	p := newPerson("sw", "wuhu", 18)
	fmt.Printf("%#v\n", p)
}

/*
*
结构体没有构造函数，自己定义； 因为结构体是值类型，如果比较复杂，值copy开销很大，所以使用指针类型
*/
func newPerson(name, city string, age int8) *person {
	return &person{
		name: name,
		city: city,
		age:  age,
	}
}
