package tfvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringRuneCount(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		runeType    runeType
		runeNumber  int
		expectError bool
	}
	tests := map[string]testCase{
		"unknown": {
			val:         types.StringUnknown(),
			runeType:    DotRune,
			runeNumber:  1,
			expectError: false,
		},
		"null": {
			val:         types.StringNull(),
			runeType:    DotRune,
			runeNumber:  1,
			expectError: false,
		},
		"DotRune_valid": {
			val:         types.StringValue("ok.ok"),
			runeType:    DotRune,
			runeNumber:  1,
			expectError: false,
		},
		"DotRune_invalid": {
			val:         types.StringValue("not ok"),
			runeType:    DotRune,
			runeNumber:  1,
			expectError: true,
		},
		"Two_DotRune_valid": {
			val:         types.StringValue("ok.ok.ok"),
			runeType:    DotRune,
			runeNumber:  2,
			expectError: false,
		},
		"Two_DotRune_invalid": {
			val:         types.StringValue("not.ok"),
			runeType:    DotRune,
			runeNumber:  2,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.StringResponse{}

			StringRuneCount(test.runeType, test.runeNumber).ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestString1DotCount(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		expectError bool
	}
	tests := map[string]testCase{
		"unknown": {
			val:         types.StringUnknown(),
			expectError: false,
		},
		"null": {
			val:         types.StringNull(),
			expectError: false,
		},
		"DotRune_valid": {
			val:         types.StringValue("ok.ok"),
			expectError: false,
		},
		"DotRune_invalid": {
			val:         types.StringValue("not ok"),
			expectError: true,
		},
		"Two_DotRune_invalid": {
			val:         types.StringValue("ok.ok.ok"),
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.StringResponse{}

			String1DotCount().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
