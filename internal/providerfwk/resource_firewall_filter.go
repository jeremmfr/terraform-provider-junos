package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

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
	_ resource.Resource                   = &firewallFilter{}
	_ resource.ResourceWithConfigure      = &firewallFilter{}
	_ resource.ResourceWithValidateConfig = &firewallFilter{}
	_ resource.ResourceWithImportState    = &firewallFilter{}
	_ resource.ResourceWithUpgradeState   = &firewallFilter{}
)

type firewallFilter struct {
	client *junos.Client
}

func newFirewallFilterResource() resource.Resource {
	return &firewallFilter{}
}

func (rsc *firewallFilter) typeName() string {
	return providerName + "_firewall_filter"
}

func (rsc *firewallFilter) junosName() string {
	return "firewall filter"
}

func (rsc *firewallFilter) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *firewallFilter) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *firewallFilter) Configure(
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

func (rsc *firewallFilter) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<name>" + junos.IDSeparator + "<family>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Filter name.",

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"family": schema.StringAttribute{
				Required:    true,
				Description: "Family where create this filter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(junos.InetW, junos.Inet6W, "any", "ccc", "mpls", "vpls", "ethernet-switching"),
				},
			},
			"interface_specific": schema.BoolAttribute{
				Optional:    true,
				Description: "Defined counters are interface specific.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"term": schema.ListNestedBlock{
				Description: "For each name of term.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Term name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 64),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"filter": schema.StringAttribute{
							Optional:    true,
							Description: "Filter to include.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 64),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"from": schema.SingleNestedBlock{
							Description: "Define match criteria.",
							Attributes: map[string]schema.Attribute{
								"address": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP source or destination address.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											tfvalidator.StringCIDRNetwork(),
										),
									},
								},
								"address_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP source or destination address not in this prefix.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											tfvalidator.StringCIDRNetwork(),
										),
									},
								},
								"destination_address": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP destination address.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											tfvalidator.StringCIDRNetwork(),
										),
									},
								},
								"destination_address_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP destination address not in this prefix.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											tfvalidator.StringCIDRNetwork(),
										),
									},
								},
								"destination_mac_address": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Destination MAC address.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.RegexMatches(regexp.MustCompile(
												`^[a-f0-9]{2}(:[a-f0-9]{2}){5}\/\d+$`),
												"must be an MAC address with mask",
											),
										),
									},
								},
								"destination_mac_address_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Destination MAC address not in this range.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.RegexMatches(regexp.MustCompile(
												`^[a-f0-9]{2}(:[a-f0-9]{2}){5}\/\d+$`),
												"must be an MAC address with mask",
											),
										),
									},
								},
								"destination_port": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match TCP/UDP destination port.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"destination_port_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match TCP/UDP destination port.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"destination_prefix_list": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP destination prefixes in named list.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 250),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
								"destination_prefix_list_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match addresses not in this prefix list.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 250),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
								"forwarding_class": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match forwarding class.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 64),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
								"forwarding_class_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match forwarding class.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 64),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
								"icmp_code": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match ICMP message code.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"icmp_code_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match ICMP message code.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"icmp_type": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match ICMP message type.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"icmp_type_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match ICMP message type.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"interface": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match interface name.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.InterfaceWithWildcardFormat),
										),
									},
								},
								"is_fragment": schema.BoolAttribute{
									Optional:    true,
									Description: "Match if packet is a fragment.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"loss_priority": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match Loss Priority.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.OneOf("high", "low", "medium-high", "medium-low"),
										),
									},
								},
								"loss_priority_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match Loss Priority.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.OneOf("high", "low", "medium-high", "medium-low"),
										),
									},
								},
								"next_header": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match next header protocol type.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"next_header_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match next header protocol type.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"packet_length": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match packet length.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.RegexMatches(regexp.MustCompile(
												`^\d+(-\d+)?$`),
												"must be an integer or a range of integers",
											),
										),
									},
								},
								"packet_length_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match packet length.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.RegexMatches(regexp.MustCompile(
												`^\d+(-\d+)?$`),
												"must be an integer or a range of integers",
											),
										),
									},
								},
								"policy_map": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match policy map.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 64),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
								"policy_map_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match policy map.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 64),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
								"port": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match TCP/UDP source or destination port.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"port_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match TCP/UDP source or destination port.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"prefix_list": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP source or destination prefixes in named list.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 250),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
								"prefix_list_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match addresses not in this prefix list.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 250),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
								"protocol": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP protocol type.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"protocol_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match IP protocol type.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"source_address": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP source address.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											tfvalidator.StringCIDRNetwork(),
										),
									},
								},
								"source_address_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP source address not in this prefix.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											tfvalidator.StringCIDRNetwork(),
										),
									},
								},
								"source_mac_address": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Source MAC address.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.RegexMatches(regexp.MustCompile(
												`^[a-f0-9]{2}(:[a-f0-9]{2}){5}\/\d+$`),
												"must be an MAC address with mask",
											),
										),
									},
								},
								"source_mac_address_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Source MAC address not in this range.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.RegexMatches(regexp.MustCompile(
												`^[a-f0-9]{2}(:[a-f0-9]{2}){5}\/\d+$`),
												"must be an MAC address with mask",
											),
										),
									},
								},
								"source_port": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match TCP/UDP source port.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"source_port_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Do not match TCP/UDP source port.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										),
									},
								},
								"source_prefix_list": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP source prefixes in named list.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 250),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
								"source_prefix_list_except": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Match IP source prefixes not in this prefix list.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 250),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
								"tcp_established": schema.BoolAttribute{
									Optional:    true,
									Description: "Match packet of an established TCP connection.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"tcp_flags": schema.StringAttribute{
									Optional:    true,
									Description: "Match TCP flags (in symbolic or hex formats).",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"tcp_initial": schema.BoolAttribute{
									Optional:    true,
									Description: "Match initial packet of a TCP connection.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
							},
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
						"then": schema.SingleNestedBlock{
							Description: "Define action to take if the `from` condition is matched.",
							Attributes: map[string]schema.Attribute{
								"action": schema.StringAttribute{
									Optional:    true,
									Description: "Action for term if needed.",
									Validators: []validator.String{
										stringvalidator.OneOf("accept", "reject", "discard", "next term"),
									},
								},
								"count": schema.StringAttribute{
									Optional:    true,
									Description: "Count the packet in the named counter.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 64),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"forwarding_class": schema.StringAttribute{
									Optional:    true,
									Description: "Classify packet to forwarding class.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 64),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"log": schema.BoolAttribute{
									Optional:    true,
									Description: "Log the packet.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"loss_priority": schema.StringAttribute{
									Optional:    true,
									Description: "Packet's loss priority.",
									Validators: []validator.String{
										stringvalidator.OneOf("high", "low", "medium-high", "medium-low"),
									},
								},
								"packet_mode": schema.BoolAttribute{
									Optional:    true,
									Description: "Bypass flow mode for the packet.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"policer": schema.StringAttribute{
									Optional:    true,
									Description: "Name of policer to use to rate-limit traffic.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"port_mirror": schema.BoolAttribute{
									Optional:    true,
									Description: "Port-mirror the packet.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"routing_instance": schema.StringAttribute{
									Optional:    true,
									Description: "Packets are directed to specified routing instance.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 63),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									},
								},
								"sample": schema.BoolAttribute{
									Optional:    true,
									Description: "Sample the packet.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"service_accounting": schema.BoolAttribute{
									Optional:    true,
									Description: "Count the packets for service accounting.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"syslog": schema.BoolAttribute{
									Optional:    true,
									Description: "System log (syslog) information about the packet.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
							},
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

type firewallFilterData struct {
	ID                types.String              `tfsdk:"id"`
	Name              types.String              `tfsdk:"name"`
	Family            types.String              `tfsdk:"family"`
	InterfaceSpecific types.Bool                `tfsdk:"interface_specific"`
	Term              []firewallFilterBlockTerm `tfsdk:"term"`
}

type firewallFilterConfig struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Family            types.String `tfsdk:"family"`
	InterfaceSpecific types.Bool   `tfsdk:"interface_specific"`
	Term              types.List   `tfsdk:"term"`
}

type firewallFilterBlockTerm struct {
	Name   types.String                      `tfsdk:"name"`
	Filter types.String                      `tfsdk:"filter"`
	From   *firewallFilterBlockTermBlockFrom `tfsdk:"from"`
	Then   *firewallFilterBlockTermBlockThen `tfsdk:"then"`
}

type firewallFilterBlockTermConfig struct {
	Name   types.String                            `tfsdk:"name"`
	Filter types.String                            `tfsdk:"filter"`
	From   *firewallFilterBlockTermBlockFromConfig `tfsdk:"from"`
	Then   *firewallFilterBlockTermBlockThen       `tfsdk:"then"`
}

func (block *firewallFilterBlockTermConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block, "Name")
}

type firewallFilterBlockTermBlockFrom struct {
	Address                     []types.String `tfsdk:"address"`
	AddressExcept               []types.String `tfsdk:"address_except"`
	DestinationAddress          []types.String `tfsdk:"destination_address"`
	DestinationAddressExcept    []types.String `tfsdk:"destination_address_except"`
	DestinationMacAddress       []types.String `tfsdk:"destination_mac_address"`
	DestinationMacAddressExcept []types.String `tfsdk:"destination_mac_address_except"`
	DestinationPort             []types.String `tfsdk:"destination_port"`
	DestinationPortExcept       []types.String `tfsdk:"destination_port_except"`
	DestinationPrefixList       []types.String `tfsdk:"destination_prefix_list"`
	DestinationPrefixListExcept []types.String `tfsdk:"destination_prefix_list_except"`
	ForwardingClass             []types.String `tfsdk:"forwarding_class"`
	ForwardingClassExcept       []types.String `tfsdk:"forwarding_class_except"`
	IcmpCode                    []types.String `tfsdk:"icmp_code"`
	IcmpCodeExcept              []types.String `tfsdk:"icmp_code_except"`
	IcmpType                    []types.String `tfsdk:"icmp_type"`
	IcmpTypeExcept              []types.String `tfsdk:"icmp_type_except"`
	Interface                   []types.String `tfsdk:"interface"`
	IsFragment                  types.Bool     `tfsdk:"is_fragment"`
	LossPriority                []types.String `tfsdk:"loss_priority"`
	LossPriorityExcept          []types.String `tfsdk:"loss_priority_except"`
	NextHeader                  []types.String `tfsdk:"next_header"`
	NextHeaderExcept            []types.String `tfsdk:"next_header_except"`
	PacketLength                []types.String `tfsdk:"packet_length"`
	PacketLengthExcept          []types.String `tfsdk:"packet_length_except"`
	PolicyMap                   []types.String `tfsdk:"policy_map"`
	PolicyMapExcept             []types.String `tfsdk:"policy_map_except"`
	Port                        []types.String `tfsdk:"port"`
	PortExcept                  []types.String `tfsdk:"port_except"`
	PrefixList                  []types.String `tfsdk:"prefix_list"`
	PrefixListExcept            []types.String `tfsdk:"prefix_list_except"`
	Protocol                    []types.String `tfsdk:"protocol"`
	ProtocolExcept              []types.String `tfsdk:"protocol_except"`
	SourceAddress               []types.String `tfsdk:"source_address"`
	SourceAddressExcept         []types.String `tfsdk:"source_address_except"`
	SourceMacAddress            []types.String `tfsdk:"source_mac_address"`
	SourceMacAddressExcept      []types.String `tfsdk:"source_mac_address_except"`
	SourcePort                  []types.String `tfsdk:"source_port"`
	SourcePortExcept            []types.String `tfsdk:"source_port_except"`
	SourcePrefixList            []types.String `tfsdk:"source_prefix_list"`
	SourcePrefixListExcept      []types.String `tfsdk:"source_prefix_list_except"`
	TCPEstablished              types.Bool     `tfsdk:"tcp_established"`
	TCPFlags                    types.String   `tfsdk:"tcp_flags"`
	TCPInitial                  types.Bool     `tfsdk:"tcp_initial"`
}

type firewallFilterBlockTermBlockFromConfig struct {
	Address                     types.Set    `tfsdk:"address"`
	AddressExcept               types.Set    `tfsdk:"address_except"`
	DestinationAddress          types.Set    `tfsdk:"destination_address"`
	DestinationAddressExcept    types.Set    `tfsdk:"destination_address_except"`
	DestinationMacAddress       types.Set    `tfsdk:"destination_mac_address"`
	DestinationMacAddressExcept types.Set    `tfsdk:"destination_mac_address_except"`
	DestinationPort             types.Set    `tfsdk:"destination_port"`
	DestinationPortExcept       types.Set    `tfsdk:"destination_port_except"`
	DestinationPrefixList       types.Set    `tfsdk:"destination_prefix_list"`
	DestinationPrefixListExcept types.Set    `tfsdk:"destination_prefix_list_except"`
	ForwardingClass             types.Set    `tfsdk:"forwarding_class"`
	ForwardingClassExcept       types.Set    `tfsdk:"forwarding_class_except"`
	IcmpCode                    types.Set    `tfsdk:"icmp_code"`
	IcmpCodeExcept              types.Set    `tfsdk:"icmp_code_except"`
	IcmpType                    types.Set    `tfsdk:"icmp_type"`
	IcmpTypeExcept              types.Set    `tfsdk:"icmp_type_except"`
	Interface                   types.Set    `tfsdk:"interface"`
	IsFragment                  types.Bool   `tfsdk:"is_fragment"`
	LossPriority                types.Set    `tfsdk:"loss_priority"`
	LossPriorityExcept          types.Set    `tfsdk:"loss_priority_except"`
	NextHeader                  types.Set    `tfsdk:"next_header"`
	NextHeaderExcept            types.Set    `tfsdk:"next_header_except"`
	PacketLength                types.Set    `tfsdk:"packet_length"`
	PacketLengthExcept          types.Set    `tfsdk:"packet_length_except"`
	PolicyMap                   types.Set    `tfsdk:"policy_map"`
	PolicyMapExcept             types.Set    `tfsdk:"policy_map_except"`
	Port                        types.Set    `tfsdk:"port"`
	PortExcept                  types.Set    `tfsdk:"port_except"`
	PrefixList                  types.Set    `tfsdk:"prefix_list"`
	PrefixListExcept            types.Set    `tfsdk:"prefix_list_except"`
	Protocol                    types.Set    `tfsdk:"protocol"`
	ProtocolExcept              types.Set    `tfsdk:"protocol_except"`
	SourceAddress               types.Set    `tfsdk:"source_address"`
	SourceAddressExcept         types.Set    `tfsdk:"source_address_except"`
	SourceMacAddress            types.Set    `tfsdk:"source_mac_address"`
	SourceMacAddressExcept      types.Set    `tfsdk:"source_mac_address_except"`
	SourcePort                  types.Set    `tfsdk:"source_port"`
	SourcePortExcept            types.Set    `tfsdk:"source_port_except"`
	SourcePrefixList            types.Set    `tfsdk:"source_prefix_list"`
	SourcePrefixListExcept      types.Set    `tfsdk:"source_prefix_list_except"`
	TCPEstablished              types.Bool   `tfsdk:"tcp_established"`
	TCPFlags                    types.String `tfsdk:"tcp_flags"`
	TCPInitial                  types.Bool   `tfsdk:"tcp_initial"`
}

func (block *firewallFilterBlockTermBlockFromConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type firewallFilterBlockTermBlockThen struct {
	Action            types.String `tfsdk:"action"`
	Count             types.String `tfsdk:"count"`
	ForwardingClass   types.String `tfsdk:"forwarding_class"`
	Log               types.Bool   `tfsdk:"log"`
	LossPriority      types.String `tfsdk:"loss_priority"`
	PacketMode        types.Bool   `tfsdk:"packet_mode"`
	Policer           types.String `tfsdk:"policer"`
	PortMirror        types.Bool   `tfsdk:"port_mirror"`
	RoutingInstance   types.String `tfsdk:"routing_instance"`
	Sample            types.Bool   `tfsdk:"sample"`
	ServiceAccounting types.Bool   `tfsdk:"service_accounting"`
	Syslog            types.Bool   `tfsdk:"syslog"`
}

func (block *firewallFilterBlockTermBlockThen) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (rsc *firewallFilter) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config firewallFilterConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Term.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("term"),
			tfdiag.MissingConfigErrSummary,
			"term block must be specified",
		)
	} else if !config.Term.IsUnknown() {
		var configTerm []firewallFilterBlockTermConfig
		asDiags := config.Term.ElementsAs(ctx, &configTerm, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		termName := make(map[string]struct{})
		for i, block := range configTerm {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := termName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple term blocks with the same name %q", name),
					)
				}
				termName[name] = struct{}{}
			}
			if block.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("term").AtListIndex(i).AtName("name"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("term block %q is empty", block.Name.ValueString()),
				)
			}
			if block.From != nil {
				if block.From.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("from block in term block %q is empty", block.Name.ValueString()),
					)
				}
				if !config.Family.IsNull() && !config.Family.IsUnknown() {
					block.From.validateWithFamily(
						ctx,
						config.Family.ValueString(),
						path.Root("term").AtListIndex(i).AtName("from"),
						resp,
					)
				}
				if !block.From.DestinationPort.IsNull() && !block.From.DestinationPort.IsUnknown() &&
					!block.From.DestinationPortExcept.IsNull() && !block.From.DestinationPortExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("destination_port"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("destination_port and destination_port_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.IcmpCode.IsNull() && !block.From.IcmpCode.IsUnknown() &&
					!block.From.IcmpCodeExcept.IsNull() && !block.From.IcmpCodeExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("icmp_code"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("icmp_code and icmp_code_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.IcmpType.IsNull() && !block.From.IcmpType.IsUnknown() &&
					!block.From.IcmpTypeExcept.IsNull() && !block.From.IcmpTypeExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("icmp_type"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("icmp_type and icmp_type_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.ForwardingClass.IsNull() && !block.From.ForwardingClass.IsUnknown() &&
					!block.From.ForwardingClassExcept.IsNull() && !block.From.ForwardingClassExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("forwarding_class"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("forwarding_class and forwarding_class_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.LossPriority.IsNull() && !block.From.LossPriority.IsUnknown() &&
					!block.From.LossPriorityExcept.IsNull() && !block.From.LossPriorityExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("loss_priority"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("loss_priority and loss_priority_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.NextHeader.IsNull() && !block.From.NextHeader.IsUnknown() &&
					!block.From.NextHeaderExcept.IsNull() && !block.From.NextHeaderExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("next_header"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("next_header and next_header_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.PacketLength.IsNull() && !block.From.PacketLength.IsUnknown() &&
					!block.From.PacketLengthExcept.IsNull() && !block.From.PacketLengthExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("packet_length"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("packet_length and packet_length_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.PolicyMap.IsNull() && !block.From.PolicyMap.IsUnknown() &&
					!block.From.PolicyMapExcept.IsNull() && !block.From.PolicyMapExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("policy_map"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("policy_map and policy_map_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.Port.IsNull() && !block.From.Port.IsUnknown() &&
					!block.From.PortExcept.IsNull() && !block.From.PortExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("port"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("port and port_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.Protocol.IsNull() && !block.From.Protocol.IsUnknown() &&
					!block.From.ProtocolExcept.IsNull() && !block.From.ProtocolExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("protocol"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("protocol and protocol_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.SourcePort.IsNull() && !block.From.SourcePort.IsUnknown() &&
					!block.From.SourcePortExcept.IsNull() && !block.From.SourcePortExcept.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("source_port"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("source_port and source_port_except cannot be configured together"+
							" in from block in term block %q", block.Name.ValueString()),
					)
				}
			}
			if block.Then != nil {
				if block.Then.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("then"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("then block in term block %q is empty", block.Name.ValueString()),
					)
				}
			}
		}
	}
}

//nolint:gocyclo
func (block *firewallFilterBlockTermBlockFromConfig) validateWithFamily(
	_ context.Context, family string, pathRoot path.Path, resp *resource.ValidateConfigResponse,
) {
	errorMessageWithFamilySuffix := fmt.Sprintf(" in from block cannot be configured with family %q", family)
	if !block.Address.IsNull() && !block.Address.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("address"),
			tfdiag.ConflictConfigErrSummary,
			"address"+errorMessageWithFamilySuffix,
		)
	}
	if !block.AddressExcept.IsNull() && !block.AddressExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("address_except"),
			tfdiag.ConflictConfigErrSummary,
			"address_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.DestinationAddress.IsNull() && !block.DestinationAddress.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W,
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_address"),
			tfdiag.ConflictConfigErrSummary,
			"destination_address"+errorMessageWithFamilySuffix,
		)
	}
	if !block.DestinationAddressExcept.IsNull() && !block.DestinationAddressExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W,
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_address_except"),
			tfdiag.ConflictConfigErrSummary,
			"destination_address_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.DestinationMacAddress.IsNull() && !block.DestinationMacAddress.IsUnknown() &&
		!slices.Contains([]string{
			"vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_mac_address"),
			tfdiag.ConflictConfigErrSummary,
			"destination_mac_address"+errorMessageWithFamilySuffix,
		)
	}
	if !block.DestinationMacAddressExcept.IsNull() && !block.DestinationMacAddressExcept.IsUnknown() &&
		!slices.Contains([]string{
			"vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_mac_address_except"),
			tfdiag.ConflictConfigErrSummary,
			"destination_mac_address_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.DestinationPort.IsNull() && !block.DestinationPort.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_port"),
			tfdiag.ConflictConfigErrSummary,
			"destination_port"+errorMessageWithFamilySuffix,
		)
	}
	if !block.DestinationPortExcept.IsNull() && !block.DestinationPortExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_port_except"),
			tfdiag.ConflictConfigErrSummary,
			"destination_port_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.DestinationPrefixList.IsNull() && !block.DestinationPrefixList.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_prefix_list"),
			tfdiag.ConflictConfigErrSummary,
			"destination_prefix_list"+errorMessageWithFamilySuffix,
		)
	}
	if !block.DestinationPrefixListExcept.IsNull() && !block.DestinationPrefixListExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_prefix_list_except"),
			tfdiag.ConflictConfigErrSummary,
			"destination_prefix_list_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.ForwardingClass.IsNull() && !block.ForwardingClass.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "any", "ccc", "mpls", "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("forwarding_class"),
			tfdiag.ConflictConfigErrSummary,
			"forwarding_class"+errorMessageWithFamilySuffix,
		)
	}
	if !block.ForwardingClassExcept.IsNull() && !block.ForwardingClassExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "any", "ccc", "mpls", "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("forwarding_class_except"),
			tfdiag.ConflictConfigErrSummary,
			"forwarding_class_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.IcmpCode.IsNull() && !block.IcmpCode.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("icmp_code"),
			tfdiag.ConflictConfigErrSummary,
			"icmp_code"+errorMessageWithFamilySuffix,
		)
	}
	if !block.IcmpCodeExcept.IsNull() && !block.IcmpCodeExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("icmp_code_except"),
			tfdiag.ConflictConfigErrSummary,
			"icmp_code_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.IcmpType.IsNull() && !block.IcmpType.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("icmp_type"),
			tfdiag.ConflictConfigErrSummary,
			"icmp_type"+errorMessageWithFamilySuffix,
		)
	}
	if !block.IcmpTypeExcept.IsNull() && !block.IcmpTypeExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("icmp_type_except"),
			tfdiag.ConflictConfigErrSummary,
			"icmp_type_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.Interface.IsNull() && !block.Interface.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "any", "mpls", "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("interface"),
			tfdiag.ConflictConfigErrSummary,
			"interface"+errorMessageWithFamilySuffix,
		)
	}
	if !block.IsFragment.IsNull() && !block.IsFragment.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("is_fragment"),
			tfdiag.ConflictConfigErrSummary,
			"is_fragment"+errorMessageWithFamilySuffix,
		)
	}
	if !block.LossPriority.IsNull() && !block.LossPriority.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "any", "ccc", "mpls", "vpls",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("loss_priority"),
			tfdiag.ConflictConfigErrSummary,
			"loss_priority"+errorMessageWithFamilySuffix,
		)
	}
	if !block.LossPriorityExcept.IsNull() && !block.LossPriorityExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "any", "ccc", "mpls", "vpls",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("loss_priority_except"),
			tfdiag.ConflictConfigErrSummary,
			"loss_priority_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.NextHeader.IsNull() && !block.NextHeader.IsUnknown() &&
		!slices.Contains([]string{
			junos.Inet6W,
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("next_header"),
			tfdiag.ConflictConfigErrSummary,
			"next_header"+errorMessageWithFamilySuffix,
		)
	}
	if !block.NextHeaderExcept.IsNull() && !block.NextHeaderExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.Inet6W,
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("next_header_except"),
			tfdiag.ConflictConfigErrSummary,
			"next_header_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.PacketLength.IsNull() && !block.PacketLength.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "any",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("packet_length"),
			tfdiag.ConflictConfigErrSummary,
			"packet_length"+errorMessageWithFamilySuffix,
		)
	}
	if !block.PacketLengthExcept.IsNull() && !block.PacketLengthExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "any",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("packet_length_except"),
			tfdiag.ConflictConfigErrSummary,
			"packet_length_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.PolicyMap.IsNull() && !block.PolicyMap.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "any", "ccc", "mpls", "vpls",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("policy_map"),
			tfdiag.ConflictConfigErrSummary,
			"policy_map"+errorMessageWithFamilySuffix,
		)
	}
	if !block.PolicyMapExcept.IsNull() && !block.PolicyMapExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "any", "ccc", "mpls", "vpls",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("policy_map_except"),
			tfdiag.ConflictConfigErrSummary,
			"policy_map_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.Port.IsNull() && !block.Port.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("port"),
			tfdiag.ConflictConfigErrSummary,
			"port"+errorMessageWithFamilySuffix,
		)
	}
	if !block.PortExcept.IsNull() && !block.PortExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("port_except"),
			tfdiag.ConflictConfigErrSummary,
			"port_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.PrefixList.IsNull() && !block.PrefixList.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("prefix_list"),
			tfdiag.ConflictConfigErrSummary,
			"prefix_list"+errorMessageWithFamilySuffix,
		)
	}
	if !block.PrefixListExcept.IsNull() && !block.PrefixListExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("prefix_list_except"),
			tfdiag.ConflictConfigErrSummary,
			"prefix_list_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.Protocol.IsNull() && !block.Protocol.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("protocol"),
			tfdiag.ConflictConfigErrSummary,
			"protocol"+errorMessageWithFamilySuffix,
		)
	}
	if !block.ProtocolExcept.IsNull() && !block.ProtocolExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("protocol_except"),
			tfdiag.ConflictConfigErrSummary,
			"protocol_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.SourceAddress.IsNull() && !block.SourceAddress.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W,
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_address"),
			tfdiag.ConflictConfigErrSummary,
			"source_address"+errorMessageWithFamilySuffix,
		)
	}
	if !block.SourceAddressExcept.IsNull() && !block.SourceAddressExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W,
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_address_except"),
			tfdiag.ConflictConfigErrSummary,
			"source_address_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.SourceMacAddress.IsNull() && !block.SourceMacAddress.IsUnknown() &&
		!slices.Contains([]string{
			"vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_mac_address"),
			tfdiag.ConflictConfigErrSummary,
			"source_mac_address"+errorMessageWithFamilySuffix,
		)
	}
	if !block.SourceMacAddressExcept.IsNull() && !block.SourceMacAddressExcept.IsUnknown() &&
		!slices.Contains([]string{
			"vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_mac_address_except"),
			tfdiag.ConflictConfigErrSummary,
			"source_mac_address_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.SourcePort.IsNull() && !block.SourcePort.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_port"),
			tfdiag.ConflictConfigErrSummary,
			"source_port"+errorMessageWithFamilySuffix,
		)
	}
	if !block.SourcePortExcept.IsNull() && !block.SourcePortExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_port_except"),
			tfdiag.ConflictConfigErrSummary,
			"source_port_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.SourcePrefixList.IsNull() && !block.SourcePrefixList.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_prefix_list"),
			tfdiag.ConflictConfigErrSummary,
			"source_prefix_list"+errorMessageWithFamilySuffix,
		)
	}
	if !block.SourcePrefixListExcept.IsNull() && !block.SourcePrefixListExcept.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_prefix_list_except"),
			tfdiag.ConflictConfigErrSummary,
			"source_prefix_list_except"+errorMessageWithFamilySuffix,
		)
	}
	if !block.TCPEstablished.IsNull() && !block.TCPEstablished.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("tcp_established"),
			tfdiag.ConflictConfigErrSummary,
			"tcp_established"+errorMessageWithFamilySuffix,
		)
	}
	if !block.TCPFlags.IsNull() && !block.TCPFlags.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "vpls", "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("tcp_flags"),
			tfdiag.ConflictConfigErrSummary,
			"tcp_flags"+errorMessageWithFamilySuffix,
		)
	}
	if !block.TCPInitial.IsNull() && !block.TCPInitial.IsUnknown() &&
		!slices.Contains([]string{
			junos.InetW, junos.Inet6W, "ethernet-switching",
		}, family) {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("tcp_initial"),
			tfdiag.ConflictConfigErrSummary,
			"tcp_initial"+errorMessageWithFamilySuffix,
		)
	}
}

func (rsc *firewallFilter) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan firewallFilterData
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
			filterExists, err := checkFirewallFilterExists(fnCtx, plan.Name.ValueString(), plan.Family.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if filterExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			filterExists, err := checkFirewallFilterExists(fnCtx, plan.Name.ValueString(), plan.Family.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !filterExists {
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

func (rsc *firewallFilter) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data firewallFilterData
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
			state.Family.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *firewallFilter) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state firewallFilterData
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

func (rsc *firewallFilter) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state firewallFilterData
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

func (rsc *firewallFilter) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data firewallFilterData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <name>"+junos.IDSeparator+"<family>)",
	)
}

func checkFirewallFilterExists(
	_ context.Context, name, family string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"firewall family " + family + " filter \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *firewallFilterData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + rscData.Family.ValueString())
}

func (rscData *firewallFilterData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *firewallFilterData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set firewall family " + rscData.Family.ValueString() + " filter \"" + rscData.Name.ValueString() + "\" "

	if rscData.InterfaceSpecific.ValueBool() {
		configSet = append(configSet, setPrefix+"interface-specific")
	}
	termName := make(map[string]struct{})
	for i, block := range rscData.Term {
		name := block.Name.ValueString()
		if _, ok := termName[name]; ok {
			return path.Root("term").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple term blocks with the same name %q", name)
		}
		termName[name] = struct{}{}

		setPrefixTerm := setPrefix + "term \"" + name + "\" "
		if v := block.Filter.ValueString(); v != "" {
			configSet = append(configSet, setPrefixTerm+"filter \""+v+"\"")
		}
		if block.From != nil {
			blockSet, pathErr, err := block.From.configSet(setPrefixTerm, path.Root("term").AtListIndex(i).AtName("from"))
			if err != nil {
				return pathErr, err
			}
			configSet = append(configSet, blockSet...)
		}
		if block.Then != nil {
			blockSet := block.Then.configSet(setPrefixTerm)
			configSet = append(configSet, blockSet...)
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *firewallFilterBlockTermBlockFrom) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, 100)
	setPrefix += "from "

	for _, v := range block.Address {
		configSet = append(configSet, setPrefix+"address "+v.ValueString())
	}
	for _, v := range block.AddressExcept {
		configSet = append(configSet, setPrefix+"address "+v.ValueString()+" except")
	}
	for _, v := range block.DestinationAddress {
		configSet = append(configSet, setPrefix+"destination-address "+v.ValueString())
	}
	for _, v := range block.DestinationAddressExcept {
		configSet = append(configSet, setPrefix+"destination-address "+v.ValueString()+" except")
	}
	for _, v := range block.DestinationMacAddress {
		configSet = append(configSet, setPrefix+"destination-mac-address "+v.ValueString())
	}
	for _, v := range block.DestinationMacAddressExcept {
		configSet = append(configSet, setPrefix+"destination-mac-address "+v.ValueString()+" except")
	}
	if len(block.DestinationPort) > 0 && len(block.DestinationPortExcept) > 0 {
		return configSet,
			pathRoot.AtName("destination_port"),
			errors.New("destination_port and destination_port_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.DestinationPort {
		configSet = append(configSet, setPrefix+"destination-port "+v.ValueString())
	}
	for _, v := range block.DestinationPortExcept {
		configSet = append(configSet, setPrefix+"destination-port-except "+v.ValueString())
	}
	for _, v := range block.DestinationPrefixList {
		configSet = append(configSet, setPrefix+"destination-prefix-list \""+v.ValueString()+"\"")
	}
	for _, v := range block.DestinationPrefixListExcept {
		configSet = append(configSet, setPrefix+"destination-prefix-list \""+v.ValueString()+"\" except")
	}
	if len(block.ForwardingClass) > 0 && len(block.ForwardingClassExcept) > 0 {
		return configSet,
			pathRoot.AtName("forwarding_class"),
			errors.New("forwarding_class and forwarding_class_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.ForwardingClass {
		configSet = append(configSet, setPrefix+"forwarding-class "+v.ValueString())
	}
	for _, v := range block.ForwardingClassExcept {
		configSet = append(configSet, setPrefix+"forwarding-class-except "+v.ValueString())
	}
	if len(block.IcmpCode) > 0 && len(block.IcmpCodeExcept) > 0 {
		return configSet,
			pathRoot.AtName("icmp_code"),
			errors.New("icmp_code and icmp_code_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.IcmpCode {
		configSet = append(configSet, setPrefix+"icmp-code "+v.ValueString())
	}
	for _, v := range block.IcmpCodeExcept {
		configSet = append(configSet, setPrefix+"icmp-code-except "+v.ValueString())
	}
	if len(block.IcmpType) > 0 && len(block.IcmpTypeExcept) > 0 {
		return configSet,
			pathRoot.AtName("icmp_type"),
			errors.New("icmp_type and icmp_type_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.IcmpType {
		configSet = append(configSet, setPrefix+"icmp-type "+v.ValueString())
	}
	for _, v := range block.IcmpTypeExcept {
		configSet = append(configSet, setPrefix+"icmp-type-except "+v.ValueString())
	}
	for _, v := range block.Interface {
		configSet = append(configSet, setPrefix+"interface "+v.ValueString())
	}
	if block.IsFragment.ValueBool() {
		configSet = append(configSet, setPrefix+"is-fragment")
	}
	if len(block.LossPriority) > 0 && len(block.LossPriorityExcept) > 0 {
		return configSet,
			pathRoot.AtName("loss_priority"),
			errors.New("loss_priority and loss_priority_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.LossPriority {
		configSet = append(configSet, setPrefix+"loss-priority "+v.ValueString())
	}
	for _, v := range block.LossPriorityExcept {
		configSet = append(configSet, setPrefix+"loss-priority-except "+v.ValueString())
	}
	if len(block.NextHeader) > 0 && len(block.NextHeaderExcept) > 0 {
		return configSet,
			pathRoot.AtName("next_header"),
			errors.New("next_header and next_header_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.NextHeader {
		configSet = append(configSet, setPrefix+"next-header "+v.ValueString())
	}
	for _, v := range block.NextHeaderExcept {
		configSet = append(configSet, setPrefix+"next-header-except "+v.ValueString())
	}
	if len(block.PacketLength) > 0 && len(block.PacketLengthExcept) > 0 {
		return configSet,
			pathRoot.AtName("packet_length"),
			errors.New("packet_length and packet_length_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.PacketLength {
		configSet = append(configSet, setPrefix+"packet-length "+v.ValueString())
	}
	for _, v := range block.PacketLengthExcept {
		configSet = append(configSet, setPrefix+"packet-length-except "+v.ValueString())
	}
	if len(block.PolicyMap) > 0 && len(block.PolicyMapExcept) > 0 {
		return configSet,
			pathRoot.AtName("policy_map"),
			errors.New("policy_map and policy_map_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.PolicyMap {
		configSet = append(configSet, setPrefix+"policy-map "+v.ValueString())
	}
	for _, v := range block.PolicyMapExcept {
		configSet = append(configSet, setPrefix+"policy-map-except "+v.ValueString())
	}
	if len(block.Port) > 0 && len(block.PortExcept) > 0 {
		return configSet,
			pathRoot.AtName("port"),
			errors.New("port and port_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.Port {
		configSet = append(configSet, setPrefix+"port "+v.ValueString())
	}
	for _, v := range block.PortExcept {
		configSet = append(configSet, setPrefix+"port-except "+v.ValueString())
	}
	for _, v := range block.PrefixList {
		configSet = append(configSet, setPrefix+"prefix-list \""+v.ValueString()+"\"")
	}
	for _, v := range block.PrefixListExcept {
		configSet = append(configSet, setPrefix+"prefix-list \""+v.ValueString()+"\" except")
	}
	if len(block.Protocol) > 0 && len(block.ProtocolExcept) > 0 {
		return configSet,
			pathRoot.AtName("protocol"),
			errors.New("protocol and protocol_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.Protocol {
		configSet = append(configSet, setPrefix+"protocol "+v.ValueString())
	}
	for _, v := range block.ProtocolExcept {
		configSet = append(configSet, setPrefix+"protocol-except "+v.ValueString())
	}
	for _, v := range block.SourceAddress {
		configSet = append(configSet, setPrefix+"source-address "+v.ValueString())
	}
	for _, v := range block.SourceAddressExcept {
		configSet = append(configSet, setPrefix+"source-address "+v.ValueString()+" except")
	}
	for _, v := range block.SourceMacAddress {
		configSet = append(configSet, setPrefix+"source-mac-address "+v.ValueString())
	}
	for _, v := range block.SourceMacAddressExcept {
		configSet = append(configSet, setPrefix+"source-mac-address "+v.ValueString()+" except")
	}
	if len(block.SourcePort) > 0 && len(block.SourcePortExcept) > 0 {
		return configSet,
			pathRoot.AtName("source_port"),
			errors.New("source_port and source_port_except cannot be configured together" +
				" in from block")
	}
	for _, v := range block.SourcePort {
		configSet = append(configSet, setPrefix+"source-port "+v.ValueString())
	}
	for _, v := range block.SourcePortExcept {
		configSet = append(configSet, setPrefix+"source-port-except "+v.ValueString())
	}
	for _, v := range block.SourcePrefixList {
		configSet = append(configSet, setPrefix+"source-prefix-list \""+v.ValueString()+"\"")
	}
	for _, v := range block.SourcePrefixListExcept {
		configSet = append(configSet, setPrefix+"source-prefix-list \""+v.ValueString()+"\" except")
	}
	if block.TCPEstablished.ValueBool() {
		if block.TCPFlags.ValueString() != "" {
			return configSet,
				pathRoot.AtName("tcp_established"),
				errors.New("tcp_established and tcp_flags cannot be configured together" +
					" in from block")
		}
		if block.TCPInitial.ValueBool() {
			return configSet,
				pathRoot.AtName("tcp_established"),
				errors.New("tcp_established and tcp_initial cannot be configured together" +
					" in from block")
		}
		configSet = append(configSet, setPrefix+"tcp-established")
	}
	if v := block.TCPFlags.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"tcp-flags \""+v+"\"")
	}
	if block.TCPInitial.ValueBool() {
		configSet = append(configSet, setPrefix+"tcp-initial")
	}

	return configSet, path.Empty(), nil
}

func (block *firewallFilterBlockTermBlockThen) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 1)
	setPrefix += "then "

	if v := block.Action.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+v)
	}
	if v := block.Count.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"count \""+v+"\"")
	}
	if v := block.ForwardingClass.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"forwarding-class \""+v+"\"")
	}
	if block.Log.ValueBool() {
		configSet = append(configSet, setPrefix+"log")
	}
	if v := block.LossPriority.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"loss-priority "+v)
	}
	if block.PacketMode.ValueBool() {
		configSet = append(configSet, setPrefix+"packet-mode")
	}
	if v := block.Policer.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"policer \""+v+"\"")
	}
	if block.PortMirror.ValueBool() {
		configSet = append(configSet, setPrefix+"port-mirror")
	}
	if v := block.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-instance "+v)
	}
	if block.Sample.ValueBool() {
		configSet = append(configSet, setPrefix+"sample")
	}
	if block.ServiceAccounting.ValueBool() {
		configSet = append(configSet, setPrefix+"service-accounting")
	}
	if block.Syslog.ValueBool() {
		configSet = append(configSet, setPrefix+"syslog")
	}

	return configSet
}

func (rscData *firewallFilterData) read(
	_ context.Context, name, family string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"firewall family " + family + " filter \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.Family = types.StringValue(family)
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
			case itemTrim == "interface-specific":
				rscData.InterfaceSpecific = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "term "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				var term firewallFilterBlockTerm
				rscData.Term, term = tfdata.ExtractBlockWithTFTypesString(
					rscData.Term, "Name", strings.Trim(name, "\""))
				term.Name = types.StringValue(strings.Trim(name, "\""))
				balt.CutPrefixInString(&itemTrim, name+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "filter "):
					term.Filter = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "from "):
					if term.From == nil {
						term.From = &firewallFilterBlockTermBlockFrom{}
					}
					term.From.read(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "then "):
					if term.Then == nil {
						term.Then = &firewallFilterBlockTermBlockThen{}
					}
					term.Then.read(itemTrim)
				}
				rscData.Term = append(rscData.Term, term)
			}
		}
	}

	return nil
}

func (block *firewallFilterBlockTermBlockFrom) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "address "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			block.AddressExcept = append(block.AddressExcept, types.StringValue(itemTrim))
		} else {
			block.Address = append(block.Address, types.StringValue(itemTrim))
		}
	case balt.CutPrefixInString(&itemTrim, "destination-address "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			block.DestinationAddressExcept = append(block.DestinationAddressExcept, types.StringValue(itemTrim))
		} else {
			block.DestinationAddress = append(block.DestinationAddress, types.StringValue(itemTrim))
		}
	case balt.CutPrefixInString(&itemTrim, "destination-mac-address "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			block.DestinationMacAddressExcept = append(block.DestinationMacAddressExcept, types.StringValue(itemTrim))
		} else {
			block.DestinationMacAddress = append(block.DestinationMacAddress, types.StringValue(itemTrim))
		}
	case balt.CutPrefixInString(&itemTrim, "destination-port "):
		block.DestinationPort = append(block.DestinationPort, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "destination-port-except "):
		block.DestinationPortExcept = append(block.DestinationPortExcept, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "destination-prefix-list "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			block.DestinationPrefixListExcept = append(block.DestinationPrefixListExcept,
				types.StringValue(strings.Trim(itemTrim, "\"")))
		} else {
			block.DestinationPrefixList = append(block.DestinationPrefixList, types.StringValue(strings.Trim(itemTrim, "\"")))
		}
	case balt.CutPrefixInString(&itemTrim, "forwarding-class "):
		block.ForwardingClass = append(block.ForwardingClass, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "forwarding-class-except "):
		block.ForwardingClassExcept = append(block.ForwardingClassExcept, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "icmp-code "):
		block.IcmpCode = append(block.IcmpCode, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "icmp-code-except "):
		block.IcmpCodeExcept = append(block.IcmpCodeExcept, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "icmp-type "):
		block.IcmpType = append(block.IcmpType, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "icmp-type-except "):
		block.IcmpTypeExcept = append(block.IcmpTypeExcept, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "interface "):
		block.Interface = append(block.Interface, types.StringValue(itemTrim))
	case itemTrim == "is-fragment":
		block.IsFragment = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "loss-priority "):
		block.LossPriority = append(block.LossPriority, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "loss-priority-except "):
		block.LossPriorityExcept = append(block.LossPriorityExcept, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "next-header "):
		block.NextHeader = append(block.NextHeader, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "next-header-except "):
		block.NextHeaderExcept = append(block.NextHeaderExcept, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "port "):
		block.Port = append(block.Port, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "port-except "):
		block.PortExcept = append(block.PortExcept, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "packet-length "):
		block.PacketLength = append(block.PacketLength, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "packet-length-except "):
		block.PacketLengthExcept = append(block.PacketLengthExcept, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "policy-map "):
		block.PolicyMap = append(block.PolicyMap, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "policy-map-except "):
		block.PolicyMapExcept = append(block.PolicyMapExcept, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		block.Protocol = append(block.Protocol, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "protocol-except "):
		block.ProtocolExcept = append(block.ProtocolExcept, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "prefix-list "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			block.PrefixListExcept = append(block.PrefixListExcept, types.StringValue(strings.Trim(itemTrim, "\"")))
		} else {
			block.PrefixList = append(block.PrefixList, types.StringValue(strings.Trim(itemTrim, "\"")))
		}
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			block.SourceAddressExcept = append(block.SourceAddressExcept, types.StringValue(itemTrim))
		} else {
			block.SourceAddress = append(block.SourceAddress, types.StringValue(itemTrim))
		}
	case balt.CutPrefixInString(&itemTrim, "source-mac-address "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			block.SourceMacAddressExcept = append(block.SourceMacAddressExcept, types.StringValue(itemTrim))
		} else {
			block.SourceMacAddress = append(block.SourceMacAddress, types.StringValue(itemTrim))
		}
	case balt.CutPrefixInString(&itemTrim, "source-port "):
		block.SourcePort = append(block.SourcePort, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "source-port-except "):
		block.SourcePortExcept = append(block.SourcePortExcept, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "source-prefix-list "):
		if balt.CutSuffixInString(&itemTrim, " except") {
			block.SourcePrefixListExcept = append(block.SourcePrefixListExcept, types.StringValue(strings.Trim(itemTrim, "\"")))
		} else {
			block.SourcePrefixList = append(block.SourcePrefixList, types.StringValue(strings.Trim(itemTrim, "\"")))
		}
	case itemTrim == "tcp-established":
		block.TCPEstablished = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "tcp-flags "):
		block.TCPFlags = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "tcp-initial":
		block.TCPInitial = types.BoolValue(true)
	}
}

func (block *firewallFilterBlockTermBlockThen) read(itemTrim string) {
	switch {
	case itemTrim == "accept",
		itemTrim == "reject",
		itemTrim == junos.DiscardW,
		itemTrim == "next term":
		block.Action = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "count "):
		block.Count = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "forwarding-class "):
		block.ForwardingClass = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "log":
		block.Log = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "loss-priority "):
		block.LossPriority = types.StringValue(itemTrim)
	case itemTrim == "packet-mode":
		block.PacketMode = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "policer "):
		block.Policer = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "port-mirror":
		block.PortMirror = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "routing-instance "):
		block.RoutingInstance = types.StringValue(itemTrim)
	case itemTrim == "sample":
		block.Sample = types.BoolValue(true)
	case itemTrim == "service-accounting":
		block.ServiceAccounting = types.BoolValue(true)
	case itemTrim == "syslog":
		block.Syslog = types.BoolValue(true)
	}
}

func (rscData *firewallFilterData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete firewall family " + rscData.Family.ValueString() + " filter \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
