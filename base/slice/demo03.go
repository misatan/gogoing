package slice

import "fmt"

func init() {
	var array = []int{10, 20, 30, 40, 50}
	var s1 = array[0:2]
	var s2 = s1
	fmt.Printf("Before s1 = %v,Pointer = %p,len = %d,cap = %d\n", s1, &s1, len(s1), cap(s1))
	fmt.Printf("Before s2 = %v,Pointer = %p,len = %d,cap = %d\n", s2, &s2, len(s2), cap(s2))
	s2 = append(s2, 50)
	s2[1] += 10
	fmt.Printf("After s1 = %v,Pointer = %p,len = %d,cap = %d\n", s1, &s1, len(s1), cap(s1))
	fmt.Printf("After s2 = %v,Pointer = %p,len = %d,cap = %d\n", s2, &s2, len(s2), cap(s2))
	fmt.Printf("After array = %v,Pointer = %p,len = %d,cap = %d\n", array, &array, len(array), cap(array))
}
