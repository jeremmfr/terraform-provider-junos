package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type removeNullBlockModifier struct{}

func (m removeNullBlockModifier) Description(ctx context.Context) string {
	return "If block is not configured, modify plan to null"
}

func (m removeNullBlockModifier) MarkdownDescription(ctx context.Context) string {
	return "If block is not configured, modify plan to null"
}

func (m removeNullBlockModifier) PlanModifyObject(
	ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse,
) {
	// Reference: https://github.com/hashicorp/terraform/issues/32460
	if req.ConfigValue.IsNull() {
		resp.PlanValue = req.ConfigValue
	}
}
