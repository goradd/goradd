package page

import (
	"bytes"
	"encoding/gob"
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
	pageEncoder = e
}

// Encoder defines objects that can be encoded into a pagestate.
type Encoder interface {
	Encode(v interface{}) error
}

// Decoder defines objects that can be decoded from a pagestate. If the object does not implement this, we will look for GobDecode support.
type Decoder interface {
	Decode(v interface{}) error
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

func (e GobSerializer) Encode(v interface{}) (err error) {
	return e.Encoder.Encode(v)
}

func (e GobDeserializer) Decode(v interface{}) (err error) {
	return e.Decoder.Decode(v)
}

