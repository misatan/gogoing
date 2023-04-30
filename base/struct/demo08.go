package main

import "fmt"

/**
函数接受者：类似java中的this
定义如下：
	func(接受者变量,接受者类型) 方法名(参数列表) (返回参数){
		函数体
	}
接受者变量：接受者中的参数变量名在命名时，官方建议使用接受者类型名的第一个小写字母，而不是self，this之类的命名。
接受者类型：接受者类型和参数类型相似，可以是指针类型也可以是非指针类型
方法名、参数列表、返回参数：与函数定义相同
*/
func main() {
	p1 := newPerson("sw", 25)
	p1.Dream()
}

type Person struct {
	name string
	age  int8
}

func newPerson(name string, age int8) *Person {
	return &Person{
		name: name,
		age:  age,
	}
}

func (p Person) Dream() {
	fmt.Printf("%s的梦想是学好go语言", p.name)
}
