package stringmap

import (
	"encoding/json"
	"io/ioutil"
)

// Converts a json map to a stringmap
// Integers come through as float64, so must be converted to int to be used
func FromJson(j []byte)  (out map[string]interface{}, err error) {
	err = json.Unmarshal(j, &out)
	return
}

func FromJsonFile(f string)  (out map[string]interface{}, err error) {
	c, err := ioutil.ReadFile(f)

	if err != nil {
		return
	}

	return FromJson(c)
}

