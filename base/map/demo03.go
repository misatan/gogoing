package main

import "fmt"

func main() {
	fmt.Println("---test1---")
	test1()
	fmt.Println("---test2---")
	test2()
	fmt.Println("---test3---")
	test3()
}

func test1() {
	scoreMap := make(map[string]int)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	scoreMap["李四"] = 60
	for k, v := range scoreMap {
		fmt.Println(k, v)
	}
}

//只想遍历key
func test2() {
	scoreMap := make(map[string]int)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	scoreMap["李四"] = 60
	for k := range scoreMap {
		fmt.Println(k)
	}
}

//只想遍历value
func test3() {
	scoreMap := make(map[string]int)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	scoreMap["李四"] = 60
	for _, v := range scoreMap {
		fmt.Println(v)
	}
}
