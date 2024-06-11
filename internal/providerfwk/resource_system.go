package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"html"
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
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &system{}
	_ resource.ResourceWithConfigure      = &system{}
	_ resource.ResourceWithValidateConfig = &system{}
	_ resource.ResourceWithImportState    = &system{}
	_ resource.ResourceWithUpgradeState   = &system{}
)

type system struct {
	client *junos.Client
}

func newSystemResource() resource.Resource {
	return &system{}
}

func (rsc *system) typeName() string {
	return providerName + "_system"
}

func (rsc *system) junosName() string {
	return "system"
}

func (rsc *system) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *system) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *system) Configure(
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

func (rsc *system) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `system`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"authentication_order": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Order in which authentication methods are invoked.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("password", "radius", "tacplus"),
					),
				},
			},
			"auto_snapshot": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable auto-snapshot when boots from alternate slice.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"default_address_selection": schema.BoolAttribute{
				Optional:    true,
				Description: "Use loopback interface as source address for locally generated packets.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"domain_name": schema.StringAttribute{
				Optional:    true,
				Description: "Domain name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
				},
			},
			"host_name": schema.StringAttribute{
				Optional:    true,
				Description: "Hostname.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
				},
			},
			"max_configuration_rollbacks": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum rollback configuration.",
				Validators: []validator.Int64{
					int64validator.Between(0, 49),
				},
			},
			"max_configurations_on_flash": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of configuration files stored on flash.",
				Validators: []validator.Int64{
					int64validator.Between(0, 49),
				},
			},
			"name_server": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "DNS name servers.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						tfvalidator.StringIPAddress(),
					),
				},
			},
			"no_multicast_echo": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable responding to ICMP echo requests sent to multicast group addresses.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_ping_record_route": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not insert IP address in ping replies.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_ping_time_stamp": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not insert time stamp in ping replies.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_redirects": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable ICMP redirects.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_redirects_ipv6": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable IPV6 ICMP redirects.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"radius_options_attributes_nas_id": schema.StringAttribute{
				Optional:    true,
				Description: "Value of NAS-ID in outgoing RADIUS packets.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 64),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"radius_options_attributes_nas_ipaddress": schema.StringAttribute{
				Optional:    true,
				Description: "Value of NAS-IP-Address in outgoing RADIUS packets.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
			"radius_options_enhanced_accounting": schema.BoolAttribute{
				Optional:    true,
				Description: "Include authentication method, remote port and user-privileges in `login` accounting.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"radius_options_password_protocol_mschapv2": schema.BoolAttribute{
				Optional:    true,
				Description: "MSCHAP version 2 for password protocol used in RADIUS packets.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"tacplus_options_authorization_time_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "TACACS+ authorization refresh time interval (minutes).",
				Validators: []validator.Int64{
					int64validator.Between(15, 1440),
				},
			},
			"tacplus_options_enhanced_accounting": schema.BoolAttribute{
				Optional:    true,
				Description: "Include authentication method, remote port and user-privileges in `login` accounting.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"tacplus_options_exclude_cmd_attribute": schema.BoolAttribute{
				Optional:    true,
				Description: "In start/stop requests, do not include `cmd` attribute.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"tacplus_options_no_cmd_attribute_value": schema.BoolAttribute{
				Optional:    true,
				Description: "In start/stop requests, set `cmd` attribute value to empty string.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"tacplus_options_service_name": schema.StringAttribute{
				Optional:    true,
				Description: "TACACS+ service name.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"tacplus_options_strict_authorization": schema.BoolAttribute{
				Optional:    true,
				Description: "Deny login if authorization request fails.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"tacplus_options_no_strict_authorization": schema.BoolAttribute{
				Optional:    true,
				Description: "Don't deny login if authorization request fails.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"tacplus_options_timestamp_and_timezone": schema.BoolAttribute{
				Optional:    true,
				Description: "In start/stop accounting packets, include `start-time`, `stop-time` and `timezone` attributes.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"time_zone": schema.StringAttribute{
				Optional:    true,
				Description: "Time zone name or POSIX-compliant time zone string.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
				},
			},
			"tracing_dest_override_syslog_host": schema.StringAttribute{
				Optional:    true,
				Description: "Send trace messages to remote syslog server.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress().IPv4Only(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"accounting": schema.SingleNestedBlock{
				Description: "Declare `accounting` configuration.",
				Attributes: map[string]schema.Attribute{
					"events": schema.SetAttribute{
						ElementType: types.StringType,
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Events to be logged.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf("change-log", "interactive-commands", "login"),
							),
						},
					},
					"destination_radius": schema.BoolAttribute{
						Optional:    true,
						Description: "Send RADIUS accounting records.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"destination_tacplus": schema.BoolAttribute{
						Optional:    true,
						Description: "Send TACACS+ accounting records.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"enhanced_avs_max": schema.Int64Attribute{
						Optional:    true,
						Description: "No. of AV pairs each of which can store a max of 250 Bytes.",
						Validators: []validator.Int64{
							int64validator.Between(7, 15),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"destination_radius_server": schema.ListNestedBlock{
						Description: "For each address, RADIUS accounting server configuration.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"address": schema.StringAttribute{
									Required:    true,
									Description: "RADIUS server address.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress(),
									},
								},
								"secret": schema.StringAttribute{
									Required:    true,
									Sensitive:   true,
									Description: "Shared secret with the RADIUS server.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"accounting_port": schema.Int64Attribute{
									Optional:    true,
									Description: "RADIUS server accounting port number.",
									Validators: []validator.Int64{
										int64validator.Between(1, 65535),
									},
								},
								"accounting_retry": schema.Int64Attribute{
									Optional:    true,
									Description: "Accounting retry attempts.",
									Validators: []validator.Int64{
										int64validator.Between(0, 100),
									},
								},
								"accounting_timeout": schema.Int64Attribute{
									Optional:    true,
									Description: "Accounting request timeout period.",
									Validators: []validator.Int64{
										int64validator.Between(0, 1000),
									},
								},
								"dynamic_request_port": schema.Int64Attribute{
									Optional:    true,
									Description: "RADIUS client dynamic request port number.",
									Validators: []validator.Int64{
										int64validator.Between(1, 65535),
									},
								},
								"max_outstanding_requests": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum requests in flight to server.",
									Validators: []validator.Int64{
										int64validator.Between(0, 2000),
									},
								},
								"port": schema.Int64Attribute{
									Optional:    true,
									Description: "RADIUS server authentication port number.",
									Validators: []validator.Int64{
										int64validator.Between(1, 65535),
									},
								},
								"preauthentication_port": schema.Int64Attribute{
									Optional:    true,
									Description: "RADIUS server preauthentication port number.",
									Validators: []validator.Int64{
										int64validator.Between(1, 65535),
									},
								},
								"preauthentication_secret": schema.StringAttribute{
									Optional:    true,
									Sensitive:   true,
									Description: "Shared secret with the RADIUS server.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"retry": schema.Int64Attribute{
									Optional:    true,
									Description: "Retry attempts.",
									Validators: []validator.Int64{
										int64validator.Between(1, 100),
									},
								},
								"routing_instance": schema.StringAttribute{
									Optional:    true,
									Description: "Routing instance.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 63),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									},
								},
								"source_address": schema.StringAttribute{
									Optional:    true,
									Description: "Use specified address as source address.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress(),
									},
								},
								"timeout": schema.Int64Attribute{
									Optional:    true,
									Description: "Request timeout period.",
									Validators: []validator.Int64{
										int64validator.Between(1, 1000),
									},
								},
							},
						},
					},
					"destination_tacplus_server": schema.ListNestedBlock{
						Description: "For each address, TACACS+ accounting server configuration.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"address": schema.StringAttribute{
									Required:    true,
									Description: "TACACS+ authentication server address.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress(),
									},
								},
								"port": schema.Int64Attribute{
									Optional:    true,
									Description: "TACACS+ authentication server port number.",
									Validators: []validator.Int64{
										int64validator.Between(1, 65535),
									},
								},
								"routing_instance": schema.StringAttribute{
									Optional:    true,
									Description: "Routing instance.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 63),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
									},
								},
								"secret": schema.StringAttribute{
									Optional:    true,
									Sensitive:   true,
									Description: "Shared secret with the authentication server.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"single_connection": schema.BoolAttribute{
									Optional:    true,
									Description: "Optimize TCP connection attempts.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"source_address": schema.StringAttribute{
									Optional:    true,
									Description: "Use specified address as source address.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress(),
									},
								},
								"timeout": schema.Int64Attribute{
									Optional:    true,
									Description: "Request timeout period.",
									Validators: []validator.Int64{
										int64validator.Between(1, 90),
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
			"archival_configuration": schema.SingleNestedBlock{
				Description: "Declare `archival configuration` configuration.",
				Attributes: map[string]schema.Attribute{
					"transfer_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Frequency at which file transfer happens (minutes).",
						Validators: []validator.Int64{
							int64validator.Between(15, 2880),
						},
					},
					"transfer_on_commit": schema.BoolAttribute{
						Optional:    true,
						Description: "Transfer after each commit.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"archive_site": schema.ListNestedBlock{
						Description: "For each url, configure archive-site destination.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"url": schema.StringAttribute{
									Required:    true,
									Description: "URLs to receive configuration files.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
										tfvalidator.StringSpaceExclusion(),
									},
								},
								"password": schema.StringAttribute{
									Optional:    true,
									Sensitive:   true,
									Description: "Password for login into the archive site.",
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
			"inet6_backup_router": schema.SingleNestedBlock{
				Description: "Declare `inet6-backup-router` configuration.",
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Address of router to use while booting.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress().IPv6Only(),
						},
					},
					"destination": schema.SetAttribute{
						ElementType: types.StringType,
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Destination networks reachable through the router.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(
								tfvalidator.StringCIDR().IPv6Only(),
							),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"internet_options": schema.SingleNestedBlock{
				Description: "Declare `internet-options` configuration.",
				Attributes: map[string]schema.Attribute{
					"gre_path_mtu_discovery": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable path MTU discovery for GRE tunnels.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_gre_path_mtu_discovery": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't enable path MTU discovery for GRE tunnels.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ipip_path_mtu_discovery": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable path MTU discovery for IP-IP tunnels.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_ipip_path_mtu_discovery": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't enable path MTU discovery for IP-IP tunnels.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ipv6_duplicate_addr_detection_transmits": schema.Int64Attribute{
						Optional:    true,
						Description: "IPv6 Duplicate address detection transmits.",
						Validators: []validator.Int64{
							int64validator.Between(0, 20),
						},
					},
					"ipv6_path_mtu_discovery": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable IPv6 Path MTU discovery.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_ipv6_path_mtu_discovery": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't enable IPv6 Path MTU discovery.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ipv6_path_mtu_discovery_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "IPv6 Path MTU Discovery timeout (minutes).",
						Validators: []validator.Int64{
							int64validator.Between(5, 71582788),
						},
					},
					"ipv6_reject_zero_hop_limit": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable dropping IPv6 packets with zero hop-limit.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_ipv6_reject_zero_hop_limit": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't enable dropping IPv6 packets with zero hop-limit.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_tcp_reset": schema.StringAttribute{
						Optional:    true,
						Description: "Do not send RST TCP packet for packets sent to non-listening ports.",
						Validators: []validator.String{
							stringvalidator.OneOf("drop-all-tcp", "drop-tcp-with-syn-only"),
						},
					},
					"no_tcp_rfc1323": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable RFC 1323 TCP extensions.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_tcp_rfc1323_paws": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable RFC 1323 Protection Against Wrapped Sequence Number extension.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"path_mtu_discovery": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable Path MTU discovery on TCP connections.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_path_mtu_discovery": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't enable Path MTU discovery on TCP connections.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"source_port_upper_limit": schema.Int64Attribute{
						Optional:    true,
						Description: "Specify upper limit of source port selection range.",
						Validators: []validator.Int64{
							int64validator.Between(5000, 65535),
						},
					},
					"source_quench": schema.BoolAttribute{
						Optional:    true,
						Description: "React to incoming ICMP Source Quench messages.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_source_quench": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't react to incoming ICMP Source Quench messages.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"tcp_drop_synfin_set": schema.BoolAttribute{
						Optional:    true,
						Description: "Drop TCP packets that have both SYN and FIN flags.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"tcp_mss": schema.Int64Attribute{
						Optional:    true,
						Description: " Maximum value of TCP MSS for IPV4 traffic (bytes).",
						Validators: []validator.Int64{
							int64validator.Between(64, 65535),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"icmpv4_rate_limit": schema.SingleNestedBlock{
						Description: "Declare `icmpv4-rate-limit` configuration.",
						Attributes: map[string]schema.Attribute{
							"bucket_size": schema.Int64Attribute{
								Optional:    true,
								Description: "ICMP rate-limiting maximum bucket size (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"packet_rate": schema.Int64Attribute{
								Optional:    true,
								Description: "ICMP rate-limiting packets earned per second.",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"icmpv6_rate_limit": schema.SingleNestedBlock{
						Description: "Declare `icmpv6-rate-limit` configuration.",
						Attributes: map[string]schema.Attribute{
							"bucket_size": schema.Int64Attribute{
								Optional:    true,
								Description: "ICMPv6 rate-limiting maximum bucket size (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"packet_rate": schema.Int64Attribute{
								Optional:    true,
								Description: "ICMPv6 rate-limiting packets earned per second.",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
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
			"license": schema.SingleNestedBlock{
				Description: "Declare `license` configuration.",
				Attributes: map[string]schema.Attribute{
					"autoupdate": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable autoupdate license keys.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"autoupdate_password": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Password for autoupdate license keys from license servers.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"autoupdate_url": schema.StringAttribute{
						Optional:    true,
						Description: "Url for autoupdate license keys from license servers.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 250),
							tfvalidator.StringDoubleQuoteExclusion(),
							tfvalidator.StringSpaceExclusion(),
						},
					},
					"renew_before_expiration": schema.Int64Attribute{
						Optional:    true,
						Description: "License renewal lead time before expiration, in days.",
						Validators: []validator.Int64{
							int64validator.Between(0, 60),
						},
					},
					"renew_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "License checking interval, in hours.",
						Validators: []validator.Int64{
							int64validator.Between(1, 336),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"login": schema.SingleNestedBlock{
				Description: "Declare `login` configuration.",
				Attributes: map[string]schema.Attribute{
					"announcement": schema.StringAttribute{
						Optional:    true,
						Description: "System announcement message (displayed after login).",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 2048),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"deny_sources_address": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Sources from which logins are denied.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(1),
								stringvalidator.Any(
									tfvalidator.StringCIDR(),
									tfvalidator.StringIPAddress(),
								),
							),
						},
					},
					"idle_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum idle time before logout (minutes).",
						Validators: []validator.Int64{
							int64validator.Between(1, 60),
						},
					},
					"message": schema.StringAttribute{
						Optional:    true,
						Description: "System login message.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"password": schema.SingleNestedBlock{
						Description: "eclare `password` configuration.",
						Attributes: map[string]schema.Attribute{
							"change_type": schema.StringAttribute{
								Optional:    true,
								Description: "Password change type.",
								Validators: []validator.String{
									stringvalidator.OneOf("character-sets", "set-transitions"),
								},
							},
							"format": schema.StringAttribute{
								Optional:    true,
								Description: "Encryption method to use for password.",
								Validators: []validator.String{
									stringvalidator.OneOf("sha1", "sha256", "sha512"),
								},
							},
							"maximum_length": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum password length for all users.",
								Validators: []validator.Int64{
									int64validator.Between(20, 128),
								},
							},
							"minimum_changes": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum number of changes in password.",
								Validators: []validator.Int64{
									int64validator.Between(1, 128),
								},
							},
							"minimum_character_changes": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum number of character changes between old and new passwords.",
								Validators: []validator.Int64{
									int64validator.Between(4, 15),
								},
							},
							"minimum_length": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum password length for all users.",
								Validators: []validator.Int64{
									int64validator.Between(6, 20),
								},
							},
							"minimum_lower_cases": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum number of lower-case class characters in password.",
								Validators: []validator.Int64{
									int64validator.Between(1, 128),
								},
							},
							"minimum_numerics": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum number of numeric class characters in password.",
								Validators: []validator.Int64{
									int64validator.Between(1, 128),
								},
							},
							"minimum_punctuations": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum number of punctuation class characters in password.",
								Validators: []validator.Int64{
									int64validator.Between(1, 128),
								},
							},
							"minimum_reuse": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum number of old passwords which should not be same as the new password.",
								Validators: []validator.Int64{
									int64validator.Between(1, 20),
								},
							},
							"minimum_upper_cases": schema.Int64Attribute{
								Optional:    true,
								Description: "Minimum number of upper-case class characters in password.",
								Validators: []validator.Int64{
									int64validator.Between(1, 128),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"retry_options": schema.SingleNestedBlock{
						Description: "Declare `retry-options` configuration.",
						Attributes: map[string]schema.Attribute{
							"backoff_factor": schema.Int64Attribute{
								Optional:    true,
								Description: "Delay factor after `backoff-threshold` password failures.",
								Validators: []validator.Int64{
									int64validator.Between(5, 10),
								},
							},
							"backoff_threshold": schema.Int64Attribute{
								Optional:    true,
								Description: "Number of password failures before delay is introduced.",
								Validators: []validator.Int64{
									int64validator.Between(1, 3),
								},
							},
							"lockout_period": schema.Int64Attribute{
								Optional:    true,
								Description: "Amount of time user account is locked after `tries_before_disconnect` failures (minutes).",
								Validators: []validator.Int64{
									int64validator.Between(1, 43200),
								},
							},
							"maximum_time": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum time the connection will remain for user to enter username and password.",
								Validators: []validator.Int64{
									int64validator.Between(20, 300),
								},
							},
							"minimum_time": schema.Int64Attribute{
								Optional:    true,
								Description: " Minimum total connection time if all attempts fail.",
								Validators: []validator.Int64{
									int64validator.Between(20, 60),
								},
							},
							"tries_before_disconnect": schema.Int64Attribute{
								Optional:    true,
								Description: "Number of times user is allowed to try password.",
								Validators: []validator.Int64{
									int64validator.Between(2, 10),
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
			"name_server_opts": schema.ListNestedBlock{
				Description: "DNS name servers with optional options.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Required:    true,
							Description: "Address of the name server.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress(),
							},
						},
						"routing_instance": schema.StringAttribute{
							Optional:    true,
							Description: "Routing instance through which the name server is reachable.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
					},
				},
			},
			"ntp": schema.SingleNestedBlock{
				Description: "Declare `ntp` configuration.",
				Attributes: map[string]schema.Attribute{
					"boot_server": schema.StringAttribute{
						Optional:    true,
						Description: "Server to query during boot sequence.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"broadcast_client": schema.BoolAttribute{
						Optional:    true,
						Description: "Listen to broadcast NTP.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"interval_range": schema.Int64Attribute{
						Optional:    true,
						Description: "Set the minpoll and maxpoll interval range.",
						Validators: []validator.Int64{
							int64validator.Between(0, 3),
						},
					},
					"multicast_client": schema.BoolAttribute{
						Optional:    true,
						Description: "Listen to multicast NTP.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"multicast_client_address": schema.StringAttribute{
						Optional:    true,
						Description: "Multicast address to listen to.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"threshold_action": schema.StringAttribute{
						Optional:    true,
						Description: "Select actions for NTP abnormal adjustment.",
						Validators: []validator.String{
							stringvalidator.OneOf("accept", "reject"),
						},
					},
					"threshold_value": schema.Int64Attribute{
						Optional:    true,
						Description: "Set the maximum threshold(sec) allowed for NTP adjustment.",
						Validators: []validator.Int64{
							int64validator.Between(1, 600),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"ports": schema.SingleNestedBlock{
				Description: "Declare `ports` configuration.",
				Attributes: map[string]schema.Attribute{
					"auxiliary_authentication_order": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Order in which authentication methods are invoked on auxiliary port.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.ValueStringsAre(
								stringvalidator.OneOf("password", "radius", "tacplus"),
							),
						},
					},
					"auxiliary_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable console on auxiliary port.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"auxiliary_insecure": schema.BoolAttribute{
						Optional:    true,
						Description: "Disallow superuser access on auxiliary port.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"auxiliary_logout_on_disconnect": schema.BoolAttribute{
						Optional:    true,
						Description: "Log out the console session when cable is unplugged.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"auxiliary_type": schema.StringAttribute{
						Optional:    true,
						Description: "Terminal type on auxiliary port.",
						Validators: []validator.String{
							stringvalidator.OneOf("ansi", "small-xterm", "vt100", "xterm"),
						},
					},
					"console_authentication_order": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Order in which authentication methods are invoked on console port.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.ValueStringsAre(
								stringvalidator.OneOf("password", "radius", "tacplus"),
							),
						},
					},
					"console_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable console on console port.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"console_insecure": schema.BoolAttribute{
						Optional:    true,
						Description: "Disallow superuser access on console port.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"console_logout_on_disconnect": schema.BoolAttribute{
						Optional:    true,
						Description: "Log out the console session when cable is unplugged.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"console_type": schema.StringAttribute{
						Optional:    true,
						Description: "Terminal type on console port.",
						Validators: []validator.String{
							stringvalidator.OneOf("ansi", "small-xterm", "vt100", "xterm"),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"services": schema.SingleNestedBlock{
				Description: "Declare `services` configuration.",
				Attributes: map[string]schema.Attribute{
					"web_management_session_idle_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Default timeout of web-management sessions (minutes).",
						Validators: []validator.Int64{
							int64validator.Between(1, 1440),
						},
					},
					"web_management_session_limit": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum number of web-management sessions to allow.",
						Validators: []validator.Int64{
							int64validator.Between(1, 1024),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"netconf_ssh": schema.SingleNestedBlock{
						Description: "Declare `netconf ssh` configuration.",
						Attributes: map[string]schema.Attribute{
							"client_alive_count_max": schema.Int64Attribute{
								Optional:    true,
								Description: "Threshold of missing client-alive responses that triggers a disconnect.",
								Validators: []validator.Int64{
									int64validator.Between(0, 255),
								},
							},
							"client_alive_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Frequency of client-alive requests (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 65535),
								},
							},
							"connection_limit": schema.Int64Attribute{
								Optional:    true,
								Description: "Limit number of simultaneous connections (connections).",
								Validators: []validator.Int64{
									int64validator.Between(1, 250),
								},
							},
							"rate_limit": schema.Int64Attribute{
								Optional:    true,
								Description: "Limit incoming connection rate (connections per minute).",
								Validators: []validator.Int64{
									int64validator.Between(1, 250),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"netconf_traceoptions": schema.SingleNestedBlock{
						Description: "Declare `netconf traceoptions` configuration.",
						Attributes: map[string]schema.Attribute{
							"file_name": schema.StringAttribute{
								Optional:    true,
								Description: "Name of file in which to write trace information.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
									tfvalidator.StringSpaceExclusion(),
									tfvalidator.StringRuneExclusion('/', '%'),
								},
							},
							"file_files": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum number of trace files.",
								Validators: []validator.Int64{
									int64validator.Between(2, 1000),
								},
							},
							"file_match": schema.StringAttribute{
								Optional:    true,
								Description: "Regular expression for lines to be logged.",
								Validators: []validator.String{
									tfvalidator.StringRegex(),
								},
							},
							"file_size": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum trace file size.",
								Validators: []validator.Int64{
									int64validator.Between(10240, 1073741824),
								},
							},
							"file_world_readable": schema.BoolAttribute{
								Optional:    true,
								Description: "Allow any user to read the log file.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"file_no_world_readable": schema.BoolAttribute{
								Optional:    true,
								Description: "Don't allow any user to read the log file.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"flag": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Tracing parameters.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.OneOf("all", "debug", "incoming", "outgoing"),
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
							"on_demand": schema.BoolAttribute{
								Optional:    true,
								Description: "Enable on-demand tracing.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"ssh": schema.SingleNestedBlock{
						Description: "Declare `ssh` configuration.",
						Attributes: map[string]schema.Attribute{
							"authentication_order": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Order in which authentication methods are invoked.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.ValueStringsAre(
										stringvalidator.OneOf("password", "radius", "tacplus"),
									),
								},
							},
							"ciphers": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Specify the ciphers allowed for protocol version 2.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.AlgorithmFormat),
									),
								},
							},
							"client_alive_count_max": schema.Int64Attribute{
								Optional:    true,
								Description: "Threshold of missing client-alive responses that triggers a disconnect.",
								Validators: []validator.Int64{
									int64validator.Between(0, 255),
								},
							},
							"client_alive_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Frequency of client-alive requests (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 65535),
								},
							},
							"connection_limit": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum number of allowed connections.",
								Validators: []validator.Int64{
									int64validator.Between(1, 250),
								},
							},
							"fingerprint_hash": schema.StringAttribute{
								Optional:    true,
								Description: "Configure hash algorithm used when displaying key fingerprints.",
								Validators: []validator.String{
									stringvalidator.OneOf("md5", "sha2-256"),
								},
							},
							"hostkey_algorithm": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Specify permissible SSH host-key algorithms.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.AlgorithmFormat),
									),
								},
							},
							"key_exchange": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Specify ssh key-exchange for Diffie-Hellman keys.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.AlgorithmFormat),
									),
								},
							},
							"log_key_changes": schema.BoolAttribute{
								Optional:    true,
								Description: "Log changes to authorized keys to syslog.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"macs": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Message Authentication Code algorithms allowed (SSHv2).",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.AlgorithmFormat),
									),
								},
							},
							"max_pre_authentication_packets": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum number of pre-authentication SSH packets per single SSH connection.",
								Validators: []validator.Int64{
									int64validator.Between(20, 2147483647),
								},
							},
							"max_sessions_per_connection": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum number of sessions per single SSH connection.",
								Validators: []validator.Int64{
									int64validator.Between(1, 65535),
								},
							},
							"no_passwords": schema.BoolAttribute{
								Optional:    true,
								Description: "Disables ssh password based authentication.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_public_keys": schema.BoolAttribute{
								Optional:    true,
								Description: "Disables ssh public key based authentication.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"port": schema.Int64Attribute{
								Optional:    true,
								Description: "Port number to accept incoming connections.",
								Validators: []validator.Int64{
									int64validator.Between(1, 65535),
								},
							},
							"protocol_version": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Specify ssh protocol versions supported.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.OneOf("v1", "v2"),
									),
								},
							},
							"rate_limit": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum number of connections per minute.",
								Validators: []validator.Int64{
									int64validator.Between(1, 250),
								},
							},
							"root_login": schema.StringAttribute{
								Optional:    true,
								Description: "Configure root access via ssh.",
								Validators: []validator.String{
									stringvalidator.OneOf("allow", "deny", "deny-password"),
								},
							},
							"tcp_forwarding": schema.BoolAttribute{
								Optional:    true,
								Description: "Allow forwarding TCP connections via SSH.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_tcp_forwarding": schema.BoolAttribute{
								Optional:    true,
								Description: "Do not allow forwarding TCP connections via SSH.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"web_management_http": schema.SingleNestedBlock{
						Description: "Enable `web-management http`.",
						Attributes: map[string]schema.Attribute{
							"interface": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Specify the name of one or more interfaces.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
										tfvalidator.String1DotCount(),
									),
								},
							},
							"port": schema.Int64Attribute{
								Optional:    true,
								Description: "Port number to connect to HTTP service.",
								Validators: []validator.Int64{
									int64validator.Between(1, 65535),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"web_management_https": schema.SingleNestedBlock{
						Description: "Declare `web-management https` configuration.",
						Attributes: map[string]schema.Attribute{
							"interface": schema.SetAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Description: "Specify the name of one or more interfaces.",
								Validators: []validator.Set{
									setvalidator.SizeAtLeast(1),
									setvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
										tfvalidator.String1DotCount(),
									),
								},
							},
							"local_certificate": schema.StringAttribute{
								Optional:    true,
								Description: "Specify the name of the certificate.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"pki_local_certificate": schema.StringAttribute{
								Optional:    true,
								Description: "Specify the name of the certificate that is generated by the PKI and authenticated by a CA.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"port": schema.Int64Attribute{
								Optional:    true,
								Description: "Port number to connect to HTTPS service.",
								Validators: []validator.Int64{
									int64validator.Between(1, 65535),
								},
							},
							"system_generated_certificate": schema.BoolAttribute{
								Optional:    true,
								Description: "Will automatically generate a self-signed certificate.",
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
			"syslog": schema.SingleNestedBlock{
				Description: "Declare `syslog` configuration.",
				Attributes: map[string]schema.Attribute{
					"log_rotate_frequency": schema.Int64Attribute{
						Optional:    true,
						Description: "Rotate log frequency (minutes).",
						Validators: []validator.Int64{
							int64validator.Between(1, 59),
						},
					},
					"source_address": schema.StringAttribute{
						Optional:    true,
						Description: "Use specified address as source address.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"time_format_millisecond": schema.BoolAttribute{
						Optional:    true,
						Description: "Include milliseconds in system log timestamp.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"time_format_year": schema.BoolAttribute{
						Optional:    true,
						Description: "Include year in system log timestamp.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"archive": schema.SingleNestedBlock{
						Description: "Declare `archive` configuration.",
						Attributes: map[string]schema.Attribute{
							"binary_data": schema.BoolAttribute{
								Optional:    true,
								Description: "ark file as if it contains binary data.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_binary_data": schema.BoolAttribute{
								Optional:    true,
								Description: "Don't mark file as if it contains binary data.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"files": schema.Int64Attribute{
								Optional:    true,
								Description: "Number of files to be archived.",
								Validators: []validator.Int64{
									int64validator.Between(1, 1000),
								},
							},
							"size": schema.Int64Attribute{
								Optional:    true,
								Description: "Size of files to be archived.",
								Validators: []validator.Int64{
									int64validator.Between(65536, 1073741824),
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
					"console": schema.SingleNestedBlock{
						Description: "Declare `console` configuration.",
						Attributes: map[string]schema.Attribute{
							"any_severity": schema.StringAttribute{
								Optional:    true,
								Description: "All facilities severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"authorization_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Authorization system severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"changelog_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Configuration change log severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"conflictlog_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Configuration conflict log severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"daemon_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Various system processes severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"dfc_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Dynamic flow capture severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"external_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Local external applications severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"firewall_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Firewall filtering system severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"ftp_severity": schema.StringAttribute{
								Optional:    true,
								Description: "FTP process severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"interactivecommands_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Commands executed by the UI severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"kernel_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Kernel severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"ntp_severity": schema.StringAttribute{
								Optional:    true,
								Description: "NTP process severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"pfe_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Packet Forwarding Engine severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"security_severity": schema.StringAttribute{
								Optional:    true,
								Description: "Security related severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
								},
							},
							"user_severity": schema.StringAttribute{
								Optional:    true,
								Description: "User processes severity.",
								Validators: []validator.String{
									stringvalidator.OneOf(junos.SyslogSeverity()...),
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

//nolint:lll
type systemData struct {
	ID                                      types.String                      `tfsdk:"id"`
	AuthenticationOrder                     []types.String                    `tfsdk:"authentication_order"`
	AutoSnapshot                            types.Bool                        `tfsdk:"auto_snapshot"`
	DefaultAddressSelection                 types.Bool                        `tfsdk:"default_address_selection"`
	DomainName                              types.String                      `tfsdk:"domain_name"`
	HostName                                types.String                      `tfsdk:"host_name"`
	MaxConfigurationRollbacks               types.Int64                       `tfsdk:"max_configuration_rollbacks"`
	MaxConfigurationsOnFlash                types.Int64                       `tfsdk:"max_configurations_on_flash"`
	NameServer                              []types.String                    `tfsdk:"name_server"`
	NoMulticastEcho                         types.Bool                        `tfsdk:"no_multicast_echo"`
	NoPingRecordRoute                       types.Bool                        `tfsdk:"no_ping_record_route"`
	NoPingTimestamp                         types.Bool                        `tfsdk:"no_ping_time_stamp"`
	NoRedirects                             types.Bool                        `tfsdk:"no_redirects"`
	NoRedirectsIPv6                         types.Bool                        `tfsdk:"no_redirects_ipv6"`
	RadiusOptionsAttributesNasID            types.String                      `tfsdk:"radius_options_attributes_nas_id"`
	RadiusOptionsAttributesNasIpaddress     types.String                      `tfsdk:"radius_options_attributes_nas_ipaddress"`
	RadiusOptionsEnhancedAccounting         types.Bool                        `tfsdk:"radius_options_enhanced_accounting"`
	RadiusOptionsPasswordProtocolMschapv2   types.Bool                        `tfsdk:"radius_options_password_protocol_mschapv2"`
	TacplusOptionsAuthorizationTimeInterval types.Int64                       `tfsdk:"tacplus_options_authorization_time_interval"`
	TacplusOptionsEnhancedAccounting        types.Bool                        `tfsdk:"tacplus_options_enhanced_accounting"`
	TacplusOptionsExcludeCmdAttribute       types.Bool                        `tfsdk:"tacplus_options_exclude_cmd_attribute"`
	TacplusOptionsNoCmdAttributeValue       types.Bool                        `tfsdk:"tacplus_options_no_cmd_attribute_value"`
	TacplusOptionsServiceName               types.String                      `tfsdk:"tacplus_options_service_name"`
	TacplusOptionsStrictAuthorization       types.Bool                        `tfsdk:"tacplus_options_strict_authorization"`
	TacplusOptionsNoStrictAuthorization     types.Bool                        `tfsdk:"tacplus_options_no_strict_authorization"`
	TacplusOptionsTimestampAndTimezone      types.Bool                        `tfsdk:"tacplus_options_timestamp_and_timezone"`
	TimeZone                                types.String                      `tfsdk:"time_zone"`
	TracingDestOverrideSyslogHost           types.String                      `tfsdk:"tracing_dest_override_syslog_host"`
	Accounting                              *systemBlockAccounting            `tfsdk:"accounting"`
	ArchivalConfiguration                   *systemBlockArchivalConfiguration `tfsdk:"archival_configuration"`
	Inet6BackupRouter                       *systemBlockInet6BackupRouter     `tfsdk:"inet6_backup_router"`
	InternetOptions                         *systemBlockInternetOptions       `tfsdk:"internet_options"`
	License                                 *systemBlockLicense               `tfsdk:"license"`
	Login                                   *systemBlockLogin                 `tfsdk:"login"`
	NameServerOpts                          []systemBlockNameServerOpts       `tfsdk:"name_server_opts"`
	Ntp                                     *systemBlockNtp                   `tfsdk:"ntp"`
	Ports                                   *systemBlockPorts                 `tfsdk:"ports"`
	Services                                *systemBlockServices              `tfsdk:"services"`
	Syslog                                  *systemBlockSyslog                `tfsdk:"syslog"`
}

//nolint:lll
type systemConfig struct {
	ID                                      types.String                            `tfsdk:"id"`
	AuthenticationOrder                     types.List                              `tfsdk:"authentication_order"`
	AutoSnapshot                            types.Bool                              `tfsdk:"auto_snapshot"`
	DefaultAddressSelection                 types.Bool                              `tfsdk:"default_address_selection"`
	DomainName                              types.String                            `tfsdk:"domain_name"`
	HostName                                types.String                            `tfsdk:"host_name"`
	MaxConfigurationRollbacks               types.Int64                             `tfsdk:"max_configuration_rollbacks"`
	MaxConfigurationsOnFlash                types.Int64                             `tfsdk:"max_configurations_on_flash"`
	NameServer                              types.List                              `tfsdk:"name_server"`
	NoMulticastEcho                         types.Bool                              `tfsdk:"no_multicast_echo"`
	NoPingRecordRoute                       types.Bool                              `tfsdk:"no_ping_record_route"`
	NoPingTimestamp                         types.Bool                              `tfsdk:"no_ping_time_stamp"`
	NoRedirects                             types.Bool                              `tfsdk:"no_redirects"`
	NoRedirectsIPv6                         types.Bool                              `tfsdk:"no_redirects_ipv6"`
	RadiusOptionsAttributesNasID            types.String                            `tfsdk:"radius_options_attributes_nas_id"`
	RadiusOptionsAttributesNasIpaddress     types.String                            `tfsdk:"radius_options_attributes_nas_ipaddress"`
	RadiusOptionsEnhancedAccounting         types.Bool                              `tfsdk:"radius_options_enhanced_accounting"`
	RadiusOptionsPasswordProtocolMschapv2   types.Bool                              `tfsdk:"radius_options_password_protocol_mschapv2"`
	TacplusOptionsAuthorizationTimeInterval types.Int64                             `tfsdk:"tacplus_options_authorization_time_interval"`
	TacplusOptionsEnhancedAccounting        types.Bool                              `tfsdk:"tacplus_options_enhanced_accounting"`
	TacplusOptionsExcludeCmdAttribute       types.Bool                              `tfsdk:"tacplus_options_exclude_cmd_attribute"`
	TacplusOptionsNoCmdAttributeValue       types.Bool                              `tfsdk:"tacplus_options_no_cmd_attribute_value"`
	TacplusOptionsServiceName               types.String                            `tfsdk:"tacplus_options_service_name"`
	TacplusOptionsStrictAuthorization       types.Bool                              `tfsdk:"tacplus_options_strict_authorization"`
	TacplusOptionsNoStrictAuthorization     types.Bool                              `tfsdk:"tacplus_options_no_strict_authorization"`
	TacplusOptionsTimestampAndTimezone      types.Bool                              `tfsdk:"tacplus_options_timestamp_and_timezone"`
	TimeZone                                types.String                            `tfsdk:"time_zone"`
	TracingDestOverrideSyslogHost           types.String                            `tfsdk:"tracing_dest_override_syslog_host"`
	Accounting                              *systemBlockAccountingConfig            `tfsdk:"accounting"`
	ArchivalConfiguration                   *systemBlockArchivalConfigurationConfig `tfsdk:"archival_configuration"`
	Inet6BackupRouter                       *systemBlockInet6BackupRouterConfig     `tfsdk:"inet6_backup_router"`
	InternetOptions                         *systemBlockInternetOptions             `tfsdk:"internet_options"`
	License                                 *systemBlockLicense                     `tfsdk:"license"`
	Login                                   *systemBlockLoginConfig                 `tfsdk:"login"`
	NameServerOpts                          types.List                              `tfsdk:"name_server_opts"`
	Ntp                                     *systemBlockNtp                         `tfsdk:"ntp"`
	Ports                                   *systemBlockPortsConfig                 `tfsdk:"ports"`
	Services                                *systemBlockServicesConfig              `tfsdk:"services"`
	Syslog                                  *systemBlockSyslog                      `tfsdk:"syslog"`
}

type systemBlockAccounting struct {
	Events                   []types.String                                       `tfsdk:"events"`
	DestinationRadius        types.Bool                                           `tfsdk:"destination_radius"`
	DestinationTacplus       types.Bool                                           `tfsdk:"destination_tacplus"`
	EnhancedAvsMax           types.Int64                                          `tfsdk:"enhanced_avs_max"`
	DestinationRadiusServer  []systemBlockAccountingBlockDestinationRadiusServer  `tfsdk:"destination_radius_server"`
	DestinationTacplusServer []systemBlockAccountingBlockDestinationTacplusServer `tfsdk:"destination_tacplus_server"`
}

type systemBlockAccountingConfig struct {
	Events                   types.Set   `tfsdk:"events"`
	DestinationRadius        types.Bool  `tfsdk:"destination_radius"`
	DestinationTacplus       types.Bool  `tfsdk:"destination_tacplus"`
	EnhancedAvsMax           types.Int64 `tfsdk:"enhanced_avs_max"`
	DestinationRadiusServer  types.List  `tfsdk:"destination_radius_server"`
	DestinationTacplusServer types.List  `tfsdk:"destination_tacplus_server"`
}

type systemBlockAccountingBlockDestinationRadiusServer struct {
	Address                 types.String `tfsdk:"address"`
	Secret                  types.String `tfsdk:"secret"`
	AccountingPort          types.Int64  `tfsdk:"accounting_port"`
	AccountingRetry         types.Int64  `tfsdk:"accounting_retry"`
	AccountingTimeout       types.Int64  `tfsdk:"accounting_timeout"`
	DynamicRequestPort      types.Int64  `tfsdk:"dynamic_request_port"`
	MaxOutstandingRequests  types.Int64  `tfsdk:"max_outstanding_requests"`
	Port                    types.Int64  `tfsdk:"port"`
	PreauthenticationPort   types.Int64  `tfsdk:"preauthentication_port"`
	PreauthenticationSecret types.String `tfsdk:"preauthentication_secret"`
	Retry                   types.Int64  `tfsdk:"retry"`
	RoutingInstance         types.String `tfsdk:"routing_instance"`
	SourceAddress           types.String `tfsdk:"source_address"`
	Timeout                 types.Int64  `tfsdk:"timeout"`
}

type systemBlockAccountingBlockDestinationTacplusServer struct {
	Address          types.String `tfsdk:"address"`
	Port             types.Int64  `tfsdk:"port"`
	RoutingInstance  types.String `tfsdk:"routing_instance"`
	Secret           types.String `tfsdk:"secret"`
	SingleConnection types.Bool   `tfsdk:"single_connection"`
	SourceAddress    types.String `tfsdk:"source_address"`
	Timeout          types.Int64  `tfsdk:"timeout"`
}

type systemBlockArchivalConfiguration struct {
	TransferInterval types.Int64                                        `tfsdk:"transfer_interval"`
	TransferOnCommit types.Bool                                         `tfsdk:"transfer_on_commit"`
	ArchiveSite      []systemBlockArchivalConfigurationBlockArchiveSite `tfsdk:"archive_site"`
}

type systemBlockArchivalConfigurationConfig struct {
	TransferInterval types.Int64 `tfsdk:"transfer_interval"`
	TransferOnCommit types.Bool  `tfsdk:"transfer_on_commit"`
	ArchiveSite      types.List  `tfsdk:"archive_site"`
}

type systemBlockArchivalConfigurationBlockArchiveSite struct {
	URL      types.String `tfsdk:"url"`
	Password types.String `tfsdk:"password"`
}

type systemBlockInet6BackupRouter struct {
	Address     types.String   `tfsdk:"address"`
	Destination []types.String `tfsdk:"destination"`
}

type systemBlockInet6BackupRouterConfig struct {
	Address     types.String `tfsdk:"address"`
	Destination types.Set    `tfsdk:"destination"`
}

//nolint:lll
type systemBlockInternetOptions struct {
	GrePathMtuDiscovery                 types.Bool                                    `tfsdk:"gre_path_mtu_discovery"`
	NoGrePathMtuDiscovery               types.Bool                                    `tfsdk:"no_gre_path_mtu_discovery"`
	IpipPathMtuDiscovery                types.Bool                                    `tfsdk:"ipip_path_mtu_discovery"`
	NoIpipPathMtuDiscovery              types.Bool                                    `tfsdk:"no_ipip_path_mtu_discovery"`
	IPv6DuplicateAddrDetectionTransmits types.Int64                                   `tfsdk:"ipv6_duplicate_addr_detection_transmits"`
	IPv6PathMtuDiscovery                types.Bool                                    `tfsdk:"ipv6_path_mtu_discovery"`
	NoIPv6PathMtuDiscovery              types.Bool                                    `tfsdk:"no_ipv6_path_mtu_discovery"`
	IPv6PathMtuDiscoveryTimeout         types.Int64                                   `tfsdk:"ipv6_path_mtu_discovery_timeout"`
	IPv6RejectZeroHopLimit              types.Bool                                    `tfsdk:"ipv6_reject_zero_hop_limit"`
	NoIPv6RejectZeroHopLimit            types.Bool                                    `tfsdk:"no_ipv6_reject_zero_hop_limit"`
	NoTCPReset                          types.String                                  `tfsdk:"no_tcp_reset"`
	NoTCPRFC1323                        types.Bool                                    `tfsdk:"no_tcp_rfc1323"`
	NoTCPRFC1323Paws                    types.Bool                                    `tfsdk:"no_tcp_rfc1323_paws"`
	PathMtuDiscovery                    types.Bool                                    `tfsdk:"path_mtu_discovery"`
	NoPathMtuDiscovery                  types.Bool                                    `tfsdk:"no_path_mtu_discovery"`
	SourcePortUpperLimit                types.Int64                                   `tfsdk:"source_port_upper_limit"`
	SourceQuench                        types.Bool                                    `tfsdk:"source_quench"`
	NoSourceQuench                      types.Bool                                    `tfsdk:"no_source_quench"`
	TCPDropSynfinSet                    types.Bool                                    `tfsdk:"tcp_drop_synfin_set"`
	TCPMss                              types.Int64                                   `tfsdk:"tcp_mss"`
	IcmpV4RateLimit                     *systemBlockInternetOptionsBlockIcmpRateLimit `tfsdk:"icmpv4_rate_limit"`
	IcmpV6RateLimit                     *systemBlockInternetOptionsBlockIcmpRateLimit `tfsdk:"icmpv6_rate_limit"`
}

func (block *systemBlockInternetOptions) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockInternetOptionsBlockIcmpRateLimit struct {
	BucketSize types.Int64 `tfsdk:"bucket_size"`
	PacketRate types.Int64 `tfsdk:"packet_rate"`
}

func (block *systemBlockInternetOptionsBlockIcmpRateLimit) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockLicense struct {
	Autoupdate            types.Bool   `tfsdk:"autoupdate"`
	AutoupdatePassword    types.String `tfsdk:"autoupdate_password"`
	AutoupdateURL         types.String `tfsdk:"autoupdate_url"`
	RenewBeforeExpiration types.Int64  `tfsdk:"renew_before_expiration"`
	RenewInterval         types.Int64  `tfsdk:"renew_interval"`
}

func (block *systemBlockLicense) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockLogin struct {
	Announcement       types.String                       `tfsdk:"announcement"`
	DenySourcesAddress []types.String                     `tfsdk:"deny_sources_address"`
	IdleTimeout        types.Int64                        `tfsdk:"idle_timeout"`
	Message            types.String                       `tfsdk:"message"`
	Password           *systemBlockLoginBlockPassword     `tfsdk:"password"`
	RetryOptions       *systemBlockLoginBlockRetryOptions `tfsdk:"retry_options"`
}

func (block *systemBlockLogin) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockLoginConfig struct {
	Announcement       types.String                       `tfsdk:"announcement"`
	DenySourcesAddress types.Set                          `tfsdk:"deny_sources_address"`
	IdleTimeout        types.Int64                        `tfsdk:"idle_timeout"`
	Message            types.String                       `tfsdk:"message"`
	Password           *systemBlockLoginBlockPassword     `tfsdk:"password"`
	RetryOptions       *systemBlockLoginBlockRetryOptions `tfsdk:"retry_options"`
}

func (block *systemBlockLoginConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockNameServerOpts struct {
	Address         types.String `tfsdk:"address"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
}

type systemBlockLoginBlockPassword struct {
	ChangeType              types.String `tfsdk:"change_type"`
	Format                  types.String `tfsdk:"format"`
	MaximumLength           types.Int64  `tfsdk:"maximum_length"`
	MinimumChanges          types.Int64  `tfsdk:"minimum_changes"`
	MinimumCharacterChanges types.Int64  `tfsdk:"minimum_character_changes"`
	MinimumLength           types.Int64  `tfsdk:"minimum_length"`
	MinimumLowerCases       types.Int64  `tfsdk:"minimum_lower_cases"`
	MinimumNumerics         types.Int64  `tfsdk:"minimum_numerics"`
	MinimumPunctuations     types.Int64  `tfsdk:"minimum_punctuations"`
	MinimumReuse            types.Int64  `tfsdk:"minimum_reuse"`
	MinimumUpperCases       types.Int64  `tfsdk:"minimum_upper_cases"`
}

func (block *systemBlockLoginBlockPassword) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockLoginBlockRetryOptions struct {
	BackoffFactor         types.Int64 `tfsdk:"backoff_factor"`
	BackoffThreshold      types.Int64 `tfsdk:"backoff_threshold"`
	LockoutPeriod         types.Int64 `tfsdk:"lockout_period"`
	MaximumTime           types.Int64 `tfsdk:"maximum_time"`
	MinimumTime           types.Int64 `tfsdk:"minimum_time"`
	TriesBeforeDisconnect types.Int64 `tfsdk:"tries_before_disconnect"`
}

func (block *systemBlockLoginBlockRetryOptions) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockNtp struct {
	BootServer             types.String `tfsdk:"boot_server"`
	BroadcastClient        types.Bool   `tfsdk:"broadcast_client"`
	IntervalRange          types.Int64  `tfsdk:"interval_range"`
	MulticastClient        types.Bool   `tfsdk:"multicast_client"`
	MulticastClientAddress types.String `tfsdk:"multicast_client_address"`
	ThresholdAction        types.String `tfsdk:"threshold_action"`
	ThresholdValue         types.Int64  `tfsdk:"threshold_value"`
}

func (block *systemBlockNtp) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockPorts struct {
	AuxiliaryAuthenticationOrder []types.String `tfsdk:"auxiliary_authentication_order"`
	AuxiliaryDisable             types.Bool     `tfsdk:"auxiliary_disable"`
	AuxiliaryInsecure            types.Bool     `tfsdk:"auxiliary_insecure"`
	AuxiliaryLogoutOnDisconnect  types.Bool     `tfsdk:"auxiliary_logout_on_disconnect"`
	AuxiliaryType                types.String   `tfsdk:"auxiliary_type"`
	ConsoleAuthenticationOrder   []types.String `tfsdk:"console_authentication_order"`
	ConsoleDisable               types.Bool     `tfsdk:"console_disable"`
	ConsoleInsecure              types.Bool     `tfsdk:"console_insecure"`
	ConsoleLogoutOnDisconnect    types.Bool     `tfsdk:"console_logout_on_disconnect"`
	ConsoleType                  types.String   `tfsdk:"console_type"`
}

func (block *systemBlockPorts) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockPortsConfig struct {
	AuxiliaryAuthenticationOrder types.List   `tfsdk:"auxiliary_authentication_order"`
	AuxiliaryDisable             types.Bool   `tfsdk:"auxiliary_disable"`
	AuxiliaryInsecure            types.Bool   `tfsdk:"auxiliary_insecure"`
	AuxiliaryLogoutOnDisconnect  types.Bool   `tfsdk:"auxiliary_logout_on_disconnect"`
	AuxiliaryType                types.String `tfsdk:"auxiliary_type"`
	ConsoleAuthenticationOrder   types.List   `tfsdk:"console_authentication_order"`
	ConsoleDisable               types.Bool   `tfsdk:"console_disable"`
	ConsoleInsecure              types.Bool   `tfsdk:"console_insecure"`
	ConsoleLogoutOnDisconnect    types.Bool   `tfsdk:"console_logout_on_disconnect"`
	ConsoleType                  types.String `tfsdk:"console_type"`
}

func (block *systemBlockPortsConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:lll
type systemBlockServices struct {
	WebManagementSessionIdleTimeout types.Int64                                  `tfsdk:"web_management_session_idle_timeout"`
	WebManagementSessionLimit       types.Int64                                  `tfsdk:"web_management_session_limit"`
	NetconfSSH                      *systemBlockServicesBlockNetconfSSH          `tfsdk:"netconf_ssh"`
	NetconfTraceoptions             *systemBlockServicesBlockNetconfTraceoptions `tfsdk:"netconf_traceoptions"`
	SSH                             *systemBlockServicesBlockSSH                 `tfsdk:"ssh"`
	WebManagementHTTP               *systemBlockServicesBlockWebManagementHTTP   `tfsdk:"web_management_http"`
	WebManagementHTTPS              *systemBlockServicesBlockWebManagementHTTPS  `tfsdk:"web_management_https"`
}

func (block *systemBlockServices) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:lll
type systemBlockServicesConfig struct {
	WebManagementSessionIdleTimeout types.Int64                                        `tfsdk:"web_management_session_idle_timeout"`
	WebManagementSessionLimit       types.Int64                                        `tfsdk:"web_management_session_limit"`
	NetconfSSH                      *systemBlockServicesBlockNetconfSSH                `tfsdk:"netconf_ssh"`
	NetconfTraceoptions             *systemBlockServicesBlockNetconfTraceoptionsConfig `tfsdk:"netconf_traceoptions"`
	SSH                             *systemBlockServicesBlockSSHConfig                 `tfsdk:"ssh"`
	WebManagementHTTP               *systemBlockServicesBlockWebManagementHTTPConfig   `tfsdk:"web_management_http"`
	WebManagementHTTPS              *systemBlockServicesBlockWebManagementHTTPSConfig  `tfsdk:"web_management_https"`
}

func (block *systemBlockServicesConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockServicesBlockNetconfSSH struct {
	ClientAliveCountMax types.Int64 `tfsdk:"client_alive_count_max"`
	ClientAliveInterval types.Int64 `tfsdk:"client_alive_interval"`
	ConnectionLimit     types.Int64 `tfsdk:"connection_limit"`
	RateLimit           types.Int64 `tfsdk:"rate_limit"`
}

func (block *systemBlockServicesBlockNetconfSSH) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockServicesBlockNetconfTraceoptions struct {
	FileName            types.String   `tfsdk:"file_name"`
	FileFiles           types.Int64    `tfsdk:"file_files"`
	FileMatch           types.String   `tfsdk:"file_match"`
	FileSize            types.Int64    `tfsdk:"file_size"`
	FileWorldReadable   types.Bool     `tfsdk:"file_world_readable"`
	FileNoWorldReadable types.Bool     `tfsdk:"file_no_world_readable"`
	Flag                []types.String `tfsdk:"flag"`
	NoRemoteTrace       types.Bool     `tfsdk:"no_remote_trace"`
	OnDemand            types.Bool     `tfsdk:"on_demand"`
}

func (block *systemBlockServicesBlockNetconfTraceoptions) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockServicesBlockNetconfTraceoptionsConfig struct {
	FileName            types.String `tfsdk:"file_name"`
	FileFiles           types.Int64  `tfsdk:"file_files"`
	FileMatch           types.String `tfsdk:"file_match"`
	FileSize            types.Int64  `tfsdk:"file_size"`
	FileWorldReadable   types.Bool   `tfsdk:"file_world_readable"`
	FileNoWorldReadable types.Bool   `tfsdk:"file_no_world_readable"`
	Flag                types.Set    `tfsdk:"flag"`
	NoRemoteTrace       types.Bool   `tfsdk:"no_remote_trace"`
	OnDemand            types.Bool   `tfsdk:"on_demand"`
}

func (block *systemBlockServicesBlockNetconfTraceoptionsConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockServicesBlockSSH struct {
	AuthenticationOrder         []types.String `tfsdk:"authentication_order"`
	Ciphers                     []types.String `tfsdk:"ciphers"`
	ClientAliveCountMax         types.Int64    `tfsdk:"client_alive_count_max"`
	ClientAliveInterval         types.Int64    `tfsdk:"client_alive_interval"`
	ConnectionLimit             types.Int64    `tfsdk:"connection_limit"`
	FingerprintHash             types.String   `tfsdk:"fingerprint_hash"`
	HostkeyAlgorithm            []types.String `tfsdk:"hostkey_algorithm"`
	KeyExchange                 []types.String `tfsdk:"key_exchange"`
	LogKeyChanges               types.Bool     `tfsdk:"log_key_changes"`
	Macs                        []types.String `tfsdk:"macs"`
	MaxPreAuthenticationPackets types.Int64    `tfsdk:"max_pre_authentication_packets"`
	MaxSessionsPerConnection    types.Int64    `tfsdk:"max_sessions_per_connection"`
	NoPasswords                 types.Bool     `tfsdk:"no_passwords"`
	NoPublicKeys                types.Bool     `tfsdk:"no_public_keys"`
	Port                        types.Int64    `tfsdk:"port"`
	ProtocolVersion             []types.String `tfsdk:"protocol_version"`
	RateLimit                   types.Int64    `tfsdk:"rate_limit"`
	RootLogin                   types.String   `tfsdk:"root_login"`
	TCPForwarding               types.Bool     `tfsdk:"tcp_forwarding"`
	NoTCPForwarding             types.Bool     `tfsdk:"no_tcp_forwarding"`
}

func (block *systemBlockServicesBlockSSH) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockServicesBlockSSHConfig struct {
	AuthenticationOrder         types.List   `tfsdk:"authentication_order"`
	Ciphers                     types.Set    `tfsdk:"ciphers"`
	ClientAliveCountMax         types.Int64  `tfsdk:"client_alive_count_max"`
	ClientAliveInterval         types.Int64  `tfsdk:"client_alive_interval"`
	ConnectionLimit             types.Int64  `tfsdk:"connection_limit"`
	FingerprintHash             types.String `tfsdk:"fingerprint_hash"`
	HostkeyAlgorithm            types.Set    `tfsdk:"hostkey_algorithm"`
	KeyExchange                 types.Set    `tfsdk:"key_exchange"`
	LogKeyChanges               types.Bool   `tfsdk:"log_key_changes"`
	Macs                        types.Set    `tfsdk:"macs"`
	MaxPreAuthenticationPackets types.Int64  `tfsdk:"max_pre_authentication_packets"`
	MaxSessionsPerConnection    types.Int64  `tfsdk:"max_sessions_per_connection"`
	NoPasswords                 types.Bool   `tfsdk:"no_passwords"`
	NoPublicKeys                types.Bool   `tfsdk:"no_public_keys"`
	Port                        types.Int64  `tfsdk:"port"`
	ProtocolVersion             types.Set    `tfsdk:"protocol_version"`
	RateLimit                   types.Int64  `tfsdk:"rate_limit"`
	RootLogin                   types.String `tfsdk:"root_login"`
	TCPForwarding               types.Bool   `tfsdk:"tcp_forwarding"`
	NoTCPForwarding             types.Bool   `tfsdk:"no_tcp_forwarding"`
}

func (block *systemBlockServicesBlockSSHConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockServicesBlockWebManagementHTTP struct {
	Interface []types.String `tfsdk:"interface"`
	Port      types.Int64    `tfsdk:"port"`
}

type systemBlockServicesBlockWebManagementHTTPConfig struct {
	Interface types.Set   `tfsdk:"interface"`
	Port      types.Int64 `tfsdk:"port"`
}

type systemBlockServicesBlockWebManagementHTTPS struct {
	Interface                  []types.String `tfsdk:"interface"`
	LocalCertificate           types.String   `tfsdk:"local_certificate"`
	PkiLocalCertificate        types.String   `tfsdk:"pki_local_certificate"`
	Port                       types.Int64    `tfsdk:"port"`
	SystemGeneratedCertificate types.Bool     `tfsdk:"system_generated_certificate"`
}

type systemBlockServicesBlockWebManagementHTTPSConfig struct {
	Interface                  types.Set    `tfsdk:"interface"`
	LocalCertificate           types.String `tfsdk:"local_certificate"`
	PkiLocalCertificate        types.String `tfsdk:"pki_local_certificate"`
	Port                       types.Int64  `tfsdk:"port"`
	SystemGeneratedCertificate types.Bool   `tfsdk:"system_generated_certificate"`
}

type systemBlockSyslog struct {
	LogRotateFrequency    types.Int64                    `tfsdk:"log_rotate_frequency"`
	SourceAddress         types.String                   `tfsdk:"source_address"`
	TimeFormatMillisecond types.Bool                     `tfsdk:"time_format_millisecond"`
	TimeFormatYear        types.Bool                     `tfsdk:"time_format_year"`
	Archive               *systemBlockSyslogBlockArchive `tfsdk:"archive"`
	Console               *systemBlockSyslogBlockConsole `tfsdk:"console"`
}

func (block *systemBlockSyslog) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemBlockSyslogBlockArchive struct {
	BinaryData      types.Bool  `tfsdk:"binary_data"`
	NoBinaryData    types.Bool  `tfsdk:"no_binary_data"`
	Files           types.Int64 `tfsdk:"files"`
	Size            types.Int64 `tfsdk:"size"`
	WorldReadable   types.Bool  `tfsdk:"world_readable"`
	NoWorldReadable types.Bool  `tfsdk:"no_world_readable"`
}

type systemBlockSyslogBlockConsole struct {
	AnySeverity                 types.String `tfsdk:"any_severity"`
	AuthorizationSeverity       types.String `tfsdk:"authorization_severity"`
	ChangelogSeverity           types.String `tfsdk:"changelog_severity"`
	ConflictlogSeverity         types.String `tfsdk:"conflictlog_severity"`
	DaemonSeverity              types.String `tfsdk:"daemon_severity"`
	DfcSeverity                 types.String `tfsdk:"dfc_severity"`
	ExternalSeverity            types.String `tfsdk:"external_severity"`
	FirewallSeverity            types.String `tfsdk:"firewall_severity"`
	FtpSeverity                 types.String `tfsdk:"ftp_severity"`
	InteractivecommandsSeverity types.String `tfsdk:"interactivecommands_severity"`
	KernelSeverity              types.String `tfsdk:"kernel_severity"`
	NtpSeverity                 types.String `tfsdk:"ntp_severity"`
	PfeSeverity                 types.String `tfsdk:"pfe_severity"`
	SecuritySeverity            types.String `tfsdk:"security_severity"`
	UserSeverity                types.String `tfsdk:"user_severity"`
}

func (block *systemBlockSyslogBlockConsole) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:gocyclo,gocognit
func (rsc *system) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config systemConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.NameServer.IsNull() && !config.NameServer.IsUnknown() &&
		!config.NameServerOpts.IsNull() && !config.NameServerOpts.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name_server"),
			tfdiag.ConflictConfigErrSummary,
			"name_server and name_server_opts cannot be configured together",
		)
	}
	if !config.TacplusOptionsExcludeCmdAttribute.IsNull() && !config.TacplusOptionsExcludeCmdAttribute.IsUnknown() &&
		!config.TacplusOptionsNoCmdAttributeValue.IsNull() && !config.TacplusOptionsNoCmdAttributeValue.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tacplus_options_exclude_cmd_attribute"),
			tfdiag.ConflictConfigErrSummary,
			"tacplus_options_exclude_cmd_attribute and tacplus_options_no_cmd_attribute_value cannot be configured together",
		)
	}
	if !config.TacplusOptionsStrictAuthorization.IsNull() && !config.TacplusOptionsStrictAuthorization.IsUnknown() &&
		!config.TacplusOptionsNoStrictAuthorization.IsNull() && !config.TacplusOptionsNoStrictAuthorization.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tacplus_options_strict_authorization"),
			tfdiag.ConflictConfigErrSummary,
			"tacplus_options_strict_authorization and tacplus_options_no_strict_authorization cannot be configured together",
		)
	}
	if config.Accounting != nil {
		if config.Accounting.Events.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("accounting").AtName("events"),
				tfdiag.MissingConfigErrSummary,
				"events must be specified in accounting block",
			)
		}
		if config.Accounting.DestinationRadius.IsNull() &&
			config.Accounting.DestinationTacplus.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("accounting").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"one of destination_radius or destination_tacplus must be specified"+
					" in accounting block",
			)
		}
		if !config.Accounting.DestinationRadiusServer.IsNull() {
			if config.Accounting.DestinationRadius.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("accounting").AtName("destination_radius_server"),
					tfdiag.MissingConfigErrSummary,
					"destination_radius must be specified with destination_radius_server"+
						" in accounting block",
				)
			}
			if !config.Accounting.DestinationRadiusServer.IsUnknown() {
				var configDestinationRadiusServer []systemBlockAccountingBlockDestinationRadiusServer
				asDiags := config.Accounting.DestinationRadiusServer.ElementsAs(ctx, &configDestinationRadiusServer, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}
				destinationRadiusServerAddress := make(map[string]struct{})
				for i, block := range configDestinationRadiusServer {
					if block.Address.IsUnknown() {
						continue
					}
					address := block.Address.ValueString()
					if _, ok := destinationRadiusServerAddress[address]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("accounting").AtName("destination_radius_server").AtListIndex(i).AtName("address"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple destination_radius_server blocks with the same address %q"+
								" in accounting block", address),
						)
					}
					destinationRadiusServerAddress[address] = struct{}{}
				}
			}
		}
		if !config.Accounting.DestinationTacplusServer.IsNull() {
			if config.Accounting.DestinationTacplus.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("accounting").AtName("destination_tacplus_server"),
					tfdiag.MissingConfigErrSummary,
					"destination_tacplus must be specified with destination_tacplus_server"+
						" in accounting block",
				)
			}
			if !config.Accounting.DestinationTacplusServer.IsUnknown() {
				var configDestinationTacplusServer []systemBlockAccountingBlockDestinationTacplusServer
				asDiags := config.Accounting.DestinationTacplusServer.ElementsAs(ctx, &configDestinationTacplusServer, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}
				destinationTacplusServerAddress := make(map[string]struct{})
				for i, block := range configDestinationTacplusServer {
					if block.Address.IsUnknown() {
						continue
					}
					address := block.Address.ValueString()
					if _, ok := destinationTacplusServerAddress[address]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("accounting").AtName("destination_tacplus_server").AtListIndex(i).AtName("address"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple destination_tacplus_server blocks with the same address %q"+
								" in accounting block", address),
						)
					}
					destinationTacplusServerAddress[address] = struct{}{}
				}
			}
		}
	}
	if config.ArchivalConfiguration != nil {
		if config.ArchivalConfiguration.ArchiveSite.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("archival_configuration").AtName("archive_site"),
				tfdiag.MissingConfigErrSummary,
				"archive_site must be specified in archival_configuration block",
			)
		} else if !config.ArchivalConfiguration.ArchiveSite.IsUnknown() {
			var configArchiveSite []systemBlockArchivalConfigurationBlockArchiveSite
			asDiags := config.ArchivalConfiguration.ArchiveSite.ElementsAs(ctx, &configArchiveSite, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}
			archiveSiteURL := make(map[string]struct{})
			for i, block := range configArchiveSite {
				if !block.URL.IsUnknown() {
					url := block.URL.ValueString()
					if _, ok := archiveSiteURL[url]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("archival_configuration").AtName("archive_site").AtListIndex(i).AtName("url"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple archive_site blocks with the same url %q"+
								" in archival_configuration block", url),
						)
					}
					archiveSiteURL[url] = struct{}{}
				}
			}
		}
		if config.ArchivalConfiguration.TransferInterval.IsNull() &&
			config.ArchivalConfiguration.TransferOnCommit.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("archival_configuration").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"one of transfer_interval or transfer_on_commit must be specified"+
					" in archival_configuration block",
			)
		}
		if !config.ArchivalConfiguration.TransferInterval.IsNull() &&
			!config.ArchivalConfiguration.TransferInterval.IsUnknown() &&
			!config.ArchivalConfiguration.TransferOnCommit.IsNull() &&
			!config.ArchivalConfiguration.TransferOnCommit.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("archival_configuration").AtName("transfer_on_commit"),
				tfdiag.ConflictConfigErrSummary,
				"only one of transfer_interval or transfer_on_commit must be specified"+
					" in archival_configuration block",
			)
		}
	}
	if config.Inet6BackupRouter != nil {
		if config.Inet6BackupRouter.Address.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("inet6_backup_router").AtName("address"),
				tfdiag.MissingConfigErrSummary,
				"address must be specified in inet6_backup_router block",
			)
		}
		if config.Inet6BackupRouter.Destination.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("inet6_backup_router").AtName("destination"),
				tfdiag.MissingConfigErrSummary,
				"destination must be specified in inet6_backup_router block",
			)
		}
	}
	if config.InternetOptions != nil {
		if config.InternetOptions.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("internet_options").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"internet_options block is empty",
			)
		}
		if !config.InternetOptions.GrePathMtuDiscovery.IsNull() &&
			!config.InternetOptions.GrePathMtuDiscovery.IsUnknown() &&
			!config.InternetOptions.NoGrePathMtuDiscovery.IsNull() &&
			!config.InternetOptions.NoGrePathMtuDiscovery.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("internet_options").AtName("gre_path_mtu_discovery"),
				tfdiag.ConflictConfigErrSummary,
				"gre_path_mtu_discovery and no_gre_path_mtu_discovery cannot be configured together"+
					" in internet_options block",
			)
		}
		if !config.InternetOptions.IpipPathMtuDiscovery.IsNull() &&
			!config.InternetOptions.IpipPathMtuDiscovery.IsUnknown() &&
			!config.InternetOptions.NoIpipPathMtuDiscovery.IsNull() &&
			!config.InternetOptions.NoIpipPathMtuDiscovery.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("internet_options").AtName("ipip_path_mtu_discovery"),
				tfdiag.ConflictConfigErrSummary,
				"ipip_path_mtu_discovery and no_ipip_path_mtu_discovery cannot be configured together"+
					" in internet_options block",
			)
		}
		if !config.InternetOptions.IPv6PathMtuDiscovery.IsNull() &&
			!config.InternetOptions.IPv6PathMtuDiscovery.IsUnknown() &&
			!config.InternetOptions.NoIPv6PathMtuDiscovery.IsNull() &&
			!config.InternetOptions.NoIPv6PathMtuDiscovery.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("internet_options").AtName("ipv6_path_mtu_discovery"),
				tfdiag.ConflictConfigErrSummary,
				"ipv6_path_mtu_discovery and no_ipv6_path_mtu_discovery cannot be configured together"+
					" in internet_options block",
			)
		}
		if !config.InternetOptions.IPv6RejectZeroHopLimit.IsNull() &&
			!config.InternetOptions.IPv6RejectZeroHopLimit.IsUnknown() &&
			!config.InternetOptions.NoIPv6RejectZeroHopLimit.IsNull() &&
			!config.InternetOptions.NoIPv6RejectZeroHopLimit.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("internet_options").AtName("ipv6_reject_zero_hop_limit"),
				tfdiag.ConflictConfigErrSummary,
				"ipv6_reject_zero_hop_limit and no_ipv6_reject_zero_hop_limit cannot be configured together"+
					" in internet_options block",
			)
		}
		if !config.InternetOptions.PathMtuDiscovery.IsNull() &&
			!config.InternetOptions.PathMtuDiscovery.IsUnknown() &&
			!config.InternetOptions.NoPathMtuDiscovery.IsNull() &&
			!config.InternetOptions.NoPathMtuDiscovery.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("internet_options").AtName("path_mtu_discovery"),
				tfdiag.ConflictConfigErrSummary,
				"path_mtu_discovery and no_path_mtu_discovery cannot be configured together"+
					" in internet_options block",
			)
		}
		if !config.InternetOptions.SourceQuench.IsNull() &&
			!config.InternetOptions.SourceQuench.IsUnknown() &&
			!config.InternetOptions.NoSourceQuench.IsNull() &&
			!config.InternetOptions.NoSourceQuench.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("internet_options").AtName("source_quench"),
				tfdiag.ConflictConfigErrSummary,
				"source_quench and no_source_quench cannot be configured together"+
					" in internet_options block",
			)
		}
		if config.InternetOptions.IcmpV4RateLimit != nil {
			if config.InternetOptions.IcmpV4RateLimit.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("internet_options").AtName("icmpv4_rate_limit").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"icmpv4_rate_limit block in internet_options block is empty",
				)
			}
		}
		if config.InternetOptions.IcmpV6RateLimit != nil {
			if config.InternetOptions.IcmpV6RateLimit.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("internet_options").AtName("icmpv6_rate_limit").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"icmpv6_rate_limit block in internet_options block is empty",
				)
			}
		}
	}
	if config.License != nil {
		if config.License.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("license").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"license block is empty",
			)
		}
		if !config.License.AutoupdateURL.IsNull() &&
			config.License.Autoupdate.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("license").AtName("autoupdate_url"),
				tfdiag.MissingConfigErrSummary,
				"autoupdate must be specified with autoupdate_url in license block",
			)
		}
		if !config.License.AutoupdatePassword.IsNull() &&
			config.License.AutoupdateURL.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("license").AtName("autoupdate_password"),
				tfdiag.MissingConfigErrSummary,
				"autoupdate_url must be specified with autoupdate_password in license block",
			)
		}
		if !config.License.RenewInterval.IsNull() &&
			config.License.RenewBeforeExpiration.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("license").AtName("renew_interval"),
				tfdiag.MissingConfigErrSummary,
				"renew_before_expiration and renew_interval must be configured together"+
					" in license block",
			)
		}
		if !config.License.RenewBeforeExpiration.IsNull() &&
			config.License.RenewInterval.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("license").AtName("renew_before_expiration"),
				tfdiag.MissingConfigErrSummary,
				"renew_before_expiration and renew_interval must be configured together"+
					" in license block",
			)
		}
	}
	if config.Login != nil {
		if config.Login.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("login").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"login block is empty",
			)
		}
		if config.Login.Password != nil {
			if config.Login.Password.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("login").AtName("password").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"password block in login block is empty",
				)
			}
		}
		if config.Login.RetryOptions != nil {
			if config.Login.RetryOptions.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("login").AtName("retry_options").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"retry_options block in login block is empty",
				)
			}
		}
	}
	if config.Ntp != nil {
		if config.Ntp.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ntp").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"ntp block is empty",
			)
		}
		if !config.Ntp.MulticastClientAddress.IsNull() &&
			config.Ntp.MulticastClient.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ntp").AtName("multicast_client_address"),
				tfdiag.MissingConfigErrSummary,
				"multicast_client must be specified with multicast_client_address"+
					" in ntp block",
			)
		}
		if !config.Ntp.ThresholdAction.IsNull() &&
			config.Ntp.ThresholdValue.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ntp").AtName("threshold_action"),
				tfdiag.MissingConfigErrSummary,
				"threshold_action and threshold_value must be configured together"+
					" in ntp block",
			)
		}
		if !config.Ntp.ThresholdValue.IsNull() &&
			config.Ntp.ThresholdAction.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ntp").AtName("threshold_value"),
				tfdiag.MissingConfigErrSummary,
				"threshold_action and threshold_value must be configured together"+
					" in ntp block",
			)
		}
	}
	if config.Ports != nil {
		if config.Ports.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ports").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"ports block is empty",
			)
		}
	}
	if config.Services != nil {
		if config.Services.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("services").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"services block is empty",
			)
		}
		if config.Services.NetconfSSH != nil {
			if config.Services.NetconfSSH.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("services").AtName("netconf_ssh").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"netconf_ssh block in services block is empty",
				)
			}
		}
		if config.Services.NetconfTraceoptions != nil {
			if config.Services.NetconfTraceoptions.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("services").AtName("netconf_traceoptions").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"netconf_traceoptions block in services block is empty",
				)
			}
			if !config.Services.NetconfTraceoptions.FileWorldReadable.IsNull() &&
				!config.Services.NetconfTraceoptions.FileWorldReadable.IsUnknown() &&
				!config.Services.NetconfTraceoptions.FileNoWorldReadable.IsNull() &&
				!config.Services.NetconfTraceoptions.FileNoWorldReadable.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("services").AtName("netconf_traceoptions").AtName("file_world_readable"),
					tfdiag.ConflictConfigErrSummary,
					"file_world_readable and file_no_world_readable cannot be configured together"+
						" in netconf_traceoptions block in services block",
				)
			}
		}
		if config.Services.SSH != nil {
			if config.Services.SSH.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("services").AtName("ssh").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"ssh block in services block is empty",
				)
			}
			if !config.Services.SSH.NoPasswords.IsNull() &&
				!config.Services.SSH.NoPasswords.IsUnknown() &&
				!config.Services.SSH.NoPublicKeys.IsNull() &&
				!config.Services.SSH.NoPublicKeys.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("services").AtName("ssh").AtName("no_passwords"),
					tfdiag.ConflictConfigErrSummary,
					"no_passwords and no_public_keys cannot be configured together"+
						" in ssh block in services block",
				)
			}
			if !config.Services.SSH.TCPForwarding.IsNull() &&
				!config.Services.SSH.TCPForwarding.IsUnknown() &&
				!config.Services.SSH.NoTCPForwarding.IsNull() &&
				!config.Services.SSH.NoTCPForwarding.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("services").AtName("ssh").AtName("tcp_forwarding"),
					tfdiag.ConflictConfigErrSummary,
					"tcp_forwarding and no_tcp_forwarding cannot be configured together"+
						" in ssh block in services block",
				)
			}
		}
		if !config.Services.WebManagementSessionIdleTimeout.IsNull() {
			if config.Services.WebManagementHTTP == nil && config.Services.WebManagementHTTPS == nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("services").AtName("web_management_session_idle_timeout"),
					tfdiag.MissingConfigErrSummary,
					"web_management_http or web_management_https block must be specified"+
						" with web_management_session_idle_timeout in services block",
				)
			}
		}
		if !config.Services.WebManagementSessionLimit.IsNull() {
			if config.Services.WebManagementHTTP == nil && config.Services.WebManagementHTTPS == nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("services").AtName("web_management_session_limit"),
					tfdiag.MissingConfigErrSummary,
					"web_management_http or web_management_https block must be specified"+
						" with web_management_session_limit in services block",
				)
			}
		}
		if config.Services.WebManagementHTTPS != nil {
			if config.Services.WebManagementHTTPS.LocalCertificate.IsNull() &&
				config.Services.WebManagementHTTPS.PkiLocalCertificate.IsNull() &&
				config.Services.WebManagementHTTPS.SystemGeneratedCertificate.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("services").AtName("web_management_https").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"one of local_certificate, pki_local_certificate or system_generated_certificate must be specified"+
						" in web_management_https block in services block",
				)
			}
			if !config.Services.WebManagementHTTPS.LocalCertificate.IsNull() &&
				!config.Services.WebManagementHTTPS.LocalCertificate.IsUnknown() {
				if (!config.Services.WebManagementHTTPS.PkiLocalCertificate.IsNull() &&
					!config.Services.WebManagementHTTPS.PkiLocalCertificate.IsUnknown()) ||
					(!config.Services.WebManagementHTTPS.SystemGeneratedCertificate.IsNull() &&
						!config.Services.WebManagementHTTPS.SystemGeneratedCertificate.IsUnknown()) {
					resp.Diagnostics.AddAttributeError(
						path.Root("services").AtName("web_management_https").AtName("local_certificate"),
						tfdiag.ConflictConfigErrSummary,
						"only one of local_certificate, pki_local_certificate or system_generated_certificate must be specified"+
							" in web_management_https block in services block",
					)
				}
			}
			if !config.Services.WebManagementHTTPS.PkiLocalCertificate.IsNull() &&
				!config.Services.WebManagementHTTPS.PkiLocalCertificate.IsUnknown() &&
				!config.Services.WebManagementHTTPS.SystemGeneratedCertificate.IsNull() &&
				!config.Services.WebManagementHTTPS.SystemGeneratedCertificate.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("services").AtName("web_management_https").AtName("pki_local_certificate"),
					tfdiag.ConflictConfigErrSummary,
					"only one of local_certificate, pki_local_certificate or system_generated_certificate must be specified"+
						" in web_management_https block in services block",
				)
			}
		}
	}
	if config.Syslog != nil {
		if config.Syslog.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("syslog").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"syslog block is empty",
			)
		}
		if config.Syslog.Archive != nil {
			if !config.Syslog.Archive.BinaryData.IsNull() && !config.Syslog.Archive.BinaryData.IsUnknown() &&
				!config.Syslog.Archive.NoBinaryData.IsNull() && !config.Syslog.Archive.NoBinaryData.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("syslog").AtName("archive").AtName("binary_data"),
					tfdiag.ConflictConfigErrSummary,
					"binary_data and no_binary_data cannot be configured together"+
						" in archive block in syslog block",
				)
			}
			if !config.Syslog.Archive.WorldReadable.IsNull() && !config.Syslog.Archive.WorldReadable.IsUnknown() &&
				!config.Syslog.Archive.NoWorldReadable.IsNull() && !config.Syslog.Archive.NoWorldReadable.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("syslog").AtName("archive").AtName("world_readable"),
					tfdiag.ConflictConfigErrSummary,
					"world_readable and no_world_readable cannot be configured together"+
						" in archive block in syslog block",
				)
			}
		}
		if config.Syslog.Console != nil {
			if config.Syslog.Console.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("syslog").AtName("console").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"console block in syslog block is empty",
				)
			}
		}
	}
}

func (rsc *system) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemData
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

func (rsc *system) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemData
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
			nameServerWithOpts := false
			for _, block := range data.NameServerOpts {
				if !block.RoutingInstance.IsNull() {
					nameServerWithOpts = true

					break
				}
			}
			if nameServerWithOpts || len(state.NameServerOpts) != 0 {
				data.NameServer = nil
			} else {
				data.NameServerOpts = nil
			}
		},
		resp,
	)
}

func (rsc *system) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemData
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

func (rsc *system) Delete(
	_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse,
) {
	// no-op
}

func (rsc *system) ImportState(
	ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	var data systemData
	if err := data.read(ctx, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}

	nameServerWithOpts := false
	for _, block := range data.NameServerOpts {
		if !block.RoutingInstance.IsNull() {
			nameServerWithOpts = true

			break
		}
	}
	if nameServerWithOpts {
		data.NameServer = nil
	} else {
		data.NameServerOpts = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rscData *systemData) fillID() {
	rscData.ID = types.StringValue("system")
}

func (rscData *systemData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set system "

	for _, v := range rscData.AuthenticationOrder {
		configSet = append(configSet, setPrefix+"authentication-order "+v.ValueString())
	}
	if rscData.AutoSnapshot.ValueBool() {
		configSet = append(configSet, setPrefix+"auto-snapshot")
	}
	if rscData.DefaultAddressSelection.ValueBool() {
		configSet = append(configSet, setPrefix+"default-address-selection")
	}
	if v := rscData.DomainName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"domain-name "+v)
	}
	if v := rscData.HostName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"host-name "+v)
	}
	if !rscData.MaxConfigurationRollbacks.IsNull() {
		configSet = append(configSet, setPrefix+"max-configuration-rollbacks "+
			utils.ConvI64toa(rscData.MaxConfigurationRollbacks.ValueInt64()))
	}
	if !rscData.MaxConfigurationsOnFlash.IsNull() {
		configSet = append(configSet, setPrefix+"max-configurations-on-flash "+
			utils.ConvI64toa(rscData.MaxConfigurationsOnFlash.ValueInt64()))
	}
	for _, v := range rscData.NameServer {
		configSet = append(configSet, setPrefix+"name-server "+v.ValueString())
	}
	if rscData.NoMulticastEcho.ValueBool() {
		configSet = append(configSet, setPrefix+"no-multicast-echo")
	}
	if rscData.NoPingRecordRoute.ValueBool() {
		configSet = append(configSet, setPrefix+"no-ping-record-route")
	}
	if rscData.NoPingTimestamp.ValueBool() {
		configSet = append(configSet, setPrefix+"no-ping-time-stamp")
	}
	if rscData.NoRedirects.ValueBool() {
		configSet = append(configSet, setPrefix+"no-redirects")
	}
	if rscData.NoRedirectsIPv6.ValueBool() {
		configSet = append(configSet, setPrefix+"no-redirects-ipv6")
	}
	if v := rscData.RadiusOptionsAttributesNasID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"radius-options attributes nas-id \""+v+"\"")
	}
	if v := rscData.RadiusOptionsAttributesNasIpaddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"radius-options attributes nas-ip-address "+v)
	}
	if rscData.RadiusOptionsEnhancedAccounting.ValueBool() {
		configSet = append(configSet, setPrefix+"radius-options enhanced-accounting")
	}
	if rscData.RadiusOptionsPasswordProtocolMschapv2.ValueBool() {
		configSet = append(configSet, setPrefix+"radius-options password-protocol mschap-v2")
	}
	if !rscData.TacplusOptionsAuthorizationTimeInterval.IsNull() {
		configSet = append(configSet, setPrefix+"tacplus-options authorization-time-interval "+
			utils.ConvI64toa(rscData.TacplusOptionsAuthorizationTimeInterval.ValueInt64()))
	}
	if rscData.TacplusOptionsEnhancedAccounting.ValueBool() {
		configSet = append(configSet, setPrefix+"tacplus-options enhanced-accounting")
	}
	if rscData.TacplusOptionsExcludeCmdAttribute.ValueBool() {
		configSet = append(configSet, setPrefix+"tacplus-options exclude-cmd-attribute")
	}
	if rscData.TacplusOptionsNoCmdAttributeValue.ValueBool() {
		configSet = append(configSet, setPrefix+"tacplus-options no-cmd-attribute-value")
	}
	if v := rscData.TacplusOptionsServiceName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"tacplus-options service-name \""+v+"\"")
	}
	if rscData.TacplusOptionsStrictAuthorization.ValueBool() {
		configSet = append(configSet, setPrefix+"tacplus-options strict-authorization")
	}
	if rscData.TacplusOptionsNoStrictAuthorization.ValueBool() {
		configSet = append(configSet, setPrefix+"tacplus-options no-strict-authorization")
	}
	if rscData.TacplusOptionsTimestampAndTimezone.ValueBool() {
		configSet = append(configSet, setPrefix+"tacplus-options timestamp-and-timezone")
	}
	if v := rscData.TimeZone.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"time-zone "+v)
	}
	if v := rscData.TracingDestOverrideSyslogHost.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"tracing destination-override syslog host "+v)
	}

	if rscData.Accounting != nil {
		blockSet, pathErr, err := rscData.Accounting.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.ArchivalConfiguration != nil {
		blockSet, pathErr, err := rscData.ArchivalConfiguration.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.Inet6BackupRouter != nil {
		configSet = append(configSet, setPrefix+"inet6-backup-router "+rscData.Inet6BackupRouter.Address.ValueString())
		for _, v := range rscData.Inet6BackupRouter.Destination {
			configSet = append(configSet, setPrefix+"inet6-backup-router destination "+v.ValueString())
		}
	}
	if rscData.InternetOptions != nil {
		if rscData.InternetOptions.isEmpty() {
			return path.Root("internet_options").AtName("*"),
				errors.New("internet_options block is empty")
		}
		blockSet, pathErr, err := rscData.InternetOptions.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.License != nil {
		if rscData.License.isEmpty() {
			return path.Root("license").AtName("*"),
				errors.New("license block is empty")
		}
		blockSet, pathErr, err := rscData.License.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.Login != nil {
		if rscData.Login.isEmpty() {
			return path.Root("login").AtName("*"),
				errors.New("login block is empty")
		}
		blockSet, pathErr, err := rscData.Login.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	for _, block := range rscData.NameServerOpts {
		configSet = append(configSet, setPrefix+"name-server "+block.Address.ValueString())
		if v := block.RoutingInstance.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"name-server "+block.Address.ValueString()+
				" routing-instance "+v)
		}
	}
	if rscData.Ntp != nil {
		if rscData.Ntp.isEmpty() {
			return path.Root("ntp").AtName("*"),
				errors.New("ntp block is empty")
		}
		blockSet, pathErr, err := rscData.Ntp.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.Ports != nil {
		if rscData.Ports.isEmpty() {
			return path.Root("ports").AtName("*"),
				errors.New("ports block is empty")
		}
		configSet = append(configSet, rscData.Ports.configSet()...)
	}
	if rscData.Services != nil {
		if rscData.Services.isEmpty() {
			return path.Root("services").AtName("*"),
				errors.New("services block is empty")
		}
		blockSet, pathErr, err := rscData.Services.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.Syslog != nil {
		if rscData.Syslog.isEmpty() {
			return path.Root("syslog").AtName("*"),
				errors.New("syslog block is empty")
		}
		blockSet, pathErr, err := rscData.Syslog.configSet()
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *systemBlockAccounting) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, len(block.Events)+1)
	setPrefix := "set system accounting "

	for _, v := range block.Events {
		configSet = append(configSet, setPrefix+"events "+v.ValueString())
	}
	if block.DestinationRadius.ValueBool() {
		configSet = append(configSet, setPrefix+"destination radius")
	}
	if block.DestinationTacplus.ValueBool() {
		configSet = append(configSet, setPrefix+"destination tacplus")
	}
	if !block.EnhancedAvsMax.IsNull() {
		configSet = append(configSet, setPrefix+"enhanced-avs-max "+
			utils.ConvI64toa(block.EnhancedAvsMax.ValueInt64()))
	}
	destinationRadiusServerAddress := make(map[string]struct{})
	for i, v := range block.DestinationRadiusServer {
		address := v.Address.ValueString()
		if _, ok := destinationRadiusServerAddress[address]; ok {
			return configSet, path.Root("accounting").AtName("destination_radius_server").AtListIndex(i).AtName("address"),
				fmt.Errorf("multiple destination_radius_server blocks with the same address %q"+
					" in accounting block", address)
		}
		destinationRadiusServerAddress[address] = struct{}{}

		configSet = append(configSet, v.configSet()...)
	}
	destinationTacplusServerAddress := make(map[string]struct{})
	for i, v := range block.DestinationTacplusServer {
		address := v.Address.ValueString()
		if _, ok := destinationTacplusServerAddress[address]; ok {
			return configSet, path.Root("accounting").AtName("destination_tacplus_server").AtListIndex(i).AtName("address"),
				fmt.Errorf("multiple destination_tacplus_server blocks with the same address %q"+
					" in accounting block", address)
		}
		destinationTacplusServerAddress[address] = struct{}{}

		configSet = append(configSet, v.configSet()...)
	}

	return configSet, path.Empty(), nil
}

func (block *systemBlockAccountingBlockDestinationRadiusServer) configSet() []string {
	setPrefix := "set system accounting destination radius server " + block.Address.ValueString() + " "

	configSet := []string{
		setPrefix + "secret \"" + block.Secret.ValueString() + "\"",
	}

	if !block.AccountingPort.IsNull() {
		configSet = append(configSet, setPrefix+"accounting-port "+
			utils.ConvI64toa(block.AccountingPort.ValueInt64()))
	}
	if !block.AccountingRetry.IsNull() {
		configSet = append(configSet, setPrefix+"accounting-retry "+
			utils.ConvI64toa(block.AccountingRetry.ValueInt64()))
	}
	if !block.AccountingTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"accounting-timeout "+
			utils.ConvI64toa(block.AccountingTimeout.ValueInt64()))
	}
	if !block.DynamicRequestPort.IsNull() {
		configSet = append(configSet, setPrefix+"dynamic-request-port "+
			utils.ConvI64toa(block.DynamicRequestPort.ValueInt64()))
	}
	if !block.MaxOutstandingRequests.IsNull() {
		configSet = append(configSet, setPrefix+"max-outstanding-requests "+
			utils.ConvI64toa(block.MaxOutstandingRequests.ValueInt64()))
	}
	if !block.Port.IsNull() {
		configSet = append(configSet, setPrefix+"port "+
			utils.ConvI64toa(block.Port.ValueInt64()))
	}
	if !block.PreauthenticationPort.IsNull() {
		configSet = append(configSet, setPrefix+"preauthentication-port "+
			utils.ConvI64toa(block.PreauthenticationPort.ValueInt64()))
	}
	if v := block.PreauthenticationSecret.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"preauthentication-secret \""+v+"\"")
	}
	if !block.Retry.IsNull() {
		configSet = append(configSet, setPrefix+"retry "+
			utils.ConvI64toa(block.Retry.ValueInt64()))
	}
	if v := block.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-instance "+v)
	}
	if v := block.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}
	if !block.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(block.Timeout.ValueInt64()))
	}

	return configSet
}

func (block *systemBlockAccountingBlockDestinationTacplusServer) configSet() []string {
	setPrefix := "set system accounting destination tacplus server " + block.Address.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if !block.Port.IsNull() {
		configSet = append(configSet, setPrefix+"port "+
			utils.ConvI64toa(block.Port.ValueInt64()))
	}
	if v := block.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-instance "+v)
	}
	if v := block.Secret.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"secret \""+v+"\"")
	}
	if block.SingleConnection.ValueBool() {
		configSet = append(configSet, setPrefix+"single-connection")
	}
	if v := block.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}
	if !block.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(block.Timeout.ValueInt64()))
	}

	return configSet
}

func (block *systemBlockArchivalConfiguration) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix := "set system archival configuration "

	archiveSiteURL := make(map[string]struct{})
	for i, v := range block.ArchiveSite {
		url := v.URL.ValueString()
		if _, ok := archiveSiteURL[url]; ok {
			return configSet, path.Root("archival_configuration").AtName("archive_site").AtListIndex(i).AtName("url"),
				fmt.Errorf("multiple archive_site blocks with the same url %q"+
					" in archival_configuration block", url)
		}
		archiveSiteURL[url] = struct{}{}
		configSet = append(configSet, setPrefix+"archive-sites \""+url+"\"")
		if password := v.Password.ValueString(); password != "" {
			configSet = append(configSet,
				setPrefix+"archive-sites \""+url+"\" password \""+password+"\"")
		}
	}
	switch {
	case !block.TransferInterval.IsNull():
		configSet = append(configSet, setPrefix+"transfer-interval "+
			utils.ConvI64toa(block.TransferInterval.ValueInt64()))
	case block.TransferOnCommit.ValueBool():
		configSet = append(configSet, setPrefix+"transfer-on-commit")
	default:
		return configSet, path.Root("archival_configuration").AtName("*"),
			errors.New("one of transfer_interval or transfer_on_commit must be specified" +
				" in archival_configuration block")
	}

	return configSet, path.Empty(), nil
}

func (block *systemBlockInternetOptions) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix := "set system internet-options "

	if block.GrePathMtuDiscovery.ValueBool() {
		configSet = append(configSet, setPrefix+"gre-path-mtu-discovery")
	}
	if block.NoGrePathMtuDiscovery.ValueBool() {
		configSet = append(configSet, setPrefix+"no-gre-path-mtu-discovery")
	}
	if block.IpipPathMtuDiscovery.ValueBool() {
		configSet = append(configSet, setPrefix+"ipip-path-mtu-discovery")
	}
	if block.NoIpipPathMtuDiscovery.ValueBool() {
		configSet = append(configSet, setPrefix+"no-ipip-path-mtu-discovery")
	}
	if !block.IPv6DuplicateAddrDetectionTransmits.IsNull() {
		configSet = append(configSet, setPrefix+"ipv6-duplicate-addr-detection-transmits "+
			utils.ConvI64toa(block.IPv6DuplicateAddrDetectionTransmits.ValueInt64()))
	}
	if block.IPv6PathMtuDiscovery.ValueBool() {
		configSet = append(configSet, setPrefix+"ipv6-path-mtu-discovery")
	}
	if block.NoIPv6PathMtuDiscovery.ValueBool() {
		configSet = append(configSet, setPrefix+"no-ipv6-path-mtu-discovery")
	}
	if !block.IPv6PathMtuDiscoveryTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"ipv6-path-mtu-discovery-timeout "+
			utils.ConvI64toa(block.IPv6PathMtuDiscoveryTimeout.ValueInt64()))
	}
	if block.IPv6RejectZeroHopLimit.ValueBool() {
		configSet = append(configSet, setPrefix+"ipv6-reject-zero-hop-limit")
	}
	if block.NoIPv6RejectZeroHopLimit.ValueBool() {
		configSet = append(configSet, setPrefix+"no-ipv6-reject-zero-hop-limit")
	}
	if v := block.NoTCPReset.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"no-tcp-reset "+v)
	}
	if block.NoTCPRFC1323.ValueBool() {
		configSet = append(configSet, setPrefix+"no-tcp-rfc1323")
	}
	if block.NoTCPRFC1323Paws.ValueBool() {
		configSet = append(configSet, setPrefix+"no-tcp-rfc1323-paws")
	}
	if block.PathMtuDiscovery.ValueBool() {
		configSet = append(configSet, setPrefix+"path-mtu-discovery")
	}
	if block.NoPathMtuDiscovery.ValueBool() {
		configSet = append(configSet, setPrefix+"no-path-mtu-discovery")
	}
	if !block.SourcePortUpperLimit.IsNull() {
		configSet = append(configSet, setPrefix+"source-port upper-limit "+
			utils.ConvI64toa(block.SourcePortUpperLimit.ValueInt64()))
	}
	if block.SourceQuench.ValueBool() {
		configSet = append(configSet, setPrefix+"source-quench")
	}
	if block.NoSourceQuench.ValueBool() {
		configSet = append(configSet, setPrefix+"no-source-quench")
	}
	if block.TCPDropSynfinSet.ValueBool() {
		configSet = append(configSet, setPrefix+"tcp-drop-synfin-set")
	}
	if !block.TCPMss.IsNull() {
		configSet = append(configSet, setPrefix+"tcp-mss "+
			utils.ConvI64toa(block.TCPMss.ValueInt64()))
	}

	if block.IcmpV4RateLimit != nil {
		if block.IcmpV4RateLimit.isEmpty() {
			return configSet, path.Root("internet_options").AtName("icmpv4_rate_limit").AtName("*"),
				errors.New("icmpv4_rate_limit block in internet_options block is empty")
		}
		if !block.IcmpV4RateLimit.BucketSize.IsNull() {
			configSet = append(configSet, setPrefix+"icmpv4-rate-limit bucket-size "+
				utils.ConvI64toa(block.IcmpV4RateLimit.BucketSize.ValueInt64()))
		}
		if !block.IcmpV4RateLimit.PacketRate.IsNull() {
			configSet = append(configSet, setPrefix+"icmpv4-rate-limit packet-rate "+
				utils.ConvI64toa(block.IcmpV4RateLimit.PacketRate.ValueInt64()))
		}
	}
	if block.IcmpV6RateLimit != nil {
		if block.IcmpV6RateLimit.isEmpty() {
			return configSet, path.Root("internet_options").AtName("icmpv6_rate_limit").AtName("*"),
				errors.New("icmpv6_rate_limit block in internet_options block is empty")
		}
		if !block.IcmpV6RateLimit.BucketSize.IsNull() {
			configSet = append(configSet, setPrefix+"icmpv6-rate-limit bucket-size "+
				utils.ConvI64toa(block.IcmpV6RateLimit.BucketSize.ValueInt64()))
		}
		if !block.IcmpV6RateLimit.PacketRate.IsNull() {
			configSet = append(configSet, setPrefix+"icmpv6-rate-limit packet-rate "+
				utils.ConvI64toa(block.IcmpV6RateLimit.PacketRate.ValueInt64()))
		}
	}

	return configSet, path.Empty(), nil
}

func (block *systemBlockLicense) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix := "set system license "

	if block.Autoupdate.ValueBool() {
		configSet = append(configSet, setPrefix+"autoupdate")
		if vURL := block.AutoupdateURL.ValueString(); vURL != "" {
			configSet = append(configSet, setPrefix+"autoupdate url \""+vURL+"\"")
			if v := block.AutoupdatePassword.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"autoupdate url \""+vURL+"\" password \""+v+"\"")
			}
		} else if block.AutoupdatePassword.ValueString() != "" {
			return configSet, path.Root("license").AtName("autoupdate_password"),
				errors.New("autoupdate_url must be specified with autoupdate_password in license block")
		}
	} else {
		if block.AutoupdateURL.ValueString() != "" {
			return configSet, path.Root("license").AtName("autoupdate_url"),
				errors.New("autoupdate must be specified with autoupdate_url in license block")
		} else if block.AutoupdatePassword.ValueString() != "" {
			return configSet, path.Root("license").AtName("autoupdate_password"),
				errors.New("autoupdate and autoupdate_url must be specified with autoupdate_password in license block")
		}
	}
	if !block.RenewBeforeExpiration.IsNull() {
		configSet = append(configSet, setPrefix+"renew before-expiration "+
			utils.ConvI64toa(block.RenewBeforeExpiration.ValueInt64()))
	}
	if !block.RenewInterval.IsNull() {
		configSet = append(configSet, setPrefix+"renew interval "+
			utils.ConvI64toa(block.RenewInterval.ValueInt64()))
	}

	return configSet, path.Empty(), nil
}

func (block *systemBlockLogin) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix := "set system login "

	if v := block.Announcement.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"announcement \""+v+"\"")
	}
	for _, v := range block.DenySourcesAddress {
		configSet = append(configSet, setPrefix+"deny-sources address "+v.ValueString())
	}
	if !block.IdleTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"idle-timeout "+
			utils.ConvI64toa(block.IdleTimeout.ValueInt64()))
	}
	if v := block.Message.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"message \""+v+"\"")
	}

	if block.Password != nil {
		if block.Password.isEmpty() {
			return configSet, path.Root("login").AtName("password").AtName("*"),
				errors.New("password block in login block is empty")
		}
		if v := block.Password.ChangeType.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"password change-type "+v)
		}
		if v := block.Password.Format.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"password format "+v)
		}
		if !block.Password.MaximumLength.IsNull() {
			configSet = append(configSet, setPrefix+"password maximum-length "+
				utils.ConvI64toa(block.Password.MaximumLength.ValueInt64()))
		}
		if !block.Password.MinimumChanges.IsNull() {
			configSet = append(configSet, setPrefix+"password minimum-changes "+
				utils.ConvI64toa(block.Password.MinimumChanges.ValueInt64()))
		}
		if !block.Password.MinimumCharacterChanges.IsNull() {
			configSet = append(configSet, setPrefix+"password minimum-character-changes "+
				utils.ConvI64toa(block.Password.MinimumCharacterChanges.ValueInt64()))
		}
		if !block.Password.MinimumLength.IsNull() {
			configSet = append(configSet, setPrefix+"password minimum-length "+
				utils.ConvI64toa(block.Password.MinimumLength.ValueInt64()))
		}
		if !block.Password.MinimumLowerCases.IsNull() {
			configSet = append(configSet, setPrefix+"password minimum-lower-cases "+
				utils.ConvI64toa(block.Password.MinimumLowerCases.ValueInt64()))
		}
		if !block.Password.MinimumNumerics.IsNull() {
			configSet = append(configSet, setPrefix+"password minimum-numerics "+
				utils.ConvI64toa(block.Password.MinimumNumerics.ValueInt64()))
		}
		if !block.Password.MinimumPunctuations.IsNull() {
			configSet = append(configSet, setPrefix+"password minimum-punctuations "+
				utils.ConvI64toa(block.Password.MinimumPunctuations.ValueInt64()))
		}
		if !block.Password.MinimumReuse.IsNull() {
			configSet = append(configSet, setPrefix+"password minimum-reuse "+
				utils.ConvI64toa(block.Password.MinimumReuse.ValueInt64()))
		}
		if !block.Password.MinimumUpperCases.IsNull() {
			configSet = append(configSet, setPrefix+"password minimum-upper-cases "+
				utils.ConvI64toa(block.Password.MinimumUpperCases.ValueInt64()))
		}
	}
	if block.RetryOptions != nil {
		if block.RetryOptions.isEmpty() {
			return configSet, path.Root("login").AtName("retry_options").AtName("*"),
				errors.New("retry_options block in login block is empty")
		}
		if !block.RetryOptions.BackoffFactor.IsNull() {
			configSet = append(configSet, setPrefix+"retry-options backoff-factor "+
				utils.ConvI64toa(block.RetryOptions.BackoffFactor.ValueInt64()))
		}
		if !block.RetryOptions.BackoffThreshold.IsNull() {
			configSet = append(configSet, setPrefix+"retry-options backoff-threshold "+
				utils.ConvI64toa(block.RetryOptions.BackoffThreshold.ValueInt64()))
		}
		if !block.RetryOptions.LockoutPeriod.IsNull() {
			configSet = append(configSet, setPrefix+"retry-options lockout-period "+
				utils.ConvI64toa(block.RetryOptions.LockoutPeriod.ValueInt64()))
		}
		if !block.RetryOptions.MaximumTime.IsNull() {
			configSet = append(configSet, setPrefix+"retry-options maximum-time "+
				utils.ConvI64toa(block.RetryOptions.MaximumTime.ValueInt64()))
		}
		if !block.RetryOptions.MinimumTime.IsNull() {
			configSet = append(configSet, setPrefix+"retry-options minimum-time "+
				utils.ConvI64toa(block.RetryOptions.MinimumTime.ValueInt64()))
		}
		if !block.RetryOptions.TriesBeforeDisconnect.IsNull() {
			configSet = append(configSet, setPrefix+"retry-options tries-before-disconnect "+
				utils.ConvI64toa(block.RetryOptions.TriesBeforeDisconnect.ValueInt64()))
		}
	}

	return configSet, path.Empty(), nil
}

func (block *systemBlockNtp) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix := "set system ntp "

	if v := block.BootServer.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"boot-server "+v)
	}
	if block.BroadcastClient.ValueBool() {
		configSet = append(configSet, setPrefix+"broadcast-client")
	}
	if !block.IntervalRange.IsNull() {
		configSet = append(configSet, setPrefix+"interval-range "+
			utils.ConvI64toa(block.IntervalRange.ValueInt64()))
	}
	if block.MulticastClient.ValueBool() {
		configSet = append(configSet, setPrefix+"multicast-client")
		if v := block.MulticastClientAddress.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"multicast-client "+v)
		}
	} else if block.MulticastClientAddress.ValueString() != "" {
		return configSet, path.Root("ntp").AtName("multicast_client_address"),
			errors.New("multicast_client must be specified with multicast_client_address in ntp block")
	}
	if !block.ThresholdValue.IsNull() {
		if v := block.ThresholdAction.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"threshold "+
				utils.ConvI64toa(block.ThresholdValue.ValueInt64())+
				" action "+v)
		} else {
			return configSet, path.Root("ntp").AtName("threshold_value"),
				errors.New("threshold_action and threshold_value must be configured together in ntp block")
		}
	} else if block.ThresholdAction.ValueString() != "" {
		return configSet, path.Root("ntp").AtName("threshold_action"),
			errors.New("threshold_action and threshold_value must be configured together in ntp block")
	}

	return configSet, path.Empty(), nil
}

func (block *systemBlockPorts) configSet() []string { // configSet
	configSet := make([]string, 0)
	setPrefix := "set system ports "

	for _, v := range block.AuxiliaryAuthenticationOrder {
		configSet = append(configSet, setPrefix+"auxiliary authentication-order "+v.ValueString())
	}
	if block.AuxiliaryDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"auxiliary disable")
	}
	if block.AuxiliaryInsecure.ValueBool() {
		configSet = append(configSet, setPrefix+"auxiliary insecure")
	}
	if block.AuxiliaryLogoutOnDisconnect.ValueBool() {
		configSet = append(configSet, setPrefix+"auxiliary log-out-on-disconnect")
	}
	if v := block.AuxiliaryType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"auxiliary type "+v)
	}
	for _, v := range block.ConsoleAuthenticationOrder {
		configSet = append(configSet, setPrefix+"console authentication-order "+v.ValueString())
	}
	if block.ConsoleDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"console disable")
	}
	if block.ConsoleInsecure.ValueBool() {
		configSet = append(configSet, setPrefix+"console insecure")
	}
	if block.ConsoleLogoutOnDisconnect.ValueBool() {
		configSet = append(configSet, setPrefix+"console log-out-on-disconnect")
	}
	if v := block.ConsoleType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"console type "+v)
	}

	return configSet
}

func (block *systemBlockServices) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix := "set system services "

	if !block.WebManagementSessionIdleTimeout.IsNull() {
		if block.WebManagementHTTP == nil && block.WebManagementHTTPS == nil {
			return configSet, path.Root("services").AtName("web_management_session_idle_timeout"),
				errors.New("web_management_http or web_management_https block must be specified" +
					" with web_management_session_idle_timeout in services block")
		}
		configSet = append(configSet, setPrefix+"web-management session idle-timeout "+
			utils.ConvI64toa(block.WebManagementSessionIdleTimeout.ValueInt64()))
	}
	if !block.WebManagementSessionLimit.IsNull() {
		if block.WebManagementHTTP == nil && block.WebManagementHTTPS == nil {
			return configSet, path.Root("services").AtName("web_management_session_limit"),
				errors.New("web_management_http or web_management_https block must be specified" +
					" with web_management_session_limit in services block")
		}
		configSet = append(configSet, setPrefix+"web-management session session-limit "+
			utils.ConvI64toa(block.WebManagementSessionLimit.ValueInt64()))
	}
	if block.NetconfSSH != nil {
		if block.NetconfSSH.isEmpty() {
			return configSet, path.Root("services").AtName("netconf_ssh").AtName("*"),
				errors.New("netconf_ssh block in services block is empty")
		}
		if !block.NetconfSSH.ClientAliveCountMax.IsNull() {
			configSet = append(configSet, setPrefix+"netconf ssh client-alive-count-max "+
				utils.ConvI64toa(block.NetconfSSH.ClientAliveCountMax.ValueInt64()))
		}
		if !block.NetconfSSH.ClientAliveInterval.IsNull() {
			configSet = append(configSet, setPrefix+"netconf ssh client-alive-interval "+
				utils.ConvI64toa(block.NetconfSSH.ClientAliveInterval.ValueInt64()))
		}
		if !block.NetconfSSH.ConnectionLimit.IsNull() {
			configSet = append(configSet, setPrefix+"netconf ssh connection-limit "+
				utils.ConvI64toa(block.NetconfSSH.ConnectionLimit.ValueInt64()))
		}
		if !block.NetconfSSH.RateLimit.IsNull() {
			configSet = append(configSet, setPrefix+"netconf ssh rate-limit "+
				utils.ConvI64toa(block.NetconfSSH.RateLimit.ValueInt64()))
		}
	}
	if block.NetconfTraceoptions != nil {
		if block.NetconfTraceoptions.isEmpty() {
			return configSet, path.Root("services").AtName("netconf_traceoptions").AtName("*"),
				errors.New("netconf_traceoptions block in services block is empty")
		}
		if v := block.NetconfTraceoptions.FileName.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"netconf traceoptions file \""+v+"\"")
		}
		if !block.NetconfTraceoptions.FileFiles.IsNull() {
			configSet = append(configSet, setPrefix+"netconf traceoptions file files "+
				utils.ConvI64toa(block.NetconfTraceoptions.FileFiles.ValueInt64()))
		}
		if v := block.NetconfTraceoptions.FileMatch.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"netconf traceoptions file match \""+v+"\"")
		}
		if !block.NetconfTraceoptions.FileSize.IsNull() {
			configSet = append(configSet, setPrefix+"netconf traceoptions file size "+
				utils.ConvI64toa(block.NetconfTraceoptions.FileSize.ValueInt64()))
		}
		if block.NetconfTraceoptions.FileWorldReadable.ValueBool() {
			configSet = append(configSet, setPrefix+"netconf traceoptions file world-readable")
		}
		if block.NetconfTraceoptions.FileNoWorldReadable.ValueBool() {
			configSet = append(configSet, setPrefix+"netconf traceoptions file no-world-readable")
		}
		for _, v := range block.NetconfTraceoptions.Flag {
			configSet = append(configSet, setPrefix+"netconf traceoptions flag "+v.ValueString())
		}
		if block.NetconfTraceoptions.NoRemoteTrace.ValueBool() {
			configSet = append(configSet, setPrefix+"netconf traceoptions no-remote-trace")
		}
		if block.NetconfTraceoptions.OnDemand.ValueBool() {
			configSet = append(configSet, setPrefix+"netconf traceoptions on-demand")
		}
	}
	if block.SSH != nil {
		if block.SSH.isEmpty() {
			return configSet, path.Root("services").AtName("ssh").AtName("*"),
				errors.New("ssh block in services block is empty")
		}
		for _, v := range block.SSH.AuthenticationOrder {
			configSet = append(configSet, setPrefix+"ssh authentication-order "+v.ValueString())
		}
		for _, v := range block.SSH.Ciphers {
			configSet = append(configSet, setPrefix+"ssh ciphers \""+v.ValueString()+"\"")
		}
		if !block.SSH.ClientAliveCountMax.IsNull() {
			configSet = append(configSet, setPrefix+"ssh client-alive-count-max "+
				utils.ConvI64toa(block.SSH.ClientAliveCountMax.ValueInt64()))
		}
		if !block.SSH.ClientAliveInterval.IsNull() {
			configSet = append(configSet, setPrefix+"ssh client-alive-interval "+
				utils.ConvI64toa(block.SSH.ClientAliveInterval.ValueInt64()))
		}
		if !block.SSH.ConnectionLimit.IsNull() {
			configSet = append(configSet, setPrefix+"ssh connection-limit "+
				utils.ConvI64toa(block.SSH.ConnectionLimit.ValueInt64()))
		}
		if v := block.SSH.FingerprintHash.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"ssh fingerprint-hash "+v)
		}
		for _, v := range block.SSH.HostkeyAlgorithm {
			configSet = append(configSet, setPrefix+"ssh hostkey-algorithm \""+v.ValueString()+"\"")
		}
		for _, v := range block.SSH.KeyExchange {
			configSet = append(configSet, setPrefix+"ssh key-exchange \""+v.ValueString()+"\"")
		}
		if block.SSH.LogKeyChanges.ValueBool() {
			configSet = append(configSet, setPrefix+"ssh log-key-changes")
		}
		for _, v := range block.SSH.Macs {
			configSet = append(configSet, setPrefix+"ssh macs \""+v.ValueString()+"\"")
		}
		if !block.SSH.MaxPreAuthenticationPackets.IsNull() {
			configSet = append(configSet, setPrefix+"ssh max-pre-authentication-packets "+
				utils.ConvI64toa(block.SSH.MaxPreAuthenticationPackets.ValueInt64()))
		}
		if !block.SSH.MaxSessionsPerConnection.IsNull() {
			configSet = append(configSet, setPrefix+"ssh max-sessions-per-connection "+
				utils.ConvI64toa(block.SSH.MaxSessionsPerConnection.ValueInt64()))
		}
		if block.SSH.NoPasswords.ValueBool() {
			configSet = append(configSet, setPrefix+"ssh no-passwords")
		}
		if block.SSH.NoPublicKeys.ValueBool() {
			configSet = append(configSet, setPrefix+"ssh no-public-keys")
		}
		if !block.SSH.Port.IsNull() {
			configSet = append(configSet, setPrefix+"ssh port "+
				utils.ConvI64toa(block.SSH.Port.ValueInt64()))
		}
		for _, v := range block.SSH.ProtocolVersion {
			configSet = append(configSet, setPrefix+"ssh protocol-version "+v.ValueString())
		}
		if !block.SSH.RateLimit.IsNull() {
			configSet = append(configSet, setPrefix+"ssh rate-limit "+
				utils.ConvI64toa(block.SSH.RateLimit.ValueInt64()))
		}
		if v := block.SSH.RootLogin.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"ssh root-login "+v)
		}
		if block.SSH.TCPForwarding.ValueBool() {
			configSet = append(configSet, setPrefix+"ssh tcp-forwarding")
		}
		if block.SSH.NoTCPForwarding.ValueBool() {
			configSet = append(configSet, setPrefix+"ssh no-tcp-forwarding")
		}
	}
	if block.WebManagementHTTP != nil {
		configSet = append(configSet, setPrefix+"web-management http")
		for _, v := range block.WebManagementHTTP.Interface {
			configSet = append(configSet, setPrefix+"web-management http interface "+v.ValueString())
		}
		if !block.WebManagementHTTP.Port.IsNull() {
			configSet = append(configSet, setPrefix+"web-management http port "+
				utils.ConvI64toa(block.WebManagementHTTP.Port.ValueInt64()))
		}
	}
	if block.WebManagementHTTPS != nil {
		for _, v := range block.WebManagementHTTPS.Interface {
			configSet = append(configSet, setPrefix+"web-management https interface "+v.ValueString())
		}
		if v := block.WebManagementHTTPS.LocalCertificate.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"web-management https local-certificate \""+v+"\"")
		}
		if v := block.WebManagementHTTPS.PkiLocalCertificate.ValueString(); v != "" {
			configSet = append(configSet,
				setPrefix+"web-management https pki-local-certificate \""+v+"\"")
		}
		if !block.WebManagementHTTPS.Port.IsNull() {
			configSet = append(configSet, setPrefix+"web-management https port "+
				utils.ConvI64toa(block.WebManagementHTTPS.Port.ValueInt64()))
		}
		if block.WebManagementHTTPS.SystemGeneratedCertificate.ValueBool() {
			configSet = append(configSet, setPrefix+"web-management https system-generated-certificate")
		}
	}

	return configSet, path.Empty(), nil
}

func (block *systemBlockSyslog) configSet() (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix := "set system syslog "

	if !block.LogRotateFrequency.IsNull() {
		configSet = append(configSet, setPrefix+"log-rotate-frequency "+
			utils.ConvI64toa(block.LogRotateFrequency.ValueInt64()))
	}
	if v := block.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}
	if block.TimeFormatMillisecond.ValueBool() {
		configSet = append(configSet, setPrefix+"time-format millisecond")
	}
	if block.TimeFormatYear.ValueBool() {
		configSet = append(configSet, setPrefix+"time-format year")
	}

	if block.Archive != nil {
		configSet = append(configSet, setPrefix+"archive")
		if block.Archive.BinaryData.ValueBool() {
			configSet = append(configSet, setPrefix+"archive binary-data")
		}
		if block.Archive.NoBinaryData.ValueBool() {
			configSet = append(configSet, setPrefix+"archive no-binary-data")
		}
		if !block.Archive.Files.IsNull() {
			configSet = append(configSet, setPrefix+"archive files "+
				utils.ConvI64toa(block.Archive.Files.ValueInt64()))
		}
		if !block.Archive.Size.IsNull() {
			configSet = append(configSet, setPrefix+"archive size "+
				utils.ConvI64toa(block.Archive.Size.ValueInt64()))
		}
		if block.Archive.WorldReadable.ValueBool() {
			configSet = append(configSet, setPrefix+"archive world-readable")
		}
		if block.Archive.NoWorldReadable.ValueBool() {
			configSet = append(configSet, setPrefix+"archive no-world-readable")
		}
	}
	if block.Console != nil {
		if block.Console.isEmpty() {
			return configSet, path.Root("syslog").AtName("console").AtName("*"),
				errors.New("console block in syslog block is empty")
		}
		if v := block.Console.AnySeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console any "+v)
		}
		if v := block.Console.AuthorizationSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console authorization "+v)
		}
		if v := block.Console.ChangelogSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console change-log "+v)
		}
		if v := block.Console.ConflictlogSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console conflict-log "+v)
		}
		if v := block.Console.DaemonSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console daemon "+v)
		}
		if v := block.Console.DfcSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console dfc "+v)
		}
		if v := block.Console.ExternalSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console external "+v)
		}
		if v := block.Console.FirewallSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console firewall "+v)
		}
		if v := block.Console.FtpSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console ftp "+v)
		}
		if v := block.Console.InteractivecommandsSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console interactive-commands "+v)
		}
		if v := block.Console.KernelSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console kernel "+v)
		}
		if v := block.Console.NtpSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console ntp "+v)
		}
		if v := block.Console.PfeSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console pfe "+v)
		}
		if v := block.Console.SecuritySeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console security "+v)
		}
		if v := block.Console.UserSeverity.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"console user "+v)
		}
	}

	return configSet, path.Empty(), nil
}

func (rscData *systemData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "authentication-order "):
				rscData.AuthenticationOrder = append(rscData.AuthenticationOrder, types.StringValue(itemTrim))
			case itemTrim == "auto-snapshot":
				rscData.AutoSnapshot = types.BoolValue(true)
			case itemTrim == "default-address-selection":
				rscData.DefaultAddressSelection = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "domain-name "):
				rscData.DomainName = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "host-name "):
				rscData.HostName = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "max-configuration-rollbacks "):
				rscData.MaxConfigurationRollbacks, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "max-configurations-on-flash "):
				rscData.MaxConfigurationsOnFlash, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "name-server "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var nameServerOpts systemBlockNameServerOpts
				rscData.NameServerOpts, nameServerOpts = tfdata.ExtractBlockWithTFTypesString(
					rscData.NameServerOpts, "Address", itemTrimFields[0],
				)
				nameServerOpts.Address = types.StringValue(itemTrimFields[0])
				rscData.NameServer = append(rscData.NameServer, types.StringValue(itemTrimFields[0]))
				if balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" routing-instance ") {
					nameServerOpts.RoutingInstance = types.StringValue(itemTrim)
				}
				rscData.NameServerOpts = append(rscData.NameServerOpts, nameServerOpts)
			case itemTrim == "no-multicast-echo":
				rscData.NoMulticastEcho = types.BoolValue(true)
			case itemTrim == "no-ping-record-route":
				rscData.NoPingRecordRoute = types.BoolValue(true)
			case itemTrim == "no-ping-time-stamp":
				rscData.NoPingTimestamp = types.BoolValue(true)
			case itemTrim == "no-redirects":
				rscData.NoRedirects = types.BoolValue(true)
			case itemTrim == "no-redirects-ipv6":
				rscData.NoRedirectsIPv6 = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "radius-options attributes nas-id "):
				rscData.RadiusOptionsAttributesNasID = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "radius-options attributes nas-ip-address "):
				rscData.RadiusOptionsAttributesNasIpaddress = types.StringValue(itemTrim)
			case itemTrim == "radius-options enhanced-accounting":
				rscData.RadiusOptionsEnhancedAccounting = types.BoolValue(true)
			case itemTrim == "radius-options password-protocol mschap-v2":
				rscData.RadiusOptionsPasswordProtocolMschapv2 = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "tacplus-options authorization-time-interval "):
				rscData.TacplusOptionsAuthorizationTimeInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "tacplus-options enhanced-accounting":
				rscData.TacplusOptionsEnhancedAccounting = types.BoolValue(true)
			case itemTrim == "tacplus-options exclude-cmd-attribute":
				rscData.TacplusOptionsExcludeCmdAttribute = types.BoolValue(true)
			case itemTrim == "tacplus-options no-cmd-attribute-value":
				rscData.TacplusOptionsNoCmdAttributeValue = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "tacplus-options service-name "):
				rscData.TacplusOptionsServiceName = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "tacplus-options strict-authorization":
				rscData.TacplusOptionsStrictAuthorization = types.BoolValue(true)
			case itemTrim == "tacplus-options no-strict-authorization":
				rscData.TacplusOptionsNoStrictAuthorization = types.BoolValue(true)
			case itemTrim == "tacplus-options timestamp-and-timezone":
				rscData.TacplusOptionsTimestampAndTimezone = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "time-zone "):
				rscData.TimeZone = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "tracing destination-override syslog host "):
				rscData.TracingDestOverrideSyslogHost = types.StringValue(itemTrim)

			case balt.CutPrefixInString(&itemTrim, "accounting "):
				if rscData.Accounting == nil {
					rscData.Accounting = &systemBlockAccounting{}
				}
				if err := rscData.Accounting.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "archival configuration "):
				if rscData.ArchivalConfiguration == nil {
					rscData.ArchivalConfiguration = &systemBlockArchivalConfiguration{}
				}
				if err := rscData.ArchivalConfiguration.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "inet6-backup-router "):
				if rscData.Inet6BackupRouter == nil {
					rscData.Inet6BackupRouter = &systemBlockInet6BackupRouter{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "destination "):
					rscData.Inet6BackupRouter.Destination = append(rscData.Inet6BackupRouter.Destination, types.StringValue(itemTrim))
				default:
					rscData.Inet6BackupRouter.Address = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "internet-options "):
				if rscData.InternetOptions == nil {
					rscData.InternetOptions = &systemBlockInternetOptions{}
				}
				if err := rscData.InternetOptions.read(itemTrim); err != nil {
					return err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, systemBlockLicense{}.junosLines()):
				if rscData.License == nil {
					rscData.License = &systemBlockLicense{}
				}
				if err := rscData.License.read(itemTrim); err != nil {
					return err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, systemBlockLogin{}.junosLines()):
				if rscData.Login == nil {
					rscData.Login = &systemBlockLogin{}
				}
				if err := rscData.Login.read(itemTrim); err != nil {
					return err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, systemBlockNtp{}.junosLines()):
				if rscData.Ntp == nil {
					rscData.Ntp = &systemBlockNtp{}
				}
				if err := rscData.Ntp.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "ports "):
				if rscData.Ports == nil {
					rscData.Ports = &systemBlockPorts{}
				}
				if err := rscData.Ports.read(itemTrim); err != nil {
					return err
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, systemBlockServices{}.junosLines()):
				if rscData.Services == nil {
					rscData.Services = &systemBlockServices{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "services web-management session idle-timeout "):
					rscData.Services.WebManagementSessionIdleTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "services web-management session session-limit "):
					rscData.Services.WebManagementSessionLimit, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case bchk.StringHasOneOfPrefixes(itemTrim, systemBlockServicesBlockNetconfSSH{}.junosLines()):
					if rscData.Services.NetconfSSH == nil {
						rscData.Services.NetconfSSH = &systemBlockServicesBlockNetconfSSH{}
					}
					if err := rscData.Services.NetconfSSH.read(itemTrim); err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "services netconf traceoptions "):
					if rscData.Services.NetconfTraceoptions == nil {
						rscData.Services.NetconfTraceoptions = &systemBlockServicesBlockNetconfTraceoptions{}
					}
					if err := rscData.Services.NetconfTraceoptions.read(itemTrim); err != nil {
						return err
					}
				case bchk.StringHasOneOfPrefixes(itemTrim, systemBlockServicesBlockSSH{}.junosLines()):
					if rscData.Services.SSH == nil {
						rscData.Services.SSH = &systemBlockServicesBlockSSH{}
					}
					if err := rscData.Services.SSH.read(itemTrim); err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "services web-management https"):
					if rscData.Services.WebManagementHTTPS == nil {
						rscData.Services.WebManagementHTTPS = &systemBlockServicesBlockWebManagementHTTPS{}
					}
					if err := rscData.Services.WebManagementHTTPS.read(itemTrim); err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "services web-management http"):
					if rscData.Services.WebManagementHTTP == nil {
						rscData.Services.WebManagementHTTP = &systemBlockServicesBlockWebManagementHTTP{}
					}
					if err := rscData.Services.WebManagementHTTP.read(itemTrim); err != nil {
						return err
					}
				}
			case bchk.StringHasOneOfPrefixes(itemTrim, systemBlockSyslog{}.junosLines()):
				if rscData.Syslog == nil {
					rscData.Syslog = &systemBlockSyslog{}
				}
				if err := rscData.Syslog.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *systemBlockAccounting) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "events "):
		block.Events = append(block.Events, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "enhanced-avs-max "):
		block.EnhancedAvsMax, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "destination radius"):
		block.DestinationRadius = types.BoolValue(true)
		if balt.CutPrefixInString(&itemTrim, " server ") {
			itemTrimFields := strings.Split(itemTrim, " ")
			var destinationRadiusServer systemBlockAccountingBlockDestinationRadiusServer
			block.DestinationRadiusServer, destinationRadiusServer = tfdata.ExtractBlockWithTFTypesString(
				block.DestinationRadiusServer, "Address", itemTrimFields[0],
			)
			destinationRadiusServer.Address = types.StringValue(itemTrimFields[0])
			balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
			if err := destinationRadiusServer.read(itemTrim); err != nil {
				return err
			}
			block.DestinationRadiusServer = append(block.DestinationRadiusServer, destinationRadiusServer)
		}
	case balt.CutPrefixInString(&itemTrim, "destination tacplus"):
		block.DestinationTacplus = types.BoolValue(true)
		if balt.CutPrefixInString(&itemTrim, " server ") {
			itemTrimFields := strings.Split(itemTrim, " ")
			var destinationTacplusServer systemBlockAccountingBlockDestinationTacplusServer
			block.DestinationTacplusServer, destinationTacplusServer = tfdata.ExtractBlockWithTFTypesString(
				block.DestinationTacplusServer, "Address", itemTrimFields[0],
			)
			destinationTacplusServer.Address = types.StringValue(itemTrimFields[0])
			balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
			if err := destinationTacplusServer.read(itemTrim); err != nil {
				return err
			}
			block.DestinationTacplusServer = append(block.DestinationTacplusServer, destinationTacplusServer)
		}
	}

	return nil
}

func (block *systemBlockAccountingBlockDestinationRadiusServer) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "secret "):
		block.Secret, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "secret")
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "accounting-port "):
		block.AccountingPort, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "accounting-retry "):
		block.AccountingRetry, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "accounting-timeout "):
		block.AccountingTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "dynamic-request-port "):
		block.DynamicRequestPort, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "max-outstanding-requests "):
		block.MaxOutstandingRequests, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "port "):
		block.Port, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "preauthentication-port "):
		block.PreauthenticationPort, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "preauthentication-secret "):
		block.PreauthenticationSecret, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "preauthentication-secret")
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "retry "):
		block.Retry, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "routing-instance "):
		block.RoutingInstance = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		block.SourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "timeout "):
		block.Timeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (block *systemBlockAccountingBlockDestinationTacplusServer) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "port "):
		block.Port, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "routing-instance "):
		block.RoutingInstance = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "secret "):
		block.Secret, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "secret")
		if err != nil {
			return err
		}
	case itemTrim == "single-connection":
		block.SingleConnection = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		block.SourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "timeout "):
		block.Timeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (block *systemBlockArchivalConfiguration) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "archive-sites "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) > 2 { // <url> password <password>
			password, err := tfdata.JunosDecode(strings.Trim(itemTrimFields[2], "\""), "password")
			if err != nil {
				return err
			}
			block.ArchiveSite = append(block.ArchiveSite, systemBlockArchivalConfigurationBlockArchiveSite{
				URL:      types.StringValue(strings.Trim(itemTrimFields[0], "\"")),
				Password: password,
			})
		} else { // <url>
			block.ArchiveSite = append(block.ArchiveSite, systemBlockArchivalConfigurationBlockArchiveSite{
				URL: types.StringValue(strings.Trim(itemTrimFields[0], "\"")),
			})
		}
	case balt.CutPrefixInString(&itemTrim, "transfer-interval "):
		block.TransferInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "transfer-on-commit":
		block.TransferOnCommit = types.BoolValue(true)
	}

	return nil
}

func (block *systemBlockInternetOptions) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "gre-path-mtu-discovery":
		block.GrePathMtuDiscovery = types.BoolValue(true)
	case itemTrim == "no-gre-path-mtu-discovery":
		block.NoGrePathMtuDiscovery = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "icmpv4-rate-limit"):
		if block.IcmpV4RateLimit == nil {
			block.IcmpV4RateLimit = &systemBlockInternetOptionsBlockIcmpRateLimit{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, " bucket-size "):
			block.IcmpV4RateLimit.BucketSize, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " packet-rate "):
			block.IcmpV4RateLimit.PacketRate, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "icmpv6-rate-limit"):
		if block.IcmpV6RateLimit == nil {
			block.IcmpV6RateLimit = &systemBlockInternetOptionsBlockIcmpRateLimit{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, " bucket-size "):
			block.IcmpV6RateLimit.BucketSize, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " packet-rate "):
			block.IcmpV6RateLimit.PacketRate, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case itemTrim == "ipip-path-mtu-discovery":
		block.IpipPathMtuDiscovery = types.BoolValue(true)
	case itemTrim == "no-ipip-path-mtu-discovery":
		block.NoIpipPathMtuDiscovery = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "ipv6-duplicate-addr-detection-transmits "):
		block.IPv6DuplicateAddrDetectionTransmits, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "ipv6-path-mtu-discovery":
		block.IPv6PathMtuDiscovery = types.BoolValue(true)
	case itemTrim == "no-ipv6-path-mtu-discovery":
		block.NoIPv6PathMtuDiscovery = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "ipv6-path-mtu-discovery-timeout "):
		block.IPv6PathMtuDiscoveryTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "ipv6-reject-zero-hop-limit":
		block.IPv6RejectZeroHopLimit = types.BoolValue(true)
	case itemTrim == "no-ipv6-reject-zero-hop-limit":
		block.NoIPv6RejectZeroHopLimit = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "no-tcp-reset "):
		block.NoTCPReset = types.StringValue(itemTrim)
	case itemTrim == "no-tcp-rfc1323":
		block.NoTCPRFC1323 = types.BoolValue(true)
	case itemTrim == "no-tcp-rfc1323-paws":
		block.NoTCPRFC1323Paws = types.BoolValue(true)
	case itemTrim == "path-mtu-discovery":
		block.PathMtuDiscovery = types.BoolValue(true)
	case itemTrim == "no-path-mtu-discovery":
		block.NoPathMtuDiscovery = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "source-port upper-limit "):
		block.SourcePortUpperLimit, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "source-quench":
		block.SourceQuench = types.BoolValue(true)
	case itemTrim == "no-source-quench":
		block.NoSourceQuench = types.BoolValue(true)
	case itemTrim == "tcp-drop-synfin-set":
		block.TCPDropSynfinSet = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "tcp-mss "):
		block.TCPMss, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (systemBlockLicense) junosLines() []string {
	return []string{
		"license autoupdate",
		"license renew",
	}
}

func (block *systemBlockLicense) read(itemTrim string) (err error) {
	itemTrim = strings.TrimPrefix(itemTrim, "license ")
	switch {
	case itemTrim == "autoupdate":
		block.Autoupdate = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "autoupdate url "):
		block.Autoupdate = types.BoolValue(true)
		url := tfdata.FirstElementOfJunosLine(itemTrim)
		block.AutoupdateURL = types.StringValue(strings.Trim(url, "\""))

		if balt.CutPrefixInString(&itemTrim, url+" password ") {
			password, err := tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "password")
			if err != nil {
				return err
			}
			block.AutoupdatePassword = password
		}
	case balt.CutPrefixInString(&itemTrim, "renew before-expiration "):
		block.RenewBeforeExpiration, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "renew interval "):
		block.RenewInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (systemBlockLogin) junosLines() []string {
	return []string{
		"login announcement",
		"login deny-sources",
		"login idle-timeout",
		"login message",
		"login password",
		"login retry-options",
	}
}

func (block *systemBlockLogin) read(itemTrim string) (err error) {
	itemTrim = strings.TrimPrefix(itemTrim, "login ")
	switch {
	case balt.CutPrefixInString(&itemTrim, "announcement "):
		block.Announcement = types.StringValue(html.UnescapeString(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "deny-sources address "):
		block.DenySourcesAddress = append(block.DenySourcesAddress, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "idle-timeout "):
		block.IdleTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "message "):
		block.Message = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "password "):
		if block.Password == nil {
			block.Password = &systemBlockLoginBlockPassword{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, "change-type "):
			block.Password.ChangeType = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "format "):
			block.Password.Format = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "maximum-length "):
			block.Password.MaximumLength, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-changes "):
			block.Password.MinimumChanges, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-character-changes "):
			block.Password.MinimumCharacterChanges, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-length "):
			block.Password.MinimumLength, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-lower-cases "):
			block.Password.MinimumLowerCases, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-numerics "):
			block.Password.MinimumNumerics, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-punctuations "):
			block.Password.MinimumPunctuations, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-reuse "):
			block.Password.MinimumReuse, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-upper-cases "):
			block.Password.MinimumUpperCases, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	case balt.CutPrefixInString(&itemTrim, "retry-options "):
		if block.RetryOptions == nil {
			block.RetryOptions = &systemBlockLoginBlockRetryOptions{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, "backoff-factor "):
			block.RetryOptions.BackoffFactor, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "backoff-threshold "):
			block.RetryOptions.BackoffThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "lockout-period "):
			block.RetryOptions.LockoutPeriod, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "maximum-time "):
			block.RetryOptions.MaximumTime, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "minimum-time "):
			block.RetryOptions.MinimumTime, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "tries-before-disconnect "):
			block.RetryOptions.TriesBeforeDisconnect, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (systemBlockNtp) junosLines() []string {
	return []string{
		"ntp boot-server",
		"ntp broadcast-client",
		"ntp interval-range",
		"ntp multicast-client",
		"ntp source-address",
		"ntp threshold",
	}
}

func (block *systemBlockNtp) read(itemTrim string) (err error) {
	itemTrim = strings.TrimPrefix(itemTrim, "ntp ")
	switch {
	case balt.CutPrefixInString(&itemTrim, "boot-server "):
		block.BootServer = types.StringValue(itemTrim)
	case itemTrim == "broadcast-client":
		block.BroadcastClient = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "interval-range "):
		block.IntervalRange, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "multicast-client"):
		block.MulticastClient = types.BoolValue(true)
		if balt.CutPrefixInString(&itemTrim, " ") {
			block.MulticastClientAddress = types.StringValue(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "threshold action "):
		block.ThresholdAction = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "threshold "):
		block.ThresholdValue, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (block *systemBlockPorts) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "auxiliary authentication-order "):
		block.AuxiliaryAuthenticationOrder = append(block.AuxiliaryAuthenticationOrder, types.StringValue(itemTrim))
	case itemTrim == "auxiliary disable":
		block.AuxiliaryDisable = types.BoolValue(true)
	case itemTrim == "auxiliary insecure":
		block.AuxiliaryInsecure = types.BoolValue(true)
	case itemTrim == "auxiliary log-out-on-disconnect":
		block.AuxiliaryLogoutOnDisconnect = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "auxiliary type "):
		block.AuxiliaryType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "console authentication-order "):
		block.ConsoleAuthenticationOrder = append(block.ConsoleAuthenticationOrder, types.StringValue(itemTrim))
	case itemTrim == "console disable":
		block.ConsoleDisable = types.BoolValue(true)
	case itemTrim == "console insecure":
		block.ConsoleInsecure = types.BoolValue(true)
	case itemTrim == "console log-out-on-disconnect":
		block.ConsoleLogoutOnDisconnect = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "console type "):
		block.ConsoleType = types.StringValue(itemTrim)
	}

	return nil
}

func (systemBlockServices) junosLines() []string {
	s := make([]string, 0, 50)
	s = append(s, "services web-management session idle-timeout")
	s = append(s, "services web-management session session-limit")
	s = append(s, systemBlockServicesBlockNetconfSSH{}.junosLines()...)
	s = append(s, "services netconf traceoptions")
	s = append(s, systemBlockServicesBlockSSH{}.junosLines()...)
	s = append(s, "services web-management http")
	s = append(s, "services web-management https")

	return s
}

func (systemBlockServicesBlockNetconfSSH) junosLines() []string {
	return []string{
		"services netconf ssh client-alive-count-max",
		"services netconf ssh client-alive-interval",
		"services netconf ssh connection-limit",
		"services netconf ssh rate-limit",
	}
}

func (block *systemBlockServicesBlockNetconfSSH) read(itemTrim string) (err error) {
	itemTrim = strings.TrimPrefix(itemTrim, "services netconf ssh ")
	switch {
	case balt.CutPrefixInString(&itemTrim, "client-alive-count-max "):
		block.ClientAliveCountMax, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "client-alive-interval "):
		block.ClientAliveInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "connection-limit "):
		block.ConnectionLimit, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "rate-limit "):
		block.RateLimit, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (block *systemBlockServicesBlockNetconfTraceoptions) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "file files "):
		block.FileFiles, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "file match "):
		block.FileMatch = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "file size "):
		switch {
		case balt.CutSuffixInString(&itemTrim, "k"):
			block.FileSize, err = tfdata.ConvAtoi64Value(itemTrim)
			block.FileSize = types.Int64Value(block.FileSize.ValueInt64() * 1024)
		case balt.CutSuffixInString(&itemTrim, "m"):
			block.FileSize, err = tfdata.ConvAtoi64Value(itemTrim)
			block.FileSize = types.Int64Value(block.FileSize.ValueInt64() * 1024 * 1024)
		case balt.CutSuffixInString(&itemTrim, "g"):
			block.FileSize, err = tfdata.ConvAtoi64Value(itemTrim)
			block.FileSize = types.Int64Value(block.FileSize.ValueInt64() * 1024 * 1024 * 1024)
		default:
			block.FileSize, err = tfdata.ConvAtoi64Value(itemTrim)
		}
		if err != nil {
			return err
		}
	case itemTrim == "file world-readable":
		block.FileWorldReadable = types.BoolValue(true)
	case itemTrim == "file no-world-readable":
		block.FileNoWorldReadable = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "file "):
		block.FileName = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "flag "):
		block.Flag = append(block.Flag, types.StringValue(itemTrim))
	case itemTrim == "no-remote-trace":
		block.NoRemoteTrace = types.BoolValue(true)
	case itemTrim == "on-demand":
		block.OnDemand = types.BoolValue(true)
	}

	return nil
}

func (systemBlockServicesBlockSSH) junosLines() []string {
	return []string{
		"services ssh authentication-order",
		"services ssh ciphers",
		"services ssh client-alive-count-max",
		"services ssh client-alive-interval",
		"services ssh connection-limit",
		"services ssh fingerprint-hash",
		"services ssh hostkey-algorithm",
		"services ssh key-exchange",
		"services ssh log-key-changes",
		"services ssh macs",
		"services ssh max-pre-authentication-packets",
		"services ssh max-sessions-per-connection",
		"services ssh no-passwords",
		"services ssh no-public-keys",
		"services ssh port",
		"services ssh protocol-version",
		"services ssh rate-limit",
		"services ssh root-login",
		"services ssh no-tcp-forwarding",
		"services ssh tcp-forwarding",
	}
}

func (block *systemBlockServicesBlockSSH) read(itemTrim string) (err error) {
	itemTrim = strings.TrimPrefix(itemTrim, "services ssh ")
	switch {
	case balt.CutPrefixInString(&itemTrim, "authentication-order "):
		block.AuthenticationOrder = append(block.AuthenticationOrder, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "ciphers "):
		block.Ciphers = append(block.Ciphers, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "client-alive-count-max "):
		block.ClientAliveCountMax, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "client-alive-interval "):
		block.ClientAliveInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "connection-limit "):
		block.ConnectionLimit, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "fingerprint-hash "):
		block.FingerprintHash = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "hostkey-algorithm "):
		block.HostkeyAlgorithm = append(block.HostkeyAlgorithm, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "key-exchange "):
		block.KeyExchange = append(block.KeyExchange, types.StringValue(strings.Trim(itemTrim, "\"")))
	case itemTrim == "log-key-changes":
		block.LogKeyChanges = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "macs "):
		block.Macs = append(block.Macs, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "max-pre-authentication-packets "):
		block.MaxPreAuthenticationPackets, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "max-sessions-per-connection "):
		block.MaxSessionsPerConnection, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "no-passwords":
		block.NoPasswords = types.BoolValue(true)
	case itemTrim == "no-public-keys":
		block.NoPublicKeys = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "port "):
		block.Port, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "protocol-version "):
		block.ProtocolVersion = append(block.ProtocolVersion, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "rate-limit "):
		block.RateLimit, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "root-login "):
		block.RootLogin = types.StringValue(itemTrim)
	case itemTrim == "tcp-forwarding":
		block.TCPForwarding = types.BoolValue(true)
	case itemTrim == "no-tcp-forwarding":
		block.NoTCPForwarding = types.BoolValue(true)
	}

	return nil
}

func (block *systemBlockServicesBlockWebManagementHTTP) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, " interface "):
		block.Interface = append(block.Interface, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, " port "):
		block.Port, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (block *systemBlockServicesBlockWebManagementHTTPS) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, " interface "):
		block.Interface = append(block.Interface, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, " local-certificate "):
		block.LocalCertificate = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, " pki-local-certificate "):
		block.PkiLocalCertificate = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, " port "):
		block.Port, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == " system-generated-certificate":
		block.SystemGeneratedCertificate = types.BoolValue(true)
	}

	return nil
}

func (systemBlockSyslog) junosLines() []string {
	return []string{
		"syslog archive",
		"syslog console ",
		"syslog log-rotate-frequency",
		"syslog source-address",
		"syslog time-format ",
	}
}

func (block *systemBlockSyslog) read(itemTrim string) (err error) {
	itemTrim = strings.TrimPrefix(itemTrim, "syslog ")
	switch {
	case balt.CutPrefixInString(&itemTrim, "log-rotate-frequency "):
		block.LogRotateFrequency, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		block.SourceAddress = types.StringValue(itemTrim)
	case itemTrim == "time-format millisecond":
		block.TimeFormatMillisecond = types.BoolValue(true)
	case itemTrim == "time-format year":
		block.TimeFormatYear = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "archive"):
		if block.Archive == nil {
			block.Archive = &systemBlockSyslogBlockArchive{}
		}
		switch {
		case itemTrim == " binary-data":
			block.Archive.BinaryData = types.BoolValue(true)
		case itemTrim == " no-binary-data":
			block.Archive.NoBinaryData = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, " files "):
			block.Archive.Files, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, " size "):
			switch {
			case balt.CutSuffixInString(&itemTrim, "k"):
				block.Archive.Size, err = tfdata.ConvAtoi64Value(itemTrim)
				block.Archive.Size = types.Int64Value(block.Archive.Size.ValueInt64() * 1024)
			case balt.CutSuffixInString(&itemTrim, "m"):
				block.Archive.Size, err = tfdata.ConvAtoi64Value(itemTrim)
				block.Archive.Size = types.Int64Value(block.Archive.Size.ValueInt64() * 1024 * 1024)
			case balt.CutSuffixInString(&itemTrim, "g"):
				block.Archive.Size, err = tfdata.ConvAtoi64Value(itemTrim)
				block.Archive.Size = types.Int64Value(block.Archive.Size.ValueInt64() * 1024 * 1024 * 1024)
			default:
				block.Archive.Size, err = tfdata.ConvAtoi64Value(itemTrim)
			}
			if err != nil {
				return err
			}
		case itemTrim == " world-readable":
			block.Archive.WorldReadable = types.BoolValue(true)
		case itemTrim == " no-world-readable":
			block.Archive.NoWorldReadable = types.BoolValue(true)
		}
	case balt.CutPrefixInString(&itemTrim, "console "):
		if block.Console == nil {
			block.Console = &systemBlockSyslogBlockConsole{}
		}
		switch {
		case balt.CutPrefixInString(&itemTrim, "any "):
			block.Console.AnySeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "authorization "):
			block.Console.AuthorizationSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "change-log "):
			block.Console.ChangelogSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "conflict-log "):
			block.Console.ConflictlogSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "daemon "):
			block.Console.DaemonSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "dfc "):
			block.Console.DfcSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "external "):
			block.Console.ExternalSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "firewall "):
			block.Console.FirewallSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "ftp "):
			block.Console.FtpSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "interactive-commands "):
			block.Console.InteractivecommandsSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "kernel "):
			block.Console.KernelSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "ntp "):
			block.Console.NtpSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "pfe "):
			block.Console.PfeSeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "security "):
			block.Console.SecuritySeverity = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "user "):
			block.Console.UserSeverity = types.StringValue(itemTrim)
		}
	}

	return nil
}

func (rscData *systemData) del(
	_ context.Context, junSess *junos.Session,
) error {
	listLinesToDelete := make([]string, 0, 100)
	listLinesToDelete = append(listLinesToDelete, "accounting")
	listLinesToDelete = append(listLinesToDelete, "archival configuration")
	listLinesToDelete = append(listLinesToDelete, "authentication-order")
	listLinesToDelete = append(listLinesToDelete, "auto-snapshot")
	listLinesToDelete = append(listLinesToDelete, "default-address-selection")
	listLinesToDelete = append(listLinesToDelete, "domain-name")
	listLinesToDelete = append(listLinesToDelete, "host-name")
	listLinesToDelete = append(listLinesToDelete, "inet6-backup-router")
	listLinesToDelete = append(listLinesToDelete, "internet-options")
	listLinesToDelete = append(listLinesToDelete, systemBlockLicense{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, systemBlockLogin{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, "max-configuration-rollbacks")
	listLinesToDelete = append(listLinesToDelete, "max-configurations-on-flash")
	listLinesToDelete = append(listLinesToDelete, systemBlockNtp{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, "name-server")
	listLinesToDelete = append(listLinesToDelete, "no-multicast-echo")
	listLinesToDelete = append(listLinesToDelete, "no-ping-record-route")
	listLinesToDelete = append(listLinesToDelete, "no-ping-time-stamp")
	listLinesToDelete = append(listLinesToDelete, "no-redirects")
	listLinesToDelete = append(listLinesToDelete, "no-redirects-ipv6")
	listLinesToDelete = append(listLinesToDelete, "ports")
	listLinesToDelete = append(listLinesToDelete, "radius-options")
	listLinesToDelete = append(listLinesToDelete, systemBlockServices{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, systemBlockSyslog{}.junosLines()...)
	listLinesToDelete = append(listLinesToDelete, "tacplus-options")
	listLinesToDelete = append(listLinesToDelete, "time-zone")
	listLinesToDelete = append(listLinesToDelete,
		"tracing destination-override syslog host",
	)

	configSet := make([]string, len(listLinesToDelete))
	delPrefix := "delete system "
	for i, line := range listLinesToDelete {
		configSet[i] = delPrefix + line
	}

	return junSess.ConfigSet(configSet)
}
