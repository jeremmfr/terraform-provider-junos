package providerfwk

import (
	"context"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

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
	_ resource.Resource                = &systemSyslogUser{}
	_ resource.ResourceWithConfigure   = &systemSyslogUser{}
	_ resource.ResourceWithImportState = &systemSyslogUser{}
)

type systemSyslogUser struct {
	client *junos.Client
}

func newSystemSyslogUserResource() resource.Resource {
	return &systemSyslogUser{}
}

func (rsc *systemSyslogUser) typeName() string {
	return providerName + "_system_syslog_user"
}

func (rsc *systemSyslogUser) junosName() string {
	return "system syslog user"
}

func (rsc *systemSyslogUser) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *systemSyslogUser) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *systemSyslogUser) Configure(
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

func (rsc *systemSyslogUser) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<username>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "Name of user to notify (or `*` for all).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					stringvalidator.Any(
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
						stringvalidator.OneOf("*"),
					),
				},
			},
			"allow_duplicates": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not suppress the repeated message.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"match": schema.StringAttribute{
				Optional:    true,
				Description: "Regular expression for lines to be logged.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"match_strings": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Matching string(s) for lines to be logged.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"any_severity": schema.StringAttribute{
				Optional:    true,
				Description: "All facilities sseverity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"authorization_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Authorization system severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"changelog_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Configuration change log severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"conflictlog_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Configuration conflict log severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"daemon_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Various system processes severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"dfc_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Dynamic flow capture severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"external_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Local external applications severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"firewall_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Firewall filtering system severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"ftp_severity": schema.StringAttribute{
				Optional:    true,
				Description: "FTP process severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"interactivecommands_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Commands executed by the UI severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"kernel_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Kernel severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"ntp_severity": schema.StringAttribute{
				Optional:    true,
				Description: "NTP process severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"pfe_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Packet Forwarding Engine severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"security_severity": schema.StringAttribute{
				Optional:    true,
				Description: "Security related severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
			"user_severity": schema.StringAttribute{
				Optional:    true,
				Description: "User processes severity.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
		},
	}
}

type systemSyslogUserData struct {
	ID                          types.String   `tfsdk:"id"`
	Username                    types.String   `tfsdk:"username"`
	AllowDuplicates             types.Bool     `tfsdk:"allow_duplicates"`
	Match                       types.String   `tfsdk:"match"`
	MatchStrings                []types.String `tfsdk:"match_strings"`
	AnySeverity                 types.String   `tfsdk:"any_severity"`
	AuthorizationSeverity       types.String   `tfsdk:"authorization_severity"`
	ChangelogSeverity           types.String   `tfsdk:"changelog_severity"`
	ConflictlogSeverity         types.String   `tfsdk:"conflictlog_severity"`
	DaemonSeverity              types.String   `tfsdk:"daemon_severity"`
	DfcSeverity                 types.String   `tfsdk:"dfc_severity"`
	ExternalSeverity            types.String   `tfsdk:"external_severity"`
	FirewallSeverity            types.String   `tfsdk:"firewall_severity"`
	FtpSeverity                 types.String   `tfsdk:"ftp_severity"`
	InteractivecommandsSeverity types.String   `tfsdk:"interactivecommands_severity"`
	KernelSeverity              types.String   `tfsdk:"kernel_severity"`
	NtpSeverity                 types.String   `tfsdk:"ntp_severity"`
	PfeSeverity                 types.String   `tfsdk:"pfe_severity"`
	SecuritySeverity            types.String   `tfsdk:"security_severity"`
	UserSeverity                types.String   `tfsdk:"user_severity"`
}

func (rsc *systemSyslogUser) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemSyslogUserData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Username.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Empty Username",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "username"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			userExists, err := checkSystemSyslogUserExists(fnCtx, plan.Username.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if userExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Username),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			userExists, err := checkSystemSyslogUserExists(fnCtx, plan.Username.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !userExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Username),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *systemSyslogUser) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemSyslogUserData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Username.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *systemSyslogUser) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemSyslogUserData
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

func (rsc *systemSyslogUser) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state systemSyslogUserData
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

func (rsc *systemSyslogUser) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data systemSyslogUserData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "username"),
	)
}

func checkSystemSyslogUserExists(
	_ context.Context, username string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system syslog user " + username + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *systemSyslogUserData) fillID() {
	rscData.ID = types.StringValue(rscData.Username.ValueString())
}

func (rscData *systemSyslogUserData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemSyslogUserData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set system syslog user " + rscData.Username.ValueString() + " "

	if rscData.AllowDuplicates.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-duplicates")
	}
	if v := rscData.Match.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"match \""+v+"\"")
	}
	for _, v := range rscData.MatchStrings {
		configSet = append(configSet, setPrefix+"match-strings \""+v.ValueString()+"\"")
	}
	if v := rscData.AnySeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"any "+v)
	}
	if v := rscData.AuthorizationSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authorization "+v)
	}
	if v := rscData.ChangelogSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"change-log "+v)
	}
	if v := rscData.ConflictlogSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"conflict-log "+v)
	}
	if v := rscData.DaemonSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"daemon "+v)
	}
	if v := rscData.DfcSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"dfc "+v)
	}
	if v := rscData.ExternalSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"external "+v)
	}
	if v := rscData.FirewallSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"firewall "+v)
	}
	if v := rscData.FtpSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ftp "+v)
	}
	if v := rscData.InteractivecommandsSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"interactive-commands "+v)
	}
	if v := rscData.KernelSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"kernel "+v)
	}
	if v := rscData.NtpSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ntp "+v)
	}
	if v := rscData.PfeSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pfe "+v)
	}
	if v := rscData.SecuritySeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"security "+v)
	}
	if v := rscData.UserSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"user "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *systemSyslogUserData) read(
	_ context.Context, username string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system syslog user " + username + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Username = types.StringValue(username)
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
			case itemTrim == "allow-duplicates":
				rscData.AllowDuplicates = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "match "):
				rscData.Match = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "match-strings "):
				rscData.MatchStrings = append(rscData.MatchStrings,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "any "):
				rscData.AnySeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "authorization "):
				rscData.AuthorizationSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "change-log "):
				rscData.ChangelogSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "conflict-log "):
				rscData.ConflictlogSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "daemon "):
				rscData.DaemonSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "dfc "):
				rscData.DfcSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "external "):
				rscData.ExternalSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "firewall "):
				rscData.FirewallSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "ftp "):
				rscData.FtpSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "interactive-commands "):
				rscData.InteractivecommandsSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "kernel "):
				rscData.KernelSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "ntp "):
				rscData.NtpSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "pfe "):
				rscData.PfeSeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "security "):
				rscData.SecuritySeverity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "user "):
				rscData.UserSeverity = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (rscData *systemSyslogUserData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete system syslog user " + rscData.Username.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
