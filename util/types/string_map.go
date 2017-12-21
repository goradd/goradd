package types

// The StringMapI (pronounced StringMappy) interface describes structures that implement the important
// and common map of strings indexed by strings which is a part of most modern languages.
type StringMapI interface {
	Set(key string, val string) (changed bool, err error)
	Get(key string) (val string)
	Has(key string) (exists bool)
	Remove(key string)
	Values()([]string)
	Keys()([]string)
	Len()(int)
	// Iter returns a channel to iterate over that can be used in a for loop
	Iter() <-chan string
	IterKeys() <-chan string
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


// Set the key to the value. Non-string values will be converted to strings
// using fmt's default conversion for the type. Returns two values: whether
// something actually changed, and if there was an error.
func (o StringMap) Set(key string, val string) (changed bool, err error) {
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

// Iter can be used with range to iterate over the string values. Order of values is not
// guaranteed. See OrderedStringMap for an ordered version.
// Note that we return a buffered channel the size of the return values so there is no blocking
func (o StringMap) Iter() <-chan string {
	c := make(chan string, o.Len())

	f := func() {
		for _, v := range o {
			c <- v
		}
		close(c)
	}
	go f()

	return c
}

// IterKeys can be used with range to iterate over the string keys. Order of values is not
// guaranteed. See OrderedStringMap for an ordered version.
// Note that we return a buffered channel the size of the return values so there is no blocking
func (o StringMap) IterKeys() <-chan string {
	c := make(chan string, o.Len())

	f := func() {
		for k := range o {
			c <- k
		}
		close(c)
	}
	go f()

	return c
}

func (o StringMap) Merge(i StringMapI) {
	if i == nil {
		return
	}
	for k := range i.IterKeys() {
		o[k] = i.Get(k)
	}
}

func (o StringMap) Equals(i StringMapI) bool {
	if i == nil {
		return false
	}
	if i.Len() != o.Len() {
		return false
	}
	var ret bool = true

	for k := range i.IterKeys() {
		if ret && (!o.Has(k) || o[k] != i.Get(k)) {
			ret = false	// don't just return because we are in a channel and we want to use up the channel
		}
	}
	return ret
}