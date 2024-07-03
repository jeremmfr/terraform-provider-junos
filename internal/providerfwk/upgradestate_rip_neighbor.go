package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *ripNeighbor) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"group": schema.StringAttribute{
						Required: true,
					},
					"ng": schema.BoolAttribute{
						Optional: true,
					},
					"routing_instance": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"any_sender": schema.BoolAttribute{
						Optional: true,
					},
					"authentication_key": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"authentication_type": schema.StringAttribute{
						Optional: true,
					},
					"check_zero": schema.BoolAttribute{
						Optional: true,
					},
					"no_check_zero": schema.BoolAttribute{
						Optional: true,
					},
					"demand_circuit": schema.BoolAttribute{
						Optional: true,
					},
					"dynamic_peers": schema.BoolAttribute{
						Optional: true,
					},
					"import": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"interface_type_p2mp": schema.BoolAttribute{
						Optional: true,
					},
					"max_retrans_time": schema.Int64Attribute{
						Optional: true,
					},
					"message_size": schema.Int64Attribute{
						Optional: true,
					},
					"metric_in": schema.Int64Attribute{
						Optional: true,
					},
					"peer": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"receive": schema.StringAttribute{
						Optional: true,
					},
					"route_timeout": schema.Int64Attribute{
						Optional: true,
					},
					"send": schema.StringAttribute{
						Optional: true,
					},
					"update_interval": schema.Int64Attribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"authentication_selective_md5": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"key_id": schema.Int64Attribute{
									Required: true,
								},
								"key": schema.StringAttribute{
									Required:  true,
									Sensitive: true,
								},
								"start_time": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
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
			StateUpgrader: upgradeRipNeighborStateV0toV1,
		},
	}
}

func upgradeRipNeighborStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                         types.String   `tfsdk:"id"`
		Name                       types.String   `tfsdk:"name"`
		Group                      types.String   `tfsdk:"group"`
		Ng                         types.Bool     `tfsdk:"ng"`
		RoutingInstance            types.String   `tfsdk:"routing_instance"`
		AnySender                  types.Bool     `tfsdk:"any_sender"`
		AuthenticationKey          types.String   `tfsdk:"authentication_key"`
		AuthenticationType         types.String   `tfsdk:"authentication_type"`
		CheckZero                  types.Bool     `tfsdk:"check_zero"`
		NoCheckZero                types.Bool     `tfsdk:"no_check_zero"`
		DemandCircuit              types.Bool     `tfsdk:"demand_circuit"`
		DynamicPeers               types.Bool     `tfsdk:"dynamic_peers"`
		Import                     []types.String `tfsdk:"import"`
		InterfaceTypeP2mp          types.Bool     `tfsdk:"interface_type_p2mp"`
		MaxRetransTime             types.Int64    `tfsdk:"max_retrans_time"`
		MessageSize                types.Int64    `tfsdk:"message_size"`
		MetricIn                   types.Int64    `tfsdk:"metric_in"`
		Peer                       []types.String `tfsdk:"peer"`
		Receive                    types.String   `tfsdk:"receive"`
		RouteTimeout               types.Int64    `tfsdk:"route_timeout"`
		Send                       types.String   `tfsdk:"send"`
		UpdateInterval             types.Int64    `tfsdk:"update_interval"`
		AuthenticationSelectiveMD5 []struct {
			KeyID     types.Int64  `tfsdk:"key_id"`
			Key       types.String `tfsdk:"key"`
			StartTime types.String `tfsdk:"start_time"`
		} `tfsdk:"authentication_selective_md5"`
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

	var dataV1 ripNeighborData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Group = dataV0.Group
	if dataV0.Ng.ValueBool() {
		dataV1.Ng = dataV0.Ng
	}
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.AnySender = dataV0.AnySender
	dataV1.AuthenticationKey = dataV0.AuthenticationKey
	dataV1.AuthenticationType = dataV0.AuthenticationType
	dataV1.CheckZero = dataV0.CheckZero
	dataV1.NoCheckZero = dataV0.NoCheckZero
	dataV1.DemandCircuit = dataV0.DemandCircuit
	dataV1.DynamicPeers = dataV0.DynamicPeers
	dataV1.Import = dataV0.Import
	dataV1.InterfaceTypeP2mp = dataV0.InterfaceTypeP2mp
	dataV1.MaxRetransTime = dataV0.MaxRetransTime
	dataV1.MessageSize = dataV0.MessageSize
	dataV1.MetricIn = dataV0.MetricIn
	dataV1.Peer = dataV0.Peer
	dataV1.Receive = dataV0.Receive
	dataV1.RouteTimeout = dataV0.RouteTimeout
	dataV1.Send = dataV0.Send
	dataV1.UpdateInterval = dataV0.UpdateInterval
	for _, block := range dataV0.AuthenticationSelectiveMD5 {
		dataV1.AuthenticationSelectiveMD5 = append(dataV1.AuthenticationSelectiveMD5,
			ripNeighborBlockAuthenticationSelectiveMd5{
				KeyID:     block.KeyID,
				Key:       block.Key,
				StartTime: block.StartTime,
			},
		)
	}
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
