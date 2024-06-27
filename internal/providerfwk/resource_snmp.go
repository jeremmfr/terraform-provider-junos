package providerfwk

import (
	"context"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource                   = &snmp{}
	_ resource.ResourceWithConfigure      = &snmp{}
	_ resource.ResourceWithValidateConfig = &snmp{}
	_ resource.ResourceWithImportState    = &snmp{}
	_ resource.ResourceWithUpgradeState   = &snmp{}
)

type snmp struct {
	client *junos.Client
}

func newSnmpResource() resource.Resource {
	return &snmp{}
}

func (rsc *snmp) typeName() string {
	return providerName + "_snmp"
}

func (rsc *snmp) junosName() string {
	return "snmp"
}

func (rsc *snmp) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *snmp) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *snmp) Configure(
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

func (rsc *snmp) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `snmp`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"clean_on_destroy": schema.BoolAttribute{
				Optional:    true,
				Description: "Clean supported lines when destroy this resource.",
			},
			"arp": schema.BoolAttribute{
				Optional:    true,
				Description: "JVision ARP.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"arp_host_name_resolution": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable host name resolution for JVision ARP.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"contact": schema.StringAttribute{
				Optional:    true,
				Description: "Contact information for administrator.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "System description.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"engine_id": schema.StringAttribute{
				Optional:    true,
				Description: "SNMPv3 engine ID.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(use-default-ip-address|use-mac-address|local .+)$`),
						"must have 'use-default-ip-address', 'use-mac-address' or 'local ...'",
					),
				},
			},
			"filter_duplicates": schema.BoolAttribute{
				Optional:    true,
				Description: "Filter requests with duplicate source address/port and request ID.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"filter_interfaces": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Regular expressions to list of interfaces that needs to be filtered.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"filter_internal_interfaces": schema.BoolAttribute{
				Optional:    true,
				Description: "Filter all internal interfaces.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"if_count_with_filter_interfaces": schema.BoolAttribute{
				Optional:    true,
				Description: "Filter interfaces config for ifNumber and ipv6Interfaces.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"interface": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Restrict SNMP requests to interfaces.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
						tfvalidator.String1DotCount(),
					),
				},
			},
			"location": schema.StringAttribute{
				Optional:    true,
				Description: "Physical location of system.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"routing_instance_access": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable SNMP routing instance.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"routing_instance_access_list": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Allow/Deny SNMP access to routing instances.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 63),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
						stringvalidator.NoneOfCaseInsensitive(junos.DefaultW),
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"health_monitor": schema.SingleNestedBlock{
				Description: "Enable `health-monitor`.",
				Attributes: map[string]schema.Attribute{
					"falling_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Falling threshold applied to all monitored objects.",
						Validators: []validator.Int64{
							int64validator.Between(0, 100),
						},
					},
					"idp": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IDP health monitor.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"idp_falling_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Falling threshold applied to all idp monitored objects.",
						Validators: []validator.Int64{
							int64validator.Between(0, 100),
						},
					},
					"idp_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Interval between idp samples.",
						Validators: []validator.Int64{
							int64validator.Between(1, 2147483647),
						},
					},
					"idp_rising_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Rising threshold applied to all monitored idp objects.",
						Validators: []validator.Int64{
							int64validator.Between(0, 100),
						},
					},
					"interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Interval between samples.",
						Validators: []validator.Int64{
							int64validator.Between(1, 2147483647),
						},
					},
					"rising_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Rising threshold applied to all monitored objects.",
						Validators: []validator.Int64{
							int64validator.Between(0, 100),
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

type snmpData struct {
	ID                          types.String            `tfsdk:"id"`
	CleanOnDestroy              types.Bool              `tfsdk:"clean_on_destroy"`
	ARP                         types.Bool              `tfsdk:"arp"`
	ARPHostNameResolution       types.Bool              `tfsdk:"arp_host_name_resolution"`
	Contact                     types.String            `tfsdk:"contact"`
	Description                 types.String            `tfsdk:"description"`
	EngineID                    types.String            `tfsdk:"engine_id"`
	FilterDuplicates            types.Bool              `tfsdk:"filter_duplicates"`
	FilterInterfaces            []types.String          `tfsdk:"filter_interfaces"`
	FilterInternalInterfaces    types.Bool              `tfsdk:"filter_internal_interfaces"`
	IfCountWithFilterInterfaces types.Bool              `tfsdk:"if_count_with_filter_interfaces"`
	Interface                   []types.String          `tfsdk:"interface"`
	Location                    types.String            `tfsdk:"location"`
	RoutingInstanceAccess       types.Bool              `tfsdk:"routing_instance_access"`
	RoutingInstanceAccessList   []types.String          `tfsdk:"routing_instance_access_list"`
	HealthMonitor               *snmpBlockHealthMonitor `tfsdk:"health_monitor"`
}

type snmpConfig struct {
	ID                          types.String            `tfsdk:"id"`
	CleanOnDestroy              types.Bool              `tfsdk:"clean_on_destroy"`
	ARP                         types.Bool              `tfsdk:"arp"`
	ARPHostNameResolution       types.Bool              `tfsdk:"arp_host_name_resolution"`
	Contact                     types.String            `tfsdk:"contact"`
	Description                 types.String            `tfsdk:"description"`
	EngineID                    types.String            `tfsdk:"engine_id"`
	FilterDuplicates            types.Bool              `tfsdk:"filter_duplicates"`
	FilterInterfaces            types.Set               `tfsdk:"filter_interfaces"`
	FilterInternalInterfaces    types.Bool              `tfsdk:"filter_internal_interfaces"`
	IfCountWithFilterInterfaces types.Bool              `tfsdk:"if_count_with_filter_interfaces"`
	Interface                   types.Set               `tfsdk:"interface"`
	Location                    types.String            `tfsdk:"location"`
	RoutingInstanceAccess       types.Bool              `tfsdk:"routing_instance_access"`
	RoutingInstanceAccessList   types.Set               `tfsdk:"routing_instance_access_list"`
	HealthMonitor               *snmpBlockHealthMonitor `tfsdk:"health_monitor"`
}

type snmpBlockHealthMonitor struct {
	FallingThreshold    types.Int64 `tfsdk:"falling_threshold"`
	Idp                 types.Bool  `tfsdk:"idp"`
	IdpFallingThreshold types.Int64 `tfsdk:"idp_falling_threshold"`
	IdpInterval         types.Int64 `tfsdk:"idp_interval"`
	IdpRisingThreshold  types.Int64 `tfsdk:"idp_rising_threshold"`
	Interval            types.Int64 `tfsdk:"interval"`
	RisingThreshold     types.Int64 `tfsdk:"rising_threshold"`
}

func (rsc *snmp) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config snmpConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.ARPHostNameResolution.IsNull() && !config.ARPHostNameResolution.IsUnknown() &&
		config.ARP.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("arp_host_name_resolution"),
			tfdiag.MissingConfigErrSummary,
			"arp must be specified with arp_host_name_resolution",
		)
	}
	if !config.RoutingInstanceAccessList.IsNull() && !config.RoutingInstanceAccessList.IsUnknown() &&
		config.RoutingInstanceAccess.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("routing_instance_access_list"),
			tfdiag.MissingConfigErrSummary,
			"routing_instance_access must be specified with routing_instance_access_list",
		)
	}
	if config.HealthMonitor != nil {
		if config.HealthMonitor.Idp.IsNull() {
			if !config.HealthMonitor.IdpFallingThreshold.IsNull() && !config.HealthMonitor.IdpFallingThreshold.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("health_monitor").AtName("idp_falling_threshold"),
					tfdiag.MissingConfigErrSummary,
					"idp must be specified with idp_falling_threshold"+
						" in health_monitor block",
				)
			}
			if !config.HealthMonitor.IdpInterval.IsNull() && !config.HealthMonitor.IdpInterval.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("health_monitor").AtName("idp_interval"),
					tfdiag.MissingConfigErrSummary,
					"idp must be specified with idp_interval"+
						" in health_monitor block",
				)
			}
			if !config.HealthMonitor.IdpRisingThreshold.IsNull() && !config.HealthMonitor.IdpRisingThreshold.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("health_monitor").AtName("idp_rising_threshold"),
					tfdiag.MissingConfigErrSummary,
					"idp must be specified with idp_rising_threshold"+
						" in health_monitor block",
				)
			}
		}
	}
}

func (rsc *snmp) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan snmpData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		nil,
		nil,
		&plan,
		resp,
	)
}

func (rsc *snmp) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data snmpData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadWithoutArg = &data
	defaultResourceRead(
		ctx,
		rsc,
		nil,
		&data,
		func() {
			data.CleanOnDestroy = state.CleanOnDestroy
		},
		resp,
	)
}

func (rsc *snmp) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state snmpData
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

func (rsc *snmp) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state snmpData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.CleanOnDestroy.ValueBool() {
		defaultResourceDelete(
			ctx,
			rsc,
			&state,
			resp,
		)
	}
}

func (rsc *snmp) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data snmpData

	var _ resourceDataReadWithoutArg = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		"",
	)
}

func (rscData *snmpData) fillID() {
	rscData.ID = types.StringValue("snmp")
}

func (rscData *snmpData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *snmpData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set snmp "

	if rscData.ARP.ValueBool() {
		configSet = append(configSet, setPrefix+"arp")
	}
	if rscData.ARPHostNameResolution.ValueBool() {
		configSet = append(configSet, setPrefix+"arp host-name-resolution")
	}
	if v := rscData.Contact.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"contact \""+v+"\"")
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := rscData.EngineID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"engine-id "+v)
	}
	if rscData.FilterDuplicates.ValueBool() {
		configSet = append(configSet, setPrefix+"filter-duplicates")
	}
	for _, v := range rscData.FilterInterfaces {
		configSet = append(configSet, setPrefix+"filter-interfaces interfaces \""+v.ValueString()+"\"")
	}
	if rscData.FilterInternalInterfaces.ValueBool() {
		configSet = append(configSet, setPrefix+"filter-interfaces all-internal-interfaces")
	}
	if rscData.IfCountWithFilterInterfaces.ValueBool() {
		configSet = append(configSet, setPrefix+"if-count-with-filter-interfaces")
	}
	for _, v := range rscData.Interface {
		configSet = append(configSet, setPrefix+"interface "+v.ValueString())
	}
	if v := rscData.Location.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"location \""+v+"\"")
	}
	if rscData.RoutingInstanceAccess.ValueBool() {
		configSet = append(configSet, setPrefix+"routing-instance-access")
	}
	for _, v := range rscData.RoutingInstanceAccessList {
		configSet = append(configSet, setPrefix+"routing-instance-access access-list \""+v.ValueString()+"\"")
	}
	if rscData.HealthMonitor != nil {
		configSet = append(configSet, setPrefix+"health-monitor")

		if !rscData.HealthMonitor.FallingThreshold.IsNull() {
			configSet = append(configSet, setPrefix+"health-monitor falling-threshold "+
				utils.ConvI64toa(rscData.HealthMonitor.FallingThreshold.ValueInt64()))
		}
		if rscData.HealthMonitor.Idp.ValueBool() {
			configSet = append(configSet, setPrefix+"health-monitor idp")
		}
		if !rscData.HealthMonitor.IdpFallingThreshold.IsNull() {
			configSet = append(configSet, setPrefix+"health-monitor idp falling-threshold "+
				utils.ConvI64toa(rscData.HealthMonitor.IdpFallingThreshold.ValueInt64()))
		}
		if !rscData.HealthMonitor.IdpInterval.IsNull() {
			configSet = append(configSet, setPrefix+"health-monitor idp interval "+
				utils.ConvI64toa(rscData.HealthMonitor.IdpInterval.ValueInt64()))
		}
		if !rscData.HealthMonitor.IdpRisingThreshold.IsNull() {
			configSet = append(configSet, setPrefix+"health-monitor idp rising-threshold "+
				utils.ConvI64toa(rscData.HealthMonitor.IdpRisingThreshold.ValueInt64()))
		}
		if !rscData.HealthMonitor.Interval.IsNull() {
			configSet = append(configSet, setPrefix+"health-monitor interval "+
				utils.ConvI64toa(rscData.HealthMonitor.Interval.ValueInt64()))
		}
		if !rscData.HealthMonitor.RisingThreshold.IsNull() {
			configSet = append(configSet, setPrefix+"health-monitor rising-threshold "+
				utils.ConvI64toa(rscData.HealthMonitor.RisingThreshold.ValueInt64()))
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *snmpData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
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
			switch {
			case itemTrim == "arp":
				rscData.ARP = types.BoolValue(true)
			case itemTrim == "arp host-name-resolution":
				rscData.ARP = types.BoolValue(true)
				rscData.ARPHostNameResolution = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "contact "):
				rscData.Contact = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "engine-id "):
				rscData.EngineID = types.StringValue(itemTrim)
			case itemTrim == "filter-duplicates":
				rscData.FilterDuplicates = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "filter-interfaces interfaces "):
				rscData.FilterInterfaces = append(rscData.FilterInterfaces,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case itemTrim == "filter-interfaces all-internal-interfaces":
				rscData.FilterInternalInterfaces = types.BoolValue(true)
			case itemTrim == "if-count-with-filter-interfaces":
				rscData.IfCountWithFilterInterfaces = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "interface "):
				rscData.Interface = append(rscData.Interface, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "location "):
				rscData.Location = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "routing-instance-access":
				rscData.RoutingInstanceAccess = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "routing-instance-access access-list "):
				rscData.RoutingInstanceAccess = types.BoolValue(true)
				rscData.RoutingInstanceAccessList = append(rscData.RoutingInstanceAccessList,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "health-monitor"):
				if rscData.HealthMonitor == nil {
					rscData.HealthMonitor = &snmpBlockHealthMonitor{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " falling-threshold "):
					rscData.HealthMonitor.FallingThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case itemTrim == " idp":
					rscData.HealthMonitor.Idp = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, " idp falling-threshold "):
					rscData.HealthMonitor.Idp = types.BoolValue(true)
					rscData.HealthMonitor.IdpFallingThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " idp interval "):
					rscData.HealthMonitor.Idp = types.BoolValue(true)
					rscData.HealthMonitor.IdpInterval, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " idp rising-threshold "):
					rscData.HealthMonitor.Idp = types.BoolValue(true)
					rscData.HealthMonitor.IdpRisingThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " interval "):
					rscData.HealthMonitor.Interval, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " rising-threshold "):
					rscData.HealthMonitor.RisingThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (rscData *snmpData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := "delete snmp "

	configSet := []string{
		delPrefix + "arp",
		delPrefix + "contact",
		delPrefix + "description",
		delPrefix + "engine-id",
		delPrefix + "filter-duplicates",
		delPrefix + "filter-interfaces",
		delPrefix + "health-monitor",
		delPrefix + "if-count-with-filter-interfaces",
		delPrefix + "interface",
		delPrefix + "location",
		delPrefix + "routing-instance-access",
	}

	return junSess.ConfigSet(configSet)
}
