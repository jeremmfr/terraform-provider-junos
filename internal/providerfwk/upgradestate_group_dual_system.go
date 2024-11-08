package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *groupDualSystem) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"apply_groups": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
				},
				Blocks: map[string]schema.Block{
					"interface_fxp0": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"family_inet_address": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"cidr_ip": schema.StringAttribute{
												Required: true,
											},
											"master_only": schema.BoolAttribute{
												Optional: true,
											},
											"preferred": schema.BoolAttribute{
												Optional: true,
											},
											"primary": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
								"family_inet6_address": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"cidr_ip": schema.StringAttribute{
												Required: true,
											},
											"master_only": schema.BoolAttribute{
												Optional: true,
											},
											"preferred": schema.BoolAttribute{
												Optional: true,
											},
											"primary": schema.BoolAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"routing_options": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Blocks: map[string]schema.Block{
								"static_route": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"destination": schema.StringAttribute{
												Required: true,
											},
											"next_hop": schema.ListAttribute{
												ElementType: types.StringType,
												Required:    true,
											},
										},
									},
								},
							},
						},
					},
					"security": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"log_source_address": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
					"system": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"host_name": schema.StringAttribute{
									Optional: true,
								},
								"backup_router_address": schema.StringAttribute{
									Optional: true,
								},
								"backup_router_destination": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"inet6_backup_router_address": schema.StringAttribute{
									Optional: true,
								},
								"inet6_backup_router_destination": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeGroupDualSystemStateV0toV1,
		},
	}
}

func upgradeGroupDualSystemStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID            types.String `tfsdk:"id"`
		Name          types.String `tfsdk:"name"`
		ApplyGroups   types.Bool   `tfsdk:"apply_groups"`
		InterfaceFXP0 []struct {
			Description       types.String `tfsdk:"description"`
			FamilyInetAddress []struct {
				CidrIP     types.String `tfsdk:"cidr_ip"`
				MasterOnly types.Bool   `tfsdk:"master_only"`
				Preferred  types.Bool   `tfsdk:"preferred"`
				Primary    types.Bool   `tfsdk:"primary"`
			} `tfsdk:"family_inet_address"`
			FamilyInet6Address []struct {
				CidrIP     types.String `tfsdk:"cidr_ip"`
				MasterOnly types.Bool   `tfsdk:"master_only"`
				Preferred  types.Bool   `tfsdk:"preferred"`
				Primary    types.Bool   `tfsdk:"primary"`
			} `tfsdk:"family_inet6_address"`
		} `tfsdk:"interface_fxp0"`
		RoutingOptions []struct {
			StaticRoute []struct {
				Destination types.String   `tfsdk:"destination"`
				NextHop     []types.String `tfsdk:"next_hop"`
			} `tfsdk:"static_route"`
		} `tfsdk:"routing_options"`
		Security []struct {
			LogSourceAddress types.String `tfsdk:"log_source_address"`
		} `tfsdk:"security"`
		System []struct {
			HostName                     types.String   `tfsdk:"host_name"`
			BackupRouterAddress          types.String   `tfsdk:"backup_router_address"`
			BackupRouterDestination      []types.String `tfsdk:"backup_router_destination"`
			Inet6BackupRouterAddress     types.String   `tfsdk:"inet6_backup_router_address"`
			Inet6BackupRouterDestination []types.String `tfsdk:"inet6_backup_router_destination"`
		} `tfsdk:"system"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 groupDualSystemData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.ApplyGroups = dataV0.ApplyGroups
	if len(dataV0.InterfaceFXP0) > 0 {
		dataV1.InterfaceFXP0 = &groupDualSystemBlockInterfaceFXP0{
			Description: dataV0.InterfaceFXP0[0].Description,
		}
		for _, blockV0 := range dataV0.InterfaceFXP0[0].FamilyInetAddress {
			dataV1.InterfaceFXP0.FamilyInetAddress = append(dataV1.InterfaceFXP0.FamilyInetAddress,
				groupDualSystemBlockInterfaceFXP0BlockFamilyAddress{
					CidrIP:     blockV0.CidrIP,
					MasterOnly: blockV0.MasterOnly,
					Preferred:  blockV0.Preferred,
					Primary:    blockV0.Primary,
				},
			)
		}
		for _, blockV0 := range dataV0.InterfaceFXP0[0].FamilyInet6Address {
			dataV1.InterfaceFXP0.FamilyInet6Address = append(dataV1.InterfaceFXP0.FamilyInet6Address,
				groupDualSystemBlockInterfaceFXP0BlockFamilyAddress{
					CidrIP:     blockV0.CidrIP,
					MasterOnly: blockV0.MasterOnly,
					Preferred:  blockV0.Preferred,
					Primary:    blockV0.Primary,
				},
			)
		}
	}
	if len(dataV0.RoutingOptions) > 0 {
		dataV1.RoutingOptions = &groupDualSystemBlockRoutingOptions{}
		for _, blockV0 := range dataV0.RoutingOptions[0].StaticRoute {
			dataV1.RoutingOptions.StaticRoute = append(dataV1.RoutingOptions.StaticRoute,
				groupDualSystemBlockRoutingOptionsBlockStaticRoute{
					Destination: blockV0.Destination,
					NextHop:     blockV0.NextHop,
				})
		}
	}
	if len(dataV0.Security) > 0 {
		dataV1.Security = &groupDualSystemBlockSecurity{
			LogSourceAddress: dataV0.Security[0].LogSourceAddress,
		}
	}
	if len(dataV0.System) > 0 {
		dataV1.System = &groupDualSystemBlockSystem{
			HostName:                     dataV0.System[0].HostName,
			BackupRouterAddress:          dataV0.System[0].BackupRouterAddress,
			BackupRouterDestination:      dataV0.System[0].BackupRouterDestination,
			Inet6BackupRouterAddress:     dataV0.System[0].Inet6BackupRouterAddress,
			Inet6BackupRouterDestination: dataV0.System[0].Inet6BackupRouterDestination,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
