package providerfwk

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
	_ resource.Resource                = &systemRadiusServer{}
	_ resource.ResourceWithConfigure   = &systemRadiusServer{}
	_ resource.ResourceWithImportState = &systemRadiusServer{}
)

type systemRadiusServer struct {
	client *junos.Client
}

func newSystemRadiusServerResource() resource.Resource {
	return &systemRadiusServer{}
}

func (rsc *systemRadiusServer) typeName() string {
	return providerName + "_system_radius_server"
}

func (rsc *systemRadiusServer) junosName() string {
	return "system radius-server"
}

func (rsc *systemRadiusServer) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *systemRadiusServer) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *systemRadiusServer) Configure(
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

func (rsc *systemRadiusServer) Schema(
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
				Description: "RADIUS server address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
			"secret": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Shared secret with the RADIUS server.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"accounting_port": schema.Int64Attribute{
				Optional:    true,
				Description: "RADIUS server accounting port number.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"accounting_retry": schema.Int64Attribute{
				Optional:    true,
				Description: "Accounting retry attempts.",
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
			},
			"accounting_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Accounting request timeout period.",
				Validators: []validator.Int64{
					int64validator.Between(0, 1000),
				},
			},
			"dynamic_request_port": schema.Int64Attribute{
				Optional:    true,
				Description: "RADIUS client dynamic request port number.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"max_outstanding_requests": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum requests in flight to server.",
				Validators: []validator.Int64{
					int64validator.Between(0, 2000),
				},
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Description: "RADIUS server authentication port number.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"preauthentication_port": schema.Int64Attribute{
				Optional:    true,
				Description: "RADIUS server preauthentication port number.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"preauthentication_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Shared secret with the RADIUS server.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"retry": schema.Int64Attribute{
				Optional:    true,
				Description: "Retry attempts.",
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
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
					int64validator.Between(1, 1000),
				},
			},
		},
	}
}

type systemRadiusServerData struct {
	ID                      types.String `tfsdk:"id"`
	Address                 types.String `tfsdk:"address"`
	Secret                  types.String `tfsdk:"secret"`
	AccountingPort          types.Int64  `tfsdk:"accounting_port"`
	AccountingRetry         types.Int64  `tfsdk:"accounting_retry"`
	AccountingTimeout       types.Int64  `tfsdk:"accounting_timeout"`
	DynamicRequestPort      types.Int64  `tfsdk:"dynamic_request_port"`
	MaxOutstandingRequests  types.Int64  `tfsdk:"max_outstanding_requests"`
	Port                    types.Int64  `tfsdk:"port"`
	PreauthenticationPort   types.Int64  `tfsdk:"preauthentication_port"`
	PreauthenticationSecret types.String `tfsdk:"preauthentication_secret"`
	Retry                   types.Int64  `tfsdk:"retry"`
	RoutingInstance         types.String `tfsdk:"routing_instance"`
	SourceAddress           types.String `tfsdk:"source_address"`
	Timeout                 types.Int64  `tfsdk:"timeout"`
}

func (rsc *systemRadiusServer) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemRadiusServerData
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
			serverExists, err := checkSystemRadiusServerExists(fnCtx, plan.Address.ValueString(), junSess)
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
			serverExists, err := checkSystemRadiusServerExists(fnCtx, plan.Address.ValueString(), junSess)
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

func (rsc *systemRadiusServer) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemRadiusServerData
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

func (rsc *systemRadiusServer) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemRadiusServerData
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

func (rsc *systemRadiusServer) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state systemRadiusServerData
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

func (rsc *systemRadiusServer) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data systemRadiusServerData

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

func checkSystemRadiusServerExists(
	_ context.Context, address string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system radius-server " + address + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *systemRadiusServerData) fillID() {
	rscData.ID = types.StringValue(rscData.Address.ValueString())
}

func (rscData *systemRadiusServerData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemRadiusServerData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set system radius-server " + rscData.Address.ValueString() + " "

	configSet := []string{
		setPrefix + "secret \"" + rscData.Secret.ValueString() + "\"",
	}

	if !rscData.AccountingPort.IsNull() {
		configSet = append(configSet, setPrefix+"accounting-port "+
			utils.ConvI64toa(rscData.AccountingPort.ValueInt64()))
	}
	if !rscData.AccountingRetry.IsNull() {
		configSet = append(configSet, setPrefix+"accounting-retry "+
			utils.ConvI64toa(rscData.AccountingRetry.ValueInt64()))
	}
	if !rscData.AccountingTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"accounting-timeout "+
			utils.ConvI64toa(rscData.AccountingTimeout.ValueInt64()))
	}
	if !rscData.DynamicRequestPort.IsNull() {
		configSet = append(configSet, setPrefix+"dynamic-request-port "+
			utils.ConvI64toa(rscData.DynamicRequestPort.ValueInt64()))
	}
	if !rscData.MaxOutstandingRequests.IsNull() {
		configSet = append(configSet, setPrefix+"max-outstanding-requests "+
			utils.ConvI64toa(rscData.MaxOutstandingRequests.ValueInt64()))
	}
	if !rscData.Port.IsNull() {
		configSet = append(configSet, setPrefix+"port "+
			utils.ConvI64toa(rscData.Port.ValueInt64()))
	}
	if !rscData.PreauthenticationPort.IsNull() {
		configSet = append(configSet, setPrefix+"preauthentication-port "+
			utils.ConvI64toa(rscData.PreauthenticationPort.ValueInt64()))
	}
	if v := rscData.PreauthenticationSecret.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"preauthentication-secret \""+v+"\"")
	}
	if !rscData.Retry.IsNull() {
		configSet = append(configSet, setPrefix+"retry "+
			utils.ConvI64toa(rscData.Retry.ValueInt64()))
	}
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-instance "+v)
	}
	if v := rscData.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}
	if !rscData.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(rscData.Timeout.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *systemRadiusServerData) read(
	_ context.Context, address string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system radius-server " + address + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "secret "):
				rscData.Secret, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "secret")
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "accounting-port "):
				rscData.AccountingPort, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "accounting-retry "):
				rscData.AccountingRetry, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "accounting-timeout "):
				rscData.AccountingTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "dynamic-request-port "):
				rscData.DynamicRequestPort, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "max-outstanding-requests "):
				rscData.MaxOutstandingRequests, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "port "):
				rscData.Port, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "preauthentication-port "):
				rscData.PreauthenticationPort, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "preauthentication-secret "):
				rscData.PreauthenticationSecret, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "preauthentication-secret")
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "retry "):
				rscData.Retry, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "routing-instance "):
				rscData.RoutingInstance = types.StringValue(itemTrim)
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

func (rscData *systemRadiusServerData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete system radius-server " + rscData.Address.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
