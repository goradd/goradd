/*
The HTML package includes general functions for manipulating html tags, comments and the like.
It includes specific functions for manipulating attributes inside of tags, including various
special attributes like styles, classes, data-* attributes, etc.

Many of the routines return a boolean to indicate whether the data actually changed. This can be used to prevent
needlessly redrawing html after setting values that had no affect on the attribute list.
*/
package html

import (
	"fmt"
	"strings"
	"errors"
	gohtml "html"
	"github.com/spekary/goradd/util/types"
	"strconv"
)

const attributeFalse = "**GORADD-FALSE**"

// Attributer is a general purpose interface for objects that return attributes based on information given.
type Attributer interface {
	Attributes(...interface{}) *Attributes
}

// An html attribute manager. Use SetAttribute to set specific attribute values, and then convert it to a string
// to get the attributes in a version embeddable in an html tag.
type Attributes struct {
	types.OrderedStringMap // Use an ordered string map so that each time we draw the attributes, they do not change order
}

// NewAttributes initializes a group of html attributes.
func NewAttributes() *Attributes {	// TODO: This should not return a pointer, since all it contains is a pointer
    return &Attributes{*types.NewOrderedStringMap()}
}

func NewAttributesFrom(i types.StringMapI) *Attributes {
	a := NewAttributes()
	a.Merge(i)
	return a
}

// SetChanged sets the value of an attribute. Looks for special attributes like "class" and "style" to do some error checking
// on them. Returns changed if something in the attribute structure changed, which is useful to determine whether to send
// the changed control to the browser.
// Returns err if the given attribute name or value is not valid.
func (m *Attributes) SetChanged(name string, v string) (changed bool, err error) {
	if strings.Contains(name, " ") {
		err = errors.New("Attribute names cannot contain spaces.")
		return
	}

	if v == attributeFalse {
		changed = m.RemoveAttribute(name)
		return
	}

	if name == "style" {
		styles := NewStyle()
		styles.SetTo(v)

		oldStyles := m.StyleMap()

		if !oldStyles.Equals(styles) {	// since maps are not ordered, we must use a special equality test. We can't just compare strings for equality here.
			changed = true
			_, err = m.OrderedStringMap.SetChanged("style", styles.String())
		}
		return
	}
	if name == "id" {
		return m.SetIDChanged(v)
	}
	if name == "class" {
		changed = m.SetClassChanged(v)
		return
	}
	if strings.HasPrefix(name, "data-") {
		return m.SetDataAttributeChanged(name[5:], v)
	}
	changed, err = m.OrderedStringMap.SetChanged(name,v)
	return
}

// Set is similar to SetChanged, but instead returns an attribute pointer so it can be chained. Will panic on errors.
// Use this when you are setting attributes using implicit strings. Set v to an empty string to create a boolean attribute.
func (m *Attributes) Set(name string, v string) *Attributes {
	_,err := m.SetChanged(name, v)
	if err != nil {
		panic(err)
	}
	return m
}

// RemoveAttribute removes the named attribute.
// Returns true if the attribute existed.
func (m *Attributes) RemoveAttribute(a string) bool {
	if m.Has(a) {
		m.Remove(a)
		return true
	}
	return false
}

// String returns the attributes escaped and encoded, ready to be placed in an html tag
// For consistency, it will output the following attributes in the following order if it finds them. Remaining tags will
// be output in random order: id, name, class
func (m *Attributes) String() string {
	var id, name, class, styles, others string
	m.Range(func (k,v string) bool {
		var str string

		if v == "" {
			str = k + " "
		} else {
			v = gohtml.EscapeString(v)
			str = fmt.Sprintf("%s=%q ", k, v)
		}

		switch k {
		case "id":id = str
		case "name":name = str
		case "class":class = str
		case "styles":styles = str
		default:others += str
		}

		return true
	})

	// put the attributes in a somewhat predictable order
	ret := id + name + class + styles + others
	ret = strings.TrimSpace(ret)

	return ret
}

// Override returns a new Attributes structure with the current attributes merged with the given attributes.
// Conflicts are won by the given overrides
func (m *Attributes) Override(i types.StringMapI) *Attributes {
	curStyles := m.StyleMap()
	newStyles := NewStyle()
	newStyles.SetTo(i.Get("style"))
	a := NewAttributesFrom(m)
	a.Merge(i)
	curStyles.Merge(newStyles)
	if curStyles.Len() > 0 {
		a.OrderedStringMap.Set("style", curStyles.String())
	}
	return a
}

// Clone returns a copy of the attributes
func (m *Attributes) Clone() *Attributes {
	return NewAttributesFrom(m)
}


// Set the id to the given value. Returns true if something changed.
func (m *Attributes) SetIDChanged(i string) (changed bool, err error) {
	if i == "" {	// empty attribute is not allowed, so its the same as removal
		changed = m.RemoveAttribute("id")
		return
	}

	if strings.ContainsAny(i, " ") {
		err = errors.New("id attributes cannot contain spaces")
		return
	}

	changed, err = m.OrderedStringMap.SetChanged("id", i)
	return
}

func (m *Attributes) SetID(i string) *Attributes {
	_,err := m.SetIDChanged(i)
	if err != nil {
		panic(err)
	}
	return m
}


// Return the value of the id attribute.
func (m *Attributes) ID() string {
	return m.Get("id")
}

// SetClass sets the class attribute to the value given.
// If you prefix the value with "+ " the given value will be appended to the end of the current class list.
// If you prefix the value with "- " the given value will be removed from an class list.
// Otherwise the current class value is replaced.
// Returns whether something actually changed or not.
// v can be multiple classes separated by a space
func (m *Attributes) SetClassChanged(v string) bool {
	if v == "" {	// empty attribute is not allowed, so its the same as removal
		m.RemoveAttribute("class")
	}

	if strings.HasPrefix(v, "+ ") {
		return m.AddClassChanged(v[2:])
	} else if strings.HasPrefix(v, "- ") {
		return m.RemoveClass(v[2:])
	}

	changed,_ := m.OrderedStringMap.SetChanged("class", v)
	return changed
}

func (m* Attributes) SetClass(v string) *Attributes {
	m.SetClassChanged(v)
	return m
}

// Use RemoveClass to remove the named class from the list of classes in the class attribute.
func (m *Attributes) RemoveClass(v string) bool {
	if m.Has("class") {
		newClass,changed := RemoveClass(m.Get("class"), v)
		if changed {
			m.OrderedStringMap.Set("class", newClass)
		}
		return changed
	}
	return false
}

// Use AddClass to add a class or classes.
// Multiple classes can be separated by spaces.
// If a class is not present, the class will be added to the end of the class list
// If a class is present, it will not be added, and the position of the current class in the list will not change
func (m *Attributes) AddClassChanged(v string) bool {
	if m.Has("class") {
		newClass,changed := AddClass(m.Get("class"), v)
		if changed {
			m.OrderedStringMap.Set("class", newClass)
		}
		return changed
	} else {
		m.OrderedStringMap.Set("class", v)
		return true
	}
}

func (m *Attributes) AddClass(v string) *Attributes {
	m.AddClassChanged(v)
	return m
}


// Return the value of the class attribute.
func (m *Attributes) Class() string {
	return m.Get("class")
}

// HasClass return true if the given class is in the class list in the class attribute.
func (m *Attributes) HasClass(c string) bool {
	var curClass string
	if curClass = m.Get("class"); curClass == "" {
		return false
	}
	f := strings.Fields(curClass)
	for _,s := range f {
		if s == c {return true}
	}
	return false
}

/*
SetDataAttribute sets the given value as an html "data-*" attribute. The named value will be retrievable in jQuery by using

	$obj.data("name");

Note: Data name cases are handled specially in jQuery. data-* attribute names are supposed to be online lower case. jQuery
converts dashed notation to camelCase. In other words, we give it a camelCase name here, it shows up in the html as
a dashed name, and then you retrieve it using javascript as camelCase again.

For example, if your html looks like this:

	<div id='test1' data-test-case="my test"></div>

You would get that value in jQuery by doing:
	$j('#test1').data('testCase');

Conversion to special html data-* name formatting is handled here automatically. So if you SetDataAttribute('testCase') here,
you can get it using .data('testCase') in jQuery
*/
func (m *Attributes) SetDataAttributeChanged(name string, v string) (changed bool, err error) {
	// validate the name
	if strings.ContainsAny(name, " !$") {
		err = errors.New("Data attribute names cannot contain spaces or $ or ! chars")
		return
	}
	suffix,err := ToDataAttr(name)
	if err == nil {
		name = "data-" + suffix
		changed, err = m.OrderedStringMap.SetChanged(name, v)
	}
	return
}


// SetDataAttribute sets the given data attribute. Note that data attribute keys must be in camelCase notation and
// connot be hyphenated. camelCase will get converted to kebab-case in html, and converted back to camelCase when
// referring to the data attribute using jQuery.data.
func (m *Attributes) SetDataAttribute(name string, v string) *Attributes {
	_, err := m.SetDataAttributeChanged(name, v)
	if err != nil {
		panic (err)
	}
	return m
}


/*
DataAttribute gets the data-* attribute value that was set previously.
Does NOT call into javascript to return a value that was set on the browser side. You need to use another
mechanism to retrieve that.
*/
func (m *Attributes) DataAttribute(name string) string {
	suffix,_ := ToDataAttr(name)
	name = "data-" + suffix
	return m.Get(name)
}

// RemoveDataAttribute removes the named data attribute. Returns true if the data attribute existed.
func (m *Attributes) RemoveDataAttribute(name string) bool {
	suffix,_ := ToDataAttr(name)
	name = "data-" + suffix
	return m.RemoveAttribute(name)
}

// HasDataAttribute returns true if the data attribute is set.
func (m *Attributes) HasDataAttribute(name string) bool {
	suffix,_ := ToDataAttr(name)
	name = "data-" + suffix
	return m.Has(name)
}

// Returns the css style string, or a blank string if there is none
func (m *Attributes) StyleString() string {
	return m.Get("style")
}

// Returns a special Style structure which lets you refer to the styles as a string map
func (m *Attributes) StyleMap() *Style {
	s := NewStyle()
	s.SetTo(m.StyleString())
	return s
}

// SetStyle sets the given style to the given value. If the value is prefixed with a plus, minus, multiply or divide, and then a space,
// it assumes that a number will follow, and the specified operation will be performed in place on the current value.
// For example, SetStyle ("height", "* 2") will double the height value without changing the unit specifier.
// When referring to a value that can be a length, you can use numeric values. In this case, "0" will be passed unchanged,
// but any other number will automatically get a "px" suffix.
func (m *Attributes) SetStyleChanged(name string, v string) (changed bool, err error) {
	s := m.StyleMap()
	changed, err = s.SetChanged(name, v)
	if err == nil {
		m.OrderedStringMap.Set("style", s.String())
	}
	return
}

func (m *Attributes) SetStyle(name string, v string) *Attributes {
	_,err := m.SetStyleChanged(name, v)
	if err != nil {
		panic (err)
	}
	return m
}

// Style gives you the value of a single style attribute value. If you want all the attributes as a style string, use
// Attribute("style").
func (m *Attributes) GetStyle(name string) string {
	s := m.StyleMap()
	return s.Get(name)
}

func (m *Attributes) HasStyle (name string) bool {
	s := m.StyleMap()
	return s.Has(name)
}

// RemoveStyle removes the style from the style list. Returns true if there was a changed.
func (m *Attributes) RemoveStyle (name string) (changed bool) {
	s := m.StyleMap()
	if s.Has(name) {
		changed = true
		s.Remove(name)
		m.OrderedStringMap.Set("style", s.String())
	}
	return changed
}


func (m *Attributes) SetDisabled(d bool) {
	if d {
		m.Set("disabled", "")
	} else {
		m.RemoveAttribute("disabled")
	}
}

func (m *Attributes) IsDisabled() bool {
	return m.Has("disabled")
}

func (m *Attributes) SetDisplay(d string) {
	m.SetStyle("display", d)
}

func (m *Attributes) IsDisplayed() bool {
	return m.GetStyle("display") != "none"
}

// AttributeString is a helper function to convert an interface type to a string that is appropriate for the value
// in the Set function.
func AttributeString(i interface{}) string {
	switch v := i.(type) {
	case fmt.Stringer:
		return v.String()
	case bool:
		if v {
			return ""	// boolean true
		} else {
			return attributeFalse // Our special value to indicate to NOT print the attribute at all
		}
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return fmt.Sprintf("%v", i)
	}
}