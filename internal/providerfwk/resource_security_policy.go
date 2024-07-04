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

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &securityPolicy{}
	_ resource.ResourceWithConfigure      = &securityPolicy{}
	_ resource.ResourceWithValidateConfig = &securityPolicy{}
	_ resource.ResourceWithImportState    = &securityPolicy{}
	_ resource.ResourceWithUpgradeState   = &securityPolicy{}
)

type securityPolicy struct {
	client *junos.Client
}

func newSecurityPolicyResource() resource.Resource {
	return &securityPolicy{}
}

func (rsc *securityPolicy) typeName() string {
	return providerName + "_security_policy"
}

func (rsc *securityPolicy) junosName() string {
	return "security policy"
}

func (rsc *securityPolicy) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityPolicy) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityPolicy) Configure(
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

func (rsc *securityPolicy) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<from_zone>" + junos.IDSeparator + "<to_zone>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"from_zone": schema.StringAttribute{
				Required:    true,
				Description: "The name of source zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"to_zone": schema.StringAttribute{
				Required:    true,
				Description: "The name of destination zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"policy": schema.ListNestedBlock{
				Description: "For each name of policy.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The name of policy.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"match_source_address": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "List of source address match.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 250),
									tfvalidator.StringDoubleQuoteExclusion(),
								),
							},
						},
						"match_destination_address": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "List of destination address match.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 250),
									tfvalidator.StringDoubleQuoteExclusion(),
								),
							},
						},
						"then": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString("permit"),
							Description: "Action of policy.",
							Validators: []validator.String{
								stringvalidator.OneOf("permit", "reject", "deny"),
							},
						},
						"count": schema.BoolAttribute{
							Optional:    true,
							Description: "Enable count.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"log_init": schema.BoolAttribute{
							Optional:    true,
							Description: "Log at session init time.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"log_close": schema.BoolAttribute{
							Optional:    true,
							Description: "Log at session close time.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"match_application": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of applications match.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 250),
									tfvalidator.StringDoubleQuoteExclusion(),
								),
							},
						},
						"match_destination_address_excluded": schema.BoolAttribute{
							Optional:    true,
							Description: "Exclude destination addresses.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"match_dynamic_application": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of dynamic application or group match.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 250),
									tfvalidator.StringDoubleQuoteExclusion(),
								),
							},
						},
						"match_source_address_excluded": schema.BoolAttribute{
							Optional:    true,
							Description: "Exclude source addresses.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"match_source_end_user_profile": schema.StringAttribute{
							Optional:    true,
							Description: "Match source end user profile (device identity profile).",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"permit_tunnel_ipsec_vpn": schema.StringAttribute{
							Optional:    true,
							Description: "Name of vpn to permit with a tunnel ipsec.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"permit_application_services": schema.SingleNestedBlock{
							Description: "Define application services for permit.",
							Attributes: map[string]schema.Attribute{
								"advanced_anti_malware_policy": schema.StringAttribute{
									Optional:    true,
									Description: "Specify advanced-anti-malware policy name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"application_firewall_rule_set": schema.StringAttribute{
									Optional:    true,
									Description: "Service rule-set name for Application firewall.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"application_traffic_control_rule_set": schema.StringAttribute{
									Optional:    true,
									Description: "Service rule-set name Application traffic control.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"gprs_gtp_profile": schema.StringAttribute{
									Optional:    true,
									Description: "Specify GPRS Tunneling Protocol profile name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"gprs_sctp_profile": schema.StringAttribute{
									Optional:    true,
									Description: "Specify GPRS stream control protocol profile name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"idp": schema.BoolAttribute{
									Optional:    true,
									Description: "Enable Intrusion detection and prevention.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"idp_policy": schema.StringAttribute{
									Optional:    true,
									Description: "Specify idp policy name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"redirect_wx": schema.BoolAttribute{
									Optional:    true,
									Description: "Set WX redirection.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"reverse_redirect_wx": schema.BoolAttribute{
									Optional:    true,
									Description: "Set WX reverse redirection.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"security_intelligence_policy": schema.StringAttribute{
									Optional:    true,
									Description: "Specify security-intelligence policy name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"utm_policy": schema.StringAttribute{
									Optional:    true,
									Description: "Specify utm policy name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
							},
							Blocks: map[string]schema.Block{
								"ssl_proxy": schema.SingleNestedBlock{
									Description: "Enable SSL Proxy.",
									Attributes: map[string]schema.Attribute{
										"profile_name": schema.StringAttribute{
											Optional:    true,
											Description: "Specify SSL proxy service profile name.",
											Validators: []validator.String{
												stringvalidator.LengthBetween(1, 250),
												tfvalidator.StringDoubleQuoteExclusion(),
											},
										},
									},
									PlanModifiers: []planmodifier.Object{
										tfplanmodifier.BlockRemoveNull(),
									},
								},
								"uac_policy": schema.SingleNestedBlock{
									Description: "Enable unified access control enforcement.",
									Attributes: map[string]schema.Attribute{
										"captive_portal": schema.StringAttribute{
											Optional:    true,
											Description: "Specify captive portal.",
											Validators: []validator.String{
												stringvalidator.LengthBetween(1, 250),
												tfvalidator.StringDoubleQuoteExclusion(),
											},
										},
									},
									PlanModifiers: []planmodifier.Object{
										tfplanmodifier.BlockRemoveNull(),
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

type securityPolicyData struct {
	ID       types.String                `tfsdk:"id"`
	FromZone types.String                `tfsdk:"from_zone"`
	ToZone   types.String                `tfsdk:"to_zone"`
	Policy   []securityPolicyBlockPolicy `tfsdk:"policy"`
}

type securityPolicyConfig struct {
	ID       types.String `tfsdk:"id"`
	FromZone types.String `tfsdk:"from_zone"`
	ToZone   types.String `tfsdk:"to_zone"`
	Policy   types.List   `tfsdk:"policy"`
}

//nolint:lll
type securityPolicyBlockPolicy struct {
	Name                            types.String                                             `tfsdk:"name"`
	MatchSourceAddress              []types.String                                           `tfsdk:"match_source_address"`
	MatchDestinationAddress         []types.String                                           `tfsdk:"match_destination_address"`
	Then                            types.String                                             `tfsdk:"then"`
	Count                           types.Bool                                               `tfsdk:"count"`
	LogInit                         types.Bool                                               `tfsdk:"log_init"`
	LogClose                        types.Bool                                               `tfsdk:"log_close"`
	MatchApplication                []types.String                                           `tfsdk:"match_application"`
	MatchDestinationAddressExcluded types.Bool                                               `tfsdk:"match_destination_address_excluded"`
	MatchDynamicApplication         []types.String                                           `tfsdk:"match_dynamic_application"`
	MatchSourceAddressExcluded      types.Bool                                               `tfsdk:"match_source_address_excluded"`
	MatchSourceEndUserProfile       types.String                                             `tfsdk:"match_source_end_user_profile"`
	PermitTunnelIpsecVpn            types.String                                             `tfsdk:"permit_tunnel_ipsec_vpn"`
	PermitApplicationServices       *securityPolicyBlockPolicyBlockPermitApplicationServices `tfsdk:"permit_application_services"`
}

//nolint:lll
type securityPolicyBlockPolicyConfig struct {
	Name                            types.String                                             `tfsdk:"name"`
	MatchSourceAddress              types.Set                                                `tfsdk:"match_source_address"`
	MatchDestinationAddress         types.Set                                                `tfsdk:"match_destination_address"`
	Then                            types.String                                             `tfsdk:"then"`
	Count                           types.Bool                                               `tfsdk:"count"`
	LogInit                         types.Bool                                               `tfsdk:"log_init"`
	LogClose                        types.Bool                                               `tfsdk:"log_close"`
	MatchApplication                types.Set                                                `tfsdk:"match_application"`
	MatchDestinationAddressExcluded types.Bool                                               `tfsdk:"match_destination_address_excluded"`
	MatchDynamicApplication         types.Set                                                `tfsdk:"match_dynamic_application"`
	MatchSourceAddressExcluded      types.Bool                                               `tfsdk:"match_source_address_excluded"`
	MatchSourceEndUserProfile       types.String                                             `tfsdk:"match_source_end_user_profile"`
	PermitTunnelIpsecVpn            types.String                                             `tfsdk:"permit_tunnel_ipsec_vpn"`
	PermitApplicationServices       *securityPolicyBlockPolicyBlockPermitApplicationServices `tfsdk:"permit_application_services"`
}

type securityPolicyBlockPolicyBlockPermitApplicationServices struct {
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
	UtmPolicy                        types.String `tfsdk:"utm_policy"`
	SSLProxy                         *struct {
		ProfileName types.String `tfsdk:"profile_name"`
	} `tfsdk:"ssl_proxy"`
	UacPolicy *struct {
		CaptivePortal types.String `tfsdk:"captive_portal"`
	} `tfsdk:"uac_policy"`
}

func (block *securityPolicyBlockPolicyBlockPermitApplicationServices) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *securityPolicyBlockPolicyBlockPermitApplicationServices) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (rsc *securityPolicy) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityPolicyConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Policy.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("policy").AtName("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one policy block must be specified",
		)
	} else if !config.Policy.IsUnknown() {
		var configPolicy []securityPolicyBlockPolicyConfig
		asDiags := config.Policy.ElementsAs(ctx, &configPolicy, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		policyName := make(map[string]struct{})
		for i, block := range configPolicy {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := policyName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple policy blocks with the same name %q", name),
					)
				}
				policyName[name] = struct{}{}
			}
			if block.MatchApplication.IsNull() && block.MatchDynamicApplication.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("policy").AtListIndex(i).AtName("name"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("at least one of match_application or match_dynamic_application "+
						"must be specified in policy %q", block.Name.ValueString()),
				)
			}
			if !block.PermitTunnelIpsecVpn.IsNull() && !block.PermitTunnelIpsecVpn.IsUnknown() &&
				!block.Then.IsNull() && !block.Then.IsUnknown() && block.Then.ValueString() != junos.PermitW {
				resp.Diagnostics.AddAttributeError(
					path.Root("policy").AtListIndex(i).AtName("then"),
					tfdiag.ConflictConfigErrSummary,
					fmt.Sprintf("then is not %q (got %q) and permit_tunnel_ipsec_vpn is set in policy %q",
						junos.PermitW, block.Then.ValueString(), block.Name.ValueString()),
				)
			}
			if block.PermitApplicationServices != nil {
				if block.PermitApplicationServices.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("permit_application_services"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("permit_application_services block is empty in policy %q", block.Name.ValueString()),
					)
				} else if block.PermitApplicationServices.hasKnownValue() &&
					!block.Then.IsNull() && !block.Then.IsUnknown() && block.Then.ValueString() != junos.PermitW {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("then"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("then is not %q (got %q) and permit_application_services is set in policy %q",
							junos.PermitW, block.Then.ValueString(), block.Name.ValueString()),
					)
				}
				if !block.PermitApplicationServices.RedirectWx.IsNull() &&
					!block.PermitApplicationServices.RedirectWx.IsUnknown() &&
					!block.PermitApplicationServices.ReverseRedirectWx.IsNull() &&
					!block.PermitApplicationServices.ReverseRedirectWx.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("redirect_wx"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("redirect_wx and reverse_redirect_wx enabled both in policy %q", block.Name.ValueString()),
					)
				}
			}
		}
	}
}

func (rsc *securityPolicy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.FromZone.ValueString() == "" || plan.ToZone.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Empty Zone",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "from_zone or to_zone"),
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
			policyExists, err := checkSecurityPolicyExists(
				fnCtx,
				plan.FromZone.ValueString(),
				plan.ToZone.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if policyExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf(rsc.junosName()+" from %q to %q already exists",
						plan.FromZone.ValueString(), plan.ToZone.ValueString()),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			policyExists, err := checkSecurityPolicyExists(
				fnCtx,
				plan.FromZone.ValueString(),
				plan.ToZone.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !policyExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					fmt.Sprintf(rsc.junosName()+" from %q to %q not exists after commit "+
						"=> check your config", plan.FromZone.ValueString(), plan.ToZone.ValueString()),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *securityPolicy) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityPolicyData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.FromZone.ValueString(),
			state.ToZone.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *securityPolicy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
			} else {
				resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigClearUnlockWarnSummary, junSess.ConfigClear())...)
	}()

	listLinesToPairPolicy, err := readSecurityPolicyTunnelPairPolicyLines(
		ctx,
		state.FromZone.ValueString(),
		state.ToZone.ValueString(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	if err := junSess.ConfigSet(listLinesToPairPolicy); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())

		return
	}
	warns, err := junSess.CommitConf(ctx, "update resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityPolicy) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityPolicyData
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

func (rsc *securityPolicy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityPolicyData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <zone>"+junos.IDSeparator+"<name>)",
	)
}

func checkSecurityPolicyExists(
	_ context.Context, fromZone, toZone string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security policies from-zone " + fromZone + " to-zone " + toZone + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityPolicyData) fillID() {
	rscData.ID = types.StringValue(rscData.FromZone.ValueString() + junos.IDSeparator + rscData.ToZone.ValueString())
}

func (rscData *securityPolicyData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityPolicyData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security policies" +
		" from-zone " + rscData.FromZone.ValueString() +
		" to-zone " + rscData.ToZone.ValueString() +
		" policy "

	policyName := make(map[string]struct{})
	for i, block := range rscData.Policy {
		name := block.Name.ValueString()
		if _, ok := policyName[name]; ok {
			return path.Root("policy").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple policy blocks with the same name %q", name)
		}
		policyName[name] = struct{}{}

		setPrefixPolicy := setPrefix + name + " "
		for _, v := range block.MatchSourceAddress {
			configSet = append(configSet, setPrefixPolicy+"match source-address \""+v.ValueString()+"\"")
		}
		for _, v := range block.MatchDestinationAddress {
			configSet = append(configSet, setPrefixPolicy+"match destination-address \""+v.ValueString()+"\"")
		}
		configSet = append(configSet, setPrefixPolicy+"then "+block.Then.ValueString())
		if block.Count.ValueBool() {
			configSet = append(configSet, setPrefixPolicy+"then count")
		}
		if block.LogInit.ValueBool() {
			configSet = append(configSet, setPrefixPolicy+"then log session-init")
		}
		if block.LogClose.ValueBool() {
			configSet = append(configSet, setPrefixPolicy+"then log session-close")
		}
		if len(block.MatchApplication) == 0 &&
			len(block.MatchDynamicApplication) == 0 {
			return path.Root("policy").AtListIndex(i).AtName("name"),
				fmt.Errorf("at least one of match_application or match_dynamic_application "+
					"must be specified in policy %q", block.Name.ValueString())
		}
		for _, v := range block.MatchApplication {
			configSet = append(configSet, setPrefixPolicy+"match application \""+v.ValueString()+"\"")
		}
		if block.MatchDestinationAddressExcluded.ValueBool() {
			configSet = append(configSet, setPrefixPolicy+"match destination-address-excluded")
		}
		for _, v := range block.MatchDynamicApplication {
			configSet = append(configSet, setPrefixPolicy+"match dynamic-application \""+v.ValueString()+"\"")
		}
		if block.MatchSourceAddressExcluded.ValueBool() {
			configSet = append(configSet, setPrefixPolicy+"match source-address-excluded")
		}
		if v := block.MatchSourceEndUserProfile.ValueString(); v != "" {
			configSet = append(configSet, setPrefixPolicy+"match source-end-user-profile \""+v+"\"")
		}
		if v := block.PermitTunnelIpsecVpn.ValueString(); v != "" {
			if block.Then.ValueString() != junos.PermitW {
				return path.Root("policy").AtListIndex(i).AtName("then"), fmt.Errorf(
					"conflict: then is not %q (got %q) and permit_tunnel_ipsec_vpn is set in policy %q",
					junos.PermitW, block.Then.ValueString(), block.Name.ValueString(),
				)
			}
			configSet = append(configSet, setPrefixPolicy+"then permit tunnel ipsec-vpn \""+
				block.PermitTunnelIpsecVpn.ValueString()+"\"")
		}
		if block.PermitApplicationServices != nil {
			if block.PermitApplicationServices.isEmpty() {
				return path.Root("policy").AtListIndex(i).AtName("permit_application_services"), fmt.Errorf(
					"permit_application_services block is empty in policy %q",
					block.Name.ValueString(),
				)
			}
			if block.Then.ValueString() != junos.PermitW {
				return path.Root("policy").AtListIndex(i).AtName("then"), fmt.Errorf(
					"conflict: then is not %q (got %q) and permit_application_services is set in policy %q",
					junos.PermitW, block.Then.ValueString(), block.Name.ValueString(),
				)
			}
			configSetAppSvc, err := block.PermitApplicationServices.configSet(setPrefixPolicy)
			if err != nil {
				return path.Root("policy").AtListIndex(i).AtName("permit_application_services"), err
			}
			if len(configSetAppSvc) == 0 {
				return path.Root("policy").AtListIndex(i).AtName("permit_application_services"), fmt.Errorf(
					"permit_application_services block is empty in policy %q",
					block.Name.ValueString(),
				)
			}
			configSet = append(configSet, configSetAppSvc...)
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityPolicyBlockPolicyBlockPermitApplicationServices) configSet(
	setPrefixPolicy string,
) (
	[]string, error,
) {
	configSet := make([]string, 0)
	setPrefixPolicyPermitAppSvc := setPrefixPolicy + "then permit application-services "

	if v := block.AdvancedAntiMalwarePolicy.ValueString(); v != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"advanced-anti-malware-policy \""+v+"\"")
	}
	if v := block.ApplicationFirewallRuleSet.ValueString(); v != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"application-firewall rule-set \""+v+"\"")
	}
	if v := block.ApplicationTrafficControlRuleSet.ValueString(); v != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"application-traffic-control rule-set \""+v+"\"")
	}
	if v := block.GprsGtpProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"gprs-gtp-profile \""+v+"\"")
	}
	if v := block.GprsSctpProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"gprs-sctp-profile \""+v+"\"")
	}
	if block.Idp.ValueBool() {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"idp")
	}
	if v := block.IdpPolicy.ValueString(); v != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"idp-policy \""+v+"\"")
	}
	if block.RedirectWx.ValueBool() && block.ReverseRedirectWx.ValueBool() {
		return configSet, errors.New("conflict: redirect_wx and reverse_redirect_wx enabled both")
	}
	if block.RedirectWx.ValueBool() {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"redirect-wx")
	}
	if block.ReverseRedirectWx.ValueBool() {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"reverse-redirect-wx")
	}
	if v := block.SecurityIntelligencePolicy.ValueString(); v != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"security-intelligence-policy \""+v+"\"")
	}
	if block.SSLProxy != nil {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"ssl-proxy")
		if v := block.SSLProxy.ProfileName.ValueString(); v != "" {
			configSet = append(configSet, setPrefixPolicyPermitAppSvc+"ssl-proxy profile-name \""+v+"\"")
		}
	}
	if block.UacPolicy != nil {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"uac-policy")
		if v := block.UacPolicy.CaptivePortal.ValueString(); v != "" {
			configSet = append(configSet, setPrefixPolicyPermitAppSvc+"uac-policy captive-portal \""+v+"\"")
		}
	}
	if v := block.UtmPolicy.ValueString(); v != "" {
		configSet = append(configSet, setPrefixPolicyPermitAppSvc+"utm-policy \""+v+"\"")
	}

	return configSet, nil
}

func (rscData *securityPolicyData) read(
	_ context.Context, fromZone, toZone string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security policies from-zone " + fromZone + " to-zone " + toZone + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.FromZone = types.StringValue(fromZone)
		rscData.ToZone = types.StringValue(toZone)
		rscData.fillID()
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "policy ") {
				itemTrimFields := strings.Split(itemTrim, " ")
				var policy securityPolicyBlockPolicy
				rscData.Policy, policy = tfdata.ExtractBlockWithTFTypesString(rscData.Policy, "Name", itemTrimFields[0])
				policy.Name = types.StringValue(itemTrimFields[0])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "match source-address "):
					policy.MatchSourceAddress = append(policy.MatchSourceAddress,
						types.StringValue(strings.Trim(itemTrim, "\"")))
				case balt.CutPrefixInString(&itemTrim, "match destination-address "):
					policy.MatchDestinationAddress = append(policy.MatchDestinationAddress,
						types.StringValue(strings.Trim(itemTrim, "\"")))
				case balt.CutPrefixInString(&itemTrim, "match application "):
					policy.MatchApplication = append(policy.MatchApplication,
						types.StringValue(strings.Trim(itemTrim, "\"")))
				case itemTrim == "match destination-address-excluded":
					policy.MatchDestinationAddressExcluded = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "match dynamic-application "):
					policy.MatchDynamicApplication = append(policy.MatchDynamicApplication,
						types.StringValue(strings.Trim(itemTrim, "\"")))
				case itemTrim == "match source-address-excluded":
					policy.MatchSourceAddressExcluded = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "match source-end-user-profile "):
					policy.MatchSourceEndUserProfile = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "then "):
					switch {
					case itemTrim == "permit",
						itemTrim == "deny",
						itemTrim == "reject":
						policy.Then = types.StringValue(itemTrim)
					case itemTrim == "count":
						policy.Count = types.BoolValue(true)
					case itemTrim == "log session-init":
						policy.LogInit = types.BoolValue(true)
					case itemTrim == "log session-close":
						policy.LogClose = types.BoolValue(true)
					case balt.CutPrefixInString(&itemTrim, "permit tunnel ipsec-vpn "):
						policy.Then = types.StringValue(junos.PermitW)
						policy.PermitTunnelIpsecVpn = types.StringValue(strings.Trim(itemTrim, "\""))
					case balt.CutPrefixInString(&itemTrim, "permit application-services "):
						policy.Then = types.StringValue(junos.PermitW)
						if policy.PermitApplicationServices == nil {
							policy.PermitApplicationServices = &securityPolicyBlockPolicyBlockPermitApplicationServices{}
						}
						policy.PermitApplicationServices.read(itemTrim)
					}
				}
				rscData.Policy = append(rscData.Policy, policy)
			}
		}
	}

	return nil
}

func (block *securityPolicyBlockPolicyBlockPermitApplicationServices) read(
	itemTrim string,
) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "advanced-anti-malware-policy "):
		block.AdvancedAntiMalwarePolicy = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "application-firewall rule-set "):
		block.ApplicationFirewallRuleSet = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "application-traffic-control rule-set "):
		block.ApplicationTrafficControlRuleSet = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "gprs-gtp-profile "):
		block.GprsGtpProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "gprs-sctp-profile "):
		block.GprsSctpProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "idp":
		block.Idp = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "idp-policy "):
		block.IdpPolicy = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "redirect-wx":
		block.RedirectWx = types.BoolValue(true)
	case itemTrim == "reverse-redirect-wx":
		block.ReverseRedirectWx = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "security-intelligence-policy "):
		block.SecurityIntelligencePolicy = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "utm-policy "):
		block.UtmPolicy = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "ssl-proxy"):
		if balt.CutPrefixInString(&itemTrim, " profile-name ") {
			block.SSLProxy = &struct {
				ProfileName types.String `tfsdk:"profile_name"`
			}{
				ProfileName: types.StringValue(strings.Trim(itemTrim, "\"")),
			}
		} else {
			block.SSLProxy = &struct {
				ProfileName types.String `tfsdk:"profile_name"`
			}{}
		}
	case balt.CutPrefixInString(&itemTrim, "uac-policy"):
		if balt.CutPrefixInString(&itemTrim, " captive-portal ") {
			block.UacPolicy = &struct {
				CaptivePortal types.String `tfsdk:"captive_portal"`
			}{
				CaptivePortal: types.StringValue(strings.Trim(itemTrim, "\"")),
			}
		} else {
			block.UacPolicy = &struct {
				CaptivePortal types.String `tfsdk:"captive_portal"`
			}{}
		}
	}
}

func readSecurityPolicyTunnelPairPolicyLines(
	_ context.Context, fromZone, toZone string, junSess *junos.Session,
) (
	[]string, error,
) {
	listLines := make([]string, 0)
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security policies from-zone " + fromZone + " to-zone " + toZone + junos.PipeDisplaySet)
	if err != nil {
		return listLines, err
	}
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			if strings.Contains(item, "then permit tunnel pair-policy ") {
				listLines = append(listLines, item)
			}
		}
	}

	return listLines, nil
}

func (rscData *securityPolicyData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security policies from-zone " + rscData.FromZone.ValueString() + " to-zone " + rscData.ToZone.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
