package tfvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = StringFormat(DefaultFormat)

type stringFormat int

const (
	DefaultFormat stringFormat = iota
	AddressNameFormat
	DNSNameFormat
	InterfaceFormat
	HexadecimalFormat
)

func (f stringFormat) invalidRune() func(rune) bool {
	switch f {
	case DefaultFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_'
		}
	case AddressNameFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '_' && r != ':' && r != '.' && r != '/'
		}
	case DNSNameFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '_' && r != '.'
		}
	case InterfaceFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '/' && r != '.' && r != ':'
		}
	case HexadecimalFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'f') && (r < 'A' || r > 'F') && (r < '0' || r > '9')
		}
	default:
		return func(r rune) bool {
			return true
		}
	}
}

func (f stringFormat) String() string {
	switch f {
	case DefaultFormat:
		return "letters, numbers, dashes and underscores"
	case AddressNameFormat:
		return "letters, numbers, dashes, dots, colons, slashes and underscores"
	case DNSNameFormat:
		return "letters, numbers, dashes, dots and underscores"
	case InterfaceFormat:
		return "letters, numbers, dashes, slashes, dots and colons"
	case HexadecimalFormat:
		return "A-F or a-f letters and numbers"
	default:
		return ""
	}
}

type StringFormatValidator struct {
	format    stringFormat
	sensitive bool
}

func StringFormat(format stringFormat) StringFormatValidator {
	return StringFormatValidator{format: format}
}

func (v StringFormatValidator) WithSensitiveData() StringFormatValidator {
	v.sensitive = true

	return v
}

func (v StringFormatValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Must be a string consisting of %s.", v.format.String())
}

func (v StringFormatValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Must be a string consisting of %s.", v.format.String())
}

func (v StringFormatValidator) ValidateString(
	ctx context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()
	if index := strings.IndexFunc(value, v.format.invalidRune()); index != -1 {
		if v.sensitive {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid String Character",
				fmt.Sprintf("string has an unauthorized character: not in %s", v.format),
			)

			return
		}

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid String Character",
			fmt.Sprintf("string has an unauthorized character: %q", value[index]),
		)

		return
	}
}
