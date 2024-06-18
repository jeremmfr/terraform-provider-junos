package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *bridgeDomain) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"routing_instance": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"community_vlans": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
					"domain_id": schema.Int64Attribute{
						Optional: true,
					},
					"domain_type_bridge": schema.BoolAttribute{
						Optional: true,
					},
					"isolated_vlan": schema.Int64Attribute{
						Optional: true,
					},
					"routing_interface": schema.StringAttribute{
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
								"decapsulate_accept_inner_vlan": schema.BoolAttribute{
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
			StateUpgrader: upgradeBridgeDomainStateV0toV1,
		},
	}
}

func upgradeBridgeDomainStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID               types.String   `tfsdk:"id"`
		Name             types.String   `tfsdk:"name"`
		RoutingInstance  types.String   `tfsdk:"routing_instance"`
		CommunityVlans   []types.String `tfsdk:"community_vlans"`
		Description      types.String   `tfsdk:"description"`
		DomainID         types.Int64    `tfsdk:"domain_id"`
		DomainTypeBridge types.Bool     `tfsdk:"domain_type_bridge"`
		IsolatedVLAN     types.Int64    `tfsdk:"isolated_vlan"`
		RoutingInterface types.String   `tfsdk:"routing_interface"`
		ServiceID        types.Int64    `tfsdk:"service_id"`
		VlanID           types.Int64    `tfsdk:"vlan_id"`
		VlanIDList       []types.String `tfsdk:"vlan_id_list"`
		Vxlan            []struct {
			Vni                        types.Int64  `tfsdk:"vni"`
			VniExtendEvpn              types.Bool   `tfsdk:"vni_extend_evpn"`
			DecapsulateAcceptInnerVlan types.Bool   `tfsdk:"decapsulate_accept_inner_vlan"`
			EncapsulateInnerVlan       types.Bool   `tfsdk:"encapsulate_inner_vlan"`
			IngressNodeReplication     types.Bool   `tfsdk:"ingress_node_replication"`
			MulticastGroup             types.String `tfsdk:"multicast_group"`
			OvsdbManaged               types.Bool   `tfsdk:"ovsdb_managed"`
			UnreachableVtepAgingTimer  types.Int64  `tfsdk:"unreachable_vtep_aging_timer"`
		} `tfsdk:"vxlan"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 bridgeDomainData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.RoutingInstance = dataV0.RoutingInstance
	dataV1.DomainTypeBridge = dataV0.DomainTypeBridge
	dataV1.CommunityVlans = dataV0.CommunityVlans
	dataV1.Description = dataV0.Description
	dataV1.DomainID = dataV0.DomainID
	dataV1.IsolatedVLAN = dataV0.IsolatedVLAN
	dataV1.RoutingInterface = dataV0.RoutingInterface
	dataV1.ServiceID = dataV0.ServiceID
	dataV1.VlanID = dataV0.VlanID
	dataV1.VlanIDList = dataV0.VlanIDList
	if len(dataV0.Vxlan) > 0 {
		dataV1.Vxlan = &bridgeDomainBlockVxlan{
			Vni:                        dataV0.Vxlan[0].Vni,
			VniExtendEvpn:              dataV0.Vxlan[0].VniExtendEvpn,
			DecapsulateAcceptInnerVlan: dataV0.Vxlan[0].DecapsulateAcceptInnerVlan,
			EncapsulateInnerVlan:       dataV0.Vxlan[0].EncapsulateInnerVlan,
			IngressNodeReplication:     dataV0.Vxlan[0].IngressNodeReplication,
			OvsdbManaged:               dataV0.Vxlan[0].OvsdbManaged,
			MulticastGroup:             dataV0.Vxlan[0].MulticastGroup,
			UnreachableVtepAgingTimer:  dataV0.Vxlan[0].UnreachableVtepAgingTimer,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
