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
							Description: "Action of policy.",
							PlanModifiers: []planmodifier.String{
								tfplanmodifier.StringDefault("permit"),
							},
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
	ID     types.String                 `tfsdk:"id"`
	Policy []securityGlobalPolicyPolicy `tfsdk:"policy"`
}

type securityGlobalPolicyConfig struct {
	ID     types.String `tfsdk:"id"`
	Policy types.List   `tfsdk:"policy"`
}

//nolint:lll
type securityGlobalPolicyPolicy struct {
	Count                           types.Bool                                     `tfsdk:"count"`
	LogInit                         types.Bool                                     `tfsdk:"log_init"`
	LogClose                        types.Bool                                     `tfsdk:"log_close"`
	MatchDestinationAddressExcluded types.Bool                                     `tfsdk:"match_destination_address_excluded"`
	MatchSourceAddressExcluded      types.Bool                                     `tfsdk:"match_source_address_excluded"`
	Name                            types.String                                   `tfsdk:"name"`
	Then                            types.String                                   `tfsdk:"then"`
	MatchSourceEndUserProfile       types.String                                   `tfsdk:"match_source_end_user_profile"`
	MatchSourceAddress              []types.String                                 `tfsdk:"match_source_address"`
	MatchDestinationAddress         []types.String                                 `tfsdk:"match_destination_address"`
	MatchFromZone                   []types.String                                 `tfsdk:"match_from_zone"`
	MatchToZone                     []types.String                                 `tfsdk:"match_to_zone"`
	MatchApplication                []types.String                                 `tfsdk:"match_application"`
	MatchDynamicApplication         []types.String                                 `tfsdk:"match_dynamic_application"`
	PermitApplicationServices       *securityPolicyPolicyPermitApplicationServices `tfsdk:"permit_application_services"`
}

//nolint:lll
type securityGlobalPolicyPolicyConfig struct {
	Count                           types.Bool                                     `tfsdk:"count"`
	LogInit                         types.Bool                                     `tfsdk:"log_init"`
	LogClose                        types.Bool                                     `tfsdk:"log_close"`
	MatchDestinationAddressExcluded types.Bool                                     `tfsdk:"match_destination_address_excluded"`
	MatchSourceAddressExcluded      types.Bool                                     `tfsdk:"match_source_address_excluded"`
	Name                            types.String                                   `tfsdk:"name"`
	Then                            types.String                                   `tfsdk:"then"`
	MatchSourceEndUserProfile       types.String                                   `tfsdk:"match_source_end_user_profile"`
	MatchSourceAddress              types.Set                                      `tfsdk:"match_source_address"`
	MatchDestinationAddress         types.Set                                      `tfsdk:"match_destination_address"`
	MatchFromZone                   types.Set                                      `tfsdk:"match_from_zone"`
	MatchToZone                     types.Set                                      `tfsdk:"match_to_zone"`
	MatchApplication                types.Set                                      `tfsdk:"match_application"`
	MatchDynamicApplication         types.Set                                      `tfsdk:"match_dynamic_application"`
	PermitApplicationServices       *securityPolicyPolicyPermitApplicationServices `tfsdk:"permit_application_services"`
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
			"Missing Configuration Error",
			"at least one policy block must be specified",
		)
	} else if !config.Policy.IsUnknown() {
		var configPolicy []securityGlobalPolicyPolicyConfig
		asDiags := config.Policy.ElementsAs(ctx, &configPolicy, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		policyName := make(map[string]struct{})
		for i, block := range configPolicy {
			if block.MatchApplication.IsNull() && block.MatchDynamicApplication.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("policy").AtListIndex(i).AtName("name"),
					"Missing Configuration Error",
					fmt.Sprintf("at least one of match_application or match_dynamic_application "+
						"must be specified in policy %q", block.Name.ValueString()),
				)
			}
			if block.Name.ValueString() != "" {
				if _, ok := policyName[block.Name.ValueString()]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("name"),
						"Duplicate Configuration Error",
						fmt.Sprintf("multiple policy blocks with the same name %q", block.Name.ValueString()),
					)
				}
				policyName[block.Name.ValueString()] = struct{}{}
			}
			if block.PermitApplicationServices != nil {
				if block.PermitApplicationServices.IsEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("permit_application_services"),
						"Missing Configuration Error",
						fmt.Sprintf("permit_application_services block is empty in policy %q", block.Name.ValueString()),
					)
				}
				if block.Then.ValueString() != "" && block.Then.ValueString() != junos.PermitW {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("then"),
						"Conflict Configuration Error",
						fmt.Sprintf("then is not %q (%q) and permit_application_services is set in policy %q",
							junos.PermitW, block.Then.ValueString(), block.Name.ValueString()),
					)
				}
				if block.PermitApplicationServices.RedirectWx.ValueBool() &&
					block.PermitApplicationServices.ReverseRedirectWx.ValueBool() {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("redirect_wx"),
						"Conflict Configuration Error",
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

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		resp.Diagnostics.AddError(
			"Compatibility Error",
			fmt.Sprintf(rsc.junosName()+" not compatible "+
				"with Junos device %q", junSess.SystemInformation.HardwareModel))

		return
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}
	var check securityGlobalPolicyData
	err = check.read(ctx, junSess)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Pre Check Error", err.Error())

		return
	}
	if len(check.Policy) > 0 {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Pre Check Error", rsc.junosName()+" already exists")

		return
	}

	if errPath, err := plan.set(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("create resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityGlobalPolicy) Read(
	ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	var data securityGlobalPolicyData
	junos.MutexLock()
	err = data.read(ctx, junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *securityGlobalPolicy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan securityGlobalPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := plan.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	if err := plan.del(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("update resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityGlobalPolicy) Delete(
	ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var empty securityGlobalPolicyData

	if rsc.client.FakeDeleteAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := empty.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	if err := empty.del(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	warns, err := junSess.CommitConf("delete resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}
}

func (rsc *securityGlobalPolicy) ImportState(
	ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	var data securityGlobalPolicyData
	if err := data.read(ctx, junSess); err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rscData *securityGlobalPolicyData) fillID() {
	rscData.ID = types.StringValue("security_global_policy")
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
			if block.PermitApplicationServices.IsEmpty() {
				return path.Root("policy").AtListIndex(i).AtName("permit_application_services"), fmt.Errorf(
					"permit_application_services block is empty in policy %q",
					block.Name.ValueString(),
				)
			}
			if block.Then.ValueString() != junos.PermitW {
				return path.Root("policy").AtListIndex(i).AtName("then"), fmt.Errorf(
					"conflict: then is not %q (%q) and permit_application_services is set in policy %q",
					junos.PermitW, block.Then.ValueString(), block.Name.ValueString(),
				)
			}
			configSetAppSvc, err := block.PermitApplicationServices.set(setPrefixPolicy)
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
) (
	err error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "security policies global" + junos.PipeDisplaySetRelative)
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
				var policy securityGlobalPolicyPolicy
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
							policy.PermitApplicationServices = &securityPolicyPolicyPermitApplicationServices{}
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
