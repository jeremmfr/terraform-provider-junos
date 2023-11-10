package tfvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringFormat(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		format      stringFormat
		expectError bool
		sensitive   bool
	}
	tests := map[string]testCase{
		"unknown": {
			val:         types.StringUnknown(),
			format:      DefaultFormat,
			expectError: false,
		},
		"null": {
			val:         types.StringNull(),
			format:      DefaultFormat,
			expectError: false,
		},
		"DefaultFormat_valid": {
			val:         types.StringValue("ok"),
			format:      DefaultFormat,
			expectError: false,
		},
		"DefaultFormat_invalid": {
			val:         types.StringValue("not ok"),
			format:      DefaultFormat,
			expectError: true,
		},
		"AddressNameFormat_valid": {
			val:         types.StringValue("ok/ok"),
			format:      AddressNameFormat,
			expectError: false,
		},
		"AddressNameFormat_invalid": {
			val:         types.StringValue("not ok"),
			format:      AddressNameFormat,
			expectError: true,
		},
		"DNSNameFormat_valid": {
			val:         types.StringValue("ok.ok"),
			format:      DNSNameFormat,
			expectError: false,
		},
		"DNSNameFormat_invalid": {
			val:         types.StringValue("not ok"),
			format:      DNSNameFormat,
			expectError: true,
		},
		"InterfaceFormat_valid": {
			val:         types.StringValue("ok.ok"),
			format:      InterfaceFormat,
			expectError: false,
		},
		"InterfaceFormat_invalid": {
			val:         types.StringValue("not ok"),
			format:      InterfaceFormat,
			expectError: true,
		},
		"HexadecimalFormat_valid": {
			val:         types.StringValue("AB01"),
			format:      HexadecimalFormat,
			expectError: false,
		},
		"HexadecimalFormat_invalid": {
			val:         types.StringValue("AZ01"),
			format:      HexadecimalFormat,
			expectError: true,
		},
		"HexadecimalFormat_sensitive_valid": {
			val:         types.StringValue("cd01"),
			format:      HexadecimalFormat,
			expectError: false,
			sensitive:   true,
		},
		"HexadecimalFormat_sensitive_invalid": {
			val:         types.StringValue("cy01"),
			format:      HexadecimalFormat,
			expectError: true,
			sensitive:   true,
		},
		"ASPathRegularExpression_valid": {
			val:         types.StringValue(".* 209 .*"),
			format:      ASPathRegularExpression,
			expectError: false,
		},
		"ASPathRegularExpression_invalid": {
			val:         types.StringValue(".* AS209 .*"),
			format:      ASPathRegularExpression,
			expectError: true,
		},
		"AlgorithmFormat_valid": {
			val:         types.StringValue("ok@ok.net"),
			format:      AlgorithmFormat,
			expectError: false,
		},
		"AlgorithmFormat_invalid": {
			val:         types.StringValue("not ok@ok.net"),
			format:      AlgorithmFormat,
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
			if test.sensitive {
				StringFormat(test.format).WithSensitiveData().ValidateString(context.TODO(), request, &response)
			} else {
				StringFormat(test.format).ValidateString(context.TODO(), request, &response)
			}

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
