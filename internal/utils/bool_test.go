package utils_test

import (
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/utils"
)

func TestParseTrue(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         string
		expectValue bool
	}
	tests := map[string]testCase{
		"unknown": {
			val:         "unknown",
			expectValue: false,
		},
		"empty": {
			val:         "",
			expectValue: false,
		},
		"1": {
			val:         "1",
			expectValue: true,
		},
		"T": {
			val:         "T",
			expectValue: true,
		},
		"true": {
			val:         "true",
			expectValue: true,
		},
		"0": {
			val:         "0",
			expectValue: false,
		},
		"F": {
			val:         "F",
			expectValue: false,
		},
		"false": {
			val:         "false",
			expectValue: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := utils.ParseTrue(test.val)

			if r != test.expectValue {
				t.Fatalf("got unexpected value: want %t, got %t", test.expectValue, r)
			}
		})
	}
}
