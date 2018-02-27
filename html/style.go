package html

import (
	"strings"
	"fmt"
	"regexp"
	"strconv"
	"errors"
	"math"
	"github.com/spekary/goradd/util/types"
)

const numericMatch = `-?[\d]*(\.[\d]+)?`

// keys for style attributes that take a number that is not a length
var nonLengthNumerics = map[string]bool{
	"volume": true,
	"speech-rate": true,
	"orphans": true,
	"widows": true,
	"pitch-range": true,
	"font-weight": true,
	"z-index": true,
	"counter-increment": true,
	"counter-reset": true,
}

// Style makes it easy to add and manipulate individual properties in a generated style sheet
// Its main use is for generating a style attribute in an HTML tag
// It implements the String interface to get the style properties as an HTML embeddable string
type Style struct {
	types.StringMap
}

// NewStyle initializes an empty Style object
func NewStyle() *Style {
	return &Style{types.NewStringMap()}
}

// SetTo receives a style encoded "style" attribute into the Style structure (e.g. "width: 4px; border: 1px solid black")
func (s *Style) SetTo(text string) (changed bool, err error) {
	s.RemoveAll()
	a := strings.Split(string(text), ";")	// break apart into pairs
	changed = false
	err = nil
	for _, value := range a {
		b:= strings.Split(value, ":")
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


// Set sets the given style to the given value. If the value is prefixed with a plus, minus, multiply or divide, and then a space,
// it assumes that a number will follow, and the specified operation will be performed in place on the current value
// For example, Set ("height", "* 2") will double the height value without changing the unit specifier
// When referring to a value that can be a length, you can use numeric values. In this case, "0" will be passed unchanged,
// but any other number will automatically get a "px" suffix.

func (s Style) SetChanged(n string, v string) (changed bool, err error) {
	if strings.Contains(n, " ") {
		err = errors.New("Attribute names cannot contain spaces.")
		return
	}
	isNumeric, _ :=  regexp.MatchString("^" + numericMatch + "$", v)

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

func (s Style) Set(n string, v string) Style {
	_, err := s.SetChanged(n,v)
	if err != nil {
		panic(err)
	}
	return s
}


func (s Style) Get(name string) string {
	return s.StringMap.Get(name)
}

// Used in the regular expression replacement function below
func opReplacer(op string, v float64) func(string) string {
	return func(cur string) string {
		if cur == "" {return ""} // bug workaround
		//fmt.Println(cur)
		f, err := strconv.ParseFloat(cur, 0)
		if err != nil {
			panic("The number detector is broken on " + cur) // this is coming directly from the regular expression match
		}
		var newVal float64
		switch op {
		case "+":
			newVal = f+v
		case "-":
			newVal = f-v
		case "*":
			newVal = f*v
		case "/":
			newVal = f/v
		default:
			panic("Unexpected operation")
		}

		// floating point operations sometimes are not accurate. This is an attempt to correct epsilons.
		return fmt.Sprintf("%v", roundFloat(newVal, 6))

	}
}

// mathOp applies the given math operation and value to all the numeric values found in the given attribute.
// Bug(r) If the operation is working on a zero, and the result is not a zero, we may get a raw number with no unit. Not a big deal, but result will use default unit of browser, which is not always px
func (s Style) mathOp (attribute string, op string, val string) (changed bool, err error) {
	cur := s.Get(attribute)
	if cur == "" {
		cur = "0"
	}

	f,err := strconv.ParseFloat(val,0)
	if err != nil {return}
	r, err := regexp.Compile(numericMatch)
	if err != nil {return}
	newStr := r.ReplaceAllStringFunc(cur, opReplacer(op, f))
	changed = s.set(attribute, newStr)
	return
}

// RemoveAll resets the style to contain no styles
func (s *Style) RemoveAll() {
	s.StringMap = types.NewStringMap()
}

// Returns the string version of the style attribute, suitable for inclusion in an HTML style tag
func (s Style) String() string {
	var text []byte

	text = s.encode()
	return string(text)
}

// Raw set and return true if changed
func (s Style) set (n string, v string) bool {
	changed,_ := s.StringMap.SetChanged(n, v)
	return changed
}

// code to take out rounding errors when doing length math
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
func (s Style) encode() (text []byte) {
	var items []string

	for key, value := range s.StringMap {
		items = append (items , key + ":" + value)
	}
	text = []byte(strings.Join(items, ";"))
	return text
}