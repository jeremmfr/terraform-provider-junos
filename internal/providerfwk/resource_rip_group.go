package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &ripGroup{}
	_ resource.ResourceWithConfigure      = &ripGroup{}
	_ resource.ResourceWithValidateConfig = &ripGroup{}
	_ resource.ResourceWithImportState    = &ripGroup{}
	_ resource.ResourceWithUpgradeState   = &ripGroup{}
)

type ripGroup struct {
	client *junos.Client
}

func newRipGroupResource() resource.Resource {
	return &ripGroup{}
}

func (rsc *ripGroup) typeName() string {
	return providerName + "_rip_group"
}

func (rsc *ripGroup) junosName() string {
	return "rip|ripng group"
}

func (rsc *ripGroup) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *ripGroup) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *ripGroup) Configure(
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

func (rsc *ripGroup) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<name>" + junos.IDSeparator + "<routing_instance>` or " +
					"`<name>" + junos.IDSeparator + "ng" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 48),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"ng": schema.BoolAttribute{
				Optional:    true,
				Description: "Protocol `ripng` instead of `rip`.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for RIP group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"demand_circuit": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable demand circuit.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"export": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Export policy.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 63),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
			"import": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Import policy.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 63),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
			"max_retrans_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum time to re-transmit a message in demand-circuit.",
				Validators: []validator.Int64{
					int64validator.Between(5, 180),
				},
			},
			"metric_out": schema.Int64Attribute{
				Optional:    true,
				Description: "Default metric of exported routes.",
				Validators: []validator.Int64{
					int64validator.Between(1, 15),
				},
			},
			"preference": schema.Int64Attribute{
				Optional:    true,
				Description: "Preference of routes learned by this group.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"route_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Delay before routes time out (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(30, 360),
				},
			},
			"update_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Interval between regular route updates (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(10, 60),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"bfd_liveness_detection": ripBlockBfdLivenessDetection{}.resourceSchema(),
		},
	}
}

type ripGroupData struct {
	ID                   types.String                  `tfsdk:"id"`
	Name                 types.String                  `tfsdk:"name"`
	Ng                   types.Bool                    `tfsdk:"ng"`
	RoutingInstance      types.String                  `tfsdk:"routing_instance"`
	DemandCircuit        types.Bool                    `tfsdk:"demand_circuit"`
	Export               []types.String                `tfsdk:"export"`
	Import               []types.String                `tfsdk:"import"`
	MaxRetransTime       types.Int64                   `tfsdk:"max_retrans_time"`
	MetricOut            types.Int64                   `tfsdk:"metric_out"`
	Preference           types.Int64                   `tfsdk:"preference"`
	RouteTimeout         types.Int64                   `tfsdk:"route_timeout"`
	UpdateInterval       types.Int64                   `tfsdk:"update_interval"`
	BfdLivenessDetection *ripBlockBfdLivenessDetection `tfsdk:"bfd_liveness_detection"`
}

func (rscData *ripGroupData) junosName() string {
	if rscData.Ng.ValueBool() {
		return "ripng group"
	}

	return "rip group"
}

type ripGroupConfig struct {
	ID                   types.String                  `tfsdk:"id"`
	Name                 types.String                  `tfsdk:"name"`
	Ng                   types.Bool                    `tfsdk:"ng"`
	RoutingInstance      types.String                  `tfsdk:"routing_instance"`
	DemandCircuit        types.Bool                    `tfsdk:"demand_circuit"`
	Export               types.List                    `tfsdk:"export"`
	Import               types.List                    `tfsdk:"import"`
	MaxRetransTime       types.Int64                   `tfsdk:"max_retrans_time"`
	MetricOut            types.Int64                   `tfsdk:"metric_out"`
	Preference           types.Int64                   `tfsdk:"preference"`
	RouteTimeout         types.Int64                   `tfsdk:"route_timeout"`
	UpdateInterval       types.Int64                   `tfsdk:"update_interval"`
	BfdLivenessDetection *ripBlockBfdLivenessDetection `tfsdk:"bfd_liveness_detection"`
}

func (rsc *ripGroup) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config ripGroupConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Ng.IsNull() && !config.Ng.IsUnknown() {
		if !config.DemandCircuit.IsNull() && !config.DemandCircuit.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("demand_circuit"),
				tfdiag.ConflictConfigErrSummary,
				"ng and demand_circuit cannot be configured together",
			)
		}
		if !config.MaxRetransTime.IsNull() && !config.MaxRetransTime.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("max_retrans_time"),
				tfdiag.ConflictConfigErrSummary,
				"ng and max_retrans_time cannot be configured together",
			)
		}
		if config.BfdLivenessDetection != nil && config.BfdLivenessDetection.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("bfd_liveness_detection"),
				tfdiag.ConflictConfigErrSummary,
				"ng and bfd_liveness_detection cannot be configured together",
			)
		}
	}
	if config.BfdLivenessDetection != nil && config.BfdLivenessDetection.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("bfd_liveness_detection").AtName("*"),
			tfdiag.MissingConfigErrSummary,
			"bfd_liveness_detection block is empty",
		)
	}
}

func (rsc *ripGroup) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan ripGroupData
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
			groupExists, err := checkRipGroupExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.Ng.ValueBool(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if groupExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsInRoutingInstanceMessage(&plan, plan.Name, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsMessage(&plan, plan.Name),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			groupExists, err := checkRipGroupExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.Ng.ValueBool(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !groupExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsInRoutingInstanceAfterCommitMessage(&plan, plan.Name, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsAfterCommitMessage(&plan, plan.Name),
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

func (rsc *ripGroup) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data ripGroupData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String1Bool1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.Ng.ValueBool(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *ripGroup) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state ripGroupData
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

func (rsc *ripGroup) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state ripGroupData
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

func (rsc *ripGroup) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	var data ripGroupData
	idSplit := strings.Split(req.ID, junos.IDSeparator)
	switch {
	case len(idSplit) < 2:
		resp.Diagnostics.AddError(
			"Bad ID Format",
			fmt.Sprintf("missing element(s) in id with separator %q", junos.IDSeparator),
		)

		return
	case len(idSplit) == 2:
		if err := data.read(ctx, idSplit[0], false, idSplit[1], junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
	default:
		if idSplit[1] != "ng" {
			resp.Diagnostics.AddError(
				"Bad ID Format",
				"id must be "+
					"<name>"+junos.IDSeparator+"<routing_instance> or "+
					"<name>"+junos.IDSeparator+"ng"+junos.IDSeparator+"<routing_instance>",
			)

			return
		}
		if err := data.read(ctx, idSplit[0], true, idSplit[2], junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be "+
				"<name>"+junos.IDSeparator+"<routing_instance> or "+
				"<name>"+junos.IDSeparator+"ng"+junos.IDSeparator+"<routing_instance>)",
		)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkRipGroupExists(
	_ context.Context, name string, ng bool, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	if ng {
		showPrefix += "protocols ripng "
	} else {
		showPrefix += "protocols rip "
	}
	showConfig, err := junSess.Command(showPrefix +
		"group \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *ripGroupData) fillID() {
	idPrefix := rscData.Name.ValueString()
	if rscData.Ng.ValueBool() {
		idPrefix += junos.IDSeparator + "ng"
	}
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(idPrefix + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(idPrefix + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *ripGroupData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *ripGroupData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	if rscData.Ng.ValueBool() {
		setPrefix += "protocols ripng "
	} else {
		setPrefix += "protocols rip "
	}
	setPrefix += "group \"" + rscData.Name.ValueString() + "\" "

	configSet := []string{
		setPrefix,
	}

	if rscData.DemandCircuit.ValueBool() {
		configSet = append(configSet, setPrefix+"demand-circuit")
	}
	for _, v := range rscData.Export {
		configSet = append(configSet, setPrefix+"export "+v.ValueString())
	}
	for _, v := range rscData.Import {
		configSet = append(configSet, setPrefix+"import "+v.ValueString())
	}
	if !rscData.MaxRetransTime.IsNull() {
		configSet = append(configSet, setPrefix+"max-retrans-time "+
			utils.ConvI64toa(rscData.MaxRetransTime.ValueInt64()))
	}
	if !rscData.MetricOut.IsNull() {
		configSet = append(configSet, setPrefix+"metric-out "+
			utils.ConvI64toa(rscData.MetricOut.ValueInt64()))
	}
	if !rscData.Preference.IsNull() {
		configSet = append(configSet, setPrefix+"preference "+
			utils.ConvI64toa(rscData.Preference.ValueInt64()))
	}
	if !rscData.RouteTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"route-timeout "+
			utils.ConvI64toa(rscData.RouteTimeout.ValueInt64()))
	}
	if !rscData.UpdateInterval.IsNull() {
		configSet = append(configSet, setPrefix+"update-interval "+
			utils.ConvI64toa(rscData.UpdateInterval.ValueInt64()))
	}
	if rscData.BfdLivenessDetection != nil {
		if rscData.BfdLivenessDetection.isEmpty() {
			return path.Root("bfd_liveness_detection").AtName("*"),
				errors.New("bfd_liveness_detection block is empty")
		}

		configSet = append(configSet, rscData.BfdLivenessDetection.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *ripGroupData) read(
	_ context.Context, name string, ng bool, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	if ng {
		showPrefix += "protocols ripng "
	} else {
		showPrefix += "protocols rip "
	}
	showConfig, err := junSess.Command(showPrefix +
		"group \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		if ng {
			rscData.Ng = types.BoolValue(true)
		}
		rscData.RoutingInstance = types.StringValue(routingInstance)
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
			case itemTrim == "demand-circuit":
				rscData.DemandCircuit = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "export "):
				rscData.Export = append(rscData.Export, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "import "):
				rscData.Import = append(rscData.Import, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "max-retrans-time "):
				rscData.MaxRetransTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "metric-out "):
				rscData.MetricOut, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "preference "):
				rscData.Preference, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "route-timeout "):
				rscData.RouteTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "update-interval "):
				rscData.UpdateInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "bfd-liveness-detection "):
				if rscData.BfdLivenessDetection == nil {
					rscData.BfdLivenessDetection = &ripBlockBfdLivenessDetection{}
				}
				if err := rscData.BfdLivenessDetection.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *ripGroupData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	if rscData.Ng.ValueBool() {
		delPrefix += "protocols ripng "
	} else {
		delPrefix += "protocols rip "
	}
	delPrefix += "group \"" + rscData.Name.ValueString() + "\" "

	configSet := []string{
		delPrefix + "bfd-liveness-detection",
		delPrefix + "demand-circuit",
		delPrefix + "export",
		delPrefix + "import",
		delPrefix + "max-retrans-time",
		delPrefix + "metric-out",
		delPrefix + "preference",
		delPrefix + "route-timeout",
		delPrefix + "update-interval",
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *ripGroupData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	if rscData.Ng.ValueBool() {
		delPrefix += "protocols ripng "
	} else {
		delPrefix += "protocols rip "
	}

	configSet := []string{
		delPrefix + "group \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
