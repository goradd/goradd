package types

import (
	"sync"
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// SafeStringMap is an implementention of the map of strings keyed by strings. It is concurrency
// safe, supports various synchronization modes, and implements the StringMapI interface
type SafeStringMap struct {
	items map[string]string
	sync.RWMutex
}

func NewSafeStringMap() *SafeStringMap {
	return &SafeStringMap{items: make(map[string]string)}
}

func NewSafeStringMapFrom(i StringMapI) *SafeStringMap {
	m := NewSafeStringMap()
	m.Merge(i)
	return m
}

func (o *SafeStringMap) SetChanged(key string, val string) (changed bool, err error) {
	o.Lock()
	defer o.Unlock()
	var ok bool

	var oldVal string

	if oldVal,ok = o.items[key]; !ok || oldVal != val {
		o.items[key] = val
		changed = true
	}
	return
}

func (o *SafeStringMap) Set(key string, val string) *SafeStringMap {
	o.Lock()
	defer o.Unlock()

	o.items[key] = val
	return o
}

// Get returns the string based on its key. If it does not exist, an empty string will be returned.
func (o *SafeStringMap) Get(key string) (val string) {
	o.Lock()
	defer o.Unlock()
	val, _ =  o.items[key]
	return
}

func (o *SafeStringMap) Has(key string) (exists bool) {
	o.Lock()
	defer o.Unlock()
	_, exists =  o.items[key]
	return
}

func (o *SafeStringMap) Remove(key string) {
	o.Lock()
	defer o.Unlock()
	delete (o.items,key)
}

// Values returns a slice of the string values
func (o *SafeStringMap) Values() []string {
	o.Lock()
	defer o.Unlock()

	vals := make ([]string, 0, len(o.items))

	for _, v := range o.items {
		vals = append(vals, v)
	}
	return vals
}

// Values returns a slice of the string keys
func (o *SafeStringMap) Keys() []string {
	o.Lock()
	defer o.Unlock()

	keys := make ([]string, 0, len(o.items))

	for k := range o.items {
		keys = append(keys, k)
	}
	return keys
}



func (o *SafeStringMap) Len() int {
	return len (o.items)
}


func (o *SafeStringMap) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err = encoder.Encode(o.items)
	data = buf.Bytes()
	return
}

func (o *SafeStringMap) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf) // Will read from network.
	err := dec.Decode(&o.items)
	return err
}

func (o *SafeStringMap) MarshalJSON() (data []byte, err error) {
	data, err = json.Marshal(o.items)
	return
}

func (o *SafeStringMap) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &o.items)
}

func (o *SafeStringMap) Merge(i StringMapI) {
	if o == nil {
		return
	}

	i.Range(func (k,v string) bool {
		o.Set(k, v)
		return true
	})
}

func (o *SafeStringMap) Equals(i StringMapI) bool {
	if i == nil {
		return o == nil
	}
	if i.Len() != o.Len() {
		return false
	}
	var ret bool = true

	o.Range(func (k,v string) bool {
		if !o.Has(k) || o.Get(k) != i.Get(k) {
			ret = false	// don't just return because we are in a channel and we want to use up the channel
			return false
		}
		return true
	})
	return ret
}

// Range will call the given function with every key and value in the order they were placed into the SafeString
// During this process, the map will be locked, so do not use a function that will be taking significant amounts of time
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o *SafeStringMap) Range(f func(key string, value string) bool) {
	if o == nil {
		return
	}

	o.Lock()
	defer o.Unlock()

	for k, v := range o.items {
		if !f(k, v) {
			break
		}
	}
}
