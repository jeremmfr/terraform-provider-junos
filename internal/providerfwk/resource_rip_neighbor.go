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
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &ripNeighbor{}
	_ resource.ResourceWithConfigure      = &ripNeighbor{}
	_ resource.ResourceWithValidateConfig = &ripNeighbor{}
	_ resource.ResourceWithImportState    = &ripNeighbor{}
	_ resource.ResourceWithUpgradeState   = &ripNeighbor{}
)

type ripNeighbor struct {
	client *junos.Client
}

func newRipNeighborResource() resource.Resource {
	return &ripNeighbor{}
}

func (rsc *ripNeighbor) typeName() string {
	return providerName + "_rip_neighbor"
}

func (rsc *ripNeighbor) junosName() string {
	return "rip|ripng neighbor"
}

func (rsc *ripNeighbor) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *ripNeighbor) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *ripNeighbor) Configure(
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

func (rsc *ripNeighbor) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<name>" + junos.IDSeparator + "<group>" + junos.IDSeparator + "<routing_instance>` or " +
					"`<name>" + junos.IDSeparator + "<group>" + junos.IDSeparator + "ng" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Logical interface name or `all`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					stringvalidator.Any(
						tfvalidator.String1DotCount(),
						stringvalidator.OneOf("all"),
					),
				},
			},
			"group": schema.StringAttribute{
				Required:    true,
				Description: "Name of RIP or RIPng group for this neighbor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 48),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"ng": schema.BoolAttribute{
				Optional:    true,
				Description: "Protocol `ripng` instead of `rip`.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for RIP neighbor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"any_sender": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable strict checks on sender address.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"authentication_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Authentication key (password).",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"authentication_type": schema.StringAttribute{
				Optional:    true,
				Description: "Authentication type.",
				Validators: []validator.String{
					stringvalidator.OneOf("md5", "none", "simple"),
				},
			},
			"check_zero": schema.BoolAttribute{
				Optional:    true,
				Description: "Check reserved fields on incoming RIPv1 packets.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_check_zero": schema.BoolAttribute{
				Optional:    true,
				Description: "Don't check reserved fields on incoming RIPv1 packets.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"demand_circuit": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable demand circuit.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"dynamic_peers": schema.BoolAttribute{
				Optional:    true,
				Description: "Learn peers dynamically on a p2mp interface.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"import": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Import policy.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 63),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
			"interface_type_p2mp": schema.BoolAttribute{
				Optional:    true,
				Description: "Point-to-multipoint link.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"max_retrans_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum time to re-transmit a message in demand-circuit.",
				Validators: []validator.Int64{
					int64validator.Between(5, 180),
				},
			},
			"message_size": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of route entries per update message.",
				Validators: []validator.Int64{
					int64validator.Between(25, 255),
				},
			},
			"metric_in": schema.Int64Attribute{
				Optional:    true,
				Description: "Metric value to add to incoming routes.",
				Validators: []validator.Int64{
					int64validator.Between(1, 15),
				},
			},
			"peer": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "P2MP peer.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringIPAddress(),
					),
				},
			},
			"receive": schema.StringAttribute{
				Optional:    true,
				Description: "Configure RIP receive options.",
				Validators: []validator.String{
					stringvalidator.OneOf("both", "none", "version-1", "version-2"),
				},
			},
			"route_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Delay before routes time out (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(30, 360),
				},
			},
			"send": schema.StringAttribute{
				Optional:    true,
				Description: "Configure RIP send options.",
				Validators: []validator.String{
					stringvalidator.OneOf("broadcast", "multicast", "none", "version-1"),
				},
			},
			"update_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Interval between regular route updates (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(10, 60),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"authentication_selective_md5": schema.ListNestedBlock{
				Description: "For each key_id, MD5 authentication key.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key_id": schema.Int64Attribute{
							Required:    true,
							Description: "Key ID for MD5 authentication.",
							Validators: []validator.Int64{
								int64validator.Between(0, 255),
							},
						},
						"key": schema.StringAttribute{
							Required:    true,
							Sensitive:   true,
							Description: "MD5 authentication key value.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"start_time": schema.StringAttribute{
							Optional:    true,
							Description: "Start time for key transmission.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(
									`^\d{4}\-\d\d?\-\d\d?\.\d{2}:\d{2}:\d{2}$`),
									"must be in the format 'YYYY-MM-DD.HH:MM:SS'",
								),
							},
						},
					},
				},
			},
			"bfd_liveness_detection": ripBlockBfdLivenessDetection{}.resourceSchema(),
		},
	}
}

type ripNeighborData struct {
	ID                         types.String                                 `tfsdk:"id"`
	Name                       types.String                                 `tfsdk:"name"`
	Group                      types.String                                 `tfsdk:"group"`
	Ng                         types.Bool                                   `tfsdk:"ng"`
	RoutingInstance            types.String                                 `tfsdk:"routing_instance"`
	AnySender                  types.Bool                                   `tfsdk:"any_sender"`
	AuthenticationKey          types.String                                 `tfsdk:"authentication_key"`
	AuthenticationType         types.String                                 `tfsdk:"authentication_type"`
	CheckZero                  types.Bool                                   `tfsdk:"check_zero"`
	NoCheckZero                types.Bool                                   `tfsdk:"no_check_zero"`
	DemandCircuit              types.Bool                                   `tfsdk:"demand_circuit"`
	DynamicPeers               types.Bool                                   `tfsdk:"dynamic_peers"`
	Import                     []types.String                               `tfsdk:"import"`
	InterfaceTypeP2mp          types.Bool                                   `tfsdk:"interface_type_p2mp"`
	MaxRetransTime             types.Int64                                  `tfsdk:"max_retrans_time"`
	MessageSize                types.Int64                                  `tfsdk:"message_size"`
	MetricIn                   types.Int64                                  `tfsdk:"metric_in"`
	Peer                       []types.String                               `tfsdk:"peer"`
	Receive                    types.String                                 `tfsdk:"receive"`
	RouteTimeout               types.Int64                                  `tfsdk:"route_timeout"`
	Send                       types.String                                 `tfsdk:"send"`
	UpdateInterval             types.Int64                                  `tfsdk:"update_interval"`
	AuthenticationSelectiveMD5 []ripNeighborBlockAuthenticationSelectiveMd5 `tfsdk:"authentication_selective_md5"`
	BfdLivenessDetection       *ripBlockBfdLivenessDetection                `tfsdk:"bfd_liveness_detection"`
}

type ripNeighborConfig struct {
	ID                         types.String                  `tfsdk:"id"`
	Name                       types.String                  `tfsdk:"name"`
	Group                      types.String                  `tfsdk:"group"`
	Ng                         types.Bool                    `tfsdk:"ng"`
	RoutingInstance            types.String                  `tfsdk:"routing_instance"`
	AnySender                  types.Bool                    `tfsdk:"any_sender"`
	AuthenticationKey          types.String                  `tfsdk:"authentication_key"`
	AuthenticationType         types.String                  `tfsdk:"authentication_type"`
	CheckZero                  types.Bool                    `tfsdk:"check_zero"`
	NoCheckZero                types.Bool                    `tfsdk:"no_check_zero"`
	DemandCircuit              types.Bool                    `tfsdk:"demand_circuit"`
	DynamicPeers               types.Bool                    `tfsdk:"dynamic_peers"`
	Import                     types.List                    `tfsdk:"import"`
	InterfaceTypeP2mp          types.Bool                    `tfsdk:"interface_type_p2mp"`
	MaxRetransTime             types.Int64                   `tfsdk:"max_retrans_time"`
	MessageSize                types.Int64                   `tfsdk:"message_size"`
	MetricIn                   types.Int64                   `tfsdk:"metric_in"`
	Peer                       types.Set                     `tfsdk:"peer"`
	Receive                    types.String                  `tfsdk:"receive"`
	RouteTimeout               types.Int64                   `tfsdk:"route_timeout"`
	Send                       types.String                  `tfsdk:"send"`
	UpdateInterval             types.Int64                   `tfsdk:"update_interval"`
	AuthenticationSelectiveMd5 types.List                    `tfsdk:"authentication_selective_md5"`
	BfdLivenessDetection       *ripBlockBfdLivenessDetection `tfsdk:"bfd_liveness_detection"`
}

type ripNeighborBlockAuthenticationSelectiveMd5 struct {
	KeyID     types.Int64  `tfsdk:"key_id"`
	Key       types.String `tfsdk:"key"`
	StartTime types.String `tfsdk:"start_time"`
}

func (rsc *ripNeighbor) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config ripNeighborConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Ng.IsNull() && !config.Ng.IsUnknown() {
		if !config.AnySender.IsNull() && !config.AnySender.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("any_sender"),
				tfdiag.ConflictConfigErrSummary,
				"ng and any_sender cannot be configured together",
			)
		}
		if !config.AuthenticationKey.IsNull() && !config.AuthenticationKey.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_key"),
				tfdiag.ConflictConfigErrSummary,
				"ng and authentication_key cannot be configured together",
			)
		}
		if !config.AuthenticationType.IsNull() && !config.AuthenticationType.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_type"),
				tfdiag.ConflictConfigErrSummary,
				"ng and authentication_type cannot be configured together",
			)
		}
		if !config.CheckZero.IsNull() && !config.CheckZero.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("check_zero"),
				tfdiag.ConflictConfigErrSummary,
				"ng and check_zero cannot be configured together",
			)
		}
		if !config.NoCheckZero.IsNull() && !config.NoCheckZero.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("no_check_zero"),
				tfdiag.ConflictConfigErrSummary,
				"ng and no_check_zero cannot be configured together",
			)
		}
		if !config.DemandCircuit.IsNull() && !config.DemandCircuit.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("demand_circuit"),
				tfdiag.ConflictConfigErrSummary,
				"ng and demand_circuit cannot be configured together",
			)
		}
		if !config.DynamicPeers.IsNull() && !config.DynamicPeers.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("dynamic_peers"),
				tfdiag.ConflictConfigErrSummary,
				"ng and dynamic_peers cannot be configured together",
			)
		}
		if !config.InterfaceTypeP2mp.IsNull() && !config.InterfaceTypeP2mp.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("interface_type_p2mp"),
				tfdiag.ConflictConfigErrSummary,
				"ng and interface_type_p2mp cannot be configured together",
			)
		}
		if !config.MaxRetransTime.IsNull() && !config.MaxRetransTime.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("max_retrans_time"),
				tfdiag.ConflictConfigErrSummary,
				"ng and max_retrans_time cannot be configured together",
			)
		}
		if !config.MessageSize.IsNull() && !config.MessageSize.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("message_size"),
				tfdiag.ConflictConfigErrSummary,
				"ng and message_size cannot be configured together",
			)
		}
		if !config.Peer.IsNull() && !config.Peer.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("peer"),
				tfdiag.ConflictConfigErrSummary,
				"ng and peer cannot be configured together",
			)
		}
		if !config.UpdateInterval.IsNull() && !config.UpdateInterval.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("update_interval"),
				tfdiag.ConflictConfigErrSummary,
				"ng and update_interval cannot be configured together",
			)
		}
		if !config.AuthenticationSelectiveMd5.IsNull() && !config.AuthenticationSelectiveMd5.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_selective_md5"),
				tfdiag.ConflictConfigErrSummary,
				"ng and authentication_selective_md5 cannot be configured together",
			)
		}
		if config.BfdLivenessDetection != nil && config.BfdLivenessDetection.hasKnownValue() {
			resp.Diagnostics.AddAttributeError(
				path.Root("bfd_liveness_detection"),
				tfdiag.ConflictConfigErrSummary,
				"ng and bfd_liveness_detection cannot be configured together",
			)
		}
	}
	if !config.AuthenticationSelectiveMd5.IsNull() && !config.AuthenticationSelectiveMd5.IsUnknown() {
		if !config.AuthenticationKey.IsNull() && !config.AuthenticationKey.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_key"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_selective_md5 and authentication_key cannot be configured together",
			)
		}
		if !config.AuthenticationType.IsNull() && !config.AuthenticationType.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_type"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_selective_md5 and authentication_type cannot be configured together",
			)
		}
	}
	if !config.CheckZero.IsNull() && !config.CheckZero.IsUnknown() &&
		!config.NoCheckZero.IsNull() && !config.NoCheckZero.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("no_check_zero"),
			tfdiag.ConflictConfigErrSummary,
			"check_zero and no_check_zero cannot be configured together",
		)
	}
	if !config.DynamicPeers.IsNull() && !config.DynamicPeers.IsUnknown() &&
		config.InterfaceTypeP2mp.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dynamic_peers"),
			tfdiag.MissingConfigErrSummary,
			"interface_type_p2mp must be specified with dynamic_peers",
		)
	}
	if !config.Peer.IsNull() && !config.Peer.IsUnknown() &&
		config.InterfaceTypeP2mp.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dynamic_peers"),
			tfdiag.MissingConfigErrSummary,
			"interface_type_p2mp must be specified with peer",
		)
	}
	if !config.AuthenticationSelectiveMd5.IsNull() && !config.AuthenticationSelectiveMd5.IsUnknown() {
		var configAuthenticationSelectiveMd5 []ripNeighborBlockAuthenticationSelectiveMd5
		asDiags := config.AuthenticationSelectiveMd5.ElementsAs(ctx, &configAuthenticationSelectiveMd5, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		authenticationSelectiveMD5KeyID := make(map[int64]struct{})
		for i, block := range configAuthenticationSelectiveMd5 {
			if block.KeyID.IsUnknown() {
				continue
			}
			keyID := block.KeyID.ValueInt64()
			if _, ok := authenticationSelectiveMD5KeyID[keyID]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("authentication_selective_md5").AtListIndex(i).AtName("key_id"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple authentication_selective_md5 blocks with the same key_id %d", keyID),
				)
			}
			authenticationSelectiveMD5KeyID[keyID] = struct{}{}
		}
	}
	if config.BfdLivenessDetection != nil && config.BfdLivenessDetection.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("bfd_liveness_detection").AtName("*"),
			tfdiag.MissingConfigErrSummary,
			"bfd_liveness_detection block is empty",
		)
	}
}

func (rsc *ripNeighbor) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan ripNeighborData
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
	if plan.Group.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("group"),
			"Empty Group",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "group"),
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
			groupExists, err := checkRipGroupExists(
				fnCtx,
				plan.Group.ValueString(),
				plan.Ng.ValueBool(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			protoPrefix := "rip "
			if plan.Ng.ValueBool() {
				protoPrefix = "ripng "
			}
			if !groupExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddAttributeError(
						path.Root("group"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf(protoPrefix+"group %q doesn't exist in routing-instance %q",
							plan.Group.ValueString(), v),
					)
				} else {
					resp.Diagnostics.AddAttributeError(
						path.Root("group"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf(protoPrefix+"group %q doesn't exist",
							plan.Group.ValueString()),
					)
				}

				return false
			}
			neighborExists, err := checkRipNeighborExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.Group.ValueString(),
				plan.Ng.ValueBool(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if neighborExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf(protoPrefix+"neighbor %q in group %q already exists in routing-instance %q",
							plan.Name.ValueString(), plan.Group.ValueString(), v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf(protoPrefix+"neighbor %q in group %q already exists",
							plan.Name.ValueString(), plan.Group.ValueString()),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			neighborExists, err := checkRipNeighborExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.Group.ValueString(),
				plan.Ng.ValueBool(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !neighborExists {
				protoPrefix := "rip "
				if plan.Ng.ValueBool() {
					protoPrefix = "ripng "
				}
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						fmt.Sprintf(protoPrefix+"neighbor %q does not exists in group %q in routing-instance %q after commit "+
							"=> check your config", plan.Name.ValueString(), plan.Group.ValueString(), v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						fmt.Sprintf(protoPrefix+"neighbor %q does not exists in group %q after commit "+
							"=> check your config", plan.Name.ValueString(), plan.Group.ValueString()),
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

func (rsc *ripNeighbor) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data ripNeighborData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String1Bool1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.Group.ValueString(),
			state.Ng.ValueBool(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *ripNeighbor) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state ripNeighborData
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

func (rsc *ripNeighbor) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state ripNeighborData
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

func (rsc *ripNeighbor) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	var data ripNeighborData
	idSplit := strings.Split(req.ID, junos.IDSeparator)
	switch {
	case len(idSplit) < 3:
		resp.Diagnostics.AddError(
			"Bad ID Format",
			fmt.Sprintf("missing element(s) in id with separator %q", junos.IDSeparator),
		)

		return
	case len(idSplit) == 3:
		if err := data.read(ctx, idSplit[0], idSplit[1], false, idSplit[2], junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
	default:
		if idSplit[2] != "ng" {
			resp.Diagnostics.AddError(
				"Bad ID Format",
				"id must be "+
					"<name>"+junos.IDSeparator+"<group>"+junos.IDSeparator+"<routing_instance> or "+
					"<name>"+junos.IDSeparator+"<group>"+junos.IDSeparator+"ng"+junos.IDSeparator+"<routing_instance>",
			)

			return
		}
		if err := data.read(ctx, idSplit[0], idSplit[1], true, idSplit[3], junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be "+
				"<name>"+junos.IDSeparator+"<group>"+junos.IDSeparator+"<routing_instance> or "+
				"<name>"+junos.IDSeparator+"<group>"+junos.IDSeparator+"ng"+junos.IDSeparator+"<routing_instance>)",
		)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkRipNeighborExists(
	_ context.Context, name, group string, ng bool, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	if ng {
		showPrefix += "protocols ripng "
	} else {
		showPrefix += "protocols rip "
	}
	showConfig, err := junSess.Command(showPrefix +
		"group \"" + group + "\" neighbor " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *ripNeighborData) fillID() {
	idPrefix := rscData.Name.ValueString() + junos.IDSeparator + rscData.Group.ValueString()
	if rscData.Ng.ValueBool() {
		idPrefix += junos.IDSeparator + "ng"
	}
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(idPrefix + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(idPrefix + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *ripNeighborData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *ripNeighborData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	if rscData.Ng.ValueBool() {
		setPrefix += "protocols ripng "
	} else {
		setPrefix += "protocols rip "
	}
	setPrefix += "group \"" + rscData.Group.ValueString() + "\" neighbor " + rscData.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if rscData.AnySender.ValueBool() {
		configSet = append(configSet, setPrefix+"any-sender")
	}
	if v := rscData.AuthenticationKey.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication-key \""+v+"\"")
	}
	if v := rscData.AuthenticationType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication-type "+v)
	}
	if rscData.CheckZero.ValueBool() {
		configSet = append(configSet, setPrefix+"check-zero")
	}
	if rscData.NoCheckZero.ValueBool() {
		configSet = append(configSet, setPrefix+"no-check-zero")
	}
	if rscData.DemandCircuit.ValueBool() {
		configSet = append(configSet, setPrefix+"demand-circuit")
	}
	if rscData.DynamicPeers.ValueBool() {
		configSet = append(configSet, setPrefix+"dynamic-peers")
	}
	for _, v := range rscData.Import {
		configSet = append(configSet, setPrefix+"import "+v.ValueString())
	}
	if rscData.InterfaceTypeP2mp.ValueBool() {
		configSet = append(configSet, setPrefix+"interface-type p2mp")
	}
	if !rscData.MaxRetransTime.IsNull() {
		configSet = append(configSet, setPrefix+"max-retrans-time "+
			utils.ConvI64toa(rscData.MaxRetransTime.ValueInt64()))
	}
	if !rscData.MessageSize.IsNull() {
		configSet = append(configSet, setPrefix+"message-size "+
			utils.ConvI64toa(rscData.MessageSize.ValueInt64()))
	}
	if !rscData.MetricIn.IsNull() {
		configSet = append(configSet, setPrefix+"metric-in "+
			utils.ConvI64toa(rscData.MetricIn.ValueInt64()))
	}
	for _, v := range rscData.Peer {
		configSet = append(configSet, setPrefix+"peer "+v.ValueString())
	}
	if v := rscData.Receive.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"receive "+v)
	}
	if !rscData.RouteTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"route-timeout "+
			utils.ConvI64toa(rscData.RouteTimeout.ValueInt64()))
	}
	if v := rscData.Send.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"send "+v)
	}
	if !rscData.UpdateInterval.IsNull() {
		configSet = append(configSet, setPrefix+"update-interval "+
			utils.ConvI64toa(rscData.UpdateInterval.ValueInt64()))
	}
	authenticationSelectiveMD5KeyID := make(map[int64]struct{})
	for i, block := range rscData.AuthenticationSelectiveMD5 {
		keyID := block.KeyID.ValueInt64()
		if _, ok := authenticationSelectiveMD5KeyID[keyID]; ok {
			return path.Root("authentication_selective_md5").AtListIndex(i).AtName("key_id"),
				fmt.Errorf("multiple authentication_selective_md5 blocks with the same key_id %d", keyID)
		}
		authenticationSelectiveMD5KeyID[keyID] = struct{}{}

		configSet = append(configSet, setPrefix+"authentication-selective-md5 "+
			utils.ConvI64toa(keyID)+" key \""+block.Key.ValueString()+"\"")
		if v := block.StartTime.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"authentication-selective-md5 "+
				utils.ConvI64toa(keyID)+" start-time "+v)
		}
	}
	if rscData.BfdLivenessDetection != nil {
		if rscData.BfdLivenessDetection.isEmpty() {
			return path.Root("bfd_liveness_detection").AtName("*"),
				errors.New("bfd_liveness_detection block is empty")
		}

		configSet = append(configSet, rscData.BfdLivenessDetection.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *ripNeighborData) read(
	_ context.Context, name, group string, ng bool, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	if ng {
		showPrefix += "protocols ripng "
	} else {
		showPrefix += "protocols rip "
	}
	showConfig, err := junSess.Command(showPrefix +
		"group \"" + group + "\" neighbor " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.Group = types.StringValue(group)
		if ng {
			rscData.Ng = types.BoolValue(true)
		}
		rscData.RoutingInstance = types.StringValue(routingInstance)
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
			case itemTrim == "any-sender":
				rscData.AnySender = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "authentication-key "):
				rscData.AuthenticationKey, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "authentication-key")
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "authentication-type "):
				rscData.AuthenticationType = types.StringValue(itemTrim)
			case itemTrim == "check-zero":
				rscData.CheckZero = types.BoolValue(true)
			case itemTrim == "demand-circuit":
				rscData.DemandCircuit = types.BoolValue(true)
			case itemTrim == "dynamic-peers":
				rscData.DynamicPeers = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "import "):
				rscData.Import = append(rscData.Import, types.StringValue(itemTrim))
			case itemTrim == "interface-type p2mp":
				rscData.InterfaceTypeP2mp = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "max-retrans-time "):
				rscData.MaxRetransTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "message-size "):
				rscData.MessageSize, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "metric-in "):
				rscData.MetricIn, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "no-check-zero":
				rscData.NoCheckZero = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "peer "):
				rscData.Peer = append(rscData.Peer, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "receive "):
				rscData.Receive = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "route-timeout "):
				rscData.RouteTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "send "):
				rscData.Send = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "update-interval "):
				rscData.UpdateInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "authentication-selective-md5 "):
				itemTrimFields := strings.Split(itemTrim, " ")
				keyID, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
				if err != nil {
					return err
				}
				var authenticationSelectiveMD5 ripNeighborBlockAuthenticationSelectiveMd5
				rscData.AuthenticationSelectiveMD5, authenticationSelectiveMD5 = tfdata.ExtractBlockWithTFTypesInt64(
					rscData.AuthenticationSelectiveMD5, "KeyID", keyID.ValueInt64(),
				)
				authenticationSelectiveMD5.KeyID = keyID
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "key "):
					authenticationSelectiveMD5.Key, err = tfdata.JunosDecode(
						strings.Trim(itemTrim, "\""),
						"authentication-selective-md5 key",
					)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "start-time "):
					authenticationSelectiveMD5.StartTime = types.StringValue(strings.Split(strings.Trim(itemTrim, "\""), " ")[0])
				}
				rscData.AuthenticationSelectiveMD5 = append(rscData.AuthenticationSelectiveMD5, authenticationSelectiveMD5)
			case balt.CutPrefixInString(&itemTrim, "bfd-liveness-detection "):
				if rscData.BfdLivenessDetection == nil {
					rscData.BfdLivenessDetection = &ripBlockBfdLivenessDetection{}
				}
				if err := rscData.BfdLivenessDetection.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *ripNeighborData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	if rscData.Ng.ValueBool() {
		delPrefix += "protocols ripng "
	} else {
		delPrefix += "protocols rip "
	}

	configSet := []string{
		delPrefix + "group \"" + rscData.Group.ValueString() + "\" neighbor " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
