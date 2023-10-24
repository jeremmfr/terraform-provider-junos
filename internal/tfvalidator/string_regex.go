package tfvalidator

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = StringRegex()

type StringRegexValidator struct{}

func StringRegex() StringRegexValidator {
	return StringRegexValidator{}
}

func (v StringRegexValidator) Description(_ context.Context) string {
	return "Must be a valid regex."
}

func (v StringRegexValidator) MarkdownDescription(_ context.Context) string {
	return "Must be a valid regex."
}

func (v StringRegexValidator) ValidateString(
	_ context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	if _, err := regexp.Compile(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid regex",
			err.Error(),
		)

		return
	}
}
