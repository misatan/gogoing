package main

import "fmt"

func main() {
	var mapSlice = make([]map[string]string, 3)
	for index, value := range mapSlice {
		fmt.Printf("index:%d,value:%v\n", index, value)
	}
	fmt.Println("before init")
	//对切片中的map元素初始化
	mapSlice[0] = make(map[string]string, 10)
	mapSlice[0]["name"] = "王五"
	mapSlice[0]["password"] = "123456"
	mapSlice[0]["addr"] = "沈阳大街"
	for index, arr := range mapSlice {
		fmt.Printf("index:%d,value:%v\n", index, arr)
	}
}
