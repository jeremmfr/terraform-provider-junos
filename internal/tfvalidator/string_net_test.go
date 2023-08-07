package tfvalidator_test

import (
	"context"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringIPAddress(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		expectError bool
		v4only      bool
		v6only      bool
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
			val:         types.StringValue("192.0.2.1"),
			expectError: false,
		},
		"valid v6": {
			val:         types.StringValue("2001:2::1"),
			expectError: false,
		},
		"empty": {
			val:         types.StringValue(""),
			expectError: true,
		},
		"invalid": {
			val:         types.StringValue("292.0.2.1"),
			expectError: true,
		},
		"valid with v4only": {
			val:         types.StringValue("192.0.2.1"),
			expectError: false,
			v4only:      true,
		},
		"valid v6 but with v4only": {
			val:         types.StringValue("2001:2::1"),
			expectError: true,
			v4only:      true,
		},
		"invalid with v4only": {
			val:         types.StringValue("292.0.2.1"),
			expectError: true,
			v4only:      true,
		},
		"valid with v6only": {
			val:         types.StringValue("2001:2::1"),
			expectError: false,
			v6only:      true,
		},
		"valid v4 but with v6only": {
			val:         types.StringValue("192.0.2.1"),
			expectError: true,
			v6only:      true,
		},
		"invalid with v6only": {
			val:         types.StringValue("2001:2:::1"),
			expectError: true,
			v6only:      true,
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

			switch {
			case test.v4only:
				tfvalidator.StringIPAddress().IPv4Only().ValidateString(context.TODO(), request, &response)
			case test.v6only:
				tfvalidator.StringIPAddress().IPv6Only().ValidateString(context.TODO(), request, &response)
			default:
				tfvalidator.StringIPAddress().ValidateString(context.TODO(), request, &response)
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

func TestStringCIDR(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		expectError bool
		v4only      bool
		v6only      bool
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
			val:         types.StringValue("192.0.2.1/24"),
			expectError: false,
		},
		"empty": {
			val:         types.StringValue(""),
			expectError: true,
		},
		"invalid": {
			val:         types.StringValue("192.0.2.1"),
			expectError: true,
		},
		"invalid mask": {
			val:         types.StringValue("192.0.2.1/256"),
			expectError: true,
		},
		"invalid ip": {
			val:         types.StringValue("192.0.2."),
			expectError: true,
		},
		"valid v4": {
			val:         types.StringValue("192.0.2.1/24"),
			expectError: false,
			v4only:      true,
		},
		"invalid v4": {
			val:         types.StringValue("192.0.2.1"),
			expectError: true,
			v4only:      true,
		},
		"valid v6 but with v4only": {
			val:         types.StringValue("2001:2::1/64"),
			expectError: true,
			v4only:      true,
		},
		"valid v6": {
			val:         types.StringValue("2001:2::1/64"),
			expectError: false,
			v6only:      true,
		},
		"invalid v6": {
			val:         types.StringValue("2001:2:::1/64"),
			expectError: true,
			v6only:      true,
		},
		"valid v4 but with v6only": {
			val:         types.StringValue("192.0.2.1/24"),
			expectError: true,
			v6only:      true,
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
			switch {
			case test.v4only:
				tfvalidator.StringCIDR().IPv4Only().ValidateString(context.TODO(), request, &response)
			case test.v6only:
				tfvalidator.StringCIDR().IPv6Only().ValidateString(context.TODO(), request, &response)
			default:
				tfvalidator.StringCIDR().ValidateString(context.TODO(), request, &response)
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

func TestStringCIDRNetwork(t *testing.T) {
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
			val:         types.StringValue("192.0.2.0/24"),
			expectError: false,
		},
		"empty": {
			val:         types.StringValue(""),
			expectError: true,
		},
		"invalid": {
			val:         types.StringValue("192.0.2.1/24"),
			expectError: true,
		},
		"invalid mask": {
			val:         types.StringValue("192.0.2.1"),
			expectError: true,
		},
		"invalid ip": {
			val:         types.StringValue("192.0.2."),
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
			tfvalidator.StringCIDRNetwork().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestStringWildcardNetwork(t *testing.T) {
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
			val:         types.StringValue("192.0.2.0/255.255.255.0"),
			expectError: false,
		},
		"invalid": {
			val:         types.StringValue("192.0.2.0/255.255.253.0"),
			expectError: true,
		},
		"invalid mask": {
			val:         types.StringValue("192.0.2.0"),
			expectError: true,
		},
		"invalid ip": {
			val:         types.StringValue("192.0.2."),
			expectError: true,
		},
		"invalid ipv6": {
			val:         types.StringValue("2001:db8::/32"),
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
			tfvalidator.StringWildcardNetwork().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestStringMACAddress(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val            types.String
		expectError    bool
		mac48ColonHexa bool
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
			val:         types.StringValue("00:00:5e:00:53:01"),
			expectError: false,
		},
		"valid with Colon-Hexadecimal validation": {
			val:            types.StringValue("00:00:5e:00:53:01"),
			expectError:    false,
			mac48ColonHexa: true,
		},
		"invalid": {
			val:         types.StringValue("00:00:5e:00:53:zz"),
			expectError: true,
		},
		"valid without Colon-Hexadecimal notation": {
			val:         types.StringValue("0000.5e00.5301"),
			expectError: false,
		},
		"valid without Colon-Hexadecimal notation but need it": {
			val:            types.StringValue("0000.5e00.5301"),
			expectError:    true,
			mac48ColonHexa: true,
		},
		"invalid with Colon-Hexadecimal validation": {
			val:            types.StringValue("0000.5e00.53zz"),
			expectError:    true,
			mac48ColonHexa: true,
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
			switch {
			case test.mac48ColonHexa:
				tfvalidator.StringMACAddress().WithMac48ColonHexa().ValidateString(context.TODO(), request, &response)
			default:
				tfvalidator.StringMACAddress().ValidateString(context.TODO(), request, &response)
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
