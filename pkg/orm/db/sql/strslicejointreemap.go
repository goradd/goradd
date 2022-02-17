package sql

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// A JoinTreeItemSliceMap combines a map with a slice so that you can range over a
// map in a predictable order. By default, the order will be the same order that items were inserted,
// i.e. a FIFO list. This is similar to how PHP arrays work.
// JoinTreeItemSliceMap implements the sort interface so you can change the order
// before ranging over the values if desired.
// It is NOT safe for concurrent use.
// The zero of this is usable immediately.
// The JoinTreeItemSliceMap satisfies the JoinTreeItemMapI interface.
type JoinTreeItemSliceMap struct {
	items map[string]*JoinTreeItem
	order []string
	lessF func(key1,key2 string, val1, val2 *JoinTreeItem) bool
}

// NewJoinTreeItemSliceMap creates a new map that maps string's to *JoinTreeItem's.
func NewJoinTreeItemSliceMap() *JoinTreeItemSliceMap {
	return new (JoinTreeItemSliceMap)
}

// NewJoinTreeItemSliceMapFrom creates a new JoinTreeItemMap from a
// JoinTreeItemMapI interface object
func NewJoinTreeItemSliceMapFrom(i JoinTreeItemMapI) *JoinTreeItemSliceMap {
	m := new (JoinTreeItemSliceMap)
	m.Merge(i)
	return m
}

// NewJoinTreeItemSliceMapFromMap creates a new JoinTreeItemSliceMap from a
// GO map[string]*JoinTreeItem object. Note that this will pass control of the given map to the
// new object. After you do this, DO NOT change the original map.
func NewJoinTreeItemSliceMapFromMap(i map[string]*JoinTreeItem) *JoinTreeItemSliceMap {
	m := NewJoinTreeItemSliceMap()
	m.items = i
	m.order = make([]string, len(m.items), len(m.items))
	j := 0
	for k := range m.items {
	    m.order[j] = k
	    j++
	}
	return m
}

// SetSortFunc sets the sort function which will determine the order of the items in the map
// on an ongoing basis. Normally, items will iterate in the order they were added.
// The sort function is a Less function, that returns true when item 1 is "less" than item 2.
// The sort function receives both the keys and values, so it can use either to decide how to sort.
func (o *JoinTreeItemSliceMap) SetSortFunc(f func(key1,key2 string, val1, val2 *JoinTreeItem) bool) *JoinTreeItemSliceMap {
    o.lessF = f
    if f != nil && len(o.order) > 0 {
        sort.Slice(o.order, func(i,j int) bool {
            return f(o.order[i], o.order[j], o.items[o.order[i]], o.items[o.order[j]])
        })
    }

    return o
}

// SortByKeys sets up the map to have its sort order sort by keys, lowest to highest
func (o *JoinTreeItemSliceMap) SortByKeys() *JoinTreeItemSliceMap {
    o.SetSortFunc(keySortJoinTreeItemSliceMap)
    return o
}

func keySortJoinTreeItemSliceMap(key1, key2 string, val1, val2 *JoinTreeItem) bool {
    return key1 < key2
}





// Set sets the given key to the given value.
// If the key already exists, the range order will not change.
func (o *JoinTreeItemSliceMap) Set(key string, val *JoinTreeItem) {
	var ok bool
	var oldVal *JoinTreeItem

	if o == nil {
	    panic("You must initialize the map before using it.")
	}

	if o.items == nil {
	    o.items = make(map[string]*JoinTreeItem)
	}

	_, ok = o.items[key]
    if o.lessF != nil {
        if ok {
            // delete old key location
            loc := sort.Search (len(o.items), func(n int) bool {
                return !o.lessF(o.order[n], key, o.items[o.order[n]], oldVal)
            })
            o.order = append(o.order[:loc], o.order[loc+1:]...)
        }

        loc := sort.Search (len(o.order), func(n int) bool {
            return o.lessF(key, o.order[n], val, o.items[o.order[n]])
        })
        // insert
        o.order = append(o.order, key)
        copy(o.order[loc+1:], o.order[loc:])
        o.order[loc] = key
    } else {
        if !ok {
            o.order = append(o.order, key)
        }
    }
    o.items[key] = val

	return
}

// SetAt sets the given key to the given value, but also inserts it at the index specified.  If the index is bigger than
// the length, it puts it at the end. Negative indexes are backwards from the end.
func (o *JoinTreeItemSliceMap) SetAt(index int, key string, val *JoinTreeItem)  {
    if o == nil {
        panic("You must initialize the map before using it.")
    }

    if o.lessF != nil {
        panic("You cannot use SetAt if you are also using a sort function.")
    }

	if index >= len(o.order) {
		o.Set(key, val)
		return
	}

	var ok bool
	var emptyKey string

	if _, ok = o.items[key]; !ok {
		if index <= -len(o.items) {
			index = 0
		}
		if index < 0 {
			index = len(o.items) + index
		}

		o.order = append(o.order, emptyKey)
		copy(o.order[index+1:], o.order[index:])
		o.order[index] = key
	}
	o.items[key] = val
    return
}

// Delete removes the item with the given key.
func (o *JoinTreeItemSliceMap) Delete(key string) {
    if o == nil {
        return
    }

    if _,ok := o.items[key]; ok {
        if o.lessF != nil {
            oldVal := o.items[key]
            loc := sort.Search (len(o.items), func(n int) bool {
                return !o.lessF(o.order[n], key, o.items[o.order[n]], oldVal)
            })
            o.order = append(o.order[:loc], o.order[loc+1:]...)
        } else {
            for i, v := range o.order {
                if v == key {
                    o.order = append(o.order[:i], o.order[i+1:]...)
                    break
                }
            }
        }
        delete(o.items, key)
    }
}

// Get returns the value based on its key. If the key does not exist, an empty value is returned.
func (o *JoinTreeItemSliceMap) Get(key string) (val *JoinTreeItem) {
    val,_ = o.Load(key)
    return
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.Map.Load()
func (o *JoinTreeItemSliceMap) Load(key string) (val *JoinTreeItem, ok bool) {
    if o == nil {
        return
    }
    if o.items != nil {
    	val, ok = o.items[key]
    }
	return
}



// Has returns true if the given key exists in the map.
func (o *JoinTreeItemSliceMap) Has(key string) (ok bool) {
    if o == nil {
        return false
    }
    if o.items != nil {
	    _, ok = o.items[key]
    }
	return
}



// GetAt returns the value based on its position. If the position is out of bounds, an empty value is returned.
func (o *JoinTreeItemSliceMap) GetAt(position int) (val *JoinTreeItem) {
    if o == nil {
        return
    }
	if position < len(o.order) && position >= 0 {
		val, _ = o.items[o.order[position]]
	}
	return
}

// GetKeyAt returns the key based on its position. If the position is out of bounds, an empty value is returned.
func (o *JoinTreeItemSliceMap) GetKeyAt(position int) (key string) {
    if o == nil {
        return
    }
	if position < len(o.order) && position >= 0 {
		key = o.order[position]
	}
	return
}

// Values returns a slice of the values in the order they were added or sorted.
func (o *JoinTreeItemSliceMap) Values() (vals []*JoinTreeItem) {
    if o == nil {
        return
    }

    if o.items != nil {
  	    vals = make([]*JoinTreeItem, len(o.order))
        for i, v := range o.order {
            vals[i] = o.items[v]
        }
    }

	return
}

// Keys returns the keys of the map, in the order they were added or sorted
func (o *JoinTreeItemSliceMap) Keys() (keys []string) {
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
func (o *JoinTreeItemSliceMap) Len() int {
    if o == nil {
        return 0
    }
    l := len(o.order)
	return l
}


// Copy will make a copy of the map and a copy of the underlying data.
func (o *JoinTreeItemSliceMap) Copy() *JoinTreeItemSliceMap {
	cp := NewJoinTreeItemSliceMap()

	o.Range(func(key string, value *JoinTreeItem) bool {

		cp.Set(key, value)
		return true
	})
	cp.lessF = o.lessF
	return cp
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
// If you are using a sort function, you must save and restore the sort function in a separate operation
// since functions are not serializable.
func (o *JoinTreeItemSliceMap) MarshalBinary() (data []byte, err error) {
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
// JoinTreeItemSliceMap
func (o *JoinTreeItemSliceMap) UnmarshalBinary(data []byte) (err error) {
    var items map[string]*JoinTreeItem
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
func (o *JoinTreeItemSliceMap) MarshalJSON() (data []byte, err error) {
	// Json objects are unordered
	data, err = json.Marshal(o.items)
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a JoinTreeItemMap.
// The JSON must start with an object.
func (o *JoinTreeItemSliceMap) UnmarshalJSON(data []byte) (err error) {
    var items map[string]*JoinTreeItem

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
func (o *JoinTreeItemSliceMap) Merge(i JoinTreeItemMapI) {
	if i != nil {
		i.Range(func(k string, v *JoinTreeItem) bool {
			o.Set(k, v)
			return true
		})
	}
}

// MergeMap merges the given standard map with the current one. The given one takes precedent on collisions.
func (o *JoinTreeItemSliceMap) MergeMap(m map[string]*JoinTreeItem) {
	if m == nil {
		return
	}

	for k,v := range m {
		o.Set(k, v)
	}
}


// Range will call the given function with every key and value in the order
// they were placed in the map, or in if you sorted the map, in your custom order.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o *JoinTreeItemSliceMap) Range(f func(key string, value *JoinTreeItem) bool) {
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
func (o *JoinTreeItemSliceMap) Equals(i JoinTreeItemMapI) bool {
	l := i.Len()
	if l == 0 {
		return o == nil
	}
	if l != o.Len() {
		return false
	}
	var ret = true

	o.Range(func(k string, v *JoinTreeItem) bool {
		if !i.Has(k) || v != i.Get(k) {
			ret = false
			return false
		}
		return true
	})
	return ret
}

func (o *JoinTreeItemSliceMap) Clear() {
    if o == nil {return}
	o.items = nil
	o.order = nil

}

func (o *JoinTreeItemSliceMap) IsNil() bool {
	return o == nil
}

func (o *JoinTreeItemSliceMap) String() string {
	var s string

	s = "{"
	o.Range(func(k string, v *JoinTreeItem) bool {
		s += fmt.Sprintf(`%#v:%#v,`, k, v)
		return true
	})
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}



func init() {
	gob.Register(new (JoinTreeItemSliceMap))
}