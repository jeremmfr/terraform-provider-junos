package providerfwk

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

var (
	_ validator.String = &stringFormatValidator{}
	_ validator.String = &stringIPAddressValidator{}
	_ validator.String = &stringCIDRNetworkValidator{}
	_ validator.String = &stringWildcardNetworkValidator{}
)

type stringFormat int

const (
	defaultFormat stringFormat = iota
	addressNameFormat
	dnsNameFormat
	interfaceFormat
)

func (f stringFormat) InvalidRune() func(rune) bool {
	switch f {
	case defaultFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_'
		}
	case addressNameFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '_' && r != ':' && r != '.' && r != '/'
		}
	case dnsNameFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '_' && r != '.'
		}
	case interfaceFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '/' && r != '.' && r != ':'
		}
	default:
		return func(r rune) bool {
			return true
		}
	}
}

func (f stringFormat) String() string {
	switch f {
	case defaultFormat:
		return "letters, numbers, dashes and underscores"
	case addressNameFormat:
		return "letters, numbers, dashes, dots, colons, slashes and underscores"
	case dnsNameFormat:
		return "letters, numbers, dashes, dots and underscores"
	case interfaceFormat:
		return "letters, numbers, dashes, slashes, dots and colons"
	default:
		return ""
	}
}

type stringFormatValidator struct {
	format stringFormat
}

func newStringFormatValidator(format stringFormat) stringFormatValidator {
	return stringFormatValidator{format: format}
}

func (v stringFormatValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Must be a string consisting of %s.", v.format.String())
}

func (v stringFormatValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Must be a string consisting of %s.", v.format.String())
}

func (v stringFormatValidator) ValidateString(
	ctx context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()
	if index := strings.IndexFunc(value, v.format.InvalidRune()); index != -1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid String Character",
			fmt.Sprintf("string has an unauthorized character: %q", value[index]),
		)

		return
	}
}

type stringIPAddressValidator struct{}

func (v stringIPAddressValidator) Description(ctx context.Context) string {
	return "Must be a valid IP address."
}

func (v stringIPAddressValidator) MarkdownDescription(ctx context.Context) string {
	return "Must be a valid IP address."
}

func (v stringIPAddressValidator) ValidateString(
	ctx context.Context, req validator.StringRequest, resp *validator.StringResponse,
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
}

type stringCIDRNetworkValidator struct{}

func (v stringCIDRNetworkValidator) Description(ctx context.Context) string {
	return "Must be a valid CIDR network."
}

func (v stringCIDRNetworkValidator) MarkdownDescription(ctx context.Context) string {
	return "Must be a valid CIDR network."
}

func (v stringCIDRNetworkValidator) ValidateString(
	ctx context.Context, req validator.StringRequest, resp *validator.StringResponse,
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
			fmt.Sprintf("%s", err),
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

type stringWildcardNetworkValidator struct{}

func (v stringWildcardNetworkValidator) Description(ctx context.Context) string {
	return "Must be a valid wildcard address."
}

func (v stringWildcardNetworkValidator) MarkdownDescription(ctx context.Context) string {
	return "Must be a valid wildcard address."
}

func (v stringWildcardNetworkValidator) ValidateString(
	ctx context.Context, req validator.StringRequest, resp *validator.StringResponse,
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
		if !bchk.InSlice(octet, []string{"255", "254", "252", "248", "240", "224", "192", "128", "0"}) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Wildcard Address",
				fmt.Sprintf("string %q for mask must be in subnet mask format, octet %q is not", valueSplit[1], octet),
			)

			return
		}
	}
}

type stringRuneExclusionValidator struct {
	runes []rune
}

func newStringRuneExclusionValidator(runes ...rune) stringRuneExclusionValidator {
	var validator stringRuneExclusionValidator
	validator.runes = append(validator.runes, runes...)

	return validator
}

func newStringDoubleQuoteExclusionValidator() stringRuneExclusionValidator {
	return newStringRuneExclusionValidator('"')
}

func newStringSpaceExclusionValidator() stringRuneExclusionValidator {
	return newStringRuneExclusionValidator('\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0)
}

func (v stringRuneExclusionValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Must be a string without characters %s .", string(v.runes))
}

func (v stringRuneExclusionValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Must be a string without characters %s .", string(v.runes))
}

func (v stringRuneExclusionValidator) ValidateString(
	ctx context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()
	for _, r := range v.runes {
		if index := strings.IndexRune(value, r); index != -1 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid String Character",
				fmt.Sprintf("string has an unauthorized character: %q", value[index]),
			)

			return
		}
	}
}
