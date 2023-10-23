package tfvalidator_test

import (
	"context"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringNumberRange(t *testing.T) {
	t.Parallel()

	nameRange := "Test"
	type testCase struct {
		val         types.String
		min         int
		max         int
		name        *string
		expectError bool
	}
	tests := map[string]testCase{
		"unknown": {
			val:         types.StringUnknown(),
			min:         10,
			max:         20,
			expectError: false,
		},
		"null": {
			val:         types.StringNull(),
			min:         10,
			max:         20,
			expectError: false,
		},
		"valid_digit": {
			val:         types.StringValue("15"),
			min:         10,
			max:         20,
			expectError: false,
		},
		"valid_range": {
			val:         types.StringValue("15-18"),
			min:         10,
			max:         20,
			expectError: false,
		},
		"valid_range_limit": {
			val:         types.StringValue("10-20"),
			min:         10,
			max:         20,
			expectError: false,
		},
		"invalid_digit": {
			val:         types.StringValue("0"),
			min:         10,
			max:         20,
			expectError: true,
		},
		"invalid_range_min": {
			val:         types.StringValue("9-18"),
			min:         10,
			max:         20,
			expectError: true,
		},
		"invalid_range_max": {
			val:         types.StringValue("10-21"),
			min:         10,
			max:         20,
			expectError: true,
		},
		"invalid_range_inv": {
			val:         types.StringValue("15-12"),
			min:         10,
			max:         20,
			expectError: true,
		},
		"invalid_range_name": {
			val:         types.StringValue("15d"),
			min:         10,
			max:         20,
			expectError: true,
			name:        &nameRange,
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
			validator := tfvalidator.StringNumberRange(test.min, test.max)
			if test.name != nil {
				validator = validator.WithNameInError(*test.name)
			}
			validator.ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
