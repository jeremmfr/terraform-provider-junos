package tfvalidator_test

import (
	"context"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringRuneExclusion(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val          types.String
		excludedRune rune
		expectError  bool
	}
	tests := map[string]testCase{
		"unknown": {
			val:          types.StringUnknown(),
			excludedRune: '/',
			expectError:  false,
		},
		"null": {
			val:          types.StringNull(),
			excludedRune: '/',
			expectError:  false,
		},
		"valid": {
			val:          types.StringValue("ok"),
			excludedRune: '/',
			expectError:  false,
		},
		"invalid": {
			val:          types.StringValue("ok/ok"),
			excludedRune: '/',
			expectError:  true,
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
			tfvalidator.StringRuneExclusion(test.excludedRune).ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestStringDoubleQuoteExclusion(t *testing.T) {
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
		"valid": {
			val:         types.StringValue("ok"),
			expectError: false,
		},
		"invalid": {
			val:         types.StringValue("\"ok\""),
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
			tfvalidator.StringDoubleQuoteExclusion().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestStringSpaceExclusion(t *testing.T) {
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
		"valid": {
			val:         types.StringValue("ok"),
			expectError: false,
		},
		"invalid": {
			val:         types.StringValue("k o"),
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
			tfvalidator.StringSpaceExclusion().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestStringDotExclusion(t *testing.T) {
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
		"valid": {
			val:         types.StringValue("ok"),
			expectError: false,
		},
		"invalid": {
			val:         types.StringValue("k.o"),
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
			tfvalidator.StringDotExclusion().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
