package types

// The StringMapI (pronounced StringMappy) interface describes structures that implement the important
// and common map of strings indexed by strings which is a part of most modern languages.
type StringMapI interface {
	SetChanged(key string, val string) (changed bool, err error)
	Get(key string) (val string)
	Has(key string) (exists bool)
	Remove(key string)
	Values()([]string)
	Keys()([]string)
	Len()(int)
	// Range will iterate ove the keys and values in the map. Pattern is taken from sync.Map
	Range(f func(key string, value string) bool)
	Merge(i StringMapI)
}

// StringMap is an implementention of a string map, which is a very common data structure.
// It implements the StringMapI interface. However, it is not concurrency safe. See SafeStringMap
// for a safe version of this
// This unsafe version is useful for initializing a static map to be used to merge into other maps
type StringMap map[string]string

func NewStringMap() StringMap {
	return make(StringMap)
}

func NewStringMapFrom(i StringMapI) StringMap {
	m := NewStringMap()
	m.Merge(i)
	return m
}


// SetChanged sets the key to the value. Returns two values: whether
// something actually changed, and if there was an error.
func (o StringMap) SetChanged(key string, val string) (changed bool, err error) {
	var ok bool
	var oldVal string

	if o == nil {
		panic("StringMap is not initialized.")
	}

	if oldVal,ok = o[key]; !ok || oldVal != val {
		o[key] = val
		changed = true
	}
	return
}

// Set will set the key to the value and return itself for easy chaining
func (o StringMap) Set(key string, val string) StringMap {
	o[key] = val
	return o
}


// Get returns the string based on its key. If it does not exist, an empty string will be returned.
func (o StringMap) Get(key string) (val string) {
	return o[key]
}

func (o StringMap) Has(key string) (exists bool) {
	_, exists =  o[key]
	return
}

func (o StringMap) Remove(key string) {
	delete (o,key)
}

func (o *StringMap) RemoveAll() {
	*o = NewStringMap()
}

// Values returns a slice of the string values
func (o StringMap) Values() []string {
	vals := make ([]string, 0, len(o))

	for _, v := range o {
		vals = append(vals, v)
	}
	return vals
}

// Keys returns a slice of the string keys
func (o StringMap) Keys() []string {
	keys := make ([]string, 0, len(o))

	for k := range o {
		keys = append(keys, k)
	}
	return keys
}



func (o StringMap) Len() int {
	return len (o)
}

// Range will call the given function with every key and value.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o StringMap) Range(f func(key string, value string) bool) {
	for k, v := range o {
		if !f(k, v) {
			break
		}
	}
}

func (o StringMap) Merge(i StringMapI) {
	if i == nil {
		return
	}
	i.Range(func(k,v string) bool {
		o[k] = i.Get(k)
		return true
	})
}

func (o StringMap) Equals(i StringMapI) bool {
	if i == nil {
		return false
	}
	if i.Len() != o.Len() {
		return false
	}
	var ret bool = true

	i.Range(func(k,v string) bool {
		if !o.Has(k) || o[k] != i.Get(k) {
			ret = false	// don't just return because we are in a channel and we want to use up the channel
			return false // stop iterating
		}
		return true
	})

	return ret
}

