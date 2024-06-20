package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
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
	_ datasource.DataSource                   = &interfaceLogicalDataSource{}
	_ datasource.DataSourceWithConfigure      = &interfaceLogicalDataSource{}
	_ datasource.DataSourceWithValidateConfig = &interfaceLogicalDataSource{}
)

type interfaceLogicalDataSource struct {
	client *junos.Client
}

func (dsc *interfaceLogicalDataSource) typeName() string {
	return providerName + "_interface_logical"
}

func (dsc *interfaceLogicalDataSource) junosName() string {
	return "logical interface"
}

func newInterfaceLogicalDataSource() datasource.DataSource {
	return &interfaceLogicalDataSource{}
}

func (dsc *interfaceLogicalDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *interfaceLogicalDataSource) Configure(
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

func (dsc *interfaceLogicalDataSource) Schema(
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
				Description: "Name of logical interface (with dot).",
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
				Description: "Logical link-layer encapsulation.",
			},
			"routing_instance": schema.StringAttribute{
				Computed:    true,
				Description: "Routing_instance where the interface is.",
			},
			"security_inbound_protocols": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The inbound protocols allowed.",
			},
			"security_inbound_services": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The inbound services allowed.",
			},
			"security_zone": schema.StringAttribute{
				Computed:    true,
				Description: "Security zone where the interface is.",
			},
			"vlan_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Virtual LAN identifier value for 802.1q VLAN tags.",
			},
			"family_inet": schema.ObjectAttribute{
				Computed:    true,
				Description: "Family inet enabled and possible configuration.",
				AttributeTypes: map[string]attr.Type{
					"filter_input":    types.StringType,
					"filter_output":   types.StringType,
					"mtu":             types.Int64Type,
					"sampling_input":  types.BoolType,
					"sampling_output": types.BoolType,
					"address": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"cidr_ip":   types.StringType,
						"preferred": types.BoolType,
						"primary":   types.BoolType,
						"vrrp_group": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
							"identifier":               types.Int64Type,
							"virtual_address":          types.ListType{}.WithElementType(types.StringType),
							"accept_data":              types.BoolType,
							"no_accept_data":           types.BoolType,
							"advertise_interval":       types.Int64Type,
							"advertisements_threshold": types.Int64Type,
							"authentication_key":       types.StringType,
							"authentication_type":      types.StringType,
							"preempt":                  types.BoolType,
							"no_preempt":               types.BoolType,
							"priority":                 types.Int64Type,
							"track_interface": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
								"interface":     types.StringType,
								"priority_cost": types.StringType,
							})),
							"track_route": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
								"route":            types.StringType,
								"routing_instance": types.StringType,
								"priority_cost":    types.StringType,
							})),
						})),
					})),
					"dhcp": types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"srx_old_option_name":                            types.BoolType,
						"client_identifier_ascii":                        types.StringType,
						"client_identifier_hexadecimal":                  types.StringType,
						"client_identifier_prefix_hostname":              types.BoolType,
						"client_identifier_prefix_routing_instance_name": types.BoolType,
						"client_identifier_use_interface_description":    types.StringType,
						"client_identifier_userid_ascii":                 types.StringType,
						"client_identifier_userid_hexadecimal":           types.StringType,
						"force_discover":                                 types.BoolType,
						"lease_time":                                     types.Int64Type,
						"lease_time_infinite":                            types.BoolType,
						"metric":                                         types.Int64Type,
						"no_dns_install":                                 types.BoolType,
						"options_no_hostname":                            types.BoolType,
						"retransmission_attempt":                         types.Int64Type,
						"retransmission_interval":                        types.Int64Type,
						"server_address":                                 types.StringType,
						"update_server":                                  types.BoolType,
						"vendor_id":                                      types.StringType,
					}),
					"rpf_check": types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"fail_filter": types.StringType,
						"mode_loose":  types.BoolType,
					}),
				},
			},
			"family_inet6": schema.ObjectAttribute{
				Computed:    true,
				Description: "Family inet6 enabled and possible configuration.",
				AttributeTypes: map[string]attr.Type{
					"dad_disable":     types.BoolType,
					"filter_input":    types.StringType,
					"filter_output":   types.StringType,
					"mtu":             types.Int64Type,
					"sampling_input":  types.BoolType,
					"sampling_output": types.BoolType,
					"address": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"cidr_ip":   types.StringType,
						"preferred": types.BoolType,
						"primary":   types.BoolType,
						"vrrp_group": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
							"identifier":                 types.Int64Type,
							"virtual_address":            types.ListType{}.WithElementType(types.StringType),
							"virtual_link_local_address": types.StringType,
							"accept_data":                types.BoolType,
							"no_accept_data":             types.BoolType,
							"advertise_interval":         types.Int64Type,
							"advertisements_threshold":   types.Int64Type,
							"preempt":                    types.BoolType,
							"no_preempt":                 types.BoolType,
							"priority":                   types.Int64Type,
							"track_interface": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
								"interface":     types.StringType,
								"priority_cost": types.StringType,
							})),
							"track_route": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
								"route":            types.StringType,
								"routing_instance": types.StringType,
								"priority_cost":    types.StringType,
							})),
						})),
					})),
					"dhcpv6_client": types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"client_identifier_duid_type":               types.StringType,
						"client_type":                               types.StringType,
						"client_ia_type_na":                         types.BoolType,
						"client_ia_type_pd":                         types.BoolType,
						"no_dns_install":                            types.BoolType,
						"prefix_delegating_preferred_prefix_length": types.Int64Type,
						"prefix_delegating_sub_prefix_length":       types.Int64Type,
						"rapid_commit":                              types.BoolType,
						"req_option":                                types.ListType{}.WithElementType(types.StringType),
						"retransmission_attempt":                    types.Int64Type,
						"update_router_advertisement_interface":     types.SetType{}.WithElementType(types.StringType),
						"update_server":                             types.BoolType,
					}),
					"rpf_check": types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"fail_filter": types.StringType,
						"mode_loose":  types.BoolType,
					}),
				},
			},
			"tunnel": schema.ObjectAttribute{
				Computed:    true,
				Description: "Tunnel parameters.",
				AttributeTypes: map[string]attr.Type{
					"destination":                  types.StringType,
					"source":                       types.StringType,
					"allow_fragmentation":          types.BoolType,
					"do_not_fragment":              types.BoolType,
					"flow_label":                   types.Int64Type,
					"path_mtu_discovery":           types.BoolType,
					"no_path_mtu_discovery":        types.BoolType,
					"routing_instance_destination": types.StringType,
					"traffic_class":                types.Int64Type,
					"ttl":                          types.Int64Type,
				},
			},
		},
	}
}

type interfaceLogicalDataSourceData struct {
	ID                       types.String                      `tfsdk:"id"`
	ConfigInterface          types.String                      `tfsdk:"config_interface"`
	Match                    types.String                      `tfsdk:"match"`
	Name                     types.String                      `tfsdk:"name"`
	Description              types.String                      `tfsdk:"description"`
	Disable                  types.Bool                        `tfsdk:"disable"`
	Encapsulation            types.String                      `tfsdk:"encapsulation"`
	RoutingInstance          types.String                      `tfsdk:"routing_instance"`
	SecurityInboundProtocols []types.String                    `tfsdk:"security_inbound_protocols"`
	SecurityInboundServices  []types.String                    `tfsdk:"security_inbound_services"`
	SecurityZone             types.String                      `tfsdk:"security_zone"`
	VlanID                   types.Int64                       `tfsdk:"vlan_id"`
	FamilyInet               *interfaceLogicalBlockFamilyInet  `tfsdk:"family_inet"`
	FamilyInet6              *interfaceLogicalBlockFamilyInet6 `tfsdk:"family_inet6"`
	Tunnel                   *interfaceLogicalBlockTunnel      `tfsdk:"tunnel"`
}

func (dsc *interfaceLogicalDataSource) ValidateConfig(
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

func (dsc *interfaceLogicalDataSource) Read(
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
		resp.Diagnostics.AddError(tfdiag.NotFoundErrSummary, "no logical interface found with arguments provided")

		return
	}

	var rscData interfaceLogicalData
	if err := rscData.read(ctx, nameFound, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ReadErrSummary, err.Error())

		return
	}

	var data interfaceLogicalDataSourceData
	data.copyFromResourceData(rscData)

	data.ConfigInterface = configInterface
	data.Match = match
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (dsc *interfaceLogicalDataSource) searchName(
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
		switch len(itemTrimFields) {
		case 0, 1, 2:
			continue
		default:
			if itemTrimFields[1] == "unit" && !slices.Contains(itemTrimFields, "ethernet-switching") {
				intConfigList = append(intConfigList, itemTrimFields[0]+"."+itemTrimFields[2])
			}
		}
	}
	intConfigList = balt.UniqueInSlice(intConfigList)
	if len(intConfigList) == 0 {
		return "", nil
	}
	if len(intConfigList) > 1 {
		return "", errors.New("too many different logical interfaces found")
	}

	return intConfigList[0], nil
}

func (dscData *interfaceLogicalDataSourceData) copyFromResourceData(rscData interfaceLogicalData) {
	dscData.ID = rscData.ID
	dscData.Name = rscData.Name
	dscData.Description = rscData.Description
	dscData.Disable = rscData.Disable
	dscData.Encapsulation = rscData.Encapsulation
	dscData.FamilyInet = rscData.FamilyInet
	dscData.FamilyInet6 = rscData.FamilyInet6
	dscData.RoutingInstance = rscData.RoutingInstance
	dscData.SecurityInboundProtocols = rscData.SecurityInboundProtocols
	dscData.SecurityInboundServices = rscData.SecurityInboundServices
	dscData.SecurityZone = rscData.SecurityZone
	dscData.Tunnel = rscData.Tunnel
	dscData.VlanID = rscData.VlanID
}
