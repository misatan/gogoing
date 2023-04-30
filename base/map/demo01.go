package mapDemo

import "fmt"

func init() {
	//使用案例1
	test1()
	//使用案例2 [支持在声明时填充元素]
	test2()
}

func test1() {
	scoreMap := make(map[string]int, 8)
	scoreMap["zhangsan"] = 90
	scoreMap["xiaoming"] = 100
	fmt.Println(scoreMap)
	fmt.Println(scoreMap["xiaoming"])
	fmt.Printf("type of a:%T\n", scoreMap)
}

func test2() {
	userInfo := map[string]string{
		"username": "misatan",
		"password": "12356",
	}
	fmt.Println(userInfo)
}
