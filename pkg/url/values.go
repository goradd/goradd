package url

import (
	"fmt"
	"net/url"
	"strconv"
)

// GetString extracts the named string out of values.
func GetString(values url.Values, name string) (value string, err error) {
	if v, ok := values[name]; !ok {
		err = fmt.Errorf("%s was not found", name)
		return
	} else {
		return v[0], nil
	}
}

// GetString extracts the named int out of values.
func GetInt(values url.Values, name string) (value int, err error) {
	if v, ok := values[name]; !ok {
		err = fmt.Errorf("%s was not found", name)
	} else {
		value, err = strconv.Atoi(v[0])
	}
	return
}

// GetString extracts the named float64 out of values.
func GetFloat64(values url.Values, name string) (value float64, err error) {
	if v, ok := values[name]; !ok {
		err = fmt.Errorf("%s was not found", name)
	} else {
		value, err = strconv.ParseFloat(v[0], 64)
	}
	return
}

// GetStrings returns the string array associated with the name
func GetStrings(values url.Values, name string) (value []string, err error) {
	if v, ok := values[name]; !ok {
		err = fmt.Errorf("%s was not found", name)
		return
	} else {
		return v, nil
	}
}
