package html

import (
	"strings"
)

// Utilities to manage class strings

// AddClass is a utility function that appends the given class name(s) to the end of the given string, if the names are not
// already in the string. Returns the new string, and a value indicating whether it changed or not.
// newClasses can have multiple space separated strings to add multiple classes. The final string returned will have no duplicates.
// Since the order of a class list in html may make a difference, you should take care in the order of the classes you add
// if this matters in your situation.
func AddClass(classes string, newClasses string) (string, bool) {
	var changed bool
	var found bool

	classArray := strings.Fields(classes)
	newClassArray := strings.Fields(newClasses)
	for _, s := range newClassArray {
		found = false
		for _, s2 := range classArray {
			if s2 == s {
				found = true
			}
		}
		if !found {
			classArray = append(classArray, s)
			changed = true
		}
	}
	return strings.Join(classArray, " "), changed
}

// HasClass searches the list of strings for the given class name. testClass can only be a single class.
func HasClass(classes string, testClass string) (found bool) {
	classArray := strings.Fields(classes)
	for _, s := range classArray {
		if s == testClass {
			found = true
			break
		}
	}
	return
}

// Use RemoveClass to remove a class from the list of classes given. You can give it more than one class to remove by
// separating the classes with spaces in the removeClass string.
func RemoveClass(class string, removeClass string) (string, bool) {
	classes := strings.Fields(class)
	removeClasses := strings.Fields(removeClass)
	ret := ""
	var removed, found bool

	for _, s := range classes {
		found = false
		for _, s2 := range removeClasses {
			if s2 == s {
				removed = true
				found = true
			}
		}
		if !found {
			ret = ret + s + " "
		}
	}

	ret = strings.TrimSpace(ret)

	return ret, removed
}

/*
Many CSS frameworks use families of classes, which are built up from a base family name. For example,
Bootstrap uses 'col-lg-6' to represent a table that is 6 units wide on large screens and Foundation
uses 'large-6' to do the same thing. This utility removes classes that start with a particular prefix
to remove whatever sizing class was specified.
Returns the resulting class list, and true if the list actually changed.
*/
func RemoveClassesWithPrefix(class string, prefix string) (string, bool) {
	classes := strings.Fields(class)
	ret := ""
	var removed bool

	for _, s := range classes {
		if strings.HasPrefix(s, prefix) {
			removed = true
		} else {
			ret = ret + s + " "
		}
	}

	ret = strings.TrimSpace(ret)

	return ret, removed
}
