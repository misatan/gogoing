package pointer

import "fmt"

func main() {
	var num int
	fmt.Printf("num addr: %p\n", &num)
	ptr := &num
	fmt.Printf("num:%d\n", num)
	*ptr = 10
	fmt.Printf("num:%d\n", num)
}
