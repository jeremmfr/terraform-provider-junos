package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityUtmProfileWebFilteringJuniperLocal) UpgradeState(
	_ context.Context,
) map[int64]resource.StateUpgrader {
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
					"custom_block_message": schema.StringAttribute{
						Optional: true,
					},
					"default_action": schema.StringAttribute{
						Optional: true,
					},
					"timeout": schema.Int64Attribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"fallback_settings": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"default": schema.StringAttribute{
									Optional: true,
								},
								"server_connectivity": schema.StringAttribute{
									Optional: true,
								},
								"timeout": schema.StringAttribute{
									Optional: true,
								},
								"too_many_requests": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSecurityUtmProfileWebFilteringJuniperLocalV0toV1,
		},
	}
}

func upgradeSecurityUtmProfileWebFilteringJuniperLocalV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                 types.String `tfsdk:"id"`
		Name               types.String `tfsdk:"name"`
		CustomBlockMessage types.String `tfsdk:"custom_block_message"`
		DefaultAction      types.String `tfsdk:"default_action"`
		Timeout            types.Int64  `tfsdk:"timeout"`
		FallbackSettings   []struct {
			Default            types.String `tfsdk:"default"`
			ServerConnectivity types.String `tfsdk:"server_connectivity"`
			Timeout            types.String `tfsdk:"timeout"`
			TooManyRequests    types.String `tfsdk:"too_many_requests"`
		} `tfsdk:"fallback_settings"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityUtmProfileWebFilteringJuniperLocalData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.CustomBlockMessage = dataV0.CustomBlockMessage
	dataV1.DefaultAction = dataV0.DefaultAction
	dataV1.Timeout = dataV0.Timeout
	if len(dataV0.FallbackSettings) > 0 {
		dataV1.FallbackSettings = &securityUtmProfileWebFilteringBlockFallbackSettings{
			Default:            dataV0.FallbackSettings[0].Default,
			ServerConnectivity: dataV0.FallbackSettings[0].ServerConnectivity,
			Timeout:            dataV0.FallbackSettings[0].Timeout,
			TooManyRequests:    dataV0.FallbackSettings[0].TooManyRequests,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
