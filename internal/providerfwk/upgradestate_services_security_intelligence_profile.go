package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *servicesSecurityIntelligenceProfile) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"category": schema.StringAttribute{
						Required: true,
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"rule": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"then_action": schema.StringAttribute{
									Required: true,
								},
								"then_log": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"match": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"threat_level": schema.ListAttribute{
												ElementType: types.Int64Type,
												Required:    true,
											},
											"feed_name": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
										},
									},
								},
							},
						},
					},
					"default_rule_then": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"action": schema.StringAttribute{
									Required: true,
								},
								"log": schema.BoolAttribute{
									Optional: true,
								},
								"no_log": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeServicesSecurityIntelligenceProfileV0toV1,
		},
	}
}

func upgradeServicesSecurityIntelligenceProfileV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID          types.String `tfsdk:"id"`
		Name        types.String `tfsdk:"name"`
		Category    types.String `tfsdk:"category"`
		Description types.String `tfsdk:"description"`
		Rule        []struct {
			Name       types.String `tfsdk:"name"`
			ThenAction types.String `tfsdk:"then_action"`
			ThenLog    types.Bool   `tfsdk:"then_log"`
			Match      []struct {
				ThreatLevel []types.Int64  `tfsdk:"threat_level"`
				FeedName    []types.String `tfsdk:"feed_name"`
			} `tfsdk:"match"`
		} `tfsdk:"rule"`
		DefaultRuleThen []struct {
			Action types.String `tfsdk:"action"`
			Log    types.Bool   `tfsdk:"log"`
			NoLog  types.Bool   `tfsdk:"no_log"`
		} `tfsdk:"default_rule_then"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 servicesSecurityIntelligenceProfileData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Category = dataV0.Category
	dataV1.Description = dataV0.Description
	for _, blockV0 := range dataV0.Rule {
		blockV1 := servicesSecurityIntelligenceProfileBlockRule{
			Name:       blockV0.Name,
			ThenAction: blockV0.ThenAction,
			ThenLog:    blockV0.ThenLog,
		}
		if len(blockV0.Match) > 0 {
			blockV1.Match = &servicesSecurityIntelligenceProfileBlockRuleBlockMatch{
				ThreatLevel: blockV0.Match[0].ThreatLevel,
				FeedName:    blockV0.Match[0].FeedName,
			}
		}
		dataV1.Rule = append(dataV1.Rule, blockV1)
	}
	if len(dataV0.DefaultRuleThen) > 0 {
		dataV1.DefaultRuleThen = &servicesSecurityIntelligenceProfileBlockDefaultRuleThen{
			Action: dataV0.DefaultRuleThen[0].Action,
			Log:    dataV0.DefaultRuleThen[0].Log,
			NoLog:  dataV0.DefaultRuleThen[0].NoLog,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
