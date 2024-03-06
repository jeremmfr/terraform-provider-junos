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
	DefaultFormatAndSpace
	AddressNameFormat
	DNSNameFormat
	InterfaceFormat
	InterfaceWithWildcardFormat
	HexadecimalFormat
	ASPathRegularExpression
	AlgorithmFormat
)

func (f stringFormat) invalidRune() func(rune) bool {
	switch f {
	case DefaultFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '_'
		}
	case DefaultFormatAndSpace:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '_' && r != ' '
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
	case InterfaceWithWildcardFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '/' && r != '.' && r != ':' && r != '*'
		}
	case HexadecimalFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'f') && (r < 'A' || r > 'F') && (r < '0' || r > '9')
		}
	case ASPathRegularExpression:
		return func(r rune) bool {
			return r != '^' && r != '$' && r != '-' && r != ',' && r != '|' && r != '*' && r != '!' &&
				r != '+' && r != '?' && r != '{' && r != '}' && r != '.' && r != '[' && r != ']' &&
				r != '(' && r != ')' && (r < '0' || r > '9') && r != ' '
		}
	case AlgorithmFormat:
		return func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') &&
				r != '-' && r != '_' && r != '@' && r != '.'
		}
	default:
		return func(_ rune) bool {
			return true
		}
	}
}

func (f stringFormat) String() string {
	switch f {
	case DefaultFormat:
		return "letters, numbers, dashes and underscores"
	case DefaultFormatAndSpace:
		return "letters, numbers, dashes, underscores and spaces"
	case AddressNameFormat:
		return "letters, numbers, dashes, dots, colons, slashes and underscores"
	case DNSNameFormat:
		return "letters, numbers, dashes, dots and underscores"
	case InterfaceFormat:
		return "letters, numbers, dashes, slashes, dots and colons"
	case InterfaceWithWildcardFormat:
		return "letters, numbers, dashes, slashes, dots, colons and asterisk"
	case HexadecimalFormat:
		return "A-F or a-f letters and numbers"
	case ASPathRegularExpression:
		return "regular expression characters, numbers and spaces"
	case AlgorithmFormat:
		return "letters, numbers, dashes, underscores, at symbol and dots"
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

func (v StringFormatValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Must be a string consisting of %s.", v.format.String())
}

func (v StringFormatValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Must be a string consisting of %s.", v.format.String())
}

func (v StringFormatValidator) ValidateString(
	_ context.Context, req validator.StringRequest, resp *validator.StringResponse,
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
