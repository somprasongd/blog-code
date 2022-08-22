package math

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	// Arrange
	input1, input2 := 4, 6
	want := 10

	// Act
	got := Add(input1, input2)

	// Assert
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(4, 6)
	}
}
func ExampleAdd() {
	result := Add(4, 6)
	fmt.Println(result)
	// Output: 10
}
