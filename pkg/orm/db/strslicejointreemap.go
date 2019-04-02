package db

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strings"

	"fmt"

)

// A joinTreeItemSliceMap combines a map with a slice so that you can range over a
// map in a predictable order. By default, the order will be the same order that items were inserted,
// i.e. a FIFO list. This is similar to how PHP arrays work.
// joinTreeItemSliceMap implements the sort interface so you can change the order
// before ranging over the values if desired.
// It is NOT safe for concurrent use.
// The zero of this is usable immediately.
// The joinTreeItemSliceMap satisfies the joinTreeItemMapI interface.
type joinTreeItemSliceMap struct {
	items map[string]*joinTreeItem
	order []string
}

// NewjoinTreeItemSliceMap creates a new map that maps string's to *joinTreeItem's.
func NewjoinTreeItemSliceMap() *joinTreeItemSliceMap {
	return new (joinTreeItemSliceMap)
}

// NewjoinTreeItemSliceMapFrom creates a new joinTreeItemMap from a
// joinTreeItemMapI interface object
func NewjoinTreeItemSliceMapFrom(i joinTreeItemMapI) *joinTreeItemSliceMap {
	m := new (joinTreeItemSliceMap)
	m.Merge(i)
	return m
}

// NewjoinTreeItemSliceMapFromMap creates a new joinTreeItemSliceMap from a
// GO map[string]*joinTreeItem object. Note that this will pass control of the given map to the
// new object. After you do this, DO NOT change the original map.
func NewjoinTreeItemSliceMapFromMap(i map[string]*joinTreeItem) *joinTreeItemSliceMap {
	m := NewjoinTreeItemSliceMap()
	m.items = i
	m.order = make([]string, len(m.items), len(m.items))
	j := 0
	for k := range m.items {
	    m.order[j] = k
	    j++
	}
	return m
}



// SetChanged sets the value, but also appends the value to the end of the list.
// It returns true if something in the map changed. If the key
// was already in the map, the order will not change, but the value will be replaced. If you want the
// order to change, you must Delete then call SetChanged.
func (o *joinTreeItemSliceMap) SetChanged(key string, val *joinTreeItem) (changed bool) {
	var ok bool
	var oldVal *joinTreeItem

	if o == nil {
	    panic("You must initialize the map before using it.")
	}

	if o.items == nil {
	    o.items = make(map[string]*joinTreeItem)
	}

	if oldVal, ok = o.items[key]; !ok || oldVal != val {
		if !ok {
			o.order = append(o.order, key)
		}
		o.items[key] = val
		changed = true
	}

	return
}

// Set sets the given key to the given value.
// If the key already exists, the range order will not change.
func (o *joinTreeItemSliceMap) Set(key string, val *joinTreeItem) {
	o.SetChanged(key, val)
}

// SetAt sets the given key to the given value, but also inserts it at the index specified.  If the index is bigger than
// the length, or -1, it is the same as Set, in that it puts it at the end. Negative indexes are backwards from the
// end, if smaller than the negative length, just inserts at the beginning.
func (o *joinTreeItemSliceMap) SetAt(index int, key string, val *joinTreeItem)  {
    if o == nil {
        panic("You must initialize the map before using it.")
    }

	if index == -1 || index >= len(o.order) {
		o.Set(key, val)
		return
	}

	var ok bool
	var emptyKey string

	if _, ok = o.items[key]; !ok {
		if index < -len(o.items) {
			index = 0
		}
		if index < 0 {
			index = len(o.items) + index + 1
		}

		o.order = append(o.order, emptyKey)
		copy(o.order[index+1:], o.order[index:])
		o.order[index] = key
	}
	o.items[key] = val
    return
}

// Delete removes the item with the given key.
func (o *joinTreeItemSliceMap) Delete(key string) {
    if o == nil {
        return
    }
	for i, v := range o.order {
		if v == key {
			o.order = append(o.order[:i], o.order[i+1:]...)
			break
		}
	}
	delete(o.items, key)
}

// Get returns the value based on its key. If the key does not exist, an empty value is returned.
func (o *joinTreeItemSliceMap) Get(key string) (val *joinTreeItem) {
    val,_ = o.Load(key)
    return
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.Map.Load()
func (o *joinTreeItemSliceMap) Load(key string) (val *joinTreeItem, ok bool) {
    if o == nil {
        return
    }
    if o.items != nil {
    	val, ok = o.items[key]
    }
	return
}



// Has returns true if the given key exists in the map.
func (o *joinTreeItemSliceMap) Has(key string) (ok bool) {
    if o == nil {
        return false
    }
    if o.items != nil {
	    _, ok = o.items[key]
    }
	return
}

// GetAt returns the value based on its position. If the position is out of bounds, an empty value is returned.
func (o *joinTreeItemSliceMap) GetAt(position int) (val *joinTreeItem) {
    if o == nil {
        return
    }
	if position < len(o.order) && position >= 0 {
		val, _ = o.items[o.order[position]]
	}
	return
}


// Values returns a slice of the values in the order they were added or sorted.
func (o *joinTreeItemSliceMap) Values() (vals []*joinTreeItem) {
    if o == nil {
        return
    }

    if o.items != nil {
  	    vals = make([]*joinTreeItem, len(o.order))
        for i, v := range o.order {
            vals[i] = o.items[v]
        }
    }

	return
}

// Keys returns the keys of the map, in the order they were added or sorted
func (o *joinTreeItemSliceMap) Keys() (keys []string) {
    if o == nil {
        return
    }

    if len(o.order) != 0 {
 	    keys = make([]string, len(o.order))
        for i, v := range o.order {
            keys[i] = v
        }
    }

	return
}

// Len returns the number of items in the map
func (o *joinTreeItemSliceMap) Len() int {
    if o == nil {
        return 0
    }
    l := len(o.order)
	return l
}

// Less is part of the interface that allows the map to be sorted by keys.
// It returns true if the value at position i should be sorted before the value at position j.
func (o *joinTreeItemSliceMap) Less(i, j int) bool {


	return o.order[i] < o.order[j]

}

// Swap is part of the interface that allows the slice to be sorted. It swaps the positions
// of the items and position i and j.
func (o *joinTreeItemSliceMap) Swap(i, j int) {
	o.order[i], o.order[j] = o.order[j], o.order[i]
}

 

// Copy will make a copy of the map and a copy of the underlying data.
func (o *joinTreeItemSliceMap) Copy() *joinTreeItemSliceMap {
	cp := NewjoinTreeItemSliceMap()

	o.Range(func(key string, value *joinTreeItem) bool {

		cp.Set(key, value)
		return true
	})
	return cp
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
func (o *joinTreeItemSliceMap) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err = encoder.Encode(o.items)
	if err == nil {
		err = encoder.Encode(o.order)
	}
	data = buf.Bytes()
	return
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a
// joinTreeItemSliceMap
func (o *joinTreeItemSliceMap) UnmarshalBinary(data []byte) (err error) {
    var items map[string]*joinTreeItem
	var order []string

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&items); err == nil {
		err = dec.Decode(&order)
	}

	if err == nil {
        o.items = items
        o.order = order
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (o *joinTreeItemSliceMap) MarshalJSON() (data []byte, err error) {
	// Json objects are unordered
	data, err = json.Marshal(o.items)
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a joinTreeItemMap.
// The JSON must start with an object.
func (o *joinTreeItemSliceMap) UnmarshalJSON(data []byte) (err error) {
    var items map[string]*joinTreeItem

	if err = json.Unmarshal(data, &items); err == nil {
        o.items = items
        // Create a default order, since these are inherently unordered
        o.order = make([]string, len(o.items))
        i := 0
        for k := range o.items {
            o.order[i] = k
            i++
        }
	}
	return
}


// Merge the given map into the current one
func (o *joinTreeItemSliceMap) Merge(i joinTreeItemMapI) {
	if i != nil {
		i.Range(func(k string, v *joinTreeItem) bool {
			o.Set(k, v)
			return true
		})
	}
}

// Range will call the given function with every key and value in the order
// they were placed in the map, or in if you sorted the map, in your custom order.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o *joinTreeItemSliceMap) Range(f func(key string, value *joinTreeItem) bool) {
	if o == nil {
		return
	}
	if  o.items != nil {
        for _, k := range o.order {
            if !f(k, o.items[k]) {
                break
            }
        }
    }
}

// Equals returns true if the map equals the given map, paying attention only to the content of the
// map and not the order.
func (o *joinTreeItemSliceMap) Equals(i joinTreeItemMapI) bool {
	l := i.Len()
	if l == 0 {
		return o == nil
	}
	if l != o.Len() {
		return false
	}
	var ret = true

	o.Range(func(k string, v *joinTreeItem) bool {
		if !i.Has(k) || v != i.Get(k) {
			ret = false
			return false
		}
		return true
	})
	return ret
}

func (o *joinTreeItemSliceMap) Clear() {
    if o == nil {return}
	o.items = nil
	o.order = nil

}

func (o *joinTreeItemSliceMap) IsNil() bool {
	return o == nil
}

func (o *joinTreeItemSliceMap) String() string {
	var s string

	s = "{"
	o.Range(func(k string, v *joinTreeItem) bool {
		s += `"` + k + `":"` +
fmt.Sprintf("%v", v)+ `",`
		return true
	})
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}



func init() {
	gob.Register(new (joinTreeItemSliceMap))
}