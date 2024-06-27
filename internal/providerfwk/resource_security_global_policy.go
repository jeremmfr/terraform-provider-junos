package providerfwk

import (
	"context"
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
	_ resource.Resource                   = &securityGlobalPolicy{}
	_ resource.ResourceWithConfigure      = &securityGlobalPolicy{}
	_ resource.ResourceWithValidateConfig = &securityGlobalPolicy{}
	_ resource.ResourceWithImportState    = &securityGlobalPolicy{}
	_ resource.ResourceWithUpgradeState   = &securityGlobalPolicy{}
)

type securityGlobalPolicy struct {
	client *junos.Client
}

func newSecurityGlobalPolicyResource() resource.Resource {
	return &securityGlobalPolicy{}
}

func (rsc *securityGlobalPolicy) typeName() string {
	return providerName + "_security_global_policy"
}

func (rsc *securityGlobalPolicy) junosName() string {
	return "security policies global"
}

func (rsc *securityGlobalPolicy) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityGlobalPolicy) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityGlobalPolicy) Configure(
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

func (rsc *securityGlobalPolicy) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `security_global_policy`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"policy": schema.ListNestedBlock{
				Description: "For each policy name.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Security policy name.",
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
						"match_from_zone": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "Match multiple source zone.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringFormat(tfvalidator.DefaultFormat),
								),
							},
						},
						"match_to_zone": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "Match multiple destination zone.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringFormat(tfvalidator.DefaultFormat),
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
					},
					Blocks: map[string]schema.Block{
						"permit_application_services": schema.SingleNestedBlock{
							Description: "Declare `permit application-services` configuration.",
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

type securityGlobalPolicyData struct {
	ID     types.String                      `tfsdk:"id"`
	Policy []securityGlobalPolicyBlockPolicy `tfsdk:"policy"`
}

type securityGlobalPolicyConfig struct {
	ID     types.String `tfsdk:"id"`
	Policy types.List   `tfsdk:"policy"`
}

//nolint:lll
type securityGlobalPolicyBlockPolicy struct {
	Name                            types.String                                             `tfsdk:"name"`
	MatchSourceAddress              []types.String                                           `tfsdk:"match_source_address"`
	MatchDestinationAddress         []types.String                                           `tfsdk:"match_destination_address"`
	MatchFromZone                   []types.String                                           `tfsdk:"match_from_zone"`
	MatchToZone                     []types.String                                           `tfsdk:"match_to_zone"`
	Then                            types.String                                             `tfsdk:"then"`
	Count                           types.Bool                                               `tfsdk:"count"`
	LogInit                         types.Bool                                               `tfsdk:"log_init"`
	LogClose                        types.Bool                                               `tfsdk:"log_close"`
	MatchApplication                []types.String                                           `tfsdk:"match_application"`
	MatchDestinationAddressExcluded types.Bool                                               `tfsdk:"match_destination_address_excluded"`
	MatchDynamicApplication         []types.String                                           `tfsdk:"match_dynamic_application"`
	MatchSourceAddressExcluded      types.Bool                                               `tfsdk:"match_source_address_excluded"`
	MatchSourceEndUserProfile       types.String                                             `tfsdk:"match_source_end_user_profile"`
	PermitApplicationServices       *securityPolicyBlockPolicyBlockPermitApplicationServices `tfsdk:"permit_application_services"`
}

//nolint:lll
type securityGlobalPolicyBlockPolicyConfig struct {
	Name                            types.String                                             `tfsdk:"name"`
	MatchSourceAddress              types.Set                                                `tfsdk:"match_source_address"`
	MatchDestinationAddress         types.Set                                                `tfsdk:"match_destination_address"`
	MatchFromZone                   types.Set                                                `tfsdk:"match_from_zone"`
	MatchToZone                     types.Set                                                `tfsdk:"match_to_zone"`
	Then                            types.String                                             `tfsdk:"then"`
	Count                           types.Bool                                               `tfsdk:"count"`
	LogInit                         types.Bool                                               `tfsdk:"log_init"`
	LogClose                        types.Bool                                               `tfsdk:"log_close"`
	MatchApplication                types.Set                                                `tfsdk:"match_application"`
	MatchDestinationAddressExcluded types.Bool                                               `tfsdk:"match_destination_address_excluded"`
	MatchDynamicApplication         types.Set                                                `tfsdk:"match_dynamic_application"`
	MatchSourceAddressExcluded      types.Bool                                               `tfsdk:"match_source_address_excluded"`
	MatchSourceEndUserProfile       types.String                                             `tfsdk:"match_source_end_user_profile"`
	PermitApplicationServices       *securityPolicyBlockPolicyBlockPermitApplicationServices `tfsdk:"permit_application_services"`
}

func (rsc *securityGlobalPolicy) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityGlobalPolicyConfig
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
		var configPolicy []securityGlobalPolicyBlockPolicyConfig
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

func (rsc *securityGlobalPolicy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityGlobalPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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
			var check securityGlobalPolicyData
			if err := check.read(fnCtx, junSess); err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if len(check.Policy) > 0 {
				resp.Diagnostics.AddError(tfdiag.DuplicateConfigErrSummary, rsc.junosName()+" already exists")

				return false
			}

			return true
		},
		nil,
		&plan,
		resp,
	)
}

func (rsc *securityGlobalPolicy) Read(
	ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse,
) {
	var data securityGlobalPolicyData

	var _ resourceDataReadWithoutArg = &data
	defaultResourceRead(
		ctx,
		rsc,
		nil,
		&data,
		nil,
		resp,
	)
}

func (rsc *securityGlobalPolicy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan securityGlobalPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceUpdate(
		ctx,
		rsc,
		&plan,
		&plan,
		resp,
	)
}

func (rsc *securityGlobalPolicy) Delete(
	ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var empty securityGlobalPolicyData

	defaultResourceDelete(
		ctx,
		rsc,
		&empty,
		resp,
	)
}

func (rsc *securityGlobalPolicy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityGlobalPolicyData

	var _ resourceDataReadWithoutArg = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		"",
	)
}

func (rscData *securityGlobalPolicyData) fillID() {
	rscData.ID = types.StringValue("security_global_policy")
}

func (rscData *securityGlobalPolicyData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityGlobalPolicyData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security policies global policy "

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
		for _, v := range block.MatchFromZone {
			configSet = append(configSet, setPrefixPolicy+"match from-zone "+v.ValueString())
		}
		for _, v := range block.MatchToZone {
			configSet = append(configSet, setPrefixPolicy+"match to-zone "+v.ValueString())
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

func (rscData *securityGlobalPolicyData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security policies global" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	rscData.fillID()
	if showConfig != junos.EmptyW {
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
				var policy securityGlobalPolicyBlockPolicy
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
				case balt.CutPrefixInString(&itemTrim, "match from-zone "):
					policy.MatchFromZone = append(policy.MatchFromZone, types.StringValue(itemTrim))
				case balt.CutPrefixInString(&itemTrim, "match to-zone "):
					policy.MatchToZone = append(policy.MatchToZone, types.StringValue(itemTrim))
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

func (rscData *securityGlobalPolicyData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security policies global",
	}

	return junSess.ConfigSet(configSet)
}
