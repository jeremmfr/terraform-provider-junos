package providerfwk

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
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
	_ resource.Resource                   = &vstpVlan{}
	_ resource.ResourceWithConfigure      = &vstpVlan{}
	_ resource.ResourceWithValidateConfig = &vstpVlan{}
	_ resource.ResourceWithImportState    = &vstpVlan{}
)

type vstpVlan struct {
	client *junos.Client
}

func newVstpVlanResource() resource.Resource {
	return &vstpVlan{}
}

func (rsc *vstpVlan) typeName() string {
	return providerName + "_vstp_vlan"
}

func (rsc *vstpVlan) junosName() string {
	return "vstp vlan"
}

func (rsc *vstpVlan) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *vstpVlan) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *vstpVlan) Configure(
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

func (rsc *vstpVlan) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<vlan_id>" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vlan_id": schema.StringAttribute{
				Required:    true,
				Description: "VLAN id or `all`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(409[0-4]|(40[0-8]|[1-3]\d\d|[1-9]\d|[1-9])\d|[1-9]|all)$`),
						"must be a VLAN id (1..4094) or all"),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for vstp protocol if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"backup_bridge_priority": schema.StringAttribute{
				Optional:    true,
				Description: "Priority of the bridge.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d\d?k$`),
						"must be a number with increments of 4k - 4k,8k,..60k",
					),
				},
			},
			"bridge_priority": schema.StringAttribute{
				Optional:    true,
				Description: "Priority of the bridge.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(0|\d\d?k)$`),
						"must be a number with increments of 4k - 0,4k,8k,..60k",
					),
				},
			},
			"forward_delay": schema.Int64Attribute{
				Optional:    true,
				Description: "Time spent in listening or learning state (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(4, 30),
				},
			},
			"hello_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Time interval between configuration BPDUs (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 10),
				},
			},
			"max_age": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum age of received protocol bpdu (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(6, 40),
				},
			},
			"system_identifier": schema.StringAttribute{
				Optional:    true,
				Description: "System identifier to represent this node.",
				Validators: []validator.String{
					tfvalidator.StringMACAddress().WithMac48ColonHexa(),
				},
			},
		},
	}
}

type vstpVlanData struct {
	ID                   types.String `tfsdk:"id"`
	VlanID               types.String `tfsdk:"vlan_id"`
	RoutingInstance      types.String `tfsdk:"routing_instance"`
	BackupBridgePriority types.String `tfsdk:"backup_bridge_priority"`
	BridgePriority       types.String `tfsdk:"bridge_priority"`
	ForwardDelay         types.Int64  `tfsdk:"forward_delay"`
	HelloTime            types.Int64  `tfsdk:"hello_time"`
	MaxAge               types.Int64  `tfsdk:"max_age"`
	SystemIdentifier     types.String `tfsdk:"system_identifier"`
}

func (rsc *vstpVlan) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config vstpVlanData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.BackupBridgePriority.IsNull() && !config.BackupBridgePriority.IsUnknown() {
		if v, err := strconv.Atoi(strings.TrimSuffix(
			config.BackupBridgePriority.ValueString(), "k",
		)); err == nil {
			if v%4 != 0 {
				resp.Diagnostics.AddAttributeError(
					path.Root("backup_bridge_priority"),
					"Bad Value Error",
					"backup_bridge_priority must be a multiple of 4k",
				)
			}
			if v < 4 || v > 60 {
				resp.Diagnostics.AddAttributeError(
					path.Root("backup_bridge_priority"),
					"Bad Value Error",
					"backup_bridge_priority must be between 4k and 60k",
				)
			}
			if !config.BridgePriority.IsNull() && !config.BridgePriority.IsUnknown() {
				if bridgePriority, err := strconv.Atoi(strings.TrimSuffix(
					config.BridgePriority.ValueString(), "k",
				)); err == nil {
					if v <= bridgePriority {
						resp.Diagnostics.AddAttributeError(
							path.Root("backup_bridge_priority"),
							"Bad Value Error",
							"backup_bridge_priority must be worse (higher value) than bridge_priority",
						)
					}
				}
			}
		}
	}
	if !config.BridgePriority.IsNull() && !config.BridgePriority.IsUnknown() {
		if v, err := strconv.Atoi(strings.TrimSuffix(
			config.BridgePriority.ValueString(), "k",
		)); err == nil {
			if v%4 != 0 {
				resp.Diagnostics.AddAttributeError(
					path.Root("bridge_priority"),
					"Bad Value Error",
					"bridge_priority must be a multiple of 4k",
				)
			}
			if v < 0 || v > 60 {
				resp.Diagnostics.AddAttributeError(
					path.Root("bridge_priority"),
					"Bad Value Error",
					"bridge_priority must be between 0 and 60k",
				)
			}
		}
	}
}

func (rsc *vstpVlan) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan vstpVlanData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.VlanID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("vlan_id"),
			"Empty Vlan ID",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "vlan_id"),
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
			vlanExists, err := checkVstpVlanExists(
				fnCtx,
				plan.VlanID.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if vlanExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsInRoutingInstanceMessage(rsc, plan.VlanID, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsMessage(rsc, plan.VlanID),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			vlanExists, err := checkVstpVlanExists(
				fnCtx,
				plan.VlanID.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !vlanExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsInRoutingInstanceAfterCommitMessage(rsc, plan.VlanID, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.VlanID),
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

func (rsc *vstpVlan) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data vstpVlanData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.VlanID.ValueString(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *vstpVlan) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state vstpVlanData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataDelWithOpts = &state
	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *vstpVlan) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state vstpVlanData
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

func (rsc *vstpVlan) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data vstpVlanData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <vlan_id>"+junos.IDSeparator+"<routing_instance>)",
	)
}

func checkVstpVlanExists(
	_ context.Context, vlanID, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols vstp vlan " + vlanID + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *vstpVlanData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.VlanID.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.VlanID.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *vstpVlanData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *vstpVlanData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "protocols vstp vlan " + rscData.VlanID.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if v := rscData.BackupBridgePriority.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"backup-bridge-priority "+v)
	}
	if v := rscData.BridgePriority.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"bridge-priority "+v)
	}
	if !rscData.ForwardDelay.IsNull() {
		configSet = append(configSet, setPrefix+"forward-delay "+
			utils.ConvI64toa(rscData.ForwardDelay.ValueInt64()))
	}
	if !rscData.HelloTime.IsNull() {
		configSet = append(configSet, setPrefix+"hello-time "+
			utils.ConvI64toa(rscData.HelloTime.ValueInt64()))
	}
	if !rscData.MaxAge.IsNull() {
		configSet = append(configSet, setPrefix+"max-age "+
			utils.ConvI64toa(rscData.MaxAge.ValueInt64()))
	}
	if v := rscData.SystemIdentifier.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"system-identifier "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *vstpVlanData) read(
	_ context.Context, vlanID, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols vstp vlan " + vlanID + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.VlanID = types.StringValue(vlanID)
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
			case balt.CutPrefixInString(&itemTrim, "backup-bridge-priority "):
				rscData.BackupBridgePriority = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "bridge-priority "):
				rscData.BridgePriority = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "forward-delay "):
				rscData.ForwardDelay, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "hello-time "):
				rscData.HelloTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "max-age "):
				rscData.MaxAge, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "system-identifier "):
				rscData.SystemIdentifier = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (rscData *vstpVlanData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "protocols vstp vlan " + rscData.VlanID.ValueString() + " "

	configSet := []string{
		delPrefix + "backup-bridge-priority",
		delPrefix + "bridge-priority",
		delPrefix + "forward-delay",
		delPrefix + "hello-time",
		delPrefix + "max-age",
		delPrefix + "system-identifier",
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *vstpVlanData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}

	configSet := []string{
		delPrefix + "protocols vstp vlan " + rscData.VlanID.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
