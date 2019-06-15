package page

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var tests1 = []struct {
	data        []byte
	stringValue string
	intValue    int
	floatValue  float64
	boolValue   bool
}{
	{[]byte(nil), "", 0, 0.0, false},
	{[]byte(""), "", 0, 0.0, false},
	{[]byte(`"a"`), "a", 0, 0.0, true},
	{[]byte("a"), "a", 0, 0.0, true},
	{[]byte("0"), "0", 0, 0.0, false},
	{[]byte("1"), "1", 1, 1.0, true},
	{[]byte("15"), "15", 15, 15.0, true},
	{[]byte("true"), "true", 0, 0.0, true},
	{[]byte(`"true"`), "true", 0, 0.0, true},
	{[]byte("false"), "false", 0, 0.0, false},
	{[]byte("-1"), "-1", -1, -1.0, true},
	{[]byte(".1"), ".1", 0, 0.1, true},
	{[]byte("-4.1"), "-4.1", -4, -4.1, true},
	{[]byte(`["a"]`), `["a"]`, 0, 0, true},
	{[]byte(`{"a":"b"}`), `{"a":"b"}`, 0, 0, true},
}

func TestValues(t *testing.T) {
	var a ActionParams

	for _, test := range tests1 {
		a.values.Event = test.data
		assert.EqualValues(t, test.stringValue, a.EventValueString())
		assert.EqualValues(t, test.intValue, a.EventValueInt())
		assert.EqualValues(t, test.floatValue, a.EventValueFloat())
		assert.EqualValues(t, test.boolValue, a.EventValueBool())
	}

	for _, test := range tests1 {
		a.values.Action = test.data
		assert.EqualValues(t, test.stringValue, a.ActionValueString())
		assert.EqualValues(t, test.intValue, a.ActionValueInt())
		assert.EqualValues(t, test.floatValue, a.ActionValueFloat())
		assert.EqualValues(t, test.boolValue, a.ActionValueBool())
	}

	for _, test := range tests1 {
		a.values.Control = test.data
		assert.EqualValues(t, test.stringValue, a.ControlValueString())
		assert.EqualValues(t, test.intValue, a.ControlValueInt())
		assert.EqualValues(t, test.floatValue, a.ControlValueFloat())
		assert.EqualValues(t, test.boolValue, a.ControlValueBool())
	}
}
