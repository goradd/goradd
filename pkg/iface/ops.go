package iface

func If(cond bool, i1, i2 interface{}) interface{} {
	if cond {
		return i1
	} else {
		return i2
	}
}
