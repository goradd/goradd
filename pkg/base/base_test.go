package base

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Bird struct {
	Base
}

// BirdI defines the virtual functions of Bird.
type BirdI interface {
	BaseI
	Call() string
}

func NewBird() *Bird {
	b := new(Bird)
	b.Init(b)
	return b
}

func (a *Bird) GetCall() string {
	return a.this().Call()
}

// Call can be overridden by a subclass.
func (a *Bird) Call() string {
	return "chirp"
}

func (a *Bird) this() BirdI {
	return a.Self().(BirdI)
}

type Duck struct {
	Bird
}

type DuckI interface {
	BirdI
}

func NewDuck() *Duck {
	d := new(Duck)
	d.Init(d)
	return d
}

func (a *Duck) Call() string {
	return "quack"
}

func TestBase(t *testing.T) {
	assert.Equal(t, "chirp", NewBird().GetCall())
	assert.Equal(t, "quack", NewDuck().GetCall())
}

func TestBase_String(t *testing.T) {
	assert.Equal(t, "base.Bird", NewBird().String())
	assert.Equal(t, "base.Duck", NewDuck().String())
}

func TestBase_Init(t *testing.T) {
	assert.Panics(t, func() {
		d := NewDuck()
		d.Init(NewBird())
	})
	assert.NotPanics(t, func() {
		d := NewDuck()
		d.Init(d)
	})

}

type errStruct struct {
	Base
}

func TestEmbedded(t *testing.T) {
	d := NewDuck()
	b := Embedded[*Bird](d)
	s := b.Call()
	assert.Equal(t, "chirp", s)
	assert.Panics(t, func() {
		Embedded[*errStruct](d)
	})
}
