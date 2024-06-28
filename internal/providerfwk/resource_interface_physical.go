package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                   = &interfacePhysical{}
	_ resource.ResourceWithConfigure      = &interfacePhysical{}
	_ resource.ResourceWithValidateConfig = &interfacePhysical{}
	_ resource.ResourceWithImportState    = &interfacePhysical{}
	_ resource.ResourceWithUpgradeState   = &interfacePhysical{}
)

type interfacePhysical struct {
	client *junos.Client
}

func newInterfacePhysicalResource() resource.Resource {
	return &interfacePhysical{}
}

func (rsc *interfacePhysical) typeName() string {
	return providerName + "_interface_physical"
}

func (rsc *interfacePhysical) junosName() string {
	return "physical interface"
}

func (rsc *interfacePhysical) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *interfacePhysical) Configure(
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

func (rsc *interfacePhysical) Schema(
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
				Description: "Name of physical interface (without dot).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					tfvalidator.StringDotExclusion(),
				},
			},
			"no_disable_on_destroy": schema.BoolAttribute{
				Optional:    true,
				Description: "When destroy this resource, delete all configurations.",
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
				Description: "Disable this interface.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"encapsulation": schema.StringAttribute{
				Optional:    true,
				Description: "Physical link-layer encapsulation.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"flexible_vlan_tagging": schema.BoolAttribute{
				Optional:    true,
				Description: "Support for no tagging, or single and double 802.1q VLAN tagging.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"gratuitous_arp_reply": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable gratuitous ARP reply.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"hold_time_down": schema.Int64Attribute{
				Optional:    true,
				Description: "Link down hold time (milliseconds).",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"hold_time_up": schema.Int64Attribute{
				Optional:    true,
				Description: "Link up hold time (milliseconds).",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"link_mode": schema.StringAttribute{
				Optional:    true,
				Description: "Link operational mode.",
				Validators: []validator.String{
					stringvalidator.OneOf("automatic", "full-duplex", "half-duplex"),
				},
			},
			"mtu": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum transmission unit.",
				Validators: []validator.Int64{
					int64validator.Between(1, 9500),
				},
			},
			"no_gratuitous_arp_reply": schema.BoolAttribute{
				Optional:    true,
				Description: "Don't enable gratuitous ARP reply.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_gratuitous_arp_request": schema.BoolAttribute{
				Optional:    true,
				Description: "Ignore gratuitous ARP request.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"speed": schema.StringAttribute{
				Optional:    true,
				Description: "Link speed.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(\d+(m|g)|2\.5g|auto|auto-10m-100m)$`),
						"must be a valid speed (10m | 100m | 1g ...)"),
				},
			},
			"storm_control": schema.StringAttribute{
				Optional:    true,
				Description: "Storm control profile name to bind.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 127),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"trunk": schema.BoolAttribute{
				Optional:    true,
				Description: "Interface mode is trunk.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"trunk_non_els": schema.BoolAttribute{
				Optional:    true,
				Description: "Port mode is trunk.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"vlan_members": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of vlan for membership for this interface.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 63),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
			"vlan_native": schema.Int64Attribute{
				Optional:    true,
				Description: "Vlan for untagged frames.",
				Validators: []validator.Int64{
					int64validator.Between(1, 4094),
				},
			},
			"vlan_native_non_els": schema.StringAttribute{
				Optional:    true,
				Description: "Vlan for untagged frames (non-ELS).",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 64),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"vlan_tagging": schema.BoolAttribute{
				Optional:    true,
				Description: "Add 802.1q VLAN tagging support.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"esi": schema.SingleNestedBlock{
				Description: "Define ESI Config parameters.",
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "ESI Mode.",
						Validators: []validator.String{
							stringvalidator.OneOf("all-active", "single-active"),
						},
					},
					"auto_derive_lacp": schema.BoolAttribute{
						Optional:    true,
						Description: "Auto-derive ESI value for the interface.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"df_election_type": schema.StringAttribute{
						Optional:    true,
						Description: "DF Election Type.",
						Validators: []validator.String{
							stringvalidator.OneOf("mod", "preference"),
						},
					},
					"identifier": schema.StringAttribute{
						Optional:    true,
						Description: "The ESI value for the interface.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^([\d\w]{2}:){9}[\d\w]{2}$`),
								"must be ten octets integer value with colon separator"),
						},
					},
					"source_bmac": schema.StringAttribute{
						Optional:    true,
						Description: "Unicast Source B-MAC address per ESI for PBB-EVPN.",
						Validators: []validator.String{
							tfvalidator.StringMACAddress().WithMac48ColonHexa(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"ether_opts": schema.SingleNestedBlock{
				Description: "Declare `ether-options` configuration.",
				Attributes:  rsc.schemaEtherOptsAttributes(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"gigether_opts": schema.SingleNestedBlock{
				Description: "Declare `gigether-options` configuration.",
				Attributes:  rsc.schemaEtherOptsAttributes(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"parent_ether_opts": schema.SingleNestedBlock{
				Description: "Declare `aggregated-ether-options` or `redundant-ether-options` configuration" +
					" (it depends on the interface `name`).",
				Attributes: map[string]schema.Attribute{
					"flow_control": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable flow control.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_flow_control": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't enable flow control.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"loopback": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable loopback.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_loopback": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't enable loopback.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"link_speed": schema.StringAttribute{
						Optional:    true,
						Description: "Link speed of individual interface that joins the AE.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"100m", "1g", "2.5g", "5g", "8g",
								"10g", "25g", "40g", "50g", "80g",
								"100g", "400g", "mixed", "oc192",
							),
						},
					},
					"minimum_bandwidth": schema.StringAttribute{
						Optional:    true,
						Description: "Minimum bandwidth configured for aggregated bundle.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^[0-9]+ (k|g|m)?bps$`),
								"must be 'N (k|g|m)?bps' format"),
						},
					},
					"minimum_links": schema.Int64Attribute{
						Optional:    true,
						Description: "Minimum number of aggregated/active links (1..64).",
						Validators: []validator.Int64{
							int64validator.Between(1, 64),
						},
					},
					"redundancy_group": schema.Int64Attribute{
						Optional:    true,
						Description: "Redundancy group of this interface (1..128) for reth interface.",
						Validators: []validator.Int64{
							int64validator.Between(1, 128),
						},
					},
					"source_address_filter": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Source address filters.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.ValueStringsAre(
								tfvalidator.StringMACAddress().WithMac48ColonHexa(),
							),
						},
					},
					"source_filtering": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable source address filtering.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"bfd_liveness_detection": schema.SingleNestedBlock{
						Description: "Declare bfd-liveness-detection in aggregated-ether-options configuration.",
						Attributes: map[string]schema.Attribute{
							"local_address": schema.StringAttribute{
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "BFD local address.",
								Validators: []validator.String{
									tfvalidator.StringIPAddress(),
								},
							},
							"authentication_algorithm": schema.StringAttribute{
								Optional:    true,
								Description: "Authentication algorithm name.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"keyed-md5",
										"keyed-sha-1",
										"meticulous-keyed-md5",
										"meticulous-keyed-sha-1",
										"simple-password",
									),
								},
							},
							"authentication_key_chain": schema.StringAttribute{
								Optional:    true,
								Description: "Authentication Key chain name.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"authentication_loose_check": schema.BoolAttribute{
								Optional:    true,
								Description: "Verify authentication only if authentication is negotiated.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"detection_time_threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "High detection-time triggering a trap (milliseconds).",
								Validators: []validator.Int64{
									int64validator.Between(1, 4294967295),
								},
							},
							"holddown_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Time to hold the session-UP notification to the client (0..255000 milliseconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 255000),
								},
							},
							"minimum_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum transmit and receive interval (1..255000 milliseconds).",
								Validators: []validator.Int64{
									int64validator.Between(1, 255000),
								},
							},
							"minimum_receive_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum receive interval (1..255000 milliseconds).",
								Validators: []validator.Int64{
									int64validator.Between(1, 255000),
								},
							},
							"multiplier": schema.Int64Attribute{
								Optional:    true,
								Description: "Detection time multiplier (1..255).",
								Validators: []validator.Int64{
									int64validator.Between(1, 255),
								},
							},
							"neighbor": schema.StringAttribute{
								Optional:    true,
								Description: "BFD neighbor address.",
								Validators: []validator.String{
									tfvalidator.StringIPAddress(),
								},
							},
							"no_adaptation": schema.BoolAttribute{
								Optional:    true,
								Description: "Disable adaptation.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"transmit_interval_minimum_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum transmit interval (1..255000 milliseconds).",
								Validators: []validator.Int64{
									int64validator.Between(1, 255000),
								},
							},
							"transmit_interval_threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "High transmit interval triggering a trap (milliseconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"version": schema.StringAttribute{
								Optional:    true,
								Description: "BFD protocol version number.",
								Validators: []validator.String{
									stringvalidator.OneOf("0", "1", "automatic"),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"lacp": schema.SingleNestedBlock{
						Description: "Declare lacp configuration.",
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								Required: false, // true when SingleNestedBlock is specified
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf("active", "passive"),
								},
							},
							"admin_key": schema.Int64Attribute{
								Optional:    true,
								Description: "Node's administrative key.",
								Validators: []validator.Int64{
									int64validator.Between(0, 65535),
								},
							},
							"periodic": schema.StringAttribute{
								Optional:    true,
								Description: "Timer interval for periodic transmission of LACP packets.",
								Validators: []validator.String{
									stringvalidator.OneOf("fast", "slow"),
								},
							},
							"sync_reset": schema.StringAttribute{
								Optional:    true,
								Description: "On minimum-link failure notify out of sync to peer.",
								Validators: []validator.String{
									stringvalidator.OneOf("disable", "enable"),
								},
							},
							"system_id": schema.StringAttribute{
								Optional:    true,
								Description: "Node's System ID, encoded as a MAC address.",
								Validators: []validator.String{
									tfvalidator.StringMACAddress().WithMac48ColonHexa(),
								},
							},
							"system_priority": schema.Int64Attribute{
								Optional:    true,
								Description: "Priority of the system (0 ... 65535).",
								Validators: []validator.Int64{
									int64validator.Between(0, 65535),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"mc_ae": schema.SingleNestedBlock{
						Description: "Multi-chassis aggregation (MC-AE) network device configuration.",
						Attributes: map[string]schema.Attribute{
							"chassis_id": schema.Int64Attribute{
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "Chassis id of MC-AE network device.",
								Validators: []validator.Int64{
									int64validator.Between(0, 1),
								},
							},
							"mc_ae_id": schema.Int64Attribute{
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "MC-AE group id.",
								Validators: []validator.Int64{
									int64validator.Between(1, 65535),
								},
							},
							"mode": schema.StringAttribute{
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "Mode of the MC-AE.",
								Validators: []validator.String{
									stringvalidator.OneOf("active-active", "active-standby"),
								},
							},
							"status_control": schema.StringAttribute{
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "Status of the MC-AE chassis.",
								Validators: []validator.String{
									stringvalidator.OneOf("active", "standby"),
								},
							},
							"enhanced_convergence": schema.BoolAttribute{
								Optional:    true,
								Description: "Optimized convergence time for mcae.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"init_delay_time": schema.Int64Attribute{
								Optional:    true,
								Description: "Init delay timer for mcae sm for min traffic loss (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(1, 6000),
								},
							},
							"redundancy_group": schema.Int64Attribute{
								Optional:    true,
								Description: "Redundancy group id.",
								Validators: []validator.Int64{
									int64validator.Between(1, 4294967294),
								},
							},
							"revert_time": schema.Int64Attribute{
								Optional:    true,
								Description: "Wait interval before performing switchover (minute).",
								Validators: []validator.Int64{
									int64validator.Between(1, 10),
								},
							},
							"switchover_mode": schema.StringAttribute{
								Optional:    true,
								Description: "Switchover mode.",
								Validators: []validator.String{
									stringvalidator.OneOf("revertive", "non-revertive"),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"events_iccp_peer_down": schema.SingleNestedBlock{
								Description: "Define behavior in the event of ICCP peer down.",
								Attributes: map[string]schema.Attribute{
									"force_icl_down": schema.BoolAttribute{
										Optional:    true,
										Description: "Bring down ICL logical interface.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"prefer_status_control_active": schema.BoolAttribute{
										Optional:    true,
										Description: "Keep this node up.",
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
		},
	}
}

func (rsc *interfacePhysical) schemaEtherOptsAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"ae_8023ad": schema.StringAttribute{
			Optional:    true,
			Description: "Name of an aggregated Ethernet interface to join.",
			Validators: []validator.String{
				stringvalidator.RegexMatches(regexp.MustCompile(
					`^ae\d+$`),
					"must be an ae interface"),
			},
		},
		"auto_negotiation": schema.BoolAttribute{
			Optional:    true,
			Description: "Enable auto-negotiation.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"no_auto_negotiation": schema.BoolAttribute{
			Optional:    true,
			Description: "Don't enable auto-negotiation.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"flow_control": schema.BoolAttribute{
			Optional:    true,
			Description: "Enable flow control.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"no_flow_control": schema.BoolAttribute{
			Optional:    true,
			Description: "Don't enable flow control.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"loopback": schema.BoolAttribute{
			Optional:    true,
			Description: "Enable loopback.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"no_loopback": schema.BoolAttribute{
			Optional:    true,
			Description: "Don't enable loopback.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"redundant_parent": schema.StringAttribute{
			Optional:    true,
			Description: "Name of a redundant ethernet interface to join.",
			Validators: []validator.String{
				stringvalidator.RegexMatches(regexp.MustCompile(
					`^reth\d+$`),
					"must be a reth interface"),
			},
		},
	}
}

type interfacePhysicalData struct {
	ID                     types.String                           `tfsdk:"id"`
	Name                   types.String                           `tfsdk:"name"`
	NoDisableOnDestroy     types.Bool                             `tfsdk:"no_disable_on_destroy"`
	Description            types.String                           `tfsdk:"description"`
	Disable                types.Bool                             `tfsdk:"disable"`
	Encapsulation          types.String                           `tfsdk:"encapsulation"`
	FlexibleVlanTagging    types.Bool                             `tfsdk:"flexible_vlan_tagging"`
	GratuitousArpReply     types.Bool                             `tfsdk:"gratuitous_arp_reply"`
	HoldTimeDown           types.Int64                            `tfsdk:"hold_time_down"`
	HoldTimeUp             types.Int64                            `tfsdk:"hold_time_up"`
	LinkMode               types.String                           `tfsdk:"link_mode"`
	Mtu                    types.Int64                            `tfsdk:"mtu"`
	NoGratuitousArpReply   types.Bool                             `tfsdk:"no_gratuitous_arp_reply"`
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

type interfacePhysicalConfig struct {
	ID                     types.String                                 `tfsdk:"id"`
	Name                   types.String                                 `tfsdk:"name"`
	NoDisableOnDestroy     types.Bool                                   `tfsdk:"no_disable_on_destroy"`
	Description            types.String                                 `tfsdk:"description"`
	Disable                types.Bool                                   `tfsdk:"disable"`
	Encapsulation          types.String                                 `tfsdk:"encapsulation"`
	FlexibleVlanTagging    types.Bool                                   `tfsdk:"flexible_vlan_tagging"`
	GratuitousArpReply     types.Bool                                   `tfsdk:"gratuitous_arp_reply"`
	HoldTimeDown           types.Int64                                  `tfsdk:"hold_time_down"`
	HoldTimeUp             types.Int64                                  `tfsdk:"hold_time_up"`
	LinkMode               types.String                                 `tfsdk:"link_mode"`
	Mtu                    types.Int64                                  `tfsdk:"mtu"`
	NoGratuitousArpReply   types.Bool                                   `tfsdk:"no_gratuitous_arp_reply"`
	NoGratuitousArpRequest types.Bool                                   `tfsdk:"no_gratuitous_arp_request"`
	Speed                  types.String                                 `tfsdk:"speed"`
	StormControl           types.String                                 `tfsdk:"storm_control"`
	Trunk                  types.Bool                                   `tfsdk:"trunk"`
	TrunkNonELS            types.Bool                                   `tfsdk:"trunk_non_els"`
	VlanMembers            types.List                                   `tfsdk:"vlan_members"`
	VlanNative             types.Int64                                  `tfsdk:"vlan_native"`
	VlanNativeNonELS       types.String                                 `tfsdk:"vlan_native_non_els"`
	VlanTagging            types.Bool                                   `tfsdk:"vlan_tagging"`
	ESI                    *interfacePhysicalBlockESI                   `tfsdk:"esi"`
	EtherOpts              *interfacePhysicalBlockEtherOpts             `tfsdk:"ether_opts"`
	GigetherOpts           *interfacePhysicalBlockEtherOpts             `tfsdk:"gigether_opts"`
	ParentEtherOpts        *interfacePhysicalBlockParentEtherOptsConfig `tfsdk:"parent_ether_opts"`
}

type interfacePhysicalBlockESI struct {
	Mode           types.String `tfsdk:"mode"`
	AutoDeriveLACP types.Bool   `tfsdk:"auto_derive_lacp"`
	DFElectionType types.String `tfsdk:"df_election_type"`
	Identifier     types.String `tfsdk:"identifier"`
	SourceBMAC     types.String `tfsdk:"source_bmac"`
}

type interfacePhysicalBlockEtherOpts struct {
	Ae8023ad          types.String `tfsdk:"ae_8023ad"`
	AutoNegotiation   types.Bool   `tfsdk:"auto_negotiation"`
	NoAutoNegotiation types.Bool   `tfsdk:"no_auto_negotiation"`
	FlowControl       types.Bool   `tfsdk:"flow_control"`
	NoFlowControl     types.Bool   `tfsdk:"no_flow_control"`
	Loopback          types.Bool   `tfsdk:"loopback"`
	NoLoopback        types.Bool   `tfsdk:"no_loopback"`
	RedundantParent   types.String `tfsdk:"redundant_parent"`
}

func (block *interfacePhysicalBlockEtherOpts) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *interfacePhysicalBlockEtherOpts) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type interfacePhysicalBlockParentEtherOpts struct {
	FlowControl          types.Bool                                                      `tfsdk:"flow_control"`
	NoFlowControl        types.Bool                                                      `tfsdk:"no_flow_control"`
	Loopback             types.Bool                                                      `tfsdk:"loopback"`
	NoLoopback           types.Bool                                                      `tfsdk:"no_loopback"`
	LinkSpeed            types.String                                                    `tfsdk:"link_speed"`
	MinimumBandwidth     types.String                                                    `tfsdk:"minimum_bandwidth"`
	MinimumLinks         types.Int64                                                     `tfsdk:"minimum_links"`
	RedundancyGroup      types.Int64                                                     `tfsdk:"redundancy_group"`
	SourceAddressFilter  []types.String                                                  `tfsdk:"source_address_filter"`
	SourceFiltering      types.Bool                                                      `tfsdk:"source_filtering"`
	BFDLivenessDetection *interfacePhysicalBlockParentEtherOptsBlockBFDLivenessDetection `tfsdk:"bfd_liveness_detection"`
	LACP                 *interfacePhysicalBlockParentEtherOptsBlockLACP                 `tfsdk:"lacp"`
	MCAE                 *interfacePhysicalBlockParentEtherOptsBlockMCAE                 `tfsdk:"mc_ae"`
}

type interfacePhysicalBlockParentEtherOptsConfig struct {
	FlowControl          types.Bool                                                      `tfsdk:"flow_control"`
	NoFlowControl        types.Bool                                                      `tfsdk:"no_flow_control"`
	Loopback             types.Bool                                                      `tfsdk:"loopback"`
	NoLoopback           types.Bool                                                      `tfsdk:"no_loopback"`
	LinkSpeed            types.String                                                    `tfsdk:"link_speed"`
	MinimumBandwidth     types.String                                                    `tfsdk:"minimum_bandwidth"`
	MinimumLinks         types.Int64                                                     `tfsdk:"minimum_links"`
	RedundancyGroup      types.Int64                                                     `tfsdk:"redundancy_group"`
	SourceAddressFilter  types.List                                                      `tfsdk:"source_address_filter"`
	SourceFiltering      types.Bool                                                      `tfsdk:"source_filtering"`
	BFDLivenessDetection *interfacePhysicalBlockParentEtherOptsBlockBFDLivenessDetection `tfsdk:"bfd_liveness_detection"`
	LACP                 *interfacePhysicalBlockParentEtherOptsBlockLACP                 `tfsdk:"lacp"`
	MCAE                 *interfacePhysicalBlockParentEtherOptsBlockMCAE                 `tfsdk:"mc_ae"`
}

func (block *interfacePhysicalBlockParentEtherOptsConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *interfacePhysicalBlockParentEtherOptsConfig) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type interfacePhysicalBlockParentEtherOptsBlockBFDLivenessDetection struct {
	LocalAddress                    types.String `tfsdk:"local_address"`
	AuthenticationAlgorithm         types.String `tfsdk:"authentication_algorithm"`
	AuthenticationKeyChain          types.String `tfsdk:"authentication_key_chain"`
	AuthenticationLooseCheck        types.Bool   `tfsdk:"authentication_loose_check"`
	DetectionTimeThreshold          types.Int64  `tfsdk:"detection_time_threshold"`
	HolddownInterval                types.Int64  `tfsdk:"holddown_interval"`
	MinimumInterval                 types.Int64  `tfsdk:"minimum_interval"`
	MinimumReceiveInterval          types.Int64  `tfsdk:"minimum_receive_interval"`
	Multiplier                      types.Int64  `tfsdk:"multiplier"`
	Neighbor                        types.String `tfsdk:"neighbor"`
	NoAdaptation                    types.Bool   `tfsdk:"no_adaptation"`
	TransmitIntervalMinimumInterval types.Int64  `tfsdk:"transmit_interval_minimum_interval"`
	TransmitIntervalThreshold       types.Int64  `tfsdk:"transmit_interval_threshold"`
	Version                         types.String `tfsdk:"version"`
}

type interfacePhysicalBlockParentEtherOptsBlockLACP struct {
	Mode           types.String `tfsdk:"mode"`
	AdminKey       types.Int64  `tfsdk:"admin_key"`
	Periodic       types.String `tfsdk:"periodic"`
	SyncReset      types.String `tfsdk:"sync_reset"`
	SystemID       types.String `tfsdk:"system_id"`
	SystemPriority types.Int64  `tfsdk:"system_priority"`
}

//nolint:lll
type interfacePhysicalBlockParentEtherOptsBlockMCAE struct {
	ChassisID           types.Int64                                                            `tfsdk:"chassis_id"`
	MCAEID              types.Int64                                                            `tfsdk:"mc_ae_id"`
	Mode                types.String                                                           `tfsdk:"mode"`
	StatusControl       types.String                                                           `tfsdk:"status_control"`
	EnhancedConvergence types.Bool                                                             `tfsdk:"enhanced_convergence"`
	InitDelayTime       types.Int64                                                            `tfsdk:"init_delay_time"`
	RedundancyGroup     types.Int64                                                            `tfsdk:"redundancy_group"`
	RevertTime          types.Int64                                                            `tfsdk:"revert_time"`
	SwitchoverMode      types.String                                                           `tfsdk:"switchover_mode"`
	EventsIccpPeerDown  *interfacePhysicalBlockParentEtherOptsBlockMCAEBlockEventsIccpPeerDown `tfsdk:"events_iccp_peer_down"`
}

func (block *interfacePhysicalBlockParentEtherOptsBlockMCAE) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type interfacePhysicalBlockParentEtherOptsBlockMCAEBlockEventsIccpPeerDown struct {
	ForceIclDown              types.Bool `tfsdk:"force_icl_down"`
	PreferStatusControlActive types.Bool `tfsdk:"prefer_status_control_active"`
}

//nolint:gocyclo
func (rsc *interfacePhysical) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config interfacePhysicalConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ESI != nil {
		if config.ESI.Mode.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("esi").AtName("mode"),
				tfdiag.MissingConfigErrSummary,
				"mode must be specified in esi block",
			)
		}
		if !config.ESI.AutoDeriveLACP.IsNull() && !config.ESI.AutoDeriveLACP.IsUnknown() &&
			!config.ESI.Identifier.IsNull() && !config.ESI.Identifier.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("esi").AtName("auto_derive_lacp"),
				tfdiag.ConflictConfigErrSummary,
				"only one of auto_derive_lacp or identifier can be specified in esi block",
			)
		}
	}
	if config.EtherOpts != nil {
		if config.EtherOpts.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ether_opts").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"ether_opts block is empty",
			)
		} else if config.EtherOpts.hasKnownValue() &&
			((config.GigetherOpts != nil && config.GigetherOpts.hasKnownValue()) ||
				(config.ParentEtherOpts != nil && config.ParentEtherOpts.hasKnownValue())) {
			resp.Diagnostics.AddAttributeError(
				path.Root("ether_opts").AtName("*"),
				tfdiag.ConflictConfigErrSummary,
				"only one of ether_opts, gigether_opts or parent_ether_opts block can be specified",
			)
		}
		if !config.EtherOpts.Ae8023ad.IsNull() && !config.EtherOpts.Ae8023ad.IsUnknown() &&
			!config.EtherOpts.RedundantParent.IsNull() && !config.EtherOpts.RedundantParent.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ether_opts").AtName("ae_8023ad"),
				tfdiag.ConflictConfigErrSummary,
				"ae_8023ad and redundant_parent cannot be configured together in ether_opts block",
			)
		}
		if !config.EtherOpts.AutoNegotiation.IsNull() && !config.EtherOpts.AutoNegotiation.IsUnknown() &&
			!config.EtherOpts.NoAutoNegotiation.IsNull() && !config.EtherOpts.NoAutoNegotiation.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ether_opts").AtName("auto_negotiation"),
				tfdiag.ConflictConfigErrSummary,
				"auto_negotiation and no_auto_negotiation cannot be configured together in ether_opts block",
			)
		}
		if !config.EtherOpts.FlowControl.IsNull() && !config.EtherOpts.FlowControl.IsUnknown() &&
			!config.EtherOpts.NoFlowControl.IsNull() && !config.EtherOpts.NoFlowControl.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ether_opts").AtName("flow_control"),
				tfdiag.ConflictConfigErrSummary,
				"flow_control and no_flow_control cannot be configured together in ether_opts block",
			)
		}
		if !config.EtherOpts.Loopback.IsNull() && !config.EtherOpts.Loopback.IsUnknown() &&
			!config.EtherOpts.NoLoopback.IsNull() && !config.EtherOpts.NoLoopback.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ether_opts").AtName("loopback"),
				tfdiag.ConflictConfigErrSummary,
				"loopback and no_loopback cannot be configured together in ether_opts block",
			)
		}
	}
	if config.GigetherOpts != nil {
		if config.GigetherOpts.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("gigether_opts").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"gigether_opts block is empty",
			)
		} else if config.GigetherOpts.hasKnownValue() &&
			((config.EtherOpts != nil && config.EtherOpts.hasKnownValue()) ||
				(config.ParentEtherOpts != nil && config.ParentEtherOpts.hasKnownValue())) {
			resp.Diagnostics.AddAttributeError(
				path.Root("gigether_opts").AtName("*"),
				tfdiag.ConflictConfigErrSummary,
				"only one of ether_opts, gigether_opts or parent_ether_opts block can be specified",
			)
		}
		if !config.GigetherOpts.Ae8023ad.IsNull() && !config.GigetherOpts.Ae8023ad.IsUnknown() &&
			!config.GigetherOpts.RedundantParent.IsNull() && !config.GigetherOpts.RedundantParent.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("gigether_opts").AtName("ae_8023ad"),
				tfdiag.ConflictConfigErrSummary,
				"ae_8023ad and redundant_parent cannot be configured together in gigether_opts block",
			)
		}
		if !config.GigetherOpts.AutoNegotiation.IsNull() && !config.GigetherOpts.AutoNegotiation.IsUnknown() &&
			!config.GigetherOpts.NoAutoNegotiation.IsNull() && !config.GigetherOpts.NoAutoNegotiation.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("gigether_opts").AtName("auto_negotiation"),
				tfdiag.ConflictConfigErrSummary,
				"auto_negotiation and no_auto_negotiation cannot be configured together in gigether_opts block",
			)
		}
		if !config.GigetherOpts.FlowControl.IsNull() && !config.GigetherOpts.FlowControl.IsUnknown() &&
			!config.GigetherOpts.NoFlowControl.IsNull() && !config.GigetherOpts.NoFlowControl.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("gigether_opts").AtName("flow_control"),
				tfdiag.ConflictConfigErrSummary,
				"flow_control and no_flow_control cannot be configured together in gigether_opts block",
			)
		}
		if !config.GigetherOpts.Loopback.IsNull() && !config.GigetherOpts.Loopback.IsUnknown() &&
			!config.GigetherOpts.NoLoopback.IsNull() && !config.GigetherOpts.NoLoopback.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("gigether_opts").AtName("loopback"),
				tfdiag.ConflictConfigErrSummary,
				"loopback and no_loopback cannot be configured together in gigether_opts block",
			)
		}
	}
	if config.ParentEtherOpts != nil {
		if config.ParentEtherOpts.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("parent_ether_opts").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"parent_ether_opts block is empty",
			)
		} else if config.ParentEtherOpts.hasKnownValue() {
			if !config.Name.IsUnknown() {
				if v := config.Name.ValueString(); !strings.HasPrefix(v, "ae") && !strings.HasPrefix(v, "reth") {
					resp.Diagnostics.AddAttributeError(
						path.Root("parent_ether_opts").AtName("*"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("parent_ether_opts not compatible with this interface %q"+
							" (need to be ae* or reth* interface)", v),
					)
				}
			}
			if (config.EtherOpts != nil && config.EtherOpts.hasKnownValue()) ||
				(config.GigetherOpts != nil && config.GigetherOpts.hasKnownValue()) {
				resp.Diagnostics.AddAttributeError(
					path.Root("parent_ether_opts").AtName("*"),
					tfdiag.ConflictConfigErrSummary,
					"only one of ether_opts, gigether_opts or parent_ether_opts block can be specified",
				)
			}
		}
		if !config.ParentEtherOpts.FlowControl.IsNull() && !config.ParentEtherOpts.FlowControl.IsUnknown() &&
			!config.ParentEtherOpts.NoFlowControl.IsNull() && !config.ParentEtherOpts.NoFlowControl.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("parent_ether_opts").AtName("flow_control"),
				tfdiag.ConflictConfigErrSummary,
				"flow_control and no_flow_control cannot be configured together in parent_ether_opts block",
			)
		}
		if !config.ParentEtherOpts.Loopback.IsNull() && !config.ParentEtherOpts.Loopback.IsUnknown() &&
			!config.ParentEtherOpts.NoLoopback.IsNull() && !config.ParentEtherOpts.NoLoopback.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("parent_ether_opts").AtName("loopback"),
				tfdiag.ConflictConfigErrSummary,
				"loopback and no_loopback cannot be configured together in parent_ether_opts block",
			)
		}
		if !config.ParentEtherOpts.MinimumBandwidth.IsNull() && !config.ParentEtherOpts.MinimumBandwidth.IsUnknown() &&
			!config.ParentEtherOpts.MinimumLinks.IsNull() && !config.ParentEtherOpts.MinimumLinks.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("parent_ether_opts").AtName("minimum_bandwidth"),
				tfdiag.ConflictConfigErrSummary,
				"minimum_bandwidth and minimum_links cannot be configured together in parent_ether_opts block",
			)
		}
		if config.ParentEtherOpts.BFDLivenessDetection != nil {
			if config.ParentEtherOpts.BFDLivenessDetection.LocalAddress.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("parent_ether_opts").AtName("bfd_liveness_detection").AtName("local_address"),
					tfdiag.MissingConfigErrSummary,
					"local_address must be specified in bfd_liveness_detection block in parent_ether_opts block",
				)
			}
		}
		if config.ParentEtherOpts.LACP != nil {
			if config.ParentEtherOpts.LACP.Mode.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("parent_ether_opts").AtName("lacp").AtName("mode"),
					tfdiag.MissingConfigErrSummary,
					"mode must be specified in lacp block in parent_ether_opts block",
				)
			}
		}
		if config.ParentEtherOpts.MCAE != nil {
			if config.ParentEtherOpts.MCAE.hasKnownValue() && !config.Name.IsUnknown() {
				if v := config.Name.ValueString(); !strings.HasPrefix(v, "ae") {
					resp.Diagnostics.AddAttributeError(
						path.Root("parent_ether_opts").AtName("mc_ae").AtName("*"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("mc_ae in parent_ether_opts block not compatible with this interface %q"+
							" (need to be ae* interface)", v),
					)
				}
			}
			if config.ParentEtherOpts.MCAE.ChassisID.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("parent_ether_opts").AtName("mc_ae").AtName("chassis_id"),
					tfdiag.MissingConfigErrSummary,
					"chassis_id must be specified in mc_ae block in parent_ether_opts block",
				)
			}
			if config.ParentEtherOpts.MCAE.MCAEID.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("parent_ether_opts").AtName("mc_ae_id").AtName("chassis_id"),
					tfdiag.MissingConfigErrSummary,
					"mc_ae_id must be specified in mc_ae block in parent_ether_opts block",
				)
			}
			if config.ParentEtherOpts.MCAE.Mode.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("parent_ether_opts").AtName("mode").AtName("chassis_id"),
					tfdiag.MissingConfigErrSummary,
					"mode must be specified in mc_ae block in parent_ether_opts block",
				)
			}
			if config.ParentEtherOpts.MCAE.StatusControl.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("parent_ether_opts").AtName("status_control").AtName("chassis_id"),
					tfdiag.MissingConfigErrSummary,
					"status_control must be specified in mc_ae block in parent_ether_opts block",
				)
			}
		}
		if !config.ParentEtherOpts.RedundancyGroup.IsNull() && !config.ParentEtherOpts.RedundancyGroup.IsUnknown() {
			if !config.Name.IsUnknown() {
				if v := config.Name.ValueString(); !strings.HasPrefix(v, "reth") {
					resp.Diagnostics.AddAttributeError(
						path.Root("parent_ether_opts").AtName("redundancy_group"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("redundancy_group in parent_ether_opts block not compatible with this interface %q"+
							" (need to be reth* interface)", v),
					)
				}
			}
		}
	}

	if config.Disable.ValueBool() && config.Description.ValueString() == "NC" {
		resp.Diagnostics.AddAttributeError(
			path.Root("disable"),
			tfdiag.ConflictConfigErrSummary,
			"disable=true and description=NC is not allowed "+
				"because the provider might consider the resource deleted",
		)
	}

	if !config.GratuitousArpReply.IsNull() && !config.GratuitousArpReply.IsUnknown() &&
		!config.NoGratuitousArpReply.IsNull() && !config.NoGratuitousArpReply.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("gratuitous_arp_reply"),
			tfdiag.ConflictConfigErrSummary,
			"gratuitous_arp_reply and no_gratuitous_arp_reply cannot be configured together",
		)
	}

	if !config.HoldTimeDown.IsNull() && config.HoldTimeUp.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hold_time_down"),
			tfdiag.MissingConfigErrSummary,
			"hold_time_down and hold_time_up must be specified together",
		)
	}
	if config.HoldTimeDown.IsNull() && !config.HoldTimeUp.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hold_time_up"),
			tfdiag.MissingConfigErrSummary,
			"hold_time_down and hold_time_up must be specified together",
		)
	}
}

func (rsc *interfacePhysical) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan interfacePhysicalData
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
	if strings.Contains(plan.Name.ValueString(), ".") {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Bad Name",
			"could not create "+rsc.junosName()+" with a dot in the name",
		)

		return
	}

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := delInterfaceNC(
			ctx,
			plan.Name.ValueString(),
			rsc.client.GroupInterfaceDelete(),
			junSess,
		); err != nil {
			resp.Diagnostics.AddError("Pre Config Set Error", err.Error())

			return
		}
		if errPath, err := plan.set(ctx, "", junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
			} else {
				resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
			}

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigClearUnlockWarnSummary, junSess.ConfigClear())...)
	}()

	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(
		ctx,
		plan.Name.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

		return
	}
	if !ncInt && !emptyInt {
		resp.Diagnostics.AddError(
			tfdiag.DuplicateConfigErrSummary,
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
			resp.Diagnostics.AddError("Pre Config Set Error", err.Error())

			return
		}
	}

	if errPath, err := plan.set(ctx, "", junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf(ctx, "create resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	ncInt, emptyInt, err = checkInterfacePhysicalNCEmpty(
		ctx,
		plan.Name.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

		return
	}
	if ncInt {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			fmt.Sprintf(rsc.junosName()+" %q always disable (NC) after commit "+
				"=> check your config", plan.Name.ValueString()),
		)

		return
	}
	if emptyInt {
		intExists, err := junSess.CheckInterfaceExists(plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

			return
		}
		if !intExists {
			resp.Diagnostics.AddError(
				tfdiag.NotFoundErrSummary,
				fmt.Sprintf(rsc.junosName()+" %q not exists and config can't found after commit"+
					"=> check your config", plan.Name.ValueString()),
			)

			return
		}
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *interfacePhysical) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data interfacePhysicalData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	defer junos.MutexUnlock()

	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(
		ctx,
		state.Name.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if ncInt {
		resp.State.RemoveResource(ctx)

		return
	}
	if emptyInt {
		intExists, err := junSess.CheckInterfaceExists(state.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !intExists {
			resp.State.RemoveResource(ctx)

			return
		}
	}

	if err := data.read(ctx, state.Name.ValueString(), junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}

	data.NoDisableOnDestroy = state.NoDisableOnDestroy
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *interfacePhysical) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state interfacePhysicalData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var oldAE string
	if state.EtherOpts != nil {
		if v := state.EtherOpts.Ae8023ad.ValueString(); v != "" {
			oldAE = v
		}
	}
	if state.GigetherOpts != nil {
		if v := state.GigetherOpts.Ae8023ad.ValueString(); v != "" {
			oldAE = v
		}
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.delOpts(ctx, junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

			return
		}
		if errPath, err := plan.set(ctx, oldAE, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
			} else {
				resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigClearUnlockWarnSummary, junSess.ConfigClear())...)
	}()

	if err := state.delOpts(ctx, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

		return
	}
	if err := state.unsetAE(ctx, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

		return
	}

	if errPath, err := plan.set(ctx, oldAE, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf(ctx, "update resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *interfacePhysical) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state interfacePhysicalData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeDeleteAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

			return
		}

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigClearUnlockWarnSummary, junSess.ConfigClear())...)
	}()

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

		return
	}
	warns, err := junSess.CommitConf(ctx, "delete resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	if !state.NoDisableOnDestroy.ValueBool() {
		intExists, err := junSess.CheckInterfaceExists(state.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Pre Disable Config Set Error", err.Error())
		} else if intExists {
			if err := addInterfaceNC(
				ctx,
				state.Name.ValueString(),
				rsc.client.GroupInterfaceDelete(),
				junSess,
			); err != nil {
				resp.Diagnostics.AddError("Disable Config Set Error", err.Error())

				return
			}
			warns, err = junSess.CommitConf(ctx, "disable(NC) resource "+rsc.typeName())
			resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

				return
			}
		}
	}
}

func (rsc *interfacePhysical) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	if strings.Count(req.ID, ".") != 0 {
		resp.Diagnostics.AddError(
			tfdiag.PreCheckErrSummary,
			fmt.Sprintf("name of interface need to doesn't have a dot, got %q", req.ID),
		)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(
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
	if emptyInt {
		intExists, err := junSess.CheckInterfaceExists(req.ID)
		if err != nil {
			resp.Diagnostics.AddError("Interface Read Error", err.Error())

			return
		}
		if !intExists {
			resp.Diagnostics.AddError(
				tfdiag.NotFoundErrSummary,
				defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "name"),
			)

			return
		}
	}

	var data interfacePhysicalData
	if err := data.read(ctx, req.ID, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rscData *interfacePhysicalData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func checkInterfacePhysicalNCEmpty(
	_ context.Context, name, groupInterfaceDelete string, junSess *junos.Session,
) (
	ncInt, // interface is set with NC config
	emtyInt bool, // interface is not set (empty)
	_ error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"interfaces " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return false, false, err
	}
	showConfigLines := make([]string, 0)
	// remove unused lines
	for _, item := range strings.Split(showConfig, "\n") {
		// show parameters root on interface exclude unit parameters (except ethernet-switching)
		if strings.HasPrefix(item, "set unit") && !strings.Contains(item, "ethernet-switching") {
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
		return false, true, nil
	}
	showConfig = strings.Join(showConfigLines, "\n")
	if groupInterfaceDelete != "" {
		if showConfig == "set apply-groups "+groupInterfaceDelete {
			return true, false, nil
		}
	}
	if showConfig == "set description NC\nset disable" ||
		showConfig == "set disable\nset description NC" {
		return true, false, nil
	}
	if showConfig == junos.EmptyW {
		return false, true, nil
	}

	return false, false, nil
}

func checkInterfacePhysicalContainsUnit(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"interfaces " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return false, err
	}
	for _, item := range strings.Split(showConfig, "\n") {
		if strings.Contains(item, junos.XMLStartTagConfigOut) {
			continue
		}
		if strings.Contains(item, junos.XMLEndTagConfigOut) {
			break
		}
		if strings.HasPrefix(item, "set unit") {
			if strings.Contains(item, "ethernet-switching") {
				continue
			}

			return true, nil
		}
	}

	return false, nil
}

func (rscData *interfacePhysicalData) unsetAE(
	ctx context.Context, junSess *junos.Session,
) error {
	var oldAE string
	switch {
	case rscData.EtherOpts != nil:
		if v := rscData.EtherOpts.Ae8023ad.ValueString(); v != "" {
			oldAE = v
		}
	case rscData.GigetherOpts != nil:
		if v := rscData.GigetherOpts.Ae8023ad.ValueString(); v != "" {
			oldAE = v
		}
	}
	if oldAE != "" {
		aggregatedCount, err := findInterfaceAggregatedCountMax(
			ctx,
			"",
			oldAE,
			rscData.Name.ValueString(),
			junSess,
		)
		if err != nil {
			return err
		}
		if aggregatedCount == "0" {
			return junSess.ConfigSet([]string{"delete chassis aggregated-devices ethernet device-count"})
		}

		return junSess.ConfigSet([]string{"set chassis aggregated-devices ethernet device-count " + aggregatedCount})
	}

	return nil
}

func (rscData *interfacePhysicalData) set(
	ctx context.Context, oldAE string, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set interfaces " + rscData.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if rscData.Disable.ValueBool() {
		if rscData.Description.ValueString() == "NC" {
			return path.Root("disable"), errors.New("disable=true and description=NC is not allowed " +
				"because the provider might consider the resource deleted")
		}
		configSet = append(configSet, setPrefix+"disable")
	}
	if rscData.ESI != nil {
		configSet = append(configSet, rscData.ESI.configSet(setPrefix)...)
	}
	if v := rscData.Encapsulation.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"encapsulation "+v)
	}
	if name := rscData.Name.ValueString(); strings.HasPrefix(name, "ae") && junSess.HasNetconf() {
		aggregatedCount, err := findInterfaceAggregatedCountMax(ctx, name, "", name, junSess)
		if err != nil {
			return path.Root("name"), err
		}
		configSet = append(configSet, "set chassis aggregated-devices ethernet device-count "+aggregatedCount)
	} else if rscData.EtherOpts != nil || rscData.GigetherOpts != nil {
		var newAE string
		var pathAE path.Path
		switch {
		case rscData.EtherOpts != nil:
			if v := rscData.EtherOpts.Ae8023ad.ValueString(); v != "" {
				pathAE = path.Root("ether_opts").AtName("ae_8023ad")
				newAE = v
				configSet = append(configSet, setPrefix+"ether-options 802.3ad "+v)
			}
			if rscData.EtherOpts.AutoNegotiation.ValueBool() {
				configSet = append(configSet, setPrefix+"ether-options auto-negotiation")
			}
			if rscData.EtherOpts.NoAutoNegotiation.ValueBool() {
				configSet = append(configSet, setPrefix+"ether-options no-auto-negotiation")
			}
			if rscData.EtherOpts.FlowControl.ValueBool() {
				configSet = append(configSet, setPrefix+"ether-options flow-control")
			}
			if rscData.EtherOpts.NoFlowControl.ValueBool() {
				configSet = append(configSet, setPrefix+"ether-options no-flow-control")
			}
			if rscData.EtherOpts.Loopback.ValueBool() {
				configSet = append(configSet, setPrefix+"ether-options loopback")
			}
			if rscData.EtherOpts.NoLoopback.ValueBool() {
				configSet = append(configSet, setPrefix+"ether-options no-loopback")
			}
			if v := rscData.EtherOpts.RedundantParent.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"ether-options redundant-parent "+v)
			}
			if !strings.HasPrefix(configSet[len(configSet)-1], setPrefix+"ether-options") {
				return path.Root("ether_opts").AtName("*"), errors.New("ether_opts block is empty")
			}
		case rscData.GigetherOpts != nil:
			if v := rscData.GigetherOpts.Ae8023ad.ValueString(); v != "" {
				pathAE = path.Root("gigether_opts").AtName("ae_8023ad")
				newAE = v
				configSet = append(configSet, setPrefix+"gigether-options 802.3ad "+v)
			}
			if rscData.GigetherOpts.AutoNegotiation.ValueBool() {
				configSet = append(configSet, setPrefix+"gigether-options auto-negotiation")
			}
			if rscData.GigetherOpts.NoAutoNegotiation.ValueBool() {
				configSet = append(configSet, setPrefix+"gigether-options no-auto-negotiation")
			}
			if rscData.GigetherOpts.FlowControl.ValueBool() {
				configSet = append(configSet, setPrefix+"gigether-options flow-control")
			}
			if rscData.GigetherOpts.NoFlowControl.ValueBool() {
				configSet = append(configSet, setPrefix+"gigether-options no-flow-control")
			}
			if rscData.GigetherOpts.Loopback.ValueBool() {
				configSet = append(configSet, setPrefix+"gigether-options loopback")
			}
			if rscData.GigetherOpts.NoLoopback.ValueBool() {
				configSet = append(configSet, setPrefix+"gigether-options no-loopback")
			}
			if v := rscData.GigetherOpts.RedundantParent.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"gigether-options redundant-parent "+v)
			}
			if !strings.HasPrefix(configSet[len(configSet)-1], setPrefix+"gigether-options") {
				return path.Root("gigether_opts").AtName("*"), errors.New("gigether_opts block is empty")
			}
		}
		if newAE != "" && junSess.HasNetconf() {
			aggregatedCount, err := findInterfaceAggregatedCountMax(ctx, newAE, oldAE, name, junSess)
			if err != nil {
				return pathAE, err
			}
			configSet = append(configSet, "set chassis aggregated-devices ethernet device-count "+aggregatedCount)
		}
	}
	if rscData.FlexibleVlanTagging.ValueBool() {
		configSet = append(configSet, setPrefix+"flexible-vlan-tagging")
	}
	if rscData.GratuitousArpReply.ValueBool() {
		configSet = append(configSet, setPrefix+"gratuitous-arp-reply")
	}
	if !rscData.HoldTimeDown.IsNull() {
		configSet = append(configSet, setPrefix+"hold-time down "+utils.ConvI64toa(rscData.HoldTimeDown.ValueInt64()))
	}
	if !rscData.HoldTimeUp.IsNull() {
		configSet = append(configSet, setPrefix+"hold-time up "+utils.ConvI64toa(rscData.HoldTimeUp.ValueInt64()))
	}
	if v := rscData.LinkMode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"link-mode "+v)
	}
	if !rscData.Mtu.IsNull() {
		configSet = append(configSet, setPrefix+"mtu "+utils.ConvI64toa(rscData.Mtu.ValueInt64()))
	}
	if rscData.NoGratuitousArpReply.ValueBool() {
		configSet = append(configSet, setPrefix+"no-gratuitous-arp-reply")
	}
	if rscData.NoGratuitousArpRequest.ValueBool() {
		configSet = append(configSet, setPrefix+"no-gratuitous-arp-request")
	}
	if rscData.ParentEtherOpts != nil {
		blockSet, pathErr, err := rscData.ParentEtherOpts.configSet(setPrefix, rscData.Name.ValueString())
		if err != nil {
			return pathErr, err
		}
		if len(blockSet) == 0 {
			return path.Root("parent_ether_opts").AtName("*"), errors.New("parent_ether_opts block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if v := rscData.Speed.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"speed "+v)
	}
	if v := rscData.StormControl.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"unit 0 family ethernet-switching storm-control \""+v+"\"")
	}
	if rscData.Trunk.ValueBool() {
		configSet = append(configSet, setPrefix+"unit 0 family ethernet-switching interface-mode trunk")
	}
	if rscData.TrunkNonELS.ValueBool() {
		configSet = append(configSet, setPrefix+"unit 0 family ethernet-switching port-mode trunk")
	}
	for _, v := range rscData.VlanMembers {
		configSet = append(configSet, setPrefix+
			"unit 0 family ethernet-switching vlan members "+v.ValueString())
	}
	if !rscData.VlanNative.IsNull() {
		configSet = append(configSet, setPrefix+"native-vlan-id "+
			utils.ConvI64toa(rscData.VlanNative.ValueInt64()))
	}
	if v := rscData.VlanNativeNonELS.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"unit 0 family ethernet-switching native-vlan-id "+v)
	}
	if rscData.VlanTagging.ValueBool() {
		configSet = append(configSet, setPrefix+"vlan-tagging")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *interfacePhysicalBlockESI) configSet(setPrefix string) []string {
	configSet := []string{
		setPrefix + "esi " + block.Mode.ValueString(),
	}

	if block.AutoDeriveLACP.ValueBool() {
		configSet = append(configSet, setPrefix+"esi auto-derive lacp")
	}
	if v := block.DFElectionType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"esi df-election-type "+v)
	}
	if v := block.Identifier.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"esi "+v)
	}
	if v := block.SourceBMAC.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"esi source-bmac "+v)
	}

	return configSet
}

func (block *interfacePhysicalBlockParentEtherOpts) configSet(
	setPrefix, interfaceName string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	switch {
	case strings.HasPrefix(interfaceName, "ae"):
		setPrefix += "aggregated-ether-options "
	case strings.HasPrefix(interfaceName, "reth"):
		setPrefix += "redundant-ether-options "
	default:
		return configSet,
			path.Root("parent_ether_opts").AtName("*"),
			fmt.Errorf("parent_ether_opts not compatible with this interface %q"+
				" (need to be ae* or reth* interface)", interfaceName)
	}

	if block.BFDLivenessDetection != nil {
		setPrefixBFDLiveDetect := setPrefix + "bfd-liveness-detection "
		configSet = append(configSet, setPrefixBFDLiveDetect+"local-address "+
			block.BFDLivenessDetection.LocalAddress.ValueString())

		if v := block.BFDLivenessDetection.AuthenticationAlgorithm.ValueString(); v != "" {
			configSet = append(configSet, setPrefixBFDLiveDetect+"authentication algorithm "+v)
		}
		if v := block.BFDLivenessDetection.AuthenticationKeyChain.ValueString(); v != "" {
			configSet = append(configSet, setPrefixBFDLiveDetect+"authentication key-chain \""+v+"\"")
		}
		if block.BFDLivenessDetection.AuthenticationLooseCheck.ValueBool() {
			configSet = append(configSet, setPrefixBFDLiveDetect+"authentication loose-check")
		}
		if !block.BFDLivenessDetection.DetectionTimeThreshold.IsNull() {
			configSet = append(configSet, setPrefixBFDLiveDetect+"detection-time threshold "+
				utils.ConvI64toa(block.BFDLivenessDetection.DetectionTimeThreshold.ValueInt64()))
		}
		if !block.BFDLivenessDetection.HolddownInterval.IsNull() {
			configSet = append(configSet, setPrefixBFDLiveDetect+"holddown-interval "+
				utils.ConvI64toa(block.BFDLivenessDetection.HolddownInterval.ValueInt64()))
		}
		if !block.BFDLivenessDetection.MinimumInterval.IsNull() {
			configSet = append(configSet, setPrefixBFDLiveDetect+"minimum-interval "+
				utils.ConvI64toa(block.BFDLivenessDetection.MinimumInterval.ValueInt64()))
		}
		if !block.BFDLivenessDetection.MinimumReceiveInterval.IsNull() {
			configSet = append(configSet, setPrefixBFDLiveDetect+"minimum-receive-interval "+
				utils.ConvI64toa(block.BFDLivenessDetection.MinimumReceiveInterval.ValueInt64()))
		}
		if !block.BFDLivenessDetection.Multiplier.IsNull() {
			configSet = append(configSet, setPrefixBFDLiveDetect+"multiplier "+
				utils.ConvI64toa(block.BFDLivenessDetection.Multiplier.ValueInt64()))
		}
		if v := block.BFDLivenessDetection.Neighbor.ValueString(); v != "" {
			configSet = append(configSet, setPrefixBFDLiveDetect+"neighbor "+v)
		}
		if block.BFDLivenessDetection.NoAdaptation.ValueBool() {
			configSet = append(configSet, setPrefixBFDLiveDetect+"no-adaptation")
		}
		if !block.BFDLivenessDetection.TransmitIntervalMinimumInterval.IsNull() {
			configSet = append(configSet, setPrefixBFDLiveDetect+"transmit-interval minimum-interval "+
				utils.ConvI64toa(block.BFDLivenessDetection.TransmitIntervalMinimumInterval.ValueInt64()))
		}
		if !block.BFDLivenessDetection.TransmitIntervalThreshold.IsNull() {
			configSet = append(configSet, setPrefixBFDLiveDetect+"transmit-interval threshold "+
				utils.ConvI64toa(block.BFDLivenessDetection.TransmitIntervalThreshold.ValueInt64()))
		}
		if v := block.BFDLivenessDetection.Version.ValueString(); v != "" {
			configSet = append(configSet, setPrefixBFDLiveDetect+"version "+v)
		}
	}
	if block.FlowControl.ValueBool() {
		configSet = append(configSet, setPrefix+"flow-control")
	}
	if block.NoFlowControl.ValueBool() {
		configSet = append(configSet, setPrefix+"no-flow-control")
	}
	if block.LACP != nil {
		setPrefixLACP := setPrefix + "lacp "
		configSet = append(configSet, setPrefixLACP+block.LACP.Mode.ValueString())
		if !block.LACP.AdminKey.IsNull() {
			configSet = append(configSet, setPrefixLACP+"admin-key "+
				utils.ConvI64toa(block.LACP.AdminKey.ValueInt64()))
		}
		if v := block.LACP.Periodic.ValueString(); v != "" {
			configSet = append(configSet, setPrefixLACP+"periodic "+v)
		}
		if v := block.LACP.SyncReset.ValueString(); v != "" {
			configSet = append(configSet, setPrefixLACP+"sync-reset "+v)
		}
		if v := block.LACP.SystemID.ValueString(); v != "" {
			configSet = append(configSet, setPrefixLACP+"system-id "+v)
		}
		if !block.LACP.SystemPriority.IsNull() {
			configSet = append(configSet, setPrefixLACP+"system-priority "+
				utils.ConvI64toa(block.LACP.SystemPriority.ValueInt64()))
		}
	}
	if block.Loopback.ValueBool() {
		configSet = append(configSet, setPrefix+"loopback")
	}
	if block.NoLoopback.ValueBool() {
		configSet = append(configSet, setPrefix+"no-loopback")
	}
	if v := block.LinkSpeed.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"link-speed "+v)
	}
	if block.MCAE != nil {
		if !strings.HasPrefix(interfaceName, "ae") {
			return configSet,
				path.Root("parent_ether_opts").AtName("mc_ae").AtName("*"),
				fmt.Errorf("mc_ae in parent_ether_opts block not compatible with this interface %q"+
					" (need to be ae* interface)", interfaceName)
		}
		configSet = append(configSet, setPrefix+"mc-ae chassis-id "+
			utils.ConvI64toa(block.MCAE.ChassisID.ValueInt64()))
		configSet = append(configSet, setPrefix+"mc-ae mc-ae-id "+
			utils.ConvI64toa(block.MCAE.MCAEID.ValueInt64()))
		configSet = append(configSet, setPrefix+"mc-ae mode "+block.MCAE.Mode.ValueString())
		configSet = append(configSet, setPrefix+"mc-ae status-control "+block.MCAE.StatusControl.ValueString())

		if block.MCAE.EnhancedConvergence.ValueBool() {
			configSet = append(configSet, setPrefix+"mc-ae enhanced-convergence")
		}
		if block.MCAE.EventsIccpPeerDown != nil {
			configSet = append(configSet, setPrefix+"mc-ae events iccp-peer-down")

			if block.MCAE.EventsIccpPeerDown.ForceIclDown.ValueBool() {
				configSet = append(configSet, setPrefix+"mc-ae events iccp-peer-down force-icl-down")
			}
			if block.MCAE.EventsIccpPeerDown.PreferStatusControlActive.ValueBool() {
				configSet = append(configSet, setPrefix+"mc-ae events iccp-peer-down prefer-status-control-active")
			}
		}
		if !block.MCAE.InitDelayTime.IsNull() {
			configSet = append(configSet, setPrefix+"mc-ae init-delay-time "+
				utils.ConvI64toa(block.MCAE.InitDelayTime.ValueInt64()))
		}
		if !block.MCAE.RedundancyGroup.IsNull() {
			configSet = append(configSet, setPrefix+"mc-ae redundancy-group "+
				utils.ConvI64toa(block.MCAE.RedundancyGroup.ValueInt64()))
		}
		if !block.MCAE.RevertTime.IsNull() {
			configSet = append(configSet, setPrefix+"mc-ae revert-time "+
				utils.ConvI64toa(block.MCAE.RevertTime.ValueInt64()))
		}
		if v := block.MCAE.SwitchoverMode.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"mc-ae switchover-mode "+v)
		}
	}
	if v := block.MinimumBandwidth.ValueString(); v != "" {
		vS := strings.Split(v, " ")
		configSet = append(configSet, setPrefix+"minimum-bandwidth bw-value "+vS[0])
		if len(vS) > 1 {
			configSet = append(configSet, setPrefix+"minimum-bandwidth bw-unit "+vS[1])
		}
	}
	if !block.MinimumLinks.IsNull() {
		configSet = append(configSet, setPrefix+"minimum-links "+
			utils.ConvI64toa(block.MinimumLinks.ValueInt64()))
	}
	if !block.RedundancyGroup.IsNull() {
		if !strings.HasPrefix(interfaceName, "reth") {
			return configSet,
				path.Root("parent_ether_opts").AtName("redundancy_group"),
				fmt.Errorf("redundancy_group in parent_ether_opts block not compatible with this interface %q"+
					" (need to be reth* interface)", interfaceName)
		}
		configSet = append(configSet, setPrefix+"redundancy-group "+
			utils.ConvI64toa(block.RedundancyGroup.ValueInt64()))
	}
	for _, v := range block.SourceAddressFilter {
		configSet = append(configSet, setPrefix+"source-address-filter "+v.ValueString())
	}
	if block.SourceFiltering.ValueBool() {
		configSet = append(configSet, setPrefix+"source-filtering")
	}

	return configSet, path.Empty(), nil
}

func addInterfaceNC(
	_ context.Context, name, groupInterfaceDelete string, junSess *junos.Session,
) (
	err error,
) {
	if groupInterfaceDelete == "" {
		err = junSess.ConfigSet([]string{"set interfaces " + name + " disable description NC"})
	} else {
		err = junSess.ConfigSet([]string{"set interfaces " + name + " apply-groups " + groupInterfaceDelete})
	}
	if err != nil {
		return err
	}

	return nil
}

func (rscData *interfacePhysicalData) read(
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
			if strings.Contains(item, " unit ") && !strings.Contains(item, "ethernet-switching") {
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
			case balt.CutPrefixInString(&itemTrim, "aggregated-ether-options "):
				if rscData.ParentEtherOpts == nil {
					rscData.ParentEtherOpts = &interfacePhysicalBlockParentEtherOpts{}
				}
				if err := rscData.ParentEtherOpts.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "redundant-ether-options "):
				if rscData.ParentEtherOpts == nil {
					rscData.ParentEtherOpts = &interfacePhysicalBlockParentEtherOpts{}
				}
				if err := rscData.ParentEtherOpts.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "disable":
				rscData.Disable = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "encapsulation "):
				rscData.Encapsulation = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "esi "):
				if rscData.ESI == nil {
					rscData.ESI = &interfacePhysicalBlockESI{}
				}
				if err := rscData.ESI.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "ether-options "):
				if rscData.EtherOpts == nil {
					rscData.EtherOpts = &interfacePhysicalBlockEtherOpts{}
				}
				rscData.EtherOpts.read(itemTrim)
			case itemTrim == "flexible-vlan-tagging":
				rscData.FlexibleVlanTagging = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "gigether-options "):
				if rscData.GigetherOpts == nil {
					rscData.GigetherOpts = &interfacePhysicalBlockEtherOpts{}
				}
				rscData.GigetherOpts.read(itemTrim)
			case itemTrim == "gratuitous-arp-reply":
				rscData.GratuitousArpReply = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "hold-time down "):
				rscData.HoldTimeDown, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "hold-time up "):
				rscData.HoldTimeUp, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "link-mode "):
				rscData.LinkMode = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "mtu "):
				rscData.Mtu, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "native-vlan-id "):
				rscData.VlanNative, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "unit 0 family ethernet-switching native-vlan-id "):
				rscData.VlanNativeNonELS = types.StringValue(itemTrim)
			case itemTrim == "no-gratuitous-arp-reply":
				rscData.NoGratuitousArpReply = types.BoolValue(true)
			case itemTrim == "no-gratuitous-arp-request":
				rscData.NoGratuitousArpRequest = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "speed "):
				rscData.Speed = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "unit 0 family ethernet-switching storm-control "):
				rscData.StormControl = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "unit 0 family ethernet-switching interface-mode trunk":
				rscData.Trunk = types.BoolValue(true)
			case itemTrim == "unit 0 family ethernet-switching port-mode trunk":
				rscData.TrunkNonELS = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "unit 0 family ethernet-switching vlan members "):
				rscData.VlanMembers = append(rscData.VlanMembers, types.StringValue(itemTrim))
			case itemTrim == "vlan-tagging":
				rscData.VlanTagging = types.BoolValue(true)
			default:
				continue
			}
		}
	}

	return nil
}

func (block *interfacePhysicalBlockESI) read(itemTrim string) error {
	identifier, err := regexp.MatchString(`^([\d\w]{2}:){9}[\d\w]{2}`, itemTrim)
	if err != nil {
		return fmt.Errorf("esi_identifier regexp error: %w", err)
	}
	switch {
	case identifier:
		block.Identifier = types.StringValue(itemTrim)
	case itemTrim == "all-active", itemTrim == "single-active":
		block.Mode = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "df-election-type "):
		block.DFElectionType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-bmac "):
		block.SourceBMAC = types.StringValue(itemTrim)
	case itemTrim == "auto-derive lacp":
		block.AutoDeriveLACP = types.BoolValue(true)
	}

	return nil
}

func (block *interfacePhysicalBlockEtherOpts) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "802.3ad "):
		block.Ae8023ad = types.StringValue(itemTrim)
	case itemTrim == "auto-negotiation":
		block.AutoNegotiation = types.BoolValue(true)
	case itemTrim == "no-auto-negotiation":
		block.NoAutoNegotiation = types.BoolValue(true)
	case itemTrim == "flow-control":
		block.FlowControl = types.BoolValue(true)
	case itemTrim == "no-flow-control":
		block.NoFlowControl = types.BoolValue(true)
	case itemTrim == "loopback":
		block.Loopback = types.BoolValue(true)
	case itemTrim == "no-loopback":
		block.NoLoopback = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "redundant-parent "):
		block.RedundantParent = types.StringValue(itemTrim)
	}
}

func (block *interfacePhysicalBlockParentEtherOpts) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "bfd-liveness-detection "):
		if block.BFDLivenessDetection == nil {
			block.BFDLivenessDetection = &interfacePhysicalBlockParentEtherOptsBlockBFDLivenessDetection{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, "local-address "):
			block.BFDLivenessDetection.LocalAddress = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "authentication algorithm "):
			block.BFDLivenessDetection.AuthenticationAlgorithm = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "authentication key-chain "):
			block.BFDLivenessDetection.AuthenticationKeyChain = types.StringValue(strings.Trim(itemTrim, "\""))
		case itemTrim == "authentication loose-check":
			block.BFDLivenessDetection.AuthenticationLooseCheck = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "detection-time threshold "):
			block.BFDLivenessDetection.DetectionTimeThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "holddown-interval "):
			block.BFDLivenessDetection.HolddownInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "minimum-interval "):
			block.BFDLivenessDetection.MinimumInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "minimum-receive-interval "):
			block.BFDLivenessDetection.MinimumReceiveInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "multiplier "):
			block.BFDLivenessDetection.Multiplier, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "neighbor "):
			block.BFDLivenessDetection.Neighbor = types.StringValue(itemTrim)
		case itemTrim == "no-adaptation":
			block.BFDLivenessDetection.NoAdaptation = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "transmit-interval minimum-interval "):
			block.BFDLivenessDetection.TransmitIntervalMinimumInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "transmit-interval threshold "):
			block.BFDLivenessDetection.TransmitIntervalThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "version "):
			block.BFDLivenessDetection.Version = types.StringValue(itemTrim)
		}
	case itemTrim == "flow-control":
		block.FlowControl = types.BoolValue(true)
	case itemTrim == "no-flow-control":
		block.NoFlowControl = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "lacp "):
		if block.LACP == nil {
			block.LACP = &interfacePhysicalBlockParentEtherOptsBlockLACP{}
		}
		switch {
		case itemTrim == "active", itemTrim == "passive":
			block.LACP.Mode = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "admin-key "):
			block.LACP.AdminKey, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "periodic "):
			block.LACP.Periodic = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "sync-reset "):
			block.LACP.SyncReset = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "system-id "):
			block.LACP.SystemID = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "system-priority "):
			block.LACP.SystemPriority, err = tfdata.ConvAtoi64Value(itemTrim)
		}
	case itemTrim == "loopback":
		block.Loopback = types.BoolValue(true)
	case itemTrim == "no-loopback":
		block.NoLoopback = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "link-speed "):
		block.LinkSpeed = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "mc-ae "):
		if block.MCAE == nil {
			block.MCAE = &interfacePhysicalBlockParentEtherOptsBlockMCAE{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, "chassis-id "):
			block.MCAE.ChassisID, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "mc-ae-id "):
			block.MCAE.MCAEID, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "mode "):
			block.MCAE.Mode = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "status-control "):
			block.MCAE.StatusControl = types.StringValue(itemTrim)
		case itemTrim == "enhanced-convergence":
			block.MCAE.EnhancedConvergence = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "events iccp-peer-down"):
			if block.MCAE.EventsIccpPeerDown == nil {
				block.MCAE.EventsIccpPeerDown = &interfacePhysicalBlockParentEtherOptsBlockMCAEBlockEventsIccpPeerDown{}
			}
			switch {
			case itemTrim == " force-icl-down":
				block.MCAE.EventsIccpPeerDown.ForceIclDown = types.BoolValue(true)
			case itemTrim == " prefer-status-control-active":
				block.MCAE.EventsIccpPeerDown.PreferStatusControlActive = types.BoolValue(true)
			}
		case balt.CutPrefixInString(&itemTrim, "init-delay-time "):
			block.MCAE.InitDelayTime, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "redundancy-group "):
			block.MCAE.RedundancyGroup, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "revert-time "):
			block.MCAE.RevertTime, err = tfdata.ConvAtoi64Value(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "switchover-mode "):
			block.MCAE.SwitchoverMode = types.StringValue(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-bandwidth bw-value "):
		block.MinimumBandwidth = types.StringValue(itemTrim + block.MinimumBandwidth.ValueString())
	case balt.CutPrefixInString(&itemTrim, "minimum-bandwidth bw-unit "):
		block.MinimumBandwidth = types.StringValue(block.MinimumBandwidth.ValueString() + " " + itemTrim)
	case balt.CutPrefixInString(&itemTrim, "minimum-links "):
		block.MinimumLinks, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "redundancy-group "):
		block.RedundancyGroup, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-address-filter "):
		block.SourceAddressFilter = append(block.SourceAddressFilter,
			types.StringValue(itemTrim),
		)
	case itemTrim == "source-filtering":
		block.SourceFiltering = types.BoolValue(true)
	}

	if err != nil {
		return err
	}

	return nil
}

func (rscData *interfacePhysicalData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := "delete interfaces " + rscData.Name.ValueString() + " "

	configSet := []string{
		delPrefix + "aggregated-ether-options",
		delPrefix + "description",
		delPrefix + "disable",
		delPrefix + "encapsulation",
		delPrefix + "esi",
		delPrefix + "ether-options",
		delPrefix + "flexible-vlan-tagging",
		delPrefix + "gigether-options",
		delPrefix + "gratuitous-arp-reply",
		delPrefix + "hold-time",
		delPrefix + "link-mode",
		delPrefix + "native-vlan-id",
		delPrefix + "no-gratuitous-arp-reply",
		delPrefix + "no-gratuitous-arp-request",
		delPrefix + "redundant-ether-options",
		delPrefix + "speed",
		delPrefix + "unit 0 family ethernet-switching interface-mode",
		delPrefix + "unit 0 family ethernet-switching native-vlan-id",
		delPrefix + "unit 0 family ethernet-switching port-mode",
		delPrefix + "unit 0 family ethernet-switching storm-control",
		delPrefix + "unit 0 family ethernet-switching vlan members",
		delPrefix + "vlan-tagging",
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *interfacePhysicalData) del(
	ctx context.Context, junSess *junos.Session,
) error {
	if junSess.HasNetconf() {
		if containsUnit, err := checkInterfacePhysicalContainsUnit(
			ctx,
			rscData.Name.ValueString(),
			junSess,
		); err != nil {
			return err
		} else if containsUnit {
			return fmt.Errorf("interface %q is used for a logical unit interface", rscData.Name.ValueString())
		}
	}

	if err := junSess.ConfigSet([]string{
		"delete interfaces " + rscData.Name.ValueString(),
	}); err != nil {
		return err
	}

	if !junSess.HasNetconf() {
		return nil
	}
	if v := rscData.Name.ValueString(); strings.HasPrefix(v, "ae") {
		aggregatedCount, err := findInterfaceAggregatedCountMax(ctx, "", v, v, junSess)
		if err != nil {
			return err
		}
		if aggregatedCount == "0" {
			err = junSess.ConfigSet([]string{"delete chassis aggregated-devices ethernet device-count"})
			if err != nil {
				return err
			}
		} else {
			err = junSess.ConfigSet([]string{"set chassis aggregated-devices ethernet device-count " + aggregatedCount})
			if err != nil {
				return err
			}
		}
	} else if rscData.EtherOpts != nil || rscData.GigetherOpts != nil {
		var aeDel string
		switch {
		case rscData.EtherOpts != nil:
			aeDel = rscData.EtherOpts.Ae8023ad.ValueString()
		case rscData.GigetherOpts != nil:
			aeDel = rscData.GigetherOpts.Ae8023ad.ValueString()
		}
		if aeDel != "" {
			lastAEchild, err := findInterfaceAggregatedLastChild(
				ctx,
				aeDel,
				rscData.Name.ValueString(),
				junSess,
			)
			if err != nil {
				return err
			}
			if lastAEchild {
				aggregatedCount, err := findInterfaceAggregatedCountMax(
					ctx,
					"",
					aeDel,
					rscData.Name.ValueString(),
					junSess,
				)
				if err != nil {
					return err
				}
				if aggregatedCount == "0" {
					err = junSess.ConfigSet([]string{"delete chassis aggregated-devices ethernet device-count"})
					if err != nil {
						return err
					}
				} else {
					err = junSess.ConfigSet([]string{"set chassis aggregated-devices ethernet device-count " + aggregatedCount})
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func delInterfaceNC(
	_ context.Context, name, groupInterfaceDelete string, junSess *junos.Session,
) error {
	delPrefix := "delete interfaces " + name + " "

	configSet := []string{
		delPrefix + "description",
		delPrefix + "disable",
	}
	if groupInterfaceDelete != "" {
		configSet = append(configSet, delPrefix+"apply-groups "+groupInterfaceDelete)
	}

	return junSess.ConfigSet(configSet)
}

func findInterfaceAggregatedLastChild(
	_ context.Context, ae, interFace string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"interfaces" + junos.PipeDisplaySetRelative)
	if err != nil {
		return false, err
	}
	lastAE := true
	for _, item := range strings.Split(showConfig, "\n") {
		if strings.HasSuffix(item, "ether-options 802.3ad "+ae) &&
			!strings.HasPrefix(item, junos.SetLS+interFace+" ") {
			lastAE = false
		}
	}

	return lastAE, nil
}

func findInterfaceAggregatedCountMax(
	ctx context.Context, newAE, oldAE, interFace string, junSess *junos.Session,
) (
	string, error,
) {
	if newAE == "" {
		newAE = "ae-1"
	}
	if oldAE == "" {
		oldAE = "ae-1"
	}

	newAENumInt, err := strconv.Atoi(strings.TrimPrefix(newAE, "ae"))
	if err != nil {
		return "", fmt.Errorf("converting ae interaface '%v' to integer: %w", newAE, err)
	}
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"interfaces" + junos.PipeDisplaySetRelative)
	if err != nil {
		return "", err
	}

	listAEFound := make([]string, 0)
	regexpAEchild := regexp.MustCompile(`ether-options 802\.3ad ae\d+$`)
	regexpAEparent := regexp.MustCompile(`^set ae\d+ `)

	for _, line := range strings.Split(showConfig, "\n") {
		aeMatchChild := regexpAEchild.MatchString(line)
		aeMatchParent := regexpAEparent.MatchString(line)

		switch {
		case aeMatchChild:
			wordsLine := strings.Fields(line)
			if interFace == oldAE {
				// findInterfaceAggregatedCountMax called for delete parent interface
				listAEFound = append(listAEFound, wordsLine[len(wordsLine)-1])
			} else if wordsLine[len(wordsLine)-1] != oldAE {
				listAEFound = append(listAEFound, wordsLine[len(wordsLine)-1])
			}
		case aeMatchParent:
			wordsLine := strings.Fields(line)
			if interFace != oldAE {
				// findInterfaceAggregatedCountMax called for child interface or new parent
				listAEFound = append(listAEFound, wordsLine[1])
			} else if wordsLine[1] != oldAE {
				listAEFound = append(listAEFound, wordsLine[1])
			}
		}
	}

	lastOldAE, err := findInterfaceAggregatedLastChild(ctx, oldAE, interFace, junSess)
	if err != nil {
		return "", err
	}
	if !lastOldAE {
		listAEFound = append(listAEFound, oldAE)
	}

	if len(listAEFound) > 0 {
		balt.SortStringsByLengthInc(listAEFound)
		lastAeInt, err := strconv.Atoi(strings.TrimPrefix(listAEFound[len(listAEFound)-1], "ae"))
		if err != nil {
			return "", fmt.Errorf("converting internal variable lastAeInt to integer: %w", err)
		}

		if lastAeInt > newAENumInt {
			return strconv.Itoa(lastAeInt + 1), nil
		}
	}

	return strconv.Itoa(newAENumInt + 1), nil
}
