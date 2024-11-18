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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &securityScreen{}
	_ resource.ResourceWithConfigure      = &securityScreen{}
	_ resource.ResourceWithValidateConfig = &securityScreen{}
	_ resource.ResourceWithImportState    = &securityScreen{}
	_ resource.ResourceWithUpgradeState   = &securityScreen{}
)

type securityScreen struct {
	client *junos.Client
}

func newSecurityScreenResource() resource.Resource {
	return &securityScreen{}
}

func (rsc *securityScreen) typeName() string {
	return providerName + "_security_screen"
}

func (rsc *securityScreen) junosName() string {
	return "security screen ids-option"
}

func (rsc *securityScreen) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityScreen) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityScreen) Configure(
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

//nolint:gochecknoglobals
var (
	securityScreenUserDefinedOptionTypeValidator = stringvalidator.RegexMatches(
		regexp.MustCompile(`^([1-9]|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])( to ([1-9]|[1-9]\d|1\d\d|2[0-4]\d|25[0-5]))?$`),
		"must match '(1..255)' or '(1..255) to (1..255)'",
	)
	securityScreenUserDefinedHeaderTypeValidator = stringvalidator.RegexMatches(
		regexp.MustCompile(`^(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])( to ([1-9]|[1-9]\d|1\d\d|2[0-4]\d|25[0-5]))?$`),
		"must match '(0..255)' or '(0..255) to (0..255)'",
	)
)

func (rsc *securityScreen) Schema(
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
				Description: "The name of screen.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"alarm_without_drop": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not drop packet, only generate alarm.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of screen.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"icmp": schema.SingleNestedBlock{
				Description: "Configure ICMP ids options.",
				Attributes: map[string]schema.Attribute{
					"fragment": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable ICMP fragment ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"icmpv6_malformed": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable ICMPv6 malformed ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"large": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable large ICMP packet (size > 1024) ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ping_death": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable ping of death ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"flood": schema.SingleNestedBlock{
						Description: "Enable ICMP flood ids option.",
						Attributes: map[string]schema.Attribute{
							"threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Threshold (ICMP packets per second).",
								Validators: []validator.Int64{
									int64validator.Between(1, 1000000),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"sweep": schema.SingleNestedBlock{
						Description: "Enable ICMP sweep ids option.",
						Attributes: map[string]schema.Attribute{
							"threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Threshold (microseconds in which 10 ICMP packets are detected).",
								Validators: []validator.Int64{
									int64validator.Between(1000, 1000000),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"ip": schema.SingleNestedBlock{
				Description: "Configure IP layer ids options.",
				Attributes: map[string]schema.Attribute{
					"bad_option": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP with bad option ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"block_frag": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP fragment blocking ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ipv6_extension_header_limit": schema.Int64Attribute{
						Optional:    true,
						Description: "Enable IPv6 extension header limit ids option.",
						Validators: []validator.Int64{
							int64validator.Between(0, 32),
						},
					},
					"ipv6_malformed_header": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IPv6 malformed header ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"loose_source_route_option": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP with loose source route ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"record_route_option": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP with record route option ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"security_option": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP with security option ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"source_route_option": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP source route ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"spoofing": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP address spoofing ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"stream_option": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP with stream option ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"strict_source_route_option": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP with strict source route ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"tear_drop": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable tear drop ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"timestamp_option": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP with timestamp option ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"unknown_protocol": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IP unknown protocol ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"ipv6_extension_header": schema.SingleNestedBlock{
						Description: "Configure ipv6 extension header ids option.",
						Attributes: map[string]schema.Attribute{
							"ah_header": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable IPv6 Authentication Header ids option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"esp_header": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable IPv6 Encapsulating Security Payload header ids option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"hip_header": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable IPv6 Host Identify Protocol header ids option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"fragment_header": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable IPv6 fragment header ids option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"mobility_header": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable IPv6 mobility header ids option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_next_header": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable IPv6 no next header ids option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"routing_header": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable IPv6 routing header ids option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"shim6_header": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable IPv6 shim header ids option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"user_defined_header_type": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "User-defined header type range.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.ValueStringsAre(
										securityScreenUserDefinedHeaderTypeValidator,
									),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"destination_header": schema.SingleNestedBlock{
								Description: "Enable IPv6 destination option header ids option.",
								Attributes: map[string]schema.Attribute{
									"home_address_option": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable home address option ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"ilnp_nonce_option": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable Identifier-Locator Network Protocol Nonce option ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"line_identification_option": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable line identification option ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"tunnel_encapsulation_limit_option": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable tunnel encapsulation limit option ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"user_defined_option_type": schema.ListAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Description: "User-defined option type range.",
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
											listvalidator.ValueStringsAre(
												securityScreenUserDefinedOptionTypeValidator,
											),
										},
									},
								},
								PlanModifiers: []planmodifier.Object{
									tfplanmodifier.BlockRemoveNull(),
								},
							},
							"hop_by_hop_header": schema.SingleNestedBlock{
								Description: "Enable IPv6 hop by hop option header ids option.",
								Attributes: map[string]schema.Attribute{
									"calipso_option": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable Common Architecture Label IPv6 Security Option ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"jumbo_payload_option": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable jumbo payload option ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"quick_start_option": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable quick start option ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"router_alert_option": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable router alert option ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"rpl_option": schema.BoolAttribute{
										Optional:    true,
										Description: "nable Routing Protocol for Low-power and Lossy networks option ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"smf_dpd_option": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable Simplified Multicast Forwarding ipv6 Duplicate Packet Detection option ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"user_defined_option_type": schema.ListAttribute{
										ElementType: types.StringType,
										Optional:    true,
										Description: "User-defined option type range.",
										Validators: []validator.List{
											listvalidator.SizeAtLeast(1),
											listvalidator.ValueStringsAre(
												securityScreenUserDefinedOptionTypeValidator,
											),
										},
									},
								},
								PlanModifiers: []planmodifier.Object{
									tfplanmodifier.BlockRemoveNull(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"tunnel": schema.SingleNestedBlock{
						Description: "Configure IP tunnel ids options.",
						Attributes: map[string]schema.Attribute{
							"bad_inner_header": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable IP tunnel bad inner header ids option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"ip_in_udp_teredo": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable IP tunnel IPinUDP Teredo ids option.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"gre": schema.SingleNestedBlock{
								Description: "Configure IP tunnel GRE ids option.",
								Attributes: map[string]schema.Attribute{
									"gre_4in4": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel GRE 4in4 ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"gre_4in6": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel GRE 4in6 ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"gre_6in4": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel GRE 6in4 ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"gre_6in6": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel GRE 6in6 ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
								},
								PlanModifiers: []planmodifier.Object{
									tfplanmodifier.BlockRemoveNull(),
								},
							},
							"ipip": schema.SingleNestedBlock{
								Description: "Configure IP tunnel IPIP ids option.",
								Attributes: map[string]schema.Attribute{
									"dslite": schema.BoolAttribute{
										Optional:    true,
										Description: " Enable IP tunnel IPIP DS-Lite ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"ipip_4in4": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel IPIP 4in4 ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"ipip_4in6": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel IPIP 4in6 ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"ipip_6in4": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel IPIP 6in4 ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"ipip_6in6": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel IPIP 6in6 ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"ipip_6over4": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel IPIP 6over4 ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"ipip_6to4relay": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel IPIP 6to4 Relay ids option.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"isatap": schema.BoolAttribute{
										Optional:    true,
										Description: "Enable IP tunnel IPIP ISATAP ids option.",
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
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"limit_session": schema.SingleNestedBlock{
				Description: "Configure limit sessions.",
				Attributes: map[string]schema.Attribute{
					"destination_ip_based": schema.Int64Attribute{
						Optional:    true,
						Description: "Limit sessions to the same destination IP.",
						Validators: []validator.Int64{
							int64validator.Between(1, 2000000),
						},
					},
					"source_ip_based": schema.Int64Attribute{
						Optional:    true,
						Description: "Limit sessions from the same source IP",
						Validators: []validator.Int64{
							int64validator.Between(1, 2000000),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"tcp": schema.SingleNestedBlock{
				Description: "Configure TCP Layer ids options.",
				Attributes: map[string]schema.Attribute{
					"fin_no_ack": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable Fin bit with no ACK bit ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"land": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable land attack ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_flag": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable TCP packet without flag ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"syn_fin": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable SYN and FIN bits set attack ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"syn_frag": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable SYN fragment ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"winnuke": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable winnuke attack ids option.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"port_scan": schema.SingleNestedBlock{
						Description: "Enable TCP port scan ids option.",
						Attributes: map[string]schema.Attribute{
							"threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Threshold (microseconds in which 10 attack packets are detected).",
								Validators: []validator.Int64{
									int64validator.Between(1000, 1000000),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"sweep": schema.SingleNestedBlock{
						Description: "Enable TCP sweep ids option.",
						Attributes: map[string]schema.Attribute{
							"threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Threshold (microseconds in which 10 TCP packets are detected).",
								Validators: []validator.Int64{
									int64validator.Between(1000, 1000000),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"syn_ack_ack_proxy": schema.SingleNestedBlock{
						Description: "Enable syn-ack-ack proxy ids option.",
						Attributes: map[string]schema.Attribute{
							"threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Threshold (un-authenticated connections).",
								Validators: []validator.Int64{
									int64validator.Between(1, 250000),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"syn_flood": schema.SingleNestedBlock{
						Description: "Enable SYN flood ids option.",
						Attributes: map[string]schema.Attribute{
							"alarm_threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Alarm threshold (requests per second).",
								Validators: []validator.Int64{
									int64validator.Between(1, 500000),
								},
							},
							"attack_threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Attack threshold (proxied requests per second).",
								Validators: []validator.Int64{
									int64validator.Between(1, 500000),
								},
							},
							"destination_threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Destination threshold (SYN pps).",
								Validators: []validator.Int64{
									int64validator.Between(4, 500000),
								},
							},
							"source_threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Source threshold (SYN pps).",
								Validators: []validator.Int64{
									int64validator.Between(4, 500000),
								},
							},
							"timeout": schema.Int64Attribute{
								Optional:    true,
								Description: "SYN flood ager timeout (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(1, 50),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"whitelist": schema.SetNestedBlock{
								Description: "For each name of white-list to declare.",
								NestedObject: schema.NestedBlockObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Required:    true,
											Description: "White-list name.",
											Validators: []validator.String{
												stringvalidator.LengthBetween(1, 32),
												tfvalidator.StringFormat(tfvalidator.DefaultFormat),
											},
										},
										"destination_address": schema.SetAttribute{
											ElementType: types.StringType,
											Optional:    true,
											Description: "Destination address.",
											Validators: []validator.Set{
												setvalidator.SizeAtLeast(1),
												setvalidator.ValueStringsAre(
													tfvalidator.StringCIDRNetwork(),
												),
											},
										},
										"source_address": schema.SetAttribute{
											ElementType: types.StringType,
											Optional:    true,
											Description: "Source address.",
											Validators: []validator.Set{
												setvalidator.SizeAtLeast(1),
												setvalidator.ValueStringsAre(
													tfvalidator.StringCIDRNetwork(),
												),
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
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"udp": schema.SingleNestedBlock{
				Description: "Configure UDP layer ids options.",
				Blocks: map[string]schema.Block{
					"flood": schema.SingleNestedBlock{
						Description: "UDP flood ids option.",
						Attributes: map[string]schema.Attribute{
							"threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Threshold (UDP packets per second).",
								Validators: []validator.Int64{
									int64validator.Between(1, 1000000),
								},
							},
							"whitelist": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "List of UDP flood white list group name.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthBetween(1, 32),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"port_scan": schema.SingleNestedBlock{
						Description: "UDP port scan ids option.",
						Attributes: map[string]schema.Attribute{
							"threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Threshold (microseconds in which 10 attack packets are detected).",
								Validators: []validator.Int64{
									int64validator.Between(1000, 1000000),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"sweep": schema.SingleNestedBlock{
						Description: "UDP sweep ids option.",
						Attributes: map[string]schema.Attribute{
							"threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Threshold (microseconds in which 10 UDP packets are detected).",
								Validators: []validator.Int64{
									int64validator.Between(1000, 1000000),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
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

type securityScreenData struct {
	ID               types.String                     `tfsdk:"id"`
	Name             types.String                     `tfsdk:"name"`
	AlarmWithoutDrop types.Bool                       `tfsdk:"alarm_without_drop"`
	Description      types.String                     `tfsdk:"description"`
	Icmp             *securityScreenBlockIcmp         `tfsdk:"icmp"`
	IP               *securityScreenBlockIP           `tfsdk:"ip"`
	LimitSession     *securityScreenBlockLimitSession `tfsdk:"limit_session"`
	TCP              *securityScreenBlockTCP          `tfsdk:"tcp"`
	UDP              *securityScreenBlockUDP          `tfsdk:"udp"`
}

func (rscData *securityScreenData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData, "ID", "Name")
}

type securityScreenConfig struct {
	ID               types.String                     `tfsdk:"id"`
	Name             types.String                     `tfsdk:"name"`
	AlarmWithoutDrop types.Bool                       `tfsdk:"alarm_without_drop"`
	Description      types.String                     `tfsdk:"description"`
	Icmp             *securityScreenBlockIcmp         `tfsdk:"icmp"`
	IP               *securityScreenBlockIPConfig     `tfsdk:"ip"`
	LimitSession     *securityScreenBlockLimitSession `tfsdk:"limit_session"`
	TCP              *securityScreenBlockTCPConfig    `tfsdk:"tcp"`
	UDP              *securityScreenBlockUDPConfig    `tfsdk:"udp"`
}

func (config *securityScreenConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(config, "ID", "Name")
}

type securityScreenBlockWithThreshold struct {
	Threshold types.Int64 `tfsdk:"threshold"`
}

type securityScreenBlockIcmp struct {
	Fragment        types.Bool                        `tfsdk:"fragment"`
	Icmpv6Malformed types.Bool                        `tfsdk:"icmpv6_malformed"`
	Large           types.Bool                        `tfsdk:"large"`
	PingDeath       types.Bool                        `tfsdk:"ping_death"`
	Flood           *securityScreenBlockWithThreshold `tfsdk:"flood"`
	Sweep           *securityScreenBlockWithThreshold `tfsdk:"sweep"`
}

func (block *securityScreenBlockIcmp) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockIP struct {
	BadOption                types.Bool                                     `tfsdk:"bad_option"`
	BlockFrag                types.Bool                                     `tfsdk:"block_frag"`
	IPv6ExtensionHeaderLimit types.Int64                                    `tfsdk:"ipv6_extension_header_limit"`
	IPv6MalformedHeader      types.Bool                                     `tfsdk:"ipv6_malformed_header"`
	LooseSourceRouteOption   types.Bool                                     `tfsdk:"loose_source_route_option"`
	RecordRouteOption        types.Bool                                     `tfsdk:"record_route_option"`
	SecurityOption           types.Bool                                     `tfsdk:"security_option"`
	SourceRouteOption        types.Bool                                     `tfsdk:"source_route_option"`
	Spoofing                 types.Bool                                     `tfsdk:"spoofing"`
	StreamOption             types.Bool                                     `tfsdk:"stream_option"`
	StrictSourceRouteOption  types.Bool                                     `tfsdk:"strict_source_route_option"`
	TearDrop                 types.Bool                                     `tfsdk:"tear_drop"`
	TimestampOption          types.Bool                                     `tfsdk:"timestamp_option"`
	UnknownProtocol          types.Bool                                     `tfsdk:"unknown_protocol"`
	IPv6ExtensionHeader      *securityScreenBlockIPBlockIPv6ExtensionHeader `tfsdk:"ipv6_extension_header"`
	Tunnel                   *securityScreenBlockIPBlockTunnel              `tfsdk:"tunnel"`
}

func (block *securityScreenBlockIP) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockIPConfig struct {
	BadOption                types.Bool                                           `tfsdk:"bad_option"`
	BlockFrag                types.Bool                                           `tfsdk:"block_frag"`
	IPv6ExtensionHeaderLimit types.Int64                                          `tfsdk:"ipv6_extension_header_limit"`
	IPv6MalformedHeader      types.Bool                                           `tfsdk:"ipv6_malformed_header"`
	LooseSourceRouteOption   types.Bool                                           `tfsdk:"loose_source_route_option"`
	RecordRouteOption        types.Bool                                           `tfsdk:"record_route_option"`
	SecurityOption           types.Bool                                           `tfsdk:"security_option"`
	SourceRouteOption        types.Bool                                           `tfsdk:"source_route_option"`
	Spoofing                 types.Bool                                           `tfsdk:"spoofing"`
	StreamOption             types.Bool                                           `tfsdk:"stream_option"`
	StrictSourceRouteOption  types.Bool                                           `tfsdk:"strict_source_route_option"`
	TearDrop                 types.Bool                                           `tfsdk:"tear_drop"`
	TimestampOption          types.Bool                                           `tfsdk:"timestamp_option"`
	UnknownProtocol          types.Bool                                           `tfsdk:"unknown_protocol"`
	IPv6ExtensionHeader      *securityScreenBlockIPBlockIPv6ExtensionHeaderConfig `tfsdk:"ipv6_extension_header"`
	Tunnel                   *securityScreenBlockIPBlockTunnel                    `tfsdk:"tunnel"`
}

func (block *securityScreenBlockIPConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:lll
type securityScreenBlockIPBlockIPv6ExtensionHeader struct {
	AhHeader              types.Bool                                                           `tfsdk:"ah_header"`
	EspHeader             types.Bool                                                           `tfsdk:"esp_header"`
	HipHeader             types.Bool                                                           `tfsdk:"hip_header"`
	FragmentHeader        types.Bool                                                           `tfsdk:"fragment_header"`
	MobilityHeader        types.Bool                                                           `tfsdk:"mobility_header"`
	NoNextHeader          types.Bool                                                           `tfsdk:"no_next_header"`
	RoutingHeader         types.Bool                                                           `tfsdk:"routing_header"`
	Shim6Header           types.Bool                                                           `tfsdk:"shim6_header"`
	UserDefinedHeaderType []types.String                                                       `tfsdk:"user_defined_header_type"`
	DestinationHeader     *securityScreenBlockIPBlockIPv6ExtensionHeaderBlockDestinationHeader `tfsdk:"destination_header"`
	HopByHopHeader        *securityScreenBlockIPBlockIPv6ExtensionHeaderBlockHopByHopHeader    `tfsdk:"hop_by_hop_header"`
}

func (block *securityScreenBlockIPBlockIPv6ExtensionHeader) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:lll
type securityScreenBlockIPBlockIPv6ExtensionHeaderConfig struct {
	AhHeader              types.Bool                                                                 `tfsdk:"ah_header"`
	EspHeader             types.Bool                                                                 `tfsdk:"esp_header"`
	HipHeader             types.Bool                                                                 `tfsdk:"hip_header"`
	FragmentHeader        types.Bool                                                                 `tfsdk:"fragment_header"`
	MobilityHeader        types.Bool                                                                 `tfsdk:"mobility_header"`
	NoNextHeader          types.Bool                                                                 `tfsdk:"no_next_header"`
	RoutingHeader         types.Bool                                                                 `tfsdk:"routing_header"`
	Shim6Header           types.Bool                                                                 `tfsdk:"shim6_header"`
	UserDefinedHeaderType types.List                                                                 `tfsdk:"user_defined_header_type"`
	DestinationHeader     *securityScreenBlockIPBlockIPv6ExtensionHeaderBlockDestinationHeaderConfig `tfsdk:"destination_header"`
	HopByHopHeader        *securityScreenBlockIPBlockIPv6ExtensionHeaderBlockHopByHopHeaderConfig    `tfsdk:"hop_by_hop_header"`
}

func (block *securityScreenBlockIPBlockIPv6ExtensionHeaderConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockIPBlockIPv6ExtensionHeaderBlockDestinationHeader struct {
	HomeAddressOption              types.Bool     `tfsdk:"home_address_option"`
	IlnpNonceOption                types.Bool     `tfsdk:"ilnp_nonce_option"`
	LineIdentificationOption       types.Bool     `tfsdk:"line_identification_option"`
	TunnelEncapsulationLimitOption types.Bool     `tfsdk:"tunnel_encapsulation_limit_option"`
	UserDefinedOptionType          []types.String `tfsdk:"user_defined_option_type"`
}

type securityScreenBlockIPBlockIPv6ExtensionHeaderBlockDestinationHeaderConfig struct {
	HomeAddressOption              types.Bool `tfsdk:"home_address_option"`
	IlnpNonceOption                types.Bool `tfsdk:"ilnp_nonce_option"`
	LineIdentificationOption       types.Bool `tfsdk:"line_identification_option"`
	TunnelEncapsulationLimitOption types.Bool `tfsdk:"tunnel_encapsulation_limit_option"`
	UserDefinedOptionType          types.List `tfsdk:"user_defined_option_type"`
}

type securityScreenBlockIPBlockIPv6ExtensionHeaderBlockHopByHopHeader struct {
	CalipsoOption         types.Bool     `tfsdk:"calipso_option"`
	JumboPayloadOption    types.Bool     `tfsdk:"jumbo_payload_option"`
	QuickStartOption      types.Bool     `tfsdk:"quick_start_option"`
	RouterAlertOption     types.Bool     `tfsdk:"router_alert_option"`
	RplOption             types.Bool     `tfsdk:"rpl_option"`
	SmfDpdOption          types.Bool     `tfsdk:"smf_dpd_option"`
	UserDefinedOptionType []types.String `tfsdk:"user_defined_option_type"`
}

type securityScreenBlockIPBlockIPv6ExtensionHeaderBlockHopByHopHeaderConfig struct {
	CalipsoOption         types.Bool `tfsdk:"calipso_option"`
	JumboPayloadOption    types.Bool `tfsdk:"jumbo_payload_option"`
	QuickStartOption      types.Bool `tfsdk:"quick_start_option"`
	RouterAlertOption     types.Bool `tfsdk:"router_alert_option"`
	RplOption             types.Bool `tfsdk:"rpl_option"`
	SmfDpdOption          types.Bool `tfsdk:"smf_dpd_option"`
	UserDefinedOptionType types.List `tfsdk:"user_defined_option_type"`
}

type securityScreenBlockIPBlockTunnel struct {
	BadInnerHeader types.Bool                                 `tfsdk:"bad_inner_header"`
	IPInUDPTeredo  types.Bool                                 `tfsdk:"ip_in_udp_teredo"`
	Gre            *securityScreenBlockIPBlockTunnelBlockGre  `tfsdk:"gre"`
	Ipip           *securityScreenBlockIPBlockTunnelBlockIpip `tfsdk:"ipip"`
}

func (block *securityScreenBlockIPBlockTunnel) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockIPBlockTunnelBlockGre struct {
	Gre4in4 types.Bool `tfsdk:"gre_4in4"`
	Gre4in6 types.Bool `tfsdk:"gre_4in6"`
	Gre6in4 types.Bool `tfsdk:"gre_6in4"`
	Gre6in6 types.Bool `tfsdk:"gre_6in6"`
}

func (block *securityScreenBlockIPBlockTunnelBlockGre) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockIPBlockTunnelBlockIpip struct {
	Dslite        types.Bool `tfsdk:"dslite"`
	Ipip4in4      types.Bool `tfsdk:"ipip_4in4"`
	Ipip4in6      types.Bool `tfsdk:"ipip_4in6"`
	Ipip6in4      types.Bool `tfsdk:"ipip_6in4"`
	Ipip6in6      types.Bool `tfsdk:"ipip_6in6"`
	Ipip6over4    types.Bool `tfsdk:"ipip_6over4"`
	Ipip6to4relay types.Bool `tfsdk:"ipip_6to4relay"`
	Isatap        types.Bool `tfsdk:"isatap"`
}

func (block *securityScreenBlockIPBlockTunnelBlockIpip) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockLimitSession struct {
	DestinationIPBased types.Int64 `tfsdk:"destination_ip_based"`
	SourceIPBased      types.Int64 `tfsdk:"source_ip_based"`
}

func (block *securityScreenBlockLimitSession) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockTCP struct {
	FinNoAck       types.Bool                           `tfsdk:"fin_no_ack"`
	Land           types.Bool                           `tfsdk:"land"`
	NoFlag         types.Bool                           `tfsdk:"no_flag"`
	SynFin         types.Bool                           `tfsdk:"syn_fin"`
	SynFrag        types.Bool                           `tfsdk:"syn_frag"`
	Winnuke        types.Bool                           `tfsdk:"winnuke"`
	PortScan       *securityScreenBlockWithThreshold    `tfsdk:"port_scan"`
	Sweep          *securityScreenBlockWithThreshold    `tfsdk:"sweep"`
	SynAckAckProxy *securityScreenBlockWithThreshold    `tfsdk:"syn_ack_ack_proxy"`
	SynFlood       *securityScreenBlockTCPBlockSynFlood `tfsdk:"syn_flood"`
}

func (block *securityScreenBlockTCP) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockTCPConfig struct {
	FinNoAck       types.Bool                                 `tfsdk:"fin_no_ack"`
	Land           types.Bool                                 `tfsdk:"land"`
	NoFlag         types.Bool                                 `tfsdk:"no_flag"`
	SynFin         types.Bool                                 `tfsdk:"syn_fin"`
	SynFrag        types.Bool                                 `tfsdk:"syn_frag"`
	Winnuke        types.Bool                                 `tfsdk:"winnuke"`
	PortScan       *securityScreenBlockWithThreshold          `tfsdk:"port_scan"`
	Sweep          *securityScreenBlockWithThreshold          `tfsdk:"sweep"`
	SynAckAckProxy *securityScreenBlockWithThreshold          `tfsdk:"syn_ack_ack_proxy"`
	SynFlood       *securityScreenBlockTCPBlockSynFloodConfig `tfsdk:"syn_flood"`
}

func (block *securityScreenBlockTCPConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockTCPBlockSynFlood struct {
	AlarmThreshold       types.Int64                                         `tfsdk:"alarm_threshold"`
	AttackThreshold      types.Int64                                         `tfsdk:"attack_threshold"`
	DestinationThreshold types.Int64                                         `tfsdk:"destination_threshold"`
	SourceThreshold      types.Int64                                         `tfsdk:"source_threshold"`
	Timeout              types.Int64                                         `tfsdk:"timeout"`
	Whitelist            []securityScreenBlockTCPBlockSynFloodBlockWhitelist `tfsdk:"whitelist"`
}

type securityScreenBlockTCPBlockSynFloodConfig struct {
	AlarmThreshold       types.Int64 `tfsdk:"alarm_threshold"`
	AttackThreshold      types.Int64 `tfsdk:"attack_threshold"`
	DestinationThreshold types.Int64 `tfsdk:"destination_threshold"`
	SourceThreshold      types.Int64 `tfsdk:"source_threshold"`
	Timeout              types.Int64 `tfsdk:"timeout"`
	Whitelist            types.Set   `tfsdk:"whitelist"`
}

type securityScreenBlockTCPBlockSynFloodBlockWhitelist struct {
	Name               types.String   `tfsdk:"name"`
	DestinationAddress []types.String `tfsdk:"destination_address"`
	SourceAddress      []types.String `tfsdk:"source_address"`
}

func (block *securityScreenBlockTCPBlockSynFloodBlockWhitelist) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block, "Name")
}

type securityScreenBlockTCPBlockSynFloodBlockWhitelistConfig struct {
	Name               types.String `tfsdk:"name"`
	DestinationAddress types.Set    `tfsdk:"destination_address"`
	SourceAddress      types.Set    `tfsdk:"source_address"`
}

func (block *securityScreenBlockTCPBlockSynFloodBlockWhitelistConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block, "Name")
}

type securityScreenBlockUDP struct {
	Flood    *securityScreenBlockUDPBlockFlood `tfsdk:"flood"`
	PortScan *securityScreenBlockWithThreshold `tfsdk:"port_scan"`
	Sweep    *securityScreenBlockWithThreshold `tfsdk:"sweep"`
}

func (block *securityScreenBlockUDP) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockUDPConfig struct {
	Flood    *securityScreenBlockUDPBlockFloodConfig `tfsdk:"flood"`
	PortScan *securityScreenBlockWithThreshold       `tfsdk:"port_scan"`
	Sweep    *securityScreenBlockWithThreshold       `tfsdk:"sweep"`
}

func (block *securityScreenBlockUDPConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityScreenBlockUDPBlockFlood struct {
	Threshold types.Int64    `tfsdk:"threshold"`
	Whitelist []types.String `tfsdk:"whitelist"`
}

type securityScreenBlockUDPBlockFloodConfig struct {
	Threshold types.Int64 `tfsdk:"threshold"`
	Whitelist types.Set   `tfsdk:"whitelist"`
}

func (rsc *securityScreen) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityScreenConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one of arguments need to be set (in addition to `name`)",
		)
	}

	if config.Icmp != nil {
		if config.Icmp.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("icmp").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"icmp block is empty",
			)
		}
	}
	if config.IP != nil {
		if config.IP.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ip").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"ip block is empty",
			)
		}
		if config.IP.IPv6ExtensionHeader != nil {
			if config.IP.IPv6ExtensionHeader.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ip").AtName("ipv6_extension_header").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"ipv6_extension_header block is empty"+
						" in ip block",
				)
			}
		}
		if config.IP.Tunnel != nil {
			if config.IP.Tunnel.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ip").AtName("tunnel").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"tunnel block is empty"+
						" in ip block",
				)
			}
			if config.IP.Tunnel.Gre != nil {
				if config.IP.Tunnel.Gre.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("ip").AtName("tunnel").AtName("gre").AtName("*"),
						tfdiag.MissingConfigErrSummary,
						"gre block is empty"+
							" in tunnel block in ip block",
					)
				}
			}
			if config.IP.Tunnel.Ipip != nil {
				if config.IP.Tunnel.Ipip.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("ip").AtName("tunnel").AtName("ipip").AtName("*"),
						tfdiag.MissingConfigErrSummary,
						"ipip block is empty"+
							" in tunnel block in ip block",
					)
				}
			}
		}
	}
	if config.LimitSession != nil {
		if config.LimitSession.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("limit_session").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"limit_session block is empty",
			)
		}
	}
	if config.TCP != nil {
		if config.TCP.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("tcp").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"tcp block is empty",
			)
		}
		if config.TCP.SynFlood != nil {
			if !config.TCP.SynFlood.Whitelist.IsNull() &&
				!config.TCP.SynFlood.Whitelist.IsUnknown() {
				var configWhitelist []securityScreenBlockTCPBlockSynFloodBlockWhitelistConfig
				asDiags := config.TCP.SynFlood.Whitelist.ElementsAs(ctx, &configWhitelist, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				whitelistName := make(map[string]struct{})
				for _, block := range configWhitelist {
					if block.isEmpty() {
						resp.Diagnostics.AddAttributeError(
							path.Root("tcp").AtName("syn_flood").AtName("whitelist"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("whitelist block %q is empty"+
								" in syn_flood block in tcp block", block.Name.ValueString()),
						)
					}
					if block.Name.IsUnknown() {
						continue
					}
					name := block.Name.ValueString()
					if _, ok := whitelistName[name]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("tcp").AtName("syn_flood").AtName("whitelist"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple whitelist blocks with the same name %q"+
								" in syn_flood block in tcp block", name),
						)
					}
					whitelistName[name] = struct{}{}
				}
			}
		}
	}
	if config.UDP != nil {
		if config.UDP.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("udp").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"udp block is empty",
			)
		}
	}
}

func (rsc *securityScreen) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityScreenData
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
			if !junSess.CheckCompatibilitySecurity() {
				resp.Diagnostics.AddError(
					tfdiag.CompatibilityErrSummary,
					rsc.junosName()+junSess.SystemInformation.NotCompatibleMsg(),
				)

				return false
			}
			screenExists, err := checkSecurityScreenExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if screenExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			screenExists, err := checkSecurityScreenExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !screenExists {
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

func (rsc *securityScreen) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityScreenData
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

func (rsc *securityScreen) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityScreenData
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

func (rsc *securityScreen) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityScreenData
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

func (rsc *securityScreen) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityScreenData

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

func checkSecurityScreenExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security screen ids-option \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityScreenData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityScreenData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityScreenData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security screen ids-option \"" + rscData.Name.ValueString() + "\" "

	if rscData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name`)")
	}

	if rscData.AlarmWithoutDrop.ValueBool() {
		configSet = append(configSet, setPrefix+"alarm-without-drop")
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}

	if rscData.Icmp != nil {
		if rscData.Icmp.isEmpty() {
			return path.Root("icmp").AtName("*"),
				errors.New("icmp block is empty")
		}

		configSet = append(configSet, rscData.Icmp.configSet(setPrefix)...)
	}
	if rscData.IP != nil {
		if rscData.IP.isEmpty() {
			return path.Root("ip").AtName("*"),
				errors.New("ip block is empty")
		}

		blockSet, pathErr, err := rscData.IP.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.LimitSession != nil {
		if rscData.LimitSession.isEmpty() {
			return path.Root("limit_session").AtName("*"),
				errors.New("limit_session block is empty")
		}

		if !rscData.LimitSession.DestinationIPBased.IsNull() {
			configSet = append(configSet, setPrefix+"limit-session destination-ip-based "+
				utils.ConvI64toa(rscData.LimitSession.DestinationIPBased.ValueInt64()))
		}
		if !rscData.LimitSession.SourceIPBased.IsNull() {
			configSet = append(configSet, setPrefix+"limit-session source-ip-based "+
				utils.ConvI64toa(rscData.LimitSession.SourceIPBased.ValueInt64()))
		}
	}
	if rscData.TCP != nil {
		if rscData.TCP.isEmpty() {
			return path.Root("tcp").AtName("*"),
				errors.New("tcp block is empty")
		}

		blockSet, pathErr, err := rscData.TCP.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.UDP != nil {
		if rscData.UDP.isEmpty() {
			return path.Root("udp").AtName("*"),
				errors.New("udp block is empty")
		}

		configSet = append(configSet, rscData.UDP.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityScreenBlockIcmp) configSet(setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "icmp "

	if block.Fragment.ValueBool() {
		configSet = append(configSet, setPrefix+"fragment")
	}
	if block.Icmpv6Malformed.ValueBool() {
		configSet = append(configSet, setPrefix+"icmpv6-malformed")
	}
	if block.Large.ValueBool() {
		configSet = append(configSet, setPrefix+"large")
	}
	if block.PingDeath.ValueBool() {
		configSet = append(configSet, setPrefix+"ping-death")
	}

	if block.Flood != nil {
		configSet = append(configSet, setPrefix+"flood")

		if !block.Flood.Threshold.IsNull() {
			configSet = append(configSet, setPrefix+"flood threshold "+
				utils.ConvI64toa(block.Flood.Threshold.ValueInt64()))
		}
	}
	if block.Sweep != nil {
		configSet = append(configSet, setPrefix+"ip-sweep")

		if !block.Sweep.Threshold.IsNull() {
			configSet = append(configSet, setPrefix+"ip-sweep threshold "+
				utils.ConvI64toa(block.Sweep.Threshold.ValueInt64()))
		}
	}

	return configSet
}

func (block *securityScreenBlockIP) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "ip "

	if block.BadOption.ValueBool() {
		configSet = append(configSet, setPrefix+"bad-option")
	}
	if block.BlockFrag.ValueBool() {
		configSet = append(configSet, setPrefix+"block-frag")
	}
	if !block.IPv6ExtensionHeaderLimit.IsNull() {
		configSet = append(configSet, setPrefix+"ipv6-extension-header-limit "+
			utils.ConvI64toa(block.IPv6ExtensionHeaderLimit.ValueInt64()))
	}
	if block.IPv6MalformedHeader.ValueBool() {
		configSet = append(configSet, setPrefix+"ipv6-malformed-header")
	}
	if block.LooseSourceRouteOption.ValueBool() {
		configSet = append(configSet, setPrefix+"loose-source-route-option")
	}
	if block.RecordRouteOption.ValueBool() {
		configSet = append(configSet, setPrefix+"record-route-option")
	}
	if block.SecurityOption.ValueBool() {
		configSet = append(configSet, setPrefix+"security-option")
	}
	if block.SourceRouteOption.ValueBool() {
		configSet = append(configSet, setPrefix+"source-route-option")
	}
	if block.Spoofing.ValueBool() {
		configSet = append(configSet, setPrefix+"spoofing")
	}
	if block.StreamOption.ValueBool() {
		configSet = append(configSet, setPrefix+"stream-option")
	}
	if block.StrictSourceRouteOption.ValueBool() {
		configSet = append(configSet, setPrefix+"strict-source-route-option")
	}
	if block.TearDrop.ValueBool() {
		configSet = append(configSet, setPrefix+"tear-drop")
	}
	if block.TimestampOption.ValueBool() {
		configSet = append(configSet, setPrefix+"timestamp-option")
	}
	if block.UnknownProtocol.ValueBool() {
		configSet = append(configSet, setPrefix+"unknown-protocol")
	}

	if block.IPv6ExtensionHeader != nil {
		if block.IPv6ExtensionHeader.isEmpty() {
			return configSet,
				path.Root("ip").AtName("ipv6_extension_header").AtName("*"),
				errors.New("ipv6_extension_header block is empty" +
					" in ip block")
		}

		configSet = append(configSet, block.IPv6ExtensionHeader.configSet(setPrefix)...)
	}
	if block.Tunnel != nil {
		if block.Tunnel.isEmpty() {
			return configSet,
				path.Root("ip").AtName("tunnel").AtName("*"),
				errors.New("tunnel block is empty" +
					" in ip block")
		}

		blockSet, pathErr, err := block.Tunnel.configSet(setPrefix)
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return configSet, path.Empty(), nil
}

func (block *securityScreenBlockIPBlockIPv6ExtensionHeader) configSet(setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "ipv6-extension-header "

	if block.AhHeader.ValueBool() {
		configSet = append(configSet, setPrefix+"AH-header")
	}
	if block.EspHeader.ValueBool() {
		configSet = append(configSet, setPrefix+"ESP-header")
	}
	if block.HipHeader.ValueBool() {
		configSet = append(configSet, setPrefix+"HIP-header")
	}
	if block.FragmentHeader.ValueBool() {
		configSet = append(configSet, setPrefix+"fragment-header")
	}
	if block.MobilityHeader.ValueBool() {
		configSet = append(configSet, setPrefix+"mobility-header")
	}
	if block.NoNextHeader.ValueBool() {
		configSet = append(configSet, setPrefix+"no-next-header")
	}
	if block.RoutingHeader.ValueBool() {
		configSet = append(configSet, setPrefix+"routing-header")
	}
	if block.Shim6Header.ValueBool() {
		configSet = append(configSet, setPrefix+"shim6-header")
	}
	for _, v := range block.UserDefinedHeaderType {
		configSet = append(configSet, setPrefix+"user-defined-header-type "+v.ValueString())
	}

	if block.DestinationHeader != nil {
		configSet = append(configSet, block.DestinationHeader.configSet(setPrefix)...)
	}
	if block.HopByHopHeader != nil {
		configSet = append(configSet, block.HopByHopHeader.configSet(setPrefix)...)
	}

	return configSet
}

func (block *securityScreenBlockIPBlockIPv6ExtensionHeaderBlockDestinationHeader) configSet(setPrefix string) []string {
	setPrefix += "destination-header "

	configSet := []string{
		setPrefix,
	}

	if block.HomeAddressOption.ValueBool() {
		configSet = append(configSet, setPrefix+"home-address-option")
	}
	if block.IlnpNonceOption.ValueBool() {
		configSet = append(configSet, setPrefix+"ILNP-nonce-option")
	}
	if block.LineIdentificationOption.ValueBool() {
		configSet = append(configSet, setPrefix+"line-identification-option")
	}
	if block.TunnelEncapsulationLimitOption.ValueBool() {
		configSet = append(configSet, setPrefix+"tunnel-encapsulation-limit-option")
	}
	for _, v := range block.UserDefinedOptionType {
		configSet = append(configSet, setPrefix+"user-defined-option-type "+v.ValueString())
	}

	return configSet
}

func (block *securityScreenBlockIPBlockIPv6ExtensionHeaderBlockHopByHopHeader) configSet(setPrefix string) []string {
	setPrefix += "hop-by-hop-header "

	configSet := []string{
		setPrefix,
	}

	if block.CalipsoOption.ValueBool() {
		configSet = append(configSet, setPrefix+"CALIPSO-option")
	}
	if block.JumboPayloadOption.ValueBool() {
		configSet = append(configSet, setPrefix+"jumbo-payload-option")
	}
	if block.QuickStartOption.ValueBool() {
		configSet = append(configSet, setPrefix+"quick-start-option")
	}
	if block.RouterAlertOption.ValueBool() {
		configSet = append(configSet, setPrefix+"router-alert-option")
	}
	if block.RplOption.ValueBool() {
		configSet = append(configSet, setPrefix+"RPL-option")
	}
	if block.SmfDpdOption.ValueBool() {
		configSet = append(configSet, setPrefix+"SMF-DPD-option")
	}
	for _, v := range block.UserDefinedOptionType {
		configSet = append(configSet, setPrefix+"user-defined-option-type "+v.ValueString())
	}

	return configSet
}

func (block *securityScreenBlockIPBlockTunnel) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "tunnel "

	if block.BadInnerHeader.ValueBool() {
		configSet = append(configSet, setPrefix+"bad-inner-header")
	}
	if block.IPInUDPTeredo.ValueBool() {
		configSet = append(configSet, setPrefix+"ip-in-udp teredo")
	}

	if block.Gre != nil {
		if block.Gre.isEmpty() {
			return configSet,
				path.Root("ip").AtName("tunnel").AtName("gre").AtName("*"),
				errors.New("gre block is empty" +
					" in tunnel block in ip block")
		}

		configSet = append(configSet, block.Gre.configSet(setPrefix)...)
	}
	if block.Ipip != nil {
		if block.Ipip.isEmpty() {
			return configSet,
				path.Root("ip").AtName("tunnel").AtName("ipip").AtName("*"),
				errors.New("ipip block is empty" +
					" in tunnel block in ip block")
		}

		configSet = append(configSet, block.Ipip.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (block *securityScreenBlockIPBlockTunnelBlockGre) configSet(setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "gre "

	if block.Gre4in4.ValueBool() {
		configSet = append(configSet, setPrefix+"gre-4in4")
	}
	if block.Gre4in6.ValueBool() {
		configSet = append(configSet, setPrefix+"gre-4in6")
	}
	if block.Gre6in4.ValueBool() {
		configSet = append(configSet, setPrefix+"gre-6in4")
	}
	if block.Gre6in6.ValueBool() {
		configSet = append(configSet, setPrefix+"gre-6in6")
	}

	return configSet
}

func (block *securityScreenBlockIPBlockTunnelBlockIpip) configSet(setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "ipip "

	if block.Dslite.ValueBool() {
		configSet = append(configSet, setPrefix+"dslite")
	}
	if block.Ipip4in4.ValueBool() {
		configSet = append(configSet, setPrefix+"ipip-4in4")
	}
	if block.Ipip4in6.ValueBool() {
		configSet = append(configSet, setPrefix+"ipip-4in6")
	}
	if block.Ipip6in4.ValueBool() {
		configSet = append(configSet, setPrefix+"ipip-6in4")
	}
	if block.Ipip6in6.ValueBool() {
		configSet = append(configSet, setPrefix+"ipip-6in6")
	}
	if block.Ipip6over4.ValueBool() {
		configSet = append(configSet, setPrefix+"ipip-6over4")
	}
	if block.Ipip6to4relay.ValueBool() {
		configSet = append(configSet, setPrefix+"ipip-6to4relay")
	}
	if block.Isatap.ValueBool() {
		configSet = append(configSet, setPrefix+"isatap")
	}

	return configSet
}

func (block *securityScreenBlockTCP) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "tcp "

	if block.FinNoAck.ValueBool() {
		configSet = append(configSet, setPrefix+"fin-no-ack")
	}
	if block.Land.ValueBool() {
		configSet = append(configSet, setPrefix+"land")
	}
	if block.NoFlag.ValueBool() {
		configSet = append(configSet, setPrefix+"tcp-no-flag")
	}
	if block.SynFin.ValueBool() {
		configSet = append(configSet, setPrefix+"syn-fin")
	}
	if block.SynFrag.ValueBool() {
		configSet = append(configSet, setPrefix+"syn-frag")
	}
	if block.Winnuke.ValueBool() {
		configSet = append(configSet, setPrefix+"winnuke")
	}

	if block.PortScan != nil {
		configSet = append(configSet, setPrefix+"port-scan")

		if !block.PortScan.Threshold.IsNull() {
			configSet = append(configSet, setPrefix+"port-scan threshold "+
				utils.ConvI64toa(block.PortScan.Threshold.ValueInt64()))
		}
	}
	if block.Sweep != nil {
		configSet = append(configSet, setPrefix+"tcp-sweep")

		if !block.Sweep.Threshold.IsNull() {
			configSet = append(configSet, setPrefix+"tcp-sweep threshold "+
				utils.ConvI64toa(block.Sweep.Threshold.ValueInt64()))
		}
	}
	if block.SynAckAckProxy != nil {
		configSet = append(configSet, setPrefix+"syn-ack-ack-proxy")

		if !block.SynAckAckProxy.Threshold.IsNull() {
			configSet = append(configSet, setPrefix+"syn-ack-ack-proxy threshold "+
				utils.ConvI64toa(block.SynAckAckProxy.Threshold.ValueInt64()))
		}
	}
	if block.SynFlood != nil {
		blockSet, pathErr, err := block.SynFlood.configSet(setPrefix)
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return configSet, path.Empty(), nil
}

func (block *securityScreenBlockTCPBlockSynFlood) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "syn-flood "

	configSet := []string{
		setPrefix,
	}

	if !block.AlarmThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"alarm-threshold "+
			utils.ConvI64toa(block.AlarmThreshold.ValueInt64()))
	}
	if !block.AttackThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"attack-threshold "+
			utils.ConvI64toa(block.AttackThreshold.ValueInt64()))
	}
	if !block.DestinationThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"destination-threshold "+
			utils.ConvI64toa(block.DestinationThreshold.ValueInt64()))
	}
	if !block.SourceThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"source-threshold "+
			utils.ConvI64toa(block.SourceThreshold.ValueInt64()))
	}
	if !block.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(block.Timeout.ValueInt64()))
	}

	whitelistName := make(map[string]struct{})
	for _, subBlock := range block.Whitelist {
		name := subBlock.Name.ValueString()
		if _, ok := whitelistName[name]; ok {
			return configSet,
				path.Root("tcp").AtName("syn_flood").AtName("whitelist"),
				fmt.Errorf("multiple whitelist blocks with the same name %q"+
					" in syn_flood block in tcp block", name)
		}
		whitelistName[name] = struct{}{}

		if subBlock.isEmpty() {
			return configSet,
				path.Root("tcp").AtName("syn_flood").AtName("whitelist"),
				fmt.Errorf("whitelist block %q is empty"+
					" in syn_flood block in tcp block", name)
		}

		for _, v := range subBlock.DestinationAddress {
			configSet = append(configSet, setPrefix+"white-list "+name+
				" destination-address "+v.ValueString())
		}
		for _, v := range subBlock.SourceAddress {
			configSet = append(configSet, setPrefix+"white-list "+name+
				" source-address "+v.ValueString())
		}
	}

	return configSet, path.Empty(), nil
}

func (block *securityScreenBlockUDP) configSet(setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "udp "

	if block.Flood != nil {
		configSet = append(configSet, setPrefix+"flood")

		if !block.Flood.Threshold.IsNull() {
			configSet = append(configSet, setPrefix+"flood threshold "+
				utils.ConvI64toa(block.Flood.Threshold.ValueInt64()))
		}
		for _, v := range block.Flood.Whitelist {
			configSet = append(configSet, setPrefix+"flood white-list "+v.ValueString())
		}
	}
	if block.PortScan != nil {
		configSet = append(configSet, setPrefix+"port-scan")

		if !block.PortScan.Threshold.IsNull() {
			configSet = append(configSet, setPrefix+"port-scan threshold "+
				utils.ConvI64toa(block.PortScan.Threshold.ValueInt64()))
		}
	}
	if block.Sweep != nil {
		configSet = append(configSet, setPrefix+"udp-sweep")

		if !block.Sweep.Threshold.IsNull() {
			configSet = append(configSet, setPrefix+"udp-sweep threshold "+
				utils.ConvI64toa(block.Sweep.Threshold.ValueInt64()))
		}
	}

	return configSet
}

func (rscData *securityScreenData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security screen ids-option \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case itemTrim == "alarm-without-drop":
				rscData.AlarmWithoutDrop = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "icmp "):
				if rscData.Icmp == nil {
					rscData.Icmp = &securityScreenBlockIcmp{}
				}

				if err := rscData.Icmp.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "ip "):
				if rscData.IP == nil {
					rscData.IP = &securityScreenBlockIP{}
				}

				if err := rscData.IP.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "limit-session "):
				if rscData.LimitSession == nil {
					rscData.LimitSession = &securityScreenBlockLimitSession{}
				}

				switch {
				case balt.CutPrefixInString(&itemTrim, "destination-ip-based "):
					rscData.LimitSession.DestinationIPBased, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "source-ip-based "):
					rscData.LimitSession.SourceIPBased, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "tcp "):
				if rscData.TCP == nil {
					rscData.TCP = &securityScreenBlockTCP{}
				}

				if err := rscData.TCP.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "udp "):
				if rscData.UDP == nil {
					rscData.UDP = &securityScreenBlockUDP{}
				}

				if err := rscData.UDP.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *securityScreenBlockIcmp) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "fragment":
		block.Fragment = types.BoolValue(true)
	case itemTrim == "icmpv6-malformed":
		block.Icmpv6Malformed = types.BoolValue(true)
	case itemTrim == "large":
		block.Large = types.BoolValue(true)
	case itemTrim == "ping-death":
		block.PingDeath = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "flood"):
		if block.Flood == nil {
			block.Flood = &securityScreenBlockWithThreshold{}
		}

		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			block.Flood.Threshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "ip-sweep"):
		if block.Sweep == nil {
			block.Sweep = &securityScreenBlockWithThreshold{}
		}

		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			block.Sweep.Threshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (block *securityScreenBlockIP) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "bad-option":
		block.BadOption = types.BoolValue(true)
	case itemTrim == "block-frag":
		block.BlockFrag = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "ipv6-extension-header-limit "):
		block.IPv6ExtensionHeaderLimit, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "ipv6-malformed-header":
		block.IPv6MalformedHeader = types.BoolValue(true)
	case itemTrim == "loose-source-route-option":
		block.LooseSourceRouteOption = types.BoolValue(true)
	case itemTrim == "record-route-option":
		block.RecordRouteOption = types.BoolValue(true)
	case itemTrim == "security-option":
		block.SecurityOption = types.BoolValue(true)
	case itemTrim == "source-route-option":
		block.SourceRouteOption = types.BoolValue(true)
	case itemTrim == "spoofing":
		block.Spoofing = types.BoolValue(true)
	case itemTrim == "stream-option":
		block.StreamOption = types.BoolValue(true)
	case itemTrim == "strict-source-route-option":
		block.StrictSourceRouteOption = types.BoolValue(true)
	case itemTrim == "tear-drop":
		block.TearDrop = types.BoolValue(true)
	case itemTrim == "timestamp-option":
		block.TimestampOption = types.BoolValue(true)
	case itemTrim == "unknown-protocol":
		block.UnknownProtocol = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "ipv6-extension-header "):
		if block.IPv6ExtensionHeader == nil {
			block.IPv6ExtensionHeader = &securityScreenBlockIPBlockIPv6ExtensionHeader{}
		}

		block.IPv6ExtensionHeader.read(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "tunnel "):
		if block.Tunnel == nil {
			block.Tunnel = &securityScreenBlockIPBlockTunnel{}
		}

		block.Tunnel.read(itemTrim)
	}

	return nil
}

func (block *securityScreenBlockIPBlockIPv6ExtensionHeader) read(itemTrim string) {
	switch {
	case itemTrim == "AH-header":
		block.AhHeader = types.BoolValue(true)
	case itemTrim == "ESP-header":
		block.EspHeader = types.BoolValue(true)
	case itemTrim == "HIP-header":
		block.HipHeader = types.BoolValue(true)
	case itemTrim == "fragment-header":
		block.FragmentHeader = types.BoolValue(true)
	case itemTrim == "mobility-header":
		block.MobilityHeader = types.BoolValue(true)
	case itemTrim == "no-next-header":
		block.NoNextHeader = types.BoolValue(true)
	case itemTrim == "routing-header":
		block.RoutingHeader = types.BoolValue(true)
	case itemTrim == "shim6-header":
		block.Shim6Header = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "user-defined-header-type "):
		block.UserDefinedHeaderType = append(block.UserDefinedHeaderType, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "destination-header"):
		if block.DestinationHeader == nil {
			block.DestinationHeader = &securityScreenBlockIPBlockIPv6ExtensionHeaderBlockDestinationHeader{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			block.DestinationHeader.read(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "hop-by-hop-header"):
		if block.HopByHopHeader == nil {
			block.HopByHopHeader = &securityScreenBlockIPBlockIPv6ExtensionHeaderBlockHopByHopHeader{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			block.HopByHopHeader.read(itemTrim)
		}
	}
}

func (block *securityScreenBlockIPBlockIPv6ExtensionHeaderBlockDestinationHeader) read(itemTrim string) {
	switch {
	case itemTrim == "home-address-option":
		block.HomeAddressOption = types.BoolValue(true)
	case itemTrim == "ILNP-nonce-option":
		block.IlnpNonceOption = types.BoolValue(true)
	case itemTrim == "line-identification-option":
		block.LineIdentificationOption = types.BoolValue(true)
	case itemTrim == "tunnel-encapsulation-limit-option":
		block.TunnelEncapsulationLimitOption = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "user-defined-option-type "):
		block.UserDefinedOptionType = append(block.UserDefinedOptionType, types.StringValue(itemTrim))
	}
}

func (block *securityScreenBlockIPBlockIPv6ExtensionHeaderBlockHopByHopHeader) read(itemTrim string) {
	switch {
	case itemTrim == "CALIPSO-option":
		block.CalipsoOption = types.BoolValue(true)
	case itemTrim == "jumbo-payload-option":
		block.JumboPayloadOption = types.BoolValue(true)
	case itemTrim == "quick-start-option":
		block.QuickStartOption = types.BoolValue(true)
	case itemTrim == "router-alert-option":
		block.RouterAlertOption = types.BoolValue(true)
	case itemTrim == "RPL-option":
		block.RplOption = types.BoolValue(true)
	case itemTrim == "SMF-DPD-option":
		block.SmfDpdOption = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "user-defined-option-type "):
		block.UserDefinedOptionType = append(block.UserDefinedOptionType, types.StringValue(itemTrim))
	}
}

func (block *securityScreenBlockIPBlockTunnel) read(itemTrim string) {
	switch {
	case itemTrim == "bad-inner-header":
		block.BadInnerHeader = types.BoolValue(true)
	case itemTrim == "ip-in-udp teredo":
		block.IPInUDPTeredo = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "gre "):
		if block.Gre == nil {
			block.Gre = &securityScreenBlockIPBlockTunnelBlockGre{}
		}

		block.Gre.read(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ipip "):
		if block.Ipip == nil {
			block.Ipip = &securityScreenBlockIPBlockTunnelBlockIpip{}
		}

		block.Ipip.read(itemTrim)
	}
}

func (block *securityScreenBlockIPBlockTunnelBlockGre) read(itemTrim string) {
	switch {
	case itemTrim == "gre-4in4":
		block.Gre4in4 = types.BoolValue(true)
	case itemTrim == "gre-4in6":
		block.Gre4in6 = types.BoolValue(true)
	case itemTrim == "gre-6in4":
		block.Gre6in4 = types.BoolValue(true)
	case itemTrim == "gre-6in6":
		block.Gre6in6 = types.BoolValue(true)
	}
}

func (block *securityScreenBlockIPBlockTunnelBlockIpip) read(itemTrim string) {
	switch {
	case itemTrim == "dslite":
		block.Dslite = types.BoolValue(true)
	case itemTrim == "ipip-4in4":
		block.Ipip4in4 = types.BoolValue(true)
	case itemTrim == "ipip-4in6":
		block.Ipip4in6 = types.BoolValue(true)
	case itemTrim == "ipip-6in4":
		block.Ipip6in4 = types.BoolValue(true)
	case itemTrim == "ipip-6in6":
		block.Ipip6in6 = types.BoolValue(true)
	case itemTrim == "ipip-6over4":
		block.Ipip6over4 = types.BoolValue(true)
	case itemTrim == "ipip-6to4relay":
		block.Ipip6to4relay = types.BoolValue(true)
	case itemTrim == "isatap":
		block.Isatap = types.BoolValue(true)
	}
}

func (block *securityScreenBlockTCP) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "fin-no-ack":
		block.FinNoAck = types.BoolValue(true)
	case itemTrim == "land":
		block.Land = types.BoolValue(true)
	case itemTrim == "tcp-no-flag":
		block.NoFlag = types.BoolValue(true)
	case itemTrim == "syn-fin":
		block.SynFin = types.BoolValue(true)
	case itemTrim == "syn-frag":
		block.SynFrag = types.BoolValue(true)
	case itemTrim == "winnuke":
		block.Winnuke = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "port-scan"):
		if block.PortScan == nil {
			block.PortScan = &securityScreenBlockWithThreshold{}
		}

		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			block.PortScan.Threshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "tcp-sweep"):
		if block.Sweep == nil {
			block.Sweep = &securityScreenBlockWithThreshold{}
		}

		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			block.Sweep.Threshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "syn-ack-ack-proxy"):
		if block.SynAckAckProxy == nil {
			block.SynAckAckProxy = &securityScreenBlockWithThreshold{}
		}

		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			block.SynAckAckProxy.Threshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "syn-flood"):
		if block.SynFlood == nil {
			block.SynFlood = &securityScreenBlockTCPBlockSynFlood{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			if err := block.SynFlood.read(itemTrim); err != nil {
				return err
			}
		}
	}

	return nil
}

func (block *securityScreenBlockTCPBlockSynFlood) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "alarm-threshold "):
		block.AlarmThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "attack-threshold "):
		block.AttackThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "destination-threshold "):
		block.DestinationThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "source-threshold "):
		block.SourceThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "timeout "):
		block.Timeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "white-list "):
		name := tfdata.FirstElementOfJunosLine(itemTrim)
		var whitelist securityScreenBlockTCPBlockSynFloodBlockWhitelist
		block.Whitelist, whitelist = tfdata.ExtractBlockWithTFTypesString(
			block.Whitelist, "Name", name,
		)
		whitelist.Name = types.StringValue(name)
		balt.CutPrefixInString(&itemTrim, name+" ")

		switch {
		case balt.CutPrefixInString(&itemTrim, "destination-address "):
			whitelist.DestinationAddress = append(whitelist.DestinationAddress, types.StringValue(itemTrim))
		case balt.CutPrefixInString(&itemTrim, "source-address "):
			whitelist.SourceAddress = append(whitelist.SourceAddress, types.StringValue(itemTrim))
		}
		block.Whitelist = append(block.Whitelist, whitelist)
	}

	return nil
}

func (block *securityScreenBlockUDP) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "flood"):
		if block.Flood == nil {
			block.Flood = &securityScreenBlockUDPBlockFlood{}
		}

		switch {
		case balt.CutPrefixInString(&itemTrim, " threshold "):
			block.Flood.Threshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " white-list "):
			block.Flood.Whitelist = append(block.Flood.Whitelist, types.StringValue(itemTrim))
		}
	case balt.CutPrefixInString(&itemTrim, "port-scan"):
		if block.PortScan == nil {
			block.PortScan = &securityScreenBlockWithThreshold{}
		}

		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			block.PortScan.Threshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "udp-sweep"):
		if block.Sweep == nil {
			block.Sweep = &securityScreenBlockWithThreshold{}
		}

		if balt.CutPrefixInString(&itemTrim, " threshold ") {
			block.Sweep.Threshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (rscData *securityScreenData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security screen ids-option \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
