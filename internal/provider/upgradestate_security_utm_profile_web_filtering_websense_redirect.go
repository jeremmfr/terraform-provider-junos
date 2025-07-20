package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) UpgradeState(
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
					"account": schema.StringAttribute{
						Optional: true,
					},
					"custom_block_message": schema.StringAttribute{
						Optional: true,
					},
					"sockets": schema.Int64Attribute{
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
					"server": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									Optional: true,
								},
								"port": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSecurityUtmProfileWebFilteringWebsenseRedirectV0toV1,
		},
	}
}

func upgradeSecurityUtmProfileWebFilteringWebsenseRedirectV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                 types.String `tfsdk:"id"`
		Name               types.String `tfsdk:"name"`
		Account            types.String `tfsdk:"account"`
		CustomBlockMessage types.String `tfsdk:"custom_block_message"`
		Sockets            types.Int64  `tfsdk:"sockets"`
		Timeout            types.Int64  `tfsdk:"timeout"`
		FallbackSettings   []struct {
			Default            types.String `tfsdk:"default"`
			ServerConnectivity types.String `tfsdk:"server_connectivity"`
			Timeout            types.String `tfsdk:"timeout"`
			TooManyRequests    types.String `tfsdk:"too_many_requests"`
		} `tfsdk:"fallback_settings"`
		Server []struct {
			Host types.String `tfsdk:"host"`
			Port types.Int64  `tfsdk:"port"`
		} `tfsdk:"server"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityUtmProfileWebFilteringWebsenseRedirectData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Account = dataV0.Account
	dataV1.CustomBlockMessage = dataV0.CustomBlockMessage
	dataV1.Sockets = dataV0.Sockets
	dataV1.Timeout = dataV0.Timeout
	if len(dataV0.FallbackSettings) > 0 {
		dataV1.FallbackSettings = &securityUtmProfileWebFilteringBlockFallbackSettings{
			Default:            dataV0.FallbackSettings[0].Default,
			ServerConnectivity: dataV0.FallbackSettings[0].ServerConnectivity,
			Timeout:            dataV0.FallbackSettings[0].Timeout,
			TooManyRequests:    dataV0.FallbackSettings[0].TooManyRequests,
		}
	}
	if len(dataV0.Server) > 0 {
		dataV1.Server = &securityUtmProfileWebFilteringWebsenseRedirectBlockServer{
			Host: dataV0.Server[0].Host,
			Port: dataV0.Server[0].Port,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
