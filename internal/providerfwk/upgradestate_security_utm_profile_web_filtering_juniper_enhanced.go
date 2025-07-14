package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) UpgradeState(
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
					"no_safe_search": schema.BoolAttribute{
						Optional: true,
					},
					"quarantine_custom_message": schema.StringAttribute{
						Optional: true,
					},
					"timeout": schema.Int64Attribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"block_message": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type_custom_redirect_url": schema.BoolAttribute{
									Optional: true,
								},
								"url": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"category": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"action": schema.StringAttribute{
									Required: true,
								},
							},
							Blocks: map[string]schema.Block{
								"reputation_action": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"site_reputation": schema.StringAttribute{
												Required: true,
											},
											"action": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
							},
						},
					},
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
					"quarantine_message": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"url": schema.StringAttribute{
									Optional: true,
								},
								"type_custom_redirect_url": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
					"site_reputation_action": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"site_reputation": schema.StringAttribute{
									Required: true,
								},
								"action": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSecurityUtmProfileWebFilteringJuniperEnhancedV0toV1,
		},
	}
}

func upgradeSecurityUtmProfileWebFilteringJuniperEnhancedV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                      types.String `tfsdk:"id"`
		Name                    types.String `tfsdk:"name"`
		CustomBlockMessage      types.String `tfsdk:"custom_block_message"`
		DefaultAction           types.String `tfsdk:"default_action"`
		NoSafeSearch            types.Bool   `tfsdk:"no_safe_search"`
		QuarantineCustomMessage types.String `tfsdk:"quarantine_custom_message"`
		Timeout                 types.Int64  `tfsdk:"timeout"`
		BlockMessage            []struct {
			TypeCustomRedirectURL types.Bool   `tfsdk:"type_custom_redirect_url"`
			URL                   types.String `tfsdk:"url"`
		} `tfsdk:"block_message"`
		Category []struct {
			Name             types.String `tfsdk:"name"   tfdata:"identifier"`
			Action           types.String `tfsdk:"action"`
			ReputationAction []struct {
				SiteReputation types.String `tfsdk:"site_reputation"`
				Action         types.String `tfsdk:"action"`
			} `tfsdk:"reputation_action"`
		} `tfsdk:"category"`
		FallbackSettings []struct {
			Default            types.String `tfsdk:"default"`
			ServerConnectivity types.String `tfsdk:"server_connectivity"`
			Timeout            types.String `tfsdk:"timeout"`
			TooManyRequests    types.String `tfsdk:"too_many_requests"`
		} `tfsdk:"fallback_settings"`
		QuarantineMessage []struct {
			TypeCustomRedirectURL types.Bool   `tfsdk:"type_custom_redirect_url"`
			URL                   types.String `tfsdk:"url"`
		} `tfsdk:"quarantine_message"`
		SiteReputationAction []struct {
			SiteReputation types.String `tfsdk:"site_reputation"`
			Action         types.String `tfsdk:"action"`
		} `tfsdk:"site_reputation_action"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityUtmProfileWebFilteringJuniperEnhancedData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.CustomBlockMessage = dataV0.CustomBlockMessage
	dataV1.DefaultAction = dataV0.DefaultAction
	dataV1.NoSafeSearch = dataV0.NoSafeSearch
	dataV1.QuarantineCustomMessage = dataV0.QuarantineCustomMessage
	dataV1.Timeout = dataV0.Timeout
	if len(dataV0.BlockMessage) > 0 {
		dataV1.BlockMessage = &securityUtmProfileWebFilteringJuniperEnhancedBlockMessage{
			TypeCustomRedirectURL: dataV0.BlockMessage[0].TypeCustomRedirectURL,
			URL:                   dataV0.BlockMessage[0].URL,
		}
	}
	for _, blockV0 := range dataV0.Category {
		blockV1 := securityUtmProfileWebFilteringJuniperEnhancedBlockCategory{
			Name:   blockV0.Name,
			Action: blockV0.Action,
		}
		for _, subBlockV0 := range blockV0.ReputationAction {
			blockV1.ReputationAction = append(blockV1.ReputationAction,
				securityUtmProfileWebFilteringJuniperEnhancedBlockReputationAction{
					SiteReputation: subBlockV0.SiteReputation,
					Action:         subBlockV0.Action,
				},
			)
		}
		dataV1.Category = append(dataV1.Category, blockV1)
	}
	if len(dataV0.FallbackSettings) > 0 {
		dataV1.FallbackSettings = &securityUtmProfileWebFilteringBlockFallbackSettings{
			Default:            dataV0.FallbackSettings[0].Default,
			ServerConnectivity: dataV0.FallbackSettings[0].ServerConnectivity,
			Timeout:            dataV0.FallbackSettings[0].Timeout,
			TooManyRequests:    dataV0.FallbackSettings[0].TooManyRequests,
		}
	}
	if len(dataV0.QuarantineMessage) > 0 {
		dataV1.QuarantineMessage = &securityUtmProfileWebFilteringJuniperEnhancedBlockMessage{
			TypeCustomRedirectURL: dataV0.QuarantineMessage[0].TypeCustomRedirectURL,
			URL:                   dataV0.QuarantineMessage[0].URL,
		}
	}
	for _, blockV0 := range dataV0.SiteReputationAction {
		dataV1.SiteReputationAction = append(dataV1.SiteReputationAction,
			securityUtmProfileWebFilteringJuniperEnhancedBlockReputationAction{
				SiteReputation: blockV0.SiteReputation,
				Action:         blockV0.Action,
			},
		)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
