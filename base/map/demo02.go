package mapDemo

import "fmt"

func init() {
	scoreMap := make(map[string]int)
	scoreMap["zhangsan"] = 90
	scoreMap["lisi"] = 100
	//如果key存在ok为true，v为对应的值，不存在ok为false，v为值类型的默认值
	v, ok := scoreMap["lisi"]
	if ok {
		fmt.Println(v)
	} else {
		fmt.Println("查无此人")
	}
}
