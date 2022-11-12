package math

import (
	"fmt"
	"testing"
)

func ExampleMinInt() {
	i, v := MinInt(50, 13, 9, 200, 7, 23)

	fmt.Printf("%d,%d", i, v)
	// Output: 4,7
}

func TestMinInt(t *testing.T) {
	i, v := MinInt[int]()
	if i != -1 || v != 0 {
		t.Error("MinInt test failed for zero elements")
	}

	i, v = MinInt[int](7)
	if i != 0 || v != 7 {
		t.Error("MinInt test failed for one element")
	}

	i, v = MinInt(-4, -5, 6)
	if i != 1 || v != -5 {
		t.Error("MinInt test failed for negative elements")
	}
}

func ExampleMaxInt() {
	i, v := MaxInt(50, 13, 9, 200, 7, 23)

	fmt.Printf("%d,%d", i, v)
	// Output: 3,200
}

func TestMaxInt(t *testing.T) {
	i, v := MaxInt[int]()
	if i != -1 || v != 0 {
		t.Error("MaxInt test failed for zero elements")
	}

	i, v = MaxInt(7)
	if i != 0 || v != 7 {
		t.Error("MaxInt test failed for one element")
	}

	i, v = MaxInt(-4, -5, -6)
	if i != 0 || v != -4 {
		t.Error("MaxInt test failed for negative elements")
	}
}

func ExampleDiffInts() {
	diffs := DiffInts(50, 13, 9, 200, 7, 23)

	fmt.Print(diffs)
	// Output: [37 4 -191 193 -16 -27]
}

func TestDiffInts(t *testing.T) {
	values := DiffInts[int]()
	if values != nil {
		t.Error("DiffInts test failed for empty list")
	}

	values = DiffInts(7)
	if fmt.Sprint(values) != "[7]" {
		t.Errorf("DiffInts test failed for single item list. Got %s", fmt.Sprint(values))
	}
}

func ExampleSumInts() {
	sums := SumInts(50, 13, -9, 200, 7, 23)

	fmt.Print(sums)
	// Output: [63 4 191 207 30 73]
}

func TestSumInts(t *testing.T) {
	values := SumInts[int]()
	if values != nil {
		t.Error("SumInts test failed for empty list")
	}

	values = SumInts(7)
	if fmt.Sprint(values) != "[7]" {
		t.Errorf("SumInts test failed for single item list. Got %s", fmt.Sprint(values))
	}
}

func ExampleSquareInt() {
	fmt.Print(SquareInt(2))
	// Output: 4
}

func ExampleCubeInt() {
	fmt.Print(CubeInt(2))
	// Output: 8
}

func ExampleSqSqInt() {
	fmt.Print(SqSqInt(2))
	// Output: 16
}

func ExampleAbsInt() {
	fmt.Print(AbsInt(-5), AbsInt(5))
	// Output: 5 5
}

func ExamplePowerInt() {
	fmt.Print(PowerInt(3, 3))
	// Output: 27
}

func TestPowerInt(t *testing.T) {
	tests := []struct {
		name  string
		base  int
		power int
		want  int
	}{
		{"-1,5", -1, 5, -1},
		{"3,3", 3, 3, 27},
		{"0,0", 0, 0, 0},
		{"5,1", 5, 1, 5},
		{"1,-2", 1, -2, 1},
		{"-1,-3", -1, -3, -1},
		{"5,0", 5, 0, 1},
		{"5,-3", 5, -3, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PowerInt(tt.base, tt.power); got != tt.want {
				t.Errorf("PowerInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
