package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityDynamicAddressName) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"description": schema.StringAttribute{
						Optional: true,
					},
					"profile_feed_name": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"profile_category": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"feed": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"property": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Required: true,
											},
											"string": schema.ListAttribute{
												ElementType: types.StringType,
												Required:    true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSecurityDynamicAddressNameV0toV1,
		},
	}
}

func upgradeSecurityDynamicAddressNameV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID              types.String `tfsdk:"id"`
		Name            types.String `tfsdk:"name"`
		Description     types.String `tfsdk:"description"`
		ProfileFeedName types.String `tfsdk:"profile_feed_name"`
		ProfileCategory []struct {
			Name     types.String `tfsdk:"name"`
			Feed     types.String `tfsdk:"feed"`
			Property []struct {
				Name   types.String   `tfsdk:"name"`
				String []types.String `tfsdk:"string"`
			} `tfsdk:"property"`
		} `tfsdk:"profile_category"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityDynamicAddressNameData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Description = dataV0.Description
	dataV1.ProfileFeedName = dataV0.ProfileFeedName
	if len(dataV0.ProfileCategory) > 0 {
		dataV1.ProfileCategory = &securityDynamicAddressNameBlockProfileCategory{
			Name: dataV0.ProfileCategory[0].Name,
			Feed: dataV0.ProfileCategory[0].Feed,
		}
		for _, blockV0 := range dataV0.ProfileCategory[0].Property {
			dataV1.ProfileCategory.Property = append(dataV1.ProfileCategory.Property,
				securityDynamicAddressNameBlockProfileCategoryBlockProperty{
					Name:   blockV0.Name,
					String: blockV0.String,
				},
			)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
