package providerfwk

import (
	"context"
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
	_ resource.Resource                 = &systemSyslogHost{}
	_ resource.ResourceWithConfigure    = &systemSyslogHost{}
	_ resource.ResourceWithImportState  = &systemSyslogHost{}
	_ resource.ResourceWithUpgradeState = &systemSyslogHost{}
)

type systemSyslogHost struct {
	client *junos.Client
}

func newSystemSyslogHostResource() resource.Resource {
	return &systemSyslogHost{}
}

func (rsc *systemSyslogHost) typeName() string {
	return providerName + "_system_syslog_host"
}

func (rsc *systemSyslogHost) junosName() string {
	return "system syslog host"
}

func (rsc *systemSyslogHost) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *systemSyslogHost) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *systemSyslogHost) Configure(
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

func (rsc *systemSyslogHost) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<host>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Required:    true,
				Description: "Host to be notified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
				},
			},
			"allow_duplicates": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not suppress the repeated message.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"exclude_hostname": schema.BoolAttribute{
				Optional:    true,
				Description: "Exclude hostname field in messages.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"explicit_priority": schema.BoolAttribute{
				Optional:    true,
				Description: "Include priority and facility in messages.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"facility_override": schema.StringAttribute{
				Optional:    true,
				Description: "Alternate facility for logging to remote host.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogFacilities()...),
				},
			},
			"log_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Prefix for all logging to this host.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
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
			"port": schema.Int64Attribute{
				Optional:    true,
				Description: "Port number.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"source_address": schema.StringAttribute{
				Optional:    true,
				Description: "Use specified address as source address.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
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
		Blocks: map[string]schema.Block{
			"structured_data": schema.SingleNestedBlock{
				Description: "Log system message in structured format.",
				Attributes: map[string]schema.Attribute{
					"brief": schema.BoolAttribute{
						Optional:    true,
						Description: "Omit English-language text from end of logged message.",
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

type systemSyslogHostData struct {
	ID                          types.String                         `tfsdk:"id"`
	Host                        types.String                         `tfsdk:"host"`
	AllowDuplicates             types.Bool                           `tfsdk:"allow_duplicates"`
	ExcludeHostname             types.Bool                           `tfsdk:"exclude_hostname"`
	ExplicitPriority            types.Bool                           `tfsdk:"explicit_priority"`
	FacilityOverride            types.String                         `tfsdk:"facility_override"`
	LogPrefix                   types.String                         `tfsdk:"log_prefix"`
	Match                       types.String                         `tfsdk:"match"`
	MatchStrings                []types.String                       `tfsdk:"match_strings"`
	Port                        types.Int64                          `tfsdk:"port"`
	SourceAddress               types.String                         `tfsdk:"source_address"`
	AnySeverity                 types.String                         `tfsdk:"any_severity"`
	AuthorizationSeverity       types.String                         `tfsdk:"authorization_severity"`
	ChangelogSeverity           types.String                         `tfsdk:"changelog_severity"`
	ConflictlogSeverity         types.String                         `tfsdk:"conflictlog_severity"`
	DaemonSeverity              types.String                         `tfsdk:"daemon_severity"`
	DfcSeverity                 types.String                         `tfsdk:"dfc_severity"`
	ExternalSeverity            types.String                         `tfsdk:"external_severity"`
	FirewallSeverity            types.String                         `tfsdk:"firewall_severity"`
	FtpSeverity                 types.String                         `tfsdk:"ftp_severity"`
	InteractivecommandsSeverity types.String                         `tfsdk:"interactivecommands_severity"`
	KernelSeverity              types.String                         `tfsdk:"kernel_severity"`
	NtpSeverity                 types.String                         `tfsdk:"ntp_severity"`
	PfeSeverity                 types.String                         `tfsdk:"pfe_severity"`
	SecuritySeverity            types.String                         `tfsdk:"security_severity"`
	UserSeverity                types.String                         `tfsdk:"user_severity"`
	StructuredData              *systemSyslogHostBlockStructuredData `tfsdk:"structured_data"`
}

type systemSyslogHostBlockStructuredData struct {
	Brief types.Bool `tfsdk:"brief"`
}

func (rsc *systemSyslogHost) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemSyslogHostData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Host.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Empty Host",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "host"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			hostExists, err := checkSystemSyslogHostExists(fnCtx, plan.Host.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if hostExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Host),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			hostExists, err := checkSystemSyslogHostExists(fnCtx, plan.Host.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !hostExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Host),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *systemSyslogHost) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemSyslogHostData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Host.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *systemSyslogHost) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemSyslogHostData
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

func (rsc *systemSyslogHost) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state systemSyslogHostData
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

func (rsc *systemSyslogHost) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data systemSyslogHostData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "host"),
	)
}

func checkSystemSyslogHostExists(
	_ context.Context, host string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system syslog host " + host + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *systemSyslogHostData) fillID() {
	rscData.ID = types.StringValue(rscData.Host.ValueString())
}

func (rscData *systemSyslogHostData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemSyslogHostData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set system syslog host " + rscData.Host.ValueString() + " "

	if rscData.AllowDuplicates.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-duplicates")
	}
	if rscData.ExcludeHostname.ValueBool() {
		configSet = append(configSet, setPrefix+"exclude-hostname")
	}
	if rscData.ExplicitPriority.ValueBool() {
		configSet = append(configSet, setPrefix+"explicit-priority")
	}
	if v := rscData.FacilityOverride.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"facility-override "+v)
	}
	if v := rscData.LogPrefix.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"log-prefix \""+v+"\"")
	}
	if v := rscData.Match.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"match \""+v+"\"")
	}
	for _, v := range rscData.MatchStrings {
		configSet = append(configSet, setPrefix+"match-strings \""+v.ValueString()+"\"")
	}
	if !rscData.Port.IsNull() {
		configSet = append(configSet, setPrefix+"port "+
			utils.ConvI64toa(rscData.Port.ValueInt64()))
	}
	if v := rscData.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
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
	if rscData.StructuredData != nil {
		configSet = append(configSet, setPrefix+"structured-data")
		if rscData.StructuredData.Brief.ValueBool() {
			configSet = append(configSet, setPrefix+"structured-data brief")
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *systemSyslogHostData) read(
	_ context.Context, host string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system syslog host " + host + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Host = types.StringValue(host)
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
			case itemTrim == "exclude-hostname":
				rscData.ExcludeHostname = types.BoolValue(true)
			case itemTrim == "explicit-priority":
				rscData.ExplicitPriority = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "facility-override "):
				rscData.FacilityOverride = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "log-prefix "):
				rscData.LogPrefix = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "match "):
				rscData.Match = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "match-strings "):
				rscData.MatchStrings = append(rscData.MatchStrings,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "port "):
				rscData.Port, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "source-address "):
				rscData.SourceAddress = types.StringValue(itemTrim)
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
			case balt.CutPrefixInString(&itemTrim, "structured-data"):
				if rscData.StructuredData == nil {
					rscData.StructuredData = &systemSyslogHostBlockStructuredData{}
				}
				if itemTrim == " brief" {
					rscData.StructuredData.Brief = types.BoolValue(true)
				}
			}
		}
	}

	return nil
}

func (rscData *systemSyslogHostData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete system syslog host " + rscData.Host.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
