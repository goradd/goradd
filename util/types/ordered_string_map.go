package types

import (
	"strings"
	"sync"
	"sort"
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// An OrderedStringMap is similar to PHP's indexed arrays. You can get the strings
// by key, or by position. When you iterate over it, you will get the results back in the
// order the strings were put in. It is safe for concurrent use.
// It satisfies the StringMapI interface and the sort.Sort.Interface interface, so its sortable
type OrderedStringMap struct {
	sync.RWMutex
	items map[string]string
	order []string
}

func NewOrderedStringMap() *OrderedStringMap {
	return &OrderedStringMap{items: make(map[string]string)}
}

func NewOrderedStringMapFrom(i StringMapI) *OrderedStringMap {
	m := NewOrderedStringMap()
	m.Merge(i)
	return m
}

// Set sets the value, but also appends the value to the end of the list for when you
// iterate over the list
func (o *OrderedStringMap) Set(key string, val string) (changed bool, err error) {
	o.Lock()
	defer o.Unlock()

	var ok bool
	var oldVal string

	if oldVal,ok = o.items[key]; !ok && oldVal != val {
		o.order = append(o.order, key)
		o.items[key] = val
		changed = true
	}
	return
}

func (o *OrderedStringMap) Remove(key string) {
	o.Lock()
	defer o.Unlock()
	for i, v := range o.order {
		if v == key {
			o.order = append(o.order[:i], o.order[i + 1:]...)
			continue
		}
	}
	delete (o.items,key)
}

// Get returns the string based on its key. If it does not exist, will return the empty string
func (o *OrderedStringMap) Get(key string) (val string) {
	o.RLock()
	defer o.RUnlock()
	val, _ =  o.items[key]
	return
}

func (o *OrderedStringMap) Has(key string) (ok bool) {
	o.Lock()
	defer o.Unlock()
	_, ok =  o.items[key]
	return
}


// GetAt returns the string based on its position
// To see if a value exists, simply test that the position is less than the length
func (o *OrderedStringMap) GetAt(position int) (val string) {
	o.RLock()
	defer o.RUnlock()
	if position >= o.Len() {
		val = ""
	} else {
		val, _ = o.items[o.order[position]]
	}
	return
}


// Join is just like strings.Join
func (o *OrderedStringMap) Join(glue string) string {
	return strings.Join(o.Values(), glue)
}

// Strings returns a slice of the strings in the order they were added
func (o *OrderedStringMap) Values() []string {
	o.RLock()
	defer o.RUnlock()
	vals := make ([]string, len(o.order))

	for i, v := range o.order {
		vals[i] = o.items[v]
	}
	return vals
}


// Keys are the keys of the strings, in the order they were added
func (o *OrderedStringMap) Keys() []string {
	o.Lock()
	defer o.Unlock()
	vals := make ([]string, len(o.order))

	for i, v := range o.order {
		vals[i] = v
	}
	return vals
}

func (o *OrderedStringMap) Len() int {
	return len (o.order)
}

// Iter can be used with range to iterate over the strings in the order in which they were added.
// Note that we return a buffered channel the size of the return values so there is no blocking
func (o *OrderedStringMap) Iter() <-chan string {
	c := make(chan string, o.Len())

	f := func() {
		o.RLock()
		defer o.RUnlock()

		for _, v := range o.order {
			c <- o.items[v]
		}
		close(c)
	}
	go f()

	return c
}

// IterKeys can be used with range to iterate over the keys in the order in which they were added.
// You can then use Get(key) to get the actual value
// Note that we return a buffered channel the size of the return values so there is no blocking
func (o *OrderedStringMap) IterKeys() <-chan string {
	c := make(chan string, o.Len())

	f := func() {
		o.RLock()
		defer o.RUnlock()

		for _, v := range o.order {
			c <- v
		}
		close(c)
	}
	go f()

	return c
}

func (o *OrderedStringMap) Less(i, j int) bool {
	o.RLock()
	defer o.RUnlock()
	return o.items[o.order[i]] < o.items[o.order[j]]
}

func (o *OrderedStringMap) Swap(i, j int)      {
	o.Lock()
	defer o.Unlock()
	o.order[i], o.order[j] = o.order[j], o.order[i]
}


// Sort by keys interface
type orderedstringbykeys struct {
	// This embedded Interface permits Reverse to use the methods of
	// another Interface implementation.
	sort.Interface
}

// A helper function to allow OrderedStringMaps to be sorted by keys
func OrderStringMapByKeys(o *OrderedStringMap) sort.Interface {
	return &orderedstringbykeys{o}
}

// A helper function to allow OrderedStringMaps to be sorted by keys
// See the IterKeys example
func (r orderedstringbykeys) Less(i, j int) bool {
	var o *OrderedStringMap = r.Interface.(*OrderedStringMap)
	o.Lock()
	defer o.Unlock()
	return o.order[i] < o.order[j]
}

func (o *OrderedStringMap) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err = encoder.Encode(o.items)
	if err == nil {
		err = encoder.Encode(o.order)
	}
	data = buf.Bytes()
	return
}

func (o *OrderedStringMap) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf) // Will read from network.
	err := dec.Decode(&o.items)
	if err == nil {
		err = dec.Decode(&o.order)
	}
	return err
}

func (o *OrderedStringMap) MarshalJSON() (data []byte, err error) {
	// Json objects are unordered
	data, err = json.Marshal(o.items)
	return
}

func (o *OrderedStringMap) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &o.items)
	if err == nil {
		// Create a default order, since these are inherently unordered
		o.order = make ([]string, len(o.items))
		i := 0
		for k := range o.items {
			o.order[i] = k
			i++
		}
	}
	return err
}

func (o *OrderedStringMap) String() string {
	var s string

	s = "{"
	for k := range o.IterKeys() {
		s += `"` + k + `":"` + o.Get(k) + `",`
	}
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}


// Merge the given string map into the current one
// Can be any kind of string map
func (o *OrderedStringMap) Merge(i StringMapI) {
	for k := range i.IterKeys() {
		o.Set(k, i.Get(k))
	}
}