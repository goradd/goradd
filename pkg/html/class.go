package html

import (
	"strings"
)

// Utilities to manage class strings

// AddAttributeValue is a utility function that appends the given space separated words to the end
// of the given string, if the words are not already in the string. This is primarily used for
// adding classes to a class attribute, but other attributes use this structure as well, like
// aria-labelledby and aria-describedby attributes.
//
// Returns the new string, and a value indicating whether it changed or not.
// The final string returned will have no duplicates.
// Since the order of a class list in html makes a difference, you should take care in the
// order of the classes you add if this matters in your situation.
func AddAttributeValue(originalValues string, newValues string) (string, bool) {
	var changed bool
	var found bool

	wordArray := strings.Fields(originalValues)
	newWordArray := strings.Fields(newValues)
	for _, s := range newWordArray {
		found = false
		for _, s2 := range wordArray {
			if s2 == s {
				found = true
			}
		}
		if !found {
			wordArray = append(wordArray, s)
			changed = true
		}
	}
	return strings.Join(wordArray, " "), changed
}

// HasWord searches the list of strings for the given word.
func HasWord(words string, testWord string) (found bool) {
	classArray := strings.Fields(words)
	for _, s := range classArray {
		if s == testWord {
			found = true
			break
		}
	}
	return
}

// Use RemoveAttributeValue to remove a value from the list of space-separated values given.
// You can give it more than one value to remove by
// separating the values with spaces in the removeValue string. This is particularly useful
// for removing a class from a class list in a class attribute.
func RemoveAttributeValue(originalValues string, removeValue string) (string, bool) {
	classes := strings.Fields(originalValues)
	removeClasses := strings.Fields(removeValue)
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
RemoveClassesWithPrefix will remove all classes from the class string with the given prefix.
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
