package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityPolicy) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"from_zone": schema.StringAttribute{
						Required: true,
					},
					"to_zone": schema.StringAttribute{
						Required: true,
					},
				},
				Blocks: map[string]schema.Block{
					"policy": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"match_source_address": schema.SetAttribute{
									ElementType: types.StringType,
									Required:    true,
								},
								"match_destination_address": schema.SetAttribute{
									ElementType: types.StringType,
									Required:    true,
								},
								"then": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"count": schema.BoolAttribute{
									Optional: true,
								},
								"log_init": schema.BoolAttribute{
									Optional: true,
								},
								"log_close": schema.BoolAttribute{
									Optional: true,
								},
								"match_application": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"match_destination_address_excluded": schema.BoolAttribute{
									Optional: true,
								},
								"match_dynamic_application": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"match_source_address_excluded": schema.BoolAttribute{
									Optional: true,
								},
								"match_source_end_user_profile": schema.StringAttribute{
									Optional: true,
								},
								"permit_tunnel_ipsec_vpn": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"permit_application_services": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"advanced_anti_malware_policy": schema.StringAttribute{
												Optional: true,
											},
											"application_firewall_rule_set": schema.StringAttribute{
												Optional: true,
											},
											"application_traffic_control_rule_set": schema.StringAttribute{
												Optional: true,
											},
											"gprs_gtp_profile": schema.StringAttribute{
												Optional: true,
											},
											"gprs_sctp_profile": schema.StringAttribute{
												Optional: true,
											},
											"idp": schema.BoolAttribute{
												Optional: true,
											},
											"idp_policy": schema.StringAttribute{
												Optional: true,
											},
											"redirect_wx": schema.BoolAttribute{
												Optional: true,
											},
											"reverse_redirect_wx": schema.BoolAttribute{
												Optional: true,
											},
											"security_intelligence_policy": schema.StringAttribute{
												Optional: true,
											},
											"utm_policy": schema.StringAttribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"ssl_proxy": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"profile_name": schema.StringAttribute{
															Optional: true,
														},
													},
												},
											},
											"uac_policy": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"captive_portal": schema.StringAttribute{
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
					},
				},
			},
			StateUpgrader: upgradeSecurityPolicyV0toV1,
		},
	}
}

func upgradeSecurityPolicyV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID       types.String `tfsdk:"id"`
		FromZone types.String `tfsdk:"from_zone"`
		ToZone   types.String `tfsdk:"to_zone"`
		Policy   []struct {
			Name                            types.String   `tfsdk:"name"`
			MatchSourceAddress              []types.String `tfsdk:"match_source_address"`
			MatchDestinationAddress         []types.String `tfsdk:"match_destination_address"`
			Then                            types.String   `tfsdk:"then"`
			Count                           types.Bool     `tfsdk:"count"`
			LogInit                         types.Bool     `tfsdk:"log_init"`
			LogClose                        types.Bool     `tfsdk:"log_close"`
			MatchApplication                []types.String `tfsdk:"match_application"`
			MatchDestinationAddressExcluded types.Bool     `tfsdk:"match_destination_address_excluded"`
			MatchDynamicApplication         []types.String `tfsdk:"match_dynamic_application"`
			MatchSourceAddressExcluded      types.Bool     `tfsdk:"match_source_address_excluded"`
			MatchSourceEndUserProfile       types.String   `tfsdk:"match_source_end_user_profile"`
			PermitTunnelIpsecVpn            types.String   `tfsdk:"permit_tunnel_ipsec_vpn"`
			PermitApplicationServices       []struct {
				AdvancedAntiMalwarePolicy        types.String `tfsdk:"advanced_anti_malware_policy"`
				ApplicationFirewallRuleSet       types.String `tfsdk:"application_firewall_rule_set"`
				ApplicationTrafficControlRuleSet types.String `tfsdk:"application_traffic_control_rule_set"`
				GprsGtpProfile                   types.String `tfsdk:"gprs_gtp_profile"`
				GprsSctpProfile                  types.String `tfsdk:"gprs_sctp_profile"`
				Idp                              types.Bool   `tfsdk:"idp"`
				IdpPolicy                        types.String `tfsdk:"idp_policy"`
				RedirectWx                       types.Bool   `tfsdk:"redirect_wx"`
				ReverseRedirectWx                types.Bool   `tfsdk:"reverse_redirect_wx"`
				SecurityIntelligencePolicy       types.String `tfsdk:"security_intelligence_policy"`
				SSLProxy                         []struct {
					ProfileName types.String `tfsdk:"profile_name"`
				} `tfsdk:"ssl_proxy"`
				UacPolicy []struct {
					CaptivePortal types.String `tfsdk:"captive_portal"`
				} `tfsdk:"uac_policy"`
				UtmPolicy types.String `tfsdk:"utm_policy"`
			} `tfsdk:"permit_application_services"`
		} `tfsdk:"policy"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityPolicyData
	dataV1.ID = dataV0.ID
	dataV1.FromZone = dataV0.FromZone
	dataV1.ToZone = dataV0.ToZone
	for _, blockV0 := range dataV0.Policy {
		blockV1 := securityPolicyBlockPolicy{
			Count:                           blockV0.Count,
			LogInit:                         blockV0.LogInit,
			LogClose:                        blockV0.LogClose,
			MatchDestinationAddressExcluded: blockV0.MatchDestinationAddressExcluded,
			MatchSourceAddressExcluded:      blockV0.MatchSourceAddressExcluded,
			Name:                            blockV0.Name,
			Then:                            blockV0.Then,
			MatchSourceEndUserProfile:       blockV0.MatchSourceEndUserProfile,
			PermitTunnelIpsecVpn:            blockV0.PermitTunnelIpsecVpn,
			MatchSourceAddress:              blockV0.MatchSourceAddress,
			MatchDestinationAddress:         blockV0.MatchDestinationAddress,
			MatchApplication:                blockV0.MatchApplication,
			MatchDynamicApplication:         blockV0.MatchDynamicApplication,
		}
		if len(blockV0.PermitApplicationServices) > 0 {
			blockV1.PermitApplicationServices = &securityPolicyBlockPolicyBlockPermitApplicationServices{
				Idp:                              blockV0.PermitApplicationServices[0].Idp,
				RedirectWx:                       blockV0.PermitApplicationServices[0].RedirectWx,
				ReverseRedirectWx:                blockV0.PermitApplicationServices[0].ReverseRedirectWx,
				AdvancedAntiMalwarePolicy:        blockV0.PermitApplicationServices[0].AdvancedAntiMalwarePolicy,
				ApplicationFirewallRuleSet:       blockV0.PermitApplicationServices[0].ApplicationFirewallRuleSet,
				ApplicationTrafficControlRuleSet: blockV0.PermitApplicationServices[0].ApplicationTrafficControlRuleSet,
				GprsGtpProfile:                   blockV0.PermitApplicationServices[0].GprsGtpProfile,
				GprsSctpProfile:                  blockV0.PermitApplicationServices[0].GprsSctpProfile,
				IdpPolicy:                        blockV0.PermitApplicationServices[0].IdpPolicy,
				SecurityIntelligencePolicy:       blockV0.PermitApplicationServices[0].SecurityIntelligencePolicy,
				UtmPolicy:                        blockV0.PermitApplicationServices[0].UtmPolicy,
			}
			if len(blockV0.PermitApplicationServices[0].SSLProxy) > 0 {
				blockV1.PermitApplicationServices.SSLProxy = &blockV0.PermitApplicationServices[0].SSLProxy[0]
			}
			if len(blockV0.PermitApplicationServices[0].UacPolicy) > 0 {
				blockV1.PermitApplicationServices.UacPolicy = &blockV0.PermitApplicationServices[0].UacPolicy[0]
			}
		}
		dataV1.Policy = append(dataV1.Policy, blockV1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
