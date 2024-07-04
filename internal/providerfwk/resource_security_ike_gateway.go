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
	_ resource.Resource                   = &securityIkeGateway{}
	_ resource.ResourceWithConfigure      = &securityIkeGateway{}
	_ resource.ResourceWithValidateConfig = &securityIkeGateway{}
	_ resource.ResourceWithImportState    = &securityIkeGateway{}
	_ resource.ResourceWithUpgradeState   = &securityIkeGateway{}
)

type securityIkeGateway struct {
	client *junos.Client
}

func newSecurityIkeGatewayResource() resource.Resource {
	return &securityIkeGateway{}
}

func (rsc *securityIkeGateway) typeName() string {
	return providerName + "_security_ike_gateway"
}

func (rsc *securityIkeGateway) junosName() string {
	return "security ike gateway"
}

func (rsc *securityIkeGateway) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityIkeGateway) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityIkeGateway) Configure(
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

func (rsc *securityIkeGateway) Schema(
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
				Description: "Label for the remote (peer) gateway.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"external_interface": schema.StringAttribute{
				Required:    true,
				Description: "Interface for IKE negotiations.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
				},
			},
			"policy": schema.StringAttribute{
				Required:    true,
				Description: "Name of the IKE policy.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"address": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Addresses or hostnames of peer:1 primary, upto 4 backups.",
				Validators: []validator.List{
					listvalidator.SizeBetween(1, 5),
					listvalidator.ValueStringsAre(
						tfvalidator.StringIPAddress(),
					),
				},
			},
			"general_ike_id": schema.BoolAttribute{
				Optional:    true,
				Description: "Accept peer IKE-ID in general.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"local_address": schema.StringAttribute{
				Optional:    true,
				Description: "Local IP for IKE negotiations.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
			"no_nat_traversal": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable IPSec NAT traversal.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Description: "Negotiate using either IKE v1 or IKE v2 protocol.",
				Validators: []validator.String{
					stringvalidator.OneOf("v1-only", "v2-only"),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"dynamic_remote": schema.SingleNestedBlock{
				Description: "Declare site to site peer with dynamic IP address.",
				Attributes: map[string]schema.Attribute{
					"connections_limit": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum number of users connected to gateway.",
						Validators: []validator.Int64{
							int64validator.Between(1, 4294967295),
						},
					},
					"hostname": schema.StringAttribute{
						Optional:    true,
						Description: "Use a fully-qualified domain name.",
						Validators: []validator.String{
							tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
						},
					},
					"ike_user_type": schema.StringAttribute{
						Optional:    true,
						Description: "Type of the IKE ID.",
						Validators: []validator.String{
							stringvalidator.OneOf("shared-ike-id", "group-ike-id"),
						},
					},
					"inet": schema.StringAttribute{
						Optional:    true,
						Description: "Use an IPV4 address to identify the dynamic peer.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"inet6": schema.StringAttribute{
						Optional:    true,
						Description: "Use an IPV6 address to identify the dynamic peer.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"reject_duplicate_connection": schema.BoolAttribute{
						Optional:    true,
						Description: "Reject new connection from duplicate IKE-id.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"user_at_hostname": schema.StringAttribute{
						Optional:    true,
						Description: "Use an e-mail address.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"distinguished_name": schema.SingleNestedBlock{
						Description: "Declare distinguished-name configuration.",
						Attributes: map[string]schema.Attribute{
							"container": schema.StringAttribute{
								Optional:    true,
								Description: "Container string for a distinguished name.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"wildcard": schema.StringAttribute{
								Optional:    true,
								Description: "Wildcard string for a distinguished name.",
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
			"aaa": schema.SingleNestedBlock{
				Description: "Use extended authentication.",
				Attributes: map[string]schema.Attribute{
					"access_profile": schema.StringAttribute{
						Optional:    true,
						Description: "Access profile that contains authentication information.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 250),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"client_password": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "AAA client password.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 128),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"client_username": schema.StringAttribute{
						Optional:    true,
						Description: "AAA client username.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 128),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"dead_peer_detection": schema.SingleNestedBlock{
				Description: "Declare RFC-3706 DPD configuration.",
				Attributes: map[string]schema.Attribute{
					"interval": schema.Int64Attribute{
						Optional:    true,
						Description: "The interval at which to send DPD.",
						Validators: []validator.Int64{
							int64validator.Between(10, 60),
						},
					},
					"send_mode": schema.StringAttribute{
						Optional:    true,
						Description: "Specify how probes are sent.",
						Validators: []validator.String{
							stringvalidator.OneOf("always-send", "optimized", "probe-idle-tunnel"),
						},
					},
					"threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum number of DPD retransmissions.",
						Validators: []validator.Int64{
							int64validator.Between(1, 5),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"local_identity": schema.SingleNestedBlock{
				Description: "Set the local IKE identity.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Type of IKE identity.",
						Validators: []validator.String{
							stringvalidator.OneOf("distinguished-name", "hostname", "inet", "inet6", "user-at-hostname"),
						},
					},
					"value": schema.StringAttribute{
						Optional:    true,
						Description: "Value for IKE identity.",
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
			"remote_identity": schema.SingleNestedBlock{
				Description: "Set the remote IKE identity.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Type of IKE identity.",
						Validators: []validator.String{
							stringvalidator.OneOf("distinguished-name", "hostname", "inet", "inet6", "user-at-hostname"),
						},
					},
					"value": schema.StringAttribute{
						Optional:    true,
						Description: "Value for IKE identity.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"distinguished_name_container": schema.StringAttribute{
						Optional:    true,
						Description: "Container string for a distinguished name.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"distinguished_name_wildcard": schema.StringAttribute{
						Optional:    true,
						Description: "Wildcard string for a distinguished name.",
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
	}
}

type securityIkeGatewayData struct {
	ID                types.String                              `tfsdk:"id"`
	Name              types.String                              `tfsdk:"name"`
	ExternalInterface types.String                              `tfsdk:"external_interface"`
	Policy            types.String                              `tfsdk:"policy"`
	Address           []types.String                            `tfsdk:"address"`
	GeneralIkeID      types.Bool                                `tfsdk:"general_ike_id"`
	LocalAddress      types.String                              `tfsdk:"local_address"`
	NoNatTraversal    types.Bool                                `tfsdk:"no_nat_traversal"`
	Version           types.String                              `tfsdk:"version"`
	Aaa               *securityIkeGatewayBlockAaa               `tfsdk:"aaa"`
	DeadPeerDetection *securityIkeGatewayBlockDeadPeerDetection `tfsdk:"dead_peer_detection"`
	DynamicRemote     *securityIkeGatewayBlockDynamicRemote     `tfsdk:"dynamic_remote"`
	LocalIdentity     *securityIkeGatewayBlockLocalIdentity     `tfsdk:"local_identity"`
	RemoteIdentity    *securityIkeGatewayBlockRemoteIdentity    `tfsdk:"remote_identity"`
}

type securityIkeGatewayConfig struct {
	ID                types.String                              `tfsdk:"id"`
	Name              types.String                              `tfsdk:"name"`
	ExternalInterface types.String                              `tfsdk:"external_interface"`
	Policy            types.String                              `tfsdk:"policy"`
	Address           types.List                                `tfsdk:"address"`
	GeneralIkeID      types.Bool                                `tfsdk:"general_ike_id"`
	LocalAddress      types.String                              `tfsdk:"local_address"`
	NoNatTraversal    types.Bool                                `tfsdk:"no_nat_traversal"`
	Version           types.String                              `tfsdk:"version"`
	Aaa               *securityIkeGatewayBlockAaa               `tfsdk:"aaa"`
	DeadPeerDetection *securityIkeGatewayBlockDeadPeerDetection `tfsdk:"dead_peer_detection"`
	DynamicRemote     *securityIkeGatewayBlockDynamicRemote     `tfsdk:"dynamic_remote"`
	LocalIdentity     *securityIkeGatewayBlockLocalIdentity     `tfsdk:"local_identity"`
	RemoteIdentity    *securityIkeGatewayBlockRemoteIdentity    `tfsdk:"remote_identity"`
}

//nolint:lll
type securityIkeGatewayBlockDynamicRemote struct {
	ConnectionsLimit          types.Int64                                                 `tfsdk:"connections_limit"`
	Hostname                  types.String                                                `tfsdk:"hostname"`
	IkeUserType               types.String                                                `tfsdk:"ike_user_type"`
	Inet                      types.String                                                `tfsdk:"inet"`
	Inet6                     types.String                                                `tfsdk:"inet6"`
	RejectDuplicateConnection types.Bool                                                  `tfsdk:"reject_duplicate_connection"`
	UserAtHostname            types.String                                                `tfsdk:"user_at_hostname"`
	DistinguishedName         *securityIkeGatewayBlockDynamicRemoteBlockDistinguishedName `tfsdk:"distinguished_name"`
}

func (block *securityIkeGatewayBlockDynamicRemote) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityIkeGatewayBlockDynamicRemoteBlockDistinguishedName struct {
	Container types.String `tfsdk:"container"`
	Wildcard  types.String `tfsdk:"wildcard"`
}

func (block *securityIkeGatewayBlockDynamicRemoteBlockDistinguishedName) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityIkeGatewayBlockAaa struct {
	AccessProfile  types.String `tfsdk:"access_profile"`
	ClientPassword types.String `tfsdk:"client_password"`
	ClientUsername types.String `tfsdk:"client_username"`
}

type securityIkeGatewayBlockDeadPeerDetection struct {
	Interval  types.Int64  `tfsdk:"interval"`
	SendMode  types.String `tfsdk:"send_mode"`
	Threshold types.Int64  `tfsdk:"threshold"`
}

type securityIkeGatewayBlockLocalIdentity struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

type securityIkeGatewayBlockRemoteIdentity struct {
	Type                       types.String `tfsdk:"type"`
	Value                      types.String `tfsdk:"value"`
	DistinguishedNameContainer types.String `tfsdk:"distinguished_name_container"`
	DistinguishedNameWildcard  types.String `tfsdk:"distinguished_name_wildcard"`
}

func (rsc *securityIkeGateway) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityIkeGatewayConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Address.IsNull() && config.DynamicRemote == nil {
		resp.Diagnostics.AddError(
			tfdiag.MissingConfigErrSummary,
			"one of address or dynamic_remote must be specified",
		)
	}
	if !config.Address.IsNull() && !config.Address.IsUnknown() &&
		config.DynamicRemote != nil && config.DynamicRemote.hasKnownValue() {
		resp.Diagnostics.AddAttributeError(
			path.Root("address"),
			tfdiag.ConflictConfigErrSummary,
			"only one of address or dynamic_remote must be specified",
		)
	}
	if config.DynamicRemote != nil && config.DynamicRemote.hasKnownValue() &&
		!config.GeneralIkeID.IsNull() && !config.GeneralIkeID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("general_ike_id"),
			tfdiag.ConflictConfigErrSummary,
			"cannot set general_ike_id if dynamic_remote is used",
		)
	}
	if config.DynamicRemote != nil {
		switch {
		case config.DynamicRemote.DistinguishedName != nil && config.DynamicRemote.DistinguishedName.hasKnownValue():
			if (!config.DynamicRemote.Hostname.IsNull() && !config.DynamicRemote.Hostname.IsUnknown()) ||
				(!config.DynamicRemote.Inet.IsNull() && !config.DynamicRemote.Inet.IsUnknown()) ||
				(!config.DynamicRemote.Inet6.IsNull() && !config.DynamicRemote.Inet6.IsUnknown()) ||
				(!config.DynamicRemote.UserAtHostname.IsNull() && !config.DynamicRemote.UserAtHostname.IsUnknown()) {
				resp.Diagnostics.AddAttributeError(
					path.Root("dynamic_remote").AtName("distinguished_name"),
					tfdiag.ConflictConfigErrSummary,
					"only one of distinguished_name, hostname, inet, inet6 or user_at_hostname "+
						"can be specified in dynamic_remote block",
				)
			}
		case !config.DynamicRemote.Hostname.IsNull() && !config.DynamicRemote.Hostname.IsUnknown():
			if (config.DynamicRemote.DistinguishedName != nil && config.DynamicRemote.DistinguishedName.hasKnownValue()) ||
				(!config.DynamicRemote.Inet.IsNull() && !config.DynamicRemote.Inet.IsUnknown()) ||
				(!config.DynamicRemote.Inet6.IsNull() && !config.DynamicRemote.Inet6.IsUnknown()) ||
				(!config.DynamicRemote.UserAtHostname.IsNull() && !config.DynamicRemote.UserAtHostname.IsUnknown()) {
				resp.Diagnostics.AddAttributeError(
					path.Root("dynamic_remote").AtName("hostname"),
					tfdiag.ConflictConfigErrSummary,
					"only one of distinguished_name, hostname, inet, inet6 or user_at_hostname "+
						"can be specified in dynamic_remote block",
				)
			}
		case !config.DynamicRemote.Inet.IsNull() && !config.DynamicRemote.Inet.IsUnknown():
			if (config.DynamicRemote.DistinguishedName != nil && config.DynamicRemote.DistinguishedName.hasKnownValue()) ||
				(!config.DynamicRemote.Hostname.IsNull() && !config.DynamicRemote.Hostname.IsUnknown()) ||
				(!config.DynamicRemote.Inet6.IsNull() && !config.DynamicRemote.Inet6.IsUnknown()) ||
				(!config.DynamicRemote.UserAtHostname.IsNull() && !config.DynamicRemote.UserAtHostname.IsUnknown()) {
				resp.Diagnostics.AddAttributeError(
					path.Root("dynamic_remote").AtName("inet"),
					tfdiag.ConflictConfigErrSummary,
					"only one of distinguished_name, hostname, inet, inet6 or user_at_hostname "+
						"can be specified in dynamic_remote block",
				)
			}
		case !config.DynamicRemote.Inet6.IsNull() && !config.DynamicRemote.Inet6.IsUnknown():
			if (config.DynamicRemote.DistinguishedName != nil && config.DynamicRemote.DistinguishedName.hasKnownValue()) ||
				(!config.DynamicRemote.Hostname.IsNull() && !config.DynamicRemote.Hostname.IsUnknown()) ||
				(!config.DynamicRemote.Inet.IsNull() && !config.DynamicRemote.Inet.IsUnknown()) ||
				(!config.DynamicRemote.UserAtHostname.IsNull() && !config.DynamicRemote.UserAtHostname.IsUnknown()) {
				resp.Diagnostics.AddAttributeError(
					path.Root("dynamic_remote").AtName("inet6"),
					tfdiag.ConflictConfigErrSummary,
					"only one of distinguished_name, hostname, inet, inet6 or user_at_hostname "+
						"can be specified in dynamic_remote block",
				)
			}
		case !config.DynamicRemote.UserAtHostname.IsNull() && !config.DynamicRemote.UserAtHostname.IsUnknown():
			if (config.DynamicRemote.DistinguishedName != nil && config.DynamicRemote.DistinguishedName.hasKnownValue()) ||
				(!config.DynamicRemote.Hostname.IsNull() && !config.DynamicRemote.Hostname.IsUnknown()) ||
				(!config.DynamicRemote.Inet.IsNull() && !config.DynamicRemote.Inet.IsUnknown()) ||
				(!config.DynamicRemote.Inet6.IsNull() && !config.DynamicRemote.Inet6.IsUnknown()) {
				resp.Diagnostics.AddAttributeError(
					path.Root("dynamic_remote").AtName("user_at_hostname"),
					tfdiag.ConflictConfigErrSummary,
					"only one of distinguished_name, hostname, inet, inet6 or user_at_hostname "+
						"can be specified in dynamic_remote block",
				)
			}
		}
	}
	if config.Aaa != nil {
		if config.Aaa.AccessProfile.IsNull() && config.Aaa.ClientUsername.IsNull() && config.Aaa.ClientPassword.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("aaa").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"one of access_profile or client_username/client_password must be specified in aaa block",
			)
		}
		if !config.Aaa.AccessProfile.IsNull() && !config.Aaa.AccessProfile.IsUnknown() &&
			((!config.Aaa.ClientUsername.IsNull() && !config.Aaa.ClientUsername.IsUnknown()) ||
				(!config.Aaa.ClientPassword.IsNull() && !config.Aaa.ClientPassword.IsUnknown())) {
			resp.Diagnostics.AddAttributeError(
				path.Root("aaa").AtName("access_profile"),
				tfdiag.ConflictConfigErrSummary,
				"only one of access_profile or client_username/client_password must be specifiedin aaa block",
			)
		}
		if config.Aaa.ClientUsername.IsNull() && !config.Aaa.ClientPassword.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("aaa").AtName("client_password"),
				tfdiag.MissingConfigErrSummary,
				"client_username and client_password must be specified together in aaa block",
			)
		}
		if !config.Aaa.ClientUsername.IsNull() && config.Aaa.ClientPassword.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("aaa").AtName("client_username"),
				tfdiag.MissingConfigErrSummary,
				"client_username and client_password must be specified together in aaa block",
			)
		}
	}
	if config.LocalIdentity != nil {
		if config.LocalIdentity.Type.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("local_identity").AtName("type"),
				tfdiag.MissingConfigErrSummary,
				"type must be specified in local_identity block",
			)
		}
		if !config.LocalIdentity.Type.IsNull() && !config.LocalIdentity.Type.IsUnknown() {
			if v := config.LocalIdentity.Type.ValueString(); v == "distinguished-name" {
				if !config.LocalIdentity.Value.IsNull() && !config.LocalIdentity.Value.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("local_identity").AtName("value"),
						tfdiag.ConflictConfigErrSummary,
						"value should not be specified when type is set to distinguished-name in local_identity block",
					)
				}
			} else {
				if config.LocalIdentity.Value.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("local_identity").AtName("type"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("value must be specified when type is set to %q in local_identity block", v),
					)
				}
			}
		}
	}
	if config.RemoteIdentity != nil {
		if config.RemoteIdentity.Type.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("remote_identity").AtName("type"),
				tfdiag.MissingConfigErrSummary,
				"type must be specified in remote_identity block",
			)
		}
		if !config.RemoteIdentity.Type.IsNull() && !config.RemoteIdentity.Type.IsUnknown() {
			if v := config.RemoteIdentity.Type.ValueString(); v == "distinguished-name" {
				if !config.RemoteIdentity.Value.IsNull() && !config.RemoteIdentity.Value.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("remote_identity").AtName("value"),
						tfdiag.ConflictConfigErrSummary,
						"value should not be specified when type is set to distinguished-name in remote_identity block",
					)
				}
			} else {
				if config.RemoteIdentity.Value.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("remote_identity").AtName("type"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("value must be specified when type is set to %q in remote_identity block", v),
					)
				}
				if (!config.RemoteIdentity.DistinguishedNameContainer.IsNull() &&
					!config.RemoteIdentity.DistinguishedNameContainer.IsUnknown()) ||
					(!config.RemoteIdentity.DistinguishedNameWildcard.IsNull() &&
						!config.RemoteIdentity.DistinguishedNameWildcard.IsUnknown()) {
					resp.Diagnostics.AddAttributeError(
						path.Root("remote_identity").AtName("type"),
						tfdiag.ConflictConfigErrSummary,
						"type must be set to distinguished-name with "+
							"distinguished_name_container and distinguished_name_wildcard in remote_identity block",
					)
				}
			}
		}
	}
}

func (rsc *securityIkeGateway) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityIkeGatewayData
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
			gatewayExists, err := checkSecurityIkeGatewayExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if gatewayExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			gatewayExists, err := checkSecurityIkeGatewayExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !gatewayExists {
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

func (rsc *securityIkeGateway) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityIkeGatewayData
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

func (rsc *securityIkeGateway) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityIkeGatewayData
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

func (rsc *securityIkeGateway) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityIkeGatewayData
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

func (rsc *securityIkeGateway) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityIkeGatewayData

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

func checkSecurityIkeGatewayExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ike gateway \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityIkeGatewayData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityIkeGatewayData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityIkeGatewayData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security ike gateway \"" + rscData.Name.ValueString() + "\" "

	configSet = append(configSet, setPrefix+"ike-policy \""+rscData.Policy.ValueString()+"\"")
	configSet = append(configSet, setPrefix+"external-interface "+rscData.ExternalInterface.ValueString())
	for _, v := range rscData.Address {
		configSet = append(configSet, setPrefix+"address "+v.ValueString())
	}
	if rscData.DynamicRemote != nil {
		if !rscData.DynamicRemote.ConnectionsLimit.IsNull() {
			configSet = append(configSet, setPrefix+"dynamic connections-limit "+
				utils.ConvI64toa(rscData.DynamicRemote.ConnectionsLimit.ValueInt64()))
		}
		if rscData.DynamicRemote.DistinguishedName != nil {
			configSet = append(configSet, setPrefix+"dynamic distinguished-name")
			if v := rscData.DynamicRemote.DistinguishedName.Container.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"dynamic distinguished-name container \""+v+"\"")
			}
			if v := rscData.DynamicRemote.DistinguishedName.Wildcard.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"dynamic distinguished-name wildcard \""+v+"\"")
			}
		}
		if v := rscData.DynamicRemote.Hostname.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"dynamic hostname "+v)
		}
		if v := rscData.DynamicRemote.IkeUserType.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"dynamic ike-user-type "+v)
		}
		if v := rscData.DynamicRemote.Inet.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"dynamic inet "+v)
		}
		if v := rscData.DynamicRemote.Inet6.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"dynamic inet6 "+v)
		}
		if rscData.DynamicRemote.RejectDuplicateConnection.ValueBool() {
			configSet = append(configSet, setPrefix+"dynamic reject-duplicate-connection")
		}
		if v := rscData.DynamicRemote.UserAtHostname.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"dynamic user-at-hostname \""+v+"\"")
		}
	}
	if rscData.Aaa != nil {
		if v := rscData.Aaa.AccessProfile.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"aaa access-profile \""+v+"\"")
		}
		if v := rscData.Aaa.ClientPassword.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"aaa client password \""+v+"\"")
		}
		if v := rscData.Aaa.ClientUsername.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"aaa client username \""+v+"\"")
		}
	}
	if rscData.DeadPeerDetection != nil {
		configSet = append(configSet, setPrefix+"dead-peer-detection")

		if !rscData.DeadPeerDetection.Interval.IsNull() {
			configSet = append(configSet, setPrefix+"dead-peer-detection interval "+
				utils.ConvI64toa(rscData.DeadPeerDetection.Interval.ValueInt64()))
		}
		if v := rscData.DeadPeerDetection.SendMode.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"dead-peer-detection "+v)
		}
		if !rscData.DeadPeerDetection.Threshold.IsNull() {
			configSet = append(configSet, setPrefix+"dead-peer-detection threshold "+
				utils.ConvI64toa(rscData.DeadPeerDetection.Threshold.ValueInt64()))
		}
	}
	if rscData.GeneralIkeID.ValueBool() {
		configSet = append(configSet, setPrefix+"general-ikeid")
	}
	if v := rscData.LocalAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"local-address "+v)
	}
	if rscData.LocalIdentity != nil {
		if v1 := rscData.LocalIdentity.Type.ValueString(); v1 == "distinguished-name" {
			if !rscData.LocalIdentity.Value.IsNull() {
				return path.Root("local_identity").AtName("value"),
					errors.New("conflict: value should not be specified " +
						"when type is set to distinguished-name in local_identity block")
			}
			configSet = append(configSet, setPrefix+"local-identity "+v1)
		} else {
			if v2 := rscData.LocalIdentity.Value.ValueString(); v2 != "" {
				configSet = append(configSet, setPrefix+"local-identity "+v1+" \""+v2+"\"")
			} else {
				return path.Root("local_identity").AtName("type"),
					fmt.Errorf("missing: value must be specified "+
						"when type is set to %q in local_identity block", v1)
			}
		}
	}
	if rscData.NoNatTraversal.ValueBool() {
		configSet = append(configSet, setPrefix+"no-nat-traversal")
	}
	if rscData.RemoteIdentity != nil {
		if v1 := rscData.RemoteIdentity.Type.ValueString(); v1 == "distinguished-name" {
			configSet = append(configSet, setPrefix+"remote-identity "+v1)
			if v2 := rscData.RemoteIdentity.DistinguishedNameContainer.ValueString(); v2 != "" {
				configSet = append(configSet, setPrefix+"remote-identity "+v1+" container \""+v2+"\"")
			}
			if v2 := rscData.RemoteIdentity.DistinguishedNameWildcard.ValueString(); v2 != "" {
				configSet = append(configSet, setPrefix+"remote-identity "+v1+" wildcard \""+v2+"\"")
			}
		} else {
			if !rscData.RemoteIdentity.DistinguishedNameContainer.IsNull() ||
				!rscData.RemoteIdentity.DistinguishedNameWildcard.IsNull() {
				return path.Root("remote_identity").AtName("type"),
					errors.New("conflict: type must be set to distinguished-name " +
						"with distinguished_name_container and distinguished_name_wildcard in remote_identity block")
			}
			if v2 := rscData.RemoteIdentity.Value.ValueString(); v2 != "" {
				configSet = append(configSet, setPrefix+"remote-identity "+v1+" \""+v2+"\"")
			} else {
				return path.Root("remote_identity").AtName("type"),
					fmt.Errorf("missing: value must be specified "+
						"when type is set to %q in remote_identity block", v1)
			}
		}
	}
	if v := rscData.Version.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"version "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityIkeGatewayData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ike gateway \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "external-interface "):
				rscData.ExternalInterface = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "ike-policy "):
				rscData.Policy = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "address "):
				rscData.Address = append(rscData.Address, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "dynamic "):
				if rscData.DynamicRemote == nil {
					rscData.DynamicRemote = &securityIkeGatewayBlockDynamicRemote{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "connections-limit "):
					rscData.DynamicRemote.ConnectionsLimit, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "distinguished-name"):
					if rscData.DynamicRemote.DistinguishedName == nil {
						rscData.DynamicRemote.DistinguishedName = &securityIkeGatewayBlockDynamicRemoteBlockDistinguishedName{}
					}
					switch {
					case balt.CutPrefixInString(&itemTrim, " container "):
						rscData.DynamicRemote.DistinguishedName.Container = types.StringValue(strings.Trim(itemTrim, "\""))
					case balt.CutPrefixInString(&itemTrim, " wildcard "):
						rscData.DynamicRemote.DistinguishedName.Wildcard = types.StringValue(strings.Trim(itemTrim, "\""))
					}
				case balt.CutPrefixInString(&itemTrim, "hostname "):
					rscData.DynamicRemote.Hostname = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "ike-user-type "):
					rscData.DynamicRemote.IkeUserType = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "inet "):
					rscData.DynamicRemote.Inet = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "inet6 "):
					rscData.DynamicRemote.Inet6 = types.StringValue(itemTrim)
				case itemTrim == "reject-duplicate-connection":
					rscData.DynamicRemote.RejectDuplicateConnection = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "user-at-hostname "):
					rscData.DynamicRemote.UserAtHostname = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			case balt.CutPrefixInString(&itemTrim, "aaa "):
				if rscData.Aaa == nil {
					rscData.Aaa = &securityIkeGatewayBlockAaa{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "access-profile "):
					rscData.Aaa.AccessProfile = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "client password "):
					rscData.Aaa.ClientPassword, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "aaa client password")
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "client username "):
					rscData.Aaa.ClientUsername = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			case balt.CutPrefixInString(&itemTrim, "dead-peer-detection"):
				if rscData.DeadPeerDetection == nil {
					rscData.DeadPeerDetection = &securityIkeGatewayBlockDeadPeerDetection{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, " interval "):
					rscData.DeadPeerDetection.Interval, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case itemTrim == " always-send":
					rscData.DeadPeerDetection.SendMode = types.StringValue("always-send")
				case itemTrim == " optimized":
					rscData.DeadPeerDetection.SendMode = types.StringValue("optimized")
				case itemTrim == " probe-idle-tunnel":
					rscData.DeadPeerDetection.SendMode = types.StringValue("probe-idle-tunnel")
				case balt.CutPrefixInString(&itemTrim, " threshold "):
					rscData.DeadPeerDetection.Threshold, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			case itemTrim == "general-ikeid":
				rscData.GeneralIkeID = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "local-address "):
				rscData.LocalAddress = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "local-identity "):
				if rscData.LocalIdentity == nil {
					rscData.LocalIdentity = &securityIkeGatewayBlockLocalIdentity{}
				}
				itemTrimFields := strings.Split(itemTrim, " ")
				rscData.LocalIdentity.Type = types.StringValue(itemTrimFields[0])
				if len(itemTrimFields) > 1 {
					rscData.LocalIdentity.Value = types.StringValue(
						strings.Trim(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" "), "\""))
				}
			case itemTrim == "no-nat-traversal":
				rscData.NoNatTraversal = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "remote-identity "):
				if rscData.RemoteIdentity == nil {
					rscData.RemoteIdentity = &securityIkeGatewayBlockRemoteIdentity{}
				}
				itemTrimFields := strings.Split(itemTrim, " ")
				rscData.RemoteIdentity.Type = types.StringValue(itemTrimFields[0])
				if len(itemTrimFields) > 1 {
					if rscData.RemoteIdentity.Type.ValueString() == "distinguished-name" {
						if itemTrimFields[1] == "container" {
							rscData.RemoteIdentity.DistinguishedNameContainer = types.StringValue(
								strings.Trim(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" "+itemTrimFields[1]+" "), "\""))
						}
						if itemTrimFields[1] == "wildcard" {
							rscData.RemoteIdentity.DistinguishedNameWildcard = types.StringValue(
								strings.Trim(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" "+itemTrimFields[1]+" "), "\""))
						}
					} else {
						rscData.RemoteIdentity.Value = types.StringValue(
							strings.Trim(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" "), "\""))
					}
				}
			case balt.CutPrefixInString(&itemTrim, "version "):
				rscData.Version = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (rscData *securityIkeGatewayData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security ike gateway \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
