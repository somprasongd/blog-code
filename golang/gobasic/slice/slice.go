package slice

import "fmt"

func Learn() {
	x := []int{1, 2, 3}
	fmt.Printf("%#v, len=%v, cap=%v\n", x, len(x), cap(x))
	x = append(x, 4)
	x = append(x, 5)

	y := len(x)

	fmt.Printf("%#v, len=%v, cap=%v\n", x, y, cap(x))

	a := x[0:2]

	fmt.Printf("%#v\n", a)
}
