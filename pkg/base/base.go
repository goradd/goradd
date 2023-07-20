// Package base describes the Base struct, which provides virtual functions and inheritance similar to
// that of C++ and Java.
package base

import (
	"fmt"
	"reflect"
)

// Base is a base structure to embed in objects that can call its functions as virtual functions. It provides support
// for a pattern of inheritance that is similar to C++ or Java.
//
// Virtual functions are functions that provide a default implementation but that can be overridden by a subclass.
// For some, the lack of virtual function support in Go is a benefit, as some claim virtual functions contribute to poor architecture.
// However, for certain problems, they can be really helpful. They are especially helpful when you want to model a system
// that has default functionality, but you want to create similar items that have small changes to that functionality.
// Application frameworks are a good example of this.
//
// To create such a structure, use the following pattern:
//
//	 type Bird struct {
//	   Base
//	 }
//
//	 // BirdI defines the virtual functions of Bird.
//	 type BirdI interface {
//	   BaseI
//	   Call() string
//	 }
//
//	 func NewBird() *Bird {
//	   b := new(Bird)
//	   b.Init(b)
//	   return b
//	 }
//
//	 func (a *Bird) PrintCall () {
//		  fmt.Print(a.this().Call())
//	 }
//
//	 // Call can be overridden by a subclass.
//	 func (a *Bird) Call () string {
//	   return "chirp"
//	 }
//
//	 func (a *Bird) this() BirdI {
//	   return a.Self().(BirdI)
//	 }
//
//	 type Duck struct {
//	   Bird
//	 }
//
//	 type DuckI interface {
//	   BirdI
//	 }
//
//	 func NewDuck() *Duck {
//	   d := new(Duck)
//	   d.Init(d)
//	   return d
//	 }
//
//	 func (a *Duck) Call () string {
//	   return "quack"
//	 }
//
// The following code will then print "quack":
//
//	NewDuck().PrintCall()
//
// Be careful when using GobDecode() to restore an encoded object. To restore virtual function ability, you should
// call Init() on the object after it is decoded.
//

type Base struct {
	self any
}

// Init captures the object being passed to it so that virtual functions can be called on it.
func (b *Base) Init(self any) {
	if b.self != nil && b.self != self {
		panic("do not initialize an object twice")
	}
	b.self = self
}

// Self returns the self as an interface object.
func (b *Base) Self() any {
	return b.self
}

// TypeOf returns the reflection Type of self.
func (b *Base) TypeOf() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(b.self)).Type()
}

// String implements the Stringer interface.
//
// This default functionality outputs the object's type. String is overridable.
func (b *Base) String() string {
	return b.TypeOf().String()
}

// BaseI is an interface representing all objects that embed the Base struct.
//
// It adds the ability for all such objects to get the object's type and a string
// representation of the object, which by default is the object's type. To be sure
// to get the object's type even if String() is overridden, use TypeOf().String().
type BaseI interface {
	// TypeOf returns the reflection type of the object. This allows superclasses to determine
	// the type of the subclassing object. The type will not be a pointer to the type, but the type itself.
	TypeOf() reflect.Type
	// String returns a string representation of the object.
	String() string
}

// Embedded returns the structure of type T that is embedded in o.
//
// T should be a pointer type to the embedded structure and be public. The returned structure will bypass
// virtual functions that are defined in o and call functions directly on the embedded structure. It will also
// have access to all the private members of o.
//
// Will panic if T is not found in o.
//
// Another way to do this, and perhaps better in certain circumstances, is to create a private method
// on both the object and interface to the object that simply returns the object.
//
// For example. using the Bird-Duck example:
//
//		type BirdI interface {
//		  BaseI
//		  Call() string
//		  self() *Bird
//		}
//
//	 func (a* Bird) self() *Bird {
//		  return a
//		}
//
// This can be very helpful in special situations, like if you want access to private members from recursive structures.
//
// For example:
//
//	  type Bird struct {
//	    Base
//	    chicks []BirdI
//	    wasFed bool
//		 }
//
// You can do this:
//
//	   func (a *Bird) FeedChicks () string {
//		    for chick := range a.self().chicks {
//		      chick.self().wasFed = true
//	     }
//		  }
//
// And the FeedChicks() function will work on any type of Bird class, without having to expose the internals of Birds,
// or create more private methods on the Bird interface just to get access to the internals of Bird.
func Embedded[T BaseI](o BaseI) T {
	var t T
	typ := reflect.TypeOf(t)
	i := findField(typ, o)
	if i == nil {
		panic(fmt.Errorf("%s does not contain %s", o.TypeOf().String(), typ.String()))
	}
	return i.(T)
}

func findField(typ reflect.Type, o BaseI) any {
	var isPtr bool
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		isPtr = true
	}

	v := reflect.ValueOf(o)
	v = reflect.Indirect(v)
	structFields := reflect.VisibleFields(v.Type())

	for _, structField := range structFields {
		if structField.Anonymous {
			if structField.Type == typ {
				// found it
				v2 := v.FieldByIndex(structField.Index)
				if isPtr {
					v2 = v2.Addr()
				}
				return v2.Interface()
			}
		} else {
			// Anonymous functions are supposed to appear first, so we take advantage of that to exit early.
			break
		}
	}
	return nil // not found
}
