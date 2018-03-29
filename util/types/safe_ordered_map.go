package types

import (
	"sync"
	"sort"
	"encoding/json"
	"bytes"
	"errors"
	"bufio"
)

// A SafeOrderedMap is similar to PHP's indexed arrays. You can get the strings
// by key, or by position. When you iterate over it, you will get the results back in the
// order the strings were put in. It is safe for concurrent use. However, since it is ordered, there are not many
// situations where you would want to do concurrent writes to the map.
// It satisfies the StringMapI interface and the sort.Sort.Interface interface, so its sortable
type SafeOrderedMap struct {
	sync.RWMutex
	items map[string]interface{}
	order []string
}

func NewSafeOrderedMap() *SafeOrderedMap {
	return &SafeOrderedMap{items: make(map[string]interface{})}
}

func (o *SafeOrderedMap) Clear() {
	o.Lock()
	defer o.Unlock()
	o.items = nil
	o.order = nil
}

// Set sets the value, but also appends the value to the end of the list for when you
// iterate over the list. If the value already exists, the order does not change
func (o *SafeOrderedMap) Set(key string, val interface{}) MapI {
	o.Lock()
	defer o.Unlock()

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


func (o *SafeOrderedMap) Remove(key string) {
	if o.items == nil {
		return
	}
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

// Get returns the string based on its key. If it does not exist, will return nil
func (o *SafeOrderedMap) Get(key string) (val interface{}) {
	o.RLock()
	defer o.RUnlock()
	if o.items == nil {
		return nil
	}

	val, _ =  o.items[key]
	return
}

// Return a string, or the default value if not found. If the value was found, but is not a string, returns false in typeOk.
func (o *SafeOrderedMap) GetString(key string) (val string, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(string)
		return
	} else {
		return "", true
	}
}

// Return a bool, or the default value if not found. If the value was found, but is not a bool, returns false in typeOk.
func (o *SafeOrderedMap) GetBool(key string) (val bool, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(bool)
		return
	} else {
		return false, true
	}
}

// Return a int, or the default value if not found. If the value was found, but is not a int, returns false in typeOk.
func (o *SafeOrderedMap) GetInt(key string) (val int, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(int)
		return
	} else {
		return 0, true
	}
}

// Return a float64, or the default value if not found. If the value was found, but is not a float64, returns false in typeOk.
func (o *SafeOrderedMap) GetFloat(key string) (val float64, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(float64)
		return
	} else {
		return 0, true
	}
}


func (o *SafeOrderedMap) Has(key string) (ok bool) {
	o.RLock()
	defer o.RUnlock()
	if o.items == nil {
		return false
	}

	_, ok =  o.items[key]
	return
}


// GetAt returns the string based on its position
// To see if a value exists, simply test that the position is less than the length
func (o *SafeOrderedMap) GetAt(position int) (val interface{}) {
	o.RLock()
	defer o.RUnlock()

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


// Values returns a slice of the values in the order they were added
func (o *SafeOrderedMap) Values() []interface{} {
	o.Lock()
	defer o.Unlock()

	vals := make ([]interface{}, len(o.order))

	for i, v := range o.order {
		vals[i] = o.items[v]
	}
	return vals
}


// Keys are the keys of the strings, in the order they were added
func (o *SafeOrderedMap) Keys() []string {
	o.Lock()
	defer o.Unlock()
	vals := make ([]string, len(o.order))

	for i, v := range o.order {
		vals[i] = v
	}
	return vals
}

func (o *SafeOrderedMap) Len() int {
	return len (o.order)
}

// Range will call the given function with every key and value in the order they were placed into the SafeOrderedMap
// During this process, the map will be locked, so do not use a function that will be taking significant amounts of time
// If f returns false, it stops the iteration. This is taken from the sync.Map thingy
func (o *SafeOrderedMap) Range(f func(key string, value interface{}) bool) {
	o.RLock()
	defer o.RUnlock()

	for _, k := range o.order {
		if !f(k, o.items[k]) {
			break
		}
	}
}

func (o *SafeOrderedMap) Less(i, j int) bool {
	o.RLock()
	defer o.RUnlock()
	return o.order[i] < o.order[j]
}

func (o *SafeOrderedMap) Swap(i, j int)      {
	o.Lock()
	defer o.Unlock()
	o.order[i], o.order[j] = o.order[j], o.order[i]
}



// Sort by keys interface
type safeOrderedByKeys struct {
	// This embedded Interface permits Reverse to use the methods of
	// another Interface implementation.
	sort.Interface
}

// A helper function to allow SafeOrderedMaps to be sorted by keys
func SafeOrderMapByKeys(o *SafeOrderedMap) sort.Interface {
	return &safeOrderedByKeys{o}
}

// A helper function to allow SafeOrderedMaps to be sorted by keys
// See the IterKeys example
func (r safeOrderedByKeys) Less(i, j int) bool {
	var o *SafeOrderedMap = r.Interface.(*SafeOrderedMap)
	o.RLock()
	defer o.RUnlock()
	return o.order[i] < o.order[j]
}


// UnmarshalJSON implements the json.Unmarshaller interface so that you can output a json structure in the same order
// it was saved in. The golang json library saves a json object as a golang map, and golang maps are not guaranteed to
// be iterated in the same order they were created. This function remedies that. It requires that you give it
// a json object that begins with a { character.
func (o *SafeOrderedMap) UnmarshalJSON(in []byte) error {

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

func (o *SafeOrderedMap) getJsonMap(dec *json.Decoder) (err error) {
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

func (o *SafeOrderedMap) getJsonToken(dec *json.Decoder) (ret interface{}, err error) {
	t, err := dec.Token()
	if err !=  nil {
		return nil, err
	}
	switch t.(type) {
	case json.Delim:
		d := t.(json.Delim)
		switch d {
		case '{':
			m := NewSafeOrderedMap()
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
func (o *SafeOrderedMap) MarshalJSON() (out []byte, err error) {
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


func (o *SafeOrderedMap) Copy() MapI {
	cp := NewSafeOrderedMap()

	o.Range(func (key string, value interface{}) bool {
		if copier,ok := value.(Copier); ok {
			value = copier.Copy()
		}
		cp.Set(key, value)
		return true
	})
	return cp
}