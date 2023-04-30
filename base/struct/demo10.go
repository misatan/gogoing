package teststruct

import "fmt"

func main() {
	p := Anonymous{
		"yes",
		18,
	}
	fmt.Printf("%#v\n", p)
	fmt.Println(p.string, p.int)
}

// Anonymous 匿名字段：类型名作为字段名；结构体要求字段名必须唯一，所以每种字段类型最多只有一个匿名字段
type Anonymous struct {
	string
	int
}
