package datetime

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {

	d := Now()
	d2 := Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())
	assert.True(t, d.Equal(d2))

	d = d.UTC()
	d2 = Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), nil)
	assert.True(t, d.Equal(d2))
}

func TestZero(t *testing.T) {

	d := NewZeroDate()
	d2 := NewDateTime()

	assert.True(t, d.Equal(d2))
}

func ExampleDateTime_MarshalJSON() {
	d := NewDateTime("2012-11-01T22:08:41+00:00")

	v,_ := d.MarshalJSON()
	fmt.Println(string(v))
	// Output: {"d":1,"goraddObject":"dt","h":22,"m":8,"mo":10,"ms":0,"s":41,"t":true,"y":2012,"z":false}
}

func ExampleDateTime_JavaScript() {
	d := NewDateTime("2012-11-01T22:08:41+00:00")

	v:= d.JavaScript()
	fmt.Println(v)
	// Output: new Date(Date.UTC(2012, 10, 1, 22, 8, 41, 0))
}

func TestTime(t *testing.T) {
	d := NewDateTime("3 04 PM", "8 41 PM")
	d2 := Time(20,41,0,0)
	assert.True(t, d.Equal(d2), d.String(), d2.String())

}