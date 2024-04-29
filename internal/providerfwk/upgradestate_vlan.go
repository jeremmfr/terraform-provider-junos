package providerfwk

import (
	"context"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *vlan) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"community_vlans": schema.SetAttribute{
						ElementType: types.Int64Type,
						Optional:    true,
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
					"forward_filter_input": schema.StringAttribute{
						Optional: true,
					},
					"forward_filter_output": schema.StringAttribute{
						Optional: true,
					},
					"forward_flood_input": schema.StringAttribute{
						Optional: true,
					},
					"isolated_vlan": schema.Int64Attribute{
						Optional: true,
					},
					"l3_interface": schema.StringAttribute{
						Optional: true,
					},
					"private_vlan": schema.StringAttribute{
						Optional: true,
					},
					"service_id": schema.Int64Attribute{
						Optional: true,
					},
					"vlan_id": schema.Int64Attribute{
						Optional: true,
					},
					"vlan_id_list": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"vxlan": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"vni": schema.Int64Attribute{
									Required: true,
								},
								"vni_extend_evpn": schema.BoolAttribute{
									Optional: true,
								},
								"encapsulate_inner_vlan": schema.BoolAttribute{
									Optional: true,
								},
								"ingress_node_replication": schema.BoolAttribute{
									Optional: true,
								},
								"multicast_group": schema.StringAttribute{
									Optional: true,
								},
								"ovsdb_managed": schema.BoolAttribute{
									Optional: true,
								},
								"unreachable_vtep_aging_timer": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeVlanV0toV1,
		},
	}
}

func upgradeVlanV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                  types.String   `tfsdk:"id"`
		Name                types.String   `tfsdk:"name"`
		CommunityVlans      []types.Int64  `tfsdk:"community_vlans"`
		Description         types.String   `tfsdk:"description"`
		ForwardFilterInput  types.String   `tfsdk:"forward_filter_input"`
		ForwardFilterOutput types.String   `tfsdk:"forward_filter_output"`
		ForwardFloodInput   types.String   `tfsdk:"forward_flood_input"`
		IsolatedVlan        types.Int64    `tfsdk:"isolated_vlan"`
		L3Interface         types.String   `tfsdk:"l3_interface"`
		PrivateVlan         types.String   `tfsdk:"private_vlan"`
		ServiceID           types.Int64    `tfsdk:"service_id"`
		VlanID              types.Int64    `tfsdk:"vlan_id"`
		VlanIDList          []types.String `tfsdk:"vlan_id_list"`
		Vxlan               []struct {
			Vni                       types.Int64  `tfsdk:"vni"`
			VniExtendEvpn             types.Bool   `tfsdk:"vni_extend_evpn"`
			EncapsulateInnerVlan      types.Bool   `tfsdk:"encapsulate_inner_vlan"`
			IngressNodeReplication    types.Bool   `tfsdk:"ingress_node_replication"`
			MulticastGroup            types.String `tfsdk:"multicast_group"`
			OvsdbManaged              types.Bool   `tfsdk:"ovsdb_managed"`
			UnreachableVtepAgingTimer types.Int64  `tfsdk:"unreachable_vtep_aging_timer"`
		} `tfsdk:"vxlan"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 vlanData
	dataV1.Name = dataV0.Name
	dataV1.RoutingInstance = types.StringValue(junos.DefaultW)
	dataV1.fillID()
	dataV1.CommunityVlans = make([]types.String, len(dataV0.CommunityVlans))
	for i, v := range dataV0.CommunityVlans {
		dataV1.CommunityVlans[i] = types.StringValue(utils.ConvI64toa(v.ValueInt64()))
	}
	dataV1.Description = dataV0.Description
	dataV1.ForwardFilterInput = dataV0.ForwardFilterInput
	dataV1.ForwardFilterOutput = dataV0.ForwardFilterOutput
	dataV1.ForwardFloodInput = dataV0.ForwardFloodInput
	if dataV0.IsolatedVlan.ValueInt64() != 0 {
		dataV1.IsolatedVlan = types.StringValue(utils.ConvI64toa(dataV0.IsolatedVlan.ValueInt64()))
	}
	dataV1.L3Interface = dataV0.L3Interface
	dataV1.PrivateVlan = dataV0.PrivateVlan
	dataV1.ServiceID = dataV0.ServiceID
	if dataV0.VlanID.ValueInt64() != 0 {
		dataV1.VlanID = types.StringValue(utils.ConvI64toa(dataV0.VlanID.ValueInt64()))
	}
	dataV1.VlanIDList = dataV0.VlanIDList
	if len(dataV0.Vxlan) > 0 {
		dataV1.Vxlan = &vlanBlockVxlan{
			Vni:                       dataV0.Vxlan[0].Vni,
			VniExtendEvpn:             dataV0.Vxlan[0].VniExtendEvpn,
			EncapsulateInnerVlan:      dataV0.Vxlan[0].EncapsulateInnerVlan,
			IngressNodeReplication:    dataV0.Vxlan[0].IngressNodeReplication,
			MulticastGroup:            dataV0.Vxlan[0].MulticastGroup,
			OvsdbManaged:              dataV0.Vxlan[0].OvsdbManaged,
			UnreachableVtepAgingTimer: dataV0.Vxlan[0].UnreachableVtepAgingTimer,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
