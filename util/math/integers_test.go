package math

import (
	"fmt"
	"testing"
)

func ExampleMath_MinInt() {
	i,v := MinInt(50,13,9,200,7,23)

	fmt.Printf("%d,%d",i,v)
	// Output: 4,7
}

func TestMinInt(t *testing.T) {
	i,v := MinInt()
	if i != -1 || v != 0 {
		t.Error("MinInt test failed for zero elements")
	}

	i,v = MinInt(7)
	if i != 0 || v != 7 {
		t.Error("MinInt test failed for one element")
	}

	i,v = MinInt(-4,-5,6)
	if i != 1 || v != -5 {
		t.Error("MinInt test failed for negative elements")
	}
}

func ExampleMath_MaxInt() {
	i,v := MaxInt(50,13,9,200,7,23)

	fmt.Printf("%d,%d",i,v)
	// Output: 3,200
}

func TestMaxInt(t *testing.T) {
	i,v := MaxInt()
	if i != -1 || v != 0 {
		t.Error("MaxInt test failed for zero elements")
	}

	i,v = MaxInt(7)
	if i != 0 || v != 7 {
		t.Error("MaxInt test failed for one element")
	}

	i,v = MaxInt(-4,-5,-6)
	if i != 0 || v != -4 {
		t.Error("MaxInt test failed for negative elements")
	}
}

func ExampleMath_DiffInts() {
	diffs := DiffInts(50,13,9,200,7,23)

	fmt.Print(diffs)
	// Output: [37 4 -191 193 -16 -27]
}

func TestDiffInts(t *testing.T) {
	values := DiffInts()
	if values != nil {
		t.Error("DiffInts test failed for empty list")
	}

	values = DiffInts(7)
	if fmt.Sprint(values) != "[7]" {
		t.Errorf("DiffInts test failed for single item list. Got %s", fmt.Sprint(values))
	}
}

func ExampleMath_SumInts() {
	sums := SumInts(50,13,-9,200,7,23)

	fmt.Print(sums)
	// Output: [63 4 191 207 30 73]
}

func TestSumInts(t *testing.T) {
	values := SumInts()
	if values != nil {
		t.Error("SumInts test failed for empty list")
	}

	values = SumInts(7)
	if fmt.Sprint(values) != "[7]" {
		t.Errorf("SumInts test failed for single item list. Got %s", fmt.Sprint(values))
	}
}

func ExampleMath_PowerInt() {
	fmt.Print(PowerInt(-1,5), PowerInt(3,3), PowerInt(0,0))
	// Output: -1 27 0
}

