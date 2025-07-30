package provider

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
	_ resource.Resource                   = &bgpNeighbor{}
	_ resource.ResourceWithConfigure      = &bgpNeighbor{}
	_ resource.ResourceWithModifyPlan     = &bgpNeighbor{}
	_ resource.ResourceWithValidateConfig = &bgpNeighbor{}
	_ resource.ResourceWithImportState    = &bgpNeighbor{}
	_ resource.ResourceWithUpgradeState   = &bgpNeighbor{}
)

type bgpNeighbor struct {
	client *junos.Client
}

func newBgpNeighborResource() resource.Resource {
	return &bgpNeighbor{}
}

func (rsc *bgpNeighbor) typeName() string {
	return providerName + "_bgp_neighbor"
}

func (rsc *bgpNeighbor) junosName() string {
	return "bgp neighbor"
}

func (rsc *bgpNeighbor) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *bgpNeighbor) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *bgpNeighbor) Configure(
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

func (rsc *bgpNeighbor) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	attributes := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			Description: "An identifier for the resource with format " +
				"`<ip>" + junos.IDSeparator + "<routing_instance>" + junos.IDSeparator + "<group>`.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"ip": schema.StringAttribute{
			Required:    true,
			Description: "IP of neighbor.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				tfvalidator.StringIPAddress(),
			},
		},
		"routing_instance": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(junos.DefaultW),
			Description: "Routing instance for bgp protocol if not root level.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 63),
				tfvalidator.StringFormat(tfvalidator.DefaultFormat),
			},
		},
		"group": schema.StringAttribute{
			Required:    true,
			Description: "Name of BGP group for this neighbor.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 250),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
	}
	maps.Copy(attributes, bgpAttrData{}.attributesSchema())

	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes:  attributes,
		Blocks:      bgpAttrData{}.blocksSchema(),
	}
}

type bgpNeighborData struct {
	bgpAttrData

	ID              types.String `tfsdk:"id"`
	IP              types.String `tfsdk:"ip"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
	Group           types.String `tfsdk:"group"`
}

type bgpNeighborConfig struct {
	bgpAttrConfig

	ID              types.String `tfsdk:"id"`
	IP              types.String `tfsdk:"ip"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
	Group           types.String `tfsdk:"group"`
}

func (rsc *bgpNeighbor) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config bgpNeighborConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.bgpAttrConfig.validateConfig(ctx, resp)
}

func (rsc *bgpNeighbor) ModifyPlan(
	ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var config, plan bgpNeighborConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.bgpAttrConfig.modifyPlan(ctx, &plan.bgpAttrConfig)

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func (rsc *bgpNeighbor) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan bgpNeighborData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.IP.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("ip"),
			"Empty ip",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "ip"),
		)

		return
	}
	if plan.Group.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("group"),
			"Empty group",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "group"),
		)

		return
	}

	if plan.AdvertiseExternal.IsUnknown() {
		plan.AdvertiseExternal = types.BoolNull()
		if plan.AdvertiseExternalConditional.ValueBool() {
			plan.AdvertiseExternal = types.BoolValue(true)
		}
	}
	if plan.MetricOutIgp.IsUnknown() {
		plan.MetricOutIgp = types.BoolNull()
		if plan.MetricOutIgpDelayMedUpdate.ValueBool() {
			plan.MetricOutIgp = types.BoolValue(true)
		}
		if !plan.MetricOutIgpOffset.IsNull() {
			plan.MetricOutIgp = types.BoolValue(true)
		}
	}
	if plan.MetricOutMinimumIgp.IsUnknown() {
		plan.MetricOutMinimumIgp = types.BoolNull()
		if !plan.MetricOutMinimumIgpOffset.IsNull() {
			plan.MetricOutMinimumIgp = types.BoolValue(true)
		}
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
			bgpGroupExists, err := checkBgpGroupExists(
				fnCtx,
				plan.Group.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !bgpGroupExists {
				resp.Diagnostics.AddAttributeError(
					path.Root("group"),
					tfdiag.PreCheckErrSummary,
					fmt.Sprintf("bgp group %q doesn't exist", plan.Group.ValueString()),
				)

				return false
			}
			bgpNeighborExists, err := checkBgpNeighborExists(
				fnCtx,
				plan.IP.ValueString(),
				plan.RoutingInstance.ValueString(),
				plan.Group.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if bgpNeighborExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf(rsc.junosName()+" %q already exists in group %q in routing-instance %q",
							plan.IP.ValueString(), plan.Group.ValueString(), v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf(rsc.junosName()+" %q already exists in group %q",
							plan.IP.ValueString(), plan.Group.ValueString()),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			bgpNeighborExists, err := checkBgpNeighborExists(
				fnCtx,
				plan.IP.ValueString(),
				plan.RoutingInstance.ValueString(),
				plan.Group.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !bgpNeighborExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						fmt.Sprintf(rsc.junosName()+" %q does not exists in group %q in routing-instance %q after commit "+
							"=> check your config", plan.IP.ValueString(), plan.Group.ValueString(), v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						fmt.Sprintf(rsc.junosName()+" %q does not exists in group %q after commit "+
							"=> check your config", plan.IP.ValueString(), plan.Group.ValueString()),
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

func (rsc *bgpNeighbor) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data bgpNeighborData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom3String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.IP.ValueString(),
			state.RoutingInstance.ValueString(),
			state.Group.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *bgpNeighbor) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state bgpNeighborData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.AdvertiseExternal.IsUnknown() {
		plan.AdvertiseExternal = types.BoolNull()
		if plan.AdvertiseExternalConditional.ValueBool() {
			plan.AdvertiseExternal = types.BoolValue(true)
		}
	}
	if plan.MetricOutIgp.IsUnknown() {
		plan.MetricOutIgp = types.BoolNull()
		if plan.MetricOutIgpDelayMedUpdate.ValueBool() {
			plan.MetricOutIgp = types.BoolValue(true)
		}
		if !plan.MetricOutIgpOffset.IsNull() {
			plan.MetricOutIgp = types.BoolValue(true)
		}
	}
	if plan.MetricOutMinimumIgp.IsUnknown() {
		plan.MetricOutMinimumIgp = types.BoolNull()
		if !plan.MetricOutMinimumIgpOffset.IsNull() {
			plan.MetricOutMinimumIgp = types.BoolValue(true)
		}
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

func (rsc *bgpNeighbor) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state bgpNeighborData
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

func (rsc *bgpNeighbor) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data bgpNeighborData

	var _ resourceDataReadFrom3String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <ip>"+junos.IDSeparator+"<routing_instance>"+junos.IDSeparator+"<group>)",
	)
}

func checkBgpNeighborExists(
	_ context.Context, ip, routingInstance, group string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols bgp group \"" + group + "\" neighbor " + ip + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *bgpNeighborData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(
			rscData.IP.ValueString() + junos.IDSeparator +
				v + junos.IDSeparator +
				rscData.Group.ValueString(),
		)
	} else {
		rscData.ID = types.StringValue(
			rscData.IP.ValueString() + junos.IDSeparator +
				junos.DefaultW + junos.IDSeparator +
				rscData.Group.ValueString(),
		)
	}
}

func (rscData *bgpNeighborData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *bgpNeighborData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "protocols bgp group \"" + rscData.Group.ValueString() + "\" neighbor " + rscData.IP.ValueString() + " "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	dataConfigSet, errPath, err := rscData.bgpAttrData.configSet(setPrefix)
	if err != nil {
		return errPath, err
	}
	configSet = append(configSet, dataConfigSet...)

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *bgpNeighborData) read(
	_ context.Context,
	ip,
	routingInstance,
	group string,
	junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols bgp group \"" + group + "\" neighbor " + ip + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.IP = types.StringValue(ip)
		if routingInstance == "" {
			rscData.RoutingInstance = types.StringValue(junos.DefaultW)
		} else {
			rscData.RoutingInstance = types.StringValue(routingInstance)
		}
		rscData.Group = types.StringValue(group)
		rscData.fillID()
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if err := rscData.bgpAttrData.read(itemTrim, junSess); err != nil {
				return err
			}
		}
	}

	return nil
}

func (rscData *bgpNeighborData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "protocols bgp group \"" + rscData.Group.ValueString() + "\" neighbor " + rscData.IP.ValueString() + " "

	return junSess.ConfigSet(rscData.bgpAttrData.configOptsToDel(delPrefix))
}

func (rscData *bgpNeighborData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}

	configSet := []string{
		delPrefix + "protocols bgp group \"" + rscData.Group.ValueString() + "\" neighbor " + rscData.IP.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
