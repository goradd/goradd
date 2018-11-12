package main

import (
	"fmt"
)

type B struct {
	b int
}

type A struct {
	B
	a int
}

type Doer interface {
	Do()
}

func (d *B) Do() {

}

func main() {
	myA := &A{}
	var d Doer

	d := myA
	d = &myA.B



	fmt.Println("Hello, playground")
}

