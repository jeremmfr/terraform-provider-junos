package tfdata_test

import (
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
)

func TestFirstElementOfJunosLine(t *testing.T) {
	t.Parallel()

	type testCase struct {
		inputStr     string
		expectOutput string
	}

	tests := map[string]testCase{
		"Simple": {
			inputStr:     `foo bar`,
			expectOutput: `foo`,
		},
		"With double quote": {
			inputStr:     `"foo" bar`,
			expectOutput: `"foo"`,
		},
		"With double quote and space": {
			inputStr:     `"foo baz" bar`,
			expectOutput: `"foo baz"`,
		},
		"With double quote and multiple spaces": {
			inputStr:     `" foo baz " bar`,
			expectOutput: `" foo baz "`,
		},
		"With double quote and space and other word with double quote": {
			inputStr:     `"foo baz" "bar qux"`,
			expectOutput: `"foo baz"`,
		},
		"One word": {
			inputStr:     `foo`,
			expectOutput: `foo`,
		},
		"One word with double quote": {
			inputStr:     `"foo"`,
			expectOutput: `"foo"`,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			output := tfdata.FirstElementOfJunosLine(test.inputStr)

			if output != test.expectOutput {
				t.Errorf("expected %s, got %s", test.expectOutput, output)
			}
		})
	}
}
