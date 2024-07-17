package providerfwk

import (
	"context"
	"errors"
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
	_ resource.Resource                   = &securityLogStream{}
	_ resource.ResourceWithConfigure      = &securityLogStream{}
	_ resource.ResourceWithValidateConfig = &securityLogStream{}
	_ resource.ResourceWithImportState    = &securityLogStream{}
	_ resource.ResourceWithUpgradeState   = &securityLogStream{}
)

type securityLogStream struct {
	client *junos.Client
}

func newSecurityLogStreamResource() resource.Resource {
	return &securityLogStream{}
}

func (rsc *securityLogStream) typeName() string {
	return providerName + "_security_log_stream"
}

func (rsc *securityLogStream) junosName() string {
	return "security log stream"
}

func (rsc *securityLogStream) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityLogStream) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityLogStream) Configure(
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

func (rsc *securityLogStream) Schema(
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
				Description: "Name of security log stream.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"category": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Selects the type of events that may be logged.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
			"filter_threat_attack": schema.BoolAttribute{
				Optional:    true,
				Description: "Threat-attack security events are logged.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"format": schema.StringAttribute{
				Optional:    true,
				Description: "Specify the log stream format.",
				Validators: []validator.String{
					stringvalidator.OneOf("binary", "sd-syslog", "syslog", "welf"),
				},
			},
			"rate_limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Rate-limit for security logs.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"severity": schema.StringAttribute{
				Optional:    true,
				Description: "Severity threshold for security logs.",
				Validators: []validator.String{
					stringvalidator.OneOf(junos.SyslogSeverity()...),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"file": schema.SingleNestedBlock{
				Description: "Configure log file options for logs in local file.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Name of local log file.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 250),
							tfvalidator.StringDoubleQuoteExclusion(),
							tfvalidator.StringSpaceExclusion(),
							tfvalidator.StringRuneExclusion('/', '%'),
						},
					},
					"allow_duplicates": schema.BoolAttribute{
						Optional:    true,
						Description: "To disable log consolidation.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"rotation": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum number of rotate files.",
						Validators: []validator.Int64{
							int64validator.Between(2, 19),
						},
					},
					"size": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum size of local log file in megabytes.",
						Validators: []validator.Int64{
							int64validator.Between(1, 3),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"host": schema.SingleNestedBlock{
				Description: "Configure destination to send security logs to.",
				Attributes: map[string]schema.Attribute{
					"ip_address": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "IP address.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 250),
							tfvalidator.StringDoubleQuoteExclusion(),
							tfvalidator.StringSpaceExclusion(),
						},
					},
					"port": schema.Int64Attribute{
						Optional:    true,
						Description: "Host port number.",
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"routing_instance": schema.StringAttribute{
						Optional:    true,
						Description: "Routing instance name.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 63),
							tfvalidator.StringFormat(tfvalidator.DefaultFormat),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"transport": schema.SingleNestedBlock{
				Description: "Set security log transport settings.",
				Attributes: map[string]schema.Attribute{
					"protocol": schema.StringAttribute{
						Optional:    true,
						Description: "Set security log transport protocol for the device.",
						Validators: []validator.String{
							stringvalidator.OneOf("tcp", "tls", "udp"),
						},
					},
					"tcp_connections": schema.Int64Attribute{
						Optional:    true,
						Description: "Set tcp connection number per-stream.",
						Validators: []validator.Int64{
							int64validator.Between(1, 5),
						},
					},
					"tls_profile": schema.StringAttribute{
						Optional:    true,
						Description: "TLS profile.",
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
		},
	}
}

type securityLogStreamData struct {
	ID                 types.String                     `tfsdk:"id"`
	Name               types.String                     `tfsdk:"name"`
	Category           []types.String                   `tfsdk:"category"`
	FilterThreatAttack types.Bool                       `tfsdk:"filter_threat_attack"`
	Format             types.String                     `tfsdk:"format"`
	RateLimit          types.Int64                      `tfsdk:"rate_limit"`
	Severity           types.String                     `tfsdk:"severity"`
	File               *securityLogStreamBlockFile      `tfsdk:"file"`
	Host               *securityLogStreamBlockHost      `tfsdk:"host"`
	Transport          *securityLogStreamBlockTransport `tfsdk:"transport"`
}

type securityLogStreamConfig struct {
	ID                 types.String                     `tfsdk:"id"`
	Name               types.String                     `tfsdk:"name"`
	Category           types.List                       `tfsdk:"category"`
	FilterThreatAttack types.Bool                       `tfsdk:"filter_threat_attack"`
	Format             types.String                     `tfsdk:"format"`
	RateLimit          types.Int64                      `tfsdk:"rate_limit"`
	Severity           types.String                     `tfsdk:"severity"`
	File               *securityLogStreamBlockFile      `tfsdk:"file"`
	Host               *securityLogStreamBlockHost      `tfsdk:"host"`
	Transport          *securityLogStreamBlockTransport `tfsdk:"transport"`
}

type securityLogStreamBlockFile struct {
	Name            types.String `tfsdk:"name"`
	AllowDuplicates types.Bool   `tfsdk:"allow_duplicates"`
	Rotation        types.Int64  `tfsdk:"rotation"`
	Size            types.Int64  `tfsdk:"size"`
}

func (block *securityLogStreamBlockFile) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityLogStreamBlockHost struct {
	IPAddress       types.String `tfsdk:"ip_address"`
	Port            types.Int64  `tfsdk:"port"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
}

func (block *securityLogStreamBlockHost) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityLogStreamBlockTransport struct {
	Protocol       types.String `tfsdk:"protocol"`
	TCPConnections types.Int64  `tfsdk:"tcp_connections"`
	TLSProfile     types.String `tfsdk:"tls_profile"`
}

func (rsc *securityLogStream) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityLogStreamConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Category.IsNull() && !config.Category.IsUnknown() &&
		!config.FilterThreatAttack.IsNull() && !config.FilterThreatAttack.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("category"),
			tfdiag.ConflictConfigErrSummary,
			"category and filter_threat_attack cannot be configured together",
		)
	}
	if config.File != nil && config.File.hasKnownValue() &&
		config.Host != nil && config.Host.hasKnownValue() {
		resp.Diagnostics.AddAttributeError(
			path.Root("file"),
			tfdiag.ConflictConfigErrSummary,
			"file and host cannot be configured together",
		)
	}
	if config.File != nil && config.File.Name.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("file").AtName("name"),
			tfdiag.MissingConfigErrSummary,
			"name must be specified in file block",
		)
	}
	if config.Host != nil && config.Host.IPAddress.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host").AtName("ip_address"),
			tfdiag.MissingConfigErrSummary,
			"ip_address must be specified in host block",
		)
	}
	if config.Transport != nil &&
		!config.Transport.Protocol.IsNull() && !config.Transport.Protocol.IsUnknown() &&
		config.Transport.Protocol.ValueString() != "tls" &&
		!config.Transport.TLSProfile.IsNull() && !config.Transport.TLSProfile.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("transport").AtName("protocol"),
			tfdiag.ConflictConfigErrSummary,
			"protocol must be 'tls' with tls_profile"+
				" in transport block",
		)
	}
}

func (rsc *securityLogStream) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityLogStreamData
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
			logStreamExists, err := checkSecurityLogStreamExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if logStreamExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			logStreamExists, err := checkSecurityLogStreamExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !logStreamExists {
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

func (rsc *securityLogStream) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityLogStreamData
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

func (rsc *securityLogStream) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityLogStreamData
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

func (rsc *securityLogStream) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityLogStreamData
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

func (rsc *securityLogStream) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityLogStreamData

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

func checkSecurityLogStreamExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security log stream " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityLogStreamData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityLogStreamData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityLogStreamData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set security log stream " + rscData.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	for _, v := range rscData.Category {
		configSet = append(configSet, setPrefix+"category "+v.ValueString())
	}
	if rscData.FilterThreatAttack.ValueBool() {
		configSet = append(configSet, setPrefix+"filter threat-attack")
	}
	if v := rscData.Format.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"format "+v)
	}
	if !rscData.RateLimit.IsNull() {
		configSet = append(configSet, setPrefix+"rate-limit "+
			utils.ConvI64toa(rscData.RateLimit.ValueInt64()))
	}
	if v := rscData.Severity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"severity "+v)
	}
	if rscData.File != nil {
		configSet = append(configSet, setPrefix+"file name \""+rscData.File.Name.ValueString()+"\"")
		if rscData.File.AllowDuplicates.ValueBool() {
			configSet = append(configSet, setPrefix+"file allow-duplicates")
		}
		if !rscData.File.Rotation.IsNull() {
			configSet = append(configSet, setPrefix+"file rotation "+
				utils.ConvI64toa(rscData.File.Rotation.ValueInt64()))
		}
		if !rscData.File.Size.IsNull() {
			configSet = append(configSet, setPrefix+"file size "+
				utils.ConvI64toa(rscData.File.Size.ValueInt64()))
		}
	}
	if rscData.Host != nil {
		configSet = append(configSet, setPrefix+"host \""+rscData.Host.IPAddress.ValueString()+"\"")
		if !rscData.Host.Port.IsNull() {
			configSet = append(configSet, setPrefix+"host port "+
				utils.ConvI64toa(rscData.Host.Port.ValueInt64()))
		}
		if v := rscData.Host.RoutingInstance.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"host routing-instance "+v)
		}
	}
	if rscData.Transport != nil {
		configSet = append(configSet, setPrefix+"transport")
		if v := rscData.Transport.Protocol.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"transport protocol "+v)
			if v != "tls" && rscData.Transport.TLSProfile.ValueString() != "" {
				return path.Root("transport").AtName("protocol"),
					errors.New("protocol must be 'tls' with tls_profile" +
						" in transport block")
			}
		}
		if !rscData.Transport.TCPConnections.IsNull() {
			configSet = append(configSet, setPrefix+"transport tcp-connections "+
				utils.ConvI64toa(rscData.Transport.TCPConnections.ValueInt64()))
		}
		if v := rscData.Transport.TLSProfile.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"transport tls-profile \""+v+"\"")
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityLogStreamData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security log stream " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "category "):
				rscData.Category = append(rscData.Category, types.StringValue(itemTrim))
			case itemTrim == "filter threat-attack":
				rscData.FilterThreatAttack = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "format "):
				rscData.Format = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "rate-limit "):
				rscData.RateLimit, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "severity "):
				rscData.Severity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "file "):
				if rscData.File == nil {
					rscData.File = &securityLogStreamBlockFile{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "name "):
					rscData.File.Name = types.StringValue(strings.Trim(itemTrim, "\""))
				case itemTrim == "allow-duplicates":
					rscData.File.AllowDuplicates = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "rotation "):
					rscData.File.Rotation, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "size "):
					rscData.File.Size, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "host "):
				if rscData.Host == nil {
					rscData.Host = &securityLogStreamBlockHost{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "port "):
					rscData.Host.Port, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "routing-instance "):
					rscData.Host.RoutingInstance = types.StringValue(itemTrim)
				default:
					rscData.Host.IPAddress = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			case balt.CutPrefixInString(&itemTrim, "transport"):
				if rscData.Transport == nil {
					rscData.Transport = &securityLogStreamBlockTransport{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " protocol "):
					rscData.Transport.Protocol = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " tcp-connections "):
					rscData.Transport.TCPConnections, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " tls-profile "):
					rscData.Transport.TLSProfile = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			}
		}
	}

	return nil
}

func (rscData *securityLogStreamData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security log stream " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
