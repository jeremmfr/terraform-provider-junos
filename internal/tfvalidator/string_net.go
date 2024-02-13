package tfvalidator

import (
	"context"
	"fmt"
	"net"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ validator.String = StringIPAddress()
	_ validator.String = StringCIDR()
	_ validator.String = StringCIDRNetwork()
	_ validator.String = StringWildcardNetwork()
)

type StringIPAddressValidator struct {
	v4only bool
	v6only bool
}

func StringIPAddress() StringIPAddressValidator {
	return StringIPAddressValidator{}
}

func (v StringIPAddressValidator) IPv4Only() StringIPAddressValidator {
	v.v4only = true

	return v
}

func (v StringIPAddressValidator) IPv6Only() StringIPAddressValidator {
	v.v6only = true

	return v
}

func (v StringIPAddressValidator) Description(_ context.Context) string {
	return "Must be a valid IP address."
}

func (v StringIPAddressValidator) MarkdownDescription(_ context.Context) string {
	return "Must be a valid IP address."
}

func (v StringIPAddressValidator) ValidateString(
	_ context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	if ip := net.ParseIP(value); ip == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid IP Address",
			fmt.Sprintf("string is not an IP address: %q", value),
		)

		return
	}
	if v.v4only {
		if strings.Contains(value, ":") {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid IPv4 Address",
				fmt.Sprintf("string is not an IPv4 address: %q", value),
			)
		}
	}
	if v.v6only {
		if !strings.Contains(value, ":") {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid IPv4 Address",
				fmt.Sprintf("string is not an IPv6 address: %q", value),
			)
		}
	}
}

type StringCIDRValidator struct {
	v4only bool
	v6only bool
}

func StringCIDR() StringCIDRValidator {
	return StringCIDRValidator{}
}

func (v StringCIDRValidator) IPv4Only() StringCIDRValidator {
	v.v4only = true

	return v
}

func (v StringCIDRValidator) IPv6Only() StringCIDRValidator {
	v.v6only = true

	return v
}

func (v StringCIDRValidator) Description(_ context.Context) string {
	return "Must be a valid CIDR."
}

func (v StringCIDRValidator) MarkdownDescription(_ context.Context) string {
	return "Must be a valid CIDR."
}

func (v StringCIDRValidator) ValidateString(
	_ context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	ipAddr, _, err := net.ParseCIDR(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR",
			err.Error(),
		)

		return
	}
	if ipAddr == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR",
			fmt.Sprintf("invalid CIDR: %q", value),
		)

		return
	}

	if v.v4only {
		if strings.Contains(value, ":") {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid IPv4 CIDR",
				fmt.Sprintf("string is not an IPv4 CIDR: %q", value),
			)
		}
	}
	if v.v6only {
		if !strings.Contains(value, ":") {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid IPv6 CIDR",
				fmt.Sprintf("string is not an IPv6 CIDR: %q", value),
			)
		}
	}
}

type StringCIDRNetworkValidator struct{}

func StringCIDRNetwork() StringCIDRNetworkValidator {
	return StringCIDRNetworkValidator{}
}

func (v StringCIDRNetworkValidator) Description(_ context.Context) string {
	return "Must be a valid CIDR network."
}

func (v StringCIDRNetworkValidator) MarkdownDescription(_ context.Context) string {
	return "Must be a valid CIDR network."
}

func (v StringCIDRNetworkValidator) ValidateString(
	_ context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR Network",
			err.Error(),
		)

		return
	}
	if ipnet == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR Network",
			fmt.Sprintf("invalid CIDR: %q", value),
		)

		return
	}

	if value != ipnet.String() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid CIDR Network",
			fmt.Sprintf("string is not a CIDR network: %q != %q", value, ipnet),
		)

		return
	}
}

func StringCIDRNetworkValidateAttribute(
	_ context.Context, strAttr types.String,
) error {
	if strAttr.IsNull() || strAttr.IsUnknown() {
		return nil
	}

	value := strAttr.ValueString()

	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		return err //nolint:wrapcheck
	}
	if ipnet == nil {
		return fmt.Errorf("invalid CIDR: %q", value)
	}

	if value != ipnet.String() {
		return fmt.Errorf("string is not a CIDR network: %q != %q", value, ipnet)
	}

	return nil
}

type StringWildcardNetworkValidator struct{}

func StringWildcardNetwork() StringWildcardNetworkValidator {
	return StringWildcardNetworkValidator{}
}

func (v StringWildcardNetworkValidator) Description(_ context.Context) string {
	return "Must be a valid wildcard address."
}

func (v StringWildcardNetworkValidator) MarkdownDescription(_ context.Context) string {
	return "Must be a valid wildcard address."
}

func (v StringWildcardNetworkValidator) ValidateString(
	_ context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	if !strings.Contains(value, "/") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Wildcard Address",
			fmt.Sprintf("string has no mask: %q", value),
		)

		return
	}
	if strings.Contains(value, ":") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Wildcard Address",
			fmt.Sprintf("string looks like a IPv6 address but is not supported: %q", value),
		)

		return
	}
	valueSplit := strings.Split(value, "/")
	ip := net.ParseIP(valueSplit[0]).To4()
	if ip == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Wildcard Address",
			fmt.Sprintf("string %q is not a valid IP address", valueSplit[0]),
		)

		return
	}
	mask := net.ParseIP(valueSplit[1]).To4()
	if mask == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Wildcard Address",
			fmt.Sprintf("string %q for mask is improperly formatted, must be in x.x.x.x notation", valueSplit[1]),
		)

		return
	}
	for _, octet := range strings.Split(mask.String(), ".") {
		if !slices.Contains([]string{"255", "254", "252", "248", "240", "224", "192", "128", "0"}, octet) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Wildcard Address",
				fmt.Sprintf("string %q for mask must be in subnet mask format, octet %q is not", valueSplit[1], octet),
			)

			return
		}
	}
}

type StringMACAddressValidator struct {
	mac48ColonHexadecimalNotation bool
}

func StringMACAddress() StringMACAddressValidator {
	return StringMACAddressValidator{}
}

func (v StringMACAddressValidator) WithMac48ColonHexa() StringMACAddressValidator {
	v.mac48ColonHexadecimalNotation = true

	return v
}

func (v StringMACAddressValidator) Description(_ context.Context) string {
	return "Must be a valid MAC address."
}

func (v StringMACAddressValidator) MarkdownDescription(_ context.Context) string {
	return "Must be a valid MAC address."
}

func (v StringMACAddressValidator) ValidateString(
	_ context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	if _, err := net.ParseMAC(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid MAC Address",
			err.Error(),
		)

		return
	}

	if v.mac48ColonHexadecimalNotation {
		if value[2] != ':' {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid MAC Address",
				fmt.Sprintf("not with Colon-Hexadecimal notation mm:mm:mm:ss:ss:ss : %q", value),
			)

			return
		}

		if (len(value)+1)%3 != 0 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid MAC Address",
				fmt.Sprintf("not a 48-bit address mm:mm:mm:ss:ss:ss : %q", value),
			)

			return
		}
		if (len(value)+1)/3 != 6 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid MAC Address",
				fmt.Sprintf("not a 48-bit address mm:mm:mm:ss:ss:ss : %q", value),
			)

			return
		}
	}
}
