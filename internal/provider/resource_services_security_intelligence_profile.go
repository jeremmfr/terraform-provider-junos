package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
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
	_ resource.Resource                   = &servicesSecurityIntelligenceProfile{}
	_ resource.ResourceWithConfigure      = &servicesSecurityIntelligenceProfile{}
	_ resource.ResourceWithValidateConfig = &servicesSecurityIntelligenceProfile{}
	_ resource.ResourceWithImportState    = &servicesSecurityIntelligenceProfile{}
	_ resource.ResourceWithUpgradeState   = &servicesSecurityIntelligenceProfile{}
)

type servicesSecurityIntelligenceProfile struct {
	client *junos.Client
}

func newServicesSecurityIntelligenceProfileResource() resource.Resource {
	return &servicesSecurityIntelligenceProfile{}
}

func (rsc *servicesSecurityIntelligenceProfile) typeName() string {
	return providerName + "_services_security_intelligence_profile"
}

func (rsc *servicesSecurityIntelligenceProfile) junosName() string {
	return "services security-intelligence profile"
}

func (rsc *servicesSecurityIntelligenceProfile) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *servicesSecurityIntelligenceProfile) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *servicesSecurityIntelligenceProfile) Configure(
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

func (rsc *servicesSecurityIntelligenceProfile) Schema(
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
				Description: "Security intelligence profile name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"category": schema.StringAttribute{
				Required:    true,
				Description: "Profile category name.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of profile.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				Description: "For each rule name.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Profile rule name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"then_action": schema.StringAttribute{
							Required:    true,
							Description: "Security intelligence profile action.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(
									`^(permit|recommended|sinkhole|block (drop|close( http (file|message|redirect-url) .+)?))$`),
									"must have valid action (permit|recommended|block...)"),
							},
						},
						"then_log": schema.BoolAttribute{
							Optional:    true,
							Description: "Log security intelligence block action.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"match": schema.SingleNestedBlock{
							Description: "Configure profile matching feed name and threat levels.",
							Attributes: map[string]schema.Attribute{
								"threat_level": schema.ListAttribute{
									ElementType: types.Int64Type,
									Required:    true,
									Description: "Profile matching threat levels, higher number is more severe.",
									Validators: []validator.List{
										listvalidator.SizeAtLeast(1),
										listvalidator.NoNullValues(),
										listvalidator.ValueInt64sAre(
											int64validator.Between(1, 10),
										),
									},
								},
								"feed_name": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Profile matching feed name.",
									Validators: []validator.List{
										listvalidator.SizeAtLeast(1),
										listvalidator.NoNullValues(),
										listvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 63),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
							},
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
							Validators: []validator.Object{
								objectvalidator.IsRequired(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.IsRequired(),
					listvalidator.SizeAtLeast(1),
				},
			},
			"default_rule_then": schema.SingleNestedBlock{
				Description: "Declare profile default rule.",
				Attributes: map[string]schema.Attribute{
					"action": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Security intelligence profile action.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^(permit|recommended|sinkhole|block (drop|close( http (file|message|redirect-url) .+)?))$`),
								"must have valid action (permit|recommended|block...)"),
						},
					},
					"log": schema.BoolAttribute{
						Optional:    true,
						Description: "Log security intelligence block action.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_log": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't log security intelligence block action.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
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

type servicesSecurityIntelligenceProfileData struct {
	ID              types.String                                             `tfsdk:"id"`
	Name            types.String                                             `tfsdk:"name"`
	Category        types.String                                             `tfsdk:"category"`
	Description     types.String                                             `tfsdk:"description"`
	Rule            []servicesSecurityIntelligenceProfileBlockRule           `tfsdk:"rule"`
	DefaultRuleThen *servicesSecurityIntelligenceProfileBlockDefaultRuleThen `tfsdk:"default_rule_then"`
}

type servicesSecurityIntelligenceProfileConfig struct {
	ID              types.String                                             `tfsdk:"id"`
	Name            types.String                                             `tfsdk:"name"`
	Category        types.String                                             `tfsdk:"category"`
	Description     types.String                                             `tfsdk:"description"`
	Rule            types.List                                               `tfsdk:"rule"`
	DefaultRuleThen *servicesSecurityIntelligenceProfileBlockDefaultRuleThen `tfsdk:"default_rule_then"`
}

type servicesSecurityIntelligenceProfileBlockRule struct {
	Name       types.String                                            `tfsdk:"name"        tfdata:"identifier"`
	ThenAction types.String                                            `tfsdk:"then_action"`
	ThenLog    types.Bool                                              `tfsdk:"then_log"`
	Match      *servicesSecurityIntelligenceProfileBlockRuleBlockMatch `tfsdk:"match"`
}

type servicesSecurityIntelligenceProfileBlockRuleBlockMatch struct {
	ThreatLevel []types.Int64  `tfsdk:"threat_level"`
	FeedName    []types.String `tfsdk:"feed_name"`
}

type servicesSecurityIntelligenceProfileBlockDefaultRuleThen struct {
	Action types.String `tfsdk:"action"`
	Log    types.Bool   `tfsdk:"log"`
	NoLog  types.Bool   `tfsdk:"no_log"`
}

func (rsc *servicesSecurityIntelligenceProfile) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config servicesSecurityIntelligenceProfileConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Rule.IsNull() &&
		!config.Rule.IsUnknown() {
		var configRule []servicesSecurityIntelligenceProfileBlockRule
		asDiags := config.Rule.ElementsAs(ctx, &configRule, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		ruleName := make(map[string]struct{})
		for i, block := range configRule {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := ruleName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("rule").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple rule blocks with the same name %q", name),
					)
				}
				ruleName[name] = struct{}{}
			}

			if !config.Category.IsNull() &&
				!config.Category.IsUnknown() &&
				!block.ThenAction.IsNull() &&
				!block.ThenAction.IsUnknown() {
				if block.ThenAction.ValueString() == "sinkhole" &&
					config.Category.ValueString() != "DNS" {
					resp.Diagnostics.AddAttributeError(
						path.Root("rule").AtListIndex(i).AtName("then_action"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("sinkhole action requires DNS category"+
							" in rules block %q", block.Name.ValueString()),
					)
				}
			}
		}
	}

	if config.DefaultRuleThen != nil {
		if config.DefaultRuleThen.Action.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("default_rule_then").AtName("action"),
				tfdiag.MissingConfigErrSummary,
				"action must be specified"+
					" in default_rule_then block",
			)
		}
		if !config.Category.IsNull() &&
			!config.Category.IsUnknown() &&
			!config.DefaultRuleThen.Action.IsNull() &&
			!config.DefaultRuleThen.Action.IsUnknown() {
			if config.DefaultRuleThen.Action.ValueString() == "sinkhole" &&
				config.Category.ValueString() != "DNS" {
				resp.Diagnostics.AddAttributeError(
					path.Root("default_rule_then").AtName("action"),
					tfdiag.ConflictConfigErrSummary,
					"sinkhole action requires DNS category"+
						" in default_rule_then block",
				)
			}
		}
		if !config.DefaultRuleThen.Log.IsNull() &&
			!config.DefaultRuleThen.Log.IsUnknown() &&
			!config.DefaultRuleThen.NoLog.IsNull() &&
			!config.DefaultRuleThen.NoLog.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("default_rule_then").AtName("log"),
				tfdiag.ConflictConfigErrSummary,
				"log and no_log cannot be configured together"+
					" in default_rule_then block",
			)
		}
	}
}

func (rsc *servicesSecurityIntelligenceProfile) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan servicesSecurityIntelligenceProfileData
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
			profileExists, err := checkServicesSecurityIntelligenceProfileExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if profileExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			profileExists, err := checkServicesSecurityIntelligenceProfileExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !profileExists {
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

func (rsc *servicesSecurityIntelligenceProfile) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data servicesSecurityIntelligenceProfileData
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

func (rsc *servicesSecurityIntelligenceProfile) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state servicesSecurityIntelligenceProfileData
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

func (rsc *servicesSecurityIntelligenceProfile) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state servicesSecurityIntelligenceProfileData
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

func (rsc *servicesSecurityIntelligenceProfile) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data servicesSecurityIntelligenceProfileData

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

func checkServicesSecurityIntelligenceProfileExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services security-intelligence profile \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *servicesSecurityIntelligenceProfileData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *servicesSecurityIntelligenceProfileData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *servicesSecurityIntelligenceProfileData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set services security-intelligence profile \"" + rscData.Name.ValueString() + "\" "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix + "category " + rscData.Category.ValueString()

	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}

	ruleName := make(map[string]struct{})
	for i, block := range rscData.Rule {
		name := block.Name.ValueString()
		if _, ok := ruleName[name]; ok {
			return path.Root("rule").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple rule blocks with the same name %q", name)
		}
		ruleName[name] = struct{}{}

		configSet = append(configSet, block.configSet(setPrefix)...)
	}
	if rscData.DefaultRuleThen != nil {
		configSet = append(configSet, rscData.DefaultRuleThen.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *servicesSecurityIntelligenceProfileBlockRule) configSet(setPrefix string) []string {
	setPrefix += "rule \"" + block.Name.ValueString() + "\" "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix + "then action " + block.ThenAction.ValueString()

	if block.ThenLog.ValueBool() {
		configSet = append(configSet, setPrefix+"then log")
	}

	if block.Match != nil {
		for _, v := range block.Match.ThreatLevel {
			configSet = append(configSet, setPrefix+"match threat-level "+
				utils.ConvI64toa(v.ValueInt64()))
		}
		for _, v := range block.Match.FeedName {
			configSet = append(configSet, setPrefix+"match feed-name \""+v.ValueString()+"\"")
		}
	}

	return configSet
}

func (block *servicesSecurityIntelligenceProfileBlockDefaultRuleThen) configSet(setPrefix string) []string {
	setPrefix += "default-rule then "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix + "action " + block.Action.ValueString()

	if block.Log.ValueBool() {
		configSet = append(configSet, setPrefix+"log")
	}
	if block.NoLog.ValueBool() {
		configSet = append(configSet, setPrefix+"no-log")
	}

	return configSet
}

func (rscData *servicesSecurityIntelligenceProfileData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services security-intelligence profile \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.fillID()
		for item := range strings.SplitSeq(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "category "):
				rscData.Category = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "default-rule then "):
				if rscData.DefaultRuleThen == nil {
					rscData.DefaultRuleThen = &servicesSecurityIntelligenceProfileBlockDefaultRuleThen{}
				}

				rscData.DefaultRuleThen.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "rule "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.Rule = tfdata.AppendPotentialNewBlock(rscData.Rule, types.StringValue(strings.Trim(name, "\"")))
				rule := &rscData.Rule[len(rscData.Rule)-1]
				balt.CutPrefixInString(&itemTrim, name+" ")

				if err := rule.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *servicesSecurityIntelligenceProfileBlockRule) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "match "):
		if block.Match == nil {
			block.Match = &servicesSecurityIntelligenceProfileBlockRuleBlockMatch{}
		}

		switch {
		case balt.CutPrefixInString(&itemTrim, "threat-level "):
			threatLevel, err := tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
			block.Match.ThreatLevel = append(block.Match.ThreatLevel, threatLevel)
		case balt.CutPrefixInString(&itemTrim, "feed-name "):
			block.Match.FeedName = append(block.Match.FeedName, types.StringValue(strings.Trim(itemTrim, "\"")))
		}
	case balt.CutPrefixInString(&itemTrim, "then action "):
		block.ThenAction = types.StringValue(itemTrim)
	case itemTrim == "then log":
		block.ThenLog = types.BoolValue(true)
	}

	return nil
}

func (block *servicesSecurityIntelligenceProfileBlockDefaultRuleThen) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "action "):
		block.Action = types.StringValue(itemTrim)
	case itemTrim == "log":
		block.Log = types.BoolValue(true)
	case itemTrim == "no-log":
		block.NoLog = types.BoolValue(true)
	}
}

func (rscData *servicesSecurityIntelligenceProfileData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete services security-intelligence profile \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
