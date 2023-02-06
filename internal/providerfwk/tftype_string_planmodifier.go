package providerfwk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stringDefaultModifier struct {
	defaultValue string
}

// If value is not configured, defaults to defaultValue.
func newStringDefaultModifier(defaultValue string) planmodifier.String {
	return stringDefaultModifier{
		defaultValue: defaultValue,
	}
}

func (m stringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %q", m.defaultValue)
}

func (m stringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.defaultValue)
}

func (m stringDefaultModifier) PlanModifyString(
	ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse,
) {
	if !req.ConfigValue.IsNull() {
		return
	}

	resp.PlanValue = types.StringValue(m.defaultValue)
}
