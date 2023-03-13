package tfvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Bool = BoolTrue()

type BoolTrueValidator struct{}

// Must be true or null, false is not supported.
func BoolTrue() BoolTrueValidator {
	return BoolTrueValidator{}
}

func (m BoolTrueValidator) Description(ctx context.Context) string {
	return "Must be true or null, false is not supported"
}

func (m BoolTrueValidator) MarkdownDescription(ctx context.Context) string {
	return "Must be true or null, false is not supported"
}

func (m BoolTrueValidator) ValidateBool(
	ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if !req.ConfigValue.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid boolean value",
			"Must be true or null, false is not supported",
		)
	}
}
