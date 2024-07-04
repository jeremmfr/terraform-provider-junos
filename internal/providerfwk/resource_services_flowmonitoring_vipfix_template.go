package providerfwk

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	_ resource.Resource                   = &servicesFlowMonitoringVIPFixTemplate{}
	_ resource.ResourceWithConfigure      = &servicesFlowMonitoringVIPFixTemplate{}
	_ resource.ResourceWithValidateConfig = &servicesFlowMonitoringVIPFixTemplate{}
	_ resource.ResourceWithImportState    = &servicesFlowMonitoringVIPFixTemplate{}
	_ resource.ResourceWithUpgradeState   = &servicesFlowMonitoringVIPFixTemplate{}
)

type servicesFlowMonitoringVIPFixTemplate struct {
	client *junos.Client
}

func newServicesFlowMonitoringVIPFixTemplateResource() resource.Resource {
	return &servicesFlowMonitoringVIPFixTemplate{}
}

func (rsc *servicesFlowMonitoringVIPFixTemplate) typeName() string {
	return providerName + "_services_flowmonitoring_vipfix_template"
}

func (rsc *servicesFlowMonitoringVIPFixTemplate) junosName() string {
	return "services flow-monitoring version-ipfix template"
}

func (rsc *servicesFlowMonitoringVIPFixTemplate) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *servicesFlowMonitoringVIPFixTemplate) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *servicesFlowMonitoringVIPFixTemplate) Configure(
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

func (rsc *servicesFlowMonitoringVIPFixTemplate) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
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
				Description: "Name of flow-monitoring version-ipfix template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Type of template.",
				Validators: []validator.String{
					stringvalidator.OneOf("bridge-template", "ipv4-template", "ipv6-template", "mpls-template"),
				},
			},
			"flow_active_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Interval after which active flow is exported (10..600).",
				Validators: []validator.Int64{
					int64validator.Between(10, 600),
				},
			},
			"flow_inactive_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Period of inactivity that marks a flow inactive (10..600).",
				Validators: []validator.Int64{
					int64validator.Between(10, 600),
				},
			},
			"flow_key_flow_direction": schema.BoolAttribute{
				Optional:    true,
				Description: "Include flow direction.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"flow_key_output_interface": schema.BoolAttribute{
				Optional:    true,
				Description: "Include output interface.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"flow_key_vlan_id": schema.BoolAttribute{
				Optional:    true,
				Description: "Include vlan ID.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"ip_template_export_extension": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Export-extension for `ipv4-template`, `ipv6-template` type.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
			"mpls_template_label_position": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Description: "One or more MPLS label positions (1..8).",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueInt64sAre(
						int64validator.Between(1, 8),
					),
				},
			},
			"nexthop_learning_enable": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable nexthop learning.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"nexthop_learning_disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable nexthop learning.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"observation_domain_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Observation Domain Id (0..255).",
				Validators: []validator.Int64{
					int64validator.Between(0, 255),
				},
			},
			"option_template_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Options template id (1024..65535).",
				Validators: []validator.Int64{
					int64validator.Between(1024, 65535),
				},
			},
			"template_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Template id (1024..65535).",
				Validators: []validator.Int64{
					int64validator.Between(1024, 65535),
				},
			},
			"tunnel_observation_ipv4": schema.BoolAttribute{
				Optional:    true,
				Description: "Tunnel observation IPv4.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"tunnel_observation_ipv6": schema.BoolAttribute{
				Optional:    true,
				Description: "Tunnel observation IPv6.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"option_refresh_rate": schema.SingleNestedBlock{
				Description: "Declare `option-refresh-rate` configuration.",
				Attributes: map[string]schema.Attribute{
					"packets": schema.Int64Attribute{
						Optional:    true,
						Description: "In number of packets (1..480000)",
						Validators: []validator.Int64{
							int64validator.Between(1, 480000),
						},
					},
					"seconds": schema.Int64Attribute{
						Optional:    true,
						Description: "In number of seconds (10..600).",
						Validators: []validator.Int64{
							int64validator.Between(10, 600),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"template_refresh_rate": schema.SingleNestedBlock{
				Description: "Declare `template-refresh-rate` configuration.",
				Attributes: map[string]schema.Attribute{
					"packets": schema.Int64Attribute{
						Optional:    true,
						Description: "In number of packets (1..480000)",
						Validators: []validator.Int64{
							int64validator.Between(1, 480000),
						},
					},
					"seconds": schema.Int64Attribute{
						Optional:    true,
						Description: "In number of seconds (10..600).",
						Validators: []validator.Int64{
							int64validator.Between(10, 600),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
		},
	}
}

type servicesFlowMonitoringVIPFixTemplateData struct {
	ID                        types.String                                          `tfsdk:"id"`
	Name                      types.String                                          `tfsdk:"name"`
	Type                      types.String                                          `tfsdk:"type"`
	FlowActiveTimeout         types.Int64                                           `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout       types.Int64                                           `tfsdk:"flow_inactive_timeout"`
	FlowKeyFlowDirection      types.Bool                                            `tfsdk:"flow_key_flow_direction"`
	FlowKeyOutputInterface    types.Bool                                            `tfsdk:"flow_key_output_interface"`
	FlowKeyVlanID             types.Bool                                            `tfsdk:"flow_key_vlan_id"`
	IPTemplateExportExtension []types.String                                        `tfsdk:"ip_template_export_extension"`
	MPLSTemplateLabelPosition []types.Int64                                         `tfsdk:"mpls_template_label_position"`
	NexthopLearningEnable     types.Bool                                            `tfsdk:"nexthop_learning_enable"`
	NexthopLearningDisable    types.Bool                                            `tfsdk:"nexthop_learning_disable"`
	ObservationDomainID       types.Int64                                           `tfsdk:"observation_domain_id"`
	OptionTemplateID          types.Int64                                           `tfsdk:"option_template_id"`
	TemplateID                types.Int64                                           `tfsdk:"template_id"`
	OptionRefreshRate         *servicesFlowMonitoringVIPFixTemplateBlockRefreshRate `tfsdk:"option_refresh_rate"`
	TemplateRefreshRate       *servicesFlowMonitoringVIPFixTemplateBlockRefreshRate `tfsdk:"template_refresh_rate"`
	TunnelObservationIPv4     types.Bool                                            `tfsdk:"tunnel_observation_ipv4"`
	TunnelObservationIPv6     types.Bool                                            `tfsdk:"tunnel_observation_ipv6"`
}

type servicesFlowMonitoringVIPFixTemplateConfig struct {
	ID                        types.String                                          `tfsdk:"id"`
	Name                      types.String                                          `tfsdk:"name"`
	Type                      types.String                                          `tfsdk:"type"`
	FlowActiveTimeout         types.Int64                                           `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout       types.Int64                                           `tfsdk:"flow_inactive_timeout"`
	FlowKeyFlowDirection      types.Bool                                            `tfsdk:"flow_key_flow_direction"`
	FlowKeyOutputInterface    types.Bool                                            `tfsdk:"flow_key_output_interface"`
	FlowKeyVlanID             types.Bool                                            `tfsdk:"flow_key_vlan_id"`
	IPTemplateExportExtension types.Set                                             `tfsdk:"ip_template_export_extension"`
	MPLSTemplateLabelPosition types.List                                            `tfsdk:"mpls_template_label_position"`
	NexthopLearningEnable     types.Bool                                            `tfsdk:"nexthop_learning_enable"`
	NexthopLearningDisable    types.Bool                                            `tfsdk:"nexthop_learning_disable"`
	ObservationDomainID       types.Int64                                           `tfsdk:"observation_domain_id"`
	OptionTemplateID          types.Int64                                           `tfsdk:"option_template_id"`
	TemplateID                types.Int64                                           `tfsdk:"template_id"`
	OptionRefreshRate         *servicesFlowMonitoringVIPFixTemplateBlockRefreshRate `tfsdk:"option_refresh_rate"`
	TemplateRefreshRate       *servicesFlowMonitoringVIPFixTemplateBlockRefreshRate `tfsdk:"template_refresh_rate"`
	TunnelObservationIPv4     types.Bool                                            `tfsdk:"tunnel_observation_ipv4"`
	TunnelObservationIPv6     types.Bool                                            `tfsdk:"tunnel_observation_ipv6"`
}

type servicesFlowMonitoringVIPFixTemplateBlockRefreshRate struct {
	Packets types.Int64 `tfsdk:"packets"`
	Seconds types.Int64 `tfsdk:"seconds"`
}

func (rsc *servicesFlowMonitoringVIPFixTemplate) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config servicesFlowMonitoringVIPFixTemplateConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.NexthopLearningEnable.IsNull() && !config.NexthopLearningEnable.IsUnknown() &&
		!config.NexthopLearningDisable.IsNull() && !config.NexthopLearningDisable.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("nexthop_learning_enable"),
			tfdiag.ConflictConfigErrSummary,
			"cannot have nexthop_learning_enable and nexthop_learning_disable at the same time",
		)
	}
	if !config.IPTemplateExportExtension.IsNull() && !config.IPTemplateExportExtension.IsUnknown() &&
		!config.Type.IsNull() && !config.Type.IsUnknown() {
		if v := config.Type.ValueString(); v != "ipv4-template" && v != "ipv6-template" {
			resp.Diagnostics.AddAttributeError(
				path.Root("ip_template_export_extension"),
				tfdiag.ConflictConfigErrSummary,
				fmt.Sprintf("ip_template_export_extension not compatible with type %q", v),
			)
		}
	}
	if !config.MPLSTemplateLabelPosition.IsNull() && !config.MPLSTemplateLabelPosition.IsUnknown() &&
		!config.Type.IsNull() && !config.Type.IsUnknown() {
		if v := config.Type.ValueString(); v != "mpls-template" {
			resp.Diagnostics.AddAttributeError(
				path.Root("mpls_template_label_position"),
				tfdiag.ConflictConfigErrSummary,
				fmt.Sprintf("mpls_template_label_position not compatible with type %q", v),
			)
		}
	}
}

func (rsc *servicesFlowMonitoringVIPFixTemplate) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan servicesFlowMonitoringVIPFixTemplateData
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
	if plan.Type.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("type"),
			"Empty Type",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "type"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			templateExists, err := checkServicesFlowMonitoringVIPFixTemplateExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if templateExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			templateExists, err := checkServicesFlowMonitoringVIPFixTemplateExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !templateExists {
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

func (rsc *servicesFlowMonitoringVIPFixTemplate) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data servicesFlowMonitoringVIPFixTemplateData
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
		nil,
		resp,
	)
}

func (rsc *servicesFlowMonitoringVIPFixTemplate) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state servicesFlowMonitoringVIPFixTemplateData
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

func (rsc *servicesFlowMonitoringVIPFixTemplate) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state servicesFlowMonitoringVIPFixTemplateData
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

func (rsc *servicesFlowMonitoringVIPFixTemplate) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data servicesFlowMonitoringVIPFixTemplateData

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

func checkServicesFlowMonitoringVIPFixTemplateExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services flow-monitoring version-ipfix template \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *servicesFlowMonitoringVIPFixTemplateData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *servicesFlowMonitoringVIPFixTemplateData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *servicesFlowMonitoringVIPFixTemplateData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set services flow-monitoring version-ipfix template \"" + rscData.Name.ValueString() + "\" "

	configSet = append(configSet, setPrefix+rscData.Type.ValueString())
	for _, v := range rscData.IPTemplateExportExtension {
		if v2 := rscData.Type.ValueString(); v2 != "ipv4-template" && v2 != "ipv6-template" {
			return path.Root("ip_template_export_extension"),
				fmt.Errorf("ip_template_export_extension not compatible with type %q", v2)
		}
		configSet = append(configSet,
			setPrefix+rscData.Type.ValueString()+" export-extension "+v.ValueString())
	}
	for _, v := range rscData.MPLSTemplateLabelPosition {
		if v2 := rscData.Type.ValueString(); v2 != "mpls-template" {
			return path.Root("mpls_template_label_position"),
				fmt.Errorf("mpls_template_label_position not compatible with type %q", v2)
		}
		if !v.IsNull() {
			configSet = append(configSet,
				setPrefix+rscData.Type.ValueString()+" label-position "+utils.ConvI64toa(v.ValueInt64()))
		}
	}
	if !rscData.FlowActiveTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"flow-active-timeout "+
			utils.ConvI64toa(rscData.FlowActiveTimeout.ValueInt64()))
	}
	if !rscData.FlowInactiveTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"flow-inactive-timeout "+
			utils.ConvI64toa(rscData.FlowInactiveTimeout.ValueInt64()))
	}
	if rscData.FlowKeyFlowDirection.ValueBool() {
		configSet = append(configSet, setPrefix+"flow-key flow-direction")
	}
	if rscData.FlowKeyOutputInterface.ValueBool() {
		configSet = append(configSet, setPrefix+"flow-key output-interface")
	}
	if rscData.FlowKeyVlanID.ValueBool() {
		configSet = append(configSet, setPrefix+"flow-key vlan-id")
	}
	if rscData.NexthopLearningEnable.ValueBool() {
		configSet = append(configSet, setPrefix+"nexthop-learning enable")
	}
	if rscData.NexthopLearningDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"nexthop-learning disable")
	}
	if !rscData.ObservationDomainID.IsNull() {
		configSet = append(configSet, setPrefix+"observation-domain-id "+
			utils.ConvI64toa(rscData.ObservationDomainID.ValueInt64()))
	}
	if rscData.OptionRefreshRate != nil {
		configSet = append(configSet, setPrefix+"option-refresh-rate")
		if !rscData.OptionRefreshRate.Packets.IsNull() {
			configSet = append(configSet, setPrefix+"option-refresh-rate packets "+
				utils.ConvI64toa(rscData.OptionRefreshRate.Packets.ValueInt64()))
		}
		if !rscData.OptionRefreshRate.Seconds.IsNull() {
			configSet = append(configSet, setPrefix+"option-refresh-rate seconds "+
				utils.ConvI64toa(rscData.OptionRefreshRate.Seconds.ValueInt64()))
		}
	}
	if !rscData.OptionTemplateID.IsNull() {
		configSet = append(configSet, setPrefix+"option-template-id "+
			utils.ConvI64toa(rscData.OptionTemplateID.ValueInt64()))
	}
	if !rscData.TemplateID.IsNull() {
		configSet = append(configSet, setPrefix+"template-id "+
			utils.ConvI64toa(rscData.TemplateID.ValueInt64()))
	}
	if rscData.TemplateRefreshRate != nil {
		configSet = append(configSet, setPrefix+"template-refresh-rate")
		if !rscData.TemplateRefreshRate.Packets.IsNull() {
			configSet = append(configSet, setPrefix+"template-refresh-rate packets "+
				utils.ConvI64toa(rscData.TemplateRefreshRate.Packets.ValueInt64()))
		}
		if !rscData.TemplateRefreshRate.Seconds.IsNull() {
			configSet = append(configSet, setPrefix+"template-refresh-rate seconds "+
				utils.ConvI64toa(rscData.TemplateRefreshRate.Seconds.ValueInt64()))
		}
	}
	if rscData.TunnelObservationIPv4.ValueBool() {
		configSet = append(configSet, setPrefix+"tunnel-observation ipv4")
	}
	if rscData.TunnelObservationIPv6.ValueBool() {
		configSet = append(configSet, setPrefix+"tunnel-observation ipv6")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *servicesFlowMonitoringVIPFixTemplateData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services flow-monitoring version-ipfix template \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case slices.Contains([]string{"bridge-template", "ipv4-template", "ipv6-template", "mpls-template"}, itemTrim):
				rscData.Type = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "ipv6-template export-extension "):
				rscData.Type = types.StringValue("ipv6-template")
				rscData.IPTemplateExportExtension = append(rscData.IPTemplateExportExtension, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "ipv4-template export-extension "):
				rscData.Type = types.StringValue("ipv4-template")
				rscData.IPTemplateExportExtension = append(rscData.IPTemplateExportExtension, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "mpls-template label-position "):
				rscData.Type = types.StringValue("mpls-template")
				number, err := tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
				rscData.MPLSTemplateLabelPosition = append(rscData.MPLSTemplateLabelPosition, number)
			case balt.CutPrefixInString(&itemTrim, "flow-active-timeout "):
				rscData.FlowActiveTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "flow-inactive-timeout "):
				rscData.FlowInactiveTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "flow-key flow-direction":
				rscData.FlowKeyFlowDirection = types.BoolValue(true)
			case itemTrim == "flow-key output-interface":
				rscData.FlowKeyOutputInterface = types.BoolValue(true)
			case itemTrim == "flow-key vlan-id":
				rscData.FlowKeyVlanID = types.BoolValue(true)
			case itemTrim == "nexthop-learning enable":
				rscData.NexthopLearningEnable = types.BoolValue(true)
			case itemTrim == "nexthop-learning disable":
				rscData.NexthopLearningDisable = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "observation-domain-id "):
				rscData.ObservationDomainID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "option-refresh-rate"):
				if rscData.OptionRefreshRate == nil {
					rscData.OptionRefreshRate = &servicesFlowMonitoringVIPFixTemplateBlockRefreshRate{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " packets "):
					rscData.OptionRefreshRate.Packets, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " seconds "):
					rscData.OptionRefreshRate.Seconds, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "option-template-id "):
				rscData.OptionTemplateID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "template-refresh-rate"):
				if rscData.TemplateRefreshRate == nil {
					rscData.TemplateRefreshRate = &servicesFlowMonitoringVIPFixTemplateBlockRefreshRate{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " packets "):
					rscData.TemplateRefreshRate.Packets, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " seconds "):
					rscData.TemplateRefreshRate.Seconds, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "template-id "):
				rscData.TemplateID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "tunnel-observation ipv4":
				rscData.TunnelObservationIPv4 = types.BoolValue(true)
			case itemTrim == "tunnel-observation ipv6":
				rscData.TunnelObservationIPv6 = types.BoolValue(true)
			}
		}
	}

	return nil
}

func (rscData *servicesFlowMonitoringVIPFixTemplateData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete services flow-monitoring version-ipfix template \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
