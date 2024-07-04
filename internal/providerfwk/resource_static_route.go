package providerfwk

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                   = &staticRoute{}
	_ resource.ResourceWithConfigure      = &staticRoute{}
	_ resource.ResourceWithValidateConfig = &staticRoute{}
	_ resource.ResourceWithImportState    = &staticRoute{}
)

type staticRoute struct {
	client *junos.Client
}

func newStaticRouteResource() resource.Resource {
	return &staticRoute{}
}

func (rsc *staticRoute) typeName() string {
	return providerName + "_static_route"
}

func (rsc *staticRoute) junosName() string {
	return "static route"
}

func (rsc *staticRoute) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *staticRoute) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *staticRoute) Configure(
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

func (rsc *staticRoute) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<destination>" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"destination": schema.StringAttribute{
				Required:    true,
				Description: "Destination prefix.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					tfvalidator.StringCIDRNetwork(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for static route.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Description: "Remove inactive route from forwarding table.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"as_path_aggregator_address": schema.StringAttribute{
				Optional:    true,
				Description: "Address of BGP system to add AGGREGATOR path attribute to route.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
			"as_path_aggregator_as_number": schema.StringAttribute{
				Optional:    true,
				Description: "AS number to add AGGREGATOR path attribute to route.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d+(\.\d+)?$`),
						"must be in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format"),
				},
			},
			"as_path_atomic_aggregate": schema.BoolAttribute{
				Optional:    true,
				Description: "Add ATOMIC_AGGREGATE path attribute to route.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"as_path_origin": schema.StringAttribute{
				Optional:    true,
				Description: "Define origin.",
				Validators: []validator.String{
					stringvalidator.OneOf("egp", "igp", "incomplete"),
				},
			},
			"as_path_path": schema.StringAttribute{
				Optional:    true,
				Description: "Path to as-path.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"community": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "BGP community.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"discard": schema.BoolAttribute{
				Optional:    true,
				Description: "Drop packets to destination; send no ICMP unreachables.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"install": schema.BoolAttribute{
				Optional:    true,
				Description: "Install route into forwarding table.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_install": schema.BoolAttribute{
				Optional:    true,
				Description: "Don't install route into forwarding table.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"metric": schema.Int64Attribute{
				Optional:    true,
				Description: "Metric for static route.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"next_hop": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Next-hop to destination.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						stringvalidator.Any(
							tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
							tfvalidator.StringIPAddress(),
						),
					),
				},
			},
			"next_table": schema.StringAttribute{
				Optional:    true,
				Description: "Next hop to another table.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"passive": schema.BoolAttribute{
				Optional:    true,
				Description: "Retain inactive route in forwarding table.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"preference": schema.Int64Attribute{
				Optional:    true,
				Description: "Preference for aggregate route.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"readvertise": schema.BoolAttribute{
				Optional:    true,
				Description: "Mark route as eligible to be readvertised.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_readvertise": schema.BoolAttribute{
				Optional:    true,
				Description: "Don't mark route as eligible to be readvertised.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"receive": schema.BoolAttribute{
				Optional:    true,
				Description: "Install a receive route for the destination.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"reject": schema.BoolAttribute{
				Optional:    true,
				Description: "Drop packets to destination; send ICMP unreachables.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"resolve": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow resolution of indirectly connected next hops.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_resolve": schema.BoolAttribute{
				Optional:    true,
				Description: "Don't allow resolution of indirectly connected next hops.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"retain": schema.BoolAttribute{
				Optional:    true,
				Description: "Always keep route in forwarding table.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_retain": schema.BoolAttribute{
				Optional:    true,
				Description: "Don't always keep route in forwarding table.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"qualified_next_hop": schema.ListNestedBlock{
				Description: "For each `next_hop` with qualifiers.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"next_hop": schema.StringAttribute{
							Required:    true,
							Description: "Next-hop with qualifiers to destination.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.Any(
									tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
									tfvalidator.StringIPAddress(),
								),
							},
						},
						"interface": schema.StringAttribute{
							Optional:    true,
							Description: "Interface of qualified next hop.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
							},
						},
						"metric": schema.Int64Attribute{
							Optional:    true,
							Description: "Metric of qualified next hop.",
							Validators: []validator.Int64{
								int64validator.Between(0, 4294967295),
							},
						},
						"preference": schema.Int64Attribute{
							Optional:    true,
							Description: "Preference of qualified next hop.",
							Validators: []validator.Int64{
								int64validator.Between(0, 4294967295),
							},
						},
					},
				},
			},
		},
	}
}

type staticRouteData struct {
	ID                       types.String                       `tfsdk:"id"`
	Destination              types.String                       `tfsdk:"destination"`
	RoutingInstance          types.String                       `tfsdk:"routing_instance"`
	Active                   types.Bool                         `tfsdk:"active"`
	ASPathAggregatorAddress  types.String                       `tfsdk:"as_path_aggregator_address"`
	ASPathAggregatorASNumber types.String                       `tfsdk:"as_path_aggregator_as_number"`
	ASPathAtomicAggregate    types.Bool                         `tfsdk:"as_path_atomic_aggregate"`
	ASPathOrigin             types.String                       `tfsdk:"as_path_origin"`
	ASPathPath               types.String                       `tfsdk:"as_path_path"`
	Community                []types.String                     `tfsdk:"community"`
	Discard                  types.Bool                         `tfsdk:"discard"`
	Install                  types.Bool                         `tfsdk:"install"`
	NoInstall                types.Bool                         `tfsdk:"no_install"`
	Metric                   types.Int64                        `tfsdk:"metric"`
	NextHop                  []types.String                     `tfsdk:"next_hop"`
	NextTable                types.String                       `tfsdk:"next_table"`
	Passive                  types.Bool                         `tfsdk:"passive"`
	Preference               types.Int64                        `tfsdk:"preference"`
	Readvertise              types.Bool                         `tfsdk:"readvertise"`
	NoReadvertise            types.Bool                         `tfsdk:"no_readvertise"`
	Receive                  types.Bool                         `tfsdk:"receive"`
	Reject                   types.Bool                         `tfsdk:"reject"`
	Resolve                  types.Bool                         `tfsdk:"resolve"`
	NoResolve                types.Bool                         `tfsdk:"no_resolve"`
	Retain                   types.Bool                         `tfsdk:"retain"`
	NoRetain                 types.Bool                         `tfsdk:"no_retain"`
	QualifiedNextHop         []staticRouteBlockQualifiedNextHop `tfsdk:"qualified_next_hop"`
}

type staticRouteConfig struct {
	ID                       types.String `tfsdk:"id"`
	Destination              types.String `tfsdk:"destination"`
	RoutingInstance          types.String `tfsdk:"routing_instance"`
	Active                   types.Bool   `tfsdk:"active"`
	ASPathAggregatorAddress  types.String `tfsdk:"as_path_aggregator_address"`
	ASPathAggregatorASNumber types.String `tfsdk:"as_path_aggregator_as_number"`
	ASPathAtomicAggregate    types.Bool   `tfsdk:"as_path_atomic_aggregate"`
	ASPathOrigin             types.String `tfsdk:"as_path_origin"`
	ASPathPath               types.String `tfsdk:"as_path_path"`
	Community                types.List   `tfsdk:"community"`
	Discard                  types.Bool   `tfsdk:"discard"`
	Install                  types.Bool   `tfsdk:"install"`
	NoInstall                types.Bool   `tfsdk:"no_install"`
	Metric                   types.Int64  `tfsdk:"metric"`
	NextHop                  types.List   `tfsdk:"next_hop"`
	NextTable                types.String `tfsdk:"next_table"`
	Passive                  types.Bool   `tfsdk:"passive"`
	Preference               types.Int64  `tfsdk:"preference"`
	Readvertise              types.Bool   `tfsdk:"readvertise"`
	NoReadvertise            types.Bool   `tfsdk:"no_readvertise"`
	Receive                  types.Bool   `tfsdk:"receive"`
	Reject                   types.Bool   `tfsdk:"reject"`
	Resolve                  types.Bool   `tfsdk:"resolve"`
	NoResolve                types.Bool   `tfsdk:"no_resolve"`
	Retain                   types.Bool   `tfsdk:"retain"`
	NoRetain                 types.Bool   `tfsdk:"no_retain"`
	QualifiedNextHop         types.List   `tfsdk:"qualified_next_hop"`
}

type staticRouteBlockQualifiedNextHop struct {
	NextHop    types.String `tfsdk:"next_hop"`
	Interface  types.String `tfsdk:"interface"`
	Metric     types.Int64  `tfsdk:"metric"`
	Preference types.Int64  `tfsdk:"preference"`
}

func (rsc *staticRoute) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config staticRouteConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Active.IsNull() && !config.Active.IsUnknown() &&
		!config.Passive.IsNull() && !config.Passive.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("active"),
			tfdiag.ConflictConfigErrSummary,
			"active and passive cannot be configured together",
		)
	}
	if !config.ASPathAggregatorASNumber.IsNull() &&
		config.ASPathAggregatorAddress.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("as_path_aggregator_as_number"),
			tfdiag.MissingConfigErrSummary,
			"as_path_aggregator_address must be specified with as_path_aggregator_as_number",
		)
	}
	if !config.ASPathAggregatorAddress.IsNull() &&
		config.ASPathAggregatorASNumber.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("as_path_aggregator_address"),
			tfdiag.MissingConfigErrSummary,
			"as_path_aggregator_as_number must be specified with as_path_aggregator_address",
		)
	}
	if !config.Discard.IsNull() && !config.Discard.IsUnknown() {
		if !config.NextHop.IsNull() && !config.NextHop.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("discard"),
				tfdiag.ConflictConfigErrSummary,
				"discard and next_hop cannot be configured together",
			)
		}
		if !config.NextTable.IsNull() && !config.NextTable.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("discard"),
				tfdiag.ConflictConfigErrSummary,
				"discard and next_table cannot be configured together",
			)
		}
		if !config.QualifiedNextHop.IsNull() && !config.QualifiedNextHop.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("discard"),
				tfdiag.ConflictConfigErrSummary,
				"discard and qualified_next_hop cannot be configured together",
			)
		}
		if !config.Receive.IsNull() && !config.Receive.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("discard"),
				tfdiag.ConflictConfigErrSummary,
				"discard and receive cannot be configured together",
			)
		}
		if !config.Reject.IsNull() && !config.Reject.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("discard"),
				tfdiag.ConflictConfigErrSummary,
				"discard and reject cannot be configured together",
			)
		}
	}
	if !config.Install.IsNull() && !config.Install.IsUnknown() &&
		!config.NoInstall.IsNull() && !config.NoInstall.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("install"),
			tfdiag.ConflictConfigErrSummary,
			"install and no_install cannot be configured together",
		)
	}
	if !config.NextHop.IsNull() && !config.NextHop.IsUnknown() {
		if !config.NextTable.IsNull() && !config.NextTable.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("next_hop"),
				tfdiag.ConflictConfigErrSummary,
				"next_hop and next_table cannot be configured together",
			)
		}
		if !config.Receive.IsNull() && !config.Receive.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("next_hop"),
				tfdiag.ConflictConfigErrSummary,
				"next_hop and receive cannot be configured together",
			)
		}
		if !config.Reject.IsNull() && !config.Reject.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("next_hop"),
				tfdiag.ConflictConfigErrSummary,
				"next_hop and reject cannot be configured together",
			)
		}
	}
	if !config.NextTable.IsNull() && !config.NextTable.IsUnknown() {
		if !config.QualifiedNextHop.IsNull() && !config.QualifiedNextHop.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("next_table"),
				tfdiag.ConflictConfigErrSummary,
				"next_table and qualified_next_hop cannot be configured together",
			)
		}
		if !config.Receive.IsNull() && !config.Receive.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("next_table"),
				tfdiag.ConflictConfigErrSummary,
				"next_table and receive cannot be configured together",
			)
		}
		if !config.Reject.IsNull() && !config.Reject.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("next_hop"),
				tfdiag.ConflictConfigErrSummary,
				"next_table and reject cannot be configured together",
			)
		}
	}
	if !config.Readvertise.IsNull() && !config.Readvertise.IsUnknown() &&
		!config.NoReadvertise.IsNull() && !config.NoReadvertise.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("readvertise"),
			tfdiag.ConflictConfigErrSummary,
			"readvertise and no_readvertise cannot be configured together",
		)
	}
	if !config.Receive.IsNull() && !config.Receive.IsUnknown() {
		if !config.QualifiedNextHop.IsNull() && !config.QualifiedNextHop.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("receive"),
				tfdiag.ConflictConfigErrSummary,
				"receive and qualified_next_hop cannot be configured together",
			)
		}
		if !config.Reject.IsNull() && !config.Reject.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("receive"),
				tfdiag.ConflictConfigErrSummary,
				"receive and reject cannot be configured together",
			)
		}
	}
	if !config.Reject.IsNull() && !config.Reject.IsUnknown() {
		if !config.QualifiedNextHop.IsNull() && !config.QualifiedNextHop.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("reject"),
				tfdiag.ConflictConfigErrSummary,
				"reject and qualified_next_hop cannot be configured together",
			)
		}
	}
	if !config.Resolve.IsNull() && !config.Resolve.IsUnknown() {
		if !config.NoResolve.IsNull() && !config.NoResolve.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("resolve"),
				tfdiag.ConflictConfigErrSummary,
				"resolve and no_resolve cannot be configured together",
			)
		}
		if !config.Retain.IsNull() && !config.Retain.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("resolve"),
				tfdiag.ConflictConfigErrSummary,
				"resolve and retain cannot be configured together",
			)
		}
		if !config.NoRetain.IsNull() && !config.NoRetain.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("resolve"),
				tfdiag.ConflictConfigErrSummary,
				"resolve and no_retain cannot be configured together",
			)
		}
	}
	if !config.Retain.IsNull() && !config.Retain.IsUnknown() &&
		!config.NoRetain.IsNull() && !config.NoRetain.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("retain"),
			tfdiag.ConflictConfigErrSummary,
			"retain and no_retain cannot be configured together",
		)
	}
	if !config.QualifiedNextHop.IsNull() && !config.QualifiedNextHop.IsUnknown() {
		var configQualifiedNextHop []staticRouteBlockQualifiedNextHop
		asDiags := config.QualifiedNextHop.ElementsAs(ctx, &configQualifiedNextHop, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		qualifiedNextHopNextHop := make(map[string]struct{})
		for i, block := range configQualifiedNextHop {
			if block.NextHop.IsUnknown() {
				continue
			}
			nextHop := block.NextHop.ValueString()
			if _, ok := qualifiedNextHopNextHop[nextHop]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("qualified_next_hop").AtListIndex(i).AtName("next_hop"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple qualified_next_hop blocks with the same next_hop %q",
						nextHop),
				)
			}
			qualifiedNextHopNextHop[nextHop] = struct{}{}
		}
	}
}

func (rsc *staticRoute) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan staticRouteData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Destination.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("destination"),
			"Empty Destination",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "destination"),
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
			routeExists, err := checkStaticRouteExists(
				fnCtx,
				plan.Destination.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if routeExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsInRoutingInstanceMessage(rsc, plan.Destination, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						defaultResourceAlreadyExistsMessage(rsc, plan.Destination),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			routeExists, err := checkStaticRouteExists(
				fnCtx,
				plan.Destination.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !routeExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsInRoutingInstanceAfterCommitMessage(rsc, plan.Destination, v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Destination),
					)
				}

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *staticRoute) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data staticRouteData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Destination.ValueString(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *staticRoute) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state staticRouteData
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

func (rsc *staticRoute) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state staticRouteData
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

func (rsc *staticRoute) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data staticRouteData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <destination>"+junos.IDSeparator+"<routing_instance>)",
	)
}

func checkStaticRouteExists(
	_ context.Context, destination, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	switch routingInstance {
	case junos.DefaultW, "":
		showPrefix += junos.RoutingOptionsWS
		if strings.Contains(destination, ":") {
			showPrefix += junos.RibInet60WS
		}
	default:
		showPrefix += junos.RoutingInstancesWS + routingInstance + " " + junos.RoutingOptionsWS
		if strings.Contains(destination, ":") {
			showPrefix += "rib " + routingInstance + ".inet6.0 "
		}
	}
	showConfig, err := junSess.Command(showPrefix +
		"static route " + destination + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *staticRouteData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Destination.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Destination.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *staticRouteData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *staticRouteData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	switch routingInstance := rscData.RoutingInstance.ValueString(); routingInstance {
	case junos.DefaultW, "":
		setPrefix += junos.RoutingOptionsWS
		if strings.Contains(rscData.Destination.ValueString(), ":") {
			setPrefix += junos.RibInet60WS
		}
	default:
		setPrefix += junos.RoutingInstancesWS + routingInstance + " " + junos.RoutingOptionsWS
		if strings.Contains(rscData.Destination.ValueString(), ":") {
			setPrefix += "rib " + routingInstance + ".inet6.0 "
		}
	}
	setPrefix += "static route " + rscData.Destination.ValueString() + " "

	if rscData.Active.ValueBool() {
		configSet = append(configSet, setPrefix+"active")
	}
	if vNumber, vAddress := rscData.ASPathAggregatorASNumber.ValueString(),
		rscData.ASPathAggregatorAddress.ValueString(); vNumber != "" && vAddress != "" {
		configSet = append(configSet, setPrefix+"as-path aggregator "+vNumber+" "+vAddress)
	}
	if rscData.ASPathAtomicAggregate.ValueBool() {
		configSet = append(configSet, setPrefix+"as-path atomic-aggregate")
	}
	if v := rscData.ASPathOrigin.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"as-path origin "+v)
	}
	if v := rscData.ASPathPath.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"as-path path \""+v+"\"")
	}
	for _, v := range rscData.Community {
		configSet = append(configSet, setPrefix+"community \""+v.ValueString()+"\"")
	}
	if rscData.Discard.ValueBool() {
		configSet = append(configSet, setPrefix+"discard")
	}
	if rscData.Install.ValueBool() {
		configSet = append(configSet, setPrefix+"install")
	}
	if rscData.NoInstall.ValueBool() {
		configSet = append(configSet, setPrefix+"no-install")
	}
	if !rscData.Metric.IsNull() {
		configSet = append(configSet, setPrefix+"metric "+
			utils.ConvI64toa(rscData.Metric.ValueInt64()))
	}
	for _, v := range rscData.NextHop {
		configSet = append(configSet, setPrefix+"next-hop "+v.ValueString())
	}
	if v := rscData.NextTable.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"next-table \""+v+"\"")
	}
	if rscData.Passive.ValueBool() {
		configSet = append(configSet, setPrefix+"passive")
	}
	if !rscData.Preference.IsNull() {
		configSet = append(configSet, setPrefix+"preference "+
			utils.ConvI64toa(rscData.Preference.ValueInt64()))
	}
	qualifiedNextHopNextHop := make(map[string]struct{})
	for i, block := range rscData.QualifiedNextHop {
		nextHop := block.NextHop.ValueString()
		if _, ok := qualifiedNextHopNextHop[nextHop]; ok {
			return path.Root("qualified_next_hop").AtListIndex(i).AtName("next_hop"),
				fmt.Errorf("multiple qualified_next_hop blocks with the same next_hop %q", nextHop)
		}
		qualifiedNextHopNextHop[nextHop] = struct{}{}

		setPrefixQualifiedNextHop := setPrefix + "qualified-next-hop " + nextHop
		configSet = append(configSet, setPrefixQualifiedNextHop)
		setPrefixQualifiedNextHop += " "
		if v := block.Interface.ValueString(); v != "" {
			configSet = append(configSet, setPrefixQualifiedNextHop+"interface "+v)
		}
		if !block.Metric.IsNull() {
			configSet = append(configSet, setPrefixQualifiedNextHop+"metric "+
				utils.ConvI64toa(block.Metric.ValueInt64()))
		}
		if !block.Preference.IsNull() {
			configSet = append(configSet, setPrefixQualifiedNextHop+"preference "+
				utils.ConvI64toa(block.Preference.ValueInt64()))
		}
	}
	if rscData.Readvertise.ValueBool() {
		configSet = append(configSet, setPrefix+"readvertise")
	}
	if rscData.NoReadvertise.ValueBool() {
		configSet = append(configSet, setPrefix+"no-readvertise")
	}
	if rscData.Receive.ValueBool() {
		configSet = append(configSet, setPrefix+"receive")
	}
	if rscData.Reject.ValueBool() {
		configSet = append(configSet, setPrefix+"reject")
	}
	if rscData.Resolve.ValueBool() {
		configSet = append(configSet, setPrefix+"resolve")
	}
	if rscData.NoResolve.ValueBool() {
		configSet = append(configSet, setPrefix+"no-resolve")
	}
	if rscData.Retain.ValueBool() {
		configSet = append(configSet, setPrefix+"retain")
	}
	if rscData.NoRetain.ValueBool() {
		configSet = append(configSet, setPrefix+"no-retain")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *staticRouteData) read(
	_ context.Context, destination, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	switch routingInstance {
	case junos.DefaultW, "":
		showPrefix += junos.RoutingOptionsWS
		if strings.Contains(destination, ":") {
			showPrefix += junos.RibInet60WS
		}
	default:
		showPrefix += junos.RoutingInstancesWS + routingInstance + " " + junos.RoutingOptionsWS
		if strings.Contains(destination, ":") {
			showPrefix += "rib " + routingInstance + ".inet6.0 "
		}
	}
	showConfig, err := junSess.Command(showPrefix +
		"static route " + destination + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Destination = types.StringValue(destination)
		if routingInstance == "" {
			rscData.RoutingInstance = types.StringValue(junos.DefaultW)
		} else {
			rscData.RoutingInstance = types.StringValue(routingInstance)
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
			case itemTrim == "active":
				rscData.Active = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "as-path aggregator "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 2 { // <as_number> <address>
					return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "as-path aggregator", itemTrim)
				}
				rscData.ASPathAggregatorASNumber = types.StringValue(itemTrimFields[0])
				rscData.ASPathAggregatorAddress = types.StringValue(itemTrimFields[1])
			case itemTrim == "as-path atomic-aggregate":
				rscData.ASPathAtomicAggregate = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "as-path origin "):
				rscData.ASPathOrigin = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "as-path path "):
				rscData.ASPathPath = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "community "):
				rscData.Community = append(rscData.Community, types.StringValue(strings.Trim(itemTrim, "\"")))
			case itemTrim == junos.DiscardW:
				rscData.Discard = types.BoolValue(true)
			case itemTrim == "install":
				rscData.Install = types.BoolValue(true)
			case itemTrim == "no-install":
				rscData.NoInstall = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "metric "):
				rscData.Metric, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "next-hop "):
				rscData.NextHop = append(rscData.NextHop, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "next-table "):
				rscData.NextTable = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "passive":
				rscData.Passive = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "preference "):
				rscData.Preference, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "qualified-next-hop "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var qualifiedNextHop staticRouteBlockQualifiedNextHop
				rscData.QualifiedNextHop, qualifiedNextHop = tfdata.ExtractBlockWithTFTypesString(
					rscData.QualifiedNextHop, "NextHop", itemTrimFields[0],
				)
				qualifiedNextHop.NextHop = types.StringValue(itemTrimFields[0])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "interface "):
					qualifiedNextHop.Interface = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "metric "):
					qualifiedNextHop.Metric, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "preference "):
					qualifiedNextHop.Preference, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
				rscData.QualifiedNextHop = append(rscData.QualifiedNextHop, qualifiedNextHop)
			case itemTrim == "readvertise":
				rscData.Readvertise = types.BoolValue(true)
			case itemTrim == "no-readvertise":
				rscData.NoReadvertise = types.BoolValue(true)
			case itemTrim == "receive":
				rscData.Receive = types.BoolValue(true)
			case itemTrim == "reject":
				rscData.Reject = types.BoolValue(true)
			case itemTrim == "resolve":
				rscData.Resolve = types.BoolValue(true)
			case itemTrim == "no-resolve":
				rscData.NoResolve = types.BoolValue(true)
			case itemTrim == "retain":
				rscData.Retain = types.BoolValue(true)
			case itemTrim == "no-retain":
				rscData.NoRetain = types.BoolValue(true)
			}
		}
	}

	return nil
}

func (rscData *staticRouteData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	switch routingInstance := rscData.RoutingInstance.ValueString(); routingInstance {
	case junos.DefaultW, "":
		delPrefix += junos.RoutingOptionsWS
		if strings.Contains(rscData.Destination.ValueString(), ":") {
			delPrefix += junos.RibInet60WS
		}
	default:
		delPrefix += junos.RoutingInstancesWS + routingInstance + " " + junos.RoutingOptionsWS
		if strings.Contains(rscData.Destination.ValueString(), ":") {
			delPrefix += "rib " + routingInstance + ".inet6.0 "
		}
	}

	configSet := []string{
		delPrefix + "static route " + rscData.Destination.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
