package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Must be true or null, false is not supported.
type boolTrueValidator struct{}

func (m boolTrueValidator) Description(ctx context.Context) string {
	return "Must be true or null, false is not supported"
}

func (m boolTrueValidator) MarkdownDescription(ctx context.Context) string {
	return "Must be true or null, false is not supported"
}

func (m boolTrueValidator) ValidateBool(
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
