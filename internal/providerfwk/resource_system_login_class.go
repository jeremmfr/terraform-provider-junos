package providerfwk

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                   = &systemLoginClass{}
	_ resource.ResourceWithConfigure      = &systemLoginClass{}
	_ resource.ResourceWithValidateConfig = &systemLoginClass{}
	_ resource.ResourceWithImportState    = &systemLoginClass{}
)

type systemLoginClass struct {
	client *junos.Client
}

func newSystemLoginClassResource() resource.Resource {
	return &systemLoginClass{}
}

func (rsc *systemLoginClass) typeName() string {
	return providerName + "_system_login_class"
}

func (rsc *systemLoginClass) junosName() string {
	return "system login class"
}

func (rsc *systemLoginClass) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *systemLoginClass) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *systemLoginClass) Configure(
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

func (rsc *systemLoginClass) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
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
				Description: "The name of system login class.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"access_end": schema.StringAttribute{
				Optional:    true,
				Description: "End time for remote access.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^([0-1]\d|2[0-3]):([0-5]\d):([0-5]\d)$`),
						"must be in the format 'HH:MM:SS'",
					),
				},
			},
			"access_start": schema.StringAttribute{
				Optional:    true,
				Description: "Start time for remote access.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^([0-1]\d|2[0-3]):([0-5]\d):([0-5]\d)$`),
						"must be in the format 'HH:MM:SS'",
					),
				},
			},
			"allow_commands": schema.StringAttribute{
				Optional:    true,
				Description: "Regular expression for commands to allow explicitly.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"allow_commands_regexps": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Object path regular expressions to allow commands.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"allow_configuration": schema.StringAttribute{
				Optional:    true,
				Description: "Regular expression for configure to allow explicitly.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"allow_configuration_regexps": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Object path regular expressions to allow.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"allow_hidden_commands": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow all hidden commands to be executed.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"allowed_days": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Day(s) of week when access is allowed.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday",
						),
					),
				},
			},
			"cli_prompt": schema.StringAttribute{
				Optional:    true,
				Description: "Cli prompt name for this class.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"configuration_breadcrumbs": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable breadcrumbs during display of configuration.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"confirm_commands": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of commands to be confirmed explicitly.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"deny_commands": schema.StringAttribute{
				Optional:    true,
				Description: "Regular expression for commands to deny explicitly.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"deny_commands_regexps": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Object path regular expressions to deny commands.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"deny_configuration": schema.StringAttribute{
				Optional:    true,
				Description: "Regular expression for configure to deny explicitly.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"deny_configuration_regexps": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Object path regular expressions to deny.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"idle_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum idle time before logout (minutes).",
				Validators: []validator.Int64{
					int64validator.Between(1, 4294967295),
				},
			},
			"logical_system": schema.StringAttribute{
				Optional:    true,
				Description: "Logical system associated with login.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"login_alarms": schema.BoolAttribute{
				Optional:    true,
				Description: "Display system alarms when logging in.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"login_script": schema.StringAttribute{
				Optional:    true,
				Description: "Execute this login-script when logging in.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"login_tip": schema.BoolAttribute{
				Optional:    true,
				Description: "Display tip when logging in.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_hidden_commands_except": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Deny all hidden commands with exemptions.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"permissions": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Set of permitted operation categories.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.NoNullValues(),
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"access", "access-control",
							"admin", "admin-control",
							"all",
							"clear",
							"configure",
							"control",
							"field",
							"firewall", "firewall-control",
							"floppy",
							"flow-tap", "flow-tap-control", "flow-tap-operation",
							"idp-profiler-operation",
							"interface", "interface-control",
							"maintenance",
							"network",
							"pgcp-session-mirroring", "pgcp-session-mirroring-control",
							"reset",
							"rollback",
							"routing", "routing-control",
							"secret", "secret-control",
							"security", "security-control",
							"shell",
							"snmp", "snmp-control",
							"storage", "storage-control",
							"system", "system-control",
							"trace", "trace-control",
							"unified-edge", "unified-edge-control",
							"view", "view-configuration",
						),
					),
				},
			},
			"security_role": schema.StringAttribute{
				Optional:    true,
				Description: "Common Criteria security role.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"audit-administrator",
						"crypto-administrator",
						"ids-administrator",
						"security-administrator",
					),
				},
			},
			"tenant": schema.StringAttribute{
				Optional:    true,
				Description: "Tenant associated with this login.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
	}
}

type systemLoginClassData struct {
	ID                        types.String   `tfsdk:"id"                          tfdata:"skip_isempty"`
	Name                      types.String   `tfsdk:"name"                        tfdata:"skip_isempty"`
	AccessEnd                 types.String   `tfsdk:"access_end"`
	AccessStart               types.String   `tfsdk:"access_start"`
	AllowCommands             types.String   `tfsdk:"allow_commands"`
	AllowCommandsRegexps      []types.String `tfsdk:"allow_commands_regexps"`
	AllowConfiguration        types.String   `tfsdk:"allow_configuration"`
	AllowConfigurationRegexps []types.String `tfsdk:"allow_configuration_regexps"`
	AllowHiddenCommands       types.Bool     `tfsdk:"allow_hidden_commands"`
	AllowedDays               []types.String `tfsdk:"allowed_days"`
	CliPrompt                 types.String   `tfsdk:"cli_prompt"`
	ConfigurationBreadcrumbs  types.Bool     `tfsdk:"configuration_breadcrumbs"`
	ConfirmCommands           []types.String `tfsdk:"confirm_commands"`
	DenyCommands              types.String   `tfsdk:"deny_commands"`
	DenyCommandsRegexps       []types.String `tfsdk:"deny_commands_regexps"`
	DenyConfiguration         types.String   `tfsdk:"deny_configuration"`
	DenyConfigurationRegexps  []types.String `tfsdk:"deny_configuration_regexps"`
	IdleTimeout               types.Int64    `tfsdk:"idle_timeout"`
	LogicalSystem             types.String   `tfsdk:"logical_system"`
	LoginAlarms               types.Bool     `tfsdk:"login_alarms"`
	LoginScript               types.String   `tfsdk:"login_script"`
	LoginTip                  types.Bool     `tfsdk:"login_tip"`
	NoHiddenCommandsExcept    []types.String `tfsdk:"no_hidden_commands_except"`
	Permissions               []types.String `tfsdk:"permissions"`
	SecurityRole              types.String   `tfsdk:"security_role"`
	Tenant                    types.String   `tfsdk:"tenant"`
}

func (rscData *systemLoginClassData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

type systemLoginClassConfig struct {
	ID                        types.String `tfsdk:"id"                          tfdata:"skip_isempty"`
	Name                      types.String `tfsdk:"name"                        tfdata:"skip_isempty"`
	AccessEnd                 types.String `tfsdk:"access_end"`
	AccessStart               types.String `tfsdk:"access_start"`
	AllowCommands             types.String `tfsdk:"allow_commands"`
	AllowCommandsRegexps      types.List   `tfsdk:"allow_commands_regexps"`
	AllowConfiguration        types.String `tfsdk:"allow_configuration"`
	AllowConfigurationRegexps types.List   `tfsdk:"allow_configuration_regexps"`
	AllowHiddenCommands       types.Bool   `tfsdk:"allow_hidden_commands"`
	AllowedDays               types.List   `tfsdk:"allowed_days"`
	CliPrompt                 types.String `tfsdk:"cli_prompt"`
	ConfigurationBreadcrumbs  types.Bool   `tfsdk:"configuration_breadcrumbs"`
	ConfirmCommands           types.List   `tfsdk:"confirm_commands"`
	DenyCommands              types.String `tfsdk:"deny_commands"`
	DenyCommandsRegexps       types.List   `tfsdk:"deny_commands_regexps"`
	DenyConfiguration         types.String `tfsdk:"deny_configuration"`
	DenyConfigurationRegexps  types.List   `tfsdk:"deny_configuration_regexps"`
	IdleTimeout               types.Int64  `tfsdk:"idle_timeout"`
	LogicalSystem             types.String `tfsdk:"logical_system"`
	LoginAlarms               types.Bool   `tfsdk:"login_alarms"`
	LoginScript               types.String `tfsdk:"login_script"`
	LoginTip                  types.Bool   `tfsdk:"login_tip"`
	NoHiddenCommandsExcept    types.List   `tfsdk:"no_hidden_commands_except"`
	Permissions               types.Set    `tfsdk:"permissions"`
	SecurityRole              types.String `tfsdk:"security_role"`
	Tenant                    types.String `tfsdk:"tenant"`
}

func (config *systemLoginClassConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(config)
}

func (rsc *systemLoginClass) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config systemLoginClassConfig
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

	if !config.AccessEnd.IsNull() &&
		!config.AccessEnd.IsUnknown() &&
		config.AccessStart.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_end"),
			tfdiag.MissingConfigErrSummary,
			"access_start must be specified with access_end",
		)
	}
	if !config.AccessStart.IsNull() &&
		!config.AccessStart.IsUnknown() &&
		config.AccessEnd.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_start"),
			tfdiag.MissingConfigErrSummary,
			"access_end must be specified with access_start",
		)
	}
	if !config.AllowCommands.IsNull() &&
		!config.AllowCommands.IsUnknown() &&
		!config.AllowCommandsRegexps.IsNull() &&
		!config.AllowCommandsRegexps.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("allow_commands"),
			tfdiag.ConflictConfigErrSummary,
			"allow_commands and allow_commands_regexps cannot be configured together",
		)
	}
	if !config.AllowConfiguration.IsNull() &&
		!config.AllowConfiguration.IsUnknown() &&
		!config.AllowConfigurationRegexps.IsNull() &&
		!config.AllowConfigurationRegexps.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("allow_configuration"),
			tfdiag.ConflictConfigErrSummary,
			"allow_configuration and allow_configuration_regexps cannot be configured together",
		)
	}
	if !config.AllowHiddenCommands.IsNull() &&
		!config.AllowHiddenCommands.IsUnknown() &&
		!config.NoHiddenCommandsExcept.IsNull() &&
		!config.NoHiddenCommandsExcept.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("allow_hidden_commands"),
			tfdiag.ConflictConfigErrSummary,
			"allow_hidden_commands and no_hidden_commands_except cannot be configured together",
		)
	}
	if !config.DenyCommands.IsNull() &&
		!config.DenyCommands.IsUnknown() &&
		!config.DenyCommandsRegexps.IsNull() &&
		!config.DenyCommandsRegexps.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("deny_commands"),
			tfdiag.ConflictConfigErrSummary,
			"deny_commands and deny_commands_regexps cannot be configured together",
		)
	}
	if !config.DenyConfiguration.IsNull() &&
		!config.DenyConfiguration.IsUnknown() &&
		!config.DenyConfigurationRegexps.IsNull() &&
		!config.DenyConfigurationRegexps.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("deny_configuration"),
			tfdiag.ConflictConfigErrSummary,
			"deny_configuration and deny_configuration_regexps cannot be configured together",
		)
	}
}

func (rsc *systemLoginClass) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemLoginClassData
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
			classExists, err := checkSystemLoginClassExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if classExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			classExists, err := checkSystemLoginClassExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !classExists {
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

func (rsc *systemLoginClass) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemLoginClassData
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

func (rsc *systemLoginClass) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemLoginClassData
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

func (rsc *systemLoginClass) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state systemLoginClassData
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

func (rsc *systemLoginClass) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data systemLoginClassData

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

func checkSystemLoginClassExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system login class " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *systemLoginClassData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *systemLoginClassData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemLoginClassData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)
	setPrefix := "set system login class " + rscData.Name.ValueString() + " "

	if rscData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name`)")
	}

	if v := rscData.AccessEnd.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"access-end \""+v+"\"")
	}
	if v := rscData.AccessStart.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"access-start \""+v+"\"")
	}
	if v := rscData.AllowCommands.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"allow-commands \""+v+"\"")
	}
	for _, v := range rscData.AllowCommandsRegexps {
		configSet = append(configSet, setPrefix+"allow-commands-regexps \""+v.ValueString()+"\"")
	}
	if v := rscData.AllowConfiguration.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"allow-configuration \""+v+"\"")
	}
	for _, v := range rscData.AllowConfigurationRegexps {
		configSet = append(configSet, setPrefix+"allow-configuration-regexps \""+v.ValueString()+"\"")
	}
	if rscData.AllowHiddenCommands.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-hidden-commands")
	}
	for _, v := range rscData.AllowedDays {
		configSet = append(configSet, setPrefix+"allowed-days "+v.ValueString())
	}
	if rscData.ConfigurationBreadcrumbs.ValueBool() {
		configSet = append(configSet, setPrefix+"configuration-breadcrumbs")
	}
	if v := rscData.CliPrompt.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"cli prompt \""+v+"\"")
	}
	for _, v := range rscData.ConfirmCommands {
		configSet = append(configSet, setPrefix+"confirm-commands \""+v.ValueString()+"\"")
	}
	if v := rscData.DenyCommands.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"deny-commands \""+v+"\"")
	}
	for _, v := range rscData.DenyCommandsRegexps {
		configSet = append(configSet, setPrefix+"deny-commands-regexps \""+v.ValueString()+"\"")
	}
	if v := rscData.DenyConfiguration.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"deny-configuration \""+v+"\"")
	}
	for _, v := range rscData.DenyConfigurationRegexps {
		configSet = append(configSet, setPrefix+"deny-configuration-regexps \""+v.ValueString()+"\"")
	}
	if !rscData.IdleTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"idle-timeout "+
			utils.ConvI64toa(rscData.IdleTimeout.ValueInt64()))
	}
	if v := rscData.LogicalSystem.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"logical-system \""+v+"\"")
	}
	if rscData.LoginAlarms.ValueBool() {
		configSet = append(configSet, setPrefix+"login-alarms")
	}
	if v := rscData.LoginScript.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"login-script \""+v+"\"")
	}
	if rscData.LoginTip.ValueBool() {
		configSet = append(configSet, setPrefix+"login-tip")
	}
	for _, v := range rscData.NoHiddenCommandsExcept {
		configSet = append(configSet, setPrefix+"no-hidden-commands except \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.Permissions {
		configSet = append(configSet, setPrefix+"permissions "+v.ValueString())
	}
	if v := rscData.SecurityRole.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"security-role "+v)
	}
	if v := rscData.Tenant.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"tenant \""+v+"\"")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *systemLoginClassData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system login class " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "access-end "):
				rscData.AccessEnd = types.StringValue(strings.Split(strings.Trim(itemTrim, "\""), " ")[0])
			case balt.CutPrefixInString(&itemTrim, "access-start "):
				rscData.AccessStart = types.StringValue(strings.Split(strings.Trim(itemTrim, "\""), " ")[0])
			case balt.CutPrefixInString(&itemTrim, "allow-commands "):
				rscData.AllowCommands = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "allow-commands-regexps "):
				rscData.AllowCommandsRegexps = append(rscData.AllowCommandsRegexps,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "allow-configuration "):
				rscData.AllowConfiguration = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "allow-configuration-regexps "):
				rscData.AllowConfigurationRegexps = append(rscData.AllowConfigurationRegexps,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case itemTrim == "allow-hidden-commands":
				rscData.AllowHiddenCommands = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "allowed-days "):
				rscData.AllowedDays = append(rscData.AllowedDays,
					types.StringValue(itemTrim))
			case itemTrim == "configuration-breadcrumbs":
				rscData.ConfigurationBreadcrumbs = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "cli prompt "):
				rscData.CliPrompt = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "confirm-commands "):
				rscData.ConfirmCommands = append(rscData.ConfirmCommands,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "deny-commands "):
				rscData.DenyCommands = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "deny-commands-regexps "):
				rscData.DenyCommandsRegexps = append(rscData.DenyCommandsRegexps,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "deny-configuration "):
				rscData.DenyConfiguration = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "deny-configuration-regexps "):
				rscData.DenyConfigurationRegexps = append(rscData.DenyConfigurationRegexps,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "idle-timeout "):
				rscData.IdleTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "logical-system "):
				rscData.LogicalSystem = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "login-alarms":
				rscData.LoginAlarms = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "login-script "):
				rscData.LoginScript = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "login-tip":
				rscData.LoginTip = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "no-hidden-commands except "):
				rscData.NoHiddenCommandsExcept = append(rscData.NoHiddenCommandsExcept,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "permissions "):
				rscData.Permissions = append(rscData.Permissions,
					types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "security-role "):
				rscData.SecurityRole = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "tenant "):
				rscData.Tenant = types.StringValue(strings.Trim(itemTrim, "\""))
			}
		}
	}

	return nil
}

func (rscData *systemLoginClassData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete system login class " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
