package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityNatDestination) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
								"application": schema.SetAttribute{
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
							},
							Blocks: map[string]schema.Block{
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
			StateUpgrader: upgradeSecurityNatDestinationV0toV1,
		},
	}
}

func upgradeSecurityNatDestinationV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Description types.String `tfsdk:"description"`
		From        []struct {
			Type  types.String   `tfsdk:"type"`
			Value []types.String `tfsdk:"value"`
		} `tfsdk:"from"`
		Rule []struct {
			Name                   types.String   `tfsdk:"name"`
			DestinationAddress     types.String   `tfsdk:"destination_address"`
			DestinationAddressName types.String   `tfsdk:"destination_address_name"`
			Application            []types.String `tfsdk:"application"`
			DestinationPort        []types.String `tfsdk:"destination_port"`
			Protocol               []types.String `tfsdk:"protocol"`
			SourceAddress          []types.String `tfsdk:"source_address"`
			SourceAddressName      []types.String `tfsdk:"source_address_name"`
			Then                   []struct {
				Type types.String `tfsdk:"type"`
				Pool types.String `tfsdk:"pool"`
			} `tfsdk:"then"`
		} `tfsdk:"rule"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityNatDestinationData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Description = dataV0.Description
	if len(dataV0.From) > 0 {
		dataV1.From = &securityNatDestinationBlockFrom{
			Type:  dataV0.From[0].Type,
			Value: dataV0.From[0].Value,
		}
	}
	for _, blockV0 := range dataV0.Rule {
		blockV1 := securityNatDestinationBlockRule{
			Name:                   blockV0.Name,
			DestinationAddress:     blockV0.DestinationAddress,
			DestinationAddressName: blockV0.DestinationAddressName,
			Application:            blockV0.Application,
			DestinationPort:        blockV0.DestinationPort,
			Protocol:               blockV0.Protocol,
			SourceAddress:          blockV0.SourceAddress,
			SourceAddressName:      blockV0.SourceAddressName,
		}
		if len(blockV0.Then) > 0 {
			blockV1.Then = &securityNatDestinationBlockRuleBlockThen{
				Type: blockV0.Then[0].Type,
				Pool: blockV0.Then[0].Pool,
			}
		}
		dataV1.Rule = append(dataV1.Rule, blockV1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
