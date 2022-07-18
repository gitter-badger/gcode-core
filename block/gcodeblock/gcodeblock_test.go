package gcodeblock

import (
	"fmt"
	"hash"
	"testing"

	"github.com/mauroalderete/gcode-cli/block"
	"github.com/mauroalderete/gcode-cli/block/internal/gcodefactory"
	"github.com/mauroalderete/gcode-cli/checksum"
	"github.com/mauroalderete/gcode-cli/gcode"
	"github.com/mauroalderete/gcode-cli/gcode/addressablegcode"
)

func TestNew(t *testing.T) {

	// N4 G92 E0*67 ;comentario
	mockLineNumber, err := addressablegcode.New[uint32]('N', 4)
	if err != nil {
		t.Errorf("got error not nil, want error nil: %v", err)
	}

	mockCommand, err := addressablegcode.New[int32]('G', 92)
	if err != nil {
		t.Errorf("got error not nil, want error nil: %v", err)
	}

	mockParam1, err := addressablegcode.New[int32]('E', 0)
	if err != nil {
		t.Errorf("got error not nil, want error nil: %v", err)
	}

	mockParameters := []gcode.Gcoder{mockParam1}

	mockChecksum, err := addressablegcode.New[uint32]('*', 67)
	if err != nil {
		t.Errorf("got error not nil, want error nil: %v", err)
	}

	mockComment := ";comentario"

	mockHash := checksum.New()

	mockGcodeFactory := &gcodefactory.GcodeFactory{}

	cases := map[string]struct {
		lineNumber         gcode.AddressableGcoder[uint32]
		command            gcode.Gcoder
		parameters         []gcode.Gcoder
		checksum           gcode.AddressableGcoder[uint32]
		comment            string
		hash               hash.Hash
		gcodeFactory       gcode.GcoderFactory
		configLineNumber   bool
		configParameters   bool
		configChecksum     bool
		configComment      bool
		configHash         bool
		configGcodeFactory bool
		valid              bool
		output             string
	}{
		"Single Word command": {
			command: mockCommand,
			output:  "G92",
			valid:   true,
		},
		"lineNumber nil": {
			command:          mockCommand,
			configLineNumber: true,
			valid:            false,
			output:           "",
		},
		"parameters nil": {
			command:          mockCommand,
			configParameters: true,
			valid:            false,
			output:           "",
		},
		"checksum nil": {
			command:        mockCommand,
			configChecksum: true,
			valid:          false,
			output:         "",
		},
		"hash nil": {
			command:    mockCommand,
			configHash: true,
			valid:      false,
			output:     "",
		},
		"gcodeFactory nil": {
			command:            mockCommand,
			configGcodeFactory: true,
			valid:              false,
			output:             "",
		},
		"+lineNumber": {
			command:          mockCommand,
			configLineNumber: true,
			lineNumber:       mockLineNumber,
			valid:            true,
			output:           "N4 G92",
		},
		"+linenumber+parameters": {
			command:          mockCommand,
			configLineNumber: true,
			lineNumber:       mockLineNumber,
			configParameters: true,
			parameters:       mockParameters,
			valid:            true,
			output:           "N4 G92 E0",
		},
		"+linenumber+parameters+checksum": {
			command:          mockCommand,
			configLineNumber: true,
			lineNumber:       mockLineNumber,
			configParameters: true,
			parameters:       mockParameters,
			configChecksum:   true,
			checksum:         mockChecksum,
			valid:            true,
			output:           "N4 G92 E0*67",
		},
		"+linenumber+parameters+checksum+comment": {
			command:          mockCommand,
			configLineNumber: true,
			lineNumber:       mockLineNumber,
			configParameters: true,
			parameters:       mockParameters,
			configChecksum:   true,
			checksum:         mockChecksum,
			configComment:    true,
			comment:          mockComment,
			valid:            true,
			output:           "N4 G92 E0*67 ;comentario",
		},
		"+linenumber+parameters+checksum+comment+hash": {
			command:          mockCommand,
			configLineNumber: true,
			lineNumber:       mockLineNumber,
			configParameters: true,
			parameters:       mockParameters,
			configChecksum:   true,
			checksum:         mockChecksum,
			configComment:    true,
			comment:          mockComment,
			configHash:       true,
			hash:             mockHash,
			valid:            true,
			output:           "N4 G92 E0*67 ;comentario",
		},
		"+linenumber+parameters+checksum+comment+hash+gcodeFactory": {
			command:            mockCommand,
			configLineNumber:   true,
			lineNumber:         mockLineNumber,
			configParameters:   true,
			parameters:         mockParameters,
			configChecksum:     true,
			checksum:           mockChecksum,
			configComment:      true,
			comment:            mockComment,
			configHash:         true,
			hash:               mockHash,
			configGcodeFactory: true,
			gcodeFactory:       mockGcodeFactory,
			valid:              true,
			output:             "N4 G92 E0*67 ;comentario",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gb, err := New(tc.command, func(config block.BlockConstructorConfigurer) error {

				if tc.configLineNumber {
					err := config.SetLineNumber(tc.lineNumber)
					if err != nil {
						return err
					}
				}

				if tc.configParameters {
					err := config.SetParameters(tc.parameters)
					if err != nil {
						return err
					}
				}

				if tc.configChecksum {
					err := config.SetChecksum(tc.checksum)
					if err != nil {
						return err
					}
				}

				if tc.configComment {
					err := config.SetComment(tc.comment)
					if err != nil {
						return err
					}
				}

				if tc.configGcodeFactory {
					err := config.SetGcodeFactory(tc.gcodeFactory)
					if err != nil {
						return err
					}
				}

				if tc.configHash {
					err := config.SetHash(tc.hash)
					if err != nil {
						return err
					}
				}

				return nil
			})

			if tc.valid {
				if err != nil {
					t.Errorf("got error %v, want error nil", err)
				}

				if gb == nil {
					t.Errorf("got gcodeBlock nil, want gcodeBlock not nil")
					return
				}

				if gb.ToLine("%l %c %p%k %m") != tc.output {
					t.Errorf("got gcodeBlock (%d)[%s], want gcodeBlock: (%d)[%s]", len(gb.ToLine("%l %c %p%k %m")), gb.ToLine("%l %c %p%k %m"), len(tc.output), tc.output)
				}
			} else {
				if err == nil {
					t.Errorf("got error nil, want error not nil")
				}

				if gb != nil {
					t.Errorf("got gcodeBlock not nil, want gcodeBlock nil")
				}
			}
		})
	}
}

func TestGcodeblogk_Parameters(t *testing.T) {

	var cases = [1]struct {
		source     string
		parameters []string
	}{
		{
			source:     "N7 G1 X2.0 Y2.0 F3000.0",
			parameters: []string{"X2.0", "Y2.0", "F3000.0"},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("(%v)", i), func(t *testing.T) {
			b, err := Parse(tc.source)
			if err != nil {
				t.Errorf("got %v, want nil error", err)
				return
			}
			if b == nil {
				t.Errorf("got nil block, want %v", tc.source)
				return
			}

			if len(b.Parameters()) != len(tc.parameters) {
				t.Errorf("got parameters size %d, want paramters size %d", len(b.Parameters()), len(tc.parameters))
				return
			}

			match := true
			for i, s := range tc.parameters {
				if b.Parameters()[i].String() != s {
					match = false
					break
				}
			}

			if !match {
				t.Errorf("got %v, want %v", b.ToLine("%p"), tc.parameters)
			}
		})
	}

}

func TestGcodeblogk_Calculate(t *testing.T) {

	var cases = [1]struct {
		source        string
		checksumValue uint32
	}{
		{
			source:        "N7 G1 X2.0 Y2.0 F3000.0",
			checksumValue: 85,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("(%v)", i), func(t *testing.T) {
			b, err := Parse(tc.source)
			if err != nil {
				t.Errorf("got %v, want nil error", err)
				return
			}
			if b == nil {
				t.Errorf("got nil block, want %v", tc.source)
				return
			}

			gc, err := b.CalculateChecksum()
			if err != nil {
				t.Errorf("got error %v, want error nil", err)
				return
			}

			if gc.Address() != tc.checksumValue {
				t.Errorf("got checksum value %d, want checksum value %d", gc.Address(), tc.checksumValue)
			}
		})
	}
}

func TestGcodeblogk_Verify(t *testing.T) {

	var cases = [1]struct {
		source string
		err    bool
		ok     bool
	}{
		{
			source: "N7 G1 X2.0 Y2.0 F3000.0",
			err:    true,
			ok:     false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("(%v)", i), func(t *testing.T) {
			b, err := Parse(tc.source)
			if err != nil {
				t.Errorf("got %v, want nil error", err)
				return
			}
			if b == nil {
				t.Errorf("got nil block, want %v", tc.source)
				return
			}

			ok, err := b.VerifyChecksum()
			if ok != tc.ok || (err != nil) != tc.err {
				t.Errorf("got error %v verified %v, want error %v verified %v", err, ok, tc.err, tc.ok)
			}
		})
	}
}
