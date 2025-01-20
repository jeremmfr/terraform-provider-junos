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
	_ resource.Resource                   = &systemServicesDhcpLocalserverGroup{}
	_ resource.ResourceWithConfigure      = &systemServicesDhcpLocalserverGroup{}
	_ resource.ResourceWithValidateConfig = &systemServicesDhcpLocalserverGroup{}
	_ resource.ResourceWithImportState    = &systemServicesDhcpLocalserverGroup{}
	_ resource.ResourceWithUpgradeState   = &systemServicesDhcpLocalserverGroup{}
)

type systemServicesDhcpLocalserverGroup struct {
	client *junos.Client
}

func newSystemServicesDhcpLocalserverGroupResource() resource.Resource {
	return &systemServicesDhcpLocalserverGroup{}
}

func (rsc *systemServicesDhcpLocalserverGroup) typeName() string {
	return providerName + "_system_services_dhcp_localserver_group"
}

func (rsc *systemServicesDhcpLocalserverGroup) junosName() string {
	return "system services dhcp-local-server group"
}

func (rsc *systemServicesDhcpLocalserverGroup) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *systemServicesDhcpLocalserverGroup) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *systemServicesDhcpLocalserverGroup) Configure(
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

func (rsc *systemServicesDhcpLocalserverGroup) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format" +
					" `<name>" + junos.IDSeparator + "<routing_instance>" + junos.IDSeparator + "<version>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Group name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
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
			"authentication_password": schema.StringAttribute{
				Optional:    true,
				Description: "DHCP authentication, username password to use.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringDoubleQuoteExclusion(),
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
			"reauthenticate_lease_renewal": schema.BoolAttribute{
				Optional:    true,
				Description: "Reauthenticate on each renew, rebind, DISCOVER or SOLICIT.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"reauthenticate_remote_id_mismatch": schema.BoolAttribute{
				Optional:    true,
				Description: "Reauthenticate on remote-id mismatch for renew, rebind and re-negotiation.",
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
		},
		Blocks: map[string]schema.Block{
			"authentication_username_include": schema.SingleNestedBlock{
				Description: "DHCP authentication, add username options.",
				Attributes:  dhcpBlockAuthenticationUsernameInclude{}.attributesSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"interface": schema.SetNestedBlock{
				Description: "For each name of interface to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Interface name.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
								stringvalidator.Any(
									tfvalidator.String1DotCount(),
									stringvalidator.OneOf("all"),
								),
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
						"exclude": schema.BoolAttribute{
							Optional:    true,
							Description: "Exclude this interface range.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
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
							Description: "Short cycle lockout min time in seconds.",
							Validators: []validator.Int64{
								int64validator.Between(1, 86400),
							},
						},
						"trace": schema.BoolAttribute{
							Optional:    true,
							Description: "Enable tracing for this interface.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"upto": schema.StringAttribute{
							Optional:    true,
							Description: "Interface up to.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
								tfvalidator.String1DotCount(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"overrides_v4": schema.SingleNestedBlock{
							Description: "DHCP override processing.",
							Attributes:  systemServicesDhcpLocalserverGroupBlockOverridesV4{}.attributesSchema(),
							Blocks:      systemServicesDhcpLocalserverGroupBlockOverridesV4{}.blocksSchema(),
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
						"overrides_v6": schema.SingleNestedBlock{
							Description: "DHCPv6 override processing.",
							Attributes:  systemServicesDhcpLocalserverGroupBlockOverridesV6{}.attributesSchema(),
							Blocks:      systemServicesDhcpLocalserverGroupBlockOverridesV6{}.blocksSchema(),
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
					},
				},
			},
			"lease_time_validation": schema.SingleNestedBlock{
				Description: "Configure lease time violation validation.",
				Attributes: map[string]schema.Attribute{
					"lease_time_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Threshold for lease time violation seconds (seconds).",
						Validators: []validator.Int64{
							int64validator.Between(60, 2147483647),
						},
					},
					"violation_action": schema.StringAttribute{
						Optional:    true,
						Description: " Lease time validation violation action.",
						Validators: []validator.String{
							stringvalidator.OneOf("override-lease", "strict"),
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
				Attributes:  systemServicesDhcpLocalserverGroupBlockOverridesV4{}.attributesSchema(),
				Blocks:      systemServicesDhcpLocalserverGroupBlockOverridesV4{}.blocksSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"overrides_v6": schema.SingleNestedBlock{
				Description: "DHCPv6 override processing.",
				Attributes:  systemServicesDhcpLocalserverGroupBlockOverridesV6{}.attributesSchema(),
				Blocks:      systemServicesDhcpLocalserverGroupBlockOverridesV6{}.blocksSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"reconfigure": schema.SingleNestedBlock{
				Description: "DHCP reconfigure processing.",
				Attributes: map[string]schema.Attribute{
					"attempts": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of reconfigure attempts before aborting.",
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
						},
					},
					"clear_on_abort": schema.BoolAttribute{
						Optional:    true,
						Description: "Delete client on reconfiguration abort.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"support_option_pd_exclude": schema.BoolAttribute{
						Optional:    true,
						Description: "Request prefix exclude option in reconfigure message.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Initial timeout value for retry.",
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
						},
					},
					"token": schema.StringAttribute{
						Optional:    true,
						Description: "Reconfigure token.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 244),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"trigger_radius_disconnect": schema.BoolAttribute{
						Optional:    true,
						Description: "Trigger DHCP reconfigure by radius initiated disconnect.",
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
	}
}

//nolint:lll
type systemServicesDhcpLocalserverGroupData struct {
	ID                                   types.String                                                `tfsdk:"id"                                       tfdata:"skip_isempty"`
	Name                                 types.String                                                `tfsdk:"name"                                     tfdata:"skip_isempty"`
	RoutingInstance                      types.String                                                `tfsdk:"routing_instance"                         tfdata:"skip_isempty"`
	Version                              types.String                                                `tfsdk:"version"                                  tfdata:"skip_isempty"`
	AccessProfile                        types.String                                                `tfsdk:"access_profile"`
	AuthenticationPassword               types.String                                                `tfsdk:"authentication_password"`
	DynamicProfile                       types.String                                                `tfsdk:"dynamic_profile"`
	DynamicProfileAggregateClients       types.Bool                                                  `tfsdk:"dynamic_profile_aggregate_clients"`
	DynamicProfileAggregateClientsAction types.String                                                `tfsdk:"dynamic_profile_aggregate_clients_action"`
	DynamicProfileUsePrimary             types.String                                                `tfsdk:"dynamic_profile_use_primary"`
	LivenessDetectionFailureAction       types.String                                                `tfsdk:"liveness_detection_failure_action"`
	ReauthenticateLeaseRenewal           types.Bool                                                  `tfsdk:"reauthenticate_lease_renewal"`
	ReauthenticateRemoteIDMismatch       types.Bool                                                  `tfsdk:"reauthenticate_remote_id_mismatch"`
	RemoteIDMismatchDisconnect           types.Bool                                                  `tfsdk:"remote_id_mismatch_disconnect"`
	RouteSuppressionAccess               types.Bool                                                  `tfsdk:"route_suppression_access"`
	RouteSuppressionAccessInternal       types.Bool                                                  `tfsdk:"route_suppression_access_internal"`
	RouteSuppressionDestination          types.Bool                                                  `tfsdk:"route_suppression_destination"`
	ServiceProfile                       types.String                                                `tfsdk:"service_profile"`
	ShortCycleProtectionLockoutMaxTime   types.Int64                                                 `tfsdk:"short_cycle_protection_lockout_max_time"`
	ShortCycleProtectionLockoutMinTime   types.Int64                                                 `tfsdk:"short_cycle_protection_lockout_min_time"`
	AuthenticationUsernameInclude        *dhcpBlockAuthenticationUsernameInclude                     `tfsdk:"authentication_username_include"`
	Interface                            []systemServicesDhcpLocalserverGroupBlockInterface          `tfsdk:"interface"`
	LeaseTimeValidation                  *systemServicesDhcpLocalserverGroupBlockLeaseTimeValidation `tfsdk:"lease_time_validation"`
	LivenessDetectionMethodBfd           *dhcpBlockLivenessDetectionMethodBfd                        `tfsdk:"liveness_detection_method_bfd"`
	LivenessDetectionMethodLayer2        *dhcpBlockLivenessDetectionMethodLayer2                     `tfsdk:"liveness_detection_method_layer2"`
	OverridesV4                          *systemServicesDhcpLocalserverGroupBlockOverridesV4         `tfsdk:"overrides_v4"`
	OverridesV6                          *systemServicesDhcpLocalserverGroupBlockOverridesV6         `tfsdk:"overrides_v6"`
	Reconfigure                          *systemServicesDhcpLocalserverGroupBlockReconfigure         `tfsdk:"reconfigure"`
}

func (rscData *systemServicesDhcpLocalserverGroupData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

//nolint:lll
type systemServicesDhcpLocalserverGroupConfig struct {
	ID                                   types.String                                                `tfsdk:"id"                                       tfdata:"skip_isempty"`
	Name                                 types.String                                                `tfsdk:"name"                                     tfdata:"skip_isempty"`
	RoutingInstance                      types.String                                                `tfsdk:"routing_instance"                         tfdata:"skip_isempty"`
	Version                              types.String                                                `tfsdk:"version"                                  tfdata:"skip_isempty"`
	AccessProfile                        types.String                                                `tfsdk:"access_profile"`
	AuthenticationPassword               types.String                                                `tfsdk:"authentication_password"`
	DynamicProfile                       types.String                                                `tfsdk:"dynamic_profile"`
	DynamicProfileAggregateClients       types.Bool                                                  `tfsdk:"dynamic_profile_aggregate_clients"`
	DynamicProfileAggregateClientsAction types.String                                                `tfsdk:"dynamic_profile_aggregate_clients_action"`
	DynamicProfileUsePrimary             types.String                                                `tfsdk:"dynamic_profile_use_primary"`
	LivenessDetectionFailureAction       types.String                                                `tfsdk:"liveness_detection_failure_action"`
	ReauthenticateLeaseRenewal           types.Bool                                                  `tfsdk:"reauthenticate_lease_renewal"`
	ReauthenticateRemoteIDMismatch       types.Bool                                                  `tfsdk:"reauthenticate_remote_id_mismatch"`
	RemoteIDMismatchDisconnect           types.Bool                                                  `tfsdk:"remote_id_mismatch_disconnect"`
	RouteSuppressionAccess               types.Bool                                                  `tfsdk:"route_suppression_access"`
	RouteSuppressionAccessInternal       types.Bool                                                  `tfsdk:"route_suppression_access_internal"`
	RouteSuppressionDestination          types.Bool                                                  `tfsdk:"route_suppression_destination"`
	ServiceProfile                       types.String                                                `tfsdk:"service_profile"`
	ShortCycleProtectionLockoutMaxTime   types.Int64                                                 `tfsdk:"short_cycle_protection_lockout_max_time"`
	ShortCycleProtectionLockoutMinTime   types.Int64                                                 `tfsdk:"short_cycle_protection_lockout_min_time"`
	AuthenticationUsernameInclude        *dhcpBlockAuthenticationUsernameInclude                     `tfsdk:"authentication_username_include"`
	Interface                            types.Set                                                   `tfsdk:"interface"`
	LeaseTimeValidation                  *systemServicesDhcpLocalserverGroupBlockLeaseTimeValidation `tfsdk:"lease_time_validation"`
	LivenessDetectionMethodBfd           *dhcpBlockLivenessDetectionMethodBfd                        `tfsdk:"liveness_detection_method_bfd"`
	LivenessDetectionMethodLayer2        *dhcpBlockLivenessDetectionMethodLayer2                     `tfsdk:"liveness_detection_method_layer2"`
	OverridesV4                          *systemServicesDhcpLocalserverGroupBlockOverridesV4Config   `tfsdk:"overrides_v4"`
	OverridesV6                          *systemServicesDhcpLocalserverGroupBlockOverridesV6Config   `tfsdk:"overrides_v6"`
	Reconfigure                          *systemServicesDhcpLocalserverGroupBlockReconfigure         `tfsdk:"reconfigure"`
}

func (config *systemServicesDhcpLocalserverGroupConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(config)
}

//nolint:lll
type systemServicesDhcpLocalserverGroupBlockInterface struct {
	Name                                 types.String                                        `tfsdk:"name"                                     tfdata:"identifier"`
	AccessProfile                        types.String                                        `tfsdk:"access_profile"`
	DynamicProfile                       types.String                                        `tfsdk:"dynamic_profile"`
	DynamicProfileAggregateClients       types.Bool                                          `tfsdk:"dynamic_profile_aggregate_clients"`
	DynamicProfileAggregateClientsAction types.String                                        `tfsdk:"dynamic_profile_aggregate_clients_action"`
	DynamicProfileUsePrimary             types.String                                        `tfsdk:"dynamic_profile_use_primary"`
	Exclude                              types.Bool                                          `tfsdk:"exclude"`
	ServiceProfile                       types.String                                        `tfsdk:"service_profile"`
	ShortCycleProtectionLockoutMaxTime   types.Int64                                         `tfsdk:"short_cycle_protection_lockout_max_time"`
	ShortCycleProtectionLockoutMinTime   types.Int64                                         `tfsdk:"short_cycle_protection_lockout_min_time"`
	Trace                                types.Bool                                          `tfsdk:"trace"`
	Upto                                 types.String                                        `tfsdk:"upto"`
	OverridesV4                          *systemServicesDhcpLocalserverGroupBlockOverridesV4 `tfsdk:"overrides_v4"`
	OverridesV6                          *systemServicesDhcpLocalserverGroupBlockOverridesV6 `tfsdk:"overrides_v6"`
}

//nolint:lll
type systemServicesDhcpLocalserverGroupBlockInterfaceConfig struct {
	Name                                 types.String                                              `tfsdk:"name"`
	AccessProfile                        types.String                                              `tfsdk:"access_profile"`
	DynamicProfile                       types.String                                              `tfsdk:"dynamic_profile"`
	DynamicProfileAggregateClients       types.Bool                                                `tfsdk:"dynamic_profile_aggregate_clients"`
	DynamicProfileAggregateClientsAction types.String                                              `tfsdk:"dynamic_profile_aggregate_clients_action"`
	DynamicProfileUsePrimary             types.String                                              `tfsdk:"dynamic_profile_use_primary"`
	Exclude                              types.Bool                                                `tfsdk:"exclude"`
	ServiceProfile                       types.String                                              `tfsdk:"service_profile"`
	ShortCycleProtectionLockoutMaxTime   types.Int64                                               `tfsdk:"short_cycle_protection_lockout_max_time"`
	ShortCycleProtectionLockoutMinTime   types.Int64                                               `tfsdk:"short_cycle_protection_lockout_min_time"`
	Trace                                types.Bool                                                `tfsdk:"trace"`
	Upto                                 types.String                                              `tfsdk:"upto"`
	OverridesV4                          *systemServicesDhcpLocalserverGroupBlockOverridesV4Config `tfsdk:"overrides_v4"`
	OverridesV6                          *systemServicesDhcpLocalserverGroupBlockOverridesV6Config `tfsdk:"overrides_v6"`
}

type systemServicesDhcpLocalserverGroupBlockLeaseTimeValidation struct {
	LeaseTimeThreshold types.Int64  `tfsdk:"lease_time_threshold"`
	ViolationAction    types.String `tfsdk:"violation_action"`
}

//nolint:lll
type systemServicesDhcpLocalserverGroupBlockOverridesV4 struct {
	AllowNoEndOption             types.Bool                                                          `tfsdk:"allow_no_end_option"`
	AsymmetricLeaseTime          types.Int64                                                         `tfsdk:"asymmetric_lease_time"`
	BootpSupport                 types.Bool                                                          `tfsdk:"bootp_support"`
	ClientDiscoverMatch          types.String                                                        `tfsdk:"client_discover_match"`
	DelayOfferDelayTime          types.Int64                                                         `tfsdk:"delay_offer_delay_time"`
	DeleteBindingOnRenegotiation types.Bool                                                          `tfsdk:"delete_binding_on_renegotiation"`
	DualStack                    types.String                                                        `tfsdk:"dual_stack"`
	IncludeOption82Forcerenew    types.Bool                                                          `tfsdk:"include_option_82_forcerenew"`
	IncludeOption82Nak           types.Bool                                                          `tfsdk:"include_option_82_nak"`
	InterfaceClientLimit         types.Int64                                                         `tfsdk:"interface_client_limit"`
	ProcessInform                types.Bool                                                          `tfsdk:"process_inform"`
	ProcessInformPool            types.String                                                        `tfsdk:"process_inform_pool"`
	ProtocolAttributes           types.String                                                        `tfsdk:"protocol_attributes"`
	DelayOfferBasedOn            []systemServicesDhcpLocalserverGroupBlockOverridesBlockDelayBasedOn `tfsdk:"delay_offer_based_on"`
}

func (block *systemServicesDhcpLocalserverGroupBlockOverridesV4) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (systemServicesDhcpLocalserverGroupBlockOverridesV4) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"allow_no_end_option": schema.BoolAttribute{
			Optional:    true,
			Description: "Allow packets without end-of-option.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"asymmetric_lease_time": schema.Int64Attribute{
			Optional:    true,
			Description: "Use a reduced lease time for the client (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(600, 86400),
			},
		},
		"bootp_support": schema.BoolAttribute{
			Optional:    true,
			Description: "Allow processing of bootp requests.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"client_discover_match": schema.StringAttribute{
			Optional:    true,
			Description: "Use incoming interface or option 60 and option 82 match criteria for DISCOVER PDU.",
			Validators: []validator.String{
				stringvalidator.OneOf("incoming-interface", "option60-and-option82"),
			},
		},
		"delay_offer_delay_time": schema.Int64Attribute{
			Optional:    true,
			Description: "Time delay between discover and offer (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(1, 30),
			},
		},
		"delete_binding_on_renegotiation": schema.BoolAttribute{
			Optional:    true,
			Description: "Delete binding on renegotiation.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"dual_stack": schema.StringAttribute{
			Optional:    true,
			Description: "Dual stack group to use.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"include_option_82_forcerenew": schema.BoolAttribute{
			Optional:    true,
			Description: "Include option-82 in FORCERENEW.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"include_option_82_nak": schema.BoolAttribute{
			Optional:    true,
			Description: "Include option-82 in NAK.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"interface_client_limit": schema.Int64Attribute{
			Optional:    true,
			Description: "Limit the number of clients allowed on an interface.",
			Validators: []validator.Int64{
				int64validator.Between(1, 500000),
			},
		},
		"process_inform": schema.BoolAttribute{
			Optional:    true,
			Description: "Process INFORM PDUs.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"process_inform_pool": schema.StringAttribute{
			Optional:    true,
			Description: "Pool name for family inet.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"protocol_attributes": schema.StringAttribute{
			Optional:    true,
			Description: "DHCPv4 attributes to use as defined under access protocol-attributes.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
	}
}

func (systemServicesDhcpLocalserverGroupBlockOverridesV4) blocksSchema() map[string]schema.Block {
	return map[string]schema.Block{
		"delay_offer_based_on": schema.SetNestedBlock{
			Description: "For each combination of block arguments, filter options for dhcp-server.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"option": schema.StringAttribute{
						Required:    true,
						Description: "Option.",
						Validators: []validator.String{
							stringvalidator.OneOf("option-60", "option-77", "option-82"),
						},
					},
					"compare": schema.StringAttribute{
						Required:    true,
						Description: "How to compare.",
						Validators: []validator.String{
							stringvalidator.OneOf("equals", "not-equals", "starts-with"),
						},
					},
					"value_type": schema.StringAttribute{
						Required:    true,
						Description: "Type of string.",
						Validators: []validator.String{
							stringvalidator.OneOf("ascii", "hexadecimal"),
						},
					},
					"value": schema.StringAttribute{
						Required:    true,
						Description: "String to compare.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 256),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
			},
		},
	}
}

type systemServicesDhcpLocalserverGroupBlockOverridesV4Config struct {
	AllowNoEndOption             types.Bool   `tfsdk:"allow_no_end_option"`
	AsymmetricLeaseTime          types.Int64  `tfsdk:"asymmetric_lease_time"`
	BootpSupport                 types.Bool   `tfsdk:"bootp_support"`
	ClientDiscoverMatch          types.String `tfsdk:"client_discover_match"`
	DelayOfferDelayTime          types.Int64  `tfsdk:"delay_offer_delay_time"`
	DeleteBindingOnRenegotiation types.Bool   `tfsdk:"delete_binding_on_renegotiation"`
	DualStack                    types.String `tfsdk:"dual_stack"`
	IncludeOption82Forcerenew    types.Bool   `tfsdk:"include_option_82_forcerenew"`
	IncludeOption82Nak           types.Bool   `tfsdk:"include_option_82_nak"`
	InterfaceClientLimit         types.Int64  `tfsdk:"interface_client_limit"`
	ProcessInform                types.Bool   `tfsdk:"process_inform"`
	ProcessInformPool            types.String `tfsdk:"process_inform_pool"`
	ProtocolAttributes           types.String `tfsdk:"protocol_attributes"`
	DelayOfferBasedOn            types.Set    `tfsdk:"delay_offer_based_on"`
}

func (block *systemServicesDhcpLocalserverGroupBlockOverridesV4Config) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *systemServicesDhcpLocalserverGroupBlockOverridesV4Config) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type systemServicesDhcpLocalserverGroupBlockOverridesBlockDelayBasedOn struct {
	Option    types.String `tfsdk:"option"`
	Compare   types.String `tfsdk:"compare"`
	ValueType types.String `tfsdk:"value_type"`
	Value     types.String `tfsdk:"value"`
}

//nolint:lll
type systemServicesDhcpLocalserverGroupBlockOverridesV6 struct {
	AlwaysAddOptionDNSServer                types.Bool                                                          `tfsdk:"always_add_option_dns_server"`
	AlwaysProcessOptionRequestOption        types.Bool                                                          `tfsdk:"always_process_option_request_option"`
	AsymmetricLeaseTime                     types.Int64                                                         `tfsdk:"asymmetric_lease_time"`
	AsymmetricPrefixLeaseTime               types.Int64                                                         `tfsdk:"asymmetric_prefix_lease_time"`
	ClientNegotiationMatchIncomingInterface types.Bool                                                          `tfsdk:"client_negotiation_match_incoming_interface"`
	DelayAdvertiseDelayTime                 types.Int64                                                         `tfsdk:"delay_advertise_delay_time"`
	DelegatedPool                           types.String                                                        `tfsdk:"delegated_pool"`
	DeleteBindingOnRenegotiation            types.Bool                                                          `tfsdk:"delete_binding_on_renegotiation"`
	DualStack                               types.String                                                        `tfsdk:"dual_stack"`
	InterfaceClientLimit                    types.Int64                                                         `tfsdk:"interface_client_limit"`
	MultiAddressEmbeddedOptionResponse      types.Bool                                                          `tfsdk:"multi_address_embedded_option_response"`
	ProcessInform                           types.Bool                                                          `tfsdk:"process_inform"`
	ProcessInformPool                       types.String                                                        `tfsdk:"process_inform_pool"`
	ProtocolAttributes                      types.String                                                        `tfsdk:"protocol_attributes"`
	RapidCommit                             types.Bool                                                          `tfsdk:"rapid_commit"`
	TopLevelStatusCode                      types.Bool                                                          `tfsdk:"top_level_status_code"`
	DelayAdvertiseBasedOn                   []systemServicesDhcpLocalserverGroupBlockOverridesBlockDelayBasedOn `tfsdk:"delay_advertise_based_on"`
}

func (block *systemServicesDhcpLocalserverGroupBlockOverridesV6) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (systemServicesDhcpLocalserverGroupBlockOverridesV6) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"always_add_option_dns_server": schema.BoolAttribute{
			Optional:    true,
			Description: "Add option-23, DNS recursive name server in Advertise and Reply.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"always_process_option_request_option": schema.BoolAttribute{
			Optional:    true,
			Description: "Always process option even after address allocation failure.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"asymmetric_lease_time": schema.Int64Attribute{
			Optional:    true,
			Description: "Use a reduced lease time for the client. In seconds (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(600, 86400),
			},
		},
		"asymmetric_prefix_lease_time": schema.Int64Attribute{
			Optional:    true,
			Description: "Use a reduced prefix lease time for the client. In seconds (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(600, 86400),
			},
		},
		"client_negotiation_match_incoming_interface": schema.BoolAttribute{
			Optional:    true,
			Description: "Use incoming interface match criteria for SOLICIT PDU.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"delay_advertise_delay_time": schema.Int64Attribute{
			Optional:    true,
			Description: "Time delay between solicit and advertise (seconds).",
			Validators: []validator.Int64{
				int64validator.Between(1, 30),
			},
		},
		"delegated_pool": schema.StringAttribute{
			Optional:    true,
			Description: "Delegated pool name for inet6.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"delete_binding_on_renegotiation": schema.BoolAttribute{
			Optional:    true,
			Description: "Delete binding on renegotiation.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"dual_stack": schema.StringAttribute{
			Optional:    true,
			Description: "Dual stack group to use.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"interface_client_limit": schema.Int64Attribute{
			Optional:    true,
			Description: "Limit the number of clients allowed on an interface.",
			Validators: []validator.Int64{
				int64validator.Between(1, 500000),
			},
		},
		"multi_address_embedded_option_response": schema.BoolAttribute{
			Optional:    true,
			Description: "If the client requests multiple addresses place the options in each address.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"process_inform": schema.BoolAttribute{
			Optional:    true,
			Description: "Process INFORM PDUs.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"process_inform_pool": schema.StringAttribute{
			Optional:    true,
			Description: "Pool name for family inet6.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"protocol_attributes": schema.StringAttribute{
			Optional:    true,
			Description: "DHCPv6 attributes to use as defined under access protocol-attributes.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 64),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"rapid_commit": schema.BoolAttribute{
			Optional:    true,
			Description: "Enable rapid commit processing.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"top_level_status_code": schema.BoolAttribute{
			Optional:    true,
			Description: "A top level status code option rather than encapsulated in IA for NoAddrsAvail in Advertise PDUs.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
	}
}

func (systemServicesDhcpLocalserverGroupBlockOverridesV6) blocksSchema() map[string]schema.Block {
	return map[string]schema.Block{
		"delay_advertise_based_on": schema.SetNestedBlock{
			Description: "For each combination of block arguments, filter options for dhcp-server.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"option": schema.StringAttribute{
						Required:    true,
						Description: "Option.",
						Validators: []validator.String{
							stringvalidator.OneOf("option-15", "option-16", "option-18", "option-37"),
						},
					},
					"compare": schema.StringAttribute{
						Required:    true,
						Description: "How to compare.",
						Validators: []validator.String{
							stringvalidator.OneOf("equals", "not-equals", "starts-with"),
						},
					},
					"value_type": schema.StringAttribute{
						Required:    true,
						Description: "Type of string.",
						Validators: []validator.String{
							stringvalidator.OneOf("ascii", "hexadecimal"),
						},
					},
					"value": schema.StringAttribute{
						Required:    true,
						Description: "String to compare.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 256),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
			},
		},
	}
}

type systemServicesDhcpLocalserverGroupBlockOverridesV6Config struct {
	AlwaysAddOptionDNSServer                types.Bool   `tfsdk:"always_add_option_dns_server"`
	AlwaysProcessOptionRequestOption        types.Bool   `tfsdk:"always_process_option_request_option"`
	AsymmetricLeaseTime                     types.Int64  `tfsdk:"asymmetric_lease_time"`
	AsymmetricPrefixLeaseTime               types.Int64  `tfsdk:"asymmetric_prefix_lease_time"`
	ClientNegotiationMatchIncomingInterface types.Bool   `tfsdk:"client_negotiation_match_incoming_interface"`
	DelayAdvertiseDelayTime                 types.Int64  `tfsdk:"delay_advertise_delay_time"`
	DelegatedPool                           types.String `tfsdk:"delegated_pool"`
	DeleteBindingOnRenegotiation            types.Bool   `tfsdk:"delete_binding_on_renegotiation"`
	DualStack                               types.String `tfsdk:"dual_stack"`
	InterfaceClientLimit                    types.Int64  `tfsdk:"interface_client_limit"`
	MultiAddressEmbeddedOptionResponse      types.Bool   `tfsdk:"multi_address_embedded_option_response"`
	ProcessInform                           types.Bool   `tfsdk:"process_inform"`
	ProcessInformPool                       types.String `tfsdk:"process_inform_pool"`
	ProtocolAttributes                      types.String `tfsdk:"protocol_attributes"`
	RapidCommit                             types.Bool   `tfsdk:"rapid_commit"`
	TopLevelStatusCode                      types.Bool   `tfsdk:"top_level_status_code"`
	DelayAdvertiseBasedOn                   types.Set    `tfsdk:"delay_advertise_based_on"`
}

func (block *systemServicesDhcpLocalserverGroupBlockOverridesV6Config) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *systemServicesDhcpLocalserverGroupBlockOverridesV6Config) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type systemServicesDhcpLocalserverGroupBlockReconfigure struct {
	Attempts                types.Int64  `tfsdk:"attempts"`
	ClearOnAbort            types.Bool   `tfsdk:"clear_on_abort"`
	SupportOptionPdExclude  types.Bool   `tfsdk:"support_option_pd_exclude"`
	Timeout                 types.Int64  `tfsdk:"timeout"`
	Token                   types.String `tfsdk:"token"`
	TriggerRadiusDisconnect types.Bool   `tfsdk:"trigger_radius_disconnect"`
}

func (rsc *systemServicesDhcpLocalserverGroup) ValidateConfig( //nolint:gocyclo
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config systemServicesDhcpLocalserverGroupConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one of arguments need to be set (in addition to `name`, `routing_instance` and `version`)",
		)
	}

	version := config.Version.ValueString()
	switch {
	case config.Version.IsUnknown():
	case version == "v4" || config.Version.IsNull():
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
		if config.OverridesV6 != nil && config.OverridesV6.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("overrides_v6"),
				tfdiag.ConflictConfigErrSummary,
				"overrides_v6 cannot be configured when version = v4",
			)
		}
	case version == "v6":
		if !config.RouteSuppressionDestination.IsNull() &&
			!config.RouteSuppressionDestination.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("route_suppression_destination"),
				tfdiag.ConflictConfigErrSummary,
				"route_suppression_destination cannot be configured when version = v6",
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
		if config.OverridesV4 != nil && config.OverridesV4.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("overrides_v4"),
				tfdiag.ConflictConfigErrSummary,
				"overrides_v4 cannot be configured when version = v6",
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
	if !config.Interface.IsNull() &&
		!config.Interface.IsUnknown() {
		var configInterface []systemServicesDhcpLocalserverGroupBlockInterfaceConfig
		asDiags := config.Interface.ElementsAs(ctx, &configInterface, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		interfaceName := make(map[string]struct{})
		for _, block := range configInterface {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := interfaceName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple interface blocks with the same name %q", name),
					)
				}
				interfaceName[name] = struct{}{}
			}

			if !block.DynamicProfileAggregateClients.IsNull() &&
				!block.DynamicProfileAggregateClients.IsUnknown() &&
				!block.DynamicProfileUsePrimary.IsNull() &&
				!block.DynamicProfileUsePrimary.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("interface"),
					tfdiag.ConflictConfigErrSummary,
					fmt.Sprintf("dynamic_profile_aggregate_clients and dynamic_profile_use_primary cannot be configured together"+
						" in interface block %q", block.Name.ValueString()),
				)
			}
			if !block.DynamicProfileAggregateClients.IsNull() &&
				!block.DynamicProfileAggregateClients.IsUnknown() &&
				block.DynamicProfile.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("interface"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("dynamic_profile must be specified with dynamic_profile_aggregate_clients"+
						" in interface block %q", block.Name.ValueString()),
				)
			}
			if !block.DynamicProfileAggregateClientsAction.IsNull() &&
				!block.DynamicProfileAggregateClientsAction.IsUnknown() &&
				block.DynamicProfileAggregateClients.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("interface"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("dynamic_profile_aggregate_clients must be specified with dynamic_profile_aggregate_clients_action"+
						" in interface block %q", block.Name.ValueString()),
				)
			}
			if !block.DynamicProfileUsePrimary.IsNull() &&
				!block.DynamicProfileUsePrimary.IsUnknown() &&
				block.DynamicProfile.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("interface"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("dynamic_profile must be specified with dynamic_profile_use_primary"+
						" in interface block %q", block.Name.ValueString()),
				)
			}
			if block.OverridesV4 != nil {
				if block.OverridesV4.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("overrides_v4 block is empty"+
							" in interface block %q", block.Name.ValueString()),
					)
				}
				if block.OverridesV4.hasKnownValue() && version == "v6" {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("overrides_v4 cannot be configured when version = v6"+
							" in interface block %q", block.Name.ValueString()),
					)
				}

				if !block.OverridesV4.DelayOfferDelayTime.IsNull() &&
					!block.OverridesV4.DelayOfferDelayTime.IsUnknown() &&
					block.OverridesV4.DelayOfferBasedOn.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("delay_offer_based_on must be specified with delay_offer_delay_time"+
							" in overrides_v4 block in interface block %q", block.Name.ValueString()),
					)
				}
			}
			if block.OverridesV6 != nil {
				if block.OverridesV6.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("overrides_v6 block is empty"+
							" in interface block %q", block.Name.ValueString()),
					)
				}
				if block.OverridesV6.hasKnownValue() &&
					(version == "v4" || config.Version.IsNull()) {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("overrides_v6 cannot be configured when version = v4"+
							" in interface block %q", block.Name.ValueString()),
					)
				}

				if !block.OverridesV6.DelayAdvertiseDelayTime.IsNull() &&
					!block.OverridesV6.DelayAdvertiseDelayTime.IsUnknown() &&
					block.OverridesV6.DelayAdvertiseBasedOn.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("delay_advertise_based_on must be specified with delay_advertise_delay_time"+
							" in overrides_v6 block in interface block %q", block.Name.ValueString()),
					)
				}
			}
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

		if !config.OverridesV4.DelayOfferDelayTime.IsNull() &&
			!config.OverridesV4.DelayOfferDelayTime.IsUnknown() &&
			config.OverridesV4.DelayOfferBasedOn.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("overrides_v4").AtName("delay_offer_delay_time"),
				tfdiag.MissingConfigErrSummary,
				"delay_offer_based_on must be specified with delay_offer_delay_time"+
					" in overrides_v4 block",
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

		if !config.OverridesV6.DelayAdvertiseDelayTime.IsNull() &&
			!config.OverridesV6.DelayAdvertiseDelayTime.IsUnknown() &&
			config.OverridesV6.DelayAdvertiseBasedOn.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("overrides_v6").AtName("delay_advertise_delay_time"),
				tfdiag.MissingConfigErrSummary,
				"delay_advertise_based_on must be specified with delay_advertise_delay_time"+
					" in overrides_v6 block",
			)
		}
	}
}

func (rsc *systemServicesDhcpLocalserverGroup) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemServicesDhcpLocalserverGroupData
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
			groupExists, err := checkSystemServicesDhcpLocalserverGroupExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				plan.Version.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if groupExists {
				if plan.Version.ValueString() == "v6" {
					if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf(
								"system services dhcp-local-server dhcpv6 group %q already exists in routing-instance %q",
								plan.Name.ValueString(), v,
							),
						)
					} else {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf(
								"system services dhcp-local-server dhcpv6 group %q already exists",
								plan.Name.ValueString(),
							),
						)
					}
				} else {
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
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			groupExists, err := checkSystemServicesDhcpLocalserverGroupExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				plan.Version.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !groupExists {
				if plan.Version.ValueString() == "v6" {
					if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf(
								"system services dhcp-local-server dhcpv6 group %q does not exists in routing-instance %q after commit "+
									"=> check your config", plan.Name.ValueString(), v,
							),
						)
					} else {
						resp.Diagnostics.AddError(
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("system services dhcp-local-server dhcpv6 group %q does not exists after commit "+
								"=> check your config", plan.Name.ValueString()),
						)
					}
				} else {
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
				}

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *systemServicesDhcpLocalserverGroup) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemServicesDhcpLocalserverGroupData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom3String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.RoutingInstance.ValueString(),
			state.Version.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *systemServicesDhcpLocalserverGroup) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemServicesDhcpLocalserverGroupData
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

func (rsc *systemServicesDhcpLocalserverGroup) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state systemServicesDhcpLocalserverGroupData
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

func (rsc *systemServicesDhcpLocalserverGroup) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data systemServicesDhcpLocalserverGroupData

	var _ resourceDataReadFrom3String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <name>"+junos.IDSeparator+"<routing_instance>"+junos.IDSeparator+"<version>)",
	)
}

func checkSystemServicesDhcpLocalserverGroupExists(
	_ context.Context, name, routingInstance, version string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showPrefix += "system services dhcp-local-server "
	if version == "v6" {
		showPrefix += "dhcpv6 "
	}
	showConfig, err := junSess.Command(showPrefix +
		"group " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *systemServicesDhcpLocalserverGroupData) fillID() {
	routingInstance := rscData.RoutingInstance.ValueString()
	version := rscData.Version.ValueString()
	if routingInstance == "" {
		routingInstance = junos.DefaultW
	}
	if version == "" {
		version = "v4"
	}

	rscData.ID = types.StringValue(rscData.Name.ValueString() +
		junos.IDSeparator + routingInstance +
		junos.IDSeparator + version)
}

func (rscData *systemServicesDhcpLocalserverGroupData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemServicesDhcpLocalserverGroupData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	if rscData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name`, `routing_instance` and `version`)")
	}

	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "system services dhcp-local-server "
	version := rscData.Version.ValueString()
	if version == "v6" {
		setPrefix += "dhcpv6 "
	} else if version != "v4" {
		version = "v4"
	}
	setPrefix += "group " + rscData.Name.ValueString() + " "

	if v := rscData.AccessProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"access-profile \""+v+"\"")
	}
	if v := rscData.AuthenticationPassword.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication password \""+v+"\"")
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
	if v := rscData.LivenessDetectionFailureAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"liveness-detection failure-action "+v)
	}
	if rscData.ReauthenticateLeaseRenewal.ValueBool() {
		configSet = append(configSet, setPrefix+"reauthenticate lease-renewal")
	}
	if rscData.ReauthenticateRemoteIDMismatch.ValueBool() {
		configSet = append(configSet, setPrefix+"reauthenticate remote-id-mismatch")
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
	interfaceName := make(map[string]struct{})
	for _, block := range rscData.Interface {
		ifaceName := block.Name.ValueString()
		if _, ok := interfaceName[ifaceName]; ok {
			return path.Root("interface"),
				fmt.Errorf("multiple interface blocks with the same name %q", ifaceName)
		}
		interfaceName[ifaceName] = struct{}{}

		blockSet, err := block.configSet(setPrefix, version)
		if err != nil {
			return path.Root("interface"), err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.LeaseTimeValidation != nil {
		configSet = append(configSet, setPrefix+"lease-time-validation")

		if !rscData.LeaseTimeValidation.LeaseTimeThreshold.IsNull() {
			configSet = append(configSet, setPrefix+"lease-time-validation lease-time-threshold "+
				utils.ConvI64toa(rscData.LeaseTimeValidation.LeaseTimeThreshold.ValueInt64()))
		}
		if v := rscData.LeaseTimeValidation.ViolationAction.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"lease-time-validation violation-action "+v)
		}
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
	if rscData.Reconfigure != nil {
		configSet = append(configSet, rscData.Reconfigure.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *systemServicesDhcpLocalserverGroupBlockInterface) configSet(
	setPrefix, version string,
) (
	[]string, // configSet
	error, // error
) {
	setPrefix += "interface " + block.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if v := block.AccessProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"access-profile \""+v+"\"")
	}
	if v := block.DynamicProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"dynamic-profile \""+v+"\"")
		if block.DynamicProfileAggregateClients.ValueBool() {
			configSet = append(configSet, setPrefix+"dynamic-profile aggregate-clients")
			if vv := block.DynamicProfileAggregateClientsAction.ValueString(); vv != "" {
				configSet = append(configSet, setPrefix+"dynamic-profile aggregate-clients "+vv)
			}
		} else if block.DynamicProfileAggregateClientsAction.ValueString() != "" {
			return configSet,
				fmt.Errorf("dynamic_profile_aggregate_clients must be specified with "+
					"dynamic_profile_aggregate_clients_action"+
					" in interface block %q", block.Name.ValueString())
		}
		if vv := block.DynamicProfileUsePrimary.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"dynamic-profile use-primary \""+vv+"\"")
		}
	} else {
		if block.DynamicProfileAggregateClients.ValueBool() {
			return configSet,
				fmt.Errorf("dynamic_profile must be specified with "+
					"dynamic_profile_aggregate_clients"+
					" in interface block %q", block.Name.ValueString())
		}
		if block.DynamicProfileAggregateClientsAction.ValueString() != "" {
			return configSet,
				fmt.Errorf("dynamic_profile must be specified with "+
					"dynamic_profile_aggregate_clients_action"+
					" in interface block %q", block.Name.ValueString())
		}
		if block.DynamicProfileUsePrimary.ValueString() != "" {
			return configSet,
				fmt.Errorf("dynamic_profile must be specified with "+
					"dynamic_profile_use_primary"+
					" in interface block %q", block.Name.ValueString())
		}
	}
	if block.Exclude.ValueBool() {
		configSet = append(configSet, setPrefix+"exclude")
	}
	if v := block.ServiceProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"service-profile \""+v+"\"")
	}
	if !block.ShortCycleProtectionLockoutMaxTime.IsNull() {
		configSet = append(configSet, setPrefix+"short-cycle-protection lockout-max-time "+
			utils.ConvI64toa(block.ShortCycleProtectionLockoutMaxTime.ValueInt64()))
	}
	if !block.ShortCycleProtectionLockoutMinTime.IsNull() {
		configSet = append(configSet, setPrefix+"short-cycle-protection lockout-min-time "+
			utils.ConvI64toa(block.ShortCycleProtectionLockoutMinTime.ValueInt64()))
	}
	if block.Trace.ValueBool() {
		configSet = append(configSet, setPrefix+"trace")
	}
	if v := block.Upto.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"upto "+v)
	}

	if block.OverridesV4 != nil {
		if block.OverridesV4.isEmpty() {
			return configSet,
				fmt.Errorf("overrides_v4 block is empty"+
					" in interface block %q", block.Name.ValueString())
		}
		if version == "v6" {
			return configSet,
				fmt.Errorf("overrides_v4 cannot be configured when version = v6"+
					" in interface block %q", block.Name.ValueString())
		}

		blockSet, _, err := block.OverridesV4.configSet(setPrefix)
		if err != nil {
			return configSet,
				fmt.Errorf(err.Error()+" in interface block %q", block.Name.ValueString())
		}
		configSet = append(configSet, blockSet...)
	}
	if block.OverridesV6 != nil {
		if block.OverridesV6.isEmpty() {
			return configSet,
				fmt.Errorf("overrides_v6 block is empty"+
					" in interface block %q", block.Name.ValueString())
		}
		if version == "v4" {
			return configSet,
				fmt.Errorf("overrides_v6 cannot be configured when version = v4"+
					" in interface block %q", block.Name.ValueString())
		}

		blockSet, _, err := block.OverridesV6.configSet(setPrefix)
		if err != nil {
			return configSet,
				fmt.Errorf(err.Error()+" in interface block %q", block.Name.ValueString())
		}
		configSet = append(configSet, blockSet...)
	}

	return configSet, nil
}

func (block *systemServicesDhcpLocalserverGroupBlockOverridesV4) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "overrides "

	if block.AllowNoEndOption.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-no-end-option")
	}
	if !block.AsymmetricLeaseTime.IsNull() {
		configSet = append(configSet, setPrefix+"asymmetric-lease-time "+
			utils.ConvI64toa(block.AsymmetricLeaseTime.ValueInt64()))
	}
	if block.BootpSupport.ValueBool() {
		configSet = append(configSet, setPrefix+"bootp-support")
	}
	if v := block.ClientDiscoverMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"client-discover-match "+v)
	}
	if !block.DelayOfferDelayTime.IsNull() {
		if len(block.DelayOfferBasedOn) == 0 {
			return configSet,
				path.Root("overrides_v4").AtName("delay_offer_delay_time"),
				errors.New("delay_offer_based_on must be specified with delay_offer_delay_time" +
					" in overrides_v4 block")
		}

		configSet = append(configSet, setPrefix+"delay-offer delay-time "+
			utils.ConvI64toa(block.DelayOfferDelayTime.ValueInt64()))
	}
	if block.DeleteBindingOnRenegotiation.ValueBool() {
		configSet = append(configSet, setPrefix+"delete-binding-on-renegotiation")
	}
	if v := block.DualStack.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"dual-stack \""+v+"\"")
	}
	if block.IncludeOption82Forcerenew.ValueBool() {
		configSet = append(configSet, setPrefix+"include-option-82 forcerenew")
	}
	if block.IncludeOption82Nak.ValueBool() {
		configSet = append(configSet, setPrefix+"include-option-82 nak")
	}
	if !block.InterfaceClientLimit.IsNull() {
		configSet = append(configSet, setPrefix+"interface-client-limit "+
			utils.ConvI64toa(block.InterfaceClientLimit.ValueInt64()))
	}
	if block.ProcessInform.ValueBool() {
		configSet = append(configSet, setPrefix+"process-inform")

		if v := block.ProcessInformPool.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"process-inform pool \""+v+"\"")
		}
	} else if !block.ProcessInformPool.IsNull() {
		return configSet,
			path.Root("overrides_v4").AtName("process_inform_pool"),
			errors.New("process_inform must be specified with process_inform_pool" +
				" in overrides_v4 block")
	}
	if v := block.ProtocolAttributes.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"protocol-attributes \""+v+"\"")
	}

	for _, subBlock := range block.DelayOfferBasedOn {
		configSet = append(configSet,
			setPrefix+"delay-offer based-on "+
				subBlock.Option.ValueString()+" "+
				subBlock.Compare.ValueString()+" "+
				subBlock.ValueType.ValueString()+" "+
				"\""+subBlock.Value.ValueString()+"\"")
	}

	return configSet, path.Empty(), nil
}

func (block *systemServicesDhcpLocalserverGroupBlockOverridesV6) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "overrides "

	if block.AlwaysAddOptionDNSServer.ValueBool() {
		configSet = append(configSet, setPrefix+"always-add-option-dns-server")
	}
	if block.AlwaysProcessOptionRequestOption.ValueBool() {
		configSet = append(configSet, setPrefix+"always-process-option-request-option")
	}
	if !block.AsymmetricLeaseTime.IsNull() {
		configSet = append(configSet, setPrefix+"asymmetric-lease-time "+
			utils.ConvI64toa(block.AsymmetricLeaseTime.ValueInt64()))
	}
	if !block.AsymmetricPrefixLeaseTime.IsNull() {
		configSet = append(configSet, setPrefix+"asymmetric-prefix-lease-time "+
			utils.ConvI64toa(block.AsymmetricPrefixLeaseTime.ValueInt64()))
	}
	if block.ClientNegotiationMatchIncomingInterface.ValueBool() {
		configSet = append(configSet, setPrefix+"client-negotiation-match incoming-interface")
	}
	if !block.DelayAdvertiseDelayTime.IsNull() {
		if len(block.DelayAdvertiseBasedOn) == 0 {
			return configSet,
				path.Root("overrides_v6").AtName("delay_advertise_delay_time"),
				errors.New("delay_advertise_based_on must be specified with delay_advertise_delay_time" +
					" in overrides_v4 block")
		}

		configSet = append(configSet, setPrefix+"delay-advertise delay-time "+
			utils.ConvI64toa(block.DelayAdvertiseDelayTime.ValueInt64()))
	}
	if v := block.DelegatedPool.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"delegated-pool \""+v+"\"")
	}
	if block.DeleteBindingOnRenegotiation.ValueBool() {
		configSet = append(configSet, setPrefix+"delete-binding-on-renegotiation")
	}
	if v := block.DualStack.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"dual-stack \""+v+"\"")
	}
	if !block.InterfaceClientLimit.IsNull() {
		configSet = append(configSet, setPrefix+"interface-client-limit "+
			utils.ConvI64toa(block.InterfaceClientLimit.ValueInt64()))
	}
	if block.MultiAddressEmbeddedOptionResponse.ValueBool() {
		configSet = append(configSet, setPrefix+"multi-address-embedded-option-response")
	}
	if block.ProcessInform.ValueBool() {
		configSet = append(configSet, setPrefix+"process-inform")

		if v := block.ProcessInformPool.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"process-inform pool \""+v+"\"")
		}
	} else if !block.ProcessInformPool.IsNull() {
		return configSet,
			path.Root("overrides_v6").AtName("process_inform_pool"),
			errors.New("process_inform must be specified with process_inform_pool" +
				" in overrides_v6 block")
	}
	if v := block.ProtocolAttributes.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"protocol-attributes \""+v+"\"")
	}
	if block.RapidCommit.ValueBool() {
		configSet = append(configSet, setPrefix+"rapid-commit")
	}
	if block.TopLevelStatusCode.ValueBool() {
		configSet = append(configSet, setPrefix+"top-level-status-code")
	}

	for _, subBlock := range block.DelayAdvertiseBasedOn {
		configSet = append(configSet,
			setPrefix+"delay-advertise based-on "+
				subBlock.Option.ValueString()+" "+
				subBlock.Compare.ValueString()+" "+
				subBlock.ValueType.ValueString()+" "+
				"\""+subBlock.Value.ValueString()+"\"")
	}

	return configSet, path.Empty(), nil
}

func (block *systemServicesDhcpLocalserverGroupBlockReconfigure) configSet(setPrefix string) []string {
	setPrefix += "reconfigure "

	configSet := []string{
		setPrefix,
	}

	if !block.Attempts.IsNull() {
		configSet = append(configSet, setPrefix+"attempts "+
			utils.ConvI64toa(block.Attempts.ValueInt64()))
	}
	if block.ClearOnAbort.ValueBool() {
		configSet = append(configSet, setPrefix+"clear-on-abort")
	}
	if block.SupportOptionPdExclude.ValueBool() {
		configSet = append(configSet, setPrefix+"support-option-pd-exclude")
	}
	if !block.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(block.Timeout.ValueInt64()))
	}
	if v := block.Token.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"token \""+v+"\"")
	}
	if block.TriggerRadiusDisconnect.ValueBool() {
		configSet = append(configSet, setPrefix+"trigger radius-disconnect")
	}

	return configSet
}

func (rscData *systemServicesDhcpLocalserverGroupData) read(
	_ context.Context, name, routingInstance, version string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showPrefix += "system services dhcp-local-server "
	if version == "v6" {
		showPrefix += "dhcpv6 "
	}
	showConfig, err := junSess.Command(showPrefix +
		"group " + name + junos.PipeDisplaySetRelative)
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
		if version == "v6" {
			rscData.Version = types.StringValue(version)
		} else {
			rscData.Version = types.StringValue("v4")
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
			case balt.CutPrefixInString(&itemTrim, "access-profile "):
				rscData.AccessProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "authentication password "):
				rscData.AuthenticationPassword = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "authentication username-include "):
				if rscData.AuthenticationUsernameInclude == nil {
					rscData.AuthenticationUsernameInclude = &dhcpBlockAuthenticationUsernameInclude{}
				}

				rscData.AuthenticationUsernameInclude.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "dynamic-profile aggregate-clients"):
				rscData.DynamicProfileAggregateClients = types.BoolValue(true)
				if balt.CutPrefixInString(&itemTrim, " ") {
					rscData.DynamicProfileAggregateClientsAction = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "dynamic-profile use-primary "):
				rscData.DynamicProfileUsePrimary = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "dynamic-profile "):
				rscData.DynamicProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "interface "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				var interFace systemServicesDhcpLocalserverGroupBlockInterface
				rscData.Interface, interFace = tfdata.ExtractBlock(rscData.Interface, types.StringValue(name))

				if balt.CutPrefixInString(&itemTrim, name+" ") {
					if err := interFace.read(itemTrim, rscData.Version.ValueString()); err != nil {
						return err
					}
				}
				rscData.Interface = append(rscData.Interface, interFace)
			case balt.CutPrefixInString(&itemTrim, "lease-time-validation"):
				if rscData.LeaseTimeValidation == nil {
					rscData.LeaseTimeValidation = &systemServicesDhcpLocalserverGroupBlockLeaseTimeValidation{}
				}

				switch {
				case balt.CutPrefixInString(&itemTrim, " lease-time-threshold "):
					rscData.LeaseTimeValidation.LeaseTimeThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " violation-action "):
					rscData.LeaseTimeValidation.ViolationAction = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "liveness-detection failure-action "):
				rscData.LivenessDetectionFailureAction = types.StringValue(itemTrim)
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
						rscData.OverridesV4 = &systemServicesDhcpLocalserverGroupBlockOverridesV4{}
					}

					err = rscData.OverridesV4.read(itemTrim)
				case "v6":
					if rscData.OverridesV6 == nil {
						rscData.OverridesV6 = &systemServicesDhcpLocalserverGroupBlockOverridesV6{}
					}

					err = rscData.OverridesV6.read(itemTrim)
				}
				if err != nil {
					return err
				}
			case itemTrim == "reauthenticate lease-renewal":
				rscData.ReauthenticateLeaseRenewal = types.BoolValue(true)
			case itemTrim == "reauthenticate remote-id-mismatch":
				rscData.ReauthenticateRemoteIDMismatch = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "reconfigure"):
				if rscData.Reconfigure == nil {
					rscData.Reconfigure = &systemServicesDhcpLocalserverGroupBlockReconfigure{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.Reconfigure.read(itemTrim); err != nil {
						return err
					}
				}
			case itemTrim == "remote-id-mismatch disconnect":
				rscData.RemoteIDMismatchDisconnect = types.BoolValue(true)
			case itemTrim == "route-suppression access":
				rscData.RouteSuppressionAccess = types.BoolValue(true)
			case itemTrim == "route-suppression access-internal":
				rscData.RouteSuppressionAccessInternal = types.BoolValue(true)
			case itemTrim == "route-suppression destination":
				rscData.RouteSuppressionDestination = types.BoolValue(true)
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
			}
		}
	}

	return nil
}

func (block *systemServicesDhcpLocalserverGroupBlockInterface) read(itemTrim, version string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "access-profile "):
		block.AccessProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "dynamic-profile "):
		switch {
		case balt.CutPrefixInString(&itemTrim, "use-primary "):
			block.DynamicProfileUsePrimary = types.StringValue(strings.Trim(itemTrim, "\""))
		case balt.CutPrefixInString(&itemTrim, "aggregate-clients"):
			block.DynamicProfileAggregateClients = types.BoolValue(true)
			if balt.CutPrefixInString(&itemTrim, " ") {
				block.DynamicProfileAggregateClientsAction = types.StringValue(itemTrim)
			}
		default:
			block.DynamicProfile = types.StringValue(strings.Trim(itemTrim, "\""))
		}
	case itemTrim == "exclude":
		block.Exclude = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "overrides "):
		switch version {
		case "v4":
			if block.OverridesV4 == nil {
				block.OverridesV4 = &systemServicesDhcpLocalserverGroupBlockOverridesV4{}
			}

			err = block.OverridesV4.read(itemTrim)
		case "v6":
			if block.OverridesV6 == nil {
				block.OverridesV6 = &systemServicesDhcpLocalserverGroupBlockOverridesV6{}
			}

			err = block.OverridesV6.read(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "service-profile "):
		block.ServiceProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "short-cycle-protection lockout-max-time "):
		block.ShortCycleProtectionLockoutMaxTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "short-cycle-protection lockout-min-time "):
		block.ShortCycleProtectionLockoutMinTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "trace":
		block.Trace = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "upto "):
		block.Upto = types.StringValue(itemTrim)
	}
	if err != nil {
		return err
	}

	return nil
}

func (block *systemServicesDhcpLocalserverGroupBlockOverridesV4) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "allow-no-end-option":
		block.AllowNoEndOption = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "asymmetric-lease-time "):
		block.AsymmetricLeaseTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "bootp-support":
		block.BootpSupport = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "client-discover-match "):
		block.ClientDiscoverMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "delay-offer based-on "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 4 { // <option> <compare> <value_type> <value>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "delay-offer based-on", itemTrim)
		}
		block.DelayOfferBasedOn = append(block.DelayOfferBasedOn,
			systemServicesDhcpLocalserverGroupBlockOverridesBlockDelayBasedOn{
				Option:    types.StringValue(itemTrimFields[0]),
				Compare:   types.StringValue(itemTrimFields[1]),
				ValueType: types.StringValue(itemTrimFields[2]),
				Value:     types.StringValue(html.UnescapeString(strings.Trim(strings.Join(itemTrimFields[3:], " "), "\""))),
			},
		)
	case balt.CutPrefixInString(&itemTrim, "delay-offer delay-time "):
		block.DelayOfferDelayTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "delete-binding-on-renegotiation":
		block.DeleteBindingOnRenegotiation = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "dual-stack "):
		block.DualStack = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "include-option-82 forcerenew":
		block.IncludeOption82Forcerenew = types.BoolValue(true)
	case itemTrim == "include-option-82 nak":
		block.IncludeOption82Nak = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "interface-client-limit "):
		block.InterfaceClientLimit, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "process-inform pool "):
		block.ProcessInform = types.BoolValue(true)
		block.ProcessInformPool = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "process-inform":
		block.ProcessInform = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "protocol-attributes "):
		block.ProtocolAttributes = types.StringValue(strings.Trim(itemTrim, "\""))
	}
	if err != nil {
		return err
	}

	return nil
}

func (block *systemServicesDhcpLocalserverGroupBlockOverridesV6) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "always-add-option-dns-server":
		block.AlwaysAddOptionDNSServer = types.BoolValue(true)
	case itemTrim == "always-process-option-request-option":
		block.AlwaysProcessOptionRequestOption = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "asymmetric-lease-time "):
		block.AsymmetricLeaseTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "asymmetric-prefix-lease-time "):
		block.AsymmetricPrefixLeaseTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "client-negotiation-match incoming-interface":
		block.ClientNegotiationMatchIncomingInterface = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "delay-advertise based-on "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 4 { // <option> <compare> <value_type> <value>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "delay-advertise based-on", itemTrim)
		}
		block.DelayAdvertiseBasedOn = append(block.DelayAdvertiseBasedOn,
			systemServicesDhcpLocalserverGroupBlockOverridesBlockDelayBasedOn{
				Option:    types.StringValue(itemTrimFields[0]),
				Compare:   types.StringValue(itemTrimFields[1]),
				ValueType: types.StringValue(itemTrimFields[2]),
				Value:     types.StringValue(html.UnescapeString(strings.Trim(strings.Join(itemTrimFields[3:], " "), "\""))),
			},
		)
	case balt.CutPrefixInString(&itemTrim, "delay-advertise delay-time "):
		block.DelayAdvertiseDelayTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "delegated-pool "):
		block.DelegatedPool = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "delete-binding-on-renegotiation":
		block.DeleteBindingOnRenegotiation = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "dual-stack "):
		block.DualStack = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "interface-client-limit "):
		block.InterfaceClientLimit, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "multi-address-embedded-option-response":
		block.MultiAddressEmbeddedOptionResponse = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "process-inform pool "):
		block.ProcessInform = types.BoolValue(true)
		block.ProcessInformPool = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "process-inform":
		block.ProcessInform = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "protocol-attributes "):
		block.ProtocolAttributes = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "rapid-commit":
		block.RapidCommit = types.BoolValue(true)
	case itemTrim == "top-level-status-code":
		block.TopLevelStatusCode = types.BoolValue(true)
	}
	if err != nil {
		return err
	}

	return nil
}

func (block *systemServicesDhcpLocalserverGroupBlockReconfigure) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "attempts "):
		block.Attempts, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "clear-on-abort":
		block.ClearOnAbort = types.BoolValue(true)
	case itemTrim == "support-option-pd-exclude":
		block.SupportOptionPdExclude = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "timeout "):
		block.Timeout, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "token "):
		block.Token = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "trigger radius-disconnect":
		block.TriggerRadiusDisconnect = types.BoolValue(true)
	}
	if err != nil {
		return err
	}

	return nil
}

func (rscData *systemServicesDhcpLocalserverGroupData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "system services dhcp-local-server "
	if rscData.Version.ValueString() == "v6" {
		delPrefix += "dhcpv6 "
	}

	configSet := []string{
		delPrefix + "group " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
