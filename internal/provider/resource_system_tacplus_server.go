package provider

import (
	"context"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &systemTacplusServer{}
	_ resource.ResourceWithConfigure      = &systemTacplusServer{}
	_ resource.ResourceWithValidateConfig = &systemTacplusServer{}
	_ resource.ResourceWithImportState    = &systemTacplusServer{}
)

type systemTacplusServer struct {
	client *junos.Client
}

func newSystemTacplusServerResource() resource.Resource {
	return &systemTacplusServer{}
}

func (rsc *systemTacplusServer) typeName() string {
	return providerName + "_system_tacplus_server"
}

func (rsc *systemTacplusServer) junosName() string {
	return "system tacplus-server"
}

func (rsc *systemTacplusServer) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *systemTacplusServer) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *systemTacplusServer) Configure(
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

func (rsc *systemTacplusServer) Schema(
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
				Description: "TACACS+ authentication server address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Description: "TACACS+ authentication server port number.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
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
			"secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Shared secret with the authentication server.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"secret_wo": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
				Description: "Shared secret with the authentication server, not stored in state.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
					stringvalidator.AlsoRequires(path.MatchRoot("secret_wo_version")),
				},
			},
			"secret_wo_version": schema.Int64Attribute{
				Optional:    true,
				Description: "Version of `secret_wo` to trigger the sending of its value.",
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("secret_wo")),
				},
			},
			"single_connection": schema.BoolAttribute{
				Optional:    true,
				Description: "Optimize TCP connection attempts.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"source_address": schema.StringAttribute{
				Optional:    true,
				Description: "Use specified address as source address.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Request timeout period.",
				Validators: []validator.Int64{
					int64validator.Between(1, 90),
				},
			},
		},
	}
}

type systemTacplusServerData struct {
	ID               types.String `tfsdk:"id"`
	Address          types.String `tfsdk:"address"`
	Port             types.Int64  `tfsdk:"port"`
	RoutingInstance  types.String `tfsdk:"routing_instance"`
	Secret           types.String `tfsdk:"secret"`
	SecretWO         types.String `tfsdk:"secret_wo"`
	SecretWOVersion  types.Int64  `tfsdk:"secret_wo_version"`
	SingleConnection types.Bool   `tfsdk:"single_connection"`
	SourceAddress    types.String `tfsdk:"source_address"`
	Timeout          types.Int64  `tfsdk:"timeout"`
}

func (rsc *systemTacplusServer) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config systemTacplusServerData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Secret.IsNull() && !config.Secret.IsUnknown() &&
		!config.SecretWO.IsNull() && !config.SecretWO.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret"),
			tfdiag.ConflictConfigErrSummary,
			"secret and secret_wo cannot be configured together",
		)
	}
}

func (rsc *systemTacplusServer) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemTacplusServerData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(plan.getWriteOnly(ctx, req.Config)...)
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
			serverExists, err := checkSystemTacplusServerExists(fnCtx, plan.Address.ValueString(), junSess)
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
			serverExists, err := checkSystemTacplusServerExists(fnCtx, plan.Address.ValueString(), junSess)
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

func (rsc *systemTacplusServer) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemTacplusServerData
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
		func() {
			data.keepWriteOnly(&state)
		},
		resp,
	)
}

func (rsc *systemTacplusServer) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemTacplusServerData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(plan.getWriteOnly(ctx, req.Config)...)
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

func (rsc *systemTacplusServer) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state systemTacplusServerData
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

func (rsc *systemTacplusServer) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data systemTacplusServerData

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

func checkSystemTacplusServerExists(
	ctx context.Context, address string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"system tacplus-server "+address+junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

// getWriteOnly read the write-only arguments from the configuration,
// their values aren't present in the plan or the state.
func (rscData *systemTacplusServerData) getWriteOnly(
	ctx context.Context, config tfsdk.Config,
) diag.Diagnostics {
	return config.GetAttribute(ctx, path.Root("secret_wo"), &rscData.SecretWO)
}

// keepWriteOnly carry over the version arguments of the write-only arguments from the state,
// and don't read the secrets in the standard arguments when the write-only ones are used.
func (rscData *systemTacplusServerData) keepWriteOnly(state *systemTacplusServerData) {
	rscData.SecretWOVersion = state.SecretWOVersion
	if !state.SecretWOVersion.IsNull() {
		rscData.Secret = types.StringNull()
	}
}

func (rscData *systemTacplusServerData) fillID() {
	rscData.ID = types.StringValue(rscData.Address.ValueString())
}

func (rscData *systemTacplusServerData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemTacplusServerData) set(
	ctx context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set system tacplus-server " + rscData.Address.ValueString() + " "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if !rscData.Port.IsNull() {
		configSet = append(configSet, setPrefix+"port "+
			utils.ConvI64toa(rscData.Port.ValueInt64()))
	}
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-instance "+v)
	}
	if v := rscData.Secret.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"secret \""+v+"\"")
	} else if v := rscData.SecretWO.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"secret \""+v+"\"")
	}
	if rscData.SingleConnection.ValueBool() {
		configSet = append(configSet, setPrefix+"single-connection")
	}
	if v := rscData.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}
	if !rscData.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(rscData.Timeout.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(ctx, configSet)
}

func (rscData *systemTacplusServerData) read(
	ctx context.Context, address string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"system tacplus-server "+address+junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Address = types.StringValue(address)
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
			case balt.CutPrefixInString(&itemTrim, "port "):
				rscData.Port, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "routing-instance "):
				rscData.RoutingInstance = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "secret "):
				rscData.Secret, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "secret")
				if err != nil {
					return err
				}
			case itemTrim == "single-connection":
				rscData.SingleConnection = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "source-address "):
				rscData.SourceAddress = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "timeout "):
				rscData.Timeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *systemTacplusServerData) del(
	ctx context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete system tacplus-server " + rscData.Address.ValueString(),
	}

	return junSess.ConfigSet(ctx, configSet)
}
