package provider

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &policyoptionsPolicyStatementDataSource{}
	_ datasource.DataSourceWithConfigure = &policyoptionsPolicyStatementDataSource{}
)

type policyoptionsPolicyStatementDataSource struct {
	client *junos.Client
}

func (dsc *policyoptionsPolicyStatementDataSource) typeName() string {
	return providerName + "_policyoptions_policy_statement"
}

func (dsc *policyoptionsPolicyStatementDataSource) junosName() string {
	return "policy-options policy-statement"
}

func (dsc *policyoptionsPolicyStatementDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newPolicyoptionsPolicyStatementDataSource() datasource.DataSource {
	return &policyoptionsPolicyStatementDataSource{}
}

func (dsc *policyoptionsPolicyStatementDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *policyoptionsPolicyStatementDataSource) Configure(
	ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*junos.Client)
	if !ok {
		unexpectedDataSourceConfigureType(ctx, req, resp)

		return
	}
	dsc.client = client
}

func (dsc *policyoptionsPolicyStatementDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get configuration from a " + dsc.junosName() + ".",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the data source with format `<name>`.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name to identify the policy.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"dynamic_db": schema.BoolAttribute{
				Computed:    true,
				Description: "Object may exist in dynamic database.",
			},
		},
		Blocks: map[string]schema.Block{
			"from": schema.SingleNestedBlock{
				Description: "Conditions to match the source of a route.",
				Attributes:  policyoptionsPolicyStatementDscBlockFromAttributesSchema(),
				Blocks:      policyoptionsPolicyStatementDscBlockFromBlocksSchema(),
			},
			"to": schema.SingleNestedBlock{
				Description: "Conditions to match the destination of a route.",
				Attributes:  policyoptionsPolicyStatementDscBlockToAttributesSchema(),
			},
			"then": schema.SingleNestedBlock{
				Description: "Actions to take if 'from' and 'to' conditions match.",
				Attributes:  policyoptionsPolicyStatementDscBlockThenAttributesSchema(),
				Blocks:      policyoptionsPolicyStatementDscBlockThenBlocksSchema(),
			},
			"term": schema.ListNestedBlock{
				Description: "For each policy term.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of term.",
						},
					},
					Blocks: map[string]schema.Block{
						"from": schema.SingleNestedBlock{
							Description: "Conditions to match the source of a route.",
							Attributes:  policyoptionsPolicyStatementDscBlockFromAttributesSchema(),
							Blocks:      policyoptionsPolicyStatementDscBlockFromBlocksSchema(),
						},
						"to": schema.SingleNestedBlock{
							Description: "Conditions to match the destination of a route.",
							Attributes:  policyoptionsPolicyStatementDscBlockToAttributesSchema(),
						},
						"then": schema.SingleNestedBlock{
							Description: "Actions to take if 'from' and 'to' conditions match.",
							Attributes:  policyoptionsPolicyStatementDscBlockThenAttributesSchema(),
							Blocks:      policyoptionsPolicyStatementDscBlockThenBlocksSchema(),
						},
					},
				},
			},
		},
	}
}

func policyoptionsPolicyStatementDscBlockFromAttributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"aggregate_contributor": schema.BoolAttribute{
			Computed:    true,
			Description: "Match more specifics of an aggregate.",
		},
		"bgp_as_path": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Name of AS path regular expression.",
		},
		"bgp_as_path_group": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Name of AS path group.",
		},
		"bgp_community": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "BGP community.",
		},
		"bgp_origin": schema.StringAttribute{
			Computed:    true,
			Description: "BGP origin attribute.",
		},
		"bgp_srte_discriminator": schema.Int64Attribute{
			Computed:    true,
			Description: "Srte discriminator.",
		},
		"color": schema.Int64Attribute{
			Computed:    true,
			Description: "Color (preference) value.",
		},
		"evpn_esi": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "ESI in EVPN Route.",
		},
		"evpn_mac_route": schema.StringAttribute{
			Computed:    true,
			Description: "EVPN Mac Route type.",
		},
		"evpn_tag": schema.SetAttribute{
			ElementType: types.Int64Type,
			Computed:    true,
			Description: "Tag in EVPN Route (0..4294967295).",
		},
		"family": schema.StringAttribute{
			Computed:    true,
			Description: "Family.",
		},
		"local_preference": schema.Int64Attribute{
			Computed:    true,
			Description: "Local preference associated with a route.",
		},
		"interface": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Interface name or address.",
		},
		"metric": schema.Int64Attribute{
			Computed:    true,
			Description: "Metric value.",
		},
		"neighbor": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Neighboring router.",
		},
		"next_hop": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Next-hop router.",
		},
		"next_hop_type_merged": schema.BoolAttribute{
			Computed:    true,
			Description: "Merged next hop.",
		},
		"ospf_area": schema.StringAttribute{
			Computed:    true,
			Description: "OSPF area identifier.",
		},
		"policy": schema.ListAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Name of policy to evaluate.",
		},
		"preference": schema.Int64Attribute{
			Computed:    true,
			Description: "Preference value.",
		},
		"prefix_list": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Prefix-lists of routes to match.",
		},
		"protocol": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Protocol from which route was learned.",
		},
		"route_type": schema.StringAttribute{
			Computed:    true,
			Description: "Route type.",
		},
		"routing_instance": schema.StringAttribute{
			Computed:    true,
			Description: "Routing protocol instance.",
		},
		"srte_color": schema.Int64Attribute{
			Computed:    true,
			Description: "Srte color.",
		},
		"state": schema.StringAttribute{
			Computed:    true,
			Description: "Route state.",
		},
		"tunnel_type": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Tunnel type.",
		},
		"validation_database": schema.StringAttribute{
			Computed:    true,
			Description: "Name to identify a validation-state.",
		},
	}
}

func policyoptionsPolicyStatementDscBlockFromBlocksSchema() map[string]schema.Block {
	return map[string]schema.Block{
		"bgp_as_path_calc_length": schema.SetNestedBlock{
			Description: "Number of BGP ASes excluding confederations.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of ASes (0..1024).",
					},
					"match": schema.StringAttribute{
						Computed:    true,
						Description: "Type of match: equal values, higher or equal values, lower or equal values.",
					},
				},
			},
		},
		"bgp_as_path_unique_count": schema.SetNestedBlock{
			Description: "Number of unique BGP ASes excluding confederations.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of ASes (0..1024).",
					},
					"match": schema.StringAttribute{
						Computed:    true,
						Description: "Type of match: equal values, higher or equal values, lower or equal values.",
					},
				},
			},
		},
		"bgp_community_count": schema.SetNestedBlock{
			Description: "Number of BGP communities.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of communities (0..1024).",
					},
					"match": schema.StringAttribute{
						Computed:    true,
						Description: "Type of match: equal values, higher or equal values, lower or equal values.",
					},
				},
			},
		},
		"next_hop_weight": schema.SetNestedBlock{
			Description: "Weight of the gateway.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"match": schema.StringAttribute{
						Computed:    true,
						Description: "Type of match for weight.",
					},
					"weight": schema.Int64Attribute{
						Computed:    true,
						Description: "Weight of the gateway (1..65535).",
					},
				},
			},
		},
		"route_filter": schema.ListNestedBlock{
			Description: "Routes to match.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"route": schema.StringAttribute{
						Computed:    true,
						Description: "IP address.",
					},
					"option": schema.StringAttribute{
						Computed:    true,
						Description: "Mask option.",
					},
					"option_value": schema.StringAttribute{
						Computed:    true,
						Description: "For options that need an argument.",
					},
				},
			},
		},
	}
}

func policyoptionsPolicyStatementDscBlockToAttributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"bgp_as_path": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Name of AS path regular expression.",
		},
		"bgp_as_path_group": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Name of AS path group.",
		},
		"bgp_community": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "BGP community.",
		},
		"bgp_origin": schema.StringAttribute{
			Computed:    true,
			Description: "BGP origin attribute.",
		},
		"family": schema.StringAttribute{
			Computed:    true,
			Description: "Family.",
		},
		"local_preference": schema.Int64Attribute{
			Computed:    true,
			Description: "Local preference associated with a route.",
		},
		"interface": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Interface name or address.",
		},
		"metric": schema.Int64Attribute{
			Computed:    true,
			Description: "Metric value.",
		},
		"neighbor": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Neighboring router.",
		},
		"next_hop": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Next-hop router.",
		},
		"ospf_area": schema.StringAttribute{
			Computed:    true,
			Description: "OSPF area identifier.",
		},
		"policy": schema.ListAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Name of policy to evaluate.",
		},
		"preference": schema.Int64Attribute{
			Computed:    true,
			Description: "Preference value.",
		},
		"protocol": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Protocol from which route was learned.",
		},
		"routing_instance": schema.StringAttribute{
			Computed:    true,
			Description: "Routing protocol instance.",
		},
	}
}

func policyoptionsPolicyStatementDscBlockThenAttributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"action": schema.StringAttribute{
			Computed:    true,
			Description: "Action `accept` or `reject`.",
		},
		"as_path_expand": schema.StringAttribute{
			Computed:    true,
			Description: "Prepend AS numbers prior to adding local-as.",
		},
		"as_path_prepend": schema.StringAttribute{
			Computed:    true,
			Description: "Prepend AS numbers to an AS path.",
		},
		"default_action": schema.StringAttribute{
			Computed:    true,
			Description: "Set default policy action.",
		},
		"load_balance": schema.StringAttribute{
			Computed:    true,
			Description: "Type of load balancing in forwarding table.",
		},
		"next": schema.StringAttribute{
			Computed:    true,
			Description: "Skip to next `policy` or `term`.",
		},
		"next_hop": schema.StringAttribute{
			Computed:    true,
			Description: "Set the address of the next-hop router.",
		},
		"origin": schema.StringAttribute{
			Computed:    true,
			Description: "BGP path origin.",
		},
	}
}

func policyoptionsPolicyStatementDscBlockThenBlocksSchema() map[string]schema.Block {
	return map[string]schema.Block{
		"community": schema.ListNestedBlock{
			Description: "For each community action.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"action": schema.StringAttribute{
						Computed:    true,
						Description: "Action on BGP community.",
					},
					"value": schema.StringAttribute{
						Computed:    true,
						Description: "Name to identify a BGP community.",
					},
				},
			},
		},
		"local_preference": schema.SingleNestedBlock{
			Description: "Declare local-preference action.",
			Attributes: map[string]schema.Attribute{
				"action": schema.StringAttribute{
					Computed:    true,
					Description: "Action on local-preference.",
				},
				"value": schema.Int64Attribute{
					Computed:    true,
					Description: "Value for action (local-preference, constant).",
				},
			},
		},
		"metric": schema.SingleNestedBlock{
			Description: "Declare metric action.",
			Attributes: map[string]schema.Attribute{
				"action": schema.StringAttribute{
					Computed:    true,
					Description: "Action on metric.",
				},
				"value": schema.Int64Attribute{
					Computed:    true,
					Description: "Value for action (metric, constant).",
				},
			},
		},
		"preference": schema.SingleNestedBlock{
			Description: "Declare preference action.",
			Attributes: map[string]schema.Attribute{
				"action": schema.StringAttribute{
					Computed:    true,
					Description: "Action on preference.",
				},
				"value": schema.Int64Attribute{
					Computed:    true,
					Description: "Value for action (preference, constant).",
				},
			},
		},
	}
}

type policyoptionsPolicyStatementDataSourceData struct {
	ID        types.String                            `tfsdk:"id"`
	Name      types.String                            `tfsdk:"name"`
	DynamicDB types.Bool                              `tfsdk:"dynamic_db"`
	From      *policyoptionsPolicyStatementBlockFrom  `tfsdk:"from"`
	To        *policyoptionsPolicyStatementBlockTo    `tfsdk:"to"`
	Then      *policyoptionsPolicyStatementBlockThen  `tfsdk:"then"`
	Term      []policyoptionsPolicyStatementBlockTerm `tfsdk:"term"`
}

func (dsc *policyoptionsPolicyStatementDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var name types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data policyoptionsPolicyStatementDataSourceData
	var rscData policyoptionsPolicyStatementData

	var _ resourceDataReadFrom1String = &rscData
	defaultDataSourceReadFromResource(
		ctx,
		dsc,
		[]string{
			name.ValueString(),
		},
		&data,
		&rscData,
		resp,
		fmt.Sprintf(dsc.junosName()+" %q doesn't exist", name.ValueString()),
	)
}

func (dscData *policyoptionsPolicyStatementDataSourceData) copyFromResourceData(data any) {
	rscData := data.(*policyoptionsPolicyStatementData)
	dscData.ID = rscData.ID
	dscData.Name = rscData.Name
	dscData.DynamicDB = rscData.DynamicDB
	dscData.From = rscData.From
	dscData.To = rscData.To
	dscData.Then = rscData.Then
	dscData.Term = rscData.Term
}
