package provider

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
	_ resource.Resource                   = &igmpSnoopingVlan{}
	_ resource.ResourceWithConfigure      = &igmpSnoopingVlan{}
	_ resource.ResourceWithValidateConfig = &igmpSnoopingVlan{}
	_ resource.ResourceWithImportState    = &igmpSnoopingVlan{}
)

type igmpSnoopingVlan struct {
	client *junos.Client
}

func newIgmpSnoopingVlanResource() resource.Resource {
	return &igmpSnoopingVlan{}
}

func (rsc *igmpSnoopingVlan) typeName() string {
	return providerName + "_igmp_snooping_vlan"
}

func (rsc *igmpSnoopingVlan) junosName() string {
	return "igmp-snooping vlan"
}

func (rsc *igmpSnoopingVlan) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *igmpSnoopingVlan) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *igmpSnoopingVlan) Configure(
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

func (rsc *igmpSnoopingVlan) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<name>" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "VLAN name or `all`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for igmp-snooping protocol if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"immediate_leave": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable immediate group leave on interfaces.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"l2_querier_source_address": schema.StringAttribute{
				Optional:    true,
				Description: "Enable L2 querier mode with source address.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress().IPv4Only(),
				},
			},
			"proxy": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable proxy mode.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"proxy_source_address": schema.StringAttribute{
				Optional:    true,
				Description: "Source IP address to use for proxy.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress().IPv4Only(),
				},
			},
			"query_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "When to send host query messages (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 1024),
				},
			},
			"query_last_member_interval": schema.StringAttribute{
				Optional:    true,
				Description: "When to send group query messages (seconds).",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d+(\.\d+)?$`),
						"must be a number with optional decimal"),
				},
			},
			"query_response_interval": schema.StringAttribute{
				Optional:    true,
				Description: "How long to wait for a host query response (seconds).",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d+(\.\d+)?$`),
						"must be a number with optional decimal"),
				},
			},
			"robust_count": schema.Int64Attribute{
				Optional:    true,
				Description: "Expected packet loss on a subnet.",
				Validators: []validator.Int64{
					int64validator.Between(2, 10),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"interface": schema.ListNestedBlock{
				Description: "For each interface name, configure interface options for IGMP.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Interface name.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
								tfvalidator.String1DotCount(),
							},
						},
						"group_limit": schema.Int64Attribute{
							Optional:    true,
							Description: "Maximum number of groups an interface can join.",
							Validators: []validator.Int64{
								int64validator.Between(0, 65535),
							},
						},
						"host_only_interface": schema.BoolAttribute{
							Optional:    true,
							Description: "Enable interface to be treated as host-side interface.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"immediate_leave": schema.BoolAttribute{
							Optional:    true,
							Description: "Enable immediate group leave on interface.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"multicast_router_interface": schema.BoolAttribute{
							Optional:    true,
							Description: "Enabling multicast-router-interface on the interface.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"static_group": schema.SetNestedBlock{
							Description: "For each static group address.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"address": schema.StringAttribute{
										Required:    true,
										Description: "IP multicast group address.",
										Validators: []validator.String{
											tfvalidator.StringIPAddress().IPv4Only(),
										},
									},
									"source": schema.StringAttribute{
										Optional:    true,
										Description: "IP multicast source address.",
										Validators: []validator.String{
											tfvalidator.StringIPAddress().IPv4Only(),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type igmpSnoopingVlanData struct {
	ID                      types.String                     `tfsdk:"id"`
	Name                    types.String                     `tfsdk:"name"`
	RoutingInstance         types.String                     `tfsdk:"routing_instance"`
	ImmediateLeave          types.Bool                       `tfsdk:"immediate_leave"`
	L2QuerierSourceAddress  types.String                     `tfsdk:"l2_querier_source_address"`
	Proxy                   types.Bool                       `tfsdk:"proxy"`
	ProxySourceAddress      types.String                     `tfsdk:"proxy_source_address"`
	QueryInterval           types.Int64                      `tfsdk:"query_interval"`
	QueryLastMemberInterval types.String                     `tfsdk:"query_last_member_interval"`
	QueryResponseInterval   types.String                     `tfsdk:"query_response_interval"`
	RobustCount             types.Int64                      `tfsdk:"robust_count"`
	Interface               []igmpSnoopingVlanBlockInterface `tfsdk:"interface"`
}

type igmpSnoopingVlanConfig struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	RoutingInstance         types.String `tfsdk:"routing_instance"`
	ImmediateLeave          types.Bool   `tfsdk:"immediate_leave"`
	L2QuerierSourceAddress  types.String `tfsdk:"l2_querier_source_address"`
	Proxy                   types.Bool   `tfsdk:"proxy"`
	ProxySourceAddress      types.String `tfsdk:"proxy_source_address"`
	QueryInterval           types.Int64  `tfsdk:"query_interval"`
	QueryLastMemberInterval types.String `tfsdk:"query_last_member_interval"`
	QueryResponseInterval   types.String `tfsdk:"query_response_interval"`
	RobustCount             types.Int64  `tfsdk:"robust_count"`
	Interface               types.List   `tfsdk:"interface"`
}

//nolint:lll
type igmpSnoopingVlanBlockInterface struct {
	Name                     types.String                                     `tfsdk:"name"                       tfdata:"identifier"`
	GroupLimit               types.Int64                                      `tfsdk:"group_limit"`
	HostOnlyInterface        types.Bool                                       `tfsdk:"host_only_interface"`
	ImmediateLeave           types.Bool                                       `tfsdk:"immediate_leave"`
	MulticastRouterInterface types.Bool                                       `tfsdk:"multicast_router_interface"`
	StaticGroup              []igmpSnoopingVlanBlockInterfaceBlockStaticGroup `tfsdk:"static_group"`
}

type igmpSnoopingVlanBlockInterfaceConfig struct {
	Name                     types.String `tfsdk:"name"`
	GroupLimit               types.Int64  `tfsdk:"group_limit"`
	HostOnlyInterface        types.Bool   `tfsdk:"host_only_interface"`
	ImmediateLeave           types.Bool   `tfsdk:"immediate_leave"`
	MulticastRouterInterface types.Bool   `tfsdk:"multicast_router_interface"`
	StaticGroup              types.Set    `tfsdk:"static_group"`
}

type igmpSnoopingVlanBlockInterfaceBlockStaticGroup struct {
	Address types.String `tfsdk:"address" tfdata:"identifier"`
	Source  types.String `tfsdk:"source"`
}

func (rsc *igmpSnoopingVlan) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config igmpSnoopingVlanConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.ProxySourceAddress.IsNull() &&
		!config.ProxySourceAddress.IsUnknown() &&
		config.Proxy.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("proxy_source_address"),
			tfdiag.MissingConfigErrSummary,
			"proxy must be specified with proxy_source_address",
		)
	}

	if !config.Interface.IsNull() && !config.Interface.IsUnknown() {
		var configInterface []igmpSnoopingVlanBlockInterfaceConfig
		asDiags := config.Interface.ElementsAs(ctx, &configInterface, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		interfaceName := make(map[string]struct{})
		for i, block := range configInterface {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := interfaceName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple interface blocks with the same name %q", name),
					)
				}
				interfaceName[name] = struct{}{}
			}

			if !block.StaticGroup.IsNull() && !block.StaticGroup.IsUnknown() {
				var configStaticGroup []igmpSnoopingVlanBlockInterfaceBlockStaticGroup
				asDiags := block.StaticGroup.ElementsAs(ctx, &configStaticGroup, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				staticGroupAddress := make(map[string]struct{})
				for _, subBlock := range configStaticGroup {
					if subBlock.Address.IsUnknown() {
						continue
					}

					address := subBlock.Address.ValueString()
					if _, ok := staticGroupAddress[address]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("interface").AtListIndex(i).AtName("static_group").AtName("*"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple static_group blocks with the same address %q"+
								" in interface block %q", address, block.Name.ValueString()),
						)
					}
					staticGroupAddress[address] = struct{}{}
				}
			}
		}
	}
}

func (rsc *igmpSnoopingVlan) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan igmpSnoopingVlanData
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
			vlanExists, err := checkIgmpSnoopingVlanExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if vlanExists {
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

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			vlanExists, err := checkIgmpSnoopingVlanExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !vlanExists {
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

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *igmpSnoopingVlan) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data igmpSnoopingVlanData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *igmpSnoopingVlan) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state igmpSnoopingVlanData
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

func (rsc *igmpSnoopingVlan) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state igmpSnoopingVlanData
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

func (rsc *igmpSnoopingVlan) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data igmpSnoopingVlanData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <name>"+junos.IDSeparator+"<routing_instance>)",
	)
}

func checkIgmpSnoopingVlanExists(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols igmp-snooping vlan " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *igmpSnoopingVlanData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *igmpSnoopingVlanData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *igmpSnoopingVlanData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "protocols igmp-snooping vlan " + rscData.Name.ValueString() + " "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if rscData.ImmediateLeave.ValueBool() {
		configSet = append(configSet, setPrefix+"immediate-leave")
	}
	if v := rscData.L2QuerierSourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"l2-querier source-address "+v)
	}
	if rscData.Proxy.ValueBool() {
		configSet = append(configSet, setPrefix+"proxy")
		if v := rscData.ProxySourceAddress.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"proxy source-address "+v)
		}
	} else if rscData.ProxySourceAddress.ValueString() != "" {
		return path.Root("proxy_source_address"),
			errors.New("proxy must be specified with proxy_source_address")
	}
	if !rscData.QueryInterval.IsNull() {
		configSet = append(configSet, setPrefix+"query-interval "+
			utils.ConvI64toa(rscData.QueryInterval.ValueInt64()))
	}
	if v := rscData.QueryLastMemberInterval.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"query-last-member-interval "+v)
	}
	if v := rscData.QueryResponseInterval.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"query-response-interval "+v)
	}
	if !rscData.RobustCount.IsNull() {
		configSet = append(configSet, setPrefix+"robust-count "+
			utils.ConvI64toa(rscData.RobustCount.ValueInt64()))
	}

	interfaceName := make(map[string]struct{})
	for i, block := range rscData.Interface {
		name := block.Name.ValueString()
		if _, ok := interfaceName[name]; ok {
			return path.Root("interface").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple interface blocks with the same name %q", name)
		}
		interfaceName[name] = struct{}{}

		blockSet, pathErr, err := block.configSet(
			setPrefix,
			path.Root("interface").AtListIndex(i),
			fmt.Sprintf(" in interface block %q", block.Name.ValueString()),
		)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *igmpSnoopingVlanBlockInterface) configSet(
	setPrefix string, pathRoot path.Path, blockErrorSuffix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "interface " + block.Name.ValueString() + " "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if !block.GroupLimit.IsNull() {
		configSet = append(configSet, setPrefix+"group-limit "+
			utils.ConvI64toa(block.GroupLimit.ValueInt64()))
	}
	if block.HostOnlyInterface.ValueBool() {
		configSet = append(configSet, setPrefix+"host-only-interface")
	}
	if block.ImmediateLeave.ValueBool() {
		configSet = append(configSet, setPrefix+"immediate-leave")
	}
	if block.MulticastRouterInterface.ValueBool() {
		configSet = append(configSet, setPrefix+"multicast-router-interface")
	}

	staticGroupAddress := make(map[string]struct{})
	for _, subBlock := range block.StaticGroup {
		address := subBlock.Address.ValueString()
		if _, ok := staticGroupAddress[address]; ok {
			return configSet,
				pathRoot.AtName("static_group").AtName("*"),
				fmt.Errorf("multiple static_group blocks with the same address %q"+
					blockErrorSuffix, address)
		}
		staticGroupAddress[address] = struct{}{}

		configSet = append(configSet, setPrefix+"static group "+address)
		if v := subBlock.Source.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"static group "+address+" source "+v)
		}
	}

	return configSet, path.Empty(), nil
}

func (rscData *igmpSnoopingVlanData) read(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols igmp-snooping vlan " + name + junos.PipeDisplaySetRelative)
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
		rscData.fillID()
		for item := range strings.SplitSeq(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "immediate-leave":
				rscData.ImmediateLeave = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "interface "):
				interfaceName := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.Interface = tfdata.AppendPotentialNewBlock(rscData.Interface, types.StringValue(interfaceName))
				interFace := &rscData.Interface[len(rscData.Interface)-1]
				balt.CutPrefixInString(&itemTrim, interfaceName+" ")

				if err := interFace.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "l2-querier source-address "):
				rscData.L2QuerierSourceAddress = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "proxy"):
				rscData.Proxy = types.BoolValue(true)
				if balt.CutPrefixInString(&itemTrim, " source-address ") {
					rscData.ProxySourceAddress = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "query-interval "):
				rscData.QueryInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "query-last-member-interval "):
				rscData.QueryLastMemberInterval = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "query-response-interval "):
				rscData.QueryResponseInterval = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "robust-count "):
				rscData.RobustCount, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *igmpSnoopingVlanBlockInterface) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "group-limit "):
		block.GroupLimit, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "host-only-interface":
		block.HostOnlyInterface = types.BoolValue(true)
	case itemTrim == "immediate-leave":
		block.ImmediateLeave = types.BoolValue(true)
	case itemTrim == "multicast-router-interface":
		block.MulticastRouterInterface = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "static group "):
		address := tfdata.FirstElementOfJunosLine(itemTrim)
		var staticGroup igmpSnoopingVlanBlockInterfaceBlockStaticGroup
		block.StaticGroup, staticGroup = tfdata.ExtractBlock(block.StaticGroup, types.StringValue(address))

		if balt.CutPrefixInString(&itemTrim, address+" source ") {
			staticGroup.Source = types.StringValue(itemTrim)
		}
		block.StaticGroup = append(block.StaticGroup, staticGroup)
	}

	return nil
}

func (rscData *igmpSnoopingVlanData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}

	configSet := []string{
		delPrefix + "protocols igmp-snooping vlan " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
