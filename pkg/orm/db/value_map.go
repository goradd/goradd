package db

import "encoding/gob"

// Copier implements the copy interface, that returns a deep copy of an object.
type Copier interface {
	Copy() interface{}
}

type ValueMap map[string]interface{}

func NewValueMap() ValueMap {
	return make(ValueMap)
}

// Copy does a deep copy and supports the deep copy interface
func (m ValueMap) Copy() interface{} {
	vm := ValueMap{}
	for k, v := range m {
		if c, ok := v.(Copier); ok {
			v = c.Copy()
		}
		vm[k] = v
	}
	return vm
}

func init() {
	gob.Register(&ValueMap{})
}

