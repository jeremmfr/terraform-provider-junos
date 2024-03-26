package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityNatStatic) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"configure_rules_singly": schema.BoolAttribute{
						Optional: true,
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
					"rule": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"destination_address": schema.StringAttribute{
									Optional: true,
								},
								"destination_address_name": schema.StringAttribute{
									Optional: true,
								},
								"destination_port": schema.Int64Attribute{
									Optional: true,
								},
								"destination_port_to": schema.Int64Attribute{
									Optional: true,
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
							Blocks: map[string]schema.Block{
								"then": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												Required: true,
											},
											"mapped_port": schema.Int64Attribute{
												Optional: true,
											},
											"mapped_port_to": schema.Int64Attribute{
												Optional: true,
											},
											"prefix": schema.StringAttribute{
												Optional: true,
											},
											"routing_instance": schema.StringAttribute{
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
			StateUpgrader: upgradeSecurityNatStaticV0toV1,
		},
	}
}

func upgradeSecurityNatStaticV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                   types.String `tfsdk:"id"`
		Name                 types.String `tfsdk:"name"`
		ConfigureRulesSingly types.Bool   `tfsdk:"configure_rules_singly"`
		Description          types.String `tfsdk:"description"`
		From                 []struct {
			Type  types.String   `tfsdk:"type"`
			Value []types.String `tfsdk:"value"`
		} `tfsdk:"from"`
		Rule []struct {
			Name                   types.String   `tfsdk:"name"`
			DestinationAddress     types.String   `tfsdk:"destination_address"`
			DestinationAddressName types.String   `tfsdk:"destination_address_name"`
			DestinationPort        types.Int64    `tfsdk:"destination_port"`
			DestiantionPortTo      types.Int64    `tfsdk:"destination_port_to"`
			SourceAddress          []types.String `tfsdk:"source_address"`
			SourceAddressName      []types.String `tfsdk:"source_address_name"`
			SourcePort             []types.String `tfsdk:"source_port"`
			Then                   []struct {
				Type            types.String `tfsdk:"type"`
				MappedPort      types.Int64  `tfsdk:"mapped_port"`
				MappedPortTo    types.Int64  `tfsdk:"mapped_port_to"`
				Prefix          types.String `tfsdk:"prefix"`
				RoutingInstance types.String `tfsdk:"routing_instance"`
			} `tfsdk:"then"`
		} `tfsdk:"rule"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityNatStaticData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.ConfigureRulesSingly = dataV0.ConfigureRulesSingly
	dataV1.Description = dataV0.Description
	if len(dataV0.From) > 0 {
		dataV1.From = &securityNatStaticBlockFrom{
			Type:  dataV0.From[0].Type,
			Value: dataV0.From[0].Value,
		}
	}
	for _, blockV0 := range dataV0.Rule {
		blockV1 := securityNatStaticBlockRule{
			Name:                   blockV0.Name,
			DestinationAddress:     blockV0.DestinationAddress,
			DestinationAddressName: blockV0.DestinationAddressName,
			DestinationPort:        blockV0.DestinationPort,
			DestiantionPortTo:      blockV0.DestiantionPortTo,
			SourceAddress:          blockV0.SourceAddress,
			SourceAddressName:      blockV0.SourceAddressName,
			SourcePort:             blockV0.SourcePort,
		}
		if len(blockV0.Then) > 0 {
			blockV1.Then = &securityNatStaticRuleBlockThen{
				Type:            blockV0.Then[0].Type,
				MappedPort:      blockV0.Then[0].MappedPort,
				MappedPortTo:    blockV0.Then[0].MappedPortTo,
				Prefix:          blockV0.Then[0].Prefix,
				RoutingInstance: blockV0.Then[0].RoutingInstance,
			}
		}
		dataV1.Rule = append(dataV1.Rule, blockV1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
