package slice

import "fmt"

func main() {
	// demo01test1()
	demo01test2()
}

// 数组函数传参,数组复制验证
func demo01test1() {
	arrayA := [2]int{100, 200}
	var arrayB [2]int

	arrayB = arrayA

	fmt.Printf("arrayA : %p , %v\n", &arrayA, arrayA)
	fmt.Printf("arrayB : %p , %v\n", &arrayB, arrayB)

	testArray(arrayA)
}

func demo01test2() {
	arrayA := []int{100, 200}
	testArrayPoint(&arrayA) //1.传数组指针
	arrayB := arrayA[:]
	testArrayPoint(&arrayB) //2.传切片
	fmt.Printf("arrayA: %p,%v\n", &arrayA, arrayA)
}

func testArrayPoint(x *[]int) {
	fmt.Printf("func Array: %p,%v\n", x, *x)
	(*x)[1] += 100
}

func testArray(x [2]int) {
	fmt.Printf("func Array : %p , %v\n", &x, x)
}
