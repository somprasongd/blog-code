package learnmap

import "fmt"

func Learn() {
	m := map[string]string{
		"a": "apple",
		"b": "banana",
	}

	fmt.Printf("%v\n", m["a"])

	m["c"] = "cranberry"

	fmt.Printf("%#v\n", m)

	delete(m, "b")

	fmt.Printf("%#v\n", m)
}
