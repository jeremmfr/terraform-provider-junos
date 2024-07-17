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
	_ resource.Resource                   = &ospfArea{}
	_ resource.ResourceWithConfigure      = &ospfArea{}
	_ resource.ResourceWithValidateConfig = &ospfArea{}
	_ resource.ResourceWithImportState    = &ospfArea{}
	_ resource.ResourceWithUpgradeState   = &ospfArea{}
)

type ospfArea struct {
	client *junos.Client
}

func newOspfAreaResource() resource.Resource {
	return &ospfArea{}
}

func (rsc *ospfArea) typeName() string {
	return providerName + "_ospf_area"
}

func (rsc *ospfArea) junosName() string {
	return "ospf|ospf3 area"
}

func (rsc *ospfArea) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *ospfArea) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *ospfArea) Configure(
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

func (rsc *ospfArea) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<area_id>" + junos.IDSeparator + "<version>" +
					junos.IDSeparator + "<routing_instance>` or " +
					"`<area_id>" + junos.IDSeparator + "<version>" +
					junos.IDSeparator + "<realm>" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"area_id": schema.StringAttribute{
				Required:    true,
				Description: "Area ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.Any(
						tfvalidator.StringIPAddress().IPv4Only(),
						stringvalidator.RegexMatches(regexp.MustCompile(
							`^\d+$`),
							"should be usually in the IP format (but a number is accepted)",
						),
					),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("v2"),
				Description: "Version of ospf.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("v2", "v3"),
				},
			},
			"realm": schema.StringAttribute{
				Optional:    true,
				Description: "OSPFv3 realm configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("ipv4-unicast", "ipv4-multicast", "ipv6-multicast"),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for ospf area.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"context_identifier": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Configure context identifier in support of edge protection.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringIPAddress().IPv4Only(),
					),
				},
			},
			"inter_area_prefix_export": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Export policy for Inter Area Prefix LSAs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"inter_area_prefix_import": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Import policy for Inter Area Prefix LSAs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"network_summary_export": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Export policy for Type 3 Summary LSAs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"network_summary_import": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Import policy for Type 3 Summary LSAs.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"no_context_identifier_advertisement": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable context identifier advertisements in this area.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"interface": schema.ListNestedBlock{
				Description: "For each interface or interface-range to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of interface or interface-range.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
							},
						},
						"authentication_simple_password": schema.StringAttribute{
							Optional:    true,
							Sensitive:   true,
							Description: "Authentication key.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"dead_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Dead interval (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"demand_circuit": schema.BoolAttribute{
							Optional:    true,
							Description: "Interface functions as a demand circuit.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"disable": schema.BoolAttribute{
							Optional:    true,
							Description: "Disable OSPF on this interface.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"dynamic_neighbors": schema.BoolAttribute{
							Optional:    true,
							Description: "Learn neighbors dynamically on a p2mp interface.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"flood_reduction": schema.BoolAttribute{
							Optional:    true,
							Description: "Enable flood reduction.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"hello_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Hello interval (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 255),
							},
						},
						"interface_type": schema.StringAttribute{
							Optional:    true,
							Description: "Type of interface.",
							Validators: []validator.String{
								stringvalidator.OneOf("nbma", "p2mp", "p2mp-over-lan", "p2p"),
							},
						},
						"ipsec_sa": schema.StringAttribute{
							Optional:    true,
							Description: "IPSec security association name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"ipv4_adjacency_segment_protected_type": schema.StringAttribute{
							Optional:    true,
							Description: "Type to define adjacency SID is eligible for protection.",
							Validators: []validator.String{
								stringvalidator.OneOf("dynamic", "index", "label"),
							},
						},
						"ipv4_adjacency_segment_protected_value": schema.StringAttribute{
							Optional:    true,
							Description: "Value for index or label to define adjacency SID is eligible for protection.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(
									`^\d+$`),
									"should be a numeric value",
								),
							},
						},
						"ipv4_adjacency_segment_unprotected_type": schema.StringAttribute{
							Optional:    true,
							Description: "Type to define adjacency SID uneligible for protection.",
							Validators: []validator.String{
								stringvalidator.OneOf("dynamic", "index", "label"),
							},
						},
						"ipv4_adjacency_segment_unprotected_value": schema.StringAttribute{
							Optional:    true,
							Description: "Value for index or label to define adjacency SID uneligible for protection.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(
									`^\d+$`),
									"should be a numeric value",
								),
							},
						},
						"link_protection": schema.BoolAttribute{
							Optional:    true,
							Description: "Protect interface from link faults only.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"metric": schema.Int64Attribute{
							Optional:    true,
							Description: "Interface metric.",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"mtu": schema.Int64Attribute{
							Optional:    true,
							Description: "Maximum OSPF packet size.",
							Validators: []validator.Int64{
								int64validator.Between(128, 65535),
							},
						},
						"no_advertise_adjacency_segment": schema.BoolAttribute{
							Optional:    true,
							Description: "Do not advertise an adjacency segment for this interface.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"no_eligible_backup": schema.BoolAttribute{
							Optional:    true,
							Description: "Not eligible to backup traffic from protected interfaces.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"no_eligible_remote_backup": schema.BoolAttribute{
							Optional:    true,
							Description: "Not eligible for Remote-LFA backup traffic from protected interfaces.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"no_interface_state_traps": schema.BoolAttribute{
							Optional:    true,
							Description: "Do not send interface state change traps.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"no_neighbor_down_notification": schema.BoolAttribute{
							Optional:    true,
							Description: "Don't inform other protocols about neighbor down events.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"node_link_protection": schema.BoolAttribute{
							Optional:    true,
							Description: "Protect interface from both link and node faults.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"passive": schema.BoolAttribute{
							Optional:    true,
							Description: "Do not run OSPF, but advertise it.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"passive_traffic_engineering_remote_node_id": schema.StringAttribute{
							Optional:    true,
							Description: "Advertise TE link information, remote address of the link.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress(),
							},
						},
						"passive_traffic_engineering_remote_node_router_id": schema.StringAttribute{
							Optional:    true,
							Description: "Advertise TE link information, TE Router-ID of the remote node.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress(),
							},
						},
						"poll_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Poll interval for NBMA interfaces.",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"priority": schema.Int64Attribute{
							Optional:    true,
							Description: "Designated router priority.",
							Validators: []validator.Int64{
								int64validator.Between(0, 255),
							},
						},
						"retransmit_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Retransmission interval (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"secondary": schema.BoolAttribute{
							Optional:    true,
							Description: "Treat interface as secondary.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"strict_bfd": schema.BoolAttribute{
							Optional:    true,
							Description: "Enable strict bfd over this interface.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"te_metric": schema.Int64Attribute{
							Optional:    true,
							Description: "Traffic engineering metric.",
							Validators: []validator.Int64{
								int64validator.Between(1, 4294967295),
							},
						},
						"transit_delay": schema.Int64Attribute{
							Optional:    true,
							Description: "Transit delay (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"authentication_md5": schema.ListNestedBlock{
							Description: "For each key_id, MD5 authentication key.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"key_id": schema.Int64Attribute{
										Required:    true,
										Description: "Key ID for MD5 authentication.",
										Validators: []validator.Int64{
											int64validator.Between(0, 255),
										},
									},
									"key": schema.StringAttribute{
										Required:    true,
										Sensitive:   true,
										Description: "MD5 authentication key value.",
										Validators: []validator.String{
											stringvalidator.LengthAtLeast(1),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"start_time": schema.StringAttribute{
										Optional:    true,
										Description: "Start time for key transmission.",
										Validators: []validator.String{
											stringvalidator.RegexMatches(regexp.MustCompile(
												`^\d{4}\-\d\d?\-\d\d?\.\d{2}:\d{2}:\d{2}$`),
												"must be in the format 'YYYY-MM-DD.HH:MM:SS'",
											),
										},
									},
								},
							},
						},
						"bandwidth_based_metrics": schema.SetNestedBlock{
							Description: "For each bandwidth, configure bandwidth based metrics.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"bandwidth": schema.StringAttribute{
										Required:    true,
										Description: "Bandwidth threshold.",
										Validators: []validator.String{
											stringvalidator.RegexMatches(regexp.MustCompile(
												`^(\d)+(m|k|g)?$`),
												`must be a bandwidth ^(\d)+(m|k|g)?$`),
										},
									},
									"metric": schema.Int64Attribute{
										Required:    true,
										Description: "Metric associated with specified bandwidth.",
										Validators: []validator.Int64{
											int64validator.Between(1, 65535),
										},
									},
								},
							},
						},
						"bfd_liveness_detection": schema.SingleNestedBlock{
							Description: "Bidirectional Forwarding Detection options.",
							Attributes: map[string]schema.Attribute{
								"authentication_algorithm": schema.StringAttribute{
									Optional:    true,
									Description: "Authentication algorithm name.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									},
								},
								"authentication_key_chain": schema.StringAttribute{
									Optional:    true,
									Description: "Authentication key chain name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 128),
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
								"full_neighbors_only": schema.BoolAttribute{
									Optional:    true,
									Description: "Setup BFD sessions only to Full neighbors.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"holddown_interval": schema.Int64Attribute{
									Optional:    true,
									Description: "Time to hold the session-UP notification to the client (milliseconds).",
									Validators: []validator.Int64{
										int64validator.Between(1, 255000),
									},
								},
								"minimum_interval": schema.Int64Attribute{
									Optional:    true,
									Description: "Minimum transmit and receive interval (milliseconds).",
									Validators: []validator.Int64{
										int64validator.Between(1, 255000),
									},
								},
								"minimum_receive_interval": schema.Int64Attribute{
									Optional:    true,
									Description: "Minimum receive interval (milliseconds).",
									Validators: []validator.Int64{
										int64validator.Between(1, 255000),
									},
								},
								"multiplier": schema.Int64Attribute{
									Optional:    true,
									Description: "Detection time multiplier.",
									Validators: []validator.Int64{
										int64validator.Between(1, 255),
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
									Description: "Minimum transmit interval (milliseconds).",
									Validators: []validator.Int64{
										int64validator.Between(1, 255000),
									},
								},
								"transmit_interval_threshold": schema.Int64Attribute{
									Optional:    true,
									Description: "High transmit interval triggering a trap (milliseconds).",
									Validators: []validator.Int64{
										int64validator.Between(1, 4294967295),
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
						"neighbor": schema.SetNestedBlock{
							Description: "For each address, configure NBMA neighbor.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"address": schema.StringAttribute{
										Required:    true,
										Description: "Address of neighbor.",
										Validators: []validator.String{
											tfvalidator.StringIPAddress(),
										},
									},
									"eligible": schema.BoolAttribute{
										Optional:    true,
										Description: "Eligible to be DR on an NBMA network.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
								},
							},
						},
					},
				},
			},
			"area_range": schema.SetNestedBlock{
				Description: "For each `range`, configure area range.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"range": schema.StringAttribute{
							Required:    true,
							Description: "Range to summarize routes in this area.",
							Validators: []validator.String{
								tfvalidator.StringCIDRNetwork(),
							},
						},
						"exact": schema.BoolAttribute{
							Optional:    true,
							Description: "Enforce exact match for advertisement of this area range.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"override_metric": schema.Int64Attribute{
							Optional:    true,
							Description: "Override the dynamic metric for this area-range.",
							Validators: []validator.Int64{
								int64validator.Between(1, 16777215),
							},
						},
						"restrict": schema.BoolAttribute{
							Optional:    true,
							Description: "Restrict advertisement of this area range.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
					},
				},
			},
			"nssa": schema.SingleNestedBlock{
				Description: "Configure a not-so-stubby area.",
				Attributes: map[string]schema.Attribute{
					"summaries": schema.BoolAttribute{
						Optional:    true,
						Description: "Flood summary LSAs into this NSSA area.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_summaries": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't flood summary LSAs into this NSSA area.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"area_range": schema.SetNestedBlock{
						Description: "For each `range`, configure area range.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"range": schema.StringAttribute{
									Required:    true,
									Description: "Range to summarize routes in this area.",
									Validators: []validator.String{
										tfvalidator.StringCIDRNetwork(),
									},
								},
								"exact": schema.BoolAttribute{
									Optional:    true,
									Description: "Enforce exact match for advertisement of this area range.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"override_metric": schema.Int64Attribute{
									Optional:    true,
									Description: "Override the dynamic metric for this area-range.",
									Validators: []validator.Int64{
										int64validator.Between(1, 16777215),
									},
								},
								"restrict": schema.BoolAttribute{
									Optional:    true,
									Description: "Restrict advertisement of this area range.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
							},
						},
					},
					"default_lsa": schema.SingleNestedBlock{
						Description: "Configure a default LSA.",
						Attributes: map[string]schema.Attribute{
							"default_metric": schema.Int64Attribute{
								Optional:    true,
								Description: "Metric for the default route in this area.",
								Validators: []validator.Int64{
									int64validator.Between(1, 16777215),
								},
							},
							"metric_type": schema.Int64Attribute{
								Optional:    true,
								Description: "External metric type for the default type 7 LSA.",
								Validators: []validator.Int64{
									int64validator.Between(1, 2),
								},
							},
							"type_7": schema.BoolAttribute{
								Optional:    true,
								Description: "Flood type 7 default LSA if no-summaries is configured.",
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
			"stub": schema.SingleNestedBlock{
				Description: "onfigure a stub area.",
				Attributes: map[string]schema.Attribute{
					"default_metric": schema.Int64Attribute{
						Optional:    true,
						Description: "Metric for the default route in this stub area.",
						Validators: []validator.Int64{
							int64validator.Between(1, 16777215),
						},
					},
					"summaries": schema.BoolAttribute{
						Optional:    true,
						Description: "Flood summary LSAs into this stub area.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_summaries": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't flood summary LSAs into this stub area.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"virtual_link": schema.SetNestedBlock{
				Description: "For each combination of `neighbor_id` and `transit_area`, configure virtual link.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"neighbor_id": schema.StringAttribute{
							Required:    true,
							Description: "Router ID of a virtual neighbor.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress().IPv4Only(),
							},
						},
						"transit_area": schema.StringAttribute{
							Required:    true,
							Description: "Transit area in common with virtual neighbor.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress().IPv4Only(),
							},
						},
						"dead_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Dead interval (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"demand_circuit": schema.BoolAttribute{
							Optional:    true,
							Description: "Interface functions as a demand circuit.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"disable": schema.BoolAttribute{
							Optional:    true,
							Description: "Disable this virtual link.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"flood_reduction": schema.BoolAttribute{
							Optional:    true,
							Description: "Enable flood reduction.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"hello_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Hello interval (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 255),
							},
						},
						"ipsec_sa": schema.StringAttribute{
							Optional:    true,
							Description: "IPSec security association name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"mtu": schema.Int64Attribute{
							Optional:    true,
							Description: "Maximum OSPF packet size.",
							Validators: []validator.Int64{
								int64validator.Between(128, 65535),
							},
						},
						"retransmit_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Retransmission interval (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"transit_delay": schema.Int64Attribute{
							Optional:    true,
							Description: "Transit delay (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
					},
				},
			},
		},
	}
}

type ospfAreaData struct {
	ID                               types.String               `tfsdk:"id"`
	AreaID                           types.String               `tfsdk:"area_id"`
	Version                          types.String               `tfsdk:"version"`
	Realm                            types.String               `tfsdk:"realm"`
	RoutingInstance                  types.String               `tfsdk:"routing_instance"`
	ContextIdentifier                []types.String             `tfsdk:"context_identifier"`
	InterAreaPrefixExport            []types.String             `tfsdk:"inter_area_prefix_export"`
	InterAreaPrefixImport            []types.String             `tfsdk:"inter_area_prefix_import"`
	NetworkSummaryExport             []types.String             `tfsdk:"network_summary_export"`
	NetworkSummaryImport             []types.String             `tfsdk:"network_summary_import"`
	NoContextIdentifierAdvertisement types.Bool                 `tfsdk:"no_context_identifier_advertisement"`
	Interface                        []ospfAreaBlockInterface   `tfsdk:"interface"`
	AreaRange                        []ospfAreaBlockAreaRange   `tfsdk:"area_range"`
	Nssa                             *ospfAreaBlockNssa         `tfsdk:"nssa"`
	Stub                             *ospfAreaBlockStub         `tfsdk:"stub"`
	VirtualLink                      []ospfAreaBlockVirtualLink `tfsdk:"virtual_link"`
}

type ospfAreaConfig struct {
	ID                               types.String             `tfsdk:"id"`
	AreaID                           types.String             `tfsdk:"area_id"`
	Version                          types.String             `tfsdk:"version"`
	Realm                            types.String             `tfsdk:"realm"`
	RoutingInstance                  types.String             `tfsdk:"routing_instance"`
	ContextIdentifier                types.Set                `tfsdk:"context_identifier"`
	InterAreaPrefixExport            types.List               `tfsdk:"inter_area_prefix_export"`
	InterAreaPrefixImport            types.List               `tfsdk:"inter_area_prefix_import"`
	NetworkSummaryExport             types.List               `tfsdk:"network_summary_export"`
	NetworkSummaryImport             types.List               `tfsdk:"network_summary_import"`
	NoContextIdentifierAdvertisement types.Bool               `tfsdk:"no_context_identifier_advertisement"`
	Interface                        types.List               `tfsdk:"interface"`
	AreaRange                        types.Set                `tfsdk:"area_range"`
	Nssa                             *ospfAreaBlockNssaConfig `tfsdk:"nssa"`
	Stub                             *ospfAreaBlockStub       `tfsdk:"stub"`
	VirtualLink                      types.Set                `tfsdk:"virtual_link"`
}

//nolint:lll
type ospfAreaBlockInterface struct {
	Name                                        types.String                                       `tfsdk:"name"`
	AuthenticationSimplePassword                types.String                                       `tfsdk:"authentication_simple_password"`
	DeadInterval                                types.Int64                                        `tfsdk:"dead_interval"`
	DemandCircuit                               types.Bool                                         `tfsdk:"demand_circuit"`
	Disable                                     types.Bool                                         `tfsdk:"disable"`
	DynamicNeighbors                            types.Bool                                         `tfsdk:"dynamic_neighbors"`
	FloodReduction                              types.Bool                                         `tfsdk:"flood_reduction"`
	HelloInterval                               types.Int64                                        `tfsdk:"hello_interval"`
	InterfaceType                               types.String                                       `tfsdk:"interface_type"`
	IpsecSA                                     types.String                                       `tfsdk:"ipsec_sa"`
	IPv4AdjacencySegmentProtectedType           types.String                                       `tfsdk:"ipv4_adjacency_segment_protected_type"`
	IPv4AdjacencySegmentProtectedValue          types.String                                       `tfsdk:"ipv4_adjacency_segment_protected_value"`
	IPv4AdjacencySegmentUnprotectedType         types.String                                       `tfsdk:"ipv4_adjacency_segment_unprotected_type"`
	IPv4AdjacencySegmentUnprotectedValue        types.String                                       `tfsdk:"ipv4_adjacency_segment_unprotected_value"`
	LinkProtection                              types.Bool                                         `tfsdk:"link_protection"`
	Metric                                      types.Int64                                        `tfsdk:"metric"`
	Mtu                                         types.Int64                                        `tfsdk:"mtu"`
	NoAdvertiseAdjacencySegment                 types.Bool                                         `tfsdk:"no_advertise_adjacency_segment"`
	NoEligibleBackup                            types.Bool                                         `tfsdk:"no_eligible_backup"`
	NoEligibleRemoteBackup                      types.Bool                                         `tfsdk:"no_eligible_remote_backup"`
	NoInterfaceStateTraps                       types.Bool                                         `tfsdk:"no_interface_state_traps"`
	NoNeighborDownNotification                  types.Bool                                         `tfsdk:"no_neighbor_down_notification"`
	NodeLinkProtection                          types.Bool                                         `tfsdk:"node_link_protection"`
	Passive                                     types.Bool                                         `tfsdk:"passive"`
	PassiveTrafficEngineeringRemoteNodeID       types.String                                       `tfsdk:"passive_traffic_engineering_remote_node_id"`
	PassiveTrafficEngineeringRemoteNodeRouterID types.String                                       `tfsdk:"passive_traffic_engineering_remote_node_router_id"`
	PollInterval                                types.Int64                                        `tfsdk:"poll_interval"`
	Priority                                    types.Int64                                        `tfsdk:"priority"`
	RetransmitInterval                          types.Int64                                        `tfsdk:"retransmit_interval"`
	Secondary                                   types.Bool                                         `tfsdk:"secondary"`
	StrictBfd                                   types.Bool                                         `tfsdk:"strict_bfd"`
	TeMetric                                    types.Int64                                        `tfsdk:"te_metric"`
	TransitDelay                                types.Int64                                        `tfsdk:"transit_delay"`
	AuthenticationMD5                           []ospfAreaBlockInterfaceBlockAuthenticationMD5     `tfsdk:"authentication_md5"`
	BandwidthBasedMetrics                       []ospfAreaBlockInterfaceBlockBandwidthBasedMetrics `tfsdk:"bandwidth_based_metrics"`
	BfdLivenessDetection                        *ospfAreaBlockInterfaceBlockBfdLivenessDetection   `tfsdk:"bfd_liveness_detection"`
	Neighbor                                    []ospfAreaBlockInterfaceBlockNeighbor              `tfsdk:"neighbor"`
}

//nolint:lll
type ospfAreaBlockInterfaceConfig struct {
	Name                                        types.String                                     `tfsdk:"name"`
	AuthenticationSimplePassword                types.String                                     `tfsdk:"authentication_simple_password"`
	DeadInterval                                types.Int64                                      `tfsdk:"dead_interval"`
	DemandCircuit                               types.Bool                                       `tfsdk:"demand_circuit"`
	Disable                                     types.Bool                                       `tfsdk:"disable"`
	DynamicNeighbors                            types.Bool                                       `tfsdk:"dynamic_neighbors"`
	FloodReduction                              types.Bool                                       `tfsdk:"flood_reduction"`
	HelloInterval                               types.Int64                                      `tfsdk:"hello_interval"`
	InterfaceType                               types.String                                     `tfsdk:"interface_type"`
	IpsecSA                                     types.String                                     `tfsdk:"ipsec_sa"`
	IPv4AdjacencySegmentProtectedType           types.String                                     `tfsdk:"ipv4_adjacency_segment_protected_type"`
	IPv4AdjacencySegmentProtectedValue          types.String                                     `tfsdk:"ipv4_adjacency_segment_protected_value"`
	IPv4AdjacencySegmentUnprotectedType         types.String                                     `tfsdk:"ipv4_adjacency_segment_unprotected_type"`
	IPv4AdjacencySegmentUnprotectedValue        types.String                                     `tfsdk:"ipv4_adjacency_segment_unprotected_value"`
	LinkProtection                              types.Bool                                       `tfsdk:"link_protection"`
	Metric                                      types.Int64                                      `tfsdk:"metric"`
	Mtu                                         types.Int64                                      `tfsdk:"mtu"`
	NoAdvertiseAdjacencySegment                 types.Bool                                       `tfsdk:"no_advertise_adjacency_segment"`
	NoEligibleBackup                            types.Bool                                       `tfsdk:"no_eligible_backup"`
	NoEligibleRemoteBackup                      types.Bool                                       `tfsdk:"no_eligible_remote_backup"`
	NoInterfaceStateTraps                       types.Bool                                       `tfsdk:"no_interface_state_traps"`
	NoNeighborDownNotification                  types.Bool                                       `tfsdk:"no_neighbor_down_notification"`
	NodeLinkProtection                          types.Bool                                       `tfsdk:"node_link_protection"`
	Passive                                     types.Bool                                       `tfsdk:"passive"`
	PassiveTrafficEngineeringRemoteNodeID       types.String                                     `tfsdk:"passive_traffic_engineering_remote_node_id"`
	PassiveTrafficEngineeringRemoteNodeRouterID types.String                                     `tfsdk:"passive_traffic_engineering_remote_node_router_id"`
	PollInterval                                types.Int64                                      `tfsdk:"poll_interval"`
	Priority                                    types.Int64                                      `tfsdk:"priority"`
	RetransmitInterval                          types.Int64                                      `tfsdk:"retransmit_interval"`
	Secondary                                   types.Bool                                       `tfsdk:"secondary"`
	StrictBfd                                   types.Bool                                       `tfsdk:"strict_bfd"`
	TeMetric                                    types.Int64                                      `tfsdk:"te_metric"`
	TransitDelay                                types.Int64                                      `tfsdk:"transit_delay"`
	AuthenticationMD5                           types.List                                       `tfsdk:"authentication_md5"`
	BandwidthBasedMetrics                       types.Set                                        `tfsdk:"bandwidth_based_metrics"`
	BfdLivenessDetection                        *ospfAreaBlockInterfaceBlockBfdLivenessDetection `tfsdk:"bfd_liveness_detection"`
	Neighbor                                    types.Set                                        `tfsdk:"neighbor"`
}

type ospfAreaBlockInterfaceBlockAuthenticationMD5 struct {
	KeyID     types.Int64  `tfsdk:"key_id"`
	Key       types.String `tfsdk:"key"`
	StartTime types.String `tfsdk:"start_time"`
}

type ospfAreaBlockInterfaceBlockBandwidthBasedMetrics struct {
	Bandwidth types.String `tfsdk:"bandwidth"`
	Metric    types.Int64  `tfsdk:"metric"`
}

type ospfAreaBlockInterfaceBlockBfdLivenessDetection struct {
	AuthenticationAlgorithm         types.String `tfsdk:"authentication_algorithm"`
	AuthenticationKeyChain          types.String `tfsdk:"authentication_key_chain"`
	AuthenticationLooseCheck        types.Bool   `tfsdk:"authentication_loose_check"`
	DetectionTimeThreshold          types.Int64  `tfsdk:"detection_time_threshold"`
	FullNeighborsOnly               types.Bool   `tfsdk:"full_neighbors_only"`
	HolddownInterval                types.Int64  `tfsdk:"holddown_interval"`
	MinimumInterval                 types.Int64  `tfsdk:"minimum_interval"`
	MinimumReceiveInterval          types.Int64  `tfsdk:"minimum_receive_interval"`
	Multiplier                      types.Int64  `tfsdk:"multiplier"`
	NoAdaptation                    types.Bool   `tfsdk:"no_adaptation"`
	TransmitIntervalMinimumInterval types.Int64  `tfsdk:"transmit_interval_minimum_interval"`
	TransmitIntervalThreshold       types.Int64  `tfsdk:"transmit_interval_threshold"`
	Version                         types.String `tfsdk:"version"`
}

func (block *ospfAreaBlockInterfaceBlockBfdLivenessDetection) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type ospfAreaBlockInterfaceBlockNeighbor struct {
	Address  types.String `tfsdk:"address"`
	Eligbile types.Bool   `tfsdk:"eligible"`
}

type ospfAreaBlockAreaRange struct {
	Range          types.String `tfsdk:"range"`
	Exact          types.Bool   `tfsdk:"exact"`
	OverrideMetric types.Int64  `tfsdk:"override_metric"`
	Restrict       types.Bool   `tfsdk:"restrict"`
}

type ospfAreaBlockNssa struct {
	Summaries   types.Bool                        `tfsdk:"summaries"`
	NoSummaries types.Bool                        `tfsdk:"no_summaries"`
	AreaRange   []ospfAreaBlockAreaRange          `tfsdk:"area_range"`
	DefaultLsa  *ospfAreaBlockNssaBlockDefaultLsa `tfsdk:"default_lsa"`
}

type ospfAreaBlockNssaConfig struct {
	Summaries   types.Bool                        `tfsdk:"summaries"`
	NoSummaries types.Bool                        `tfsdk:"no_summaries"`
	AreaRange   types.Set                         `tfsdk:"area_range"`
	DefaultLsa  *ospfAreaBlockNssaBlockDefaultLsa `tfsdk:"default_lsa"`
}

func (block *ospfAreaBlockNssaConfig) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type ospfAreaBlockNssaBlockDefaultLsa struct {
	DefaultMetric types.Int64 `tfsdk:"default_metric"`
	MetricType    types.Int64 `tfsdk:"metric_type"`
	Type7         types.Bool  `tfsdk:"type_7"`
}

type ospfAreaBlockStub struct {
	DefaultMetric types.Int64 `tfsdk:"default_metric"`
	Summaries     types.Bool  `tfsdk:"summaries"`
	NoSummaries   types.Bool  `tfsdk:"no_summaries"`
}

func (block *ospfAreaBlockStub) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type ospfAreaBlockVirtualLink struct {
	NeighborID         types.String `tfsdk:"neighbor_id"`
	TransitArea        types.String `tfsdk:"transit_area"`
	DeadInterval       types.Int64  `tfsdk:"dead_interval"`
	DemandCircuit      types.Bool   `tfsdk:"demand_circuit"`
	Disable            types.Bool   `tfsdk:"disable"`
	FloodReduction     types.Bool   `tfsdk:"flood_reduction"`
	HelloInterval      types.Int64  `tfsdk:"hello_interval"`
	IpsecSA            types.String `tfsdk:"ipsec_sa"`
	Mtu                types.Int64  `tfsdk:"mtu"`
	RetransmitInterval types.Int64  `tfsdk:"retransmit_interval"`
	TransitDelay       types.Int64  `tfsdk:"transit_delay"`
}

func (rsc *ospfArea) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config ospfAreaConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Version.IsNull() && !config.Version.IsUnknown() {
		switch config.Version.ValueString() {
		case "v2":
			if !config.Realm.IsNull() && !config.Realm.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("realm"),
					tfdiag.ConflictConfigErrSummary,
					"realm cannot be configured when version = v2",
				)
			}
			if !config.InterAreaPrefixExport.IsNull() && !config.InterAreaPrefixExport.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("inter_area_prefix_export"),
					tfdiag.ConflictConfigErrSummary,
					"inter_area_prefix_export cannot be configured when version = v2",
				)
			}
			if !config.InterAreaPrefixImport.IsNull() && !config.InterAreaPrefixImport.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("inter_area_prefix_import"),
					tfdiag.ConflictConfigErrSummary,
					"inter_area_prefix_import cannot be configured when version = v2",
				)
			}
		case "v3":
			if !config.NetworkSummaryExport.IsNull() && !config.NetworkSummaryExport.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("network_summary_export"),
					tfdiag.ConflictConfigErrSummary,
					"network_summary_export cannot be configured when version = v3",
				)
			}
			if !config.NetworkSummaryImport.IsNull() && !config.NetworkSummaryImport.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("network_summary_import"),
					tfdiag.ConflictConfigErrSummary,
					"network_summary_import cannot be configured when version = v3",
				)
			}
		}
	}
	if !config.ContextIdentifier.IsNull() && !config.ContextIdentifier.IsUnknown() &&
		!config.NoContextIdentifierAdvertisement.IsNull() && !config.NoContextIdentifierAdvertisement.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("context_identifier"),
			tfdiag.ConflictConfigErrSummary,
			"context_identifier and no_context_identifier_advertisement cannot be configured together",
		)
	}
	if config.Interface.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("interface"),
			tfdiag.MissingConfigErrSummary,
			"interface block must be specified",
		)
	} else if !config.Interface.IsUnknown() {
		var configInterface []ospfAreaBlockInterfaceConfig
		asDiags := config.Interface.ElementsAs(ctx, &configInterface, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		interfaceName := make(map[string]struct{})
		for i, block := range configInterface {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := interfaceName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple interface blocks with the same name %q", name),
					)
				}
				interfaceName[name] = struct{}{}
			}
			if !block.AuthenticationSimplePassword.IsNull() && !block.AuthenticationSimplePassword.IsUnknown() &&
				!block.AuthenticationMD5.IsNull() && !block.AuthenticationMD5.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("interface").AtListIndex(i).AtName("authentication_simple_password"),
					tfdiag.ConflictConfigErrSummary,
					fmt.Sprintf("authentication_simple_password and authentication_md5 cannot be configured together"+
						" in interface block %q", block.Name.ValueString()),
				)
			}
			if !block.IPv4AdjacencySegmentProtectedValue.IsNull() && !block.IPv4AdjacencySegmentProtectedValue.IsUnknown() &&
				block.IPv4AdjacencySegmentProtectedType.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("interface").AtListIndex(i).AtName("ipv4_adjacency_segment_protected_value"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("ipv4_adjacency_segment_protected_type must be specified with ipv4_adjacency_segment_protected_value"+
						" in interface block %q", block.Name.ValueString()),
				)
			}
			if !block.IPv4AdjacencySegmentUnprotectedValue.IsNull() && !block.IPv4AdjacencySegmentUnprotectedValue.IsUnknown() &&
				block.IPv4AdjacencySegmentUnprotectedType.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("interface").AtListIndex(i).AtName("ipv4_adjacency_segment_unprotected_value"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("ipv4_adjacency_segment_unprotected_type must be specified"+
						" with ipv4_adjacency_segment_unprotected_value"+
						" in interface block %q", block.Name.ValueString()),
				)
			}
			if block.Passive.IsNull() {
				if !block.PassiveTrafficEngineeringRemoteNodeID.IsNull() &&
					!block.PassiveTrafficEngineeringRemoteNodeID.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface").AtListIndex(i).AtName("passive_traffic_engineering_remote_node_id"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("passive must be specified with passive_traffic_engineering_remote_node_id"+
							" in interface block %q", block.Name.ValueString()),
					)
				}
				if !block.PassiveTrafficEngineeringRemoteNodeRouterID.IsNull() &&
					!block.PassiveTrafficEngineeringRemoteNodeRouterID.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface").AtListIndex(i).AtName("passive_traffic_engineering_remote_node_router_id"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("passive must be specified with passive_traffic_engineering_remote_node_router_id"+
							" in interface block %q", block.Name.ValueString()),
					)
				}
			}

			if !block.AuthenticationMD5.IsNull() && !block.AuthenticationMD5.IsUnknown() {
				var configAuthenticationMD5 []ospfAreaBlockInterfaceBlockAuthenticationMD5
				asDiags := block.AuthenticationMD5.ElementsAs(ctx, &configAuthenticationMD5, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				authenticationMD5KeyID := make(map[int64]struct{})
				for ii, blockAuthenticationMD5 := range configAuthenticationMD5 {
					if blockAuthenticationMD5.KeyID.IsUnknown() {
						continue
					}
					keyID := blockAuthenticationMD5.KeyID.ValueInt64()
					if _, ok := authenticationMD5KeyID[keyID]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("interface").AtListIndex(i).AtName("authentication_md5").AtListIndex(ii).AtName("key_id"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple authentication_md5 blocks with the same key_id %d"+
								" in interface block %q", keyID, block.Name.ValueString()),
						)
					}
					authenticationMD5KeyID[keyID] = struct{}{}
				}
			}
			if !block.BandwidthBasedMetrics.IsNull() && !block.BandwidthBasedMetrics.IsUnknown() {
				var configBandwidthBasedMetrics []ospfAreaBlockInterfaceBlockBandwidthBasedMetrics
				asDiags := block.BandwidthBasedMetrics.ElementsAs(ctx, &configBandwidthBasedMetrics, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				bandwidthBasedMetricsBandwidth := make(map[string]struct{})
				for _, blockBandwidthBasedMetrics := range configBandwidthBasedMetrics {
					if blockBandwidthBasedMetrics.Bandwidth.IsUnknown() {
						continue
					}
					bandwidth := blockBandwidthBasedMetrics.Bandwidth.ValueString()
					if _, ok := bandwidthBasedMetricsBandwidth[bandwidth]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("interface").AtListIndex(i).AtName("bandwidth_based_metrics"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple bandwidth_based_metrics blocks with the same bandwidth %q"+
								" in interface block %q", bandwidth, block.Name.ValueString()),
						)
					}
					bandwidthBasedMetricsBandwidth[bandwidth] = struct{}{}
				}
			}
			if block.BfdLivenessDetection != nil {
				if block.BfdLivenessDetection.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface").AtListIndex(i).AtName("bfd_liveness_detection").AtName("*"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("bfd_liveness_detection block is empty"+
							" in interface block %q", block.Name.ValueString()),
					)
				}
			}
			if !block.Neighbor.IsNull() && !block.Neighbor.IsUnknown() {
				var configNeighbor []ospfAreaBlockInterfaceBlockNeighbor
				asDiags := block.Neighbor.ElementsAs(ctx, &configNeighbor, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				neighborAddress := make(map[string]struct{})
				for _, blockNeighbor := range configNeighbor {
					if blockNeighbor.Address.IsUnknown() {
						continue
					}
					address := blockNeighbor.Address.ValueString()
					if _, ok := neighborAddress[address]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("interface").AtListIndex(i).AtName("neighbor"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple neighbor blocks with the same address %q"+
								" in interface block %q", address, block.Name.ValueString()),
						)
					}
					neighborAddress[address] = struct{}{}
				}
			}
		}
	}
	if !config.AreaRange.IsNull() && !config.AreaRange.IsUnknown() {
		var configAreaRange []ospfAreaBlockAreaRange
		asDiags := config.AreaRange.ElementsAs(ctx, &configAreaRange, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		areaRangeRange := make(map[string]struct{})
		for _, block := range configAreaRange {
			if block.Range.IsUnknown() {
				continue
			}
			rangeValue := block.Range.ValueString()
			if _, ok := areaRangeRange[rangeValue]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("area_range"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple area_range blocks with the same range %q", rangeValue),
				)
			}
			areaRangeRange[rangeValue] = struct{}{}
		}
	}
	if config.Nssa != nil && config.Nssa.hasKnownValue() &&
		config.Stub != nil && config.Stub.hasKnownValue() {
		resp.Diagnostics.AddAttributeError(
			path.Root("nssa"),
			tfdiag.ConflictConfigErrSummary,
			"nssa and stub cannot be configured together",
		)
	}
	if config.Nssa != nil {
		if !config.Nssa.Summaries.IsNull() && !config.Nssa.Summaries.IsUnknown() &&
			!config.Nssa.NoSummaries.IsNull() && !config.Nssa.NoSummaries.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nssa").AtName("summaries"),
				tfdiag.ConflictConfigErrSummary,
				"summaries and no_summaries cannot be configured together"+
					" in nssa block",
			)
		}
		if !config.Nssa.AreaRange.IsNull() && !config.Nssa.AreaRange.IsUnknown() {
			var configAreaRange []ospfAreaBlockAreaRange
			asDiags := config.Nssa.AreaRange.ElementsAs(ctx, &configAreaRange, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			areaRangeRange := make(map[string]struct{})
			for _, block := range configAreaRange {
				if block.Range.IsUnknown() {
					continue
				}
				rangeValue := block.Range.ValueString()
				if _, ok := areaRangeRange[rangeValue]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("nssa").AtName("area_range"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple area_range blocks with the same range %q"+
							" in nssa block", rangeValue),
					)
				}
				areaRangeRange[rangeValue] = struct{}{}
			}
		}
	}
	if config.Stub != nil {
		if !config.Stub.Summaries.IsNull() && !config.Stub.Summaries.IsUnknown() &&
			!config.Stub.NoSummaries.IsNull() && !config.Stub.NoSummaries.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("stub").AtName("summaries"),
				tfdiag.ConflictConfigErrSummary,
				"summaries and no_summaries cannot be configured together"+
					" in stub block",
			)
		}
	}
	if !config.VirtualLink.IsNull() && !config.VirtualLink.IsUnknown() {
		var configVirtualLink []ospfAreaBlockVirtualLink
		asDiags := config.VirtualLink.ElementsAs(ctx, &configVirtualLink, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		virtualLinkNeighborIDTransitArea := make(map[string]struct{})
		for _, block := range configVirtualLink {
			if block.NeighborID.IsUnknown() {
				continue
			}
			if block.TransitArea.IsUnknown() {
				continue
			}
			neighborID := block.NeighborID.ValueString()
			transitArea := block.TransitArea.ValueString()
			if _, ok := virtualLinkNeighborIDTransitArea[neighborID+junos.IDSeparator+transitArea]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("virtual_link"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple virtual_link blocks with the same neighbor_id %q and transit_area %q",
						neighborID, transitArea),
				)
			}
			virtualLinkNeighborIDTransitArea[neighborID+junos.IDSeparator+transitArea] = struct{}{}
		}
	}
}

func (rsc *ospfArea) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan ospfAreaData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.AreaID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("area_id"),
			"Empty Area ID",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "area_id"),
		)

		return
	}
	if version := plan.Version.ValueString(); version == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("version"),
			"Empty Version",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "version"),
		)

		return
	} else if version == "v2" && plan.Realm.ValueString() != "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("realm"),
			tfdiag.ConflictConfigErrSummary,
			"realm cannot be configured when version = v2",
		)
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
			areaExists, err := checkOspfAreaExists(
				fnCtx,
				plan.AreaID.ValueString(),
				plan.Version.ValueString(),
				plan.Realm.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if areaExists {
				var inRoutingInstanceMessage string
				if v := plan.RoutingInstance.ValueString(); v != "" {
					inRoutingInstanceMessage = fmt.Sprintf(" in routing-instance %q", v)
				}
				switch plan.Version.ValueString() {
				case "v2":
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("ospf area %q already exists"+inRoutingInstanceMessage, plan.AreaID.ValueString()),
					)
				case "v3":
					if realm := plan.Realm.ValueString(); realm != "" {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("ospf3 realm %q area %q already exists"+inRoutingInstanceMessage, realm, plan.AreaID.ValueString()),
						)
					} else {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("ospf3 area %q already exists"+inRoutingInstanceMessage, realm, plan.AreaID.ValueString()),
						)
					}
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			areaExists, err := checkOspfAreaExists(
				fnCtx,
				plan.AreaID.ValueString(),
				plan.Version.ValueString(),
				plan.Realm.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !areaExists {
				var inRoutingInstanceMessage string
				if v := plan.RoutingInstance.ValueString(); v != "" {
					inRoutingInstanceMessage = fmt.Sprintf("in routing-instance %q", v)
				}
				switch plan.Version.ValueString() {
				case "v2":
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						fmt.Sprintf("ospf area %q does not exists after commit "+inRoutingInstanceMessage+
							"=> check your config", plan.AreaID.ValueString()),
					)
				case "v3":
					if realm := plan.Realm.ValueString(); realm != "" {
						resp.Diagnostics.AddError(
							tfdiag.NotFoundErrSummary,
							fmt.Sprintf("ospf3 realm %q area %q does not exists after commit "+inRoutingInstanceMessage+
								"=> check your config", realm, plan.AreaID.ValueString()),
						)
					} else {
						resp.Diagnostics.AddError(
							tfdiag.NotFoundErrSummary,
							fmt.Sprintf("ospf3 area %q does not exists after commit "+inRoutingInstanceMessage+
								"=> check your config", realm, plan.AreaID.ValueString()),
						)
					}
				}

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *ospfArea) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data ospfAreaData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom4String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.AreaID.ValueString(),
			state.Version.ValueString(),
			state.Realm.ValueString(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *ospfArea) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state ospfAreaData
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

func (rsc *ospfArea) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state ospfAreaData
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

func (rsc *ospfArea) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	var data ospfAreaData
	idSplit := strings.Split(req.ID, junos.IDSeparator)
	switch {
	case len(idSplit) < 3:
		resp.Diagnostics.AddError(
			"Bad ID Format",
			fmt.Sprintf("missing element(s) in id with separator %q", junos.IDSeparator),
		)

		return
	case len(idSplit) == 4:
		if err := data.read(ctx, idSplit[0], idSplit[1], idSplit[2], idSplit[3], junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
	default:
		if err := data.read(ctx, idSplit[0], idSplit[1], "", idSplit[2], junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be "+
				"<aread_id>"+junos.IDSeparator+"<version>"+junos.IDSeparator+"<routing_instance> or "+
				"<aread_id>"+junos.IDSeparator+"<version>"+junos.IDSeparator+"<realm>"+junos.IDSeparator+"<routing_instance>)",
		)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkOspfAreaExists(
	_ context.Context, areaID, version, realm, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	if version == "v3" {
		showPrefix += "protocols " + junos.OspfV3 + " "
	} else {
		showPrefix += "protocols " + junos.OspfV2 + " "
		if realm != "" {
			return false, errors.New("realm can't set if version != v3")
		}
	}
	if realm != "" {
		showPrefix += "realm " + realm + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"area " + areaID + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *ospfAreaData) fillID() {
	routingInstance := junos.DefaultW
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		routingInstance = v
	}
	if realm := rscData.Realm.ValueString(); realm != "" {
		rscData.ID = types.StringValue(
			rscData.AreaID.ValueString() + junos.IDSeparator +
				rscData.Version.ValueString() + junos.IDSeparator +
				realm + junos.IDSeparator +
				routingInstance,
		)
	} else {
		rscData.ID = types.StringValue(
			rscData.AreaID.ValueString() + junos.IDSeparator +
				rscData.Version.ValueString() + junos.IDSeparator +
				routingInstance,
		)
	}
}

func (rscData *ospfAreaData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *ospfAreaData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	if rscData.Version.ValueString() == "v3" {
		setPrefix += "protocols " + junos.OspfV3 + " "
	} else {
		setPrefix += "protocols " + junos.OspfV2 + " "
		if rscData.Realm.ValueString() != "" {
			return path.Root("realm"), errors.New("realm can't set if version != v3")
		}
	}
	if v := rscData.Realm.ValueString(); v != "" {
		setPrefix += "realm " + v + " "
	}
	setPrefix += "area " + rscData.AreaID.ValueString() + " "

	for _, v := range rscData.ContextIdentifier {
		configSet = append(configSet, setPrefix+"context-identifier "+v.ValueString())
	}
	for _, v := range rscData.InterAreaPrefixExport {
		configSet = append(configSet, setPrefix+"inter-area-prefix-export \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.InterAreaPrefixImport {
		configSet = append(configSet, setPrefix+"inter-area-prefix-import \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.NetworkSummaryExport {
		configSet = append(configSet, setPrefix+"network-summary-export \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.NetworkSummaryImport {
		configSet = append(configSet, setPrefix+"network-summary-import \""+v.ValueString()+"\"")
	}
	if rscData.NoContextIdentifierAdvertisement.ValueBool() {
		configSet = append(configSet, setPrefix+"no-context-identifier-advertisement")
	}

	interfaceName := make(map[string]struct{})
	for i, block := range rscData.Interface {
		name := block.Name.ValueString()
		if _, ok := interfaceName[name]; ok {
			return path.Root("interface").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple interface blocks with the same name %q", name)
		}
		interfaceName[name] = struct{}{}

		blockSet, pathErr, err := block.configSet(setPrefix, path.Root("interface").AtListIndex(i))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	areaRangeRange := make(map[string]struct{})
	for _, block := range rscData.AreaRange {
		rangeValue := block.Range.ValueString()
		if _, ok := areaRangeRange[rangeValue]; ok {
			return path.Root("area_range"),
				fmt.Errorf("multiple area_range blocks with the same range %q", rangeValue)
		}
		areaRangeRange[rangeValue] = struct{}{}

		configSet = append(configSet, block.configSet(setPrefix)...)
	}
	if rscData.Nssa != nil {
		blockSet, pathErr, err := rscData.Nssa.configSet(setPrefix, path.Root("nssa"))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.Stub != nil {
		configSet = append(configSet, setPrefix+"stub")

		if !rscData.Stub.DefaultMetric.IsNull() {
			configSet = append(configSet, setPrefix+"stub default-metric "+
				utils.ConvI64toa(rscData.Stub.DefaultMetric.ValueInt64()))
		}
		if rscData.Stub.Summaries.ValueBool() {
			configSet = append(configSet, setPrefix+"stub summaries")
		}
		if rscData.Stub.NoSummaries.ValueBool() {
			configSet = append(configSet, setPrefix+"stub no-summaries")
		}
	}
	virtualLinkNeighborIDTransitArea := make(map[string]struct{})
	for _, block := range rscData.VirtualLink {
		neighborID := block.NeighborID.ValueString()
		transitArea := block.TransitArea.ValueString()
		if _, ok := virtualLinkNeighborIDTransitArea[neighborID+junos.IDSeparator+transitArea]; ok {
			return path.Root("virtual_link"),
				fmt.Errorf("multiple virtual_link blocks with the same neighbor_id %q and transit_area %q", neighborID, transitArea)
		}
		virtualLinkNeighborIDTransitArea[neighborID+junos.IDSeparator+transitArea] = struct{}{}

		configSet = append(configSet, block.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *ospfAreaBlockInterface) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "interface " + block.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if v := block.AuthenticationSimplePassword.ValueString(); v != "" {
		if len(block.AuthenticationMD5) > 0 {
			return configSet,
				pathRoot.AtName("authentication_md5"),
				fmt.Errorf("authentication_simple_password and authentication_md5 cannot be configured together"+
					" in interface block %q", block.Name.ValueString())
		}
		configSet = append(configSet, setPrefix+"authentication simple-password \""+v+"\"")
	}
	if !block.DeadInterval.IsNull() {
		configSet = append(configSet, setPrefix+"dead-interval "+
			utils.ConvI64toa(block.DeadInterval.ValueInt64()))
	}
	if block.DemandCircuit.ValueBool() {
		configSet = append(configSet, setPrefix+"demand-circuit")
	}
	if block.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if block.DynamicNeighbors.ValueBool() {
		configSet = append(configSet, setPrefix+"dynamic-neighbors")
	}
	if block.FloodReduction.ValueBool() {
		configSet = append(configSet, setPrefix+"flood-reduction")
	}
	if !block.HelloInterval.IsNull() {
		configSet = append(configSet, setPrefix+"hello-interval "+
			utils.ConvI64toa(block.HelloInterval.ValueInt64()))
	}
	if v := block.InterfaceType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"interface-type "+v)
	}
	if v := block.IpsecSA.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ipsec-sa \""+v+"\"")
	}
	if t := block.IPv4AdjacencySegmentProtectedType.ValueString(); t != "" {
		if v := block.IPv4AdjacencySegmentProtectedValue.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"ipv4-adjacency-segment protected "+t+" \""+v+"\"")
		} else {
			configSet = append(configSet, setPrefix+"ipv4-adjacency-segment protected "+t)
		}
	} else if block.IPv4AdjacencySegmentProtectedValue.ValueString() != "" {
		return configSet,
			pathRoot.AtName("ipv4_adjacency_segment_protected_value"),
			fmt.Errorf("ipv4_adjacency_segment_protected_type must be specified with ipv4_adjacency_segment_protected_value"+
				" in interface block %q", block.Name.ValueString())
	}
	if t := block.IPv4AdjacencySegmentUnprotectedType.ValueString(); t != "" {
		if v := block.IPv4AdjacencySegmentUnprotectedValue.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"ipv4-adjacency-segment unprotected "+t+" \""+v+"\"")
		} else {
			configSet = append(configSet, setPrefix+"ipv4-adjacency-segment unprotected "+t)
		}
	} else if block.IPv4AdjacencySegmentUnprotectedValue.ValueString() != "" {
		return configSet,
			pathRoot.AtName("ipv4_adjacency_segment_unprotected_value"),
			fmt.Errorf("ipv4_adjacency_segment_unprotected_type must be specified with ipv4_adjacency_segment_unprotected_value"+
				" in interface block %q", block.Name.ValueString())
	}
	if block.LinkProtection.ValueBool() {
		configSet = append(configSet, setPrefix+"link-protection")
	}
	if !block.Metric.IsNull() {
		configSet = append(configSet, setPrefix+"metric "+
			utils.ConvI64toa(block.Metric.ValueInt64()))
	}
	if !block.Mtu.IsNull() {
		configSet = append(configSet, setPrefix+"mtu "+
			utils.ConvI64toa(block.Mtu.ValueInt64()))
	}
	if block.NoAdvertiseAdjacencySegment.ValueBool() {
		configSet = append(configSet, setPrefix+"no-advertise-adjacency-segment")
	}
	if block.NoEligibleBackup.ValueBool() {
		configSet = append(configSet, setPrefix+"no-eligible-backup")
	}
	if block.NoEligibleRemoteBackup.ValueBool() {
		configSet = append(configSet, setPrefix+"no-eligible-remote-backup")
	}
	if block.NoInterfaceStateTraps.ValueBool() {
		configSet = append(configSet, setPrefix+"no-interface-state-traps")
	}
	if block.NoNeighborDownNotification.ValueBool() {
		configSet = append(configSet, setPrefix+"no-neighbor-down-notification")
	}
	if block.NodeLinkProtection.ValueBool() {
		configSet = append(configSet, setPrefix+"node-link-protection")
	}
	if block.Passive.ValueBool() {
		configSet = append(configSet, setPrefix+"passive")
		if v := block.PassiveTrafficEngineeringRemoteNodeID.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"passive traffic-engineering remote-node-id "+v)
		}
		if v := block.PassiveTrafficEngineeringRemoteNodeRouterID.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"passive traffic-engineering remote-node-router-id "+v)
		}
	} else {
		if block.PassiveTrafficEngineeringRemoteNodeID.ValueString() != "" {
			return configSet,
				pathRoot.AtName("passive_traffic_engineering_remote_node_id"),
				fmt.Errorf("passive must be specified with passive_traffic_engineering_remote_node_id"+
					" in interface block %q", block.Name.ValueString())
		}
		if block.PassiveTrafficEngineeringRemoteNodeRouterID.ValueString() != "" {
			return configSet,
				pathRoot.AtName("passive_traffic_engineering_remote_node_router_id"),
				fmt.Errorf("passive must be specified with passive_traffic_engineering_remote_node_router_id"+
					" in interface block %q", block.Name.ValueString())
		}
	}
	if !block.PollInterval.IsNull() {
		configSet = append(configSet, setPrefix+"poll-interval "+
			utils.ConvI64toa(block.PollInterval.ValueInt64()))
	}
	if !block.Priority.IsNull() {
		configSet = append(configSet, setPrefix+"priority "+
			utils.ConvI64toa(block.Priority.ValueInt64()))
	}
	if !block.RetransmitInterval.IsNull() {
		configSet = append(configSet, setPrefix+"retransmit-interval "+
			utils.ConvI64toa(block.RetransmitInterval.ValueInt64()))
	}
	if block.Secondary.ValueBool() {
		configSet = append(configSet, setPrefix+"secondary")
	}
	if block.StrictBfd.ValueBool() {
		configSet = append(configSet, setPrefix+"strict-bfd")
	}
	if !block.TeMetric.IsNull() {
		configSet = append(configSet, setPrefix+"te-metric "+
			utils.ConvI64toa(block.TeMetric.ValueInt64()))
	}
	if !block.TransitDelay.IsNull() {
		configSet = append(configSet, setPrefix+"transit-delay "+
			utils.ConvI64toa(block.TransitDelay.ValueInt64()))
	}

	authenticationMD5KeyID := make(map[int64]struct{})
	for i, blockAuthenticationMD5 := range block.AuthenticationMD5 {
		keyID := blockAuthenticationMD5.KeyID.ValueInt64()
		if _, ok := authenticationMD5KeyID[keyID]; ok {
			return configSet,
				pathRoot.AtName("authentication_md5").AtListIndex(i).AtName("key_id"),
				fmt.Errorf("multiple authentication_md5 blocks with the same key_id %d"+
					" in interface block %q", keyID, block.Name.ValueString())
		}
		authenticationMD5KeyID[keyID] = struct{}{}

		configSet = append(configSet, setPrefix+"authentication md5 "+
			utils.ConvI64toa(keyID)+" key \""+blockAuthenticationMD5.Key.ValueString()+"\"")
		if v := blockAuthenticationMD5.StartTime.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"authentication md5 "+
				utils.ConvI64toa(keyID)+" start-time "+v)
		}
	}
	bandwidthBasedMetricsBandwidth := make(map[string]struct{})
	for _, blockBandwidthBasedMetrics := range block.BandwidthBasedMetrics {
		bandwidth := blockBandwidthBasedMetrics.Bandwidth.ValueString()
		if _, ok := bandwidthBasedMetricsBandwidth[bandwidth]; ok {
			return configSet,
				pathRoot.AtName("bandwidth_based_metrics"),
				fmt.Errorf("multiple bandwidth_based_metrics blocks with the same bandwidth %q"+
					" in interface block %q", bandwidth, block.Name.ValueString())
		}
		bandwidthBasedMetricsBandwidth[bandwidth] = struct{}{}

		configSet = append(configSet, setPrefix+"bandwidth-based-metrics bandwidth "+bandwidth+
			" metric "+utils.ConvI64toa(blockBandwidthBasedMetrics.Metric.ValueInt64()))
	}
	if block.BfdLivenessDetection != nil {
		if block.BfdLivenessDetection.isEmpty() {
			return configSet,
				path.Root("bfd_liveness_detection").AtName("*"),
				fmt.Errorf("bfd_liveness_detection block is empty"+
					" in interface block %q", block.Name.ValueString())
		}

		configSet = append(configSet, block.BfdLivenessDetection.configSet(setPrefix)...)
	}
	neighborAddress := make(map[string]struct{})
	for _, blockNeighbor := range block.Neighbor {
		address := blockNeighbor.Address.ValueString()
		if _, ok := neighborAddress[address]; ok {
			return configSet,
				pathRoot.AtName("neighbor"),
				fmt.Errorf("multiple neighbor blocks with the same address %q"+
					" in interface block %q", address, block.Name.ValueString())
		}
		neighborAddress[address] = struct{}{}

		configSet = append(configSet, setPrefix+"neighbor "+address)
		if blockNeighbor.Eligbile.ValueBool() {
			configSet = append(configSet, setPrefix+"neighbor "+address+" eligible")
		}
	}

	return configSet, path.Empty(), nil
}

func (block *ospfAreaBlockInterfaceBlockBfdLivenessDetection) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 1)
	setPrefix += "bfd-liveness-detection "

	if v := block.AuthenticationAlgorithm.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication algorithm "+v)
	}
	if v := block.AuthenticationKeyChain.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication key-chain \""+v+"\"")
	}
	if block.AuthenticationLooseCheck.ValueBool() {
		configSet = append(configSet, setPrefix+"authentication loose-check")
	}
	if !block.DetectionTimeThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"detection-time threshold "+
			utils.ConvI64toa(block.DetectionTimeThreshold.ValueInt64()))
	}
	if block.FullNeighborsOnly.ValueBool() {
		configSet = append(configSet, setPrefix+"full-neighbors-only")
	}
	if !block.HolddownInterval.IsNull() {
		configSet = append(configSet, setPrefix+"holddown-interval "+
			utils.ConvI64toa(block.HolddownInterval.ValueInt64()))
	}
	if !block.MinimumInterval.IsNull() {
		configSet = append(configSet, setPrefix+"minimum-interval "+
			utils.ConvI64toa(block.MinimumInterval.ValueInt64()))
	}
	if !block.MinimumReceiveInterval.IsNull() {
		configSet = append(configSet, setPrefix+"minimum-receive-interval "+
			utils.ConvI64toa(block.MinimumReceiveInterval.ValueInt64()))
	}
	if !block.Multiplier.IsNull() {
		configSet = append(configSet, setPrefix+"multiplier "+
			utils.ConvI64toa(block.Multiplier.ValueInt64()))
	}
	if block.NoAdaptation.ValueBool() {
		configSet = append(configSet, setPrefix+"no-adaptation")
	}
	if !block.TransmitIntervalMinimumInterval.IsNull() {
		configSet = append(configSet, setPrefix+"transmit-interval minimum-interval "+
			utils.ConvI64toa(block.TransmitIntervalMinimumInterval.ValueInt64()))
	}
	if !block.TransmitIntervalThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"transmit-interval threshold "+
			utils.ConvI64toa(block.TransmitIntervalThreshold.ValueInt64()))
	}
	if v := block.Version.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"version "+v)
	}

	return configSet
}

func (block *ospfAreaBlockAreaRange) configSet(setPrefix string) []string {
	setPrefix += "area-range " + block.Range.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if block.Exact.ValueBool() {
		configSet = append(configSet, setPrefix+"exact")
	}
	if !block.OverrideMetric.IsNull() {
		configSet = append(configSet, setPrefix+"override-metric "+
			utils.ConvI64toa(block.OverrideMetric.ValueInt64()))
	}
	if block.Restrict.ValueBool() {
		configSet = append(configSet, setPrefix+"restrict")
	}

	return configSet
}

func (block *ospfAreaBlockNssa) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "nssa "

	configSet := []string{
		setPrefix,
	}

	if block.Summaries.ValueBool() {
		configSet = append(configSet, setPrefix+"summaries")
	}
	if block.NoSummaries.ValueBool() {
		configSet = append(configSet, setPrefix+"no-summaries")
	}
	areaRangeRange := make(map[string]struct{})
	for _, blockAreaRange := range block.AreaRange {
		rangeValue := blockAreaRange.Range.ValueString()
		if _, ok := areaRangeRange[rangeValue]; ok {
			return configSet,
				pathRoot.AtName("area_range"),
				fmt.Errorf("multiple area_range blocks with the same range %q"+
					" in nssa block", rangeValue)
		}
		areaRangeRange[rangeValue] = struct{}{}

		configSet = append(configSet, blockAreaRange.configSet(setPrefix)...)
	}
	if block.DefaultLsa != nil {
		configSet = append(configSet, setPrefix+"default-lsa")

		if !block.DefaultLsa.DefaultMetric.IsNull() {
			configSet = append(configSet, setPrefix+"default-lsa default-metric "+
				utils.ConvI64toa(block.DefaultLsa.DefaultMetric.ValueInt64()))
		}
		if !block.DefaultLsa.MetricType.IsNull() {
			configSet = append(configSet, setPrefix+"default-lsa metric-type "+
				utils.ConvI64toa(block.DefaultLsa.MetricType.ValueInt64()))
		}
		if block.DefaultLsa.Type7.ValueBool() {
			configSet = append(configSet, setPrefix+"default-lsa type-7")
		}
	}

	return configSet, path.Empty(), nil
}

func (block *ospfAreaBlockVirtualLink) configSet(setPrefix string) []string {
	setPrefix += "virtual-link" +
		" neighbor-id " + block.NeighborID.ValueString() +
		" transit-area " + block.TransitArea.ValueString() +
		" "

	configSet := []string{
		setPrefix,
	}

	if !block.DeadInterval.IsNull() {
		configSet = append(configSet, setPrefix+"dead-interval "+
			utils.ConvI64toa(block.DeadInterval.ValueInt64()))
	}
	if block.DemandCircuit.ValueBool() {
		configSet = append(configSet, setPrefix+"demand-circuit")
	}
	if block.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if block.FloodReduction.ValueBool() {
		configSet = append(configSet, setPrefix+"flood-reduction")
	}
	if !block.HelloInterval.IsNull() {
		configSet = append(configSet, setPrefix+"hello-interval "+
			utils.ConvI64toa(block.HelloInterval.ValueInt64()))
	}
	if v := block.IpsecSA.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ipsec-sa \""+v+"\"")
	}
	if !block.Mtu.IsNull() {
		configSet = append(configSet, setPrefix+"mtu "+
			utils.ConvI64toa(block.Mtu.ValueInt64()))
	}
	if !block.RetransmitInterval.IsNull() {
		configSet = append(configSet, setPrefix+"retransmit-interval "+
			utils.ConvI64toa(block.RetransmitInterval.ValueInt64()))
	}
	if !block.TransitDelay.IsNull() {
		configSet = append(configSet, setPrefix+"transit-delay "+
			utils.ConvI64toa(block.TransitDelay.ValueInt64()))
	}

	return configSet
}

func (rscData *ospfAreaData) read(
	_ context.Context, areaID, version, realm, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	if version == "v3" {
		showPrefix += "protocols " + junos.OspfV3 + " "
	} else {
		showPrefix += "protocols " + junos.OspfV2 + " "
		if realm != "" {
			return errors.New("realm can't set if version != v3")
		}
	}
	if realm != "" {
		showPrefix += "realm " + realm + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"area " + areaID + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.AreaID = types.StringValue(areaID)
		if version == "v3" {
			rscData.Version = types.StringValue(version)
		} else {
			rscData.Version = types.StringValue("v2")
		}
		if realm != "" {
			rscData.Realm = types.StringValue(realm)
		}
		if routingInstance != "" {
			rscData.RoutingInstance = types.StringValue(routingInstance)
		} else {
			rscData.RoutingInstance = types.StringValue(junos.DefaultW)
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
			case balt.CutPrefixInString(&itemTrim, "context-identifier "):
				rscData.ContextIdentifier = append(rscData.ContextIdentifier,
					types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "inter-area-prefix-export "):
				rscData.InterAreaPrefixExport = append(rscData.InterAreaPrefixExport,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "inter-area-prefix-import "):
				rscData.InterAreaPrefixImport = append(rscData.InterAreaPrefixImport,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "network-summary-export "):
				rscData.NetworkSummaryExport = append(rscData.NetworkSummaryExport,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "network-summary-import "):
				rscData.NetworkSummaryImport = append(rscData.NetworkSummaryImport,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case itemTrim == "no-context-identifier-advertisement":
				rscData.NoContextIdentifierAdvertisement = types.BoolValue(true)

			case balt.CutPrefixInString(&itemTrim, "interface "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var interFace ospfAreaBlockInterface
				rscData.Interface, interFace = tfdata.ExtractBlockWithTFTypesString(
					rscData.Interface, "Name", itemTrimFields[0],
				)
				interFace.Name = types.StringValue(itemTrimFields[0])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")

				if err := interFace.read(itemTrim); err != nil {
					return err
				}
				rscData.Interface = append(rscData.Interface, interFace)
			case balt.CutPrefixInString(&itemTrim, "area-range "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var areaRange ospfAreaBlockAreaRange
				rscData.AreaRange, areaRange = tfdata.ExtractBlockWithTFTypesString(
					rscData.AreaRange, "Range", itemTrimFields[0],
				)
				areaRange.Range = types.StringValue(itemTrimFields[0])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")

				if err := areaRange.read(itemTrim); err != nil {
					return err
				}
				rscData.AreaRange = append(rscData.AreaRange, areaRange)
			case balt.CutPrefixInString(&itemTrim, "nssa"):
				if rscData.Nssa == nil {
					rscData.Nssa = &ospfAreaBlockNssa{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.Nssa.read(itemTrim); err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "stub"):
				if rscData.Stub == nil {
					rscData.Stub = &ospfAreaBlockStub{}
				}

				switch {
				case balt.CutPrefixInString(&itemTrim, " default-metric "):
					rscData.Stub.DefaultMetric, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case itemTrim == " summaries":
					rscData.Stub.Summaries = types.BoolValue(true)
				case itemTrim == " no-summaries":
					rscData.Stub.NoSummaries = types.BoolValue(true)
				}
			case balt.CutPrefixInString(&itemTrim, "virtual-link "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 4 { // neighbor-id <neighbor_id> transit-area <transit_area>
					return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "virtual-link", itemTrim)
				}
				var virtualLink ospfAreaBlockVirtualLink
				rscData.VirtualLink, virtualLink = tfdata.ExtractBlockWith2TFTypesString(
					rscData.VirtualLink, "NeighborID", itemTrimFields[1], "TransitArea", itemTrimFields[3],
				)
				virtualLink.NeighborID = types.StringValue(itemTrimFields[1])
				virtualLink.TransitArea = types.StringValue(itemTrimFields[3])
				balt.CutPrefixInString(&itemTrim, "neighbor-id "+itemTrimFields[1]+" transit-area "+itemTrimFields[3]+" ")

				switch {
				case balt.CutPrefixInString(&itemTrim, "dead-interval "):
					virtualLink.DeadInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				case itemTrim == "demand-circuit":
					virtualLink.DemandCircuit = types.BoolValue(true)
				case itemTrim == "disable":
					virtualLink.Disable = types.BoolValue(true)
				case itemTrim == "flood-reduction":
					virtualLink.FloodReduction = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "hello-interval "):
					virtualLink.HelloInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "ipsec-sa "):
					virtualLink.IpsecSA = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "mtu "):
					virtualLink.Mtu, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "retransmit-interval "):
					virtualLink.RetransmitInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "transit-delay "):
					virtualLink.TransitDelay, err = tfdata.ConvAtoi64Value(itemTrim)
				}
				if err != nil {
					return err
				}
				rscData.VirtualLink = append(rscData.VirtualLink, virtualLink)
			}
		}
	}

	return nil
}

func (block *ospfAreaBlockInterface) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication simple-password "):
		block.AuthenticationSimplePassword, err = tfdata.JunosDecode(
			strings.Trim(itemTrim, "\""),
			"authentication simple-password",
		)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "dead-interval "):
		block.DeadInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "demand-circuit":
		block.DemandCircuit = types.BoolValue(true)
	case itemTrim == junos.DisableW:
		block.Disable = types.BoolValue(true)
	case itemTrim == "dynamic-neighbors":
		block.DynamicNeighbors = types.BoolValue(true)
	case itemTrim == "flood-reduction":
		block.FloodReduction = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "hello-interval "):
		block.HelloInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "interface-type "):
		block.InterfaceType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ipsec-sa "):
		block.IpsecSA = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "ipv4-adjacency-segment protected "):
		itemTrimFields := strings.Split(itemTrim, " ")
		block.IPv4AdjacencySegmentProtectedType = types.StringValue(itemTrimFields[0])
		if len(itemTrimFields) > 1 { // <type> <value>
			block.IPv4AdjacencySegmentProtectedValue = types.StringValue(itemTrimFields[1])
		}
	case balt.CutPrefixInString(&itemTrim, "ipv4-adjacency-segment unprotected "):
		itemTrimFields := strings.Split(itemTrim, " ")
		block.IPv4AdjacencySegmentUnprotectedType = types.StringValue(itemTrimFields[0])
		if len(itemTrimFields) > 1 { // <type> <value>
			block.IPv4AdjacencySegmentUnprotectedValue = types.StringValue(itemTrimFields[1])
		}
	case itemTrim == "link-protection":
		block.LinkProtection = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "metric "):
		block.Metric, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "mtu "):
		block.Mtu, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "no-advertise-adjacency-segment":
		block.NoAdvertiseAdjacencySegment = types.BoolValue(true)
	case itemTrim == "no-eligible-backup":
		block.NoEligibleBackup = types.BoolValue(true)
	case itemTrim == "no-eligible-remote-backup":
		block.NoEligibleRemoteBackup = types.BoolValue(true)
	case itemTrim == "no-interface-state-traps":
		block.NoInterfaceStateTraps = types.BoolValue(true)
	case itemTrim == "no-neighbor-down-notification":
		block.NoNeighborDownNotification = types.BoolValue(true)
	case itemTrim == "node-link-protection":
		block.NodeLinkProtection = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "passive"):
		block.Passive = types.BoolValue(true)
		switch {
		case balt.CutPrefixInString(&itemTrim, " traffic-engineering remote-node-id "):
			block.PassiveTrafficEngineeringRemoteNodeID = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, " traffic-engineering remote-node-router-id "):
			block.PassiveTrafficEngineeringRemoteNodeRouterID = types.StringValue(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "poll-interval "):
		block.PollInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "priority "):
		block.Priority, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "retransmit-interval "):
		block.RetransmitInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "secondary":
		block.Secondary = types.BoolValue(true)
	case itemTrim == "strict-bfd":
		block.StrictBfd = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "te-metric "):
		block.TeMetric, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "transit-delay "):
		block.TransitDelay, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}

	case balt.CutPrefixInString(&itemTrim, "authentication md5 "):
		itemTrimFields := strings.Split(itemTrim, " ")
		keyID, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
		if err != nil {
			return err
		}
		var authenticationMD5 ospfAreaBlockInterfaceBlockAuthenticationMD5
		block.AuthenticationMD5, authenticationMD5 = tfdata.ExtractBlockWithTFTypesInt64(
			block.AuthenticationMD5, "KeyID", keyID.ValueInt64(),
		)
		authenticationMD5.KeyID = keyID
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")

		switch {
		case balt.CutPrefixInString(&itemTrim, "key "):
			authenticationMD5.Key, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "authentication md5 key")
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "start-time "):
			authenticationMD5.StartTime = types.StringValue(strings.Split(strings.Trim(itemTrim, "\""), " ")[0])
		}
		block.AuthenticationMD5 = append(block.AuthenticationMD5, authenticationMD5)
	case balt.CutPrefixInString(&itemTrim, "bandwidth-based-metrics bandwidth "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 3 { // <bandwidth> metric <metric>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "bandwidth-based-metrics bandwidth", itemTrim)
		}
		metric, err := tfdata.ConvAtoi64Value(itemTrimFields[2])
		if err != nil {
			return err
		}

		block.BandwidthBasedMetrics = append(block.BandwidthBasedMetrics, ospfAreaBlockInterfaceBlockBandwidthBasedMetrics{
			Bandwidth: types.StringValue(itemTrimFields[0]),
			Metric:    metric,
		})
	case balt.CutPrefixInString(&itemTrim, "bfd-liveness-detection "):
		if block.BfdLivenessDetection == nil {
			block.BfdLivenessDetection = &ospfAreaBlockInterfaceBlockBfdLivenessDetection{}
		}

		if err := block.BfdLivenessDetection.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "neighbor "):
		itemTrimFields := strings.Split(itemTrim, " ") // <address> (eligible)?
		block.Neighbor = append(block.Neighbor, ospfAreaBlockInterfaceBlockNeighbor{
			Address: types.StringValue(itemTrimFields[0]),
		})
		if len(itemTrimFields) > 1 && itemTrimFields[1] == "eligible" {
			block.Neighbor[len(block.Neighbor)-1].Eligbile = types.BoolValue(true)
		}
	}

	return nil
}

func (block *ospfAreaBlockInterfaceBlockBfdLivenessDetection) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication algorithm "):
		block.AuthenticationAlgorithm = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "authentication key-chain "):
		block.AuthenticationKeyChain = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "authentication loose-check":
		block.AuthenticationLooseCheck = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "detection-time threshold "):
		block.DetectionTimeThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "full-neighbors-only":
		block.FullNeighborsOnly = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "holddown-interval "):
		block.HolddownInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-interval "):
		block.MinimumInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "minimum-receive-interval "):
		block.MinimumReceiveInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "multiplier "):
		block.Multiplier, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "no-adaptation":
		block.NoAdaptation = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "transmit-interval minimum-interval "):
		block.TransmitIntervalMinimumInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "transmit-interval threshold "):
		block.TransmitIntervalThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "version "):
		block.Version = types.StringValue(itemTrim)
	}

	return nil
}

func (block *ospfAreaBlockAreaRange) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "exact":
		block.Exact = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "override-metric "):
		block.OverrideMetric, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "restrict":
		block.Restrict = types.BoolValue(true)
	}

	return nil
}

func (block *ospfAreaBlockNssa) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "summaries":
		block.Summaries = types.BoolValue(true)
	case itemTrim == "no-summaries":
		block.NoSummaries = types.BoolValue(true)

	case balt.CutPrefixInString(&itemTrim, "area-range "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var areaRange ospfAreaBlockAreaRange
		block.AreaRange, areaRange = tfdata.ExtractBlockWithTFTypesString(
			block.AreaRange, "Range", itemTrimFields[0],
		)
		areaRange.Range = types.StringValue(itemTrimFields[0])
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")

		if err := areaRange.read(itemTrim); err != nil {
			return err
		}
		block.AreaRange = append(block.AreaRange, areaRange)
	case balt.CutPrefixInString(&itemTrim, "default-lsa"):
		if block.DefaultLsa == nil {
			block.DefaultLsa = &ospfAreaBlockNssaBlockDefaultLsa{}
		}

		switch {
		case balt.CutPrefixInString(&itemTrim, " default-metric "):
			block.DefaultLsa.DefaultMetric, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " metric-type "):
			block.DefaultLsa.MetricType, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case itemTrim == " type-7":
			block.DefaultLsa.Type7 = types.BoolValue(true)
		}
	}

	return nil
}

func (rscData *ospfAreaData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	if rscData.Version.ValueString() == "v3" {
		delPrefix += "protocols " + junos.OspfV3 + " "
	} else {
		delPrefix += "protocols " + junos.OspfV2 + " "
		if rscData.Realm.ValueString() != "" {
			return errors.New("realm can't set if version != v3")
		}
	}
	if v := rscData.Realm.ValueString(); v != "" {
		delPrefix += "realm " + v + " "
	}

	configSet := []string{
		delPrefix + "area " + rscData.AreaID.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
