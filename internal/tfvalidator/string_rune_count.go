package tfvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ validator.String = StringRuneCount(DotRune, 1)
	_ validator.String = String1DotCount()
)

type runeType int

const (
	DotRune runeType = iota
)

type StringRuneCountValidator struct {
	char   runeType
	number int
}

func StringRuneCount(char runeType, number int) StringRuneCountValidator {
	if number < 1 {
		number = 1
	}

	return StringRuneCountValidator{
		char:   char,
		number: number,
	}
}

func String1DotCount() StringRuneCountValidator {
	return StringRuneCountValidator{
		char:   DotRune,
		number: 1,
	}
}

func (t runeType) Singular() string {
	switch t {
	case DotRune:
		return "dot"
	default:
		return ""
	}
}

func (t runeType) Plural() string {
	switch t {
	case DotRune:
		return "dots"
	default:
		return ""
	}
}

func (t runeType) Rune() rune {
	switch t {
	case DotRune:
		return '.'
	default:
		return '\u0000'
	}
}

func (v StringRuneCountValidator) Description(_ context.Context) string {
	if v.number > 1 {
		return fmt.Sprintf("Must be a string with %d %s.", v.number, v.char.Plural())
	}

	return fmt.Sprintf("Must be a string with %d %s.", v.number, v.char.Singular())
}

func (v StringRuneCountValidator) MarkdownDescription(_ context.Context) string {
	if v.number > 1 {
		return fmt.Sprintf("Must be a string with %d %s.", v.number, v.char.Plural())
	}

	return fmt.Sprintf("Must be a string with %d %s.", v.number, v.char.Singular())
}

func (v StringRuneCountValidator) ValidateString(
	_ context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()
	if strings.Count(value, string(v.char.Rune())) != v.number {
		if v.number > 1 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing Character",
				fmt.Sprintf("string must have %d %s", v.number, v.char.Plural()),
			)

			return
		}

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing Character",
			fmt.Sprintf("string must have %d %s", v.number, v.char.Singular()),
		)

		return
	}
}
