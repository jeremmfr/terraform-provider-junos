package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *ospf) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"version": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"routing_instance": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"disable": schema.BoolAttribute{
						Optional: true,
					},
					"domain_id": schema.StringAttribute{
						Optional: true,
					},
					"export": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"external_preference": schema.Int64Attribute{
						Optional: true,
					},
					"forwarding_address_to_broadcast": schema.BoolAttribute{
						Optional: true,
					},
					"import": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"labeled_preference": schema.Int64Attribute{
						Optional: true,
					},
					"lsa_refresh_interval": schema.Int64Attribute{
						Optional: true,
					},
					"no_nssa_abr": schema.BoolAttribute{
						Optional: true,
					},
					"no_rfc1583": schema.BoolAttribute{
						Optional: true,
					},
					"preference": schema.Int64Attribute{
						Optional: true,
					},
					"prefix_export_limit": schema.Int64Attribute{
						Optional: true,
					},
					"reference_bandwidth": schema.StringAttribute{
						Optional: true,
					},
					"rib_group": schema.StringAttribute{
						Optional: true,
					},
					"sham_link": schema.BoolAttribute{
						Optional: true,
					},
					"sham_link_local": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"database_protection": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"maximum_lsa": schema.Int64Attribute{
									Required: true,
								},
								"ignore_count": schema.Int64Attribute{
									Optional: true,
								},
								"ignore_time": schema.Int64Attribute{
									Optional: true,
								},
								"reset_time": schema.Int64Attribute{
									Optional: true,
								},
								"warning_only": schema.BoolAttribute{
									Optional: true,
								},
								"warning_threshold": schema.Int64Attribute{
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
								"helper_disable": schema.BoolAttribute{
									Optional: true,
								},
								"helper_disable_type": schema.StringAttribute{
									Optional: true,
								},
								"no_strict_lsa_checking": schema.BoolAttribute{
									Optional: true,
								},
								"notify_duration": schema.Int64Attribute{
									Optional: true,
								},
								"restart_duration": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"overload": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"allow_route_leaking": schema.BoolAttribute{
									Optional: true,
								},
								"as_external": schema.BoolAttribute{
									Optional: true,
								},
								"stub_network": schema.BoolAttribute{
									Optional: true,
								},
								"timeout": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"spf_options": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"delay": schema.Int64Attribute{
									Optional: true,
								},
								"holddown": schema.Int64Attribute{
									Optional: true,
								},
								"no_ignore_our_externals": schema.BoolAttribute{
									Optional: true,
								},
								"rapid_runs": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeOspfV0toV1,
		},
	}
}

func upgradeOspfV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                           types.String   `tfsdk:"id"`
		Version                      types.String   `tfsdk:"version"`
		RoutingInstance              types.String   `tfsdk:"routing_instance"`
		Disable                      types.Bool     `tfsdk:"disable"`
		DomainID                     types.String   `tfsdk:"domain_id"`
		Export                       []types.String `tfsdk:"export"`
		ExternalPreference           types.Int64    `tfsdk:"external_preference"`
		ForwardingAddressToBroadcast types.Bool     `tfsdk:"forwarding_address_to_broadcast"`
		Import                       []types.String `tfsdk:"import"`
		LabeledPreference            types.Int64    `tfsdk:"labeled_preference"`
		LsaRefreshInterval           types.Int64    `tfsdk:"lsa_refresh_interval"`
		NoNssaAbr                    types.Bool     `tfsdk:"no_nssa_abr"`
		NoRfc1583                    types.Bool     `tfsdk:"no_rfc1583"`
		Preference                   types.Int64    `tfsdk:"preference"`
		PrefixExportLimit            types.Int64    `tfsdk:"prefix_export_limit"`
		ReferenceBandwidth           types.String   `tfsdk:"reference_bandwidth"`
		RibGroup                     types.String   `tfsdk:"rib_group"`
		ShamLink                     types.Bool     `tfsdk:"sham_link"`
		ShamLinkLocal                types.String   `tfsdk:"sham_link_local"`
		DatabaseProtection           []struct {
			MaximumLsa       types.Int64 `tfsdk:"maximum_lsa"`
			IgnoreCount      types.Int64 `tfsdk:"ignore_count"`
			IgnoreTime       types.Int64 `tfsdk:"ignore_time"`
			ResetTime        types.Int64 `tfsdk:"reset_time"`
			WarningOnly      types.Bool  `tfsdk:"warning_only"`
			WarningThreshold types.Int64 `tfsdk:"warning_threshold"`
		} `tfsdk:"database_protection"`
		GracefulRestart []struct {
			Disable             types.Bool   `tfsdk:"disable"`
			HelperDisable       types.Bool   `tfsdk:"helper_disable"`
			HelperDisableType   types.String `tfsdk:"helper_disable_type"`
			NoStrictLsaChecking types.Bool   `tfsdk:"no_strict_lsa_checking"`
			NotifyDuration      types.Int64  `tfsdk:"notify_duration"`
			RestartDuration     types.Int64  `tfsdk:"restart_duration"`
		} `tfsdk:"graceful_restart"`
		Overload []struct {
			AllowRouteLeaking types.Bool  `tfsdk:"allow_route_leaking"`
			ASExternal        types.Bool  `tfsdk:"as_external"`
			StubNetwork       types.Bool  `tfsdk:"stub_network"`
			Timeout           types.Int64 `tfsdk:"timeout"`
		} `tfsdk:"overload"`
		SpfOptions []struct {
			Delay                types.Int64 `tfsdk:"delay"`
			Holddown             types.Int64 `tfsdk:"holddown"`
			NoIgnoreOurExternals types.Bool  `tfsdk:"no_ignore_our_externals"`
			RapidRuns            types.Int64 `tfsdk:"rapid_runs"`
		} `tfsdk:"spf_options"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 ospfData
	dataV1.ID = dataV0.ID
	dataV1.Version = dataV0.Version
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.Disable = dataV0.Disable
	dataV1.DomainID = dataV0.DomainID
	dataV1.Export = dataV0.Export
	dataV1.ExternalPreference = dataV0.ExternalPreference
	dataV1.ForwardingAddressToBroadcast = dataV0.ForwardingAddressToBroadcast
	dataV1.Import = dataV0.Import
	dataV1.LabeledPreference = dataV0.LabeledPreference
	dataV1.LsaRefreshInterval = dataV0.LsaRefreshInterval
	dataV1.NoNssaAbr = dataV0.NoNssaAbr
	dataV1.NoRfc1583 = dataV0.NoRfc1583
	dataV1.Preference = dataV0.Preference
	dataV1.PrefixExportLimit = dataV0.PrefixExportLimit
	dataV1.ReferenceBandwidth = dataV0.ReferenceBandwidth
	dataV1.RibGroup = dataV0.RibGroup
	dataV1.ShamLink = dataV0.ShamLink
	dataV1.ShamLinkLocal = dataV0.ShamLinkLocal
	if len(dataV0.DatabaseProtection) > 0 {
		dataV1.DatabaseProtection = &ospfBlockDatabaseProtection{
			MaximumLsa:       dataV0.DatabaseProtection[0].MaximumLsa,
			IgnoreCount:      dataV0.DatabaseProtection[0].IgnoreCount,
			IgnoreTime:       dataV0.DatabaseProtection[0].IgnoreTime,
			ResetTime:        dataV0.DatabaseProtection[0].ResetTime,
			WarningOnly:      dataV0.DatabaseProtection[0].WarningOnly,
			WarningThreshold: dataV0.DatabaseProtection[0].WarningThreshold,
		}
	}
	if len(dataV0.GracefulRestart) > 0 {
		dataV1.GracefulRestart = &ospfBlockGracefulRestart{
			Disable:             dataV0.GracefulRestart[0].Disable,
			HelperDisable:       dataV0.GracefulRestart[0].HelperDisable,
			HelperDisableType:   dataV0.GracefulRestart[0].HelperDisableType,
			NoStrictLsaChecking: dataV0.GracefulRestart[0].NoStrictLsaChecking,
			NotifyDuration:      dataV0.GracefulRestart[0].NotifyDuration,
			RestartDuration:     dataV0.GracefulRestart[0].RestartDuration,
		}
	}
	if len(dataV0.Overload) > 0 {
		dataV1.Overload = &ospfBlockOverload{
			AllowRouteLeaking: dataV0.Overload[0].AllowRouteLeaking,
			ASExternal:        dataV0.Overload[0].ASExternal,
			StubNetwork:       dataV0.Overload[0].StubNetwork,
			Timeout:           dataV0.Overload[0].Timeout,
		}
	}
	if len(dataV0.SpfOptions) > 0 {
		dataV1.SpfOptions = &ospfBlockSpfOptions{
			Delay:                dataV0.SpfOptions[0].Delay,
			Holddown:             dataV0.SpfOptions[0].Holddown,
			NoIgnoreOurExternals: dataV0.SpfOptions[0].NoIgnoreOurExternals,
			RapidRuns:            dataV0.SpfOptions[0].RapidRuns,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
