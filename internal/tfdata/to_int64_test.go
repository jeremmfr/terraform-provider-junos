package tfdata_test

import (
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConvAtoi64Value(t *testing.T) {
	t.Parallel()

	type testCase struct {
		inputStr     string
		expectOutput types.Int64
		expectError  bool
	}

	tests := map[string]testCase{
		"Valid": {
			inputStr:     "10",
			expectOutput: types.Int64Value(10),
			expectError:  false,
		},
		"Zero": {
			inputStr:     "0",
			expectOutput: types.Int64Value(0),
			expectError:  false,
		},
		"Invalid": {
			inputStr:     "!0",
			expectOutput: types.Int64Null(),
			expectError:  true,
		},
		"Empty": {
			inputStr:     "",
			expectOutput: types.Int64Null(),
			expectError:  true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			output, err := tfdata.ConvAtoi64Value(test.inputStr)
			if err != nil {
				if !test.expectError {
					t.Errorf("got unexpected error: %s", err)
				}
			}
			if !output.Equal(test.expectOutput) {
				t.Errorf("expected %s, got %s", test.expectOutput.String(), output.String())
			}
		})
	}
}
