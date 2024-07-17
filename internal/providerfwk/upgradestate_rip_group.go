package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *ripGroup) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"ng": schema.BoolAttribute{
						Optional: true,
					},
					"routing_instance": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"demand_circuit": schema.BoolAttribute{
						Optional: true,
					},
					"export": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"import": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"max_retrans_time": schema.Int64Attribute{
						Optional: true,
					},
					"metric_out": schema.Int64Attribute{
						Optional: true,
					},
					"preference": schema.Int64Attribute{
						Optional: true,
					},
					"route_timeout": schema.Int64Attribute{
						Optional: true,
					},
					"update_interval": schema.Int64Attribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"bfd_liveness_detection": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"authentication_algorithm": schema.StringAttribute{
									Optional: true,
								},
								"authentication_key_chain": schema.StringAttribute{
									Optional: true,
								},
								"authentication_loose_check": schema.BoolAttribute{
									Optional: true,
								},
								"detection_time_threshold": schema.Int64Attribute{
									Optional: true,
								},
								"minimum_interval": schema.Int64Attribute{
									Optional: true,
								},
								"minimum_receive_interval": schema.Int64Attribute{
									Optional: true,
								},
								"multiplier": schema.Int64Attribute{
									Optional: true,
								},
								"no_adaptation": schema.BoolAttribute{
									Optional: true,
								},
								"transmit_interval_minimum_interval": schema.Int64Attribute{
									Optional: true,
								},
								"transmit_interval_threshold": schema.Int64Attribute{
									Optional: true,
								},
								"version": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeRipGroupStateV0toV1,
		},
	}
}

func upgradeRipGroupStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                   types.String   `tfsdk:"id"`
		Name                 types.String   `tfsdk:"name"`
		Ng                   types.Bool     `tfsdk:"ng"`
		RoutingInstance      types.String   `tfsdk:"routing_instance"`
		DemandCircuit        types.Bool     `tfsdk:"demand_circuit"`
		Export               []types.String `tfsdk:"export"`
		Import               []types.String `tfsdk:"import"`
		MaxRetransTime       types.Int64    `tfsdk:"max_retrans_time"`
		MetricOut            types.Int64    `tfsdk:"metric_out"`
		Preference           types.Int64    `tfsdk:"preference"`
		RouteTimeout         types.Int64    `tfsdk:"route_timeout"`
		UpdateInterval       types.Int64    `tfsdk:"update_interval"`
		BfdLivenessDetection []struct {
			AuthenticationAlgorithm         types.String `tfsdk:"authentication_algorithm"`
			AuthenticationKeyChain          types.String `tfsdk:"authentication_key_chain"`
			AuthenticationLooseCheck        types.Bool   `tfsdk:"authentication_loose_check"`
			DetectionTimeThreshold          types.Int64  `tfsdk:"detection_time_threshold"`
			MinimumInterval                 types.Int64  `tfsdk:"minimum_interval"`
			MinimumReceiveInterval          types.Int64  `tfsdk:"minimum_receive_interval"`
			Multiplier                      types.Int64  `tfsdk:"multiplier"`
			NoAdaptation                    types.Bool   `tfsdk:"no_adaptation"`
			TransmitIntervalMinimumInterval types.Int64  `tfsdk:"transmit_interval_minimum_interval"`
			TransmitIntervalThreshold       types.Int64  `tfsdk:"transmit_interval_threshold"`
			Version                         types.String `tfsdk:"version"`
		} `tfsdk:"bfd_liveness_detection"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 ripGroupData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	if dataV0.Ng.ValueBool() {
		dataV1.Ng = dataV0.Ng
	}
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.DemandCircuit = dataV0.DemandCircuit
	dataV1.Export = dataV0.Export
	dataV1.Import = dataV0.Import
	dataV1.MaxRetransTime = dataV0.MaxRetransTime
	dataV1.MetricOut = dataV0.MetricOut
	dataV1.Preference = dataV0.Preference
	dataV1.RouteTimeout = dataV0.RouteTimeout
	dataV1.UpdateInterval = dataV0.UpdateInterval

	if len(dataV0.BfdLivenessDetection) > 0 {
		dataV1.BfdLivenessDetection = &ripBlockBfdLivenessDetection{
			AuthenticationLooseCheck:        dataV0.BfdLivenessDetection[0].AuthenticationLooseCheck,
			AuthenticationAlgorithm:         dataV0.BfdLivenessDetection[0].AuthenticationAlgorithm,
			AuthenticationKeyChain:          dataV0.BfdLivenessDetection[0].AuthenticationKeyChain,
			DetectionTimeThreshold:          dataV0.BfdLivenessDetection[0].DetectionTimeThreshold,
			MinimumInterval:                 dataV0.BfdLivenessDetection[0].MinimumInterval,
			MinimumReceiveInterval:          dataV0.BfdLivenessDetection[0].MinimumReceiveInterval,
			Multiplier:                      dataV0.BfdLivenessDetection[0].Multiplier,
			NoAdaptation:                    dataV0.BfdLivenessDetection[0].NoAdaptation,
			TransmitIntervalMinimumInterval: dataV0.BfdLivenessDetection[0].TransmitIntervalMinimumInterval,
			TransmitIntervalThreshold:       dataV0.BfdLivenessDetection[0].TransmitIntervalThreshold,
			Version:                         dataV0.BfdLivenessDetection[0].Version,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
