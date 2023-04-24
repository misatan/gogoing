package sliceDemo

import "fmt"

func main() {
	test1()
	test2()
}

func test1() {
	array := []int{10, 20, 30, 40}
	slice := make([]int, 6)
	n := copy(slice, array)
	fmt.Println(n, slice)
}

func test2() {
	slice := make([]byte, 3)
	n := copy(slice, "abcdef")
	fmt.Println(n, slice)
}

/**
打印结果：
	4 [10 20 30 40 0 0]
	3 [97,98,99]
*/
