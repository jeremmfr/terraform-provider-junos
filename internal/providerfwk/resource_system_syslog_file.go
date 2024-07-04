package providerfwk

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
	_ resource.Resource                   = &systemSyslogFile{}
	_ resource.ResourceWithConfigure      = &systemSyslogFile{}
	_ resource.ResourceWithValidateConfig = &systemSyslogFile{}
	_ resource.ResourceWithImportState    = &systemSyslogFile{}
	_ resource.ResourceWithUpgradeState   = &systemSyslogFile{}
)

type systemSyslogFile struct {
	client *junos.Client
}

func newSystemSyslogFileResource() resource.Resource {
	return &systemSyslogFile{}
}

func (rsc *systemSyslogFile) typeName() string {
	return providerName + "_system_syslog_file"
}

func (rsc *systemSyslogFile) junosName() string {
	return "system syslog file"
}

func (rsc *systemSyslogFile) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *systemSyslogFile) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *systemSyslogFile) Configure(
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

func (rsc *systemSyslogFile) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<filename>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"filename": schema.StringAttribute{
				Required:    true,
				Description: "Name of file in which to log data.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
					tfvalidator.StringSpaceExclusion(),
					tfvalidator.StringRuneExclusion('/', '%'),
				},
			},
			"allow_duplicates": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not suppress the repeated message.",
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
		Blocks: map[string]schema.Block{
			"archive": schema.SingleNestedBlock{
				Description: "Define parameters for archiving log messages.",
				Attributes: map[string]schema.Attribute{
					"binary_data": schema.BoolAttribute{
						Optional:    true,
						Description: "Mark file as if it contains binary data.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_binary_data": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't mark file as if it contains binary data.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"files": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of files to be archived.",
						Validators: []validator.Int64{
							int64validator.Between(1, 1000),
						},
					},
					"size": schema.Int64Attribute{
						Optional:    true,
						Description: "Size of files to be archived (bytes).",
						Validators: []validator.Int64{
							int64validator.Between(65536, 1073741824),
						},
					},
					"start_time": schema.StringAttribute{
						Optional:    true,
						Description: "Start time for file transmission (YYYY-MM-DD.HH:MM:SS).",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^\d{4}\-\d\d?\-\d\d?\.\d{2}:\d{2}:\d{2}$`),
								"must be in the format 'YYYY-MM-DD.HH:MM:SS'",
							),
						},
					},
					"transfer_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Frequency at which to transfer files to archive sites (minutes).",
						Validators: []validator.Int64{
							int64validator.Between(5, 2880),
						},
					},
					"world_readable": schema.BoolAttribute{
						Optional:    true,
						Description: "Allow any user to read the log file.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_world_readable": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't allow any user to read the log file.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"sites": schema.ListNestedBlock{
						Description: "For each url, configure an archive site.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"url": schema.StringAttribute{
									Required:    true,
									Description: "Primary or failover URLs to receive archive files.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"password": schema.StringAttribute{
									Optional:    true,
									Sensitive:   true,
									Description: "Password for login into the archive site.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"routing_instance": schema.StringAttribute{
									Optional:    true,
									Description: "Routing instance.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 63),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									},
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
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

type systemSyslogFileData struct {
	ID                          types.String                         `tfsdk:"id"`
	Filename                    types.String                         `tfsdk:"filename"`
	AllowDuplicates             types.Bool                           `tfsdk:"allow_duplicates"`
	ExplicitPriority            types.Bool                           `tfsdk:"explicit_priority"`
	Match                       types.String                         `tfsdk:"match"`
	MatchStrings                []types.String                       `tfsdk:"match_strings"`
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
	Archive                     *systemSyslogFileBlockArchive        `tfsdk:"archive"`
	StructuredData              *systemSyslogFileBlockStructuredData `tfsdk:"structured_data"`
}

type systemSyslogFileConfig struct {
	ID                          types.String                         `tfsdk:"id"`
	Filename                    types.String                         `tfsdk:"filename"`
	AllowDuplicates             types.Bool                           `tfsdk:"allow_duplicates"`
	ExplicitPriority            types.Bool                           `tfsdk:"explicit_priority"`
	Match                       types.String                         `tfsdk:"match"`
	MatchStrings                types.List                           `tfsdk:"match_strings"`
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
	Archive                     *systemSyslogFileBlockArchiveConfig  `tfsdk:"archive"`
	StructuredData              *systemSyslogFileBlockStructuredData `tfsdk:"structured_data"`
}

type systemSyslogFileBlockArchive struct {
	BinaryData       types.Bool                               `tfsdk:"binary_data"`
	NoBinaryData     types.Bool                               `tfsdk:"no_binary_data"`
	Files            types.Int64                              `tfsdk:"files"`
	Size             types.Int64                              `tfsdk:"size"`
	StartTime        types.String                             `tfsdk:"start_time"`
	TransferInterval types.Int64                              `tfsdk:"transfer_interval"`
	WorldReadable    types.Bool                               `tfsdk:"world_readable"`
	NoWorldReadable  types.Bool                               `tfsdk:"no_world_readable"`
	Sites            []systemSyslogFileBlockArchiveBlockSites `tfsdk:"sites"`
}

type systemSyslogFileBlockArchiveConfig struct {
	BinaryData       types.Bool   `tfsdk:"binary_data"`
	NoBinaryData     types.Bool   `tfsdk:"no_binary_data"`
	Files            types.Int64  `tfsdk:"files"`
	Size             types.Int64  `tfsdk:"size"`
	StartTime        types.String `tfsdk:"start_time"`
	TransferInterval types.Int64  `tfsdk:"transfer_interval"`
	WorldReadable    types.Bool   `tfsdk:"world_readable"`
	NoWorldReadable  types.Bool   `tfsdk:"no_world_readable"`
	Sites            types.List   `tfsdk:"sites"`
}

type systemSyslogFileBlockArchiveBlockSites struct {
	URL             types.String `tfsdk:"url"`
	Password        types.String `tfsdk:"password"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
}

type systemSyslogFileBlockStructuredData struct {
	Brief types.Bool `tfsdk:"brief"`
}

func (rsc *systemSyslogFile) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config systemSyslogFileConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Archive != nil {
		if !config.Archive.BinaryData.IsNull() && !config.Archive.BinaryData.IsUnknown() &&
			!config.Archive.NoBinaryData.IsNull() && !config.Archive.NoBinaryData.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("archive").AtName("binary_data"),
				tfdiag.ConflictConfigErrSummary,
				"binary_data and no_binary_data cannot be configured together"+
					" in archive block",
			)
		}
		if !config.Archive.WorldReadable.IsNull() && !config.Archive.WorldReadable.IsUnknown() &&
			!config.Archive.NoWorldReadable.IsNull() && !config.Archive.NoWorldReadable.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("archive").AtName("world_readable"),
				tfdiag.ConflictConfigErrSummary,
				"world_readable and no_world_readable cannot be configured together"+
					" in archive block",
			)
		}
		if config.Archive.Sites.IsNull() {
			if !config.Archive.StartTime.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("archive").AtName("start_time"),
					tfdiag.MissingConfigErrSummary,
					"sites must be specified with start_time in archive block",
				)
			}
			if !config.Archive.TransferInterval.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("archive").AtName("transfer_interval"),
					tfdiag.MissingConfigErrSummary,
					"sites must be specified with transfer_interval in archive block",
				)
			}
		} else if !config.Archive.Sites.IsUnknown() {
			var configSites []systemSyslogFileBlockArchiveBlockSites
			asDiags := config.Archive.Sites.ElementsAs(ctx, &configSites, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			sitesURL := make(map[string]struct{})
			for i, blockSites := range configSites {
				if blockSites.URL.IsUnknown() {
					continue
				}
				url := blockSites.URL.ValueString()
				if _, ok := sitesURL[url]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("archive").AtName("sites").AtListIndex(i).AtName("url"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple sites blocks with the same url %q in archive block", url),
					)
				}
				sitesURL[url] = struct{}{}
			}
		}
	}
}

func (rsc *systemSyslogFile) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemSyslogFileData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Filename.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("filename"),
			"Empty Filename",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "filename"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			fileExists, err := checkSystemSyslogFileExists(fnCtx, plan.Filename.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if fileExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Filename),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			fileExists, err := checkSystemSyslogFileExists(fnCtx, plan.Filename.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !fileExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Filename),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *systemSyslogFile) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemSyslogFileData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Filename.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *systemSyslogFile) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemSyslogFileData
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

func (rsc *systemSyslogFile) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state systemSyslogFileData
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

func (rsc *systemSyslogFile) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data systemSyslogFileData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "filename"),
	)
}

func checkSystemSyslogFileExists(
	_ context.Context, filename string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system syslog file \"" + filename + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *systemSyslogFileData) fillID() {
	rscData.ID = types.StringValue(rscData.Filename.ValueString())
}

func (rscData *systemSyslogFileData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemSyslogFileData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set system syslog file \"" + rscData.Filename.ValueString() + "\" "

	if rscData.AllowDuplicates.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-duplicates")
	}
	if rscData.ExplicitPriority.ValueBool() {
		configSet = append(configSet, setPrefix+"explicit-priority")
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
	if rscData.Archive != nil {
		configSet = append(configSet, setPrefix+"archive")

		if rscData.Archive.BinaryData.ValueBool() {
			configSet = append(configSet, setPrefix+"archive binary-data")
		}
		if rscData.Archive.NoBinaryData.ValueBool() {
			configSet = append(configSet, setPrefix+"archive no-binary-data")
		}
		if !rscData.Archive.Files.IsNull() {
			configSet = append(configSet, setPrefix+"archive files "+
				utils.ConvI64toa(rscData.Archive.Files.ValueInt64()))
		}
		if !rscData.Archive.Size.IsNull() {
			configSet = append(configSet, setPrefix+"archive size "+
				utils.ConvI64toa(rscData.Archive.Size.ValueInt64()))
		}
		if v := rscData.Archive.StartTime.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"archive start-time "+v)
		}
		if !rscData.Archive.TransferInterval.IsNull() {
			configSet = append(configSet, setPrefix+"archive transfer-interval "+
				utils.ConvI64toa(rscData.Archive.TransferInterval.ValueInt64()))
		}
		if rscData.Archive.WorldReadable.ValueBool() {
			configSet = append(configSet, setPrefix+"archive world-readable")
		}
		if rscData.Archive.NoWorldReadable.ValueBool() {
			configSet = append(configSet, setPrefix+"archive no-world-readable")
		}
		sitesURL := make(map[string]struct{})
		for i, blockSites := range rscData.Archive.Sites {
			url := blockSites.URL.ValueString()
			if _, ok := sitesURL[url]; ok {
				return path.Root("archive").AtName("sites").AtListIndex(i).AtName("url"),
					fmt.Errorf("multiple sites blocks with the same url %q in archive block",
						url)
			}
			sitesURL[url] = struct{}{}

			setPrefixArchiveSites := setPrefix + "archive archive-sites \"" + url + "\""
			configSet = append(configSet, setPrefixArchiveSites)
			if v := blockSites.Password.ValueString(); v != "" {
				configSet = append(configSet, setPrefixArchiveSites+" password \""+v+"\"")
			}
			if v := blockSites.RoutingInstance.ValueString(); v != "" {
				configSet = append(configSet, setPrefixArchiveSites+" routing-instance "+v)
			}
		}
	}
	if rscData.StructuredData != nil {
		configSet = append(configSet, setPrefix+"structured-data")
		if rscData.StructuredData.Brief.ValueBool() {
			configSet = append(configSet, setPrefix+"structured-data brief")
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *systemSyslogFileData) read(
	_ context.Context, filename string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system syslog file \"" + filename + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Filename = types.StringValue(filename)
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
			case itemTrim == "explicit-priority":
				rscData.ExplicitPriority = types.BoolValue(true)
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
			case balt.CutPrefixInString(&itemTrim, "archive"):
				if rscData.Archive == nil {
					rscData.Archive = &systemSyslogFileBlockArchive{}
				}
				switch {
				case itemTrim == " binary-data":
					rscData.Archive.BinaryData = types.BoolValue(true)
				case itemTrim == " no-binary-data":
					rscData.Archive.NoBinaryData = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, " files "):
					rscData.Archive.Files, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " size "):
					switch {
					case balt.CutSuffixInString(&itemTrim, "k"):
						rscData.Archive.Size, err = tfdata.ConvAtoi64Value(itemTrim)
						rscData.Archive.Size = types.Int64Value(rscData.Archive.Size.ValueInt64() * 1024)
					case balt.CutSuffixInString(&itemTrim, "m"):
						rscData.Archive.Size, err = tfdata.ConvAtoi64Value(itemTrim)
						rscData.Archive.Size = types.Int64Value(rscData.Archive.Size.ValueInt64() * 1024 * 1024)
					case balt.CutSuffixInString(&itemTrim, "g"):
						rscData.Archive.Size, err = tfdata.ConvAtoi64Value(itemTrim)
						rscData.Archive.Size = types.Int64Value(rscData.Archive.Size.ValueInt64() * 1024 * 1024 * 1024)
					default:
						rscData.Archive.Size, err = tfdata.ConvAtoi64Value(itemTrim)
					}
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " transfer-interval "):
					rscData.Archive.TransferInterval, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " start-time "):
					rscData.Archive.StartTime = types.StringValue(strings.Split(strings.Trim(itemTrim, "\""), " ")[0])
				case itemTrim == " world-readable":
					rscData.Archive.WorldReadable = types.BoolValue(true)
				case itemTrim == " no-world-readable":
					rscData.Archive.NoWorldReadable = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, " archive-sites "):
					url := tfdata.FirstElementOfJunosLine(itemTrim)
					var sites systemSyslogFileBlockArchiveBlockSites
					rscData.Archive.Sites, sites = tfdata.ExtractBlockWithTFTypesString(
						rscData.Archive.Sites, "URL", strings.Trim(url, "\""),
					)
					sites.URL = types.StringValue(strings.Trim(url, "\""))
					balt.CutPrefixInString(&itemTrim, url+" ")
					switch {
					case balt.CutPrefixInString(&itemTrim, "password "):
						sites.Password, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "password")
						if err != nil {
							return err
						}
					case balt.CutPrefixInString(&itemTrim, "routing-instance "):
						sites.RoutingInstance = types.StringValue(itemTrim)
					}
					rscData.Archive.Sites = append(rscData.Archive.Sites, sites)
				}
			case balt.CutPrefixInString(&itemTrim, "structured-data"):
				if rscData.StructuredData == nil {
					rscData.StructuredData = &systemSyslogFileBlockStructuredData{}
				}
				if itemTrim == " brief" {
					rscData.StructuredData.Brief = types.BoolValue(true)
				}
			}
		}
	}

	return nil
}

func (rscData *systemSyslogFileData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete system syslog file \"" + rscData.Filename.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
