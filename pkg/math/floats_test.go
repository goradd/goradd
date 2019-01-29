package math

import "testing"

func TestRoundInt(t *testing.T) {
	tests := []struct {
		name string
		arg float64
		want int
	}{
		{"0.9", 0.9, 1},
		{"0.5", 0.5, 1},
		{"0.49", 0.49, 0},
		{"1.01", 1.01, 1},
		{"1.49", 1.49, 1},
		{"1.5", 1.5, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundInt(tt.arg); got != tt.want {
				t.Errorf("RoundInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundFloat(t *testing.T) {
	type args struct {
		f      float64
		digits int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"0.9", args{0.9, 0}, 1},
		{"0.5", args{0.5, 0}, 1},
		{"0.49", args{0.49, 0}, 0},
		{"1.01", args{1.01, 0}, 1},
		{"1.49", args{1.49, 0}, 1},
		{"1.5", args{1.5, 0}, 2},
		{"0.9", args{0.9, 1}, 0.9},
		{"0.5", args{0.5, 1}, 0.5},
		{"0.45", args{0.45, 1}, 0.5},
		{"1.01", args{1.01, 1}, 1},
		{"1.49", args{1.49, 1}, 1.5},
		{"1.495", args{1.495, 2}, 1.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundFloat(tt.args.f, tt.args.digits); got != tt.want {
				t.Errorf("RoundFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
