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
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &interfaceLogical{}
	_ resource.ResourceWithConfigure      = &interfaceLogical{}
	_ resource.ResourceWithModifyPlan     = &interfaceLogical{}
	_ resource.ResourceWithValidateConfig = &interfaceLogical{}
	_ resource.ResourceWithImportState    = &interfaceLogical{}
	_ resource.ResourceWithUpgradeState   = &interfaceLogical{}
)

type interfaceLogical struct {
	client *junos.Client
}

func newInterfaceLogicalResource() resource.Resource {
	return &interfaceLogical{}
}

func (rsc *interfaceLogical) typeName() string {
	return providerName + "_interface_logical"
}

func (rsc *interfaceLogical) junosName() string {
	return "logical interface"
}

func (rsc *interfaceLogical) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *interfaceLogical) Configure(
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

func (rsc *interfaceLogical) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
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
				Description: "Name of logical interface (with dot).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					tfvalidator.String1DotCount(),
				},
			},
			"st0_also_on_destroy": schema.BoolAttribute{
				Optional: true,
				Description: "When destroy this resource, if the name has prefix `st0.`, " +
					"delete all configurations (not keep empty st0 interface).",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description for interface.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable this logical interface.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Description: "Add this interface in routing_instance.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					stringvalidator.NoneOfCaseInsensitive(junos.DefaultW),
				},
			},
			"security_inbound_protocols": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The inbound protocols allowed.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
			"security_inbound_services": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The inbound services allowed.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
			"security_zone": schema.StringAttribute{
				Optional:    true,
				Description: "Add this interface in a security zone.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"vlan_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Virtual LAN identifier value for 802.1q VLAN tags.",
				Validators: []validator.Int64{
					int64validator.Between(1, 4094),
				},
			},
			"vlan_no_compute": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable the automatic compute of the `vlan_id` argument when not set.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"family_inet": schema.SingleNestedBlock{
				Description: "Enable family inet and add configurations if specified.",
				Attributes: map[string]schema.Attribute{
					"filter_input": schema.StringAttribute{
						Optional:    true,
						Description: "Filter to be applied to received packets.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"filter_output": schema.StringAttribute{
						Optional:    true,
						Description: "Filter to be applied to transmitted packets.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"mtu": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum transmission unit.",
						Validators: []validator.Int64{
							int64validator.Between(1, 9500),
						},
					},
					"sampling_input": schema.BoolAttribute{
						Optional:    true,
						Description: "Sample all packets input on this interface.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"sampling_output": schema.BoolAttribute{
						Optional:    true,
						Description: "Sample all packets output on this interface.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"address": schema.ListNestedBlock{
						Description: "For each IPv4 address to declare.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"cidr_ip": schema.StringAttribute{
									Required:    true,
									Description: "IPv4 address in CIDR format.",
									Validators: []validator.String{
										tfvalidator.StringCIDR().IPv4Only(),
									},
								},
								"preferred": schema.BoolAttribute{
									Optional:    true,
									Description: "Preferred address on interface.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"primary": schema.BoolAttribute{
									Optional:    true,
									Description: "Candidate for primary address in system.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
							},
							Blocks: map[string]schema.Block{
								"vrrp_group": schema.ListNestedBlock{
									Description: "For each vrrp group to declare.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"identifier": schema.Int64Attribute{
												Required:    true,
												Description: "ID for vrrp.",
												Validators: []validator.Int64{
													int64validator.Between(1, 255),
												},
											},
											"virtual_address": schema.ListAttribute{
												ElementType: types.StringType,
												Required:    true,
												Description: "Virtual IPv4 addresses.",
												Validators: []validator.List{
													listvalidator.SizeAtLeast(1),
													listvalidator.ValueStringsAre(
														tfvalidator.StringIPAddress().IPv4Only(),
													),
												},
											},
											"accept_data": schema.BoolAttribute{
												Optional:    true,
												Description: "Accept packets destined for virtual IP address.",
												Validators: []validator.Bool{
													tfvalidator.BoolTrue(),
												},
											},
											"no_accept_data": schema.BoolAttribute{
												Optional:    true,
												Description: "Don't accept packets destined for virtual IP address.",
												Validators: []validator.Bool{
													tfvalidator.BoolTrue(),
												},
											},
											"advertise_interval": schema.Int64Attribute{
												Optional:    true,
												Description: "Advertisement interval (seconds).",
												Validators: []validator.Int64{
													int64validator.Between(1, 255),
												},
											},
											"advertisements_threshold": schema.Int64Attribute{
												Optional:    true,
												Description: "Number of vrrp advertisements missed before declaring master down.",
												Validators: []validator.Int64{
													int64validator.Between(1, 15),
												},
											},
											"authentication_key": schema.StringAttribute{
												Optional:    true,
												Sensitive:   true,
												Description: "Authentication key.",
												Validators: []validator.String{
													stringvalidator.LengthBetween(1, 16),
													tfvalidator.StringDoubleQuoteExclusion(),
												},
											},
											"authentication_type": schema.StringAttribute{
												Optional:    true,
												Description: "Authentication type.",
												Validators: []validator.String{
													stringvalidator.OneOf("md5", "simple"),
												},
											},
											"preempt": schema.BoolAttribute{
												Optional:    true,
												Description: "Allow preemption.",
												Validators: []validator.Bool{
													tfvalidator.BoolTrue(),
												},
											},
											"no_preempt": schema.BoolAttribute{
												Optional:    true,
												Description: "Don't allow preemption.",
												Validators: []validator.Bool{
													tfvalidator.BoolTrue(),
												},
											},
											"priority": schema.Int64Attribute{
												Optional:    true,
												Description: "Virtual router election priority.",
												Validators: []validator.Int64{
													int64validator.Between(1, 255),
												},
											},
										},
										Blocks: map[string]schema.Block{
											"track_interface": schema.ListNestedBlock{
												Description: "For each interface to track in VRRP group.",
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"interface": schema.StringAttribute{
															Required:    true,
															Description: "Name of interface.",
															Validators: []validator.String{
																stringvalidator.LengthAtLeast(1),
																tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
															},
														},
														"priority_cost": schema.Int64Attribute{
															Required:    true,
															Description: "Value to subtract from priority when interface is down.",
															Validators: []validator.Int64{
																int64validator.Between(1, 254),
															},
														},
													},
												},
											},
											"track_route": schema.ListNestedBlock{
												Description: "For each route to track in VRRP group.",
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"route": schema.StringAttribute{
															Required:    true,
															Description: "Route address.",
															Validators: []validator.String{
																tfvalidator.StringCIDR(),
															},
														},
														"routing_instance": schema.StringAttribute{
															Required:    true,
															Description: "Routing instance to which route belongs, or `default`.",
															Validators: []validator.String{
																stringvalidator.LengthBetween(1, 63),
																tfvalidator.StringFormat(tfvalidator.DefaultFormat),
															},
														},
														"priority_cost": schema.Int64Attribute{
															Required:    true,
															Description: "Value to subtract from priority when route is down.",
															Validators: []validator.Int64{
																int64validator.Between(1, 254),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"dhcp": schema.SingleNestedBlock{
						Description: "Enable DHCP client and configuration.",
						Attributes: map[string]schema.Attribute{
							"srx_old_option_name": schema.BoolAttribute{
								Optional:    true,
								Description: "For configuration, use the old option name `dhcp-client` instead of `dhcp`.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"client_identifier_ascii": schema.StringAttribute{
								Optional:    true,
								Description: "Client identifier as an ASCII string.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"client_identifier_hexadecimal": schema.StringAttribute{
								Optional:    true,
								Description: "Client identifier as a hexadecimal string.",
								Validators: []validator.String{
									stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]+$`),
										"must be hexadecimal digits (0-9, a-f, A-F)"),
								},
							},
							"client_identifier_prefix_hostname": schema.BoolAttribute{
								Optional:    true,
								Description: "Add prefix router host name to client-id option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"client_identifier_prefix_routing_instance_name": schema.BoolAttribute{
								Optional:    true,
								Description: "Add prefix routing instance name to client-id option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"client_identifier_use_interface_description": schema.StringAttribute{
								Optional:    true,
								Description: "Use the interface description.",
								Validators: []validator.String{
									stringvalidator.OneOf("device", "logical"),
								},
							},
							"client_identifier_userid_ascii": schema.StringAttribute{
								Optional:    true,
								Description: "Add user id as an ASCII string to client-id option.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"client_identifier_userid_hexadecimal": schema.StringAttribute{
								Optional:    true,
								Description: "Add user id as a hexadecimal string to client-id option.",
								Validators: []validator.String{
									stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]+$`),
										"must be hexadecimal digits (0-9, a-f, A-F)"),
								},
							},
							"force_discover": schema.BoolAttribute{
								Optional:    true,
								Description: "Send DHCPDISCOVER after DHCPREQUEST retransmission failure.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"lease_time": schema.Int64Attribute{
								Optional:    true,
								Description: "Lease time in seconds requested in DHCP client protocol packet.",
								Validators: []validator.Int64{
									int64validator.Between(60, 2147483647),
								},
							},
							"lease_time_infinite": schema.BoolAttribute{
								Optional:    true,
								Description: "Lease never expires.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"metric": schema.Int64Attribute{
								Optional:    true,
								Description: "Client initiated default-route metric.",
								Validators: []validator.Int64{
									int64validator.Between(0, 255),
								},
							},
							"no_dns_install": schema.BoolAttribute{
								Optional:    true,
								Description: "Do not install DNS information learned from DHCP server.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"options_no_hostname": schema.BoolAttribute{
								Optional:    true,
								Description: "Do not carry hostname (RFC option code is 12) in packet.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"retransmission_attempt": schema.Int64Attribute{
								Optional:    true,
								Description: "Number of attempts to retransmit the DHCP client protocol packet.",
								Validators: []validator.Int64{
									int64validator.Between(0, 50000),
								},
							},
							"retransmission_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Number of seconds between successive retransmission.",
								Validators: []validator.Int64{
									int64validator.Between(4, 64),
								},
							},
							"server_address": schema.StringAttribute{
								Optional:    true,
								Description: "DHCP Server-address.",
								Validators: []validator.String{
									tfvalidator.StringIPAddress().IPv4Only(),
								},
							},
							"update_server": schema.BoolAttribute{
								Optional:    true,
								Description: "Propagate TCP/IP settings to DHCP server.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"vendor_id": schema.StringAttribute{
								Optional:    true,
								Description: "Vendor class id for the DHCP Client.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 60),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
						},
					},
					"rpf_check": schema.SingleNestedBlock{
						Description: "Enable reverse-path-forwarding checks on this interface.",
						Attributes: map[string]schema.Attribute{
							"fail_filter": schema.StringAttribute{
								Optional:    true,
								Description: "Name of filter applied to packets failing RPF check.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"mode_loose": schema.BoolAttribute{
								Optional:    true,
								Description: "Use reverse-path-forwarding loose mode instead the strict mode.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
					},
				},
			},
			"family_inet6": schema.SingleNestedBlock{
				Description: "Enable family inet6 and add configurations if specified.",
				Attributes: map[string]schema.Attribute{
					"dad_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable duplicate-address-detection.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"filter_input": schema.StringAttribute{
						Optional:    true,
						Description: "Filter to be applied to received packets.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"filter_output": schema.StringAttribute{
						Optional:    true,
						Description: "Filter to be applied to transmitted packets.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"mtu": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum transmission unit.",
						Validators: []validator.Int64{
							int64validator.Between(1, 9500),
						},
					},
					"sampling_input": schema.BoolAttribute{
						Optional:    true,
						Description: "Sample all packets input on this interface.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"sampling_output": schema.BoolAttribute{
						Optional:    true,
						Description: "Sample all packets output on this interface.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"address": schema.ListNestedBlock{
						Description: " For each IPv6 address to declare.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"cidr_ip": schema.StringAttribute{
									Required:    true,
									Description: "IPv6 address in CIDR format.",
									Validators: []validator.String{
										tfvalidator.StringCIDR().IPv6Only(),
									},
								},
								"preferred": schema.BoolAttribute{
									Optional:    true,
									Description: "Preferred address on interface.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"primary": schema.BoolAttribute{
									Optional:    true,
									Description: "Candidate for primary address in system.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
							},
							Blocks: map[string]schema.Block{
								"vrrp_group": schema.ListNestedBlock{
									Description: "For each vrrp group to declare.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"identifier": schema.Int64Attribute{
												Required:    true,
												Description: "ID for vrrp.",
												Validators: []validator.Int64{
													int64validator.Between(1, 255),
												},
											},
											"virtual_address": schema.ListAttribute{
												ElementType: types.StringType,
												Required:    true,
												Description: "Virtual IPv6 addresses.",
												Validators: []validator.List{
													listvalidator.SizeAtLeast(1),
													listvalidator.ValueStringsAre(
														tfvalidator.StringIPAddress().IPv6Only(),
													),
												},
											},
											"virtual_link_local_address": schema.StringAttribute{
												Required:    true,
												Description: "Address IPv6 for Virtual link-local addresses.",
												Validators: []validator.String{
													tfvalidator.StringIPAddress().IPv6Only(),
												},
											},
											"accept_data": schema.BoolAttribute{
												Optional:    true,
												Description: "Accept packets destined for virtual IP address.",
												Validators: []validator.Bool{
													tfvalidator.BoolTrue(),
												},
											},
											"no_accept_data": schema.BoolAttribute{
												Optional:    true,
												Description: "Don't accept packets destined for virtual IP address.",
												Validators: []validator.Bool{
													tfvalidator.BoolTrue(),
												},
											},
											"advertise_interval": schema.Int64Attribute{
												Optional:    true,
												Description: "Advertisement interval (seconds).",
												Validators: []validator.Int64{
													int64validator.Between(100, 40000),
												},
											},
											"advertisements_threshold": schema.Int64Attribute{
												Optional:    true,
												Description: "Number of vrrp advertisements missed before declaring master down.",
												Validators: []validator.Int64{
													int64validator.Between(1, 15),
												},
											},
											"preempt": schema.BoolAttribute{
												Optional:    true,
												Description: "Allow preemption.",
												Validators: []validator.Bool{
													tfvalidator.BoolTrue(),
												},
											},
											"no_preempt": schema.BoolAttribute{
												Optional:    true,
												Description: "Don't allow preemption.",
												Validators: []validator.Bool{
													tfvalidator.BoolTrue(),
												},
											},
											"priority": schema.Int64Attribute{
												Optional:    true,
												Description: "Virtual router election priority.",
												Validators: []validator.Int64{
													int64validator.Between(1, 255),
												},
											},
										},
										Blocks: map[string]schema.Block{
											"track_interface": schema.ListNestedBlock{
												Description: "For each interface to track in VRRP group.",
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"interface": schema.StringAttribute{
															Required:    true,
															Description: "Name of interface.",
															Validators: []validator.String{
																stringvalidator.LengthAtLeast(1),
																tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
															},
														},
														"priority_cost": schema.Int64Attribute{
															Required:    true,
															Description: "Value to subtract from priority when interface is down.",
															Validators: []validator.Int64{
																int64validator.Between(1, 254),
															},
														},
													},
												},
											},
											"track_route": schema.ListNestedBlock{
												Description: "For each route to track in VRRP group.",
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"route": schema.StringAttribute{
															Required:    true,
															Description: "Route address.",
															Validators: []validator.String{
																tfvalidator.StringCIDR(),
															},
														},
														"routing_instance": schema.StringAttribute{
															Required:    true,
															Description: "Routing instance to which route belongs, or `default`.",
															Validators: []validator.String{
																stringvalidator.LengthBetween(1, 63),
																tfvalidator.StringFormat(tfvalidator.DefaultFormat),
															},
														},
														"priority_cost": schema.Int64Attribute{
															Required:    true,
															Description: "Value to subtract from priority when route is down.",
															Validators: []validator.Int64{
																int64validator.Between(1, 254),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"dhcpv6_client": schema.SingleNestedBlock{
						Description: "Enable DHCP client and configuration.",
						Attributes: map[string]schema.Attribute{
							"client_identifier_duid_type": schema.StringAttribute{
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "DUID identifying a client.",
								Validators: []validator.String{
									stringvalidator.OneOf("duid-ll", "duid-llt", "vendor"),
								},
							},
							"client_type": schema.StringAttribute{
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "DHCPv6 client type.",
								Validators: []validator.String{
									stringvalidator.OneOf("autoconfig", "stateful"),
								},
							},
							"client_ia_type_na": schema.BoolAttribute{
								Optional:    true,
								Description: "DHCPv6 client identity association type Non-temporary Address.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"client_ia_type_pd": schema.BoolAttribute{
								Optional:    true,
								Description: "DHCPv6 client identity association type Prefix Address.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_dns_install": schema.BoolAttribute{
								Optional:    true,
								Description: "Do not install DNS information learned from DHCP server.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"prefix_delegating_preferred_prefix_length": schema.Int64Attribute{
								Optional:    true,
								Description: "Client preferred prefix length.",
								Validators: []validator.Int64{
									int64validator.Between(0, 64),
								},
							},
							"prefix_delegating_sub_prefix_length": schema.Int64Attribute{
								Optional:    true,
								Description: "The sub prefix length for LAN interfaces.",
								Validators: []validator.Int64{
									int64validator.Between(1, 127),
								},
							},
							"rapid_commit": schema.BoolAttribute{
								Optional:    true,
								Description: "Option is used to signal the use of the two message exchange for address assignment.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"req_option": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "DHCPV6 client requested option configuration.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									),
								},
							},
							"retransmission_attempt": schema.Int64Attribute{
								Optional:    true,
								Description: "Number of attempts to retransmit the DHCPV6 client protocol packet.",
								Validators: []validator.Int64{
									int64validator.Between(0, 9),
								},
							},
							"update_router_advertisement_interface": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Interfaces on which to delegate prefix.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
										tfvalidator.String1DotCount(),
									),
								},
							},
							"update_server": schema.BoolAttribute{
								Optional:    true,
								Description: "Propagate TCP/IP settings to DHCP server.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
					},
					"rpf_check": schema.SingleNestedBlock{
						Description: "Enable reverse-path-forwarding checks on this interface.",
						Attributes: map[string]schema.Attribute{
							"fail_filter": schema.StringAttribute{
								Optional:    true,
								Description: "Name of filter applied to packets failing RPF check.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"mode_loose": schema.BoolAttribute{
								Optional:    true,
								Description: "Use reverse-path-forwarding loose mode instead the strict mode.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
					},
				},
			},
			"tunnel": schema.SingleNestedBlock{
				Description: "Tunnel parameters.",
				Attributes: map[string]schema.Attribute{
					"destination": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Tunnel destination.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"source": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Tunnel source.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"allow_fragmentation": schema.BoolAttribute{
						Optional:    true,
						Description: "Do not set DF bit on packets.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"do_not_fragment": schema.BoolAttribute{
						Optional:    true,
						Description: "Set DF bit on packets.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"flow_label": schema.Int64Attribute{
						Optional:    true,
						Description: "Flow label field of IP6-header.",
						Validators: []validator.Int64{
							int64validator.Between(0, 1048575),
						},
					},
					"path_mtu_discovery": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable path MTU discovery for tunnels.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_path_mtu_discovery": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't enable path MTU discovery for tunnels.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"routing_instance_destination": schema.StringAttribute{
						Optional:    true,
						Description: "Routing instance to which tunnel ends belong.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 63),
							tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							stringvalidator.NoneOfCaseInsensitive(junos.DefaultW),
						},
					},
					"traffic_class": schema.Int64Attribute{
						Optional:    true,
						Description: "TOS/Traffic class field of IP-header",
						Validators: []validator.Int64{
							int64validator.Between(0, 255),
						},
					},
					"ttl": schema.Int64Attribute{
						Optional:    true,
						Description: "Time to live",
						Validators: []validator.Int64{
							int64validator.Between(1, 255),
						},
					},
				},
			},
		},
	}
}

type interfaceLogicalData struct {
	St0AlsoOnDestroy         types.Bool                        `tfsdk:"st0_also_on_destroy"`
	VlanNoCompute            types.Bool                        `tfsdk:"vlan_no_compute"`
	Disable                  types.Bool                        `tfsdk:"disable"`
	ID                       types.String                      `tfsdk:"id"`
	Name                     types.String                      `tfsdk:"name"`
	Description              types.String                      `tfsdk:"description"`
	RoutingInstance          types.String                      `tfsdk:"routing_instance"`
	SecurityInboundProtocols []types.String                    `tfsdk:"security_inbound_protocols"`
	SecurityInboundServices  []types.String                    `tfsdk:"security_inbound_services"`
	SecurityZone             types.String                      `tfsdk:"security_zone"`
	VlanID                   types.Int64                       `tfsdk:"vlan_id"`
	FamilyInet               *interfaceLogicalBlockFamilyInet  `tfsdk:"family_inet"`
	FamilyInet6              *interfaceLogicalBlockFamilyInet6 `tfsdk:"family_inet6"`
	Tunnel                   *interfaceLogicalBlockTunnel      `tfsdk:"tunnel"`
}

type interfaceLogicalConfig struct {
	St0AlsoOnDestroy         types.Bool                              `tfsdk:"st0_also_on_destroy"`
	VlanNoCompute            types.Bool                              `tfsdk:"vlan_no_compute"`
	Disable                  types.Bool                              `tfsdk:"disable"`
	ID                       types.String                            `tfsdk:"id"`
	Name                     types.String                            `tfsdk:"name"`
	Description              types.String                            `tfsdk:"description"`
	RoutingInstance          types.String                            `tfsdk:"routing_instance"`
	SecurityInboundProtocols types.Set                               `tfsdk:"security_inbound_protocols"`
	SecurityInboundServices  types.Set                               `tfsdk:"security_inbound_services"`
	SecurityZone             types.String                            `tfsdk:"security_zone"`
	VlanID                   types.Int64                             `tfsdk:"vlan_id"`
	FamilyInet               *interfaceLogicalBlockFamilyInetConfig  `tfsdk:"family_inet"`
	FamilyInet6              *interfaceLogicalBlockFamilyInet6Config `tfsdk:"family_inet6"`
	Tunnel                   *interfaceLogicalBlockTunnel            `tfsdk:"tunnel"`
}

type interfaceLogicalBlockFamilyInet struct {
	SamplingInput  types.Bool                                    `tfsdk:"sampling_input"`
	SamplingOutput types.Bool                                    `tfsdk:"sampling_output"`
	FilterInput    types.String                                  `tfsdk:"filter_input"`
	FilterOutput   types.String                                  `tfsdk:"filter_output"`
	Mtu            types.Int64                                   `tfsdk:"mtu"`
	Address        []interfaceLogicalBlockFamilyInetBlockAddress `tfsdk:"address"`
	DHCP           *interfaceLogicalBlockFamilyInetBlockDhcp     `tfsdk:"dhcp"`
	RPFCheck       *interfaceLogicalBlockFamilyBlockRPFCheck     `tfsdk:"rpf_check"`
}

type interfaceLogicalBlockFamilyInetConfig struct {
	SamplingInput  types.Bool                                `tfsdk:"sampling_input"`
	SamplingOutput types.Bool                                `tfsdk:"sampling_output"`
	FilterInput    types.String                              `tfsdk:"filter_input"`
	FilterOutput   types.String                              `tfsdk:"filter_output"`
	Mtu            types.Int64                               `tfsdk:"mtu"`
	Address        types.List                                `tfsdk:"address"`
	DHCP           *interfaceLogicalBlockFamilyInetBlockDhcp `tfsdk:"dhcp"`
	RPFCheck       *interfaceLogicalBlockFamilyBlockRPFCheck `tfsdk:"rpf_check"`
}

type interfaceLogicalBlockFamilyBlockRPFCheck struct {
	FailFilter types.String `tfsdk:"fail_filter"`
	ModeLoose  types.Bool   `tfsdk:"mode_loose"`
}

type interfaceLogicalBlockFamilyInetBlockAddress struct {
	Preferred types.Bool                                                  `tfsdk:"preferred"`
	Primary   types.Bool                                                  `tfsdk:"primary"`
	CidrIP    types.String                                                `tfsdk:"cidr_ip"`
	VRRPGroup []interfaceLogicalBlockFamilyInetBlockAddressBlockVRRPGroup `tfsdk:"vrrp_group"`
}

type interfaceLogicalBlockFamilyInetBlockAddressConfig struct {
	Preferred types.Bool   `tfsdk:"preferred"`
	Primary   types.Bool   `tfsdk:"primary"`
	CidrIP    types.String `tfsdk:"cidr_ip"`
	VRRPGroup types.List   `tfsdk:"vrrp_group"`
}

//nolint:lll
type interfaceLogicalBlockFamilyInetBlockAddressBlockVRRPGroup struct {
	AcceptData              types.Bool                                                                 `tfsdk:"accept_data"`
	NoAcceptData            types.Bool                                                                 `tfsdk:"no_accept_data"`
	Preempt                 types.Bool                                                                 `tfsdk:"preempt"`
	NoPreempt               types.Bool                                                                 `tfsdk:"no_preempt"`
	Identifier              types.Int64                                                                `tfsdk:"identifier"`
	VirtualAddress          []types.String                                                             `tfsdk:"virtual_address"`
	AdvertiseInterval       types.Int64                                                                `tfsdk:"advertise_interval"`
	AdvertisementsThreshold types.Int64                                                                `tfsdk:"advertisements_threshold"`
	AuthenticationKey       types.String                                                               `tfsdk:"authentication_key"`
	AuthenticationType      types.String                                                               `tfsdk:"authentication_type"`
	Priority                types.Int64                                                                `tfsdk:"priority"`
	TrackInterface          []interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackInterface `tfsdk:"track_interface"`
	TrackRoute              []interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackRoute     `tfsdk:"track_route"`
}

type interfaceLogicalBlockFamilyInetBlockAddressBlockVRRPGroupConfig struct {
	AcceptData              types.Bool   `tfsdk:"accept_data"`
	NoAcceptData            types.Bool   `tfsdk:"no_accept_data"`
	Preempt                 types.Bool   `tfsdk:"preempt"`
	NoPreempt               types.Bool   `tfsdk:"no_preempt"`
	Identifier              types.Int64  `tfsdk:"identifier"`
	VirtualAddress          types.List   `tfsdk:"virtual_address"`
	AdvertiseInterval       types.Int64  `tfsdk:"advertise_interval"`
	AdvertisementsThreshold types.Int64  `tfsdk:"advertisements_threshold"`
	AuthenticationKey       types.String `tfsdk:"authentication_key"`
	AuthenticationType      types.String `tfsdk:"authentication_type"`
	Priority                types.Int64  `tfsdk:"priority"`
	TrackInterface          types.List   `tfsdk:"track_interface"`
	TrackRoute              types.List   `tfsdk:"track_route"`
}

type interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackInterface struct {
	Interface    types.String `tfsdk:"interface"`
	PriorityCost types.Int64  `tfsdk:"priority_cost"`
}

type interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackRoute struct {
	Route           types.String `tfsdk:"route"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
	PriorityCost    types.Int64  `tfsdk:"priority_cost"`
}

type interfaceLogicalBlockFamilyInetBlockDhcp struct {
	SrxOldOptionName                          types.Bool   `tfsdk:"srx_old_option_name"`
	ClientIdentifierPrefixHostname            types.Bool   `tfsdk:"client_identifier_prefix_hostname"`
	ClientIdentifierPrefixRoutingInstanceName types.Bool   `tfsdk:"client_identifier_prefix_routing_instance_name"`
	ForceDiscover                             types.Bool   `tfsdk:"force_discover"`
	LeaseTimeInfinite                         types.Bool   `tfsdk:"lease_time_infinite"`
	NoDNSInstall                              types.Bool   `tfsdk:"no_dns_install"`
	OptionsNoHostname                         types.Bool   `tfsdk:"options_no_hostname"`
	UpdateServer                              types.Bool   `tfsdk:"update_server"`
	ClientIdentifierASCII                     types.String `tfsdk:"client_identifier_ascii"`
	ClientIdentifierHexadecimal               types.String `tfsdk:"client_identifier_hexadecimal"`
	ClientIdentifierUseInterfaceDescription   types.String `tfsdk:"client_identifier_use_interface_description"`
	ClientIdentifierUseridASCII               types.String `tfsdk:"client_identifier_userid_ascii"`
	ClientIdentifierUseridHexadecimal         types.String `tfsdk:"client_identifier_userid_hexadecimal"`
	LeaseTime                                 types.Int64  `tfsdk:"lease_time"`
	Metric                                    types.Int64  `tfsdk:"metric"`
	RetransmissionAttempt                     types.Int64  `tfsdk:"retransmission_attempt"`
	RetransmissionInterval                    types.Int64  `tfsdk:"retransmission_interval"`
	ServerAddress                             types.String `tfsdk:"server_address"`
	VendorID                                  types.String `tfsdk:"vendor_id"`
}

type interfaceLogicalBlockFamilyInet6 struct {
	DadDisable     types.Bool                                         `tfsdk:"dad_disable"`
	SamplingInput  types.Bool                                         `tfsdk:"sampling_input"`
	SamplingOutput types.Bool                                         `tfsdk:"sampling_output"`
	FilterInput    types.String                                       `tfsdk:"filter_input"`
	FilterOutput   types.String                                       `tfsdk:"filter_output"`
	Mtu            types.Int64                                        `tfsdk:"mtu"`
	Address        []interfaceLogicalBlockFamilyInet6BlockAddress     `tfsdk:"address"`
	DHCPv6Client   *interfaceLogicalBlockFamilyInet6BlockDhcpV6Client `tfsdk:"dhcpv6_client"`
	RPFCheck       *interfaceLogicalBlockFamilyBlockRPFCheck          `tfsdk:"rpf_check"`
}

type interfaceLogicalBlockFamilyInet6Config struct {
	DadDisable     types.Bool                                               `tfsdk:"dad_disable"`
	SamplingInput  types.Bool                                               `tfsdk:"sampling_input"`
	SamplingOutput types.Bool                                               `tfsdk:"sampling_output"`
	FilterInput    types.String                                             `tfsdk:"filter_input"`
	FilterOutput   types.String                                             `tfsdk:"filter_output"`
	Mtu            types.Int64                                              `tfsdk:"mtu"`
	Address        types.List                                               `tfsdk:"address"`
	DHCPv6Client   *interfaceLogicalBlockFamilyInet6BlockDhcpV6ClientConfig `tfsdk:"dhcpv6_client"`
	RPFCheck       *interfaceLogicalBlockFamilyBlockRPFCheck                `tfsdk:"rpf_check"`
}

type interfaceLogicalBlockFamilyInet6BlockAddress struct {
	Preferred types.Bool                                                   `tfsdk:"preferred"`
	Primary   types.Bool                                                   `tfsdk:"primary"`
	CidrIP    types.String                                                 `tfsdk:"cidr_ip"`
	VRRPGroup []interfaceLogicalBlockFamilyInet6BlockAddressBlockVRRPGroup `tfsdk:"vrrp_group"`
}

type interfaceLogicalBlockFamilyInet6BlockAddressConfig struct {
	Preferred types.Bool   `tfsdk:"preferred"`
	Primary   types.Bool   `tfsdk:"primary"`
	CidrIP    types.String `tfsdk:"cidr_ip"`
	VRRPGroup types.List   `tfsdk:"vrrp_group"`
}

//nolint:lll
type interfaceLogicalBlockFamilyInet6BlockAddressBlockVRRPGroup struct {
	AcceptData              types.Bool                                                                 `tfsdk:"accept_data"`
	NoAcceptData            types.Bool                                                                 `tfsdk:"no_accept_data"`
	Preempt                 types.Bool                                                                 `tfsdk:"preempt"`
	NoPreempt               types.Bool                                                                 `tfsdk:"no_preempt"`
	Identifier              types.Int64                                                                `tfsdk:"identifier"`
	VirtualAddress          []types.String                                                             `tfsdk:"virtual_address"`
	VirutalLinkLocalAddress types.String                                                               `tfsdk:"virtual_link_local_address"`
	AdvertiseInterval       types.Int64                                                                `tfsdk:"advertise_interval"`
	AdvertisementsThreshold types.Int64                                                                `tfsdk:"advertisements_threshold"`
	Priority                types.Int64                                                                `tfsdk:"priority"`
	TrackInterface          []interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackInterface `tfsdk:"track_interface"`
	TrackRoute              []interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackRoute     `tfsdk:"track_route"`
}

type interfaceLogicalBlockFamilyInet6BlockAddressBlockVRRPGroupConfig struct {
	AcceptData              types.Bool   `tfsdk:"accept_data"`
	NoAcceptData            types.Bool   `tfsdk:"no_accept_data"`
	Preempt                 types.Bool   `tfsdk:"preempt"`
	NoPreempt               types.Bool   `tfsdk:"no_preempt"`
	Identifier              types.Int64  `tfsdk:"identifier"`
	VirtualAddress          types.List   `tfsdk:"virtual_address"`
	VirutalLinkLocalAddress types.String `tfsdk:"virtual_link_local_address"`
	AdvertiseInterval       types.Int64  `tfsdk:"advertise_interval"`
	AdvertisementsThreshold types.Int64  `tfsdk:"advertisements_threshold"`
	Priority                types.Int64  `tfsdk:"priority"`
	TrackInterface          types.List   `tfsdk:"track_interface"`
	TrackRoute              types.List   `tfsdk:"track_route"`
}

type interfaceLogicalBlockFamilyInet6BlockDhcpV6Client struct {
	ClientIATypeNA                        types.Bool     `tfsdk:"client_ia_type_na"`
	ClientIATypePD                        types.Bool     `tfsdk:"client_ia_type_pd"`
	NoDNSInstall                          types.Bool     `tfsdk:"no_dns_install"`
	RapidCommit                           types.Bool     `tfsdk:"rapid_commit"`
	ClientIdentifierDuidType              types.String   `tfsdk:"client_identifier_duid_type"`
	ClientType                            types.String   `tfsdk:"client_type"`
	PrefixDelegatingPreferredPrefixLength types.Int64    `tfsdk:"prefix_delegating_preferred_prefix_length"`
	PrefixDelegatingSubPrefixLength       types.Int64    `tfsdk:"prefix_delegating_sub_prefix_length"`
	ReqOption                             []types.String `tfsdk:"req_option"`
	RetransmissionAttempt                 types.Int64    `tfsdk:"retransmission_attempt"`
	UpdateRouterAdvertisementInterface    []types.String `tfsdk:"update_router_advertisement_interface"`
	UpdateServer                          types.Bool     `tfsdk:"update_server"`
}

type interfaceLogicalBlockFamilyInet6BlockDhcpV6ClientConfig struct {
	ClientIATypeNA                        types.Bool   `tfsdk:"client_ia_type_na"`
	ClientIATypePD                        types.Bool   `tfsdk:"client_ia_type_pd"`
	NoDNSInstall                          types.Bool   `tfsdk:"no_dns_install"`
	RapidCommit                           types.Bool   `tfsdk:"rapid_commit"`
	ClientIdentifierDuidType              types.String `tfsdk:"client_identifier_duid_type"`
	ClientType                            types.String `tfsdk:"client_type"`
	PrefixDelegatingPreferredPrefixLength types.Int64  `tfsdk:"prefix_delegating_preferred_prefix_length"`
	PrefixDelegatingSubPrefixLength       types.Int64  `tfsdk:"prefix_delegating_sub_prefix_length"`
	ReqOption                             types.Set    `tfsdk:"req_option"`
	RetransmissionAttempt                 types.Int64  `tfsdk:"retransmission_attempt"`
	UpdateRouterAdvertisementInterface    types.Set    `tfsdk:"update_router_advertisement_interface"`
	UpdateServer                          types.Bool   `tfsdk:"update_server"`
}

type interfaceLogicalBlockTunnel struct {
	AllowFragmentation         types.Bool   `tfsdk:"allow_fragmentation"`
	DoNotFragment              types.Bool   `tfsdk:"do_not_fragment"`
	PathMtuDiscovery           types.Bool   `tfsdk:"path_mtu_discovery"`
	NoPathMtuDiscovery         types.Bool   `tfsdk:"no_path_mtu_discovery"`
	Destination                types.String `tfsdk:"destination"`
	Source                     types.String `tfsdk:"source"`
	FlowLabel                  types.Int64  `tfsdk:"flow_label"`
	RoutingInstanceDestination types.String `tfsdk:"routing_instance_destination"`
	TrafficClass               types.Int64  `tfsdk:"traffic_class"`
	TTL                        types.Int64  `tfsdk:"ttl"`
}

//nolint:gocognit
func (rsc *interfaceLogical) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config interfaceLogicalConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.SecurityZone.IsNull() {
		if !config.SecurityInboundProtocols.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("security_inbound_protocols"),
				"Missing Configuration Error",
				"security_zone must be specified with security_inbound_protocols",
			)
		}
		if !config.SecurityInboundServices.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("security_inbound_services"),
				"Missing Configuration Error",
				"security_zone must be specified with security_inbound_services",
			)
		}
	}
	if config.FamilyInet != nil {
		if config.FamilyInet.DHCP != nil {
			if !config.FamilyInet.Address.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet").AtName("dhcp").AtName("*"),
					"Conflict Configuration Error",
					"cannot set dhcp block if address block is used in family_inet block",
				)
			}
			if !config.FamilyInet.DHCP.ClientIdentifierASCII.IsNull() &&
				!config.FamilyInet.DHCP.ClientIdentifierHexadecimal.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet").AtName("dhcp").AtName("client_identifier_ascii"),
					"Conflict Configuration Error",
					"client_identifier_ascii and client_identifier_hexadecimal cannot be configured together "+
						"in dhcp block in family_inet block",
				)
			}
			if !config.FamilyInet.DHCP.ClientIdentifierUseridASCII.IsNull() &&
				!config.FamilyInet.DHCP.ClientIdentifierUseridHexadecimal.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet").AtName("dhcp").AtName("client_identifier_userid_ascii"),
					"Conflict Configuration Error",
					"client_identifier_userid_ascii and client_identifier_userid_hexadecimal cannot be configured together "+
						"in dhcp block in family_inet block",
				)
			}
			if !config.FamilyInet.DHCP.LeaseTime.IsNull() &&
				!config.FamilyInet.DHCP.LeaseTimeInfinite.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet").AtName("dhcp").AtName("lease_time"),
					"Conflict Configuration Error",
					"lease_time and lease_time_infinite cannot be configured together "+
						"in dhcp block in family_inet block",
				)
			}
		}
		if !config.FamilyInet.Address.IsNull() && !config.FamilyInet.Address.IsUnknown() {
			var configAddress []interfaceLogicalBlockFamilyInetBlockAddressConfig
			asDiags := config.FamilyInet.Address.ElementsAs(ctx, &configAddress, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			addressCIDRIP := make(map[string]struct{})
			for i, address := range configAddress {
				if !address.CidrIP.IsUnknown() {
					if _, ok := addressCIDRIP[address.CidrIP.ValueString()]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("family_inet").AtName("address").AtListIndex(i).AtName("cidr_ip"),
							"Conflict Configuration Error",
							fmt.Sprintf("multiple address blocks with the same cidr_ip %q in family_inet block",
								address.CidrIP.ValueString()),
						)
					}
					addressCIDRIP[address.CidrIP.ValueString()] = struct{}{}
				}
				if !address.VRRPGroup.IsNull() && !address.VRRPGroup.IsUnknown() {
					var configVRRPGroup []interfaceLogicalBlockFamilyInetBlockAddressBlockVRRPGroupConfig
					asDiags := address.VRRPGroup.ElementsAs(ctx, &configVRRPGroup, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}
					vrrpGroupID := make(map[int64]struct{})
					for ii, vrrpGroup := range configVRRPGroup {
						if !vrrpGroup.Identifier.IsUnknown() {
							if _, ok := vrrpGroupID[vrrpGroup.Identifier.ValueInt64()]; ok {
								resp.Diagnostics.AddAttributeError(
									path.Root("family_inet").AtName("address").AtListIndex(i).
										AtName("vrrp_group").AtListIndex(ii).AtName("identifier"),
									"Conflict Configuration Error",
									fmt.Sprintf("multiple vrrp_group blocks with the same identifier %d in address block %q in family_inet block",
										vrrpGroup.Identifier.ValueInt64(), address.CidrIP.ValueString()),
								)
							}
							vrrpGroupID[vrrpGroup.Identifier.ValueInt64()] = struct{}{}
						}

						if !vrrpGroup.TrackInterface.IsNull() && !vrrpGroup.TrackInterface.IsUnknown() {
							var configTrackInterface []interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackInterface
							asDiags := vrrpGroup.TrackInterface.ElementsAs(ctx, &configTrackInterface, false)
							if asDiags.HasError() {
								resp.Diagnostics.Append(asDiags...)

								return
							}
							trackInterfaceInterface := make(map[string]struct{})
							for iii, trackInterface := range configTrackInterface {
								if trackInterface.Interface.IsUnknown() {
									continue
								}
								if _, ok := trackInterfaceInterface[trackInterface.Interface.ValueString()]; ok {
									resp.Diagnostics.AddAttributeError(
										path.Root("family_inet").AtName("address").AtListIndex(i).
											AtName("vrrp_group").AtListIndex(ii).
											AtName("track_interface").AtListIndex(iii).AtName("interface"),
										"Conflict Configuration Error",
										fmt.Sprintf("multiple track_interface blocks with the same interface %q "+
											"in vrrp_group block %d in address block %q in family_inet block",
											trackInterface.Interface.ValueString(), vrrpGroup.Identifier.ValueInt64(), address.CidrIP.ValueString()),
									)
								}
								trackInterfaceInterface[trackInterface.Interface.ValueString()] = struct{}{}
							}
						}
						if !vrrpGroup.TrackRoute.IsNull() && !vrrpGroup.TrackRoute.IsUnknown() {
							var configTrackRoute []interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackRoute
							asDiags := vrrpGroup.TrackRoute.ElementsAs(ctx, &configTrackRoute, false)
							if asDiags.HasError() {
								resp.Diagnostics.Append(asDiags...)

								return
							}
							trackRouteRoute := make(map[string]struct{})
							for iii, trackRoute := range configTrackRoute {
								if trackRoute.Route.IsUnknown() {
									continue
								}
								if _, ok := trackRouteRoute[trackRoute.Route.ValueString()]; ok {
									resp.Diagnostics.AddAttributeError(
										path.Root("family_inet").AtName("address").AtListIndex(i).
											AtName("vrrp_group").AtListIndex(ii).
											AtName("track_route").AtListIndex(iii).AtName("route"),
										"Conflict Configuration Error",
										fmt.Sprintf("multiple track_route blocks with the same route %q "+
											"in vrrp_group block %d in address block %q in family_inet block",
											trackRoute.Route.ValueString(), vrrpGroup.Identifier.ValueInt64(), address.CidrIP.ValueString()),
									)
								}
								trackRouteRoute[trackRoute.Route.ValueString()] = struct{}{}
							}
						}
					}
				}
			}
		}
	}
	if config.FamilyInet6 != nil {
		if config.FamilyInet6.DHCPv6Client != nil {
			if config.FamilyInet6.DHCPv6Client.ClientIdentifierDuidType.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet6").AtName("dhcpv6_client").AtName("client_identifier_duid_type"),
					"Missing Configuration Error",
					"client_identifier_duid_type must be specified in dhcpv6_client block in family_inet6 block",
				)
			}
			if config.FamilyInet6.DHCPv6Client.ClientType.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet6").AtName("dhcpv6_client").AtName("client_type"),
					"Missing Configuration Error",
					"client_type must be specified in dhcpv6_client block in family_inet6 block",
				)
			}
			if !config.FamilyInet6.Address.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet6").AtName("dhcpv6_client").AtName("*"),
					"Conflict Configuration Error",
					"cannot set dhcpv6_client block if address block is used in family_inet6 block",
				)
			}
			if config.FamilyInet6.DHCPv6Client.ClientIATypeNA.IsNull() &&
				config.FamilyInet6.DHCPv6Client.ClientIATypePD.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet6").AtName("dhcpv6_client").AtName("*"),
					"Missing Configuration Error",
					"at least one client_ia_type_na or client_ia_type_pd must be specified",
				)
			}
		}
		if !config.FamilyInet6.Address.IsNull() && !config.FamilyInet6.Address.IsUnknown() {
			var configAddress []interfaceLogicalBlockFamilyInet6BlockAddressConfig
			asDiags := config.FamilyInet6.Address.ElementsAs(ctx, &configAddress, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			addressCIDRIP := make(map[string]struct{})
			for i, address := range configAddress {
				if !address.CidrIP.IsUnknown() {
					if _, ok := addressCIDRIP[address.CidrIP.ValueString()]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("family_inet6").AtName("address").AtListIndex(i).AtName("cidr_ip"),
							"Conflict Configuration Error",
							fmt.Sprintf("multiple address blocks with the same cidr_ip %q in family_inet6 block",
								address.CidrIP.ValueString()),
						)
					}
					addressCIDRIP[address.CidrIP.ValueString()] = struct{}{}
				}
				if !address.VRRPGroup.IsNull() && !address.VRRPGroup.IsUnknown() {
					var configVRRPGroup []interfaceLogicalBlockFamilyInet6BlockAddressBlockVRRPGroupConfig
					asDiags := address.VRRPGroup.ElementsAs(ctx, &configVRRPGroup, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}
					vrrpGroupID := make(map[int64]struct{})
					for ii, vrrpGroup := range configVRRPGroup {
						if !vrrpGroup.Identifier.IsUnknown() {
							if _, ok := vrrpGroupID[vrrpGroup.Identifier.ValueInt64()]; ok {
								resp.Diagnostics.AddAttributeError(
									path.Root("family_inet6").AtName("address").AtListIndex(i).
										AtName("vrrp_group").AtListIndex(ii).AtName("identifier"),
									"Conflict Configuration Error",
									fmt.Sprintf("multiple vrrp_group blocks with the same identifier %d in address block %q in family_inet6 block",
										vrrpGroup.Identifier.ValueInt64(), address.CidrIP.ValueString()),
								)
							}
							vrrpGroupID[vrrpGroup.Identifier.ValueInt64()] = struct{}{}
						}

						if !vrrpGroup.TrackInterface.IsNull() && !vrrpGroup.TrackInterface.IsUnknown() {
							var configTrackInterface []interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackInterface
							asDiags := vrrpGroup.TrackInterface.ElementsAs(ctx, &configTrackInterface, false)
							if asDiags.HasError() {
								resp.Diagnostics.Append(asDiags...)

								return
							}
							trackInterfaceInterface := make(map[string]struct{})
							for iii, trackInterface := range configTrackInterface {
								if trackInterface.Interface.IsUnknown() {
									continue
								}
								if _, ok := trackInterfaceInterface[trackInterface.Interface.ValueString()]; ok {
									resp.Diagnostics.AddAttributeError(
										path.Root("family_inet6").AtName("address").AtListIndex(i).
											AtName("vrrp_group").AtListIndex(ii).
											AtName("track_interface").AtListIndex(iii).AtName("interface"),
										"Conflict Configuration Error",
										fmt.Sprintf("multiple track_interface blocks with the same interface %q "+
											"in vrrp_group block %d in address block %q in family_inet6 block",
											trackInterface.Interface.ValueString(), vrrpGroup.Identifier.ValueInt64(), address.CidrIP.ValueString()),
									)
								}
								trackInterfaceInterface[trackInterface.Interface.ValueString()] = struct{}{}
							}
						}
						if !vrrpGroup.TrackRoute.IsNull() && !vrrpGroup.TrackRoute.IsUnknown() {
							var configTrackRoute []interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackRoute
							asDiags := vrrpGroup.TrackRoute.ElementsAs(ctx, &configTrackRoute, false)
							if asDiags.HasError() {
								resp.Diagnostics.Append(asDiags...)

								return
							}
							trackRouteRoute := make(map[string]struct{})
							for iii, trackRoute := range configTrackRoute {
								if trackRoute.Route.IsUnknown() {
									continue
								}
								if _, ok := trackRouteRoute[trackRoute.Route.ValueString()]; ok {
									resp.Diagnostics.AddAttributeError(
										path.Root("family_inet6").AtName("address").AtListIndex(i).
											AtName("vrrp_group").AtListIndex(ii).
											AtName("track_route").AtListIndex(iii).AtName("route"),
										"Conflict Configuration Error",
										fmt.Sprintf("multiple track_route blocks with the same route %q "+
											"in vrrp_group block %d in address block %q in family_inet6 block",
											trackRoute.Route.ValueString(), vrrpGroup.Identifier.ValueInt64(), address.CidrIP.ValueString()),
									)
								}
								trackRouteRoute[trackRoute.Route.ValueString()] = struct{}{}
							}
						}
					}
				}
			}
		}
	}

	if config.Tunnel != nil {
		if config.Tunnel.Destination.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("tunnel").AtName("destination"),
				"Missing Configuration Error",
				"destination must be specified in tunnel block",
			)
		}
		if config.Tunnel.Source.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("tunnel").AtName("source"),
				"Missing Configuration Error",
				"source must be specified in tunnel block",
			)
		}
		if !config.Tunnel.AllowFragmentation.IsNull() &&
			!config.Tunnel.DoNotFragment.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("tunnel").AtName("allow_fragmentation"),
				"Conflict Configuration Error",
				"allow_fragmentation and do_not_fragment cannot be configured together "+
					"in tunnel block",
			)
		}
		if !config.Tunnel.PathMtuDiscovery.IsNull() &&
			!config.Tunnel.NoPathMtuDiscovery.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("tunnel").AtName("path_mtu_discovery"),
				"Conflict Configuration Error",
				"path_mtu_discovery and no_path_mtu_discovery cannot be configured together "+
					"in tunnel block",
			)
		}
	}
}

func (rsc *interfaceLogical) ModifyPlan(
	ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var config, plan interfaceLogicalConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.VlanID.IsNull() {
		if config.VlanNoCompute.IsNull() {
			plan.computeVlanID()
		} else if plan.VlanNoCompute.ValueBool() {
			plan.VlanID = types.Int64Null()
		}
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func (rsc *interfaceLogical) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan interfaceLogicalData
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
	if !strings.Contains(plan.Name.ValueString(), ".") {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Bad Name",
			"could not create "+rsc.junosName()+" with name without a dot",
		)

		return
	}

	if plan.VlanID.IsUnknown() {
		if plan.VlanNoCompute.ValueBool() {
			plan.VlanID = types.Int64Null()
		} else {
			plan.computeVlanID()
		}
	}

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := delInterfaceNC(
			ctx,
			plan.Name.ValueString(),
			rsc.client.GroupInterfaceDelete(),
			junSess,
		); err != nil {
			resp.Diagnostics.AddError("Config Pre Set Error", err.Error())

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

	ncInt, emptyInt, _, err := checkInterfaceLogicalNCEmpty(
		ctx,
		plan.Name.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Pre Check Error", err.Error())

		return
	}
	if !ncInt && !emptyInt {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError(
			"Duplicate Configuration Error",
			fmt.Sprintf(rsc.junosName()+" %q already configured", plan.Name.ValueString()),
		)

		return
	}
	if ncInt {
		if err := delInterfaceNC(
			ctx,
			plan.Name.ValueString(),
			rsc.client.GroupInterfaceDelete(),
			junSess,
		); err != nil {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddError("Config Pre Set Error", err.Error())

			return
		}
	}

	if v := plan.SecurityZone.ValueString(); v != "" {
		if !junSess.CheckCompatibilitySecurity() {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddAttributeError(
				path.Root("security_zone"),
				"Compatibility Error",
				fmt.Sprintf("security zone arguments not compatible "+
					"with Junos device %q", junSess.SystemInformation.HardwareModel),
			)

			return
		}
		zonesExists, err := checkSecurityZonesExists(ctx, v, junSess)
		if err != nil {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddError("Pre Check Error", err.Error())

			return
		}
		if !zonesExists {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddAttributeError(
				path.Root("security_zone"),
				"Missing Configuration Error",
				fmt.Sprintf("security zone %q doesn't exist", v),
			)

			return
		}
	}

	if v := plan.RoutingInstance.ValueString(); v != "" {
		instanceExists, err := checkRoutingInstanceExists(ctx, v, junSess)
		if err != nil {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddError("Pre Check Error", err.Error())

			return
		}
		if !instanceExists {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddAttributeError(
				path.Root("routing_instance"),
				"Missing Configuration Error",
				fmt.Sprintf("routing instance %q doesn't exist", v),
			)

			return
		}
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

	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(
		ctx,
		plan.Name.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError("Post Check Error", err.Error())

		return
	}
	if ncInt {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf(rsc.junosName()+" %q always disable (NC) after commit "+
				"=> check your config", plan.Name.ValueString()),
		)

		return
	}
	if emptyInt && !setInt {
		intExists, err := junSess.CheckInterfaceExists(plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Post Check Error", err.Error())

			return
		}
		if !intExists {
			resp.Diagnostics.AddError(
				"Not Found Error",
				fmt.Sprintf(rsc.junosName()+" %q not exists and config can't found after commit"+
					"=> check your config", plan.Name.ValueString()),
			)

			return
		}
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *interfaceLogical) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data interfaceLogicalData
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
	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(
		ctx,
		state.Name.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		junos.MutexUnlock()
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}
	if ncInt {
		junos.MutexUnlock()
		resp.State.RemoveResource(ctx)

		return
	}
	if emptyInt && !setInt {
		intExists, err := junSess.CheckInterfaceExists(state.Name.ValueString())
		if err != nil {
			junos.MutexUnlock()
			resp.Diagnostics.AddError("Config Read Error", err.Error())

			return
		}
		if !intExists {
			junos.MutexUnlock()
			resp.State.RemoveResource(ctx)

			return
		}
	}

	err = data.read(ctx, state.Name.ValueString(), junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}

	data.St0AlsoOnDestroy = state.St0AlsoOnDestroy
	data.VlanNoCompute = state.VlanNoCompute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *interfaceLogical) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state interfaceLogicalData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.VlanID.IsUnknown() {
		if plan.VlanNoCompute.ValueBool() {
			plan.VlanID = types.Int64Null()
		} else {
			plan.computeVlanID()
		}
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.delOpts(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}
		if v := state.SecurityZone.ValueString(); v != "" {
			if err := state.delZone(ctx, junSess); err != nil {
				resp.Diagnostics.AddError("Config Del Error", err.Error())

				return
			}
		}
		if v := state.RoutingInstance.ValueString(); v != "" && v != plan.RoutingInstance.ValueString() {
			if err := state.delRoutingInstance(ctx, junSess); err != nil {
				resp.Diagnostics.AddError("Config Del Error", err.Error())

				return
			}
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

	if vSte := state.SecurityZone.ValueString(); vSte != "" {
		if vSte != "" {
			if err := state.delZone(ctx, junSess); err != nil {
				resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
				resp.Diagnostics.AddError("Config Del Error", err.Error())

				return
			}
		}
	}
	if vPln := plan.SecurityZone.ValueString(); vPln != "" && vPln != state.SecurityZone.ValueString() {
		if !junSess.CheckCompatibilitySecurity() {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddAttributeError(
				path.Root("security_zone"),
				"Compatibility Error",
				fmt.Sprintf("security zone arguments not compatible "+
					"with Junos device %q", junSess.SystemInformation.HardwareModel),
			)

			return
		}
		zonesExists, err := checkSecurityZonesExists(ctx, vPln, junSess)
		if err != nil {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddError("Pre Check Error", err.Error())

			return
		}
		if !zonesExists {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddAttributeError(
				path.Root("security_zone"),
				"Missing Configuration Error",
				fmt.Sprintf("security zone %q doesn't exist", vPln),
			)

			return
		}
	}

	if vSte, vPln := state.RoutingInstance.ValueString(), plan.RoutingInstance.ValueString(); vSte != vPln {
		if vSte != "" {
			if err := state.delRoutingInstance(ctx, junSess); err != nil {
				resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
				resp.Diagnostics.AddError("Config Del Error", err.Error())

				return
			}
		}
		if vPln != "" {
			instanceExists, err := checkRoutingInstanceExists(ctx, vPln, junSess)
			if err != nil {
				resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
				resp.Diagnostics.AddError("Pre Check Error", err.Error())

				return
			}
			if !instanceExists {
				resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
				resp.Diagnostics.AddAttributeError(
					path.Root("routing_instance"),
					"Missing Configuration Error",
					fmt.Sprintf("routing instance %q doesn't exist", vPln),
				)

				return
			}
		}
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

func (rsc *interfaceLogical) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state interfaceLogicalData
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

func (rsc *interfaceLogical) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	if strings.Count(req.ID, ".") != 1 {
		resp.Diagnostics.AddError(
			"Pre Check Error",
			fmt.Sprintf("name of interface need to have a dot, got %q", req.ID),
		)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(
		ctx,
		req.ID,
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError("Interface Read Error", err.Error())

		return
	}
	if ncInt {
		resp.Diagnostics.AddError(
			"Disable Error",
			fmt.Sprintf("interface %q is disabled (NC), import is not possible", req.ID),
		)

		return
	}
	if emptyInt && !setInt {
		intExists, err := junSess.CheckInterfaceExists(req.ID)
		if err != nil {
			resp.Diagnostics.AddError("Interface Read Error", err.Error())

			return
		}
		if !intExists {
			resp.Diagnostics.AddError(
				"Not Found Error",
				fmt.Sprintf("don't find "+rsc.junosName()+" with id %q "+
					"(id must be <name>)", req.ID),
			)

			return
		}
	}

	var data interfaceLogicalData
	if err := data.read(ctx, req.ID, junSess); err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}

	if data.VlanID.IsNull() {
		intCut := strings.Split(req.ID, ".")
		if !bchk.InSlice(intCut[0], []string{junos.St0Word, "irb", "vlan"}) &&
			intCut[1] != "0" {
			data.VlanNoCompute = types.BoolValue(true)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rscData *interfaceLogicalData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscCfg *interfaceLogicalConfig) computeVlanID() {
	if !rscCfg.VlanID.IsUnknown() {
		return
	}
	if rscCfg.Name.IsUnknown() {
		return
	}
	rscCfg.VlanID = types.Int64Null()

	intCut := strings.Split(rscCfg.Name.ValueString(), ".")
	if len(intCut) < 2 {
		return
	}
	if bchk.InSlice(
		intCut[0],
		[]string{
			junos.St0Word,
			"irb",
			"vlan",
		},
	) {
		return
	}
	v, err := tfdata.ConvAtoi64Value(intCut[1])
	if err == nil && v.ValueInt64() >= 1 && v.ValueInt64() <= 4094 {
		rscCfg.VlanID = v
	}
}

func (rscData *interfaceLogicalData) computeVlanID() {
	if !rscData.VlanID.IsUnknown() {
		return
	}
	rscData.VlanID = types.Int64Null()

	intCut := strings.Split(rscData.Name.ValueString(), ".")
	if len(intCut) < 2 {
		return
	}
	if bchk.InSlice(
		intCut[0],
		[]string{
			junos.St0Word,
			"irb",
			"vlan",
		},
	) {
		return
	}
	v, err := tfdata.ConvAtoi64Value(intCut[1])
	if err == nil && v.ValueInt64() >= 1 && v.ValueInt64() <= 4094 {
		rscData.VlanID = v
	}
}

func checkInterfaceLogicalNCEmpty(
	_ context.Context, name, groupInterfaceDelete string, junSess *junos.Session,
) (
	ncInt, // interface is set with NC config
	emtyInt, // interface is emty not set or just with set
	justSet bool, // interface is empty with set
	_ error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"interfaces " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return false, false, false, err
	}
	showConfigLines := make([]string, 0)
	// remove unused lines
	for _, item := range strings.Split(showConfig, "\n") {
		// exclude ethernet-switching (parameters in junos_interface_physical)
		if strings.Contains(item, "ethernet-switching") {
			continue
		}
		if strings.Contains(item, junos.XMLStartTagConfigOut) {
			continue
		}
		if strings.Contains(item, junos.XMLEndTagConfigOut) {
			break
		}
		if item == "" {
			continue
		}
		showConfigLines = append(showConfigLines, item)
	}
	if len(showConfigLines) == 0 {
		return false, true, true, nil
	}
	showConfig = strings.Join(showConfigLines, "\n")
	if groupInterfaceDelete != "" {
		if showConfig == "set apply-groups "+groupInterfaceDelete {
			return true, false, false, nil
		}
	}
	if showConfig == "set description NC\nset disable" ||
		showConfig == "set disable\nset description NC" {
		return true, false, false, nil
	}
	switch {
	case showConfig == junos.SetLS:
		return false, true, true, nil
	case showConfig == junos.EmptyW:
		return false, true, false, nil
	default:
		return false, false, false, nil
	}
}

func (rscData *interfaceLogicalData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	intCut := strings.Split(rscData.Name.ValueString(), ".")
	if len(intCut) != 2 {
		return path.Root("name"),
			fmt.Errorf("the name %q doesn't contain one dot", rscData.Name.ValueString())
	}

	setPrefix := "set interfaces " + rscData.Name.ValueString() + " "
	configSet := []string{
		setPrefix,
	}

	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if rscData.Disable.ValueBool() {
		if rscData.Description.ValueString() == "NC" {
			return path.Root("disable"), fmt.Errorf("disable=true and description=NC is not allowed " +
				"because the provider might consider the resource deleted")
		}
		configSet = append(configSet, setPrefix+"disable")
	}
	if rscData.FamilyInet != nil {
		configSet = append(configSet, setPrefix+"family inet")

		addressCIDRIP := make(map[string]struct{})
		for i, address := range rscData.FamilyInet.Address {
			if _, ok := addressCIDRIP[address.CidrIP.ValueString()]; ok {
				return path.Root("family_inet").AtName("address").AtListIndex(i).AtName("cidr_ip"),
					fmt.Errorf("multiple address blocks with the same cidr_ip %q in family_inet block", address.CidrIP.ValueString())
			}
			addressCIDRIP[address.CidrIP.ValueString()] = struct{}{}
			blockSet, pathErr, err := address.set(setPrefix, path.Root("family_inet").AtName("address").AtListIndex(i))
			if err != nil {
				return pathErr, err
			}
			configSet = append(configSet, blockSet...)
		}
		if rscData.FamilyInet.DHCP != nil {
			configSet = append(configSet, rscData.FamilyInet.DHCP.set(setPrefix)...)
		}
		if v := rscData.FamilyInet.FilterInput.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"family inet filter input \""+v+"\"")
		}
		if v := rscData.FamilyInet.FilterOutput.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"family inet filter output \""+v+"\"")
		}
		if !rscData.FamilyInet.Mtu.IsNull() {
			configSet = append(configSet, setPrefix+"family inet mtu "+utils.ConvI64toa(rscData.FamilyInet.Mtu.ValueInt64()))
		}
		if rscData.FamilyInet.RPFCheck != nil {
			configSet = append(configSet, setPrefix+"family inet rpf-check")

			if v := rscData.FamilyInet.RPFCheck.FailFilter.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"family inet rpf-check fail-filter \""+v+"\"")
			}
			if rscData.FamilyInet.RPFCheck.ModeLoose.ValueBool() {
				configSet = append(configSet, setPrefix+"family inet rpf-check mode loose")
			}
		}
		if rscData.FamilyInet.SamplingInput.ValueBool() {
			configSet = append(configSet, setPrefix+"family inet sampling input")
		}
		if rscData.FamilyInet.SamplingOutput.ValueBool() {
			configSet = append(configSet, setPrefix+"family inet sampling output")
		}
	}
	if rscData.FamilyInet6 != nil {
		configSet = append(configSet, setPrefix+"family inet6")

		addressCIDRIP := make(map[string]struct{})
		for i, address := range rscData.FamilyInet6.Address {
			if _, ok := addressCIDRIP[address.CidrIP.ValueString()]; ok {
				return path.Root("family_inet6").AtName("address").AtListIndex(i).AtName("cidr_ip"),
					fmt.Errorf("multiple address blocks with the same cidr_ip %q in family_inet6 block", address.CidrIP.ValueString())
			}
			addressCIDRIP[address.CidrIP.ValueString()] = struct{}{}
			blockSet, pathErr, err := address.set(setPrefix, path.Root("family_inet6").AtName("address").AtListIndex(i))
			if err != nil {
				return pathErr, err
			}
			configSet = append(configSet, blockSet...)
		}
		if rscData.FamilyInet6.DHCPv6Client != nil {
			configSet = append(configSet, rscData.FamilyInet6.DHCPv6Client.set(setPrefix)...)
		}
		if rscData.FamilyInet6.DadDisable.ValueBool() {
			configSet = append(configSet, setPrefix+"family inet6 dad-disable")
		}
		if v := rscData.FamilyInet6.FilterInput.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"family inet6 filter input \""+v+"\"")
		}
		if v := rscData.FamilyInet6.FilterOutput.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"family inet6 filter output \""+v+"\"")
		}
		if !rscData.FamilyInet6.Mtu.IsNull() {
			configSet = append(configSet, setPrefix+"family inet6 mtu "+utils.ConvI64toa(rscData.FamilyInet6.Mtu.ValueInt64()))
		}
		if rscData.FamilyInet6.RPFCheck != nil {
			configSet = append(configSet, setPrefix+"family inet6 rpf-check")

			if v := rscData.FamilyInet6.RPFCheck.FailFilter.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"family inet6 rpf-check fail-filter \""+v+"\"")
			}
			if rscData.FamilyInet6.RPFCheck.ModeLoose.ValueBool() {
				configSet = append(configSet, setPrefix+"family inet6 rpf-check mode loose")
			}
		}
		if rscData.FamilyInet6.SamplingInput.ValueBool() {
			configSet = append(configSet, setPrefix+"family inet6 sampling input")
		}
		if rscData.FamilyInet6.SamplingOutput.ValueBool() {
			configSet = append(configSet, setPrefix+"family inet6 sampling output")
		}
	}
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, junos.SetRoutingInstances+v+" interface "+rscData.Name.ValueString())
	}
	if securityZone := rscData.SecurityZone.ValueString(); securityZone != "" {
		configSet = append(configSet, "set security zones security-zone "+securityZone+
			" interfaces "+rscData.Name.ValueString())

		for _, v := range rscData.SecurityInboundProtocols {
			configSet = append(configSet, "set security zones security-zone "+securityZone+
				" interfaces "+rscData.Name.ValueString()+" host-inbound-traffic protocols "+v.ValueString())
		}
		for _, v := range rscData.SecurityInboundServices {
			configSet = append(configSet, "set security zones security-zone "+securityZone+
				" interfaces "+rscData.Name.ValueString()+" host-inbound-traffic system-services "+v.ValueString())
		}
	}
	if rscData.Tunnel != nil {
		configSet = append(configSet, setPrefix+"tunnel destination "+rscData.Tunnel.Destination.ValueString())
		configSet = append(configSet, setPrefix+"tunnel source "+rscData.Tunnel.Source.ValueString())
		if rscData.Tunnel.AllowFragmentation.ValueBool() {
			configSet = append(configSet, setPrefix+"tunnel allow-fragmentation")
		}
		if rscData.Tunnel.DoNotFragment.ValueBool() {
			configSet = append(configSet, setPrefix+"tunnel do-not-fragment")
		}
		if !rscData.Tunnel.FlowLabel.IsNull() {
			configSet = append(configSet, setPrefix+"tunnel flow-label "+
				utils.ConvI64toa(rscData.Tunnel.FlowLabel.ValueInt64()))
		}
		if rscData.Tunnel.PathMtuDiscovery.ValueBool() {
			configSet = append(configSet, setPrefix+"tunnel path-mtu-discovery")
		}
		if rscData.Tunnel.NoPathMtuDiscovery.ValueBool() {
			configSet = append(configSet, setPrefix+"tunnel no-path-mtu-discovery")
		}
		if v := rscData.Tunnel.RoutingInstanceDestination.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"tunnel routing-instance destination "+v)
		}
		if !rscData.Tunnel.TrafficClass.IsNull() {
			configSet = append(configSet, setPrefix+"tunnel traffic-class "+
				utils.ConvI64toa(rscData.Tunnel.TrafficClass.ValueInt64()))
		}
		if !rscData.Tunnel.TTL.IsNull() {
			configSet = append(configSet, setPrefix+"tunnel ttl "+
				utils.ConvI64toa(rscData.Tunnel.TTL.ValueInt64()))
		}
	}
	if !rscData.VlanID.IsNull() {
		configSet = append(configSet, setPrefix+"vlan-id "+utils.ConvI64toa(rscData.VlanID.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *interfaceLogicalBlockFamilyInetBlockAddress) set(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "family inet address " + block.CidrIP.ValueString()
	configSet := []string{
		setPrefix,
	}

	if block.Preferred.ValueBool() {
		configSet = append(configSet, setPrefix+" preferred")
	}
	if block.Primary.ValueBool() {
		configSet = append(configSet, setPrefix+" primary")
	}
	vrrpGroupID := make(map[int64]struct{})
	for i, vrrpGroup := range block.VRRPGroup {
		if strings.HasPrefix(setPrefix, "set interfaces st0.") {
			return configSet, pathRoot.AtName("vrrp_group").AtListIndex(i).AtName("*"),
				fmt.Errorf("vrrp not available on st0 interface")
		}

		if _, ok := vrrpGroupID[vrrpGroup.Identifier.ValueInt64()]; ok {
			return configSet, pathRoot.AtName("vrrp_group").AtListIndex(i).AtName("identifier"),
				fmt.Errorf("multiple vrrp_group blocks with the same identifier %d in address block %q in family_inet block",
					vrrpGroup.Identifier.ValueInt64(), block.CidrIP.ValueString())
		}
		vrrpGroupID[vrrpGroup.Identifier.ValueInt64()] = struct{}{}

		setPrefixVRRPGroup := setPrefix + " vrrp-group " + utils.ConvI64toa(vrrpGroup.Identifier.ValueInt64()) + " "
		for _, v := range vrrpGroup.VirtualAddress {
			configSet = append(configSet, setPrefixVRRPGroup+"virtual-address "+v.ValueString())
		}
		if !vrrpGroup.AdvertiseInterval.IsNull() {
			configSet = append(configSet, setPrefixVRRPGroup+"advertise-interval "+
				utils.ConvI64toa(vrrpGroup.AdvertiseInterval.ValueInt64()))
		}
		if v := vrrpGroup.AuthenticationKey.ValueString(); v != "" {
			configSet = append(configSet, setPrefixVRRPGroup+"authentication-key \""+v+"\"")
		}
		if v := vrrpGroup.AuthenticationType.ValueString(); v != "" {
			configSet = append(configSet, setPrefixVRRPGroup+"authentication-type "+v)
		}
		if vrrpGroup.AcceptData.ValueBool() {
			configSet = append(configSet, setPrefixVRRPGroup+"accept-data")
		}
		if vrrpGroup.NoAcceptData.ValueBool() {
			configSet = append(configSet, setPrefixVRRPGroup+"no-accept-data")
		}
		if !vrrpGroup.AdvertisementsThreshold.IsNull() {
			configSet = append(configSet, setPrefixVRRPGroup+"advertisements-threshold "+
				utils.ConvI64toa(vrrpGroup.AdvertisementsThreshold.ValueInt64()))
		}
		if vrrpGroup.Preempt.ValueBool() {
			configSet = append(configSet, setPrefixVRRPGroup+"preempt")
		}
		if vrrpGroup.NoPreempt.ValueBool() {
			configSet = append(configSet, setPrefixVRRPGroup+"no-preempt")
		}
		if !vrrpGroup.Priority.IsNull() {
			configSet = append(configSet, setPrefixVRRPGroup+"priority "+
				utils.ConvI64toa(vrrpGroup.Priority.ValueInt64()))
		}
		trackInterfaceInterface := make(map[string]struct{})
		for ii, trackInterface := range vrrpGroup.TrackInterface {
			if _, ok := trackInterfaceInterface[trackInterface.Interface.ValueString()]; ok {
				return configSet,
					pathRoot.AtName("vrrp_group").AtListIndex(i).AtName("track_interface").AtListIndex(ii).AtName("interface"),
					fmt.Errorf("multiple track_interface blocks with the same interface %q "+
						"in vrrp_group block %d in address block %q in family_inet block",
						trackInterface.Interface.ValueString(),
						vrrpGroup.Identifier.ValueInt64(),
						block.CidrIP.ValueString(),
					)
			}
			trackInterfaceInterface[trackInterface.Interface.ValueString()] = struct{}{}

			configSet = append(configSet, setPrefixVRRPGroup+"track interface "+trackInterface.Interface.ValueString()+
				" priority-cost "+utils.ConvI64toa(trackInterface.PriorityCost.ValueInt64()))
		}
		trackRouteRoute := make(map[string]struct{})
		for ii, trackRoute := range vrrpGroup.TrackRoute {
			if _, ok := trackRouteRoute[trackRoute.Route.ValueString()]; ok {
				return configSet,
					pathRoot.AtName("vrrp_group").AtListIndex(i).AtName("track_route").AtListIndex(ii).AtName("route"),
					fmt.Errorf("multiple track_route blocks with the same route %q "+
						"in vrrp_group block %d in address block %q in family_inet block",
						trackRoute.Route.ValueString(),
						vrrpGroup.Identifier.ValueInt64(),
						block.CidrIP.ValueString(),
					)
			}
			trackRouteRoute[trackRoute.Route.ValueString()] = struct{}{}

			configSet = append(configSet, setPrefixVRRPGroup+"track route "+trackRoute.Route.ValueString()+
				" routing-instance "+trackRoute.RoutingInstance.ValueString()+
				" priority-cost "+utils.ConvI64toa(trackRoute.PriorityCost.ValueInt64()))
		}
	}

	return configSet, path.Empty(), nil
}

func (block *interfaceLogicalBlockFamilyInet6BlockAddress) set(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "family inet6 address " + block.CidrIP.ValueString()
	configSet := []string{
		setPrefix,
	}

	if block.Preferred.ValueBool() {
		configSet = append(configSet, setPrefix+" preferred")
	}
	if block.Primary.ValueBool() {
		configSet = append(configSet, setPrefix+" primary")
	}
	vrrpGroupID := make(map[int64]struct{})
	for i, vrrpGroup := range block.VRRPGroup {
		if strings.HasPrefix(setPrefix, "set interfaces st0.") {
			return configSet, pathRoot.AtName("vrrp_group").AtListIndex(i).AtName("*"),
				fmt.Errorf("vrrp not available on st0 interface")
		}

		if _, ok := vrrpGroupID[vrrpGroup.Identifier.ValueInt64()]; ok {
			return configSet, pathRoot.AtName("vrrp_group").AtListIndex(i).AtName("identifier"),
				fmt.Errorf("multiple blocks vrrp_group with the same identifier %d in address block %q in family_inet6 block",
					vrrpGroup.Identifier.ValueInt64(), block.CidrIP.ValueString())
		}
		vrrpGroupID[vrrpGroup.Identifier.ValueInt64()] = struct{}{}

		setPrefixVRRPGroup := setPrefix + " vrrp-inet6-group " + utils.ConvI64toa(vrrpGroup.Identifier.ValueInt64()) + " "
		for _, v := range vrrpGroup.VirtualAddress {
			configSet = append(configSet, setPrefixVRRPGroup+"virtual-inet6-address "+v.ValueString())
		}
		configSet = append(configSet,
			setPrefixVRRPGroup+"virtual-link-local-address "+vrrpGroup.VirutalLinkLocalAddress.ValueString())

		if !vrrpGroup.AdvertiseInterval.IsNull() {
			configSet = append(configSet, setPrefixVRRPGroup+"inet6-advertise-interval "+
				utils.ConvI64toa(vrrpGroup.AdvertiseInterval.ValueInt64()))
		}
		if vrrpGroup.AcceptData.ValueBool() {
			configSet = append(configSet, setPrefixVRRPGroup+"accept-data")
		}
		if vrrpGroup.NoAcceptData.ValueBool() {
			configSet = append(configSet, setPrefixVRRPGroup+"no-accept-data")
		}
		if !vrrpGroup.AdvertisementsThreshold.IsNull() {
			configSet = append(configSet, setPrefixVRRPGroup+"advertisements-threshold "+
				utils.ConvI64toa(vrrpGroup.AdvertisementsThreshold.ValueInt64()))
		}
		if vrrpGroup.Preempt.ValueBool() {
			configSet = append(configSet, setPrefixVRRPGroup+"preempt")
		}
		if vrrpGroup.NoPreempt.ValueBool() {
			configSet = append(configSet, setPrefixVRRPGroup+"no-preempt")
		}
		if !vrrpGroup.Priority.IsNull() {
			configSet = append(configSet, setPrefixVRRPGroup+"priority "+
				utils.ConvI64toa(vrrpGroup.Priority.ValueInt64()))
		}
		trackInterfaceInterface := make(map[string]struct{})
		for ii, trackInterface := range vrrpGroup.TrackInterface {
			if _, ok := trackInterfaceInterface[trackInterface.Interface.ValueString()]; ok {
				return configSet,
					pathRoot.AtName("vrrp_group").AtListIndex(i).AtName("track_interface").AtListIndex(ii).AtName("interface"),
					fmt.Errorf("multiple track_interface blocks with the same interface %q "+
						"in vrrp_group block %d in address block %q in family_inet6 block",
						trackInterface.Interface.ValueString(),
						vrrpGroup.Identifier.ValueInt64(),
						block.CidrIP.ValueString(),
					)
			}
			trackInterfaceInterface[trackInterface.Interface.ValueString()] = struct{}{}

			configSet = append(configSet, setPrefixVRRPGroup+"track interface "+trackInterface.Interface.ValueString()+
				" priority-cost "+utils.ConvI64toa(trackInterface.PriorityCost.ValueInt64()))
		}
		trackRouteRoute := make(map[string]struct{})
		for ii, trackRoute := range vrrpGroup.TrackRoute {
			if _, ok := trackRouteRoute[trackRoute.Route.ValueString()]; ok {
				return configSet,
					pathRoot.AtName("vrrp_group").AtListIndex(i).AtName("track_route").AtListIndex(ii).AtName("route"),
					fmt.Errorf("multiple track_route blocks with the same route %q "+
						"in vrrp_group block %d in address block %q in family_inet6 block",
						trackRoute.Route.ValueString(),
						vrrpGroup.Identifier.ValueInt64(),
						block.CidrIP.ValueString(),
					)
			}
			trackRouteRoute[trackRoute.Route.ValueString()] = struct{}{}

			configSet = append(configSet, setPrefixVRRPGroup+"track route "+trackRoute.Route.ValueString()+
				" routing-instance "+trackRoute.RoutingInstance.ValueString()+
				" priority-cost "+utils.ConvI64toa(trackRoute.PriorityCost.ValueInt64()))
		}
	}

	return configSet, path.Empty(), nil
}

func (block *interfaceLogicalBlockFamilyInetBlockDhcp) set(setPrefix string) []string {
	setPrefix += "family inet dhcp"
	if block.SrxOldOptionName.ValueBool() {
		setPrefix += "-client "
	} else {
		setPrefix += " "
	}
	configSet := []string{
		setPrefix,
	}

	if v := block.ClientIdentifierASCII.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"client-identifier ascii \""+v+"\"")
	}
	if v := block.ClientIdentifierHexadecimal.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"client-identifier hexadecimal "+v)
	}
	if block.ClientIdentifierPrefixHostname.ValueBool() {
		configSet = append(configSet, setPrefix+"client-identifier prefix host-name")
	}
	if block.ClientIdentifierPrefixRoutingInstanceName.ValueBool() {
		configSet = append(configSet, setPrefix+"client-identifier prefix routing-instance-name")
	}
	if v := block.ClientIdentifierUseInterfaceDescription.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"client-identifier use-interface-description "+v)
	}
	if v := block.ClientIdentifierUseridASCII.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"client-identifier user-id ascii \""+v+"\"")
	}
	if v := block.ClientIdentifierUseridHexadecimal.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"client-identifier user-id hexadecimal "+v)
	}
	if block.ForceDiscover.ValueBool() {
		configSet = append(configSet, setPrefix+"force-discover")
	}
	if !block.LeaseTime.IsNull() {
		configSet = append(configSet, setPrefix+"lease-time "+
			utils.ConvI64toa(block.LeaseTime.ValueInt64()))
	}
	if block.LeaseTimeInfinite.ValueBool() {
		configSet = append(configSet, setPrefix+"lease-time infinite")
	}
	if !block.Metric.IsNull() {
		configSet = append(configSet, setPrefix+"metric "+
			utils.ConvI64toa(block.Metric.ValueInt64()))
	}
	if block.NoDNSInstall.ValueBool() {
		configSet = append(configSet, setPrefix+"no-dns-install")
	}
	if block.OptionsNoHostname.ValueBool() {
		configSet = append(configSet, setPrefix+"options no-hostname")
	}
	if !block.RetransmissionAttempt.IsNull() {
		configSet = append(configSet, setPrefix+"retransmission-attempt "+
			utils.ConvI64toa(block.RetransmissionAttempt.ValueInt64()))
	}
	if !block.RetransmissionInterval.IsNull() {
		configSet = append(configSet, setPrefix+"retransmission-interval "+
			utils.ConvI64toa(block.RetransmissionInterval.ValueInt64()))
	}
	if v := block.ServerAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"server-address "+v)
	}
	if block.UpdateServer.ValueBool() {
		configSet = append(configSet, setPrefix+"update-server")
	}
	if v := block.VendorID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"vendor-id \""+v+"\"")
	}

	return configSet
}

func (block *interfaceLogicalBlockFamilyInet6BlockDhcpV6Client) set(setPrefix string) []string {
	setPrefix += "family inet6 dhcpv6-client "
	configSet := []string{
		setPrefix + "client-identifier duid-type " + block.ClientIdentifierDuidType.ValueString(),
		setPrefix + "client-type " + block.ClientType.ValueString(),
	}

	if block.ClientIATypeNA.ValueBool() {
		configSet = append(configSet, setPrefix+"client-ia-type ia-na")
	}
	if block.ClientIATypePD.ValueBool() {
		configSet = append(configSet, setPrefix+"client-ia-type ia-pd")
	}
	if block.NoDNSInstall.ValueBool() {
		configSet = append(configSet, setPrefix+"no-dns-install")
	}
	if !block.PrefixDelegatingPreferredPrefixLength.IsNull() {
		configSet = append(configSet, setPrefix+"prefix-delegating preferred-prefix-length "+
			utils.ConvI64toa(block.PrefixDelegatingPreferredPrefixLength.ValueInt64()))
	}
	if !block.PrefixDelegatingSubPrefixLength.IsNull() {
		configSet = append(configSet, setPrefix+"prefix-delegating sub-prefix-length "+
			utils.ConvI64toa(block.PrefixDelegatingSubPrefixLength.ValueInt64()))
	}
	if block.RapidCommit.ValueBool() {
		configSet = append(configSet, setPrefix+"rapid-commit")
	}
	for _, v := range block.ReqOption {
		configSet = append(configSet, setPrefix+"req-option "+v.ValueString())
	}
	if !block.RetransmissionAttempt.IsNull() {
		configSet = append(configSet, setPrefix+"retransmission-attempt "+
			utils.ConvI64toa(block.RetransmissionAttempt.ValueInt64()))
	}
	for _, v := range block.UpdateRouterAdvertisementInterface {
		configSet = append(configSet, setPrefix+"update-router-advertisement interface "+v.ValueString())
	}
	if block.UpdateServer.ValueBool() {
		configSet = append(configSet, setPrefix+"update-server")
	}

	return configSet
}

func (rscData *interfaceLogicalData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"interfaces " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}

	rscData.Name = types.StringValue(name)
	rscData.fillID()
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			// exclude ethernet-switching (parameters in junos_interface_physical)
			if strings.Contains(item, "ethernet-switching") {
				continue
			}
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
			case itemTrim == "disable":
				rscData.Disable = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "family inet6"):
				if rscData.FamilyInet6 == nil {
					rscData.FamilyInet6 = &interfaceLogicalBlockFamilyInet6{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " address "):
					itemTrimFields := strings.Split(itemTrim, " ")
					var address interfaceLogicalBlockFamilyInet6BlockAddress
					rscData.FamilyInet6.Address, address = tfdata.ExtractBlockWithTFTypesString(
						rscData.FamilyInet6.Address, "CidrIP", itemTrimFields[0],
					)
					address.CidrIP = types.StringValue(itemTrimFields[0])

					if len(itemTrimFields) > 1 {
						balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
						if err := address.read(itemTrim); err != nil {
							return err
						}
					}
					rscData.FamilyInet6.Address = append(rscData.FamilyInet6.Address, address)
				case balt.CutPrefixInString(&itemTrim, " dhcpv6-client "):
					if rscData.FamilyInet6.DHCPv6Client == nil {
						rscData.FamilyInet6.DHCPv6Client = &interfaceLogicalBlockFamilyInet6BlockDhcpV6Client{}
					}
					if err := rscData.FamilyInet6.DHCPv6Client.read(itemTrim); err != nil {
						return err
					}
				case itemTrim == " dad-disable":
					rscData.FamilyInet6.DadDisable = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, " filter input "):
					rscData.FamilyInet6.FilterInput = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, " filter output "):
					rscData.FamilyInet6.FilterOutput = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, " mtu "):
					rscData.FamilyInet6.Mtu, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " rpf-check"):
					if rscData.FamilyInet6.RPFCheck == nil {
						rscData.FamilyInet6.RPFCheck = &interfaceLogicalBlockFamilyBlockRPFCheck{}
					}
					switch {
					case balt.CutPrefixInString(&itemTrim, " fail-filter "):
						rscData.FamilyInet6.RPFCheck.FailFilter = types.StringValue(strings.Trim(itemTrim, "\""))
					case itemTrim == " mode loose":
						rscData.FamilyInet6.RPFCheck.ModeLoose = types.BoolValue(true)
					}
				case itemTrim == " sampling input":
					rscData.FamilyInet6.SamplingInput = types.BoolValue(true)
				case itemTrim == " sampling output":
					rscData.FamilyInet6.SamplingOutput = types.BoolValue(true)
				}
			case balt.CutPrefixInString(&itemTrim, "family inet"):
				if rscData.FamilyInet == nil {
					rscData.FamilyInet = &interfaceLogicalBlockFamilyInet{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " address "):
					itemTrimFields := strings.Split(itemTrim, " ")
					var address interfaceLogicalBlockFamilyInetBlockAddress
					rscData.FamilyInet.Address, address = tfdata.ExtractBlockWithTFTypesString(
						rscData.FamilyInet.Address, "CidrIP", itemTrimFields[0],
					)
					address.CidrIP = types.StringValue(itemTrimFields[0])

					if len(itemTrimFields) > 1 {
						balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
						if err := address.read(itemTrim); err != nil {
							return err
						}
					}
					rscData.FamilyInet.Address = append(rscData.FamilyInet.Address, address)
				case strings.HasPrefix(itemTrim, " dhcp"):
					if rscData.FamilyInet.DHCP == nil {
						rscData.FamilyInet.DHCP = &interfaceLogicalBlockFamilyInetBlockDhcp{}
						if strings.HasPrefix(itemTrim, " dhcp-client") {
							rscData.FamilyInet.DHCP.SrxOldOptionName = types.BoolValue(true)
						}
					}
					if balt.CutPrefixInString(&itemTrim, " dhcp ") || balt.CutPrefixInString(&itemTrim, " dhcp-client ") {
						if err := rscData.FamilyInet.DHCP.read(itemTrim); err != nil {
							return err
						}
					}
				case balt.CutPrefixInString(&itemTrim, " filter input "):
					rscData.FamilyInet.FilterInput = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, " filter output "):
					rscData.FamilyInet.FilterOutput = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, " mtu "):
					rscData.FamilyInet.Mtu, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " rpf-check"):
					if rscData.FamilyInet.RPFCheck == nil {
						rscData.FamilyInet.RPFCheck = &interfaceLogicalBlockFamilyBlockRPFCheck{}
					}
					switch {
					case balt.CutPrefixInString(&itemTrim, " fail-filter "):
						rscData.FamilyInet.RPFCheck.FailFilter = types.StringValue(strings.Trim(itemTrim, "\""))
					case itemTrim == " mode loose":
						rscData.FamilyInet.RPFCheck.ModeLoose = types.BoolValue(true)
					}
				case itemTrim == " sampling input":
					rscData.FamilyInet.SamplingInput = types.BoolValue(true)
				case itemTrim == " sampling output":
					rscData.FamilyInet.SamplingOutput = types.BoolValue(true)
				}
			case balt.CutPrefixInString(&itemTrim, "tunnel "):
				if rscData.Tunnel == nil {
					rscData.Tunnel = &interfaceLogicalBlockTunnel{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "destination "):
					rscData.Tunnel.Destination = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "source "):
					rscData.Tunnel.Source = types.StringValue(itemTrim)
				case itemTrim == "allow-fragmentation":
					rscData.Tunnel.AllowFragmentation = types.BoolValue(true)
				case itemTrim == "do-not-fragment":
					rscData.Tunnel.DoNotFragment = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "flow-label "):
					rscData.Tunnel.FlowLabel, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case itemTrim == "no-path-mtu-discovery":
					rscData.Tunnel.NoPathMtuDiscovery = types.BoolValue(true)
				case itemTrim == "path-mtu-discovery":
					rscData.Tunnel.PathMtuDiscovery = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "routing-instance destination "):
					rscData.Tunnel.RoutingInstanceDestination = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "traffic-class "):
					rscData.Tunnel.TrafficClass, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "ttl "):
					rscData.Tunnel.TTL, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "vlan-id "):
				rscData.VlanID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			default:
				continue
			}
		}
	}
	showConfigRoutingInstances, err := junSess.Command(junos.CmdShowConfig +
		"routing-instances" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	regexpInt := regexp.MustCompile(`set \S+ interface ` + name + `$`)
	for _, item := range strings.Split(showConfigRoutingInstances, "\n") {
		intMatch := regexpInt.MatchString(item)
		if intMatch {
			rscData.RoutingInstance = types.StringValue(
				strings.TrimPrefix(
					strings.TrimSuffix(
						item, " interface "+name,
					),
					junos.SetLS,
				),
			)

			break
		}
	}
	if junSess.CheckCompatibilitySecurity() {
		showConfigSecurityZones, err := junSess.Command(junos.CmdShowConfig + "security zones" + junos.PipeDisplaySetRelative)
		if err != nil {
			return err
		}
		regexpInts := regexp.MustCompile(`set security-zone \S+ interfaces ` + name + `( host-inbound-traffic .*)?$`)
		for _, item := range strings.Split(showConfigSecurityZones, "\n") {
			intMatch := regexpInts.MatchString(item)
			if intMatch {
				itemTrimFields := strings.Split(strings.TrimPrefix(item, "set security-zone "), " ")
				rscData.SecurityZone = types.StringValue(itemTrimFields[0])
				if err := rscData.readSecurityInboundTraffic(name, junSess); err != nil {
					return err
				}

				break
			}
		}
	}

	return nil
}

func (block *interfaceLogicalBlockFamilyInetBlockAddress) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "primary":
		block.Primary = types.BoolValue(true)
	case itemTrim == "preferred":
		block.Preferred = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "vrrp-group "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var vrrpGroup interfaceLogicalBlockFamilyInetBlockAddressBlockVRRPGroup
		vrrpGroupIdentifier, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
		if err != nil {
			return err
		}
		block.VRRPGroup, vrrpGroup = tfdata.ExtractBlockWithTFTypesInt64(
			block.VRRPGroup, "Identifier", vrrpGroupIdentifier.ValueInt64(),
		)
		vrrpGroup.Identifier = vrrpGroupIdentifier
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")

		switch {
		case balt.CutPrefixInString(&itemTrim, "virtual-address "):
			vrrpGroup.VirtualAddress = append(vrrpGroup.VirtualAddress, types.StringValue(itemTrim))
		case itemTrim == "accept-data":
			vrrpGroup.AcceptData = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "advertise-interval "):
			vrrpGroup.AdvertiseInterval, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "advertisements-threshold "):
			vrrpGroup.AdvertisementsThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "authentication-key "):
			vrrpGroup.AuthenticationKey, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "authentication-key")
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "authentication-type "):
			vrrpGroup.AuthenticationType = types.StringValue(itemTrim)
		case itemTrim == "no-accept-data":
			vrrpGroup.NoAcceptData = types.BoolValue(true)
		case itemTrim == "no-preempt":
			vrrpGroup.NoPreempt = types.BoolValue(true)
		case itemTrim == "preempt":
			vrrpGroup.Preempt = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "priority "):
			vrrpGroup.Priority, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "track interface "):
			itemTrackFields := strings.Split(itemTrim, " ")
			if len(itemTrackFields) < 3 { // <interface> priority-cost <priority_cost>
				return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "track interface", itemTrim)
			}
			cost, err := tfdata.ConvAtoi64Value(itemTrackFields[2])
			if err != nil {
				return err
			}
			vrrpGroup.TrackInterface = append(vrrpGroup.TrackInterface,
				interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackInterface{
					Interface:    types.StringValue(itemTrackFields[0]),
					PriorityCost: cost,
				},
			)
		case balt.CutPrefixInString(&itemTrim, "track route "):
			itemTrackFields := strings.Split(itemTrim, " ")
			if len(itemTrackFields) < 5 { // <route> routing-instance <routing_instance> priority-cost <priority_cost>
				return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "track route", itemTrim)
			}
			cost, err := tfdata.ConvAtoi64Value(itemTrackFields[4])
			if err != nil {
				return err
			}
			vrrpGroup.TrackRoute = append(vrrpGroup.TrackRoute,
				interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackRoute{
					Route:           types.StringValue(itemTrackFields[0]),
					RoutingInstance: types.StringValue(itemTrackFields[2]),
					PriorityCost:    cost,
				},
			)
		}
		block.VRRPGroup = append(block.VRRPGroup, vrrpGroup)
	}

	return nil
}

func (block *interfaceLogicalBlockFamilyInetBlockDhcp) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "client-identifier ascii "):
		block.ClientIdentifierASCII = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "client-identifier hexadecimal "):
		block.ClientIdentifierHexadecimal = types.StringValue(itemTrim)
	case itemTrim == "client-identifier prefix host-name":
		block.ClientIdentifierPrefixHostname = types.BoolValue(true)
	case itemTrim == "client-identifier prefix routing-instance-name":
		block.ClientIdentifierPrefixRoutingInstanceName = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "client-identifier use-interface-description "):
		block.ClientIdentifierUseInterfaceDescription = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "client-identifier user-id ascii "):
		block.ClientIdentifierUseridASCII = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "client-identifier user-id hexadecimal "):
		block.ClientIdentifierUseridHexadecimal = types.StringValue(itemTrim)
	case itemTrim == "force-discover":
		block.ForceDiscover = types.BoolValue(true)
	case itemTrim == "lease-time infinite":
		block.LeaseTimeInfinite = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "lease-time "):
		block.LeaseTime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "metric "):
		block.Metric, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "no-dns-install":
		block.NoDNSInstall = types.BoolValue(true)
	case itemTrim == "options no-hostname":
		block.OptionsNoHostname = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "retransmission-attempt "):
		block.RetransmissionAttempt, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "retransmission-interval "):
		block.RetransmissionInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "server-address "):
		block.ServerAddress = types.StringValue(itemTrim)
	case itemTrim == "update-server":
		block.UpdateServer = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "vendor-id "):
		block.VendorID = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return nil
}

func (block *interfaceLogicalBlockFamilyInet6BlockAddress) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "primary":
		block.Primary = types.BoolValue(true)
	case itemTrim == "preferred":
		block.Preferred = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "vrrp-inet6-group "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var vrrpGroup interfaceLogicalBlockFamilyInet6BlockAddressBlockVRRPGroup
		vrrpGroupIdentifier, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
		if err != nil {
			return err
		}
		block.VRRPGroup, vrrpGroup = tfdata.ExtractBlockWithTFTypesInt64(
			block.VRRPGroup, "Identifier", vrrpGroupIdentifier.ValueInt64(),
		)
		vrrpGroup.Identifier = vrrpGroupIdentifier
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")

		switch {
		case balt.CutPrefixInString(&itemTrim, "virtual-inet6-address "):
			vrrpGroup.VirtualAddress = append(vrrpGroup.VirtualAddress, types.StringValue(itemTrim))
		case balt.CutPrefixInString(&itemTrim, "virtual-link-local-address "):
			vrrpGroup.VirutalLinkLocalAddress = types.StringValue(itemTrim)
		case itemTrim == "accept-data":
			vrrpGroup.AcceptData = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "inet6-advertise-interval "):
			vrrpGroup.AdvertiseInterval, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "advertisements-threshold "):
			vrrpGroup.AdvertisementsThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case itemTrim == "no-accept-data":
			vrrpGroup.NoAcceptData = types.BoolValue(true)
		case itemTrim == "no-preempt":
			vrrpGroup.NoPreempt = types.BoolValue(true)
		case itemTrim == "preempt":
			vrrpGroup.Preempt = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "priority "):
			vrrpGroup.Priority, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "track interface "):
			itemTrackFields := strings.Split(itemTrim, " ")
			if len(itemTrackFields) < 3 { // <interface> priority-cost <priority_cost>
				return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "track interface", itemTrim)
			}
			cost, err := tfdata.ConvAtoi64Value(itemTrackFields[2])
			if err != nil {
				return err
			}
			vrrpGroup.TrackInterface = append(vrrpGroup.TrackInterface,
				interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackInterface{
					Interface:    types.StringValue(itemTrackFields[0]),
					PriorityCost: cost,
				},
			)
		case balt.CutPrefixInString(&itemTrim, "track route "):
			itemTrackFields := strings.Split(itemTrim, " ")
			if len(itemTrackFields) < 5 { // <route> routing-instance <routing_instance> priority-cost <priority_cost>
				return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "track route", itemTrim)
			}
			cost, err := tfdata.ConvAtoi64Value(itemTrackFields[4])
			if err != nil {
				return err
			}
			vrrpGroup.TrackRoute = append(vrrpGroup.TrackRoute,
				interfaceLogicalBlockFamilyBlockAddressBlockVRRPGroupBlockTrackRoute{
					Route:           types.StringValue(itemTrackFields[0]),
					RoutingInstance: types.StringValue(itemTrackFields[2]),
					PriorityCost:    cost,
				},
			)
		}
		block.VRRPGroup = append(block.VRRPGroup, vrrpGroup)
	}

	return nil
}

func (block *interfaceLogicalBlockFamilyInet6BlockDhcpV6Client) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "client-identifier duid-type "):
		block.ClientIdentifierDuidType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "client-type "):
		block.ClientType = types.StringValue(itemTrim)
	case itemTrim == "client-ia-type ia-na":
		block.ClientIATypeNA = types.BoolValue(true)
	case itemTrim == "client-ia-type ia-pd":
		block.ClientIATypePD = types.BoolValue(true)
	case itemTrim == "no-dns-install":
		block.NoDNSInstall = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "prefix-delegating preferred-prefix-length "):
		block.PrefixDelegatingPreferredPrefixLength, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "prefix-delegating sub-prefix-length "):
		block.PrefixDelegatingSubPrefixLength, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "rapid-commit":
		block.RapidCommit = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "req-option "):
		block.ReqOption = append(block.ReqOption, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "retransmission-attempt "):
		block.RetransmissionAttempt, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "update-router-advertisement interface "):
		block.UpdateRouterAdvertisementInterface = append(
			block.UpdateRouterAdvertisementInterface,
			types.StringValue(itemTrim),
		)
	case itemTrim == "update-server":
		block.UpdateServer = types.BoolValue(true)
	}

	return nil
}

func (rscData *interfaceLogicalData) readSecurityInboundTraffic(
	name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + rscData.SecurityZone.ValueString() +
		" interfaces " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}

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
			case balt.CutPrefixInString(&itemTrim, "host-inbound-traffic protocols "):
				rscData.SecurityInboundProtocols = append(rscData.SecurityInboundProtocols, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "host-inbound-traffic system-services "):
				rscData.SecurityInboundServices = append(rscData.SecurityInboundServices, types.StringValue(itemTrim))
			}
		}
	}

	return nil
}

func (rscData *interfaceLogicalData) del(ctx context.Context, junSess *junos.Session) error {
	configSet := []string{
		"delete interfaces " + rscData.Name.ValueString(),
	}
	if strings.HasPrefix(rscData.Name.ValueString(), "st0.") && !rscData.St0AlsoOnDestroy.ValueBool() {
		// interface totally delete by
		// - junos_interface_st0_unit resource
		// else there is an interface st0.x empty
		configSet = append(configSet,
			"set interfaces "+rscData.Name.ValueString(),
		)
	}
	if err := junSess.ConfigSet(configSet); err != nil {
		return err
	}
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		if err := rscData.delRoutingInstance(ctx, junSess); err != nil {
			return err
		}
	}
	if v := rscData.SecurityZone.ValueString(); v != "" {
		if err := rscData.delZone(ctx, junSess); err != nil {
			return err
		}
	}

	return nil
}

func (rscData *interfaceLogicalData) delOpts(_ context.Context, junSess *junos.Session) error {
	delPrefix := "delete interfaces " + rscData.Name.ValueString() + " "
	configSet := []string{
		delPrefix + "description",
		delPrefix + "disable",
		delPrefix + "family inet",
		delPrefix + "family inet6",
		delPrefix + "tunnel",
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *interfaceLogicalData) delZone(_ context.Context, junSess *junos.Session) error {
	configSet := []string{
		"delete security zones security-zone " + rscData.SecurityZone.ValueString() +
			" interfaces " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *interfaceLogicalData) delRoutingInstance(_ context.Context, junSess *junos.Session) error {
	configSet := []string{
		junos.DelRoutingInstances + rscData.RoutingInstance.ValueString() +
			" interface " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
