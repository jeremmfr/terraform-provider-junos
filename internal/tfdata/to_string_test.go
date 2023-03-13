package tfdata_test

import (
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestJunosDecode(t *testing.T) {
	t.Parallel()

	type testCase struct {
		inputStr     string
		expectOutput types.String
		expectError  bool
	}

	tests := map[string]testCase{
		"Valid": {
			inputStr:     "$9$1HFIyKXxdsgJ-VH.Pfn6lKMXdsZUi5Qnikfz",
			expectOutput: types.StringValue("testPassWord"),
			expectError:  false,
		},
		"Invalid": {
			inputStr:     "$9aaa",
			expectOutput: types.StringNull(),
			expectError:  true,
		},
		"Empty": {
			inputStr:     "",
			expectOutput: types.StringNull(),
			expectError:  true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			output, err := tfdata.JunosDecode(test.inputStr, "Message")
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
