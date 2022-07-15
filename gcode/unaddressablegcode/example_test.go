package unaddressablegcode_test

import (
	"fmt"

	"github.com/mauroalderete/gcode-cli/gcode/unaddressablegcode"
)

func ExampleNew() {

	gc, err := unaddressablegcode.New('X')
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s", gc)
	// Output: X
}

func ExampleGcode_HasAddress() {

	gc, err := unaddressablegcode.New('X')
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	fmt.Printf("%v\n", gc.HasAddress())

	// Output: false
}

func ExampleGcode_String() {

	gc, err := unaddressablegcode.New('X')
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	fmt.Println(gc)

	// Output: X
}

func ExampleGcode_Word() {

	gc, err := unaddressablegcode.New('D')
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	fmt.Println(gc.Word())

	// Output: 68
}
