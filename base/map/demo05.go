package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) //初始化随机数种子

	var scoreMap = make(map[string]int, 200)

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("stu%02d", i) //生成stu开头的字符串
		value := rand.Intn(100)          //生成0-99的随机数
		scoreMap[key] = value
	}

	//去除map中的所有key存入切片keys
	var keys = make([]string, 0, 200)
	for k := range scoreMap {
		keys = append(keys, k)
	}
	//对切片进行排序
	sort.Strings(keys)
	//按照排序后的key遍历map（for range 遍历切片 index,value）
	for _, key := range keys {
		fmt.Println(key, scoreMap[key])
	}
}
