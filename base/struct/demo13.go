package main

import (
	"encoding/json"
	"fmt"
)

/*
*
结构体字段可见性:

	结构体字段【大写】开头表示公开可访问，小写表示私有（仅在定义当前结构体的包可以访问）

结构体序列化与反序列化
*/
func main() {
	c := &Class{
		Title:    "101",
		Students: make([]*Student, 0, 200),
	}

	for i := 0; i < 10; i++ {
		stu := &Student{
			Name:   fmt.Sprintf("stu%02d", i),
			Gender: "男",
			ID:     i,
		}
		c.Students = append(c.Students, stu)
	}
	//JSON序列化：结构体 ——> JSON字符串
	data, err := json.Marshal(c)
	if err != nil {
		fmt.Println("json marshal failed")
		return
	}
	fmt.Printf("json:%s\n", data)
	//JSON反序列化：JSON字符串 ——> 结构体
	str := `{"Title":"101","Students":[{"ID":0,"Gender":"男","Name":"stu00"},{"ID":1,"Gender":"男","Name":"stu01"},{"ID":2,"Gender":"男","Name":"stu02"},{"ID":3,"Gender":"男","Name":"stu03"},{"ID":4,"Gender":"男","Name":"stu04"},{"ID":5,"Gender":"男","Name":"stu05"},{"ID":6,"Gender":"男","Name":"stu06"},{"ID":7,"Gender":"男","Name":"stu07"},{"ID":8,"Gender":"男","Name":"stu08"},{"ID":9,"Gender":"男","Name":"stu09"}]}`
	c1 := &Class{}
	err0 := json.Unmarshal([]byte(str), c1)
	if err0 != nil {
		fmt.Println("json unmarshal failed")
	}
	fmt.Printf("%#v\n", c1)
}

// Student 学生
type Student struct {
	ID     int
	Gender string
	Name   string
}

// Class 班级
type Class struct {
	Title    string
	Students []*Student
}
