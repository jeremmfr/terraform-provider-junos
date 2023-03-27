package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityIkeGateway) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"external_interface": schema.StringAttribute{
						Required: true,
					},
					"policy": schema.StringAttribute{
						Required: true,
					},
					"address": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"general_ike_id": schema.BoolAttribute{
						Optional: true,
					},
					"local_address": schema.StringAttribute{
						Optional: true,
					},
					"no_nat_traversal": schema.BoolAttribute{
						Optional: true,
					},
					"version": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"dynamic_remote": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"connections_limit": schema.Int64Attribute{
									Optional: true,
								},
								"hostname": schema.StringAttribute{
									Optional: true,
								},
								"ike_user_type": schema.StringAttribute{
									Optional: true,
								},
								"inet": schema.StringAttribute{
									Optional: true,
								},
								"inet6": schema.StringAttribute{
									Optional: true,
								},
								"reject_duplicate_connection": schema.BoolAttribute{
									Optional: true,
								},
								"user_at_hostname": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"distinguished_name": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"container": schema.StringAttribute{
												Optional: true,
											},
											"wildcard": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"aaa": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"access_profile": schema.StringAttribute{
									Optional: true,
								},
								"client_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"client_username": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"dead_peer_detection": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"interval": schema.Int64Attribute{
									Optional: true,
								},
								"send_mode": schema.StringAttribute{
									Optional: true,
								},
								"threshold": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"local_identity": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required: true,
								},
								"value": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"remote_identity": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required: true,
								},
								"value": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSecurityPolicyV0toV1,
		},
	}
}

func upgradeSecurityPolicyV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	//nolint:lll
	type modelV0 struct {
		GeneralIkeID      types.Bool                                 `tfsdk:"general_ike_id"`
		NoNatTraversal    types.Bool                                 `tfsdk:"no_nat_traversal"`
		ID                types.String                               `tfsdk:"id"`
		Name              types.String                               `tfsdk:"name"`
		ExternalInterface types.String                               `tfsdk:"external_interface"`
		Policy            types.String                               `tfsdk:"policy"`
		Address           []types.String                             `tfsdk:"address"`
		LocalAddress      types.String                               `tfsdk:"local_address"`
		Version           types.String                               `tfsdk:"version"`
		Aaa               []securityIkeGatewayBlockAaa               `tfsdk:"aaa"`
		DeadPeerDetection []securityIkeGatewayBlockDeadPeerDetection `tfsdk:"dead_peer_detection"`
		DynamicRemote     []struct {
			ConnectionsLimit          types.Int64                                                  `tfsdk:"connections_limit"`
			Hostname                  types.String                                                 `tfsdk:"hostname"`
			IkeUserType               types.String                                                 `tfsdk:"ike_user_type"`
			Inet                      types.String                                                 `tfsdk:"inet"`
			Inet6                     types.String                                                 `tfsdk:"inet6"`
			RejectDuplicateConnection types.Bool                                                   `tfsdk:"reject_duplicate_connection"`
			UserAtHostname            types.String                                                 `tfsdk:"user_at_hostname"`
			DistinguishedName         []securityIkeGatewayBlockDynamicRemoteBlockDistinguishedName `tfsdk:"distinguished_name"`
		} `tfsdk:"dynamic_remote"`
		LocalIdentity  []securityIkeGatewayBlockLocalIdentity  `tfsdk:"local_identity"`
		RemoteIdentity []securityIkeGatewayBlockRemoteIdentity `tfsdk:"remote_identity"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityIkeGatewayData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.ExternalInterface = dataV0.ExternalInterface
	dataV1.Policy = dataV0.Policy
	dataV1.Address = dataV0.Address
	dataV1.GeneralIkeID = dataV0.GeneralIkeID
	dataV1.LocalAddress = dataV0.LocalAddress
	dataV1.NoNatTraversal = dataV0.NoNatTraversal
	dataV1.Version = dataV0.Version
	if len(dataV0.DynamicRemote) > 0 {
		dataV1.DynamicRemote = &securityIkeGatewayBlockDynamicRemote{}
		dataV1.DynamicRemote.ConnectionsLimit = dataV0.DynamicRemote[0].ConnectionsLimit
		dataV1.DynamicRemote.Hostname = dataV0.DynamicRemote[0].Hostname
		dataV1.DynamicRemote.IkeUserType = dataV0.DynamicRemote[0].IkeUserType
		dataV1.DynamicRemote.Inet = dataV0.DynamicRemote[0].Inet
		dataV1.DynamicRemote.Inet6 = dataV0.DynamicRemote[0].Inet6
		dataV1.DynamicRemote.RejectDuplicateConnection = dataV0.DynamicRemote[0].RejectDuplicateConnection
		dataV1.DynamicRemote.UserAtHostname = dataV0.DynamicRemote[0].UserAtHostname
		if len(dataV0.DynamicRemote[0].DistinguishedName) > 0 {
			dataV1.DynamicRemote.DistinguishedName = &securityIkeGatewayBlockDynamicRemoteBlockDistinguishedName{}
			dataV1.DynamicRemote.DistinguishedName.Container = dataV0.DynamicRemote[0].DistinguishedName[0].Container
			dataV1.DynamicRemote.DistinguishedName.Wildcard = dataV0.DynamicRemote[0].DistinguishedName[0].Wildcard
		}
	}
	if len(dataV0.Aaa) > 0 {
		dataV1.Aaa = &dataV0.Aaa[0]
	}
	if len(dataV0.DeadPeerDetection) > 0 {
		dataV1.DeadPeerDetection = &dataV0.DeadPeerDetection[0]
	}
	if len(dataV0.LocalIdentity) > 0 {
		dataV1.LocalIdentity = &dataV0.LocalIdentity[0]
	}
	if len(dataV0.RemoteIdentity) > 0 {
		dataV1.RemoteIdentity = &dataV0.RemoteIdentity[0]
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
