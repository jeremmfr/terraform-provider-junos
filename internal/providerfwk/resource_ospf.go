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
	_ resource.Resource                   = &ospf{}
	_ resource.ResourceWithConfigure      = &ospf{}
	_ resource.ResourceWithValidateConfig = &ospf{}
	_ resource.ResourceWithImportState    = &ospf{}
	_ resource.ResourceWithUpgradeState   = &ospf{}
)

type ospf struct {
	client *junos.Client
}

func newOspfResource() resource.Resource {
	return &ospf{}
}

func (rsc *ospf) typeName() string {
	return providerName + "_ospf"
}

func (rsc *ospf) junosName() string {
	return "protocols ospf|ospf3"
}

func (rsc *ospf) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *ospf) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *ospf) Configure(
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

func (rsc *ospf) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<version>" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("v2"),
				Description: "Version of ospf.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("v2", "v3"),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for ospf protocol if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable OSPF.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"domain_id": schema.StringAttribute{
				Optional:    true,
				Description: "Configure domain ID.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"export": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Export policy.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"external_preference": schema.Int64Attribute{
				Optional:    true,
				Description: "Preference of external routes.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"forwarding_address_to_broadcast": schema.BoolAttribute{
				Optional:    true,
				Description: "Set forwarding address in Type 5 LSA in broadcast network.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"import": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Import policy (for external routes or setting priority).",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"labeled_preference": schema.Int64Attribute{
				Optional:    true,
				Description: "Preference of labeled routes.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"lsa_refresh_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "SA refresh interval (minutes).",
				Validators: []validator.Int64{
					int64validator.Between(25, 50),
				},
			},
			"no_nssa_abr": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable full NSSA functionality at ABR.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"no_rfc1583": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable RFC1583 compatibility.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"preference": schema.Int64Attribute{
				Optional:    true,
				Description: "Preference of internal routes.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"prefix_export_limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of prefixes that can be exported.",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"reference_bandwidth": schema.StringAttribute{
				Optional:    true,
				Description: "Bandwidth for calculating metric defaults.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(\d)+(m|k|g)?$`),
						`must be a bandwidth ^(\d)+(m|k|g)?$`),
				},
			},
			"rib_group": schema.StringAttribute{
				Optional:    true,
				Description: "Routing table group for importing OSPF routes.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"sham_link": schema.BoolAttribute{
				Optional:    true,
				Description: "Configure parameters for sham links.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"sham_link_local": schema.StringAttribute{
				Optional:    true,
				Description: "Local sham link endpoint address.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"database_protection": schema.SingleNestedBlock{
				Description: "Declare `database-protection` configuration.",
				Attributes: map[string]schema.Attribute{
					"maximum_lsa": schema.Int64Attribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Maximum allowed non self-generated LSAs.",
						Validators: []validator.Int64{
							int64validator.Between(1, 1000000),
						},
					},
					"ignore_count": schema.Int64Attribute{
						Optional:    true,
						Description: "Maximum number of times to go into ignore state.",
						Validators: []validator.Int64{
							int64validator.Between(1, 32),
						},
					},
					"ignore_time": schema.Int64Attribute{
						Optional:    true,
						Description: "Time to stay in ignore state and ignore all neighbors.",
						Validators: []validator.Int64{
							int64validator.Between(30, 3600),
						},
					},
					"reset_time": schema.Int64Attribute{
						Optional:    true,
						Description: "Time after which the ignore count gets reset to zero.",
						Validators: []validator.Int64{
							int64validator.Between(60, 86400),
						},
					},
					"warning_only": schema.BoolAttribute{
						Optional:    true,
						Description: "Emit only a warning when LSA maximum limit is exceeded.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"warning_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Percentage of LSA maximum above which to trigger warning (percent).",
						Validators: []validator.Int64{
							int64validator.Between(30, 100),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"graceful_restart": schema.SingleNestedBlock{
				Description: "Declare `graceful-restart` configuration.",
				Attributes: map[string]schema.Attribute{
					"disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable OSPF graceful restart capability.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"helper_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable graceful restart helper capability.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"helper_disable_type": schema.StringAttribute{
						Optional:    true,
						Description: "Disable graceful restart helper capability for specific type.",
						Validators: []validator.String{
							stringvalidator.OneOf("both", "restart-signaling", "standard"),
						},
					},
					"no_strict_lsa_checking": schema.BoolAttribute{
						Optional:    true,
						Description: "Do not abort graceful helper mode upon LSA changes.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"notify_duration": schema.Int64Attribute{
						Optional:    true,
						Description: "Time to send all max-aged grace LSAs (seconds).",
						Validators: []validator.Int64{
							int64validator.Between(1, 3600),
						},
					},
					"restart_duration": schema.Int64Attribute{
						Optional:    true,
						Description: "Time for all neighbors to become full (seconds).",
						Validators: []validator.Int64{
							int64validator.Between(1, 3600),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"overload": schema.SingleNestedBlock{
				Description: "Set the overload mode (repel transit traffic).",
				Attributes: map[string]schema.Attribute{
					"allow_route_leaking": schema.BoolAttribute{
						Optional:    true,
						Description: "Allow routes to be leaked when overload is configured.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"as_external": schema.BoolAttribute{
						Optional:    true,
						Description: "Advertise As External with maximum usable metric.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"stub_network": schema.BoolAttribute{
						Optional:    true,
						Description: "Advertise Stub Network with maximum metric.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Time after which overload mode is reset (seconds).",
						Validators: []validator.Int64{
							int64validator.Between(60, 1800),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"spf_options": schema.SingleNestedBlock{
				Description: "Declare `spf-options` configuration.",
				Attributes: map[string]schema.Attribute{
					"delay": schema.Int64Attribute{
						Optional:    true,
						Description: "Time to wait before running an SPF (milliseconds).",
						Validators: []validator.Int64{
							int64validator.Between(50, 8000),
						},
					},
					"holddown": schema.Int64Attribute{
						Optional:    true,
						Description: "Time to hold down before running an SPF (milliseconds).",
						Validators: []validator.Int64{
							int64validator.Between(2000, 20000),
						},
					},
					"no_ignore_our_externals": schema.BoolAttribute{
						Optional:    true,
						Description: "Do not ignore self-generated external and NSSA LSAs.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"rapid_runs": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of maximum rapid SPF runs before holddown.",
						Validators: []validator.Int64{
							int64validator.Between(1, 10),
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

type ospfData struct {
	ID                           types.String                 `tfsdk:"id"`
	Version                      types.String                 `tfsdk:"version"`
	RoutingInstance              types.String                 `tfsdk:"routing_instance"`
	Disable                      types.Bool                   `tfsdk:"disable"`
	DomainID                     types.String                 `tfsdk:"domain_id"`
	Export                       []types.String               `tfsdk:"export"`
	ExternalPreference           types.Int64                  `tfsdk:"external_preference"`
	ForwardingAddressToBroadcast types.Bool                   `tfsdk:"forwarding_address_to_broadcast"`
	Import                       []types.String               `tfsdk:"import"`
	LabeledPreference            types.Int64                  `tfsdk:"labeled_preference"`
	LsaRefreshInterval           types.Int64                  `tfsdk:"lsa_refresh_interval"`
	NoNssaAbr                    types.Bool                   `tfsdk:"no_nssa_abr"`
	NoRfc1583                    types.Bool                   `tfsdk:"no_rfc1583"`
	Preference                   types.Int64                  `tfsdk:"preference"`
	PrefixExportLimit            types.Int64                  `tfsdk:"prefix_export_limit"`
	ReferenceBandwidth           types.String                 `tfsdk:"reference_bandwidth"`
	RibGroup                     types.String                 `tfsdk:"rib_group"`
	ShamLink                     types.Bool                   `tfsdk:"sham_link"`
	ShamLinkLocal                types.String                 `tfsdk:"sham_link_local"`
	DatabaseProtection           *ospfBlockDatabaseProtection `tfsdk:"database_protection"`
	GracefulRestart              *ospfBlockGracefulRestart    `tfsdk:"graceful_restart"`
	Overload                     *ospfBlockOverload           `tfsdk:"overload"`
	SpfOptions                   *ospfBlockSpfOptions         `tfsdk:"spf_options"`
}

type ospfConfig struct {
	ID                           types.String                 `tfsdk:"id"`
	Version                      types.String                 `tfsdk:"version"`
	RoutingInstance              types.String                 `tfsdk:"routing_instance"`
	Disable                      types.Bool                   `tfsdk:"disable"`
	DomainID                     types.String                 `tfsdk:"domain_id"`
	Export                       types.List                   `tfsdk:"export"`
	ExternalPreference           types.Int64                  `tfsdk:"external_preference"`
	ForwardingAddressToBroadcast types.Bool                   `tfsdk:"forwarding_address_to_broadcast"`
	Import                       types.List                   `tfsdk:"import"`
	LabeledPreference            types.Int64                  `tfsdk:"labeled_preference"`
	LsaRefreshInterval           types.Int64                  `tfsdk:"lsa_refresh_interval"`
	NoNssaAbr                    types.Bool                   `tfsdk:"no_nssa_abr"`
	NoRfc1583                    types.Bool                   `tfsdk:"no_rfc1583"`
	Preference                   types.Int64                  `tfsdk:"preference"`
	PrefixExportLimit            types.Int64                  `tfsdk:"prefix_export_limit"`
	ReferenceBandwidth           types.String                 `tfsdk:"reference_bandwidth"`
	RibGroup                     types.String                 `tfsdk:"rib_group"`
	ShamLink                     types.Bool                   `tfsdk:"sham_link"`
	ShamLinkLocal                types.String                 `tfsdk:"sham_link_local"`
	DatabaseProtection           *ospfBlockDatabaseProtection `tfsdk:"database_protection"`
	GracefulRestart              *ospfBlockGracefulRestart    `tfsdk:"graceful_restart"`
	Overload                     *ospfBlockOverload           `tfsdk:"overload"`
	SpfOptions                   *ospfBlockSpfOptions         `tfsdk:"spf_options"`
}

type ospfBlockDatabaseProtection struct {
	MaximumLsa       types.Int64 `tfsdk:"maximum_lsa"`
	IgnoreCount      types.Int64 `tfsdk:"ignore_count"`
	IgnoreTime       types.Int64 `tfsdk:"ignore_time"`
	ResetTime        types.Int64 `tfsdk:"reset_time"`
	WarningOnly      types.Bool  `tfsdk:"warning_only"`
	WarningThreshold types.Int64 `tfsdk:"warning_threshold"`
}

type ospfBlockGracefulRestart struct {
	Disable             types.Bool   `tfsdk:"disable"`
	HelperDisable       types.Bool   `tfsdk:"helper_disable"`
	HelperDisableType   types.String `tfsdk:"helper_disable_type"`
	NoStrictLsaChecking types.Bool   `tfsdk:"no_strict_lsa_checking"`
	NotifyDuration      types.Int64  `tfsdk:"notify_duration"`
	RestartDuration     types.Int64  `tfsdk:"restart_duration"`
}

func (block *ospfBlockGracefulRestart) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type ospfBlockOverload struct {
	AllowRouteLeaking types.Bool  `tfsdk:"allow_route_leaking"`
	ASExternal        types.Bool  `tfsdk:"as_external"`
	StubNetwork       types.Bool  `tfsdk:"stub_network"`
	Timeout           types.Int64 `tfsdk:"timeout"`
}

type ospfBlockSpfOptions struct {
	Delay                types.Int64 `tfsdk:"delay"`
	Holddown             types.Int64 `tfsdk:"holddown"`
	NoIgnoreOurExternals types.Bool  `tfsdk:"no_ignore_our_externals"`
	RapidRuns            types.Int64 `tfsdk:"rapid_runs"`
}

func (block *ospfBlockSpfOptions) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (rsc *ospf) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config ospfConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.DomainID.IsNull() && !config.DomainID.IsUnknown() {
		if config.RoutingInstance.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("domain_id"),
				tfdiag.MissingConfigErrSummary,
				"routing_instance must be specified with domain_id",
			)
		} else if !config.RoutingInstance.IsUnknown() &&
			config.RoutingInstance.ValueString() == junos.DefaultW {
			resp.Diagnostics.AddAttributeError(
				path.Root("domain_id"),
				tfdiag.ConflictConfigErrSummary,
				fmt.Sprintf("routing_instance cannot be %q with domain_id", junos.DefaultW),
			)
		}
	}
	if !config.ShamLinkLocal.IsNull() && !config.ShamLinkLocal.IsUnknown() &&
		config.ShamLink.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("sham_link_local"),
			tfdiag.MissingConfigErrSummary,
			"sham_link must be specified with sham_link_local",
		)
	}

	if config.GracefulRestart != nil {
		if config.GracefulRestart.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("graceful_restart"),
				tfdiag.MissingConfigErrSummary,
				"graceful_restart block is empty",
			)
		}
		if !config.GracefulRestart.HelperDisableType.IsNull() && !config.GracefulRestart.HelperDisableType.IsUnknown() &&
			config.GracefulRestart.HelperDisable.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("graceful_restart").AtName("helper_disable_type"),
				tfdiag.MissingConfigErrSummary,
				"helper_disable must be specified with helper_disable_type"+
					" in graceful_restart block",
			)
		}
		if !config.GracefulRestart.NoStrictLsaChecking.IsNull() && !config.GracefulRestart.NoStrictLsaChecking.IsUnknown() &&
			!config.GracefulRestart.HelperDisable.IsNull() && !config.GracefulRestart.HelperDisable.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("graceful_restart").AtName("no_strict_lsa_checking"),
				tfdiag.ConflictConfigErrSummary,
				"no_strict_lsa_checking and helper_disable cannot be configured together"+
					" in graceful_restart block",
			)
		}
	}
	if config.SpfOptions != nil {
		if config.SpfOptions.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("spf_options"),
				tfdiag.MissingConfigErrSummary,
				"spf_options block is empty",
			)
		}
	}
}

func (rsc *ospf) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan ospfData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Version.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("version"),
			"Empty Version",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "version"),
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

			return true
		},
		nil,
		&plan,
		resp,
	)
}

func (rsc *ospf) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data ospfData
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

	err = data.read(ctx, state.Version.ValueString(), state.RoutingInstance.ValueString(), junSess)
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

func (rsc *ospf) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state ospfData
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

func (rsc *ospf) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state ospfData
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

func (rsc *ospf) ImportState(
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
	if idList[0] != "v2" && idList[0] != "v3" {
		resp.Diagnostics.AddError(
			"Bad ID Format",
			fmt.Sprintf("%q is not a valid version", idList[0]),
		)

		return
	}
	if idList[1] != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, idList[1], junSess)
		if err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !instanceExists {
			resp.Diagnostics.AddError(
				tfdiag.NotFoundErrSummary,
				fmt.Sprintf("routing instance %q doesn't exist", idList[1]),
			)

			return
		}
	}

	var data ospfData
	if err := data.read(ctx, idList[0], idList[1], junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if data.nullID() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be <version>_-_<routing_instance>)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rscData *ospfData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Version.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Version.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *ospfData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *ospfData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	routingInstance := rscData.RoutingInstance.ValueString()
	if routingInstance != "" && routingInstance != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	ospfVersion := junos.OspfV2
	if rscData.Version.ValueString() == "v3" {
		ospfVersion = junos.OspfV3
	}
	setPrefix += "protocols " + ospfVersion + " "

	if rscData.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if v := rscData.DomainID.ValueString(); v != "" {
		if routingInstance == "" || routingInstance == junos.DefaultW {
			return path.Root("domain_id"),
				fmt.Errorf("domain_id cannot be configured when routing_instance = %q", junos.DefaultW)
		}
		configSet = append(configSet, setPrefix+"domain-id \""+v+"\"")
	}
	for _, v := range rscData.Export {
		configSet = append(configSet, setPrefix+"export \""+v.ValueString()+"\"")
	}
	if !rscData.ExternalPreference.IsNull() {
		configSet = append(configSet, setPrefix+"external-preference "+
			utils.ConvI64toa(rscData.ExternalPreference.ValueInt64()))
	}
	if rscData.ForwardingAddressToBroadcast.ValueBool() {
		configSet = append(configSet, setPrefix+"forwarding-address-to-broadcast")
	}
	for _, v := range rscData.Import {
		configSet = append(configSet, setPrefix+"import \""+v.ValueString()+"\"")
	}
	if !rscData.LabeledPreference.IsNull() {
		configSet = append(configSet, setPrefix+"labeled-preference "+
			utils.ConvI64toa(rscData.LabeledPreference.ValueInt64()))
	}
	if !rscData.LsaRefreshInterval.IsNull() {
		configSet = append(configSet, setPrefix+"lsa-refresh-interval "+
			utils.ConvI64toa(rscData.LsaRefreshInterval.ValueInt64()))
	}
	if rscData.NoNssaAbr.ValueBool() {
		configSet = append(configSet, setPrefix+"no-nssa-abr")
	}
	if rscData.NoRfc1583.ValueBool() {
		configSet = append(configSet, setPrefix+"no-rfc-1583")
	}
	if !rscData.Preference.IsNull() {
		configSet = append(configSet, setPrefix+"preference "+
			utils.ConvI64toa(rscData.Preference.ValueInt64()))
	}
	if !rscData.PrefixExportLimit.IsNull() {
		configSet = append(configSet, setPrefix+"prefix-export-limit "+
			utils.ConvI64toa(rscData.PrefixExportLimit.ValueInt64()))
	}
	if v := rscData.ReferenceBandwidth.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"reference-bandwidth "+v)
	}
	if v := rscData.RibGroup.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"rib-group \""+v+"\"")
	}
	if rscData.ShamLink.ValueBool() {
		configSet = append(configSet, setPrefix+"sham-link")
		if v := rscData.ShamLinkLocal.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"sham-link local "+v)
		}
	} else if rscData.ShamLinkLocal.ValueString() != "" {
		return path.Root("sham_link_local"),
			errors.New("sham_link must be specified with sham_link_local")
	}

	if rscData.DatabaseProtection != nil {
		configSet = append(configSet, rscData.DatabaseProtection.configSet(setPrefix)...)
	}
	if rscData.GracefulRestart != nil {
		if rscData.GracefulRestart.isEmpty() {
			return path.Root("graceful_restart").AtName("*"),
				errors.New("graceful_restart block is empty")
		}

		blockSet, pathErr, err := rscData.GracefulRestart.configSet(setPrefix, path.Root("graceful_restart"))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.Overload != nil {
		configSet = append(configSet, rscData.Overload.configSet(setPrefix)...)
	}
	if rscData.SpfOptions != nil {
		if rscData.SpfOptions.isEmpty() {
			return path.Root("spf_options").AtName("*"),
				errors.New("spf_options block is empty")
		}

		configSet = append(configSet, rscData.SpfOptions.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *ospfBlockDatabaseProtection) configSet(setPrefix string) []string {
	setPrefix += "database-protection "

	configSet := []string{
		setPrefix + "maximum-lsa " + utils.ConvI64toa(block.MaximumLsa.ValueInt64()),
	}

	if !block.IgnoreCount.IsNull() {
		configSet = append(configSet, setPrefix+"ignore-count "+
			utils.ConvI64toa(block.IgnoreCount.ValueInt64()))
	}
	if !block.IgnoreTime.IsNull() {
		configSet = append(configSet, setPrefix+"ignore-time "+
			utils.ConvI64toa(block.IgnoreTime.ValueInt64()))
	}
	if !block.ResetTime.IsNull() {
		configSet = append(configSet, setPrefix+"reset-time "+
			utils.ConvI64toa(block.ResetTime.ValueInt64()))
	}
	if block.WarningOnly.ValueBool() {
		configSet = append(configSet, setPrefix+"warning-only")
	}
	if !block.WarningThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"warning-threshold "+
			utils.ConvI64toa(block.WarningThreshold.ValueInt64()))
	}

	return configSet
}

func (block *ospfBlockGracefulRestart) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "graceful-restart "

	if block.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if block.HelperDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"helper-disable")
		if v := block.HelperDisableType.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"helper-disable "+v)
		}
	} else if block.HelperDisableType.ValueString() != "" {
		return configSet,
			pathRoot.AtName("helper_disable_type"),
			errors.New("helper_disable must be specified with helper_disable_type" +
				" in graceful_restart block")
	}
	if block.NoStrictLsaChecking.ValueBool() {
		configSet = append(configSet, setPrefix+"no-strict-lsa-checking")
	}
	if !block.NotifyDuration.IsNull() {
		configSet = append(configSet, setPrefix+"notify-duration "+
			utils.ConvI64toa(block.NotifyDuration.ValueInt64()))
	}
	if !block.RestartDuration.IsNull() {
		configSet = append(configSet, setPrefix+"restart-duration "+
			utils.ConvI64toa(block.RestartDuration.ValueInt64()))
	}

	return configSet, path.Empty(), nil
}

func (block *ospfBlockOverload) configSet(setPrefix string) []string {
	setPrefix += "overload "

	configSet := []string{
		setPrefix,
	}

	if block.AllowRouteLeaking.ValueBool() {
		configSet = append(configSet, setPrefix+"allow-route-leaking")
	}
	if block.ASExternal.ValueBool() {
		configSet = append(configSet, setPrefix+"as-external")
	}
	if block.StubNetwork.ValueBool() {
		configSet = append(configSet, setPrefix+"stub-network")
	}
	if !block.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(block.Timeout.ValueInt64()))
	}

	return configSet
}

func (block *ospfBlockSpfOptions) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 1)
	setPrefix += "spf-options "

	if !block.Delay.IsNull() {
		configSet = append(configSet, setPrefix+"delay "+
			utils.ConvI64toa(block.Delay.ValueInt64()))
	}
	if !block.Holddown.IsNull() {
		configSet = append(configSet, setPrefix+"holddown "+
			utils.ConvI64toa(block.Holddown.ValueInt64()))
	}
	if block.NoIgnoreOurExternals.ValueBool() {
		configSet = append(configSet, setPrefix+"no-ignore-our-externals")
	}
	if !block.RapidRuns.IsNull() {
		configSet = append(configSet, setPrefix+"rapid-runs "+
			utils.ConvI64toa(block.RapidRuns.ValueInt64()))
	}

	return configSet
}

func (rscData *ospfData) read(
	_ context.Context, version, routingInstance string, junSess *junos.Session,
) error {
	ospfVersion := junos.OspfV2
	if version == "v3" {
		ospfVersion = junos.OspfV3
	}
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols " + ospfVersion + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if version == "v3" {
		rscData.Version = types.StringValue(version)
	} else {
		rscData.Version = types.StringValue("v2")
	}
	if routingInstance == "" {
		rscData.RoutingInstance = types.StringValue(junos.DefaultW)
	} else {
		rscData.RoutingInstance = types.StringValue(routingInstance)
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
			case itemTrim == "disable":
				rscData.Disable = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "domain-id "):
				rscData.DomainID = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "export "):
				rscData.Export = append(rscData.Export, types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "external-preference "):
				rscData.ExternalPreference, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "forwarding-address-to-broadcast":
				rscData.ForwardingAddressToBroadcast = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "import "):
				rscData.Import = append(rscData.Import, types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "labeled-preference "):
				rscData.LabeledPreference, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "lsa-refresh-interval "):
				rscData.LsaRefreshInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "no-nssa-abr":
				rscData.NoNssaAbr = types.BoolValue(true)
			case itemTrim == "no-rfc-1583":
				rscData.NoRfc1583 = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "preference "):
				rscData.Preference, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "prefix-export-limit "):
				rscData.PrefixExportLimit, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "reference-bandwidth "):
				rscData.ReferenceBandwidth = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "rib-group "):
				rscData.RibGroup = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "sham-link"):
				rscData.ShamLink = types.BoolValue(true)
				if balt.CutPrefixInString(&itemTrim, " local ") {
					rscData.ShamLinkLocal = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "database-protection "):
				if rscData.DatabaseProtection == nil {
					rscData.DatabaseProtection = &ospfBlockDatabaseProtection{}
				}

				if err := rscData.DatabaseProtection.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "graceful-restart "):
				if rscData.GracefulRestart == nil {
					rscData.GracefulRestart = &ospfBlockGracefulRestart{}
				}

				if err := rscData.GracefulRestart.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "overload"):
				if rscData.Overload == nil {
					rscData.Overload = &ospfBlockOverload{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.Overload.read(itemTrim); err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "spf-options "):
				if rscData.SpfOptions == nil {
					rscData.SpfOptions = &ospfBlockSpfOptions{}
				}

				if err := rscData.SpfOptions.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *ospfBlockDatabaseProtection) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "ignore-count "):
		block.IgnoreCount, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "ignore-time "):
		block.IgnoreTime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "maximum-lsa "):
		block.MaximumLsa, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "reset-time "):
		block.ResetTime, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "warning-only":
		block.WarningOnly = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "warning-threshold "):
		block.WarningThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (block *ospfBlockGracefulRestart) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "disable":
		block.Disable = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "helper-disable"):
		block.HelperDisable = types.BoolValue(true)
		if balt.CutPrefixInString(&itemTrim, " ") {
			block.HelperDisableType = types.StringValue(itemTrim)
		}
	case itemTrim == "no-strict-lsa-checking":
		block.NoStrictLsaChecking = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "notify-duration "):
		block.NotifyDuration, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "restart-duration "):
		block.RestartDuration, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (block *ospfBlockOverload) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "allow-route-leaking":
		block.AllowRouteLeaking = types.BoolValue(true)
	case itemTrim == "as-external":
		block.ASExternal = types.BoolValue(true)
	case itemTrim == "stub-network":
		block.StubNetwork = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "timeout "):
		block.Timeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (block *ospfBlockSpfOptions) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "delay "):
		block.Delay, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "holddown "):
		block.Holddown, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "no-ignore-our-externals":
		block.NoIgnoreOurExternals = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "rapid-runs "):
		block.RapidRuns, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rscData *ospfData) del(
	_ context.Context, junSess *junos.Session,
) error {
	ospfVersion := junos.OspfV2
	if rscData.Version.ValueString() == "v3" {
		ospfVersion = junos.OspfV3
	}
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "protocols " + ospfVersion + " "

	configSet := []string{
		delPrefix + "database-protection",
		delPrefix + "disable",
		delPrefix + "domain-id",
		delPrefix + "export",
		delPrefix + "external-preference",
		delPrefix + "forwarding-address-to-broadcast",
		delPrefix + "graceful-restart",
		delPrefix + "import",
		delPrefix + "labeled-preference",
		delPrefix + "lsa-refresh-interval",
		delPrefix + "no-nssa-abr",
		delPrefix + "no-rfc-1583",
		delPrefix + "overload",
		delPrefix + "preference",
		delPrefix + "prefix-export-limit",
		delPrefix + "reference-bandwidth",
		delPrefix + "rib-group",
		delPrefix + "sham-link",
		delPrefix + "spf-options",
	}

	return junSess.ConfigSet(configSet)
}
