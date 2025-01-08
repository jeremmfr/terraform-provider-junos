package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &accessAddressAssignmentPool{}
	_ resource.ResourceWithConfigure      = &accessAddressAssignmentPool{}
	_ resource.ResourceWithValidateConfig = &accessAddressAssignmentPool{}
	_ resource.ResourceWithImportState    = &accessAddressAssignmentPool{}
	_ resource.ResourceWithUpgradeState   = &accessAddressAssignmentPool{}
)

type accessAddressAssignmentPool struct {
	client *junos.Client
}

func newAccessAddressAssignmentPoolResource() resource.Resource {
	return &accessAddressAssignmentPool{}
}

func (rsc *accessAddressAssignmentPool) typeName() string {
	return providerName + "_access_address_assignment_pool"
}

func (rsc *accessAddressAssignmentPool) junosName() string {
	return "access address-assignment pool"
}

func (rsc *accessAddressAssignmentPool) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *accessAddressAssignmentPool) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *accessAddressAssignmentPool) Configure(
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

func (rsc *accessAddressAssignmentPool) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<name>" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Address pool name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"active_drain": schema.BoolAttribute{
				Optional:    true,
				Description: "Notify client of pool active drain mode.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"hold_down": schema.BoolAttribute{
				Optional:    true,
				Description: "Place pool in passive drain mode.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"link": schema.StringAttribute{
				Optional:    true,
				Description: "Address pool link name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"family": schema.SingleNestedBlock{
				Description: "Configure address family (`inet` or `inet6`).",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Type of family.",
						Validators: []validator.String{
							stringvalidator.OneOf("inet", "inet6"),
						},
					},
					"network": schema.StringAttribute{
						Required:    true,
						Description: "Network address of pool.",
						Validators: []validator.String{
							tfvalidator.StringCIDRNetwork(),
						},
					},
					"excluded_address": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Excluded Addresses.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.NoNullValues(),
							setvalidator.ValueStringsAre(
								tfvalidator.StringIPAddress(),
							),
						},
					},
					"xauth_attributes_primary_dns": schema.StringAttribute{
						Optional:    true,
						Description: "Specify the primary-dns IP address.",
						Validators: []validator.String{
							stringvalidator.Any(
								tfvalidator.StringCIDR().IPv4Only(),
								tfvalidator.StringIPAddress().IPv6Only(),
							),
						},
					},
					"xauth_attributes_primary_wins": schema.StringAttribute{
						Optional:    true,
						Description: "Specify the primary-wins IP address.",
						Validators: []validator.String{
							tfvalidator.StringCIDR().IPv4Only(),
						},
					},
					"xauth_attributes_secondary_dns": schema.StringAttribute{
						Optional:    true,
						Description: "Specify the secondary-dns IP address.",
						Validators: []validator.String{
							stringvalidator.Any(
								tfvalidator.StringCIDR().IPv4Only(),
								tfvalidator.StringIPAddress().IPv6Only(),
							),
						},
					},
					"xauth_attributes_secondary_wins": schema.StringAttribute{
						Optional:    true,
						Description: "Specify the secondary-wins IP address.",
						Validators: []validator.String{
							tfvalidator.StringCIDR().IPv4Only(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"dhcp_attributes": schema.SingleNestedBlock{
						Description: "DHCP options and match criteria.",
						Attributes: map[string]schema.Attribute{
							"boot_file": schema.StringAttribute{
								Optional:    true,
								Description: "Boot filename advertised to clients.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"boot_server": schema.StringAttribute{
								Optional:    true,
								Description: "Boot server advertised to clients.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
								},
							},
							"dns_server": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "IPv6 domain name servers available to the client.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.NoNullValues(),
									listvalidator.ValueStringsAre(
										tfvalidator.StringIPAddress().IPv6Only(),
									),
								},
							},
							"domain_name": schema.StringAttribute{
								Optional:    true,
								Description: "Domain name advertised to clients.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
								},
							},
							"exclude_prefix_len": schema.Int64Attribute{
								Optional:    true,
								Description: "Length of IPv6 prefix to be excluded from delegated prefix.",
								Validators: []validator.Int64{
									int64validator.Between(1, 128),
								},
							},
							"grace_period": schema.Int64Attribute{
								Optional:    true,
								Description: "Grace period for leases (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"maximum_lease_time": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum lease time advertised to clients (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"maximum_lease_time_infinite": schema.BoolAttribute{
								Optional:    true,
								Description: "Lease time can be infinite.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"name_server": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "IPv4 domain name servers available to the client.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.NoNullValues(),
									listvalidator.ValueStringsAre(
										tfvalidator.StringIPAddress().IPv4Only(),
									),
								},
							},
							"netbios_node_type": schema.StringAttribute{
								Optional:    true,
								Description: "Type of NETBIOS node advertised to clients.",
								Validators: []validator.String{
									stringvalidator.OneOf("b-node", "h-node", "m-node", "p-node"),
								},
							},
							"next_server": schema.StringAttribute{
								Optional:    true,
								Description: "Next server that clients need to contact.",
								Validators: []validator.String{
									tfvalidator.StringIPAddress().IPv4Only(),
								},
							},
							"option": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "DHCP option.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.NoNullValues(),
									setvalidator.ValueStringsAre(
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^\d+ (array )?`+
												`(byte|flag|hex-string|integer|ip-address|short|string|unsigned-integer|unsigned-short) .*$`),
											`need to match '^\d+ (array )?"+
												"(byte|flag|hex-string|integer|ip-address|short|string|unsigned-integer|unsigned-short) .*$'`,
										),
									),
								},
							},
							"preferred_lifetime": schema.Int64Attribute{
								Optional:    true,
								Description: "Preferred lifetime advertised to clients (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"preferred_lifetime_infinite": schema.BoolAttribute{
								Optional:    true,
								Description: "Lease time can be infinite.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"propagate_ppp_settings": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "PPP interface name for propagating DNS/WINS settings.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.NoNullValues(),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
									),
								},
							},
							"propagate_settings": schema.StringAttribute{
								Optional:    true,
								Description: "Interface name for propagating TCP/IP Settings to pool.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"router": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Routers advertised to clients.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.NoNullValues(),
									listvalidator.ValueStringsAre(
										tfvalidator.StringIPAddress().IPv4Only(),
									),
								},
							},
							"server_identifier": schema.StringAttribute{
								Optional:    true,
								Description: "Server Identifier - IP address value.",
								Validators: []validator.String{
									tfvalidator.StringIPAddress().IPv4Only(),
								},
							},
							"sip_server_inet_address": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "SIP servers list of IPv4 addresses available to the client.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.NoNullValues(),
									listvalidator.ValueStringsAre(
										tfvalidator.StringIPAddress().IPv4Only(),
									),
								},
							},
							"sip_server_inet_domain_name": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "SIP server domain name available to clients.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.NoNullValues(),
									listvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									),
								},
							},
							"sip_server_inet6_address": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "SIP Servers list of IPv6 addresses available to the client.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.NoNullValues(),
									listvalidator.ValueStringsAre(
										tfvalidator.StringIPAddress().IPv6Only(),
									),
								},
							},
							"sip_server_inet6_domain_name": schema.StringAttribute{
								Optional:    true,
								Description: "SIP server domain name available to clients.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
								},
							},
							"t1_percentage": schema.Int64Attribute{
								Optional:    true,
								Description: "T1 time as percentage of preferred lifetime or max lease (percent)",
								Validators: []validator.Int64{
									int64validator.Between(0, 100),
								},
							},
							"t1_renewal_time": schema.Int64Attribute{
								Optional:    true,
								Description: "T1 renewal time (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"t2_percentage": schema.Int64Attribute{
								Optional:    true,
								Description: "T2 time as percentage of preferred lifetime or max lease (percent).",
								Validators: []validator.Int64{
									int64validator.Between(0, 100),
								},
							},
							"t2_rebinding_time": schema.Int64Attribute{
								Optional:    true,
								Description: "T2 rebinding time (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"tftp_server": schema.StringAttribute{
								Optional:    true,
								Description: "TFTP server IP address advertised to clients.",
								Validators: []validator.String{
									tfvalidator.StringIPAddress().IPv4Only(),
								},
							},
							"valid_lifetime": schema.Int64Attribute{
								Optional:    true,
								Description: "Valid lifetime advertised to clients (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"valid_lifetime_infinite": schema.BoolAttribute{
								Optional:    true,
								Description: "Lease time can be infinite.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"wins_server": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "WINS name servers.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.NoNullValues(),
									listvalidator.ValueStringsAre(
										tfvalidator.StringIPAddress().IPv4Only(),
									),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"option_match_82_circuit_id": schema.ListNestedBlock{
								Description: "Circuit ID portion of the option 82.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"value": schema.StringAttribute{
											Required:    true,
											Description: "Match value.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
												tfvalidator.StringDoubleQuoteExclusion(),
											},
										},
										"range": schema.StringAttribute{
											Required:    true,
											Description: "Range name.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
												tfvalidator.StringDoubleQuoteExclusion(),
											},
										},
									},
								},
							},
							"option_match_82_remote_id": schema.ListNestedBlock{
								Description: "Remote ID portion of the option 82.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"value": schema.StringAttribute{
											Required:    true,
											Description: "Match value.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
												tfvalidator.StringDoubleQuoteExclusion(),
											},
										},
										"range": schema.StringAttribute{
											Required:    true,
											Description: "Range name.",
											Validators: []validator.String{
												stringvalidator.LengthAtLeast(1),
												tfvalidator.StringDoubleQuoteExclusion(),
											},
										},
									},
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"excluded_range": schema.ListNestedBlock{
						Description: "For each name of excluded address range to declare.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "Range name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 63),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									},
								},
								"low": schema.StringAttribute{
									Required:    true,
									Description: "Lower limit of excluded address range.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress(),
									},
								},
								"high": schema.StringAttribute{
									Required:    true,
									Description: "Upper limit of excluded address range.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress(),
									},
								},
							},
						},
					},
					"host": schema.ListNestedBlock{
						Description: "For each name of host to declare.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "Hostname for static reservations.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"ip_address": schema.StringAttribute{
									Required:    true,
									Description: "Reserved address.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress(),
									},
								},
								"hardware_address": schema.StringAttribute{
									Optional:    true,
									Description: "Hardware address.",
									Validators: []validator.String{
										tfvalidator.StringMACAddress(),
									},
								},
								"user_name": schema.BoolAttribute{
									Optional:    true,
									Description: "Set subscriber user name as host identifier.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
							},
						},
					},
					"inet_range": schema.ListNestedBlock{
						Description: "For each name of address range to declare.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "Range name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 63),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									},
								},
								"low": schema.StringAttribute{
									Required:    true,
									Description: "Lower limit of address range.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv4Only(),
									},
								},
								"high": schema.StringAttribute{
									Required:    true,
									Description: "Upper limit of address range.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv4Only(),
									},
								},
							},
						},
					},
					"inet6_range": schema.ListNestedBlock{
						Description: "For each name of address range to declare.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "Range name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 63),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									},
								},
								"low": schema.StringAttribute{
									Optional:    true,
									Description: "Lower limit of IPv6 address range.",
									Validators: []validator.String{
										tfvalidator.StringCIDR().IPv6Only(),
									},
								},
								"high": schema.StringAttribute{
									Optional:    true,
									Description: "Upper limit of IPv6 address range.",
									Validators: []validator.String{
										tfvalidator.StringCIDR().IPv6Only(),
									},
								},
								"prefix_length": schema.Int64Attribute{
									Optional:    true,
									Description: "IPv6 delegated prefix length.",
									Validators: []validator.Int64{
										int64validator.Between(1, 128),
									},
								},
							},
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

type accessAddressAssignmentPoolData struct {
	ID              types.String                            `tfsdk:"id"`
	Name            types.String                            `tfsdk:"name"`
	RoutingInstance types.String                            `tfsdk:"routing_instance"`
	ActiveDrain     types.Bool                              `tfsdk:"active_drain"`
	HoldDown        types.Bool                              `tfsdk:"hold_down"`
	Link            types.String                            `tfsdk:"link"`
	Family          *accessAddressAssignmentPoolBlockFamily `tfsdk:"family"`
}

type accessAddressAssignmentPoolConfig struct {
	ID              types.String                                  `tfsdk:"id"`
	Name            types.String                                  `tfsdk:"name"`
	RoutingInstance types.String                                  `tfsdk:"routing_instance"`
	ActiveDrain     types.Bool                                    `tfsdk:"active_drain"`
	HoldDown        types.Bool                                    `tfsdk:"hold_down"`
	Link            types.String                                  `tfsdk:"link"`
	Family          *accessAddressAssignmentPoolBlockFamilyConfig `tfsdk:"family"`
}

//nolint:lll
type accessAddressAssignmentPoolBlockFamily struct {
	Type                         types.String                                               `tfsdk:"type"`
	Network                      types.String                                               `tfsdk:"network"`
	ExcludedAddress              []types.String                                             `tfsdk:"excluded_address"`
	XauthAttributesPrimaryDNS    types.String                                               `tfsdk:"xauth_attributes_primary_dns"`
	XauthAttributesPrimaryWins   types.String                                               `tfsdk:"xauth_attributes_primary_wins"`
	XauthAttributesSecondaryDNS  types.String                                               `tfsdk:"xauth_attributes_secondary_dns"`
	XauthAttributesSecondaryWins types.String                                               `tfsdk:"xauth_attributes_secondary_wins"`
	DhcpAttributes               *accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributes `tfsdk:"dhcp_attributes"`
	ExcludedRange                []accessAddressAssignmentPoolBlockFamilyBlockExcludedRange `tfsdk:"excluded_range"`
	Host                         []accessAddressAssignmentPoolBlockFamilyBlockHost          `tfsdk:"host"`
	InetRange                    []accessAddressAssignmentPoolBlockFamilyBlockInetRange     `tfsdk:"inet_range"`
	Inet6Range                   []accessAddressAssignmentPoolBlockFamilyBlockInet6Range    `tfsdk:"inet6_range"`
}

//nolint:lll
type accessAddressAssignmentPoolBlockFamilyConfig struct {
	Type                         types.String                                                     `tfsdk:"type"`
	Network                      types.String                                                     `tfsdk:"network"`
	ExcludedAddress              types.Set                                                        `tfsdk:"excluded_address"`
	XauthAttributesPrimaryDNS    types.String                                                     `tfsdk:"xauth_attributes_primary_dns"`
	XauthAttributesPrimaryWins   types.String                                                     `tfsdk:"xauth_attributes_primary_wins"`
	XauthAttributesSecondaryDNS  types.String                                                     `tfsdk:"xauth_attributes_secondary_dns"`
	XauthAttributesSecondaryWins types.String                                                     `tfsdk:"xauth_attributes_secondary_wins"`
	DhcpAttributes               *accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesConfig `tfsdk:"dhcp_attributes"`
	ExcludedRange                types.List                                                       `tfsdk:"excluded_range"`
	Host                         types.List                                                       `tfsdk:"host"`
	InetRange                    types.List                                                       `tfsdk:"inet_range"`
	Inet6Range                   types.List                                                       `tfsdk:"inet6_range"`
}

//nolint:lll
type accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributes struct {
	BootFile                  types.String                                                                  `tfsdk:"boot_file"`
	BootServer                types.String                                                                  `tfsdk:"boot_server"`
	DNSServer                 []types.String                                                                `tfsdk:"dns_server"`
	DomainName                types.String                                                                  `tfsdk:"domain_name"`
	ExcludePrefixLen          types.Int64                                                                   `tfsdk:"exclude_prefix_len"`
	GracePeriod               types.Int64                                                                   `tfsdk:"grace_period"`
	MaximumLeaseTime          types.Int64                                                                   `tfsdk:"maximum_lease_time"`
	MaximumLeaseTimeInfinite  types.Bool                                                                    `tfsdk:"maximum_lease_time_infinite"`
	NameServer                []types.String                                                                `tfsdk:"name_server"`
	NetbiosNodeType           types.String                                                                  `tfsdk:"netbios_node_type"`
	NextServer                types.String                                                                  `tfsdk:"next_server"`
	Option                    []types.String                                                                `tfsdk:"option"`
	PreferredLifetime         types.Int64                                                                   `tfsdk:"preferred_lifetime"`
	PreferredLifetimeInfinite types.Bool                                                                    `tfsdk:"preferred_lifetime_infinite"`
	PropagatePppSettings      []types.String                                                                `tfsdk:"propagate_ppp_settings"`
	PropagateSettings         types.String                                                                  `tfsdk:"propagate_settings"`
	Router                    []types.String                                                                `tfsdk:"router"`
	ServerIdentifier          types.String                                                                  `tfsdk:"server_identifier"`
	SIPServerInetAddress      []types.String                                                                `tfsdk:"sip_server_inet_address"`
	SIPServerInetDomainName   []types.String                                                                `tfsdk:"sip_server_inet_domain_name"`
	SIPServerInet6Address     []types.String                                                                `tfsdk:"sip_server_inet6_address"`
	SIPServerInet6DomainName  types.String                                                                  `tfsdk:"sip_server_inet6_domain_name"`
	T1Percentage              types.Int64                                                                   `tfsdk:"t1_percentage"`
	T1RenewalTime             types.Int64                                                                   `tfsdk:"t1_renewal_time"`
	T2Percentage              types.Int64                                                                   `tfsdk:"t2_percentage"`
	T2RebindingTime           types.Int64                                                                   `tfsdk:"t2_rebinding_time"`
	TftpServer                types.String                                                                  `tfsdk:"tftp_server"`
	ValidLifetime             types.Int64                                                                   `tfsdk:"valid_lifetime"`
	ValidLifetimeInfinite     types.Bool                                                                    `tfsdk:"valid_lifetime_infinite"`
	WinsServer                []types.String                                                                `tfsdk:"wins_server"`
	OptionMatch82CircuitID    []accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesBlockOptionMatch82 `tfsdk:"option_match_82_circuit_id"`
	OptionMatch82RemoteID     []accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesBlockOptionMatch82 `tfsdk:"option_match_82_remote_id"`
}

func (block *accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributes) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesConfig struct {
	BootFile                  types.String `tfsdk:"boot_file"`
	BootServer                types.String `tfsdk:"boot_server"`
	DNSServer                 types.List   `tfsdk:"dns_server"`
	DomainName                types.String `tfsdk:"domain_name"`
	ExcludePrefixLen          types.Int64  `tfsdk:"exclude_prefix_len"`
	GracePeriod               types.Int64  `tfsdk:"grace_period"`
	MaximumLeaseTime          types.Int64  `tfsdk:"maximum_lease_time"`
	MaximumLeaseTimeInfinite  types.Bool   `tfsdk:"maximum_lease_time_infinite"`
	NameServer                types.List   `tfsdk:"name_server"`
	NetbiosNodeType           types.String `tfsdk:"netbios_node_type"`
	NextServer                types.String `tfsdk:"next_server"`
	Option                    types.Set    `tfsdk:"option"`
	PreferredLifetime         types.Int64  `tfsdk:"preferred_lifetime"`
	PreferredLifetimeInfinite types.Bool   `tfsdk:"preferred_lifetime_infinite"`
	PropagatePppSettings      types.Set    `tfsdk:"propagate_ppp_settings"`
	PropagateSettings         types.String `tfsdk:"propagate_settings"`
	Router                    types.List   `tfsdk:"router"`
	ServerIdentifier          types.String `tfsdk:"server_identifier"`
	SIPServerInetAddress      types.List   `tfsdk:"sip_server_inet_address"`
	SIPServerInetDomainName   types.List   `tfsdk:"sip_server_inet_domain_name"`
	SIPServerInet6Address     types.List   `tfsdk:"sip_server_inet6_address"`
	SIPServerInet6DomainName  types.String `tfsdk:"sip_server_inet6_domain_name"`
	T1Percentage              types.Int64  `tfsdk:"t1_percentage"`
	T1RenewalTime             types.Int64  `tfsdk:"t1_renewal_time"`
	T2Percentage              types.Int64  `tfsdk:"t2_percentage"`
	T2RebindingTime           types.Int64  `tfsdk:"t2_rebinding_time"`
	TftpServer                types.String `tfsdk:"tftp_server"`
	ValidLifetime             types.Int64  `tfsdk:"valid_lifetime"`
	ValidLifetimeInfinite     types.Bool   `tfsdk:"valid_lifetime_infinite"`
	WinsServer                types.List   `tfsdk:"wins_server"`
	OptionMatch82CircuitID    types.List   `tfsdk:"option_match_82_circuit_id"`
	OptionMatch82RemoteID     types.List   `tfsdk:"option_match_82_remote_id"`
}

func (block *accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesBlockOptionMatch82 struct {
	Value types.String `tfsdk:"value"`
	Range types.String `tfsdk:"range"`
}

type accessAddressAssignmentPoolBlockFamilyBlockExcludedRange struct {
	Name types.String `tfsdk:"name" tfdata:"identifier"`
	Low  types.String `tfsdk:"low"`
	High types.String `tfsdk:"high"`
}

type accessAddressAssignmentPoolBlockFamilyBlockHost struct {
	Name            types.String `tfsdk:"name"             tfdata:"identifier"`
	HardwareAddress types.String `tfsdk:"hardware_address"`
	IPAddress       types.String `tfsdk:"ip_address"`
	UserName        types.Bool   `tfsdk:"user_name"`
}

type accessAddressAssignmentPoolBlockFamilyBlockInetRange struct {
	Name types.String `tfsdk:"name" tfdata:"identifier"`
	Low  types.String `tfsdk:"low"`
	High types.String `tfsdk:"high"`
}

type accessAddressAssignmentPoolBlockFamilyBlockInet6Range struct {
	Name         types.String `tfsdk:"name"          tfdata:"identifier"`
	Low          types.String `tfsdk:"low"`
	High         types.String `tfsdk:"high"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
}

func (rsc *accessAddressAssignmentPool) ValidateConfig( //nolint:gocognit,gocyclo
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config accessAddressAssignmentPoolConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Family != nil {
		if !config.Family.Type.IsNull() &&
			!config.Family.Type.IsUnknown() {
			switch config.Family.Type.ValueString() {
			case junos.InetW:
				if !config.Family.Network.IsNull() &&
					!config.Family.Network.IsUnknown() {
					ipV4Validator := tfvalidator.StringCIDR().IPv4Only()
					attrReq := validator.StringRequest{
						Path:        path.Root("family").AtName("network"),
						ConfigValue: config.Family.Network,
					}
					attrResp := validator.StringResponse{}
					ipV4Validator.ValidateString(ctx, attrReq, &attrResp)
					resp.Diagnostics.Append(attrResp.Diagnostics.Errors()...)
				}
				if !config.Family.ExcludedAddress.IsNull() &&
					!config.Family.ExcludedAddress.IsUnknown() {
					ipV4Validator := tfvalidator.StringIPAddress().IPv4Only()
					var excludedAddress []types.String
					asDiags := config.Family.ExcludedAddress.ElementsAs(ctx, &excludedAddress, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}
					for _, v := range excludedAddress {
						attrReq := validator.StringRequest{
							Path:        path.Root("family").AtName("excluded_address"),
							ConfigValue: v,
						}
						attrResp := validator.StringResponse{}
						ipV4Validator.ValidateString(ctx, attrReq, &attrResp)
						resp.Diagnostics.Append(attrResp.Diagnostics.Errors()...)
					}
				}
				if !config.Family.XauthAttributesPrimaryDNS.IsNull() &&
					!config.Family.XauthAttributesPrimaryDNS.IsUnknown() {
					ipV4Validator := tfvalidator.StringCIDR().IPv4Only()
					attrReq := validator.StringRequest{
						Path:        path.Root("family").AtName("xauth_attributes_primary_dns"),
						ConfigValue: config.Family.XauthAttributesPrimaryDNS,
					}
					attrResp := validator.StringResponse{}
					ipV4Validator.ValidateString(ctx, attrReq, &attrResp)
					resp.Diagnostics.Append(attrResp.Diagnostics.Errors()...)
				}
				if !config.Family.XauthAttributesSecondaryDNS.IsNull() &&
					!config.Family.XauthAttributesSecondaryDNS.IsUnknown() {
					ipV4Validator := tfvalidator.StringCIDR().IPv4Only()
					attrReq := validator.StringRequest{
						Path:        path.Root("family").AtName("xauth_attributes_secondary_dns"),
						ConfigValue: config.Family.XauthAttributesPrimaryDNS,
					}
					attrResp := validator.StringResponse{}
					ipV4Validator.ValidateString(ctx, attrReq, &attrResp)
					resp.Diagnostics.Append(attrResp.Diagnostics.Errors()...)
				}
				if config.Family.DhcpAttributes != nil {
					if !config.Family.DhcpAttributes.DNSServer.IsNull() &&
						!config.Family.DhcpAttributes.DNSServer.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("dhcp_attributes").AtName("dns_server"),
							tfdiag.ConflictConfigErrSummary,
							"dns_server cannot be configured when type = inet"+
								" in dhcp_attributes block in family block",
						)
					}
					if !config.Family.DhcpAttributes.ExcludePrefixLen.IsNull() &&
						!config.Family.DhcpAttributes.ExcludePrefixLen.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("dhcp_attributes").AtName("exclude_prefix_len"),
							tfdiag.ConflictConfigErrSummary,
							"exclude_prefix_len cannot be configured when type = inet"+
								" in dhcp_attributes block in family block",
						)
					}
					if !config.Family.DhcpAttributes.PreferredLifetime.IsNull() &&
						!config.Family.DhcpAttributes.PreferredLifetime.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("dhcp_attributes").AtName("preferred_lifetime"),
							tfdiag.ConflictConfigErrSummary,
							"preferred_lifetime cannot be configured when type = inet"+
								" in dhcp_attributes block in family block",
						)
					}
					if !config.Family.DhcpAttributes.PreferredLifetimeInfinite.IsNull() &&
						!config.Family.DhcpAttributes.PreferredLifetimeInfinite.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("dhcp_attributes").AtName("preferred_lifetime_infinite"),
							tfdiag.ConflictConfigErrSummary,
							"preferred_lifetime_infinite cannot be configured when type = inet"+
								" in dhcp_attributes block in family block",
						)
					}
					if !config.Family.DhcpAttributes.SIPServerInet6Address.IsNull() &&
						!config.Family.DhcpAttributes.SIPServerInet6Address.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("dhcp_attributes").AtName("sip_server_inet6_address"),
							tfdiag.ConflictConfigErrSummary,
							"sip_server_inet6_address cannot be configured when type = inet"+
								" in dhcp_attributes block in family block",
						)
					}
					if !config.Family.DhcpAttributes.SIPServerInet6DomainName.IsNull() &&
						!config.Family.DhcpAttributes.SIPServerInet6DomainName.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("dhcp_attributes").AtName("sip_server_inet6_domain_name"),
							tfdiag.ConflictConfigErrSummary,
							"sip_server_inet6_domain_name cannot be configured when type = inet"+
								" in dhcp_attributes block in family block",
						)
					}
					if !config.Family.DhcpAttributes.ValidLifetime.IsNull() &&
						!config.Family.DhcpAttributes.ValidLifetime.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("dhcp_attributes").AtName("valid_lifetime"),
							tfdiag.ConflictConfigErrSummary,
							"valid_lifetime cannot be configured when type = inet"+
								" in dhcp_attributes block in family block",
						)
					}
					if !config.Family.DhcpAttributes.ValidLifetimeInfinite.IsNull() &&
						!config.Family.DhcpAttributes.ValidLifetimeInfinite.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("dhcp_attributes").AtName("valid_lifetime_infinite"),
							tfdiag.ConflictConfigErrSummary,
							"valid_lifetime_infinite cannot be configured when type = inet"+
								" in dhcp_attributes block in family block",
						)
					}
				}
				if !config.Family.Inet6Range.IsNull() &&
					!config.Family.Inet6Range.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("inet6_range").AtName("*"),
						tfdiag.ConflictConfigErrSummary,
						"inet6_range cannot be configured when type = inet"+
							" in family block",
					)
				}
			case junos.Inet6W:
				if !config.Family.Network.IsNull() &&
					!config.Family.Network.IsUnknown() {
					ipV6Validator := tfvalidator.StringCIDR().IPv6Only()
					attrReq := validator.StringRequest{
						Path:        path.Root("family").AtName("network"),
						ConfigValue: config.Family.Network,
					}
					attrResp := validator.StringResponse{}
					ipV6Validator.ValidateString(ctx, attrReq, &attrResp)
					resp.Diagnostics.Append(attrResp.Diagnostics.Errors()...)
				}
				if !config.Family.ExcludedAddress.IsNull() &&
					!config.Family.ExcludedAddress.IsUnknown() {
					ipV6Validator := tfvalidator.StringIPAddress().IPv6Only()
					var excludedAddress []types.String
					asDiags := config.Family.ExcludedAddress.ElementsAs(ctx, &excludedAddress, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}
					for _, v := range excludedAddress {
						attrReq := validator.StringRequest{
							Path:        path.Root("family").AtName("excluded_address"),
							ConfigValue: v,
						}
						attrResp := validator.StringResponse{}
						ipV6Validator.ValidateString(ctx, attrReq, &attrResp)
						resp.Diagnostics.Append(attrResp.Diagnostics.Errors()...)
					}
				}
				if !config.Family.XauthAttributesPrimaryDNS.IsNull() &&
					!config.Family.XauthAttributesPrimaryDNS.IsUnknown() {
					ipV6Validator := tfvalidator.StringIPAddress().IPv6Only()
					attrReq := validator.StringRequest{
						Path:        path.Root("family").AtName("xauth_attributes_primary_dns"),
						ConfigValue: config.Family.XauthAttributesPrimaryDNS,
					}
					attrResp := validator.StringResponse{}
					ipV6Validator.ValidateString(ctx, attrReq, &attrResp)
					resp.Diagnostics.Append(attrResp.Diagnostics.Errors()...)
				}
				if !config.Family.XauthAttributesPrimaryWins.IsNull() &&
					!config.Family.XauthAttributesPrimaryWins.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("xauth_attributes_primary_wins"),
						tfdiag.ConflictConfigErrSummary,
						"xauth_attributes_primary_wins cannot be configured when type = inet6"+
							" in family block",
					)
				}
				if !config.Family.XauthAttributesSecondaryDNS.IsNull() &&
					!config.Family.XauthAttributesSecondaryDNS.IsUnknown() {
					ipV6Validator := tfvalidator.StringIPAddress().IPv6Only()
					attrReq := validator.StringRequest{
						Path:        path.Root("family").AtName("xauth_attributes_secondary_dns"),
						ConfigValue: config.Family.XauthAttributesSecondaryDNS,
					}
					attrResp := validator.StringResponse{}
					ipV6Validator.ValidateString(ctx, attrReq, &attrResp)
					resp.Diagnostics.Append(attrResp.Diagnostics.Errors()...)
				}
				if !config.Family.XauthAttributesSecondaryWins.IsNull() &&
					!config.Family.XauthAttributesSecondaryWins.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("xauth_attributes_secondary_wins"),
						tfdiag.ConflictConfigErrSummary,
						"xauth_attributes_secondary_wins cannot be configured when type = inet6"+
							" in family block",
					)
				}
				if !config.Family.InetRange.IsNull() &&
					!config.Family.InetRange.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("inet_range").AtName("*"),
						tfdiag.ConflictConfigErrSummary,
						"inet_range cannot be configured when type = inet6"+
							" in family block",
					)
				}
			}
		}
		if config.Family.DhcpAttributes != nil {
			if config.Family.DhcpAttributes.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family").AtName("dhcp_attributes").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"dhcp_attributes block is empty"+
						" in family block",
				)
			}
			if !config.Family.DhcpAttributes.MaximumLeaseTime.IsNull() &&
				!config.Family.DhcpAttributes.MaximumLeaseTime.IsUnknown() {
				if !config.Family.DhcpAttributes.MaximumLeaseTimeInfinite.IsNull() &&
					!config.Family.DhcpAttributes.MaximumLeaseTimeInfinite.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("maximum_lease_time"),
						tfdiag.ConflictConfigErrSummary,
						"maximum_lease_time and maximum_lease_time_infinite cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
				if !config.Family.DhcpAttributes.PreferredLifetime.IsNull() &&
					!config.Family.DhcpAttributes.PreferredLifetime.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("maximum_lease_time"),
						tfdiag.ConflictConfigErrSummary,
						"maximum_lease_time and preferred_lifetime cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
				if !config.Family.DhcpAttributes.PreferredLifetimeInfinite.IsNull() &&
					!config.Family.DhcpAttributes.PreferredLifetimeInfinite.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("maximum_lease_time"),
						tfdiag.ConflictConfigErrSummary,
						"maximum_lease_time and preferred_lifetime_infinite cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
				if !config.Family.DhcpAttributes.ValidLifetime.IsNull() &&
					!config.Family.DhcpAttributes.ValidLifetime.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("maximum_lease_time"),
						tfdiag.ConflictConfigErrSummary,
						"maximum_lease_time and valid_lifetime cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
				if !config.Family.DhcpAttributes.ValidLifetimeInfinite.IsNull() &&
					!config.Family.DhcpAttributes.ValidLifetimeInfinite.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("maximum_lease_time"),
						tfdiag.ConflictConfigErrSummary,
						"maximum_lease_time and valid_lifetime_infinite cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
			}
			if !config.Family.DhcpAttributes.MaximumLeaseTimeInfinite.IsNull() &&
				!config.Family.DhcpAttributes.MaximumLeaseTimeInfinite.IsUnknown() {
				if !config.Family.DhcpAttributes.PreferredLifetime.IsNull() &&
					!config.Family.DhcpAttributes.PreferredLifetime.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("maximum_lease_time_infinite"),
						tfdiag.ConflictConfigErrSummary,
						"maximum_lease_time_infinite and preferred_lifetime cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
				if !config.Family.DhcpAttributes.PreferredLifetimeInfinite.IsNull() &&
					!config.Family.DhcpAttributes.PreferredLifetimeInfinite.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("maximum_lease_time_infinite"),
						tfdiag.ConflictConfigErrSummary,
						"maximum_lease_time_infinite and preferred_lifetime_infinite cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
				if !config.Family.DhcpAttributes.ValidLifetime.IsNull() &&
					!config.Family.DhcpAttributes.ValidLifetime.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("maximum_lease_time_infinite"),
						tfdiag.ConflictConfigErrSummary,
						"maximum_lease_time_infinite and valid_lifetime cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
				if !config.Family.DhcpAttributes.ValidLifetimeInfinite.IsNull() &&
					!config.Family.DhcpAttributes.ValidLifetimeInfinite.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("maximum_lease_time_infinite"),
						tfdiag.ConflictConfigErrSummary,
						"maximum_lease_time_infinite and valid_lifetime_infinite cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
			}
			if !config.Family.DhcpAttributes.PreferredLifetime.IsNull() &&
				!config.Family.DhcpAttributes.PreferredLifetime.IsUnknown() &&
				!config.Family.DhcpAttributes.PreferredLifetimeInfinite.IsNull() &&
				!config.Family.DhcpAttributes.PreferredLifetimeInfinite.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family").AtName("dhcp_attributes").AtName("preferred_lifetime"),
					tfdiag.ConflictConfigErrSummary,
					"preferred_lifetime and preferred_lifetime_infinite cannot be configured together"+
						" in dhcp_attributes block in family block",
				)
			}
			if !config.Family.DhcpAttributes.ValidLifetime.IsNull() &&
				!config.Family.DhcpAttributes.ValidLifetime.IsUnknown() &&
				!config.Family.DhcpAttributes.ValidLifetimeInfinite.IsNull() &&
				!config.Family.DhcpAttributes.ValidLifetimeInfinite.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family").AtName("dhcp_attributes").AtName("valid_lifetime"),
					tfdiag.ConflictConfigErrSummary,
					"valid_lifetime and valid_lifetime_infinite cannot be configured together"+
						" in dhcp_attributes block in family block",
				)
			}
			if !config.Family.DhcpAttributes.OptionMatch82CircuitID.IsNull() &&
				!config.Family.DhcpAttributes.OptionMatch82CircuitID.IsUnknown() {
				var configOptionMatch82CircuitID []accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesBlockOptionMatch82
				asDiags := config.Family.DhcpAttributes.OptionMatch82CircuitID.ElementsAs(ctx, &configOptionMatch82CircuitID, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				optionMatch82CircuitIDValue := make(map[string]struct{})
				for i, block := range configOptionMatch82CircuitID {
					if block.Value.IsUnknown() {
						continue
					}
					value := block.Value.ValueString()
					if _, ok := optionMatch82CircuitIDValue[value]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("dhcp_attributes").
								AtName("option_match_82_circuit_id").AtListIndex(i).AtName("value"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple option_match_82_circuit_id blocks with the same value %q"+
								" in dhcp_attributes block in family block", value),
						)
					}
					optionMatch82CircuitIDValue[value] = struct{}{}
				}
			}
			if !config.Family.DhcpAttributes.OptionMatch82RemoteID.IsNull() &&
				!config.Family.DhcpAttributes.OptionMatch82RemoteID.IsUnknown() {
				var configOptionMatch82RemoteID []accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesBlockOptionMatch82
				asDiags := config.Family.DhcpAttributes.OptionMatch82RemoteID.ElementsAs(ctx, &configOptionMatch82RemoteID, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				optionMatch82RemoteIDValue := make(map[string]struct{})
				for i, block := range configOptionMatch82RemoteID {
					if block.Value.IsUnknown() {
						continue
					}
					value := block.Value.ValueString()
					if _, ok := optionMatch82RemoteIDValue[value]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("dhcp_attributes").
								AtName("option_match_82_remote_id").AtListIndex(i).AtName("value"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple option_match_82_remote_id blocks with the same value %q"+
								" in dhcp_attributes block in family block", value),
						)
					}
					optionMatch82RemoteIDValue[value] = struct{}{}
				}
			}
			if !config.Family.DhcpAttributes.T1Percentage.IsNull() &&
				!config.Family.DhcpAttributes.T1Percentage.IsUnknown() {
				if !config.Family.DhcpAttributes.T1RenewalTime.IsNull() &&
					!config.Family.DhcpAttributes.T1RenewalTime.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("t1_percentage"),
						tfdiag.ConflictConfigErrSummary,
						"t1_percentage and t1_renewal_time cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
				if !config.Family.DhcpAttributes.T2RebindingTime.IsNull() &&
					!config.Family.DhcpAttributes.T2RebindingTime.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("t1_percentage"),
						tfdiag.ConflictConfigErrSummary,
						"t1_percentage and t2_rebinding_time cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
			}
			if !config.Family.DhcpAttributes.T1RenewalTime.IsNull() &&
				!config.Family.DhcpAttributes.T1RenewalTime.IsUnknown() {
				if !config.Family.DhcpAttributes.T2Percentage.IsNull() &&
					!config.Family.DhcpAttributes.T2Percentage.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("t1_renewal_time"),
						tfdiag.ConflictConfigErrSummary,
						"t1_renewal_time and t2_percentage cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
			}
			if !config.Family.DhcpAttributes.T2Percentage.IsNull() &&
				!config.Family.DhcpAttributes.T2Percentage.IsUnknown() {
				if !config.Family.DhcpAttributes.T2RebindingTime.IsNull() &&
					!config.Family.DhcpAttributes.T2RebindingTime.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("dhcp_attributes").AtName("t2_percentage"),
						tfdiag.ConflictConfigErrSummary,
						"t2_percentage and t2_rebinding_time cannot be configured together"+
							" in dhcp_attributes block in family block",
					)
				}
			}
		}
		if !config.Family.ExcludedRange.IsNull() &&
			!config.Family.ExcludedRange.IsUnknown() {
			var configExcludedRange []accessAddressAssignmentPoolBlockFamilyBlockExcludedRange
			asDiags := config.Family.ExcludedRange.ElementsAs(ctx, &configExcludedRange, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			excludedRangeName := make(map[string]struct{})
			for i, block := range configExcludedRange {
				if !block.Name.IsUnknown() {
					name := block.Name.ValueString()
					if _, ok := excludedRangeName[name]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("excluded_range").AtListIndex(i).AtName("name"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple excluded_range blocks with the same name %q"+
								" in family block", name),
						)
					}
					excludedRangeName[name] = struct{}{}
				}
				if !config.Family.Type.IsNull() &&
					!config.Family.Type.IsUnknown() {
					ipValidator := tfvalidator.StringIPAddress().IPv4Only()
					if config.Family.Type.ValueString() == junos.Inet6W {
						ipValidator = tfvalidator.StringIPAddress().IPv6Only()
					}

					attrReq := validator.StringRequest{
						Path:        path.Root("family").AtName("excluded_range").AtListIndex(i).AtName("low"),
						ConfigValue: block.Low,
					}
					attrResp := validator.StringResponse{}
					ipValidator.ValidateString(ctx, attrReq, &attrResp)
					resp.Diagnostics.Append(attrResp.Diagnostics.Errors()...)

					attrReq2 := validator.StringRequest{
						Path:        path.Root("family").AtName("excluded_range").AtListIndex(i).AtName("high"),
						ConfigValue: block.High,
					}
					attrResp2 := validator.StringResponse{}
					ipValidator.ValidateString(ctx, attrReq2, &attrResp2)
					resp.Diagnostics.Append(attrResp2.Diagnostics.Errors()...)
				}
			}
		}
		if !config.Family.Host.IsNull() &&
			!config.Family.Host.IsUnknown() {
			var configHost []accessAddressAssignmentPoolBlockFamilyBlockHost
			asDiags := config.Family.Host.ElementsAs(ctx, &configHost, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			hostName := make(map[string]struct{})
			for i, block := range configHost {
				if !block.Name.IsUnknown() {
					name := block.Name.ValueString()
					if _, ok := hostName[name]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("host").AtListIndex(i).AtName("name"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple host blocks with the same name %q"+
								" in family block", name),
						)
					}
					hostName[name] = struct{}{}
				}
				if !config.Family.Type.IsNull() &&
					!config.Family.Type.IsUnknown() {
					valid := tfvalidator.StringIPAddress().IPv4Only()
					if config.Family.Type.ValueString() == junos.Inet6W {
						valid = tfvalidator.StringIPAddress().IPv6Only()
					}

					attrReq := validator.StringRequest{
						Path:        path.Root("family").AtName("host").AtListIndex(i).AtName("ip_address"),
						ConfigValue: block.IPAddress,
					}
					attrResp := validator.StringResponse{}
					valid.ValidateString(ctx, attrReq, &attrResp)
					resp.Diagnostics.Append(attrResp.Diagnostics.Errors()...)
				}
				if block.HardwareAddress.IsNull() &&
					block.UserName.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("host").AtListIndex(i).AtName("name"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("hardware_address or user_name must be specified"+
							" in host block %q in family block", block.Name.ValueString()),
					)
				}
				if !block.HardwareAddress.IsNull() &&
					!block.HardwareAddress.IsUnknown() &&
					!block.UserName.IsNull() &&
					!block.UserName.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("host").AtListIndex(i).AtName("hardware_address"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("hardware_address and user_name cannot be configured together"+
							" in host block %q in family block", block.Name.ValueString()),
					)
				}
			}
		}
		if !config.Family.InetRange.IsNull() &&
			!config.Family.InetRange.IsUnknown() {
			var configInetRange []accessAddressAssignmentPoolBlockFamilyBlockInetRange
			asDiags := config.Family.InetRange.ElementsAs(ctx, &configInetRange, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			inetRangeName := make(map[string]struct{})
			for i, block := range configInetRange {
				if block.Name.IsUnknown() {
					continue
				}
				name := block.Name.ValueString()
				if _, ok := inetRangeName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("inet_range").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple inet_range blocks with the same name %q"+
							" in family block", name),
					)
				}
				inetRangeName[name] = struct{}{}
			}
		}
		if !config.Family.Inet6Range.IsNull() &&
			!config.Family.Inet6Range.IsUnknown() {
			var configInet6Range []accessAddressAssignmentPoolBlockFamilyBlockInet6Range
			asDiags := config.Family.Inet6Range.ElementsAs(ctx, &configInet6Range, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			inet6RangeName := make(map[string]struct{})
			for i, block := range configInet6Range {
				if !block.Name.IsUnknown() {
					name := block.Name.ValueString()
					if _, ok := inet6RangeName[name]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("name"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple inet6_range blocks with the same name %q"+
								" in family block", name),
						)
					}
					inet6RangeName[name] = struct{}{}
				}
				if !block.Low.IsNull() &&
					!block.Low.IsUnknown() {
					if !block.PrefixLength.IsNull() &&
						!block.PrefixLength.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("prefix_length"),
							tfdiag.ConflictConfigErrSummary,
							fmt.Sprintf("prefix_length and low cannot be configured together"+
								" in inet6_range block %q in family block", block.Name.ValueString()),
						)
					}
					if block.High.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("low"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("high must be specified with low"+
								" in inet6_range block %q in family block", block.Name.ValueString()),
						)
					}
				}
				if !block.High.IsNull() &&
					!block.High.IsUnknown() {
					if !block.PrefixLength.IsNull() &&
						!block.PrefixLength.IsUnknown() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("prefix_length"),
							tfdiag.ConflictConfigErrSummary,
							fmt.Sprintf("prefix_length and high cannot be configured together"+
								" in inet6_range block %q in family block", block.Name.ValueString()),
						)
					}
					if block.Low.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("high"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("low must be specified with high"+
								" in inet6_range block %q in family block", block.Name.ValueString()),
						)
					}
				}
				if block.Low.IsNull() &&
					block.High.IsNull() &&
					block.PrefixLength.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("name"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("prefix_length or combination of low and high must be specified"+
							" in inet6_range block %q in family block", block.Name.ValueString()),
					)
				}
			}
		}
	}
}

func (rsc *accessAddressAssignmentPool) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan accessAddressAssignmentPoolData
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
			poolExists, err := checkAccessAddressAssignmentPoolExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if poolExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsInRoutingInstanceMessage(rsc, plan.Name, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsMessage(rsc, plan.Name),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			poolExists, err := checkAccessAddressAssignmentPoolExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !poolExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsInRoutingInstanceAfterCommitMessage(rsc, plan.Name, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Name),
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

func (rsc *accessAddressAssignmentPool) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data accessAddressAssignmentPoolData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *accessAddressAssignmentPool) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state accessAddressAssignmentPoolData
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

func (rsc *accessAddressAssignmentPool) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state accessAddressAssignmentPoolData
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

func (rsc *accessAddressAssignmentPool) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data accessAddressAssignmentPoolData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <name>"+junos.IDSeparator+"<routing_instance>)",
	)
}

func checkAccessAddressAssignmentPoolExists(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"access address-assignment pool " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *accessAddressAssignmentPoolData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *accessAddressAssignmentPoolData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *accessAddressAssignmentPoolData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "access address-assignment pool " + rscData.Name.ValueString() + " "

	if rscData.ActiveDrain.ValueBool() {
		configSet = append(configSet, setPrefix+"active-drain")
	}
	if rscData.HoldDown.ValueBool() {
		configSet = append(configSet, setPrefix+"hold-down")
	}
	if v := rscData.Link.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"link "+v)
	}
	if rscData.Family != nil {
		blockSet, pathErr, err := rscData.Family.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *accessAddressAssignmentPoolBlockFamily) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, 1)
	setPrefix += "family "

	familyType := block.Type.ValueString()
	if familyType == junos.Inet6W {
		setPrefix += "inet6 "
		configSet = append(configSet, setPrefix+"prefix "+block.Network.ValueString())
	} else {
		familyType = junos.InetW
		setPrefix += "inet "
		configSet = append(configSet, setPrefix+"network "+block.Network.ValueString())
	}

	for _, v := range block.ExcludedAddress {
		configSet = append(configSet, setPrefix+"excluded-address "+v.ValueString())
	}
	if v := block.XauthAttributesPrimaryDNS.ValueString(); v != "" {
		if familyType == junos.Inet6W {
			configSet = append(configSet, setPrefix+"xauth-attributes primary-dns-ipv6 "+v)
		} else {
			configSet = append(configSet, setPrefix+"xauth-attributes primary-dns "+v)
		}
	}
	if v := block.XauthAttributesPrimaryWins.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"xauth-attributes primary-wins "+v)
	}
	if v := block.XauthAttributesSecondaryDNS.ValueString(); v != "" {
		if familyType == junos.Inet6W {
			configSet = append(configSet, setPrefix+"xauth-attributes secondary-dns-ipv6 "+v)
		} else {
			configSet = append(configSet, setPrefix+"xauth-attributes secondary-dns "+v)
		}
	}
	if v := block.XauthAttributesSecondaryWins.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"xauth-attributes secondary-wins "+v)
	}

	if block.DhcpAttributes != nil {
		if block.DhcpAttributes.isEmpty() {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("*"),
				errors.New("dhcp_attributes block is empty" +
					" in family block")
		}

		blockSet, pathErr, err := block.DhcpAttributes.configSet(setPrefix, familyType)
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	excludedRangeName := make(map[string]struct{})
	for i, subBlock := range block.ExcludedRange {
		name := subBlock.Name.ValueString()
		if _, ok := excludedRangeName[name]; ok {
			return configSet,
				path.Root("family").AtName("excluded_range").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple excluded_range blocks with the same name %q"+
					" in family block", name)
		}
		excludedRangeName[name] = struct{}{}

		configSet = append(configSet, setPrefix+"excluded-range "+name+" low "+subBlock.Low.ValueString())
		configSet = append(configSet, setPrefix+"excluded-range "+name+" high "+subBlock.High.ValueString())
	}
	hostName := make(map[string]struct{})
	for i, subBlock := range block.Host {
		name := subBlock.Name.ValueString()
		if _, ok := hostName[name]; ok {
			return configSet,
				path.Root("family").AtName("host").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple host blocks with the same name %q"+
					" in family block", name)
		}
		hostName[name] = struct{}{}

		configSet = append(configSet, setPrefix+"host \""+name+"\" ip-address "+subBlock.IPAddress.ValueString())
		if v := subBlock.HardwareAddress.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"host \""+name+"\" hardware-address "+v)
		}
		if subBlock.UserName.ValueBool() {
			configSet = append(configSet, setPrefix+"host \""+name+"\" user-name")
		}
	}
	inetRangeName := make(map[string]struct{})
	for i, subBlock := range block.InetRange {
		if familyType == junos.Inet6W {
			return configSet,
				path.Root("family").AtName("inet_range").AtListIndex(i).AtName("name"),
				errors.New("inet_range cannot be configured when type = inet6" +
					" in family block")
		}
		name := subBlock.Name.ValueString()
		if _, ok := inetRangeName[name]; ok {
			return configSet,
				path.Root("family").AtName("inet_range").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple inet_range blocks with the same name %q"+
					" in family block", name)
		}
		inetRangeName[name] = struct{}{}

		configSet = append(configSet, setPrefix+"range "+name+" low "+subBlock.Low.ValueString())
		configSet = append(configSet, setPrefix+"range "+name+" high "+subBlock.High.ValueString())
	}
	for i, subBlock := range block.Inet6Range {
		if familyType == junos.InetW {
			return configSet,
				path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("name"),
				errors.New("inet6_range cannot be configured when type = inet" +
					" in family block")
		}
		name := subBlock.Name.ValueString()
		if _, ok := inetRangeName[name]; ok {
			return configSet,
				path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple inet6_range blocks with the same name %q"+
					" in family block", name)
		}
		inetRangeName[name] = struct{}{}
		switch {
		case !subBlock.PrefixLength.IsNull():
			if !subBlock.Low.IsNull() {
				return configSet,
					path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("prefix_length"),
					fmt.Errorf("prefix_length and low cannot be configured together"+
						" in inet6_range block %q in family block", name)
			}
			if !subBlock.High.IsNull() {
				return configSet,
					path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("prefix_length"),
					fmt.Errorf("prefix_length and high cannot be configured together"+
						" in inet6_range block %q in family block", name)
			}

			configSet = append(configSet, setPrefix+"range "+name+" prefix-length "+
				utils.ConvI64toa(subBlock.PrefixLength.ValueInt64()))
		case !subBlock.Low.IsNull() && !subBlock.High.IsNull():
			configSet = append(configSet, setPrefix+"range "+name+" low "+subBlock.Low.ValueString())
			configSet = append(configSet, setPrefix+"range "+name+" high "+subBlock.High.ValueString())
		default:
			return configSet,
				path.Root("family").AtName("inet6_range").AtListIndex(i).AtName("*"),
				fmt.Errorf("prefix_length or combination of low and high must be specified"+
					" in inet6_range block %q in family block", name)
		}
	}

	return configSet, path.Empty(), nil
}

func (block *accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributes) configSet(
	setPrefix,
	familyType string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "dhcp-attributes "

	if v := block.BootFile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"boot-file \""+v+"\"")
	}
	if v := block.BootServer.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"boot-server "+v)
	}
	for _, v := range block.DNSServer {
		if familyType == junos.InetW {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("dns_server"),
				errors.New("dns_server cannot be configured when type = inet" +
					" in dhcp_attributes block in family block")
		}

		configSet = append(configSet, setPrefix+"dns-server "+v.ValueString())
	}
	if v := block.DomainName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"domain-name "+v)
	}
	if !block.ExcludePrefixLen.IsNull() {
		if familyType == junos.InetW {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("exclude_prefix_len"),
				errors.New("exclude_prefix_len cannot be configured when type = inet" +
					" in dhcp_attributes block in family block")
		}

		configSet = append(configSet, setPrefix+"exclude-prefix-len "+
			utils.ConvI64toa(block.ExcludePrefixLen.ValueInt64()))
	}
	if !block.GracePeriod.IsNull() {
		configSet = append(configSet, setPrefix+"grace-period "+
			utils.ConvI64toa(block.GracePeriod.ValueInt64()))
	}
	if !block.MaximumLeaseTime.IsNull() {
		configSet = append(configSet, setPrefix+"maximum-lease-time "+
			utils.ConvI64toa(block.MaximumLeaseTime.ValueInt64()))
	}
	if block.MaximumLeaseTimeInfinite.ValueBool() {
		configSet = append(configSet, setPrefix+"maximum-lease-time infinite")
	}
	for _, v := range block.NameServer {
		configSet = append(configSet, setPrefix+"name-server "+v.ValueString())
	}
	if v := block.NetbiosNodeType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"netbios-node-type "+v)
	}
	if v := block.NextServer.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"next-server "+v)
	}
	for _, v := range block.Option {
		configSet = append(configSet, setPrefix+"option "+v.ValueString())
	}
	if !block.PreferredLifetime.IsNull() {
		if familyType == junos.InetW {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("preferred_lifetime"),
				errors.New("preferred_lifetime cannot be configured when type = inet" +
					" in dhcp_attributes block in family block")
		}

		configSet = append(configSet, setPrefix+"preferred-lifetime "+
			utils.ConvI64toa(block.PreferredLifetime.ValueInt64()))
	}
	if block.PreferredLifetimeInfinite.ValueBool() {
		if familyType == junos.InetW {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("preferred_lifetime_infinite"),
				errors.New("preferred_lifetime_infinite cannot be configured when type = inet" +
					" in dhcp_attributes block in family block")
		}

		configSet = append(configSet, setPrefix+"preferred-lifetime infinite")
	}
	for _, v := range block.PropagatePppSettings {
		configSet = append(configSet, setPrefix+"propagate-ppp-settings "+v.ValueString())
	}
	if v := block.PropagateSettings.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"propagate-settings \""+v+"\"")
	}
	for _, v := range block.Router {
		configSet = append(configSet, setPrefix+"router "+v.ValueString())
	}
	if v := block.ServerIdentifier.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"server-identifier "+v)
	}
	for _, v := range block.SIPServerInetAddress {
		configSet = append(configSet, setPrefix+"sip-server ip-address "+v.ValueString())
	}
	for _, v := range block.SIPServerInetDomainName {
		configSet = append(configSet, setPrefix+"sip-server name \""+v.ValueString()+"\"")
	}
	for _, v := range block.SIPServerInet6Address {
		if familyType == junos.InetW {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("sip_server_inet6_address"),
				errors.New("sip_server_inet6_address cannot be configured when type = inet" +
					" in dhcp_attributes block in family block")
		}

		configSet = append(configSet, setPrefix+"sip-server-address "+v.ValueString())
	}

	if v := block.SIPServerInet6DomainName.ValueString(); v != "" {
		if familyType == junos.InetW {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("sip_server_inet6_domain_name"),
				errors.New("sip_server_inet6_domain_name cannot be configured when type = inet" +
					" in dhcp_attributes block in family block")
		}

		configSet = append(configSet, setPrefix+"sip-server-domain-name \""+v+"\"")
	}
	if !block.T1Percentage.IsNull() {
		configSet = append(configSet, setPrefix+"t1-percentage "+
			utils.ConvI64toa(block.T1Percentage.ValueInt64()))
	}
	if !block.T1RenewalTime.IsNull() {
		configSet = append(configSet, setPrefix+"t1-renewal-time "+
			utils.ConvI64toa(block.T1RenewalTime.ValueInt64()))
	}
	if !block.T2Percentage.IsNull() {
		configSet = append(configSet, setPrefix+"t2-percentage "+
			utils.ConvI64toa(block.T2Percentage.ValueInt64()))
	}
	if !block.T2RebindingTime.IsNull() {
		configSet = append(configSet, setPrefix+"t2-rebinding-time "+
			utils.ConvI64toa(block.T2RebindingTime.ValueInt64()))
	}
	if v := block.TftpServer.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"tftp-server "+v)
	}
	if !block.ValidLifetime.IsNull() {
		if familyType == junos.InetW {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("valid_lifetime"),
				errors.New("valid_lifetime cannot be configured when type = inet" +
					" in dhcp_attributes block in family block")
		}

		configSet = append(configSet, setPrefix+"valid-lifetime "+
			utils.ConvI64toa(block.ValidLifetime.ValueInt64()))
	}
	if block.ValidLifetimeInfinite.ValueBool() {
		if familyType == junos.InetW {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("valid_lifetime_infinite"),
				errors.New("valid_lifetime_infinite cannot be configured when type = inet" +
					" in dhcp_attributes block in family block")
		}

		configSet = append(configSet, setPrefix+"valid-lifetime infinite")
	}
	for _, v := range block.WinsServer {
		configSet = append(configSet, setPrefix+"wins-server "+v.ValueString())
	}

	optionMatch82CircuitIDValue := make(map[string]struct{})
	for i, subBlock := range block.OptionMatch82CircuitID {
		value := subBlock.Value.ValueString()
		if _, ok := optionMatch82CircuitIDValue[value]; ok {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("option_match_82_circuit_id").AtListIndex(i).AtName("value"),
				fmt.Errorf("multiple option_match_82_circuit_id blocks with the same value %q"+
					" in dhcp_attributes block in family block", value)
		}
		optionMatch82CircuitIDValue[value] = struct{}{}

		configSet = append(configSet, setPrefix+"option-match option-82"+
			" circuit-id \""+value+"\""+
			" range \""+subBlock.Range.ValueString()+"\"")
	}
	optionMatch82RemoteIDValue := make(map[string]struct{})
	for i, subBlock := range block.OptionMatch82RemoteID {
		value := subBlock.Value.ValueString()
		if _, ok := optionMatch82RemoteIDValue[value]; ok {
			return configSet,
				path.Root("family").AtName("dhcp_attributes").AtName("option_match_82_remote_id").AtListIndex(i).AtName("value"),
				fmt.Errorf("multiple option_match_82_remote_id blocks with the same value %q"+
					" in dhcp_attributes block in family block", value)
		}
		optionMatch82RemoteIDValue[value] = struct{}{}

		configSet = append(configSet, setPrefix+"option-match option-82"+
			" remote-id \""+value+"\""+
			" range \""+subBlock.Range.ValueString()+"\"")
	}

	return configSet, path.Empty(), nil
}

func (rscData *accessAddressAssignmentPoolData) read(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"access address-assignment pool " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
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
			case itemTrim == "active-drain":
				rscData.ActiveDrain = types.BoolValue(true)
			case itemTrim == "hold-down":
				rscData.HoldDown = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "link "):
				rscData.Link = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "family "):
				familyType := tfdata.FirstElementOfJunosLine(itemTrim)
				if rscData.Family == nil {
					rscData.Family = &accessAddressAssignmentPoolBlockFamily{
						Type: types.StringValue(familyType),
					}
				}

				if err := rscData.Family.read(strings.TrimPrefix(itemTrim, familyType+" ")); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *accessAddressAssignmentPoolBlockFamily) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "network "):
		block.Network = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "prefix "):
		block.Network = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "excluded-address "):
		block.ExcludedAddress = append(block.ExcludedAddress, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "xauth-attributes primary-dns "):
		block.XauthAttributesPrimaryDNS = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "xauth-attributes primary-dns-ipv6 "):
		block.XauthAttributesPrimaryDNS = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "xauth-attributes primary-wins "):
		block.XauthAttributesPrimaryWins = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "xauth-attributes secondary-dns "):
		block.XauthAttributesSecondaryDNS = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "xauth-attributes secondary-dns-ipv6 "):
		block.XauthAttributesSecondaryDNS = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "xauth-attributes secondary-wins "):
		block.XauthAttributesSecondaryWins = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "dhcp-attributes "):
		if block.DhcpAttributes == nil {
			block.DhcpAttributes = &accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributes{}
		}

		if err := block.DhcpAttributes.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "excluded-range "):
		name := tfdata.FirstElementOfJunosLine(itemTrim)
		block.ExcludedRange = tfdata.AppendPotentialNewBlock(block.ExcludedRange, types.StringValue(name))
		excludedRange := &block.ExcludedRange[len(block.ExcludedRange)-1]
		balt.CutPrefixInString(&itemTrim, name+" ")

		switch {
		case balt.CutPrefixInString(&itemTrim, "low "):
			excludedRange.Low = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "high "):
			excludedRange.High = types.StringValue(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "host "):
		name := tfdata.FirstElementOfJunosLine(itemTrim)
		block.Host = tfdata.AppendPotentialNewBlock(block.Host, types.StringValue(strings.Trim(name, "\"")))
		host := &block.Host[len(block.Host)-1]
		balt.CutPrefixInString(&itemTrim, name+" ")

		switch {
		case balt.CutPrefixInString(&itemTrim, "ip-address "):
			host.IPAddress = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "hardware-address "):
			host.HardwareAddress = types.StringValue(itemTrim)
		case itemTrim == "user-name":
			host.UserName = types.BoolValue(true)
		}
	case balt.CutPrefixInString(&itemTrim, "range "):
		if block.Type.ValueString() == junos.InetW {
			name := tfdata.FirstElementOfJunosLine(itemTrim)
			block.InetRange = tfdata.AppendPotentialNewBlock(block.InetRange, types.StringValue(name))
			inetRange := &block.InetRange[len(block.InetRange)-1]
			balt.CutPrefixInString(&itemTrim, name+" ")

			switch {
			case balt.CutPrefixInString(&itemTrim, "low "):
				inetRange.Low = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "high "):
				inetRange.High = types.StringValue(itemTrim)
			}
		} else if block.Type.ValueString() == junos.Inet6W {
			name := tfdata.FirstElementOfJunosLine(itemTrim)
			block.Inet6Range = tfdata.AppendPotentialNewBlock(block.Inet6Range, types.StringValue(name))
			inet6Range := &block.Inet6Range[len(block.Inet6Range)-1]
			balt.CutPrefixInString(&itemTrim, name+" ")

			switch {
			case balt.CutPrefixInString(&itemTrim, "low "):
				inet6Range.Low = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "high "):
				inet6Range.High = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "prefix-length "):
				inet6Range.PrefixLength, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributes) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "boot-file "):
		block.BootFile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "boot-server "):
		block.BootServer = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "dns-server "):
		block.DNSServer = append(block.DNSServer, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "domain-name "):
		block.DomainName = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "exclude-prefix-len "):
		block.ExcludePrefixLen, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "grace-period "):
		block.GracePeriod, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "maximum-lease-time infinite":
		block.MaximumLeaseTimeInfinite = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "maximum-lease-time "):
		block.MaximumLeaseTime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "name-server "):
		block.NameServer = append(block.NameServer, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "netbios-node-type "):
		block.NetbiosNodeType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "next-server "):
		block.NextServer = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "option "):
		block.Option = append(block.Option, types.StringValue(strings.Trim(itemTrim, "\"")))
	case itemTrim == "preferred-lifetime infinite":
		block.PreferredLifetimeInfinite = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "preferred-lifetime "):
		block.PreferredLifetime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "propagate-ppp-settings "):
		block.PropagatePppSettings = append(block.PropagatePppSettings, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "propagate-settings "):
		block.PropagateSettings = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "router "):
		block.Router = append(block.Router, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "server-identifier "):
		block.ServerIdentifier = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "sip-server ip-address "):
		block.SIPServerInetAddress = append(block.SIPServerInetAddress, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "sip-server name "):
		block.SIPServerInetDomainName = append(block.SIPServerInetDomainName, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "sip-server-address "):
		block.SIPServerInet6Address = append(block.SIPServerInet6Address, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "sip-server-domain-name "):
		block.SIPServerInet6DomainName = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "t1-percentage "):
		block.T1Percentage, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "t1-renewal-time "):
		block.T1RenewalTime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "t2-percentage "):
		block.T2Percentage, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "t2-rebinding-time "):
		block.T2RebindingTime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "tftp-server "):
		block.TftpServer = types.StringValue(itemTrim)
	case itemTrim == "valid-lifetime infinite":
		block.ValidLifetimeInfinite = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "valid-lifetime "):
		block.ValidLifetime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "wins-server "):
		block.WinsServer = append(block.WinsServer, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "option-match option-82 circuit-id "):
		value := tfdata.FirstElementOfJunosLine(itemTrim)
		subBlock := accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesBlockOptionMatch82{
			Value: types.StringValue(strings.Trim(value, "\"")),
		}
		balt.CutPrefixInString(&itemTrim, value+" ")

		if balt.CutPrefixInString(&itemTrim, "range ") {
			subBlock.Range = types.StringValue(strings.Trim(itemTrim, "\""))
		}
		block.OptionMatch82CircuitID = append(block.OptionMatch82CircuitID, subBlock)
	case balt.CutPrefixInString(&itemTrim, "option-match option-82 remote-id "):
		value := tfdata.FirstElementOfJunosLine(itemTrim)
		subBlock := accessAddressAssignmentPoolBlockFamilyBlockDhcpAttributesBlockOptionMatch82{
			Value: types.StringValue(strings.Trim(value, "\"")),
		}
		balt.CutPrefixInString(&itemTrim, value+" ")

		if balt.CutPrefixInString(&itemTrim, "range ") {
			subBlock.Range = types.StringValue(strings.Trim(itemTrim, "\""))
		}
		block.OptionMatch82RemoteID = append(block.OptionMatch82RemoteID, subBlock)
	}

	return nil
}

func (rscData *accessAddressAssignmentPoolData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}

	configSet := []string{
		delPrefix + "access address-assignment pool " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
