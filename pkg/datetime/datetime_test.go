package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestStringConversion(t *testing.T) {
	d := Now()

	d2 := NewDateTime(d.String())
	assert.True(t, d.Equal(d2))

	d3 := NewDateTime(d.UTC().String())
	assert.True(t, d.Equal(d3))

}

func TestZero(t *testing.T) {

	d := NewZeroDate()
	d2 := NewDateTime()

	assert.True(t, d.Equal(d2))
}

func ExampleDateTime_MarshalJSON() {
	d := NewDateTime("2012-11-01T22:08:41+00:00")
	v, _ := d.MarshalJSON()
	fmt.Println(string(v))
	// Output: {"d":1,"goraddObject":"date","h":22,"m":8,"mo":10,"ms":0,"s":41,"t":false,"y":2012,"z":false}
}

func ExampleDateTime_JavaScript() {
	d := NewDateTime("2012-11-01T22:08:41+00:00")
	v := d.JavaScript()
	fmt.Println(v)
	// Output: new Date(2012, 10, 1, 22, 8, 41, 0)
}

func TestTimeOnly(t *testing.T) {
	d := NewDateTime("8 41 PM", "3 04 PM")
	d2 := Time(20, 41, 0, 0)
	assert.True(t, d.Equal(d2), d.String(), d2.String())

}

func TestTZ(t *testing.T) {
	// Testing that daylight savings time affects the actual time
	loc, _ := time.LoadLocation("America/New_York")
	d := Date(2012, January, 1, 4, 0, 0, 0, loc)
	d2 := Date(2012, July, 1, 5, 0, 0, 0, loc)
	h1 := d.UTC().Hour()
	h2 := d2.UTC().Hour()
	assert.Equal(t, h1, h2)

	z, _ := d.Zone()
	assert.Equal(t, z, "EST")
	z, _ = d2.Zone()
	assert.Equal(t, z, "EDT")
}

func TestFormattedRead(t *testing.T) {
	d := NewDateTime("19/2/2018", EuroDate)
	assert.Equal(t, February, d.Month())
	assert.Equal(t, 19, d.Day())
	d = NewDateTime("3:04 PM", UsTime)
	assert.Equal(t, 15, d.Hour())
	assert.Equal(t, 4, d.Minute())
	d = NewDateTime("2/19/2019 3:04 pm", UsDateTime)
	assert.Equal(t, February, d.Month())
	assert.Equal(t, 19, d.Day())
	assert.Equal(t, 15, d.Hour())
	assert.Equal(t, 4, d.Minute())
}
