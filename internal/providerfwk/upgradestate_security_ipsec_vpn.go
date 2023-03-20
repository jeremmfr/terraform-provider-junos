package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityIpsecVpn) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"bind_interface": schema.StringAttribute{
						Optional: true,
					},
					"df_bit": schema.StringAttribute{
						Optional: true,
					},
					"establish_tunnels": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"ike": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"gateway": schema.StringAttribute{
									Required: true,
								},
								"policy": schema.StringAttribute{
									Required: true,
								},
								"identity_local": schema.StringAttribute{
									Optional: true,
								},
								"identity_remote": schema.StringAttribute{
									Optional: true,
								},
								"identity_service": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"traffic_selector": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"local_ip": schema.StringAttribute{
									Required: true,
								},
								"remote_ip": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
					"vpn_monitor": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"destination_ip": schema.StringAttribute{
									Optional: true,
								},
								"optimized": schema.BoolAttribute{
									Optional: true,
								},
								"source_interface": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"source_interface_auto": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSecurityIpsecVpnV0toV1,
		},
	}
}

func upgradeSecurityIpsecVpnV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID               types.String                      `tfsdk:"id"`
		Name             types.String                      `tfsdk:"name"`
		BindInterface    types.String                      `tfsdk:"bind_interface"`
		DfBit            types.String                      `tfsdk:"df_bit"`
		EstablishTunnels types.String                      `tfsdk:"establish_tunnels"`
		Ike              []securityIpsecVpnIke             `tfsdk:"ike"`
		TrafficSelector  []securityIpsecVpnTrafficSelector `tfsdk:"traffic_selector"`
		VpnMonitor       []securityIpsecVpnVpnMonitor      `tfsdk:"vpn_monitor"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityIpsecVpnData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.BindInterface = dataV0.BindInterface
	dataV1.DfBit = dataV0.DfBit
	dataV1.EstablishTunnels = dataV0.EstablishTunnels
	if len(dataV0.Ike) > 0 {
		dataV1.Ike = &dataV0.Ike[0]
	}
	dataV1.TrafficSelector = dataV0.TrafficSelector
	if len(dataV0.VpnMonitor) > 0 {
		dataV1.VpnMonitor = &dataV0.VpnMonitor[0]
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
