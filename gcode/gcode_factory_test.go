package gcode_test

import (
	"fmt"
	"testing"

	"github.com/mauroalderete/gcode-cli/gcode"
	"github.com/mauroalderete/gcode-cli/gcode/addressablegcode"
	"github.com/mauroalderete/gcode-cli/gcode/unaddressablegcode"
)

type GcodeFactory struct{}

// NewGcode is the constructor to instance a Gcode struct.
//
// Receive a word that represents the letter of the command of a gcode.
//
// If the word is an unknown symbol it returns nil with an error description.
func (g *GcodeFactory) NewGcode(word byte) (gcode.Gcoder, error) {
	ng, err := unaddressablegcode.New(word)
	if err != nil {
		return nil, err
	}

	return ng, nil
}

func (g *GcodeFactory) NewAddressableGcodeUint32(word byte, address uint32) (gcode.AddresableGcoder[uint32], error) {

	ng, err := addressablegcode.New(word, address)
	if err != nil {
		return nil, err
	}

	return ng, nil
}

func (g *GcodeFactory) NewAddressableGcodeInt32(word byte, address int32) (gcode.AddresableGcoder[int32], error) {

	ng, err := addressablegcode.New(word, address)
	if err != nil {
		return nil, err
	}

	return ng, nil
}

func (g *GcodeFactory) NewAddressableGcodeFloat32(word byte, address float32) (gcode.AddresableGcoder[float32], error) {

	ng, err := addressablegcode.New(word, address)
	if err != nil {
		return nil, err
	}

	return ng, nil
}

func (g *GcodeFactory) NewAddressableGcodeString(word byte, address string) (gcode.AddresableGcoder[string], error) {

	ng, err := addressablegcode.New(word, address)
	if err != nil {
		return nil, err
	}

	return ng, nil
}

func TestGcodeFactoryNewGcode(t *testing.T) {
	cases := map[string]struct {
		input byte
		valid bool
	}{
		"eval_W":   {'W', true},
		"eval_X":   {'X', true},
		"eval_N":   {'N', true},
		"eval_+":   {'+', false},
		"eval_\\t": {'\t', false},
		"eval_\"":  {'"', false},
	}

	gcodeFactory := &GcodeFactory{}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gc, err := gcodeFactory.NewGcode(tc.input)

			if tc.valid {
				if err != nil {
					t.Errorf("got error %v, want error nil", err)
					return
				}
				if gc == nil {
					t.Errorf("got gcode nil, want gcode %s", string(tc.input))
					return
				}
				if gc.String() != string(tc.input) {
					t.Errorf("got gcode %s, want gcode %s", gc, string(tc.input))
				}
			} else {
				if err == nil {
					t.Errorf("got error %v, want error nil", err)
				}
				if gc != nil {
					t.Errorf("got gcode %s, want gcode nil", gc.String())
				}
			}
		})
	}
}

func TestGcodeAddressableFactoryNewGcodeAddressable(t *testing.T) {
	gcodeFactory := &GcodeFactory{}

	t.Run("addressable gcode uint32", func(t *testing.T) {
		cases := map[string]struct {
			word    byte
			address uint32
			valid   bool
		}{
			"eval_W0":   {'W', 0, true},
			"eval_X1":   {'X', 1, true},
			"eval_N2":   {'N', 2, true},
			"eval_+3":   {'+', 3, false},
			"eval_\\t4": {'\t', 4, false},
			"eval_\"5":  {'"', 5, false},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				gc, err := gcodeFactory.NewAddressableGcodeUint32(tc.word, tc.address)
				if tc.valid {
					if err != nil {
						t.Errorf("got error %v, want error nil", err)
						return
					}
					if gc.String() != fmt.Sprintf("%s%d", string(tc.word), tc.address) {
						t.Errorf("got gcode %s, want gcode %s%d", gc, string(tc.word), tc.address)
					}
				} else {
					if err == nil {
						t.Errorf("got error nil, want error not nil")
						return
					}
				}
			})
		}
	})

	t.Run("addressable gcode int32", func(t *testing.T) {
		cases := map[string]struct {
			word    byte
			address int32
			valid   bool
		}{
			"eval_W0":   {'W', -1, true},
			"eval_X1":   {'X', 0, true},
			"eval_N2":   {'N', 1, true},
			"eval_+3":   {'+', 3, false},
			"eval_\\t4": {'\t', 4, false},
			"eval_\"5":  {'"', 5, false},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				gc, err := gcodeFactory.NewAddressableGcodeInt32(tc.word, tc.address)
				if tc.valid {
					if err != nil {
						t.Errorf("got error %v, want error nil", err)
						return
					}
					if gc.String() != fmt.Sprintf("%s%d", string(tc.word), tc.address) {
						t.Errorf("got gcode %s, want gcode %s%d", gc, string(tc.word), tc.address)
					}
				} else {
					if err == nil {
						t.Errorf("got error nil, want error not nil")
						return
					}
				}
			})
		}
	})

	t.Run("addressable gcode float32", func(t *testing.T) {
		cases := map[string]struct {
			word    byte
			address float32
			valid   bool
		}{
			"eval_W0":     {'W', 0, true},
			"eval_X1.1":   {'X', 1.1, true},
			"eval_N2.2":   {'N', 2.2, true},
			"eval_+3.3":   {'+', 3.3, false},
			"eval_\\t4.4": {'\t', 4.4, false},
			"eval_\"5.5":  {'"', 5.5, false},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				gc, err := gcodeFactory.NewAddressableGcodeFloat32(tc.word, tc.address)
				if tc.valid {
					if err != nil {
						t.Errorf("got error %v, want error nil", err)
						return
					}
					if gc.String() != fmt.Sprintf("%s%.1f", string(tc.word), tc.address) {
						t.Errorf("got gcode %s, want gcode %s%.1f", gc, string(tc.word), tc.address)
					}
				} else {
					if err == nil {
						t.Errorf("got error nil, want error not nil")
						return
					}
				}
			})
		}
	})

	t.Run("addressable gcode string", func(t *testing.T) {
		cases := map[string]struct {
			word    byte
			address string
			valid   bool
		}{
			"eval_W":   {'W', "\"Hola mundo\"", true},
			"eval_X":   {'X', "\"Hola \"\"mundo\"\"\"", true},
			"eval_N":   {'N', "\"Hola mundo\"", true},
			"eval_+":   {'+', "\"Hola mundo\"", false},
			"eval_\\t": {'\t', "\"Hola mundo\"", false},
			"eval_\"":  {'"', "\"Hola mundo\"", false},
			"eval_W2":  {'W', "Hola mundo\"", false},
			"eval_X2":  {'X', "\"Hola \"mundo\"", false},
			"eval_N2":  {'N', "\"Hola mundo", false},
			"eval_W3":  {'W', "\"\tHola mundo\"", false},
			"eval_X3":  {'X', "?", false},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				gc, err := gcodeFactory.NewAddressableGcodeString(tc.word, tc.address)
				if tc.valid {
					if err != nil {
						t.Errorf("got error %v, want error nil", err)
						return
					}
					if gc.String() != fmt.Sprintf("%s%s", string(tc.word), tc.address) {
						t.Errorf("got gcode %s, want gcode %s%s", gc, string(tc.word), tc.address)
					}
				} else {
					if err == nil {
						t.Errorf("got error nil, want error not nil")
						return
					}
				}
			})
		}
	})
}
