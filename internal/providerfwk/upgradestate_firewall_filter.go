package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *firewallFilter) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"family": schema.StringAttribute{
						Required: true,
					},
					"interface_specific": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"term": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"filter": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"from": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"address": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"address_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_address": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_address_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_port": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_port_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_prefix_list": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_prefix_list_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"icmp_code": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"icmp_code_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"icmp_type": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"icmp_type_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"is_fragment": schema.BoolAttribute{
												Optional: true,
											},
											"next_header": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"next_header_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"port": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"port_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"prefix_list": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"prefix_list_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"protocol": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"protocol_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_address": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_address_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_port": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_port_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_prefix_list": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_prefix_list_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"tcp_established": schema.BoolAttribute{
												Optional: true,
											},
											"tcp_flags": schema.StringAttribute{
												Optional: true,
											},
											"tcp_initial": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
								"then": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.StringAttribute{
												Optional: true,
											},
											"count": schema.StringAttribute{
												Optional: true,
											},
											"log": schema.BoolAttribute{
												Optional: true,
											},
											"packet_mode": schema.BoolAttribute{
												Optional: true,
											},
											"policer": schema.StringAttribute{
												Optional: true,
											},
											"port_mirror": schema.BoolAttribute{
												Optional: true,
											},
											"routing_instance": schema.StringAttribute{
												Optional: true,
											},
											"sample": schema.BoolAttribute{
												Optional: true,
											},
											"service_accounting": schema.BoolAttribute{
												Optional: true,
											},
											"syslog": schema.BoolAttribute{
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
			StateUpgrader: upgradeFirewallFilterStateV0toV1,
		},
	}
}

func upgradeFirewallFilterStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		InterfaceSpecific types.Bool   `tfsdk:"interface_specific"`
		ID                types.String `tfsdk:"id"`
		Name              types.String `tfsdk:"name"`
		Family            types.String `tfsdk:"family"`
		Term              []struct {
			Name   types.String                       `tfsdk:"name"`
			Filter types.String                       `tfsdk:"filter"`
			From   []firewallFilterBlockTermBlockFrom `tfsdk:"from"`
			Then   []firewallFilterBlockTermBlockThen `tfsdk:"then"`
		} `tfsdk:"term"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 firewallFilterData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.InterfaceSpecific = dataV0.InterfaceSpecific
	dataV1.Family = dataV0.Family
	for _, blockV0 := range dataV0.Term {
		blockV1 := firewallFilterBlockTerm{
			Name:   blockV0.Name,
			Filter: blockV0.Filter,
		}
		if len(blockV0.From) > 0 {
			blockV1.From = &blockV0.From[0]
		}
		if len(blockV0.Then) > 0 {
			blockV1.Then = &blockV0.Then[0]
		}
		dataV1.Term = append(dataV1.Term, blockV1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
