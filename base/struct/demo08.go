package main

import "fmt"

/*
*
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
	p1 := newInstance("sw", 25)
	fmt.Printf("p1 age = %v\n", p1.age) //p1 age = 25
	p1.Dream()
	fmt.Printf("p1 age = %v\n", p1.age) //p1 age = 25
	p1.Sleep()
	fmt.Printf("p1 age = %v\n", p1.age) //p1 age = 18
}

func newInstance(name string, age int8) *Person {
	return &Person{
		name: name,
		age:  age,
	}
}

// Dream 值类型接受者	修改不影响原值
func (p Person) Dream() {
	fmt.Printf("%s的梦想是学好go语言\n", p.name)
	p.age = 20
}

// Sleep 指针类型接受者	修改影响原值
func (p *Person) Sleep() {
	fmt.Printf("%s要睡觉了\n", p.name)
	p.age = 18
}
