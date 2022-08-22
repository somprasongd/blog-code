package advfunc

import "fmt"

func LearnHOFParam() {
	s := HOFGreeting(func() string {
		return "Ball"
	})

	fmt.Println(s)
}

func HOFGreeting(nameFn func() string) string {
	return fmt.Sprintf("Hello %s", nameFn())
}
