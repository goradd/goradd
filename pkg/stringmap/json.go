package stringmap

import (
	"encoding/json"
	"io/ioutil"
)

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

