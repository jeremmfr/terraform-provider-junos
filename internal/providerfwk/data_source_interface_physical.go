package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource                   = &interfacePhysicalDataSource{}
	_ datasource.DataSourceWithConfigure      = &interfacePhysicalDataSource{}
	_ datasource.DataSourceWithValidateConfig = &interfacePhysicalDataSource{}
)

type interfacePhysicalDataSource struct {
	client *junos.Client
}

func (dsc *interfacePhysicalDataSource) typeName() string {
	return providerName + "_interface_physical"
}

func (dsc *interfacePhysicalDataSource) junosName() string {
	return "physical interface"
}

func newInterfacePhysicalDataSource() datasource.DataSource {
	return &interfacePhysicalDataSource{}
}

func (dsc *interfacePhysicalDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *interfacePhysicalDataSource) Configure(
	ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*junos.Client)
	if !ok {
		unexpectedDataSourceConfigureType(ctx, req, resp)

		return
	}
	dsc.client = client
}

func (dsc *interfacePhysicalDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get configuration from a " + dsc.junosName() + ".",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the data source with format `<name>`.",
			},
			"config_interface": schema.StringAttribute{
				Optional:    true,
				Description: "Specifies the interface part for search.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"match": schema.StringAttribute{
				Optional:    true,
				Description: "Regex string to filter lines and find only one interface.",
				Validators: []validator.String{
					tfvalidator.StringRegex(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of physical interface (without dot).",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description for interface.",
			},
			"disable": schema.BoolAttribute{
				Computed:    true,
				Description: "Interface disabled.",
			},
			"encapsulation": schema.StringAttribute{
				Computed:    true,
				Description: "Physical link-layer encapsulation.",
			},
			"flexible_vlan_tagging": schema.BoolAttribute{
				Computed:    true,
				Description: "Support for no tagging, or single and double 802.1q VLAN tagging.",
			},
			"gratuitous_arp_reply": schema.BoolAttribute{
				Computed:    true,
				Description: "Enable gratuitous ARP reply.",
			},
			"hold_time_down": schema.Int64Attribute{
				Computed:    true,
				Description: "Link down hold time (milliseconds).",
			},
			"hold_time_up": schema.Int64Attribute{
				Computed:    true,
				Description: "Link up hold time (milliseconds).",
			},
			"link_mode": schema.StringAttribute{
				Computed:    true,
				Description: "Link operational mode.",
			},
			"mtu": schema.Int64Attribute{
				Computed:    true,
				Description: "Maximum transmission unit.",
			},
			"no_gratuitous_arp_reply": schema.BoolAttribute{
				Computed:    true,
				Description: "Don't enable gratuitous ARP reply.",
			},
			"no_gratuitous_arp_request": schema.BoolAttribute{
				Computed:    true,
				Description: "Ignore gratuitous ARP request.",
			},
			"speed": schema.StringAttribute{
				Computed:    true,
				Description: "Link speed.",
			},
			"storm_control": schema.StringAttribute{
				Computed:    true,
				Description: "Storm control profile name to bind.",
			},
			"trunk": schema.BoolAttribute{
				Computed:    true,
				Description: "Interface mode is trunk.",
			},
			"trunk_non_els": schema.BoolAttribute{
				Computed:    true,
				Description: "Port mode is trunk.",
			},
			"vlan_members": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of vlan membership for this interface.",
			},
			"vlan_native": schema.Int64Attribute{
				Computed:    true,
				Description: "Vlan for untagged frames.",
			},
			"vlan_native_non_els": schema.StringAttribute{
				Computed:    true,
				Description: "Vlan for untagged frames (non-ELS).",
			},
			"vlan_tagging": schema.BoolAttribute{
				Computed:    true,
				Description: "802.1q VLAN tagging support.",
			},
			"esi": schema.ObjectAttribute{
				Computed:    true,
				Description: "ESI Config parameters.",
				AttributeTypes: map[string]attr.Type{
					"mode":             types.StringType,
					"auto_derive_lacp": types.BoolType,
					"df_election_type": types.StringType,
					"identifier":       types.StringType,
					"source_bmac":      types.StringType,
				},
			},
			"ether_opts": schema.ObjectAttribute{
				Computed:    true,
				Description: "The `ether-options` configuration.",
				AttributeTypes: map[string]attr.Type{
					"ae_8023ad":           types.StringType,
					"auto_negotiation":    types.BoolType,
					"no_auto_negotiation": types.BoolType,
					"flow_control":        types.BoolType,
					"no_flow_control":     types.BoolType,
					"loopback":            types.BoolType,
					"no_loopback":         types.BoolType,
					"redundant_parent":    types.StringType,
				},
			},
			"gigether_opts": schema.ObjectAttribute{
				Computed:    true,
				Description: "The `gigether-options` configuration.",
				AttributeTypes: map[string]attr.Type{
					"ae_8023ad":           types.StringType,
					"auto_negotiation":    types.BoolType,
					"no_auto_negotiation": types.BoolType,
					"flow_control":        types.BoolType,
					"no_flow_control":     types.BoolType,
					"loopback":            types.BoolType,
					"no_loopback":         types.BoolType,
					"redundant_parent":    types.StringType,
				},
			},
			"parent_ether_opts": schema.ObjectAttribute{
				Computed:    true,
				Description: "The `aggregated-ether-options` or `redundant-ether-options` configuration.",
				AttributeTypes: map[string]attr.Type{
					"flow_control":          types.BoolType,
					"no_flow_control":       types.BoolType,
					"loopback":              types.BoolType,
					"no_loopback":           types.BoolType,
					"link_speed":            types.StringType,
					"minimum_bandwidth":     types.StringType,
					"minimum_links":         types.Int64Type,
					"redundancy_group":      types.Int64Type,
					"source_address_filter": types.ListType{}.WithElementType(types.StringType),
					"source_filtering":      types.BoolType,
					"bfd_liveness_detection": types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"local_address":                      types.StringType,
						"authentication_algorithm":           types.StringType,
						"authentication_key_chain":           types.StringType,
						"authentication_loose_check":         types.BoolType,
						"detection_time_threshold":           types.Int64Type,
						"holddown_interval":                  types.Int64Type,
						"minimum_interval":                   types.Int64Type,
						"minimum_receive_interval":           types.Int64Type,
						"multiplier":                         types.Int64Type,
						"neighbor":                           types.StringType,
						"no_adaptation":                      types.BoolType,
						"transmit_interval_minimum_interval": types.Int64Type,
						"transmit_interval_threshold":        types.Int64Type,
						"version":                            types.StringType,
					}),
					"lacp": types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"mode":            types.StringType,
						"admin_key":       types.Int64Type,
						"periodic":        types.StringType,
						"sync_reset":      types.StringType,
						"system_id":       types.StringType,
						"system_priority": types.Int64Type,
					}),
					"mc_ae": types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"chassis_id":           types.Int64Type,
						"mc_ae_id":             types.Int64Type,
						"mode":                 types.StringType,
						"status_control":       types.StringType,
						"enhanced_convergence": types.BoolType,
						"init_delay_time":      types.Int64Type,
						"redundancy_group":     types.Int64Type,
						"revert_time":          types.Int64Type,
						"switchover_mode":      types.StringType,
						"events_iccp_peer_down": types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
							"force_icl_down":               types.BoolType,
							"prefer_status_control_active": types.BoolType,
						}),
					}),
				},
			},
		},
	}
}

type interfacePhysicalDataSourceData struct {
	ID                     types.String                           `tfsdk:"id"`
	ConfigInterface        types.String                           `tfsdk:"config_interface"`
	Match                  types.String                           `tfsdk:"match"`
	Name                   types.String                           `tfsdk:"name"`
	Description            types.String                           `tfsdk:"description"`
	Disable                types.Bool                             `tfsdk:"disable"`
	Encapsulation          types.String                           `tfsdk:"encapsulation"`
	FlexibleVlanTagging    types.Bool                             `tfsdk:"flexible_vlan_tagging"`
	GratuitousArpReply     types.Bool                             `tfsdk:"gratuitous_arp_reply"`
	NoGratuitousArpReply   types.Bool                             `tfsdk:"no_gratuitous_arp_reply"`
	HoldTimeDown           types.Int64                            `tfsdk:"hold_time_down"`
	HoldTimeUp             types.Int64                            `tfsdk:"hold_time_up"`
	LinkMode               types.String                           `tfsdk:"link_mode"`
	Mtu                    types.Int64                            `tfsdk:"mtu"`
	NoGratuitousArpRequest types.Bool                             `tfsdk:"no_gratuitous_arp_request"`
	Speed                  types.String                           `tfsdk:"speed"`
	StormControl           types.String                           `tfsdk:"storm_control"`
	Trunk                  types.Bool                             `tfsdk:"trunk"`
	TrunkNonELS            types.Bool                             `tfsdk:"trunk_non_els"`
	VlanMembers            []types.String                         `tfsdk:"vlan_members"`
	VlanNative             types.Int64                            `tfsdk:"vlan_native"`
	VlanNativeNonELS       types.String                           `tfsdk:"vlan_native_non_els"`
	VlanTagging            types.Bool                             `tfsdk:"vlan_tagging"`
	ESI                    *interfacePhysicalBlockESI             `tfsdk:"esi"`
	EtherOpts              *interfacePhysicalBlockEtherOpts       `tfsdk:"ether_opts"`
	GigetherOpts           *interfacePhysicalBlockEtherOpts       `tfsdk:"gigether_opts"`
	ParentEtherOpts        *interfacePhysicalBlockParentEtherOpts `tfsdk:"parent_ether_opts"`
}

func (dsc *interfacePhysicalDataSource) ValidateConfig(
	ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse,
) {
	var configInterface, match types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("config_interface"), &configInterface)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("match"), &match)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if configInterface.IsNull() && match.IsNull() {
		resp.Diagnostics.AddError(
			tfdiag.MissingConfigErrSummary,
			"one of config_interface or match must be specified",
		)
	}
}

func (dsc *interfacePhysicalDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var configInterface, match types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("config_interface"), &configInterface)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("match"), &match)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if configInterface.ValueString() == "" && match.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Empty Argument",
			"could not search "+dsc.junosName()+" with empty config_interface and match",
		)

		return
	}

	junSess, err := dsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	defer junos.MutexUnlock()

	nameFound, err := dsc.searchName(
		ctx,
		configInterface.ValueString(),
		match.ValueString(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ReadErrSummary, err.Error())

		return
	}
	if nameFound == "" {
		resp.Diagnostics.AddError(tfdiag.NotFoundErrSummary, "no physical interface found with arguments provided")

		return
	}

	var rscData interfacePhysicalData
	if err := rscData.read(ctx, nameFound, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ReadErrSummary, err.Error())

		return
	}

	var data interfacePhysicalDataSourceData
	data.copyFromResourceData(rscData)

	data.ConfigInterface = configInterface
	data.Match = match
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (dsc *interfacePhysicalDataSource) searchName(
	_ context.Context, configInterface, match string, junSess *junos.Session,
) (string, error) {
	intConfigList := make([]string, 0)
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"interfaces " + configInterface + junos.PipeDisplaySet)
	if err != nil {
		return "", err
	}
	for _, item := range strings.Split(showConfig, "\n") {
		if strings.Contains(item, junos.XMLStartTagConfigOut) {
			continue
		}
		if strings.Contains(item, junos.XMLEndTagConfigOut) {
			break
		}
		if item == "" {
			continue
		}
		itemTrim := strings.TrimPrefix(item, "set interfaces ")
		matched, err := regexp.MatchString(match, itemTrim)
		if err != nil {
			return "", fmt.Errorf("matching with regexp %q: %w", match, err)
		}
		if !matched {
			continue
		}
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) == 0 {
			continue
		}
		intConfigList = append(intConfigList, itemTrimFields[0])
	}
	intConfigList = balt.UniqueInSlice(intConfigList)
	if len(intConfigList) == 0 {
		return "", nil
	}
	if len(intConfigList) > 1 {
		return "", errors.New("too many different physical interfaces found")
	}

	return intConfigList[0], nil
}

func (dscData *interfacePhysicalDataSourceData) copyFromResourceData(rscData interfacePhysicalData) {
	dscData.ID = rscData.ID
	dscData.Name = rscData.Name
	dscData.Description = rscData.Description
	dscData.Disable = rscData.Disable
	dscData.Encapsulation = rscData.Encapsulation
	dscData.ESI = rscData.ESI
	dscData.EtherOpts = rscData.EtherOpts
	dscData.FlexibleVlanTagging = rscData.FlexibleVlanTagging
	dscData.GigetherOpts = rscData.GigetherOpts
	dscData.GratuitousArpReply = rscData.GratuitousArpReply
	dscData.HoldTimeDown = rscData.HoldTimeDown
	dscData.HoldTimeUp = rscData.HoldTimeUp
	dscData.LinkMode = rscData.LinkMode
	dscData.Mtu = rscData.Mtu
	dscData.NoGratuitousArpReply = rscData.NoGratuitousArpReply
	dscData.NoGratuitousArpRequest = rscData.NoGratuitousArpRequest
	dscData.ParentEtherOpts = rscData.ParentEtherOpts
	dscData.Speed = rscData.Speed
	dscData.StormControl = rscData.StormControl
	dscData.Trunk = rscData.Trunk
	dscData.TrunkNonELS = rscData.TrunkNonELS
	dscData.VlanMembers = rscData.VlanMembers
	dscData.VlanNative = rscData.VlanNative
	dscData.VlanNativeNonELS = rscData.VlanNativeNonELS
	dscData.VlanTagging = rscData.VlanTagging
}
