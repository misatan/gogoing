package main

import "fmt"

func main() {
	//嵌套结构体
	user := User{
		"sw",
		"male",
		Address{
			"anhui",
			"wuhu",
			"2000",
		},
	}
	fmt.Printf("user:%#v\n", user)
	//嵌套匿名结构体
	var user0 User
	user0.Name = "frc"
	user0.Gender = "female"
	user0.Address.province = "anhui"
	user0.Address.city = "hefei"
	fmt.Printf("user0:%#v\n", user0)
	//嵌套结构体,字段冲突
	var user1 EmailUser
	user1.Name = "ssw"
	user1.Gender = "male"
	user1.Address.CreateTime = "2000"
	user1.Email.CreateTime = "2000"
	fmt.Printf("user1: %#v\n", user1)
}

// Address 地址结构体
type Address struct {
	province   string
	city       string
	CreateTime string
}

// User 用户结构体(嵌套结构体)
type User struct {
	Name    string
	Gender  string
	Address Address
}

// Email 嵌套结构体的字段名冲突
type Email struct {
	Account    string
	CreateTime string
}

type EmailUser struct {
	Name   string
	Gender string
	Address
	Email
}
