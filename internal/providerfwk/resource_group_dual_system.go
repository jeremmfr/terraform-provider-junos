package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &groupDualSystem{}
	_ resource.ResourceWithConfigure      = &groupDualSystem{}
	_ resource.ResourceWithValidateConfig = &groupDualSystem{}
	_ resource.ResourceWithImportState    = &groupDualSystem{}
	_ resource.ResourceWithUpgradeState   = &groupDualSystem{}
)

type groupDualSystem struct {
	client *junos.Client
}

func newGroupDualSystemResource() resource.Resource {
	return &groupDualSystem{}
}

func (rsc *groupDualSystem) typeName() string {
	return providerName + "_group_dual_system"
}

func (rsc *groupDualSystem) junosName() string {
	return "groups node0|node1|re0|re1"
}

func (rsc *groupDualSystem) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *groupDualSystem) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *groupDualSystem) Configure(
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

func (rsc *groupDualSystem) Schema(
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
				Description: "Name of group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("node0", "node1", "re0", "re1"),
				},
			},
			"apply_groups": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Apply the group.",
			},
		},
		Blocks: map[string]schema.Block{
			"interface_fxp0": schema.SingleNestedBlock{
				Description: "Configure `fxp0` interface.",
				Attributes: map[string]schema.Attribute{
					"description": schema.StringAttribute{
						Optional:    true,
						Description: "Description for interface.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 900),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"family_inet_address": schema.ListNestedBlock{
						Description: "For each IPv4 address to declare.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"cidr_ip": schema.StringAttribute{
									Required:    true,
									Description: "IPv4 address in CIDR format.",
									Validators: []validator.String{
										tfvalidator.StringCIDR().IPv4Only(),
									},
								},
								"master_only": schema.BoolAttribute{
									Optional:    true,
									Description: "Master management IP address.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"preferred": schema.BoolAttribute{
									Optional:    true,
									Description: "Preferred address on interface.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"primary": schema.BoolAttribute{
									Optional:    true,
									Description: "Candidate for primary address in system.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
							},
						},
					},
					"family_inet6_address": schema.ListNestedBlock{
						Description: "For each IPv6 address to declare.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"cidr_ip": schema.StringAttribute{
									Required:    true,
									Description: "IPv6 address in CIDR format.",
									Validators: []validator.String{
										tfvalidator.StringCIDR().IPv6Only(),
									},
								},
								"master_only": schema.BoolAttribute{
									Optional:    true,
									Description: "Master management IP address.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"preferred": schema.BoolAttribute{
									Optional:    true,
									Description: "Preferred address on interface.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"primary": schema.BoolAttribute{
									Optional:    true,
									Description: "Candidate for primary address in system.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
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
			"routing_options": schema.SingleNestedBlock{
				Description: "Configure `routing-options` block.",
				Blocks: map[string]schema.Block{
					"static_route": schema.ListNestedBlock{
						Description: "For each destination to declare.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"destination": schema.StringAttribute{
									Required:    true,
									Description: "The destination for static route.",
									Validators: []validator.String{
										tfvalidator.StringCIDRNetwork(),
									},
								},
								"next_hop": schema.ListAttribute{
									ElementType: types.StringType,
									Required:    true,
									Description: "List of next-hop.",
									Validators: []validator.List{
										listvalidator.SizeAtLeast(1),
										listvalidator.NoNullValues(),
										listvalidator.ValueStringsAre(
											stringvalidator.LengthAtLeast(1),
											stringvalidator.Any(
												tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
												tfvalidator.StringIPAddress(),
											),
										),
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
			"security": schema.SingleNestedBlock{
				Description: "Configure `security` block.",
				Attributes: map[string]schema.Attribute{
					"log_source_address": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Source IP address used when exporting security logs.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"system": schema.SingleNestedBlock{
				Description: "Configure `system` block.",
				Attributes: map[string]schema.Attribute{
					"host_name": schema.StringAttribute{
						Optional:    true,
						Description: "Hostname.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 255),
							tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
						},
					},
					"backup_router_address": schema.StringAttribute{
						Optional:    true,
						Description: "IPv4 address backup router.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress().IPv4Only(),
						},
					},
					"backup_router_destination": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Destinations network reachable through the IPv4 router.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.NoNullValues(),
							setvalidator.ValueStringsAre(
								tfvalidator.StringCIDR().IPv4Only(),
							),
						},
					},
					"inet6_backup_router_address": schema.StringAttribute{
						Optional:    true,
						Description: "IPv6 address backup router.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress().IPv6Only(),
						},
					},
					"inet6_backup_router_destination": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Destinations network reachable through the IPv6 router.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.NoNullValues(),
							setvalidator.ValueStringsAre(
								tfvalidator.StringCIDR().IPv6Only(),
							),
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

type groupDualSystemData struct {
	ID             types.String                        `tfsdk:"id"              tfdata:"skip_isempty"`
	Name           types.String                        `tfsdk:"name"            tfdata:"skip_isempty"`
	ApplyGroups    types.Bool                          `tfsdk:"apply_groups"    tfdata:"skip_isempty"`
	InterfaceFXP0  *groupDualSystemBlockInterfaceFXP0  `tfsdk:"interface_fxp0"`
	RoutingOptions *groupDualSystemBlockRoutingOptions `tfsdk:"routing_options"`
	Security       *groupDualSystemBlockSecurity       `tfsdk:"security"`
	System         *groupDualSystemBlockSystem         `tfsdk:"system"`
}

func (rscData *groupDualSystemData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

type groupDualSystemConfig struct {
	ID             types.String                              `tfsdk:"id"              tfdata:"skip_isempty"`
	Name           types.String                              `tfsdk:"name"            tfdata:"skip_isempty"`
	ApplyGroups    types.Bool                                `tfsdk:"apply_groups"    tfdata:"skip_isempty"`
	InterfaceFXP0  *groupDualSystemBlockInterfaceFXP0Config  `tfsdk:"interface_fxp0"`
	RoutingOptions *groupDualSystemBlockRoutingOptionsConfig `tfsdk:"routing_options"`
	Security       *groupDualSystemBlockSecurity             `tfsdk:"security"`
	System         *groupDualSystemSystemConfig              `tfsdk:"system"`
}

func (config *groupDualSystemConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(config)
}

type groupDualSystemBlockInterfaceFXP0 struct {
	Description        types.String                                          `tfsdk:"description"`
	FamilyInetAddress  []groupDualSystemBlockInterfaceFXP0BlockFamilyAddress `tfsdk:"family_inet_address"`
	FamilyInet6Address []groupDualSystemBlockInterfaceFXP0BlockFamilyAddress `tfsdk:"family_inet6_address"`
}

func (block *groupDualSystemBlockInterfaceFXP0) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type groupDualSystemBlockInterfaceFXP0Config struct {
	Description        types.String `tfsdk:"description"`
	FamilyInetAddress  types.List   `tfsdk:"family_inet_address"`
	FamilyInet6Address types.List   `tfsdk:"family_inet6_address"`
}

func (block *groupDualSystemBlockInterfaceFXP0Config) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type groupDualSystemBlockInterfaceFXP0BlockFamilyAddress struct {
	CidrIP     types.String `tfsdk:"cidr_ip"     tfdata:"identifier"`
	MasterOnly types.Bool   `tfsdk:"master_only"`
	Preferred  types.Bool   `tfsdk:"preferred"`
	Primary    types.Bool   `tfsdk:"primary"`
}

type groupDualSystemBlockRoutingOptions struct {
	StaticRoute []groupDualSystemBlockRoutingOptionsBlockStaticRoute `tfsdk:"static_route"`
}

type groupDualSystemBlockRoutingOptionsConfig struct {
	StaticRoute types.List `tfsdk:"static_route"`
}

type groupDualSystemBlockRoutingOptionsBlockStaticRoute struct {
	Destination types.String   `tfsdk:"destination" tfdata:"identifier"`
	NextHop     []types.String `tfsdk:"next_hop"`
}

type groupDualSystemBlockSecurity struct {
	LogSourceAddress types.String `tfsdk:"log_source_address"`
}

type groupDualSystemBlockSystem struct {
	HostName                     types.String   `tfsdk:"host_name"`
	BackupRouterAddress          types.String   `tfsdk:"backup_router_address"`
	BackupRouterDestination      []types.String `tfsdk:"backup_router_destination"`
	Inet6BackupRouterAddress     types.String   `tfsdk:"inet6_backup_router_address"`
	Inet6BackupRouterDestination []types.String `tfsdk:"inet6_backup_router_destination"`
}

func (block *groupDualSystemBlockSystem) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type groupDualSystemSystemConfig struct {
	HostName                     types.String `tfsdk:"host_name"`
	BackupRouterAddress          types.String `tfsdk:"backup_router_address"`
	BackupRouterDestination      types.Set    `tfsdk:"backup_router_destination"`
	Inet6BackupRouterAddress     types.String `tfsdk:"inet6_backup_router_address"`
	Inet6BackupRouterDestination types.Set    `tfsdk:"inet6_backup_router_destination"`
}

func (rscData *groupDualSystemSystemConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

func (rsc *groupDualSystem) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config groupDualSystemConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one of arguments need to be set (in addition to `name` and `apply_groups`)",
		)
	}

	if config.InterfaceFXP0 != nil {
		if config.InterfaceFXP0.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("interface_fxp0"),
				tfdiag.MissingConfigErrSummary,
				"interface_fxp0 block is empty",
			)
		}
		if !config.InterfaceFXP0.FamilyInetAddress.IsNull() &&
			!config.InterfaceFXP0.FamilyInetAddress.IsUnknown() {
			var familyInetAddress []groupDualSystemBlockInterfaceFXP0BlockFamilyAddress
			asDiags := config.InterfaceFXP0.FamilyInetAddress.ElementsAs(ctx, &familyInetAddress, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			familyInetAddressCidrIP := make(map[string]struct{})
			for i, block := range familyInetAddress {
				if block.CidrIP.IsUnknown() {
					continue
				}
				cidrIP := block.CidrIP.ValueString()
				if _, ok := familyInetAddressCidrIP[cidrIP]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface_fxp0").AtName("family_inet_address").AtListIndex(i).AtName("cidr_ip"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple family_inet_address blocks with the same cidr_ip %q"+
							" in interface_fxp0 block", cidrIP),
					)
				}
				familyInetAddressCidrIP[cidrIP] = struct{}{}
			}
		}
		if !config.InterfaceFXP0.FamilyInet6Address.IsNull() &&
			!config.InterfaceFXP0.FamilyInet6Address.IsUnknown() {
			var familyInet6Address []groupDualSystemBlockInterfaceFXP0BlockFamilyAddress
			asDiags := config.InterfaceFXP0.FamilyInet6Address.ElementsAs(ctx, &familyInet6Address, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			familyInet6AddressCidrIP := make(map[string]struct{})
			for i, block := range familyInet6Address {
				if block.CidrIP.IsUnknown() {
					continue
				}
				cidrIP := block.CidrIP.ValueString()
				if _, ok := familyInet6AddressCidrIP[cidrIP]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface_fxp0").AtName("family_inet6_address").AtListIndex(i).AtName("cidr_ip"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple family_inet6_address blocks with the same cidr_ip %q"+
							" in interface_fxp0 block", cidrIP),
					)
				}
				familyInet6AddressCidrIP[cidrIP] = struct{}{}
			}
		}
	}
	if config.RoutingOptions != nil {
		if config.RoutingOptions.StaticRoute.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("routing_options").AtName("static_route"),
				tfdiag.MissingConfigErrSummary,
				"static_route must be specified in routing_options block",
			)
		}
	}
	if config.Security != nil {
		if config.Security.LogSourceAddress.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("security").AtName("log_source_address"),
				tfdiag.MissingConfigErrSummary,
				"log_source_address must be specified in security block",
			)
		}
	}
	if config.System != nil {
		if config.System.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("system"),
				tfdiag.MissingConfigErrSummary,
				"system block is empty",
			)
		}
	}
}

func (rsc *groupDualSystem) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan groupDualSystemData
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
			groupExists, err := checkGroupDualSystemExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if groupExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			groupExists, err := checkGroupDualSystemExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !groupExists {
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

func (rsc *groupDualSystem) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data groupDualSystemData
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

func (rsc *groupDualSystem) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state groupDualSystemData
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

func (rsc *groupDualSystem) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state groupDualSystemData
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

func (rsc *groupDualSystem) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data groupDualSystemData

	if !slices.Contains([]string{"node0", "node1", "re0", "re1"}, req.ID) {
		resp.Diagnostics.AddError(
			"Bad ID Format",
			fmt.Sprintf("invalid group id %q (id must be <name>)", req.ID),
		)

		return
	}

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

func checkGroupDualSystemExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"groups " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *groupDualSystemData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *groupDualSystemData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *groupDualSystemData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)

	if rscData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name` and `apply_groups`)")
	}

	if rscData.ApplyGroups.ValueBool() {
		if strings.HasPrefix(rscData.Name.ValueString(), "node") {
			configSet = append(configSet, "set apply-groups \"${node}\"")
		} else {
			configSet = append(configSet, "set apply-groups "+rscData.Name.ValueString())
		}
	}

	setPrefix := "set groups " + rscData.Name.ValueString() + " "
	if rscData.InterfaceFXP0 != nil {
		if rscData.InterfaceFXP0.isEmpty() {
			return path.Root("interface_fxp0"),
				errors.New("interface_fxp0 block is empty")
		}

		blockSet, pathErr, err := rscData.InterfaceFXP0.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.RoutingOptions != nil {
		staticRouteDestination := make(map[string]struct{})
		for i, block := range rscData.RoutingOptions.StaticRoute {
			destination := block.Destination.ValueString()
			if _, ok := staticRouteDestination[destination]; ok {
				return path.Root("routing_options").AtName("static_route").AtListIndex(i).AtName("destination"),
					fmt.Errorf("multiple static_route blocks with the same destination %q"+
						"in routing_options block", destination)
			}
			staticRouteDestination[destination] = struct{}{}

			for _, v := range block.NextHop {
				configSet = append(configSet, setPrefix+junos.RoutingOptionsWS+"static route "+
					destination+" next-hop "+v.ValueString())
			}
		}
	}
	if rscData.Security != nil {
		if v := rscData.Security.LogSourceAddress.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"security log source-address "+v)
		}
	}
	if rscData.System != nil {
		if rscData.System.isEmpty() {
			return path.Root("system"),
				errors.New("system block is empty")
		}
		configSet = append(configSet, rscData.System.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *groupDualSystemBlockInterfaceFXP0) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, 100)
	setPrefix += "interfaces fxp0 "

	if v := block.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}

	familyInetAddressCidrIP := make(map[string]struct{})
	for i, block := range block.FamilyInetAddress {
		cidrIP := block.CidrIP.ValueString()
		if _, ok := familyInetAddressCidrIP[cidrIP]; ok {
			return configSet,
				path.Root("interface_fxp0").AtName("family_inet_address").AtListIndex(i).AtName("cidr_ip"),
				fmt.Errorf("multiple family_inet_address blocks with the same cidr_ip %q"+
					"in interface_fxp0 block", cidrIP)
		}
		familyInetAddressCidrIP[cidrIP] = struct{}{}

		configSet = append(configSet, block.configSet(setPrefix+"unit 0 family inet ")...)
	}
	familyInet6AddressCidrIP := make(map[string]struct{})
	for i, block := range block.FamilyInet6Address {
		cidrIP := block.CidrIP.ValueString()
		if _, ok := familyInet6AddressCidrIP[cidrIP]; ok {
			return configSet,
				path.Root("interface_fxp0").AtName("family_inet6_address").AtListIndex(i).AtName("cidr_ip"),
				fmt.Errorf("multiple family_inet6_address blocks with the same cidr_ip %q"+
					"in interface_fxp0 block", cidrIP)
		}
		familyInet6AddressCidrIP[cidrIP] = struct{}{}

		configSet = append(configSet, block.configSet(setPrefix+"unit 0 family inet6 ")...)
	}

	return configSet, path.Empty(), nil
}

func (block *groupDualSystemBlockInterfaceFXP0BlockFamilyAddress) configSet(setPrefix string) []string {
	setPrefix += "address " + block.CidrIP.ValueString() + " "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if block.MasterOnly.ValueBool() {
		configSet = append(configSet, setPrefix+"master-only")
	}
	if block.Preferred.ValueBool() {
		configSet = append(configSet, setPrefix+"preferred")
	}
	if block.Primary.ValueBool() {
		configSet = append(configSet, setPrefix+"primary")
	}

	return configSet
}

func (block *groupDualSystemBlockSystem) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)
	setPrefix += "system "

	if v := block.HostName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"host-name \""+v+"\"")
	}
	if v := block.BackupRouterAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"backup-router "+v)
	}
	for _, v := range block.BackupRouterDestination {
		configSet = append(configSet, setPrefix+"backup-router destination "+v.ValueString())
	}
	if v := block.Inet6BackupRouterAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"inet6-backup-router "+v)
	}
	for _, v := range block.Inet6BackupRouterDestination {
		configSet = append(configSet, setPrefix+"inet6-backup-router destination "+v.ValueString())
	}

	return configSet
}

func (rscData *groupDualSystemData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"groups " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "interfaces fxp0 "):
				if rscData.InterfaceFXP0 == nil {
					rscData.InterfaceFXP0 = &groupDualSystemBlockInterfaceFXP0{}
				}

				rscData.InterfaceFXP0.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, junos.RoutingOptionsWS+"static route "):
				if rscData.RoutingOptions == nil {
					rscData.RoutingOptions = &groupDualSystemBlockRoutingOptions{}
				}

				destination := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.RoutingOptions.StaticRoute = tfdata.AppendPotentialNewBlock(
					rscData.RoutingOptions.StaticRoute, types.StringValue(destination),
				)
				staticRoute := &rscData.RoutingOptions.StaticRoute[len(rscData.RoutingOptions.StaticRoute)-1]

				if balt.CutPrefixInString(&itemTrim, destination+" next-hop ") {
					staticRoute.NextHop = append(staticRoute.NextHop, types.StringValue(itemTrim))
				}
			case balt.CutPrefixInString(&itemTrim, "security "):
				if rscData.Security == nil {
					rscData.Security = &groupDualSystemBlockSecurity{}
				}

				if balt.CutPrefixInString(&itemTrim, "log source-address ") {
					rscData.Security.LogSourceAddress = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "system "):
				if rscData.System == nil {
					rscData.System = &groupDualSystemBlockSystem{}
				}

				rscData.System.read(itemTrim)
			}
		}
	}

	rscData.ApplyGroups = types.BoolValue(false)
	showConfigApplyGroups, err := junSess.Command(junos.CmdShowConfig +
		"apply-groups" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfigApplyGroups != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.fillID()
		for _, item := range strings.Split(showConfigApplyGroups, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch itemTrim {
			case name, name + " ":
				rscData.ApplyGroups = types.BoolValue(true)
			case "\"${node}\"", "\"${node}\" ":
				if strings.HasPrefix(name, "node") {
					rscData.ApplyGroups = types.BoolValue(true)
				}
			}
		}
	}

	return nil
}

func (block *groupDualSystemBlockInterfaceFXP0) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "description "):
		block.Description = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "unit 0 family inet address "):
		cidrIP := tfdata.FirstElementOfJunosLine(itemTrim)
		block.FamilyInetAddress = tfdata.AppendPotentialNewBlock(block.FamilyInetAddress, types.StringValue(cidrIP))
		familyInetAddress := &block.FamilyInetAddress[len(block.FamilyInetAddress)-1]

		if balt.CutPrefixInString(&itemTrim, cidrIP+" ") {
			familyInetAddress.read(itemTrim)
		}
	case balt.CutPrefixInString(&itemTrim, "unit 0 family inet6 address "):
		cidrIP := tfdata.FirstElementOfJunosLine(itemTrim)
		block.FamilyInet6Address = tfdata.AppendPotentialNewBlock(block.FamilyInet6Address, types.StringValue(cidrIP))
		familyInet6Address := &block.FamilyInet6Address[len(block.FamilyInet6Address)-1]

		if balt.CutPrefixInString(&itemTrim, cidrIP+" ") {
			familyInet6Address.read(itemTrim)
		}
	}
}

func (block *groupDualSystemBlockInterfaceFXP0BlockFamilyAddress) read(itemTrim string) {
	switch {
	case strings.HasSuffix(itemTrim, "master-only"):
		block.MasterOnly = types.BoolValue(true)
	case strings.HasSuffix(itemTrim, "preferred"):
		block.Preferred = types.BoolValue(true)
	case strings.HasSuffix(itemTrim, "primary"):
		block.Primary = types.BoolValue(true)
	}
}

func (block *groupDualSystemBlockSystem) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "host-name "):
		block.HostName = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "backup-router destination "):
		block.BackupRouterDestination = append(block.BackupRouterDestination, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "backup-router "):
		block.BackupRouterAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "inet6-backup-router destination "):
		block.Inet6BackupRouterDestination = append(block.Inet6BackupRouterDestination, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "inet6-backup-router "):
		block.Inet6BackupRouterAddress = types.StringValue(itemTrim)
	}
}

func (rscData *groupDualSystemData) del(
	_ context.Context, junSess *junos.Session,
) error {
	name := rscData.Name.ValueString()

	configSet := []string{
		junos.DeleteLS + "groups " + name,
	}
	if rscData.ApplyGroups.ValueBool() {
		if strings.HasPrefix(name, "node") {
			configSet = append(configSet, junos.DeleteLS+"apply-groups \"${node}\"")
		} else {
			configSet = append(configSet, junos.DeleteLS+"apply-groups "+name)
		}
	}

	return junSess.ConfigSet(configSet)
}
