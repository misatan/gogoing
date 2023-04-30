package main

import (
	"fmt"
	"reflect"
)

type P struct {
	Name string `json:"name" bson:"Name"`
	Age  int8   `json:"age" bson:"Age"`
}

func main() {
	p := P{
		"sw",
		18,
	}

	pType := reflect.TypeOf(p)
	fieldName, isOk := pType.FieldByName("Name")
	if isOk {
		jsonTag := fieldName.Tag.Get("json")
		bsonTag := fieldName.Tag.Get("bson")
		fmt.Println("Name json Tag=", jsonTag, ",bson Tag=", bsonTag)
	} else {
		fmt.Println("no Name field")
	}

	fieldName, isOk = pType.FieldByName("Age")
	if isOk {
		jsonTag := fieldName.Tag.Get("json")
		bsonTag := fieldName.Tag.Get("bson")
		fmt.Println("Age json Tag=", jsonTag, ",bson Tag=", bsonTag)
	} else {
		fmt.Println("no Age field")
	}
}
