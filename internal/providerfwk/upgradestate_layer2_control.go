package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *layer2Control) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"nonstop_bridging": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"bpdu_block": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"disable_timeout": schema.Int64Attribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"interface": schema.SetNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Required: true,
											},
											"disable": schema.BoolAttribute{
												Optional: true,
											},
											"drop": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"mac_rewrite_interface": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"enable_all_ifl": schema.BoolAttribute{
									Optional: true,
								},
								"protocol": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeLayer2ControlV0toV1,
		},
	}
}

func upgradeLayer2ControlV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID              types.String `tfsdk:"id"`
		NonstopBridging types.Bool   `tfsdk:"nonstop_bridging"`
		BpduBlock       []struct {
			DisableTimeout types.Int64 `tfsdk:"disable_timeout"`
			Interface      []struct {
				Name    types.String `tfsdk:"name"`
				Disable types.Bool   `tfsdk:"disable"`
				Drop    types.Bool   `tfsdk:"drop"`
			} `tfsdk:"interface"`
		} `tfsdk:"bpdu_block"`
		MacRewriteInterface []struct {
			Name         types.String   `tfsdk:"name"`
			EnableAllIfl types.Bool     `tfsdk:"enable_all_ifl"`
			Protocol     []types.String `tfsdk:"protocol"`
		} `tfsdk:"mac_rewrite_interface"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 layer2ControlData
	dataV1.ID = dataV0.ID
	dataV1.NonstopBridging = dataV0.NonstopBridging
	if len(dataV0.BpduBlock) > 0 {
		dataV1.BpduBlock = &layer2ControlBlockBpduBlock{
			DisableTimeout: dataV0.BpduBlock[0].DisableTimeout,
		}
		for _, blockV0 := range dataV0.BpduBlock[0].Interface {
			dataV1.BpduBlock.Interface = append(dataV1.BpduBlock.Interface,
				layer2ControlBlockBpduBlockBlockInterface{
					Name:    blockV0.Name,
					Disable: blockV0.Disable,
					Drop:    blockV0.Drop,
				},
			)
		}
	}
	for _, blockV0 := range dataV0.MacRewriteInterface {
		dataV1.MacRewriteInterface = append(dataV1.MacRewriteInterface,
			layer2ControlBlockMacRewriteInterface{
				Name:         blockV0.Name,
				EnableAllIfl: blockV0.EnableAllIfl,
				Protocol:     blockV0.Protocol,
			},
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
