/*
The HTML package includes general functions for manipulating html tags, comments and the like.
It includes specific functions for manipulating attributes inside of tags, including various
special attributes like styles, classes, data-* attributes, etc.

Many of the routines return a boolean to indicate whether the data actually changed. This can be used to prevent
needlessly redrawing html after setting values that had no affect on the attribute list.
*/
package html

import (
	"errors"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	gohtml "html"
	"strconv"
	"strings"
)

const attributeFalse = "**GORADD-FALSE**"

// Attributer is a general purpose interface for objects that return attributes based on information given.
type Attributer interface {
	Attributes(...interface{}) *Attributes
}

// An html attribute manager. Use SetAttribute to set specific attribute values, and then convert it to a string
// to get the attributes in a version embeddable in an html tag.
type Attributes struct {
	maps.StringSliceMap // Use an ordered string map so that each time we draw the attributes, they do not change order
}

// NewAttributes initializes a group of html attributes.
func NewAttributes() *Attributes { // TODO: This should not return a pointer, since all it contains is a pointer???
	return &Attributes{*maps.NewStringSliceMap()}
}

// NewAttributesFrom creates new attributes from the given string map.
func NewAttributesFrom(i maps.StringMapI) *Attributes {
	a := NewAttributes()
	a.Merge(i)
	return a
}

// NewAttributesFromMap returns new attributes base on the given map.
func NewAttributesFromMap(i map[string]string) *Attributes {
	return &Attributes{*maps.NewStringSliceMapFromMap(i)}
}

// Copy returns a copy of the attributes.
func (a *Attributes) Copy() *Attributes {
	if a == nil {
		return nil
	}
	return NewAttributesFrom(a)
}

// SetChanged sets the value of an attribute. Looks for special attributes like "class" and "style" to do some error checking
// on them. Returns changed if something in the attribute structure changed, which is useful to determine whether to send
// the changed control to the browser.
// Returns err if the given attribute name or value is not valid.
func (a *Attributes) SetChanged(name string, v string) (changed bool, err error) {
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
		styles.SetTo(v)

		oldStyles := a.StyleMap()

		if !oldStyles.Equals(styles) { // since maps are not ordered, we must use a special equality test. We can't just compare strings for equality here.
			changed = true
			a.StringSliceMap.Set("style", styles.String())
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
	changed = a.StringSliceMap.SetChanged(name, v)
	return
}

// Set is similar to SetChanged, but instead returns an attribute pointer so it can be chained. Will panic on errors.
// Use this when you are setting attributes using implicit strings. Set v to an empty string to create a boolean attribute.
func (a *Attributes) Set(name string, v string) *Attributes {
	_, err := a.SetChanged(name, v)
	if err != nil {
		panic(err)
	}
	return a
}

// RemoveAttribute removes the named attribute.
// Returns true if the attribute existed.
func (a *Attributes) RemoveAttribute(name string) bool {
	if a.Has(name) {
		a.Delete(name)
		return true
	}
	return false
}

// String returns the attributes escaped and encoded, ready to be placed in an html tag
// For consistency, it will output the following attributes in the following order if it finds them. Remaining tags will
// be output in random order: id, name, class
func (a *Attributes) String() string {
	var id, name, class, styles, others string
	a.Range(func(k, v string) bool {
		var str string

		if v == "" {
			str = k + " "
		} else {
			v = gohtml.EscapeString(v)
			str = fmt.Sprintf("%s=%q ", k, v)
		}

		switch k {
		case "id":
			id = str
		case "name":
			name = str
		case "class":
			class = str
		case "styles":
			styles = str
		default:
			others += str
		}

		return true
	})

	// put the attributes in a somewhat predictable order
	ret := id + name + class + styles + others
	ret = strings.TrimSpace(ret)

	return ret
}

// Override returns a new Attributes structure with the current attributes merged with the given attributes.
// Conflicts are won by the given overrides. Styles will be merged as well.
func (a *Attributes) Override(i maps.StringMapI) *Attributes {
	attr := a.Copy()
	attr.Merge(i)
	return attr
}

// Merge merges the given attributes into the current attributes. Conflicts are won by the passed in map.
// Styles are merged as well, so that if both the passed in map and the current map have a styles attribute, the
// actual style properties will get merged together.
func (a *Attributes) Merge(i maps.StringMapI) {
	curStyles := a.StyleMap()
	newStyles := NewStyle()
	newStyles.SetTo(i.Get("style"))
	curStyles.Merge(newStyles)
	a.StringSliceMap.Merge(i)
	if curStyles.Len() > 0 {
		a.StringSliceMap.Set("style", curStyles.String())
	}
}

// MergeMap merges the given attributes into the current attributes. Conflicts are won by the passed in map.
// Styles are merged as well, so that if both the passed in map and the current map have a styles attribute, the
// actual style properties will get merged together.
func (a *Attributes) MergeMap(m map[string]string) {
	curStyles := a.StyleMap()
	if styles,ok := m["style"]; ok {
		newStyles := NewStyle()
		newStyles.SetTo(styles)
		curStyles.Merge(newStyles)
	}
	for k,v := range m {
		a.StringSliceMap.Set(k, v)
	}

	if curStyles.Len() > 0 {
		a.StringSliceMap.Set("style", curStyles.String())
	}
}



// SetIDChanged sets the id to the given value and returns true if something changed.
// In other words, if you set the id to the same value that it currently is, it will return false.
// It will return an error if you attempt to set the id to an illegal value.
func (a *Attributes) SetIDChanged(i string) (changed bool, err error) {
	if i == "" { // empty attribute is not allowed, so its the same as removal
		changed = a.RemoveAttribute("id")
		return
	}

	if strings.ContainsAny(i, " ") {
		err = errors.New("id attributes cannot contain spaces")
		return
	}

	changed = a.StringSliceMap.SetChanged("id", i)
	return
}

// SetID sets the id attribute to the given value
func (a *Attributes) SetID(i string) *Attributes {
	_, err := a.SetIDChanged(i)
	if err != nil {
		panic(err)
	}
	return a
}

// ID returns the value of the id attribute.
func (a *Attributes) ID() string {
	return a.Get("id")
}

// SetClass sets the class attribute to the value given.
// If you prefix the value with "+ " the given value will be appended to the end of the current class list.
// If you prefix the value with "- " the given value will be removed from an class list.
// Otherwise the current class value is replaced.
// Returns whether something actually changed or not.
// v can be multiple classes separated by a space
func (a *Attributes) SetClassChanged(v string) bool {
	if v == "" { // empty attribute is not allowed, so its the same as removal
		a.RemoveAttribute("class")
	}

	if strings.HasPrefix(v, "+ ") {
		return a.AddClassChanged(v[2:])
	} else if strings.HasPrefix(v, "- ") {
		return a.RemoveClass(v[2:])
	}

	changed := a.StringSliceMap.SetChanged("class", v)
	return changed
}

// SetClass will set the class to the given value, and return the attributes so you can chain calls.
func (a *Attributes) SetClass(v string) *Attributes {
	a.SetClassChanged(v)
	return a
}

// Use RemoveClass to remove the named class from the list of classes in the class attribute.
func (a *Attributes) RemoveClass(v string) bool {
	if a.Has("class") {
		newClass, changed := RemoveAttributeValue(a.Get("class"), v)
		if changed {
			a.StringSliceMap.Set("class", newClass)
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
func (a *Attributes) RemoveClassesWithPrefix(v string) bool {
	if a.Has("class") {
		newClass, changed := RemoveClassesWithPrefix(a.Get("class"), v)
		if changed {
			a.StringSliceMap.Set("class", newClass)
		}
		return changed
	}
	return false
}

// AddAttributeValueChanged adds the given space separated values to the end of the values in the
// given attribute, removing duplicates and returning true if the attribute was changed at all.
// An example of a place to use this is the aria-labelledby attribute, which can take multiple
// space-separated id numbers.
func (a *Attributes) AddAttributeValueChanged(attr string, values string) bool {
	if values == "" {
		return false // nothing to add
	}
	if a.Has(attr) {
		newValues, changed := AddAttributeValue(a.Get(attr), values)
		if changed {
			a.StringSliceMap.Set(attr, newValues)
		}
		return changed
	} else {
		a.StringSliceMap.Set(attr, values)
		return true
	}
}

// AddAttributeValue adds space separated values to the end of an attribute value.
// If a value is not present, the value will be added to the end of the value list.
// If a value is present, it will not be added, and the position of the current value in the list will not change.
func (a *Attributes) AddAttributeValue(attr string, value string) *Attributes {
	a.AddAttributeValueChanged(attr, value)
	return a
}

// AddClassChanged is similar to AddClass, but will return true if the class changed at all.
func (a *Attributes) AddClassChanged(v string) bool {
	return a.AddAttributeValueChanged("class", v)
}

// AddClass adds a class or classes. Multiple classes can be separated by spaces.
// If a class is not present, the class will be added to the end of the class list.
// If a class is present, it will not be added, and the position of the current class in the list will not change.
func (a *Attributes) AddClass(v string) *Attributes {
	a.AddClassChanged(v)
	return a
}

// Return the value of the class attribute.
func (a *Attributes) Class() string {
	return a.Get("class")
}

// HasAttributeValue returns true if the given value exists in the space-separated attribute value.
func (a *Attributes) HasAttributeValue(attr string, value string) bool {
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

// HasClass returns true if the given class is in the class list in the class attribute.
func (a *Attributes) HasClass(c string) bool {
	return a.HasAttributeValue("class", c)
}

//SetDataAttribute sets the given value as an html "data-*" attribute. The named value will be retrievable in jQuery by using
//
//	$obj.data("name");
//
//Note: Data name cases are handled specially. data-* attribute names are supposed to be lower kebab case. Javascript
//converts dashed notation to camelCase when converting html attributes into object properties.
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
// you can get it using .data('testCase') in jQuery
func (a *Attributes) SetDataAttributeChanged(name string, v string) (changed bool, err error) {
	// validate the name
	if strings.ContainsAny(name, " !$") {
		err = errors.New("data attribute names cannot contain spaces or $ or ! chars")
		return
	}
	suffix, err := ToDataAttr(name)
	if err == nil {
		name = "data-" + suffix
		changed = a.StringSliceMap.SetChanged(name, v)
	}
	return
}

// SetDataAttribute sets the given data attribute. Note that data attribute keys must be in camelCase notation and
// connot be hyphenated. camelCase will get converted to kebab-case in html, and converted back to camelCase when
// referring to the data attribute using .data().
func (a *Attributes) SetDataAttribute(name string, v string) *Attributes {
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
func (a *Attributes) DataAttribute(name string) string {
	suffix, _ := ToDataAttr(name)
	name = "data-" + suffix
	return a.Get(name)
}

// RemoveDataAttribute removes the named data attribute. Returns true if the data attribute existed.
func (a *Attributes) RemoveDataAttribute(name string) bool {
	suffix, _ := ToDataAttr(name)
	name = "data-" + suffix
	return a.RemoveAttribute(name)
}

// HasDataAttribute returns true if the data attribute is set.
func (a *Attributes) HasDataAttribute(name string) bool {
	suffix, _ := ToDataAttr(name)
	name = "data-" + suffix
	return a.Has(name)
}

// Returns the css style string, or a blank string if there is none
func (a *Attributes) StyleString() string {
	return a.Get("style")
}

// Returns a special Style structure which lets you refer to the styles as a string map
func (a *Attributes) StyleMap() *Style {
	s := NewStyle()
	s.SetTo(a.StyleString())
	return s
}

// SetStyle sets the given style to the given value. If the value is prefixed with a plus, minus, multiply or divide, and then a space,
// it assumes that a number will follow, and the specified operation will be performed in place on the current value.
// For example, SetStyle ("height", "* 2") will double the height value without changing the unit specifier.
// When referring to a value that can be a length, you can use numeric values. In this case, "0" will be passed unchanged,
// but any other number will automatically get a "px" suffix.
func (a *Attributes) SetStyleChanged(name string, v string) (changed bool, err error) {
	s := a.StyleMap()
	changed, err = s.SetChanged(name, v)
	if err == nil {
		a.StringSliceMap.Set("style", s.String())
	}
	return
}

func (a *Attributes) SetStyle(name string, v string) *Attributes {
	_, err := a.SetStyleChanged(name, v)
	if err != nil {
		panic(err)
	}
	return a
}

// SetStyle merges the given styles with the current styles. The given style wins on collision.
func (a *Attributes) SetStyles(s *Style) *Attributes {
	styles := a.StyleMap()
	styles.Merge(s)
	a.StringSliceMap.Set("style", styles.String())
	return a
}

// SetStylesTo sets the styles using a traditional css style string with colon and semicolon separatators
func (a *Attributes) SetStylesTo(s string) *Attributes {
	styles := a.StyleMap()
	styles.SetTo(s)
	a.StringSliceMap.Set("style", styles.String())
	return a
}

// GetStyle gives you the value of a single style attribute value. If you want all the attributes as a style string, use
// StyleString().
func (a *Attributes) GetStyle(name string) string {
	s := a.StyleMap()
	return s.Get(name)
}

// HasStyle returns true if the given style is set to any value, and false if not.
func (a *Attributes) HasStyle(name string) bool {
	s := a.StyleMap()
	return s.Has(name)
}

// RemoveStyle removes the style from the style list. Returns true if there was a change.
func (a *Attributes) RemoveStyle(name string) (changed bool) {
	s := a.StyleMap()
	if s.Has(name) {
		changed = true
		s.Delete(name)
		a.StringSliceMap.Set("style", s.String())
	}
	return changed
}

// SetDisabled sets the "disabled" attribute to the given value.
func (a *Attributes) SetDisabled(d bool) *Attributes {
	if d {
		a.Set("disabled", "")
	} else {
		a.RemoveAttribute("disabled")
	}
	return a
}

// IsDisabled returns true if the "disabled" attribute is set to true.
func (a *Attributes) IsDisabled() bool {
	return a.Has("disabled")
}

// SetDisplay sets the "display" attribute to the given value.
func (a *Attributes) SetDisplay(d string) *Attributes {
	a.SetStyle("display", d)
	return a
}

// IsDisplayed returns true if the "display" attribute is not set, or if it is set, if its not set to "none".
func (a *Attributes) IsDisplayed() bool {
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
		return fmt.Sprintf("%v", i)
	}
}

type AttributeCreator map[string]string

func (c AttributeCreator) Create() *Attributes {
	return NewAttributesFromMap(c)
}