/*
The HTML package includes general functions for manipulating html tags, comments and the like.
It includes specific functions for manipulating attributes inside of tags, including various
special attributes like styles, classes, data-* attributes, etc.

Many of the routines return a boolean to indicate whether the data actually changed. This can be used to prevent
needlessly redrawing html after setting values that had no affect on the attribute list.
*/
package html

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	gohtml "html"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const attributeFalse = "**GORADD-FALSE**"

// Attributer is a general purpose interface for objects that return attributes based on information given.
type Attributer interface {
	Attributes(...interface{}) Attributes
}

// An html attribute manager. Use SetAttribute to set specific attribute values, and then convert it to a string
// to get the attributes in a version embeddable in an html tag.
type Attributes map[string]string

// NewAttributes initializes a group of html attributes.
func NewAttributes() Attributes {
	return make(map[string]string)
}

// NewAttributesFrom creates new attributes from the given string map.
func NewAttributesFrom(i interface{}) Attributes {
	a := NewAttributes()
	a.Merge(i)
	return a
}

// Copy returns a copy of the attributes.
func (a Attributes) Copy() Attributes {
	if a == nil {
		return nil
	}
	return NewAttributesFrom(a)
}

func (a Attributes) Len() int {
	if a == nil {
		return 0
	}
	return len(a)
}

func (a Attributes) Has(attr string) bool {
	if a == nil {
		return false
	}
	_,ok := a[attr]
	return ok
}

func (a Attributes) Get(attr string) string {
	return a[attr]
}

func (a Attributes) Delete(attr string) {
	delete(a, attr)
}

// SetChanged sets the value of an attribute. Looks for special attributes like "class" and "style" to do some error checking
// on them. Returns changed if something in the attribute structure changed, which is useful to determine whether to send
// the changed control to the browser.
// Returns err if the given attribute name or value is not valid.
func (a Attributes) SetChanged(name string, v string) (changed bool, err error) {
	if strings.Contains(name, " ") {
		err = errors.New("attribute names cannot contain spaces")
		return
	}

	if v == attributeFalse {
		changed = a.RemoveAttribute(name)
		return
	}

	if name == "style" {
		styles := NewStyle()
		_, err = styles.SetTo(v)
		if err != nil {
			return
		}

		oldStyles := a.StyleMap()

		if !reflect.DeepEqual(oldStyles, styles) { // since maps are not ordered, we must use a special equality test. We can't just compare strings for equality here.
			changed = true
			a["style"] = styles.String()
		}
		return
	}
	if name == "id" {
		return a.SetIDChanged(v)
	}
	if name == "class" {
		changed = a.SetClassChanged(v)
		return
	}
	if strings.HasPrefix(name, "data-") {
		return a.SetDataAttributeChanged(name[5:], v)
	}
	changed = a.set(name, v)
	return
}

// set is a raw set and return true if changed
func (a Attributes) set(k string, v string) bool {
	oldVal, existed := a[k]
	a[k] = v
	return !existed || oldVal != v
}

// Set is similar to SetChanged, but instead returns an attribute pointer so it can be chained. Will panic on errors.
// Use this when you are setting attributes using implicit strings. Set v to an empty string to create a boolean attribute.
func (a Attributes) Set(name string, v string) Attributes {
	if a == nil {
		a = NewAttributes()
	}
	_, err := a.SetChanged(name, v)
	if err != nil {
		panic(err)
	}
	return a
}

// RemoveAttribute removes the named attribute.
// Returns true if the attribute existed.
func (a Attributes) RemoveAttribute(name string) bool {
	if a.Has(name) {
		a.Delete(name)
		return true
	}
	return false
}

// This is a helper to sort the attribute keys so that special attributes
// are returned in front
var attrSpecialSort = map[string] int {
	"id" : 1,
	"class" : 2,
	"style" : 3,
	// keep name and value together
	"name" : 4,
	"value" : 5,
	// keep src and alt together
	"src" : 6,
	"alt" : 7,
	// keep width and height together
	"width" : 8,
	"height" : 9,
}

func (a Attributes) sortedKeys() []string {
	keys := make([]string, len(a), len(a))
	idx := 0
	for k := range a {
		keys[idx] = k
		idx++
	}
	sort.Slice(keys, func(i1,i2 int) bool {
		k1 := keys[i1]
		k2 := keys[i2]
		v1,ok1 := attrSpecialSort[k1]
		v2,ok2 := attrSpecialSort[k2]
		if ok1 {
			if ok2 {
				return k1 < k2
			}
			return true
		} else if ok2 { // and !ok1
			return false
		} else { // !ok1 && !ok2
			return v1 < v2
		}
	})
	return keys
}

// String returns the attributes escaped and encoded, ready to be placed in an html tag
// For consistency, it will use attrSpecialSort to order the keys. Remaining keys will
// be output in random order.
func (a Attributes) String() string {
	var str string

	if a == nil {
		return ""
	}

	for _,k := range a.sortedKeys() {
		v := a[k]
		if v == "" {
			str += k + " "
		} else {
			v = gohtml.EscapeString(v)
			str += fmt.Sprintf("%s=%q ", k, v)
		}
	}

	return strings.TrimSpace(str)
}

func (a Attributes) Range(f func(key string, value string) bool) {
	if a == nil {
		return
	}
	for _,k := range a.sortedKeys() {
		if !f(k, a[k]) {
			break
		}
	}
}


// Override returns a new Attributes structure with the current attributes merged with the given attributes.
// Conflicts are won by the given overrides. Styles will be merged as well.
func (a Attributes) Override(i interface{}) Attributes {
	attr := a.Copy()
	attr.Merge(i)
	return attr
}

// Merge merges the given attributes into the current attributes. Conflicts are won by the passed in map.
// Styles are merged as well, so that if both the passed in map and the current map have a styles attribute, the
// actual style properties will get merged together.
func (a Attributes) Merge(i interface{}) {
	if i == nil {
		return
	}
	if a == nil {
		a = NewAttributes()
	}
	switch m := i.(type) {
	case map[string]string:
		for k,v := range m {
			if k == "style" {
				if v2,ok := a[k]; ok {
					v = MergeStyleStrings(v2, v)
				}
			}
			a[k] = v
		}
	case Attributes:
		for k,v := range m {
			if k == "style" {
				if v2,ok := a[k]; ok {
					v = MergeStyleStrings(v2, v)
				}
			}
			a[k] = v
		}
	case maps.StringMapI:
		m.Range(func(k,v string) bool {
			if k == "style" {
				if v2,ok := a[k]; ok {
					v = MergeStyleStrings(v2, v)
				}
			}
			a[k] = v
			return true
		})
	case string:
		a2 := getAttributesFromTemplate(m)
		// Merge class instead of over-write in this case
		if a2.Has("class") {
			a.AddClass(a2.Class())
			a2.RemoveAttribute("class")
		}
		a.Merge(a2)
	}
}

// SetIDChanged sets the id to the given value and returns true if something changed.
// In other words, if you set the id to the same value that it currently is, it will return false.
// It will return an error if you attempt to set the id to an illegal value.
func (a Attributes) SetIDChanged(i string) (changed bool, err error) {
	if i == "" { // empty attribute is not allowed, so its the same as removal
		changed = a.RemoveAttribute("id")
		return
	}

	if strings.ContainsAny(i, " ") {
		err = errors.New("id attributes cannot contain spaces")
		return
	}

	changed = a.set("id", i)
	return
}

// SetID sets the id attribute to the given value
func (a Attributes) SetID(i string) Attributes {
	_, err := a.SetIDChanged(i)
	if err != nil {
		panic(err)
	}
	return a
}

// ID returns the value of the id attribute.
func (a Attributes) ID() string {
	if a == nil {
		return ""
	}
	return a.Get("id")
}

// SetClass sets the class attribute to the value given.
// If you prefix the value with "+ " the given value will be appended to the end of the current class list.
// If you prefix the value with "- " the given value will be removed from an class list.
// Otherwise the current class value is replaced.
// Returns whether something actually changed or not.
// v can be multiple classes separated by a space
func (a Attributes) SetClassChanged(v string) bool {
	if v == "" { // empty attribute is not allowed, so its the same as removal
		a.RemoveAttribute("class")
	}

	if strings.HasPrefix(v, "+ ") {
		return a.AddClassChanged(v[2:])
	} else if strings.HasPrefix(v, "- ") {
		return a.RemoveClass(v[2:])
	}

	changed := a.set("class", v)
	return changed
}

// SetClass will set the class to the given value, and return the attributes so you can chain calls.
func (a Attributes) SetClass(v string) Attributes {
	a.SetClassChanged(v)
	return a
}

// Use RemoveClass to remove the named class from the list of classes in the class attribute.
func (a Attributes) RemoveClass(v string) bool {
	if a.Has("class") {
		newClass, changed := RemoveAttributeValue(a.Get("class"), v)
		if changed {
			a.set("class", newClass)
		}
		return changed
	}
	return false
}

// Use RemoveClassesWithPrefix to remove classes with the given prefix.
// Many CSS frameworks use families of classes, which are built up from a base family name. For example,
// Bootstrap uses 'col-lg-6' to represent a table that is 6 units wide on large screens and Foundation
// uses 'large-6' to do the same thing. This utility removes classes that start with a particular prefix
// to remove whatever sizing class was specified.
//Returns true if the list actually changed.
func (a Attributes) RemoveClassesWithPrefix(v string) bool {
	if a.Has("class") {
		newClass, changed := RemoveClassesWithPrefix(a.Get("class"), v)
		if changed {
			a.set("class", newClass)
		}
		return changed
	}
	return false
}

// AddAttributeValueChanged adds the given space separated values to the end of the values in the
// given attribute, removing duplicates and returning true if the attribute was changed at all.
// An example of a place to use this is the aria-labelledby attribute, which can take multiple
// space-separated id numbers.
func (a Attributes) AddAttributeValueChanged(attr string, values string) bool {
	if values == "" {
		return false // nothing to add
	}
	if a.Has(attr) {
		newValues, changed := AddAttributeValue(a.Get(attr), values)
		if changed {
			a.set(attr, newValues)
		}
		return changed
	} else {
		a.set(attr, values)
		return true
	}
}

// AddAttributeValue adds space separated values to the end of an attribute value.
// If a value is not present, the value will be added to the end of the value list.
// If a value is present, it will not be added, and the position of the current value in the list will not change.
func (a Attributes) AddAttributeValue(attr string, value string) Attributes {
	a.AddAttributeValueChanged(attr, value)
	return a
}

// AddClassChanged is similar to AddClass, but will return true if the class changed at all.
func (a Attributes) AddClassChanged(v string) bool {
	return a.AddAttributeValueChanged("class", v)
}

// AddClass adds a class or classes. Multiple classes can be separated by spaces.
// If a class is not present, the class will be added to the end of the class list.
// If a class is present, it will not be added, and the position of the current class in the list will not change.
func (a Attributes) AddClass(v string) Attributes {
	a.AddClassChanged(v)
	return a
}

// Return the value of the class attribute.
func (a Attributes) Class() string {
	return a.Get("class")
}

// HasAttributeValue returns true if the given value exists in the space-separated attribute value.
func (a Attributes) HasAttributeValue(attr string, value string) bool {
	var curValue string
	if curValue = a.Get(attr); curValue == "" {
		return false
	}
	f := strings.Fields(curValue)
	for _, s := range f {
		if s == value {
			return true
		}
	}
	return false
}

// ControlHasClass returns true if the given class is in the class list in the class attribute.
func (a Attributes) HasClass(c string) bool {
	return a.HasAttributeValue("class", c)
}

// SetDataAttributeChanged sets the given value as an html "data-*" attribute.
// The named value will be retrievable in javascript by using
//
//	$obj.dataset.valname;
//
// Note: Data name cases are handled specially. data-* attribute names are supposed to be lower kebab case. Javascript
// converts dashed notation to camelCase when converting html attributes into object properties.
// In other words, we give it a camelCase name here, it shows up in the html as
// a kebab-case name, and then you retrieve it using javascript as camelCase again.
//
// For example, if your html looks like this:
//
//	<div id='test1' data-test-case="my test"></div>
//
// You would get that value in javascript by doing:
//	g$('test1').data('testCase');
//
// Conversion to special html data-* name formatting is handled here automatically. So if you SetDataAttribute('testCase') here,
// you can get it using .dataset.testCase in javascript
func (a Attributes) SetDataAttributeChanged(name string, v string) (changed bool, err error) {
	// validate the name
	if strings.ContainsAny(name, " !$") {
		err = errors.New("data attribute names cannot contain spaces or $ or ! chars")
		return
	}
	suffix, err := ToDataAttr(name)
	if err == nil {
		name = "data-" + suffix
		changed = a.set(name, v)
	}
	return
}

// SetDataAttribute sets the given data attribute. Note that data attribute keys must be in camelCase notation and
// connot be hyphenated. camelCase will get converted to kebab-case in html, and converted back to camelCase when
// referring to the data attribute using .data().
func (a Attributes) SetDataAttribute(name string, v string) Attributes {
	_, err := a.SetDataAttributeChanged(name, v)
	if err != nil {
		panic(err)
	}
	return a
}

/*
DataAttribute gets the data-* attribute value that was set previously.
This does NOT call into javascript to return a value that was set on the browser side. You need to use another
mechanism to retrieve that.
*/
func (a Attributes) DataAttribute(name string) string {
	suffix, _ := ToDataAttr(name)
	name = "data-" + suffix
	return a.Get(name)
}

// RemoveDataAttribute removes the named data attribute. Returns true if the data attribute existed.
func (a Attributes) RemoveDataAttribute(name string) bool {
	suffix, _ := ToDataAttr(name)
	name = "data-" + suffix
	return a.RemoveAttribute(name)
}

// HasDataAttribute returns true if the data attribute is set.
func (a Attributes) HasDataAttribute(name string) bool {
	suffix, _ := ToDataAttr(name)
	name = "data-" + suffix
	return a.Has(name)
}

// Returns the css style string, or a blank string if there is none
func (a Attributes) StyleString() string {
	return a.Get("style")
}

// Returns a special Style structure which lets you refer to the styles as a string map
func (a Attributes) StyleMap() Style {
	s := NewStyle()
	s.SetTo(a.StyleString())
	return s
}

// SetStyle sets the given style to the given value. If the value is prefixed with a plus, minus, multiply or divide, and then a space,
// it assumes that a number will follow, and the specified operation will be performed in place on the current value.
// For example, SetStyle ("height", "* 2") will double the height value without changing the unit specifier.
// When referring to a value that can be a length, you can use numeric values. In this case, "0" will be passed unchanged,
// but any other number will automatically get a "px" suffix.
func (a Attributes) SetStyleChanged(name string, v string) (changed bool, err error) {
	s := a.StyleMap()
	changed, err = s.SetChanged(name, v)
	if err == nil {
		a.set("style", s.String())
	}
	return
}

func (a Attributes) SetStyle(name string, v string) Attributes {
	_, err := a.SetStyleChanged(name, v)
	if err != nil {
		panic(err)
	}
	return a
}

// SetStyle merges the given styles with the current styles. The given style wins on collision.
func (a Attributes) SetStyles(s Style) Attributes {
	styles := a.StyleMap()
	styles.Merge(s)
	a.set("style", styles.String())
	return a
}

// SetStylesTo sets the styles using a traditional css style string with colon and semicolon separatators
func (a Attributes) SetStylesTo(s string) Attributes {
	styles := a.StyleMap()
	if _,err := styles.SetTo(s); err != nil {
		return a
	}
	a.set("style", styles.String())
	return a
}

// GetStyle gives you the value of a single style attribute value. If you want all the attributes as a style string, use
// StyleString().
func (a Attributes) GetStyle(name string) string {
	s := a.StyleMap()
	return s.Get(name)
}

// HasStyle returns true if the given style is set to any value, and false if not.
func (a Attributes) HasStyle(name string) bool {
	s := a.StyleMap()
	return s.Has(name)
}

// RemoveStyle removes the style from the style list. Returns true if there was a change.
func (a Attributes) RemoveStyle(name string) (changed bool) {
	s := a.StyleMap()
	if s.Has(name) {
		changed = true
		s.Delete(name)
		a.set("style", s.String())
	}
	return changed
}

// SetDisabled sets the "disabled" attribute to the given value.
func (a Attributes) SetDisabled(d bool) Attributes {
	if d {
		a.Set("disabled", "")
	} else {
		a.RemoveAttribute("disabled")
	}
	return a
}

// IsDisabled returns true if the "disabled" attribute is set to true.
func (a Attributes) IsDisabled() bool {
	return a.Has("disabled")
}

// SetDisplay sets the "display" attribute to the given value.
func (a Attributes) SetDisplay(d string) Attributes {
	a.SetStyle("display", d)
	return a
}

// IsDisplayed returns true if the "display" attribute is not set, or if it is set, if its not set to "none".
func (a Attributes) IsDisplayed() bool {
	return a.GetStyle("display") != "none"
}

// AttributeString is a helper function to convert an interface type to a string that is appropriate for the value
// in the Set function.
func AttributeString(i interface{}) string {
	switch v := i.(type) {
	case fmt.Stringer:
		return v.String()
	case bool:
		if v {
			return "" // boolean true
		} else {
			return attributeFalse // Our special value to indicate to NOT print the attribute at all
		}
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return fmt.Sprint(i)
	}
}

// getAttributesFromTemplate returns Attributes extracted from a string in the form
// of name="value"
func getAttributesFromTemplate(s string)  Attributes {
	pairs := templateMatcher.FindAllString(s, -1)
	if len(pairs) == 0 {
		return nil
	}
	a := NewAttributes()
	for _,pair := range pairs {
		kv := strings.Split(pair, "=")
		val := kv[1][1:len(kv[1])-1] // remove quotes
		a.Set(kv[0], val)
	}
	return a
}

/*
type AttributeCreator map[string]string

func (c AttributeCreator) Create() Attributes {
	return Attributes(c)
}
*/
var templateMatcher *regexp.Regexp
func init() {
	gob.Register(Attributes{})
	templateMatcher = regexp.MustCompile(`\w+=".*?"`)
}