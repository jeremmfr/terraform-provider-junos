package provider

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
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                   = &securityUtmPolicy{}
	_ resource.ResourceWithConfigure      = &securityUtmPolicy{}
	_ resource.ResourceWithValidateConfig = &securityUtmPolicy{}
	_ resource.ResourceWithImportState    = &securityUtmPolicy{}
	_ resource.ResourceWithUpgradeState   = &securityUtmPolicy{}
)

type securityUtmPolicy struct {
	client *junos.Client
}

func newSecurityUtmPolicyResource() resource.Resource {
	return &securityUtmPolicy{}
}

func (rsc *securityUtmPolicy) typeName() string {
	return providerName + "_security_utm_policy"
}

func (rsc *securityUtmPolicy) junosName() string {
	return "security utm utm-policy"
}

func (rsc *securityUtmPolicy) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityUtmPolicy) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityUtmPolicy) Configure(
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

func (rsc *securityUtmPolicy) Schema(
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
				Description: "The name of security utm utm-policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 29),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"anti_spam_smtp_profile": schema.StringAttribute{
				Optional:    true,
				Description: "Name of anti-spam profile.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"web_filtering_profile": schema.StringAttribute{
				Optional:    true,
				Description: "Web-filtering HTTP profile (local, enhanced, websense).",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"anti_virus": schema.SingleNestedBlock{
				Description: "Configure for utm anti-virus profile.",
				Attributes: map[string]schema.Attribute{
					"ftp_download_profile": schema.StringAttribute{
						Optional:    true,
						Description: "FTP download anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"ftp_upload_profile": schema.StringAttribute{
						Optional:    true,
						Description: "FTP upload anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"http_profile": schema.StringAttribute{
						Optional:    true,
						Description: "HTTP anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"imap_profile": schema.StringAttribute{
						Optional:    true,
						Description: "IMAP anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"pop3_profile": schema.StringAttribute{
						Optional:    true,
						Description: "POP3 anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"smtp_profile": schema.StringAttribute{
						Optional:    true,
						Description: "SMTP anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"content_filtering": schema.SingleNestedBlock{
				Description: "Configure for utm content-filtering profile.",
				Attributes: map[string]schema.Attribute{
					"ftp_download_profile": schema.StringAttribute{
						Optional:    true,
						Description: "FTP download content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"ftp_upload_profile": schema.StringAttribute{
						Optional:    true,
						Description: "FTP upload content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"http_profile": schema.StringAttribute{
						Optional:    true,
						Description: "HTTP content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"imap_profile": schema.StringAttribute{
						Optional:    true,
						Description: "HTTP content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"pop3_profile": schema.StringAttribute{
						Optional:    true,
						Description: "POP3 content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"smtp_profile": schema.StringAttribute{
						Optional:    true,
						Description: "SMTP content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"content_filtering_rule_set": schema.ListNestedBlock{
				Description: "UTM CF Rule Set.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "UTM CF Rule-set name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 29),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"rule": schema.ListNestedBlock{
							Description: "UTM CF Rule.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Required:    true,
										Description: "UTM CF Rule name.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 29),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"match_applications": schema.ListAttribute{
										ElementType: types.StringType,
										Required:    true,
										Description: "List of applications to be inspected.",
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
											listvalidator.NoNullValues(),
											listvalidator.ValueStringsAre(
												stringvalidator.OneOf("any", "ftp", "http", "imap", "pop3", "smtp"),
											),
										},
									},
									"match_direction": schema.StringAttribute{
										Required:    true,
										Description: "Direction of the content to be inspected.",
										Validators: []validator.String{
											stringvalidator.OneOf("any", "download", "upload"),
										},
									},
									"match_file_types": schema.ListAttribute{
										ElementType: types.StringType,
										Required:    true,
										Description: "List of file-types in match criteria.",
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
											listvalidator.NoNullValues(),
											listvalidator.ValueStringsAre(
												stringvalidator.LengthBetween(1, 32),
												tfvalidator.StringFormat(tfvalidator.DefaultFormat),
											),
										},
									},
									"then_action": schema.StringAttribute{
										Optional:    true,
										Description: "Configure then action.",
										Validators: []validator.String{
											stringvalidator.OneOf("block", "close-client", "close-client-and-server", "close-server", "no-action"),
										},
									},
									"then_notification_log": schema.BoolAttribute{
										Optional:    true,
										Description: "Generate security event if content is blocked by rule.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
								},
								Blocks: map[string]schema.Block{
									"then_notification_endpoint": schema.SingleNestedBlock{
										Description: "Endpoint notification options for the content filtering action taken.",
										Attributes: map[string]schema.Attribute{
											"custom_message": schema.StringAttribute{
												Optional:    true,
												Description: "Custom notification message.",
												Validators: []validator.String{
													stringvalidator.LengthBetween(1, 512),
													tfvalidator.StringDoubleQuoteExclusion(),
												},
											},
											"notify_mail_sender": schema.BoolAttribute{
												Optional:    true,
												Description: "Notify mail sender.",
											},
											"type": schema.StringAttribute{
												Optional:    true,
												Description: "Endpoint notification type.",
												Validators: []validator.String{
													stringvalidator.OneOf("message", "protocol-only"),
												},
											},
										},
										PlanModifiers: []planmodifier.Object{
											tfplanmodifier.BlockRemoveNull(),
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.IsRequired(),
								listvalidator.SizeAtLeast(1),
							},
						},
					},
				},
			},
			"traffic_sessions_per_client": schema.SingleNestedBlock{
				Description: " Configure for traffic option session per client.",
				Attributes: map[string]schema.Attribute{
					"limit": schema.Int64Attribute{
						Optional:    true,
						Description: "Sessions limit.",
						Validators: []validator.Int64{
							int64validator.Between(0, 2000),
						},
					},
					"over_limit": schema.StringAttribute{
						Optional:    true,
						Description: "Over limit action.",
						Validators: []validator.String{
							stringvalidator.OneOf("block", "log-and-permit"),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
		},
	}
}

//nolint:lll
type securityUtmPolicyData struct {
	ID                       types.String                                    `tfsdk:"id"                          tfdata:"skip_isempty"`
	Name                     types.String                                    `tfsdk:"name"                        tfdata:"skip_isempty"`
	AntiSpamSMTPProfile      types.String                                    `tfsdk:"anti_spam_smtp_profile"`
	WebFilteringProfile      types.String                                    `tfsdk:"web_filtering_profile"`
	AntiVirus                *securityUtmPolicyBlockProtocolProfile          `tfsdk:"anti_virus"`
	ContentFiltering         *securityUtmPolicyBlockProtocolProfile          `tfsdk:"content_filtering"`
	ContentFilteringRuleSet  []securityUtmPolicyBlockContentFilteringRuleSet `tfsdk:"content_filtering_rule_set"`
	TrafficSessionsPerClient *securityUtmPolicyBlockTrafficSessionsPerClient `tfsdk:"traffic_sessions_per_client"`
}

func (rscData *securityUtmPolicyData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

//nolint:lll
type securityUtmPolicyConfig struct {
	ID                       types.String                                    `tfsdk:"id"                          tfdata:"skip_isempty"`
	Name                     types.String                                    `tfsdk:"name"                        tfdata:"skip_isempty"`
	AntiSpamSMTPProfile      types.String                                    `tfsdk:"anti_spam_smtp_profile"`
	WebFilteringProfile      types.String                                    `tfsdk:"web_filtering_profile"`
	AntiVirus                *securityUtmPolicyBlockProtocolProfile          `tfsdk:"anti_virus"`
	ContentFiltering         *securityUtmPolicyBlockProtocolProfile          `tfsdk:"content_filtering"`
	ContentFilteringRuleSet  types.List                                      `tfsdk:"content_filtering_rule_set"`
	TrafficSessionsPerClient *securityUtmPolicyBlockTrafficSessionsPerClient `tfsdk:"traffic_sessions_per_client"`
}

func (config *securityUtmPolicyConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(config)
}

type securityUtmPolicyBlockProtocolProfile struct {
	FTPDownloadProfile types.String `tfsdk:"ftp_download_profile"`
	FTPUploadProfile   types.String `tfsdk:"ftp_upload_profile"`
	HTTPProfile        types.String `tfsdk:"http_profile"`
	IMAPProfile        types.String `tfsdk:"imap_profile"`
	POP3Profile        types.String `tfsdk:"pop3_profile"`
	SMTPProfile        types.String `tfsdk:"smtp_profile"`
}

func (block *securityUtmPolicyBlockProtocolProfile) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *securityUtmPolicyBlockProtocolProfile) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityUtmPolicyBlockContentFilteringRuleSet struct {
	Name types.String                                             `tfsdk:"name" tfdata:"identifier"`
	Rule []securityUtmPolicyBlockContentFilteringRuleSetBlockRule `tfsdk:"rule"`
}

type securityUtmPolicyBlockContentFilteringRuleSetConfig struct {
	Name types.String `tfsdk:"name"`
	Rule types.List   `tfsdk:"rule"`
}

//nolint:lll
type securityUtmPolicyBlockContentFilteringRuleSetBlockRule struct {
	Name                     types.String                                                                         `tfsdk:"name"                       tfdata:"identifier"`
	MatchApplications        []types.String                                                                       `tfsdk:"match_applications"`
	MatchDirection           types.String                                                                         `tfsdk:"match_direction"`
	MatchFileTypes           []types.String                                                                       `tfsdk:"match_file_types"`
	ThenAction               types.String                                                                         `tfsdk:"then_action"`
	ThenNotificationLog      types.Bool                                                                           `tfsdk:"then_notification_log"`
	ThenNotificationEndpoint *securityUtmPolicyBlockContentFilteringRuleSetBlockRuleBlockThenNotificationEndpoint `tfsdk:"then_notification_endpoint"`
}

//nolint:lll
type securityUtmPolicyBlockContentFilteringRuleSetBlockRuleConfig struct {
	Name                     types.String                                                                         `tfsdk:"name"`
	MatchApplications        types.List                                                                           `tfsdk:"match_applications"`
	MatchDirection           types.String                                                                         `tfsdk:"match_direction"`
	MatchFileTypes           types.List                                                                           `tfsdk:"match_file_types"`
	ThenAction               types.String                                                                         `tfsdk:"then_action"`
	ThenNotificationLog      types.Bool                                                                           `tfsdk:"then_notification_log"`
	ThenNotificationEndpoint *securityUtmPolicyBlockContentFilteringRuleSetBlockRuleBlockThenNotificationEndpoint `tfsdk:"then_notification_endpoint"`
}

type securityUtmPolicyBlockContentFilteringRuleSetBlockRuleBlockThenNotificationEndpoint struct {
	CustomMessage    types.String `tfsdk:"custom_message"`
	NotifyMailSender types.Bool   `tfsdk:"notify_mail_sender"`
	Type             types.String `tfsdk:"type"`
}

func (block *securityUtmPolicyBlockContentFilteringRuleSetBlockRuleBlockThenNotificationEndpoint) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityUtmPolicyBlockTrafficSessionsPerClient struct {
	Limit     types.Int64  `tfsdk:"limit"`
	OverLimit types.String `tfsdk:"over_limit"`
}

func (block *securityUtmPolicyBlockTrafficSessionsPerClient) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (rsc *securityUtmPolicy) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityUtmPolicyConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one of arguments need to be set (in addition to `name`)",
		)
	}

	if config.AntiVirus != nil &&
		config.AntiVirus.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("anti_virus").AtName("*"),
			tfdiag.MissingConfigErrSummary,
			"anti_virus block is empty",
		)
	}
	if config.ContentFiltering != nil &&
		config.ContentFiltering.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("content_filtering").AtName("*"),
			tfdiag.MissingConfigErrSummary,
			"content_filtering block is empty",
		)
	}
	if config.ContentFiltering != nil &&
		config.ContentFiltering.hasKnownValue() &&
		!config.ContentFilteringRuleSet.IsNull() &&
		!config.ContentFilteringRuleSet.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("content_filtering").AtName("*"),
			tfdiag.ConflictConfigErrSummary,
			"content_filtering and content_filtering_rule_set cannot be configured together",
		)
	}
	if !config.ContentFilteringRuleSet.IsNull() &&
		!config.ContentFilteringRuleSet.IsUnknown() {
		var configContentFilteringRuleSet []securityUtmPolicyBlockContentFilteringRuleSetConfig
		asDiags := config.ContentFilteringRuleSet.ElementsAs(ctx, &configContentFilteringRuleSet, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		contentFilteringRuleSetName := make(map[string]struct{})
		for i, block := range configContentFilteringRuleSet {
			if !block.Name.IsUnknown() {
				ruleSetName := block.Name.ValueString()
				if _, ok := contentFilteringRuleSetName[ruleSetName]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("content_filtering_rule_set").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple content_filtering_rule_set blocks with the same name %q", ruleSetName),
					)
				}
				contentFilteringRuleSetName[ruleSetName] = struct{}{}
			}
			if !block.Rule.IsNull() &&
				!block.Rule.IsUnknown() {
				var configRule []securityUtmPolicyBlockContentFilteringRuleSetBlockRuleConfig
				asDiags := block.Rule.ElementsAs(ctx, &configRule, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				ruleName := make(map[string]struct{})
				for ii, subBlock := range configRule {
					if !subBlock.Name.IsUnknown() {
						name := subBlock.Name.ValueString()
						if _, ok := ruleName[name]; ok {
							resp.Diagnostics.AddAttributeError(
								path.Root("content_filtering_rule_set").AtListIndex(i).
									AtName("rule").AtListIndex(ii).
									AtName("name"),
								tfdiag.DuplicateConfigErrSummary,
								fmt.Sprintf("multiple rule blocks with the same name %q"+
									" in content_filtering_rule_set block %q", name, block.Name.ValueString()),
							)
						}
						ruleName[name] = struct{}{}
					}

					if !subBlock.MatchApplications.IsNull() &&
						!subBlock.MatchApplications.IsUnknown() {
						var matchApplications []types.String
						asDiags := subBlock.MatchApplications.ElementsAs(ctx, &matchApplications, false)
						if asDiags.HasError() {
							resp.Diagnostics.Append(asDiags...)

							return
						}

						for _, v := range matchApplications {
							if v.ValueString() == "any" {
								if len(matchApplications) != 1 {
									resp.Diagnostics.AddAttributeError(
										path.Root("content_filtering").AtListIndex(i).
											AtName("rule").AtListIndex(ii).
											AtName("match_applications"),
										tfdiag.ConflictConfigErrSummary,
										fmt.Sprintf("match_applications %q is configured, other applications are not supported"+
											" in rule block %q in content_filtering_rule_set block %q",
											v.ValueString(), subBlock.Name.ValueString(), block.Name.ValueString()),
									)
								}

								break
							}
						}
					}

					if subBlock.ThenAction.IsNull() &&
						subBlock.ThenNotificationLog.IsNull() &&
						subBlock.ThenNotificationEndpoint == nil {
						resp.Diagnostics.AddAttributeError(
							path.Root("content_filtering_rule_set").AtListIndex(i).
								AtName("rule").AtListIndex(ii).
								AtName("name"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("at least one then_* attribute must be specified"+
								" in rule block %q in content_filtering_rule_set block %q",
								subBlock.Name.ValueString(), block.Name.ValueString()),
						)
					}
					if subBlock.ThenNotificationEndpoint != nil &&
						subBlock.ThenNotificationEndpoint.isEmpty() {
						resp.Diagnostics.AddAttributeError(
							path.Root("content_filtering_rule_set").AtListIndex(i).
								AtName("rule").AtListIndex(ii).
								AtName("then_notification_endpoint"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("then_notification_endpoint block is empty"+
								" in rule block %q in content_filtering_rule_set block %q",
								subBlock.Name.ValueString(), block.Name.ValueString()),
						)
					}
				}
			}
		}
	}
	if config.TrafficSessionsPerClient != nil &&
		config.TrafficSessionsPerClient.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("traffic_sessions_per_client").AtName("*"),
			tfdiag.MissingConfigErrSummary,
			"traffic_sessions_per_client block is empty",
		)
	}
}

func (rsc *securityUtmPolicy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityUtmPolicyData
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
			policyExists, err := checkSecurityUtmPolicyExists(fnCtx, plan.Name.ValueString(), junSess)
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
			policyExists, err := checkSecurityUtmPolicyExists(fnCtx, plan.Name.ValueString(), junSess)
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

func (rsc *securityUtmPolicy) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityUtmPolicyData
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

func (rsc *securityUtmPolicy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityUtmPolicyData
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

func (rsc *securityUtmPolicy) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityUtmPolicyData
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

func (rsc *securityUtmPolicy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityUtmPolicyData

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

func checkSecurityUtmPolicyExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm utm-policy \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityUtmPolicyData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityUtmPolicyData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityUtmPolicyData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	if rscData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name`)")
	}

	configSet := make([]string, 0, 100)
	setPrefix := "set security utm utm-policy \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.AntiSpamSMTPProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"anti-spam smtp-profile \""+v+"\"")
	}
	if v := rscData.WebFilteringProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"web-filtering http-profile \""+v+"\"")
	}

	if rscData.AntiVirus != nil {
		if rscData.AntiVirus.isEmpty() {
			return path.Root("anti_virus").AtName("*"),
				errors.New("anti_virus block is empty")
		}

		configSet = append(configSet, rscData.AntiVirus.configSet(setPrefix+"anti-virus ")...)
	}
	if rscData.ContentFiltering != nil {
		if rscData.ContentFiltering.isEmpty() {
			return path.Root("content_filtering").AtName("*"),
				errors.New("content_filtering block is empty")
		}

		configSet = append(configSet, rscData.ContentFiltering.configSet(setPrefix+"content-filtering ")...)
	}
	contentFilteringRuleSetName := make(map[string]struct{})
	for i, block := range rscData.ContentFilteringRuleSet {
		name := block.Name.ValueString()
		if _, ok := contentFilteringRuleSetName[name]; ok {
			return path.Root("content_filtering_rule_set").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple content_filtering_rule_set blocks with the same name %q", name)
		}
		contentFilteringRuleSetName[name] = struct{}{}

		blockSet, pathErr, err := block.configSet(
			setPrefix,
			path.Root("content_filtering_rule_set").AtListIndex(i),
		)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.TrafficSessionsPerClient != nil {
		if rscData.TrafficSessionsPerClient.isEmpty() {
			return path.Root("traffic_sessions_per_client").AtName("*"),
				errors.New("traffic_sessions_per_client block is empty")
		}

		configSet = append(configSet, rscData.TrafficSessionsPerClient.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityUtmPolicyBlockProtocolProfile) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)

	if v := block.FTPDownloadProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ftp download-profile \""+v+"\"")
	}
	if v := block.FTPUploadProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ftp upload-profile \""+v+"\"")
	}
	if v := block.HTTPProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"http-profile \""+v+"\"")
	}
	if v := block.IMAPProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"imap-profile \""+v+"\"")
	}
	if v := block.POP3Profile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pop3-profile \""+v+"\"")
	}
	if v := block.SMTPProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"smtp-profile \""+v+"\"")
	}

	return configSet
}

func (block *securityUtmPolicyBlockContentFilteringRuleSet) configSet(
	setPrefix string,
	pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, 100)
	setPrefix += "content-filtering rule-set \"" + block.Name.ValueString() + "\" "

	ruleName := make(map[string]struct{})
	for i, v := range block.Rule {
		name := v.Name.ValueString()
		if _, ok := ruleName[name]; ok {
			return configSet,
				pathRoot.AtName("rule").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple rule blocks with the same name %q"+
					" in content_filtering_rule_set block %q", name, block.Name.ValueString())
		}
		ruleName[name] = struct{}{}

		blockSet, pathErr, err := v.configSet(
			setPrefix,
			pathRoot.AtName("rule").AtListIndex(i),
			fmt.Sprintf(" in content_filtering_rule_set block %q", block.Name.ValueString()),
		)
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return configSet, path.Empty(), nil
}

func (block *securityUtmPolicyBlockContentFilteringRuleSetBlockRule) configSet(
	setPrefix string,
	pathRoot path.Path,
	blockErrorSuffix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, 100)
	setPrefix += "rule \"" + block.Name.ValueString() + "\" "

	for _, v := range block.MatchApplications {
		configSet = append(configSet, setPrefix+"match applications "+v.ValueString())
	}
	configSet = append(configSet, setPrefix+"match direction "+block.MatchDirection.ValueString())
	for _, v := range block.MatchFileTypes {
		configSet = append(configSet, setPrefix+"match file-types "+v.ValueString())
	}
	if block.ThenAction.IsNull() &&
		block.ThenNotificationLog.IsNull() &&
		(block.ThenNotificationEndpoint == nil || block.ThenNotificationEndpoint.isEmpty()) {
		return configSet,
			pathRoot.AtName("name"),
			fmt.Errorf("at least one then_* attribute must be specified"+
				" in rule block %q"+blockErrorSuffix, block.Name.ValueString())
	}
	if v := block.ThenAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"then action "+v)
	}
	if block.ThenNotificationLog.ValueBool() {
		configSet = append(configSet, setPrefix+"then notification log")
	}
	if block.ThenNotificationEndpoint != nil {
		if block.ThenNotificationEndpoint.isEmpty() {
			return configSet,
				pathRoot.AtName("then_notification_endpoint").AtName("*"),
				fmt.Errorf("then_notification_endpoint block is empty"+
					" in rule block %q"+blockErrorSuffix, block.Name.ValueString())
		}

		configSet = append(configSet, block.ThenNotificationEndpoint.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (block *securityUtmPolicyBlockContentFilteringRuleSetBlockRuleBlockThenNotificationEndpoint) configSet(
	setPrefix string,
) []string {
	configSet := make([]string, 0, 100)
	setPrefix += "then notification endpoint "

	if v := block.CustomMessage.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"custom-message \""+v+"\"")
	}
	if !block.NotifyMailSender.IsNull() {
		if block.NotifyMailSender.ValueBool() {
			configSet = append(configSet, setPrefix+"notify-mail-sender")
		} else {
			configSet = append(configSet, setPrefix+"no-notify-mail-sender")
		}
	}
	if v := block.Type.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"type "+v)
	}

	return configSet
}

func (block *securityUtmPolicyBlockTrafficSessionsPerClient) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)
	setPrefix += "traffic-options sessions-per-client "

	if !block.Limit.IsNull() {
		configSet = append(configSet, setPrefix+"limit "+
			utils.ConvI64toa(block.Limit.ValueInt64()))
	}
	if v := block.OverLimit.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"over-limit "+v)
	}

	return configSet
}

func (rscData *securityUtmPolicyData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm utm-policy \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "anti-spam smtp-profile "):
				rscData.AntiSpamSMTPProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "anti-virus "):
				if rscData.AntiVirus == nil {
					rscData.AntiVirus = &securityUtmPolicyBlockProtocolProfile{}
				}

				rscData.AntiVirus.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "content-filtering rule-set "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.ContentFilteringRuleSet = tfdata.AppendPotentialNewBlock(
					rscData.ContentFilteringRuleSet, types.StringValue(strings.Trim(name, "\"")),
				)
				contentFilteringRuleSet := &rscData.ContentFilteringRuleSet[len(rscData.ContentFilteringRuleSet)-1]
				balt.CutPrefixInString(&itemTrim, name+" ")

				contentFilteringRuleSet.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "content-filtering "):
				if rscData.ContentFiltering == nil {
					rscData.ContentFiltering = &securityUtmPolicyBlockProtocolProfile{}
				}

				rscData.ContentFiltering.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "traffic-options sessions-per-client "):
				if rscData.TrafficSessionsPerClient == nil {
					rscData.TrafficSessionsPerClient = &securityUtmPolicyBlockTrafficSessionsPerClient{}
				}

				if err := rscData.TrafficSessionsPerClient.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "web-filtering http-profile "):
				rscData.WebFilteringProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			}
		}
	}

	return nil
}

func (block *securityUtmPolicyBlockProtocolProfile) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "ftp download-profile "):
		block.FTPDownloadProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "ftp upload-profile "):
		block.FTPUploadProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "http-profile "):
		block.HTTPProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "imap-profile "):
		block.IMAPProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "pop3-profile "):
		block.POP3Profile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "smtp-profile "):
		block.SMTPProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	}
}

func (block *securityUtmPolicyBlockContentFilteringRuleSet) read(itemTrim string) {
	if balt.CutPrefixInString(&itemTrim, "rule ") {
		name := tfdata.FirstElementOfJunosLine(itemTrim)
		block.Rule = tfdata.AppendPotentialNewBlock(block.Rule, types.StringValue(strings.Trim(name, "\"")))
		rule := &block.Rule[len(block.Rule)-1]
		balt.CutPrefixInString(&itemTrim, name+" ")

		rule.read(itemTrim)
	}
}

func (block *securityUtmPolicyBlockContentFilteringRuleSetBlockRule) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "match applications "):
		block.MatchApplications = append(block.MatchApplications, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "match direction "):
		block.MatchDirection = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "match file-types "):
		block.MatchFileTypes = append(block.MatchFileTypes, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "then action "):
		block.ThenAction = types.StringValue(itemTrim)
	case itemTrim == "then notification log":
		block.ThenNotificationLog = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "then notification endpoint "):
		if block.ThenNotificationEndpoint == nil {
			block.ThenNotificationEndpoint = &securityUtmPolicyBlockContentFilteringRuleSetBlockRuleBlockThenNotificationEndpoint{} //nolint:lll
		}

		block.ThenNotificationEndpoint.read(itemTrim)
	}
}

func (block *securityUtmPolicyBlockContentFilteringRuleSetBlockRuleBlockThenNotificationEndpoint) read(itemTrim string) { //nolint:lll
	switch {
	case balt.CutPrefixInString(&itemTrim, "custom-message "):
		block.CustomMessage = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "notify-mail-sender":
		block.NotifyMailSender = types.BoolValue(true)
	case itemTrim == "no-notify-mail-sender":
		block.NotifyMailSender = types.BoolValue(false)
	case balt.CutPrefixInString(&itemTrim, "type "):
		block.Type = types.StringValue(itemTrim)
	}
}

func (block *securityUtmPolicyBlockTrafficSessionsPerClient) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "limit "):
		block.Limit, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "over-limit "):
		block.OverLimit = types.StringValue(itemTrim)
	}

	return err
}

func (rscData *securityUtmPolicyData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security utm utm-policy \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
