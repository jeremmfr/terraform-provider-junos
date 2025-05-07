package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *servicesUserIdentificationADAccessDomain) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"user_name": schema.StringAttribute{
						Required: true,
					},
					"user_password": schema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
				},
				Blocks: map[string]schema.Block{
					"domain_controller": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"address": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
					"ip_user_mapping_discovery_wmi": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"event_log_scanning_interval": schema.Int64Attribute{
									Optional: true,
								},
								"initial_event_log_timespan": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"user_group_mapping_ldap": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"base": schema.StringAttribute{
									Required: true,
								},
								"address": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"auth_algo_simple": schema.BoolAttribute{
									Optional: true,
								},
								"ssl": schema.BoolAttribute{
									Optional: true,
								},
								"user_name": schema.StringAttribute{
									Optional: true,
								},
								"user_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeServicesUserIdentificationADAccessDomainV0toV1,
		},
	}
}

func upgradeServicesUserIdentificationADAccessDomainV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID               types.String `tfsdk:"id"`
		Name             types.String `tfsdk:"name"`
		UserName         types.String `tfsdk:"user_name"`
		UserPassword     types.String `tfsdk:"user_password"`
		DomainController []struct {
			Name    types.String `tfsdk:"name"`
			Address types.String `tfsdk:"address"`
		} `tfsdk:"domain_controller"`
		IPUserMappingDiscoveryWmi []struct {
			EventLogScanningInterval types.Int64 `tfsdk:"event_log_scanning_interval"`
			InitialEventLogTimespan  types.Int64 `tfsdk:"initial_event_log_timespan"`
		} `tfsdk:"ip_user_mapping_discovery_wmi"`
		UserGroupMappingLdap []struct {
			Base           types.String   `tfsdk:"base"`
			Address        []types.String `tfsdk:"address"`
			AuthAlgoSimple types.Bool     `tfsdk:"auth_algo_simple"`
			Ssl            types.Bool     `tfsdk:"ssl"`
			UserName       types.String   `tfsdk:"user_name"`
			UserPassword   types.String   `tfsdk:"user_password"`
		} `tfsdk:"user_group_mapping_ldap"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 servicesUserIdentificationADAccessDomainData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.UserName = dataV0.UserName
	dataV1.UserPassword = dataV0.UserPassword
	for _, blockV0 := range dataV0.DomainController {
		dataV1.DomainController = append(dataV1.DomainController,
			servicesUserIdentificationADAccessDomainBlockDomainController{
				Name:    blockV0.Name,
				Address: blockV0.Address,
			},
		)
	}
	if len(dataV0.IPUserMappingDiscoveryWmi) > 0 {
		dataV1.IPUserMappingDiscoveryWmi = &servicesUserIdentificationADAccessDomainBlockIPUserMappingDiscoveryWmi{
			EventLogScanningInterval: dataV0.IPUserMappingDiscoveryWmi[0].EventLogScanningInterval,
			InitialEventLogTimespan:  dataV0.IPUserMappingDiscoveryWmi[0].InitialEventLogTimespan,
		}
	}
	if len(dataV0.UserGroupMappingLdap) > 0 {
		dataV1.UserGroupMappingLdap = &servicesUserIdentificationADAccessDomainBlockUserGroupMappingLdap{
			Base:           dataV0.UserGroupMappingLdap[0].Base,
			Address:        dataV0.UserGroupMappingLdap[0].Address,
			AuthAlgoSimple: dataV0.UserGroupMappingLdap[0].AuthAlgoSimple,
			Ssl:            dataV0.UserGroupMappingLdap[0].Ssl,
			UserName:       dataV0.UserGroupMappingLdap[0].UserName,
			UserPassword:   dataV0.UserGroupMappingLdap[0].UserPassword,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
