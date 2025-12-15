package tfplanmodifier

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = StringUseNullStateForUnknownModifier{}

type StringUseNullStateForUnknownModifier struct{}

// StringUseNullStateForUnknown returns a plan modifier that copies a null prior state
// value into the planned value.
func StringUseNullStateForUnknown() StringUseNullStateForUnknownModifier {
	return StringUseNullStateForUnknownModifier{}
}

func (m StringUseNullStateForUnknownModifier) Description(_ context.Context) string {
	return "If there is a resource state and attribute is null in state, keep null in plan."
}

func (m StringUseNullStateForUnknownModifier) MarkdownDescription(_ context.Context) string {
	return "If there is a resource state and attribute is null in state, keep null in plan."
}

// PlanModifyString implements the plan modification logic.
func (m StringUseNullStateForUnknownModifier) PlanModifyString(
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

// StringRemoveBlankLines returns a plan modifier that removes all blank lines
// (empty lines and lines containing only whitespace) from multiline string values.
func StringRemoveBlankLines() planmodifier.String {
	return stringRemoveBlankLines{}
}

type stringRemoveBlankLines struct{}

func (m stringRemoveBlankLines) Description(_ context.Context) string {
	return "Removes blank lines (empty lines and lines containing only whitespace) from multiline strings."
}

func (m stringRemoveBlankLines) MarkdownDescription(_ context.Context) string {
	return "Removes blank lines (empty lines and lines containing only whitespace) from multiline strings."
}

// PlanModifyString implements the plan modification logic.
func (m stringRemoveBlankLines) PlanModifyString(
	_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse,
) {
	// Do nothing if there is a null planned value.
	if req.PlanValue.IsNull() {
		return
	}

	// Do nothing if there is an unknown planned value.
	if req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if value doesn't have multiple lines.
	if !strings.Contains(req.PlanValue.ValueString(), "\n") {
		return
	}

	newValue := strings.Builder{}
	for item := range strings.SplitSeq(req.PlanValue.ValueString(), "\n") {
		if strings.TrimSpace(item) == "" {
			continue
		}

		_, _ = newValue.WriteString(item + "\n")
	}

	resp.PlanValue = types.StringValue(newValue.String())
}
