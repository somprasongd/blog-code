package advfunc

import "fmt"

func LearnClosure() {
	counterFn := newCounterFunc()
	fmt.Println(counterFn())
	fmt.Println(counterFn())
	fmt.Println(counterFn())
}

func newCounterFunc() func() int {
	var i int

	return func() int {
		i++
		return i
	}
}
