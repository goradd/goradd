package datetime

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestNow(t *testing.T) {

	d := Now()
	d2 := Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())
	assert.True(t, d.Equal(d2))
	assert.True(t, d2.IsTimestamp())

	d = d.UTC()
	d2 = Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), nil)
	assert.True(t, d.Equal(d2))
	assert.True(t, d.IsTimestamp())
	assert.False(t, d2.IsTimestamp())
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
	// Output: {"d":1,"goraddObject":"dt","h":22,"m":8,"mo":10,"ms":0,"s":41,"t":false,"y":2012,"z":false}
}

func ExampleDateTime_JavaScript() {
	d := NewDateTime("2012-11-01T22:08:41+00:00")
	v:= d.JavaScript()
	fmt.Println(v)
	// Output: new Date(2012, 10, 1, 22, 8, 41, 0)
}

func TestTime(t *testing.T) {
	d := NewDateTime("8 41 PM", "3 04 PM")
	d2 := Time(20,41,0,0)
	assert.True(t, d.Equal(d2), d.String(), d2.String())

}

func TestTZ(t *testing.T) {
	// Testing that daylight savings time affects the actual time
	loc, _ := time.LoadLocation("America/New_York")
	d := Date(2012, January, 1, 4,0,0,0,loc)
	d2 := Date(2012, July, 1, 5,0,0,0,loc)
	h1 := d.UTC().Hour()
	h2 := d2.UTC().Hour()
	assert.Equal(t, h1, h2)

	z,_ := d.Zone()
	assert.Equal(t, z, "EST")
	z,_ = d2.Zone()
	assert.Equal(t, z, "EDT")
}