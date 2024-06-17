package providerfwk

import (
	"context"
	"errors"
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
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &security{}
	_ resource.ResourceWithConfigure      = &security{}
	_ resource.ResourceWithValidateConfig = &security{}
	_ resource.ResourceWithImportState    = &security{}
	_ resource.ResourceWithUpgradeState   = &security{}
)

type security struct {
	client *junos.Client
}

func newSecurityResource() resource.Resource {
	return &security{}
}

func (rsc *security) typeName() string {
	return providerName + "_security"
}

func (rsc *security) junosName() string {
	return "security"
}

func (rsc *security) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *security) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *security) Configure(
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

func (rsc *security) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `security`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"clean_on_destroy": schema.BoolAttribute{
				Optional:    true,
				Description: "Clean supported lines when destroy this resource.",
			},
		},
		Blocks: map[string]schema.Block{
			"alg": schema.SingleNestedBlock{
				Description: "Declare `alg` configuration.",
				Attributes: map[string]schema.Attribute{
					"dns_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable dns alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ftp_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable ftp alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"h323_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable h323 alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"mgcp_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable mgcp alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"msrpc_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable msrpc alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"pptp_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable pptp alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"rsh_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable rsh alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"rtsp_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable rtsp alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"sccp_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable sccp alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"sip_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable sip alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"sql_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable sql alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"sunrpc_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable sunrpc alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"talk_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable talk alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"tftp_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable tftp alg.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"flow": schema.SingleNestedBlock{
				Description: "Declare `flow` configuration.",
				Attributes: map[string]schema.Attribute{
					"allow_dns_reply": schema.BoolAttribute{
						Optional:    true,
						Description: "Allow unmatched incoming DNS reply packet.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"allow_embedded_icmp": schema.BoolAttribute{
						Optional:    true,
						Description: "Allow embedded ICMP packets not matching a session to pass through.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"allow_reverse_ecmp": schema.BoolAttribute{
						Optional:    true,
						Description: "Allow reverse ECMP route lookup.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"enable_reroute_uniform_link_check_nat": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable reroute check with uniform link and NAT check.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"force_ip_reassembly": schema.BoolAttribute{
						Optional:    true,
						Description: "Force to reassemble ip fragments.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ipsec_performance_acceleration": schema.BoolAttribute{
						Optional:    true,
						Description: "Accelerate the IPSec traffic performance.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"mcast_buffer_enhance": schema.BoolAttribute{
						Optional:    true,
						Description: "Allow to hold more packets during multicast session creation.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"pending_sess_queue_length": schema.StringAttribute{
						Optional:    true,
						Description: "Maximum queued length per pending session.",
						Validators: []validator.String{
							stringvalidator.OneOf("high", "moderate", "normal"),
						},
					},
					"preserve_incoming_fragment_size": schema.BoolAttribute{
						Optional:    true,
						Description: "Preserve incoming fragment size for egress MTU.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"route_change_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Timeout value for route change to nonexistent route (6..1800 seconds).",
						Validators: []validator.Int64{
							int64validator.Between(6, 1800),
						},
					},
					"syn_flood_protection_mode": schema.StringAttribute{
						Optional:    true,
						Description: "TCP SYN flood protection mode.",
						Validators: []validator.String{
							stringvalidator.OneOf("syn-cookie", "syn-proxy"),
						},
					},
					"sync_icmp_session": schema.BoolAttribute{
						Optional:    true,
						Description: "Allow icmp sessions to sync to peer node.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"advanced_options": schema.SingleNestedBlock{
						Description: "Declare `flow advanced-options` configuration.",
						Attributes: map[string]schema.Attribute{
							"drop_matching_link_local_address": schema.BoolAttribute{
								Optional:    true,
								Description: "Drop matching link local address.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"drop_matching_reserved_ip_address": schema.BoolAttribute{
								Optional:    true,
								Description: "Drop matching reserved source IP address.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"reverse_route_packet_mode_vr": schema.BoolAttribute{
								Optional:    true,
								Description: "Allow reverse route lookup with packet mode vr.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"aging": schema.SingleNestedBlock{
						Description: "Declare `flow aging` configuration.",
						Attributes: map[string]schema.Attribute{
							"early_ageout": schema.Int64Attribute{
								Optional:    true,
								Description: "Delay before device declares session invalid.",
								Validators: []validator.Int64{
									int64validator.Between(1, 65535),
								},
							},
							"high_watermark": schema.Int64Attribute{
								Optional:    true,
								Description: "Percentage of session-table capacity at which aggressive aging-out starts.",
								Validators: []validator.Int64{
									int64validator.Between(0, 100),
								},
							},
							"low_watermark": schema.Int64Attribute{
								Optional:    true,
								Description: "Percentage of session-table capacity at which aggressive aging-out ends.",
								Validators: []validator.Int64{
									int64validator.Between(0, 100),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"ethernet_switching": schema.SingleNestedBlock{
						Description: "Declare `flow ethernet-switching` configuration.",
						Attributes: map[string]schema.Attribute{
							"block_non_ip_all": schema.BoolAttribute{
								Optional:    true,
								Description: "Block all non-IP and non-ARP traffic including broadcast/multicast.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"bypass_non_ip_unicast": schema.BoolAttribute{
								Optional:    true,
								Description: "Allow all non-IP (including unicast) traffic.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"bpdu_vlan_flooding": schema.BoolAttribute{
								Optional:    true,
								Description: "Set 802.1D BPDU flooding based on VLAN.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"no_packet_flooding": schema.SingleNestedBlock{
								Description: "Stop IP flooding, send ARP/ICMP to trigger MAC learning.",
								Attributes: map[string]schema.Attribute{
									"no_trace_route": schema.BoolAttribute{
										Optional:    true,
										Description: "Don't send ICMP to trigger MAC learning.",
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
					"tcp_mss": schema.SingleNestedBlock{
						Description: "Declare `flow tcp-mss` configuration.",
						Attributes: map[string]schema.Attribute{
							"all_tcp_mss": schema.Int64Attribute{
								Optional:    true,
								Description: "Enable MSS override for all packets with this value.",
								Validators: []validator.Int64{
									int64validator.Between(64, 65535),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"gre_in": schema.SingleNestedBlock{
								Description: "Enable MSS override for all GRE packets coming out of an IPSec tunnel.",
								Attributes: map[string]schema.Attribute{
									"mss": schema.Int64Attribute{
										Optional:    true,
										Description: "MSS Value.",
										Validators: []validator.Int64{
											int64validator.Between(64, 65535),
										},
									},
								},
								PlanModifiers: []planmodifier.Object{
									tfplanmodifier.BlockRemoveNull(),
								},
							},
							"gre_out": schema.SingleNestedBlock{
								Description: "Enable MSS override for all GRE packets entering an IPsec tunnel.",
								Attributes: map[string]schema.Attribute{
									"mss": schema.Int64Attribute{
										Optional:    true,
										Description: "MSS Value.",
										Validators: []validator.Int64{
											int64validator.Between(64, 65535),
										},
									},
								},
								PlanModifiers: []planmodifier.Object{
									tfplanmodifier.BlockRemoveNull(),
								},
							},
							"ipsec_vpn": schema.SingleNestedBlock{
								Description: "Enable MSS override for all packets entering IPSec tunnel.",
								Attributes: map[string]schema.Attribute{
									"mss": schema.Int64Attribute{
										Optional:    true,
										Description: "MSS Value.",
										Validators: []validator.Int64{
											int64validator.Between(64, 65535),
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
					"tcp_session": schema.SingleNestedBlock{
						Description: "Declare `flow tcp-session` configuration.",
						Attributes: map[string]schema.Attribute{
							"fin_invalidate_session": schema.BoolAttribute{
								Optional:    true,
								Description: "Immediately end session on receipt of fin (FIN) segment.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"maximum_window": schema.StringAttribute{
								Optional:    true,
								Description: "Maximum TCP proxy scaled receive window.",
								Validators: []validator.String{
									stringvalidator.OneOf("64K", "128K", "256K", "512K", "1M"),
								},
							},
							"no_sequence_check": schema.BoolAttribute{
								Optional:    true,
								Description: "Disable sequence-number checking.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_syn_check": schema.BoolAttribute{
								Optional:    true,
								Description: "Disable creation-time SYN-flag check.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_syn_check_in_tunnel": schema.BoolAttribute{
								Optional:    true,
								Description: "Disable creation-time SYN-flag check for tunnel packets.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"rst_invalidate_session": schema.BoolAttribute{
								Optional:    true,
								Description: "Immediately end session on receipt of reset (RST) segment.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"rst_sequence_check": schema.BoolAttribute{
								Optional:    true,
								Description: "Check sequence number in reset (RST) segment.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"strict_syn_check": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable strict syn check.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"tcp_initial_timeout": schema.Int64Attribute{
								Optional:    true,
								Description: "Timeout for TCP session when initialization fails (4..300 seconds).",
								Validators: []validator.Int64{
									int64validator.Between(4, 300),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"time_wait_state": schema.SingleNestedBlock{
								Description: "Declare session timeout value in time-wait state.",
								Attributes: map[string]schema.Attribute{
									"apply_to_half_close_state": schema.BoolAttribute{
										Optional:    true,
										Description: "Apply time-wait-state timeout to half-close state.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"session_ageout": schema.BoolAttribute{
										Optional:    true,
										Description: "Allow session to ageout using service based timeout values.",
										Validators: []validator.Bool{
											tfvalidator.BoolTrue(),
										},
									},
									"session_timeout": schema.Int64Attribute{
										Optional:    true,
										Description: "Configure session timeout value for time-wait state (2..600 seconds).",
										Validators: []validator.Int64{
											int64validator.Between(2, 600),
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
			"forwarding_options": schema.SingleNestedBlock{
				Description: "Declare `forwarding-options` configuration.",
				Attributes: map[string]schema.Attribute{
					"inet6_mode": schema.StringAttribute{
						Optional:    true,
						Description: "Forwarding mode for inet6 family.",
						Validators: []validator.String{
							stringvalidator.OneOf("drop", "flow-based", "packet-based"),
						},
					},
					"iso_mode_packet_based": schema.BoolAttribute{
						Optional:    true,
						Description: "Forwarding mode packet-based for iso family.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"mpls_mode": schema.StringAttribute{
						Optional:    true,
						Description: "Forwarding mode for mpls family.",
						Validators: []validator.String{
							stringvalidator.OneOf("flow-based", "packet-based"),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"forwarding_process": schema.SingleNestedBlock{
				Description: "Declare `forwarding-process` configuration.",
				Attributes: map[string]schema.Attribute{
					"enhanced_services_mode": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable enhanced application services mode.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"idp_security_package": schema.SingleNestedBlock{
				Description: "Declare `idp security-package` configuration.",
				Attributes: map[string]schema.Attribute{
					"automatic_enable": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable scheduled download and update.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"automatic_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Automatic interval (1..336 hours).",
						Validators: []validator.Int64{
							int64validator.Between(1, 336),
						},
					},
					"automatic_start_time": schema.StringAttribute{
						Optional:    true,
						Description: "Automatic start time (YYYY-MM-DD.HH:MM:SS).",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^\d{4}\-\d\d?\-\d\d?\.\d{2}:\d{2}:\d{2}$`),
								"must be in the format 'YYYY-MM-DD.HH:MM:SS'",
							),
						},
					},
					"install_ignore_version_check": schema.BoolAttribute{
						Optional:    true,
						Description: "Skip version check when attack database gets installed.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"proxy_profile": schema.StringAttribute{
						Optional:    true,
						Description: "Proxy profile of security package download.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 64),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"source_address": schema.StringAttribute{
						Optional:    true,
						Description: "Source address to be used for sending download request.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"url": schema.StringAttribute{
						Optional:    true,
						Description: "URL of Security package download.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"idp_sensor_configuration": schema.SingleNestedBlock{
				Description: "Declare `idp sensor-configuration` configuration.",
				Attributes: map[string]schema.Attribute{
					"log_cache_size": schema.Int64Attribute{
						Optional:    true,
						Description: "Log cache size.",
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"security_configuration_protection_mode": schema.StringAttribute{
						Optional:    true,
						Description: "Enable security protection mode.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"datacenter",
								"datacenter-full",
								"perimeter",
								"perimeter-full",
							),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"log_suppression": schema.SingleNestedBlock{
						Description: "Enable `log suppression`.",
						Attributes: map[string]schema.Attribute{
							"disable": schema.BoolAttribute{
								Optional:    true,
								Description: "Disable log suppression.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"include_destination_address": schema.BoolAttribute{
								Optional:    true,
								Description: "Include destination address while performing a log suppression.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_include_destination_address": schema.BoolAttribute{
								Optional:    true,
								Description: "Don't include destination address while performing a log suppression.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"max_logs_operate": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum logs can be operate on.",
								Validators: []validator.Int64{
									int64validator.Between(256, 65536),
								},
							},
							"max_time_report": schema.Int64Attribute{
								Optional:    true,
								Description: "Time after suppressed logs will be reported (1..60).",
								Validators: []validator.Int64{
									int64validator.Between(1, 60),
								},
							},
							"start_log": schema.Int64Attribute{
								Optional:    true,
								Description: "Suppression start log (1..128).",
								Validators: []validator.Int64{
									int64validator.Between(1, 128),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"packet_log": schema.SingleNestedBlock{
						Description: "Declare `packet-log` configuration.",
						Attributes: map[string]schema.Attribute{
							"source_address": schema.StringAttribute{
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "Source IP address used to transport packetlog to a host.",
								Validators: []validator.String{
									tfvalidator.StringIPAddress(),
								},
							},
							"host_address": schema.StringAttribute{
								Optional:    true,
								Description: "Destination host to send packetlog to.",
								Validators: []validator.String{
									tfvalidator.StringIPAddress(),
								},
							},
							"host_port": schema.Int64Attribute{
								Optional:    true,
								Description: "Destination UDP port number.",
								Validators: []validator.Int64{
									int64validator.Between(1, 65535),
								},
							},
							"max_sessions": schema.Int64Attribute{
								Optional:    true,
								Description: "Max num of sessions in unit(%).",
								Validators: []validator.Int64{
									int64validator.Between(1, 100),
								},
							},
							"threshold_logging_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Interval of logs for max limit session/memory reached in minutes (1..60).",
								Validators: []validator.Int64{
									int64validator.Between(1, 60),
								},
							},
							"total_memory": schema.Int64Attribute{
								Optional:    true,
								Description: "Total memory unit(%).",
								Validators: []validator.Int64{
									int64validator.Between(1, 100),
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
			"ike_traceoptions": schema.SingleNestedBlock{
				Description: "Declare `ike traceoptions` configuration.",
				Attributes: map[string]schema.Attribute{
					"flag": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Tracing parameters for IKE.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							),
						},
					},
					"no_remote_trace": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable remote tracing.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"rate_limit": schema.Int64Attribute{
						Optional:    true,
						Description: "Limit the incoming rate of trace messages.",
						Validators: []validator.Int64{
							int64validator.Between(0, 4294967295),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"file": schema.SingleNestedBlock{
						Description: "Declare `file` configuration.",
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Optional:    true,
								Description: "Name of file in which to write trace information.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
									tfvalidator.StringSpaceExclusion(),
									tfvalidator.StringRuneExclusion('/', '%'),
								},
							},
							"files": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum number of trace files (2..1000).",
								Validators: []validator.Int64{
									int64validator.Between(2, 1000),
								},
							},
							"match": schema.StringAttribute{
								Optional:    true,
								Description: "Regular expression for lines to be logged.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"size": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum trace file size.",
								Validators: []validator.Int64{
									int64validator.Between(10240, 1073741824),
								},
							},
							"world_readable": schema.BoolAttribute{
								Optional:    true,
								Description: "Allow any user to read the log file.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_world_readable": schema.BoolAttribute{
								Optional:    true,
								Description: "Don't allow any user to read the log file.",
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
			"log": schema.SingleNestedBlock{
				Description: "Declare `log` configuration.",
				Attributes: map[string]schema.Attribute{
					"disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable security logging for the device.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"event_rate": schema.Int64Attribute{
						Optional:    true,
						Description: "Control plane event rate (0..1500 logs per second).",
						Validators: []validator.Int64{
							int64validator.Between(0, 1500),
						},
					},
					"facility_override": schema.StringAttribute{
						Optional:    true,
						Description: "Alternate facility for logging to remote host.",
						Validators: []validator.String{
							stringvalidator.OneOf(junos.SyslogFacilities()...),
						},
					},
					"format": schema.StringAttribute{
						Optional:    true,
						Description: "Set security log format for the device.",
						Validators: []validator.String{
							stringvalidator.OneOf("binary", "sd-syslog", "syslog"),
						},
					},
					"max_database_record": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum records in database.",
						Validators: []validator.Int64{
							int64validator.Between(0, 1000000),
						},
					},
					"mode": schema.StringAttribute{
						Optional:    true,
						Description: "Controls how security logs are processed and exported.",
						Validators: []validator.String{
							stringvalidator.OneOf("event", "stream"),
						},
					},
					"rate_cap": schema.Int64Attribute{
						Optional:    true,
						Description: "Data plane event rate (0..5000 logs per second).",
						Validators: []validator.Int64{
							int64validator.Between(0, 5000),
						},
					},
					"report": schema.BoolAttribute{
						Optional:    true,
						Description: "Set security log report settings.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"source_address": schema.StringAttribute{
						Optional:    true,
						Description: "Source ip address used when exporting security logs.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"source_interface": schema.StringAttribute{
						Optional:    true,
						Description: "Source interface used when exporting security logs.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
							tfvalidator.String1DotCount(),
						},
					},
					"utc_timestamp": schema.BoolAttribute{
						Optional:    true,
						Description: "Use UTC time for security log timestamps.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"file": schema.SingleNestedBlock{
						Description: "Declare `security log file` configuration.",
						Attributes: map[string]schema.Attribute{
							"files": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum number of binary log files (2..10).",
								Validators: []validator.Int64{
									int64validator.Between(2, 10),
								},
							},
							"name": schema.StringAttribute{
								Optional:    true,
								Description: "Name of binary log file.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
									tfvalidator.StringSpaceExclusion(),
									tfvalidator.StringRuneExclusion('/', '%'),
								},
							},
							"path": schema.StringAttribute{
								Optional:    true,
								Description: "Path to binary log files.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"size": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum size of binary log file in megabytes (1..10).",
								Validators: []validator.Int64{
									int64validator.Between(1, 10),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"transport": schema.SingleNestedBlock{
						Description: "Declare `security log transport` configuration.",
						Attributes: map[string]schema.Attribute{
							"protocol": schema.StringAttribute{
								Optional:    true,
								Description: "Set security log transport protocol for the device.",
								Validators: []validator.String{
									stringvalidator.OneOf("tcp", "tls", "udp"),
								},
							},
							"tcp_connections": schema.Int64Attribute{
								Optional:    true,
								Description: "Set tcp connection number per-stream (1..5).",
								Validators: []validator.Int64{
									int64validator.Between(1, 5),
								},
							},
							"tls_profile": schema.StringAttribute{
								Optional:    true,
								Description: "TLS profile.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
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
			"nat_source": schema.SingleNestedBlock{
				Description: "Declare `nat source` configuration.",
				Attributes: map[string]schema.Attribute{
					"address_persistent": schema.BoolAttribute{
						Optional:    true,
						Description: "Allow source address to maintain same translation.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"interface_port_overloading_factor": schema.Int64Attribute{
						Optional:    true,
						Description: "Port overloading factor for interface NAT.",
						Validators: []validator.Int64{
							int64validator.Between(0, 65535),
						},
					},
					"interface_port_overloading_off": schema.BoolAttribute{
						Optional:    true,
						Description: "Turn off interface port over-loading.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"pool_default_port_range": schema.Int64Attribute{
						Optional:    true,
						Description: "Configure Source NAT default port range lower limit.",
						Validators: []validator.Int64{
							int64validator.Between(1024, 63487),
						},
					},
					"pool_default_port_range_to": schema.Int64Attribute{
						Optional:    true,
						Description: "Configure Source NAT default port range upper limit.",
						Validators: []validator.Int64{
							int64validator.Between(1024, 63487),
						},
					},
					"pool_default_twin_port_range": schema.Int64Attribute{
						Optional:    true,
						Description: "Configure Source NAT default twin port range lower limit.",
						Validators: []validator.Int64{
							int64validator.Between(63488, 65535),
						},
					},
					"pool_default_twin_port_range_to": schema.Int64Attribute{
						Optional:    true,
						Description: "Configure Source NAT default twin port range upper limit.",
						Validators: []validator.Int64{
							int64validator.Between(63488, 65535),
						},
					},
					"pool_utilization_alarm_clear_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Clear threshold for pool utilization alarm (40..100).",
						Validators: []validator.Int64{
							int64validator.Between(40, 100),
						},
					},
					"pool_utilization_alarm_raise_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Raise threshold for pool utilization alarm (50..100).",
						Validators: []validator.Int64{
							int64validator.Between(50, 100),
						},
					},
					"port_randomization_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable Source NAT port randomization.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"session_drop_hold_down": schema.Int64Attribute{
						Optional:    true,
						Description: "Session drop hold down time (30..28800).",
						Validators: []validator.Int64{
							int64validator.Between(30, 28800),
						},
					},
					"session_persistence_scan": schema.BoolAttribute{
						Optional:    true,
						Description: "Allow source to maintain session when session scan.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"policies": schema.SingleNestedBlock{
				Description: "Declare `policies` configuration.",
				Attributes: map[string]schema.Attribute{
					"policy_rematch": schema.BoolAttribute{
						Optional:    true,
						Description: "Can be specified to allow session to remain open when an associated security policy is modified.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"policy_rematch_extensive": schema.BoolAttribute{
						Optional: true,
						Description: "Can be specified to allow session to remain open " +
							"when an associated security policy is modified, renamed, deactivated, or deleted.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"user_identification_auth_source": schema.SingleNestedBlock{
				Description: "Declare `user-identification authentication-source` configuration.",
				Attributes: map[string]schema.Attribute{
					"ad_auth_priority": schema.Int64Attribute{
						Optional:    true,
						Description: "Active-directory-authentication-table priority.",
						Validators: []validator.Int64{
							int64validator.Between(0, 65535),
						},
					},
					"aruba_clearpass_priority": schema.Int64Attribute{
						Optional:    true,
						Description: "ClearPass-authentication-table priority.",
						Validators: []validator.Int64{
							int64validator.Between(0, 65535),
						},
					},
					"firewall_auth_priority": schema.Int64Attribute{
						Optional:    true,
						Description: "Firewall-authentication priority.",
						Validators: []validator.Int64{
							int64validator.Between(0, 65535),
						},
					},
					"local_auth_priority": schema.Int64Attribute{
						Optional:    true,
						Description: "Local-authentication-table priority.",
						Validators: []validator.Int64{
							int64validator.Between(0, 65535),
						},
					},
					"unified_access_control_priority": schema.Int64Attribute{
						Optional:    true,
						Description: "Unified-access-control priority.",
						Validators: []validator.Int64{
							int64validator.Between(0, 65535),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"utm": schema.SingleNestedBlock{
				Description: "Declare `utm` configuration.",
				Attributes: map[string]schema.Attribute{
					"feature_profile_web_filtering_type": schema.StringAttribute{
						Optional:    true,
						Description: "Configuring feature-profile web-filtering type.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"juniper-enhanced",
								"juniper-local",
								"web-filtering-none",
								"websense-redirect",
							),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"feature_profile_web_filtering_juniper_enhanced_server": schema.SingleNestedBlock{
						Description: "Declare `utm feature-profile web-filtering juniper-enhanced server` configuration.",
						Attributes: map[string]schema.Attribute{
							"host": schema.StringAttribute{
								Optional:    true,
								Description: "Server host IP address or string host name.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"port": schema.Int64Attribute{
								Optional:    true,
								Description: "Server port.",
								Validators: []validator.Int64{
									int64validator.Between(1, 65535),
								},
							},
							"proxy_profile": schema.StringAttribute{
								Optional:    true,
								Description: "Proxy profile.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 64),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"routing_instance": schema.StringAttribute{
								Optional:    true,
								Description: "Routing instance name.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringFormat(tfvalidator.DefaultFormat),
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

type securityData struct {
	ID                           types.String                               `tfsdk:"id"`
	CleanOnDestroy               types.Bool                                 `tfsdk:"clean_on_destroy"`
	Alg                          *securityBlockAlg                          `tfsdk:"alg"`
	Flow                         *securityBlockFlow                         `tfsdk:"flow"`
	ForwardingOptions            *securityBlockForwardingOptions            `tfsdk:"forwarding_options"`
	ForwardingProcess            *securityBlockForwardingProcess            `tfsdk:"forwarding_process"`
	IdpSecurityPackage           *securityBlockIdpSecurityPackage           `tfsdk:"idp_security_package"`
	IdpSensorConfiguration       *securityBlockIdpSensorConfiguration       `tfsdk:"idp_sensor_configuration"`
	IkeTraceoptions              *securityBlockIkeTraceoptions              `tfsdk:"ike_traceoptions"`
	Log                          *securityBlockLog                          `tfsdk:"log"`
	NatSource                    *securityBlockNatSource                    `tfsdk:"nat_source"`
	Policies                     *securityBlockPolicies                     `tfsdk:"policies"`
	UserIdentificationAuthSource *securityBlockUserIdentificationAuthSource `tfsdk:"user_identification_auth_source"`
	Utm                          *securityBlockUtm                          `tfsdk:"utm"`
}

type securityConfig struct {
	ID                           types.String                               `tfsdk:"id"`
	CleanOnDestroy               types.Bool                                 `tfsdk:"clean_on_destroy"`
	Alg                          *securityBlockAlg                          `tfsdk:"alg"`
	Flow                         *securityBlockFlow                         `tfsdk:"flow"`
	ForwardingOptions            *securityBlockForwardingOptions            `tfsdk:"forwarding_options"`
	ForwardingProcess            *securityBlockForwardingProcess            `tfsdk:"forwarding_process"`
	IdpSecurityPackage           *securityBlockIdpSecurityPackage           `tfsdk:"idp_security_package"`
	IdpSensorConfiguration       *securityBlockIdpSensorConfiguration       `tfsdk:"idp_sensor_configuration"`
	IkeTraceoptions              *securityBlockIkeTraceoptionsConfig        `tfsdk:"ike_traceoptions"`
	Log                          *securityBlockLog                          `tfsdk:"log"`
	NatSource                    *securityBlockNatSource                    `tfsdk:"nat_source"`
	Policies                     *securityBlockPolicies                     `tfsdk:"policies"`
	UserIdentificationAuthSource *securityBlockUserIdentificationAuthSource `tfsdk:"user_identification_auth_source"`
	Utm                          *securityBlockUtm                          `tfsdk:"utm"`
}

type securityBlockAlg struct {
	DNSDisable    types.Bool `tfsdk:"dns_disable"`
	FtpDisable    types.Bool `tfsdk:"ftp_disable"`
	H323Disable   types.Bool `tfsdk:"h323_disable"`
	MgcpDisable   types.Bool `tfsdk:"mgcp_disable"`
	MsrpcDisable  types.Bool `tfsdk:"msrpc_disable"`
	PptpDisable   types.Bool `tfsdk:"pptp_disable"`
	RshDisable    types.Bool `tfsdk:"rsh_disable"`
	RtspDisable   types.Bool `tfsdk:"rtsp_disable"`
	SccpDisable   types.Bool `tfsdk:"sccp_disable"`
	SIPDisable    types.Bool `tfsdk:"sip_disable"`
	SQLDisable    types.Bool `tfsdk:"sql_disable"`
	SunrpcDisable types.Bool `tfsdk:"sunrpc_disable"`
	TalkDisable   types.Bool `tfsdk:"talk_disable"`
	TftpDisable   types.Bool `tfsdk:"tftp_disable"`
}

func (block *securityBlockAlg) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:lll
type securityBlockFlow struct {
	AllowDNSReply                    types.Bool                               `tfsdk:"allow_dns_reply"`
	AllowEmbeddedIcmp                types.Bool                               `tfsdk:"allow_embedded_icmp"`
	AllowReverseEcmp                 types.Bool                               `tfsdk:"allow_reverse_ecmp"`
	EnableRerouteUniformLinkCheckNat types.Bool                               `tfsdk:"enable_reroute_uniform_link_check_nat"`
	ForceIPReassembly                types.Bool                               `tfsdk:"force_ip_reassembly"`
	IpsecPerformanceAcceleration     types.Bool                               `tfsdk:"ipsec_performance_acceleration"`
	McastBufferEnhance               types.Bool                               `tfsdk:"mcast_buffer_enhance"`
	PendingSessQueueLength           types.String                             `tfsdk:"pending_sess_queue_length"`
	PreserveIncomingFragmentSize     types.Bool                               `tfsdk:"preserve_incoming_fragment_size"`
	RouteChangeTimeout               types.Int64                              `tfsdk:"route_change_timeout"`
	SynFloodProtectionMode           types.String                             `tfsdk:"syn_flood_protection_mode"`
	SyncIcmpSession                  types.Bool                               `tfsdk:"sync_icmp_session"`
	AdvancedOptions                  *securityBlockFlowBlockAdvancedOptions   `tfsdk:"advanced_options"`
	Aging                            *securityBlockFlowBlockAging             `tfsdk:"aging"`
	EthernetSwitching                *securityBlockFlowBlockEthernetSwitching `tfsdk:"ethernet_switching"`
	TCPMss                           *securityBlockFlowBlockTCPMss            `tfsdk:"tcp_mss"`
	TCPSession                       *securityBlockFlowBlockTCPSession        `tfsdk:"tcp_session"`
}

func (block *securityBlockFlow) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockFlowBlockAdvancedOptions struct {
	DropMatchingLinkLocalAddress  types.Bool `tfsdk:"drop_matching_link_local_address"`
	DropMatchingReservedIPAddress types.Bool `tfsdk:"drop_matching_reserved_ip_address"`
	ReverseRoutePacketModeVR      types.Bool `tfsdk:"reverse_route_packet_mode_vr"`
}

func (block *securityBlockFlowBlockAdvancedOptions) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockFlowBlockAging struct {
	EarlyAgeout   types.Int64 `tfsdk:"early_ageout"`
	HighWatermark types.Int64 `tfsdk:"high_watermark"`
	LowWatermark  types.Int64 `tfsdk:"low_watermark"`
}

func (block *securityBlockFlowBlockAging) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockFlowBlockEthernetSwitching struct {
	BlockNonIPAll      types.Bool `tfsdk:"block_non_ip_all"`
	BypassNonIPUnicast types.Bool `tfsdk:"bypass_non_ip_unicast"`
	BpduVlanFlooding   types.Bool `tfsdk:"bpdu_vlan_flooding"`
	NoPacketFlooding   *struct {
		NoTraceRoute types.Bool `tfsdk:"no_trace_route"`
	} `tfsdk:"no_packet_flooding"`
}

func (block *securityBlockFlowBlockEthernetSwitching) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockFlowBlockTCPMss struct {
	AllTCPMss types.Int64 `tfsdk:"all_tcp_mss"`
	GreIn     *struct {
		Mss types.Int64 `tfsdk:"mss"`
	} `tfsdk:"gre_in"`
	GreOut *struct {
		Mss types.Int64 `tfsdk:"mss"`
	} `tfsdk:"gre_out"`
	IpsecVpn *struct {
		Mss types.Int64 `tfsdk:"mss"`
	} `tfsdk:"ipsec_vpn"`
}

func (block *securityBlockFlowBlockTCPMss) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockFlowBlockTCPSession struct {
	FinInvalidateSession types.Bool                                          `tfsdk:"fin_invalidate_session"`
	MaximumWindow        types.String                                        `tfsdk:"maximum_window"`
	NoSequenceCheck      types.Bool                                          `tfsdk:"no_sequence_check"`
	NoSynCheck           types.Bool                                          `tfsdk:"no_syn_check"`
	NoSynCheckInTunnel   types.Bool                                          `tfsdk:"no_syn_check_in_tunnel"`
	RstInvalidateSession types.Bool                                          `tfsdk:"rst_invalidate_session"`
	RstSequenceCheck     types.Bool                                          `tfsdk:"rst_sequence_check"`
	StrictSynCheck       types.Bool                                          `tfsdk:"strict_syn_check"`
	TCPInitialTimeout    types.Int64                                         `tfsdk:"tcp_initial_timeout"`
	TimeWaitState        *securityBlockFlowBlockTCPSessionBlockTimeWaitState `tfsdk:"time_wait_state"`
}

func (block *securityBlockFlowBlockTCPSession) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockFlowBlockTCPSessionBlockTimeWaitState struct {
	ApplyToHalfCloseState types.Bool  `tfsdk:"apply_to_half_close_state"`
	SessionAgeout         types.Bool  `tfsdk:"session_ageout"`
	SessionTimeout        types.Int64 `tfsdk:"session_timeout"`
}

type securityBlockForwardingOptions struct {
	Inet6Mode          types.String `tfsdk:"inet6_mode"`
	IsoModePacketBased types.Bool   `tfsdk:"iso_mode_packet_based"`
	MplsMode           types.String `tfsdk:"mpls_mode"`
}

func (block *securityBlockForwardingOptions) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockForwardingProcess struct {
	EnhancedServicesMode types.Bool `tfsdk:"enhanced_services_mode"`
}

func (block *securityBlockForwardingProcess) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockIdpSecurityPackage struct {
	AutomaticEnable           types.Bool   `tfsdk:"automatic_enable"`
	AutomaticInterval         types.Int64  `tfsdk:"automatic_interval"`
	AutomaticStartTime        types.String `tfsdk:"automatic_start_time"`
	InstallIgnoreVersionCheck types.Bool   `tfsdk:"install_ignore_version_check"`
	ProxyProfile              types.String `tfsdk:"proxy_profile"`
	SourceAddress             types.String `tfsdk:"source_address"`
	URL                       types.String `tfsdk:"url"`
}

func (block *securityBlockIdpSecurityPackage) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:lll
type securityBlockIdpSensorConfiguration struct {
	LogCacheSize                        types.Int64                                             `tfsdk:"log_cache_size"`
	SecurityConfigurationProtectionMode types.String                                            `tfsdk:"security_configuration_protection_mode"`
	LogSuppression                      *securityBlockIdpSensorConfigurationBlockLogSuppression `tfsdk:"log_suppression"`
	PacketLog                           *securityBlockIdpSensorConfigurationBlockPacketLog      `tfsdk:"packet_log"`
}

func (block *securityBlockIdpSensorConfiguration) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockIdpSensorConfigurationBlockLogSuppression struct {
	Disable                     types.Bool  `tfsdk:"disable"`
	IncludeDestinationAddress   types.Bool  `tfsdk:"include_destination_address"`
	NoIncludeDestinationAddress types.Bool  `tfsdk:"no_include_destination_address"`
	MaxLogsOperate              types.Int64 `tfsdk:"max_logs_operate"`
	MaxTimeReport               types.Int64 `tfsdk:"max_time_report"`
	StartLog                    types.Int64 `tfsdk:"start_log"`
}

type securityBlockIdpSensorConfigurationBlockPacketLog struct {
	SourceAddress            types.String `tfsdk:"source_address"`
	HostAddress              types.String `tfsdk:"host_address"`
	HostPort                 types.Int64  `tfsdk:"host_port"`
	MaxSessions              types.Int64  `tfsdk:"max_sessions"`
	ThresholdLoggingInterval types.Int64  `tfsdk:"threshold_logging_interval"`
	TotalMemory              types.Int64  `tfsdk:"total_memory"`
}

type securityBlockIkeTraceoptions struct {
	Flag          []types.String                         `tfsdk:"flag"`
	NoRemoteTrace types.Bool                             `tfsdk:"no_remote_trace"`
	RateLimit     types.Int64                            `tfsdk:"rate_limit"`
	File          *securityBlockIkeTraceoptionsBlockFile `tfsdk:"file"`
}

func (block *securityBlockIkeTraceoptions) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockIkeTraceoptionsConfig struct {
	Flag          types.Set                              `tfsdk:"flag"`
	NoRemoteTrace types.Bool                             `tfsdk:"no_remote_trace"`
	RateLimit     types.Int64                            `tfsdk:"rate_limit"`
	File          *securityBlockIkeTraceoptionsBlockFile `tfsdk:"file"`
}

func (block *securityBlockIkeTraceoptionsConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockIkeTraceoptionsBlockFile struct {
	Name            types.String `tfsdk:"name"`
	Files           types.Int64  `tfsdk:"files"`
	Match           types.String `tfsdk:"match"`
	Size            types.Int64  `tfsdk:"size"`
	WorldReadable   types.Bool   `tfsdk:"world_readable"`
	NoWorldReadable types.Bool   `tfsdk:"no_world_readable"`
}

func (block *securityBlockIkeTraceoptionsBlockFile) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockLog struct {
	Disable           types.Bool                      `tfsdk:"disable"`
	EventRate         types.Int64                     `tfsdk:"event_rate"`
	FacilityOverride  types.String                    `tfsdk:"facility_override"`
	Format            types.String                    `tfsdk:"format"`
	MaxDatabaseRecord types.Int64                     `tfsdk:"max_database_record"`
	Mode              types.String                    `tfsdk:"mode"`
	RateCap           types.Int64                     `tfsdk:"rate_cap"`
	Report            types.Bool                      `tfsdk:"report"`
	SourceAddress     types.String                    `tfsdk:"source_address"`
	SourceInterface   types.String                    `tfsdk:"source_interface"`
	UtcTimestamp      types.Bool                      `tfsdk:"utc_timestamp"`
	File              *securityBlockLogBlockFile      `tfsdk:"file"`
	Transport         *securityBlockLogBlockTransport `tfsdk:"transport"`
}

func (block *securityBlockLog) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockLogBlockFile struct {
	Files types.Int64  `tfsdk:"files"`
	Name  types.String `tfsdk:"name"`
	Path  types.String `tfsdk:"path"`
	Size  types.Int64  `tfsdk:"size"`
}

func (block *securityBlockLogBlockFile) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockLogBlockTransport struct {
	Protocol       types.String `tfsdk:"protocol"`
	TCPConnections types.Int64  `tfsdk:"tcp_connections"`
	TLSProfile     types.String `tfsdk:"tls_profile"`
}

type securityBlockNatSource struct {
	AddressPersistent                  types.Bool  `tfsdk:"address_persistent"`
	InterfacePortOverloadingFactor     types.Int64 `tfsdk:"interface_port_overloading_factor"`
	InterfacePortOverloadingOff        types.Bool  `tfsdk:"interface_port_overloading_off"`
	PoolDefaultPortRange               types.Int64 `tfsdk:"pool_default_port_range"`
	PoolDefaultPortRangeTo             types.Int64 `tfsdk:"pool_default_port_range_to"`
	PoolDefaultTwinPortRange           types.Int64 `tfsdk:"pool_default_twin_port_range"`
	PoolDefaultTwinPortRangeTo         types.Int64 `tfsdk:"pool_default_twin_port_range_to"`
	PoolUtilizationAlarmClearThreshold types.Int64 `tfsdk:"pool_utilization_alarm_clear_threshold"`
	PoolUtilizationAlarmRaiseThreshold types.Int64 `tfsdk:"pool_utilization_alarm_raise_threshold"`
	PortRandomizationDisable           types.Bool  `tfsdk:"port_randomization_disable"`
	SessionDropHoldDown                types.Int64 `tfsdk:"session_drop_hold_down"`
	SessionPersistenceScan             types.Bool  `tfsdk:"session_persistence_scan"`
}

func (block *securityBlockNatSource) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockPolicies struct {
	PolicyRematch          types.Bool `tfsdk:"policy_rematch"`
	PolicyRematchExtensive types.Bool `tfsdk:"policy_rematch_extensive"`
}

func (block *securityBlockPolicies) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockUserIdentificationAuthSource struct {
	ADAuthPriority               types.Int64 `tfsdk:"ad_auth_priority"`
	ArubaClearpassPriority       types.Int64 `tfsdk:"aruba_clearpass_priority"`
	FirewallAuthPriority         types.Int64 `tfsdk:"firewall_auth_priority"`
	LocalAuthPriority            types.Int64 `tfsdk:"local_auth_priority"`
	UnifiedAccessControlPriority types.Int64 `tfsdk:"unified_access_control_priority"`
}

func (block *securityBlockUserIdentificationAuthSource) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:lll
type securityBlockUtm struct {
	FeatureProfileWebFilteringType                  types.String                                                          `tfsdk:"feature_profile_web_filtering_type"`
	FeatureProfileWebFilteringJuniperEnhancedServer *securityBlockUtmBlockFeatureProfileWebFilteringJuniperEnhancedServer `tfsdk:"feature_profile_web_filtering_juniper_enhanced_server"`
}

func (block *securityBlockUtm) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityBlockUtmBlockFeatureProfileWebFilteringJuniperEnhancedServer struct {
	Host            types.String `tfsdk:"host"`
	Port            types.Int64  `tfsdk:"port"`
	ProxyProfile    types.String `tfsdk:"proxy_profile"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
}

func (rsc *security) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Alg != nil {
		if config.Alg.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("alg").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"alg block is empty",
			)
		}
	}

	if config.Flow != nil {
		if config.Flow.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("flow").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"flow block is empty",
			)
		}
		if config.Flow.AdvancedOptions != nil {
			if config.Flow.AdvancedOptions.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("flow").AtName("advanced_options").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"advanced_options block is empty in flow block",
				)
			}
		}
		if config.Flow.Aging != nil {
			if config.Flow.Aging.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("flow").AtName("aging").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"aging block is empty in flow block",
				)
			}
		}
		if config.Flow.EthernetSwitching != nil {
			if config.Flow.EthernetSwitching.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("flow").AtName("ethernet_switching").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"ethernet_switching block is empty in flow block",
				)
			}
			if !config.Flow.EthernetSwitching.BlockNonIPAll.IsNull() &&
				!config.Flow.EthernetSwitching.BlockNonIPAll.IsUnknown() &&
				!config.Flow.EthernetSwitching.BypassNonIPUnicast.IsNull() &&
				!config.Flow.EthernetSwitching.BypassNonIPUnicast.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("flow").AtName("ethernet_switching").AtName("block_non_ip_all"),
					tfdiag.ConflictConfigErrSummary,
					"block_non_ip_all and bypass_non_ip_unicast can't be true in same time "+
						"in ethernet_switching block in flow block",
				)
			}
		}
		if config.Flow.TCPMss != nil {
			if config.Flow.TCPMss.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("flow").AtName("tcp_mss").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"tcp_mss block is empty in flow block",
				)
			}
		}
		if config.Flow.TCPSession != nil {
			if config.Flow.TCPSession.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("flow").AtName("tcp_session").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"tcp_session block is empty in flow block",
				)
			}
			if !config.Flow.TCPSession.StrictSynCheck.IsNull() && !config.Flow.TCPSession.StrictSynCheck.IsUnknown() {
				if !config.Flow.TCPSession.NoSynCheck.IsNull() && !config.Flow.TCPSession.NoSynCheck.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("flow").AtName("tcp_session").AtName("no_syn_check"),
						tfdiag.ConflictConfigErrSummary,
						"no_syn_check and strict_syn_check can't be true in same time "+
							"in tcp_session block in flow block",
					)
				}
				if !config.Flow.TCPSession.NoSynCheckInTunnel.IsNull() && !config.Flow.TCPSession.NoSynCheckInTunnel.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("flow").AtName("tcp_session").AtName("no_syn_check_in_tunnel"),
						tfdiag.ConflictConfigErrSummary,
						"no_syn_check_in_tunnel and strict_syn_check can't be true in same time "+
							"in tcp_session block in flow block",
					)
				}
			}
			if config.Flow.TCPSession.TimeWaitState != nil {
				if !config.Flow.TCPSession.TimeWaitState.SessionAgeout.IsNull() &&
					!config.Flow.TCPSession.TimeWaitState.SessionAgeout.IsUnknown() &&
					!config.Flow.TCPSession.TimeWaitState.SessionTimeout.IsNull() &&
					!config.Flow.TCPSession.TimeWaitState.SessionTimeout.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("flow").AtName("tcp_session").AtName("time_wait_state").AtName("session_ageout"),
						tfdiag.ConflictConfigErrSummary,
						"session_ageout and session_timeout can't be set in same time "+
							"in time_wait_state block in tcp_session block in flow block",
					)
				}
			}
		}
	}

	if config.ForwardingOptions != nil {
		if config.ForwardingOptions.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("forwarding_options").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"forwarding_options block is empty",
			)
		}
	}

	if config.ForwardingProcess != nil {
		if config.ForwardingProcess.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("forwarding_process").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"forwarding_process block is empty",
			)
		}
	}

	if config.IdpSecurityPackage != nil {
		if config.IdpSecurityPackage.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("idp_security_package").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"idp_security_package block is empty",
			)
		}
	}

	if config.IdpSensorConfiguration != nil {
		if config.IdpSensorConfiguration.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("idp_sensor_configuration").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"idp_sensor_configuration block is empty",
			)
		}
		if config.IdpSensorConfiguration.LogSuppression != nil {
			if !config.IdpSensorConfiguration.LogSuppression.IncludeDestinationAddress.IsNull() &&
				!config.IdpSensorConfiguration.LogSuppression.IncludeDestinationAddress.IsUnknown() &&
				!config.IdpSensorConfiguration.LogSuppression.NoIncludeDestinationAddress.IsNull() &&
				!config.IdpSensorConfiguration.LogSuppression.NoIncludeDestinationAddress.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("idp_sensor_configuration").AtName("log_suppression").AtName("include_destination_address"),
					tfdiag.ConflictConfigErrSummary,
					"include_destination_address and no_include_destination_address can't be true in same time "+
						"in idp_sensor_configuration block in log_suppression block",
				)
			}
		}
		if config.IdpSensorConfiguration.PacketLog != nil {
			if config.IdpSensorConfiguration.PacketLog.SourceAddress.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("idp_sensor_configuration").AtName("packet_log").AtName("source_address"),
					tfdiag.MissingConfigErrSummary,
					"source_address must be specified in packet_log block in idp_sensor_configuration block",
				)
			}
			if !config.IdpSensorConfiguration.PacketLog.HostPort.IsNull() &&
				config.IdpSensorConfiguration.PacketLog.HostAddress.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("idp_sensor_configuration").AtName("packet_log").AtName("host_port"),
					tfdiag.MissingConfigErrSummary,
					"host_address must be specified with host_port in packet_log block in idp_sensor_configuration block",
				)
			}
		}
	}

	if config.IkeTraceoptions != nil {
		if config.IkeTraceoptions.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ike_traceoptions").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"ike_traceoptions block is empty",
			)
		}

		if config.IkeTraceoptions.File != nil {
			if config.IkeTraceoptions.File.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ike_traceoptions").AtName("file").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"file block is empty in ike_traceoptions block",
				)
			}
			if !config.IkeTraceoptions.File.WorldReadable.IsNull() && !config.IkeTraceoptions.File.WorldReadable.IsUnknown() &&
				!config.IkeTraceoptions.File.NoWorldReadable.IsNull() && !config.IkeTraceoptions.File.NoWorldReadable.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ike_traceoptions").AtName("file").AtName("world_readable"),
					tfdiag.ConflictConfigErrSummary,
					"world_readable and no_world_readable can't be true in same time "+
						"in file block in ike_traceoptions block",
				)
			}
		}
	}

	if config.Log != nil {
		if config.Log.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("log").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"log block is empty",
			)
		}
		if config.Log.File != nil {
			if config.Log.File.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("log").AtName("file").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"file block is empty in log block",
				)
			}
		}
		if !config.Log.SourceAddress.IsNull() && !config.Log.SourceAddress.IsUnknown() &&
			!config.Log.SourceInterface.IsNull() && !config.Log.SourceInterface.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("log").AtName("source_address"),
				tfdiag.ConflictConfigErrSummary,
				"source_address and source_interface can't be set in same time in log block",
			)
		}
	}

	if config.NatSource != nil {
		if config.NatSource.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nat_source").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"nat_source block is empty",
			)
		}
		if !config.NatSource.InterfacePortOverloadingFactor.IsNull() &&
			!config.NatSource.InterfacePortOverloadingFactor.IsUnknown() &&
			!config.NatSource.InterfacePortOverloadingOff.IsNull() &&
			!config.NatSource.InterfacePortOverloadingOff.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nat_source").AtName("interface_port_overloading_off"),
				tfdiag.ConflictConfigErrSummary,
				"interface_port_overloading_off and interface_port_overloading_factor cannot be configured together "+
					"in nat_source block",
			)
		}
		if !config.NatSource.PoolDefaultPortRangeTo.IsNull() &&
			config.NatSource.PoolDefaultPortRange.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nat_source").AtName("pool_default_port_range_to"),
				tfdiag.MissingConfigErrSummary,
				"pool_default_port_range must be specified with pool_default_port_range_to in nat_source block",
			)
		}
		if !config.NatSource.PoolDefaultPortRange.IsNull() &&
			config.NatSource.PoolDefaultPortRangeTo.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nat_source").AtName("pool_default_port_range"),
				tfdiag.MissingConfigErrSummary,
				"pool_default_port_range_to must be specified with pool_default_port_range in nat_source block",
			)
		}
		if !config.NatSource.PoolDefaultTwinPortRangeTo.IsNull() &&
			config.NatSource.PoolDefaultTwinPortRange.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nat_source").AtName("pool_default_twin_port_range_to"),
				tfdiag.MissingConfigErrSummary,
				"pool_default_twin_port_range must be specified with pool_default_twin_port_range_to in nat_source block",
			)
		}
		if !config.NatSource.PoolDefaultTwinPortRange.IsNull() &&
			config.NatSource.PoolDefaultTwinPortRangeTo.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nat_source").AtName("pool_default_twin_port_range"),
				tfdiag.MissingConfigErrSummary,
				"pool_default_twin_port_range_to must be specified with pool_default_twin_port_range in nat_source block",
			)
		}
		if !config.NatSource.PoolUtilizationAlarmClearThreshold.IsNull() &&
			config.NatSource.PoolUtilizationAlarmRaiseThreshold.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nat_source").AtName("pool_utilization_alarm_clear_threshold"),
				tfdiag.MissingConfigErrSummary,
				"pool_utilization_alarm_raise_threshold must be specified with pool_utilization_alarm_clear_threshold "+
					"in nat_source block",
			)
		}
		if !config.NatSource.PoolUtilizationAlarmClearThreshold.IsNull() &&
			!config.NatSource.PoolUtilizationAlarmClearThreshold.IsUnknown() &&
			!config.NatSource.PoolUtilizationAlarmRaiseThreshold.IsNull() &&
			!config.NatSource.PoolUtilizationAlarmRaiseThreshold.IsUnknown() {
			if config.NatSource.PoolUtilizationAlarmClearThreshold.ValueInt64() >
				config.NatSource.PoolUtilizationAlarmRaiseThreshold.ValueInt64() {
				resp.Diagnostics.AddAttributeError(
					path.Root("nat_source").AtName("pool_utilization_alarm_clear_threshold"),
					tfdiag.ConflictConfigErrSummary,
					"pool_utilization_alarm_clear_threshold must be larger than "+
						"pool_utilization_alarm_raise_threshold in nat_source block",
				)
			}
		}
	}

	if config.Policies != nil {
		if config.Policies.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("policies").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"policies block is empty",
			)
		}
		if !config.Policies.PolicyRematch.IsNull() && !config.Policies.PolicyRematch.IsUnknown() &&
			!config.Policies.PolicyRematchExtensive.IsNull() && !config.Policies.PolicyRematchExtensive.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("policies").AtName("policy_rematch"),
				tfdiag.ConflictConfigErrSummary,
				"policy_rematch and policy_rematch_extensive can't be true in same time in policies block",
			)
		}
	}

	if config.UserIdentificationAuthSource != nil {
		if config.UserIdentificationAuthSource.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("user_identification_auth_source").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"user_identification_auth_source block is empty",
			)
		}
	}

	if config.Utm != nil {
		if config.Utm.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("utm").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"utm block is empty",
			)
		}
	}
}

func (rsc *security) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(_ context.Context, junSess *junos.Session) bool {
			if !junSess.CheckCompatibilitySecurity() {
				resp.Diagnostics.AddError(
					tfdiag.CompatibilityErrSummary,
					rsc.junosName()+junSess.SystemInformation.NotCompatibleMsg(),
				)

				return false
			}

			return true
		},
		nil,
		&plan,
		resp,
	)
}

func (rsc *security) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityData
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

func (rsc *security) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityData
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

func (rsc *security) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityData
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

func (rsc *security) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityData

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

func (rscData *securityData) fillID() {
	rscData.ID = types.StringValue("security")
}

func (rscData *securityData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security "

	if rscData.Alg != nil {
		if rscData.Alg.isEmpty() {
			return path.Root("alg").AtName("*"),
				errors.New("alg block is empty")
		}

		configSet = append(configSet, rscData.Alg.configSet()...)
	}
	if rscData.Flow != nil {
		if rscData.Flow.isEmpty() {
			return path.Root("flow").AtName("*"),
				errors.New("flow block is empty")
		}

		blockSet, pathErr, err := rscData.Flow.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.ForwardingOptions != nil {
		if rscData.ForwardingOptions.isEmpty() {
			return path.Root("forwarding_options").AtName("*"),
				errors.New("forwarding_options block is empty")
		}

		configSet = append(configSet, rscData.ForwardingOptions.configSet()...)
	}
	if rscData.ForwardingProcess != nil {
		if rscData.ForwardingProcess.isEmpty() {
			return path.Root("forwarding_process").AtName("*"),
				errors.New("forwarding_process block is empty")
		}

		if rscData.ForwardingProcess.EnhancedServicesMode.ValueBool() {
			configSet = append(configSet, setPrefix+"forwarding-process enhanced-services-mode")
		}
	}
	if rscData.IdpSecurityPackage != nil {
		if rscData.IdpSecurityPackage.isEmpty() {
			return path.Root("idp_security_package").AtName("*"),
				errors.New("idp_security_package block is empty")
		}

		configSet = append(configSet, rscData.IdpSecurityPackage.configSet()...)
	}
	if rscData.IdpSensorConfiguration != nil {
		if rscData.IdpSensorConfiguration.isEmpty() {
			return path.Root("idp_sensor_configuration").AtName("*"),
				errors.New("idp_sensor_configuration block is empty")
		}

		configSet = append(configSet, rscData.IdpSensorConfiguration.configSet()...)
	}
	if rscData.IkeTraceoptions != nil {
		if rscData.IkeTraceoptions.isEmpty() {
			return path.Root("ike_traceoptions").AtName("*"),
				errors.New("ike_traceoptions block is empty")
		}

		blockSet, pathErr, err := rscData.IkeTraceoptions.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.Log != nil {
		if rscData.Log.isEmpty() {
			return path.Root("log").AtName("*"),
				errors.New("log block is empty")
		}

		blockSet, pathErr, err := rscData.Log.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.NatSource != nil {
		if rscData.NatSource.isEmpty() {
			return path.Root("nat_source").AtName("*"),
				errors.New("nat_source block is empty")
		}

		configSet = append(configSet, rscData.NatSource.configSet()...)
	}
	if rscData.Policies != nil {
		if rscData.Policies.isEmpty() {
			return path.Root("policies").AtName("*"),
				errors.New("policies block is empty")
		}

		if rscData.Policies.PolicyRematch.ValueBool() {
			configSet = append(configSet, setPrefix+"policies policy-rematch")
		}
		if rscData.Policies.PolicyRematchExtensive.ValueBool() {
			configSet = append(configSet, setPrefix+"policies policy-rematch extensive")
		}
	}
	if rscData.UserIdentificationAuthSource != nil {
		if rscData.UserIdentificationAuthSource.isEmpty() {
			return path.Root("user_identification_auth_source").AtName("*"),
				errors.New("user_identification_auth_source block is empty")
		}

		configSet = append(configSet, rscData.UserIdentificationAuthSource.configSet()...)
	}
	if rscData.Utm != nil {
		if rscData.Utm.isEmpty() {
			return path.Root("utm").AtName("*"),
				errors.New("utm block is empty")
		}

		configSet = append(configSet, rscData.Utm.configSet()...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityBlockAlg) configSet() []string {
	configSet := make([]string, 0)
	setPrefix := "set security alg "

	if block.DNSDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"dns disable")
	}
	if block.FtpDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"ftp disable")
	}
	if block.H323Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"h323 disable")
	}
	if block.MgcpDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"mgcp disable")
	}
	if block.MsrpcDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"msrpc disable")
	}
	if block.PptpDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"pptp disable")
	}
	if block.RshDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"rsh disable")
	}
	if block.RtspDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"rtsp disable")
	}
	if block.SccpDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"sccp disable")
	}
	if block.SIPDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"sip disable")
	}
	if block.SQLDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"sql disable")
	}
	if block.SunrpcDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"sunrpc disable")
	}
	if block.TalkDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"talk disable")
	}
	if block.TftpDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"tftp disable")
	}

	return configSet
}

func (block *securityBlockFlow) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix := "set security flow "

	if block.AdvancedOptions != nil {
		if block.AdvancedOptions.isEmpty() {
			return configSet, path.Root("flow").AtName("advanced_options").AtName("*"),
				errors.New("advanced_options block is empty in flow block")
		}
		if block.AdvancedOptions.DropMatchingLinkLocalAddress.ValueBool() {
			configSet = append(configSet, setPrefix+"advanced-options drop-matching-link-local-address")
		}
		if block.AdvancedOptions.DropMatchingReservedIPAddress.ValueBool() {
			configSet = append(configSet, setPrefix+"advanced-options drop-matching-reserved-ip-address")
		}
		if block.AdvancedOptions.ReverseRoutePacketModeVR.ValueBool() {
			configSet = append(configSet, setPrefix+"advanced-options reverse-route-packet-mode-vr")
		}
	}
	if block.Aging != nil {
		if block.Aging.isEmpty() {
			return configSet, path.Root("flow").AtName("aging").AtName("*"),
				errors.New("aging block is empty in flow block")
		}
		if !block.Aging.EarlyAgeout.IsNull() {
			configSet = append(configSet, setPrefix+"aging early-ageout "+
				utils.ConvI64toa(block.Aging.EarlyAgeout.ValueInt64()))
		}
		if !block.Aging.HighWatermark.IsNull() {
			configSet = append(configSet, setPrefix+"aging high-watermark "+
				utils.ConvI64toa(block.Aging.HighWatermark.ValueInt64()))
		}
		if !block.Aging.LowWatermark.IsNull() {
			configSet = append(configSet, setPrefix+"aging low-watermark "+
				utils.ConvI64toa(block.Aging.LowWatermark.ValueInt64()))
		}
	}
	if block.AllowDNSReply.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-dns-reply")
	}
	if block.AllowEmbeddedIcmp.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-embedded-icmp")
	}
	if block.AllowReverseEcmp.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-reverse-ecmp")
	}
	if block.EnableRerouteUniformLinkCheckNat.ValueBool() {
		configSet = append(configSet, setPrefix+"enable-reroute-uniform-link-check nat")
	}
	if block.EthernetSwitching != nil {
		if block.EthernetSwitching.isEmpty() {
			return configSet, path.Root("flow").AtName("ethernet_switching").AtName("*"),
				errors.New("ethernet_switching block is empty in flow block")
		}
		if block.EthernetSwitching.BlockNonIPAll.ValueBool() {
			configSet = append(configSet, setPrefix+"ethernet-switching block-non-ip-all")
		}
		if block.EthernetSwitching.BypassNonIPUnicast.ValueBool() {
			configSet = append(configSet, setPrefix+"ethernet-switching bypass-non-ip-unicast")
		}
		if block.EthernetSwitching.BpduVlanFlooding.ValueBool() {
			configSet = append(configSet, setPrefix+"ethernet-switching bpdu-vlan-flooding")
		}
		if block.EthernetSwitching.NoPacketFlooding != nil {
			configSet = append(configSet, setPrefix+"ethernet-switching no-packet-flooding")
			if block.EthernetSwitching.NoPacketFlooding.NoTraceRoute.ValueBool() {
				configSet = append(configSet, setPrefix+"ethernet-switching no-packet-flooding no-trace-route")
			}
		}
	}
	if block.ForceIPReassembly.ValueBool() {
		configSet = append(configSet, setPrefix+"force-ip-reassembly")
	}
	if block.IpsecPerformanceAcceleration.ValueBool() {
		configSet = append(configSet, setPrefix+"ipsec-performance-acceleration")
	}
	if block.McastBufferEnhance.ValueBool() {
		configSet = append(configSet, setPrefix+"mcast-buffer-enhance")
	}
	if v := block.PendingSessQueueLength.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pending-sess-queue-length "+v)
	}
	if block.PreserveIncomingFragmentSize.ValueBool() {
		configSet = append(configSet, setPrefix+"preserve-incoming-fragment-size")
	}
	if !block.RouteChangeTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"route-change-timeout "+
			utils.ConvI64toa(block.RouteChangeTimeout.ValueInt64()))
	}
	if v := block.SynFloodProtectionMode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"syn-flood-protection-mode "+v)
	}
	if block.SyncIcmpSession.ValueBool() {
		configSet = append(configSet, setPrefix+"sync-icmp-session")
	}
	if block.TCPMss != nil {
		if block.TCPMss.isEmpty() {
			return configSet, path.Root("flow").AtName("tcp_mss").AtName("*"),
				errors.New("tcp_mss block is empty in flow block")
		}
		if !block.TCPMss.AllTCPMss.IsNull() {
			configSet = append(configSet, setPrefix+"tcp-mss all-tcp mss "+
				utils.ConvI64toa(block.TCPMss.AllTCPMss.ValueInt64()))
		}
		if block.TCPMss.GreIn != nil {
			configSet = append(configSet, setPrefix+"tcp-mss gre-in")
			if !block.TCPMss.GreIn.Mss.IsNull() {
				configSet = append(configSet, setPrefix+"tcp-mss gre-in mss "+
					utils.ConvI64toa(block.TCPMss.GreIn.Mss.ValueInt64()))
			}
		}
		if block.TCPMss.GreOut != nil {
			configSet = append(configSet, setPrefix+"tcp-mss gre-out")
			if !block.TCPMss.GreOut.Mss.IsNull() {
				configSet = append(configSet, setPrefix+"tcp-mss gre-out mss "+
					utils.ConvI64toa(block.TCPMss.GreOut.Mss.ValueInt64()))
			}
		}
		if block.TCPMss.IpsecVpn != nil {
			configSet = append(configSet, setPrefix+"tcp-mss ipsec-vpn")
			if !block.TCPMss.IpsecVpn.Mss.IsNull() {
				configSet = append(configSet, setPrefix+"tcp-mss ipsec-vpn mss "+
					utils.ConvI64toa(block.TCPMss.IpsecVpn.Mss.ValueInt64()))
			}
		}
	}
	if block.TCPSession != nil {
		if block.TCPSession.isEmpty() {
			return configSet, path.Root("flow").AtName("tcp_session").AtName("*"),
				errors.New("tcp_session block is empty in flow block")
		}

		if block.TCPSession.FinInvalidateSession.ValueBool() {
			configSet = append(configSet, setPrefix+"tcp-session fin-invalidate-session")
		}
		if v := block.TCPSession.MaximumWindow.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"tcp-session maximum-window "+v)
		}
		if block.TCPSession.NoSequenceCheck.ValueBool() {
			configSet = append(configSet, setPrefix+"tcp-session no-sequence-check")
		}
		if block.TCPSession.NoSynCheck.ValueBool() {
			configSet = append(configSet, setPrefix+"tcp-session no-syn-check")
		}
		if block.TCPSession.NoSynCheckInTunnel.ValueBool() {
			configSet = append(configSet, setPrefix+"tcp-session no-syn-check-in-tunnel")
		}
		if block.TCPSession.RstInvalidateSession.ValueBool() {
			configSet = append(configSet, setPrefix+"tcp-session rst-invalidate-session")
		}
		if block.TCPSession.RstSequenceCheck.ValueBool() {
			configSet = append(configSet, setPrefix+"tcp-session rst-sequence-check")
		}
		if block.TCPSession.StrictSynCheck.ValueBool() {
			configSet = append(configSet, setPrefix+"tcp-session strict-syn-check")
		}
		if !block.TCPSession.TCPInitialTimeout.IsNull() {
			configSet = append(configSet, setPrefix+"tcp-session tcp-initial-timeout "+
				utils.ConvI64toa(block.TCPSession.TCPInitialTimeout.ValueInt64()))
		}
		if block.TCPSession.TimeWaitState != nil {
			configSet = append(configSet, setPrefix+"tcp-session time-wait-state")
			if block.TCPSession.TimeWaitState.ApplyToHalfCloseState.ValueBool() {
				configSet = append(configSet, setPrefix+"tcp-session time-wait-state apply-to-half-close-state")
			}
			if block.TCPSession.TimeWaitState.SessionAgeout.ValueBool() {
				configSet = append(configSet, setPrefix+"tcp-session time-wait-state session-ageout")
			}
			if !block.TCPSession.TimeWaitState.SessionTimeout.IsNull() {
				configSet = append(configSet, setPrefix+"tcp-session time-wait-state session-timeout "+
					utils.ConvI64toa(block.TCPSession.TimeWaitState.SessionTimeout.ValueInt64()))
			}
		}
	}

	return configSet, path.Empty(), nil
}

func (block *securityBlockForwardingOptions) configSet() []string {
	configSet := make([]string, 0)
	setPrefix := "set security forwarding-options "

	if v := block.Inet6Mode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"family inet6 mode "+v)
	}
	if block.IsoModePacketBased.ValueBool() {
		configSet = append(configSet, setPrefix+"family iso mode packet-based")
	}
	if v := block.MplsMode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"family mpls mode "+v)
	}

	return configSet
}

func (block *securityBlockIdpSecurityPackage) configSet() []string {
	configSet := make([]string, 0)
	setPrefix := "set security idp security-package "

	if block.AutomaticEnable.ValueBool() {
		configSet = append(configSet, setPrefix+"automatic enable")
	}
	if !block.AutomaticInterval.IsNull() {
		configSet = append(configSet, setPrefix+"automatic interval "+
			utils.ConvI64toa(block.AutomaticInterval.ValueInt64()))
	}
	if v := block.AutomaticStartTime.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"automatic start-time "+v)
	}
	if block.InstallIgnoreVersionCheck.ValueBool() {
		configSet = append(configSet, setPrefix+"install ignore-version-check")
	}
	if v := block.ProxyProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"proxy-profile \""+v+"\"")
	}
	if v := block.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}
	if v := block.URL.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"url \""+v+"\"")
	}

	return configSet
}

func (block *securityBlockIdpSensorConfiguration) configSet() []string {
	configSet := make([]string, 0)
	setPrefix := "set security idp sensor-configuration "

	if !block.LogCacheSize.IsNull() {
		configSet = append(configSet, setPrefix+"log cache-size "+
			utils.ConvI64toa(block.LogCacheSize.ValueInt64()))
	}
	if block.LogSuppression != nil {
		configSet = append(configSet, setPrefix+"log suppression")

		if block.LogSuppression.Disable.ValueBool() {
			configSet = append(configSet, setPrefix+"log suppression disable")
		}
		if block.LogSuppression.IncludeDestinationAddress.ValueBool() {
			configSet = append(configSet, setPrefix+"log suppression include-destination-address")
		}
		if block.LogSuppression.NoIncludeDestinationAddress.ValueBool() {
			configSet = append(configSet, setPrefix+"log suppression no-include-destination-address")
		}
		if !block.LogSuppression.MaxLogsOperate.IsNull() {
			configSet = append(configSet, setPrefix+"log suppression max-logs-operate "+
				utils.ConvI64toa(block.LogSuppression.MaxLogsOperate.ValueInt64()))
		}
		if !block.LogSuppression.MaxTimeReport.IsNull() {
			configSet = append(configSet, setPrefix+"log suppression max-time-report "+
				utils.ConvI64toa(block.LogSuppression.MaxTimeReport.ValueInt64()))
		}
		if !block.LogSuppression.StartLog.IsNull() {
			configSet = append(configSet, setPrefix+"log suppression start-log "+
				utils.ConvI64toa(block.LogSuppression.StartLog.ValueInt64()))
		}
	}
	if block.PacketLog != nil {
		configSet = append(configSet, setPrefix+"packet-log source-address "+block.PacketLog.SourceAddress.ValueString())

		if v := block.PacketLog.HostAddress.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"packet-log host "+v)
		}
		if !block.PacketLog.HostPort.IsNull() {
			configSet = append(configSet, setPrefix+"packet-log host port "+
				utils.ConvI64toa(block.PacketLog.HostPort.ValueInt64()))
		}
		if !block.PacketLog.MaxSessions.IsNull() {
			configSet = append(configSet, setPrefix+"packet-log max-sessions "+
				utils.ConvI64toa(block.PacketLog.MaxSessions.ValueInt64()))
		}
		if !block.PacketLog.ThresholdLoggingInterval.IsNull() {
			configSet = append(configSet, setPrefix+"packet-log threshold-logging-interval "+
				utils.ConvI64toa(block.PacketLog.ThresholdLoggingInterval.ValueInt64()))
		}
		if !block.PacketLog.TotalMemory.IsNull() {
			configSet = append(configSet, setPrefix+"packet-log total-memory "+
				utils.ConvI64toa(block.PacketLog.TotalMemory.ValueInt64()))
		}
	}
	if v := block.SecurityConfigurationProtectionMode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"security-configuration protection-mode "+v)
	}

	return configSet
}

func (block *securityBlockIkeTraceoptions) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix := "set security ike traceoptions "

	if block.File != nil {
		if block.File.isEmpty() {
			return configSet, path.Root("ike_traceoptions").AtName("file").AtName("*"),
				errors.New("file block is empty in ike_traceoptions block")
		}

		if v := block.File.Name.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"file \""+v+"\"")
		}
		if !block.File.Files.IsNull() {
			configSet = append(configSet, setPrefix+"file files "+
				utils.ConvI64toa(block.File.Files.ValueInt64()))
		}
		if v := block.File.Match.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"file match \""+v+"\"")
		}
		if !block.File.Size.IsNull() {
			configSet = append(configSet, setPrefix+"file size "+
				utils.ConvI64toa(block.File.Size.ValueInt64()))
		}
		if block.File.WorldReadable.ValueBool() && block.File.NoWorldReadable.ValueBool() {
			return configSet,
				path.Root("ike_traceoptions").AtName("file").AtName("world_readable"),
				errors.New("world_readable and no_world_readable can't be true in same time " +
					"in file block in ike_traceoptions block")
		}
		if block.File.WorldReadable.ValueBool() {
			configSet = append(configSet, setPrefix+"file world-readable")
		}
		if block.File.NoWorldReadable.ValueBool() {
			configSet = append(configSet, setPrefix+"file no-world-readable")
		}
	}
	for _, v := range block.Flag {
		configSet = append(configSet, setPrefix+"flag "+v.ValueString())
	}
	if block.NoRemoteTrace.ValueBool() {
		configSet = append(configSet, setPrefix+"no-remote-trace")
	}
	if !block.RateLimit.IsNull() {
		configSet = append(configSet, setPrefix+"rate-limit "+
			utils.ConvI64toa(block.RateLimit.ValueInt64()))
	}

	return configSet, path.Empty(), nil
}

func (block *securityBlockLog) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix := "set security log "

	if block.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if !block.EventRate.IsNull() {
		configSet = append(configSet, setPrefix+"event-rate "+
			utils.ConvI64toa(block.EventRate.ValueInt64()))
	}
	if v := block.FacilityOverride.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"facility-override "+v)
	}
	if block.File != nil {
		if block.File.isEmpty() {
			return configSet, path.Root("log").AtName("file").AtName("*"),
				errors.New("file block is empty in log block")
		}

		if !block.File.Files.IsNull() {
			configSet = append(configSet, setPrefix+"file files "+
				utils.ConvI64toa(block.File.Files.ValueInt64()))
		}
		if v := block.File.Name.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"file name \""+v+"\"")
		}
		if v := block.File.Path.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"file path \""+v+"\"")
		}
		if !block.File.Size.IsNull() {
			configSet = append(configSet, setPrefix+"file size "+
				utils.ConvI64toa(block.File.Size.ValueInt64()))
		}
	}
	if v := block.Format.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"format "+v)
	}
	if !block.MaxDatabaseRecord.IsNull() {
		configSet = append(configSet, setPrefix+"max-database-record "+
			utils.ConvI64toa(block.MaxDatabaseRecord.ValueInt64()))
	}
	if v := block.Mode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"mode "+v)
	}
	if !block.RateCap.IsNull() {
		configSet = append(configSet, setPrefix+"rate-cap "+
			utils.ConvI64toa(block.RateCap.ValueInt64()))
	}
	if block.Report.ValueBool() {
		configSet = append(configSet, setPrefix+"report")
	}
	if v := block.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}
	if v := block.SourceInterface.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-interface "+v)
	}
	if block.Transport != nil {
		configSet = append(configSet, setPrefix+"transport")

		if v := block.Transport.Protocol.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"transport protocol "+v)
		}
		if !block.Transport.TCPConnections.IsNull() {
			configSet = append(configSet, setPrefix+"transport tcp-connections "+
				utils.ConvI64toa(block.Transport.TCPConnections.ValueInt64()))
		}
		if v := block.Transport.TLSProfile.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"transport tls-profile \""+v+"\"")
		}
	}
	if block.UtcTimestamp.ValueBool() {
		configSet = append(configSet, setPrefix+"utc-timestamp")
	}

	return configSet, path.Empty(), nil
}

func (block *securityBlockNatSource) configSet() []string {
	configSet := make([]string, 0)
	setPrefix := "set security nat source "

	if block.AddressPersistent.ValueBool() {
		configSet = append(configSet, setPrefix+"address-persistent")
	}
	if !block.InterfacePortOverloadingFactor.IsNull() {
		configSet = append(configSet, setPrefix+"interface port-overloading-factor "+
			utils.ConvI64toa(block.InterfacePortOverloadingFactor.ValueInt64()))
	}
	if block.InterfacePortOverloadingOff.ValueBool() {
		configSet = append(configSet, setPrefix+"interface port-overloading off")
	}
	if !block.PoolDefaultPortRange.IsNull() {
		configSet = append(configSet, setPrefix+"pool-default-port-range "+
			utils.ConvI64toa(block.PoolDefaultPortRange.ValueInt64()))
	}
	if !block.PoolDefaultPortRangeTo.IsNull() {
		configSet = append(configSet, setPrefix+"pool-default-port-range to "+
			utils.ConvI64toa(block.PoolDefaultPortRangeTo.ValueInt64()))
	}
	if !block.PoolDefaultTwinPortRange.IsNull() {
		configSet = append(configSet, setPrefix+"pool-default-twin-port-range "+
			utils.ConvI64toa(block.PoolDefaultTwinPortRange.ValueInt64()))
	}
	if !block.PoolDefaultTwinPortRangeTo.IsNull() {
		configSet = append(configSet, setPrefix+"pool-default-twin-port-range to "+
			utils.ConvI64toa(block.PoolDefaultTwinPortRangeTo.ValueInt64()))
	}
	if !block.PoolUtilizationAlarmClearThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"pool-utilization-alarm clear-threshold "+
			utils.ConvI64toa(block.PoolUtilizationAlarmClearThreshold.ValueInt64()))
	}
	if !block.PoolUtilizationAlarmRaiseThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"pool-utilization-alarm raise-threshold "+
			utils.ConvI64toa(block.PoolUtilizationAlarmRaiseThreshold.ValueInt64()))
	}
	if block.PortRandomizationDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"port-randomization disable")
	}
	if !block.SessionDropHoldDown.IsNull() {
		configSet = append(configSet, setPrefix+"session-drop-hold-down "+
			utils.ConvI64toa(block.SessionDropHoldDown.ValueInt64()))
	}
	if block.SessionPersistenceScan.ValueBool() {
		configSet = append(configSet, setPrefix+"session-persistence-scan")
	}

	return configSet
}

func (block *securityBlockUserIdentificationAuthSource) configSet() []string {
	configSet := make([]string, 0)
	setPrefix := "set security user-identification authentication-source "

	if !block.ADAuthPriority.IsNull() {
		configSet = append(configSet, setPrefix+"active-directory-authentication-table priority "+
			utils.ConvI64toa(block.ADAuthPriority.ValueInt64()))
	}
	if !block.ArubaClearpassPriority.IsNull() {
		configSet = append(configSet, setPrefix+"aruba-clearpass priority "+
			utils.ConvI64toa(block.ArubaClearpassPriority.ValueInt64()))
	}
	if !block.FirewallAuthPriority.IsNull() {
		configSet = append(configSet, setPrefix+"firewall-authentication priority "+
			utils.ConvI64toa(block.FirewallAuthPriority.ValueInt64()))
	}
	if !block.LocalAuthPriority.IsNull() {
		configSet = append(configSet, setPrefix+"local-authentication-table priority "+
			utils.ConvI64toa(block.LocalAuthPriority.ValueInt64()))
	}
	if !block.UnifiedAccessControlPriority.IsNull() {
		configSet = append(configSet, setPrefix+"unified-access-control priority "+
			utils.ConvI64toa(block.UnifiedAccessControlPriority.ValueInt64()))
	}

	return configSet
}

func (block *securityBlockUtm) configSet() []string {
	configSet := make([]string, 0)
	setPrefix := "set security utm "

	if v := block.FeatureProfileWebFilteringType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"feature-profile web-filtering type "+v)
	}
	if block.FeatureProfileWebFilteringJuniperEnhancedServer != nil {
		configSet = append(configSet, setPrefix+"feature-profile web-filtering juniper-enhanced server")

		if v := block.FeatureProfileWebFilteringJuniperEnhancedServer.Host.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"feature-profile web-filtering juniper-enhanced server host \""+v+"\"")
		}
		if !block.FeatureProfileWebFilteringJuniperEnhancedServer.Port.IsNull() {
			configSet = append(configSet, setPrefix+"feature-profile web-filtering juniper-enhanced server port "+
				utils.ConvI64toa(block.FeatureProfileWebFilteringJuniperEnhancedServer.Port.ValueInt64()))
		}
		if v := block.FeatureProfileWebFilteringJuniperEnhancedServer.ProxyProfile.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+
				"feature-profile web-filtering juniper-enhanced server proxy-profile \""+v+"\"")
		}
		if v := block.FeatureProfileWebFilteringJuniperEnhancedServer.RoutingInstance.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"feature-profile web-filtering juniper-enhanced server routing-instance "+v)
		}
	}

	return configSet
}

func (rscData *securityData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security" + junos.PipeDisplaySetRelative)
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
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockAlg{}.junosLines()):
				if rscData.Alg == nil {
					rscData.Alg = &securityBlockAlg{}
				}
				rscData.Alg.read(itemTrim)
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockFlow{}.junosLines()):
				if rscData.Flow == nil {
					rscData.Flow = &securityBlockFlow{}
				}
				if err := rscData.Flow.read(itemTrim); err != nil {
					return err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockForwardingOptions{}.junosLines()):
				if rscData.ForwardingOptions == nil {
					rscData.ForwardingOptions = &securityBlockForwardingOptions{}
				}
				rscData.ForwardingOptions.read(itemTrim)
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockForwardingProcess{}.junosLines()):
				if rscData.ForwardingProcess == nil {
					rscData.ForwardingProcess = &securityBlockForwardingProcess{}
				}
				if itemTrim == "forwarding-process enhanced-services-mode" {
					rscData.ForwardingProcess.EnhancedServicesMode = types.BoolValue(true)
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockIdpSecurityPackage{}.junosLines()):
				if rscData.IdpSecurityPackage == nil {
					rscData.IdpSecurityPackage = &securityBlockIdpSecurityPackage{}
				}
				if err := rscData.IdpSecurityPackage.read(itemTrim); err != nil {
					return err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockIdpSensorConfiguration{}.junosLines()):
				if rscData.IdpSensorConfiguration == nil {
					rscData.IdpSensorConfiguration = &securityBlockIdpSensorConfiguration{}
				}
				if err := rscData.IdpSensorConfiguration.read(itemTrim); err != nil {
					return err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockLog{}.junosLines()):
				if rscData.Log == nil {
					rscData.Log = &securityBlockLog{}
				}
				if err := rscData.Log.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "ike traceoptions "):
				if rscData.IkeTraceoptions == nil {
					rscData.IkeTraceoptions = &securityBlockIkeTraceoptions{}
				}
				if err := rscData.IkeTraceoptions.read(itemTrim); err != nil {
					return err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockNatSource{}.junosLines()):
				if rscData.NatSource == nil {
					rscData.NatSource = &securityBlockNatSource{}
				}
				if err := rscData.NatSource.read(itemTrim); err != nil {
					return err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockPolicies{}.junosLines()):
				if rscData.Policies == nil {
					rscData.Policies = &securityBlockPolicies{}
				}
				if itemTrim == "policies policy-rematch" {
					rscData.Policies.PolicyRematch = types.BoolValue(true)
				}
				if itemTrim == "policies policy-rematch extensive" {
					rscData.Policies.PolicyRematchExtensive = types.BoolValue(true)
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockUserIdentificationAuthSource{}.junosLines()):
				if rscData.UserIdentificationAuthSource == nil {
					rscData.UserIdentificationAuthSource = &securityBlockUserIdentificationAuthSource{}
				}
				if err := rscData.UserIdentificationAuthSource.read(itemTrim); err != nil {
					return err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, securityBlockUtm{}.junosLines()):
				if rscData.Utm == nil {
					rscData.Utm = &securityBlockUtm{}
				}
				if err := rscData.Utm.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (securityBlockAlg) junosLines() []string {
	return []string{
		"alg dns disable",
		"alg ftp disable",
		"alg h323 disable",
		"alg mgcp disable",
		"alg msrpc disable",
		"alg pptp disable",
		"alg rsh disable",
		"alg rtsp disable",
		"alg sccp disable",
		"alg sip disable",
		"alg sql disable",
		"alg sunrpc disable",
		"alg talk disable",
		"alg tftp disable",
	}
}

func (block *securityBlockAlg) read(itemTrim string) {
	balt.CutPrefixInString(&itemTrim, "alg ")

	if itemTrim == "dns disable" {
		block.DNSDisable = types.BoolValue(true)
	}
	if itemTrim == "ftp disable" {
		block.FtpDisable = types.BoolValue(true)
	}
	if itemTrim == "h323 disable" {
		block.H323Disable = types.BoolValue(true)
	}
	if itemTrim == "mgcp disable" {
		block.MgcpDisable = types.BoolValue(true)
	}
	if itemTrim == "msrpc disable" {
		block.MsrpcDisable = types.BoolValue(true)
	}
	if itemTrim == "pptp disable" {
		block.PptpDisable = types.BoolValue(true)
	}
	if itemTrim == "rsh disable" {
		block.RshDisable = types.BoolValue(true)
	}
	if itemTrim == "rtsp disable" {
		block.RtspDisable = types.BoolValue(true)
	}
	if itemTrim == "sccp disable" {
		block.SccpDisable = types.BoolValue(true)
	}
	if itemTrim == "sip disable" {
		block.SIPDisable = types.BoolValue(true)
	}
	if itemTrim == "sql disable" {
		block.SQLDisable = types.BoolValue(true)
	}
	if itemTrim == "sunrpc disable" {
		block.SunrpcDisable = types.BoolValue(true)
	}
	if itemTrim == "talk disable" {
		block.TalkDisable = types.BoolValue(true)
	}
	if itemTrim == "tftp disable" {
		block.TftpDisable = types.BoolValue(true)
	}
}

func (securityBlockFlow) junosLines() []string {
	return []string{
		"flow advanced-options",
		"flow aging",
		"flow allow-dns-reply",
		"flow allow-embedded-icmp",
		"flow allow-reverse-ecmp",
		"flow enable-reroute-uniform-link-check",
		"flow ethernet-switching",
		"flow force-ip-reassembly",
		"flow ipsec-performance-acceleration",
		"flow mcast-buffer-enhance",
		"flow pending-sess-queue-length",
		"flow preserve-incoming-fragment-size",
		"flow route-change-timeout",
		"flow syn-flood-protection-mode",
		"flow sync-icmp-session",
		"flow tcp-mss",
		"flow tcp-session",
	}
}

func (block *securityBlockFlow) read(itemTrim string) (err error) {
	balt.CutPrefixInString(&itemTrim, "flow ")

	switch {
	case balt.CutPrefixInString(&itemTrim, "advanced-options"):
		if block.AdvancedOptions == nil {
			block.AdvancedOptions = &securityBlockFlowBlockAdvancedOptions{}
		}
		switch {
		case itemTrim == " drop-matching-link-local-address":
			block.AdvancedOptions.DropMatchingLinkLocalAddress = types.BoolValue(true)
		case itemTrim == " drop-matching-reserved-ip-address":
			block.AdvancedOptions.DropMatchingReservedIPAddress = types.BoolValue(true)
		case itemTrim == " reverse-route-packet-mode-vr":
			block.AdvancedOptions.ReverseRoutePacketModeVR = types.BoolValue(true)
		}
	case balt.CutPrefixInString(&itemTrim, "aging"):
		if block.Aging == nil {
			block.Aging = &securityBlockFlowBlockAging{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, " early-ageout "):
			block.Aging.EarlyAgeout, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " high-watermark "):
			block.Aging.HighWatermark, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " low-watermark "):
			block.Aging.LowWatermark, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case itemTrim == "allow-dns-reply":
		block.AllowDNSReply = types.BoolValue(true)
	case itemTrim == "allow-embedded-icmp":
		block.AllowEmbeddedIcmp = types.BoolValue(true)
	case itemTrim == "allow-reverse-ecmp":
		block.AllowReverseEcmp = types.BoolValue(true)
	case itemTrim == "enable-reroute-uniform-link-check nat":
		block.EnableRerouteUniformLinkCheckNat = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "ethernet-switching"):
		if block.EthernetSwitching == nil {
			block.EthernetSwitching = &securityBlockFlowBlockEthernetSwitching{}
		}
		switch {
		case itemTrim == " block-non-ip-all":
			block.EthernetSwitching.BlockNonIPAll = types.BoolValue(true)
		case itemTrim == " bypass-non-ip-unicast":
			block.EthernetSwitching.BypassNonIPUnicast = types.BoolValue(true)
		case itemTrim == " bpdu-vlan-flooding":
			block.EthernetSwitching.BpduVlanFlooding = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, " no-packet-flooding"):
			if block.EthernetSwitching.NoPacketFlooding == nil {
				block.EthernetSwitching.NoPacketFlooding = &struct {
					NoTraceRoute types.Bool `tfsdk:"no_trace_route"`
				}{}
			}
			if itemTrim == " no-trace-route" {
				block.EthernetSwitching.NoPacketFlooding.NoTraceRoute = types.BoolValue(true)
			}
		}
	case itemTrim == "force-ip-reassembly":
		block.ForceIPReassembly = types.BoolValue(true)
	case itemTrim == "ipsec-performance-acceleration":
		block.IpsecPerformanceAcceleration = types.BoolValue(true)
	case itemTrim == "mcast-buffer-enhance":
		block.McastBufferEnhance = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "pending-sess-queue-length "):
		block.PendingSessQueueLength = types.StringValue(itemTrim)
	case itemTrim == "preserve-incoming-fragment-size":
		block.PreserveIncomingFragmentSize = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "route-change-timeout "):
		block.RouteChangeTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "syn-flood-protection-mode "):
		block.SynFloodProtectionMode = types.StringValue(itemTrim)
	case itemTrim == "sync-icmp-session":
		block.SyncIcmpSession = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "tcp-mss "):
		if block.TCPMss == nil {
			block.TCPMss = &securityBlockFlowBlockTCPMss{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, "all-tcp mss "):
			block.TCPMss.AllTCPMss, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "gre-in"):
			if block.TCPMss.GreIn == nil {
				block.TCPMss.GreIn = &struct {
					Mss types.Int64 `tfsdk:"mss"`
				}{}
			}
			if balt.CutPrefixInString(&itemTrim, " mss ") {
				block.TCPMss.GreIn.Mss, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		case balt.CutPrefixInString(&itemTrim, "gre-out"):
			if block.TCPMss.GreOut == nil {
				block.TCPMss.GreOut = &struct {
					Mss types.Int64 `tfsdk:"mss"`
				}{}
			}
			if balt.CutPrefixInString(&itemTrim, " mss ") {
				block.TCPMss.GreOut.Mss, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		case balt.CutPrefixInString(&itemTrim, "ipsec-vpn"):
			if block.TCPMss.IpsecVpn == nil {
				block.TCPMss.IpsecVpn = &struct {
					Mss types.Int64 `tfsdk:"mss"`
				}{}
			}
			if balt.CutPrefixInString(&itemTrim, " mss ") {
				block.TCPMss.IpsecVpn.Mss, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	case balt.CutPrefixInString(&itemTrim, "tcp-session "):
		if block.TCPSession == nil {
			block.TCPSession = &securityBlockFlowBlockTCPSession{}
		}
		switch {
		case itemTrim == "fin-invalidate-session":
			block.TCPSession.FinInvalidateSession = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "maximum-window "):
			block.TCPSession.MaximumWindow = types.StringValue(itemTrim)
		case itemTrim == "no-sequence-check":
			block.TCPSession.NoSequenceCheck = types.BoolValue(true)
		case itemTrim == "no-syn-check":
			block.TCPSession.NoSynCheck = types.BoolValue(true)
		case itemTrim == "no-syn-check-in-tunnel":
			block.TCPSession.NoSynCheckInTunnel = types.BoolValue(true)
		case itemTrim == "rst-invalidate-session":
			block.TCPSession.RstInvalidateSession = types.BoolValue(true)
		case itemTrim == "rst-sequence-check":
			block.TCPSession.RstSequenceCheck = types.BoolValue(true)
		case itemTrim == "strict-syn-check":
			block.TCPSession.StrictSynCheck = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "tcp-initial-timeout "):
			block.TCPSession.TCPInitialTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "time-wait-state"):
			if block.TCPSession.TimeWaitState == nil {
				block.TCPSession.TimeWaitState = &securityBlockFlowBlockTCPSessionBlockTimeWaitState{}
			}
			switch {
			case itemTrim == " apply-to-half-close-state":
				block.TCPSession.TimeWaitState.ApplyToHalfCloseState = types.BoolValue(true)
			case itemTrim == " session-ageout":
				block.TCPSession.TimeWaitState.SessionAgeout = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, " session-timeout "):
				block.TCPSession.TimeWaitState.SessionTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (securityBlockForwardingOptions) junosLines() []string {
	return []string{
		"forwarding-options family mpls mode",
		"forwarding-options family inet6 mode",
		"forwarding-options family iso mode",
	}
}

func (block *securityBlockForwardingOptions) read(itemTrim string) {
	balt.CutPrefixInString(&itemTrim, "forwarding-options ")

	switch {
	case balt.CutPrefixInString(&itemTrim, "family inet6 mode "):
		block.Inet6Mode = types.StringValue(itemTrim)
	case itemTrim == "family iso mode packet-based":
		block.IsoModePacketBased = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "family mpls mode "):
		block.MplsMode = types.StringValue(itemTrim)
	}
}

func (securityBlockForwardingProcess) junosLines() []string {
	return []string{
		"forwarding-process enhanced-services-mode",
	}
}

func (securityBlockIdpSecurityPackage) junosLines() []string {
	return []string{
		"idp security-package automatic",
		"idp security-package install",
		"idp security-package proxy-profile",
		"idp security-package source-address",
		"idp security-package url",
	}
}

func (block *securityBlockIdpSecurityPackage) read(itemTrim string) (err error) {
	balt.CutPrefixInString(&itemTrim, "idp security-package ")

	switch {
	case itemTrim == "automatic enable":
		block.AutomaticEnable = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "automatic interval "):
		block.AutomaticInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "automatic start-time "):
		block.AutomaticStartTime = types.StringValue(strings.Split(strings.Trim(itemTrim, "\""), " ")[0])
	case itemTrim == "install ignore-version-check":
		block.InstallIgnoreVersionCheck = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "proxy-profile "):
		block.ProxyProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		block.SourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "url "):
		block.URL = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return nil
}

func (securityBlockIdpSensorConfiguration) junosLines() []string {
	return []string{
		"idp sensor-configuration log",
		"idp sensor-configuration packet-log",
		"idp sensor-configuration security-configuration",
	}
}

func (block *securityBlockIdpSensorConfiguration) read(itemTrim string) (err error) {
	balt.CutPrefixInString(&itemTrim, "idp sensor-configuration ")

	switch {
	case balt.CutPrefixInString(&itemTrim, "log cache-size "):
		block.LogCacheSize, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "log suppression"):
		if block.LogSuppression == nil {
			block.LogSuppression = &securityBlockIdpSensorConfigurationBlockLogSuppression{}
		}
		switch {
		case itemTrim == " disable":
			block.LogSuppression.Disable = types.BoolValue(true)
		case itemTrim == " include-destination-address":
			block.LogSuppression.IncludeDestinationAddress = types.BoolValue(true)
		case itemTrim == " no-include-destination-address":
			block.LogSuppression.NoIncludeDestinationAddress = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, " max-logs-operate "):
			block.LogSuppression.MaxLogsOperate, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " max-time-report "):
			block.LogSuppression.MaxTimeReport, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " start-log "):
			block.LogSuppression.StartLog, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "packet-log "):
		if block.PacketLog == nil {
			block.PacketLog = &securityBlockIdpSensorConfigurationBlockPacketLog{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, "source-address "):
			block.PacketLog.SourceAddress = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "host port "):
			block.PacketLog.HostPort, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "host "):
			block.PacketLog.HostAddress = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "max-sessions "):
			block.PacketLog.MaxSessions, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "threshold-logging-interval "):
			block.PacketLog.ThresholdLoggingInterval, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "total-memory "):
			block.PacketLog.TotalMemory, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "security-configuration protection-mode "):
		block.SecurityConfigurationProtectionMode = types.StringValue(itemTrim)
	}

	return nil
}

func (securityBlockIkeTraceoptions) junosLines() []string {
	return []string{
		"ike traceoptions",
	}
}

func (block *securityBlockIkeTraceoptions) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "file"):
		if block.File == nil {
			block.File = &securityBlockIkeTraceoptionsBlockFile{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, " files "):
			block.File.Files, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " match "):
			block.File.Match = types.StringValue(strings.Trim(itemTrim, "\""))
		case balt.CutPrefixInString(&itemTrim, " size "):
			switch {
			case balt.CutSuffixInString(&itemTrim, "k"):
				block.File.Size, err = tfdata.ConvAtoi64Value(itemTrim)
				block.File.Size = types.Int64Value(block.File.Size.ValueInt64() * 1024)
			case balt.CutSuffixInString(&itemTrim, "m"):
				block.File.Size, err = tfdata.ConvAtoi64Value(itemTrim)
				block.File.Size = types.Int64Value(block.File.Size.ValueInt64() * 1024 * 1024)
			case balt.CutSuffixInString(&itemTrim, "g"):
				block.File.Size, err = tfdata.ConvAtoi64Value(itemTrim)
				block.File.Size = types.Int64Value(block.File.Size.ValueInt64() * 1024 * 1024 * 1024)
			default:
				block.File.Size, err = tfdata.ConvAtoi64Value(itemTrim)
			}
			if err != nil {
				return err
			}
		case itemTrim == " world-readable":
			block.File.WorldReadable = types.BoolValue(true)
		case itemTrim == " no-world-readable":
			block.File.NoWorldReadable = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, " "):
			block.File.Name = types.StringValue(strings.Trim(itemTrim, "\""))
		}
	case balt.CutPrefixInString(&itemTrim, "flag "):
		block.Flag = append(block.Flag, types.StringValue(itemTrim))
	case itemTrim == "no-remote-trace":
		block.NoRemoteTrace = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "rate-limit "):
		block.RateLimit, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (securityBlockLog) junosLines() []string {
	return []string{
		"log disable",
		"log event-rate",
		"log facility-override",
		"log file",
		"log format",
		"log max-database-record",
		"log mode",
		"log rate-cap",
		"log report",
		"log source-address",
		"log source-interface",
		"log transport",
		"log utc-timestamp",
	}
}

func (block *securityBlockLog) read(itemTrim string) (err error) {
	balt.CutPrefixInString(&itemTrim, "log ")

	switch {
	case itemTrim == junos.DisableW:
		block.Disable = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "event-rate "):
		block.EventRate, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "facility-override "):
		block.FacilityOverride = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "file"):
		if block.File == nil {
			block.File = &securityBlockLogBlockFile{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, " files "):
			block.File.Files, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " name "):
			block.File.Name = types.StringValue(strings.Trim(itemTrim, "\""))
		case balt.CutPrefixInString(&itemTrim, " path "):
			block.File.Path = types.StringValue(strings.Trim(itemTrim, "\""))
		case balt.CutPrefixInString(&itemTrim, " size "):
			block.File.Size, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "format "):
		block.Format = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "max-database-record "):
		block.MaxDatabaseRecord, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "mode "):
		block.Mode = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "rate-cap "):
		block.RateCap, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "report":
		block.Report = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		block.SourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-interface "):
		block.SourceInterface = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "transport"):
		if block.Transport == nil {
			block.Transport = &securityBlockLogBlockTransport{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, " protocol "):
			block.Transport.Protocol = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, " tcp-connections "):
			block.Transport.TCPConnections, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " tls-profile "):
			block.Transport.TLSProfile = types.StringValue(strings.Trim(itemTrim, "\""))
		}
	case itemTrim == "utc-timestamp":
		block.UtcTimestamp = types.BoolValue(true)
	}

	return nil
}

func (securityBlockNatSource) junosLines() []string {
	return []string{
		"nat source address-persistent",
		"nat source interface port-overloading",
		"nat source interface port-overloading-factor",
		"nat source pool-default-port-range",
		"nat source pool-default-twin-port-range",
		"nat source pool-utilization-alarm",
		"nat source port-randomization",
		"nat source session-drop-hold-down",
		"nat source session-persistence-scan",
	}
}

func (block *securityBlockNatSource) read(itemTrim string) (err error) {
	balt.CutPrefixInString(&itemTrim, "nat source ")

	switch {
	case itemTrim == "address-persistent":
		block.AddressPersistent = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "interface port-overloading-factor "):
		block.InterfacePortOverloadingFactor, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "interface port-overloading off":
		block.InterfacePortOverloadingOff = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "pool-default-port-range to "):
		block.PoolDefaultPortRangeTo, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "pool-default-port-range "):
		block.PoolDefaultPortRange, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "pool-default-twin-port-range to "):
		block.PoolDefaultTwinPortRangeTo, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "pool-default-twin-port-range "):
		block.PoolDefaultTwinPortRange, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "pool-utilization-alarm clear-threshold "):
		block.PoolUtilizationAlarmClearThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "pool-utilization-alarm raise-threshold "):
		block.PoolUtilizationAlarmRaiseThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "port-randomization disable":
		block.PortRandomizationDisable = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "session-drop-hold-down "):
		block.SessionDropHoldDown, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "session-persistence-scan":
		block.SessionPersistenceScan = types.BoolValue(true)
	}
	if err != nil {
		return err
	}

	return nil
}

func (securityBlockPolicies) junosLines() []string {
	return []string{
		"policies policy-rematch",
	}
}

func (securityBlockUserIdentificationAuthSource) junosLines() []string {
	return []string{
		"user-identification authentication-source active-directory-authentication-table",
		"user-identification authentication-source aruba-clearpass",
		"user-identification authentication-source firewall-authentication",
		"user-identification authentication-source local-authentication-table",
		"user-identification authentication-source unified-access-control",
	}
}

func (block *securityBlockUserIdentificationAuthSource) read(itemTrim string) (err error) {
	balt.CutPrefixInString(&itemTrim, "user-identification authentication-source ")

	switch {
	case balt.CutPrefixInString(&itemTrim, "active-directory-authentication-table priority "):
		block.ADAuthPriority, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "aruba-clearpass priority "):
		block.ArubaClearpassPriority, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "firewall-authentication priority "):
		block.FirewallAuthPriority, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "local-authentication-table priority "):
		block.LocalAuthPriority, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "unified-access-control priority "):
		block.UnifiedAccessControlPriority, err = tfdata.ConvAtoi64Value(itemTrim)
	}
	if err != nil {
		return err
	}

	return nil
}

func (securityBlockUtm) junosLines() []string {
	return []string{
		"utm feature-profile web-filtering type",
		"utm feature-profile web-filtering juniper-enhanced server",
	}
}

func (block *securityBlockUtm) read(itemTrim string) (err error) {
	balt.CutPrefixInString(&itemTrim, "utm ")

	switch {
	case balt.CutPrefixInString(&itemTrim, "feature-profile web-filtering type "):
		block.FeatureProfileWebFilteringType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "feature-profile web-filtering juniper-enhanced server"):
		if block.FeatureProfileWebFilteringJuniperEnhancedServer == nil {
			block.FeatureProfileWebFilteringJuniperEnhancedServer = &securityBlockUtmBlockFeatureProfileWebFilteringJuniperEnhancedServer{} //nolint:lll
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, " host "):
			block.FeatureProfileWebFilteringJuniperEnhancedServer.Host = types.StringValue(strings.Trim(itemTrim, "\""))
		case balt.CutPrefixInString(&itemTrim, " port "):
			block.FeatureProfileWebFilteringJuniperEnhancedServer.Port, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " proxy-profile "):
			block.FeatureProfileWebFilteringJuniperEnhancedServer.ProxyProfile = types.StringValue(strings.Trim(itemTrim, "\""))
		case balt.CutPrefixInString(&itemTrim, " routing-instance "):
			block.FeatureProfileWebFilteringJuniperEnhancedServer.RoutingInstance = types.StringValue(itemTrim)
		}
	}

	return nil
}

func (rscData *securityData) del(
	_ context.Context, junSess *junos.Session,
) error {
	listLinesToDelete := make([]string, 0, 50)

	listLinesToDelete = append(listLinesToDelete, securityBlockAlg{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockFlow{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockForwardingOptions{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockForwardingProcess{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockIdpSecurityPackage{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockIdpSensorConfiguration{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockIkeTraceoptions{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockLog{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockNatSource{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockPolicies{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockUserIdentificationAuthSource{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, securityBlockUtm{}.junosLines()...)

	configSet := make([]string, len(listLinesToDelete))
	delPrefix := "delete security "
	for i, line := range listLinesToDelete {
		configSet[i] = delPrefix + line
	}

	return junSess.ConfigSet(configSet)
}
