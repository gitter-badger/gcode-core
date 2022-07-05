package gcode_test

import (
	"fmt"
	"testing"

	"github.com/mauroalderete/gcode-skew-transform-cli/pkg/address"
	"github.com/mauroalderete/gcode-skew-transform-cli/pkg/gcode"
	"github.com/mauroalderete/gcode-skew-transform-cli/pkg/word"
)

func TestNewGcode(t *testing.T) {
	t.Run("valids", func(t *testing.T) {

		cases := []struct {
			word byte
		}{
			{'W'},
			{'X'},
			{'N'},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("(%v)", i), func(t *testing.T) {
				gc, err := gcode.NewGcode(c.word)
				if err != nil {
					t.Errorf("got %v, want X12", err)
					return
				}
				if gc == nil {
					t.Errorf("got nil gcode, want %v", c.word)
					return
				}
				if gc.String() != string(c.word) {
					t.Errorf("got %s, want %v", gc, c.word)
				}
			})
		}
	})

	t.Run("invalids", func(t *testing.T) {

		t.Run("word", func(t *testing.T) {
			cases := []struct {
				word byte
				err  error
			}{
				{'+', &word.WordInvalidValueError{Value: '+'}},
				{'\t', &word.WordInvalidValueError{Value: '\t'}},
				{'"', &word.WordInvalidValueError{Value: '"'}},
			}

			for i, c := range cases {
				t.Run(fmt.Sprintf("(%v)", i), func(t *testing.T) {
					_, err := gcode.NewGcode(c.word)
					if err.Error() != c.err.Error() {
						t.Errorf("got %v, want %v", err, c.err)
					}
				})
			}
		})
	})
}

func TestNewGcodeAddressable(t *testing.T) {
	t.Run("valids", func(t *testing.T) {

		t.Run("address integer", func(t *testing.T) {
			gc, err := gcode.NewGcodeAddressable[int32]('X', 12)
			if err != nil {
				t.Errorf("got %v, want X12", err)
				return
			}
			if gc.String() != "X12" {
				t.Errorf("got %s, want X12", gc)
			}
		})

		t.Run("address float", func(t *testing.T) {
			gc, err := gcode.NewGcodeAddressable[float32]('X', 12.3)
			if err != nil {
				t.Errorf("got %v, want X12.3", err)
				return
			}
			if gc.String() != "X12.3" {
				t.Errorf("got %s, want X12.3", gc)
			}
		})

		t.Run("address string", func(t *testing.T) {
			gc, err := gcode.NewGcodeAddressable('X', "\"lorem ipsu\"")
			if err != nil {
				t.Errorf("got %v, want X\"lorem ipsu\"", err)
				return
			}
			if gc.String() != "X\"lorem ipsu\"" {
				t.Errorf("got %s, want X\"lorem ipsu\"", gc)
			}
		})
	})

	t.Run("invalids", func(t *testing.T) {

		t.Run("word", func(t *testing.T) {
			cases := []struct {
				word    byte
				address int32
				err     error
			}{
				{'+', 2, &word.WordInvalidValueError{Value: '+'}},
				{'\t', 2, &word.WordInvalidValueError{Value: '\t'}},
				{'"', 2, &word.WordInvalidValueError{Value: '"'}},
			}

			for i, c := range cases {
				t.Run(fmt.Sprintf("(%v)", i), func(t *testing.T) {
					_, err := gcode.NewGcodeAddressable(c.word, c.address)
					if err.Error() != c.err.Error() {
						t.Errorf("got %v, want %v", err, c.err)
					}
				})
			}
		})

		t.Run("address string", func(t *testing.T) {
			cases := []struct {
				word    byte
				address string
				err     error
			}{
				{'X', "", &address.AddressStringTooShortError{Value: ""}},
				{'X', "\"\t\"", &address.AddressStringContainInvalidCharsError{Value: "\"\t\""}},
				{'X', "\"\"\"", &address.AddressStringQuoteError{Value: ""}},
			}

			for i, c := range cases {
				t.Run(fmt.Sprintf("(%v)", i), func(t *testing.T) {
					_, err := gcode.NewGcodeAddressable(c.word, c.address)
					if err.Error() != c.err.Error() {
						t.Errorf("got %v, want %v", err, c.err)
					}
				})
			}
		})
	})
}

func ExampleNewGcode() {

	gc, err := gcode.NewGcode('X')
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s", gc)

	// Output: X
}

func ExampleNewGcodeAddressable() {

	gc, err := gcode.NewGcodeAddressable[float32]('X', 12.3)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s", gc)

	// Output: X12.3
}

func ExampleGcode_HasAddress() {

	gc, err := gcode.NewGcode('X')
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	fmt.Printf("%v\n", gc.HasAddress())

	// Output: false
}

func ExampleGcode_HasAddress_second() {

	gca, err := gcode.NewGcodeAddressable('X', "\"Hola mundo!\"")
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	fmt.Printf("%v\n", gca.HasAddress())

	// Output: true
}

func ExampleGcode_HasAddress_third() {

	gca, err := gcode.NewGcodeAddressable('X', "\"Hola mundo!\"")
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	var gc gcode.Gcoder = gca

	fmt.Printf("%v\n", gc.HasAddress())

	// Output: true
}

func ExampleGcode_String() {

	gc, err := gcode.NewGcode('X')
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	fmt.Println(gc)

	// Output: X
}

func ExampleGcode_String_second() {

	gca, err := gcode.NewGcodeAddressable('X', "\"Hola mundo!\"")
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	fmt.Println(gca)

	// Output: X"Hola mundo!"
}

func ExampleGcode_Word() {
	gc, err := gcode.NewGcode('D')
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	var w word.Word = gc.Word()

	fmt.Println(w.String())

	// Output: D
}

func ExampleGcode_Word_second() {
	gca, err := gcode.NewGcodeAddressable[int32]('D', 0)
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	var w word.Word = gca.Word()

	fmt.Println(w.String())

	// Output: D
}

func ExampleGcodeAddressable_Address() {
	gc, err := gcode.NewGcodeAddressable[int32]('N', 66555)
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	fmt.Println(gc.Address().String())

	// Output: 66555
}

func ExampleGcodeAddressable_Address_second() {
	gca, err := gcode.NewGcodeAddressable[int32]('N', 66555)
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
		return
	}

	var gc gcode.Gcoder = gca

	if !gc.HasAddress() {
		fmt.Println("Ups! gcode not contain address")
		return
	}

	if value, ok := gc.(*gcode.GcodeAddressable[int32]); ok {
		fmt.Printf("the int32 address recovered is %v\n", value.Address().String())
	}

	// Output: the int32 address recovered is 66555
}
