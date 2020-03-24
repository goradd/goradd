package html

import (
	//"testing"
	"fmt"
	"strconv"
	"testing"
)

func ExampleAddClass() {
	classes, changed := AddAttributeValue("myClass1 myClass2", "myClass1 myClass3")
	fmt.Println(classes + ":" + strconv.FormatBool(changed))
	//Output: myClass1 myClass2 myClass3:true
}

func ExampleRemoveClass() {
	classes, changed := RemoveAttributeValue("myClass1 myClass2", "myClass1 myClass3")
	fmt.Println(classes + ":" + strconv.FormatBool(changed))
	//Output: myClass2:true
}

func ExampleHasClass() {
	found := HasWord("myClass31 myClass2", "myClass3")
	fmt.Println(strconv.FormatBool(found))
	//Output: false
}

func ExampleRemoveClassesWithPrefix() {
	classes, changed := RemoveClassesWithPrefix("col-6 col-brk col4-other", "col-")
	fmt.Println(classes + ":" + strconv.FormatBool(changed))
	//Output: col4-other:true
}

func TestAddClass(t *testing.T) {
	classes, changed := AddAttributeValue("myClass1", "myClass1")
	if classes != "myClass1" {
		t.Errorf("AddAttributeValue expected (%q), got (%q)", "myClass1", classes)
	}
	if changed {
		t.Errorf("AddAttributeValue expected no change, got change")
	}

	classes, changed = AddAttributeValue("", "myClass1")
	if classes != "myClass1" {
		t.Errorf("AddAttributeValue expected (%q), got (%q)", "myClass1", classes)
	}
	if !changed {
		t.Errorf("AddAttributeValue expected change, got no change")
	}

	classes, changed = AddAttributeValue("myClass1", "")
	if classes != "myClass1" {
		t.Errorf("AddAttributeValue expected (%q), got (%q)", "myClass1", classes)
	}
	if changed {
		t.Errorf("AddAttributeValue expected no change, got change")
	}

	// Removes extra spaces
	classes, changed = AddAttributeValue(" myClass1  myClass2", "")
	if classes != "myClass1 myClass2" {
		t.Errorf("AddAttributeValue expected (%q), got (%q)", "myClass1 myClass2", classes)
	}
	if changed {
		t.Errorf("AddAttributeValue expected no change, got change")
	}

	if !HasWord("a b c", "b") {
		t.Errorf("ControlHasClass failed to find a class")
	}

}
