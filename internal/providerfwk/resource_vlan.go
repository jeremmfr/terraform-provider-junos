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
	_ resource.Resource                   = &vlan{}
	_ resource.ResourceWithConfigure      = &vlan{}
	_ resource.ResourceWithValidateConfig = &vlan{}
	_ resource.ResourceWithImportState    = &vlan{}
	_ resource.ResourceWithUpgradeState   = &vlan{}
)

type vlan struct {
	client *junos.Client
}

func newVlanResource() resource.Resource {
	return &vlan{}
}

func (rsc *vlan) typeName() string {
	return providerName + "_vlan"
}

func (rsc *vlan) junosName() string {
	return "vlans"
}

func (rsc *vlan) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *vlan) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *vlan) Configure(
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

func (rsc *vlan) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<name>_-_<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of VLAN.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 64),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for vlan if not root level.",
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
				Description: "List of VLAN id or name of community vlans for primary vlan.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 64),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of VLAN.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"forward_filter_input": schema.StringAttribute{
				Optional:    true,
				Description: "Input filter to apply for forwarded packets.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"forward_filter_output": schema.StringAttribute{
				Optional:    true,
				Description: "Output filter to apply for forwarded packets.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"forward_flood_input": schema.StringAttribute{
				Optional:    true,
				Description: "Input filter to apply for ethernet switching flood packets.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"isolated_vlan": schema.StringAttribute{
				Optional:    true,
				Description: "VLAN id or name of isolated vlan for primary vlan.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"l3_interface": schema.StringAttribute{
				Optional:    true,
				Description: "L3 interface name for this VLAN.",
				Validators: []validator.String{
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(irb|vlan)\.`),
						"must start with 'irb.' or 'vlan.'"),
				},
			},
			"no_arp_suppression": schema.BoolAttribute{
				Optional:    true,
				Description: "Turn off ARP suppression.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"private_vlan": schema.StringAttribute{
				Optional:    true,
				Description: "Type of secondary VLAN for private vlan.",
				Validators: []validator.String{
					stringvalidator.OneOf("community", "isolated"),
				},
			},
			"service_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Service id.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"vlan_id": schema.StringAttribute{
				Optional:    true,
				Description: "802.1q VLAN id or `all` or `none`.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(409[0-4]|(40[0-8]|[1-3]\d\d|[1-9]\d|[1-9])\d|[1-9]|all|none)$`),
						"must be a VLAN id (1..4094) or all or none"),
				},
			},
			"vlan_id_list": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of 802.1q VLAN id.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(
							`^(409[0-4]|(40[0-8]|[1-3]\d\d|[1-9]\d|[1-9])\d|[1-9])`+
								`(-(409[0-4]|(40[0-8]|[1-3]\d\d|[1-9]\d|[1-9])\d|[1-9]))?$`),
							"must be a VLAN id (1..4094) or a range of VLAN id (1..4094)-(1..4094)"),
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"vxlan": schema.SingleNestedBlock{
				Description: "Declare vxlan configuration.",
				Attributes: map[string]schema.Attribute{
					"vni": schema.Int64Attribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "VXLAN identifier.",
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
						Description: "Multicast group registered for VXLAN segment.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress().IPv4Only(),
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
						Description: "Configure vlan specific static remote VXLAN tunnel endpoints.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(
								tfvalidator.StringIPAddress().IPv4Only(),
							),
						},
					},
					"translation_vni": schema.Int64Attribute{
						Optional:    true,
						Description: "Translated VXLAN identifier.",
						Validators: []validator.Int64{
							int64validator.Between(1, 16777214),
						},
					},
					"unreachable_vtep_aging_timer": schema.Int64Attribute{
						Optional:    true,
						Description: "Unreachable VXLAN tunnel endpoint removal timer (seconds).",
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

type vlanData struct {
	ID                  types.String    `tfsdk:"id"`
	Name                types.String    `tfsdk:"name"`
	RoutingInstance     types.String    `tfsdk:"routing_instance"`
	CommunityVlans      []types.String  `tfsdk:"community_vlans"`
	Description         types.String    `tfsdk:"description"`
	ForwardFilterInput  types.String    `tfsdk:"forward_filter_input"`
	ForwardFilterOutput types.String    `tfsdk:"forward_filter_output"`
	ForwardFloodInput   types.String    `tfsdk:"forward_flood_input"`
	IsolatedVlan        types.String    `tfsdk:"isolated_vlan"`
	L3Interface         types.String    `tfsdk:"l3_interface"`
	NoARPSuppression    types.Bool      `tfsdk:"no_arp_suppression"`
	PrivateVlan         types.String    `tfsdk:"private_vlan"`
	ServiceID           types.Int64     `tfsdk:"service_id"`
	VlanID              types.String    `tfsdk:"vlan_id"`
	VlanIDList          []types.String  `tfsdk:"vlan_id_list"`
	Vxlan               *vlanBlockVxlan `tfsdk:"vxlan"`
}

type vlanConfig struct {
	ID                  types.String          `tfsdk:"id"`
	Name                types.String          `tfsdk:"name"`
	RoutingInstance     types.String          `tfsdk:"routing_instance"`
	CommunityVlans      types.Set             `tfsdk:"community_vlans"`
	Description         types.String          `tfsdk:"description"`
	ForwardFilterInput  types.String          `tfsdk:"forward_filter_input"`
	ForwardFilterOutput types.String          `tfsdk:"forward_filter_output"`
	ForwardFloodInput   types.String          `tfsdk:"forward_flood_input"`
	IsolatedVlan        types.String          `tfsdk:"isolated_vlan"`
	L3Interface         types.String          `tfsdk:"l3_interface"`
	NoARPSuppression    types.Bool            `tfsdk:"no_arp_suppression"`
	PrivateVlan         types.String          `tfsdk:"private_vlan"`
	ServiceID           types.Int64           `tfsdk:"service_id"`
	VlanID              types.String          `tfsdk:"vlan_id"`
	VlanIDList          types.Set             `tfsdk:"vlan_id_list"`
	Vxlan               *vlanBlockVxlanConfig `tfsdk:"vxlan"`
}

func (rscConfig *vlanConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscConfig, "ID", "Name", "RoutingInstance")
}

type vlanBlockVxlan struct {
	Vni                       types.Int64    `tfsdk:"vni"`
	VniExtendEvpn             types.Bool     `tfsdk:"vni_extend_evpn"`
	EncapsulateInnerVlan      types.Bool     `tfsdk:"encapsulate_inner_vlan"`
	IngressNodeReplication    types.Bool     `tfsdk:"ingress_node_replication"`
	MulticastGroup            types.String   `tfsdk:"multicast_group"`
	OvsdbManaged              types.Bool     `tfsdk:"ovsdb_managed"`
	StaticRemoteVtepList      []types.String `tfsdk:"static_remote_vtep_list"`
	TranslationVni            types.Int64    `tfsdk:"translation_vni"`
	UnreachableVtepAgingTimer types.Int64    `tfsdk:"unreachable_vtep_aging_timer"`
}

type vlanBlockVxlanConfig struct {
	Vni                       types.Int64  `tfsdk:"vni"`
	VniExtendEvpn             types.Bool   `tfsdk:"vni_extend_evpn"`
	EncapsulateInnerVlan      types.Bool   `tfsdk:"encapsulate_inner_vlan"`
	IngressNodeReplication    types.Bool   `tfsdk:"ingress_node_replication"`
	MulticastGroup            types.String `tfsdk:"multicast_group"`
	OvsdbManaged              types.Bool   `tfsdk:"ovsdb_managed"`
	StaticRemoteVtepList      types.Set    `tfsdk:"static_remote_vtep_list"`
	TranslationVni            types.Int64  `tfsdk:"translation_vni"`
	UnreachableVtepAgingTimer types.Int64  `tfsdk:"unreachable_vtep_aging_timer"`
}

func (rsc *vlan) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config vlanConfig
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
		if !config.Vxlan.IngressNodeReplication.IsNull() && !config.Vxlan.IngressNodeReplication.IsUnknown() &&
			!config.Vxlan.MulticastGroup.IsNull() && !config.Vxlan.MulticastGroup.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ingress_node_replication"),
				tfdiag.ConflictConfigErrSummary,
				"ingress_node_replication and multicast_group cannot be configured together",
			)
		}
		if !config.Vxlan.IngressNodeReplication.IsNull() && !config.Vxlan.IngressNodeReplication.IsUnknown() &&
			!config.Vxlan.OvsdbManaged.IsNull() && !config.Vxlan.OvsdbManaged.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ingress_node_replication"),
				tfdiag.ConflictConfigErrSummary,
				"ingress_node_replication and ovsdb_managed cannot be configured together",
			)
		}
		if !config.Vxlan.MulticastGroup.IsNull() && !config.Vxlan.MulticastGroup.IsUnknown() &&
			!config.Vxlan.OvsdbManaged.IsNull() && !config.Vxlan.OvsdbManaged.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("multicast_group"),
				tfdiag.ConflictConfigErrSummary,
				"multicast_group and ovsdb_managed cannot be configured together",
			)
		}
	}
}

func (rsc *vlan) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan vlanData
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
			vlanExists, err := checkVlanExists(fnCtx, plan.Name.ValueString(), plan.RoutingInstance.ValueString(), junSess)
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
			vlanExists, err := checkVlanExists(fnCtx, plan.Name.ValueString(), plan.RoutingInstance.ValueString(), junSess)
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

func (rsc *vlan) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data vlanData
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

func (rsc *vlan) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state vlanData
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

func (rsc *vlan) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state vlanData
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

func (rsc *vlan) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	var data vlanData
	idSplit := strings.Split(req.ID, junos.IDSeparator)
	if len(idSplit) > 1 {
		if err := data.read(ctx, idSplit[0], idSplit[1], junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
	} else {
		if err := data.read(ctx, idSplit[0], junos.DefaultW, junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be <name> or <name>"+junos.IDSeparator+"<routing_instance>)",
		)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkVlanExists(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"vlans " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *vlanData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *vlanData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *vlanData) set(
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
	setPrefix += "vlans " + rscData.Name.ValueString() + " "

	for _, v := range rscData.CommunityVlans {
		configSet = append(configSet, setPrefix+"community-vlans "+v.ValueString())
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := rscData.ForwardFilterInput.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"forwarding-options filter input \""+v+"\"")
	}
	if v := rscData.ForwardFilterOutput.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"forwarding-options filter output \""+v+"\"")
	}
	if v := rscData.ForwardFloodInput.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"forwarding-options flood input \""+v+"\"")
	}
	if v := rscData.IsolatedVlan.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"isolated-vlan "+v)
	}
	if v := rscData.L3Interface.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"l3-interface "+v)
	}
	if rscData.NoARPSuppression.ValueBool() {
		configSet = append(configSet, setPrefix+"no-arp-suppression")
	}
	if v := rscData.PrivateVlan.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"private-vlan "+v)
	}
	if !rscData.ServiceID.IsNull() {
		configSet = append(configSet, setPrefix+"service-id "+
			utils.ConvI64toa(rscData.ServiceID.ValueInt64()))
	}
	if v := rscData.VlanID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"vlan-id "+v)
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
		if !rscData.Vxlan.TranslationVni.IsNull() {
			configSet = append(configSet, setPrefix+"vxlan translation-vni "+
				utils.ConvI64toa(rscData.Vxlan.TranslationVni.ValueInt64()))
		}
		if !rscData.Vxlan.UnreachableVtepAgingTimer.IsNull() {
			configSet = append(configSet, setPrefix+"vxlan unreachable-vtep-aging-timer "+
				utils.ConvI64toa(rscData.Vxlan.UnreachableVtepAgingTimer.ValueInt64()))
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *vlanData) read(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"vlans " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "forwarding-options filter input "):
				rscData.ForwardFilterInput = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "forwarding-options filter output "):
				rscData.ForwardFilterOutput = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "forwarding-options flood input "):
				rscData.ForwardFloodInput = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "isolated-vlan "):
				rscData.IsolatedVlan = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "l3-interface "):
				rscData.L3Interface = types.StringValue(itemTrim)
			case itemTrim == "no-arp-suppression":
				rscData.NoARPSuppression = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "private-vlan "):
				rscData.PrivateVlan = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "service-id "):
				rscData.ServiceID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "vlan-id "):
				rscData.VlanID = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "vlan-id-list "):
				rscData.VlanIDList = append(rscData.VlanIDList, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "vxlan "):
				if rscData.Vxlan == nil {
					rscData.Vxlan = &vlanBlockVxlan{}
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
							}
						}
					}
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
				case balt.CutPrefixInString(&itemTrim, "translation-vni "):
					rscData.Vxlan.TranslationVni, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
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

func (rscData *vlanData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}

	configSet := []string{
		delPrefix + "vlans " + rscData.Name.ValueString(),
	}
	if rscData.Vxlan != nil && rscData.Vxlan.VniExtendEvpn.ValueBool() {
		configSet = append(configSet, delPrefix+"protocols evpn extended-vni-list "+
			utils.ConvI64toa(rscData.Vxlan.Vni.ValueInt64()))
	}

	return junSess.ConfigSet(configSet)
}
