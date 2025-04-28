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
	_ resource.Resource                   = &securityIdpCustomAttack{}
	_ resource.ResourceWithConfigure      = &securityIdpCustomAttack{}
	_ resource.ResourceWithValidateConfig = &securityIdpCustomAttack{}
	_ resource.ResourceWithImportState    = &securityIdpCustomAttack{}
	_ resource.ResourceWithUpgradeState   = &securityIdpCustomAttack{}
)

type securityIdpCustomAttack struct {
	client *junos.Client
}

func newSecurityIdpCustomAttackResource() resource.Resource {
	return &securityIdpCustomAttack{}
}

func (rsc *securityIdpCustomAttack) typeName() string {
	return providerName + "_security_idp_custom_attack"
}

func (rsc *securityIdpCustomAttack) junosName() string {
	return "security idp custom-attack"
}

func (rsc *securityIdpCustomAttack) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityIdpCustomAttack) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityIdpCustomAttack) Configure(
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

func (rsc *securityIdpCustomAttack) Schema(
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
				Description: "Custom attack name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 60),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"severity": schema.StringAttribute{
				Required:    true,
				Description: "Select the severity that matches the lethality of this attack on your network.",
				Validators: []validator.String{
					stringvalidator.OneOf("critical", "info", "major", "minor", "warning"),
				},
			},
			"recommended_action": schema.StringAttribute{
				Optional:    true,
				Description: "Recommended action.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"close",
						"close-client",
						"close-server",
						"drop",
						"drop-packet",
						"ignore",
						"none",
					),
				},
			},
			"time_binding_count": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of times this attack is to be triggered.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"time_binding_scope": schema.StringAttribute{
				Optional:    true,
				Description: "Scope within which the count occurs.",
				Validators: []validator.String{
					stringvalidator.OneOf("destination", "peer", "source"),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"attack_type_anomaly": schema.SingleNestedBlock{
				Description: "Configure type of attack: Protocol anomaly.",
				Attributes:  securityIdpCustomAttackBlockAttackTypeAnomaly{}.attributesSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"attack_type_chain": schema.SingleNestedBlock{
				Description: "Configure type of attack: Chain attack.",
				Attributes: map[string]schema.Attribute{
					"expression": schema.StringAttribute{
						Optional:    true,
						Description: "Boolean Expression.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"order": schema.BoolAttribute{
						Optional:    true,
						Description: "Attacks should match in the order in which they are defined.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"protocol_binding": schema.StringAttribute{
						Optional:    true,
						Description: "Protocol binding over which attack will be detected.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(
								`^(application|icmp|ip|rpc|tcp|udp)`),
								"must have valid protocol (application|icmp|ip|rpc|tcp|udp) with optional option",
							),
						},
					},
					"reset": schema.BoolAttribute{
						Optional:    true,
						Description: "Repeat match should generate a new alert.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"scope": schema.StringAttribute{
						Optional:    true,
						Description: "Scope of the attack.",
						Validators: []validator.String{
							stringvalidator.OneOf("session", "transaction"),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"member": schema.ListNestedBlock{
						Description: "For each name of member attack to declare.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(2),
						},
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "Custom attack name.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
							},
							Blocks: map[string]schema.Block{
								"attack_type_anomaly": schema.SingleNestedBlock{
									Description: "Configure type of attack: Protocol anomaly.",
									Attributes:  securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly{}.attributesSchema(),
									PlanModifiers: []planmodifier.Object{
										tfplanmodifier.BlockRemoveNull(),
									},
								},
								"attack_type_signature": schema.SingleNestedBlock{
									Description: "Configure type of attack: Signature based attack.",
									Attributes:  securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature{}.attributesSchema(), //nolint:lll
									Blocks:      securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature{}.blocksSchema(),
									PlanModifiers: []planmodifier.Object{
										tfplanmodifier.BlockRemoveNull(),
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
			"attack_type_signature": schema.SingleNestedBlock{
				Description: "Configure type of attack: Signature based attack.",
				Attributes:  securityIdpCustomAttackBlockAttackTypeSignature{}.attributesSchema(),
				Blocks:      securityIdpCustomAttackBlockAttackTypeSignature{}.blocksSchema(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
		},
	}
}

type securityIdpCustomAttackData struct {
	ID                  types.String                                     `tfsdk:"id"`
	Name                types.String                                     `tfsdk:"name"`
	RecommendedAction   types.String                                     `tfsdk:"recommended_action"`
	Severity            types.String                                     `tfsdk:"severity"`
	TimeBindingCount    types.Int64                                      `tfsdk:"time_binding_count"`
	TimeBindingScope    types.String                                     `tfsdk:"time_binding_scope"`
	AttackTypeAnomaly   *securityIdpCustomAttackBlockAttackTypeAnomaly   `tfsdk:"attack_type_anomaly"`
	AttackTypeChain     *securityIdpCustomAttackBlockAttackTypeChain     `tfsdk:"attack_type_chain"`
	AttackTypeSignature *securityIdpCustomAttackBlockAttackTypeSignature `tfsdk:"attack_type_signature"`
}

type securityIdpCustomAttackConfig struct {
	ID                  types.String                                           `tfsdk:"id"`
	Name                types.String                                           `tfsdk:"name"`
	RecommendedAction   types.String                                           `tfsdk:"recommended_action"`
	Severity            types.String                                           `tfsdk:"severity"`
	TimeBindingCount    types.Int64                                            `tfsdk:"time_binding_count"`
	TimeBindingScope    types.String                                           `tfsdk:"time_binding_scope"`
	AttackTypeAnomaly   *securityIdpCustomAttackBlockAttackTypeAnomaly         `tfsdk:"attack_type_anomaly"`
	AttackTypeChain     *securityIdpCustomAttackBlockAttackTypeChainConfig     `tfsdk:"attack_type_chain"`
	AttackTypeSignature *securityIdpCustomAttackBlockAttackTypeSignatureConfig `tfsdk:"attack_type_signature"`
}

type securityIdpCustomAttackBlockAttackTypeAnomaly struct {
	securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly
	Service types.String `tfsdk:"service"`
}

func (block *securityIdpCustomAttackBlockAttackTypeAnomaly) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (securityIdpCustomAttackBlockAttackTypeAnomaly) attributesSchema() map[string]schema.Attribute {
	attributes := securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly{}.attributesSchema()
	attributes["service"] = schema.StringAttribute{
		Required:    false, // true when SingleNestedBlock is specified
		Optional:    true,
		Description: "Service name.",
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			tfvalidator.StringDoubleQuoteExclusion(),
		},
	}

	return attributes
}

type securityIdpCustomAttackBlockAttackTypeChain struct {
	Expression      types.String                                             `tfsdk:"expression"`
	Order           types.Bool                                               `tfsdk:"order"`
	ProtocolBinding types.String                                             `tfsdk:"protocol_binding"`
	Reset           types.Bool                                               `tfsdk:"reset"`
	Scope           types.String                                             `tfsdk:"scope"`
	Member          []securityIdpCustomAttackBlockAttackTypeChainBlockMember `tfsdk:"member"`
}

type securityIdpCustomAttackBlockAttackTypeChainConfig struct {
	Expression      types.String `tfsdk:"expression"`
	Order           types.Bool   `tfsdk:"order"`
	ProtocolBinding types.String `tfsdk:"protocol_binding"`
	Reset           types.Bool   `tfsdk:"reset"`
	Scope           types.String `tfsdk:"scope"`
	Member          types.List   `tfsdk:"member"`
}

func (block *securityIdpCustomAttackBlockAttackTypeChainConfig) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

//nolint:lll
type securityIdpCustomAttackBlockAttackTypeChainBlockMember struct {
	Name                types.String                                                                    `tfsdk:"name"                  tfdata:"identifier,skip_isempty"`
	AttackTypeAnomaly   *securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly   `tfsdk:"attack_type_anomaly"`
	AttackTypeSignature *securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature `tfsdk:"attack_type_signature"`
}

func (block *securityIdpCustomAttackBlockAttackTypeChainBlockMember) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:lll
type securityIdpCustomAttackBlockAttackTypeChainBlockMemberConfig struct {
	Name                types.String                                                                          `tfsdk:"name"                  tfdata:"skip_isempty"`
	AttackTypeAnomaly   *securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly         `tfsdk:"attack_type_anomaly"`
	AttackTypeSignature *securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignatureConfig `tfsdk:"attack_type_signature"`
}

func (block *securityIdpCustomAttackBlockAttackTypeChainBlockMemberConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly struct {
	Direction types.String `tfsdk:"direction"`
	Test      types.String `tfsdk:"test"`
	Shellcode types.String `tfsdk:"shellcode"`
}

func (block *securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly) attributesSchema() map[string]schema.Attribute { //nolint:lll
	return map[string]schema.Attribute{
		"direction": schema.StringAttribute{
			Required:    false, // true when SingleNestedBlock is specified
			Optional:    true,
			Description: "Connection direction of the attack.",
			Validators: []validator.String{
				stringvalidator.OneOf("any", "client-to-server", "server-to-client"),
			},
		},
		"test": schema.StringAttribute{
			Required:    false, // true when SingleNestedBlock is specified
			Optional:    true,
			Description: "Protocol anomaly condition to be checked.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"shellcode": schema.StringAttribute{
			Optional:    true,
			Description: "Specify shellcode flag for this attack.",
			Validators: []validator.String{
				stringvalidator.OneOf("all", "intel", "no-shellcode", "sparc"),
			},
		},
	}
}

type securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature struct {
	Context        types.String                                                      `tfsdk:"context"`
	Direction      types.String                                                      `tfsdk:"direction"`
	Negate         types.Bool                                                        `tfsdk:"negate"`
	Pattern        types.String                                                      `tfsdk:"pattern"`
	PatternPcre    types.String                                                      `tfsdk:"pattern_pcre"`
	Regexp         types.String                                                      `tfsdk:"regexp"`
	Shellcode      types.String                                                      `tfsdk:"shellcode"`
	ProtocolIcmp   *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp `tfsdk:"protocol_icmp"`
	ProtocolIcmpv6 *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp `tfsdk:"protocol_icmpv6"`
	ProtocolIPv4   *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4 `tfsdk:"protocol_ipv4"`
	ProtocolIPv6   *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv6 `tfsdk:"protocol_ipv6"`
	ProtocolTCP    *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCP  `tfsdk:"protocol_tcp"`
	ProtocolUDP    *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP  `tfsdk:"protocol_udp"`
}

func (securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature) attributesSchema() map[string]schema.Attribute { //nolint:lll
	return map[string]schema.Attribute{
		"context": schema.StringAttribute{
			Required:    false, // true when SingleNestedBlock is specified
			Optional:    true,
			Description: "Context.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"direction": schema.StringAttribute{
			Required:    false, // true when SingleNestedBlock is specified
			Optional:    true,
			Description: "Connection direction of the attack.",
			Validators: []validator.String{
				stringvalidator.OneOf("any", "client-to-server", "server-to-client"),
			},
		},
		"negate": schema.BoolAttribute{
			Optional:    true,
			Description: "Trigger the attack if condition is not met.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"pattern": schema.StringAttribute{
			Optional:    true,
			Description: "Pattern is the signature of the attack you want to detect.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"pattern_pcre": schema.StringAttribute{
			Optional:    true,
			Description: "Attack signature pattern in PCRE format.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 511),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"regexp": schema.StringAttribute{
			Optional:    true,
			Description: "Regular expression used for matching repetition of patterns.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"shellcode": schema.StringAttribute{
			Optional:    true,
			Description: "Specify shellcode flag for this attack.",
			Validators: []validator.String{
				stringvalidator.OneOf("all", "intel", "no-shellcode", "sparc"),
			},
		},
	}
}

func (securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature) blocksSchema() map[string]schema.Block { //nolint:lll
	return map[string]schema.Block{
		"protocol_icmp": schema.SingleNestedBlock{
			Description: "ICMP protocol parameters.",
			Attributes:  securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp{}.attributesSchema(false),
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"protocol_icmpv6": schema.SingleNestedBlock{
			Description: "ICMPv6 protocol parameters.",
			Attributes:  securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp{}.attributesSchema(true),
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"protocol_ipv4": schema.SingleNestedBlock{
			Description: "IPv4 protocol parameters",
			Attributes: map[string]schema.Attribute{
				"checksum_validate_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for validate checksum field against calculated checksum.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"checksum_validate_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for validate checksum field against calculated checksum.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"destination_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for destination IP-address.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"destination_value": schema.StringAttribute{
					Optional:    true,
					Description: "Value for destination IP-address.",
					Validators: []validator.String{
						tfvalidator.StringIPAddress().IPv4Only(),
					},
				},
				"identification_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for fragment identification.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"identification_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for fragment identification.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"ihl_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for header length in words.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"ihl_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for header length in words.",
					Validators: []validator.Int64{
						int64validator.Between(0, 15),
					},
				},
				"ip_flags": schema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Description: "IP Flag bits.",
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
						setvalidator.NoNullValues(),
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf("df", "mf", "rb", "no-df", "no-mf", "no-rb"),
						),
					},
				},
				"protocol_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for transport layer protocol.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"protocol_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for transport layer protocol.",
					Validators: []validator.Int64{
						int64validator.Between(0, 255),
					},
				},
				"source_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for source IP-address.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"source_value": schema.StringAttribute{
					Optional:    true,
					Description: "Value for source IP-address.",
					Validators: []validator.String{
						tfvalidator.StringIPAddress().IPv4Only(),
					},
				},
				"tos_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for type of service.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"tos_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for type of service.",
					Validators: []validator.Int64{
						int64validator.Between(0, 255),
					},
				},
				"total_length_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for total length of IP datagram.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"total_length_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for total length of IP datagram.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"ttl_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for time to live.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"ttl_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for time to live.",
					Validators: []validator.Int64{
						int64validator.Between(0, 255),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"protocol_ipv6": schema.SingleNestedBlock{
			Description: "IPv6 protocol parameters.",
			Attributes: map[string]schema.Attribute{
				"destination_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for destination IP-address.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"destination_value": schema.StringAttribute{
					Optional:    true,
					Description: "Value for destination IP-address.",
					Validators: []validator.String{
						tfvalidator.StringIPAddress().IPv6Only(),
					},
				},
				"extension_header_destination_option_home_address_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for home address of the mobile node in destination option extension header.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"extension_header_destination_option_home_address_value": schema.StringAttribute{
					Optional:    true,
					Description: "Value for home address of the mobile node in destination option extension header.",
					Validators: []validator.String{
						tfvalidator.StringIPAddress().IPv6Only(),
					},
				},
				"extension_header_destination_option_type_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for header type in destination option extension header.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"extension_header_destination_option_type_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for header type in  destination option extension header.",
					Validators: []validator.Int64{
						int64validator.Between(0, 255),
					},
				},
				"extension_header_routing_header_type_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for header type in routing extension header.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"extension_header_routing_header_type_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for header type in routing extension header.",
					Validators: []validator.Int64{
						int64validator.Between(0, 255),
					},
				},
				"flow_label_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for flow label identification.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"flow_label_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for flow label identification.",
					Validators: []validator.Int64{
						int64validator.Between(0, 1048575),
					},
				},
				"hop_limit_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for hop limit.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"hop_limit_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for hop limit.",
					Validators: []validator.Int64{
						int64validator.Between(0, 255),
					},
				},
				"next_header_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for the header following the basic IPv6 header.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"next_header_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for the header following the basic IPv6 header.",
					Validators: []validator.Int64{
						int64validator.Between(0, 255),
					},
				},
				"payload_length_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for length of the payload in the IPv6 datagram.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"payload_length_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for length of the payload in the IPv6 datagram.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"source_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for source IP-address.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"source_value": schema.StringAttribute{
					Optional:    true,
					Description: "Value for source IP-address.",
					Validators: []validator.String{
						tfvalidator.StringIPAddress().IPv6Only(),
					},
				},
				"traffic_class_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for traffic class. Similar to TOS in IPv4.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"traffic_class_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for traffic class. Similar to TOS in IPv4.",
					Validators: []validator.Int64{
						int64validator.Between(0, 255),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"protocol_tcp": schema.SingleNestedBlock{
			Description: "TCP protocol parameters.",
			Attributes: map[string]schema.Attribute{
				"ack_number_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for acknowledgement number.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"ack_number_value": schema.Int64Attribute{
					Optional:    true,
					Description: " Value for acknowledgement number.",
					Validators: []validator.Int64{
						int64validator.Between(0, 4294967295),
					},
				},
				"checksum_validate_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for validate checksum field against calculated checksum.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"checksum_validate_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for validate checksum field against calculated checksum.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"data_length_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for size of IP datagram subtracted by TCP header length.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"data_length_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for size of IP datagram subtracted by TCP header length.",
					Validators: []validator.Int64{
						int64validator.Between(2, 255),
					},
				},
				"destination_port_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for destination port.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"destination_port_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for destination port.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"header_length_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for header length in words.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"header_length_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for header length in words.",
					Validators: []validator.Int64{
						int64validator.Between(0, 15),
					},
				},
				"mss_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for maximum segment size.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"mss_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for maximum segment size.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"option_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for kind.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"option_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for kind.",
					Validators: []validator.Int64{
						int64validator.Between(0, 255),
					},
				},
				"reserved_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for three reserved bits.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"reserved_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for three reserved bits.",
					Validators: []validator.Int64{
						int64validator.Between(0, 7),
					},
				},
				"sequence_number_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for sequence number.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"sequence_number_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for sequence number.",
					Validators: []validator.Int64{
						int64validator.Between(0, 4294967295),
					},
				},
				"source_port_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for source port.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"source_port_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for source port.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"tcp_flags": schema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Description: "TCP header flags.",
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
						setvalidator.NoNullValues(),
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf(
								"ack", "fin", "psh", "r1", "r2", "rst", "syn", "urg",
								"no-ack", "no-fin", "no-psh", "no-r1", "no-r2", "no-rst", "no-syn", "no-urg",
							),
						),
					},
				},
				"urgent_pointer_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for urgent pointer.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"urgent_pointer_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for urgent pointer.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"window_scale_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for window scale.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"window_scale_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for sindow scale.",
					Validators: []validator.Int64{
						int64validator.Between(0, 255),
					},
				},
				"window_size_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for window size.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"window_size_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for window size.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"protocol_udp": schema.SingleNestedBlock{
			Description: "UDP protocol parameters.",
			Attributes: map[string]schema.Attribute{
				"checksum_validate_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for validate checksum field against calculated checksum.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"checksum_validate_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for validate checksum field against calculated checksum.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"data_length_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for size of IP datagram subtracted by UDP header length.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"data_length_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for size of IP datagram subtracted by UDP header length.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"destination_port_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for destination port.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"destination_port_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for destination port.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
				"source_port_match": schema.StringAttribute{
					Optional:    true,
					Description: "Condition for source port.",
					Validators: []validator.String{
						stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
					},
				},
				"source_port_value": schema.Int64Attribute{
					Optional:    true,
					Description: "Value for source port.",
					Validators: []validator.Int64{
						int64validator.Between(0, 65535),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
	}
}

type securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignatureConfig struct {
	Context        types.String                                                            `tfsdk:"context"`
	Direction      types.String                                                            `tfsdk:"direction"`
	Negate         types.Bool                                                              `tfsdk:"negate"`
	Pattern        types.String                                                            `tfsdk:"pattern"`
	PatternPcre    types.String                                                            `tfsdk:"pattern_pcre"`
	Regexp         types.String                                                            `tfsdk:"regexp"`
	Shellcode      types.String                                                            `tfsdk:"shellcode"`
	ProtocolIcmp   *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp       `tfsdk:"protocol_icmp"`
	ProtocolIcmpv6 *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp       `tfsdk:"protocol_icmpv6"`
	ProtocolIPv4   *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4Config `tfsdk:"protocol_ipv4"`
	ProtocolIPv6   *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv6       `tfsdk:"protocol_ipv6"`
	ProtocolTCP    *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCPConfig  `tfsdk:"protocol_tcp"`
	ProtocolUDP    *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP        `tfsdk:"protocol_udp"`
}

func (block *securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignatureConfig) hasKnownValue() bool { //nolint:lll
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityIdpCustomAttackBlockAttackTypeSignature struct {
	securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature
	ProtocolBinding types.String `tfsdk:"protocol_binding"`
}

func (securityIdpCustomAttackBlockAttackTypeSignature) attributesSchema() map[string]schema.Attribute {
	attributes := securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature{}.attributesSchema()
	attributes["protocol_binding"] = schema.StringAttribute{
		Optional:    true,
		Description: "Protocol binding over which attack will be detected.",
		Validators: []validator.String{
			stringvalidator.RegexMatches(regexp.MustCompile(
				`^(application|icmp|ip|rpc|tcp|udp)`),
				"must have valid protocol (application|icmp|ip|rpc|tcp|udp) with optional option",
			),
		},
	}

	return attributes
}

func (securityIdpCustomAttackBlockAttackTypeSignature) blocksSchema() map[string]schema.Block {
	return securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature{}.blocksSchema()
}

type securityIdpCustomAttackBlockAttackTypeSignatureConfig struct {
	securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignatureConfig
	ProtocolBinding types.String `tfsdk:"protocol_binding"`
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureConfig) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp struct {
	ChecksumValidateMatch types.String `tfsdk:"checksum_validate_match"`
	ChecksumValidateValue types.Int64  `tfsdk:"checksum_validate_value"`
	CodeMatch             types.String `tfsdk:"code_match"`
	CodeValue             types.Int64  `tfsdk:"code_value"`
	DataLengthMatch       types.String `tfsdk:"data_length_match"`
	DataLengthValue       types.Int64  `tfsdk:"data_length_value"`
	IdentificationMatch   types.String `tfsdk:"identification_match"`
	IdentificationValue   types.Int64  `tfsdk:"identification_value"`
	SequenceNumberMatch   types.String `tfsdk:"sequence_number_match"`
	SequenceNumberValue   types.Int64  `tfsdk:"sequence_number_value"`
	TypeMatch             types.String `tfsdk:"type_match"`
	TypeValue             types.Int64  `tfsdk:"type_value"`
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp) attributesSchema(
	v6 bool,
) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"checksum_validate_match": schema.StringAttribute{
			Optional:    true,
			Description: "Condition for validate checksum field against calculated checksum.",
			Validators: []validator.String{
				stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
			},
		},
		"checksum_validate_value": schema.Int64Attribute{
			Optional:    true,
			Description: "Value for validate checksum field against calculated checksum.",
			Validators: []validator.Int64{
				int64validator.Between(0, 65535),
			},
		},
		"code_match": schema.StringAttribute{
			Optional:    true,
			Description: "Condition for code field.",
			Validators: []validator.String{
				stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
			},
		},
		"code_value": schema.Int64Attribute{
			Optional:    true,
			Description: "Value for code field.",
			Validators: []validator.Int64{
				int64validator.Between(0, 255),
			},
		},
		"data_length_match": schema.StringAttribute{
			Optional:    true,
			Description: "Condition for size of IP datagram subtracted by ICMP header length.",
			Validators: []validator.String{
				stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
			},
		},
		"data_length_value": schema.Int64Attribute{
			Optional:    true,
			Description: "Value for size of IP datagram subtracted by ICMP header length.",
			Validators: func() []validator.Int64 {
				if v6 {
					return []validator.Int64{
						int64validator.Between(0, 255),
					}
				}

				return []validator.Int64{
					int64validator.Between(0, 65535),
				}
			}(),
		},
		"identification_match": schema.StringAttribute{
			Optional:    true,
			Description: "Condition for identifier in echo request/reply.",
			Validators: []validator.String{
				stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
			},
		},
		"identification_value": schema.Int64Attribute{
			Optional:    true,
			Description: "Value for identifier in echo request/reply.",
			Validators: []validator.Int64{
				int64validator.Between(0, 65535),
			},
		},
		"sequence_number_match": schema.StringAttribute{
			Optional:    true,
			Description: "Condition for sequence number.",
			Validators: []validator.String{
				stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
			},
		},
		"sequence_number_value": schema.Int64Attribute{
			Optional:    true,
			Description: "Value for sequence number.",
			Validators: []validator.Int64{
				int64validator.Between(0, 65535),
			},
		},
		"type_match": schema.StringAttribute{
			Optional:    true,
			Description: "Condition for type.",
			Validators: []validator.String{
				stringvalidator.OneOf("equal", "greater-than", "less-than", "not-equal"),
			},
		},
		"type_value": schema.Int64Attribute{
			Optional:    true,
			Description: "Value for type.",
			Validators: []validator.Int64{
				int64validator.Between(0, 255),
			},
		},
	}
}

type securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4 struct {
	ChecksumValidateMatch types.String   `tfsdk:"checksum_validate_match"`
	ChecksumValidateValue types.Int64    `tfsdk:"checksum_validate_value"`
	DestinationMatch      types.String   `tfsdk:"destination_match"`
	DestinationValue      types.String   `tfsdk:"destination_value"`
	IdentificationMatch   types.String   `tfsdk:"identification_match"`
	IdentificationValue   types.Int64    `tfsdk:"identification_value"`
	IhlMatch              types.String   `tfsdk:"ihl_match"`
	IhlValue              types.Int64    `tfsdk:"ihl_value"`
	IPFlags               []types.String `tfsdk:"ip_flags"`
	ProtocolMatch         types.String   `tfsdk:"protocol_match"`
	ProtocolValue         types.Int64    `tfsdk:"protocol_value"`
	SourceMatch           types.String   `tfsdk:"source_match"`
	SourceValue           types.String   `tfsdk:"source_value"`
	TosMatch              types.String   `tfsdk:"tos_match"`
	TosValue              types.Int64    `tfsdk:"tos_value"`
	TotalLengthMatch      types.String   `tfsdk:"total_length_match"`
	TotalLengthValue      types.Int64    `tfsdk:"total_length_value"`
	TTLMatch              types.String   `tfsdk:"ttl_match"`
	TTLValue              types.Int64    `tfsdk:"ttl_value"`
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4Config struct {
	ChecksumValidateMatch types.String `tfsdk:"checksum_validate_match"`
	ChecksumValidateValue types.Int64  `tfsdk:"checksum_validate_value"`
	DestinationMatch      types.String `tfsdk:"destination_match"`
	DestinationValue      types.String `tfsdk:"destination_value"`
	IdentificationMatch   types.String `tfsdk:"identification_match"`
	IdentificationValue   types.Int64  `tfsdk:"identification_value"`
	IhlMatch              types.String `tfsdk:"ihl_match"`
	IhlValue              types.Int64  `tfsdk:"ihl_value"`
	IPFlags               types.Set    `tfsdk:"ip_flags"`
	ProtocolMatch         types.String `tfsdk:"protocol_match"`
	ProtocolValue         types.Int64  `tfsdk:"protocol_value"`
	SourceMatch           types.String `tfsdk:"source_match"`
	SourceValue           types.String `tfsdk:"source_value"`
	TosMatch              types.String `tfsdk:"tos_match"`
	TosValue              types.Int64  `tfsdk:"tos_value"`
	TotalLengthMatch      types.String `tfsdk:"total_length_match"`
	TotalLengthValue      types.Int64  `tfsdk:"total_length_value"`
	TTLMatch              types.String `tfsdk:"ttl_match"`
	TTLValue              types.Int64  `tfsdk:"ttl_value"`
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4Config) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:lll
type securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv6 struct {
	DestinationMatch                                 types.String `tfsdk:"destination_match"`
	DestinationValue                                 types.String `tfsdk:"destination_value"`
	ExtensionHeaderDestinationOptionHomeAddressMatch types.String `tfsdk:"extension_header_destination_option_home_address_match"`
	ExtensionHeaderDestinationOptionHomeAddressValue types.String `tfsdk:"extension_header_destination_option_home_address_value"`
	ExtensionHeaderDestinationOptionTypeMatch        types.String `tfsdk:"extension_header_destination_option_type_match"`
	ExtensionHeaderDestinationOptionTypeValue        types.Int64  `tfsdk:"extension_header_destination_option_type_value"`
	ExtensionHeaderRoutingHeaderTypeMatch            types.String `tfsdk:"extension_header_routing_header_type_match"`
	ExtensionHeaderRoutingHeaderTypeValue            types.Int64  `tfsdk:"extension_header_routing_header_type_value"`
	FlowLabelMatch                                   types.String `tfsdk:"flow_label_match"`
	FlowLabelValue                                   types.Int64  `tfsdk:"flow_label_value"`
	HopLimitMatch                                    types.String `tfsdk:"hop_limit_match"`
	HopLimitValue                                    types.Int64  `tfsdk:"hop_limit_value"`
	NextHeaderMatch                                  types.String `tfsdk:"next_header_match"`
	NextHeaderValue                                  types.Int64  `tfsdk:"next_header_value"`
	PayloadLengthMatch                               types.String `tfsdk:"payload_length_match"`
	PayloadLengthValue                               types.Int64  `tfsdk:"payload_length_value"`
	SourceMatch                                      types.String `tfsdk:"source_match"`
	SourceValue                                      types.String `tfsdk:"source_value"`
	TrafficClassMatch                                types.String `tfsdk:"traffic_class_match"`
	TrafficClassValue                                types.Int64  `tfsdk:"traffic_class_value"`
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv6) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCP struct {
	AckNumberMatch        types.String   `tfsdk:"ack_number_match"`
	AckNumberValue        types.Int64    `tfsdk:"ack_number_value"`
	ChecksumValidateMatch types.String   `tfsdk:"checksum_validate_match"`
	ChecksumValidateValue types.Int64    `tfsdk:"checksum_validate_value"`
	DataLengthMatch       types.String   `tfsdk:"data_length_match"`
	DataLengthValue       types.Int64    `tfsdk:"data_length_value"`
	DestinationPortMatch  types.String   `tfsdk:"destination_port_match"`
	DestinationPortValue  types.Int64    `tfsdk:"destination_port_value"`
	HeaderLengthMatch     types.String   `tfsdk:"header_length_match"`
	HeaderLengthValue     types.Int64    `tfsdk:"header_length_value"`
	MssMatch              types.String   `tfsdk:"mss_match"`
	MssValue              types.Int64    `tfsdk:"mss_value"`
	OptionMatch           types.String   `tfsdk:"option_match"`
	OptionValue           types.Int64    `tfsdk:"option_value"`
	ReservedMatch         types.String   `tfsdk:"reserved_match"`
	ReservedValue         types.Int64    `tfsdk:"reserved_value"`
	SequenceNumberMatch   types.String   `tfsdk:"sequence_number_match"`
	SequenceNumberValue   types.Int64    `tfsdk:"sequence_number_value"`
	SourcePortMatch       types.String   `tfsdk:"source_port_match"`
	SourcePortValue       types.Int64    `tfsdk:"source_port_value"`
	TCPFlags              []types.String `tfsdk:"tcp_flags"`
	UrgentPointerMatch    types.String   `tfsdk:"urgent_pointer_match"`
	UrgentPointerValue    types.Int64    `tfsdk:"urgent_pointer_value"`
	WindowScaleMatch      types.String   `tfsdk:"window_scale_match"`
	WindowScaleValue      types.Int64    `tfsdk:"window_scale_value"`
	WindowSizeMatch       types.String   `tfsdk:"window_size_match"`
	WindowSizeValue       types.Int64    `tfsdk:"window_size_value"`
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCP) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCPConfig struct {
	AckNumberMatch        types.String `tfsdk:"ack_number_match"`
	AckNumberValue        types.Int64  `tfsdk:"ack_number_value"`
	ChecksumValidateMatch types.String `tfsdk:"checksum_validate_match"`
	ChecksumValidateValue types.Int64  `tfsdk:"checksum_validate_value"`
	DataLengthMatch       types.String `tfsdk:"data_length_match"`
	DataLengthValue       types.Int64  `tfsdk:"data_length_value"`
	DestinationPortMatch  types.String `tfsdk:"destination_port_match"`
	DestinationPortValue  types.Int64  `tfsdk:"destination_port_value"`
	HeaderLengthMatch     types.String `tfsdk:"header_length_match"`
	HeaderLengthValue     types.Int64  `tfsdk:"header_length_value"`
	MssMatch              types.String `tfsdk:"mss_match"`
	MssValue              types.Int64  `tfsdk:"mss_value"`
	OptionMatch           types.String `tfsdk:"option_match"`
	OptionValue           types.Int64  `tfsdk:"option_value"`
	ReservedMatch         types.String `tfsdk:"reserved_match"`
	ReservedValue         types.Int64  `tfsdk:"reserved_value"`
	SequenceNumberMatch   types.String `tfsdk:"sequence_number_match"`
	SequenceNumberValue   types.Int64  `tfsdk:"sequence_number_value"`
	SourcePortMatch       types.String `tfsdk:"source_port_match"`
	SourcePortValue       types.Int64  `tfsdk:"source_port_value"`
	TCPFlags              types.Set    `tfsdk:"tcp_flags"`
	UrgentPointerMatch    types.String `tfsdk:"urgent_pointer_match"`
	UrgentPointerValue    types.Int64  `tfsdk:"urgent_pointer_value"`
	WindowScaleMatch      types.String `tfsdk:"window_scale_match"`
	WindowScaleValue      types.Int64  `tfsdk:"window_scale_value"`
	WindowSizeMatch       types.String `tfsdk:"window_size_match"`
	WindowSizeValue       types.Int64  `tfsdk:"window_size_value"`
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCPConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCPConfig) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP struct {
	ChecksumValidateMatch types.String `tfsdk:"checksum_validate_match"`
	ChecksumValidateValue types.Int64  `tfsdk:"checksum_validate_value"`
	DataLengthMatch       types.String `tfsdk:"data_length_match"`
	DataLengthValue       types.Int64  `tfsdk:"data_length_value"`
	DestinationPortMatch  types.String `tfsdk:"destination_port_match"`
	DestinationPortValue  types.Int64  `tfsdk:"destination_port_value"`
	SourcePortMatch       types.String `tfsdk:"source_port_match"`
	SourcePortValue       types.Int64  `tfsdk:"source_port_value"`
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

func (rsc *securityIdpCustomAttack) ValidateConfig( //nolint:gocyclo,gocognit
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityIdpCustomAttackConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.AttackTypeAnomaly == nil &&
		config.AttackTypeChain == nil &&
		config.AttackTypeSignature == nil {
		resp.Diagnostics.AddError(
			tfdiag.MissingConfigErrSummary,
			"one of attack_type_anomaly, attack_type_chain or attack_type_signature must be specified",
		)
	}
	if config.AttackTypeAnomaly != nil &&
		config.AttackTypeAnomaly.hasKnownValue() &&
		config.AttackTypeChain != nil &&
		config.AttackTypeChain.hasKnownValue() {
		resp.Diagnostics.AddAttributeError(
			path.Root("attack_type_anomaly"),
			tfdiag.ConflictConfigErrSummary,
			"attack_type_anomaly and attack_type_chain cannot be configured together",
		)
	}
	if config.AttackTypeAnomaly != nil &&
		config.AttackTypeAnomaly.hasKnownValue() &&
		config.AttackTypeSignature != nil &&
		config.AttackTypeSignature.hasKnownValue() {
		resp.Diagnostics.AddAttributeError(
			path.Root("attack_type_anomaly"),
			tfdiag.ConflictConfigErrSummary,
			"attack_type_anomaly and attack_type_signature cannot be configured together",
		)
	}
	if config.AttackTypeChain != nil &&
		config.AttackTypeChain.hasKnownValue() &&
		config.AttackTypeSignature != nil &&
		config.AttackTypeSignature.hasKnownValue() {
		resp.Diagnostics.AddAttributeError(
			path.Root("attack_type_chain"),
			tfdiag.ConflictConfigErrSummary,
			"attack_type_chain and attack_type_signature cannot be configured together",
		)
	}

	if config.AttackTypeAnomaly != nil {
		if config.AttackTypeAnomaly.Direction.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("attack_type_anomaly").AtName("direction"),
				tfdiag.MissingConfigErrSummary,
				"direction must be specified"+
					" in attack_type_anomaly block",
			)
		}
		if config.AttackTypeAnomaly.Service.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("attack_type_anomaly").AtName("service"),
				tfdiag.MissingConfigErrSummary,
				"service must be specified"+
					" in attack_type_anomaly block",
			)
		}
		if config.AttackTypeAnomaly.Test.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("attack_type_anomaly").AtName("test"),
				tfdiag.MissingConfigErrSummary,
				"test must be specified"+
					" in attack_type_anomaly block",
			)
		}
	}
	if config.AttackTypeChain != nil {
		if config.AttackTypeChain.Member.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("attack_type_chain").AtName("member"),
				tfdiag.MissingConfigErrSummary,
				"member must be specified in attack_type_chain block",
			)
		} else if !config.AttackTypeChain.Member.IsUnknown() {
			var configMember []securityIdpCustomAttackBlockAttackTypeChainBlockMemberConfig
			asDiags := config.AttackTypeChain.Member.ElementsAs(ctx, &configMember, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			memberName := make(map[string]struct{})
			for i, block := range configMember {
				if !block.Name.IsUnknown() {
					name := block.Name.ValueString()
					if _, ok := memberName[name]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("attack_type_chain").AtName("member").AtListIndex(i).AtName("name"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple member blocks with the same name %q"+
								" in attack_type_chain block", name),
						)
					}
					memberName[name] = struct{}{}
				}

				if block.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("attack_type_chain").AtName("member").AtListIndex(i).AtName("name"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("one of attack_type_anomaly or attack_type_signature must be specified"+
							" in member block %q in attack_type_chain block", block.Name.ValueString()),
					)
				}
				if block.AttackTypeAnomaly != nil &&
					block.AttackTypeSignature != nil &&
					block.AttackTypeAnomaly.hasKnownValue() &&
					block.AttackTypeSignature.hasKnownValue() {
					resp.Diagnostics.AddAttributeError(
						path.Root("attack_type_chain").AtName("member").AtListIndex(i).AtName("attack_type_anomaly"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("attack_type_anomaly and attack_type_signature cannot be configured together"+
							" in member block %q in attack_type_chain block", block.Name.ValueString()),
					)
				}

				if block.AttackTypeAnomaly != nil {
					if block.AttackTypeAnomaly.Direction.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("attack_type_chain").AtName("member").AtListIndex(i).
								AtName("attack_type_anomaly").AtName("direction"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("direction must be specified"+
								" in attack_type_anomaly block in member block %q in attack_type_chain block", block.Name.ValueString()),
						)
					}
					if block.AttackTypeAnomaly.Test.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("attack_type_chain").AtName("member").AtListIndex(i).
								AtName("attack_type_anomaly").AtName("test"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("test must be specified"+
								" in attack_type_anomaly block in member block %q in attack_type_chain block", block.Name.ValueString()),
						)
					}
				}
				if block.AttackTypeSignature != nil {
					if block.AttackTypeSignature.Context.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("attack_type_chain").AtName("member").AtListIndex(i).
								AtName("attack_type_signature").AtName("context"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("context must be specified"+
								" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
						)
					}
					if block.AttackTypeSignature.Direction.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("attack_type_chain").AtName("member").AtListIndex(i).
								AtName("attack_type_signature").AtName("direction"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("direction must be specified"+
								" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
						)
					}
					if block.AttackTypeSignature.ProtocolIcmp != nil {
						if block.AttackTypeSignature.ProtocolIcmpv6 != nil &&
							block.AttackTypeSignature.ProtocolIcmpv6.hasKnownValue() &&
							block.AttackTypeSignature.ProtocolIcmp.hasKnownValue() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_icmp"),
								tfdiag.ConflictConfigErrSummary,
								fmt.Sprintf("protocol_icmp and protocol_icmpv6 cannot be configured together"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						}
						if block.AttackTypeSignature.ProtocolTCP != nil &&
							block.AttackTypeSignature.ProtocolTCP.hasKnownValue() &&
							block.AttackTypeSignature.ProtocolIcmp.hasKnownValue() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_icmp"),
								tfdiag.ConflictConfigErrSummary,
								fmt.Sprintf("protocol_icmp and protocol_tcp cannot be configured together"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						}
						if block.AttackTypeSignature.ProtocolUDP != nil &&
							block.AttackTypeSignature.ProtocolUDP.hasKnownValue() &&
							block.AttackTypeSignature.ProtocolIcmp.hasKnownValue() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_icmp"),
								tfdiag.ConflictConfigErrSummary,
								fmt.Sprintf("protocol_icmp and protocol_udp cannot be configured together"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						}

						if block.AttackTypeSignature.ProtocolIcmp.isEmpty() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_icmp").AtName("*"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("protocol_icmp block is empty"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						} else {
							block.AttackTypeSignature.ProtocolIcmp.validateConfig(
								ctx,
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_icmp"),
								fmt.Sprintf(" in protocol_icmp block in attack_type_signature block"+
									" in member block %q in attack_type_chain block", block.Name.ValueString()),
								resp,
							)
						}
					}

					if block.AttackTypeSignature.ProtocolIcmpv6 != nil {
						if block.AttackTypeSignature.ProtocolTCP != nil &&
							block.AttackTypeSignature.ProtocolTCP.hasKnownValue() &&
							block.AttackTypeSignature.ProtocolIcmpv6.hasKnownValue() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_icmpv6"),
								tfdiag.ConflictConfigErrSummary,
								fmt.Sprintf("protocol_icmpv6 and protocol_tcp cannot be configured together"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						}
						if block.AttackTypeSignature.ProtocolUDP != nil &&
							block.AttackTypeSignature.ProtocolUDP.hasKnownValue() &&
							block.AttackTypeSignature.ProtocolIcmpv6.hasKnownValue() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_icmpv6"),
								tfdiag.ConflictConfigErrSummary,
								fmt.Sprintf("protocol_icmpv6 and protocol_udp cannot be configured together"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						}

						if block.AttackTypeSignature.ProtocolIcmpv6.isEmpty() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_icmpv6").AtName("*"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("protocol_icmpv6 block is empty"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						} else {
							block.AttackTypeSignature.ProtocolIcmpv6.validateConfig(
								ctx,
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_icmpv6"),
								fmt.Sprintf(" in protocol_icmpv6 block in attack_type_signature block"+
									" in member block %q in attack_type_chain block", block.Name.ValueString()),
								resp,
							)
						}
					}

					if block.AttackTypeSignature.ProtocolIPv4 != nil {
						if block.AttackTypeSignature.ProtocolIPv4.isEmpty() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_ipv4").AtName("*"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("protocol_ipv4 block is empty"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						} else {
							block.AttackTypeSignature.ProtocolIPv4.validateConfig(
								ctx,
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_ipv4"),
								fmt.Sprintf(" in protocol_ipv4 block in attack_type_signature block"+
									" in member block %q in attack_type_chain block", block.Name.ValueString()),
								resp,
							)
						}
					}
					if block.AttackTypeSignature.ProtocolIPv6 != nil {
						if block.AttackTypeSignature.ProtocolIPv6.isEmpty() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_ipv6").AtName("*"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("protocol_ipv6 block is empty"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						} else {
							block.AttackTypeSignature.ProtocolIPv6.validateConfig(
								ctx,
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_ipv6"),
								fmt.Sprintf(" in protocol_ipv6 block in attack_type_signature block"+
									" in member block %q in attack_type_chain block", block.Name.ValueString()),
								resp,
							)
						}
					}
					if block.AttackTypeSignature.ProtocolTCP != nil {
						if block.AttackTypeSignature.ProtocolUDP != nil &&
							block.AttackTypeSignature.ProtocolUDP.hasKnownValue() &&
							block.AttackTypeSignature.ProtocolTCP.hasKnownValue() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_tcp"),
								tfdiag.ConflictConfigErrSummary,
								fmt.Sprintf("protocol_tcp and protocol_udp cannot be configured together"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						}

						if block.AttackTypeSignature.ProtocolTCP.isEmpty() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_tcp").AtName("*"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("protocol_tcp block is empty"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						} else {
							block.AttackTypeSignature.ProtocolTCP.validateConfig(
								ctx,
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_tcp"),
								fmt.Sprintf(" in protocol_tcp block in attack_type_signature block"+
									" in member block %q in attack_type_chain block", block.Name.ValueString()),
								resp,
							)
						}
					}

					if block.AttackTypeSignature.ProtocolUDP != nil {
						if block.AttackTypeSignature.ProtocolUDP.isEmpty() {
							resp.Diagnostics.AddAttributeError(
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_udp").AtName("*"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("protocol_udp block is empty"+
									" in attack_type_signature block in member block %q in attack_type_chain block", block.Name.ValueString()),
							)
						} else {
							block.AttackTypeSignature.ProtocolUDP.validateConfig(
								ctx,
								path.Root("attack_type_chain").AtName("member").AtListIndex(i).
									AtName("attack_type_signature").AtName("protocol_udp"),
								fmt.Sprintf(" in protocol_udp block in attack_type_signature block"+
									" in member block %q in attack_type_chain block", block.Name.ValueString()),
								resp,
							)
						}
					}
				}
			}
		}
	}
	if config.AttackTypeSignature != nil {
		if config.AttackTypeSignature.Context.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("attack_type_signature").AtName("context"),
				tfdiag.MissingConfigErrSummary,
				"context must be specified"+
					" in attack_type_signature block",
			)
		}
		if config.AttackTypeSignature.Direction.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("attack_type_signature").AtName("direction"),
				tfdiag.MissingConfigErrSummary,
				"direction must be specified"+
					" in attack_type_signature block",
			)
		}
		if config.AttackTypeSignature.ProtocolIcmp != nil {
			if config.AttackTypeSignature.ProtocolIcmpv6 != nil &&
				config.AttackTypeSignature.ProtocolIcmpv6.hasKnownValue() &&
				config.AttackTypeSignature.ProtocolIcmp.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_icmp"),
					tfdiag.ConflictConfigErrSummary,
					"protocol_icmp and protocol_icmpv6 cannot be configured together"+
						" in attack_type_signature block",
				)
			}
			if config.AttackTypeSignature.ProtocolTCP != nil &&
				config.AttackTypeSignature.ProtocolTCP.hasKnownValue() &&
				config.AttackTypeSignature.ProtocolIcmp.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_icmp"),
					tfdiag.ConflictConfigErrSummary,
					"protocol_icmp and protocol_tcp cannot be configured together"+
						" in attack_type_signature block",
				)
			}
			if config.AttackTypeSignature.ProtocolUDP != nil &&
				config.AttackTypeSignature.ProtocolUDP.hasKnownValue() &&
				config.AttackTypeSignature.ProtocolIcmp.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_icmp"),
					tfdiag.ConflictConfigErrSummary,
					"protocol_icmp and protocol_udp cannot be configured together"+
						" in attack_type_signature block",
				)
			}

			if config.AttackTypeSignature.ProtocolIcmp.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_icmp").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"protocol_icmp block is empty"+
						" in attack_type_signature block",
				)
			} else {
				config.AttackTypeSignature.ProtocolIcmp.validateConfig(
					ctx,
					path.Root("attack_type_signature").AtName("protocol_icmp"),
					" in protocol_icmp block in attack_type_signature block",
					resp,
				)
			}
		}

		if config.AttackTypeSignature.ProtocolIcmpv6 != nil {
			if config.AttackTypeSignature.ProtocolTCP != nil &&
				config.AttackTypeSignature.ProtocolTCP.hasKnownValue() &&
				config.AttackTypeSignature.ProtocolIcmpv6.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_icmpv6"),
					tfdiag.ConflictConfigErrSummary,
					"protocol_icmpv6 and protocol_tcp cannot be configured together"+
						" in attack_type_signature block",
				)
			}
			if config.AttackTypeSignature.ProtocolUDP != nil &&
				config.AttackTypeSignature.ProtocolUDP.hasKnownValue() &&
				config.AttackTypeSignature.ProtocolIcmpv6.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_icmpv6"),
					tfdiag.ConflictConfigErrSummary,
					"protocol_icmpv6 and protocol_udp cannot be configured together"+
						" in attack_type_signature block",
				)
			}

			if config.AttackTypeSignature.ProtocolIcmpv6.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_icmpv6").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"protocol_icmpv6 block is empty"+
						" in attack_type_signature block",
				)
			} else {
				config.AttackTypeSignature.ProtocolIcmpv6.validateConfig(
					ctx,
					path.Root("attack_type_signature").AtName("protocol_icmpv6"),
					" in protocol_icmpv6 block in attack_type_signature block",
					resp,
				)
			}
		}

		if config.AttackTypeSignature.ProtocolIPv4 != nil {
			if config.AttackTypeSignature.ProtocolIPv4.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_ipv4").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"protocol_ipv4 block is empty"+
						" in attack_type_signature block",
				)
			} else {
				config.AttackTypeSignature.ProtocolIPv4.validateConfig(
					ctx,
					path.Root("attack_type_signature").AtName("protocol_ipv4"),
					" in protocol_icmpv6 block in attack_type_signature block",
					resp,
				)
			}
		}
		if config.AttackTypeSignature.ProtocolIPv6 != nil {
			if config.AttackTypeSignature.ProtocolIPv6.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_ipv6").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"protocol_ipv6 block is empty"+
						" in attack_type_signature block",
				)
			} else {
				config.AttackTypeSignature.ProtocolIPv6.validateConfig(
					ctx,
					path.Root("attack_type_signature").AtName("protocol_ipv6"),
					" in protocol_ipv6 block in attack_type_signature block",
					resp,
				)
			}
		}
		if config.AttackTypeSignature.ProtocolTCP != nil {
			if config.AttackTypeSignature.ProtocolUDP != nil &&
				config.AttackTypeSignature.ProtocolUDP.hasKnownValue() &&
				config.AttackTypeSignature.ProtocolTCP.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_tcp"),
					tfdiag.ConflictConfigErrSummary,
					"protocol_tcp and protocol_udp cannot be configured together"+
						" in attack_type_signature block",
				)
			}

			if config.AttackTypeSignature.ProtocolTCP.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_tcp").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"protocol_tcp block is empty"+
						" in attack_type_signature block",
				)
			} else {
				config.AttackTypeSignature.ProtocolTCP.validateConfig(
					ctx,
					path.Root("attack_type_signature").AtName("protocol_tcp"),
					" in protocol_tcp block in attack_type_signature block",
					resp,
				)
			}
		}

		if config.AttackTypeSignature.ProtocolUDP != nil {
			if config.AttackTypeSignature.ProtocolUDP.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("attack_type_signature").AtName("protocol_udp").AtName("*"),
					tfdiag.MissingConfigErrSummary,
					"protocol_udp block is empty"+
						" in attack_type_signature block",
				)
			} else {
				config.AttackTypeSignature.ProtocolUDP.validateConfig(
					ctx,
					path.Root("attack_type_signature").AtName("protocol_udp"),
					" in protocol_udp block in attack_type_signature block",
					resp,
				)
			}
		}
	}
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp) validateConfig(
	_ context.Context,
	pathRoot path.Path,
	blockErrorSuffix string,
	resp *resource.ValidateConfigResponse,
) {
	if !block.ChecksumValidateMatch.IsNull() &&
		!block.ChecksumValidateMatch.IsUnknown() &&
		block.ChecksumValidateValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("checksum_validate_match"),
			tfdiag.MissingConfigErrSummary,
			"checksum_validate_value must be specified with checksum_validate_match"+blockErrorSuffix,
		)
	}
	if !block.ChecksumValidateValue.IsNull() &&
		!block.ChecksumValidateValue.IsUnknown() &&
		block.ChecksumValidateMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("checksum_validate_value"),
			tfdiag.MissingConfigErrSummary,
			"checksum_validate_match must be specified with checksum_validate_value"+blockErrorSuffix,
		)
	}
	if !block.CodeMatch.IsNull() &&
		!block.CodeMatch.IsUnknown() &&
		block.CodeValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("code_match"),
			tfdiag.MissingConfigErrSummary,
			"code_value must be specified with code_match"+blockErrorSuffix,
		)
	}
	if !block.CodeValue.IsNull() &&
		!block.CodeValue.IsUnknown() &&
		block.CodeMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("code_value"),
			tfdiag.MissingConfigErrSummary,
			"code_match must be specified with code_value"+blockErrorSuffix,
		)
	}
	if !block.DataLengthMatch.IsNull() &&
		!block.DataLengthMatch.IsUnknown() &&
		block.DataLengthValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("data_length_match"),
			tfdiag.MissingConfigErrSummary,
			"data_length_value must be specified with data_length_match"+blockErrorSuffix,
		)
	}
	if !block.DataLengthValue.IsNull() &&
		!block.DataLengthValue.IsUnknown() &&
		block.DataLengthMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("data_length_value"),
			tfdiag.MissingConfigErrSummary,
			"data_length_match must be specified with data_length_value"+blockErrorSuffix,
		)
	}
	if !block.IdentificationMatch.IsNull() &&
		!block.IdentificationMatch.IsUnknown() &&
		block.IdentificationValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("identification_match"),
			tfdiag.MissingConfigErrSummary,
			"identification_value must be specified with identification_match"+blockErrorSuffix,
		)
	}
	if !block.IdentificationValue.IsNull() &&
		!block.IdentificationValue.IsUnknown() &&
		block.IdentificationMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("identification_value"),
			tfdiag.MissingConfigErrSummary,
			"identification_match must be specified with identification_value"+blockErrorSuffix,
		)
	}
	if !block.SequenceNumberMatch.IsNull() &&
		!block.SequenceNumberMatch.IsUnknown() &&
		block.SequenceNumberValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("sequence_number_match"),
			tfdiag.MissingConfigErrSummary,
			"sequence_number_value must be specified with sequence_number_match"+blockErrorSuffix,
		)
	}
	if !block.SequenceNumberValue.IsNull() &&
		!block.SequenceNumberValue.IsUnknown() &&
		block.SequenceNumberMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("sequence_number_value"),
			tfdiag.MissingConfigErrSummary,
			"sequence_number_match must be specified with sequence_number_value"+blockErrorSuffix,
		)
	}
	if !block.TypeMatch.IsNull() &&
		!block.TypeMatch.IsUnknown() &&
		block.TypeValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("type_match"),
			tfdiag.MissingConfigErrSummary,
			"type_value must be specified with type_match"+blockErrorSuffix,
		)
	}
	if !block.TypeValue.IsNull() &&
		!block.TypeValue.IsUnknown() &&
		block.TypeMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("type_value"),
			tfdiag.MissingConfigErrSummary,
			"type_match must be specified with type_value"+blockErrorSuffix,
		)
	}
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4Config) validateConfig(
	_ context.Context,
	pathRoot path.Path,
	blockErrorSuffix string,
	resp *resource.ValidateConfigResponse,
) {
	if !block.ChecksumValidateMatch.IsNull() &&
		!block.ChecksumValidateMatch.IsUnknown() &&
		block.ChecksumValidateValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("checksum_validate_match"),
			tfdiag.MissingConfigErrSummary,
			"checksum_validate_value must be specified with checksum_validate_match"+blockErrorSuffix,
		)
	}
	if !block.ChecksumValidateValue.IsNull() &&
		!block.ChecksumValidateValue.IsUnknown() &&
		block.ChecksumValidateMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("checksum_validate_value"),
			tfdiag.MissingConfigErrSummary,
			"checksum_validate_match must be specified with checksum_validate_value"+blockErrorSuffix,
		)
	}
	if !block.DestinationMatch.IsNull() &&
		!block.DestinationMatch.IsUnknown() &&
		block.DestinationValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_match"),
			tfdiag.MissingConfigErrSummary,
			"destination_value must be specified with destination_match"+blockErrorSuffix,
		)
	}
	if !block.DestinationValue.IsNull() &&
		!block.DestinationValue.IsUnknown() &&
		block.DestinationMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_value"),
			tfdiag.MissingConfigErrSummary,
			"destination_match must be specified with destination_value"+blockErrorSuffix,
		)
	}
	if !block.IdentificationMatch.IsNull() &&
		!block.IdentificationMatch.IsUnknown() &&
		block.IdentificationValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("identification_match"),
			tfdiag.MissingConfigErrSummary,
			"identification_value must be specified with identification_match"+blockErrorSuffix,
		)
	}
	if !block.IdentificationValue.IsNull() &&
		!block.IdentificationValue.IsUnknown() &&
		block.IdentificationMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("identification_value"),
			tfdiag.MissingConfigErrSummary,
			"identification_match must be specified with identification_value"+blockErrorSuffix,
		)
	}
	if !block.IhlMatch.IsNull() &&
		!block.IhlMatch.IsUnknown() &&
		block.IhlValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("ihl_match"),
			tfdiag.MissingConfigErrSummary,
			"ihl_value must be specified with ihl_match"+blockErrorSuffix,
		)
	}
	if !block.IhlValue.IsNull() &&
		!block.IhlValue.IsUnknown() &&
		block.IhlMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("ihl_value"),
			tfdiag.MissingConfigErrSummary,
			"ihl_match must be specified with ihl_value"+blockErrorSuffix,
		)
	}
	if !block.ProtocolMatch.IsNull() &&
		!block.ProtocolMatch.IsUnknown() &&
		block.ProtocolValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("protocol_match"),
			tfdiag.MissingConfigErrSummary,
			"protocol_value must be specified with protocol_match"+blockErrorSuffix,
		)
	}
	if !block.ProtocolValue.IsNull() &&
		!block.ProtocolValue.IsUnknown() &&
		block.ProtocolMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("protocol_value"),
			tfdiag.MissingConfigErrSummary,
			"protocol_match must be specified with protocol_value"+blockErrorSuffix,
		)
	}
	if !block.SourceMatch.IsNull() &&
		!block.SourceMatch.IsUnknown() &&
		block.SourceValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_match"),
			tfdiag.MissingConfigErrSummary,
			"source_value must be specified with source_match"+blockErrorSuffix,
		)
	}
	if !block.SourceValue.IsNull() &&
		!block.SourceValue.IsUnknown() &&
		block.SourceMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_value"),
			tfdiag.MissingConfigErrSummary,
			"source_match must be specified with source_value"+blockErrorSuffix,
		)
	}
	if !block.TosMatch.IsNull() &&
		!block.TosMatch.IsUnknown() &&
		block.TosValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("tos_match"),
			tfdiag.MissingConfigErrSummary,
			"tos_value must be specified with tos_match"+blockErrorSuffix,
		)
	}
	if !block.TosValue.IsNull() &&
		!block.TosValue.IsUnknown() &&
		block.TosMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("tos_value"),
			tfdiag.MissingConfigErrSummary,
			"tos_match must be specified with tos_value"+blockErrorSuffix,
		)
	}
	if !block.TotalLengthMatch.IsNull() &&
		!block.TotalLengthMatch.IsUnknown() &&
		block.TotalLengthValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("total_length_match"),
			tfdiag.MissingConfigErrSummary,
			"total_length_value must be specified with total_length_match"+blockErrorSuffix,
		)
	}
	if !block.TotalLengthValue.IsNull() &&
		!block.TotalLengthValue.IsUnknown() &&
		block.TotalLengthMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("total_length_value"),
			tfdiag.MissingConfigErrSummary,
			"total_length_match must be specified with total_length_value"+blockErrorSuffix,
		)
	}
	if !block.TTLMatch.IsNull() &&
		!block.TTLMatch.IsUnknown() &&
		block.TTLValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("ttl_match"),
			tfdiag.MissingConfigErrSummary,
			"ttl_value must be specified with ttl_match"+blockErrorSuffix,
		)
	}
	if !block.TTLValue.IsNull() &&
		!block.TTLValue.IsUnknown() &&
		block.TTLMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("ttl_value"),
			tfdiag.MissingConfigErrSummary,
			"ttl_match must be specified with ttl_value"+blockErrorSuffix,
		)
	}
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv6) validateConfig(
	_ context.Context,
	pathRoot path.Path,
	blockErrorSuffix string,
	resp *resource.ValidateConfigResponse,
) {
	if !block.DestinationMatch.IsNull() &&
		!block.DestinationMatch.IsUnknown() &&
		block.DestinationValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_match"),
			tfdiag.MissingConfigErrSummary,
			"destination_value must be specified with destination_match"+blockErrorSuffix,
		)
	}
	if !block.DestinationValue.IsNull() &&
		!block.DestinationValue.IsUnknown() &&
		block.DestinationMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_value"),
			tfdiag.MissingConfigErrSummary,
			"destination_match must be specified with destination_value"+blockErrorSuffix,
		)
	}
	if !block.ExtensionHeaderDestinationOptionHomeAddressMatch.IsNull() &&
		!block.ExtensionHeaderDestinationOptionHomeAddressMatch.IsUnknown() &&
		block.ExtensionHeaderDestinationOptionHomeAddressValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("extension_header_destination_option_home_address_match"),
			tfdiag.MissingConfigErrSummary,
			"extension_header_destination_option_home_address_value must be specified"+
				" with extension_header_destination_option_home_address_match"+blockErrorSuffix,
		)
	}
	if !block.ExtensionHeaderDestinationOptionHomeAddressValue.IsNull() &&
		!block.ExtensionHeaderDestinationOptionHomeAddressValue.IsUnknown() &&
		block.ExtensionHeaderDestinationOptionHomeAddressMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("extension_header_destination_option_home_address_value"),
			tfdiag.MissingConfigErrSummary,
			"extension_header_destination_option_home_address_match must be specified"+
				" with extension_header_destination_option_home_address_value"+blockErrorSuffix,
		)
	}
	if !block.ExtensionHeaderDestinationOptionTypeMatch.IsNull() &&
		!block.ExtensionHeaderDestinationOptionTypeMatch.IsUnknown() &&
		block.ExtensionHeaderDestinationOptionTypeValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("extension_header_destination_option_type_match"),
			tfdiag.MissingConfigErrSummary,
			"extension_header_destination_option_type_value must be specified"+
				" with extension_header_destination_option_type_match"+blockErrorSuffix,
		)
	}
	if !block.ExtensionHeaderDestinationOptionTypeValue.IsNull() &&
		!block.ExtensionHeaderDestinationOptionTypeValue.IsUnknown() &&
		block.ExtensionHeaderDestinationOptionTypeMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("extension_header_destination_option_type_value"),
			tfdiag.MissingConfigErrSummary,
			"extension_header_destination_option_type_match must be specified"+
				" with extension_header_destination_option_type_value"+blockErrorSuffix,
		)
	}
	if !block.ExtensionHeaderRoutingHeaderTypeMatch.IsNull() &&
		!block.ExtensionHeaderRoutingHeaderTypeMatch.IsUnknown() &&
		block.ExtensionHeaderRoutingHeaderTypeValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("extension_header_routing_header_type_match"),
			tfdiag.MissingConfigErrSummary,
			"extension_header_routing_header_type_value must be specified"+
				" with extension_header_routing_header_type_match"+blockErrorSuffix,
		)
	}
	if !block.ExtensionHeaderRoutingHeaderTypeValue.IsNull() &&
		!block.ExtensionHeaderRoutingHeaderTypeValue.IsUnknown() &&
		block.ExtensionHeaderRoutingHeaderTypeMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("extension_header_routing_header_type_value"),
			tfdiag.MissingConfigErrSummary,
			"extension_header_routing_header_type_match must be specified"+
				" with extension_header_routing_header_type_value"+blockErrorSuffix,
		)
	}
	if !block.FlowLabelMatch.IsNull() &&
		!block.FlowLabelMatch.IsUnknown() &&
		block.FlowLabelValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("flow_label_match"),
			tfdiag.MissingConfigErrSummary,
			"flow_label_value must be specified with flow_label_match"+blockErrorSuffix,
		)
	}
	if !block.FlowLabelValue.IsNull() &&
		!block.FlowLabelValue.IsUnknown() &&
		block.FlowLabelMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("flow_label_value"),
			tfdiag.MissingConfigErrSummary,
			"flow_label_match must be specified with flow_label_value"+blockErrorSuffix,
		)
	}
	if !block.HopLimitMatch.IsNull() &&
		!block.HopLimitMatch.IsUnknown() &&
		block.HopLimitValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("hop_limit_match"),
			tfdiag.MissingConfigErrSummary,
			"hop_limit_value must be specified with hop_limit_match"+blockErrorSuffix,
		)
	}
	if !block.HopLimitValue.IsNull() &&
		!block.HopLimitValue.IsUnknown() &&
		block.HopLimitMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("hop_limit_value"),
			tfdiag.MissingConfigErrSummary,
			"hop_limit_match must be specified with hop_limit_value"+blockErrorSuffix,
		)
	}
	if !block.NextHeaderMatch.IsNull() &&
		!block.NextHeaderMatch.IsUnknown() &&
		block.NextHeaderValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("next_header_match"),
			tfdiag.MissingConfigErrSummary,
			"next_header_value must be specified with next_header_match"+blockErrorSuffix,
		)
	}
	if !block.NextHeaderValue.IsNull() &&
		!block.NextHeaderValue.IsUnknown() &&
		block.NextHeaderMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("next_header_value"),
			tfdiag.MissingConfigErrSummary,
			"next_header_match must be specified with next_header_value"+blockErrorSuffix,
		)
	}
	if !block.PayloadLengthMatch.IsNull() &&
		!block.PayloadLengthMatch.IsUnknown() &&
		block.PayloadLengthValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("payload_length_match"),
			tfdiag.MissingConfigErrSummary,
			"payload_length_value must be specified with payload_length_match"+blockErrorSuffix,
		)
	}
	if !block.PayloadLengthValue.IsNull() &&
		!block.PayloadLengthValue.IsUnknown() &&
		block.PayloadLengthMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("payload_length_value"),
			tfdiag.MissingConfigErrSummary,
			"payload_length_match must be specified with payload_length_value"+blockErrorSuffix,
		)
	}
	if !block.SourceMatch.IsNull() &&
		!block.SourceMatch.IsUnknown() &&
		block.SourceValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_match"),
			tfdiag.MissingConfigErrSummary,
			"source_value must be specified with source_match"+blockErrorSuffix,
		)
	}
	if !block.SourceValue.IsNull() &&
		!block.SourceValue.IsUnknown() &&
		block.SourceMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_value"),
			tfdiag.MissingConfigErrSummary,
			"source_match must be specified with source_value"+blockErrorSuffix,
		)
	}
	if !block.TrafficClassMatch.IsNull() &&
		!block.TrafficClassMatch.IsUnknown() &&
		block.TrafficClassValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("traffic_class_match"),
			tfdiag.MissingConfigErrSummary,
			"traffic_class_value must be specified with traffic_class_match"+blockErrorSuffix,
		)
	}
	if !block.TrafficClassValue.IsNull() &&
		!block.TrafficClassValue.IsUnknown() &&
		block.TrafficClassMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("traffic_class_value"),
			tfdiag.MissingConfigErrSummary,
			"traffic_class_match must be specified with traffic_class_value"+blockErrorSuffix,
		)
	}
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCPConfig) validateConfig(
	_ context.Context,
	pathRoot path.Path,
	blockErrorSuffix string,
	resp *resource.ValidateConfigResponse,
) {
	if !block.AckNumberMatch.IsNull() &&
		!block.AckNumberMatch.IsUnknown() &&
		block.AckNumberValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("ack_number_match"),
			tfdiag.MissingConfigErrSummary,
			"ack_number_value must be specified with ack_number_match"+blockErrorSuffix,
		)
	}
	if !block.AckNumberValue.IsNull() &&
		!block.AckNumberValue.IsUnknown() &&
		block.AckNumberMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("ack_number_value"),
			tfdiag.MissingConfigErrSummary,
			"ack_number_match must be specified with ack_number_value"+blockErrorSuffix,
		)
	}
	if !block.ChecksumValidateMatch.IsNull() &&
		!block.ChecksumValidateMatch.IsUnknown() &&
		block.ChecksumValidateValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("checksum_validate_match"),
			tfdiag.MissingConfigErrSummary,
			"checksum_validate_value must be specified with checksum_validate_match"+blockErrorSuffix,
		)
	}
	if !block.ChecksumValidateValue.IsNull() &&
		!block.ChecksumValidateValue.IsUnknown() &&
		block.ChecksumValidateMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("checksum_validate_value"),
			tfdiag.MissingConfigErrSummary,
			"checksum_validate_match must be specified with checksum_validate_value"+blockErrorSuffix,
		)
	}
	if !block.DataLengthMatch.IsNull() &&
		!block.DataLengthMatch.IsUnknown() &&
		block.DataLengthValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("data_length_match"),
			tfdiag.MissingConfigErrSummary,
			"data_length_value must be specified with data_length_match"+blockErrorSuffix,
		)
	}
	if !block.DataLengthValue.IsNull() &&
		!block.DataLengthValue.IsUnknown() &&
		block.DataLengthMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("data_length_value"),
			tfdiag.MissingConfigErrSummary,
			"data_length_match must be specified with data_length_value"+blockErrorSuffix,
		)
	}
	if !block.DestinationPortMatch.IsNull() &&
		!block.DestinationPortMatch.IsUnknown() &&
		block.DestinationPortValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_port_match"),
			tfdiag.MissingConfigErrSummary,
			"destination_port_value must be specified with destination_port_match"+blockErrorSuffix,
		)
	}
	if !block.DestinationPortValue.IsNull() &&
		!block.DestinationPortValue.IsUnknown() &&
		block.DestinationPortMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_port_value"),
			tfdiag.MissingConfigErrSummary,
			"destination_port_match must be specified with destination_port_value"+blockErrorSuffix,
		)
	}
	if !block.HeaderLengthMatch.IsNull() &&
		!block.HeaderLengthMatch.IsUnknown() &&
		block.HeaderLengthValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("header_length_match"),
			tfdiag.MissingConfigErrSummary,
			"header_length_value must be specified with header_length_match"+blockErrorSuffix,
		)
	}
	if !block.HeaderLengthValue.IsNull() &&
		!block.HeaderLengthValue.IsUnknown() &&
		block.HeaderLengthMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("header_length_value"),
			tfdiag.MissingConfigErrSummary,
			"header_length_match must be specified with header_length_value"+blockErrorSuffix,
		)
	}
	if !block.MssMatch.IsNull() &&
		!block.MssMatch.IsUnknown() &&
		block.MssValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("mss_match"),
			tfdiag.MissingConfigErrSummary,
			"mss_value must be specified with mss_match"+blockErrorSuffix,
		)
	}
	if !block.MssValue.IsNull() &&
		!block.MssValue.IsUnknown() &&
		block.MssMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("mss_value"),
			tfdiag.MissingConfigErrSummary,
			"mss_match must be specified with mss_value"+blockErrorSuffix,
		)
	}
	if !block.OptionMatch.IsNull() &&
		!block.OptionMatch.IsUnknown() &&
		block.OptionValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("option_match"),
			tfdiag.MissingConfigErrSummary,
			"option_value must be specified with option_match"+blockErrorSuffix,
		)
	}
	if !block.OptionValue.IsNull() &&
		!block.OptionValue.IsUnknown() &&
		block.OptionMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("option_value"),
			tfdiag.MissingConfigErrSummary,
			"option_match must be specified with option_value"+blockErrorSuffix,
		)
	}
	if !block.ReservedMatch.IsNull() &&
		!block.ReservedMatch.IsUnknown() &&
		block.ReservedValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("reserved_match"),
			tfdiag.MissingConfigErrSummary,
			"reserved_value must be specified with reserved_match"+blockErrorSuffix,
		)
	}
	if !block.ReservedValue.IsNull() &&
		!block.ReservedValue.IsUnknown() &&
		block.ReservedMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("reserved_value"),
			tfdiag.MissingConfigErrSummary,
			"reserved_match must be specified with reserved_value"+blockErrorSuffix,
		)
	}
	if !block.SequenceNumberMatch.IsNull() &&
		!block.SequenceNumberMatch.IsUnknown() &&
		block.SequenceNumberValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("sequence_number_match"),
			tfdiag.MissingConfigErrSummary,
			"sequence_number_value must be specified with sequence_number_match"+blockErrorSuffix,
		)
	}
	if !block.SequenceNumberValue.IsNull() &&
		!block.SequenceNumberValue.IsUnknown() &&
		block.SequenceNumberMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("sequence_number_value"),
			tfdiag.MissingConfigErrSummary,
			"sequence_number_match must be specified with sequence_number_value"+blockErrorSuffix,
		)
	}
	if !block.SourcePortMatch.IsNull() &&
		!block.SourcePortMatch.IsUnknown() &&
		block.SourcePortValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_port_match"),
			tfdiag.MissingConfigErrSummary,
			"source_port_value must be specified with source_port_match"+blockErrorSuffix,
		)
	}
	if !block.SourcePortValue.IsNull() &&
		!block.SourcePortValue.IsUnknown() &&
		block.SourcePortMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_port_value"),
			tfdiag.MissingConfigErrSummary,
			"source_port_match must be specified with source_port_value"+blockErrorSuffix,
		)
	}
	if !block.UrgentPointerMatch.IsNull() &&
		!block.UrgentPointerMatch.IsUnknown() &&
		block.UrgentPointerValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("urgent_pointer_match"),
			tfdiag.MissingConfigErrSummary,
			"urgent_pointer_value must be specified with urgent_pointer_match"+blockErrorSuffix,
		)
	}
	if !block.UrgentPointerValue.IsNull() &&
		!block.UrgentPointerValue.IsUnknown() &&
		block.UrgentPointerMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("urgent_pointer_value"),
			tfdiag.MissingConfigErrSummary,
			"urgent_pointer_match must be specified with urgent_pointer_value"+blockErrorSuffix,
		)
	}
	if !block.WindowScaleMatch.IsNull() &&
		!block.WindowScaleMatch.IsUnknown() &&
		block.WindowScaleValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("window_scale_match"),
			tfdiag.MissingConfigErrSummary,
			"window_scale_value must be specified with window_scale_match"+blockErrorSuffix,
		)
	}
	if !block.WindowScaleValue.IsNull() &&
		!block.WindowScaleValue.IsUnknown() &&
		block.WindowScaleMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("window_scale_value"),
			tfdiag.MissingConfigErrSummary,
			"window_scale_match must be specified with window_scale_value"+blockErrorSuffix,
		)
	}
	if !block.WindowSizeMatch.IsNull() &&
		!block.WindowSizeMatch.IsUnknown() &&
		block.WindowSizeValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("window_size_match"),
			tfdiag.MissingConfigErrSummary,
			"window_size_value must be specified with window_size_match"+blockErrorSuffix,
		)
	}
	if !block.WindowSizeValue.IsNull() &&
		!block.WindowSizeValue.IsUnknown() &&
		block.WindowSizeMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("window_size_value"),
			tfdiag.MissingConfigErrSummary,
			"window_size_match must be specified with window_size_value"+blockErrorSuffix,
		)
	}
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP) validateConfig(
	_ context.Context,
	pathRoot path.Path,
	blockErrorSuffix string,
	resp *resource.ValidateConfigResponse,
) {
	if !block.ChecksumValidateMatch.IsNull() &&
		!block.ChecksumValidateMatch.IsUnknown() &&
		block.ChecksumValidateValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("checksum_validate_match"),
			tfdiag.MissingConfigErrSummary,
			"checksum_validate_value must be specified with checksum_validate_match"+blockErrorSuffix,
		)
	}
	if !block.ChecksumValidateValue.IsNull() &&
		!block.ChecksumValidateValue.IsUnknown() &&
		block.ChecksumValidateMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("checksum_validate_value"),
			tfdiag.MissingConfigErrSummary,
			"checksum_validate_match must be specified with checksum_validate_value"+blockErrorSuffix,
		)
	}
	if !block.DataLengthMatch.IsNull() &&
		!block.DataLengthMatch.IsUnknown() &&
		block.DataLengthValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("data_length_match"),
			tfdiag.MissingConfigErrSummary,
			"data_length_value must be specified with data_length_match"+blockErrorSuffix,
		)
	}
	if !block.DataLengthValue.IsNull() &&
		!block.DataLengthValue.IsUnknown() &&
		block.DataLengthMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("data_length_value"),
			tfdiag.MissingConfigErrSummary,
			"data_length_match must be specified with data_length_value"+blockErrorSuffix,
		)
	}
	if !block.DestinationPortMatch.IsNull() &&
		!block.DestinationPortMatch.IsUnknown() &&
		block.DestinationPortValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_port_match"),
			tfdiag.MissingConfigErrSummary,
			"destination_port_value must be specified with destination_port_match"+blockErrorSuffix,
		)
	}
	if !block.DestinationPortValue.IsNull() &&
		!block.DestinationPortValue.IsUnknown() &&
		block.DestinationPortMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("destination_port_value"),
			tfdiag.MissingConfigErrSummary,
			"destination_port_match must be specified with destination_port_value"+blockErrorSuffix,
		)
	}
	if !block.SourcePortMatch.IsNull() &&
		!block.SourcePortMatch.IsUnknown() &&
		block.SourcePortValue.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_port_match"),
			tfdiag.MissingConfigErrSummary,
			"source_port_value must be specified with source_port_match"+blockErrorSuffix,
		)
	}
	if !block.SourcePortValue.IsNull() &&
		!block.SourcePortValue.IsUnknown() &&
		block.SourcePortMatch.IsNull() {
		resp.Diagnostics.AddAttributeError(
			pathRoot.AtName("source_port_value"),
			tfdiag.MissingConfigErrSummary,
			"source_port_match must be specified with source_port_value"+blockErrorSuffix,
		)
	}
}

func (rsc *securityIdpCustomAttack) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityIdpCustomAttackData
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
			attackExists, err := checkSecurityIdpCustomAttackExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if attackExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			attackExists, err := checkSecurityIdpCustomAttackExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !attackExists {
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

func (rsc *securityIdpCustomAttack) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityIdpCustomAttackData
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

func (rsc *securityIdpCustomAttack) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityIdpCustomAttackData
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

func (rsc *securityIdpCustomAttack) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityIdpCustomAttackData
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

func (rsc *securityIdpCustomAttack) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityIdpCustomAttackData

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

func checkSecurityIdpCustomAttackExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security idp custom-attack \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityIdpCustomAttackData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityIdpCustomAttackData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityIdpCustomAttackData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set security idp custom-attack \"" + rscData.Name.ValueString() + "\" "

	configSet := []string{
		setPrefix + "severity " + rscData.Severity.ValueString(),
	}

	if v := rscData.RecommendedAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"recommended-action "+v)
	}
	if !rscData.TimeBindingCount.IsNull() {
		configSet = append(configSet, setPrefix+"time-binding count "+
			utils.ConvI64toa(rscData.TimeBindingCount.ValueInt64()))
	}
	if v := rscData.TimeBindingScope.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"time-binding scope "+v)
	}

	if rscData.AttackTypeAnomaly != nil {
		configSet = append(configSet, rscData.AttackTypeAnomaly.configSet(setPrefix)...)
	}
	if rscData.AttackTypeChain != nil {
		blockSet, pathErr, err := rscData.AttackTypeChain.configSet(setPrefix, path.Root("attack_type_chain"))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.AttackTypeSignature != nil {
		blockSet, pathErr, err := rscData.AttackTypeSignature.configSet(setPrefix, path.Root("attack_type_signature"))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityIdpCustomAttackBlockAttackTypeAnomaly) configSet(setPrefix string) []string {
	configSet := []string{
		setPrefix + "attack-type anomaly service \"" + block.Service.ValueString() + "\"",
	}

	configSet = append(configSet,
		block.securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly.configSet(setPrefix)...)

	return configSet
}

func (block *securityIdpCustomAttackBlockAttackTypeChain) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, 100)
	setPrefix += "attack-type chain "

	if v := block.Expression.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"expression \""+v+"\"")
	}
	if block.Order.ValueBool() {
		configSet = append(configSet, setPrefix+"order")
	}
	if v := block.ProtocolBinding.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"protocol-binding "+v)
	}
	if block.Reset.ValueBool() {
		configSet = append(configSet, setPrefix+"reset")
	}
	if v := block.Scope.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"scope "+v)
	}

	memberName := make(map[string]struct{})
	for i, subBlock := range block.Member {
		name := subBlock.Name.ValueString()
		if _, ok := memberName[name]; ok {
			return configSet,
				path.Root("attack_type_chain").AtName("member").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple member blocks with the same name %q"+
					" in attack_type_chain block", name)
		}
		memberName[name] = struct{}{}

		blockSet, pathErr, err := subBlock.configSet(setPrefix, pathRoot.AtName("member").AtListIndex(i))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return configSet, path.Empty(), nil
}

func (block *securityIdpCustomAttackBlockAttackTypeChainBlockMember) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, 100)
	setPrefix += "member \"" + block.Name.ValueString() + "\" "

	if block.isEmpty() {
		return configSet,
			pathRoot.AtName("name"),
			fmt.Errorf("one of attack_type_anomaly or attack_type_signature must be specified"+
				" in member block %q in attack_type_chain block", block.Name.ValueString())
	}
	if block.AttackTypeAnomaly != nil && block.AttackTypeSignature != nil {
		return configSet,
			pathRoot.AtName("attack_type_anomaly"),
			fmt.Errorf("attack_type_anomaly and attack_type_signature cannot be configured together"+
				" in member block %q in attack_type_chain block", block.Name.ValueString())
	}

	if block.AttackTypeAnomaly != nil {
		configSet = append(configSet, block.AttackTypeAnomaly.configSet(setPrefix)...)
	}
	if block.AttackTypeSignature != nil {
		blockSet, pathErr, err := block.AttackTypeSignature.configSet(setPrefix, pathRoot.AtName("attack_type_signature"))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return configSet, path.Empty(), nil
}

func (block *securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly) configSet(setPrefix string) []string { //nolint:lll
	setPrefix += "attack-type anomaly "

	configSet := []string{
		setPrefix + "direction " + block.Direction.ValueString(),
		setPrefix + "test \"" + block.Test.ValueString() + "\"",
	}

	if v := block.Shellcode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"shellcode "+v)
	}

	return configSet
}

func (block *securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "attack-type signature "

	configSet := []string{
		setPrefix + "context \"" + block.Context.ValueString() + "\"",
		setPrefix + "direction " + block.Direction.ValueString(),
	}

	if block.Negate.ValueBool() {
		configSet = append(configSet, setPrefix+"negate")
	}
	if v := block.Pattern.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pattern \""+v+"\"")
	}
	if v := block.PatternPcre.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pattern-pcre \""+v+"\"")
	}
	if v := block.Regexp.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"regexp \""+v+"\"")
	}
	if v := block.Shellcode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"shellcode "+v)
	}

	if block.ProtocolIcmp != nil {
		if block.ProtocolIcmpv6 != nil {
			return configSet,
				pathRoot.AtName("protocol_icmp"),
				errors.New("protocol_icmp and protocol_icmpv6 cannot be configured together")
		}
		if block.ProtocolTCP != nil {
			return configSet,
				pathRoot.AtName("protocol_icmp"),
				errors.New("protocol_icmp and protocol_tcp cannot be configured together")
		}
		if block.ProtocolUDP != nil {
			return configSet,
				pathRoot.AtName("protocol_icmp"),
				errors.New("protocol_icmp and protocol_udp cannot be configured together")
		}
		if block.ProtocolIcmp.isEmpty() {
			return configSet,
				pathRoot.AtName("protocol_icmp").AtName("*"),
				errors.New("protocol_icmp block is empty")
		}

		configSet = append(configSet, block.ProtocolIcmp.configSet(setPrefix, false)...)
	}
	if block.ProtocolIcmpv6 != nil {
		if block.ProtocolTCP != nil {
			return configSet,
				pathRoot.AtName("protocol_icmpv6"),
				errors.New("protocol_icmpv6 and protocol_tcp cannot be configured together")
		}
		if block.ProtocolUDP != nil {
			return configSet,
				pathRoot.AtName("protocol_icmpv6"),
				errors.New("protocol_icmpv6 and protocol_udp cannot be configured together")
		}
		if block.ProtocolIcmpv6.isEmpty() {
			return configSet,
				pathRoot.AtName("protocol_icmpv6").AtName("*"),
				errors.New("protocol_icmpv6 block is empty")
		}

		configSet = append(configSet, block.ProtocolIcmpv6.configSet(setPrefix, true)...)
	}
	if block.ProtocolIPv4 != nil {
		if block.ProtocolIPv4.isEmpty() {
			return configSet,
				pathRoot.AtName("protocol_ipv4").AtName("*"),
				errors.New("protocol_ipv4 block is empty")
		}

		configSet = append(configSet, block.ProtocolIPv4.configSet(setPrefix)...)
	}
	if block.ProtocolIPv6 != nil {
		if block.ProtocolIPv6.isEmpty() {
			return configSet,
				pathRoot.AtName("protocol_ipv6").AtName("*"),
				errors.New("protocol_ipv6 block is empty")
		}

		configSet = append(configSet, block.ProtocolIPv6.configSet(setPrefix)...)
	}
	if block.ProtocolTCP != nil {
		if block.ProtocolUDP != nil {
			return configSet,
				pathRoot.AtName("protocol_tcp"),
				errors.New("protocol_tcp and protocol_udp cannot be configured together")
		}
		if block.ProtocolTCP.isEmpty() {
			return configSet,
				pathRoot.AtName("protocol_tcp").AtName("*"),
				errors.New("protocol_tcp block is empty")
		}

		configSet = append(configSet, block.ProtocolTCP.configSet(setPrefix)...)
	}
	if block.ProtocolUDP != nil {
		if block.ProtocolUDP.isEmpty() {
			return configSet,
				pathRoot.AtName("protocol_udp").AtName("*"),
				errors.New("protocol_udp block is empty")
		}

		configSet = append(configSet, block.ProtocolUDP.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (block *securityIdpCustomAttackBlockAttackTypeSignature) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet, pathErr, err := block.securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature.configSet(setPrefix, pathRoot) //nolint:lll
	if err != nil {
		return configSet, pathErr, err
	}

	if v := block.ProtocolBinding.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"attack-type signature protocol-binding "+v)
	}

	return configSet, path.Empty(), nil
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp) configSet(
	setPrefix string, v6 bool,
) []string {
	configSet := make([]string, 0, 100)
	switch v6 {
	case false:
		setPrefix += "protocol icmp "
	case true:
		setPrefix += "protocol icmpv6 "
	}

	if v := block.ChecksumValidateMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"checksum-validate match "+v)
	}
	if !block.ChecksumValidateValue.IsNull() {
		configSet = append(configSet, setPrefix+"checksum-validate value "+
			utils.ConvI64toa(block.ChecksumValidateValue.ValueInt64()))
	}
	if v := block.CodeMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"code match "+v)
	}
	if !block.CodeValue.IsNull() {
		configSet = append(configSet, setPrefix+"code value "+
			utils.ConvI64toa(block.CodeValue.ValueInt64()))
	}
	if v := block.DataLengthMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"data-length match "+v)
	}
	if !block.DataLengthValue.IsNull() {
		configSet = append(configSet, setPrefix+"data-length value "+
			utils.ConvI64toa(block.DataLengthValue.ValueInt64()))
	}
	if v := block.IdentificationMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"identification match "+v)
	}
	if !block.IdentificationValue.IsNull() {
		configSet = append(configSet, setPrefix+"identification value "+
			utils.ConvI64toa(block.IdentificationValue.ValueInt64()))
	}
	if v := block.SequenceNumberMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"sequence-number match "+v)
	}
	if !block.SequenceNumberValue.IsNull() {
		configSet = append(configSet, setPrefix+"sequence-number value "+
			utils.ConvI64toa(block.SequenceNumberValue.ValueInt64()))
	}
	if v := block.TypeMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"type match "+v)
	}
	if !block.TypeValue.IsNull() {
		configSet = append(configSet, setPrefix+"type value "+
			utils.ConvI64toa(block.TypeValue.ValueInt64()))
	}

	return configSet
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)
	setPrefix += "protocol ipv4 "

	if v := block.ChecksumValidateMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"checksum-validate match "+v)
	}
	if !block.ChecksumValidateValue.IsNull() {
		configSet = append(configSet, setPrefix+"checksum-validate value "+
			utils.ConvI64toa(block.ChecksumValidateValue.ValueInt64()))
	}
	if v := block.DestinationMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"destination match "+v)
	}
	if v := block.DestinationValue.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"destination value "+v)
	}
	if v := block.IdentificationMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"identification match "+v)
	}
	if !block.IdentificationValue.IsNull() {
		configSet = append(configSet, setPrefix+"identification value "+
			utils.ConvI64toa(block.IdentificationValue.ValueInt64()))
	}
	if v := block.IhlMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ihl match "+v)
	}
	if !block.IhlValue.IsNull() {
		configSet = append(configSet, setPrefix+"ihl value "+
			utils.ConvI64toa(block.IhlValue.ValueInt64()))
	}
	for _, v := range block.IPFlags {
		configSet = append(configSet, setPrefix+"ip-flags "+v.ValueString())
	}
	if v := block.ProtocolMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"protocol match "+v)
	}
	if !block.ProtocolValue.IsNull() {
		configSet = append(configSet, setPrefix+"protocol value "+
			utils.ConvI64toa(block.ProtocolValue.ValueInt64()))
	}
	if v := block.SourceMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source match "+v)
	}
	if v := block.SourceValue.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source value "+v)
	}
	if v := block.TosMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"tos match "+v)
	}
	if !block.TosValue.IsNull() {
		configSet = append(configSet, setPrefix+"tos value "+
			utils.ConvI64toa(block.TosValue.ValueInt64()))
	}
	if v := block.TotalLengthMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"total-length match "+v)
	}
	if !block.TotalLengthValue.IsNull() {
		configSet = append(configSet, setPrefix+"total-length value "+
			utils.ConvI64toa(block.TotalLengthValue.ValueInt64()))
	}
	if v := block.TTLMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ttl match "+v)
	}
	if !block.TTLValue.IsNull() {
		configSet = append(configSet, setPrefix+"ttl value "+
			utils.ConvI64toa(block.TTLValue.ValueInt64()))
	}

	return configSet
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv6) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)
	setPrefix += "protocol ipv6 "

	if v := block.DestinationMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"destination match "+v)
	}
	if v := block.DestinationValue.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"destination value "+v)
	}
	if v := block.ExtensionHeaderDestinationOptionHomeAddressMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"extension-header destination-option home-address match "+v)
	}
	if v := block.ExtensionHeaderDestinationOptionHomeAddressValue.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"extension-header destination-option home-address value "+v)
	}
	if v := block.ExtensionHeaderDestinationOptionTypeMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"extension-header destination-option option-type match "+v)
	}
	if !block.ExtensionHeaderDestinationOptionTypeValue.IsNull() {
		configSet = append(configSet, setPrefix+"extension-header destination-option option-type value "+
			utils.ConvI64toa(block.ExtensionHeaderDestinationOptionTypeValue.ValueInt64()))
	}
	if v := block.ExtensionHeaderRoutingHeaderTypeMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"extension-header routing-header header-type match "+v)
	}
	if !block.ExtensionHeaderRoutingHeaderTypeValue.IsNull() {
		configSet = append(configSet, setPrefix+"extension-header routing-header header-type value "+
			utils.ConvI64toa(block.ExtensionHeaderRoutingHeaderTypeValue.ValueInt64()))
	}
	if v := block.FlowLabelMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"flow-label match "+v)
	}
	if !block.FlowLabelValue.IsNull() {
		configSet = append(configSet, setPrefix+"flow-label value "+
			utils.ConvI64toa(block.FlowLabelValue.ValueInt64()))
	}
	if v := block.HopLimitMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"hop-limit match "+v)
	}
	if !block.HopLimitValue.IsNull() {
		configSet = append(configSet, setPrefix+"hop-limit value "+
			utils.ConvI64toa(block.HopLimitValue.ValueInt64()))
	}
	if v := block.NextHeaderMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"next-header match "+v)
	}
	if !block.NextHeaderValue.IsNull() {
		configSet = append(configSet, setPrefix+"next-header value "+
			utils.ConvI64toa(block.NextHeaderValue.ValueInt64()))
	}
	if v := block.PayloadLengthMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"payload-length match "+v)
	}
	if !block.PayloadLengthValue.IsNull() {
		configSet = append(configSet, setPrefix+"payload-length value "+
			utils.ConvI64toa(block.PayloadLengthValue.ValueInt64()))
	}
	if v := block.SourceMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source match "+v)
	}
	if v := block.SourceValue.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source value "+v)
	}
	if v := block.TrafficClassMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"traffic-class match "+v)
	}
	if !block.TrafficClassValue.IsNull() {
		configSet = append(configSet, setPrefix+"traffic-class value "+
			utils.ConvI64toa(block.TrafficClassValue.ValueInt64()))
	}

	return configSet
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCP) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)
	setPrefix += "protocol tcp "

	if v := block.AckNumberMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ack-number match "+v)
	}
	if !block.AckNumberValue.IsNull() {
		configSet = append(configSet, setPrefix+"ack-number value "+
			utils.ConvI64toa(block.AckNumberValue.ValueInt64()))
	}
	if v := block.ChecksumValidateMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"checksum-validate match "+v)
	}
	if !block.ChecksumValidateValue.IsNull() {
		configSet = append(configSet, setPrefix+"checksum-validate value "+
			utils.ConvI64toa(block.ChecksumValidateValue.ValueInt64()))
	}
	if v := block.DataLengthMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"data-length match "+v)
	}
	if !block.DataLengthValue.IsNull() {
		configSet = append(configSet, setPrefix+"data-length value "+
			utils.ConvI64toa(block.DataLengthValue.ValueInt64()))
	}
	if v := block.DestinationPortMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"destination-port match "+v)
	}
	if !block.DestinationPortValue.IsNull() {
		configSet = append(configSet, setPrefix+"destination-port value "+
			utils.ConvI64toa(block.DestinationPortValue.ValueInt64()))
	}
	if v := block.HeaderLengthMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"header-length match "+v)
	}
	if !block.HeaderLengthValue.IsNull() {
		configSet = append(configSet, setPrefix+"header-length value "+
			utils.ConvI64toa(block.HeaderLengthValue.ValueInt64()))
	}
	if v := block.MssMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"mss match "+v)
	}
	if !block.MssValue.IsNull() {
		configSet = append(configSet, setPrefix+"mss value "+
			utils.ConvI64toa(block.MssValue.ValueInt64()))
	}
	if v := block.OptionMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"option match "+v)
	}
	if !block.OptionValue.IsNull() {
		configSet = append(configSet, setPrefix+"option value "+
			utils.ConvI64toa(block.OptionValue.ValueInt64()))
	}
	if v := block.ReservedMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"reserved match "+v)
	}
	if !block.ReservedValue.IsNull() {
		configSet = append(configSet, setPrefix+"reserved value "+
			utils.ConvI64toa(block.ReservedValue.ValueInt64()))
	}
	if v := block.SequenceNumberMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"sequence-number match "+v)
	}
	if !block.SequenceNumberValue.IsNull() {
		configSet = append(configSet, setPrefix+"sequence-number value "+
			utils.ConvI64toa(block.SequenceNumberValue.ValueInt64()))
	}
	if v := block.SourcePortMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-port match "+v)
	}
	if !block.SourcePortValue.IsNull() {
		configSet = append(configSet, setPrefix+"source-port value "+
			utils.ConvI64toa(block.SourcePortValue.ValueInt64()))
	}
	for _, v := range block.TCPFlags {
		configSet = append(configSet, setPrefix+"tcp-flags "+v.ValueString())
	}
	if v := block.UrgentPointerMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"urgent-pointer match "+v)
	}
	if !block.UrgentPointerValue.IsNull() {
		configSet = append(configSet, setPrefix+"urgent-pointer value "+
			utils.ConvI64toa(block.UrgentPointerValue.ValueInt64()))
	}
	if v := block.WindowScaleMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"window-scale match "+v)
	}
	if !block.WindowScaleValue.IsNull() {
		configSet = append(configSet, setPrefix+"window-scale value "+
			utils.ConvI64toa(block.WindowScaleValue.ValueInt64()))
	}
	if v := block.WindowSizeMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"window-size match "+v)
	}
	if !block.WindowSizeValue.IsNull() {
		configSet = append(configSet, setPrefix+"window-size value "+
			utils.ConvI64toa(block.WindowSizeValue.ValueInt64()))
	}

	return configSet
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)
	setPrefix += "protocol udp "

	if v := block.ChecksumValidateMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"checksum-validate match "+v)
	}
	if !block.ChecksumValidateValue.IsNull() {
		configSet = append(configSet, setPrefix+"checksum-validate value "+
			utils.ConvI64toa(block.ChecksumValidateValue.ValueInt64()))
	}
	if v := block.DataLengthMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"data-length match "+v)
	}
	if !block.DataLengthValue.IsNull() {
		configSet = append(configSet, setPrefix+"data-length value "+
			utils.ConvI64toa(block.DataLengthValue.ValueInt64()))
	}
	if v := block.DestinationPortMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"destination-port match "+v)
	}
	if !block.DestinationPortValue.IsNull() {
		configSet = append(configSet, setPrefix+"destination-port value "+
			utils.ConvI64toa(block.DestinationPortValue.ValueInt64()))
	}
	if v := block.SourcePortMatch.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-port match "+v)
	}
	if !block.SourcePortValue.IsNull() {
		configSet = append(configSet, setPrefix+"source-port value "+
			utils.ConvI64toa(block.SourcePortValue.ValueInt64()))
	}

	return configSet
}

func (rscData *securityIdpCustomAttackData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security idp custom-attack \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "severity "):
				rscData.Severity = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "recommended-action "):
				rscData.RecommendedAction = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "time-binding count "):
				rscData.TimeBindingCount, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "time-binding scope "):
				rscData.TimeBindingScope = types.StringValue(itemTrim)

			case balt.CutPrefixInString(&itemTrim, "attack-type anomaly "):
				if rscData.AttackTypeAnomaly == nil {
					rscData.AttackTypeAnomaly = &securityIdpCustomAttackBlockAttackTypeAnomaly{}
				}
				rscData.AttackTypeAnomaly.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "attack-type chain "):
				if rscData.AttackTypeChain == nil {
					rscData.AttackTypeChain = &securityIdpCustomAttackBlockAttackTypeChain{}
				}
				if err := rscData.AttackTypeChain.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "attack-type signature "):
				if rscData.AttackTypeSignature == nil {
					rscData.AttackTypeSignature = &securityIdpCustomAttackBlockAttackTypeSignature{}
				}
				if err := rscData.AttackTypeSignature.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *securityIdpCustomAttackBlockAttackTypeAnomaly) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "service "):
		block.Service = types.StringValue(strings.Trim(itemTrim, "\""))
	default:
		block.securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly.read(itemTrim)
	}
}

func (block *securityIdpCustomAttackBlockAttackTypeChain) read(itemTrim string) error {
	switch {
	case balt.CutPrefixInString(&itemTrim, "expression "):
		block.Expression = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "order":
		block.Order = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "protocol-binding "):
		block.ProtocolBinding = types.StringValue(itemTrim)
	case itemTrim == "reset":
		block.Reset = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "scope "):
		block.Scope = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "member "):
		name := tfdata.FirstElementOfJunosLine(itemTrim)
		block.Member = tfdata.AppendPotentialNewBlock(block.Member, types.StringValue(strings.Trim(name, "\"")))
		member := &block.Member[len(block.Member)-1]
		balt.CutPrefixInString(&itemTrim, name+" ")

		switch {
		case balt.CutPrefixInString(&itemTrim, "attack-type anomaly "):
			if member.AttackTypeAnomaly == nil {
				member.AttackTypeAnomaly = &securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly{}
			}

			member.AttackTypeAnomaly.read(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "attack-type signature "):
			if member.AttackTypeSignature == nil {
				member.AttackTypeSignature = &securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature{}
			}

			if err := member.AttackTypeSignature.read(itemTrim); err != nil {
				return err
			}
		}
	}

	return nil
}

func (block *securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "direction "):
		block.Direction = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "test "):
		block.Test = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "shellcode "):
		block.Shellcode = types.StringValue(itemTrim)
	}
}

func (block *securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature) read(itemTrim string) error { //nolint:lll
	switch {
	case balt.CutPrefixInString(&itemTrim, "context "):
		block.Context = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "direction "):
		block.Direction = types.StringValue(itemTrim)
	case itemTrim == "negate":
		block.Negate = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "pattern "):
		block.Pattern = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "pattern-pcre "):
		block.PatternPcre = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "regexp "):
		block.Regexp = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "shellcode "):
		block.Shellcode = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "protocol icmp "):
		if block.ProtocolIcmp == nil {
			block.ProtocolIcmp = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp{}
		}

		if err := block.ProtocolIcmp.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "protocol icmpv6 "):
		if block.ProtocolIcmpv6 == nil {
			block.ProtocolIcmpv6 = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp{}
		}

		if err := block.ProtocolIcmpv6.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "protocol ipv4 "):
		if block.ProtocolIPv4 == nil {
			block.ProtocolIPv4 = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4{}
		}

		if err := block.ProtocolIPv4.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "protocol ipv6 "):
		if block.ProtocolIPv6 == nil {
			block.ProtocolIPv6 = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv6{}
		}

		if err := block.ProtocolIPv6.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "protocol tcp "):
		if block.ProtocolTCP == nil {
			block.ProtocolTCP = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCP{}
		}

		if err := block.ProtocolTCP.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "protocol udp "):
		if block.ProtocolUDP == nil {
			block.ProtocolUDP = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP{}
		}

		if err := block.ProtocolUDP.read(itemTrim); err != nil {
			return err
		}
	}

	return nil
}

func (block *securityIdpCustomAttackBlockAttackTypeSignature) read(itemTrim string) error {
	switch {
	case balt.CutPrefixInString(&itemTrim, "protocol-binding "):
		block.ProtocolBinding = types.StringValue(itemTrim)
	default:
		if err := block.securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature.read(itemTrim); err != nil { //nolint:lll
			return err
		}
	}

	return nil
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "checksum-validate match "):
		block.ChecksumValidateMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "checksum-validate value "):
		block.ChecksumValidateValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "code match "):
		block.CodeMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "code value "):
		block.CodeValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "data-length match "):
		block.DataLengthMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "data-length value "):
		block.DataLengthValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "identification match "):
		block.IdentificationMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "identification value "):
		block.IdentificationValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "sequence-number match "):
		block.SequenceNumberMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "sequence-number value "):
		block.SequenceNumberValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "type match "):
		block.TypeMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "type value "):
		block.TypeValue, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "checksum-validate match "):
		block.ChecksumValidateMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "checksum-validate value "):
		block.ChecksumValidateValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination match "):
		block.DestinationMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination value "):
		block.DestinationValue = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "identification match "):
		block.IdentificationMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "identification value "):
		block.IdentificationValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ihl match "):
		block.IhlMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ihl value "):
		block.IhlValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ip-flags "):
		block.IPFlags = append(block.IPFlags, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "protocol match "):
		block.ProtocolMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "protocol value "):
		block.ProtocolValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source match "):
		block.SourceMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source value "):
		block.SourceValue = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "tos match "):
		block.TosMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "tos value "):
		block.TosValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "total-length match "):
		block.TotalLengthMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "total-length value "):
		block.TotalLengthValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ttl match "):
		block.TTLMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ttl value "):
		block.TTLValue, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv6) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "destination match "):
		block.DestinationMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination value "):
		block.DestinationValue = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "extension-header destination-option home-address match "):
		block.ExtensionHeaderDestinationOptionHomeAddressMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "extension-header destination-option home-address value "):
		block.ExtensionHeaderDestinationOptionHomeAddressValue = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "extension-header destination-option option-type match "):
		block.ExtensionHeaderDestinationOptionTypeMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "extension-header destination-option option-type value "):
		block.ExtensionHeaderDestinationOptionTypeValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "extension-header routing-header header-type match "):
		block.ExtensionHeaderRoutingHeaderTypeMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "extension-header routing-header header-type value "):
		block.ExtensionHeaderRoutingHeaderTypeValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "flow-label match "):
		block.FlowLabelMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "flow-label value "):
		block.FlowLabelValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "hop-limit match "):
		block.HopLimitMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "hop-limit value "):
		block.HopLimitValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "next-header match "):
		block.NextHeaderMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "next-header value "):
		block.NextHeaderValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "payload-length match "):
		block.PayloadLengthMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "payload-length value "):
		block.PayloadLengthValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source match "):
		block.SourceMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source value "):
		block.SourceValue = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "traffic-class match "):
		block.TrafficClassMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "traffic-class value "):
		block.TrafficClassValue, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCP) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "ack-number match "):
		block.AckNumberMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ack-number value "):
		block.AckNumberValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "checksum-validate match "):
		block.ChecksumValidateMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "checksum-validate value "):
		block.ChecksumValidateValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "data-length match "):
		block.DataLengthMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "data-length value "):
		block.DataLengthValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination-port match "):
		block.DestinationPortMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination-port value "):
		block.DestinationPortValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "header-length match "):
		block.HeaderLengthMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "header-length value "):
		block.HeaderLengthValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "mss match "):
		block.MssMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "mss value "):
		block.MssValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "option match "):
		block.OptionMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "option value "):
		block.OptionValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "reserved match "):
		block.ReservedMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "reserved value "):
		block.ReservedValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "sequence-number match "):
		block.SequenceNumberMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "sequence-number value "):
		block.SequenceNumberValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-port match "):
		block.SourcePortMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-port value "):
		block.SourcePortValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "tcp-flags "):
		block.TCPFlags = append(block.TCPFlags, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "urgent-pointer match "):
		block.UrgentPointerMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "urgent-pointer value "):
		block.UrgentPointerValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "window-scale match "):
		block.WindowScaleMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "window-scale value "):
		block.WindowScaleValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "window-size match "):
		block.WindowSizeMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "window-size value "):
		block.WindowSizeValue, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (block *securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "checksum-validate match "):
		block.ChecksumValidateMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "checksum-validate value "):
		block.ChecksumValidateValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "data-length match "):
		block.DataLengthMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "data-length value "):
		block.DataLengthValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination-port match "):
		block.DestinationPortMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination-port value "):
		block.DestinationPortValue, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-port match "):
		block.SourcePortMatch = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-port value "):
		block.SourcePortValue, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (rscData *securityIdpCustomAttackData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security idp custom-attack \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
