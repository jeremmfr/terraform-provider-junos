package tfplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// StringUseStateNullForUnknown returns a plan modifier that copies a null prior state
// value into the planned value.
func StringUseStateNullForUnknown() planmodifier.String {
	return stringUseStateNullForUnknownModifier{}
}

type stringUseStateNullForUnknownModifier struct{}

func (m stringUseStateNullForUnknownModifier) Description(_ context.Context) string {
	return "If there is a resource state and attribute is null in state, keep null in plan."
}

func (m stringUseStateNullForUnknownModifier) MarkdownDescription(_ context.Context) string {
	return "If there is a resource state and attribute is null in state, keep null in plan."
}

// PlanModifyString implements the plan modification logic.
func (m stringUseStateNullForUnknownModifier) PlanModifyString(
	_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse,
) {
	// Do nothing if there is no resource state.
	if req.State.Raw.IsNull() {
		return
	}
	// Do nothing if there is state value.
	if !req.StateValue.IsNull() {
		return
	}
	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}
	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	resp.PlanValue = req.StateValue
}
