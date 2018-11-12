package page

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

// These objects and functions are helpers in the page serialization process. Serialization is a big nut to crack in Go,
// with many opinions and options on how to implement. This current implementation tries to be flexible and supportable first, before fast.
// As goradd matures, this can evolve into something that is more optimized. It is essentially implemented as a service that is
// initialized at startup time.
// This is only needed when we are serializing pages. Not needed on single machine implementations that keeps the pagecache in memory.

var pageEncoder PageEncoderI

type PageEncoderI interface {
	NewEncoder(b *bytes.Buffer) Encoder
	NewDecoder(b *bytes.Buffer) Decoder
}

func SetPageEncoder(e PageEncoderI) {
	if pageEncoder != nil {
		panic("Only set the page encoder when the application is initialized, and only once.")
	}
	pageEncoder = e
}

// Encoder defines objects that can be encoded into a pagestate.
type Encoder interface {
	Encode(v interface{}) error
	EncodeControl(v ControlI) error
}

// Decoder defines objects that can be decoded from a pagestate. If the object does not implement this, we will look for GobDecode support.
type Decoder interface {
	Decode(v interface{}) error
	DecodeControl(p *Page) (ControlI,error)
}

// Serializable defines the interface that allows an object to be encodable using a pre-set encoder. This saves time
// on memory allocations/deallocations, which might be extensive.
// Controls are Serializable by default. Other objects that contain controls, or that are not gob.Encoders should implement
// this as well if they are part of the pagestate.
type Serializable interface {
	Serialize(e Encoder) error
	Deserialize(d Decoder) error
}

type GobPageEncoder struct {
}

type GobSerializer struct {
	*gob.Encoder
	//*json.Encoder
}

type GobDeserializer struct {
	*gob.Decoder
	//*json.Decoder
}

func (e GobPageEncoder) NewEncoder(b *bytes.Buffer) Encoder {
	return &GobSerializer{gob.NewEncoder(b)}
}

func (e GobPageEncoder) NewDecoder(b *bytes.Buffer) Decoder {
	return &GobDeserializer{gob.NewDecoder(b)}
}


// Serialize sends
func (e GobSerializer) Encode(v interface{}) (err error) {
	switch v2 := v.(type) {
	case ControlI:
		panic("call EncodeControl instead")
	case *Page:
		return v2.Encode(e)
	case Serializable:
		if err = e.Encoder.Encode(&v2); err !=nil { // essentially encodes an empty object
			return
		}
		return v2.Serialize(e)
	default:
		return e.Encoder.Encode(v2) // use the standard gob encoder
	}
	return nil
}


type ControlLoc int8
const (
	ControlIsNil ControlLoc = iota
	ControlIsHere
	ControlIsInPage
)

func (e GobSerializer) EncodeControl(c ControlI) (err error) {
	vi := reflect.ValueOf(c)

	if vi.IsNil() {
		if err = e.Encoder.Encode(ControlIsNil); err != nil {
			return
		}
	} else if c.control().encoded {
		// send just true and an id, as we will be getting the control from the control store during decode
		if err = e.Encoder.Encode(ControlIsInPage); err != nil {
			return
		}
		if err = e.Encoder.Encode(c.ID()); err != nil {
			return
		}
	} else {
		// send false and encode the control
		if err = e.Encoder.Encode(ControlIsHere); err != nil {
			return
		}

		// This is a bit of a hack. We know that controls implement the GobEncode method, but return nil so that
		// no content is actually encoded. We do this so that we can use the GobSerializer to essentially serialize
		// a completely empty control so that we can fill it using our own Serialize/Deserialize methods, rather than GobEncode
		// which requires a memory allocation for each little object.
		// However, this interferes with being able to encode collections of things that might have a control in them.
		// There is no simple solution to this.
		if err = e.Encoder.Encode(&c); err != nil { // essentially encodes an empty control
			return
		}


		if !c.ΩisSerializer(c) {
			v := vi.Elem()
			// No direct Serialize method, so we will attempt to handle serialization ourselves by serializing the individual members
			count := v.NumField()
			for i := 0; i < count; i++ {
				f := v.Field(i)
				var ci ControlI

				// if f is an embedded control, we have to get the interface to the control as if it was a pointer
				if f.Kind() == reflect.Struct {
					if ci2, ok := f.Addr().Interface().(ControlI); ok {
						ci = ci2
					}
				} else if ci2, ok := f.Interface().(ControlI); ok {
					ci = ci2
				}

				if ci != nil {
					if err = e.EncodeControl(ci); err != nil {
						return
					}
				} else {
					iface := f.Interface()
					if err = e.Encode(&iface); err != nil {
						return
					}
				}
			}
			c.control().encoded = true
		} else {
			if err = c.Serialize(e); err != nil {
				return
			}
		}
	}
	return nil
}


func (e GobDeserializer) Decode(v interface{}) (err error) {
	switch v2 := v.(type) {
	case *ControlI:
		panic ("call DecodeControl instead")
	case ControlI:
		panic ("call DecodeControl instead with the address of a ControlI var")
	case *Page:
		if err = v2.Decode(e); err != nil {
			return
		}

	case *Serializable:
		if err = e.Decoder.Decode(v2); err != nil {
			return
		}
		if err = (*v2).Deserialize(e); err != nil {
			return
		}

	default:
		return e.Decoder.Decode(v)
	}
	return nil
}

// DecodeControl decodes a control from the stream and returns it as a ControlI.
func (e GobDeserializer) DecodeControl(p *Page) (c ControlI, err error) {
	var loc ControlLoc
	if err = e.Decoder.Decode(&loc); err != nil {
		return
	}

	if loc == ControlIsNil {
		return
	} else if loc == ControlIsInPage {
		var id string
		if err = e.Decoder.Decode(&id); err != nil {
			return
		}
		c = p.GetControl(id)
	} else {
		if err = e.Decoder.Decode(&c); err != nil {
			return
		}

		if !c.ΩisSerializer(c) {
			v := reflect.ValueOf(c).Elem()
			count := v.NumField()
			for i := 0; i < count; i++ {
				f := v.Field(i)

				if _, ok := f.Addr().Interface().(ControlI); ok {
					if ci, err := e.DecodeControl(p); err != nil {
						return nil, err
					} else if ci != nil {
						f.Set(reflect.ValueOf(ci).Elem())
					}
				} else if _, ok := f.Interface().(ControlI); ok {
					if ci, err := e.DecodeControl(p); err != nil {
						return nil, err
					} else if ci != nil {
						f.Set(reflect.ValueOf(ci).Convert(f.Type()))
					}
				} else {
					var iface interface{}
					if err = e.Decode(&iface); err != nil {
						return
					}
					f.Set(reflect.ValueOf(iface).Convert(f.Type()))
				}
			}
			c.control().encoded = true
		} else {
			if err = c.Deserialize(e,p); err != nil {
				return
			}
		}

		c.control().page = p
		p.controlRegistry.Set(c.control().ID(), c)
	}
	return
}

