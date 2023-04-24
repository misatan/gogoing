package sliceDemo

import "fmt"

func init() {
	s := []int{}
	for i := 0; i < 4098; i++ {
		var oldCap = cap(s)
		temp := s
		s = append(s, i)
		var newCap = cap(s)
		if oldCap != newCap {
			fmt.Printf("oldPointer:%p,newPoint:%p,oldCap:%d,newCap:%d\n", &temp, &s, oldCap, newCap)
		}
	}
}
