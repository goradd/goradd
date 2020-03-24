package strings

import "strconv"

// Convenience scripts for number conversion

func Atoui(s string) uint {
	v,_ := strconv.ParseUint(s, 10, 0)
	return uint(v)
}

func Atoui64(s string) uint64 {
	v,_ := strconv.ParseUint(s, 10, 64)
	return v
}

func Atoi64(s string) int64 {
	v,_ := strconv.ParseInt(s, 10, 64)
	return v
}


