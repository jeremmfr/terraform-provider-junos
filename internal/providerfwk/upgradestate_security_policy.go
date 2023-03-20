package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
			StateUpgrader: upgradesecurityPolicyV0toV1,
		},
	}
}

//nolint:lll
func upgradesecurityPolicyV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID       types.String `tfsdk:"id"`
		FromZone types.String `tfsdk:"from_zone"`
		ToZone   types.String `tfsdk:"to_zone"`
		Policy   []struct {
			Count                           types.Bool     `tfsdk:"count"`
			LogInit                         types.Bool     `tfsdk:"log_init"`
			LogClose                        types.Bool     `tfsdk:"log_close"`
			MatchDestinationAddressExcluded types.Bool     `tfsdk:"match_destination_address_excluded"`
			MatchSourceAddressExcluded      types.Bool     `tfsdk:"match_source_address_excluded"`
			Name                            types.String   `tfsdk:"name"`
			Then                            types.String   `tfsdk:"then"`
			MatchSourceEndUserProfile       types.String   `tfsdk:"match_source_end_user_profile"`
			PermitTunnelIpsecVpn            types.String   `tfsdk:"permit_tunnel_ipsec_vpn"`
			MatchSourceAddress              []types.String `tfsdk:"match_source_address"`
			MatchDestinationAddress         []types.String `tfsdk:"match_destination_address"`
			MatchApplication                []types.String `tfsdk:"match_application"`
			MatchDynamicApplication         []types.String `tfsdk:"match_dynamic_application"`
			PermitApplicationServices       []struct {
				Idp                              types.Bool   `tfsdk:"idp"`
				RedirectWx                       types.Bool   `tfsdk:"redirect_wx"`
				ReverseRedirectWx                types.Bool   `tfsdk:"reverse_redirect_wx"`
				AdvancedAntiMalwarePolicy        types.String `tfsdk:"advanced_anti_malware_policy"`
				ApplicationFirewallRuleSet       types.String `tfsdk:"application_firewall_rule_set"`
				ApplicationTrafficControlRuleSet types.String `tfsdk:"application_traffic_control_rule_set"`
				GprsGtpProfile                   types.String `tfsdk:"gprs_gtp_profile"`
				GprsSctpProfile                  types.String `tfsdk:"gprs_sctp_profile"`
				IdpPolicy                        types.String `tfsdk:"idp_policy"`
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
		var blockV1 securityPolicyPolicy
		blockV1.Count = blockV0.Count
		blockV1.LogInit = blockV0.LogInit
		blockV1.LogClose = blockV0.LogClose
		blockV1.MatchDestinationAddressExcluded = blockV0.MatchDestinationAddressExcluded
		blockV1.MatchSourceAddressExcluded = blockV0.MatchSourceAddressExcluded
		blockV1.Name = blockV0.Name
		blockV1.Then = blockV0.Then
		blockV1.MatchSourceEndUserProfile = blockV0.MatchSourceEndUserProfile
		blockV1.PermitTunnelIpsecVpn = blockV0.PermitTunnelIpsecVpn
		blockV1.MatchSourceAddress = blockV0.MatchSourceAddress
		blockV1.MatchDestinationAddress = blockV0.MatchDestinationAddress
		blockV1.MatchApplication = blockV0.MatchApplication
		blockV1.MatchDynamicApplication = blockV0.MatchDynamicApplication
		if len(blockV0.PermitApplicationServices) > 0 {
			blockV1.PermitApplicationServices = &securityPolicyPolicyPermitApplicationServices{}

			blockV1.PermitApplicationServices.Idp = blockV0.PermitApplicationServices[0].Idp
			blockV1.PermitApplicationServices.RedirectWx = blockV0.PermitApplicationServices[0].RedirectWx
			blockV1.PermitApplicationServices.ReverseRedirectWx = blockV0.PermitApplicationServices[0].ReverseRedirectWx
			blockV1.PermitApplicationServices.AdvancedAntiMalwarePolicy = blockV0.PermitApplicationServices[0].AdvancedAntiMalwarePolicy
			blockV1.PermitApplicationServices.ApplicationFirewallRuleSet = blockV0.PermitApplicationServices[0].ApplicationFirewallRuleSet
			blockV1.PermitApplicationServices.ApplicationTrafficControlRuleSet = blockV0.PermitApplicationServices[0].ApplicationTrafficControlRuleSet
			blockV1.PermitApplicationServices.GprsGtpProfile = blockV0.PermitApplicationServices[0].GprsGtpProfile
			blockV1.PermitApplicationServices.GprsSctpProfile = blockV0.PermitApplicationServices[0].GprsSctpProfile
			blockV1.PermitApplicationServices.IdpPolicy = blockV0.PermitApplicationServices[0].IdpPolicy
			blockV1.PermitApplicationServices.SecurityIntelligencePolicy = blockV0.PermitApplicationServices[0].SecurityIntelligencePolicy
			if len(blockV0.PermitApplicationServices[0].SSLProxy) > 0 {
				blockV1.PermitApplicationServices.SSLProxy = &struct {
					ProfileName basetypes.StringValue `tfsdk:"profile_name"`
				}{
					ProfileName: blockV0.PermitApplicationServices[0].SSLProxy[0].ProfileName,
				}
			}
			if len(blockV0.PermitApplicationServices[0].UacPolicy) > 0 {
				blockV1.PermitApplicationServices.UacPolicy = &struct {
					CaptivePortal types.String `tfsdk:"captive_portal"`
				}{
					CaptivePortal: blockV0.PermitApplicationServices[0].UacPolicy[0].CaptivePortal,
				}
			}
			blockV1.PermitApplicationServices.UtmPolicy = blockV0.PermitApplicationServices[0].UtmPolicy
		}
		dataV1.Policy = append(dataV1.Policy, blockV1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
