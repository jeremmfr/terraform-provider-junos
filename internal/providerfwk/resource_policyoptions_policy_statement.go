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
	_ resource.Resource                   = &policyoptionsPolicyStatement{}
	_ resource.ResourceWithConfigure      = &policyoptionsPolicyStatement{}
	_ resource.ResourceWithValidateConfig = &policyoptionsPolicyStatement{}
	_ resource.ResourceWithImportState    = &policyoptionsPolicyStatement{}
	_ resource.ResourceWithUpgradeState   = &policyoptionsPolicyStatement{}
)

type policyoptionsPolicyStatement struct {
	client *junos.Client
}

func newPolicyoptionsPolicyStatementResource() resource.Resource {
	return &policyoptionsPolicyStatement{}
}

func (rsc *policyoptionsPolicyStatement) typeName() string {
	return providerName + "_policyoptions_policy_statement"
}

func (rsc *policyoptionsPolicyStatement) junosName() string {
	return "policy-options policy-statement"
}

func (rsc *policyoptionsPolicyStatement) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *policyoptionsPolicyStatement) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *policyoptionsPolicyStatement) Configure(
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

func (rsc *policyoptionsPolicyStatement) Schema(
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
				Description: "Name to identify the policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"add_it_to_forwarding_table_export": schema.BoolAttribute{
				Optional:    true,
				Description: "Add this policy in `routing-options forwarding-table export` list.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"dynamic_db": schema.BoolAttribute{
				Optional:    true,
				Description: "Object may exist in dynamic database.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"from": schema.SingleNestedBlock{
				Description: "Conditions to match the source of a route.",
				Attributes:  rsc.schemaFromAttributes(),
				Blocks:      rsc.schemaFromBlocks(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"to": schema.SingleNestedBlock{
				Description: "Conditions to match the destination of a route.",
				Attributes:  rsc.schemaToAttributes(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"then": schema.SingleNestedBlock{
				Description: "Actions to take if 'from' and 'to' conditions match.",
				Attributes:  rsc.schemaThenAttributes(),
				Blocks:      rsc.schemaThenBlocks(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"term": schema.ListNestedBlock{
				Description: "For each policy term.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of term.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 250),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"from": schema.SingleNestedBlock{
							Description: "Conditions to match the source of a route.",
							Attributes:  rsc.schemaFromAttributes(),
							Blocks:      rsc.schemaFromBlocks(),
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
						"to": schema.SingleNestedBlock{
							Description: "Conditions to match the destination of a route.",
							Attributes:  rsc.schemaToAttributes(),
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
						"then": schema.SingleNestedBlock{
							Description: "Actions to take if 'from' and 'to' conditions match",
							Attributes:  rsc.schemaThenAttributes(),
							Blocks:      rsc.schemaThenBlocks(),
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
					},
				},
			},
		},
	}
}

func (rsc *policyoptionsPolicyStatement) schemaFromAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"aggregate_contributor": schema.BoolAttribute{
			Optional:    true,
			Description: "Match more specifics of an aggregate.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"bgp_as_path": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Name of AS path regular expression.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"bgp_as_path_group": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Name of AS path group.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"bgp_community": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "BGP community.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"bgp_origin": schema.StringAttribute{
			Optional:    true,
			Description: "BGP origin attribute.",
			Validators: []validator.String{
				stringvalidator.OneOf("egp", "igp", "incomplete"),
			},
		},
		"bgp_srte_discriminator": schema.Int64Attribute{
			Optional:    true,
			Description: "Srte discriminator.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"color": schema.Int64Attribute{
			Optional:    true,
			Description: "Color (preference) value.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"evpn_esi": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "ESI in EVPN Route.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^([\d\w]{2}:){9}[\d\w]{2}$`),
						"bad format or length"),
				),
			},
		},
		"evpn_mac_route": schema.StringAttribute{
			Optional:    true,
			Description: "EVPN Mac Route type.",
			Validators: []validator.String{
				stringvalidator.OneOf("mac-ipv4", "mac-ipv6", "mac-only"),
			},
		},
		"evpn_tag": schema.SetAttribute{
			ElementType: types.Int64Type,
			Optional:    true,
			Description: "Tag in EVPN Route (0..4294967295).",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueInt64sAre(
					int64validator.Between(0, 4294967295),
				),
			},
		},
		"family": schema.StringAttribute{
			Optional:    true,
			Description: "Family.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					"evpn", "inet", "inet-mdt", "inet-mvpn", "inet-vpn",
					"inet6", "inet6-mvpn", "inet6-vpn", "iso",
				),
			},
		},
		"local_preference": schema.Int64Attribute{
			Optional:    true,
			Description: "Local preference associated with a route.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"interface": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Interface name or address.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
					stringvalidator.Any(
						tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
						tfvalidator.StringIPAddress(),
					),
				),
			},
		},
		"metric": schema.Int64Attribute{
			Optional:    true,
			Description: "Metric value.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"neighbor": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Neighboring router.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					tfvalidator.StringIPAddress(),
				),
			},
		},
		"next_hop": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Next-hop router.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					tfvalidator.StringIPAddress(),
				),
			},
		},
		"next_hop_type_merged": schema.BoolAttribute{
			Optional:    true,
			Description: "Merged next hop.",
			Validators: []validator.Bool{
				tfvalidator.BoolTrue(),
			},
		},
		"ospf_area": schema.StringAttribute{
			Optional:    true,
			Description: "OSPF area identifier.",
			Validators: []validator.String{
				tfvalidator.StringIPAddress().IPv4Only(),
			},
		},
		"policy": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Name of policy to evaluate.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"preference": schema.Int64Attribute{
			Optional:    true,
			Description: "Preference value.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"prefix_list": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Prefix-lists of routes to match.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"protocol": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Protocol from which route was learned.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				),
			},
		},
		"route_type": schema.StringAttribute{
			Optional:    true,
			Description: "Route type.",
			Validators: []validator.String{
				stringvalidator.OneOf("external", "internal"),
			},
		},
		"routing_instance": schema.StringAttribute{
			Optional:    true,
			Description: "Routing protocol instance.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 63),
				tfvalidator.StringFormat(tfvalidator.DefaultFormat),
			},
		},
		"srte_color": schema.Int64Attribute{
			Optional:    true,
			Description: "Srte color.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"state": schema.StringAttribute{
			Optional:    true,
			Description: "Route state.",
			Validators: []validator.String{
				stringvalidator.OneOf("active", "inactive"),
			},
		},
		"tunnel_type": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Tunnel type.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf("gre", "ipip", "udp"),
				),
			},
		},
		"validation_database": schema.StringAttribute{
			Optional:    true,
			Description: "Name to identify a validation-state.",
			Validators: []validator.String{
				stringvalidator.OneOf("invalid", "unknown", "valid"),
			},
		},
	}
}

func (rsc *policyoptionsPolicyStatement) schemaFromBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"bgp_as_path_calc_length": schema.SetNestedBlock{
			Description: "Number of BGP ASes excluding confederations.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"count": schema.Int64Attribute{
						Required:    true,
						Description: "Number of ASes (0..1024).",
						Validators: []validator.Int64{
							int64validator.Between(0, 1024),
						},
					},
					"match": schema.StringAttribute{
						Required:    true,
						Description: "Type of match: equal values, higher or equal values, lower or equal values.",
						Validators: []validator.String{
							stringvalidator.OneOf("equal", "orhigher", "orlower"),
						},
					},
				},
			},
		},
		"bgp_as_path_unique_count": schema.SetNestedBlock{
			Description: "Number of unique BGP ASes excluding confederations.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"count": schema.Int64Attribute{
						Required:    true,
						Description: "Number of ASes (0..1024).",
						Validators: []validator.Int64{
							int64validator.Between(0, 1024),
						},
					},
					"match": schema.StringAttribute{
						Required:    true,
						Description: "Type of match: equal values, higher or equal values, lower or equal values.",
						Validators: []validator.String{
							stringvalidator.OneOf("equal", "orhigher", "orlower"),
						},
					},
				},
			},
		},
		"bgp_community_count": schema.SetNestedBlock{
			Description: "Number of BGP communities.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"count": schema.Int64Attribute{
						Required:    true,
						Description: "Number of communities (0..1024).",
						Validators: []validator.Int64{
							int64validator.Between(0, 1024),
						},
					},
					"match": schema.StringAttribute{
						Required:    true,
						Description: "Type of match: equal values, higher or equal values, lower or equal values.",
						Validators: []validator.String{
							stringvalidator.OneOf("equal", "orhigher", "orlower"),
						},
					},
				},
			},
		},
		"next_hop_weight": schema.SetNestedBlock{
			Description: "Weight of the gateway.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"match": schema.StringAttribute{
						Required:    true,
						Description: "Type of match for weight.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"equal", "greater-than", "greater-than-equal", "less-than", "less-than-equal",
							),
						},
					},
					"weight": schema.Int64Attribute{
						Required:    true,
						Description: "Weight of the gateway (1..65535).",
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
				},
			},
		},
		"route_filter": schema.ListNestedBlock{
			Description: "Routes to match.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"route": schema.StringAttribute{
						Required:    true,
						Description: "IP address.",
						Validators: []validator.String{
							tfvalidator.StringCIDRNetwork(),
						},
					},
					"option": schema.StringAttribute{
						Required:    true,
						Description: "Mask option.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"address-mask", "exact", "longer", "orlonger", "prefix-length-range", "through", "upto",
							),
						},
					},
					"option_value": schema.StringAttribute{
						Optional:    true,
						Description: "For options that need an argument.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
		},
	}
}

func (rsc *policyoptionsPolicyStatement) schemaToAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"bgp_as_path": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Name of AS path regular expression.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"bgp_as_path_group": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Name of AS path group.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"bgp_community": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "BGP community.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"bgp_origin": schema.StringAttribute{
			Optional:    true,
			Description: "BGP origin attribute.",
			Validators: []validator.String{
				stringvalidator.OneOf("egp", "igp", "incomplete"),
			},
		},
		"family": schema.StringAttribute{
			Optional:    true,
			Description: "Family.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					"evpn", "inet", "inet-mdt", "inet-mvpn", "inet-vpn",
					"inet6", "inet6-mvpn", "inet6-vpn", "iso",
				),
			},
		},
		"local_preference": schema.Int64Attribute{
			Optional:    true,
			Description: "Local preference associated with a route.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"interface": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Interface name or address.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
					stringvalidator.Any(
						tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
						tfvalidator.StringIPAddress(),
					),
				),
			},
		},
		"metric": schema.Int64Attribute{
			Optional:    true,
			Description: "Metric value.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"neighbor": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Neighboring router.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					tfvalidator.StringIPAddress(),
				),
			},
		},
		"next_hop": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Next-hop router.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					tfvalidator.StringIPAddress(),
				),
			},
		},
		"ospf_area": schema.StringAttribute{
			Optional:    true,
			Description: "OSPF area identifier.",
			Validators: []validator.String{
				tfvalidator.StringIPAddress().IPv4Only(),
			},
		},
		"policy": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Name of policy to evaluate.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"preference": schema.Int64Attribute{
			Optional:    true,
			Description: "Preference value.",
			Validators: []validator.Int64{
				int64validator.Between(0, 4294967295),
			},
		},
		"protocol": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Protocol from which route was learned.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				),
			},
		},
		"routing_instance": schema.StringAttribute{
			Optional:    true,
			Description: "Routing protocol instance.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 63),
				tfvalidator.StringFormat(tfvalidator.DefaultFormat),
			},
		},
	}
}

func (rsc *policyoptionsPolicyStatement) schemaThenAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"action": schema.StringAttribute{
			Optional:    true,
			Description: "Action `accept` or `reject`.",
			Validators: []validator.String{
				stringvalidator.OneOf("accept", "reject"),
			},
		},
		"as_path_expand": schema.StringAttribute{
			Optional:    true,
			Description: "Prepend AS numbers prior to adding local-as.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"as_path_prepend": schema.StringAttribute{
			Optional:    true,
			Description: "Prepend AS numbers to an AS path.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
		"default_action": schema.StringAttribute{
			Optional:    true,
			Description: "Set default policy action.",
			Validators: []validator.String{
				stringvalidator.OneOf("accept", "reject"),
			},
		},
		"load_balance": schema.StringAttribute{
			Optional:    true,
			Description: "Type of load balancing in forwarding table.",
			Validators: []validator.String{
				stringvalidator.OneOf("per-packet", "consistent-hash"),
			},
		},
		"next": schema.StringAttribute{
			Optional:    true,
			Description: "Skip to next `policy` or `term`.",
			Validators: []validator.String{
				stringvalidator.OneOf("policy", "term"),
			},
		},
		"next_hop": schema.StringAttribute{
			Optional:    true,
			Description: "Set the address of the next-hop router.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.Any(
					tfvalidator.StringIPAddress(),
					stringvalidator.OneOf("discard", "next-table", "peer-address", "reject", "self"),
				),
			},
		},
		"origin": schema.StringAttribute{
			Optional:    true,
			Description: "BGP path origin.",
			Validators: []validator.String{
				stringvalidator.OneOf("egp", "igp", "incomplete"),
			},
		},
	}
}

func (rsc *policyoptionsPolicyStatement) schemaThenBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"community": schema.ListNestedBlock{
			Description: "For each community action.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"action": schema.StringAttribute{
						Required:    true,
						Description: "Action on BGP community.",
						Validators: []validator.String{
							stringvalidator.OneOf("add", "delete", "set"),
						},
					},
					"value": schema.StringAttribute{
						Required:    true,
						Description: "Name to identify a BGP community.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 250),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
			},
		},
		"local_preference": schema.SingleNestedBlock{
			Description: "Declare local-preference action.",
			Attributes: map[string]schema.Attribute{
				"action": schema.StringAttribute{
					Required:    false, // true when SingleNestedBlock is specified
					Optional:    true,
					Description: "Action on local-preference.",
					Validators: []validator.String{
						stringvalidator.OneOf("add", "subtract", "none"),
					},
				},
				"value": schema.Int64Attribute{
					Required:    false, // true when SingleNestedBlock is specified
					Optional:    true,
					Description: "Value for action (local-preference, constant).",
					Validators: []validator.Int64{
						int64validator.Between(0, 4294967295),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"metric": schema.SingleNestedBlock{
			Description: "Declare metric action.",
			Attributes: map[string]schema.Attribute{
				"action": schema.StringAttribute{
					Required:    false, // true when SingleNestedBlock is specified
					Optional:    true,
					Description: "Action on metric.",
					Validators: []validator.String{
						stringvalidator.OneOf("add", "subtract", "none"),
					},
				},
				"value": schema.Int64Attribute{
					Required:    false, // true when SingleNestedBlock is specified
					Optional:    true,
					Description: "Value for action (metric, constant).",
					Validators: []validator.Int64{
						int64validator.Between(0, 4294967295),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
		"preference": schema.SingleNestedBlock{
			Description: "Declare preference action.",
			Attributes: map[string]schema.Attribute{
				"action": schema.StringAttribute{
					Required:    false, // true when SingleNestedBlock is specified
					Optional:    true,
					Description: "Action on preference.",
					Validators: []validator.String{
						stringvalidator.OneOf("add", "subtract", "none"),
					},
				},
				"value": schema.Int64Attribute{
					Required:    false, // true when SingleNestedBlock is specified
					Optional:    true,
					Description: "Value for action (preference, constant).",
					Validators: []validator.Int64{
						int64validator.Between(0, 4294967295),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				tfplanmodifier.BlockRemoveNull(),
			},
		},
	}
}

type policyoptionsPolicyStatementData struct {
	ID                           types.String                            `tfsdk:"id"`
	Name                         types.String                            `tfsdk:"name"`
	AddItToForwardingTableExport types.Bool                              `tfsdk:"add_it_to_forwarding_table_export"`
	DynamicDB                    types.Bool                              `tfsdk:"dynamic_db"`
	From                         *policyoptionsPolicyStatementBlockFrom  `tfsdk:"from"`
	To                           *policyoptionsPolicyStatementBlockTo    `tfsdk:"to"`
	Then                         *policyoptionsPolicyStatementBlockThen  `tfsdk:"then"`
	Term                         []policyoptionsPolicyStatementBlockTerm `tfsdk:"term"`
}

type policyoptionsPolicyStatementConfig struct {
	ID                           types.String                                 `tfsdk:"id"`
	Name                         types.String                                 `tfsdk:"name"`
	AddItToForwardingTableExport types.Bool                                   `tfsdk:"add_it_to_forwarding_table_export"`
	DynamicDB                    types.Bool                                   `tfsdk:"dynamic_db"`
	From                         *policyoptionsPolicyStatementBlockFromConfig `tfsdk:"from"`
	To                           *policyoptionsPolicyStatementBlockToConfig   `tfsdk:"to"`
	Then                         *policyoptionsPolicyStatementBlockThenConfig `tfsdk:"then"`
	Term                         types.List                                   `tfsdk:"term"`
}

type policyoptionsPolicyStatementBlockTerm struct {
	Name types.String                           `tfsdk:"name"`
	From *policyoptionsPolicyStatementBlockFrom `tfsdk:"from"`
	To   *policyoptionsPolicyStatementBlockTo   `tfsdk:"to"`
	Then *policyoptionsPolicyStatementBlockThen `tfsdk:"then"`
}

func (block *policyoptionsPolicyStatementBlockTerm) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block, "Name")
}

type policyoptionsPolicyStatementBlockTermConfig struct {
	Name types.String                                 `tfsdk:"name"`
	From *policyoptionsPolicyStatementBlockFromConfig `tfsdk:"from"`
	To   *policyoptionsPolicyStatementBlockToConfig   `tfsdk:"to"`
	Then *policyoptionsPolicyStatementBlockThenConfig `tfsdk:"then"`
}

func (block *policyoptionsPolicyStatementBlockTermConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block, "Name")
}

type policyoptionsPolicyStatementBlockFrom struct {
	AggregateContributor types.Bool                                              `tfsdk:"aggregate_contributor"`
	BgpASPath            []types.String                                          `tfsdk:"bgp_as_path"`
	BgpASPathGroup       []types.String                                          `tfsdk:"bgp_as_path_group"`
	BgpCommunity         []types.String                                          `tfsdk:"bgp_community"`
	BgpOrigin            types.String                                            `tfsdk:"bgp_origin"`
	BgpSrteDiscriminator types.Int64                                             `tfsdk:"bgp_srte_discriminator"`
	Color                types.Int64                                             `tfsdk:"color"`
	EvpnESI              []types.String                                          `tfsdk:"evpn_esi"`
	EvpnMACRoute         types.String                                            `tfsdk:"evpn_mac_route"`
	EvpnTag              []types.Int64                                           `tfsdk:"evpn_tag"`
	Family               types.String                                            `tfsdk:"family"`
	LocalPreference      types.Int64                                             `tfsdk:"local_preference"`
	Interface            []types.String                                          `tfsdk:"interface"`
	Metric               types.Int64                                             `tfsdk:"metric"`
	Neighbor             []types.String                                          `tfsdk:"neighbor"`
	NextHop              []types.String                                          `tfsdk:"next_hop"`
	NextHopTypeMerged    types.Bool                                              `tfsdk:"next_hop_type_merged"`
	OspfArea             types.String                                            `tfsdk:"ospf_area"`
	Policy               []types.String                                          `tfsdk:"policy"`
	Preference           types.Int64                                             `tfsdk:"preference"`
	PrefixList           []types.String                                          `tfsdk:"prefix_list"`
	Protocol             []types.String                                          `tfsdk:"protocol"`
	RouteType            types.String                                            `tfsdk:"route_type"`
	RoutingInstance      types.String                                            `tfsdk:"routing_instance"`
	SrteColor            types.Int64                                             `tfsdk:"srte_color"`
	State                types.String                                            `tfsdk:"state"`
	TunnelType           []types.String                                          `tfsdk:"tunnel_type"`
	ValidationDatabase   types.String                                            `tfsdk:"validation_database"`
	BgpASPathCalcLength  []policyoptionsPolicyStatementBlockFromBlockCountMatch  `tfsdk:"bgp_as_path_calc_length"`
	BgpASPathUniqueCount []policyoptionsPolicyStatementBlockFromBlockCountMatch  `tfsdk:"bgp_as_path_unique_count"`
	BgpCommunityCount    []policyoptionsPolicyStatementBlockFromBlockCountMatch  `tfsdk:"bgp_community_count"`
	NextHopWeight        []policyoptionsPolicyStatementBlockFromBlockMatchWeight `tfsdk:"next_hop_weight"`
	RouteFilter          []policyoptionsPolicyStatementBlockFromBlockRouteFilter `tfsdk:"route_filter"`
}

func (block *policyoptionsPolicyStatementBlockFrom) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type policyoptionsPolicyStatementBlockFromConfig struct {
	AggregateContributor types.Bool   `tfsdk:"aggregate_contributor"`
	BgpASPath            types.Set    `tfsdk:"bgp_as_path"`
	BgpASPathGroup       types.Set    `tfsdk:"bgp_as_path_group"`
	BgpCommunity         types.Set    `tfsdk:"bgp_community"`
	BgpOrigin            types.String `tfsdk:"bgp_origin"`
	BgpSrteDiscriminator types.Int64  `tfsdk:"bgp_srte_discriminator"`
	Color                types.Int64  `tfsdk:"color"`
	EvpnESI              types.Set    `tfsdk:"evpn_esi"`
	EvpnMACRoute         types.String `tfsdk:"evpn_mac_route"`
	EvpnTag              types.Set    `tfsdk:"evpn_tag"`
	Family               types.String `tfsdk:"family"`
	LocalPreference      types.Int64  `tfsdk:"local_preference"`
	Interface            types.Set    `tfsdk:"interface"`
	Metric               types.Int64  `tfsdk:"metric"`
	Neighbor             types.Set    `tfsdk:"neighbor"`
	NextHop              types.Set    `tfsdk:"next_hop"`
	NextHopTypeMerged    types.Bool   `tfsdk:"next_hop_type_merged"`
	OspfArea             types.String `tfsdk:"ospf_area"`
	Policy               types.List   `tfsdk:"policy"`
	Preference           types.Int64  `tfsdk:"preference"`
	PrefixList           types.Set    `tfsdk:"prefix_list"`
	Protocol             types.Set    `tfsdk:"protocol"`
	RouteType            types.String `tfsdk:"route_type"`
	RoutingInstance      types.String `tfsdk:"routing_instance"`
	SrteColor            types.Int64  `tfsdk:"srte_color"`
	State                types.String `tfsdk:"state"`
	TunnelType           types.Set    `tfsdk:"tunnel_type"`
	ValidationDatabase   types.String `tfsdk:"validation_database"`
	BgpASPathCalcLength  types.Set    `tfsdk:"bgp_as_path_calc_length"`
	BgpASPathUniqueCount types.Set    `tfsdk:"bgp_as_path_unique_count"`
	BgpCommunityCount    types.Set    `tfsdk:"bgp_community_count"`
	NextHopWeight        types.Set    `tfsdk:"next_hop_weight"`
	RouteFilter          types.List   `tfsdk:"route_filter"`
}

func (block *policyoptionsPolicyStatementBlockFromConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type policyoptionsPolicyStatementBlockFromBlockCountMatch struct {
	Count types.Int64  `tfsdk:"count"`
	Match types.String `tfsdk:"match"`
}

type policyoptionsPolicyStatementBlockFromBlockMatchWeight struct {
	Match  types.String `tfsdk:"match"`
	Weight types.Int64  `tfsdk:"weight"`
}

type policyoptionsPolicyStatementBlockFromBlockRouteFilter struct {
	Route       types.String `tfsdk:"route"`
	Option      types.String `tfsdk:"option"`
	OptionValue types.String `tfsdk:"option_value"`
}

type policyoptionsPolicyStatementBlockTo struct {
	BgpASPath       []types.String `tfsdk:"bgp_as_path"`
	BgpASPathGroup  []types.String `tfsdk:"bgp_as_path_group"`
	BgpCommunity    []types.String `tfsdk:"bgp_community"`
	BgpOrigin       types.String   `tfsdk:"bgp_origin"`
	Family          types.String   `tfsdk:"family"`
	LocalPreference types.Int64    `tfsdk:"local_preference"`
	Interface       []types.String `tfsdk:"interface"`
	Metric          types.Int64    `tfsdk:"metric"`
	Neighbor        []types.String `tfsdk:"neighbor"`
	NextHop         []types.String `tfsdk:"next_hop"`
	OspfArea        types.String   `tfsdk:"ospf_area"`
	Policy          []types.String `tfsdk:"policy"`
	Preference      types.Int64    `tfsdk:"preference"`
	Protocol        []types.String `tfsdk:"protocol"`
	RoutingInstance types.String   `tfsdk:"routing_instance"`
}

func (block *policyoptionsPolicyStatementBlockTo) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type policyoptionsPolicyStatementBlockToConfig struct {
	BgpASPath       types.Set    `tfsdk:"bgp_as_path"`
	BgpASPathGroup  types.Set    `tfsdk:"bgp_as_path_group"`
	BgpCommunity    types.Set    `tfsdk:"bgp_community"`
	BgpOrigin       types.String `tfsdk:"bgp_origin"`
	Family          types.String `tfsdk:"family"`
	LocalPreference types.Int64  `tfsdk:"local_preference"`
	Interface       types.Set    `tfsdk:"interface"`
	Metric          types.Int64  `tfsdk:"metric"`
	Neighbor        types.Set    `tfsdk:"neighbor"`
	NextHop         types.Set    `tfsdk:"next_hop"`
	OspfArea        types.String `tfsdk:"ospf_area"`
	Policy          types.List   `tfsdk:"policy"`
	Preference      types.Int64  `tfsdk:"preference"`
	Protocol        types.Set    `tfsdk:"protocol"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
}

func (block *policyoptionsPolicyStatementBlockToConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type policyoptionsPolicyStatementBlockThen struct {
	Action          types.String                                                `tfsdk:"action"`
	ASPathExpand    types.String                                                `tfsdk:"as_path_expand"`
	ASPathPrepend   types.String                                                `tfsdk:"as_path_prepend"`
	DefaultAction   types.String                                                `tfsdk:"default_action"`
	LoadBalance     types.String                                                `tfsdk:"load_balance"`
	Next            types.String                                                `tfsdk:"next"`
	NextHop         types.String                                                `tfsdk:"next_hop"`
	Origin          types.String                                                `tfsdk:"origin"`
	Community       []policyoptionsPolicyStatementBlockThenBlockActionValue     `tfsdk:"community"`
	LocalPreference *policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"local_preference"`
	Metric          *policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"metric"`
	Preference      *policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"preference"`
}

func (block *policyoptionsPolicyStatementBlockThen) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type policyoptionsPolicyStatementBlockThenConfig struct {
	Action          types.String                                                `tfsdk:"action"`
	ASPathExpand    types.String                                                `tfsdk:"as_path_expand"`
	ASPathPrepend   types.String                                                `tfsdk:"as_path_prepend"`
	DefaultAction   types.String                                                `tfsdk:"default_action"`
	LoadBalance     types.String                                                `tfsdk:"load_balance"`
	Next            types.String                                                `tfsdk:"next"`
	NextHop         types.String                                                `tfsdk:"next_hop"`
	Origin          types.String                                                `tfsdk:"origin"`
	Community       types.List                                                  `tfsdk:"community"`
	LocalPreference *policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"local_preference"`
	Metric          *policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"metric"`
	Preference      *policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"preference"`
}

func (block *policyoptionsPolicyStatementBlockThenConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type policyoptionsPolicyStatementBlockThenBlockActionValue struct {
	Action types.String `tfsdk:"action"`
	Value  types.String `tfsdk:"value"`
}

type policyoptionsPolicyStatementBlockThenBlockActionValueInt64 struct {
	Action types.String `tfsdk:"action"`
	Value  types.Int64  `tfsdk:"value"`
}

//nolint:gocognit,gocyclo
func (rsc *policyoptionsPolicyStatement) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config policyoptionsPolicyStatementConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.DynamicDB.IsNull() &&
		config.From == nil &&
		config.To == nil &&
		config.Then == nil &&
		config.Term.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one of dynamic_db, from, to, then or term block must be specified",
		)
	}

	if config.From != nil {
		if config.From.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("from").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"from block is empty",
			)
		}
		if !config.From.BgpASPathCalcLength.IsNull() && !config.From.BgpASPathCalcLength.IsUnknown() {
			var bgpASPathCalcLength []policyoptionsPolicyStatementBlockFromBlockCountMatch
			asDiags := config.From.BgpASPathCalcLength.ElementsAs(ctx, &bgpASPathCalcLength, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			bgpASPathCalcLengthCount := make(map[int64]struct{})
			for _, v := range bgpASPathCalcLength {
				if !v.Count.IsNull() && !v.Count.IsUnknown() {
					count := v.Count.ValueInt64()
					if _, ok := bgpASPathCalcLengthCount[count]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("from").AtName("bgp_as_path_calc_length"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple bgp_as_path_calc_length blocks with the same count %d"+
								" in from block", count),
						)
					}
					bgpASPathCalcLengthCount[count] = struct{}{}
				}
			}
		}
		if !config.From.BgpASPathUniqueCount.IsNull() && !config.From.BgpASPathUniqueCount.IsUnknown() {
			var bgpASPathUniqueCount []policyoptionsPolicyStatementBlockFromBlockCountMatch
			asDiags := config.From.BgpASPathUniqueCount.ElementsAs(ctx, &bgpASPathUniqueCount, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			bgpASPathUniqueCountCount := make(map[int64]struct{})
			for _, v := range bgpASPathUniqueCount {
				if !v.Count.IsNull() && !v.Count.IsUnknown() {
					count := v.Count.ValueInt64()
					if _, ok := bgpASPathUniqueCountCount[count]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("from").AtName("bgp_as_path_unique_count"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple bgp_as_path_unique_count blocks with the same count %d"+
								" in from block", count),
						)
					}
					bgpASPathUniqueCountCount[count] = struct{}{}
				}
			}
		}
		if !config.From.BgpCommunityCount.IsNull() && !config.From.BgpCommunityCount.IsUnknown() {
			var bgpCommunityCount []policyoptionsPolicyStatementBlockFromBlockCountMatch
			asDiags := config.From.BgpCommunityCount.ElementsAs(ctx, &bgpCommunityCount, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			bgpCommunityCountCount := make(map[int64]struct{})
			for _, v := range bgpCommunityCount {
				if !v.Count.IsNull() && !v.Count.IsUnknown() {
					count := v.Count.ValueInt64()
					if _, ok := bgpCommunityCountCount[count]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("from").AtName("bgp_community_count"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple bgp_community_count blocks with the same count %d"+
								" in from block", count),
						)
					}
					bgpCommunityCountCount[count] = struct{}{}
				}
			}
		}
		if !config.From.NextHopWeight.IsNull() && !config.From.NextHopWeight.IsUnknown() {
			var nextHopWeight []policyoptionsPolicyStatementBlockFromBlockMatchWeight
			asDiags := config.From.NextHopWeight.ElementsAs(ctx, &nextHopWeight, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			nextHopWeightBlock := make(map[string]struct{})
			for _, v := range nextHopWeight {
				if !v.Match.IsNull() && !v.Match.IsUnknown() &&
					!v.Weight.IsNull() && !v.Weight.IsUnknown() {
					values := v.Match.ValueString() + " " + utils.ConvI64toa(v.Weight.ValueInt64())
					if _, ok := nextHopWeightBlock[values]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("from").AtName("next_hop_weight"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple next_hop_weight blocks with the same argument values %q"+
								" in from block", values),
						)
					}
					nextHopWeightBlock[values] = struct{}{}
				}
			}
		}
		if !config.From.RouteFilter.IsNull() && !config.From.RouteFilter.IsUnknown() {
			var routeFilter []policyoptionsPolicyStatementBlockFromBlockRouteFilter
			asDiags := config.From.RouteFilter.ElementsAs(ctx, &routeFilter, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			routeFilterBlock := make(map[string]struct{})
			for ii, v := range routeFilter {
				if !v.Route.IsNull() && !v.Route.IsUnknown() &&
					!v.Option.IsNull() && !v.Option.IsUnknown() {
					values := v.Route.ValueString() + " " + v.Option.ValueString()
					if !v.OptionValue.IsNull() {
						if v.OptionValue.IsUnknown() {
							continue
						}
						values += " " + v.OptionValue.ValueString()
					}
					if _, ok := routeFilterBlock[values]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("from").AtName("route_filter").AtListIndex(ii).AtName("route"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple route_filter blocks with the same argument values %q"+
								" in from block", values),
						)
					}
					routeFilterBlock[values] = struct{}{}
				}
			}
		}
	}

	if config.To != nil {
		if config.To.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("to").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"to block is empty",
			)
		}
	}

	if config.Then != nil {
		if config.Then.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("then").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"then block is empty",
			)
		}
		if !config.Then.Community.IsNull() && !config.Then.Community.IsUnknown() {
			var community []policyoptionsPolicyStatementBlockThenBlockActionValue
			asDiags := config.Then.Community.ElementsAs(ctx, &community, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			communityBlock := make(map[string]struct{})
			for i, v := range community {
				if !v.Action.IsNull() && !v.Action.IsUnknown() &&
					!v.Value.IsNull() && !v.Value.IsUnknown() {
					values := v.Action.ValueString() + " " + v.Value.ValueString()
					if _, ok := communityBlock[values]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("then").AtName("community").AtListIndex(i).AtName("action"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple community blocks with the same argument values %q"+
								" in then block", values),
						)
					}
					communityBlock[values] = struct{}{}
				}
			}
		}
		if config.Then.LocalPreference != nil {
			if config.Then.LocalPreference.Action.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("then").AtName("local_preference").AtName("action"),
					tfdiag.MissingConfigErrSummary,
					"action must be specified in then.local_preference block",
				)
			}
			if config.Then.LocalPreference.Value.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("then").AtName("local_preference").AtName("value"),
					tfdiag.MissingConfigErrSummary,
					"value must be specified in then.local_preference block",
				)
			}
		}
		if config.Then.Metric != nil {
			if config.Then.Metric.Action.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("then").AtName("metric").AtName("action"),
					tfdiag.MissingConfigErrSummary,
					"action must be specified in then.metric block",
				)
			}
			if config.Then.Metric.Value.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("then").AtName("metric").AtName("value"),
					tfdiag.MissingConfigErrSummary,
					"value must be specified in then.metric block",
				)
			}
		}
		if config.Then.Preference != nil {
			if config.Then.Preference.Action.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("then").AtName("preference").AtName("action"),
					tfdiag.MissingConfigErrSummary,
					"action must be specified in then.preference block",
				)
			}
			if config.Then.Preference.Value.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("then").AtName("preference").AtName("value"),
					tfdiag.MissingConfigErrSummary,
					"value must be specified in then.preference block",
				)
			}
		}
	}

	if !config.Term.IsNull() && !config.Term.IsUnknown() {
		var term []policyoptionsPolicyStatementBlockTermConfig
		asDiags := config.Term.ElementsAs(ctx, &term, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		termName := make(map[string]struct{})
		for i, block := range term {
			if block.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("term").AtListIndex(i).AtName("name"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("term block %q is empty", block.Name.ValueString()),
				)
			}
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := termName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple term blocks with the same name %q", name),
					)
				}
				termName[name] = struct{}{}
			}
			if block.From != nil {
				if block.From.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("from").AtName("*"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("from block is empty in term block %q", block.Name.ValueString()),
					)
				}
				if !block.From.BgpASPathCalcLength.IsNull() && !block.From.BgpASPathCalcLength.IsUnknown() {
					var bgpASPathCalcLength []policyoptionsPolicyStatementBlockFromBlockCountMatch
					asDiags := block.From.BgpASPathCalcLength.ElementsAs(ctx, &bgpASPathCalcLength, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}

					bgpASPathCalcLengthCount := make(map[int64]struct{})
					for _, v := range bgpASPathCalcLength {
						if !v.Count.IsNull() && !v.Count.IsUnknown() {
							count := v.Count.ValueInt64()
							if _, ok := bgpASPathCalcLengthCount[count]; ok {
								resp.Diagnostics.AddAttributeError(
									path.Root("term").AtListIndex(i).AtName("from").AtName("bgp_as_path_calc_length"),
									tfdiag.DuplicateConfigErrSummary,
									fmt.Sprintf("multiple bgp_as_path_calc_length blocks with the same count %d"+
										" in from block in term block %q", count, block.Name.ValueString()),
								)
							}
							bgpASPathCalcLengthCount[count] = struct{}{}
						}
					}
				}
				if !block.From.BgpASPathUniqueCount.IsNull() && !block.From.BgpASPathUniqueCount.IsUnknown() {
					var bgpASPathUniqueCount []policyoptionsPolicyStatementBlockFromBlockCountMatch
					asDiags := block.From.BgpASPathUniqueCount.ElementsAs(ctx, &bgpASPathUniqueCount, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}

					bgpASPathUniqueCountCount := make(map[int64]struct{})
					for _, v := range bgpASPathUniqueCount {
						if !v.Count.IsNull() && !v.Count.IsUnknown() {
							count := v.Count.ValueInt64()
							if _, ok := bgpASPathUniqueCountCount[count]; ok {
								resp.Diagnostics.AddAttributeError(
									path.Root("term").AtListIndex(i).AtName("from").AtName("bgp_as_path_unique_count"),
									tfdiag.DuplicateConfigErrSummary,
									fmt.Sprintf("multiple bgp_as_path_unique_count blocks with the same count %d"+
										" in from block in term block %q", count, block.Name.ValueString()),
								)
							}
							bgpASPathUniqueCountCount[count] = struct{}{}
						}
					}
				}
				if !block.From.BgpCommunityCount.IsNull() && !block.From.BgpCommunityCount.IsUnknown() {
					var bgpCommunityCount []policyoptionsPolicyStatementBlockFromBlockCountMatch
					asDiags := block.From.BgpCommunityCount.ElementsAs(ctx, &bgpCommunityCount, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}

					bgpCommunityCountCount := make(map[int64]struct{})
					for _, v := range bgpCommunityCount {
						if !v.Count.IsNull() && !v.Count.IsUnknown() {
							count := v.Count.ValueInt64()
							if _, ok := bgpCommunityCountCount[count]; ok {
								resp.Diagnostics.AddAttributeError(
									path.Root("term").AtListIndex(i).AtName("from").AtName("bgp_community_count"),
									tfdiag.DuplicateConfigErrSummary,
									fmt.Sprintf("multiple bgp_community_count blocks with the same count %d"+
										" in from block in term block %q", count, block.Name.ValueString()),
								)
							}
							bgpCommunityCountCount[count] = struct{}{}
						}
					}
				}
				if !block.From.NextHopWeight.IsNull() && !block.From.NextHopWeight.IsUnknown() {
					var nextHopWeight []policyoptionsPolicyStatementBlockFromBlockMatchWeight
					asDiags := block.From.NextHopWeight.ElementsAs(ctx, &nextHopWeight, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}

					nextHopWeightBlock := make(map[string]struct{})
					for _, v := range nextHopWeight {
						if !v.Match.IsNull() && !v.Match.IsUnknown() &&
							!v.Weight.IsNull() && !v.Weight.IsUnknown() {
							values := v.Match.ValueString() + " " + utils.ConvI64toa(v.Weight.ValueInt64())
							if _, ok := nextHopWeightBlock[values]; ok {
								resp.Diagnostics.AddAttributeError(
									path.Root("term").AtListIndex(i).AtName("from").AtName("next_hop_weight"),
									tfdiag.DuplicateConfigErrSummary,
									fmt.Sprintf("multiple next_hop_weight blocks with the same argument values %q"+
										" in from block in term block %q", values, block.Name.ValueString()),
								)
							}
							nextHopWeightBlock[values] = struct{}{}
						}
					}
				}
				if !block.From.RouteFilter.IsNull() && !block.From.RouteFilter.IsUnknown() {
					var routeFilter []policyoptionsPolicyStatementBlockFromBlockRouteFilter
					asDiags := block.From.RouteFilter.ElementsAs(ctx, &routeFilter, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}

					routeFilterBlock := make(map[string]struct{})
					for ii, v := range routeFilter {
						if !v.Route.IsNull() && !v.Route.IsUnknown() &&
							!v.Option.IsNull() && !v.Option.IsUnknown() {
							values := v.Route.ValueString() + " " + v.Option.ValueString()
							if !v.OptionValue.IsNull() {
								if v.OptionValue.IsUnknown() {
									continue
								}
								values += " " + v.OptionValue.ValueString()
							}
							if _, ok := routeFilterBlock[values]; ok {
								resp.Diagnostics.AddAttributeError(
									path.Root("term").AtListIndex(i).AtName("from").AtName("route_filter").AtListIndex(ii).AtName("route"),
									tfdiag.DuplicateConfigErrSummary,
									fmt.Sprintf("multiple route_filter blocks with the same argument values %q"+
										" in from block in term block %q", values, block.Name.ValueString()),
								)
							}
							routeFilterBlock[values] = struct{}{}
						}
					}
				}
			}
			if block.To != nil {
				if block.To.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("to").AtName("*"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("to block is empty in term block %q", block.Name.ValueString()),
					)
				}
			}
			if block.Then != nil {
				if block.Then.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("then").AtName("*"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("then block is empty in term block %q", block.Name.ValueString()),
					)
				}
				if !block.Then.Community.IsNull() && !block.Then.Community.IsUnknown() {
					var community []policyoptionsPolicyStatementBlockThenBlockActionValue
					asDiags := block.Then.Community.ElementsAs(ctx, &community, false)
					if asDiags.HasError() {
						resp.Diagnostics.Append(asDiags...)

						return
					}

					communityBlock := make(map[string]struct{})
					for ii, v := range community {
						if !v.Action.IsNull() && !v.Action.IsUnknown() &&
							!v.Value.IsNull() && !v.Value.IsUnknown() {
							values := v.Action.ValueString() + " " + v.Value.ValueString()
							if _, ok := communityBlock[values]; ok {
								resp.Diagnostics.AddAttributeError(
									path.Root("term").AtListIndex(i).AtName("then").AtName("community").AtListIndex(ii).AtName("action"),
									tfdiag.DuplicateConfigErrSummary,
									fmt.Sprintf("multiple community blocks with the same argument values %q"+
										" in then block in term block %q", values, block.Name.ValueString()),
								)
							}
							communityBlock[values] = struct{}{}
						}
					}
				}
				if block.Then.LocalPreference != nil {
					if block.Then.LocalPreference.Action.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("term").AtListIndex(i).AtName("then").AtName("local_preference").AtName("action"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("action must be specified in then.local_preference block"+
								" in term block %q", block.Name.ValueString()),
						)
					}
					if block.Then.LocalPreference.Value.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("term").AtListIndex(i).AtName("then").AtName("local_preference").AtName("value"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("value must be specified in then.local_preference block"+
								" in term block %q", block.Name.ValueString()),
						)
					}
				}
				if block.Then.Metric != nil {
					if block.Then.Metric.Action.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("term").AtListIndex(i).AtName("then").AtName("metric").AtName("action"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("action must be specified in then.metric block"+
								" in term block %q", block.Name.ValueString()),
						)
					}
					if block.Then.Metric.Value.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("term").AtListIndex(i).AtName("then").AtName("metric").AtName("value"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("value must be specified in then.metric block"+
								" in term block %q", block.Name.ValueString()),
						)
					}
				}
				if block.Then.Preference != nil {
					if block.Then.Preference.Action.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("term").AtListIndex(i).AtName("then").AtName("preference").AtName("action"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("action must be specified in then.preference block"+
								" in term block %q", block.Name.ValueString()),
						)
					}
					if block.Then.Preference.Value.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("term").AtListIndex(i).AtName("then").AtName("preference").AtName("value"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("value must be specified in then.preference block"+
								" in term block %q", block.Name.ValueString()),
						)
					}
				}
			}
		}
	}
}

func (rsc *policyoptionsPolicyStatement) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan policyoptionsPolicyStatementData
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
			policyStatementExists, err := checkPolicyoptionsPolicyStatementExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if policyStatementExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			policyStatementExists, err := checkPolicyoptionsPolicyStatementExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !policyStatementExists {
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

func (rsc *policyoptionsPolicyStatement) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data policyoptionsPolicyStatementData
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
			if !state.AddItToForwardingTableExport.ValueBool() {
				data.AddItToForwardingTableExport = types.BoolNull()
			}
		},
		resp,
	)
}

func (rsc *policyoptionsPolicyStatement) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state policyoptionsPolicyStatementData
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

func (rsc *policyoptionsPolicyStatement) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state policyoptionsPolicyStatementData
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

func (rsc *policyoptionsPolicyStatement) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data policyoptionsPolicyStatementData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "name"),
	)

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("add_it_to_forwarding_table_export"), types.BoolNull())...)
}

func checkPolicyoptionsPolicyStatementExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"policy-options policy-statement \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *policyoptionsPolicyStatementData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *policyoptionsPolicyStatementData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *policyoptionsPolicyStatementData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set policy-options policy-statement \"" + rscData.Name.ValueString() + "\" "

	if rscData.DynamicDB.ValueBool() {
		configSet = append(configSet, setPrefix+"dynamic-db")
	}
	if rscData.From != nil {
		if rscData.From.isEmpty() {
			return path.Root("from").AtName("*"),
				errors.New("from block is empty")
		}
		blockSet, pathErr, err := rscData.From.configSet(setPrefix, path.Root("from"))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.To != nil {
		if rscData.To.isEmpty() {
			return path.Root("to").AtName("*"),
				errors.New("to block is empty")
		}
		configSet = append(configSet, rscData.To.configSet(setPrefix)...)
	}
	if rscData.Then != nil {
		if rscData.Then.isEmpty() {
			return path.Root("then").AtName("*"),
				errors.New("then block is empty")
		}
		blockSet, pathErr, err := rscData.Then.configSet(setPrefix, path.Root("then"))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	termName := make(map[string]struct{})
	for i, block := range rscData.Term {
		name := block.Name.ValueString()
		if _, ok := termName[name]; ok {
			return path.Root("term").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple term blocks with the same name %q", name)
		}
		termName[name] = struct{}{}
		if block.isEmpty() {
			return path.Root("term").AtListIndex(i).AtName("name"),
				fmt.Errorf("term block %q is empty", name)
		}

		setPrefixTerm := setPrefix + "term \"" + name + "\" "
		if block.From != nil {
			if block.From.isEmpty() {
				return path.Root("term").AtListIndex(i).AtName("from").AtName("*"),
					fmt.Errorf("from block is empty in term block %q", name)
			}
			blockSet, pathErr, err := block.From.configSet(setPrefixTerm, path.Root("term").AtListIndex(i).AtName("from"))
			if err != nil {
				return pathErr, err
			}
			configSet = append(configSet, blockSet...)
		}
		if block.To != nil {
			if block.To.isEmpty() {
				return path.Root("term").AtListIndex(i).AtName("to").AtName("*"),
					fmt.Errorf("to block is empty in term block %q", name)
			}
			configSet = append(configSet, block.To.configSet(setPrefixTerm)...)
		}
		if block.Then != nil {
			if block.Then.isEmpty() {
				return path.Root("term").AtListIndex(i).AtName("then").AtName("*"),
					fmt.Errorf("then block is empty in term block %q", name)
			}
			blockSet, pathErr, err := block.Then.configSet(setPrefixTerm, path.Root("term").AtListIndex(i).AtName("then"))
			if err != nil {
				return pathErr, err
			}
			configSet = append(configSet, blockSet...)
		}
	}
	if rscData.AddItToForwardingTableExport.ValueBool() {
		configSet = append(configSet,
			"set routing-options forwarding-table export \""+rscData.Name.ValueString()+"\"",
		)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *policyoptionsPolicyStatementBlockFrom) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "from "

	if block.AggregateContributor.ValueBool() {
		configSet = append(configSet, setPrefix+"aggregate-contributor")
	}
	for _, v := range block.BgpASPath {
		configSet = append(configSet, setPrefix+"as-path \""+v.ValueString()+"\"")
	}
	bgpASPathCalcLengthCount := make(map[int64]struct{})
	for _, v := range block.BgpASPathCalcLength {
		count := v.Count.ValueInt64()
		if _, ok := bgpASPathCalcLengthCount[count]; ok {
			return configSet,
				pathRoot.AtName("bgp_as_path_calc_length"),
				fmt.Errorf("multiple bgp_as_path_calc_length blocks with the same count %d in from block", count)
		}
		bgpASPathCalcLengthCount[count] = struct{}{}

		configSet = append(configSet,
			setPrefix+"as-path-calc-length "+utils.ConvI64toa(count)+" "+v.Match.ValueString())
	}
	for _, v := range block.BgpASPathGroup {
		configSet = append(configSet, setPrefix+"as-path-group \""+v.ValueString()+"\"")
	}
	bgpASPathUniqueCountCount := make(map[int64]struct{})
	for _, v := range block.BgpASPathUniqueCount {
		count := v.Count.ValueInt64()
		if _, ok := bgpASPathUniqueCountCount[count]; ok {
			return configSet,
				pathRoot.AtName("bgp_as_path_unique_count"),
				fmt.Errorf("multiple bgp_as_path_unique_count blocks with the same count %d in from block", count)
		}
		bgpASPathUniqueCountCount[count] = struct{}{}

		configSet = append(configSet,
			setPrefix+"as-path-unique-count "+utils.ConvI64toa(count)+" "+v.Match.ValueString())
	}
	for _, v := range block.BgpCommunity {
		configSet = append(configSet, setPrefix+"community \""+v.ValueString()+"\"")
	}
	bgpCommunityCountCount := make(map[int64]struct{})
	for _, v := range block.BgpCommunityCount {
		count := v.Count.ValueInt64()
		if _, ok := bgpCommunityCountCount[count]; ok {
			return configSet,
				pathRoot.AtName("bgp_community_count"),
				fmt.Errorf("multiple bgp_community_count blocks with the same count %d in from block", count)
		}
		bgpCommunityCountCount[count] = struct{}{}

		configSet = append(configSet,
			setPrefix+"community-count "+utils.ConvI64toa(count)+" "+v.Match.ValueString())
	}
	if v := block.BgpOrigin.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"origin "+v)
	}
	if !block.BgpSrteDiscriminator.IsNull() {
		configSet = append(configSet, setPrefix+"bgp-srte-discriminator "+
			utils.ConvI64toa(block.BgpSrteDiscriminator.ValueInt64()))
	}
	if !block.Color.IsNull() {
		configSet = append(configSet, setPrefix+"color "+
			utils.ConvI64toa(block.Color.ValueInt64()))
	}
	for _, v := range block.EvpnESI {
		configSet = append(configSet, setPrefix+"evpn-esi "+v.ValueString())
	}
	if v := block.EvpnMACRoute.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"evpn-mac-route "+v)
	}
	for _, v := range block.EvpnTag {
		configSet = append(configSet, setPrefix+"evpn-tag "+
			utils.ConvI64toa(v.ValueInt64()))
	}
	if v := block.Family.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"family "+v)
	}
	if !block.LocalPreference.IsNull() {
		configSet = append(configSet, setPrefix+"local-preference "+
			utils.ConvI64toa(block.LocalPreference.ValueInt64()))
	}
	for _, v := range block.Interface {
		configSet = append(configSet, setPrefix+"interface "+v.ValueString())
	}
	if !block.Metric.IsNull() {
		configSet = append(configSet, setPrefix+"metric "+
			utils.ConvI64toa(block.Metric.ValueInt64()))
	}
	for _, v := range block.Neighbor {
		configSet = append(configSet, setPrefix+"neighbor "+v.ValueString())
	}
	for _, v := range block.NextHop {
		configSet = append(configSet, setPrefix+"next-hop "+v.ValueString())
	}
	if block.NextHopTypeMerged.ValueBool() {
		configSet = append(configSet, setPrefix+"next-hop-type merged")
	}
	nextHopWeightBlock := make(map[string]struct{})
	for _, v := range block.NextHopWeight {
		values := v.Match.ValueString() + " " + utils.ConvI64toa(v.Weight.ValueInt64())
		if _, ok := nextHopWeightBlock[values]; ok {
			return configSet,
				pathRoot.AtName("next_hop_weight"),
				fmt.Errorf("multiple next_hop_weight blocks with the same argument values %q in from block", values)
		}
		nextHopWeightBlock[values] = struct{}{}

		configSet = append(configSet,
			setPrefix+"nexthop-weight "+v.Match.ValueString()+" "+utils.ConvI64toa(v.Weight.ValueInt64()))
	}
	if v := block.OspfArea.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"area "+v)
	}
	for _, v := range block.Policy {
		configSet = append(configSet, setPrefix+"policy \""+v.ValueString()+"\"")
	}
	if !block.Preference.IsNull() {
		configSet = append(configSet, setPrefix+"preference "+
			utils.ConvI64toa(block.Preference.ValueInt64()))
	}
	for _, v := range block.PrefixList {
		configSet = append(configSet, setPrefix+"prefix-list \""+v.ValueString()+"\"")
	}
	for _, v := range block.Protocol {
		configSet = append(configSet, setPrefix+"protocol "+v.ValueString())
	}
	routeFilterBlock := make(map[string]struct{})
	for i, v := range block.RouteFilter {
		values := v.Route.ValueString() + " " + v.Option.ValueString() + " " + v.OptionValue.ValueString()
		if _, ok := routeFilterBlock[values]; ok {
			return configSet,
				pathRoot.AtName("route_filter").AtListIndex(i).AtName("route"),
				fmt.Errorf("multiple route_filter blocks with the same argument values %q in from block", values)
		}
		routeFilterBlock[values] = struct{}{}

		setRoutFilter := setPrefix + "route-filter " +
			v.Route.ValueString() + " " + v.Option.ValueString()
		if v2 := v.OptionValue.ValueString(); v2 != "" {
			setRoutFilter += " " + v2
		}
		configSet = append(configSet, setRoutFilter)
	}
	if v := block.RouteType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"route-type "+v)
	}
	if v := block.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"instance "+v)
	}
	if !block.SrteColor.IsNull() {
		configSet = append(configSet, setPrefix+"srte-color "+
			utils.ConvI64toa(block.SrteColor.ValueInt64()))
	}
	if v := block.State.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"state "+v)
	}
	for _, v := range block.TunnelType {
		configSet = append(configSet, setPrefix+"tunnel-type "+v.ValueString())
	}
	if v := block.ValidationDatabase.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"validation-database "+v)
	}

	return configSet, path.Empty(), nil
}

func (block *policyoptionsPolicyStatementBlockTo) configSet(setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "to "

	for _, v := range block.BgpASPath {
		configSet = append(configSet, setPrefix+"as-path \""+v.ValueString()+"\"")
	}
	for _, v := range block.BgpASPathGroup {
		configSet = append(configSet, setPrefix+"as-path-group \""+v.ValueString()+"\"")
	}
	for _, v := range block.BgpCommunity {
		configSet = append(configSet, setPrefix+"community \""+v.ValueString()+"\"")
	}
	if v := block.BgpOrigin.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"origin "+v)
	}
	if v := block.Family.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"family "+v)
	}
	if !block.LocalPreference.IsNull() {
		configSet = append(configSet, setPrefix+"local-preference "+
			utils.ConvI64toa(block.LocalPreference.ValueInt64()))
	}
	for _, v := range block.Interface {
		configSet = append(configSet, setPrefix+"interface "+v.ValueString())
	}
	if !block.Metric.IsNull() {
		configSet = append(configSet, setPrefix+"metric "+
			utils.ConvI64toa(block.Metric.ValueInt64()))
	}
	for _, v := range block.Neighbor {
		configSet = append(configSet, setPrefix+"neighbor "+v.ValueString())
	}
	for _, v := range block.NextHop {
		configSet = append(configSet, setPrefix+"next-hop "+v.ValueString())
	}
	if v := block.OspfArea.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"area "+v)
	}
	for _, v := range block.Policy {
		configSet = append(configSet, setPrefix+"policy \""+v.ValueString()+"\"")
	}
	if !block.Preference.IsNull() {
		configSet = append(configSet, setPrefix+"preference "+
			utils.ConvI64toa(block.Preference.ValueInt64()))
	}
	for _, v := range block.Protocol {
		configSet = append(configSet, setPrefix+"protocol "+v.ValueString())
	}
	if v := block.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"instance "+v)
	}

	return configSet
}

func (block *policyoptionsPolicyStatementBlockThen) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "then "

	if v := block.Action.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+v)
	}
	if v := block.ASPathExpand.ValueString(); v != "" {
		if strings.HasPrefix(v, "last-as") {
			configSet = append(configSet, setPrefix+"as-path-expand "+v)
		} else {
			configSet = append(configSet, setPrefix+"as-path-expand \""+v+"\"")
		}
	}
	if v := block.ASPathPrepend.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"as-path-prepend \""+v+"\"")
	}
	communityBlock := make(map[string]struct{})
	for i, v := range block.Community {
		values := v.Action.ValueString() + " " + v.Value.ValueString()
		if _, ok := communityBlock[values]; ok {
			return configSet,
				pathRoot.AtName("community").AtListIndex(i).AtName("action"),
				fmt.Errorf("multiple community blocks with the same argument values %q in then block", values)
		}
		communityBlock[values] = struct{}{}

		configSet = append(configSet, setPrefix+
			"community "+v.Action.ValueString()+" \""+v.Value.ValueString()+"\"")
	}
	if v := block.DefaultAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"default-action "+v)
	}
	if v := block.LoadBalance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"load-balance "+v)
	}
	if block.LocalPreference != nil {
		if block.LocalPreference.Action.ValueString() == "none" {
			configSet = append(configSet, setPrefix+"local-preference "+
				utils.ConvI64toa(block.LocalPreference.Value.ValueInt64()))
		} else {
			configSet = append(configSet, setPrefix+"local-preference "+
				block.LocalPreference.Action.ValueString()+" "+
				utils.ConvI64toa(block.LocalPreference.Value.ValueInt64()))
		}
	}
	if block.Metric != nil {
		if block.Metric.Action.ValueString() == "none" {
			configSet = append(configSet, setPrefix+"metric "+
				utils.ConvI64toa(block.Metric.Value.ValueInt64()))
		} else {
			configSet = append(configSet, setPrefix+"metric "+
				block.Metric.Action.ValueString()+" "+
				utils.ConvI64toa(block.Metric.Value.ValueInt64()))
		}
	}
	if v := block.Next.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"next "+v)
	}
	if v := block.NextHop.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"next-hop "+v)
	}
	if v := block.Origin.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"origin "+v)
	}
	if block.Preference != nil {
		if block.Preference.Action.ValueString() == "none" {
			configSet = append(configSet, setPrefix+"preference "+
				utils.ConvI64toa(block.Preference.Value.ValueInt64()))
		} else {
			configSet = append(configSet, setPrefix+"preference "+
				block.Preference.Action.ValueString()+" "+
				utils.ConvI64toa(block.Preference.Value.ValueInt64()))
		}
	}

	return configSet, path.Empty(), nil
}

func (rscData *policyoptionsPolicyStatementData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"policy-options policy-statement \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case itemTrim == "dynamic-db":
				rscData.DynamicDB = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "term "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				var term policyoptionsPolicyStatementBlockTerm
				rscData.Term, term = tfdata.ExtractBlockWithTFTypesString(
					rscData.Term, "Name", strings.Trim(name, "\""))
				term.Name = types.StringValue(strings.Trim(name, "\""))
				balt.CutPrefixInString(&itemTrim, name+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "from "):
					if term.From == nil {
						term.From = &policyoptionsPolicyStatementBlockFrom{}
					}
					if err := term.From.read(itemTrim); err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "to "):
					if term.To == nil {
						term.To = &policyoptionsPolicyStatementBlockTo{}
					}
					if err := term.To.read(itemTrim); err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "then "):
					if term.Then == nil {
						term.Then = &policyoptionsPolicyStatementBlockThen{}
					}
					if err := term.Then.read(itemTrim); err != nil {
						return err
					}
				}
				rscData.Term = append(rscData.Term, term)
			case balt.CutPrefixInString(&itemTrim, "from "):
				if rscData.From == nil {
					rscData.From = &policyoptionsPolicyStatementBlockFrom{}
				}
				if err := rscData.From.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "to "):
				if rscData.To == nil {
					rscData.To = &policyoptionsPolicyStatementBlockTo{}
				}
				if err := rscData.To.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "then "):
				if rscData.Then == nil {
					rscData.Then = &policyoptionsPolicyStatementBlockThen{}
				}
				if err := rscData.Then.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	showConfigForwardingTableExport, err := junSess.Command(junos.CmdShowConfig +
		junos.RoutingOptionsWS + "forwarding-table export" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfigForwardingTableExport != junos.EmptyW {
		for _, item := range strings.Split(showConfigForwardingTableExport, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			balt.CutSuffixInString(&itemTrim, " ")
			if itemTrim == name || itemTrim == "\""+name+"\"" {
				rscData.AddItToForwardingTableExport = types.BoolValue(true)
			}
		}
	}

	return nil
}

func (block *policyoptionsPolicyStatementBlockFrom) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "aggregate-contributor":
		block.AggregateContributor = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "as-path "):
		block.BgpASPath = append(block.BgpASPath, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "as-path-calc-length "):
		itemTrimFields := strings.Split(itemTrim, " ")
		count, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
		if err != nil {
			return err
		}
		block.BgpASPathCalcLength = append(block.BgpASPathCalcLength,
			policyoptionsPolicyStatementBlockFromBlockCountMatch{
				Count: count,
				Match: types.StringValue(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" ")),
			},
		)
	case balt.CutPrefixInString(&itemTrim, "as-path-group "):
		block.BgpASPathGroup = append(block.BgpASPathGroup, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "as-path-unique-count "):
		itemTrimFields := strings.Split(itemTrim, " ")
		count, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
		if err != nil {
			return err
		}
		block.BgpASPathUniqueCount = append(block.BgpASPathUniqueCount,
			policyoptionsPolicyStatementBlockFromBlockCountMatch{
				Count: count,
				Match: types.StringValue(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" ")),
			},
		)
	case balt.CutPrefixInString(&itemTrim, "community "):
		block.BgpCommunity = append(block.BgpCommunity, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "community-count "):
		itemTrimFields := strings.Split(itemTrim, " ")
		count, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
		if err != nil {
			return err
		}
		block.BgpCommunityCount = append(block.BgpCommunityCount,
			policyoptionsPolicyStatementBlockFromBlockCountMatch{
				Count: count,
				Match: types.StringValue(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" ")),
			},
		)
	case balt.CutPrefixInString(&itemTrim, "origin "):
		block.BgpOrigin = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "bgp-srte-discriminator "):
		block.BgpSrteDiscriminator, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "color "):
		block.Color, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "evpn-esi "):
		block.EvpnESI = append(block.EvpnESI, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "evpn-mac-route "):
		block.EvpnMACRoute = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "evpn-tag "):
		evpnTag, err := tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
		block.EvpnTag = append(block.EvpnTag, evpnTag)
	case balt.CutPrefixInString(&itemTrim, "family "):
		block.Family = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "local-preference "):
		block.LocalPreference, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "interface "):
		block.Interface = append(block.Interface, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "metric "):
		block.Metric, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "neighbor "):
		block.Neighbor = append(block.Neighbor, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "next-hop "):
		block.NextHop = append(block.NextHop, types.StringValue(itemTrim))
	case itemTrim == "next-hop-type merged":
		block.NextHopTypeMerged = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "nexthop-weight "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 2 { // <match> <weight>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "nexthop-weight", itemTrim)
		}
		weight, err := tfdata.ConvAtoi64Value(itemTrimFields[1])
		if err != nil {
			return err
		}
		block.NextHopWeight = append(block.NextHopWeight,
			policyoptionsPolicyStatementBlockFromBlockMatchWeight{
				Match:  types.StringValue(itemTrimFields[0]),
				Weight: weight,
			})
	case balt.CutPrefixInString(&itemTrim, "area "):
		block.OspfArea = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "policy "):
		block.Policy = append(block.Policy, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "preference "):
		block.Preference, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "prefix-list "):
		block.PrefixList = append(block.PrefixList, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		block.Protocol = append(block.Protocol, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "route-filter "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 2 { // <route> <option> <option_value>?
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "route-filter", itemTrim)
		}
		routeFilter := policyoptionsPolicyStatementBlockFromBlockRouteFilter{
			Route:  types.StringValue(itemTrimFields[0]),
			Option: types.StringValue(itemTrimFields[1]),
		}
		if len(itemTrimFields) > 2 {
			routeFilter.OptionValue = types.StringValue(itemTrimFields[2])
		}
		block.RouteFilter = append(block.RouteFilter, routeFilter)
	case balt.CutPrefixInString(&itemTrim, "route-type "):
		block.RouteType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "instance "):
		block.RoutingInstance = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "srte-color "):
		block.SrteColor, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "state "):
		block.State = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "tunnel-type "):
		block.TunnelType = append(block.TunnelType, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "validation-database "):
		block.ValidationDatabase = types.StringValue(itemTrim)
	}

	return nil
}

func (block *policyoptionsPolicyStatementBlockTo) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "as-path "):
		block.BgpASPath = append(block.BgpASPath, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "as-path-group "):
		block.BgpASPathGroup = append(block.BgpASPathGroup, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "community "):
		block.BgpCommunity = append(block.BgpCommunity, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "origin "):
		block.BgpOrigin = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "family "):
		block.Family = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "local-preference "):
		block.LocalPreference, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "interface "):
		block.Interface = append(block.Interface, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "metric "):
		block.Metric, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "neighbor "):
		block.Neighbor = append(block.Neighbor, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "next-hop "):
		block.NextHop = append(block.NextHop, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "area "):
		block.OspfArea = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "policy "):
		block.Policy = append(block.Policy, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "preference "):
		block.Preference, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		block.Protocol = append(block.Protocol, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "instance "):
		block.RoutingInstance = types.StringValue(itemTrim)
	}

	return nil
}

func (block *policyoptionsPolicyStatementBlockThen) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "accept", itemTrim == "reject":
		block.Action = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "as-path-expand "):
		block.ASPathExpand = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "as-path-prepend "):
		block.ASPathPrepend = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "community "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 2 { // <action> <value>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "community", itemTrim)
		}
		block.Community = append(block.Community, policyoptionsPolicyStatementBlockThenBlockActionValue{
			Action: types.StringValue(itemTrimFields[0]),
			Value:  types.StringValue(strings.Trim(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" "), "\"")),
		})
	case balt.CutPrefixInString(&itemTrim, "default-action "):
		block.DefaultAction = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "load-balance "):
		block.LoadBalance = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "local-preference "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) == 1 { // <value>
			value, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
			if err != nil {
				return err
			}
			block.LocalPreference = &policyoptionsPolicyStatementBlockThenBlockActionValueInt64{
				Action: types.StringValue("none"),
				Value:  value,
			}
		} else { // <action> <value>
			value, err := tfdata.ConvAtoi64Value(itemTrimFields[1])
			if err != nil {
				return err
			}
			block.LocalPreference = &policyoptionsPolicyStatementBlockThenBlockActionValueInt64{
				Action: types.StringValue(itemTrimFields[0]),
				Value:  value,
			}
		}
	case balt.CutPrefixInString(&itemTrim, "metric "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) == 1 { // <value>
			value, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
			if err != nil {
				return err
			}
			block.Metric = &policyoptionsPolicyStatementBlockThenBlockActionValueInt64{
				Action: types.StringValue("none"),
				Value:  value,
			}
		} else { // <action> <value>
			value, err := tfdata.ConvAtoi64Value(itemTrimFields[1])
			if err != nil {
				return err
			}
			block.Metric = &policyoptionsPolicyStatementBlockThenBlockActionValueInt64{
				Action: types.StringValue(itemTrimFields[0]),
				Value:  value,
			}
		}
	case balt.CutPrefixInString(&itemTrim, "next "):
		block.Next = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "next-hop "):
		block.NextHop = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "origin "):
		block.Origin = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "preference "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) == 1 { // <value>
			value, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
			if err != nil {
				return err
			}
			block.Preference = &policyoptionsPolicyStatementBlockThenBlockActionValueInt64{
				Action: types.StringValue("none"),
				Value:  value,
			}
		} else { // <action> <value>
			value, err := tfdata.ConvAtoi64Value(itemTrimFields[1])
			if err != nil {
				return err
			}
			block.Preference = &policyoptionsPolicyStatementBlockThenBlockActionValueInt64{
				Action: types.StringValue(itemTrimFields[0]),
				Value:  value,
			}
		}
	}

	return nil
}

func (rscData *policyoptionsPolicyStatementData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete policy-options policy-statement \"" + rscData.Name.ValueString() + "\"",
	}
	if rscData.AddItToForwardingTableExport.ValueBool() {
		configSet = append(configSet,
			"delete routing-options forwarding-table export \""+rscData.Name.ValueString()+"\"",
		)
	}

	return junSess.ConfigSet(configSet)
}
