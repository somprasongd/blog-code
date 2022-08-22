package variadic

import "fmt"

func Learn() {
	x := []int{1, 2, 3, 4, 5}

	y := sum(x...)

	fmt.Println(y)
}

func sum(x ...int) int {
	s := 0
	for _, v := range x {
		s += v
	}
	return s
}
