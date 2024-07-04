package providerfwk

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	_ resource.Resource                   = &routingInstance{}
	_ resource.ResourceWithConfigure      = &routingInstance{}
	_ resource.ResourceWithValidateConfig = &routingInstance{}
	_ resource.ResourceWithImportState    = &routingInstance{}
)

type routingInstance struct {
	client *junos.Client
}

func newRoutingInstanceResource() resource.Resource {
	return &routingInstance{}
}

func (rsc *routingInstance) typeName() string {
	return providerName + "_routing_instance"
}

func (rsc *routingInstance) junosName() string {
	return "routing instance"
}

func (rsc *routingInstance) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *routingInstance) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *routingInstance) Configure(
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

func (rsc *routingInstance) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
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
				Description: "The name of routing instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					stringvalidator.NoneOfCaseInsensitive(junos.DefaultW),
				},
			},
			"configure_rd_vrfopts_singly": schema.BoolAttribute{
				Optional:    true,
				Description: "Configure `route-distinguisher` and `vrf-*` options in other resource (like `junos_evpn`).",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"configure_type_singly": schema.BoolAttribute{
				Optional:    true,
				Description: "Configure `instance-type` option in other resource (like `junos_evpn`).",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("virtual-router"),
				Description: "Type of routing instance.",
				Validators: []validator.String{
					tfvalidator.StringSpaceExclusion(),
				},
			},
			"as": schema.StringAttribute{
				Optional:    true,
				Description: "Autonomous system number in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d+(\.\d+)?$`),
						"must be in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of routing instance.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"instance_export": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Export policy for instance RIBs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"instance_import": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Import policy for instance RIBs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"remote_vtep_list": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Configure static remote VXLAN tunnel endpoints.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringIPAddress().IPv4Only(),
					),
				},
			},
			"remote_vtep_v6_list": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Configure static ipv6 remote VXLAN tunnel endpoints.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringIPAddress().IPv6Only(),
					),
				},
			},
			"route_distinguisher": schema.StringAttribute{
				Optional:    true,
				Description: "Route distinguisher for this instance.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(\d|\.)+L?:\d+$`),
						"must be use format 'x:y' where 'x' is an AS number followed by an optional 'L' (To indicate 4 byte AS), "+
							"or an IP address and 'y' is a number. e.g. 123456L:100"),
				},
			},
			"router_id": schema.StringAttribute{
				Optional:    true,
				Description: "Router identifier.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress().IPv4Only(),
				},
			},
			"vrf_export": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Export policy for VRF instance RIBs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"vrf_import": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Import policy for VRF instance RIBs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"vrf_target": schema.StringAttribute{
				Optional:    true,
				Description: "Target community to use in import and export.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^target:(\d|\.)+L?:\d+$`),
						"must be use format 'target:x:y' where 'x' is an AS number followed by an optional 'L' (To indicate 4 byte AS), "+
							"or an IP address and 'y' is a number. e.g. target:123456L:100"),
				},
			},
			"vrf_target_auto": schema.BoolAttribute{
				Optional:    true,
				Description: "Auto derive import and export target community from BGP AS & L2.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"vrf_target_export": schema.StringAttribute{
				Optional:    true,
				Description: "Target community to use when marking routes on export.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^target:(\d|\.)+L?:\d+$`),
						"must be use format 'target:x:y' where 'x' is an AS number followed by an optional 'L' (To indicate 4 byte AS), "+
							"or an IP address and 'y' is a number. e.g. target:123456L:100"),
				},
			},
			"vrf_target_import": schema.StringAttribute{
				Optional:    true,
				Description: "Target community to use when filtering on import.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^target:(\d|\.)+L?:\d+$`),
						"must be use format 'target:x:y' where 'x' is an AS number followed by an optional 'L' (To indicate 4 byte AS), "+
							"or an IP address and 'y' is a number. e.g. target:123456L:100"),
				},
			},
			"vtep_source_interface": schema.StringAttribute{
				Optional:    true,
				Description: "Source layer-3 IFL for VXLAN.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					tfvalidator.String1DotCount(),
				},
			},
		},
	}
}

type routingInstanceData struct {
	ID                      types.String   `tfsdk:"id"`
	Name                    types.String   `tfsdk:"name"`
	ConfigureRDVrfOptSingly types.Bool     `tfsdk:"configure_rd_vrfopts_singly"`
	ConfigureTypeSingly     types.Bool     `tfsdk:"configure_type_singly"`
	Type                    types.String   `tfsdk:"type"`
	AS                      types.String   `tfsdk:"as"`
	Description             types.String   `tfsdk:"description"`
	InstanceExport          []types.String `tfsdk:"instance_export"`
	InstanceImport          []types.String `tfsdk:"instance_import"`
	RemoteVtepList          []types.String `tfsdk:"remote_vtep_list"`
	RemoteVtepV6List        []types.String `tfsdk:"remote_vtep_v6_list"`
	RouteDistinguisher      types.String   `tfsdk:"route_distinguisher"`
	RouterID                types.String   `tfsdk:"router_id"`
	VRFExport               []types.String `tfsdk:"vrf_export"`
	VRFImport               []types.String `tfsdk:"vrf_import"`
	VRFTarget               types.String   `tfsdk:"vrf_target"`
	VRFTargetAuto           types.Bool     `tfsdk:"vrf_target_auto"`
	VRFTargetExport         types.String   `tfsdk:"vrf_target_export"`
	VRFTargetImport         types.String   `tfsdk:"vrf_target_import"`
	VTEPSourceInterface     types.String   `tfsdk:"vtep_source_interface"`
	Interface               []types.String `tfsdk:"-"` // to data source
}

type routingInstanceConfig struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	ConfigureRDVrfOptSingly types.Bool   `tfsdk:"configure_rd_vrfopts_singly"`
	ConfigureTypeSingly     types.Bool   `tfsdk:"configure_type_singly"`
	Type                    types.String `tfsdk:"type"`
	AS                      types.String `tfsdk:"as"`
	Description             types.String `tfsdk:"description"`
	InstanceExport          types.List   `tfsdk:"instance_export"`
	InstanceImport          types.List   `tfsdk:"instance_import"`
	RemoteVtepList          types.Set    `tfsdk:"remote_vtep_list"`
	RemoteVtepV6List        types.Set    `tfsdk:"remote_vtep_v6_list"`
	RouteDistinguisher      types.String `tfsdk:"route_distinguisher"`
	RouterID                types.String `tfsdk:"router_id"`
	VRFExport               types.List   `tfsdk:"vrf_export"`
	VRFImport               types.List   `tfsdk:"vrf_import"`
	VRFTarget               types.String `tfsdk:"vrf_target"`
	VRFTargetAuto           types.Bool   `tfsdk:"vrf_target_auto"`
	VRFTargetExport         types.String `tfsdk:"vrf_target_export"`
	VRFTargetImport         types.String `tfsdk:"vrf_target_import"`
	VTEPSourceInterface     types.String `tfsdk:"vtep_source_interface"`
}

func (rsc *routingInstance) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config routingInstanceConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ConfigureRDVrfOptSingly.ValueBool() &&
		((!config.RouteDistinguisher.IsNull() && !config.RouteDistinguisher.IsUnknown()) ||
			(!config.VRFExport.IsNull() && !config.VRFExport.IsUnknown()) ||
			(!config.VRFImport.IsNull() && !config.VRFImport.IsUnknown()) ||
			(!config.VRFTarget.IsNull() && !config.VRFTarget.IsUnknown()) ||
			(!config.VRFTargetAuto.IsNull() && !config.VRFTargetAuto.IsUnknown()) ||
			(!config.VRFTargetExport.IsNull() && !config.VRFTargetExport.IsUnknown()) ||
			(!config.VRFTargetImport.IsNull() && !config.VRFTargetImport.IsUnknown())) {
		resp.Diagnostics.AddAttributeError(
			path.Root("configure_rd_vrfopts_singly"),
			tfdiag.ConflictConfigErrSummary,
			"cannot have configure_rd_vrfopts_singly and want to configure route-distinguisher or vrf options at the same time",
		)
	}
	if !config.ConfigureTypeSingly.IsNull() {
		if config.Type.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("configure_type_singly"),
				tfdiag.MissingConfigErrSummary,
				"type must specified with empty string when configure_type_singly is enabled",
			)
		} else if !config.ConfigureTypeSingly.IsUnknown() &&
			!config.Type.IsUnknown() && config.Type.ValueString() != "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("type"),
				tfdiag.ConflictConfigErrSummary,
				"type must specified with empty string when configure_type_singly is enabled",
			)
		}
	}
}

func (rsc *routingInstance) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan routingInstanceData
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
	if plan.ConfigureTypeSingly.ValueBool() && plan.Type.ValueString() != "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("type"),
			tfdiag.ConflictConfigErrSummary,
			"type must specified with empty string when configure_type_singly is enabled",
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			instanceExists, err := checkRoutingInstanceExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if instanceExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			instanceExists, err := checkRoutingInstanceExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !instanceExists {
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

func (rsc *routingInstance) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data routingInstanceData
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
		func() {
			data.ConfigureRDVrfOptSingly = state.ConfigureRDVrfOptSingly
			if data.ConfigureRDVrfOptSingly.ValueBool() {
				data.RouteDistinguisher = types.StringNull()
				data.VRFExport = nil
				data.VRFImport = nil
				data.VRFTarget = types.StringNull()
				data.VRFTargetAuto = types.BoolNull()
				data.VRFTargetExport = types.StringNull()
				data.VRFTargetImport = types.StringNull()
			}
			data.ConfigureTypeSingly = state.ConfigureTypeSingly
			if data.ConfigureTypeSingly.ValueBool() {
				data.Type = types.StringValue("")
			}
		},
		resp,
	)
}

func (rsc *routingInstance) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state routingInstanceData
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

func (rsc *routingInstance) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state routingInstanceData
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

func (rsc *routingInstance) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data routingInstanceData

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

func checkRoutingInstanceExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		junos.RoutingInstancesWS + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *routingInstanceData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *routingInstanceData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *routingInstanceData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS + junos.RoutingInstancesWS + rscData.Name.ValueString() + " "

	if rscData.ConfigureTypeSingly.ValueBool() {
		if rscData.Type.ValueString() != "" {
			return path.Root("type"),
				errors.New("if `configure_type_singly` = true, `type` need to be set to empty value to avoid confusion")
		}
	} else {
		if v := rscData.Type.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"instance-type "+v)
		}
	}
	if !rscData.ConfigureRDVrfOptSingly.ValueBool() {
		if v := rscData.RouteDistinguisher.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"route-distinguisher "+v)
		}
		for _, v := range rscData.VRFExport {
			configSet = append(configSet, setPrefix+"vrf-export \""+v.ValueString()+"\"")
		}
		for _, v := range rscData.VRFImport {
			configSet = append(configSet, setPrefix+"vrf-import \""+v.ValueString()+"\"")
		}
		if v := rscData.VRFTarget.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"vrf-target "+v)
		}
		if rscData.VRFTargetAuto.ValueBool() {
			configSet = append(configSet, setPrefix+"vrf-target auto")
		}
		if v := rscData.VRFTargetExport.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"vrf-target export "+v)
		}
		if v := rscData.VRFTargetImport.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"vrf-target import "+v)
		}
	}
	if v := rscData.AS.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+junos.RoutingOptionsWS+"autonomous-system "+v)
	}
	if v := rscData.RouterID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+junos.RoutingOptionsWS+"router-id "+v)
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	for _, v := range rscData.InstanceExport {
		configSet = append(configSet, setPrefix+junos.RoutingOptionsWS+"instance-export \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.InstanceImport {
		configSet = append(configSet, setPrefix+junos.RoutingOptionsWS+"instance-import \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.RemoteVtepList {
		configSet = append(configSet, setPrefix+"remote-vtep-list "+v.ValueString())
	}
	for _, v := range rscData.RemoteVtepV6List {
		configSet = append(configSet, setPrefix+"remote-vtep-v6-list "+v.ValueString())
	}
	if v := rscData.VTEPSourceInterface.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"vtep-source-interface "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *routingInstanceData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		junos.RoutingInstancesWS + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "instance-type "):
				rscData.Type = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "remote-vtep-list "):
				rscData.RemoteVtepList = append(rscData.RemoteVtepList, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "remote-vtep-v6-list "):
				rscData.RemoteVtepV6List = append(rscData.RemoteVtepV6List, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "route-distinguisher "):
				rscData.RouteDistinguisher = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, junos.RoutingOptionsWS+"autonomous-system "):
				rscData.AS = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, junos.RoutingOptionsWS+"instance-export "):
				rscData.InstanceExport = append(rscData.InstanceExport,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, junos.RoutingOptionsWS+"instance-import "):
				rscData.InstanceImport = append(rscData.InstanceImport,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, junos.RoutingOptionsWS+"router-id "):
				rscData.RouterID = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "vrf-export "):
				rscData.VRFExport = append(rscData.VRFExport,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "vrf-import "):
				rscData.VRFImport = append(rscData.VRFImport,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case itemTrim == "vrf-target auto":
				rscData.VRFTargetAuto = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "vrf-target export "):
				rscData.VRFTargetExport = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "vrf-target import "):
				rscData.VRFTargetImport = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "vrf-target "):
				rscData.VRFTarget = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "vtep-source-interface "):
				rscData.VTEPSourceInterface = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "interface "):
				rscData.Interface = append(rscData.Interface,
					types.StringValue(strings.Split(itemTrim, " ")[0]))
			}
		}
	}

	return nil
}

func (rscData *routingInstanceData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS + junos.RoutingInstancesWS + rscData.Name.ValueString() + " "

	configSet := []string{
		delPrefix + "description",
		delPrefix + junos.RoutingOptionsWS + "autonomous-system",
		delPrefix + junos.RoutingOptionsWS + "instance-export",
		delPrefix + junos.RoutingOptionsWS + "instance-import",
		delPrefix + "remote-vtep-list",
		delPrefix + "remote-vtep-v6-list",
		delPrefix + junos.RoutingOptionsWS + "router-id",
		delPrefix + "vtep-source-interface",
	}
	if !rscData.ConfigureTypeSingly.ValueBool() {
		configSet = append(configSet, delPrefix+"instance-type")
	}
	if !rscData.ConfigureRDVrfOptSingly.ValueBool() {
		configSet = append(configSet,
			delPrefix+"route-distinguisher",
			delPrefix+"vrf-export",
			delPrefix+"vrf-import",
			delPrefix+"vrf_target",
		)
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *routingInstanceData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		junos.DeleteLS + junos.RoutingInstancesWS + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
