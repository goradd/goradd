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

// SetChanged sets the value, but also appends the value to the end of the list for when you
// iterate over the list. Returns whether something changed, and if an error occurred. If the key
// was already in the map, the order will not change, but the value will be replaced. If you want the
// order to change, you must Remove then SetChanged
func (o *OrderedStringMap) SetChanged(key string, val string) (changed bool, err error) {
	o.Lock()
	defer o.Unlock()

	var ok bool
	var oldVal string

	if oldVal,ok = o.items[key]; !ok || oldVal != val {
		if !ok {
			o.order = append(o.order, key)
		}
		o.items[key] = val
		changed = true
	}
	return
}

// Set sets the given key to the given value, and returns the OrderedStringMap for chaining
func (o *OrderedStringMap) Set(key string, val string) *OrderedStringMap {
	o.SetChanged(key, val)
	return o
}

// SetAt sets the given key to the given value, but also inserts it at the index specified.  If the index is bigger than
// the length, or -1, it is the same as Set, in that it puts it at the end. Negative indexes are backwards from the
// end, if smaller than the negative length, just inserts at the beginning.
func (o *OrderedStringMap) SetAt(index int, key string, val string) *OrderedStringMap {
	if index == -1 || index >= len(o.items) {
		return o.Set(key, val)
	}

	var ok bool

	if _,ok = o.items[key]; !ok  {
		if index < -len(o.items) {
			index = 0
		}
		if index < 0 {
			index = len(o.items) + index + 1
		}

		o.order = append(o.order, "")
		copy(o.order[index+1:], o.order[index:])
		o.order[index] = key
	}
	o.items[key] = val
	return o
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
	o.Range(func(k,v string) bool {
		s += `"` + k + `":"` + o.Get(k) + `",`
		return true
	})
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}


// Merge the given string map into the current one
// Can be any kind of string map
func (o *OrderedStringMap) Merge(i StringMapI) {
	i.Range(func(k,v string) bool {
		o.Set(k, v)
		return true
	})
}

// Range will call the given function with every key and value in the order they were placed into the OrderedMap
// During this process, the map will be locked, so do not use a function that will be taking significant amounts of time
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o *OrderedStringMap) Range(f func(key string, value string) bool) {
	if o == nil {
		return
	}
	for _, k := range o.order {
		if !f(k, o.items[k]) {
			break
		}
	}
}


// Equals returns true if the map equals the given map, paying attention only to the content of the map and not the order.
func (o *OrderedStringMap) Equals(i StringMapI) bool {
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

func (o *OrderedStringMap) Clear() {
	o.Lock()
	defer o.Unlock()

	o.items = make(map[string]string)
	o.order = nil
}

func (o *OrderedStringMap) IsNil() bool {
	return o == nil
}