package db

import (
	"fmt"
	"strconv"
)

type numberIn interface {
	int | uint | int32 | uint32 | int64 | uint64 | float32 | float64
}

func getMinOption(v interface{}, opt interface{}) (interface{}, error) {
	return getMinMaxOption(v, opt, true)
}

func getMaxOption(v interface{}, opt interface{}) (interface{}, error) {
	return getMinMaxOption(v, opt, false)
}

func getMinMaxOption(curVal interface{}, newValIn interface{}, isMin bool) (newValOut interface{}, err error) {
	switch curValTyped := curVal.(type) {
	case int64:
		return getMinMaxOptionTyped(curValTyped, newValIn, isMin)
	case uint64:
		return getMinMaxOptionTyped(curValTyped, newValIn, isMin)
	case float64:
		return getMinMaxOptionTyped(curValTyped, newValIn, isMin)
	default:
		return curVal, fmt.Errorf("invalid min or max default value")
	}
}

func getMinMaxOptionTyped[I numberIn](curVal I, newValIn interface{}, isMin bool) (newValOut I, err error) {
	switch newValTyped := newValIn.(type) {
	case float64: // returned by json conversion of number
		if isMin {
			if newValTyped < float64(curVal) {
				return curVal, fmt.Errorf("min value is less than what the data type allows")
			}
		} else {
			if newValTyped > float64(curVal) {
				return curVal, fmt.Errorf("max value is more than what the data type allows")
			}
		}
		return I(newValTyped), nil
	case string:
		i, err2 := parseVal(curVal, newValTyped)
		if err2 != nil {
			return curVal, err2
		}
		if isMin {
			if i < curVal {
				return curVal, fmt.Errorf("min value is less than what the data type allows")
			}
		} else {
			if i > curVal {
				return curVal, fmt.Errorf("max value is more than what the data type allows")
			}
		}
		return i, nil
	default:
		return curVal, fmt.Errorf("value type must be either numeric or string")
	}
}

func parseVal[C numberIn](v C, s string) (C, error) {
	var i interface{} = v

	switch i.(type) {
	case int64:
		out, err := strconv.ParseInt(s, 10, 64)
		return C(out), err
	case uint64:
		out, err := strconv.ParseUint(s, 10, 64)
		return C(out), err
	case float64:
		out, err := strconv.ParseFloat(s, 64)
		return C(out), err
	default:
		return v, fmt.Errorf("invalid default value")
	}
}
