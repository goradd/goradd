package math

type ints interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// MinInt returns the minimum value from a slice of ints, and the zero-based index of that value.
// The index will be -1 if there are no values given.
func MinInt[M ints](values ...M) (index int, value M) {
	if len(values) == 0 {
		return -1, 0
	}
	value = values[0]

	for i, v := range values[1:] {
		if v < value {
			value = v
			index = i + 1
		}
	}
	return
}

// MaxInt returns the maximum value from a slice of ints, and the index of that value.
// The index will be -1 if no items are given.
func MaxInt[M ints](values ...M) (index int, value M) {
	if len(values) == 0 {
		return -1, 0
	}
	value = values[0]

	for i, v := range values[1:] {
		if v > value {
			value = v
			index = i + 1
		}
	}
	return
}

// DiffInts returns the difference between each item, and the next item in the list. The last diff is between the last item
// and the first item. Returns the differences as a slice. If there is only one item, it will return just the one item, and
// not diff it with itself.
func DiffInts[M ints](values ...M) (diffs []M) {
	if len(values) == 0 {
		return nil
	}
	diffs = make([]M, len(values), len(values))
	if len(values) == 1 {
		diffs[0] = values[0]
		return
	}

	for i, v := range values[:len(values)-1] {
		diffs[i] = v - values[i+1]
	}
	diffs[len(values)-1] = values[len(values)-1] - values[0]
	return
}

// SumInts returns the sum between each item, and the next item in the slice. The last sum is between the last item
// and the first item. If there is only one item, it will just return the one item, and not sum it with itself.
func SumInts[M ints](values ...M) (sums []M) {
	if len(values) == 0 {
		return nil
	}
	sums = make([]M, len(values), len(values))
	if len(values) == 1 {
		sums[0] = values[0]
		return
	}

	for i, v := range values[:len(values)-1] {
		sums[i] = values[i+1] + v
	}
	sums[len(values)-1] = values[len(values)-1] + values[0]
	return
}

// SquareInt returns the square of the given integer
func SquareInt[M ints](a M) M {
	return a * a
}

// CubeInt returns the cube of the given integer
func CubeInt[M ints](a M) M {
	return a * a * a
}

// SqSqInt returns the given integer to its fourth power
func SqSqInt[M ints](a M) M {
	return a * a * a * a
}

// PowerInt returns the given base integer raised to the given power.
//
// If your base integer is zero, the result will be zero regardless of the power.
// Fractions will return as zero.
func PowerInt[M ints](base, power M) M {
	if base == 0 {
		return 0 // ignore undefined situations of power <= 0 in this case, since there is no int NaN
	}
	if power < 0 {
		if base > 1 || int(base) < -1 {
			return 0
		} else if base == 1 {
			return 1
		} else {
			// base is -1
			return PowerInt(base, -power)
		}
	} else if power == 0 {
		return 1
	}

	v := base
	var i M
	for i = 1; i < power; i++ {
		v = v * base
	}
	return v
}

// AbsInt returns the absolute value of the given int.
func AbsInt[M ints](a M) M {
	if a < 0 {
		return -a
	}
	return a
}
