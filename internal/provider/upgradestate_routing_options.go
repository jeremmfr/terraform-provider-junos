package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *routingOptions) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"clean_on_destroy": schema.BoolAttribute{
						Optional: true,
					},
					"forwarding_table_export_configure_singly": schema.BoolAttribute{
						Optional: true,
					},
					"instance_export": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"instance_import": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"router_id": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"autonomous_system": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"number": schema.StringAttribute{
									Required: true,
								},
								"asdot_notation": schema.BoolAttribute{
									Optional: true,
								},
								"loops": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"forwarding_table": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"chain_composite_max_label_count": schema.Int64Attribute{
									Optional: true,
								},
								"chained_composite_next_hop_ingress": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"chained_composite_next_hop_transit": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"dynamic_list_next_hop": schema.BoolAttribute{
									Optional: true,
								},
								"ecmp_fast_reroute": schema.BoolAttribute{
									Optional: true,
								},
								"no_ecmp_fast_reroute": schema.BoolAttribute{
									Optional: true,
								},
								"export": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"indirect_next_hop": schema.BoolAttribute{
									Optional: true,
								},
								"no_indirect_next_hop": schema.BoolAttribute{
									Optional: true,
								},
								"indirect_next_hop_change_acknowledgements": schema.BoolAttribute{
									Optional: true,
								},
								"no_indirect_next_hop_change_acknowledgements": schema.BoolAttribute{
									Optional: true,
								},
								"krt_nexthop_ack_timeout": schema.Int64Attribute{
									Optional: true,
								},
								"remnant_holdtime": schema.Int64Attribute{
									Optional: true,
								},
								"unicast_reverse_path": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"graceful_restart": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"disable": schema.BoolAttribute{
									Optional: true,
								},
								"restart_duration": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeRoutingOptionsV0toV1,
		},
	}
}

func upgradeRoutingOptionsV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                                   types.String   `tfsdk:"id"`
		CleanOnDestroy                       types.Bool     `tfsdk:"clean_on_destroy"`
		ForwardingTableExportConfigureSingly types.Bool     `tfsdk:"forwarding_table_export_configure_singly"`
		InstanceExport                       []types.String `tfsdk:"instance_export"`
		InstanceImport                       []types.String `tfsdk:"instance_import"`
		RouterID                             types.String   `tfsdk:"router_id"`
		AutonomousSystem                     []struct {
			Number        types.String `tfsdk:"number"`
			ASdotNotation types.Bool   `tfsdk:"asdot_notation"`
			Loops         types.Int64  `tfsdk:"loops"`
		} `tfsdk:"autonomous_system"`
		ForwardingTable []struct {
			ChainCompositeMaxLabelCount             types.Int64    `tfsdk:"chain_composite_max_label_count"`
			ChainedCompositeNextHopIngress          []types.String `tfsdk:"chained_composite_next_hop_ingress"`
			ChainedCompositeNextHopTransit          []types.String `tfsdk:"chained_composite_next_hop_transit"`
			DynamicListNextHop                      types.Bool     `tfsdk:"dynamic_list_next_hop"`
			EcmpFastReroute                         types.Bool     `tfsdk:"ecmp_fast_reroute"`
			NoEcmpFastReroute                       types.Bool     `tfsdk:"no_ecmp_fast_reroute"`
			Export                                  []types.String `tfsdk:"export"`
			IndirectNextHop                         types.Bool     `tfsdk:"indirect_next_hop"`
			NoIndirectNextHop                       types.Bool     `tfsdk:"no_indirect_next_hop"`
			IndirectNextHopChangeAcknowledgements   types.Bool     `tfsdk:"indirect_next_hop_change_acknowledgements"`
			NoIndirectNextHopChangeAcknowledgements types.Bool     `tfsdk:"no_indirect_next_hop_change_acknowledgements"`
			KrtNexthopAckTimeout                    types.Int64    `tfsdk:"krt_nexthop_ack_timeout"`
			RemnantHoldtime                         types.Int64    `tfsdk:"remnant_holdtime"`
			UnicastReversePath                      types.String   `tfsdk:"unicast_reverse_path"`
		} `tfsdk:"forwarding_table"`
		GracefulRestart []struct {
			Disable         types.Bool  `tfsdk:"disable"`
			RestartDuration types.Int64 `tfsdk:"restart_duration"`
		} `tfsdk:"graceful_restart"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 routingOptionsData
	dataV1.ID = dataV0.ID
	if dataV0.CleanOnDestroy.ValueBool() {
		dataV1.CleanOnDestroy = dataV0.CleanOnDestroy
	}
	if dataV0.ForwardingTableExportConfigureSingly.ValueBool() {
		dataV1.ForwardingTableExportConfigureSingly = dataV0.ForwardingTableExportConfigureSingly
	}
	dataV1.InstanceExport = dataV0.InstanceExport
	dataV1.InstanceImport = dataV0.InstanceImport
	dataV1.RouterID = dataV0.RouterID
	if len(dataV0.AutonomousSystem) > 0 {
		dataV1.AutonomousSystem = &routingOptionsBlockAutonomousSystem{
			Number:        dataV0.AutonomousSystem[0].Number,
			ASdotNotation: dataV0.AutonomousSystem[0].ASdotNotation,
			Loops:         dataV0.AutonomousSystem[0].Loops,
		}
	}
	if len(dataV0.ForwardingTable) > 0 {
		dataV1.ForwardingTable = &routingOptionsBlockForwardingTable{
			ChainCompositeMaxLabelCount:             dataV0.ForwardingTable[0].ChainCompositeMaxLabelCount,
			ChainedCompositeNextHopIngress:          dataV0.ForwardingTable[0].ChainedCompositeNextHopIngress,
			ChainedCompositeNextHopTransit:          dataV0.ForwardingTable[0].ChainedCompositeNextHopTransit,
			DynamicListNextHop:                      dataV0.ForwardingTable[0].DynamicListNextHop,
			EcmpFastReroute:                         dataV0.ForwardingTable[0].EcmpFastReroute,
			NoEcmpFastReroute:                       dataV0.ForwardingTable[0].NoEcmpFastReroute,
			Export:                                  dataV0.ForwardingTable[0].Export,
			IndirectNextHop:                         dataV0.ForwardingTable[0].IndirectNextHop,
			NoIndirectNextHop:                       dataV0.ForwardingTable[0].NoIndirectNextHop,
			IndirectNextHopChangeAcknowledgements:   dataV0.ForwardingTable[0].IndirectNextHopChangeAcknowledgements,
			NoIndirectNextHopChangeAcknowledgements: dataV0.ForwardingTable[0].NoIndirectNextHopChangeAcknowledgements,
			KrtNexthopAckTimeout:                    dataV0.ForwardingTable[0].KrtNexthopAckTimeout,
			RemnantHoldtime:                         dataV0.ForwardingTable[0].RemnantHoldtime,
			UnicastReversePath:                      dataV0.ForwardingTable[0].UnicastReversePath,
		}
	}
	if len(dataV0.GracefulRestart) > 0 {
		dataV1.GracefulRestart = &routingOptionsBlockGracefulRestart{
			Disable:         dataV0.GracefulRestart[0].Disable,
			RestartDuration: dataV0.GracefulRestart[0].RestartDuration,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
