package main

import "fmt"

/*
*

	if 布尔表达式 {
		执行体
	}else if 布尔表达式{

		执行体
	}else{

		执行体
	}
*/
func main() {
	ifTest(100)
}

func ifTest(num int) {
	if num < 0 {
		fmt.Println("num < 0,num:", num)
	} else if num < 100 {
		fmt.Println("num < 100,num:", num)
	} else {
		fmt.Println("num >= 100,num:", num)
	}
}
