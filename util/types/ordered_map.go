package types

import (
	"sort"
	"encoding/json"
	"bytes"
	"errors"
	"bufio"
)

// An OrderedMap is similar to PHP's indexed arrays. You can get the strings
// by key, or by position. When you iterate over it, you will get the results back in the
// order the strings were put in. It is not safe for concurrent use, but since it is ordered, you will not likely
// want to concurrently add to it.
// It satisfies the StringMapI interface and the sort.Sort.Interface interface, so its sortable
type OrderedMap struct {
	items map[string]interface{}
	order []string
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{items: make(map[string]interface{})}
}

func (o *OrderedMap) Clear() {
	o.items = nil
	o.order = nil
}

// Set sets the value, but also appends the value to the end of the list for when you
// iterate over the list. If the value already exists, the value is replaced, the order does not change. If you want
// the order to change in this situation, you must Remove then Set.
func (o *OrderedMap) Set(key string, val interface{}) MapI {
	if o.items == nil {
		o.items = make(map[string]interface{})
	}

	var ok bool

	if _,ok = o.items[key]; !ok  {
		o.order = append(o.order, key)
	}
	o.items[key] = val
	return o
}

// SetAt sets the given key to the given value, but also inserts it at the index specified.  If the index is bigger than
// the length, or -1, it is the same as Set, in that it puts it at the end. Negative indexes are backwards from the
// end, if smaller than the negative length, just inserts at the beginning.
func (o *OrderedMap) SetAt(index int, key string, val interface{}) MapI {
	if index == -1 || index >= len(o.items) {
		return o.Set(key, val)
	}

	if o.items == nil {
		o.items = make(map[string]interface{})
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


// Remove an item from the list by key. If the item does not exist, nothing happens.
func (o *OrderedMap) Remove(key string) {
	if o.items == nil {
		return
	}
	for i, v := range o.order {
		if v == key {
			o.order = append(o.order[:i], o.order[i + 1:]...)
			continue
		}
	}
	delete (o.items,key)
}

func (o *OrderedMap) RemoveAt(offset int) {
	if offset < 0 || offset >= len(o.order) {
		return
	}
	o.Remove(o.order[offset])
}

// Get returns the string based on its key. If it does not exist, will return nil
func (o *OrderedMap) Get(key string) (val interface{}) {
	if o.items == nil {
		return nil
	}

	val, _ =  o.items[key]
	return
}

// Find returns the offest of the key in the list, or -1 if not found.
func (o *OrderedMap) Find(key string) int {
	for i,k := range o.order {
		if k == key {
			return i
		}
	}
	return -1
}


// Return a string, or the default value if not found. If the value was found, but is not a string, returns false in typeOk.
func (o *OrderedMap) GetString(key string) (val string, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(string)
		return
	} else {
		return "", true
	}
}

// Return a bool, or the default value if not found. If the value was found, but is not a bool, returns false in typeOk.
func (o *OrderedMap) GetBool(key string) (val bool, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(bool)
		return
	} else {
		return false, true
	}
}

// Return a int, or the default value if not found. If the value was found, but is not a int, returns false in typeOk.
func (o *OrderedMap) GetInt(key string) (val int, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(int)
		return
	} else {
		return 0, true
	}
}

// Return a float64, or the default value if not found. If the value was found, but is not a float64, returns false in typeOk.
func (o *OrderedMap) GetFloat(key string) (val float64, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(float64)
		return
	} else {
		return 0, true
	}
}


func (o *OrderedMap) Has(key string) (ok bool) {
	if o.items == nil {
		return false
	}

	_, ok =  o.items[key]
	return
}


// GetAt returns the string based on its position
// To see if a value exists, simply test that the position is less than the length
func (o *OrderedMap) GetAt(position int) (val interface{}) {
	if o.items == nil {
		return nil
	}
	if position >= o.Len() {
		val = ""
	} else {
		val, _ = o.items[o.order[position]]
	}
	return
}

// Return a string, or the default value if not found. If the value was not found, or is not a string, returns false in ok.
func (o *OrderedMap) GetStringAt(position int) (val string, ok bool) {
	if position < 0 || position >= len(o.order) {
		return "", false
	}

	val,ok = o.GetString(o.order[position])
	return
}

// Return a bool, or the default value if not found. If the value was found, but is not a bool, returns false in ok.
func (o *OrderedMap) GetBoolAt(position int) (val bool, ok bool) {
	if position < 0 || position >= len(o.order) {
		return false, false
	}
	val,ok = o.GetBool(o.order[position])
	return
}

// Return a int, or the default value if not found. If the value was found, but is not an int, returns false in ok.
func (o *OrderedMap) GetIntAt(position int) (val int, ok bool) {
	if position < 0 || position >= len(o.order) {
		return 0, false
	}
	val,ok = o.GetInt(o.order[position])
	return
}

// Return a float, or the default value if not found. If the value was found, but is not a float, returns false in ok.
func (o *OrderedMap) GetFloatAt(position int) (val float64, ok bool) {
	if position < 0 || position >= len(o.order) {
		return 0, false
	}
	val,ok = o.GetFloat(o.order[position])
	return
}

// Values returns a slice of the values in the order they were added
func (o *OrderedMap) Values() []interface{} {
	vals := make ([]interface{}, len(o.order))

	for i, v := range o.order {
		vals[i] = o.items[v]
	}
	return vals
}


// Keys are the keys of the items, in the order they were added
func (o *OrderedMap) Keys() []string {
	vals := make ([]string, len(o.order))
	copy(vals, o.order)
	return vals
}

func (o *OrderedMap) Len() int {
	return len (o.order)
}


// Range will call the given function with every key and value in the order they were placed into the OrderedMap
// During this process, the map will be locked, so do not use a function that will be taking significant amounts of time
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o *OrderedMap) Range(f func(key string, value interface{}) bool) {
	for _, k := range o.order {
		if !f(k, o.items[k]) {
			break
		}
	}
}

func (o *OrderedMap) Less(i, j int) bool {
	return o.order[i] < o.order[j]
}

func (o *OrderedMap) Swap(i, j int)      {
	o.order[i], o.order[j] = o.order[j], o.order[i]
}



// Sort by keys interface
type orderedbykeys struct {
	// This embedded Interface permits Reverse to use the methods of
	// another Interface implementation.
	sort.Interface
}

// A helper function to allow OrderedMaps to be sorted by keys
func OrderMapByKeys(o *OrderedMap) sort.Interface {
	return &orderedbykeys{o}
}

// A helper function to allow OrderedMaps to be sorted by keys
func (r orderedbykeys) Less(i, j int) bool {
	var o *OrderedMap = r.Interface.(*OrderedMap)
	return o.order[i] < o.order[j]
}


// UnmarshalJSON implements the json.Unmarshaller interface so that you can output a json structure in the same order
// it was saved in. The golang json library saves a json object as a golang map, and golang maps are not guaranteed to
// be iterated in the same order they were created. This function remedies that. It requires that you give it
// a json object that begins with a { character.
func (o *OrderedMap) UnmarshalJSON(in []byte) error {

	b := bytes.TrimSpace(in)

	dec := json.NewDecoder(bytes.NewReader(b))
	t,err := dec.Token()
	if err != nil {
		return err
	}
	if d, ok := t.(json.Delim); !ok {
		return errors.New("Must be a json object that starts with a '{'.")
	} else if d != '{' {
		return errors.New("Must be a json object that starts with a '{'.")
	}

	return o.getJsonMap(dec)
}

func (o *OrderedMap) getJsonMap(dec *json.Decoder) (err error) {
	var key string
	var ok bool
	var t json.Token
	var value interface{}
	//var d rune

	for dec.More() {
		t, err = dec.Token()
		if key, ok = t.(string); !ok {
			return errors.New("Must be an object with string keys.")
		}

		value,err = o.getJsonToken(dec)

		if err != nil {
			return err
		}

		o.Set(key, value)
	}
	return nil
}

func (o *OrderedMap) getJsonToken(dec *json.Decoder) (ret interface{}, err error) {
	t, err := dec.Token()
	if err !=  nil {
		return nil, err
	}
	switch t.(type) {
	case json.Delim:
		d := t.(json.Delim)
		switch d {
		case '{':
			m := NewOrderedMap()
			err = m.getJsonMap(dec)
			return m, err
		case '[':
			a := []interface{}{}
			for dec.More() {
				a2, err := o.getJsonToken(dec)
				if err != nil {
					return nil, err
				} else {
					a = append(a, a2)
				}
			}
			return a, nil
			//dec.Token() // should be closed paren
		default:
			return
		}

	default:
		ret = t
		return
	}
	return
}

// MarshalJSON implements the json.Marshaller interface to all order_maps to be output as structs in the same order they were saved.
func (o *OrderedMap) MarshalJSON() (out []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	writer.WriteString("{")

	o.Range(func(k string,v interface{}) bool {
		var b2 []byte
		writer.WriteString("\"" + k + "\":")
		if b2, err = json.Marshal(v); err != nil {
			return false
		}
		writer.Write(b2)
		writer.WriteString(",")
		return true
	})
	writer.WriteString("}")
	writer.Flush()

	if err != nil {
		return nil, err
	}
	out = b.Bytes()

	out = append(out[:len(out) - 2], out[len(out)-1]) // get rid of comma
	return out, nil
}


type Copier interface {
	Copy() interface{}
}
func (o *OrderedMap) Copy() MapI {
	cp := NewOrderedMap()

	o.Range(func (key string, value interface{}) bool {
		if copier,ok := value.(Copier); ok {
			value = copier.Copy()
		}
		cp.Set(key, value)
		return true
	})
	return cp
}