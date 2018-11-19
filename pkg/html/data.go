package html

import (
	"errors"
	"regexp"
	"strings"
)

/*
 Helper function to convert a name from camel case to using dashes to separated words.
 data-* html attributes have special conversion rules. Attribute names should always be lower case. Dashes in the
 name get converted to camel case javascript variable names by jQuery.
 For example, if you want to pass the value with key name "testVar" to javascript by printing it in
 the html, you would use this function to help convert it to "data-test-var", after which you can retrieve
 in javascript by calling ".data('testVar')". on the object.
 This will also test for the existence of a camel case string it cannot handle due to a bug in jQuery
*/
func ToDataAttr(s string) (string, error) {
	if matched, _ := regexp.MatchString("^[^a-z]|[A-Z][A-Z]|\\W", s); matched {
		err := errors.New("This is not an acceptable camelCase name.")
		return s, err
	}
	re, err := regexp.Compile("[A-Z]")
	if err == nil {
		s = re.ReplaceAllStringFunc(s, func(s2 string) string { return "-" + strings.ToLower(s2) })
	}

	return strings.TrimSpace(strings.TrimPrefix(s, "-")), err
}

/*
 Helper function to convert a name from data attribute naming convention to camel case.
 data-* html attributes have special conversion rules. Key names should always be lower case. Dashes in the
 name get converted to camel case javascript variable names by jQuery.
 For example, if you want to pass the value with key name "testVar" to javascript by printing it in
 the html, you would use this function to help convert it to "data-test-var", after which you can retrieve
 in jQuery by calling ".data('testVar')". on the object.
*/
func ToDataJqKey(s string) (string, error) {
	if matched, _ := regexp.MatchString("[A-Z]|[^a-z0-9-]", s); matched {
		err := errors.New("This is not an acceptable kabob-case name.")
		return s, err
	}

	pieces := strings.Split(s, "-")
	var ret string
	for i, p := range pieces {
		if len(p) == 1 {
			err := errors.New("Due to a jQuery bug, individual kabob words must be at least 2 characters long.")
			return s, err
		}
		if i != 0 {
			p = strings.Title(p)
		}
		ret += p
	}
	return ret, nil
}
