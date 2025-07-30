package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityIdpPolicy) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
				},
				Blocks: map[string]schema.Block{
					"exempt_rule": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"description": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"match": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"custom_attack": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"custom_attack_group": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_address": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_address_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"dynamic_attack_group": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"from_zone": schema.StringAttribute{
												Optional: true,
											},
											"predefined_attack": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"predefined_attack_group": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_address": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_address_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"to_zone": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"ips_rule": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"description": schema.StringAttribute{
									Optional: true,
								},
								"terminal": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"match": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"application": schema.StringAttribute{
												Optional: true,
											},
											"custom_attack": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"custom_attack_group": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_address": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"destination_address_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"dynamic_attack_group": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"from_zone": schema.StringAttribute{
												Optional: true,
											},
											"predefined_attack": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"predefined_attack_group": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_address": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"source_address_except": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"to_zone": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"then": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.StringAttribute{
												Required: true,
											},
											"action_cos_forwarding_class": schema.StringAttribute{
												Optional: true,
											},
											"action_dscp_code_point": schema.Int64Attribute{
												Optional: true,
											},
											"ip_action": schema.StringAttribute{
												Optional: true,
											},
											"ip_action_log": schema.BoolAttribute{
												Optional: true,
											},
											"ip_action_log_create": schema.BoolAttribute{
												Optional: true,
											},
											"ip_action_refresh_timeout": schema.BoolAttribute{
												Optional: true,
											},
											"ip_action_target": schema.StringAttribute{
												Optional: true,
											},
											"ip_action_timeout": schema.Int64Attribute{
												Optional: true,
											},
											"notification": schema.BoolAttribute{
												Optional: true,
											},
											"notification_log_attacks": schema.BoolAttribute{
												Optional: true,
											},
											"notification_log_attacks_alert": schema.BoolAttribute{
												Optional: true,
											},
											"notification_packet_log": schema.BoolAttribute{
												Optional: true,
											},
											"notification_packet_log_post_attack": schema.Int64Attribute{
												Optional: true,
											},
											"notification_packet_log_post_attack_timeout": schema.Int64Attribute{
												Optional: true,
											},
											"notification_packet_log_pre_attack": schema.Int64Attribute{
												Optional: true,
											},
											"severity": schema.StringAttribute{
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
			StateUpgrader: upgradeSecurityIdpPolicyV0toV1,
		},
	}
}

func upgradeSecurityIdpPolicyV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID         types.String `tfsdk:"id"`
		Name       types.String `tfsdk:"name"`
		ExemptRule []struct {
			Name        types.String `tfsdk:"name"`
			Description types.String `tfsdk:"description"`
			Match       []struct {
				CustomAttack             []types.String `tfsdk:"custom_attack"`
				CustomAttackGroup        []types.String `tfsdk:"custom_attack_group"`
				DestinationAddress       []types.String `tfsdk:"destination_address"`
				DestinationAddressExcept []types.String `tfsdk:"destination_address_except"`
				DynamicAttackGroup       []types.String `tfsdk:"dynamic_attack_group"`
				FromZone                 types.String   `tfsdk:"from_zone"`
				PredefinedAttack         []types.String `tfsdk:"predefined_attack"`
				PredefinedAttackGroup    []types.String `tfsdk:"predefined_attack_group"`
				SourceAddress            []types.String `tfsdk:"source_address"`
				SourceAddressExcept      []types.String `tfsdk:"source_address_except"`
				ToZone                   types.String   `tfsdk:"to_zone"`
			} `tfsdk:"match"`
		} `tfsdk:"exempt_rule"`
		IpsRule []struct {
			Name        types.String `tfsdk:"name"`
			Description types.String `tfsdk:"description"`
			Terminal    types.Bool   `tfsdk:"terminal"`
			Match       []struct {
				Application              types.String   `tfsdk:"application"`
				CustomAttack             []types.String `tfsdk:"custom_attack"`
				CustomAttackGroup        []types.String `tfsdk:"custom_attack_group"`
				DestinationAddress       []types.String `tfsdk:"destination_address"`
				DestinationAddressExcept []types.String `tfsdk:"destination_address_except"`
				DynamicAttackGroup       []types.String `tfsdk:"dynamic_attack_group"`
				FromZone                 types.String   `tfsdk:"from_zone"`
				PredefinedAttack         []types.String `tfsdk:"predefined_attack"`
				PredefinedAttackGroup    []types.String `tfsdk:"predefined_attack_group"`
				SourceAddress            []types.String `tfsdk:"source_address"`
				SourceAddressExcept      []types.String `tfsdk:"source_address_except"`
				ToZone                   types.String   `tfsdk:"to_zone"`
			} `tfsdk:"match"`
			Then []struct {
				Action                                 types.String `tfsdk:"action"`
				ActionCosForwardingClass               types.String `tfsdk:"action_cos_forwarding_class"`
				ActionDscpCodePoint                    types.Int64  `tfsdk:"action_dscp_code_point"`
				IPAction                               types.String `tfsdk:"ip_action"`
				IPActionLog                            types.Bool   `tfsdk:"ip_action_log"`
				IPActionLogCreate                      types.Bool   `tfsdk:"ip_action_log_create"`
				IPActionRefreshTimeout                 types.Bool   `tfsdk:"ip_action_refresh_timeout"`
				IPActionTarget                         types.String `tfsdk:"ip_action_target"`
				IPActionTimeout                        types.Int64  `tfsdk:"ip_action_timeout"`
				Notification                           types.Bool   `tfsdk:"notification"`
				NotificationLogAttacks                 types.Bool   `tfsdk:"notification_log_attacks"`
				NotificationLogAttacksAlert            types.Bool   `tfsdk:"notification_log_attacks_alert"`
				NotificationPacketLog                  types.Bool   `tfsdk:"notification_packet_log"`
				NotificationPacketLogPostAttack        types.Int64  `tfsdk:"notification_packet_log_post_attack"`
				NotificationPacketLogPostAttackTimeout types.Int64  `tfsdk:"notification_packet_log_post_attack_timeout"`
				NotificationPacketLogPreAttack         types.Int64  `tfsdk:"notification_packet_log_pre_attack"`
				Severity                               types.String `tfsdk:"severity"`
			} `tfsdk:"then"`
		} `tfsdk:"ips_rule"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityIdpPolicyData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	for _, blockV0 := range dataV0.ExemptRule {
		blockV1 := securityIdpPolicyBlockExemptRule{
			Name:        blockV0.Name,
			Description: blockV0.Description,
		}
		if len(blockV0.Match) > 0 {
			blockV1.Match = &securityIdpPolicyBlockExemptRuleBlockMatch{
				CustomAttack:             blockV0.Match[0].CustomAttack,
				CustomAttackGroup:        blockV0.Match[0].CustomAttackGroup,
				DestinationAddress:       blockV0.Match[0].DestinationAddress,
				DestinationAddressExcept: blockV0.Match[0].DestinationAddressExcept,
				DynamicAttackGroup:       blockV0.Match[0].DynamicAttackGroup,
				FromZone:                 blockV0.Match[0].FromZone,
				PredefinedAttack:         blockV0.Match[0].PredefinedAttack,
				PredefinedAttackGroup:    blockV0.Match[0].PredefinedAttackGroup,
				SourceAddress:            blockV0.Match[0].SourceAddress,
				SourceAddressExcept:      blockV0.Match[0].SourceAddressExcept,
				ToZone:                   blockV0.Match[0].ToZone,
			}
		}
		dataV1.ExemptRule = append(dataV1.ExemptRule, blockV1)
	}
	for _, blockV0 := range dataV0.IpsRule {
		blockV1 := securityIdpPolicyBlockIpsRule{
			Name:        blockV0.Name,
			Description: blockV0.Description,
			Terminal:    blockV0.Terminal,
		}
		if len(blockV0.Match) > 0 {
			blockV1.Match = &securityIdpPolicyBlockIpsRuleBlockMatch{
				Application: blockV0.Match[0].Application,
				securityIdpPolicyBlockExemptRuleBlockMatch: securityIdpPolicyBlockExemptRuleBlockMatch{
					CustomAttack:             blockV0.Match[0].CustomAttack,
					CustomAttackGroup:        blockV0.Match[0].CustomAttackGroup,
					DestinationAddress:       blockV0.Match[0].DestinationAddress,
					DestinationAddressExcept: blockV0.Match[0].DestinationAddressExcept,
					DynamicAttackGroup:       blockV0.Match[0].DynamicAttackGroup,
					FromZone:                 blockV0.Match[0].FromZone,
					PredefinedAttack:         blockV0.Match[0].PredefinedAttack,
					PredefinedAttackGroup:    blockV0.Match[0].PredefinedAttackGroup,
					SourceAddress:            blockV0.Match[0].SourceAddress,
					SourceAddressExcept:      blockV0.Match[0].SourceAddressExcept,
					ToZone:                   blockV0.Match[0].ToZone,
				},
			}
		}
		if len(blockV0.Then) > 0 {
			blockV1.Then = &securityIdpPolicyBlockIpsRuleBlockThen{
				Action:                                 blockV0.Then[0].Action,
				ActionCosForwardingClass:               blockV0.Then[0].ActionCosForwardingClass,
				ActionDscpCodePoint:                    blockV0.Then[0].ActionDscpCodePoint,
				IPAction:                               blockV0.Then[0].IPAction,
				IPActionLog:                            blockV0.Then[0].IPActionLog,
				IPActionLogCreate:                      blockV0.Then[0].IPActionLogCreate,
				IPActionRefreshTimeout:                 blockV0.Then[0].IPActionRefreshTimeout,
				IPActionTarget:                         blockV0.Then[0].IPActionTarget,
				IPActionTimeout:                        blockV0.Then[0].IPActionTimeout,
				Notification:                           blockV0.Then[0].Notification,
				NotificationLogAttacks:                 blockV0.Then[0].NotificationLogAttacks,
				NotificationLogAttacksAlert:            blockV0.Then[0].NotificationLogAttacksAlert,
				NotificationPacketLog:                  blockV0.Then[0].NotificationPacketLog,
				NotificationPacketLogPostAttack:        blockV0.Then[0].NotificationPacketLogPostAttack,
				NotificationPacketLogPostAttackTimeout: blockV0.Then[0].NotificationPacketLogPostAttackTimeout,
				NotificationPacketLogPreAttack:         blockV0.Then[0].NotificationPacketLogPreAttack,
				Severity:                               blockV0.Then[0].Severity,
			}
		}
		dataV1.IpsRule = append(dataV1.IpsRule, blockV1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
