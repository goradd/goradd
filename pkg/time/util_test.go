package time

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDayDiff(t *testing.T) {
	tests := []struct {
		name string
		dt1 time.Time
		dt2 time.Time
		want int
	}{
		{"very close day 1", NewDateTime(2022,4,2,23,59,59,0), NewDateTime(2022,4,3,0,0,0,0), -1},
		{"very far day 1", NewDateTime(2022,4,2,0,0,0,0), NewDateTime(2022,4,3,23,59,59,0), -1},
		{"big range", NewDateTime(1900,1,1,0,0,0,0), NewDateTime(1600,1,1,0,0,0,0), 300 * 365 + 73},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DayDiff(tt.dt1, tt.dt2); got != tt.want {
				t.Errorf("DayDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLeap(t *testing.T) {
	tests := []struct {
		name string
		year int
		want bool
	}{
		{"yes", 2000, true},
		{"no", 1999, false},
		{"1600", 1600, true},
		{"1700", 1700, false},
		{"1800", 1800, false},
		{"1500", 1500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLeap(tt.year); got != tt.want {
				t.Errorf("IsLeap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeOnly(t *testing.T) {
	t1 := NewDateTime(2022,1,3,6,0,9,0)
	t2 := TimeOnly(t1)
	assert.Equal(t, t2.Year(), 0)
	assert.Equal(t, t2.Second(), 9)
}

func TestDateOnly(t *testing.T) {
	t1 := NewDateTime(2022,1,3,6,0,9,0)
	t2 := DateOnly(t1)
	assert.Equal(t, t2.Year(), 2022)
	assert.Equal(t, t2.Second(), 0)
}