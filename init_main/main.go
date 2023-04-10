package main

import "fmt"

var cons1 int

var (
	cons2 float32
	cons3 string
)

func init() {
	cons1 = 100
	cons2 = 3.14
	cons3 = "test_str"
}

func main() {
	fmt.Println("cons1 = ", cons1, ",cons2 = ", cons2, "cons3 = ", cons3)
}
