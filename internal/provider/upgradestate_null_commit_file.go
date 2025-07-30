package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *nullCommitFile) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"filename": schema.StringAttribute{
						Required: true,
					},
					"append_lines": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"clear_file_after_commit": schema.BoolAttribute{
						Optional: true,
					},
					"triggers": schema.MapAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
			StateUpgrader: upgradeNullCommitFileStateV0toV1,
		},
	}
}

func upgradeNullCommitFileStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                   types.String   `tfsdk:"id"`
		Filename             types.String   `tfsdk:"filename"`
		AppendLines          []types.String `tfsdk:"append_lines"`
		ClearFileAfterCommit types.Bool     `tfsdk:"clear_file_after_commit"`
		Triggers             types.Map      `tfsdk:"triggers"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 nullCommitFileData
	dataV1.ID = dataV0.ID
	dataV1.Filename = dataV0.Filename
	dataV1.AppendLines = dataV0.AppendLines
	dataV1.ClearFileAfterCommit = dataV0.ClearFileAfterCommit

	attributeTypes := make(map[string]attr.Type)
	for k, v := range dataV0.Triggers.Elements() {
		attributeTypes[k] = v.Type(ctx)
	}
	newTriggers, objDiags := types.ObjectValue(attributeTypes, dataV0.Triggers.Elements())
	if objDiags.HasError() {
		resp.Diagnostics.Append(objDiags...)

		return
	}
	dataV1.Triggers = types.DynamicValue(newTriggers)

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
