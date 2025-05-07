package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &securityIdpPolicy{}
	_ resource.ResourceWithConfigure      = &securityIdpPolicy{}
	_ resource.ResourceWithValidateConfig = &securityIdpPolicy{}
	_ resource.ResourceWithImportState    = &securityIdpPolicy{}
	_ resource.ResourceWithUpgradeState   = &securityIdpPolicy{}
)

type securityIdpPolicy struct {
	client *junos.Client
}

func newSecurityIdpPolicyResource() resource.Resource {
	return &securityIdpPolicy{}
}

func (rsc *securityIdpPolicy) typeName() string {
	return providerName + "_security_idp_policy"
}

func (rsc *securityIdpPolicy) junosName() string {
	return "security idp idp-policy"
}

func (rsc *securityIdpPolicy) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityIdpPolicy) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityIdpPolicy) Configure(
	ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*junos.Client)
	if !ok {
		unexpectedResourceConfigureType(ctx, req, resp)

		return
	}
	rsc.client = client
}

func (rsc *securityIdpPolicy) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "IDP policy name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"exempt_rule": schema.ListNestedBlock{
				Description: "For each name, configure exempt rule.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Rule name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 250),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Rule description.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"match": schema.SingleNestedBlock{
							Description: "Rule match criteria.",
							Attributes:  securityIdpPolicyBlockExemptRuleBlockMatch{}.attributesSchema(),
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
					},
				},
			},
			"ips_rule": schema.ListNestedBlock{
				Description: "For each name, configure IPS rule.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Rule name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 250),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Rule description.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"terminal": schema.BoolAttribute{
							Optional:    true,
							Description: "Set/Unset terminal flag.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"match": schema.SingleNestedBlock{
							Description: "Rule match criteria.",
							Attributes:  securityIdpPolicyBlockIpsRuleBlockMatch{}.attributesSchema(),
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
						"then": schema.SingleNestedBlock{
							Description: "eclare `then` configuration.",
							Attributes: map[string]schema.Attribute{
								"action": schema.StringAttribute{
									Required:    true,
									Description: "Action.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"class-of-service",
											"close-client",
											"close-client-and-server",
											"close-server",
											"drop-connection",
											"drop-packet",
											"ignore-connection",
											"mark-diffserv",
											"no-action",
											"recommended",
										),
									},
								},
								"action_cos_forwarding_class": schema.StringAttribute{
									Optional:    true,
									Description: "Forwarding class for outgoing packets.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 64),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"action_dscp_code_point": schema.Int64Attribute{
									Optional:    true,
									Description: "Codepoint value.",
									Validators: []validator.Int64{
										int64validator.Between(0, 63),
									},
								},
								"ip_action": schema.StringAttribute{
									Optional:    true,
									Description: "IP-action.",
									Validators: []validator.String{
										stringvalidator.OneOf("ip-block", "ip-close", "ip-notify"),
									},
								},
								"ip_action_log": schema.BoolAttribute{
									Optional:    true,
									Description: "Log IP action taken.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"ip_action_log_create": schema.BoolAttribute{
									Optional:    true,
									Description: "Log IP action creation.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"ip_action_refresh_timeout": schema.BoolAttribute{
									Optional:    true,
									Description: "Refresh timeout when future connections match installed ip-action filter.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"ip_action_target": schema.StringAttribute{
									Optional:    true,
									Description: "IP-action target.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"destination-address",
											"service",
											"source-address",
											"source-zone",
											"source-zone-address",
											"zone-service",
										),
									},
								},
								"ip_action_timeout": schema.Int64Attribute{
									Optional:    true,
									Description: "Number of seconds IP action should remain effective.",
									Validators: []validator.Int64{
										int64validator.Between(0, 64800),
									},
								},
								"notification": schema.BoolAttribute{
									Optional:    true,
									Description: "Configure notification.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"notification_log_attacks": schema.BoolAttribute{
									Optional:    true,
									Description: "Enable attack logging.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"notification_log_attacks_alert": schema.BoolAttribute{
									Optional:    true,
									Description: "Set alert flag in attack log.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"notification_packet_log": schema.BoolAttribute{
									Optional:    true,
									Description: "Enable packet-log.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"notification_packet_log_post_attack": schema.Int64Attribute{
									Optional:    true,
									Description: "No of packets to capture after attack.",
									Validators: []validator.Int64{
										int64validator.Between(0, 255),
									},
								},
								"notification_packet_log_post_attack_timeout": schema.Int64Attribute{
									Optional:    true,
									Description: "Timeout (seconds) after attack before stopping packet capture.",
									Validators: []validator.Int64{
										int64validator.Between(0, 1800),
									},
								},
								"notification_packet_log_pre_attack": schema.Int64Attribute{
									Optional:    true,
									Description: "No of packets to capture before attack.",
									Validators: []validator.Int64{
										int64validator.Between(1, 255),
									},
								},
								"severity": schema.StringAttribute{
									Optional:    true,
									Description: "Set rule severity level.",
									Validators: []validator.String{
										stringvalidator.OneOf("critical", "info", "major", "minor", "warning"),
									},
								},
							},
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
					},
				},
			},
		},
	}
}

type securityIdpPolicyData struct {
	ID         types.String                       `tfsdk:"id"`
	Name       types.String                       `tfsdk:"name"`
	ExemptRule []securityIdpPolicyBlockExemptRule `tfsdk:"exempt_rule"`
	IpsRule    []securityIdpPolicyBlockIpsRule    `tfsdk:"ips_rule"`
}

type securityIdpPolicyConfig struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ExemptRule types.List   `tfsdk:"exempt_rule"`
	IpsRule    types.List   `tfsdk:"ips_rule"`
}

type securityIdpPolicyBlockExemptRule struct {
	Name        types.String                                `tfsdk:"name"        tfdata:"identifier"`
	Description types.String                                `tfsdk:"description"`
	Match       *securityIdpPolicyBlockExemptRuleBlockMatch `tfsdk:"match"`
}

type securityIdpPolicyBlockExemptRuleConfig struct {
	Name        types.String                                      `tfsdk:"name"`
	Description types.String                                      `tfsdk:"description"`
	Match       *securityIdpPolicyBlockExemptRuleBlockMatchConfig `tfsdk:"match"`
}

type securityIdpPolicyBlockExemptRuleBlockMatch struct {
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
}

func (block *securityIdpPolicyBlockExemptRuleBlockMatch) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (securityIdpPolicyBlockExemptRuleBlockMatch) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"custom_attack": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Match custom attacks.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.NoNullValues(),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"custom_attack_group": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Match custom attack groups.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.NoNullValues(),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"destination_address": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Match destination address.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.NoNullValues(),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"destination_address_except": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Don't match destination address.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.NoNullValues(),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"dynamic_attack_group": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Match dynamic attack groups.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.NoNullValues(),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"from_zone": schema.StringAttribute{
			Optional:    true,
			Description: "Match from zone.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"predefined_attack": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Match predefined attacks.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.NoNullValues(),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"predefined_attack_group": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Match predefined attack groups.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.NoNullValues(),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"source_address": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Match source address.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.NoNullValues(),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"source_address_except": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Don't match source address.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.NoNullValues(),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"to_zone": schema.StringAttribute{
			Optional:    true,
			Description: "Match to zone.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
	}
}

type securityIdpPolicyBlockExemptRuleBlockMatchConfig struct {
	CustomAttack             types.Set    `tfsdk:"custom_attack"`
	CustomAttackGroup        types.Set    `tfsdk:"custom_attack_group"`
	DestinationAddress       types.Set    `tfsdk:"destination_address"`
	DestinationAddressExcept types.Set    `tfsdk:"destination_address_except"`
	DynamicAttackGroup       types.Set    `tfsdk:"dynamic_attack_group"`
	FromZone                 types.String `tfsdk:"from_zone"`
	PredefinedAttack         types.Set    `tfsdk:"predefined_attack"`
	PredefinedAttackGroup    types.Set    `tfsdk:"predefined_attack_group"`
	SourceAddress            types.Set    `tfsdk:"source_address"`
	SourceAddressExcept      types.Set    `tfsdk:"source_address_except"`
	ToZone                   types.String `tfsdk:"to_zone"`
}

func (block *securityIdpPolicyBlockExemptRuleBlockMatchConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityIdpPolicyBlockIpsRule struct {
	Name        types.String                             `tfsdk:"name"        tfdata:"identifier"`
	Description types.String                             `tfsdk:"description"`
	Terminal    types.Bool                               `tfsdk:"terminal"`
	Match       *securityIdpPolicyBlockIpsRuleBlockMatch `tfsdk:"match"`
	Then        *securityIdpPolicyBlockIpsRuleBlockThen  `tfsdk:"then"`
}

type securityIdpPolicyBlockIpsRuleConfig struct {
	Name        types.String                                   `tfsdk:"name"`
	Description types.String                                   `tfsdk:"description"`
	Terminal    types.Bool                                     `tfsdk:"terminal"`
	Match       *securityIdpPolicyBlockIpsRuleBlockMatchConfig `tfsdk:"match"`
	Then        *securityIdpPolicyBlockIpsRuleBlockThen        `tfsdk:"then"`
}

type securityIdpPolicyBlockIpsRuleBlockMatch struct {
	securityIdpPolicyBlockExemptRuleBlockMatch
	Application types.String `tfsdk:"application"`
}

func (securityIdpPolicyBlockIpsRuleBlockMatch) attributesSchema() map[string]schema.Attribute {
	attributes := securityIdpPolicyBlockExemptRuleBlockMatch{}.attributesSchema()
	attributes["application"] = schema.StringAttribute{
		Optional:    true,
		Description: "Specify application or application-set name to match.",
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			tfvalidator.StringDoubleQuoteExclusion(),
		},
	}

	return attributes
}

func (block *securityIdpPolicyBlockIpsRuleBlockMatch) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityIdpPolicyBlockIpsRuleBlockMatchConfig struct {
	securityIdpPolicyBlockExemptRuleBlockMatchConfig
	Application types.String `tfsdk:"application"`
}

func (block *securityIdpPolicyBlockIpsRuleBlockMatchConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityIdpPolicyBlockIpsRuleBlockThen struct {
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
}

func (rsc *securityIdpPolicy) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityIdpPolicyConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.ExemptRule.IsNull() &&
		!config.ExemptRule.IsUnknown() {
		var configExemptRule []securityIdpPolicyBlockExemptRuleConfig
		asDiags := config.ExemptRule.ElementsAs(ctx, &configExemptRule, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		exemptRuleName := make(map[string]struct{})
		for i, block := range configExemptRule {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := exemptRuleName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("exempt_rule").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple exempt_rule blocks with the same name %q", name),
					)
				}
				exemptRuleName[name] = struct{}{}
			}

			if block.Match == nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("exempt_rule").AtListIndex(i).AtName("match"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("match block must be specified"+
						" in exempt_rule block %q", block.Name.ValueString()),
				)
			} else {
				if block.Match.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("exempt_rule").AtListIndex(i).AtName("match").AtName("*"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("match block is empty"+
							" in exempt_rule block %q", block.Name.ValueString()),
					)
				}
				if !block.Match.DestinationAddress.IsNull() &&
					!block.Match.DestinationAddress.IsUnknown() &&
					!block.Match.DestinationAddressExcept.IsNull() &&
					!block.Match.DestinationAddressExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("exempt_rule").AtListIndex(i).AtName("match").AtName("destination_address"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("destination_address and destination_address_except cannot be configured together"+
							" in match block in exempt_rule block %q", block.Name.ValueString()),
					)
				}
				if !block.Match.SourceAddress.IsNull() &&
					!block.Match.SourceAddress.IsUnknown() &&
					!block.Match.SourceAddressExcept.IsNull() &&
					!block.Match.SourceAddressExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("exempt_rule").AtListIndex(i).AtName("match").AtName("source_address"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("source_address and source_address_except cannot be configured together"+
							" in match block in exempt_rule block %q", block.Name.ValueString()),
					)
				}
			}
		}
	}
	if !config.IpsRule.IsNull() &&
		!config.IpsRule.IsUnknown() {
		var configIpsRule []securityIdpPolicyBlockIpsRuleConfig
		asDiags := config.IpsRule.ElementsAs(ctx, &configIpsRule, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		ipsRuleName := make(map[string]struct{})
		for i, block := range configIpsRule {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := ipsRuleName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("ips_rule").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple ips_rule blocks with the same name %q", name),
					)
				}
				ipsRuleName[name] = struct{}{}
			}

			if block.Match == nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("ips_rule").AtListIndex(i).AtName("match"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("match block must be specified"+
						" in ips_rule block %q", block.Name.ValueString()),
				)
			} else {
				if block.Match.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("ips_rule").AtListIndex(i).AtName("match").AtName("*"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("match block is empty"+
							" in ips_rule block %q", block.Name.ValueString()),
					)
				}
				if !block.Match.DestinationAddress.IsNull() &&
					!block.Match.DestinationAddress.IsUnknown() &&
					!block.Match.DestinationAddressExcept.IsNull() &&
					!block.Match.DestinationAddressExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("ips_rule").AtListIndex(i).AtName("match").AtName("destination_address"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("destination_address and destination_address_except cannot be configured together"+
							" in match block in ips_rule block %q", block.Name.ValueString()),
					)
				}
				if !block.Match.SourceAddress.IsNull() &&
					!block.Match.SourceAddress.IsUnknown() &&
					!block.Match.SourceAddressExcept.IsNull() &&
					!block.Match.SourceAddressExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("ips_rule").AtListIndex(i).AtName("match").AtName("source_address"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("source_address and source_address_except cannot be configured together"+
							" in match block in ips_rule block %q", block.Name.ValueString()),
					)
				}
			}
			if block.Then == nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("ips_rule").AtListIndex(i).AtName("then"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("then block must be specified"+
						" in ips_rule block %q", block.Name.ValueString()),
				)
			} else {
				if !block.Then.Action.IsNull() &&
					!block.Then.Action.IsUnknown() {
					action := block.Then.Action.ValueString()

					if !block.Then.ActionCosForwardingClass.IsNull() &&
						!block.Then.ActionCosForwardingClass.IsUnknown() &&
						action != "class-of-service" {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("action_cos_forwarding_class"),
							tfdiag.ConflictConfigErrSummary,
							fmt.Sprintf("action_cos_forwarding_class cannot be configured when action != class-of-service"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
					if !block.Then.ActionDscpCodePoint.IsNull() &&
						!block.Then.ActionDscpCodePoint.IsUnknown() &&
						action != "class-of-service" && action != "mark-diffserv" {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("action_dscp_code_point"),
							tfdiag.ConflictConfigErrSummary,
							fmt.Sprintf("action_dscp_code_point cannot be configured when action != class-of-service and mark-diffserv"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
					if action == "class-of-service" &&
						block.Then.ActionCosForwardingClass.IsNull() &&
						block.Then.ActionDscpCodePoint.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("action"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("at least action_cos_forwarding_class or action_dscp_code_point"+
								" must be specified when action = class-of-service"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
					if action == "mark-diffserv" &&
						block.Then.ActionDscpCodePoint.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("action"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("action_dscp_code_point must be specified when action = mark-diffserv"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
				}
				if block.Then.IPAction.IsNull() {
					if !block.Then.IPActionLog.IsNull() &&
						!block.Then.IPActionLog.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("ip_action_log"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("ip_action must be specified with ip_action_log"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
					if !block.Then.IPActionLogCreate.IsNull() &&
						!block.Then.IPActionLogCreate.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("ip_action_log_create"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("ip_action must be specified with ip_action_log_create"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
					if !block.Then.IPActionRefreshTimeout.IsNull() &&
						!block.Then.IPActionRefreshTimeout.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("ip_action_refresh_timeout"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("ip_action must be specified with ip_action_refresh_timeout"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
					if !block.Then.IPActionTarget.IsNull() &&
						!block.Then.IPActionTarget.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("ip_action_target"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("ip_action must be specified with ip_action_target"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
					if !block.Then.IPActionTimeout.IsNull() &&
						!block.Then.IPActionTimeout.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("ip_action_timeout"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("ip_action must be specified with ip_action_timeout"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
				}
				if block.Then.Notification.IsNull() {
					if !block.Then.NotificationLogAttacks.IsNull() &&
						!block.Then.NotificationLogAttacks.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("notification_log_attacks"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("notification must be specified with notification_log_attacks"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
					if !block.Then.NotificationPacketLog.IsNull() &&
						!block.Then.NotificationPacketLog.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("notification_packet_log"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("notification must be specified with notification_packet_log"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
				}
				if block.Then.NotificationLogAttacks.IsNull() &&
					!block.Then.NotificationLogAttacksAlert.IsNull() &&
					!block.Then.NotificationLogAttacksAlert.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("notification_log_attacks_alert"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("notification_log_attacks must be specified with notification_log_attacks_alert"+
							" in then block in ips_rule block %q", block.Name.ValueString()),
					)
				}
				if block.Then.NotificationPacketLog.IsNull() {
					if !block.Then.NotificationPacketLogPostAttack.IsNull() &&
						!block.Then.NotificationPacketLogPostAttack.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("notification_packet_log_post_attack"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("notification_packet_log must be specified with notification_packet_log_post_attack"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
					if !block.Then.NotificationPacketLogPostAttackTimeout.IsNull() &&
						!block.Then.NotificationPacketLogPostAttackTimeout.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("notification_packet_log_post_attack_timeout"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("notification_packet_log must be specified with notification_packet_log_post_attack_timeout"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
					if !block.Then.NotificationPacketLogPreAttack.IsNull() &&
						!block.Then.NotificationPacketLogPreAttack.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("ips_rule").AtListIndex(i).AtName("then").AtName("notification_packet_log_pre_attack"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("notification_packet_log must be specified with notification_packet_log_pre_attack"+
								" in then block in ips_rule block %q", block.Name.ValueString()),
						)
					}
				}
			}
		}
	}
}

func (rsc *securityIdpPolicy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityIdpPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Empty Name",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "name"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			if !junSess.CheckCompatibilitySecurity() {
				resp.Diagnostics.AddError(
					tfdiag.CompatibilityErrSummary,
					rsc.junosName()+junSess.SystemInformation.NotCompatibleMsg(),
				)

				return false
			}
			policyExists, err := checkSecurityIdpPolicyExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if policyExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			policyExists, err := checkSecurityIdpPolicyExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !policyExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *securityIdpPolicy) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityIdpPolicyData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *securityIdpPolicy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityIdpPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *securityIdpPolicy) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityIdpPolicyData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceDelete(
		ctx,
		rsc,
		&state,
		resp,
	)
}

func (rsc *securityIdpPolicy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityIdpPolicyData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "name"),
	)
}

func checkSecurityIdpPolicyExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security idp idp-policy \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityIdpPolicyData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityIdpPolicyData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityIdpPolicyData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set security idp idp-policy \"" + rscData.Name.ValueString() + "\" "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	exemptRuleName := make(map[string]struct{})
	for i, block := range rscData.ExemptRule {
		name := block.Name.ValueString()
		if _, ok := exemptRuleName[name]; ok {
			return path.Root("exempt_rule").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple exempt_rule blocks with the same name %q", name)
		}
		exemptRuleName[name] = struct{}{}

		blockSet, pathErr, err := block.configSet(setPrefix, path.Root("exempt_rule").AtListIndex(i))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	ipsRuleName := make(map[string]struct{})
	for i, block := range rscData.IpsRule {
		name := block.Name.ValueString()
		if _, ok := ipsRuleName[name]; ok {
			return path.Root("ips_rule").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple ips_rule blocks with the same name %q", name)
		}
		ipsRuleName[name] = struct{}{}

		blockSet, pathErr, err := block.configSet(setPrefix, path.Root("ips_rule").AtListIndex(i))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityIdpPolicyBlockExemptRule) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, 100)
	setPrefix += "rulebase-exempt rule \"" + block.Name.ValueString() + "\" "

	if v := block.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if block.Match != nil {
		if block.Match.isEmpty() {
			return configSet,
				pathRoot.AtName("match").AtName("*"),
				errors.New("match block is empty")
		}

		configSet = append(configSet, block.Match.configSet(setPrefix)...)
	} else {
		return configSet,
			pathRoot.AtName("match"),
			errors.New("match block must be specified")
	}

	return configSet, path.Empty(), nil
}

func (block *securityIdpPolicyBlockExemptRuleBlockMatch) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)
	setPrefix += "match "

	for _, v := range block.CustomAttack {
		configSet = append(configSet, setPrefix+"attacks custom-attacks \""+v.ValueString()+"\"")
	}
	for _, v := range block.CustomAttackGroup {
		configSet = append(configSet, setPrefix+"attacks custom-attack-groups \""+v.ValueString()+"\"")
	}
	for _, v := range block.DestinationAddress {
		configSet = append(configSet, setPrefix+"destination-address \""+v.ValueString()+"\"")
	}
	for _, v := range block.DestinationAddressExcept {
		configSet = append(configSet, setPrefix+"destination-except \""+v.ValueString()+"\"")
	}
	for _, v := range block.DynamicAttackGroup {
		configSet = append(configSet, setPrefix+"attacks dynamic-attack-groups \""+v.ValueString()+"\"")
	}
	if v := block.FromZone.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"from-zone \""+v+"\"")
	}
	for _, v := range block.PredefinedAttack {
		configSet = append(configSet, setPrefix+"attacks predefined-attacks \""+v.ValueString()+"\"")
	}
	for _, v := range block.PredefinedAttackGroup {
		configSet = append(configSet, setPrefix+"attacks predefined-attack-groups \""+v.ValueString()+"\"")
	}
	for _, v := range block.SourceAddress {
		configSet = append(configSet, setPrefix+"source-address \""+v.ValueString()+"\"")
	}
	for _, v := range block.SourceAddressExcept {
		configSet = append(configSet, setPrefix+"source-except \""+v.ValueString()+"\"")
	}
	if v := block.ToZone.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"to-zone \""+v+"\"")
	}

	return configSet
}

func (block *securityIdpPolicyBlockIpsRule) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, 100)
	setPrefix += "rulebase-ips rule \"" + block.Name.ValueString() + "\" "

	if v := block.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if block.Terminal.ValueBool() {
		configSet = append(configSet, setPrefix+"terminal")
	}

	if block.Match != nil {
		if block.Match.isEmpty() {
			return configSet,
				pathRoot.AtName("match").AtName("*"),
				errors.New("match block is empty")
		}

		if v := block.Match.Application.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"match application \""+v+"\"")
		}
		configSet = append(
			configSet,
			block.Match.securityIdpPolicyBlockExemptRuleBlockMatch.configSet(setPrefix)...,
		)
	} else {
		return configSet,
			pathRoot.AtName("match"),
			errors.New("match block must be specified")
	}
	if block.Then != nil {
		action := block.Then.Action.ValueString()
		configSet = append(configSet, setPrefix+"then action "+action)

		if v := block.Then.ActionCosForwardingClass.ValueString(); v != "" {
			if action != "class-of-service" {
				return configSet,
					pathRoot.AtName("then").AtName("action_cos_forwarding_class"),
					errors.New("action_cos_forwarding_class cannot be configured when action != class-of-service")
			}

			configSet = append(configSet, setPrefix+"then action "+action+" forwarding-class \""+v+"\"")
		} else if action == "class-of-service" && block.Then.ActionDscpCodePoint.IsNull() {
			return configSet,
				pathRoot.AtName("then").AtName("action"),
				errors.New("at least action_cos_forwarding_class or action_dscp_code_point" +
					" must be specified when action = class-of-service")
		}
		if !block.Then.ActionDscpCodePoint.IsNull() {
			switch {
			case action == "class-of-service":
				configSet = append(configSet, setPrefix+"then action "+action+" dscp-code-point "+
					utils.ConvI64toa(block.Then.ActionDscpCodePoint.ValueInt64()))
			case action == "mark-diffserv":
				configSet = append(configSet, setPrefix+"then action "+action+" "+
					utils.ConvI64toa(block.Then.ActionDscpCodePoint.ValueInt64()))
			default:
				return configSet,
					pathRoot.AtName("then").AtName("action_dscp_code_point"),
					errors.New("action_dscp_code_point cannot be configured when action != class-of-service and mark-diffserv")
			}
		} else if action == "mark-diffserv" {
			return configSet,
				pathRoot.AtName("then").AtName("action"),
				errors.New("action_dscp_code_point must be specified when action = mark-diffserv")
		}
		if v := block.Then.IPAction.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"then ip-action "+v)

			if block.Then.IPActionLog.ValueBool() {
				configSet = append(configSet, setPrefix+"then ip-action log")
			}
			if block.Then.IPActionLogCreate.ValueBool() {
				configSet = append(configSet, setPrefix+"then ip-action log-create")
			}
			if block.Then.IPActionRefreshTimeout.ValueBool() {
				configSet = append(configSet, setPrefix+"then ip-action refresh-timeout")
			}
			if vv := block.Then.IPActionTarget.ValueString(); vv != "" {
				configSet = append(configSet, setPrefix+"then ip-action target "+vv)
			}
			if !block.Then.IPActionTimeout.IsNull() {
				configSet = append(configSet, setPrefix+"then ip-action timeout "+
					utils.ConvI64toa(block.Then.IPActionTimeout.ValueInt64()))
			}
		} else {
			if !block.Then.IPActionLog.IsNull() {
				return configSet,
					pathRoot.AtName("then").AtName("ip_action_log"),
					errors.New("ip_action must be specified with ip_action_log")
			}
			if !block.Then.IPActionLogCreate.IsNull() {
				return configSet,
					pathRoot.AtName("then").AtName("ip_action_log_create"),
					errors.New("ip_action must be specified with ip_action_log_create")
			}
			if !block.Then.IPActionRefreshTimeout.IsNull() {
				return configSet,
					pathRoot.AtName("then").AtName("ip_action_refresh_timeout"),
					errors.New("ip_action must be specified with ip_action_refresh_timeout")
			}
			if !block.Then.IPActionTarget.IsNull() {
				return configSet,
					pathRoot.AtName("then").AtName("ip_action_target"),
					errors.New("ip_action must be specified with ip_action_target")
			}
			if !block.Then.IPActionTimeout.IsNull() {
				return configSet,
					pathRoot.AtName("then").AtName("ip_action_timeout"),
					errors.New("ip_action must be specified with ip_action_timeout")
			}
		}
		if block.Then.Notification.ValueBool() {
			configSet = append(configSet, setPrefix+"then notification")

			if block.Then.NotificationLogAttacks.ValueBool() {
				configSet = append(configSet, setPrefix+"then notification log-attacks")

				if block.Then.NotificationLogAttacksAlert.ValueBool() {
					configSet = append(configSet, setPrefix+"then notification log-attacks alert")
				}
			} else if block.Then.NotificationLogAttacksAlert.ValueBool() {
				return configSet,
					pathRoot.AtName("then").AtName("notification_log_attacks_alert"),
					errors.New("notification_log_attacks must be specified with notification_log_attacks_alert")
			}
			if block.Then.NotificationPacketLog.ValueBool() {
				configSet = append(configSet, setPrefix+"then notification packet-log")

				if !block.Then.NotificationPacketLogPostAttack.IsNull() {
					configSet = append(configSet, setPrefix+"then notification packet-log post-attack "+
						utils.ConvI64toa(block.Then.NotificationPacketLogPostAttack.ValueInt64()))
				}
				if !block.Then.NotificationPacketLogPostAttackTimeout.IsNull() {
					configSet = append(configSet, setPrefix+"then notification packet-log post-attack-timeout "+
						utils.ConvI64toa(block.Then.NotificationPacketLogPostAttackTimeout.ValueInt64()))
				}
				if !block.Then.NotificationPacketLogPreAttack.IsNull() {
					configSet = append(configSet, setPrefix+"then notification packet-log pre-attack "+
						utils.ConvI64toa(block.Then.NotificationPacketLogPreAttack.ValueInt64()))
				}
			} else {
				if !block.Then.NotificationPacketLogPostAttack.IsNull() {
					return configSet,
						pathRoot.AtName("then").AtName("notification_packet_log_post_attack"),
						errors.New("notification_packet_log must be specified with notification_packet_log_post_attack")
				}
				if !block.Then.NotificationPacketLogPostAttackTimeout.IsNull() {
					return configSet,
						pathRoot.AtName("then").AtName("notification_packet_log_post_attack_timeout"),
						errors.New("notification_packet_log must be specified with notification_packet_log_post_attack_timeout")
				}
				if !block.Then.NotificationPacketLogPreAttack.IsNull() {
					return configSet,
						pathRoot.AtName("then").AtName("notification_packet_log_pre_attack"),
						errors.New("notification_packet_log must be specified with notification_packet_log_pre_attack")
				}
			}
		} else {
			if !block.Then.NotificationLogAttacks.IsNull() {
				return configSet,
					pathRoot.AtName("then").AtName("notification_log_attacks"),
					errors.New("notification must be specified with notification_log_attacks")
			}
			if !block.Then.NotificationPacketLog.IsNull() {
				return configSet,
					pathRoot.AtName("then").AtName("notification_packet_log"),
					errors.New("notification must be specified with notification_packet_log")
			}
		}
		if v := block.Then.Severity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"then severity "+v)
		}
	} else {
		return configSet,
			pathRoot.AtName("then"),
			errors.New("then block must be specified")
	}

	return configSet, path.Empty(), nil
}

func (rscData *securityIdpPolicyData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security idp idp-policy \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.fillID()
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "rulebase-exempt rule "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.ExemptRule = tfdata.AppendPotentialNewBlock(rscData.ExemptRule, types.StringValue(strings.Trim(name, "\"")))
				exemptRule := &rscData.ExemptRule[len(rscData.ExemptRule)-1]
				balt.CutPrefixInString(&itemTrim, name+" ")

				exemptRule.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "rulebase-ips rule "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.IpsRule = tfdata.AppendPotentialNewBlock(rscData.IpsRule, types.StringValue(strings.Trim(name, "\"")))
				ipsRule := &rscData.IpsRule[len(rscData.IpsRule)-1]
				balt.CutPrefixInString(&itemTrim, name+" ")

				if err := ipsRule.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *securityIdpPolicyBlockExemptRule) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "description "):
		block.Description = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "match "):
		if block.Match == nil {
			block.Match = &securityIdpPolicyBlockExemptRuleBlockMatch{}
		}

		block.Match.read(itemTrim)
	}
}

func (block *securityIdpPolicyBlockExemptRuleBlockMatch) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "attacks custom-attacks "):
		block.CustomAttack = append(block.CustomAttack,
			types.StringValue(strings.Trim(itemTrim, "\"")),
		)
	case balt.CutPrefixInString(&itemTrim, "attacks custom-attack-groups "):
		block.CustomAttackGroup = append(block.CustomAttackGroup,
			types.StringValue(strings.Trim(itemTrim, "\"")),
		)
	case balt.CutPrefixInString(&itemTrim, "destination-address "):
		block.DestinationAddress = append(block.DestinationAddress,
			types.StringValue(strings.Trim(itemTrim, "\"")),
		)
	case balt.CutPrefixInString(&itemTrim, "destination-except "):
		block.DestinationAddressExcept = append(block.DestinationAddressExcept,
			types.StringValue(strings.Trim(itemTrim, "\"")),
		)
	case balt.CutPrefixInString(&itemTrim, "attacks dynamic-attack-groups "):
		block.DynamicAttackGroup = append(block.DynamicAttackGroup,
			types.StringValue(strings.Trim(itemTrim, "\"")),
		)
	case balt.CutPrefixInString(&itemTrim, "from-zone "):
		block.FromZone = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "attacks predefined-attacks "):
		block.PredefinedAttack = append(block.PredefinedAttack,
			types.StringValue(strings.Trim(itemTrim, "\"")),
		)
	case balt.CutPrefixInString(&itemTrim, "attacks predefined-attack-groups "):
		block.PredefinedAttackGroup = append(block.PredefinedAttackGroup,
			types.StringValue(strings.Trim(itemTrim, "\"")),
		)
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		block.SourceAddress = append(block.SourceAddress,
			types.StringValue(strings.Trim(itemTrim, "\"")),
		)
	case balt.CutPrefixInString(&itemTrim, "source-except "):
		block.SourceAddressExcept = append(block.SourceAddressExcept,
			types.StringValue(strings.Trim(itemTrim, "\"")),
		)
	case balt.CutPrefixInString(&itemTrim, "to-zone "):
		block.ToZone = types.StringValue(strings.Trim(itemTrim, "\""))
	}
}

func (block *securityIdpPolicyBlockIpsRule) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "description "):
		block.Description = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "terminal":
		block.Terminal = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "match "):
		if block.Match == nil {
			block.Match = &securityIdpPolicyBlockIpsRuleBlockMatch{}
		}

		switch {
		case balt.CutPrefixInString(&itemTrim, "application "):
			block.Match.Application = types.StringValue(strings.Trim(itemTrim, "\""))
		default:
			block.Match.securityIdpPolicyBlockExemptRuleBlockMatch.read(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "then "):
		if block.Then == nil {
			block.Then = &securityIdpPolicyBlockIpsRuleBlockThen{}
		}

		switch {
		case balt.CutPrefixInString(&itemTrim, "action "):
			switch {
			case balt.CutPrefixInString(&itemTrim, "class-of-service "):
				block.Then.Action = types.StringValue("class-of-service")
				switch {
				case balt.CutPrefixInString(&itemTrim, "forwarding-class "):
					block.Then.ActionCosForwardingClass = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "dscp-code-point "):
					block.Then.ActionDscpCodePoint, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "mark-diffserv "):
				block.Then.Action = types.StringValue("mark-diffserv")
				block.Then.ActionDscpCodePoint, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			default:
				block.Then.Action = types.StringValue(itemTrim)
			}
		case balt.CutPrefixInString(&itemTrim, "ip-action "):
			switch {
			case itemTrim == "log":
				block.Then.IPActionLog = types.BoolValue(true)
			case itemTrim == "log-create":
				block.Then.IPActionLogCreate = types.BoolValue(true)
			case itemTrim == "refresh-timeout":
				block.Then.IPActionRefreshTimeout = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "target "):
				block.Then.IPActionTarget = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "timeout "):
				block.Then.IPActionTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			default:
				block.Then.IPAction = types.StringValue(itemTrim)
			}
		case balt.CutPrefixInString(&itemTrim, "notification"):
			block.Then.Notification = types.BoolValue(true)
			switch {
			case balt.CutPrefixInString(&itemTrim, " log-attacks"):
				block.Then.NotificationLogAttacks = types.BoolValue(true)
				if itemTrim == " alert" {
					block.Then.NotificationLogAttacksAlert = types.BoolValue(true)
				}
			case balt.CutPrefixInString(&itemTrim, " packet-log"):
				block.Then.NotificationPacketLog = types.BoolValue(true)
				switch {
				case balt.CutPrefixInString(&itemTrim, " post-attack "):
					block.Then.NotificationPacketLogPostAttack, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " post-attack-timeout "):
					block.Then.NotificationPacketLogPostAttackTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " pre-attack "):
					block.Then.NotificationPacketLogPreAttack, err = tfdata.ConvAtoi64Value(itemTrim)
				}
				if err != nil {
					return err
				}
			}
		case balt.CutPrefixInString(&itemTrim, "severity "):
			block.Then.Severity = types.StringValue(itemTrim)
		}
	}

	return nil
}

func (rscData *securityIdpPolicyData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security idp idp-policy \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
