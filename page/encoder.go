package page

import (
	"encoding/gob"
	"bytes"
)

// These objects and functions are helpers in the page serialization process. Serialization is a big nut to crack in Go,
// with many opinions and options on how to implement. This current implementation tries to be flexible and supportable first, before fast.
// As goradd matures, this can evolve into something that is more optimized. It is essentially implemented as a service that is
// initialized at startup time.
// This is only needed when we are serializing pages. Not needed on single machine implementations that keep the pagecache in memory.


var pageEncoder PageEncoderI

type PageEncoderI interface {
	NewEncoder(b *bytes.Buffer) Encoder
	NewDecoder (b *bytes.Buffer) Decoder
}

func SetPageEncoder(e PageEncoderI) {
	if pageEncoder != nil {
		panic("Only set the page encoder when the application is initialized, and only once.")
	}
	pageEncoder = e
}


// Encoder defines objects that can be encoded into a pagestate. If the object does not implement this, we will look for MarshalBinary support,
// and finally just encode the exported members.
type Encoder interface {
	Encode(v interface{}) error
}

// Decoder defines objects that can be decoded from a pagestate. If the object does not implement this, we will look for MarshalBinary support,
// and finally just decode using the exported members.
type Decoder interface {
	Decode (v interface{}) error
}

// Encodable defines the interface that all controls and any object inside a control should implement in order to serialize itself to the page cache
type Encodable interface {
	Encode (e Encoder) error
	Decode (d Decoder) error
}

type GobPageEncoder struct {
}

type GobEncoder struct {
	*gob.Encoder
}

type GobDecoder struct {
	*gob.Decoder
}

func (e GobPageEncoder) NewEncoder(b *bytes.Buffer) Encoder {
	return &GobEncoder{gob.NewEncoder(b)}
}

func (e GobPageEncoder) NewDecoder(b *bytes.Buffer) Decoder {
	return &GobDecoder{gob.NewDecoder(b)}
}

func (e GobEncoder) Encode(v interface{}) error {
	switch v2 := v.(type) {
	case Encodable:
		// TODO: Output the type. Type must be registered.
		return v2.Encode(e)
	default:
		return e.Encoder.Encode(v2)
	}
}

func (e GobDecoder) Decode(v interface{}) error {
	switch v2 := v.(type) {
	case Encodable: // assume this is a pointer to an interface
		// TODO: Retrieve the type, and then create a new object with that type, then decode into that variable, then set the v to that. Will need a registry of types to do that, like gob.
		return v2.Decode(e)
	default:
		return e.Decoder.Decode(v)
	}
}