package util

import (
	"math"
)

// Use RoundInt to round a float to an integer
func RoundInt(f float64) int {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int(f + math.Copysign(0.5, f))
}

// Use RoundFloat to round a float to another float with the given number of digits. Can be used to fix precision
// errors

func RoundFloat(f float64, digits int) float64 {
	f = f * math.Pow10(digits)
	if math.Abs(f) < 0.5 {
		return 0
	}
	v := int(f + math.Copysign(0.5, f))
	f = float64(v) / math.Pow10(digits)
	return f
}

func MinInt(is ...int) int {
	var min int = is[0]
	for _, i := range is[1:] {
		if i < min {
			min = i
		}
	}
	return min
}

func MaxInt(is ...int) int {
	var max int = is[0]
	for _, i := range is[1:] {
		if i > max {
			max = i
		}
	}
	return max
}
