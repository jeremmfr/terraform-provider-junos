package tfvalidator_test

import (
	"context"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBoolTrue(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Bool
		expectError bool
	}
	tests := map[string]testCase{
		"unknown": {
			val:         types.BoolUnknown(),
			expectError: false,
		},
		"null": {
			val:         types.BoolNull(),
			expectError: false,
		},
		"valid": {
			val:         types.BoolValue(true),
			expectError: false,
		},
		"invalid": {
			val:         types.BoolValue(false),
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.BoolRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.BoolResponse{}
			tfvalidator.BoolTrue().ValidateBool(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
