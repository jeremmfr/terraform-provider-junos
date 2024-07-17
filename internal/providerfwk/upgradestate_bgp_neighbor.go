package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *bgpNeighbor) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"ip": schema.StringAttribute{
						Required: true,
					},
					"routing_instance": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"group": schema.StringAttribute{
						Required: true,
					},
					"accept_remote_nexthop": schema.BoolAttribute{
						Optional: true,
					},
					"advertise_external": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"advertise_external_conditional": schema.BoolAttribute{
						Optional: true,
					},
					"advertise_inactive": schema.BoolAttribute{
						Optional: true,
					},
					"advertise_peer_as": schema.BoolAttribute{
						Optional: true,
					},
					"no_advertise_peer_as": schema.BoolAttribute{
						Optional: true,
					},
					"as_override": schema.BoolAttribute{
						Optional: true,
					},
					"authentication_algorithm": schema.StringAttribute{
						Optional: true,
					},
					"authentication_key": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"authentication_key_chain": schema.StringAttribute{
						Optional: true,
					},
					"cluster": schema.StringAttribute{
						Optional: true,
					},
					"damping": schema.BoolAttribute{
						Optional: true,
					},
					"export": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"hold_time": schema.Int64Attribute{
						Optional: true,
					},
					"import": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"keep_all": schema.BoolAttribute{
						Optional: true,
					},
					"keep_none": schema.BoolAttribute{
						Optional: true,
					},
					"local_address": schema.StringAttribute{
						Optional: true,
					},
					"local_as": schema.StringAttribute{
						Optional: true,
					},
					"local_as_alias": schema.BoolAttribute{
						Optional: true,
					},
					"local_as_loops": schema.Int64Attribute{
						Optional: true,
					},
					"local_as_no_prepend_global_as": schema.BoolAttribute{
						Optional: true,
					},
					"local_as_private": schema.BoolAttribute{
						Optional: true,
					},
					"local_interface": schema.StringAttribute{
						Optional: true,
					},
					"local_preference": schema.Int64Attribute{
						Optional: true,
					},
					"log_updown": schema.BoolAttribute{
						Optional: true,
					},
					"metric_out": schema.Int64Attribute{
						Optional: true,
					},
					"metric_out_igp": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"metric_out_igp_delay_med_update": schema.BoolAttribute{
						Optional: true,
					},
					"metric_out_igp_offset": schema.Int64Attribute{
						Optional: true,
					},
					"metric_out_minimum_igp": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"metric_out_minimum_igp_offset": schema.Int64Attribute{
						Optional: true,
					},
					"mtu_discovery": schema.BoolAttribute{
						Optional: true,
					},
					"multihop": schema.BoolAttribute{
						Optional: true,
					},
					"out_delay": schema.Int64Attribute{
						Optional: true,
					},
					"passive": schema.BoolAttribute{
						Optional: true,
					},
					"peer_as": schema.StringAttribute{
						Optional: true,
					},
					"preference": schema.Int64Attribute{
						Optional: true,
					},
					"remove_private": schema.BoolAttribute{
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
								"holddown_interval": schema.Int64Attribute{
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
								"session_mode": schema.StringAttribute{
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
					"bgp_multipath": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"allow_protection": schema.BoolAttribute{
									Optional: true,
								},
								"disable": schema.BoolAttribute{
									Optional: true,
								},
								"multiple_as": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"family_evpn": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"nlri_type": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
							},
							Blocks: map[string]schema.Block{
								"accepted_prefix_limit": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: bgpBlockFamilyBlockPrefixLimit{}.resourceSchemaAttributes(),
									},
								},
								"prefix_limit": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: bgpBlockFamilyBlockPrefixLimit{}.resourceSchemaAttributes(),
									},
								},
							},
						},
					},
					"family_inet": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"nlri_type": schema.StringAttribute{
									Required: true,
								},
							},
							Blocks: map[string]schema.Block{
								"accepted_prefix_limit": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: bgpBlockFamilyBlockPrefixLimit{}.resourceSchemaAttributes(),
									},
								},
								"prefix_limit": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: bgpBlockFamilyBlockPrefixLimit{}.resourceSchemaAttributes(),
									},
								},
							},
						},
					},
					"family_inet6": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"nlri_type": schema.StringAttribute{
									Required: true,
								},
							},
							Blocks: map[string]schema.Block{
								"accepted_prefix_limit": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: bgpBlockFamilyBlockPrefixLimit{}.resourceSchemaAttributes(),
									},
								},
								"prefix_limit": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: bgpBlockFamilyBlockPrefixLimit{}.resourceSchemaAttributes(),
									},
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
								"restart_time": schema.Int64Attribute{
									Optional: true,
								},
								"stale_route_time": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeBgpNeighborStateV0toV1,
		},
	}
}

func upgradeBgpNeighborStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                           types.String   `tfsdk:"id"`
		IP                           types.String   `tfsdk:"ip"`
		RoutingInstance              types.String   `tfsdk:"routing_instance"`
		Group                        types.String   `tfsdk:"group"`
		AcceptRemoteNexthop          types.Bool     `tfsdk:"accept_remote_nexthop"`
		AdvertiseExternal            types.Bool     `tfsdk:"advertise_external"`
		AdvertiseExternalConditional types.Bool     `tfsdk:"advertise_external_conditional"`
		AdvertiseInactive            types.Bool     `tfsdk:"advertise_inactive"`
		AdvertisePeerAS              types.Bool     `tfsdk:"advertise_peer_as"`
		NoAdvertisePeerAS            types.Bool     `tfsdk:"no_advertise_peer_as"`
		ASOverride                   types.Bool     `tfsdk:"as_override"`
		AuthenticationAlgorithm      types.String   `tfsdk:"authentication_algorithm"`
		AuthenticationKey            types.String   `tfsdk:"authentication_key"`
		AuthenticationKeyChain       types.String   `tfsdk:"authentication_key_chain"`
		Cluster                      types.String   `tfsdk:"cluster"`
		Damping                      types.Bool     `tfsdk:"damping"`
		Export                       []types.String `tfsdk:"export"`
		HoldTime                     types.Int64    `tfsdk:"hold_time"`
		Import                       []types.String `tfsdk:"import"`
		KeepAll                      types.Bool     `tfsdk:"keep_all"`
		KeepNone                     types.Bool     `tfsdk:"keep_none"`
		LocalAddress                 types.String   `tfsdk:"local_address"`
		LocalAS                      types.String   `tfsdk:"local_as"`
		LocalASAlias                 types.Bool     `tfsdk:"local_as_alias"`
		LocalASLoops                 types.Int64    `tfsdk:"local_as_loops"`
		LocalASNoPrependGlobalAS     types.Bool     `tfsdk:"local_as_no_prepend_global_as"`
		LocalASPrivate               types.Bool     `tfsdk:"local_as_private"`
		LocalInterface               types.String   `tfsdk:"local_interface"`
		LocalPreference              types.Int64    `tfsdk:"local_preference"`
		LogUpdown                    types.Bool     `tfsdk:"log_updown"`
		MetricOut                    types.Int64    `tfsdk:"metric_out"`
		MetricOutIgp                 types.Bool     `tfsdk:"metric_out_igp"`
		MetricOutIgpDelayMedUpdate   types.Bool     `tfsdk:"metric_out_igp_delay_med_update"`
		MetricOutIgpOffset           types.Int64    `tfsdk:"metric_out_igp_offset"`
		MetricOutMinimumIgp          types.Bool     `tfsdk:"metric_out_minimum_igp"`
		MetricOutMinimumIgpOffset    types.Int64    `tfsdk:"metric_out_minimum_igp_offset"`
		MtuDiscovery                 types.Bool     `tfsdk:"mtu_discovery"`
		Multihop                     types.Bool     `tfsdk:"multihop"`
		OutDelay                     types.Int64    `tfsdk:"out_delay"`
		Passive                      types.Bool     `tfsdk:"passive"`
		PeerAS                       types.String   `tfsdk:"peer_as"`
		Preference                   types.Int64    `tfsdk:"preference"`
		RemovePrivate                types.Bool     `tfsdk:"remove_private"`
		BfdLivenessDetection         []struct {
			AuthenticationAlgorithm         types.String `tfsdk:"authentication_algorithm"`
			AuthenticationKeyChain          types.String `tfsdk:"authentication_key_chain"`
			AuthenticationLooseCheck        types.Bool   `tfsdk:"authentication_loose_check"`
			DetectionTimeThreshold          types.Int64  `tfsdk:"detection_time_threshold"`
			HolddownInterval                types.Int64  `tfsdk:"holddown_interval"`
			MinimumInterval                 types.Int64  `tfsdk:"minimum_interval"`
			MinimumReceiveInterval          types.Int64  `tfsdk:"minimum_receive_interval"`
			Multiplier                      types.Int64  `tfsdk:"multiplier"`
			SessionMode                     types.String `tfsdk:"session_mode"`
			TransmitIntervalMinimumInterval types.Int64  `tfsdk:"transmit_interval_minimum_interval"`
			TransmitIntervalThreshold       types.Int64  `tfsdk:"transmit_interval_threshold"`
			Version                         types.String `tfsdk:"version"`
		} `tfsdk:"bfd_liveness_detection"`
		BgpMultipath []struct {
			AllowProtection types.Bool `tfsdk:"allow_protection"`
			Disable         types.Bool `tfsdk:"disable"`
			MultipleAS      types.Bool `tfsdk:"multiple_as"`
		} `tfsdk:"bgp_multipath"`
		FamilyEvpn []struct {
			NlriType            types.String                     `tfsdk:"nlri_type"`
			AcceptedPrefixLimit []bgpBlockFamilyBlockPrefixLimit `tfsdk:"accepted_prefix_limit"`
			PrefixLimit         []bgpBlockFamilyBlockPrefixLimit `tfsdk:"prefix_limit"`
		} `tfsdk:"family_evpn"`
		FamilyInet []struct {
			NlriType            types.String                     `tfsdk:"nlri_type"`
			AcceptedPrefixLimit []bgpBlockFamilyBlockPrefixLimit `tfsdk:"accepted_prefix_limit"`
			PrefixLimit         []bgpBlockFamilyBlockPrefixLimit `tfsdk:"prefix_limit"`
		} `tfsdk:"family_inet"`
		FamilyInet6 []struct {
			NlriType            types.String                     `tfsdk:"nlri_type"`
			AcceptedPrefixLimit []bgpBlockFamilyBlockPrefixLimit `tfsdk:"accepted_prefix_limit"`
			PrefixLimit         []bgpBlockFamilyBlockPrefixLimit `tfsdk:"prefix_limit"`
		} `tfsdk:"family_inet6"`
		GracefulRestart []struct {
			Disable        types.Bool  `tfsdk:"disable"`
			RestartTime    types.Int64 `tfsdk:"restart_time"`
			StaleRouteTime types.Int64 `tfsdk:"stale_route_time"`
		} `tfsdk:"graceful_restart"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 bgpNeighborData
	dataV1.ID = dataV0.ID
	dataV1.IP = dataV0.IP
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.Group = dataV0.Group
	dataV1.AcceptRemoteNexthop = dataV0.AcceptRemoteNexthop
	dataV1.AdvertiseExternal = dataV0.AdvertiseExternal
	dataV1.AdvertiseExternalConditional = dataV0.AdvertiseExternalConditional
	dataV1.AdvertiseInactive = dataV0.AdvertiseInactive
	dataV1.AdvertisePeerAS = dataV0.AdvertisePeerAS
	dataV1.NoAdvertisePeerAS = dataV0.NoAdvertisePeerAS
	dataV1.ASOverride = dataV0.ASOverride
	dataV1.AuthenticationAlgorithm = dataV0.AuthenticationAlgorithm
	dataV1.AuthenticationKey = dataV0.AuthenticationKey
	dataV1.AuthenticationKeyChain = dataV0.AuthenticationKeyChain
	dataV1.Cluster = dataV0.Cluster
	dataV1.Damping = dataV0.Damping
	dataV1.Export = dataV0.Export
	dataV1.HoldTime = dataV0.HoldTime
	dataV1.Import = dataV0.Import
	dataV1.KeepAll = dataV0.KeepAll
	dataV1.KeepNone = dataV0.KeepNone
	dataV1.LocalAddress = dataV0.LocalAddress
	dataV1.LocalAS = dataV0.LocalAS
	dataV1.LocalASAlias = dataV0.LocalASAlias
	dataV1.LocalASLoops = dataV0.LocalASLoops
	dataV1.LocalASNoPrependGlobalAS = dataV0.LocalASNoPrependGlobalAS
	dataV1.LocalASPrivate = dataV0.LocalASPrivate
	dataV1.LocalInterface = dataV0.LocalInterface
	dataV1.LocalPreference = dataV0.LocalPreference
	dataV1.LogUpdown = dataV0.LogUpdown
	dataV1.MetricOut = dataV0.MetricOut
	dataV1.MetricOutIgp = dataV0.MetricOutIgp
	dataV1.MetricOutIgpDelayMedUpdate = dataV0.MetricOutIgpDelayMedUpdate
	dataV1.MetricOutIgpOffset = dataV0.MetricOutIgpOffset
	dataV1.MetricOutMinimumIgp = dataV0.MetricOutMinimumIgp
	dataV1.MetricOutMinimumIgpOffset = dataV0.MetricOutMinimumIgpOffset
	dataV1.MtuDiscovery = dataV0.MtuDiscovery
	dataV1.Multihop = dataV0.Multihop
	dataV1.OutDelay = dataV0.OutDelay
	dataV1.Passive = dataV0.Passive
	dataV1.PeerAS = dataV0.PeerAS
	dataV1.Preference = dataV0.Preference
	dataV1.RemovePrivate = dataV0.RemovePrivate
	if len(dataV0.BfdLivenessDetection) > 0 {
		dataV1.BfdLivenessDetection = &bgpBlockBfdLivenessDetection{
			AuthenticationLooseCheck:        dataV0.BfdLivenessDetection[0].AuthenticationLooseCheck,
			AuthenticationAlgorithm:         dataV0.BfdLivenessDetection[0].AuthenticationAlgorithm,
			AuthenticationKeyChain:          dataV0.BfdLivenessDetection[0].AuthenticationKeyChain,
			DetectionTimeThreshold:          dataV0.BfdLivenessDetection[0].DetectionTimeThreshold,
			HolddownInterval:                dataV0.BfdLivenessDetection[0].HolddownInterval,
			MinimumInterval:                 dataV0.BfdLivenessDetection[0].MinimumInterval,
			MinimumReceiveInterval:          dataV0.BfdLivenessDetection[0].MinimumReceiveInterval,
			Multiplier:                      dataV0.BfdLivenessDetection[0].Multiplier,
			SessionMode:                     dataV0.BfdLivenessDetection[0].SessionMode,
			TransmitIntervalMinimumInterval: dataV0.BfdLivenessDetection[0].TransmitIntervalMinimumInterval,
			TransmitIntervalThreshold:       dataV0.BfdLivenessDetection[0].TransmitIntervalThreshold,
			Version:                         dataV0.BfdLivenessDetection[0].Version,
		}
	}
	if len(dataV0.BgpMultipath) > 0 {
		dataV1.BgpMultipath = &bgpBlockBgpMultipath{
			AllowProtection: dataV0.BgpMultipath[0].AllowProtection,
			Disable:         dataV0.BgpMultipath[0].Disable,
			MultipleAS:      dataV0.BgpMultipath[0].MultipleAS,
		}
	}
	for _, blockV0 := range dataV0.FamilyEvpn {
		blockV1 := bgpBlockFamily{
			NlriType: blockV0.NlriType,
		}
		if len(blockV0.AcceptedPrefixLimit) > 0 {
			blockV1.AcceptedPrefixLimit = &blockV0.AcceptedPrefixLimit[0]
		}
		if len(blockV0.PrefixLimit) > 0 {
			blockV1.PrefixLimit = &blockV0.PrefixLimit[0]
		}
		dataV1.FamilyEvpn = append(dataV1.FamilyEvpn, blockV1)
	}
	for _, blockV0 := range dataV0.FamilyInet {
		blockV1 := bgpBlockFamily{
			NlriType: blockV0.NlriType,
		}
		if len(blockV0.AcceptedPrefixLimit) > 0 {
			blockV1.AcceptedPrefixLimit = &blockV0.AcceptedPrefixLimit[0]
		}
		if len(blockV0.PrefixLimit) > 0 {
			blockV1.PrefixLimit = &blockV0.PrefixLimit[0]
		}
		dataV1.FamilyInet = append(dataV1.FamilyInet, blockV1)
	}
	for _, blockV0 := range dataV0.FamilyInet6 {
		blockV1 := bgpBlockFamily{
			NlriType: blockV0.NlriType,
		}
		if len(blockV0.AcceptedPrefixLimit) > 0 {
			blockV1.AcceptedPrefixLimit = &blockV0.AcceptedPrefixLimit[0]
		}
		if len(blockV0.PrefixLimit) > 0 {
			blockV1.PrefixLimit = &blockV0.PrefixLimit[0]
		}
		dataV1.FamilyInet6 = append(dataV1.FamilyInet6, blockV1)
	}
	if len(dataV0.GracefulRestart) > 0 {
		dataV1.GracefulRestart = &bgpBlockGracefulRestart{
			Disable:        dataV0.GracefulRestart[0].Disable,
			RestartTime:    dataV0.GracefulRestart[0].RestartTime,
			StaleRouteTime: dataV0.GracefulRestart[0].StaleRouteTime,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
