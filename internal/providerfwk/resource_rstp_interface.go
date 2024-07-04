package providerfwk

import (
	"context"
	"fmt"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &rstpInterface{}
	_ resource.ResourceWithConfigure      = &rstpInterface{}
	_ resource.ResourceWithValidateConfig = &rstpInterface{}
	_ resource.ResourceWithImportState    = &rstpInterface{}
)

type rstpInterface struct {
	client *junos.Client
}

func newRstpInterfaceResource() resource.Resource {
	return &rstpInterface{}
}

func (rsc *rstpInterface) typeName() string {
	return providerName + "_rstp_interface"
}

func (rsc *rstpInterface) junosName() string {
	return "rstp interface"
}

func (rsc *rstpInterface) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *rstpInterface) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *rstpInterface) Configure(
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

func (rsc *rstpInterface) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<name>" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Interface name or `all`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					tfvalidator.StringDotExclusion(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for rstp protocol if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"access_trunk": schema.BoolAttribute{
				Optional:    true,
				Description: "Send/Receive untagged RSTP BPDUs on this interface.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"bpdu_timeout_action_alarm": schema.BoolAttribute{
				Optional:    true,
				Description: "Generate an alarm on BPDU expiry (Loop Protect).",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"bpdu_timeout_action_block": schema.BoolAttribute{
				Optional:    true,
				Description: "Block the interface on BPDU expiry (Loop Protect).",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"cost": schema.Int64Attribute{
				Optional:    true,
				Description: "Cost of the interface.",
				Validators: []validator.Int64{
					int64validator.Between(1, 200000000),
				},
			},
			"edge": schema.BoolAttribute{
				Optional:    true,
				Description: "Port is an edge port.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"mode": schema.StringAttribute{
				Optional:    true,
				Description: "Interface mode (P2P or shared).",
				Validators: []validator.String{
					stringvalidator.OneOf("point-to-point", "shared"),
				},
			},
			"no_root_port": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not allow the interface to become root (Root Protect).",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"priority": schema.Int64Attribute{
				Optional:    true,
				Description: "Interface priority (in increments of 16).",
				Validators: []validator.Int64{
					int64validator.Between(0, 240),
				},
			},
		},
	}
}

type rstpInterfaceData struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	RoutingInstance        types.String `tfsdk:"routing_instance"`
	AccessTrunk            types.Bool   `tfsdk:"access_trunk"`
	BpduTimeoutActionAlarm types.Bool   `tfsdk:"bpdu_timeout_action_alarm"`
	BpduTimeoutActionBlock types.Bool   `tfsdk:"bpdu_timeout_action_block"`
	Cost                   types.Int64  `tfsdk:"cost"`
	Edge                   types.Bool   `tfsdk:"edge"`
	Mode                   types.String `tfsdk:"mode"`
	NoRootPort             types.Bool   `tfsdk:"no_root_port"`
	Priority               types.Int64  `tfsdk:"priority"`
}

func (rsc *rstpInterface) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config rstpInterfaceData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Priority.IsNull() && !config.Priority.IsUnknown() {
		if config.Priority.ValueInt64()%16 != 0 {
			resp.Diagnostics.AddAttributeError(
				path.Root("priority"),
				"Bad Value Error",
				"priority must be a multiple of 16",
			)
		}
	}
}

func (rsc *rstpInterface) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan rstpInterfaceData
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
			interfaceExists, err := checkRstpInterfaceExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if interfaceExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsInRoutingInstanceMessage(rsc, plan.Name, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsMessage(rsc, plan.Name),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			interfaceExists, err := checkRstpInterfaceExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !interfaceExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsInRoutingInstanceAfterCommitMessage(rsc, plan.Name, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Name),
					)
				}

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *rstpInterface) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data rstpInterfaceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *rstpInterface) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state rstpInterfaceData
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

func (rsc *rstpInterface) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state rstpInterfaceData
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

func (rsc *rstpInterface) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data rstpInterfaceData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <name>"+junos.IDSeparator+"<routing_instance>)",
	)
}

func checkRstpInterfaceExists(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols rstp interface " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *rstpInterfaceData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *rstpInterfaceData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *rstpInterfaceData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "protocols rstp interface " + rscData.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if rscData.AccessTrunk.ValueBool() {
		configSet = append(configSet, setPrefix+"access-trunk")
	}
	if rscData.BpduTimeoutActionAlarm.ValueBool() {
		configSet = append(configSet, setPrefix+"bpdu-timeout-action alarm")
	}
	if rscData.BpduTimeoutActionBlock.ValueBool() {
		configSet = append(configSet, setPrefix+"bpdu-timeout-action block")
	}
	if !rscData.Cost.IsNull() {
		configSet = append(configSet, setPrefix+"cost "+
			utils.ConvI64toa(rscData.Cost.ValueInt64()))
	}
	if rscData.Edge.ValueBool() {
		configSet = append(configSet, setPrefix+"edge")
	}
	if v := rscData.Mode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"mode "+v)
	}
	if rscData.NoRootPort.ValueBool() {
		configSet = append(configSet, setPrefix+"no-root-port")
	}
	if !rscData.Priority.IsNull() {
		configSet = append(configSet, setPrefix+"priority "+
			utils.ConvI64toa(rscData.Priority.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *rstpInterfaceData) read(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols rstp interface " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		if routingInstance == "" {
			rscData.RoutingInstance = types.StringValue(junos.DefaultW)
		} else {
			rscData.RoutingInstance = types.StringValue(routingInstance)
		}
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
			case itemTrim == "access-trunk":
				rscData.AccessTrunk = types.BoolValue(true)
			case itemTrim == "bpdu-timeout-action alarm":
				rscData.BpduTimeoutActionAlarm = types.BoolValue(true)
			case itemTrim == "bpdu-timeout-action block":
				rscData.BpduTimeoutActionBlock = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "cost "):
				rscData.Cost, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "edge":
				rscData.Edge = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "mode "):
				rscData.Mode = types.StringValue(itemTrim)
			case itemTrim == "no-root-port":
				rscData.NoRootPort = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "priority "):
				rscData.Priority, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *rstpInterfaceData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}

	configSet := []string{
		delPrefix + "protocols rstp interface " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
