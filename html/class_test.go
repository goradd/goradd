package html

import (
	//"testing"
	"fmt"
	"strconv"
	"testing"
)


func ExampleAddClass() {
	classes, changed := AddClass("myClass1 myClass2", "myClass1 myClass3")
	fmt.Println(classes + ":" + strconv.FormatBool(changed))
	//Output: myClass1 myClass2 myClass3:true
}

func ExampleRemoveClass() {
	classes, changed := RemoveClass("myClass1 myClass2", "myClass1 myClass3")
	fmt.Println(classes + ":" + strconv.FormatBool(changed))
	//Output: myClass2:true
}

func ExampleHasClass() {
	found := HasClass("myClass31 myClass2", "myClass3")
	fmt.Println(strconv.FormatBool(found))
	//Output: false
}

func ExampleRemoveClassesWithPrefix() {
	classes, changed := RemoveClassesWithPrefix("col-6 col-brk col4-other", "col-")
	fmt.Println(classes + ":" + strconv.FormatBool(changed))
	//Output: col4-other:true
}

func TestAddClass(t *testing.T) {
	classes, changed := AddClass("myClass1", "myClass1")
	if classes != "myClass1" {
		t.Errorf("AddClass expected (%q), got (%q)", "myClass1", classes)
	}
	if changed {
		t.Errorf("AddClass expected no change, got change")
	}

	classes, changed = AddClass("", "myClass1")
	if classes != "myClass1" {
		t.Errorf("AddClass expected (%q), got (%q)", "myClass1", classes)
	}
	if !changed {
		t.Errorf("AddClass expected change, got no change")
	}

	classes, changed = AddClass("myClass1", "")
	if classes != "myClass1" {
		t.Errorf("AddClass expected (%q), got (%q)", "myClass1", classes)
	}
	if changed {
		t.Errorf("AddClass expected no change, got change")
	}

	// Removes extra spaces
	classes, changed = AddClass(" myClass1  myClass2", "")
	if classes != "myClass1 myClass2" {
		t.Errorf("AddClass expected (%q), got (%q)", "myClass1 myClass2", classes)
	}
	if changed {
		t.Errorf("AddClass expected no change, got change")
	}

	if !HasClass("a b c", "b") {
		t.Errorf("HasClass failed to find a class")
	}

}
