package tfplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var (
	_ planmodifier.Object = BlockRemoveNull()
	_ planmodifier.Object = BlockSetUnsetRequireReplace()
)

type BlockRemoveNullModifier struct{}

func BlockRemoveNull() BlockRemoveNullModifier {
	return BlockRemoveNullModifier{}
}

func (m BlockRemoveNullModifier) Description(_ context.Context) string {
	return "If block is not configured, modify plan to null."
}

func (m BlockRemoveNullModifier) MarkdownDescription(_ context.Context) string {
	return "If block is not configured, modify plan to null."
}

func (m BlockRemoveNullModifier) PlanModifyObject(
	_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse,
) {
	// Reference: https://github.com/hashicorp/terraform/issues/32460
	if req.ConfigValue.IsNull() {
		resp.PlanValue = req.ConfigValue
	}
}

type BlockSetUnsetRequireReplaceModifier struct{}

func BlockSetUnsetRequireReplace() BlockSetUnsetRequireReplaceModifier {
	return BlockSetUnsetRequireReplaceModifier{}
}

func (m BlockSetUnsetRequireReplaceModifier) Description(_ context.Context) string {
	return "If the presence of block changes, Terraform will destroy and recreate the resource."
}

func (m BlockSetUnsetRequireReplaceModifier) MarkdownDescription(_ context.Context) string {
	return "If the presence of block changes, Terraform will destroy and recreate the resource."
}

func (m BlockSetUnsetRequireReplaceModifier) PlanModifyObject(
	_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse,
) {
	// Do not replace on resource creation.
	if req.State.Raw.IsNull() {
		return
	}

	// Do not replace on resource destroy.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Reference: https://github.com/hashicorp/terraform/issues/32460
	if req.ConfigValue.IsNull() {
		resp.PlanValue = req.ConfigValue
	}

	// Replace if add block (null in state, not null in plan)
	if req.StateValue.IsNull() && !req.PlanValue.IsNull() {
		resp.RequiresReplace = true
	}

	// Replace if remove block (null in plan, not null in state)
	if req.PlanValue.IsNull() && !req.StateValue.IsNull() {
		resp.RequiresReplace = true
	}
}
