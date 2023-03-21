package providerfwk

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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

func newRoutingInstance() resource.Resource {
	return &routingInstance{}
}

func (rsc *routingInstance) typeName() string {
	return providerName + "_routing_instance"
}

func (rsc *routingInstance) junosName() string {
	return "routing instance"
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
		Description: "Provides a " + rsc.junosName() + ".",
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
				Description: "Type of routing instance.",
				PlanModifiers: []planmodifier.String{
					tfplanmodifier.StringDefault("virtual-router"),
				},
				Validators: []validator.String{
					tfvalidator.StringSpaceExclusion(),
				},
			},
			"as": schema.StringAttribute{
				Optional:    true,
				Description: "Autonomous system number in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^\d+(\.\d+)?$`),
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
			"route_distinguisher": schema.StringAttribute{
				Optional:    true,
				Description: "Route distinguisher for this instance.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^(\d|\.)+L?:\d+$`),
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
					stringvalidator.RegexMatches(regexp.MustCompile(`^target:(\d|\.)+L?:\d+$`),
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
					stringvalidator.RegexMatches(regexp.MustCompile(`^target:(\d|\.)+L?:\d+$`),
						"must be use format 'target:x:y' where 'x' is an AS number followed by an optional 'L' (To indicate 4 byte AS), "+
							"or an IP address and 'y' is a number. e.g. target:123456L:100"),
				},
			},
			"vrf_target_import": schema.StringAttribute{
				Optional:    true,
				Description: "Target community to use when filtering on import.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^target:(\d|\.)+L?:\d+$`),
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
	ConfigureRDVrfOptSingly types.Bool     `tfsdk:"configure_rd_vrfopts_singly"`
	ConfigureTypeSingly     types.Bool     `tfsdk:"configure_type_singly"`
	VRFTargetAuto           types.Bool     `tfsdk:"vrf_target_auto"`
	ID                      types.String   `tfsdk:"id"`
	Name                    types.String   `tfsdk:"name"`
	Type                    types.String   `tfsdk:"type"`
	AS                      types.String   `tfsdk:"as"`
	Description             types.String   `tfsdk:"description"`
	InstanceExport          []types.String `tfsdk:"instance_export"`
	InstanceImport          []types.String `tfsdk:"instance_import"`
	RouteDistinguisher      types.String   `tfsdk:"route_distinguisher"`
	RouterID                types.String   `tfsdk:"router_id"`
	VRFExport               []types.String `tfsdk:"vrf_export"`
	VRFImport               []types.String `tfsdk:"vrf_import"`
	VRFTarget               types.String   `tfsdk:"vrf_target"`
	VRFTargetExport         types.String   `tfsdk:"vrf_target_export"`
	VRFTargetImport         types.String   `tfsdk:"vrf_target_import"`
	VTEPSourceInterface     types.String   `tfsdk:"vtep_source_interface"`
	Interface               []types.String `tfsdk:"-"` // to data source
}

type routingInstanceConfig struct {
	ConfigureRDVrfOptSingly types.Bool   `tfsdk:"configure_rd_vrfopts_singly"`
	ConfigureTypeSingly     types.Bool   `tfsdk:"configure_type_singly"`
	VRFTargetAuto           types.Bool   `tfsdk:"vrf_target_auto"`
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Type                    types.String `tfsdk:"type"`
	AS                      types.String `tfsdk:"as"`
	Description             types.String `tfsdk:"description"`
	InstanceExport          types.List   `tfsdk:"instance_export"`
	InstanceImport          types.List   `tfsdk:"instance_import"`
	RouteDistinguisher      types.String `tfsdk:"route_distinguisher"`
	RouterID                types.String `tfsdk:"router_id"`
	VRFExport               types.List   `tfsdk:"vrf_export"`
	VRFImport               types.List   `tfsdk:"vrf_import"`
	VRFTarget               types.String `tfsdk:"vrf_target"`
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

	if !config.ConfigureRDVrfOptSingly.IsNull() &&
		(!config.RouteDistinguisher.IsNull() ||
			!config.VRFExport.IsNull() ||
			!config.VRFImport.IsNull() ||
			!config.VRFTarget.IsNull() ||
			!config.VRFTargetAuto.IsNull() ||
			!config.VRFTargetExport.IsNull() ||
			!config.VRFTargetImport.IsNull()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("configure_rd_vrfopts_singly"),
			"Conflict Configuration Error",
			"cannot have configure_rd_vrfopts_singly and want to configure route-distinguisher or vrf options at the same time",
		)
	}
	if !config.ConfigureTypeSingly.IsNull() {
		if config.Type.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("configure_type_singly"),
				"Missing Configuration Error",
				"type must specified with empty string when configure_type_singly is enabled",
			)
		} else if !config.Type.IsUnknown() && config.Type.ValueString() != "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("type"),
				"Conflict Configuration Error",
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
			"could not create "+rsc.junosName()+" with empty name",
		)

		return
	}
	if plan.ConfigureTypeSingly.ValueBool() && plan.Type.ValueString() != "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("type"),
			"Conflict Configuration Error",
			"type must specified with empty string when configure_type_singly is enabled",
		)

		return
	}

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	instanceExists, err := checkRoutingInstanceExists(ctx, plan.Name.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Pre Check Error", err.Error())

		return
	}
	if instanceExists {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError(
			"Duplicate Configuration Error",
			fmt.Sprintf(rsc.junosName()+" %q already exists", plan.Name.ValueString()),
		)

		return
	}

	if errPath, err := plan.set(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("create resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	instanceExists, err = checkRoutingInstanceExists(ctx, plan.Name.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.AddError("Post Check Error", err.Error())

		return
	}
	if !instanceExists {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf(rsc.junosName()+" %q does not exists after commit "+
				"=> check your config", plan.Name.ValueString()),
		)

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *routingInstance) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data routingInstanceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	err = data.read(ctx, state.Name.ValueString(), junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}
	if data.ID.IsNull() {
		resp.State.RemoveResource(ctx)

		return
	}

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
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.delOpts(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	if err := state.delOpts(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("update resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *routingInstance) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state routingInstanceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeDeleteAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	warns, err := junSess.CommitConf("delete resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}
}

func (rsc *routingInstance) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	var data routingInstanceData
	if err := data.read(ctx, req.ID, junSess); err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf("don't find "+rsc.junosName()+" with id %q "+
				"(id must be <name>)", req.ID),
		)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkRoutingInstanceExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + name + junos.PipeDisplaySet)
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

func (rscData *routingInstanceData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetRoutingInstances + rscData.Name.ValueString() + " "

	if rscData.ConfigureTypeSingly.ValueBool() {
		if rscData.Type.ValueString() != "" {
			return path.Root("type"),
				fmt.Errorf("if `configure_type_singly` = true, `type` need to be set to empty value to avoid confusion")
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
		configSet = append(configSet, setPrefix+"routing-options autonomous-system "+v)
	}
	if v := rscData.RouterID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-options router-id "+v)
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	for _, v := range rscData.InstanceExport {
		configSet = append(configSet, setPrefix+"routing-options instance-export \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.InstanceImport {
		configSet = append(configSet, setPrefix+"routing-options instance-import \""+v.ValueString()+"\"")
	}
	if v := rscData.VTEPSourceInterface.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"vtep-source-interface "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *routingInstanceData) read(
	_ context.Context, name string, junSess *junos.Session,
) (
	err error,
) {
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
			case balt.CutPrefixInString(&itemTrim, "route-distinguisher "):
				rscData.RouteDistinguisher = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "routing-options autonomous-system "):
				rscData.AS = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "routing-options instance-export "):
				rscData.InstanceExport = append(rscData.InstanceExport,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "routing-options instance-import "):
				rscData.InstanceImport = append(rscData.InstanceImport,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "routing-options router-id "):
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
	configSet := make([]string, 0)
	setPrefix := junos.DelRoutingInstances + rscData.Name.ValueString() + " "
	configSet = append(configSet,
		setPrefix+"description",
		setPrefix+"routing-options autonomous-system",
		setPrefix+"routing-options instance-export",
		setPrefix+"routing-options instance-import",
		setPrefix+"routing-options router-id",
		setPrefix+"vtep-source-interface",
	)
	if !rscData.ConfigureTypeSingly.ValueBool() {
		configSet = append(configSet, setPrefix+"instance-type")
	}
	if !rscData.ConfigureRDVrfOptSingly.ValueBool() {
		configSet = append(configSet,
			setPrefix+"route-distinguisher",
			setPrefix+"vrf-export",
			setPrefix+"vrf-import",
			setPrefix+"vrf_target",
		)
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *routingInstanceData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		junos.DelRoutingInstances + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
