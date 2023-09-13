package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityNatStaticRule) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"rule_set": schema.StringAttribute{
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
			StateUpgrader: upgradeSecurityNatStaticRuleV0toV1,
		},
	}
}

func upgradeSecurityNatStaticRuleV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                     types.String   `tfsdk:"id"`
		Name                   types.String   `tfsdk:"name"`
		RuleSet                types.String   `tfsdk:"rule_set"`
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
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityNatStaticRuleData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.RuleSet = dataV0.RuleSet
	dataV1.DestinationAddress = dataV0.DestinationAddress
	dataV1.DestinationAddressName = dataV0.DestinationAddressName
	dataV1.DestinationPort = dataV0.DestinationPort
	dataV1.DestiantionPortTo = dataV0.DestiantionPortTo
	dataV1.SourceAddress = dataV0.SourceAddress
	dataV1.SourceAddressName = dataV0.SourceAddressName
	dataV1.SourcePort = dataV0.SourcePort
	if len(dataV0.Then) > 0 {
		dataV1.Then = &securityNatStaticRuleBlockThen{
			Type:            dataV0.Then[0].Type,
			MappedPort:      dataV0.Then[0].MappedPort,
			MappedPortTo:    dataV0.Then[0].MappedPortTo,
			Prefix:          dataV0.Then[0].Prefix,
			RoutingInstance: dataV0.Then[0].RoutingInstance,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
