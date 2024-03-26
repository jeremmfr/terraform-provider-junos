package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *evpn) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"routing_instance": schema.StringAttribute{
						Optional: true,
					},
					"encapsulation": schema.StringAttribute{
						Required: true,
					},
					"default_gateway": schema.StringAttribute{
						Optional: true,
					},
					"multicast_mode": schema.StringAttribute{
						Optional: true,
					},
					"routing_instance_evpn": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"switch_or_ri_options": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"route_distinguisher": schema.StringAttribute{
									Required: true,
								},
								"vrf_export": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"vrf_import": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"vrf_target": schema.StringAttribute{
									Optional: true,
								},
								"vrf_target_auto": schema.BoolAttribute{
									Optional: true,
								},
								"vrf_target_export": schema.StringAttribute{
									Optional: true,
								},
								"vrf_target_import": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeEvpnV0toV1,
		},
	}
}

func upgradeEvpnV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                  types.String `tfsdk:"id"`
		RoutingInstance     types.String `tfsdk:"routing_instance"`
		Encapsulation       types.String `tfsdk:"encapsulation"`
		DefaultGateway      types.String `tfsdk:"default_gateway"`
		MulticastMode       types.String `tfsdk:"multicast_mode"`
		RoutingInstanceEvpn types.Bool   `tfsdk:"routing_instance_evpn"`
		SwitchOrRIOptions   []struct {
			RouteDistinguisher types.String   `tfsdk:"route_distinguisher"`
			VRFExport          []types.String `tfsdk:"vrf_export"`
			VRFImport          []types.String `tfsdk:"vrf_import"`
			VRFTarget          types.String   `tfsdk:"vrf_target"`
			VRFTargetAuto      types.Bool     `tfsdk:"vrf_target_auto"`
			VRFTargetExport    types.String   `tfsdk:"vrf_target_export"`
			VRFTargetImport    types.String   `tfsdk:"vrf_target_import"`
		} `tfsdk:"switch_or_ri_options"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 evpnData
	dataV1.ID = dataV0.ID
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.RoutingInstanceEvpn = dataV0.RoutingInstanceEvpn
	dataV1.Encapsulation = dataV0.Encapsulation
	dataV1.DefaultGateway = dataV0.DefaultGateway
	dataV1.MulticastMode = dataV0.MulticastMode
	if len(dataV0.SwitchOrRIOptions) > 0 {
		dataV1.SwitchOrRIOptions = &evpnBlockSwitchOrRIOptions{
			VRFTargetAuto:      dataV0.SwitchOrRIOptions[0].VRFTargetAuto,
			RouteDistinguisher: dataV0.SwitchOrRIOptions[0].RouteDistinguisher,
			VRFExport:          dataV0.SwitchOrRIOptions[0].VRFExport,
			VRFImport:          dataV0.SwitchOrRIOptions[0].VRFImport,
			VRFTarget:          dataV0.SwitchOrRIOptions[0].VRFTarget,
			VRFTargetExport:    dataV0.SwitchOrRIOptions[0].VRFTargetExport,
			VRFTargetImport:    dataV0.SwitchOrRIOptions[0].VRFTargetImport,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
