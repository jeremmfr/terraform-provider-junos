package providerfwk

import (
	"context"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &bridgeDomain{}
	_ resource.ResourceWithConfigure      = &bridgeDomain{}
	_ resource.ResourceWithValidateConfig = &bridgeDomain{}
	_ resource.ResourceWithImportState    = &bridgeDomain{}
	_ resource.ResourceWithUpgradeState   = &bridgeDomain{}
)

type bridgeDomain struct {
	client *junos.Client
}

func newBridgeDomainResource() resource.Resource {
	return &bridgeDomain{}
}

func (rsc *bridgeDomain) typeName() string {
	return providerName + "_bridge_domain"
}

func (rsc *bridgeDomain) junosName() string {
	return "bridge-domain"
}

func (rsc *bridgeDomain) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *bridgeDomain) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *bridgeDomain) Configure(
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

func (rsc *bridgeDomain) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
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
				Description: "Bridge domain name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"community_vlans": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of Community VLANs for private vlan bridge domain.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringNumberRange(1, 4094).WithNameInError("VLAN"),
					),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of bridge domain.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"domain_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Domain-id for auto derived Route Target.",
				Validators: []validator.Int64{
					int64validator.Between(1, 15),
				},
			},
			"domain_type_bridge": schema.BoolAttribute{
				Optional:    true,
				Description: "Forwarding instance.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"interface": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Interface for this bridge domain.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
						tfvalidator.String1DotCount(),
					),
				},
			},
			"isolated_vlan": schema.Int64Attribute{
				Optional:    true,
				Description: "Isolated VLAN ID for private vlan bridge domain.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"routing_interface": schema.StringAttribute{
				Optional:    true,
				Description: "Routing interface name for this bridge-domain.",
				Validators: []validator.String{
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(irb|vlan)\.`),
						"must start with 'irb.' or 'vlan.'"),
				},
			},
			"service_id": schema.Int64Attribute{
				Optional: true,
				Description: "Service id required if bridge-domain is of type MC-AE and" +
					"vlan-id all or vlan-id none or vlan-tags.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"vlan_id": schema.Int64Attribute{
				Optional:    true,
				Description: "IEEE 802.1q VLAN identifier for bridging domain.",
				Validators: []validator.Int64{
					int64validator.Between(1, 4094),
				},
			},
			"vlan_id_list": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Create bridge-domain for each of the vlan-id specified in the vlan-id-list.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringNumberRange(1, 4094).WithNameInError("VLAN"),
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"vxlan": schema.SingleNestedBlock{
				Description: "Declare vxlan options.",
				Attributes: map[string]schema.Attribute{
					"vni": schema.Int64Attribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "VXLAN identifier (0..16777214).",
						Validators: []validator.Int64{
							int64validator.Between(0, 16777214),
						},
					},
					"vni_extend_evpn": schema.BoolAttribute{
						Optional:    true,
						Description: "Extend VNI to EVPN.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"decapsulate_accept_inner_vlan": schema.BoolAttribute{
						Optional:    true,
						Description: "Accept VXLAN packets with inner VLAN.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"encapsulate_inner_vlan": schema.BoolAttribute{
						Optional:    true,
						Description: "Retain inner VLAN in the packet.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ingress_node_replication": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable ingress node replication.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"multicast_group": schema.StringAttribute{
						Optional:    true,
						Description: "CIDR for Multicast group registered for VXLAN segment.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"ovsdb_managed": schema.BoolAttribute{
						Optional:    true,
						Description: "Bridge-domain is managed remotely via VXLAN OVSDB Controller.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"static_remote_vtep_list": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Configure bridge domain specific static remote VXLAN tunnel endpoints.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(
								tfvalidator.StringIPAddress().IPv4Only(),
							),
						},
					},
					"unreachable_vtep_aging_timer": schema.Int64Attribute{
						Optional:    true,
						Description: "Unreachable VXLAN tunnel endpoint removal timer.",
						Validators: []validator.Int64{
							int64validator.Between(300, 1800),
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

type bridgeDomainData struct {
	ID               types.String            `tfsdk:"id"`
	Name             types.String            `tfsdk:"name"`
	RoutingInstance  types.String            `tfsdk:"routing_instance"`
	CommunityVlans   []types.String          `tfsdk:"community_vlans"`
	Description      types.String            `tfsdk:"description"`
	DomainID         types.Int64             `tfsdk:"domain_id"`
	DomainTypeBridge types.Bool              `tfsdk:"domain_type_bridge"`
	Interface        []types.String          `tfsdk:"interface"`
	IsolatedVLAN     types.Int64             `tfsdk:"isolated_vlan"`
	RoutingInterface types.String            `tfsdk:"routing_interface"`
	ServiceID        types.Int64             `tfsdk:"service_id"`
	VlanID           types.Int64             `tfsdk:"vlan_id"`
	VlanIDList       []types.String          `tfsdk:"vlan_id_list"`
	Vxlan            *bridgeDomainBlockVxlan `tfsdk:"vxlan"`
}

type bridgeDomainConfig struct {
	ID               types.String                  `tfsdk:"id"`
	Name             types.String                  `tfsdk:"name"`
	RoutingInstance  types.String                  `tfsdk:"routing_instance"`
	CommunityVlans   types.Set                     `tfsdk:"community_vlans"`
	Description      types.String                  `tfsdk:"description"`
	DomainID         types.Int64                   `tfsdk:"domain_id"`
	DomainTypeBridge types.Bool                    `tfsdk:"domain_type_bridge"`
	Interface        types.Set                     `tfsdk:"interface"`
	IsolatedVLAN     types.Int64                   `tfsdk:"isolated_vlan"`
	RoutingInterface types.String                  `tfsdk:"routing_interface"`
	ServiceID        types.Int64                   `tfsdk:"service_id"`
	VlanID           types.Int64                   `tfsdk:"vlan_id"`
	VlanIDList       types.Set                     `tfsdk:"vlan_id_list"`
	Vxlan            *bridgeDomainBlockVxlanConfig `tfsdk:"vxlan"`
}

func (rscConfig *bridgeDomainConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscConfig, "ID", "Name", "RoutingInstance")
}

type bridgeDomainBlockVxlan struct {
	Vni                        types.Int64    `tfsdk:"vni"`
	VniExtendEvpn              types.Bool     `tfsdk:"vni_extend_evpn"`
	DecapsulateAcceptInnerVlan types.Bool     `tfsdk:"decapsulate_accept_inner_vlan"`
	EncapsulateInnerVlan       types.Bool     `tfsdk:"encapsulate_inner_vlan"`
	IngressNodeReplication     types.Bool     `tfsdk:"ingress_node_replication"`
	MulticastGroup             types.String   `tfsdk:"multicast_group"`
	OvsdbManaged               types.Bool     `tfsdk:"ovsdb_managed"`
	StaticRemoteVtepList       []types.String `tfsdk:"static_remote_vtep_list"`
	UnreachableVtepAgingTimer  types.Int64    `tfsdk:"unreachable_vtep_aging_timer"`
}

type bridgeDomainBlockVxlanConfig struct {
	Vni                        types.Int64  `tfsdk:"vni"`
	VniExtendEvpn              types.Bool   `tfsdk:"vni_extend_evpn"`
	DecapsulateAcceptInnerVlan types.Bool   `tfsdk:"decapsulate_accept_inner_vlan"`
	EncapsulateInnerVlan       types.Bool   `tfsdk:"encapsulate_inner_vlan"`
	IngressNodeReplication     types.Bool   `tfsdk:"ingress_node_replication"`
	MulticastGroup             types.String `tfsdk:"multicast_group"`
	OvsdbManaged               types.Bool   `tfsdk:"ovsdb_managed"`
	StaticRemoteVtepList       types.Set    `tfsdk:"static_remote_vtep_list"`
	UnreachableVtepAgingTimer  types.Int64  `tfsdk:"unreachable_vtep_aging_timer"`
}

func (rsc *bridgeDomain) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config bridgeDomainConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one of arguments need to be set (in addition to `name` and `routing_instance`)",
		)
	}

	if !config.VlanID.IsNull() && !config.VlanID.IsUnknown() &&
		!config.VlanIDList.IsNull() && !config.VlanIDList.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("vlan_id"),
			tfdiag.ConflictConfigErrSummary,
			"vlan_id and vlan_id_list cannot be configured together",
		)
	}
	if config.Vxlan != nil {
		if config.Vxlan.Vni.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("vxlan").AtName("vni"),
				tfdiag.MissingConfigErrSummary,
				"vni must be specified in vxlan block",
			)
		}
	}
}

func (rsc *bridgeDomain) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan bridgeDomainData
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
			if !junSess.CheckCompatibilityRouter() {
				resp.Diagnostics.AddError(
					tfdiag.CompatibilityErrSummary,
					fmt.Sprintf(rsc.junosName()+" not compatible "+
						"with Junos device %q", junSess.SystemInformation.HardwareModel),
				)

				return false
			}
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
			domainExists, err := checkBridgeDomainExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if domainExists {
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
			domainExists, err := checkBridgeDomainExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !domainExists {
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

func (rsc *bridgeDomain) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data bridgeDomainData
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

func (rsc *bridgeDomain) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state bridgeDomainData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataDelWithOpts = &state
	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *bridgeDomain) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state bridgeDomainData
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

func (rsc *bridgeDomain) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data bridgeDomainData

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

func checkBridgeDomainExists(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"bridge-domains \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *bridgeDomainData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *bridgeDomainData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *bridgeDomainData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefixProtocolsEvpn := setPrefix + "protocols evpn "
	setPrefix += "bridge-domains \"" + rscData.Name.ValueString() + "\" "

	for _, v := range rscData.CommunityVlans {
		configSet = append(configSet, setPrefix+"community-vlans "+v.ValueString())
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if !rscData.DomainID.IsNull() {
		configSet = append(configSet, setPrefix+"domain-id "+
			utils.ConvI64toa(rscData.DomainID.ValueInt64()))
	}
	if rscData.DomainTypeBridge.ValueBool() {
		configSet = append(configSet, setPrefix+"domain-type bridge")
	}
	for _, v := range rscData.Interface {
		configSet = append(configSet, setPrefix+"interface "+v.ValueString())
	}
	if !rscData.IsolatedVLAN.IsNull() {
		configSet = append(configSet, setPrefix+"isolated-vlan "+
			utils.ConvI64toa(rscData.IsolatedVLAN.ValueInt64()))
	}
	if v := rscData.RoutingInterface.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-interface "+v)
	}
	if !rscData.ServiceID.IsNull() {
		configSet = append(configSet, setPrefix+"service-id "+
			utils.ConvI64toa(rscData.ServiceID.ValueInt64()))
	}
	if !rscData.VlanID.IsNull() {
		configSet = append(configSet, setPrefix+"vlan-id "+
			utils.ConvI64toa(rscData.VlanID.ValueInt64()))
	}
	for _, v := range rscData.VlanIDList {
		configSet = append(configSet, setPrefix+"vlan-id-list "+v.ValueString())
	}
	if rscData.Vxlan != nil {
		configSet = append(configSet, setPrefix+"vxlan vni "+
			utils.ConvI64toa(rscData.Vxlan.Vni.ValueInt64()))

		if rscData.Vxlan.VniExtendEvpn.ValueBool() {
			configSet = append(configSet, setPrefixProtocolsEvpn+"extended-vni-list "+
				utils.ConvI64toa(rscData.Vxlan.Vni.ValueInt64()))
		}
		if rscData.Vxlan.DecapsulateAcceptInnerVlan.ValueBool() {
			configSet = append(configSet, setPrefix+"vxlan decapsulate-accept-inner-vlan")
		}
		if rscData.Vxlan.EncapsulateInnerVlan.ValueBool() {
			configSet = append(configSet, setPrefix+"vxlan encapsulate-inner-vlan")
		}
		if rscData.Vxlan.IngressNodeReplication.ValueBool() {
			configSet = append(configSet, setPrefix+"vxlan ingress-node-replication")
		}
		if v := rscData.Vxlan.MulticastGroup.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"vxlan multicast-group "+v)
		}
		if rscData.Vxlan.OvsdbManaged.ValueBool() {
			configSet = append(configSet, setPrefix+"vxlan ovsdb-managed")
		}
		for _, v := range rscData.Vxlan.StaticRemoteVtepList {
			configSet = append(configSet, setPrefix+"vxlan static-remote-vtep-list "+v.ValueString())
		}
		if !rscData.Vxlan.UnreachableVtepAgingTimer.IsNull() {
			configSet = append(configSet, setPrefix+"vxlan unreachable-vtep-aging-timer "+
				utils.ConvI64toa(rscData.Vxlan.UnreachableVtepAgingTimer.ValueInt64()))
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *bridgeDomainData) read(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"bridge-domains \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "community-vlans "):
				rscData.CommunityVlans = append(rscData.CommunityVlans, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "domain-id "):
				rscData.DomainID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "domain-type bridge":
				rscData.DomainTypeBridge = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "interface "):
				rscData.Interface = append(rscData.Interface, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "isolated-vlan "):
				rscData.IsolatedVLAN, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "routing-interface "):
				rscData.RoutingInterface = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "service-id "):
				rscData.ServiceID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "vlan-id "):
				rscData.VlanID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "vlan-id-list "):
				rscData.VlanIDList = append(rscData.VlanIDList, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "vxlan "):
				if rscData.Vxlan == nil {
					rscData.Vxlan = &bridgeDomainBlockVxlan{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "vni "):
					rscData.Vxlan.Vni, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
					showConfigEvpn, err := junSess.Command(showPrefix + "protocols evpn" + junos.PipeDisplaySetRelative)
					if err != nil {
						return err
					}
					if showConfigEvpn != junos.EmptyW {
						for _, itemEvpn := range strings.Split(showConfigEvpn, "\n") {
							if strings.Contains(itemEvpn, junos.XMLStartTagConfigOut) {
								continue
							}
							if strings.Contains(itemEvpn, junos.XMLEndTagConfigOut) {
								break
							}
							if itemEvpn == junos.SetLS+"extended-vni-list "+itemTrim {
								rscData.Vxlan.VniExtendEvpn = types.BoolValue(true)

								break
							}
						}
					}
				case itemTrim == "decapsulate-accept-inner-vlan":
					rscData.Vxlan.DecapsulateAcceptInnerVlan = types.BoolValue(true)
				case itemTrim == "encapsulate-inner-vlan":
					rscData.Vxlan.EncapsulateInnerVlan = types.BoolValue(true)
				case itemTrim == "ingress-node-replication":
					rscData.Vxlan.IngressNodeReplication = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "multicast-group "):
					rscData.Vxlan.MulticastGroup = types.StringValue(itemTrim)
				case itemTrim == "ovsdb-managed":
					rscData.Vxlan.OvsdbManaged = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "static-remote-vtep-list "):
					rscData.Vxlan.StaticRemoteVtepList = append(rscData.Vxlan.StaticRemoteVtepList, types.StringValue(itemTrim))
				case balt.CutPrefixInString(&itemTrim, "unreachable-vtep-aging-timer "):
					rscData.Vxlan.UnreachableVtepAgingTimer, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (rscData *bridgeDomainData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefixProtocolsEvpn := delPrefix + "protocols evpn "
	delPrefix += "bridge-domains \"" + rscData.Name.ValueString() + "\" "

	configSet := []string{
		delPrefix + "community-vlans",
		delPrefix + "description",
		delPrefix + "domain-id",
		delPrefix + "domain-type",
		delPrefix + "isolated-vlan",
		delPrefix + "routing-interface",
		delPrefix + "service-id",
		delPrefix + "vlan-id",
		delPrefix + "vlan-id-list",
		delPrefix + "vxlan",
	}
	for _, v := range rscData.Interface {
		configSet = append(configSet, delPrefix+"interface "+v.ValueString())
	}
	if rscData.Vxlan != nil {
		if rscData.Vxlan.VniExtendEvpn.ValueBool() {
			configSet = append(configSet, delPrefixProtocolsEvpn+
				"extended-vni-list "+utils.ConvI64toa(rscData.Vxlan.Vni.ValueInt64()))
		}
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *bridgeDomainData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}

	configSet := []string{
		delPrefix + "bridge-domains \"" + rscData.Name.ValueString() + "\"",
	}
	if rscData.Vxlan != nil {
		if rscData.Vxlan.VniExtendEvpn.ValueBool() {
			configSet = append(configSet, delPrefix+
				"protocols evpn extended-vni-list "+utils.ConvI64toa(rscData.Vxlan.Vni.ValueInt64()))
		}
	}

	return junSess.ConfigSet(configSet)
}
