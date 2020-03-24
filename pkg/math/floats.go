package math

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
