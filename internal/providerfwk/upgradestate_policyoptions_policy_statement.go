package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *policyoptionsPolicyStatement) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"add_it_to_forwarding_table_export": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"from": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: rsc.schemaFromAttributes(),
							Blocks:     rsc.schemaFromBlocks(),
						},
					},
					"to": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: rsc.schemaToAttributes(),
						},
					},
					"then": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: rsc.schemaThenAttributes(),
							Blocks: map[string]schema.Block{
								"community": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.StringAttribute{
												Required: true,
											},
											"value": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
								"local_preference": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.StringAttribute{
												Required: true,
											},
											"value": schema.Int64Attribute{
												Required: true,
											},
										},
									},
								},
								"metric": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.StringAttribute{
												Required: true,
											},
											"value": schema.Int64Attribute{
												Required: true,
											},
										},
									},
								},
								"preference": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.StringAttribute{
												Required: true,
											},
											"value": schema.Int64Attribute{
												Required: true,
											},
										},
									},
								},
							},
						},
					},
					"term": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
							},
							Blocks: map[string]schema.Block{
								"from": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: rsc.schemaFromAttributes(),
										Blocks:     rsc.schemaFromBlocks(),
									},
								},
								"to": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: rsc.schemaToAttributes(),
									},
								},
								"then": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: rsc.schemaThenAttributes(),
										Blocks: map[string]schema.Block{
											"community": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"action": schema.StringAttribute{
															Required: true,
														},
														"value": schema.StringAttribute{
															Required: true,
														},
													},
												},
											},
											"local_preference": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"action": schema.StringAttribute{
															Required: true,
														},
														"value": schema.Int64Attribute{
															Required: true,
														},
													},
												},
											},
											"metric": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"action": schema.StringAttribute{
															Required: true,
														},
														"value": schema.Int64Attribute{
															Required: true,
														},
													},
												},
											},
											"preference": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"action": schema.StringAttribute{
															Required: true,
														},
														"value": schema.Int64Attribute{
															Required: true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradePolicyoptionsPolicyStatementStateV0toV1,
		},
	}
}

func upgradePolicyoptionsPolicyStatementStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		AddItToForwardingTableExport types.Bool                              `tfsdk:"add_it_to_forwarding_table_export"`
		ID                           types.String                            `tfsdk:"id"`
		Name                         types.String                            `tfsdk:"name"`
		From                         []policyoptionsPolicyStatementBlockFrom `tfsdk:"from"`
		To                           []policyoptionsPolicyStatementBlockTo   `tfsdk:"to"`
		Then                         []struct {
			Action          types.String                                                 `tfsdk:"action"`
			ASPathExpand    types.String                                                 `tfsdk:"as_path_expand"`
			ASPathPrepend   types.String                                                 `tfsdk:"as_path_prepend"`
			DefaultAction   types.String                                                 `tfsdk:"default_action"`
			LoadBalance     types.String                                                 `tfsdk:"load_balance"`
			Next            types.String                                                 `tfsdk:"next"`
			NextHop         types.String                                                 `tfsdk:"next_hop"`
			Origin          types.String                                                 `tfsdk:"origin"`
			Community       []policyoptionsPolicyStatementBlockThenBlockActionValue      `tfsdk:"community"`
			LocalPreference []policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"local_preference"`
			Metric          []policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"metric"`
			Preference      []policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"preference"`
		} `tfsdk:"then"`
		Term []struct {
			Name types.String                            `tfsdk:"name"`
			From []policyoptionsPolicyStatementBlockFrom `tfsdk:"from"`
			To   []policyoptionsPolicyStatementBlockTo   `tfsdk:"to"`
			Then []struct {
				Action          types.String                                                 `tfsdk:"action"`
				ASPathExpand    types.String                                                 `tfsdk:"as_path_expand"`
				ASPathPrepend   types.String                                                 `tfsdk:"as_path_prepend"`
				DefaultAction   types.String                                                 `tfsdk:"default_action"`
				LoadBalance     types.String                                                 `tfsdk:"load_balance"`
				Next            types.String                                                 `tfsdk:"next"`
				NextHop         types.String                                                 `tfsdk:"next_hop"`
				Origin          types.String                                                 `tfsdk:"origin"`
				Community       []policyoptionsPolicyStatementBlockThenBlockActionValue      `tfsdk:"community"`
				LocalPreference []policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"local_preference"`
				Metric          []policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"metric"`
				Preference      []policyoptionsPolicyStatementBlockThenBlockActionValueInt64 `tfsdk:"preference"`
			} `tfsdk:"then"`
		} `tfsdk:"term"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 policyoptionsPolicyStatementData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.AddItToForwardingTableExport = dataV0.AddItToForwardingTableExport
	if !dataV1.AddItToForwardingTableExport.IsNull() && !dataV1.AddItToForwardingTableExport.ValueBool() {
		dataV1.AddItToForwardingTableExport = types.BoolNull()
	}
	if len(dataV0.From) > 0 {
		dataV1.From = &dataV0.From[0]
	}
	if len(dataV0.To) > 0 {
		dataV1.To = &dataV0.To[0]
	}
	if len(dataV0.Then) > 0 {
		dataV1.Then = &policyoptionsPolicyStatementBlockThen{
			Action:        dataV0.Then[0].Action,
			ASPathExpand:  dataV0.Then[0].ASPathExpand,
			ASPathPrepend: dataV0.Then[0].ASPathPrepend,
			DefaultAction: dataV0.Then[0].DefaultAction,
			LoadBalance:   dataV0.Then[0].LoadBalance,
			Next:          dataV0.Then[0].Next,
			NextHop:       dataV0.Then[0].NextHop,
			Origin:        dataV0.Then[0].Origin,
			Community:     dataV0.Then[0].Community,
		}
		if len(dataV0.Then[0].LocalPreference) > 0 {
			dataV1.Then.LocalPreference = &dataV0.Then[0].LocalPreference[0]
		}
		if len(dataV0.Then[0].Metric) > 0 {
			dataV1.Then.Metric = &dataV0.Then[0].Metric[0]
		}
		if len(dataV0.Then[0].Preference) > 0 {
			dataV1.Then.Preference = &dataV0.Then[0].Preference[0]
		}
	}
	for _, blockV0 := range dataV0.Term {
		blockV1 := policyoptionsPolicyStatementBlockTerm{
			Name: blockV0.Name,
		}
		if len(blockV0.From) > 0 {
			blockV1.From = &blockV0.From[0]
		}
		if len(blockV0.To) > 0 {
			blockV1.To = &blockV0.To[0]
		}
		if len(dataV0.Then) > 0 {
			blockV1.Then = &policyoptionsPolicyStatementBlockThen{
				Action:        blockV0.Then[0].Action,
				ASPathExpand:  blockV0.Then[0].ASPathExpand,
				ASPathPrepend: blockV0.Then[0].ASPathPrepend,
				DefaultAction: blockV0.Then[0].DefaultAction,
				LoadBalance:   blockV0.Then[0].LoadBalance,
				Next:          blockV0.Then[0].Next,
				NextHop:       blockV0.Then[0].NextHop,
				Origin:        blockV0.Then[0].Origin,
				Community:     blockV0.Then[0].Community,
			}
			if len(blockV0.Then[0].LocalPreference) > 0 {
				blockV1.Then.LocalPreference = &blockV0.Then[0].LocalPreference[0]
			}
			if len(blockV0.Then[0].Metric) > 0 {
				blockV1.Then.Metric = &blockV0.Then[0].Metric[0]
			}
			if len(dataV0.Then[0].Preference) > 0 {
				blockV1.Then.Preference = &blockV0.Then[0].Preference[0]
			}
		}
		dataV1.Term = append(dataV1.Term, blockV1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
