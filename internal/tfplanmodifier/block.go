package tfplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.Object = BlockRemoveNull()

type BlockRemoveNullModifier struct{}

func BlockRemoveNull() BlockRemoveNullModifier {
	return BlockRemoveNullModifier{}
}

func (m BlockRemoveNullModifier) Description(ctx context.Context) string {
	return "If block is not configured, modify plan to null"
}

func (m BlockRemoveNullModifier) MarkdownDescription(ctx context.Context) string {
	return "If block is not configured, modify plan to null"
}

func (m BlockRemoveNullModifier) PlanModifyObject(
	ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse,
) {
	// Reference: https://github.com/hashicorp/terraform/issues/32460
	if req.ConfigValue.IsNull() {
		resp.PlanValue = req.ConfigValue
	}
}
