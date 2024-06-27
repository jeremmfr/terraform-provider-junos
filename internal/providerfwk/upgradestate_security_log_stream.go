package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityLogStream) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"category": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"filter_threat_attack": schema.BoolAttribute{
						Optional: true,
					},
					"format": schema.StringAttribute{
						Optional: true,
					},
					"rate_limit": schema.Int64Attribute{
						Optional: true,
					},
					"severity": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"file": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"allow_duplicates": schema.BoolAttribute{
									Optional: true,
								},
								"rotation": schema.Int64Attribute{
									Optional: true,
								},
								"size": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"host": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ip_address": schema.StringAttribute{
									Required: true,
								},
								"port": schema.Int64Attribute{
									Optional: true,
								},
								"routing_instance": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSecurityLogStreamV0toV1,
		},
	}
}

func upgradeSecurityLogStreamV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                 types.String   `tfsdk:"id"`
		Name               types.String   `tfsdk:"name"`
		Category           []types.String `tfsdk:"category"`
		FilterThreatAttack types.Bool     `tfsdk:"filter_threat_attack"`
		Format             types.String   `tfsdk:"format"`
		RateLimit          types.Int64    `tfsdk:"rate_limit"`
		Severity           types.String   `tfsdk:"severity"`
		File               []struct {
			Name            types.String `tfsdk:"name"`
			AllowDuplicates types.Bool   `tfsdk:"allow_duplicates"`
			Rotation        types.Int64  `tfsdk:"rotation"`
			Size            types.Int64  `tfsdk:"size"`
		} `tfsdk:"file"`
		Host []struct {
			IPAddress       types.String `tfsdk:"ip_address"`
			Port            types.Int64  `tfsdk:"port"`
			RoutingInstance types.String `tfsdk:"routing_instance"`
		} `tfsdk:"host"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityLogStreamData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Category = dataV0.Category
	dataV1.FilterThreatAttack = dataV0.FilterThreatAttack
	dataV1.Format = dataV0.Format
	dataV1.RateLimit = dataV0.RateLimit
	dataV1.Severity = dataV0.Severity
	if len(dataV0.File) > 0 {
		dataV1.File = &securityLogStreamBlockFile{
			Name:            dataV0.File[0].Name,
			AllowDuplicates: dataV0.File[0].AllowDuplicates,
			Rotation:        dataV0.File[0].Rotation,
			Size:            dataV0.File[0].Size,
		}
	}

	if len(dataV0.Host) > 0 {
		dataV1.Host = &securityLogStreamBlockHost{
			IPAddress:       dataV0.Host[0].IPAddress,
			Port:            dataV0.Host[0].Port,
			RoutingInstance: dataV0.Host[0].RoutingInstance,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
