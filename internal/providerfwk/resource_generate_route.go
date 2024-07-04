package providerfwk

import (
	"context"
	"fmt"
	"regexp"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &generateRoute{}
	_ resource.ResourceWithConfigure      = &generateRoute{}
	_ resource.ResourceWithValidateConfig = &generateRoute{}
	_ resource.ResourceWithImportState    = &generateRoute{}
)

type generateRoute struct {
	client *junos.Client
}

func newGenerateRouteResource() resource.Resource {
	return &generateRoute{}
}

func (rsc *generateRoute) typeName() string {
	return providerName + "_generate_route"
}

func (rsc *generateRoute) junosName() string {
	return "generate route"
}

func (rsc *generateRoute) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *generateRoute) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *generateRoute) Configure(
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

func (rsc *generateRoute) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<destination>" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"destination": schema.StringAttribute{
				Required:    true,
				Description: "Destination prefix.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					tfvalidator.StringCIDRNetwork(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for generate route.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Description: "Remove inactive route from forwarding table.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"as_path_aggregator_address": schema.StringAttribute{
				Optional:    true,
				Description: "Address of BGP system to add AGGREGATOR path attribute to route.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
			"as_path_aggregator_as_number": schema.StringAttribute{
				Optional:    true,
				Description: "AS number to add AGGREGATOR path attribute to route.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d+(\.\d+)?$`),
						"must be in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format"),
				},
			},
			"as_path_atomic_aggregate": schema.BoolAttribute{
				Optional:    true,
				Description: "Add ATOMIC_AGGREGATE path attribute to route.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"as_path_origin": schema.StringAttribute{
				Optional:    true,
				Description: "Define origin.",
				Validators: []validator.String{
					stringvalidator.OneOf("egp", "igp", "incomplete"),
				},
			},
			"as_path_path": schema.StringAttribute{
				Optional:    true,
				Description: "Path to as-path.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"brief": schema.BoolAttribute{
				Optional:    true,
				Description: "Include longest common sequences from contributing paths.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"community": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "BGP community.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"discard": schema.BoolAttribute{
				Optional:    true,
				Description: "Drop packets to destination; send no ICMP unreachables.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"full": schema.BoolAttribute{
				Optional:    true,
				Description: "Include all AS numbers from all contributing paths.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"metric": schema.Int64Attribute{
				Optional:    true,
				Description: "Metric for generate route.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"next_table": schema.StringAttribute{
				Optional:    true,
				Description: "Next hop to another table.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"passive": schema.BoolAttribute{
				Optional:    true,
				Description: "Retain inactive route in forwarding table.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"policy": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Policy filter.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"preference": schema.Int64Attribute{
				Optional:    true,
				Description: "Preference for aggregate route.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
		},
	}
}

type generateRouteData struct {
	ID                       types.String   `tfsdk:"id"`
	Destination              types.String   `tfsdk:"destination"`
	RoutingInstance          types.String   `tfsdk:"routing_instance"`
	Active                   types.Bool     `tfsdk:"active"`
	ASPathAggregatorAddress  types.String   `tfsdk:"as_path_aggregator_address"`
	ASPathAggregatorASNumber types.String   `tfsdk:"as_path_aggregator_as_number"`
	ASPathAtomicAggregate    types.Bool     `tfsdk:"as_path_atomic_aggregate"`
	ASPathOrigin             types.String   `tfsdk:"as_path_origin"`
	ASPathPath               types.String   `tfsdk:"as_path_path"`
	Brief                    types.Bool     `tfsdk:"brief"`
	Community                []types.String `tfsdk:"community"`
	Discard                  types.Bool     `tfsdk:"discard"`
	Full                     types.Bool     `tfsdk:"full"`
	Metric                   types.Int64    `tfsdk:"metric"`
	NextTable                types.String   `tfsdk:"next_table"`
	Passive                  types.Bool     `tfsdk:"passive"`
	Policy                   []types.String `tfsdk:"policy"`
	Preference               types.Int64    `tfsdk:"preference"`
}

type generateRouteConfig struct {
	ID                       types.String `tfsdk:"id"`
	Destination              types.String `tfsdk:"destination"`
	RoutingInstance          types.String `tfsdk:"routing_instance"`
	Active                   types.Bool   `tfsdk:"active"`
	ASPathAggregatorAddress  types.String `tfsdk:"as_path_aggregator_address"`
	ASPathAggregatorASNumber types.String `tfsdk:"as_path_aggregator_as_number"`
	ASPathAtomicAggregate    types.Bool   `tfsdk:"as_path_atomic_aggregate"`
	ASPathOrigin             types.String `tfsdk:"as_path_origin"`
	ASPathPath               types.String `tfsdk:"as_path_path"`
	Brief                    types.Bool   `tfsdk:"brief"`
	Community                types.List   `tfsdk:"community"`
	Discard                  types.Bool   `tfsdk:"discard"`
	Full                     types.Bool   `tfsdk:"full"`
	Metric                   types.Int64  `tfsdk:"metric"`
	NextTable                types.String `tfsdk:"next_table"`
	Passive                  types.Bool   `tfsdk:"passive"`
	Policy                   types.List   `tfsdk:"policy"`
	Preference               types.Int64  `tfsdk:"preference"`
}

func (rsc *generateRoute) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config generateRouteConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Active.IsNull() && !config.Active.IsUnknown() &&
		!config.Passive.IsNull() && !config.Passive.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("active"),
			tfdiag.ConflictConfigErrSummary,
			"active and passive cannot be configured together",
		)
	}
	if !config.ASPathAggregatorASNumber.IsNull() &&
		config.ASPathAggregatorAddress.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("as_path_aggregator_as_number"),
			tfdiag.MissingConfigErrSummary,
			"as_path_aggregator_address must be specified with as_path_aggregator_as_number",
		)
	}
	if !config.ASPathAggregatorAddress.IsNull() &&
		config.ASPathAggregatorASNumber.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("as_path_aggregator_address"),
			tfdiag.MissingConfigErrSummary,
			"as_path_aggregator_as_number must be specified with as_path_aggregator_address",
		)
	}
	if !config.Brief.IsNull() && !config.Brief.IsUnknown() &&
		!config.Full.IsNull() && !config.Full.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("brief"),
			tfdiag.ConflictConfigErrSummary,
			"brief and full cannot be configured together",
		)
	}
	if !config.Discard.IsNull() && !config.Discard.IsUnknown() &&
		!config.NextTable.IsNull() && !config.NextTable.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("discard"),
			tfdiag.ConflictConfigErrSummary,
			"discard and next_table cannot be configured together",
		)
	}
}

func (rsc *generateRoute) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan generateRouteData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Destination.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("destination"),
			"Empty Destination",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "destination"),
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
			routeExists, err := checkGenerateRouteExists(
				fnCtx,
				plan.Destination.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if routeExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsInRoutingInstanceMessage(rsc, plan.Destination, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsMessage(rsc, plan.Destination),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			routeExists, err := checkGenerateRouteExists(
				fnCtx,
				plan.Destination.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !routeExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsInRoutingInstanceAfterCommitMessage(rsc, plan.Destination, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Destination),
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

func (rsc *generateRoute) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data generateRouteData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Destination.ValueString(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *generateRoute) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state generateRouteData
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

func (rsc *generateRoute) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state generateRouteData
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

func (rsc *generateRoute) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data generateRouteData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <destination>"+junos.IDSeparator+"<routing_instance>)",
	)
}

func checkGenerateRouteExists(
	_ context.Context, destination, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	switch routingInstance {
	case junos.DefaultW, "":
		showPrefix += junos.RoutingOptionsWS
		if strings.Contains(destination, ":") {
			showPrefix += junos.RibInet60WS
		}
	default:
		showPrefix += junos.RoutingInstancesWS + routingInstance + " " + junos.RoutingOptionsWS
		if strings.Contains(destination, ":") {
			showPrefix += "rib " + routingInstance + ".inet6.0 "
		}
	}
	showConfig, err := junSess.Command(showPrefix +
		"generate route " + destination + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *generateRouteData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Destination.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Destination.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *generateRouteData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *generateRouteData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := junos.SetLS
	switch routingInstance := rscData.RoutingInstance.ValueString(); routingInstance {
	case junos.DefaultW, "":
		setPrefix += junos.RoutingOptionsWS
		if strings.Contains(rscData.Destination.ValueString(), ":") {
			setPrefix += junos.RibInet60WS
		}
	default:
		setPrefix += junos.RoutingInstancesWS + routingInstance + " " + junos.RoutingOptionsWS
		if strings.Contains(rscData.Destination.ValueString(), ":") {
			setPrefix += "rib " + routingInstance + ".inet6.0 "
		}
	}
	setPrefix += "generate route " + rscData.Destination.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if rscData.Active.ValueBool() {
		configSet = append(configSet, setPrefix+"active")
	}
	if vNumber, vAddress := rscData.ASPathAggregatorASNumber.ValueString(),
		rscData.ASPathAggregatorAddress.ValueString(); vNumber != "" && vAddress != "" {
		configSet = append(configSet, setPrefix+"as-path aggregator "+vNumber+" "+vAddress)
	}
	if rscData.ASPathAtomicAggregate.ValueBool() {
		configSet = append(configSet, setPrefix+"as-path atomic-aggregate")
	}
	if v := rscData.ASPathOrigin.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"as-path origin "+v)
	}
	if v := rscData.ASPathPath.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"as-path path \""+v+"\"")
	}
	if rscData.Brief.ValueBool() {
		configSet = append(configSet, setPrefix+"brief")
	}
	for _, v := range rscData.Community {
		configSet = append(configSet, setPrefix+"community \""+v.ValueString()+"\"")
	}
	if rscData.Discard.ValueBool() {
		configSet = append(configSet, setPrefix+"discard")
	}
	if rscData.Full.ValueBool() {
		configSet = append(configSet, setPrefix+"full")
	}
	if !rscData.Metric.IsNull() {
		configSet = append(configSet, setPrefix+"metric "+
			utils.ConvI64toa(rscData.Metric.ValueInt64()))
	}
	if v := rscData.NextTable.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"next-table \""+v+"\"")
	}
	if rscData.Passive.ValueBool() {
		configSet = append(configSet, setPrefix+"passive")
	}
	for _, v := range rscData.Policy {
		configSet = append(configSet, setPrefix+"policy \""+v.ValueString()+"\"")
	}
	if !rscData.Preference.IsNull() {
		configSet = append(configSet, setPrefix+"preference "+
			utils.ConvI64toa(rscData.Preference.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *generateRouteData) read(
	_ context.Context, destination, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	switch routingInstance {
	case junos.DefaultW, "":
		showPrefix += junos.RoutingOptionsWS
		if strings.Contains(destination, ":") {
			showPrefix += junos.RibInet60WS
		}
	default:
		showPrefix += junos.RoutingInstancesWS + routingInstance + " " + junos.RoutingOptionsWS
		if strings.Contains(destination, ":") {
			showPrefix += "rib " + routingInstance + ".inet6.0 "
		}
	}
	showConfig, err := junSess.Command(showPrefix +
		"generate route " + destination + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Destination = types.StringValue(destination)
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
			case itemTrim == "active":
				rscData.Active = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "as-path aggregator "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 2 { // <as_number> <address>
					return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "as-path aggregator", itemTrim)
				}
				rscData.ASPathAggregatorASNumber = types.StringValue(itemTrimFields[0])
				rscData.ASPathAggregatorAddress = types.StringValue(itemTrimFields[1])
			case itemTrim == "as-path atomic-aggregate":
				rscData.ASPathAtomicAggregate = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "as-path origin "):
				rscData.ASPathOrigin = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "as-path path "):
				rscData.ASPathPath = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "brief":
				rscData.Brief = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "community "):
				rscData.Community = append(rscData.Community, types.StringValue(strings.Trim(itemTrim, "\"")))
			case itemTrim == junos.DiscardW:
				rscData.Discard = types.BoolValue(true)
			case itemTrim == "full":
				rscData.Full = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "metric "):
				rscData.Metric, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "next-table "):
				rscData.NextTable = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "passive":
				rscData.Passive = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "policy "):
				rscData.Policy = append(rscData.Policy, types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "preference "):
				rscData.Preference, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *generateRouteData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	switch routingInstance := rscData.RoutingInstance.ValueString(); routingInstance {
	case junos.DefaultW, "":
		delPrefix += junos.RoutingOptionsWS
		if strings.Contains(rscData.Destination.ValueString(), ":") {
			delPrefix += junos.RibInet60WS
		}
	default:
		delPrefix += junos.RoutingInstancesWS + routingInstance + " " + junos.RoutingOptionsWS
		if strings.Contains(rscData.Destination.ValueString(), ":") {
			delPrefix += "rib " + routingInstance + ".inet6.0 "
		}
	}

	configSet := []string{
		delPrefix + "generate route " + rscData.Destination.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
