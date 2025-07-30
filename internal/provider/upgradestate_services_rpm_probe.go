package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *servicesRpmProbe) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"delegate_probes": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"test": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"data_fill": schema.StringAttribute{
									Optional: true,
								},
								"data_size": schema.Int64Attribute{
									Optional: true,
								},
								"destination_interface": schema.StringAttribute{
									Optional: true,
								},
								"destination_port": schema.Int64Attribute{
									Optional: true,
								},
								"dscp_code_points": schema.StringAttribute{
									Optional: true,
								},
								"hardware_timestamp": schema.BoolAttribute{
									Optional: true,
								},
								"history_size": schema.Int64Attribute{
									Optional: true,
								},
								"inet6_source_address": schema.StringAttribute{
									Optional: true,
								},
								"moving_average_size": schema.Int64Attribute{
									Optional: true,
								},
								"one_way_hardware_timestamp": schema.BoolAttribute{
									Optional: true,
								},
								"probe_count": schema.Int64Attribute{
									Optional: true,
								},
								"probe_interval": schema.Int64Attribute{
									Optional: true,
								},
								"probe_type": schema.StringAttribute{
									Optional: true,
								},
								"routing_instance": schema.StringAttribute{
									Optional: true,
								},
								"source_address": schema.StringAttribute{
									Optional: true,
								},
								"target_type": schema.StringAttribute{
									Optional: true,
								},
								"target_value": schema.StringAttribute{
									Optional: true,
								},
								"test_interval": schema.Int64Attribute{
									Optional: true,
								},
								"traps": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"ttl": schema.Int64Attribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"rpm_scale": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"tests_count": schema.Int64Attribute{
												Required: true,
											},
											"destination_interface": schema.StringAttribute{
												Optional: true,
											},
											"destination_subunit_cnt": schema.Int64Attribute{
												Optional: true,
											},
											"source_address_base": schema.StringAttribute{
												Optional: true,
											},
											"source_count": schema.Int64Attribute{
												Optional: true,
											},
											"source_step": schema.StringAttribute{
												Optional: true,
											},
											"source_inet6_address_base": schema.StringAttribute{
												Optional: true,
											},
											"source_inet6_count": schema.Int64Attribute{
												Optional: true,
											},
											"source_inet6_step": schema.StringAttribute{
												Optional: true,
											},
											"target_address_base": schema.StringAttribute{
												Optional: true,
											},
											"target_count": schema.Int64Attribute{
												Optional: true,
											},
											"target_step": schema.StringAttribute{
												Optional: true,
											},
											"target_inet6_address_base": schema.StringAttribute{
												Optional: true,
											},
											"target_inet6_count": schema.Int64Attribute{
												Optional: true,
											},
											"target_inet6_step": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"thresholds": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"egress_time": schema.Int64Attribute{
												Optional: true,
											},
											"ingress_time": schema.Int64Attribute{
												Optional: true,
											},
											"jitter_egress": schema.Int64Attribute{
												Optional: true,
											},
											"jitter_ingress": schema.Int64Attribute{
												Optional: true,
											},
											"jitter_rtt": schema.Int64Attribute{
												Optional: true,
											},
											"rtt": schema.Int64Attribute{
												Optional: true,
											},
											"std_dev_egress": schema.Int64Attribute{
												Optional: true,
											},
											"std_dev_ingress": schema.Int64Attribute{
												Optional: true,
											},
											"std_dev_rtt": schema.Int64Attribute{
												Optional: true,
											},
											"successive_loss": schema.Int64Attribute{
												Optional: true,
											},
											"total_loss": schema.Int64Attribute{
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
			StateUpgrader: upgradeServicesRpmProbeV0toV1,
		},
	}
}

func upgradeServicesRpmProbeV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID             types.String `tfsdk:"id"`
		Name           types.String `tfsdk:"name"`
		DelegateProbes types.Bool   `tfsdk:"delegate_probes"`
		Test           []struct {
			Name                    types.String   `tfsdk:"name"`
			DataFill                types.String   `tfsdk:"data_fill"`
			DataSize                types.Int64    `tfsdk:"data_size"`
			DestinationInterface    types.String   `tfsdk:"destination_interface"`
			DestinationPort         types.Int64    `tfsdk:"destination_port"`
			DscpCodePoints          types.String   `tfsdk:"dscp_code_points"`
			HardwareTimestamp       types.Bool     `tfsdk:"hardware_timestamp"`
			HistorySize             types.Int64    `tfsdk:"history_size"`
			Inet6SourceAddress      types.String   `tfsdk:"inet6_source_address"`
			MovingAverageSize       types.Int64    `tfsdk:"moving_average_size"`
			OneWayHardwareTimestamp types.Bool     `tfsdk:"one_way_hardware_timestamp"`
			ProbeCount              types.Int64    `tfsdk:"probe_count"`
			ProbeInterval           types.Int64    `tfsdk:"probe_interval"`
			ProbeType               types.String   `tfsdk:"probe_type"`
			RoutingInstance         types.String   `tfsdk:"routing_instance"`
			SourceAddress           types.String   `tfsdk:"source_address"`
			TargetType              types.String   `tfsdk:"target_type"`
			TargetValue             types.String   `tfsdk:"target_value"`
			TestInterval            types.Int64    `tfsdk:"test_interval"`
			Traps                   []types.String `tfsdk:"traps"`
			TTL                     types.Int64    `tfsdk:"ttl"`
			RpmScale                []struct {
				TestsCount             types.Int64  `tfsdk:"tests_count"`
				DestinationInterface   types.String `tfsdk:"destination_interface"`
				DestinationSubunitCnt  types.Int64  `tfsdk:"destination_subunit_cnt"`
				SourceAddressBase      types.String `tfsdk:"source_address_base"`
				SourceCount            types.Int64  `tfsdk:"source_count"`
				SourceStep             types.String `tfsdk:"source_step"`
				SourceInet6AddressBase types.String `tfsdk:"source_inet6_address_base"`
				SourceInet6Count       types.Int64  `tfsdk:"source_inet6_count"`
				SourceInet6Step        types.String `tfsdk:"source_inet6_step"`
				TargetAddressBase      types.String `tfsdk:"target_address_base"`
				TargetCount            types.Int64  `tfsdk:"target_count"`
				TargetStep             types.String `tfsdk:"target_step"`
				TargetInet6AddressBase types.String `tfsdk:"target_inet6_address_base"`
				TargetInet6Count       types.Int64  `tfsdk:"target_inet6_count"`
				TargetInet6Step        types.String `tfsdk:"target_inet6_step"`
			} `tfsdk:"rpm_scale"`
			Thresholds []struct {
				EgressTime     types.Int64 `tfsdk:"egress_time"`
				IngressTime    types.Int64 `tfsdk:"ingress_time"`
				JitterEgress   types.Int64 `tfsdk:"jitter_egress"`
				JitterIngress  types.Int64 `tfsdk:"jitter_ingress"`
				JitterRtt      types.Int64 `tfsdk:"jitter_rtt"`
				Rtt            types.Int64 `tfsdk:"rtt"`
				StdDevEgress   types.Int64 `tfsdk:"std_dev_egress"`
				StdDevIngress  types.Int64 `tfsdk:"std_dev_ingress"`
				StdDevRtt      types.Int64 `tfsdk:"std_dev_rtt"`
				SuccessiveLoss types.Int64 `tfsdk:"successive_loss"`
				TotalLoss      types.Int64 `tfsdk:"total_loss"`
			} `tfsdk:"thresholds"`
		} `tfsdk:"test"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 servicesRpmProbeData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.DelegateProbes = dataV0.DelegateProbes
	for _, blockV0 := range dataV0.Test {
		blockV1 := servicesRpmProbeBlockTest{
			Name:                    blockV0.Name,
			DataFill:                blockV0.DataFill,
			DataSize:                blockV0.DataSize,
			DestinationInterface:    blockV0.DestinationInterface,
			DestinationPort:         blockV0.DestinationPort,
			DscpCodePoints:          blockV0.DscpCodePoints,
			HardwareTimestamp:       blockV0.HardwareTimestamp,
			HistorySize:             blockV0.HistorySize,
			Inet6SourceAddress:      blockV0.Inet6SourceAddress,
			MovingAverageSize:       blockV0.MovingAverageSize,
			OneWayHardwareTimestamp: blockV0.OneWayHardwareTimestamp,
			ProbeCount:              blockV0.ProbeCount,
			ProbeInterval:           blockV0.ProbeInterval,
			ProbeType:               blockV0.ProbeType,
			RoutingInstance:         blockV0.RoutingInstance,
			SourceAddress:           blockV0.SourceAddress,
			TargetType:              blockV0.TargetType,
			TargetValue:             blockV0.TargetValue,
			TestInterval:            blockV0.TestInterval,
			Traps:                   blockV0.Traps,
			TTL:                     blockV0.TTL,
		}
		if len(blockV0.RpmScale) > 0 {
			blockV1.RpmScale = &servicesRpmProbeBlockTestBlockRpmScale{
				TestsCount:             blockV0.RpmScale[0].TestsCount,
				DestinationInterface:   blockV0.RpmScale[0].DestinationInterface,
				DestinationSubunitCnt:  blockV0.RpmScale[0].DestinationSubunitCnt,
				SourceAddressBase:      blockV0.RpmScale[0].SourceAddressBase,
				SourceCount:            blockV0.RpmScale[0].SourceCount,
				SourceStep:             blockV0.RpmScale[0].SourceStep,
				SourceInet6AddressBase: blockV0.RpmScale[0].SourceInet6AddressBase,
				SourceInet6Count:       blockV0.RpmScale[0].SourceInet6Count,
				SourceInet6Step:        blockV0.RpmScale[0].SourceInet6Step,
				TargetAddressBase:      blockV0.RpmScale[0].TargetAddressBase,
				TargetCount:            blockV0.RpmScale[0].TargetCount,
				TargetStep:             blockV0.RpmScale[0].TargetStep,
				TargetInet6AddressBase: blockV0.RpmScale[0].TargetInet6AddressBase,
				TargetInet6Count:       blockV0.RpmScale[0].TargetInet6Count,
				TargetInet6Step:        blockV0.RpmScale[0].TargetInet6Step,
			}
		}
		if len(blockV0.Thresholds) > 0 {
			blockV1.Thresholds = &servicesRpmProbeBlockTestBlockThresholds{
				EgressTime:     blockV0.Thresholds[0].EgressTime,
				IngressTime:    blockV0.Thresholds[0].IngressTime,
				JitterEgress:   blockV0.Thresholds[0].JitterEgress,
				JitterIngress:  blockV0.Thresholds[0].JitterIngress,
				JitterRtt:      blockV0.Thresholds[0].JitterRtt,
				Rtt:            blockV0.Thresholds[0].Rtt,
				StdDevEgress:   blockV0.Thresholds[0].StdDevEgress,
				StdDevIngress:  blockV0.Thresholds[0].StdDevIngress,
				StdDevRtt:      blockV0.Thresholds[0].StdDevRtt,
				SuccessiveLoss: blockV0.Thresholds[0].SuccessiveLoss,
				TotalLoss:      blockV0.Thresholds[0].TotalLoss,
			}
		}
		dataV1.Test = append(dataV1.Test, blockV1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
