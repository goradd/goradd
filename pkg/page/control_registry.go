package page

import (
	"hash/fnv"
	"reflect"
)

var controlRegistry = make(map[uint64]reflect.Type)
var controlRegistryIds = make(map[reflect.Type]uint64)

// RegisterControl registers the control for the serialize/deserialize process. You should call this
// for each control from an init() function.
//
// As a control is added to the registry, it is assigned an id. That id is used to identify a control
// in the serialization and deserialization process.  We make a significant attempt to prevent the
// addition of controls to an application from causing a change in these ids, since an id change will
// also cause the current page cache to be invalidated. We use a hashing function, and a collision detector
// to do that. If a collision is detected, it will panic, and you should change the hash salt and try again,
// as well as bump the cache version to invalidate the cache.
func RegisterControl(i ControlI) {
	typ := reflect.Indirect(reflect.ValueOf(i)).Type()
	if _, ok := controlRegistryIds[typ]; ok {
		panic("Registering duplicate control")
	}
	hash := fnv.New64()
	n := typ.Name()
	if n == "" {
		panic("type problem")
	}
	_,_ = hash.Write([]byte(ControlRegistrySalt))
	_,_ = hash.Write([]byte(typ.PkgPath()))
	_,_ = hash.Write([]byte(typ.Name()))
	id := hash.Sum64()
	if t,ok := controlRegistry[id]; ok {
		panic("The control registry has detected a collision. " +
			t.Name() + " has collided with " + typ.Name() + ". " +
		"This is a very rare situation, but needs " +
			"to be fixed. To fix it, change the ControlRegistrySalt value, and also change the " +
			"PageCacheVersionID")
	}
	controlRegistry[id] = typ
	controlRegistryIds[typ] = id
}

func controlRegistryID(i ControlI) uint64 {
	val := reflect.Indirect(reflect.ValueOf(i))
	typ := val.Type()
	id, ok := controlRegistryIds[typ]
	if !ok {
		panic("ControlBase type is not registered: " + typ.String())
	}
	return id
}

func createRegisteredControl(registryID uint64, p *Page) ControlI {
	typ := controlRegistry[registryID]
	v := reflect.New(typ)
	c := v.Interface().(ControlI)
	c.control().Self = c
	c.control().page = p
	return c
}

func controlIsRegistered(i interface{}) bool {
	typ := reflect.Indirect(reflect.ValueOf(i)).Type()
	_,ok := controlRegistryIds[typ]
	return ok
}


