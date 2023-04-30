package teststruct

import "fmt"

func main() {
	//结构体的内存布局
	/**

	 */
	type test struct {
		a int8
		b int8
		c int8
		d int8
	}

	n := test{
		1, 2, 3, 4,
	}

	fmt.Printf("n.a %p\n", &n.a)
	fmt.Printf("n.b %p\n", &n.b)
	fmt.Printf("n.c %p\n", &n.c)
	fmt.Printf("n.d %p\n", &n.d)
}
