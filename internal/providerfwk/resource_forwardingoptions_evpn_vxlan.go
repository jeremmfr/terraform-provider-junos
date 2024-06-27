package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &forwardingoptionsEvpnVxlan{}
	_ resource.ResourceWithConfigure      = &forwardingoptionsEvpnVxlan{}
	_ resource.ResourceWithValidateConfig = &forwardingoptionsEvpnVxlan{}
	_ resource.ResourceWithImportState    = &forwardingoptionsEvpnVxlan{}
)

type forwardingoptionsEvpnVxlan struct {
	client *junos.Client
}

func newForwardingoptionsEvpnVxlanResource() resource.Resource {
	return &forwardingoptionsEvpnVxlan{}
}

func (rsc *forwardingoptionsEvpnVxlan) typeName() string {
	return providerName + "_forwardingoptions_evpn_vxlan"
}

func (rsc *forwardingoptionsEvpnVxlan) junosName() string {
	return "forwarding-options evpn-vxlan"
}

func (rsc *forwardingoptionsEvpnVxlan) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *forwardingoptionsEvpnVxlan) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *forwardingoptionsEvpnVxlan) Configure(
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

func (rsc *forwardingoptionsEvpnVxlan) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"shared_tunnels": schema.BoolAttribute{
				Optional:    true,
				Description: "Create VTEP tunnels to EVPN PE.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
	}
}

type forwardingoptionsEvpnVxlanData struct {
	ID              types.String `tfsdk:"id"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
	SharedTunnels   types.Bool   `tfsdk:"shared_tunnels"`
}

func (rscData *forwardingoptionsEvpnVxlanData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData, "ID", "RoutingInstance")
}

func (rsc *forwardingoptionsEvpnVxlan) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config forwardingoptionsEvpnVxlanData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("routing_instance"),
			tfdiag.MissingConfigErrSummary,
			"at least one of arguments need to be set (in addition to `routing_instance`)",
		)
	}
}

func (rsc *forwardingoptionsEvpnVxlan) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan forwardingoptionsEvpnVxlanData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
				instanceExists, err := checkRoutingInstanceExists(fnCtx, v, junSess)
				if err != nil {
					resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

					return false
				}
				if !instanceExists {
					resp.Diagnostics.AddAttributeError(
						path.Root("routing_instance"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("routing instance %q doesn't exist", v),
					)

					return false
				}
			}

			return true
		},
		nil,
		&plan,
		resp,
	)
}

func (rsc *forwardingoptionsEvpnVxlan) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data forwardingoptionsEvpnVxlanData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	defer junos.MutexUnlock()

	if v := state.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, v, junSess)
		if err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !instanceExists {
			resp.State.RemoveResource(ctx)

			return
		}
	}
	if err := data.read(ctx, state.RoutingInstance.ValueString(), junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if data.ID.IsNull() {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *forwardingoptionsEvpnVxlan) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state forwardingoptionsEvpnVxlanData
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

func (rsc *forwardingoptionsEvpnVxlan) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state forwardingoptionsEvpnVxlanData
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

func (rsc *forwardingoptionsEvpnVxlan) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	var data forwardingoptionsEvpnVxlanData
	if req.ID != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, req.ID, junSess)
		if err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !instanceExists {
			resp.Diagnostics.AddError(
				tfdiag.NotFoundErrSummary,
				fmt.Sprintf("routing instance %q doesn't exist", req.ID),
			)

			return
		}
	}
	if err := data.read(ctx, req.ID, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rscData *forwardingoptionsEvpnVxlanData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(v)
	} else {
		rscData.ID = types.StringValue(junos.DefaultW)
	}
}

func (rscData *forwardingoptionsEvpnVxlanData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "forwarding-options evpn-vxlan "

	if rscData.SharedTunnels.ValueBool() {
		configSet = append(configSet, setPrefix+"shared-tunnels")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *forwardingoptionsEvpnVxlanData) read(
	_ context.Context, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"forwarding-options evpn-vxlan" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if routingInstance == "" {
		rscData.RoutingInstance = types.StringValue(junos.DefaultW)
	} else {
		rscData.RoutingInstance = types.StringValue(routingInstance)
	}
	rscData.fillID()
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if itemTrim == "shared-tunnels" {
				rscData.SharedTunnels = types.BoolValue(true)
			}
		}
	}

	return nil
}

func (rscData *forwardingoptionsEvpnVxlanData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "forwarding-options evpn-vxlan "

	configSet := []string{
		delPrefix + "shared-tunnels",
	}

	return junSess.ConfigSet(configSet)
}
