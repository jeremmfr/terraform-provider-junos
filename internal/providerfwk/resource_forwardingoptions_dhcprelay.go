package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource                   = &forwardingoptionsDhcprelay{}
	_ resource.ResourceWithConfigure      = &forwardingoptionsDhcprelay{}
	_ resource.ResourceWithValidateConfig = &forwardingoptionsDhcprelay{}
	_ resource.ResourceWithImportState    = &forwardingoptionsDhcprelay{}
	_ resource.ResourceWithUpgradeState   = &forwardingoptionsDhcprelay{}
)

type forwardingoptionsDhcprelay struct {
	client *junos.Client
}

func newForwardingoptionsDhcprelayResource() resource.Resource {
	return &forwardingoptionsDhcprelay{}
}

func (rsc *forwardingoptionsDhcprelay) typeName() string {
	return providerName + "_forwardingoptions_dhcprelay"
}

func (rsc *forwardingoptionsDhcprelay) junosName() string {
	return "forwarding-options dhcp-relay"
}

func (rsc *forwardingoptionsDhcprelay) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *forwardingoptionsDhcprelay) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *forwardingoptionsDhcprelay) Configure(
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

func (rsc *forwardingoptionsDhcprelay) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Configure static configuration in `" + rsc.junosName() + " [dhcpv6]` block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<routing_instance>_-_<version>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Version for DHCP or DHCPv6.",
				Default:     stringdefault.StaticString("v4"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("v4", "v6"),
				},
			},
			"access_profile": schema.StringAttribute{
				Optional:    true,
				Description: "Access profile to use for AAA services.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"active_server_group": schema.StringAttribute{
				Optional:    true,
				Description: "Name of DHCP server group.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"active_server_group_allow_server_change": schema.BoolAttribute{
				Optional:    true,
				Description: "Accept DHCP-ACK from any server in this group.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"arp_inspection": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable Dynamic ARP Inspection.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"authentication_password": schema.StringAttribute{
				Optional:    true,
				Description: "DHCP authentication, username password to use.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"client_response_ttl": schema.Int64Attribute{
				Optional:    true,
				Description: "P time-to-live value to set in responses to client.",
				Validators: []validator.Int64{
					int64validator.Between(1, 255),
				},
			},
			"duplicate_clients_in_subnet": schema.StringAttribute{
				Optional:    true,
				Description: "Allow duplicate clients in a subnet.",
				Validators: []validator.String{
					stringvalidator.OneOf("incoming-interface", "option-82"),
				},
			},
			"duplicate_clients_incoming_interface": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow duplicate clients on different underlying interfaces.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"dynamic_profile": schema.StringAttribute{
				Optional:    true,
				Description: "Dynamic profile to use.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 80),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"dynamic_profile_aggregate_clients": schema.BoolAttribute{
				Optional:    true,
				Description: "Aggregate client profiles.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"dynamic_profile_aggregate_clients_action": schema.StringAttribute{
				Optional:    true,
				Description: "Merge or replace the client dynamic profiles.",
				Validators: []validator.String{
					stringvalidator.OneOf("merge", "replace"),
				},
			},
			"dynamic_profile_use_primary": schema.StringAttribute{
				Optional:    true,
				Description: "Dynamic profile to use on the primary interface.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 80),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"exclude_relay_agent_identifier": schema.BoolAttribute{
				Optional:    true,
				Description: "Exclude relay agent identifier from packets to server.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"forward_only": schema.BoolAttribute{
				Optional:    true,
				Description: "Forward DHCP packets without creating binding.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"forward_only_replies": schema.BoolAttribute{
				Optional:    true,
				Description: "Forward-only replies from server to appropriate routing instance based on options.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"forward_only_routing_instance": schema.StringAttribute{
				Optional:    true,
				Description: "Name of routing instance to forward-only.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"forward_snooped_clients": schema.StringAttribute{
				Optional:    true,
				Description: "Forward snooped (unicast) packets.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"all-interfaces",
						"configured-interfaces",
						"non-configured-interfaces",
					),
				},
			},
			"liveness_detection_failure_action": schema.StringAttribute{
				Optional:    true,
				Description: "Liveness detection failure action options.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"clear-binding",
						"clear-binding-if-interface-up",
						"log-only",
					),
				},
			},
			"maximum_hop_count": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of hops per packet.",
				Validators: []validator.Int64{
					int64validator.Between(1, 16),
				},
			},
			"minimum_wait_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Minimum number of seconds before requests are forwarded.",
				Validators: []validator.Int64{
					int64validator.Between(0, 30000),
				},
			},
			"no_snoop": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not snoop DHCP packets.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"persistent_storage_automatic": schema.BoolAttribute{
				Optional:    true,
				Description: "Trigger to enable flat file storage.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"relay_agent_option_79": schema.BoolAttribute{
				Optional:    true,
				Description: "Add the client MAC address to the Relay Forward header.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"remote_id_mismatch_disconnect": schema.BoolAttribute{
				Optional:    true,
				Description: "Disconnect session on remote-id mismatch.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"route_suppression_access": schema.BoolAttribute{
				Optional:    true,
				Description: "Suppress access route addition.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"route_suppression_access_internal": schema.BoolAttribute{
				Optional:    true,
				Description: "Suppress access-internal route addition.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"route_suppression_destination": schema.BoolAttribute{
				Optional:    true,
				Description: "Suppress destination route addition.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"server_match_default_action": schema.StringAttribute{
				Optional:    true,
				Description: "Server match default action.",
				Validators: []validator.String{
					stringvalidator.OneOf("create-relay-entry", "forward-only"),
				},
			},
			"server_response_time": schema.Int64Attribute{
				Optional: true,
				Description: "Number of seconds in a period of activity between" +
					" the last server response and an unanswered request.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"service_profile": schema.StringAttribute{
				Optional:    true,
				Description: "Dynamic profile to use for default service activation.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"short_cycle_protection_lockout_max_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Short cycle lockout max time in seconds.",
				Validators: []validator.Int64{
					int64validator.Between(1, 86400),
				},
			},
			"short_cycle_protection_lockout_min_time": schema.Int64Attribute{
				Optional:    true,
				Description: "hort cycle lockout min time in seconds.",
				Validators: []validator.Int64{
					int64validator.Between(1, 86400),
				},
			},
			"source_ip_change": schema.BoolAttribute{
				Optional:    true,
				Description: "Use address of egress interface as source ip.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"vendor_specific_information_host_name": schema.BoolAttribute{
				Optional:    true,
				Description: "DHCPv6 option 17 vendor-specific processing, add router host name.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"vendor_specific_information_location": schema.BoolAttribute{
				Optional: true,
				Description: "DHCPv6 option 17 vendor-specific processing," +
					" add location information expressed as interface name format.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"active_leasequery": schema.SingleNestedBlock{
				Description: "DHCP active leasequery configuration.",
				Attributes: map[string]schema.Attribute{
					"idle_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Idle timeout in seconds (seconds).",
						Validators: []validator.Int64{
							int64validator.Between(10, 3600),
						},
					},
					"peer_address": schema.StringAttribute{
						Optional:    true,
						Description: "Server ip address.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Read/write timeout in seconds (seconds).",
						Validators: []validator.Int64{
							int64validator.Between(10, 3600),
						},
					},
					"topology_discover": schema.BoolAttribute{
						Optional:    true,
						Description: "Topology discovery.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"authentication_username_include": schema.SingleNestedBlock{
				Description: "DHCP authentication, add username options.",
				Attributes:  dhcpBlockAuthenticationUsernameInclude{}.attributesSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"bulk_leasequery": schema.SingleNestedBlock{
				Description: "DHCP bulk leasequery configuration.",
				Attributes: map[string]schema.Attribute{
					"attempts": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of retry attempts.",
						Validators: []validator.Int64{
							int64validator.Between(1, 720),
						},
					},
					"timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of seconds.",
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
						},
					},
					"trigger_automatic": schema.BoolAttribute{
						Optional:    true,
						Description: "Trigger automatically for bulk leasequery.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"lease_time_validation": schema.SingleNestedBlock{
				Description: "Configure lease time violation validation.",
				Attributes:  forwardingoptionsDhcprelayBlockLeaseTimeValidation{}.attributesSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"leasequery": schema.SingleNestedBlock{
				Description: "DHCP leasequery configuration.",
				Attributes: map[string]schema.Attribute{
					"attempts": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of retry attempts.",
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
						},
					},
					"timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of seconds.",
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"liveness_detection_method_bfd": schema.SingleNestedBlock{
				Description: "Liveness detection method BFD options.",
				Attributes:  dhcpBlockLivenessDetectionMethodBfd{}.attributesSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"liveness_detection_method_layer2": schema.SingleNestedBlock{
				Description: "Liveness detection method address resolution options.",
				Attributes:  dhcpBlockLivenessDetectionMethodLayer2{}.attributesSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"overrides_v4": schema.SingleNestedBlock{
				Description: "DHCP override processing.",
				Attributes:  forwardingoptionsDhcprelayBlockOverridesV4{}.attributesSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"overrides_v6": schema.SingleNestedBlock{
				Description: "DHCPv6 override processing.",
				Attributes:  forwardingoptionsDhcprelayBlockOverridesV6{}.attributesSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"relay_agent_interface_id": schema.SingleNestedBlock{
				Description: "DHCPv6 interface-id option processing.",
				Attributes:  forwardingoptionsDhcprelayBlockRelayAgentInterfaceID{}.attributesSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"relay_agent_remote_id": schema.SingleNestedBlock{
				Description: "DHCPv6 remote-id option processing.",
				Attributes:  forwardingoptionsDhcprelayBlockRelayAgentRemoteID{}.attributesSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"relay_option": schema.SingleNestedBlock{
				Description: "DHCP option processing.",
				Attributes:  forwardingoptionsDhcprelayBlockRelayOption{}.attributesSchema(),
				Blocks:      forwardingoptionsDhcprelayBlockRelayOption{}.blocksSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"relay_option_82": schema.SingleNestedBlock{
				Description: "DHCP option-82 processing.",
				Attributes:  forwardingoptionsDhcprelayBlockRelayOption82{}.attributesSchema(),
				Blocks:      forwardingoptionsDhcprelayBlockRelayOption82{}.blocksSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"server_match_address": schema.SetNestedBlock{
				Description: "For each `address`, server match processing.",
				NestedObject: schema.NestedBlockObject{
					Attributes: forwardingoptionsDhcprelayBlockServerMatchAddress{}.attributesSchema(),
				},
			},
			"server_match_duid": schema.SetNestedBlock{
				Description: "For each combination of `compare`, `value_type` and `value` arguments, match duid processing.",
				NestedObject: schema.NestedBlockObject{
					Attributes: forwardingoptionsDhcprelayBlockServerMatchDuid{}.attributesSchema(),
				},
			},
		},
	}
}

//nolint:lll
type forwardingoptionsDhcprelayData struct {
	ID                                   types.String                                          `tfsdk:"id"`
	RoutingInstance                      types.String                                          `tfsdk:"routing_instance"`
	Version                              types.String                                          `tfsdk:"version"`
	AccessProfile                        types.String                                          `tfsdk:"access_profile"`
	ActiveServerGroup                    types.String                                          `tfsdk:"active_server_group"`
	ActiveServerGroupAllowServerChange   types.Bool                                            `tfsdk:"active_server_group_allow_server_change"`
	ARPInspection                        types.Bool                                            `tfsdk:"arp_inspection"`
	AuthenticationPassword               types.String                                          `tfsdk:"authentication_password"`
	ClientResponseTTL                    types.Int64                                           `tfsdk:"client_response_ttl"`
	DuplicateClientsInSubnet             types.String                                          `tfsdk:"duplicate_clients_in_subnet"`
	DuplicateClientsIncomingInterface    types.Bool                                            `tfsdk:"duplicate_clients_incoming_interface"`
	DynamicProfile                       types.String                                          `tfsdk:"dynamic_profile"`
	DynamicProfileAggregateClients       types.Bool                                            `tfsdk:"dynamic_profile_aggregate_clients"`
	DynamicProfileAggregateClientsAction types.String                                          `tfsdk:"dynamic_profile_aggregate_clients_action"`
	DynamicProfileUsePrimary             types.String                                          `tfsdk:"dynamic_profile_use_primary"`
	ExcludeRelayAgentIdentifier          types.Bool                                            `tfsdk:"exclude_relay_agent_identifier"`
	ForwardOnly                          types.Bool                                            `tfsdk:"forward_only"`
	ForwardOnlyReplies                   types.Bool                                            `tfsdk:"forward_only_replies"`
	ForwardOnlyRoutingInstance           types.String                                          `tfsdk:"forward_only_routing_instance"`
	ForwardSnoopedClients                types.String                                          `tfsdk:"forward_snooped_clients"`
	LivenessDetectionFailureAction       types.String                                          `tfsdk:"liveness_detection_failure_action"`
	MaximumHopCount                      types.Int64                                           `tfsdk:"maximum_hop_count"`
	MinimumWaitTime                      types.Int64                                           `tfsdk:"minimum_wait_time"`
	NoSnoop                              types.Bool                                            `tfsdk:"no_snoop"`
	PersistentStorageAutomatic           types.Bool                                            `tfsdk:"persistent_storage_automatic"`
	RelayAgentOption79                   types.Bool                                            `tfsdk:"relay_agent_option_79"`
	RemoteIDMismatchDisconnect           types.Bool                                            `tfsdk:"remote_id_mismatch_disconnect"`
	RouteSuppressionAccess               types.Bool                                            `tfsdk:"route_suppression_access"`
	RouteSuppressionAccessInternal       types.Bool                                            `tfsdk:"route_suppression_access_internal"`
	RouteSuppressionDestination          types.Bool                                            `tfsdk:"route_suppression_destination"`
	ServerMatchDefaultAction             types.String                                          `tfsdk:"server_match_default_action"`
	ServerResponseTime                   types.Int64                                           `tfsdk:"server_response_time"`
	ServiceProfile                       types.String                                          `tfsdk:"service_profile"`
	ShortCycleProtectionLockoutMaxTime   types.Int64                                           `tfsdk:"short_cycle_protection_lockout_max_time"`
	ShortCycleProtectionLockoutMinTime   types.Int64                                           `tfsdk:"short_cycle_protection_lockout_min_time"`
	SourceIPChange                       types.Bool                                            `tfsdk:"source_ip_change"`
	VendorSpecificInformationHostName    types.Bool                                            `tfsdk:"vendor_specific_information_host_name"`
	VendorSpecificInformationLocation    types.Bool                                            `tfsdk:"vendor_specific_information_location"`
	ActiveLeasequery                     *forwardingoptionsDhcprelayBlockActiveLeasequery      `tfsdk:"active_leasequery"`
	AuthenticationUsernameInclude        *dhcpBlockAuthenticationUsernameInclude               `tfsdk:"authentication_username_include"`
	BulkLeasequery                       *forwardingoptionsDhcprelayBlockBulkLeasequery        `tfsdk:"bulk_leasequery"`
	LeaseTimeValidation                  *forwardingoptionsDhcprelayBlockLeaseTimeValidation   `tfsdk:"lease_time_validation"`
	Leasequery                           *forwardingoptionsDhcprelayBlockLeasequery            `tfsdk:"leasequery"`
	LivenessDetectionMethodBfd           *dhcpBlockLivenessDetectionMethodBfd                  `tfsdk:"liveness_detection_method_bfd"`
	LivenessDetectionMethodLayer2        *dhcpBlockLivenessDetectionMethodLayer2               `tfsdk:"liveness_detection_method_layer2"`
	OverridesV4                          *forwardingoptionsDhcprelayBlockOverridesV4           `tfsdk:"overrides_v4"`
	OverridesV6                          *forwardingoptionsDhcprelayBlockOverridesV6           `tfsdk:"overrides_v6"`
	RelayAgentInterfaceID                *forwardingoptionsDhcprelayBlockRelayAgentInterfaceID `tfsdk:"relay_agent_interface_id"`
	RelayAgentRemoteID                   *forwardingoptionsDhcprelayBlockRelayAgentRemoteID    `tfsdk:"relay_agent_remote_id"`
	RelayOption                          *forwardingoptionsDhcprelayBlockRelayOption           `tfsdk:"relay_option"`
	RelayOption82                        *forwardingoptionsDhcprelayBlockRelayOption82         `tfsdk:"relay_option_82"`
	ServerMatchAddress                   []forwardingoptionsDhcprelayBlockServerMatchAddress   `tfsdk:"server_match_address"`
	ServerMatchDuid                      []forwardingoptionsDhcprelayBlockServerMatchDuid      `tfsdk:"server_match_duid"`
}

//nolint:lll
type forwardingoptionsDhcprelayConfig struct {
	ID                                   types.String                                          `tfsdk:"id"`
	RoutingInstance                      types.String                                          `tfsdk:"routing_instance"`
	Version                              types.String                                          `tfsdk:"version"`
	AccessProfile                        types.String                                          `tfsdk:"access_profile"`
	ActiveServerGroup                    types.String                                          `tfsdk:"active_server_group"`
	ActiveServerGroupAllowServerChange   types.Bool                                            `tfsdk:"active_server_group_allow_server_change"`
	ARPInspection                        types.Bool                                            `tfsdk:"arp_inspection"`
	AuthenticationPassword               types.String                                          `tfsdk:"authentication_password"`
	ClientResponseTTL                    types.Int64                                           `tfsdk:"client_response_ttl"`
	DuplicateClientsInSubnet             types.String                                          `tfsdk:"duplicate_clients_in_subnet"`
	DuplicateClientsIncomingInterface    types.Bool                                            `tfsdk:"duplicate_clients_incoming_interface"`
	DynamicProfile                       types.String                                          `tfsdk:"dynamic_profile"`
	DynamicProfileAggregateClients       types.Bool                                            `tfsdk:"dynamic_profile_aggregate_clients"`
	DynamicProfileAggregateClientsAction types.String                                          `tfsdk:"dynamic_profile_aggregate_clients_action"`
	DynamicProfileUsePrimary             types.String                                          `tfsdk:"dynamic_profile_use_primary"`
	ExcludeRelayAgentIdentifier          types.Bool                                            `tfsdk:"exclude_relay_agent_identifier"`
	ForwardOnly                          types.Bool                                            `tfsdk:"forward_only"`
	ForwardOnlyReplies                   types.Bool                                            `tfsdk:"forward_only_replies"`
	ForwardOnlyRoutingInstance           types.String                                          `tfsdk:"forward_only_routing_instance"`
	ForwardSnoopedClients                types.String                                          `tfsdk:"forward_snooped_clients"`
	LivenessDetectionFailureAction       types.String                                          `tfsdk:"liveness_detection_failure_action"`
	MaximumHopCount                      types.Int64                                           `tfsdk:"maximum_hop_count"`
	MinimumWaitTime                      types.Int64                                           `tfsdk:"minimum_wait_time"`
	NoSnoop                              types.Bool                                            `tfsdk:"no_snoop"`
	PersistentStorageAutomatic           types.Bool                                            `tfsdk:"persistent_storage_automatic"`
	RelayAgentOption79                   types.Bool                                            `tfsdk:"relay_agent_option_79"`
	RemoteIDMismatchDisconnect           types.Bool                                            `tfsdk:"remote_id_mismatch_disconnect"`
	RouteSuppressionAccess               types.Bool                                            `tfsdk:"route_suppression_access"`
	RouteSuppressionAccessInternal       types.Bool                                            `tfsdk:"route_suppression_access_internal"`
	RouteSuppressionDestination          types.Bool                                            `tfsdk:"route_suppression_destination"`
	ServerMatchDefaultAction             types.String                                          `tfsdk:"server_match_default_action"`
	ServerResponseTime                   types.Int64                                           `tfsdk:"server_response_time"`
	ServiceProfile                       types.String                                          `tfsdk:"service_profile"`
	ShortCycleProtectionLockoutMaxTime   types.Int64                                           `tfsdk:"short_cycle_protection_lockout_max_time"`
	ShortCycleProtectionLockoutMinTime   types.Int64                                           `tfsdk:"short_cycle_protection_lockout_min_time"`
	SourceIPChange                       types.Bool                                            `tfsdk:"source_ip_change"`
	VendorSpecificInformationHostName    types.Bool                                            `tfsdk:"vendor_specific_information_host_name"`
	VendorSpecificInformationLocation    types.Bool                                            `tfsdk:"vendor_specific_information_location"`
	ActiveLeasequery                     *forwardingoptionsDhcprelayBlockActiveLeasequery      `tfsdk:"active_leasequery"`
	AuthenticationUsernameInclude        *dhcpBlockAuthenticationUsernameInclude               `tfsdk:"authentication_username_include"`
	BulkLeasequery                       *forwardingoptionsDhcprelayBlockBulkLeasequery        `tfsdk:"bulk_leasequery"`
	LeaseTimeValidation                  *forwardingoptionsDhcprelayBlockLeaseTimeValidation   `tfsdk:"lease_time_validation"`
	Leasequery                           *forwardingoptionsDhcprelayBlockLeasequery            `tfsdk:"leasequery"`
	LivenessDetectionMethodBfd           *dhcpBlockLivenessDetectionMethodBfd                  `tfsdk:"liveness_detection_method_bfd"`
	LivenessDetectionMethodLayer2        *dhcpBlockLivenessDetectionMethodLayer2               `tfsdk:"liveness_detection_method_layer2"`
	OverridesV4                          *forwardingoptionsDhcprelayBlockOverridesV4           `tfsdk:"overrides_v4"`
	OverridesV6                          *forwardingoptionsDhcprelayBlockOverridesV6           `tfsdk:"overrides_v6"`
	RelayAgentInterfaceID                *forwardingoptionsDhcprelayBlockRelayAgentInterfaceID `tfsdk:"relay_agent_interface_id"`
	RelayAgentRemoteID                   *forwardingoptionsDhcprelayBlockRelayAgentRemoteID    `tfsdk:"relay_agent_remote_id"`
	RelayOption                          *forwardingoptionsDhcprelayBlockRelayOptionConfig     `tfsdk:"relay_option"`
	RelayOption82                        *forwardingoptionsDhcprelayBlockRelayOption82         `tfsdk:"relay_option_82"`
	ServerMatchAddress                   types.Set                                             `tfsdk:"server_match_address"`
	ServerMatchDuid                      types.Set                                             `tfsdk:"server_match_duid"`
}

type forwardingoptionsDhcprelayBlockActiveLeasequery struct {
	IdleTimeout      types.Int64  `tfsdk:"idle_timeout"`
	PeerAddress      types.String `tfsdk:"peer_address"`
	Timeout          types.Int64  `tfsdk:"timeout"`
	TopologyDiscover types.Bool   `tfsdk:"topology_discover"`
}

type forwardingoptionsDhcprelayBlockBulkLeasequery struct {
	Attempts         types.Int64 `tfsdk:"attempts"`
	Timeout          types.Int64 `tfsdk:"timeout"`
	TriggerAutomatic types.Bool  `tfsdk:"trigger_automatic"`
}

type forwardingoptionsDhcprelayBlockLeasequery struct {
	Attempts types.Int64 `tfsdk:"attempts"`
	Timeout  types.Int64 `tfsdk:"timeout"`
}

func (rsc *forwardingoptionsDhcprelay) ValidateConfig( //nolint:gocognit,gocyclo
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config forwardingoptionsDhcprelayConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	version := config.Version.ValueString()
	switch {
	case config.Version.IsUnknown():
	case version == "v4" || config.Version.IsNull():
		if !config.DuplicateClientsIncomingInterface.IsNull() &&
			!config.DuplicateClientsIncomingInterface.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("duplicate_clients_incoming_interface"),
				tfdiag.ConflictConfigErrSummary,
				"duplicate_clients_incoming_interface cannot be configured when version = v4",
			)
		}
		if !config.ExcludeRelayAgentIdentifier.IsNull() &&
			!config.ExcludeRelayAgentIdentifier.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("exclude_relay_agent_identifier"),
				tfdiag.ConflictConfigErrSummary,
				"exclude_relay_agent_identifier cannot be configured when version = v4",
			)
		}
		if !config.RelayAgentOption79.IsNull() &&
			!config.RelayAgentOption79.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("relay_agent_option_79"),
				tfdiag.ConflictConfigErrSummary,
				"relay_agent_option_79 cannot be configured when version = v4",
			)
		}
		if !config.RouteSuppressionAccess.IsNull() &&
			!config.RouteSuppressionAccess.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("route_suppression_access"),
				tfdiag.ConflictConfigErrSummary,
				"route_suppression_access cannot be configured when version = v4",
			)
		}
		if !config.ServiceProfile.IsNull() &&
			!config.ServiceProfile.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("service_profile"),
				tfdiag.ConflictConfigErrSummary,
				"service_profile cannot be configured when version = v4",
			)
		}
		if !config.VendorSpecificInformationHostName.IsNull() &&
			!config.VendorSpecificInformationHostName.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("vendor_specific_information_host_name"),
				tfdiag.ConflictConfigErrSummary,
				"vendor_specific_information_host_name cannot be configured when version = v4",
			)
		}
		if !config.VendorSpecificInformationLocation.IsNull() &&
			!config.VendorSpecificInformationLocation.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("vendor_specific_information_location"),
				tfdiag.ConflictConfigErrSummary,
				"vendor_specific_information_location cannot be configured when version = v4",
			)
		}

		if config.AuthenticationUsernameInclude != nil {
			if !config.AuthenticationUsernameInclude.RelayAgentInterfaceID.IsNull() &&
				!config.AuthenticationUsernameInclude.RelayAgentInterfaceID.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("authentication_username_include").AtName("relay_agent_interface_id"),
					tfdiag.ConflictConfigErrSummary,
					"relay_agent_interface_id cannot be configured when version = v4"+
						" in authentication_username_include block",
				)
			}
			if !config.AuthenticationUsernameInclude.RelayAgentRemoteID.IsNull() &&
				!config.AuthenticationUsernameInclude.RelayAgentRemoteID.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("authentication_username_include").AtName("relay_agent_remote_id"),
					tfdiag.ConflictConfigErrSummary,
					"relay_agent_remote_id cannot be configured when version = v4"+
						" in authentication_username_include block",
				)
			}
			if !config.AuthenticationUsernameInclude.RelayAgentSubscriberID.IsNull() &&
				!config.AuthenticationUsernameInclude.RelayAgentSubscriberID.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("authentication_username_include").AtName("relay_agent_subscriber_id"),
					tfdiag.ConflictConfigErrSummary,
					"relay_agent_subscriber_id cannot be configured when version = v4"+
						" in authentication_username_include block",
				)
			}
		}
		if config.BulkLeasequery != nil {
			if !config.BulkLeasequery.TriggerAutomatic.IsNull() &&
				!config.BulkLeasequery.TriggerAutomatic.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("bulk_leasequery").AtName("trigger_automatic"),
					tfdiag.ConflictConfigErrSummary,
					"trigger_automatic cannot be configured when version = v4"+
						" in bulk_leasequery block",
				)
			}
		}
		if config.OverridesV6 != nil && config.OverridesV6.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("overrides_v6"),
				tfdiag.ConflictConfigErrSummary,
				"overrides_v6 cannot be configured when version = v4",
			)
		}
		if config.RelayAgentInterfaceID != nil && config.RelayAgentInterfaceID.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("relay_agent_interface_id"),
				tfdiag.ConflictConfigErrSummary,
				"relay_agent_interface_id cannot be configured when version = v4",
			)
		}
		if config.RelayAgentRemoteID != nil && config.RelayAgentRemoteID.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("relay_agent_remote_id"),
				tfdiag.ConflictConfigErrSummary,
				"relay_agent_remote_id cannot be configured when version = v4",
			)
		}
		if config.RelayOption != nil {
			if !config.RelayOption.OptionOrder.IsNull() &&
				!config.RelayOption.OptionOrder.IsUnknown() {
				var configRelayOptionOptionOrder []types.String
				asDiags := config.RelayOption.OptionOrder.ElementsAs(ctx, &configRelayOptionOptionOrder, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}
				for _, v := range configRelayOptionOptionOrder {
					if vv := v.ValueString(); vv == "15" || vv == "16" {
						resp.Diagnostics.AddAttributeError(
							path.Root("relay_option").AtName("option_order"),
							tfdiag.ConflictConfigErrSummary,
							"option_order cannot be configured with 15 or 16 when version = v4"+
								" in relay_option block",
						)
					}
				}
			}
			if !config.RelayOption.Option15.IsNull() &&
				!config.RelayOption.Option15.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("relay_option").AtName("option_15"),
					tfdiag.ConflictConfigErrSummary,
					"option_15 cannot be configured when version = v4"+
						" in relay_option block",
				)
			}
			if config.RelayOption.Option15DefaultAction != nil &&
				config.RelayOption.Option15DefaultAction.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("relay_option").AtName("option_15_default_action"),
					tfdiag.ConflictConfigErrSummary,
					"option_15_default_action cannot be configured when version = v4"+
						" in relay_option block",
				)
			}
			if !config.RelayOption.Option16.IsNull() &&
				!config.RelayOption.Option16.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("relay_option").AtName("option_16"),
					tfdiag.ConflictConfigErrSummary,
					"option_16 cannot be configured when version = v4"+
						" in relay_option block",
				)
			}
			if config.RelayOption.Option16DefaultAction != nil &&
				config.RelayOption.Option16DefaultAction.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("relay_option").AtName("option_16_default_action"),
					tfdiag.ConflictConfigErrSummary,
					"option_16_default_action cannot be configured when version = v4"+
						" in relay_option block",
				)
			}
		}
		if !config.ServerMatchDuid.IsNull() &&
			!config.ServerMatchDuid.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("server_match_duid"),
				tfdiag.ConflictConfigErrSummary,
				"server_match_duid cannot be configured when version = v4",
			)
		}
	case version == "v6":
		if !config.ActiveServerGroupAllowServerChange.IsNull() &&
			!config.ActiveServerGroupAllowServerChange.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("active_server_group_allow_server_change"),
				tfdiag.ConflictConfigErrSummary,
				"active_server_group_allow_server_change cannot be configured when version = v6",
			)
		}
		if !config.ARPInspection.IsNull() &&
			!config.ARPInspection.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("arp_inspection"),
				tfdiag.ConflictConfigErrSummary,
				"arp_inspection cannot be configured when version = v6",
			)
		}
		if !config.ClientResponseTTL.IsNull() &&
			!config.ClientResponseTTL.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("client_response_ttl"),
				tfdiag.ConflictConfigErrSummary,
				"client_response_ttl cannot be configured when version = v6",
			)
		}
		if !config.DuplicateClientsInSubnet.IsNull() &&
			!config.DuplicateClientsInSubnet.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("duplicate_clients_in_subnet"),
				tfdiag.ConflictConfigErrSummary,
				"duplicate_clients_in_subnet cannot be configured when version = v6",
			)
		}
		if !config.MaximumHopCount.IsNull() &&
			!config.MaximumHopCount.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("maximum_hop_count"),
				tfdiag.ConflictConfigErrSummary,
				"maximum_hop_count cannot be configured when version = v6",
			)
		}
		if !config.MinimumWaitTime.IsNull() &&
			!config.MinimumWaitTime.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("minimum_wait_time"),
				tfdiag.ConflictConfigErrSummary,
				"minimum_wait_time cannot be configured when version = v6",
			)
		}
		if !config.RouteSuppressionDestination.IsNull() &&
			!config.RouteSuppressionDestination.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("route_suppression_destination"),
				tfdiag.ConflictConfigErrSummary,
				"route_suppression_destination cannot be configured when version = v6",
			)
		}
		if !config.SourceIPChange.IsNull() &&
			!config.SourceIPChange.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("source_ip_change"),
				tfdiag.ConflictConfigErrSummary,
				"source_ip_change cannot be configured when version = v6",
			)
		}

		if config.AuthenticationUsernameInclude != nil {
			if !config.AuthenticationUsernameInclude.Option60.IsNull() &&
				!config.AuthenticationUsernameInclude.Option60.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("authentication_username_include").AtName("option_60"),
					tfdiag.ConflictConfigErrSummary,
					"option_60 cannot be configured when version = v6"+
						" in authentication_username_include block",
				)
			}
			if !config.AuthenticationUsernameInclude.Option82.IsNull() &&
				!config.AuthenticationUsernameInclude.Option82.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("authentication_username_include").AtName("option_82"),
					tfdiag.ConflictConfigErrSummary,
					"option_82 cannot be configured when version = v6"+
						" in authentication_username_include block",
				)
			}
			if !config.AuthenticationUsernameInclude.Option82CircuitID.IsNull() &&
				!config.AuthenticationUsernameInclude.Option82CircuitID.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("authentication_username_include").AtName("option_82_circuit_id"),
					tfdiag.ConflictConfigErrSummary,
					"option_82_circuit_id cannot be configured when version = v6"+
						" in authentication_username_include block",
				)
			}
			if !config.AuthenticationUsernameInclude.Option82RemoteID.IsNull() &&
				!config.AuthenticationUsernameInclude.Option82RemoteID.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("authentication_username_include").AtName("option_82_remote_id"),
					tfdiag.ConflictConfigErrSummary,
					"option_82_remote_id cannot be configured when version = v6"+
						" in authentication_username_include block",
				)
			}
		}
		if config.BulkLeasequery != nil &&
			!config.BulkLeasequery.Attempts.IsNull() &&
			!config.BulkLeasequery.Attempts.IsUnknown() &&
			config.BulkLeasequery.Attempts.ValueInt64() > 10 {
			resp.Diagnostics.AddAttributeError(
				path.Root("bulk_leasequery").AtName("attempts"),
				tfdiag.ConflictConfigErrSummary,
				"attempts must be in range (1..10) when version = v6"+
					" in bulk_leasequery block",
			)
		}
		if config.OverridesV4 != nil && config.OverridesV4.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("overrides_v4"),
				tfdiag.ConflictConfigErrSummary,
				"overrides_v4 cannot be configured when version = v6",
			)
		}
		if config.RelayOption != nil {
			if !config.RelayOption.OptionOrder.IsNull() &&
				!config.RelayOption.OptionOrder.IsUnknown() {
				if !config.RelayOption.OptionOrder.IsNull() &&
					!config.RelayOption.OptionOrder.IsUnknown() {
					var configRelayOptionOptionOrder []types.String
					asDiags := config.RelayOption.OptionOrder.ElementsAs(ctx, &configRelayOptionOptionOrder, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}
					for _, v := range configRelayOptionOptionOrder {
						if vv := v.ValueString(); vv == "60" || vv == "77" {
							resp.Diagnostics.AddAttributeError(
								path.Root("relay_option").AtName("option_order"),
								tfdiag.ConflictConfigErrSummary,
								"option_order cannot be configured with 60 or 77 when version = v6"+
									" in relay_option block",
							)
						}
					}
				}
			}
			if !config.RelayOption.Option60.IsNull() &&
				!config.RelayOption.Option60.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("relay_option").AtName("option_60"),
					tfdiag.ConflictConfigErrSummary,
					"option_60 cannot be configured when version = v6"+
						" in relay_option block",
				)
			}
			if config.RelayOption.Option60DefaultAction != nil &&
				config.RelayOption.Option60DefaultAction.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("relay_option").AtName("option_60_default_action"),
					tfdiag.ConflictConfigErrSummary,
					"option_60_default_action cannot be configured when version = v6"+
						" in relay_option block",
				)
			}
			if !config.RelayOption.Option77.IsNull() &&
				!config.RelayOption.Option77.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("relay_option").AtName("option_77"),
					tfdiag.ConflictConfigErrSummary,
					"option_77 cannot be configured when version = v6"+
						" in relay_option block",
				)
			}
			if config.RelayOption.Option77DefaultAction != nil &&
				config.RelayOption.Option77DefaultAction.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("relay_option").AtName("option_77_default_action"),
					tfdiag.ConflictConfigErrSummary,
					"option_77_default_action cannot be configured when version = v6"+
						" in relay_option block",
				)
			}
		}
		if config.RelayOption82 != nil && config.RelayOption82.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("relay_option_82"),
				tfdiag.ConflictConfigErrSummary,
				"relay_option_82 cannot be configured when version = v6",
			)
		}
	}

	if !config.DynamicProfileAggregateClients.IsNull() &&
		!config.DynamicProfileAggregateClients.IsUnknown() &&
		!config.DynamicProfileUsePrimary.IsNull() &&
		!config.DynamicProfileUsePrimary.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dynamic_profile_aggregate_clients"),
			tfdiag.ConflictConfigErrSummary,
			"dynamic_profile_aggregate_clients and dynamic_profile_use_primary cannot be configured together",
		)
	}
	if !config.DynamicProfileAggregateClients.IsNull() &&
		!config.DynamicProfileAggregateClients.IsUnknown() &&
		config.DynamicProfile.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dynamic_profile_aggregate_clients"),
			tfdiag.MissingConfigErrSummary,
			"dynamic_profile must be specified with dynamic_profile_aggregate_clients",
		)
	}
	if !config.DynamicProfileAggregateClientsAction.IsNull() &&
		!config.DynamicProfileAggregateClientsAction.IsUnknown() &&
		config.DynamicProfileAggregateClients.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dynamic_profile_aggregate_clients_action"),
			tfdiag.MissingConfigErrSummary,
			"dynamic_profile_aggregate_clients must be specified with dynamic_profile_aggregate_clients_action",
		)
	}
	if !config.DynamicProfileUsePrimary.IsNull() &&
		!config.DynamicProfileUsePrimary.IsUnknown() &&
		config.DynamicProfile.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dynamic_profile_use_primary"),
			tfdiag.MissingConfigErrSummary,
			"dynamic_profile must be specified with dynamic_profile_use_primary",
		)
	}
	if !config.ForwardOnlyRoutingInstance.IsNull() &&
		!config.ForwardOnlyRoutingInstance.IsUnknown() &&
		config.ForwardOnly.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("forward_only_routing_instance"),
			tfdiag.MissingConfigErrSummary,
			"forward_only must be specified with forward_only_routing_instance",
		)
	}
	if !config.ShortCycleProtectionLockoutMaxTime.IsNull() &&
		!config.ShortCycleProtectionLockoutMaxTime.IsUnknown() &&
		config.ShortCycleProtectionLockoutMinTime.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("short_cycle_protection_lockout_max_time"),
			tfdiag.MissingConfigErrSummary,
			"short_cycle_protection_lockout_min_time must be specified with short_cycle_protection_lockout_max_time",
		)
	}
	if !config.ShortCycleProtectionLockoutMinTime.IsNull() &&
		!config.ShortCycleProtectionLockoutMinTime.IsUnknown() &&
		config.ShortCycleProtectionLockoutMaxTime.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("short_cycle_protection_lockout_min_time"),
			tfdiag.MissingConfigErrSummary,
			"short_cycle_protection_lockout_max_time must be specified with short_cycle_protection_lockout_min_time",
		)
	}

	if config.AuthenticationUsernameInclude != nil {
		if !config.AuthenticationUsernameInclude.Option82CircuitID.IsNull() &&
			!config.AuthenticationUsernameInclude.Option82CircuitID.IsUnknown() &&
			config.AuthenticationUsernameInclude.Option82.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_username_include").AtName("option_82_circuit_id"),
				tfdiag.MissingConfigErrSummary,
				"option_82 must be specified with option_82_circuit_id"+
					" in authentication_username_include block",
			)
		}
		if !config.AuthenticationUsernameInclude.Option82RemoteID.IsNull() &&
			!config.AuthenticationUsernameInclude.Option82RemoteID.IsUnknown() &&
			config.AuthenticationUsernameInclude.Option82.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_username_include").AtName("option_82_remote_id"),
				tfdiag.MissingConfigErrSummary,
				"option_82 must be specified with option_82_remote_id"+
					" in authentication_username_include block",
			)
		}
	}
	if config.LivenessDetectionMethodBfd != nil {
		if config.LivenessDetectionMethodBfd.hasKnownValue() &&
			config.LivenessDetectionMethodLayer2 != nil && config.LivenessDetectionMethodLayer2.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("liveness_detection_method_bfd").AtName("*"),
				tfdiag.ConflictConfigErrSummary,
				"liveness_detection_method_bfd and liveness_detection_method_layer2 cannot be configured together",
			)
		}
		if config.LivenessDetectionMethodBfd.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("liveness_detection_method_bfd").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"liveness_detection_method_bfd block is empty",
			)
		}
	}
	if config.LivenessDetectionMethodLayer2 != nil {
		if config.LivenessDetectionMethodLayer2.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("liveness_detection_method_layer2").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"liveness_detection_method_layer2 block is empty",
			)
		}
	}
	if config.OverridesV4 != nil {
		if config.OverridesV4.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("overrides_v4").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"overrides_v4 block is empty",
			)
		}
	}
	if config.OverridesV6 != nil {
		if config.OverridesV6.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("overrides_v6").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"overrides_v6 block is empty",
			)
		}
	}
	if config.RelayOption != nil {
		if config.RelayOption.Option15DefaultAction != nil &&
			config.RelayOption.Option15DefaultAction.Action.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("relay_option").AtName("option_15_default_action").AtName("action"),
				tfdiag.MissingConfigErrSummary,
				"action must be specified"+
					" in option_15_default_action block in relay_option block",
			)
		}
		if config.RelayOption.Option16DefaultAction != nil &&
			config.RelayOption.Option16DefaultAction.Action.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("relay_option").AtName("option_16_default_action").AtName("action"),
				tfdiag.MissingConfigErrSummary,
				"action must be specified"+
					" in option_16_default_action block in relay_option block",
			)
		}
		if config.RelayOption.Option60DefaultAction != nil &&
			config.RelayOption.Option60DefaultAction.Action.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("relay_option").AtName("option_60_default_action").AtName("action"),
				tfdiag.MissingConfigErrSummary,
				"action must be specified"+
					" in option_60_default_action block in relay_option block",
			)
		}
		if config.RelayOption.Option77DefaultAction != nil &&
			config.RelayOption.Option77DefaultAction.Action.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("relay_option").AtName("option_77_default_action").AtName("action"),
				tfdiag.MissingConfigErrSummary,
				"action must be specified"+
					" in option_77_default_action block in relay_option block",
			)
		}
	}
	if !config.ServerMatchAddress.IsNull() &&
		!config.ServerMatchAddress.IsUnknown() {
		var configServerMatchAddress []forwardingoptionsDhcprelayBlockServerMatchAddress
		asDiags := config.ServerMatchAddress.ElementsAs(ctx, &configServerMatchAddress, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		serverMatchAddressAddress := make(map[string]struct{})
		for _, block := range configServerMatchAddress {
			if block.Address.IsUnknown() {
				continue
			}
			address := block.Address.ValueString()
			if _, ok := serverMatchAddressAddress[address]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("server_match_address"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple server_match_address blocks with the same address %q", address),
				)
			}
			serverMatchAddressAddress[address] = struct{}{}
		}
	}
	if !config.ServerMatchDuid.IsNull() &&
		!config.ServerMatchDuid.IsUnknown() {
		var configServerMatchDuid []forwardingoptionsDhcprelayBlockServerMatchDuid
		asDiags := config.ServerMatchDuid.ElementsAs(ctx, &configServerMatchDuid, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		serverMatchDuidBlock := make(map[string]struct{})
		for _, block := range configServerMatchDuid {
			if block.Compare.IsUnknown() {
				continue
			}
			if block.ValueType.IsUnknown() {
				continue
			}
			if block.Value.IsUnknown() {
				continue
			}
			blockString := block.Compare.ValueString() +
				junos.IDSeparator + block.ValueType.ValueString() +
				junos.IDSeparator + block.Value.ValueString()
			if _, ok := serverMatchDuidBlock[blockString]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("server_match_duid"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple server_match_duid blocks with the same compare %q, value_type %q, value %q",
						block.Compare.ValueString(), block.ValueType.ValueString(), block.Value.ValueString()),
				)
			}
			serverMatchDuidBlock[blockString] = struct{}{}
		}
	}
}

func (rsc *forwardingoptionsDhcprelay) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan forwardingoptionsDhcprelayData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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

			return true
		},
		nil,
		&plan,
		resp,
	)
}

func (rsc *forwardingoptionsDhcprelay) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data forwardingoptionsDhcprelayData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	if v := state.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, v, junSess)
		if err != nil {
			junos.MutexUnlock()
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !instanceExists {
			junos.MutexUnlock()
			resp.State.RemoveResource(ctx)

			return
		}
	}

	err = data.read(ctx, state.RoutingInstance.ValueString(), state.Version.ValueString(), junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}

	if data.nullID() {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *forwardingoptionsDhcprelay) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state forwardingoptionsDhcprelayData
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

func (rsc *forwardingoptionsDhcprelay) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state forwardingoptionsDhcprelayData
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

func (rsc *forwardingoptionsDhcprelay) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	idList := strings.Split(req.ID, junos.IDSeparator)
	if len(idList) < 2 {
		resp.Diagnostics.AddError(
			"Bad ID Format",
			fmt.Sprintf("missing element(s) in id with separator %q", junos.IDSeparator),
		)

		return
	}
	if idList[0] != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, idList[0], junSess)
		if err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !instanceExists {
			resp.Diagnostics.AddError(
				tfdiag.NotFoundErrSummary,
				fmt.Sprintf("routing instance %q doesn't exist", idList[0]),
			)

			return
		}
	}
	if idList[1] != "v4" && idList[1] != "v6" {
		resp.Diagnostics.AddError(
			"Bad ID Format",
			fmt.Sprintf("%q is not a valid version", idList[1]),
		)

		return
	}

	var data forwardingoptionsDhcprelayData
	if err := data.read(ctx, idList[0], idList[1], junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if data.nullID() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be <routing_instance>_-_<version>)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rscData *forwardingoptionsDhcprelayData) fillID() {
	routingInstance := rscData.RoutingInstance.ValueString()
	version := rscData.Version.ValueString()
	if routingInstance == "" {
		routingInstance = junos.DefaultW
	}
	if version == "" {
		version = "v4"
	}

	rscData.ID = types.StringValue(routingInstance + junos.IDSeparator + version)
}

func (rscData *forwardingoptionsDhcprelayData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *forwardingoptionsDhcprelayData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "forwarding-options dhcp-relay "
	version := rscData.Version.ValueString()
	if version == "v6" {
		setPrefix += "dhcpv6 "
	} else if version != "v4" {
		version = "v4"
	}

	if v := rscData.AccessProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"access-profile \""+v+"\"")
	}
	if v := rscData.ActiveServerGroup.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"active-server-group \""+v+"\"")
	}
	if rscData.ActiveServerGroupAllowServerChange.ValueBool() {
		if version == "v6" {
			return path.Root("active_server_group_allow_server_change"),
				errors.New("active_server_group_allow_server_change cannot be configured when version = v6")
		}

		configSet = append(configSet, setPrefix+"active-server-group allow-server-change")
	}
	if rscData.ARPInspection.ValueBool() {
		configSet = append(configSet, setPrefix+"arp-inspection")
	}
	if v := rscData.AuthenticationPassword.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication password \""+v+"\"")
	}
	if !rscData.ClientResponseTTL.IsNull() {
		if version == "v6" {
			return path.Root("client_response_ttl"),
				errors.New("client_response_ttl cannot be configured when version = v6")
		}

		configSet = append(configSet, setPrefix+"client-response-ttl "+
			utils.ConvI64toa(rscData.ClientResponseTTL.ValueInt64()))
	}
	if v := rscData.DuplicateClientsInSubnet.ValueString(); v != "" {
		if version == "v6" {
			return path.Root("duplicate_clients_in_subnet"),
				errors.New("duplicate_clients_in_subnet cannot be configured when version = v6")
		}

		configSet = append(configSet, setPrefix+"duplicate-clients-in-subnet "+v)
	}
	if rscData.DuplicateClientsIncomingInterface.ValueBool() {
		if version == "v4" {
			return path.Root("duplicate_clients_incoming_interface"),
				errors.New("duplicate_clients_incoming_interface cannot be configured when version = v4")
		}

		configSet = append(configSet, setPrefix+"duplicate-clients incoming-interface")
	}
	if v := rscData.DynamicProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"dynamic-profile \""+v+"\"")
		if rscData.DynamicProfileAggregateClients.ValueBool() {
			configSet = append(configSet, setPrefix+"dynamic-profile aggregate-clients")
			if vv := rscData.DynamicProfileAggregateClientsAction.ValueString(); vv != "" {
				configSet = append(configSet, setPrefix+"dynamic-profile aggregate-clients "+vv)
			}
		} else if rscData.DynamicProfileAggregateClientsAction.ValueString() != "" {
			return path.Root("dynamic_profile_aggregate_clients_action"),
				errors.New("dynamic_profile_aggregate_clients must be specified with " +
					"dynamic_profile_aggregate_clients_action")
		}
		if vv := rscData.DynamicProfileUsePrimary.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"dynamic-profile use-primary \""+vv+"\"")
		}
	} else {
		if rscData.DynamicProfileAggregateClients.ValueBool() {
			return path.Root("dynamic_profile_aggregate_clients"),
				errors.New("dynamic_profile must be specified with " +
					"dynamic_profile_aggregate_clients")
		}
		if rscData.DynamicProfileAggregateClientsAction.ValueString() != "" {
			return path.Root("dynamic_profile_aggregate_clients_action"),
				errors.New("dynamic_profile must be specified with " +
					"dynamic_profile_aggregate_clients_action")
		}
		if rscData.DynamicProfileUsePrimary.ValueString() != "" {
			return path.Root("dynamic_profile_use_primary"),
				errors.New("dynamic_profile must be specified with " +
					"dynamic_profile_use_primary")
		}
	}
	if rscData.ExcludeRelayAgentIdentifier.ValueBool() {
		if version == "v4" {
			return path.Root("exclude_relay_agent_identifier"),
				errors.New("exclude_relay_agent_identifier cannot be configured when version = v4")
		}

		configSet = append(configSet, setPrefix+"exclude-relay-agent-identifier")
	}
	if rscData.ForwardOnly.ValueBool() {
		configSet = append(configSet, setPrefix+"forward-only")
		if v := rscData.ForwardOnlyRoutingInstance.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"forward-only routing-instance "+v)
		}
	} else if rscData.ForwardOnlyRoutingInstance.ValueString() != "" {
		return path.Root("forward_only_routing_instance"),
			errors.New("forward_only must be specified with forward_only_routing_instance")
	}
	if rscData.ForwardOnlyReplies.ValueBool() {
		configSet = append(configSet, setPrefix+"forward-only-replies")
	}
	if v := rscData.ForwardSnoopedClients.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"forward-snooped-clients "+v)
	}
	if v := rscData.LivenessDetectionFailureAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"liveness-detection failure-action "+v)
	}
	if !rscData.MaximumHopCount.IsNull() {
		if version == "v6" {
			return path.Root("maximum_hop_count"),
				errors.New("maximum_hop_count cannot be configured when version = v6")
		}

		configSet = append(configSet, setPrefix+"maximum-hop-count "+
			utils.ConvI64toa(rscData.MaximumHopCount.ValueInt64()))
	}
	if !rscData.MinimumWaitTime.IsNull() {
		if version == "v6" {
			return path.Root("minimum_wait_time"),
				errors.New("minimum_wait_time cannot be configured when version = v6")
		}

		configSet = append(configSet, setPrefix+"minimum-wait-time "+
			utils.ConvI64toa(rscData.MinimumWaitTime.ValueInt64()))
	}
	if rscData.NoSnoop.ValueBool() {
		configSet = append(configSet, setPrefix+"no-snoop")
	}
	if rscData.PersistentStorageAutomatic.ValueBool() {
		configSet = append(configSet, setPrefix+"persistent-storage automatic")
	}
	if rscData.RelayAgentOption79.ValueBool() {
		if version == "v4" {
			return path.Root("relay_agent_option_79"),
				errors.New("relay_agent_option_79 cannot be configured when version = v4")
		}

		configSet = append(configSet, setPrefix+"relay-agent-option-79")
	}
	if rscData.RemoteIDMismatchDisconnect.ValueBool() {
		configSet = append(configSet, setPrefix+"remote-id-mismatch disconnect")
	}
	if rscData.RouteSuppressionAccess.ValueBool() {
		if version == "v4" {
			return path.Root("route_suppression_access"),
				errors.New("route_suppression_access cannot be configured when version = v4")
		}

		configSet = append(configSet, setPrefix+"route-suppression access")
	}
	if rscData.RouteSuppressionAccessInternal.ValueBool() {
		configSet = append(configSet, setPrefix+"route-suppression access-internal")
	}
	if rscData.RouteSuppressionDestination.ValueBool() {
		if version == "v6" {
			return path.Root("route_suppression_destination"),
				errors.New("route_suppression_destination cannot be configured when version = v6")
		}

		configSet = append(configSet, setPrefix+"route-suppression destination")
	}
	if v := rscData.ServerMatchDefaultAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"server-match default-action "+v)
	}
	if !rscData.ServerResponseTime.IsNull() {
		configSet = append(configSet, setPrefix+"server-response-time "+
			utils.ConvI64toa(rscData.ServerResponseTime.ValueInt64()))
	}
	if v := rscData.ServiceProfile.ValueString(); v != "" {
		if version == "v4" {
			return path.Root("service_profile"),
				errors.New("service_profile cannot be configured when version = v4")
		}

		configSet = append(configSet, setPrefix+"service-profile \""+v+"\"")
	}
	if !rscData.ShortCycleProtectionLockoutMaxTime.IsNull() {
		configSet = append(configSet, setPrefix+"short-cycle-protection lockout-max-time "+
			utils.ConvI64toa(rscData.ShortCycleProtectionLockoutMaxTime.ValueInt64()))
	}
	if !rscData.ShortCycleProtectionLockoutMinTime.IsNull() {
		configSet = append(configSet, setPrefix+"short-cycle-protection lockout-min-time "+
			utils.ConvI64toa(rscData.ShortCycleProtectionLockoutMinTime.ValueInt64()))
	}
	if rscData.SourceIPChange.ValueBool() {
		if version == "v6" {
			return path.Root("source_ip_change"),
				errors.New("source_ip_change cannot be configured when version = v6")
		}

		configSet = append(configSet, setPrefix+"source-ip-change")
	}
	if rscData.VendorSpecificInformationHostName.ValueBool() {
		if version == "v4" {
			return path.Root("vendor_specific_information_host_name"),
				errors.New("vendor_specific_information_host_name cannot be configured when version = v4")
		}

		configSet = append(configSet, setPrefix+"vendor-specific-information host-name")
	}
	if rscData.VendorSpecificInformationLocation.ValueBool() {
		if version == "v4" {
			return path.Root("vendor_specific_information_location"),
				errors.New("vendor_specific_information_location cannot be configured when version = v4")
		}

		configSet = append(configSet, setPrefix+"vendor-specific-information location")
	}

	if rscData.ActiveLeasequery != nil {
		configSet = append(configSet, rscData.ActiveLeasequery.configSet(setPrefix)...)
	}
	if rscData.AuthenticationUsernameInclude != nil {
		if rscData.AuthenticationUsernameInclude.isEmpty() {
			return path.Root("authentication_username_include"),
				errors.New("authentication_username_include block is empty")
		}

		blockSet, pathErr, err := rscData.AuthenticationUsernameInclude.configSet(setPrefix, version)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.BulkLeasequery != nil {
		blockSet, pathErr, err := rscData.BulkLeasequery.configSet(setPrefix, version)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.LeaseTimeValidation != nil {
		configSet = append(configSet, rscData.LeaseTimeValidation.configSet(setPrefix)...)
	}
	if rscData.Leasequery != nil {
		configSet = append(configSet, rscData.Leasequery.configSet(setPrefix)...)
	}
	if rscData.LivenessDetectionMethodBfd != nil {
		if rscData.LivenessDetectionMethodBfd.isEmpty() {
			return path.Root("liveness_detection_method_bfd"),
				errors.New("liveness_detection_method_bfd block is empty")
		}

		configSet = append(configSet, rscData.LivenessDetectionMethodBfd.configSet(setPrefix)...)
	}
	if rscData.LivenessDetectionMethodLayer2 != nil {
		if rscData.LivenessDetectionMethodLayer2.isEmpty() {
			return path.Root("liveness_detection_method_layer2"),
				errors.New("liveness_detection_method_layer2 block is empty")
		}

		configSet = append(configSet, rscData.LivenessDetectionMethodLayer2.configSet(setPrefix)...)
	}
	if rscData.OverridesV4 != nil {
		if rscData.OverridesV4.isEmpty() {
			return path.Root("overrides_v4"),
				errors.New("overrides_v4 block is empty")
		}
		if version == "v6" {
			return path.Root("overrides_v4"),
				errors.New("overrides_v4 cannot be configured when version = v6")
		}

		blockSet, pathErr, err := rscData.OverridesV4.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.OverridesV6 != nil {
		if rscData.OverridesV6.isEmpty() {
			return path.Root("overrides_v6"),
				errors.New("overrides_v6 block is empty")
		}
		if version == "v4" {
			return path.Root("overrides_v6"),
				errors.New("overrides_v6 cannot be configured when version = v4")
		}

		blockSet, pathErr, err := rscData.OverridesV6.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.RelayAgentInterfaceID != nil {
		if version == "v4" {
			return path.Root("relay_agent_interface_id"),
				errors.New("relay_agent_interface_id cannot be configured when version = v4")
		}

		blockSet, pathErr, err := rscData.RelayAgentInterfaceID.configSet(setPrefix, path.Root("relay_agent_interface_id"))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.RelayAgentRemoteID != nil {
		if version == "v4" {
			return path.Root("relay_agent_remote_id"),
				errors.New("relay_agent_remote_id cannot be configured when version = v4")
		}

		blockSet, pathErr, err := rscData.RelayAgentRemoteID.configSet(setPrefix, path.Root("relay_agent_remote_id"))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.RelayOption != nil {
		blockSet, pathErr, err := rscData.RelayOption.configSet(setPrefix, version)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.RelayOption82 != nil {
		if version == "v6" {
			return path.Root("relay_option_82"),
				errors.New("relay_option_82 cannot be configured when version = v6")
		}

		configSet = append(configSet, rscData.RelayOption82.configSet(setPrefix)...)
	}
	serverMatchAddressAddress := make(map[string]struct{})
	for _, block := range rscData.ServerMatchAddress {
		address := block.Address.ValueString()
		if _, ok := serverMatchAddressAddress[address]; ok {
			return path.Root("server_match_address"),
				fmt.Errorf("multiple server_match_address blocks with the same address %q", address)
		}
		serverMatchAddressAddress[address] = struct{}{}

		configSet = append(configSet, block.configSet(setPrefix)...)
	}
	serverMatchDuidBlock := make(map[string]struct{})
	for _, block := range rscData.ServerMatchDuid {
		if version == "v4" {
			return path.Root("server_match_duid"),
				errors.New("server_match_duid cannot be configured when version = v4")
		}

		blockString := block.Compare.ValueString() +
			junos.IDSeparator + block.ValueType.ValueString() +
			junos.IDSeparator + block.Value.ValueString()
		if _, ok := serverMatchDuidBlock[blockString]; ok {
			return path.Root("server_match_duid"),

				fmt.Errorf("multiple server_match_duid blocks with the same compare %q, value_type %q, value %q",
					block.Compare.ValueString(), block.ValueType.ValueString(), block.Value.ValueString())
		}
		serverMatchDuidBlock[blockString] = struct{}{}

		configSet = append(configSet, block.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *forwardingoptionsDhcprelayBlockActiveLeasequery) configSet(setPrefix string) []string {
	setPrefix += "active-leasequery "

	configSet := []string{
		setPrefix,
	}

	if !block.IdleTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"idle-timeout "+
			utils.ConvI64toa(block.IdleTimeout.ValueInt64()))
	}
	if v := block.PeerAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"peer-address "+v)
	}
	if !block.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(block.Timeout.ValueInt64()))
	}
	if block.TopologyDiscover.ValueBool() {
		configSet = append(configSet, setPrefix+"topology-discover")
	}

	return configSet
}

func (block *forwardingoptionsDhcprelayBlockBulkLeasequery) configSet(
	setPrefix,
	version string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "bulk-leasequery "

	configSet := []string{
		setPrefix,
	}

	if !block.Attempts.IsNull() {
		if version == "v6" &&
			block.Attempts.ValueInt64() > 10 {
			return configSet,
				path.Root("bulk_leasequery").AtName("attempts"),
				errors.New("attempts must be in range (1..10) when version = v6" +
					" in bulk_leasequery block")
		}

		configSet = append(configSet, setPrefix+"attempts "+
			utils.ConvI64toa(block.Attempts.ValueInt64()))
	}
	if !block.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(block.Timeout.ValueInt64()))
	}
	if block.TriggerAutomatic.ValueBool() {
		if version == "v4" {
			return configSet,
				path.Root("bulk_leasequery").AtName("trigger_automatic"),
				errors.New("trigger_automatic cannot be configured when version = v4" +
					" in bulk_leasequery block")
		}

		configSet = append(configSet, setPrefix+"trigger automatic")
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsDhcprelayBlockLeasequery) configSet(setPrefix string) []string {
	setPrefix += "leasequery "

	configSet := []string{
		setPrefix,
	}

	if !block.Attempts.IsNull() {
		configSet = append(configSet, setPrefix+"attempts "+
			utils.ConvI64toa(block.Attempts.ValueInt64()))
	}
	if !block.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(block.Timeout.ValueInt64()))
	}

	return configSet
}

func (rscData *forwardingoptionsDhcprelayData) read( //nolint:gocognit
	_ context.Context, routingInstance, version string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	var showSuffix string
	if version == "v6" {
		showSuffix = "dhcpv6"
	}
	showConfig, err := junSess.Command(showPrefix +
		"forwarding-options dhcp-relay " + showSuffix + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if routingInstance == "" {
		rscData.RoutingInstance = types.StringValue(junos.DefaultW)
	} else {
		rscData.RoutingInstance = types.StringValue(routingInstance)
	}
	if version == "v6" {
		rscData.Version = types.StringValue(version)
	} else {
		rscData.Version = types.StringValue("v4")
	}
	rscData.fillID()
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			switch {
			case balt.CutPrefixInString(&itemTrim, "access-profile "):
				rscData.AccessProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "active-server-group allow-server-change":
				rscData.ActiveServerGroupAllowServerChange = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "active-server-group "):
				rscData.ActiveServerGroup = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "arp-inspection":
				rscData.ARPInspection = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "authentication password "):
				rscData.AuthenticationPassword = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "client-response-ttl "):
				rscData.ClientResponseTTL, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "duplicate-clients-in-subnet "):
				rscData.DuplicateClientsInSubnet = types.StringValue(itemTrim)
			case itemTrim == "duplicate-clients incoming-interface":
				rscData.DuplicateClientsIncomingInterface = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "dynamic-profile aggregate-clients"):
				rscData.DynamicProfileAggregateClients = types.BoolValue(true)
				if balt.CutPrefixInString(&itemTrim, " ") {
					rscData.DynamicProfileAggregateClientsAction = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "dynamic-profile use-primary "):
				rscData.DynamicProfileUsePrimary = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "dynamic-profile "):
				rscData.DynamicProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "exclude-relay-agent-identifier":
				rscData.ExcludeRelayAgentIdentifier = types.BoolValue(true)
			case itemTrim == "forward-only-replies":
				rscData.ForwardOnlyReplies = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "forward-only"):
				rscData.ForwardOnly = types.BoolValue(true)
				if balt.CutPrefixInString(&itemTrim, " routing-instance ") {
					rscData.ForwardOnlyRoutingInstance = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			case balt.CutPrefixInString(&itemTrim, "forward-snooped-clients "):
				rscData.ForwardSnoopedClients = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "liveness-detection failure-action "):
				rscData.LivenessDetectionFailureAction = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "maximum-hop-count "):
				rscData.MaximumHopCount, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "minimum-wait-time "):
				rscData.MinimumWaitTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "no-snoop":
				rscData.NoSnoop = types.BoolValue(true)
			case itemTrim == "persistent-storage automatic":
				rscData.PersistentStorageAutomatic = types.BoolValue(true)
			case itemTrim == "relay-agent-option-79":
				rscData.RelayAgentOption79 = types.BoolValue(true)
			case itemTrim == "remote-id-mismatch disconnect":
				rscData.RemoteIDMismatchDisconnect = types.BoolValue(true)
			case itemTrim == "route-suppression access":
				rscData.RouteSuppressionAccess = types.BoolValue(true)
			case itemTrim == "route-suppression access-internal":
				rscData.RouteSuppressionAccessInternal = types.BoolValue(true)
			case itemTrim == "route-suppression destination":
				rscData.RouteSuppressionDestination = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "server-match default-action "):
				rscData.ServerMatchDefaultAction = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "server-response-time "):
				rscData.ServerResponseTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "service-profile "):
				rscData.ServiceProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "short-cycle-protection lockout-max-time "):
				rscData.ShortCycleProtectionLockoutMaxTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "short-cycle-protection lockout-min-time "):
				rscData.ShortCycleProtectionLockoutMinTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "source-ip-change":
				rscData.SourceIPChange = types.BoolValue(true)
			case itemTrim == "vendor-specific-information host-name":
				rscData.VendorSpecificInformationHostName = types.BoolValue(true)
			case itemTrim == "vendor-specific-information location":
				rscData.VendorSpecificInformationLocation = types.BoolValue(true)

			case balt.CutPrefixInString(&itemTrim, "active-leasequery"):
				if rscData.ActiveLeasequery == nil {
					rscData.ActiveLeasequery = &forwardingoptionsDhcprelayBlockActiveLeasequery{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.ActiveLeasequery.read(itemTrim); err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "authentication username-include "):
				if rscData.AuthenticationUsernameInclude == nil {
					rscData.AuthenticationUsernameInclude = &dhcpBlockAuthenticationUsernameInclude{}
				}

				rscData.AuthenticationUsernameInclude.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "bulk-leasequery"):
				if rscData.BulkLeasequery == nil {
					rscData.BulkLeasequery = &forwardingoptionsDhcprelayBlockBulkLeasequery{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.BulkLeasequery.read(itemTrim); err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "lease-time-validation"):
				if rscData.LeaseTimeValidation == nil {
					rscData.LeaseTimeValidation = &forwardingoptionsDhcprelayBlockLeaseTimeValidation{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.LeaseTimeValidation.read(itemTrim); err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "leasequery"):
				if rscData.Leasequery == nil {
					rscData.Leasequery = &forwardingoptionsDhcprelayBlockLeasequery{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.Leasequery.read(itemTrim); err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "liveness-detection method bfd "):
				if rscData.LivenessDetectionMethodBfd == nil {
					rscData.LivenessDetectionMethodBfd = &dhcpBlockLivenessDetectionMethodBfd{}
				}

				if err := rscData.LivenessDetectionMethodBfd.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "liveness-detection method layer2-liveness-detection "):
				if rscData.LivenessDetectionMethodLayer2 == nil {
					rscData.LivenessDetectionMethodLayer2 = &dhcpBlockLivenessDetectionMethodLayer2{}
				}

				if err := rscData.LivenessDetectionMethodLayer2.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "overrides "):
				switch rscData.Version.ValueString() {
				case "v4":
					if rscData.OverridesV4 == nil {
						rscData.OverridesV4 = &forwardingoptionsDhcprelayBlockOverridesV4{}
					}

					err = rscData.OverridesV4.read(itemTrim)
				case "v6":
					if rscData.OverridesV6 == nil {
						rscData.OverridesV6 = &forwardingoptionsDhcprelayBlockOverridesV6{}
					}

					err = rscData.OverridesV6.read(itemTrim)
				}
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "relay-agent-interface-id"):
				if rscData.RelayAgentInterfaceID == nil {
					rscData.RelayAgentInterfaceID = &forwardingoptionsDhcprelayBlockRelayAgentInterfaceID{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					rscData.RelayAgentInterfaceID.read(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "relay-agent-remote-id"):
				if rscData.RelayAgentRemoteID == nil {
					rscData.RelayAgentRemoteID = &forwardingoptionsDhcprelayBlockRelayAgentRemoteID{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					rscData.RelayAgentRemoteID.read(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "relay-option-82"):
				if rscData.RelayOption82 == nil {
					rscData.RelayOption82 = &forwardingoptionsDhcprelayBlockRelayOption82{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					rscData.RelayOption82.read(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "relay-option"):
				if rscData.RelayOption == nil {
					rscData.RelayOption = &forwardingoptionsDhcprelayBlockRelayOption{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.RelayOption.read(itemTrim); err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "server-match address "):
				var block forwardingoptionsDhcprelayBlockServerMatchAddress
				if err := block.read(itemTrim); err != nil {
					return err
				}
				rscData.ServerMatchAddress = append(rscData.ServerMatchAddress, block)
			case balt.CutPrefixInString(&itemTrim, "server-match duid "):
				var block forwardingoptionsDhcprelayBlockServerMatchDuid
				if err := block.read(itemTrim); err != nil {
					return err
				}
				rscData.ServerMatchDuid = append(rscData.ServerMatchDuid, block)
			}
		}
	}

	return nil
}

func (block *forwardingoptionsDhcprelayBlockActiveLeasequery) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "idle-timeout "):
		block.IdleTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "peer-address "):
		block.PeerAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "timeout "):
		block.Timeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "topology-discover":
		block.TopologyDiscover = types.BoolValue(true)
	}

	return err
}

func (block *forwardingoptionsDhcprelayBlockBulkLeasequery) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "attempts "):
		block.Attempts, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "timeout "):
		block.Timeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "trigger automatic":
		block.TriggerAutomatic = types.BoolValue(true)
	}

	return err
}

func (block *forwardingoptionsDhcprelayBlockLeasequery) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "attempts "):
		block.Attempts, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "timeout "):
		block.Timeout, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (rscData *forwardingoptionsDhcprelayData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "forwarding-options dhcp-relay "
	if rscData.Version.ValueString() == "v6" {
		delPrefix += "dhcpv6 "
	}

	configSet := []string{
		delPrefix + "access-profile",
		delPrefix + "active-leasequery",
		delPrefix + "active-server-group",
		delPrefix + "arp-inspection",
		delPrefix + "authentication",
		delPrefix + "bulk-leasequery",
		delPrefix + "client-response-ttl",
		delPrefix + "duplicate-clients-in-subnet",
		delPrefix + "duplicate-clients",
		delPrefix + "dynamic-profile",
		delPrefix + "exclude-relay-agent-identifier",
		delPrefix + "forward-only",
		delPrefix + "forward-only-replies",
		delPrefix + "forward-snooped-clients",
		delPrefix + "lease-time-validation",
		delPrefix + "leasequery",
		delPrefix + "liveness-detection",
		delPrefix + "maximum-hop-count",
		delPrefix + "minimum-wait-time",
		delPrefix + "no-snoop",
		delPrefix + "overrides",
		delPrefix + "persistent-storage",
		delPrefix + "relay-agent-interface-id",
		delPrefix + "relay-agent-option-79",
		delPrefix + "relay-agent-remote-id",
		delPrefix + "relay-option",
		delPrefix + "relay-option-82",
		delPrefix + "remote-id-mismatch",
		delPrefix + "route-suppression",
		delPrefix + "server-match",
		delPrefix + "server-response-time",
		delPrefix + "service-profile",
		delPrefix + "short-cycle-protection",
		delPrefix + "source-ip-change",
		delPrefix + "vendor-specific-information",
	}

	return junSess.ConfigSet(configSet)
}
