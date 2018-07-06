package types

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)


type Getter interface {
	Get(key string) (val interface{})
}

// MapI is a generic map interface that can store interface{} objects. It has helpers to save and restore built-in types too.
// See Map, OrderedMap, Map, and SafeOrderedMap for various implementations of the interface
type MapI interface {
	Clear()
	Set(key string, val interface{}) MapI
	Remove(key string)
	Get(key string) (val interface{})
	GetString(key string) (val string, typeOk bool)
	GetBool(key string) (val bool, typeOk bool)
	GetInt(key string) (val int, typeOk bool)
	GetFloat(key string) (val float64, typeOk bool)
	Has(key string) (ok bool)
	Values() []interface{}
	Keys() []string
	Len() int
	Range(f func(key string, value interface{}) bool)
	Copy() MapI
}

// Map is your basic GoMap with a read/write mutex so that it can read and write concurrently.
// Go now has a sync.Map item, but that is primarily for situations of high read contention on a large amount of
// cores.
type Map struct {
	items map[string]interface{}
}

func NewMap() *Map {
	return &Map{items: make(map[string]interface{})}
}

func (o *Map) Clear() {
	o.items = make(map[string]interface{})
}

// Set sets the value
func (o *Map) Set(key string, val interface{}) MapI {
	if o.items == nil {
		o.items = make(map[string]interface{})
	}

	o.items[key] = val
	return o
}

func (o *Map) Remove(key string) {
	if o.items == nil {
		return
	}
	delete(o.items, key)
}

// Get returns the string based on its key. If it does not exist, will return a nil interface{}
func (o *Map) Get(key string) (val interface{}) {
	if o.items == nil {
		return nil
	}
	val, _ = o.items[key]
	return
}

// Return a string, or the default value if not found. If the value was found, but is not a string, returns false in typeOk.
func (o *Map) GetString(key string) (val string, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(string)
		return
	} else {
		return "", true
	}
}

// Return a bool, or the default value if not found. If the value was found, but is not a bool, returns false in typeOk.
func (o *Map) GetBool(key string) (val bool, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(bool)
		return
	} else {
		return false, true
	}
}

// Return a int, or the default value if not found. If the value was found, but is not a int, returns false in typeOk.
func (o *Map) GetInt(key string) (val int, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(int)
		return
	} else {
		return 0, true
	}
}

// Return a float64, or the default value if not found. If the value was found, but is not a float64, returns false in typeOk.
func (o *Map) GetFloat(key string) (val float64, typeOk bool) {
	if v := o.Get(key); v != nil {
		val, typeOk = v.(float64)
		return
	} else {
		return 0, true
	}
}

func (o *Map) Has(key string) (ok bool) {
	if o.items == nil {
		return false
	}

	_, ok = o.items[key]
	return
}

// Values returns a slice of the values
func (o *Map) Values() []interface{} {
	vals := make([]interface{}, 0, len(o.items))

	for _, v := range o.items {
		vals = append(vals, v)
	}
	return vals
}

// Keys returns a slice of they keys
func (o *Map) Keys() []string {
	vals := make([]string, 0, len(o.items))

	for i := range o.items {
		vals = append(vals, i)
	}
	return vals
}

func (o *Map) Len() int {
	return len(o.items)
}

// Range will call the given function with every key and value in the Map
// During this process, the map will be locked, so do not use a function that will be taking significant amounts of time
// If f returns false, it stops the iteration. This is taken from the sync.Map.
func (o *Map) Range(f func(key string, value interface{}) bool) {
	for k, v := range o.items {
		if !f(k, v) {
			break
		}
	}
}

func (o *Map) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

	enc := gob.NewEncoder(&b)
	err := enc.Encode(o.items)

	return b.Bytes(), err
}

func (o *Map) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	err := dec.Decode(&o.items)
	return err
}

// UnmarshalJSON implements the json.Unmarshaller interface so that you can output a json structure in the same order
// it was saved in. The golang json library saves a json object as a golang map, and golang maps are not guaranteed to
// be iterated in the same order they were created. This function remedies that. It requires that you give it
// a json object that begins with a { character.
func (o *Map) UnmarshalJSON(in []byte) error {

	b := bytes.TrimSpace(in)

	dec := json.NewDecoder(bytes.NewReader(b))
	t, err := dec.Token()
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

func (o *Map) getJsonMap(dec *json.Decoder) (err error) {
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

		value, err = o.getJsonToken(dec)

		if err != nil {
			return err
		}

		o.Set(key, value)
	}
	return nil
}

func (o *Map) getJsonToken(dec *json.Decoder) (ret interface{}, err error) {
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t.(type) {
	case json.Delim:
		d := t.(json.Delim)
		switch d {
		case '{':
			m := NewMap()
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
func (o *Map) MarshalJSON() (out []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	writer.WriteString("{")

	o.Range(func(k string, v interface{}) bool {
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

	out = append(out[:len(out)-2], out[len(out)-1]) // get rid of comma
	return out, nil
}

func (o *Map) Copy() MapI {
	cp := NewMap()

	o.Range(func(key string, value interface{}) bool {
		if copier, ok := value.(Copier); ok {
			value = copier.Copy()
		}
		cp.Set(key, value)
		return true
	})
	return cp
}

func (o *Map) IsNil() bool {
	return o == nil
}

func init() {
	gob.Register(NewMap())
}
