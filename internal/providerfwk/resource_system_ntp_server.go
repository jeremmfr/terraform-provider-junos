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
	_ resource.Resource                   = &systemNtpServer{}
	_ resource.ResourceWithConfigure      = &systemNtpServer{}
	_ resource.ResourceWithValidateConfig = &systemNtpServer{}
	_ resource.ResourceWithImportState    = &systemNtpServer{}
)

type systemNtpServer struct {
	client *junos.Client
}

func newSystemNtpServerResource() resource.Resource {
	return &systemNtpServer{}
}

func (rsc *systemNtpServer) typeName() string {
	return providerName + "_system_ntp_server"
}

func (rsc *systemNtpServer) junosName() string {
	return "system ntp server"
}

func (rsc *systemNtpServer) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *systemNtpServer) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *systemNtpServer) Configure(
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

func (rsc *systemNtpServer) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<address>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "Address of server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
			"key": schema.Int64Attribute{
				Optional:    true,
				Description: "Authentication key.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65534),
				},
			},
			"prefer": schema.BoolAttribute{
				Optional:    true,
				Description: "Prefer this peer_serv.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Description: "Routing instance through which server is reachable.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"version": schema.Int64Attribute{
				Optional:    true,
				Description: "NTP version to use.",
				Validators: []validator.Int64{
					int64validator.Between(1, 4),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"nts": schema.SingleNestedBlock{
				Description: "Enable NTS protocol for this server.",
				Attributes: map[string]schema.Attribute{
					"remote_identity_distinguished_name_container": schema.StringAttribute{
						Optional:    true,
						Description: "Container string for distinguished name of server to remote identity of server for verification.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"remote_identity_distinguished_name_wildcard": schema.StringAttribute{
						Optional:    true,
						Description: "Wildcard string for distinguished name of server to remote identity of server for verification.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"remote_identity_hostname": schema.StringAttribute{
						Optional:    true,
						Description: "Fully-qualified domain name to remote identity of server for verification.",
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

type systemNtpServerData struct {
	ID              types.String             `tfsdk:"id"`
	Address         types.String             `tfsdk:"address"`
	Key             types.Int64              `tfsdk:"key"`
	Prefer          types.Bool               `tfsdk:"prefer"`
	RoutingInstance types.String             `tfsdk:"routing_instance"`
	Version         types.Int64              `tfsdk:"version"`
	Nts             *systemNtpServerBlockNts `tfsdk:"nts"`
}

type systemNtpServerBlockNts struct {
	RemoteIdentityDistinguishedNameContainer types.String `tfsdk:"remote_identity_distinguished_name_container"`
	RemoteIdentityDistinguishedNameWildcard  types.String `tfsdk:"remote_identity_distinguished_name_wildcard"`
	RemoteIdentityHostname                   types.String `tfsdk:"remote_identity_hostname"`
}

func (rsc *systemNtpServer) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config systemNtpServerData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Nts != nil {
		if !config.Nts.RemoteIdentityDistinguishedNameContainer.IsNull() &&
			!config.Nts.RemoteIdentityDistinguishedNameContainer.IsUnknown() &&
			!config.Nts.RemoteIdentityDistinguishedNameWildcard.IsNull() &&
			!config.Nts.RemoteIdentityDistinguishedNameWildcard.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nts").AtName("remote_identity_distinguished_name_container"),
				tfdiag.ConflictConfigErrSummary,
				"remote_identity_distinguished_name_container and remote_identity_distinguished_name_wildcard"+
					" cannot be configured together",
			)
		}
		if !config.Nts.RemoteIdentityDistinguishedNameContainer.IsNull() &&
			!config.Nts.RemoteIdentityDistinguishedNameContainer.IsUnknown() &&
			!config.Nts.RemoteIdentityHostname.IsNull() &&
			!config.Nts.RemoteIdentityHostname.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nts").AtName("remote_identity_distinguished_name_container"),
				tfdiag.ConflictConfigErrSummary,
				"remote_identity_distinguished_name_container and remote_identity_hostname"+
					" cannot be configured together",
			)
		}
		if !config.Nts.RemoteIdentityDistinguishedNameWildcard.IsNull() &&
			!config.Nts.RemoteIdentityDistinguishedNameWildcard.IsUnknown() &&
			!config.Nts.RemoteIdentityHostname.IsNull() &&
			!config.Nts.RemoteIdentityHostname.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nts").AtName("remote_identity_distinguished_name_wildcard"),
				tfdiag.ConflictConfigErrSummary,
				"remote_identity_distinguished_name_wildcard and remote_identity_hostname"+
					" cannot be configured together",
			)
		}
	}
}

func (rsc *systemNtpServer) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemNtpServerData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Address.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("address"),
			"Empty Address",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "address"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			serverExists, err := checkSystemNtpServerExists(fnCtx, plan.Address.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if serverExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Address),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			serverExists, err := checkSystemNtpServerExists(fnCtx, plan.Address.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !serverExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Address),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *systemNtpServer) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemNtpServerData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Address.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *systemNtpServer) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemNtpServerData
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

func (rsc *systemNtpServer) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state systemNtpServerData
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

func (rsc *systemNtpServer) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data systemNtpServerData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "address"),
	)
}

func checkSystemNtpServerExists(
	_ context.Context, address string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system ntp server " + address + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *systemNtpServerData) fillID() {
	rscData.ID = types.StringValue(rscData.Address.ValueString())
}

func (rscData *systemNtpServerData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemNtpServerData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set system ntp server " + rscData.Address.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if !rscData.Key.IsNull() {
		configSet = append(configSet, setPrefix+"key "+
			utils.ConvI64toa(rscData.Key.ValueInt64()))
	}
	if rscData.Prefer.ValueBool() {
		configSet = append(configSet, setPrefix+"prefer")
	}
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-instance "+v)
	}
	if !rscData.Version.IsNull() {
		configSet = append(configSet, setPrefix+"version "+
			utils.ConvI64toa(rscData.Version.ValueInt64()))
	}

	if rscData.Nts != nil {
		configSet = append(configSet, setPrefix+"nts")

		if v := rscData.Nts.RemoteIdentityDistinguishedNameContainer.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"nts remote-identity distinguished-name container \""+v+"\"")
		}
		if v := rscData.Nts.RemoteIdentityDistinguishedNameWildcard.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"nts remote-identity distinguished-name wildcard \""+v+"\"")
		}
		if v := rscData.Nts.RemoteIdentityHostname.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"nts remote-identity hostname \""+v+"\"")
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *systemNtpServerData) read(
	_ context.Context, address string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system ntp server " + address + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Address = types.StringValue(address)
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
			case balt.CutPrefixInString(&itemTrim, "key "):
				rscData.Key, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "nts"):
				if rscData.Nts == nil {
					rscData.Nts = &systemNtpServerBlockNts{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					switch {
					case balt.CutPrefixInString(&itemTrim, "remote-identity distinguished-name container "):
						rscData.Nts.RemoteIdentityDistinguishedNameContainer = types.StringValue(strings.Trim(itemTrim, "\""))
					case balt.CutPrefixInString(&itemTrim, "remote-identity distinguished-name wildcard "):
						rscData.Nts.RemoteIdentityDistinguishedNameWildcard = types.StringValue(strings.Trim(itemTrim, "\""))
					case balt.CutPrefixInString(&itemTrim, "remote-identity hostname "):
						rscData.Nts.RemoteIdentityHostname = types.StringValue(strings.Trim(itemTrim, "\""))
					}
				}
			case itemTrim == "prefer":
				rscData.Prefer = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "routing-instance "):
				rscData.RoutingInstance = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "version "):
				rscData.Version, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *systemNtpServerData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete system ntp server " + rscData.Address.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
