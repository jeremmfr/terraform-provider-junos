package tfvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ validator.String = StringRuneExclusion('/')
	_ validator.String = StringDoubleQuoteExclusion()
	_ validator.String = StringSpaceExclusion()
)

type StringRuneExclusionValidator struct {
	runes []rune
}

func StringRuneExclusion(runes ...rune) StringRuneExclusionValidator {
	var v StringRuneExclusionValidator
	v.runes = append(v.runes, runes...)

	return v
}

func StringDoubleQuoteExclusion() StringRuneExclusionValidator {
	return StringRuneExclusion('"')
}

func StringSpaceExclusion() StringRuneExclusionValidator {
	return StringRuneExclusion('\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0)
}

func StringDotExclusion() StringRuneExclusionValidator {
	return StringRuneExclusion('.')
}

func (v StringRuneExclusionValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Must be a string without characters %s .", string(v.runes))
}

func (v StringRuneExclusionValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Must be a string without characters %s .", string(v.runes))
}

func (v StringRuneExclusionValidator) ValidateString(
	_ context.Context, req validator.StringRequest, resp *validator.StringResponse,
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
