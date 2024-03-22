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
			StateUpgrader: upgradeSecurityIkeGatewayV0toV1,
		},
	}
}

func upgradeSecurityIkeGatewayV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                types.String   `tfsdk:"id"`
		Name              types.String   `tfsdk:"name"`
		ExternalInterface types.String   `tfsdk:"external_interface"`
		Policy            types.String   `tfsdk:"policy"`
		Address           []types.String `tfsdk:"address"`
		GeneralIkeID      types.Bool     `tfsdk:"general_ike_id"`
		LocalAddress      types.String   `tfsdk:"local_address"`
		NoNatTraversal    types.Bool     `tfsdk:"no_nat_traversal"`
		Version           types.String   `tfsdk:"version"`
		DynamicRemote     []struct {
			ConnectionsLimit          types.Int64  `tfsdk:"connections_limit"`
			Hostname                  types.String `tfsdk:"hostname"`
			IkeUserType               types.String `tfsdk:"ike_user_type"`
			Inet                      types.String `tfsdk:"inet"`
			Inet6                     types.String `tfsdk:"inet6"`
			RejectDuplicateConnection types.Bool   `tfsdk:"reject_duplicate_connection"`
			UserAtHostname            types.String `tfsdk:"user_at_hostname"`
			DistinguishedName         []struct {
				Container types.String `tfsdk:"container"`
				Wildcard  types.String `tfsdk:"wildcard"`
			} `tfsdk:"distinguished_name"`
		} `tfsdk:"dynamic_remote"`
		Aaa []struct {
			AccessProfile  types.String `tfsdk:"access_profile"`
			ClientPassword types.String `tfsdk:"client_password"`
			ClientUsername types.String `tfsdk:"client_username"`
		} `tfsdk:"aaa"`
		DeadPeerDetection []struct {
			Interval  types.Int64  `tfsdk:"interval"`
			SendMode  types.String `tfsdk:"send_mode"`
			Threshold types.Int64  `tfsdk:"threshold"`
		} `tfsdk:"dead_peer_detection"`
		LocalIdentity []struct {
			Type  types.String `tfsdk:"type"`
			Value types.String `tfsdk:"value"`
		} `tfsdk:"local_identity"`
		RemoteIdentity []struct {
			Type  types.String `tfsdk:"type"`
			Value types.String `tfsdk:"value"`
		} `tfsdk:"remote_identity"`
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
		dataV1.DynamicRemote = &securityIkeGatewayBlockDynamicRemote{
			ConnectionsLimit:          dataV0.DynamicRemote[0].ConnectionsLimit,
			Hostname:                  dataV0.DynamicRemote[0].Hostname,
			IkeUserType:               dataV0.DynamicRemote[0].IkeUserType,
			Inet:                      dataV0.DynamicRemote[0].Inet,
			Inet6:                     dataV0.DynamicRemote[0].Inet6,
			RejectDuplicateConnection: dataV0.DynamicRemote[0].RejectDuplicateConnection,
			UserAtHostname:            dataV0.DynamicRemote[0].UserAtHostname,
		}
		if len(dataV0.DynamicRemote[0].DistinguishedName) > 0 {
			dataV1.DynamicRemote.DistinguishedName = &securityIkeGatewayBlockDynamicRemoteBlockDistinguishedName{
				Container: dataV0.DynamicRemote[0].DistinguishedName[0].Container,
				Wildcard:  dataV0.DynamicRemote[0].DistinguishedName[0].Wildcard,
			}
		}
	}
	if len(dataV0.Aaa) > 0 {
		dataV1.Aaa = &securityIkeGatewayBlockAaa{
			AccessProfile:  dataV0.Aaa[0].AccessProfile,
			ClientPassword: dataV0.Aaa[0].ClientPassword,
			ClientUsername: dataV0.Aaa[0].ClientUsername,
		}
	}
	if len(dataV0.DeadPeerDetection) > 0 {
		dataV1.DeadPeerDetection = &securityIkeGatewayBlockDeadPeerDetection{
			Interval:  dataV0.DeadPeerDetection[0].Interval,
			SendMode:  dataV0.DeadPeerDetection[0].SendMode,
			Threshold: dataV0.DeadPeerDetection[0].Threshold,
		}
	}
	if len(dataV0.LocalIdentity) > 0 {
		dataV1.LocalIdentity = &securityIkeGatewayBlockLocalIdentity{
			Type:  dataV0.LocalIdentity[0].Type,
			Value: dataV0.LocalIdentity[0].Value,
		}
	}
	if len(dataV0.RemoteIdentity) > 0 {
		dataV1.RemoteIdentity = &securityIkeGatewayBlockRemoteIdentity{
			Type:  dataV0.RemoteIdentity[0].Type,
			Value: dataV0.RemoteIdentity[0].Value,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
