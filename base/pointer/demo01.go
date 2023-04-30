package pointer

import "fmt"

func main() {
	a := 10
	b := &a
	fmt.Printf("a:%d,addr:%p\n", a, &a)
	fmt.Printf("b:%p,addr:%p\n", b, &b)
	fmt.Printf("type of b:%T\n", b)
	c := *b
	fmt.Printf("type of c:%T\n", c)
	fmt.Printf("value of c:%v\n", c)
}
