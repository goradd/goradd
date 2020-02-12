package datetime

func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func NumLeaps(year int) int {
	return year / 4 - year / 100 + year / 400
}

