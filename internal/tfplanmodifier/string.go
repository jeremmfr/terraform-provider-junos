package tfplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = StringDefault("")

type StringDefaultModifier struct {
	defaultValue string
}

// If value is not configured, defaults to defaultValue.
func StringDefault(defaultValue string) StringDefaultModifier {
	return StringDefaultModifier{
		defaultValue: defaultValue,
	}
}

func (m StringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %q", m.defaultValue)
}

func (m StringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.defaultValue)
}

func (m StringDefaultModifier) PlanModifyString(
	ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse,
) {
	if !req.ConfigValue.IsNull() {
		return
	}
	if req.Plan.Raw.IsNull() {
		return
	}

	resp.PlanValue = types.StringValue(m.defaultValue)
}
