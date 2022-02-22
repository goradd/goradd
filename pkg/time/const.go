package time

import "time"

// Const  is a constant that represents a time to be instantiated later.
type Const int

const (
	// Zero will generate a zero time
	Zero = Const(0)
	// Now will generate a current time
	Now = Const(1)
	// Current is a synonym for Now
	Current = Const(1)
)

// Time returns the time corresponding to the given constant
func (c Const) Time() (t time.Time) {
	switch c {
	case Zero:
		t = time.Time{}
	case Now:
		t = time.Now().UTC()
	}
	return t
}
