package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityNatSource) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"from": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required: true,
								},
								"value": schema.SetAttribute{
									ElementType: types.StringType,
									Required:    true,
								},
							},
						},
					},
					"to": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required: true,
								},
								"value": schema.SetAttribute{
									ElementType: types.StringType,
									Required:    true,
								},
							},
						},
					},
					"rule": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
							},
							Blocks: map[string]schema.Block{
								"match": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"application": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_address": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_address_name": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_port": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"protocol": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_address": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_address_name": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_port": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
										},
									},
								},
								"then": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												Required: true,
											},
											"pool": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSecurityNatSourceV0toV1,
		},
	}
}

func upgradeSecurityNatSourceV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID          types.String                   `tfsdk:"id"`
		Name        types.String                   `tfsdk:"name"`
		Description types.String                   `tfsdk:"description"`
		From        []securityNatSourceBlockFromTo `tfsdk:"from"`
		To          []securityNatSourceBlockFromTo `tfsdk:"to"`
		Rule        []struct {
			Name  types.String                           `tfsdk:"name"`
			Match []securityNatSourceBlockRuleBlockMatch `tfsdk:"match"`
			Then  []securityNatSourceBlockRuleBlockThen  `tfsdk:"then"`
		} `tfsdk:"rule"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityNatSourceData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Description = dataV0.Description
	if len(dataV0.From) > 0 {
		dataV1.From = &dataV0.From[0]
	}
	if len(dataV0.To) > 0 {
		dataV1.To = &dataV0.To[0]
	}
	for _, blockV0 := range dataV0.Rule {
		blockV1 := securityNatSourceBlockRule{
			Name: blockV0.Name,
		}
		if len(blockV0.Match) > 0 {
			blockV1.Match = &blockV0.Match[0]
		}
		if len(blockV0.Then) > 0 {
			blockV1.Then = &blockV0.Then[0]
		}
		dataV1.Rule = append(dataV1.Rule, blockV1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
