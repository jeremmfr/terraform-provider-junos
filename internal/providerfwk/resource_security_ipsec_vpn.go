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
	_ resource.Resource                   = &securityIpsecVpn{}
	_ resource.ResourceWithConfigure      = &securityIpsecVpn{}
	_ resource.ResourceWithModifyPlan     = &securityIpsecVpn{}
	_ resource.ResourceWithValidateConfig = &securityIpsecVpn{}
	_ resource.ResourceWithImportState    = &securityIpsecVpn{}
	_ resource.ResourceWithUpgradeState   = &securityIpsecVpn{}
)

type securityIpsecVpn struct {
	client *junos.Client
}

func newSecurityIpsecVpnResource() resource.Resource {
	return &securityIpsecVpn{}
}

func (rsc *securityIpsecVpn) typeName() string {
	return providerName + "_security_ipsec_vpn"
}

func (rsc *securityIpsecVpn) junosName() string {
	return "security ipsec vpn"
}

func (rsc *securityIpsecVpn) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityIpsecVpn) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityIpsecVpn) Configure(
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

func (rsc *securityIpsecVpn) Schema(
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
				Description: "The name of vpn.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"bind_interface": schema.StringAttribute{
				Optional:    true,
				Description: "Interface to bind vpn for route-based vpn.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^st0\.`),
						"must be a secure tunnel interface"),
				},
			},
			"copy_outer_dscp": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable copying outer IP header DSCP and ECN to inner IP header.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"df_bit": schema.StringAttribute{
				Optional:    true,
				Description: "Specifies how to handle the Don't Fragment bit.",
				Validators: []validator.String{
					stringvalidator.OneOf("clear", "copy", "set"),
				},
			},
			"establish_tunnels": schema.StringAttribute{
				Optional:    true,
				Description: "When the VPN comes up.",
				Validators: []validator.String{
					stringvalidator.OneOf("immediately", "on-traffic"),
				},
			},
			"multi_sa_forwarding_class": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Negotiate multiple SAs with forwarding-classes.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 32),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"ike": schema.SingleNestedBlock{
				Description: "Declare IKE-keyed configuration.",
				Attributes: map[string]schema.Attribute{
					"gateway": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "The name of security IKE gateway (phase-1).",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 32),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"policy": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "The name of IPSec policy.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 32),
							tfvalidator.StringFormat(tfvalidator.DefaultFormat),
						},
					},
					"identity_local": schema.StringAttribute{
						Optional:    true,
						Description: "IPSec proxy-id local parameter.",
						Validators: []validator.String{
							tfvalidator.StringCIDR(),
						},
					},
					"identity_remote": schema.StringAttribute{
						Optional:    true,
						Description: "IPSec proxy-id remote parameter.",
						Validators: []validator.String{
							tfvalidator.StringCIDR(),
						},
					},
					"identity_service": schema.StringAttribute{
						Optional:    true,
						Description: "IPSec proxy-id service parameter.",
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
			"manual": schema.SingleNestedBlock{
				Description: "Define a manual security association.",
				Attributes: map[string]schema.Attribute{
					"external_interface": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "External interface for the security association.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
						},
					},
					"protocol": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Define an IPSec protocol for the security association.",
						Validators: []validator.String{
							stringvalidator.OneOf("ah", "esp"),
						},
					},
					"spi": schema.Int64Attribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Define security parameter index (256..16639).",
						Validators: []validator.Int64{
							int64validator.Between(256, 16639),
						},
					},
					"authentication_algorithm": schema.StringAttribute{
						Optional:    true,
						Description: "Define authentication algorithm.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringSpaceExclusion(),
						},
					},
					"authentication_key_hexa": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Define an authentication key with format as hexadecimal.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringFormat(tfvalidator.HexadecimalFormat).WithSensitiveData(),
						},
					},
					"authentication_key_text": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Define an authentication key with format as text.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"encryption_algorithm": schema.StringAttribute{
						Optional:    true,
						Description: "Define encryption algorithm.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringSpaceExclusion(),
						},
					},
					"encryption_key_hexa": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Define an encryption key with format as hexadecimal.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringFormat(tfvalidator.HexadecimalFormat).WithSensitiveData(),
						},
					},
					"encryption_key_text": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Define an encryption key with format as text.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"gateway": schema.StringAttribute{
						Optional:    true,
						Description: "Define the IPSec peer.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"traffic_selector": schema.ListNestedBlock{
				Description: "For each name of traffic-selector to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of traffic-selector.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 31),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"local_ip": schema.StringAttribute{
							Required:    true,
							Description: "CIDR for IP addresses of local traffic-selector.",
							Validators: []validator.String{
								tfvalidator.StringCIDR(),
							},
						},
						"remote_ip": schema.StringAttribute{
							Required:    true,
							Description: "CIDR for IP addresses of remote traffic-selector.",
							Validators: []validator.String{
								tfvalidator.StringCIDR(),
							},
						},
					},
				},
			},
			"udp_encapsulate": schema.SingleNestedBlock{
				Description: "UDP encapsulation of IPsec data traffic.",
				Attributes: map[string]schema.Attribute{
					"dest_port": schema.Int64Attribute{
						Optional:    true,
						Description: "UDP destination port (1025..65536)",
						Validators: []validator.Int64{
							int64validator.Between(1025, 65536),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"vpn_monitor": schema.SingleNestedBlock{
				Description: "Declare VPN monitor liveness configuration.",
				Attributes: map[string]schema.Attribute{
					"destination_ip": schema.StringAttribute{
						Optional:    true,
						Description: "IP destination for monitor message.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"optimized": schema.BoolAttribute{
						Optional:    true,
						Description: "Optimize for scalability.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"source_interface": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Set source interface for monitor message.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
						},
					},
					"source_interface_auto": schema.BoolAttribute{
						Optional:    true,
						Description: "Compute the source_interface to 'bind_interface'.",
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

type securityIpsecVpnData struct {
	ID                     types.String                           `tfsdk:"id"`
	Name                   types.String                           `tfsdk:"name"`
	BindInterface          types.String                           `tfsdk:"bind_interface"`
	CopyOuterDscp          types.Bool                             `tfsdk:"copy_outer_dscp"`
	DfBit                  types.String                           `tfsdk:"df_bit"`
	EstablishTunnels       types.String                           `tfsdk:"establish_tunnels"`
	Ike                    *securityIpsecVpnBlockIke              `tfsdk:"ike"`
	Manual                 *securityIpsecVpnBlockManual           `tfsdk:"manual"`
	MultiSaForwardingClass []types.String                         `tfsdk:"multi_sa_forwarding_class"`
	TrafficSelector        []securityIpsecVpnBlockTrafficSelector `tfsdk:"traffic_selector"`
	UDPEncapsulate         *securityIpsecVpnBlockUDPEncapsulate   `tfsdk:"udp_encapsulate"`
	VpnMonitor             *securityIpsecVpnBlockVpnMonitor       `tfsdk:"vpn_monitor"`
}

type securityIpsecVpnConfig struct {
	ID                     types.String                         `tfsdk:"id"`
	Name                   types.String                         `tfsdk:"name"`
	BindInterface          types.String                         `tfsdk:"bind_interface"`
	CopyOuterDscp          types.Bool                           `tfsdk:"copy_outer_dscp"`
	DfBit                  types.String                         `tfsdk:"df_bit"`
	EstablishTunnels       types.String                         `tfsdk:"establish_tunnels"`
	Ike                    *securityIpsecVpnBlockIke            `tfsdk:"ike"`
	Manual                 *securityIpsecVpnBlockManual         `tfsdk:"manual"`
	MultiSaForwardingClass types.Set                            `tfsdk:"multi_sa_forwarding_class"`
	TrafficSelector        types.List                           `tfsdk:"traffic_selector"`
	UDPEncapsulate         *securityIpsecVpnBlockUDPEncapsulate `tfsdk:"udp_encapsulate"`
	VpnMonitor             *securityIpsecVpnBlockVpnMonitor     `tfsdk:"vpn_monitor"`
}

type securityIpsecVpnBlockIke struct {
	Gateway         types.String `tfsdk:"gateway"`
	Policy          types.String `tfsdk:"policy"`
	IdentityLocal   types.String `tfsdk:"identity_local"`
	IdentityRemote  types.String `tfsdk:"identity_remote"`
	IdentityService types.String `tfsdk:"identity_service"`
}

func (block *securityIpsecVpnBlockIke) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityIpsecVpnBlockManual struct {
	ExternalInterface       types.String `tfsdk:"external_interface"`
	Protocol                types.String `tfsdk:"protocol"`
	Spi                     types.Int64  `tfsdk:"spi"`
	AuthenticationAlgorithm types.String `tfsdk:"authentication_algorithm"`
	AuthenticationKeyHexa   types.String `tfsdk:"authentication_key_hexa"`
	AuthenticationKeyText   types.String `tfsdk:"authentication_key_text"`
	EncryptionAlgorithm     types.String `tfsdk:"encryption_algorithm"`
	EncryptionKeyHexa       types.String `tfsdk:"encryption_key_hexa"`
	EncryptionKeyText       types.String `tfsdk:"encryption_key_text"`
	Gateway                 types.String `tfsdk:"gateway"`
}

func (block *securityIpsecVpnBlockManual) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityIpsecVpnBlockTrafficSelector struct {
	Name     types.String `tfsdk:"name"`
	LocalIP  types.String `tfsdk:"local_ip"`
	RemoteIP types.String `tfsdk:"remote_ip"`
}

type securityIpsecVpnBlockUDPEncapsulate struct {
	DestPort types.Int64 `tfsdk:"dest_port"`
}

type securityIpsecVpnBlockVpnMonitor struct {
	DestinationIP       types.String `tfsdk:"destination_ip"`
	Optimized           types.Bool   `tfsdk:"optimized"`
	SourceInterface     types.String `tfsdk:"source_interface"`
	SourceInterfaceAuto types.Bool   `tfsdk:"source_interface_auto"`
}

func (block *securityIpsecVpnBlockVpnMonitor) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (rsc *securityIpsecVpn) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityIpsecVpnConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Ike == nil && config.Manual == nil {
		resp.Diagnostics.AddError(
			tfdiag.MissingConfigErrSummary,
			"one of ike or manual must be specified",
		)
	}
	if config.Ike != nil && config.Ike.hasKnownValue() &&
		config.Manual != nil && config.Manual.hasKnownValue() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ike").AtName("*"),
			tfdiag.ConflictConfigErrSummary,
			"only one of ike or manual must be specified",
		)
	}
	if config.Ike != nil {
		if config.Ike.Gateway.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ike").AtName("gateway"),
				tfdiag.MissingConfigErrSummary,
				"gateway must be specified in ike block",
			)
		}
		if config.Ike.Policy.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ike").AtName("policy"),
				tfdiag.MissingConfigErrSummary,
				"policy must be specified in ike block",
			)
		}
	}
	if config.Manual != nil {
		if config.Manual.hasKnownValue() &&
			!config.EstablishTunnels.IsNull() && !config.EstablishTunnels.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("establish_tunnels"),
				tfdiag.ConflictConfigErrSummary,
				"cannot set establish_tunnels if manual is used",
			)
		}
		if config.Manual.ExternalInterface.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("manual").AtName("external_interface"),
				tfdiag.MissingConfigErrSummary,
				"external_interface must be specified in manual block",
			)
		}
		if config.Manual.Protocol.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("manual").AtName("protocol"),
				tfdiag.MissingConfigErrSummary,
				"protocol must be specified in manual block",
			)
		} else if !config.Manual.Protocol.IsUnknown() {
			if v := config.Manual.Protocol.ValueString(); v == "ah" {
				if config.Manual.AuthenticationAlgorithm.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("manual").AtName("protocol"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("authentication_algorithm must be specified "+
							"with protocol set to %q in manual block", v),
					)
				}
			} else if v == "esp" {
				if config.Manual.AuthenticationAlgorithm.IsNull() &&
					config.Manual.EncryptionAlgorithm.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("manual").AtName("protocol"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("at least one of authentication_algorithm or encryption_algorithm must be specified "+
							"with protocol set to %q in manual block", v),
					)
				}
			}
		}
		if config.Manual.Spi.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("manual").AtName("spi"),
				tfdiag.MissingConfigErrSummary,
				"spi must be specified in manual block",
			)
		}
		if !config.Manual.AuthenticationAlgorithm.IsNull() {
			if config.Manual.AuthenticationKeyHexa.IsNull() &&
				config.Manual.AuthenticationKeyText.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("manual").AtName("authentication_algorithm"),
					tfdiag.MissingConfigErrSummary,
					"one of authentication_key_hexa or authentication_key_text must be specified "+
						"when authentication_algorithm is specified in manual block",
				)
			}
		}
		if !config.Manual.AuthenticationKeyHexa.IsNull() && !config.Manual.AuthenticationKeyHexa.IsUnknown() &&
			!config.Manual.AuthenticationKeyText.IsNull() && !config.Manual.AuthenticationKeyText.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("manual").AtName("authentication_key_text"),
				tfdiag.ConflictConfigErrSummary,
				"only one of authentication_key_hexa or authentication_key_text can be specified in manual block",
			)
		}
		if !config.Manual.EncryptionAlgorithm.IsNull() {
			if config.Manual.EncryptionKeyHexa.IsNull() &&
				config.Manual.EncryptionKeyText.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("manual").AtName("encryption_algorithm"),
					tfdiag.MissingConfigErrSummary,
					"one of encryption_key_hexa or encryption_key_text must be specified "+
						"when encryption_algorithm is specified in manual block",
				)
			}
		}
		if !config.Manual.EncryptionKeyHexa.IsNull() && !config.Manual.EncryptionKeyHexa.IsUnknown() &&
			!config.Manual.EncryptionKeyText.IsNull() && !config.Manual.EncryptionKeyText.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("manual").AtName("encryption_key_text"),
				tfdiag.ConflictConfigErrSummary,
				"only one of encryption_key_hexa or encryption_key_text can be specified in manual block",
			)
		}
	}
	if !config.TrafficSelector.IsNull() && !config.TrafficSelector.IsUnknown() {
		if config.Ike != nil {
			if !config.Ike.IdentityLocal.IsNull() && !config.Ike.IdentityLocal.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ike").AtName("identity_local"),
					tfdiag.ConflictConfigErrSummary,
					"ike.identity_local should not be specified when traffic_selector is used",
				)
			}
			if !config.Ike.IdentityRemote.IsNull() && !config.Ike.IdentityRemote.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ike").AtName("identity_remote"),
					tfdiag.ConflictConfigErrSummary,
					"ike.identity_remote should not be specified when traffic_selector is used",
				)
			}
			if !config.Ike.IdentityService.IsNull() && !config.Ike.IdentityService.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ike").AtName("identity_service"),
					tfdiag.ConflictConfigErrSummary,
					"ike.identity_service should not be specified when traffic_selector is used",
				)
			}
		}
		if config.VpnMonitor != nil && config.VpnMonitor.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("vpn_monitor").AtName("*"),
				tfdiag.ConflictConfigErrSummary,
				"vpn_monitor should not be specified when traffic_selector is used",
			)
		}

		var configTrafficSelector []securityIpsecVpnBlockTrafficSelector
		asDiags := config.TrafficSelector.ElementsAs(ctx, &configTrafficSelector, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		names := make(map[string]struct{})
		for i, block := range configTrafficSelector {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := names[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("traffic_selector").AtListIndex(i).AtName("name"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple traffic_selector blocks with the same name %q", name),
				)
			}
			names[name] = struct{}{}
		}
	}
	if !config.MultiSaForwardingClass.IsNull() && !config.MultiSaForwardingClass.IsUnknown() &&
		config.VpnMonitor != nil && config.VpnMonitor.hasKnownValue() {
		resp.Diagnostics.AddAttributeError(
			path.Root("vpn_monitor").AtName("*"),
			tfdiag.ConflictConfigErrSummary,
			"vpn_monitor should not be specified when multi_sa_forwarding_class is specified",
		)
	}
	if config.VpnMonitor != nil {
		if !config.VpnMonitor.SourceInterface.IsNull() && !config.VpnMonitor.SourceInterface.IsUnknown() &&
			!config.VpnMonitor.SourceInterfaceAuto.IsNull() && !config.VpnMonitor.SourceInterfaceAuto.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("vpn_monitor").AtName("source_interface_auto"),
				tfdiag.ConflictConfigErrSummary,
				"source_interface_auto should not be specified when source_interface is specified",
			)
		}
	}
}

func (rsc *securityIpsecVpn) ModifyPlan(
	ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var config, plan securityIpsecVpnConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.VpnMonitor != nil && plan.VpnMonitor != nil {
		if config.VpnMonitor.SourceInterface.IsNull() {
			if !config.VpnMonitor.SourceInterfaceAuto.IsNull() {
				if plan.BindInterface.IsNull() {
					plan.VpnMonitor.SourceInterface = types.StringNull()
				} else if !plan.BindInterface.IsUnknown() &&
					plan.VpnMonitor.SourceInterfaceAuto.ValueBool() {
					plan.VpnMonitor.SourceInterface = types.StringValue(plan.BindInterface.ValueString())
				}
			} else {
				plan.VpnMonitor.SourceInterface = types.StringNull()
			}
		}
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func (rsc *securityIpsecVpn) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityIpsecVpnData
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
	if plan.VpnMonitor != nil {
		if plan.VpnMonitor.SourceInterface.IsUnknown() {
			if plan.VpnMonitor.SourceInterfaceAuto.ValueBool() {
				if plan.BindInterface.IsNull() {
					plan.VpnMonitor.SourceInterface = types.StringNull()
				} else {
					plan.VpnMonitor.SourceInterface = types.StringValue(plan.BindInterface.ValueString())
				}
			} else {
				plan.VpnMonitor.SourceInterface = types.StringNull()
			}
		}
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
			vpnExists, err := checkSecurityIpsecVpnExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if vpnExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			vpnExists, err := checkSecurityIpsecVpnExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !vpnExists {
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

func (rsc *securityIpsecVpn) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityIpsecVpnData
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
		func() {
			if data.VpnMonitor != nil && state.VpnMonitor != nil {
				data.VpnMonitor.SourceInterfaceAuto = state.VpnMonitor.SourceInterfaceAuto
			}
		},
		resp,
	)
}

func (rsc *securityIpsecVpn) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityIpsecVpnData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.VpnMonitor != nil {
		if plan.VpnMonitor.SourceInterface.IsUnknown() {
			if plan.VpnMonitor.SourceInterfaceAuto.ValueBool() {
				if plan.BindInterface.IsNull() {
					plan.VpnMonitor.SourceInterface = types.StringNull()
				} else {
					plan.VpnMonitor.SourceInterface = types.StringValue(plan.BindInterface.ValueString())
				}
			} else {
				plan.VpnMonitor.SourceInterface = types.StringNull()
			}
		}
	}

	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *securityIpsecVpn) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityIpsecVpnData
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

func (rsc *securityIpsecVpn) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityIpsecVpnData

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

func checkSecurityIpsecVpnExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ipsec vpn \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityIpsecVpnData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityIpsecVpnData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityIpsecVpnData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security ipsec vpn \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.BindInterface.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"bind-interface "+v)
	}
	if rscData.CopyOuterDscp.ValueBool() {
		configSet = append(configSet, setPrefix+"copy-outer-dscp")
	}
	if v := rscData.DfBit.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"df-bit "+v)
	}
	if v := rscData.EstablishTunnels.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"establish-tunnels "+v)
	}
	if rscData.Ike != nil {
		if v := rscData.Ike.Gateway.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"ike gateway \""+v+"\"")
		} else {
			return path.Root("ike").AtName("gateway"), errors.New("missing: gateway must be not empty in ike block")
		}
		if v := rscData.Ike.Policy.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"ike ipsec-policy "+v)
		} else {
			return path.Root("ike").AtName("policy"), errors.New("missing: policy must be not empty in ike block")
		}
		if v := rscData.Ike.IdentityLocal.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"ike proxy-identity local "+v)
		}
		if v := rscData.Ike.IdentityRemote.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"ike proxy-identity remote "+v)
		}
		if v := rscData.Ike.IdentityService.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"ike proxy-identity service \""+v+"\"")
		}
	}
	if rscData.Manual != nil {
		if v := rscData.Manual.ExternalInterface.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"manual external-interface "+v)
		} else {
			return path.Root("manual").AtName("external_interface"),
				errors.New("missing: external_interface must be not empty in manual block")
		}
		if v := rscData.Manual.Protocol.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"manual protocol "+v)
		} else {
			return path.Root("manual").AtName("protocol"),
				errors.New("missing: protocol must be not empty in manual block")
		}
		configSet = append(configSet, setPrefix+"manual spi "+utils.ConvI64toa(rscData.Manual.Spi.ValueInt64()))
		if v := rscData.Manual.AuthenticationAlgorithm.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"manual authentication algorithm "+v)
		}
		if v := rscData.Manual.AuthenticationKeyHexa.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"manual authentication key hexadecimal "+v)
		}
		if v := rscData.Manual.AuthenticationKeyText.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"manual authentication key ascii-text \""+v+"\"")
		}
		if v := rscData.Manual.EncryptionAlgorithm.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"manual encryption algorithm "+v)
		}
		if v := rscData.Manual.EncryptionKeyHexa.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"manual encryption key hexadecimal "+v)
		}
		if v := rscData.Manual.EncryptionKeyText.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"manual encryption key ascii-text \""+v+"\"")
		}
		if v := rscData.Manual.Gateway.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"manual gateway "+v)
		}
	}
	trafficSelectorNames := make(map[string]struct{})
	for i, block := range rscData.TrafficSelector {
		name := block.Name.ValueString()
		if _, ok := trafficSelectorNames[name]; ok {
			return path.Root("traffic_selector").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple traffic_selector blocks with the same name %q", name)
		}
		trafficSelectorNames[name] = struct{}{}

		configSet = append(configSet,
			setPrefix+"traffic-selector \""+block.Name.ValueString()+"\" local-ip "+block.LocalIP.ValueString())
		configSet = append(configSet,
			setPrefix+"traffic-selector \""+block.Name.ValueString()+"\" remote-ip "+block.RemoteIP.ValueString())
	}
	for _, v := range rscData.MultiSaForwardingClass {
		configSet = append(configSet, setPrefix+"multi-sa forwarding-class \""+v.ValueString()+"\"")
	}
	if rscData.UDPEncapsulate != nil {
		configSet = append(configSet, setPrefix+"udp-encapsulate")
		if !rscData.UDPEncapsulate.DestPort.IsNull() {
			configSet = append(configSet, setPrefix+"udp-encapsulate dest-port "+
				utils.ConvI64toa(rscData.UDPEncapsulate.DestPort.ValueInt64()))
		}
	}
	if rscData.VpnMonitor != nil {
		configSet = append(configSet, setPrefix+"vpn-monitor")
		if v := rscData.VpnMonitor.DestinationIP.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"vpn-monitor destination-ip "+v)
		}
		if rscData.VpnMonitor.Optimized.ValueBool() {
			configSet = append(configSet, setPrefix+"vpn-monitor optimized")
		}
		if v := rscData.VpnMonitor.SourceInterface.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"vpn-monitor source-interface "+v)
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityIpsecVpnData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ipsec vpn \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "bind-interface "):
				rscData.BindInterface = types.StringValue(itemTrim)
			case itemTrim == "copy-outer-dscp":
				rscData.CopyOuterDscp = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "df-bit "):
				rscData.DfBit = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "establish-tunnels "):
				rscData.EstablishTunnels = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "ike "):
				if rscData.Ike == nil {
					rscData.Ike = &securityIpsecVpnBlockIke{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "gateway "):
					rscData.Ike.Gateway = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "ipsec-policy "):
					rscData.Ike.Policy = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "proxy-identity local "):
					rscData.Ike.IdentityLocal = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "proxy-identity remote "):
					rscData.Ike.IdentityRemote = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "proxy-identity service "):
					rscData.Ike.IdentityService = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			case balt.CutPrefixInString(&itemTrim, "manual "):
				if rscData.Manual == nil {
					rscData.Manual = &securityIpsecVpnBlockManual{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "external-interface "):
					rscData.Manual.ExternalInterface = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "protocol "):
					rscData.Manual.Protocol = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "spi "):
					rscData.Manual.Spi, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "authentication algorithm "):
					rscData.Manual.AuthenticationAlgorithm = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "authentication key hexadecimal "):
					rscData.Manual.AuthenticationKeyHexa, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""),
						"authentication key hexadecimal")
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "authentication key ascii-text "):
					rscData.Manual.AuthenticationKeyText, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""),
						"authentication key ascii-text")
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "encryption algorithm "):
					rscData.Manual.EncryptionAlgorithm = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "encryption key hexadecimal "):
					rscData.Manual.EncryptionKeyHexa, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""),
						"encryption key hexadecimal")
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "encryption key ascii-text "):
					rscData.Manual.EncryptionKeyText, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""),
						"encryption key ascii-text")
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "gateway "):
					rscData.Manual.Gateway = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			case balt.CutPrefixInString(&itemTrim, "multi-sa forwarding-class "):
				rscData.MultiSaForwardingClass = append(rscData.MultiSaForwardingClass,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "udp-encapsulate"):
				if rscData.UDPEncapsulate == nil {
					rscData.UDPEncapsulate = &securityIpsecVpnBlockUDPEncapsulate{}
				}
				if balt.CutPrefixInString(&itemTrim, " dest-port ") {
					rscData.UDPEncapsulate.DestPort, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "traffic-selector "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				var trafficSelector securityIpsecVpnBlockTrafficSelector
				rscData.TrafficSelector, trafficSelector = tfdata.ExtractBlockWithTFTypesString(
					rscData.TrafficSelector, "Name", strings.Trim(name, "\""))
				trafficSelector.Name = types.StringValue(strings.Trim(name, "\""))
				balt.CutPrefixInString(&itemTrim, name+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "local-ip "):
					trafficSelector.LocalIP = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "remote-ip "):
					trafficSelector.RemoteIP = types.StringValue(itemTrim)
				}
				rscData.TrafficSelector = append(rscData.TrafficSelector, trafficSelector)
			case balt.CutPrefixInString(&itemTrim, "vpn-monitor "):
				if rscData.VpnMonitor == nil {
					rscData.VpnMonitor = &securityIpsecVpnBlockVpnMonitor{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "destination-ip "):
					rscData.VpnMonitor.DestinationIP = types.StringValue(itemTrim)
				case itemTrim == "optimized":
					rscData.VpnMonitor.Optimized = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "source-interface "):
					rscData.VpnMonitor.SourceInterface = types.StringValue(itemTrim)
				}
			}
		}
	}

	return nil
}

func (rscData *securityIpsecVpnData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security ipsec vpn \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
