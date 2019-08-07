package html

import (
	"errors"
	"fmt"
	maps2 "github.com/goradd/goradd/pkg/maps"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const numericMatch = `-?[\d]*(\.[\d]+)?`

// keys for style attributes that take a number that is not a length
var nonLengthNumerics = map[string]bool{
	"volume":            true,
	"speech-rate":       true,
	"orphans":           true,
	"widows":            true,
	"pitch-range":       true,
	"font-weight":       true,
	"z-index":           true,
	"counter-increment": true,
	"counter-reset":     true,
}

// Style makes it easy to add and manipulate individual properties in a generated style sheet
// Its main use is for generating a style attribute in an HTML tag
// It implements the String interface to get the style properties as an HTML embeddable string
type Style map[string]string

// NewStyle initializes an empty Style object
func NewStyle() Style {
	return make(map[string]string)
}

func NewStyleFromMap(m map[string]string) Style {
	s := NewStyle()
	for k,v := range m {
		s[k] = v
	}
	return s
}

func (s Style) Merge(m Style) {
	for k,v := range m {
		s[k] = v
	}
}

func (s Style) Len() int {
	if s == nil {
		return 0
	}
	return len(s)
}

func (s Style) Has(prop string) bool {
	if s == nil {
		return false
	}
	_,ok := s[prop]
	return ok
}

func (s Style) Get(prop string) string {
	return s[prop]
}

func (s Style) Delete(prop string) {
	delete(s, prop)
}

// SetTo receives a style encoded "style" attribute into the Style structure (e.g. "width: 4px; border: 1px solid black")
func (s Style) SetTo(text string) (changed bool, err error) {
	s.RemoveAll()
	a := strings.Split(string(text), ";") // break apart into pairs
	changed = false
	err = nil
	for _, value := range a {
		b := strings.Split(value, ":")
		if len(b) != 2 {
			err = errors.New("Css must be a name/value pair separated by a colon. '" + string(text) + "' was given.")
			return
		}
		newChange, newErr := s.SetChanged(strings.TrimSpace(b[0]), strings.TrimSpace(b[1]))
		if newErr != nil {
			err = newErr
			return
		}
		changed = changed || newChange
	}
	return
}

// SetChanged sets the given style to the given value. If the value is prefixed with a plus, minus, multiply or divide, and then a space,
// it assumes that a number will follow, and the specified operation will be performed in place on the current value
// For example, Set ("height", "* 2") will double the height value without changing the unit specifier
// When referring to a value that can be a length, you can use numeric values. In this case, "0" will be passed unchanged,
// but any other number will automatically get a "px" suffix.

func (s Style) SetChanged(n string, v string) (changed bool, err error) {
	if strings.Contains(n, " ") {
		err = errors.New("attribute names cannot contain spaces")
		return
	}
	isNumeric, _ := regexp.MatchString("^"+numericMatch+"$", v)

	if strings.HasPrefix(v, "+ ") ||
		strings.HasPrefix(v, "- ") || // the space here distinguishes between a math operation and a negative value
		strings.HasPrefix(v, "* ") ||
		strings.HasPrefix(v, "/ ") {

		return s.mathOp(n, v[0:1], v[2:])
	}

	if v == "0" {
		changed = s.set(n, v)
		return
	}

	if isNumeric {
		if !nonLengthNumerics[n] {
			v = v + "px"
		}
		changed = s.set(n, v)
		return
	}

	changed = s.set(n, v)
	return
}

// Set is like SetChanged, but returns the Style for chaining.
// It will also allocate a style if passed a nil style, and return it
func (s Style) Set(n string, v string) Style {
	if s == nil {
		s = NewStyle()
	}
	_, err := s.SetChanged(n, v)
	if err != nil {
		panic(err)
	}
	return s
}

// opReplacer is used in the regular expression replacement function below
func opReplacer(op string, v float64) func(string) string {
	return func(cur string) string {
		if cur == "" {
			return ""
		} // bug workaround
		//fmt.Println(cur)
		f, err := strconv.ParseFloat(cur, 0)
		if err != nil {
			panic("The number detector is broken on " + cur) // this is coming directly from the regular expression match
		}
		var newVal float64
		switch op {
		case "+":
			newVal = f + v
		case "-":
			newVal = f - v
		case "*":
			newVal = f * v
		case "/":
			newVal = f / v
		default:
			panic("Unexpected operation")
		}

		// floating point operations sometimes are not accurate. This is an attempt to correct epsilons.
		return fmt.Sprintf("%v", roundFloat(newVal, 6))

	}
}

// mathOp applies the given math operation and value to all the numeric values found in the given attribute.
// Bug(r) If the operation is working on a zero, and the result is not a zero, we may get a raw number with no unit. Not a big deal, but result will use default unit of browser, which is not always px
func (s Style) mathOp(attribute string, op string, val string) (changed bool, err error) {
	cur := s.Get(attribute)
	if cur == "" {
		cur = "0"
	}

	f, err := strconv.ParseFloat(val, 0)
	if err != nil {
		return
	}
	r, err := regexp.Compile(numericMatch)
	if err != nil {
		return
	}
	newStr := r.ReplaceAllStringFunc(cur, opReplacer(op, f))
	changed = s.set(attribute, newStr)
	return
}

// RemoveAll resets the style to contain no styles
func (s Style) RemoveAll() {
	for k := range s {
		delete(s, k)
	}
}

// String returns the string version of the style attribute, suitable for inclusion in an HTML style tag
func (s Style) String() string {
	return s.encode()
}

// set is a raw set and return true if changed
func (s Style) set(k string, v string) bool {
	oldVal, existed := s[k]
	s[k] = v
	return !existed || oldVal != v
}

// roundFloat takes out rounding errors when doing length math
func roundFloat(f float64, digits int) float64 {
	f = f * math.Pow10(digits)
	if math.Abs(f) < 0.5 {
		return 0
	}
	v := int(f + math.Copysign(0.5, f))
	f = float64(v) / math.Pow10(digits)
	return f
}

// encode will output a text version of the style, suitable for inclusion in an html "style" attribute.
// it will sort the keys so that they are presented in a consistent and testable way.
func (s Style) encode() (text string) {
	keys := maps2.SortedKeys(s)

	for i, k := range keys {
		if i > 0 {
			text += ";"
		}
		text += k + ":" + s.Get(k)
	}
	return text
}

// StyleString converts an interface type that is being used to set a style value to a string that can be fed into
// the SetStyle* functions
func StyleString(i interface{}) string {
	var sValue string
	switch v := i.(type) {
	case int:
		sValue = fmt.Sprintf("%dpx", v)
	case float32:
		sValue = fmt.Sprintf("%fpx", v)
	case float64:
		sValue = fmt.Sprintf("%fpx", v)
	case string:
		sValue = v
	case fmt.Stringer:
		sValue = v.String()
	default:
		sValue = fmt.Sprintf("%v", v)
	}
	return sValue
}

type StyleCreator map[string]string

func (c StyleCreator) Create() Style {
	return NewStyleFromMap(c)
}