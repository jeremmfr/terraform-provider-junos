package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *eventoptionsPolicy) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"events": schema.SetAttribute{
						ElementType: types.StringType,
					},
				},
				Blocks: map[string]schema.Block{
					"then": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ignore": schema.BoolAttribute{
									Optional: true,
								},
								"priority_override_facility": schema.StringAttribute{
									Optional: true,
								},
								"priority_override_severity": schema.StringAttribute{
									Optional: true,
								},
								"raise_trap": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"change_configuration": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"commands": schema.ListAttribute{
												ElementType: types.StringType,
												Required:    true,
											},
											"commit_options_check": schema.BoolAttribute{
												Optional: true,
											},
											"commit_options_check_synchronize": schema.BoolAttribute{
												Optional: true,
											},
											"commit_options_force": schema.BoolAttribute{
												Optional: true,
											},
											"commit_options_log": schema.StringAttribute{
												Optional: true,
											},
											"commit_options_synchronize": schema.BoolAttribute{
												Optional: true,
											},
											"retry_count": schema.Int64Attribute{
												Optional: true,
											},
											"retry_interval": schema.Int64Attribute{
												Optional: true,
											},
											"user_name": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"event_script": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"filename": schema.StringAttribute{
												Required: true,
											},
											"output_filename": schema.StringAttribute{
												Optional: true,
											},
											"output_format": schema.StringAttribute{
												Optional: true,
											},
											"user_name": schema.StringAttribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"arguments": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Required: true,
														},
														"value": schema.StringAttribute{
															Required: true,
														},
													},
												},
											},
											"destination": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Required: true,
														},
														"retry_count": schema.Int64Attribute{
															Optional: true,
														},
														"retry_interval": schema.Int64Attribute{
															Optional: true,
														},
														"transfer_delay": schema.Int64Attribute{
															Optional: true,
														},
													},
												},
											},
										},
									},
								},
								"execute_commands": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"commands": schema.ListAttribute{
												ElementType: types.StringType,
												Required:    true,
											},
											"output_filename": schema.StringAttribute{
												Optional: true,
											},
											"output_format": schema.StringAttribute{
												Optional: true,
											},
											"user_name": schema.StringAttribute{
												Optional: true,
											},
										},
										Blocks: map[string]schema.Block{
											"destination": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Required: true,
														},
														"retry_count": schema.Int64Attribute{
															Optional: true,
														},
														"retry_interval": schema.Int64Attribute{
															Optional: true,
														},
														"transfer_delay": schema.Int64Attribute{
															Optional: true,
														},
													},
												},
											},
										},
									},
								},
								"upload": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"filename": schema.StringAttribute{
												Required: true,
											},
											"destination": schema.StringAttribute{
												Required: true,
											},
											"retry_count": schema.Int64Attribute{
												Optional: true,
											},
											"retry_interval": schema.Int64Attribute{
												Optional: true,
											},
											"transfer_delay": schema.Int64Attribute{
												Optional: true,
											},
											"user_name": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
					"attributes_match": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"from": schema.StringAttribute{
									Required: true,
								},
								"compare": schema.StringAttribute{
									Required: true,
								},
								"to": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
					"within": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"time_interval": schema.Int64Attribute{
									Required: true,
								},
								"events": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"not_events": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
								"trigger_count": schema.Int64Attribute{
									Optional: true,
								},
								"trigger_when": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeEventoptionsPolicyV0toV1,
		},
	}
}

func upgradeEventoptionsPolicyV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID     types.String   `tfsdk:"id"`
		Name   types.String   `tfsdk:"name"`
		Events []types.String `tfsdk:"events"`
		Then   []struct {
			Ignore                   types.Bool   `tfsdk:"ignore"`
			PriorityOverrideFacility types.String `tfsdk:"priority_override_facility"`
			PriorityOverrideSeverity types.String `tfsdk:"priority_override_severity"`
			RaiseTrap                types.Bool   `tfsdk:"raise_trap"`
			ChangeConfiguration      []struct {
				Commands                      []types.String `tfsdk:"commands"`
				CommitOptionsCheck            types.Bool     `tfsdk:"commit_options_check"`
				CommitOptionsCheckSynchronize types.Bool     `tfsdk:"commit_options_check_synchronize"`
				CommitOptionsForce            types.Bool     `tfsdk:"commit_options_force"`
				CommitOptionsLog              types.String   `tfsdk:"commit_options_log"`
				CommitOptionsSynchronize      types.Bool     `tfsdk:"commit_options_synchronize"`
				RetryCount                    types.Int64    `tfsdk:"retry_count"`
				RetryInterval                 types.Int64    `tfsdk:"retry_interval"`
				Username                      types.String   `tfsdk:"user_name"`
			} `tfsdk:"change_configuration"`
			EventScript []struct {
				Filename       types.String `tfsdk:"filename"`
				OutputFilename types.String `tfsdk:"output_filename"`
				OutputFormat   types.String `tfsdk:"output_format"`
				Username       types.String `tfsdk:"user_name"`
				Arguments      []struct {
					Name  types.String `tfsdk:"name"`
					Value types.String `tfsdk:"value"`
				} `tfsdk:"arguments"`
				Destination []struct {
					Name          types.String `tfsdk:"name"`
					RetryCount    types.Int64  `tfsdk:"retry_count"`
					RetryInterval types.Int64  `tfsdk:"retry_interval"`
					TransferDelay types.Int64  `tfsdk:"transfer_delay"`
				} `tfsdk:"destination"`
			} `tfsdk:"event_script"`
			ExecuteCommands []struct {
				Commands       []types.String `tfsdk:"commands"`
				OutputFilename types.String   `tfsdk:"output_filename"`
				OutputFormat   types.String   `tfsdk:"output_format"`
				Username       types.String   `tfsdk:"user_name"`
				Destination    []struct {
					Name          types.String `tfsdk:"name"`
					RetryCount    types.Int64  `tfsdk:"retry_count"`
					RetryInterval types.Int64  `tfsdk:"retry_interval"`
					TransferDelay types.Int64  `tfsdk:"transfer_delay"`
				} `tfsdk:"destination"`
			} `tfsdk:"execute_commands"`
			Upload []struct {
				Filename      types.String `tfsdk:"filename"`
				Destination   types.String `tfsdk:"destination"`
				RetryCount    types.Int64  `tfsdk:"retry_count"`
				RetryInterval types.Int64  `tfsdk:"retry_interval"`
				TransferDelay types.Int64  `tfsdk:"transfer_delay"`
				Username      types.String `tfsdk:"user_name"`
			} `tfsdk:"upload"`
		} `tfsdk:"then"`
		AttributesMatch []struct {
			From    types.String `tfsdk:"from"`
			Compare types.String `tfsdk:"compare"`
			To      types.String `tfsdk:"to"`
		} `tfsdk:"attributes_match"`
		Within []struct {
			TimeInterval types.Int64    `tfsdk:"time_interval"`
			Events       []types.String `tfsdk:"events"`
			NotEvents    []types.String `tfsdk:"not_events"`
			TriggerCount types.Int64    `tfsdk:"trigger_count"`
			TriggerWhen  types.String   `tfsdk:"trigger_when"`
		} `tfsdk:"within"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 eventoptionsPolicyData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Events = dataV0.Events
	if len(dataV0.Then) > 0 {
		dataV1.Then = &eventoptionsPolicyBlockThen{
			Ignore:                   dataV0.Then[0].Ignore,
			RaiseTrap:                dataV0.Then[0].RaiseTrap,
			PriorityOverrideFacility: dataV0.Then[0].PriorityOverrideFacility,
			PriorityOverrideSeverity: dataV0.Then[0].PriorityOverrideSeverity,
		}
		if len(dataV0.Then[0].ChangeConfiguration) > 0 {
			dataV1.Then.ChangeConfiguration = &eventoptionsPolicyBlockThenBlockChangeConfigurtion{
				CommitOptionsCheck:            dataV0.Then[0].ChangeConfiguration[0].CommitOptionsCheck,
				CommitOptionsCheckSynchronize: dataV0.Then[0].ChangeConfiguration[0].CommitOptionsCheckSynchronize,
				CommitOptionsForce:            dataV0.Then[0].ChangeConfiguration[0].CommitOptionsForce,
				CommitOptionsSynchronize:      dataV0.Then[0].ChangeConfiguration[0].CommitOptionsSynchronize,
				Commands:                      dataV0.Then[0].ChangeConfiguration[0].Commands,
				CommitOptionsLog:              dataV0.Then[0].ChangeConfiguration[0].CommitOptionsLog,
				RetryCount:                    dataV0.Then[0].ChangeConfiguration[0].RetryCount,
				RetryInterval:                 dataV0.Then[0].ChangeConfiguration[0].RetryInterval,
				Username:                      dataV0.Then[0].ChangeConfiguration[0].Username,
			}
		}
		for _, blockV0 := range dataV0.Then[0].EventScript {
			blockV1 := eventoptionsPolicyBlockThenBlockEventScript{
				Filename:       blockV0.Filename,
				OutputFilename: blockV0.OutputFilename,
				OutputFormat:   blockV0.OutputFormat,
				Username:       blockV0.Username,
			}
			for _, subBlockV0 := range blockV0.Arguments {
				blockV1.Arguments = append(blockV1.Arguments, eventoptionsPolicyBlockThenBlockEventScriptBlockArguments{
					Name:  subBlockV0.Name,
					Value: subBlockV0.Value,
				})
			}
			if len(blockV0.Destination) > 0 {
				blockV1.Destination = &eventoptionsPolicyBlockThenBlockDestination{
					Name:          blockV0.Destination[0].Name,
					RetryCount:    blockV0.Destination[0].RetryCount,
					RetryInterval: blockV0.Destination[0].RetryInterval,
					TransferDelay: blockV0.Destination[0].TransferDelay,
				}
			}
			dataV1.Then.EventScript = append(dataV1.Then.EventScript, blockV1)
		}
		if len(dataV0.Then[0].ExecuteCommands) > 0 {
			dataV1.Then.ExecuteCommands = &eventoptionsPolicyBlockThenBlockExecuteCommands{
				Commands:       dataV0.Then[0].ExecuteCommands[0].Commands,
				OutputFilename: dataV0.Then[0].ExecuteCommands[0].OutputFilename,
				OutputFormat:   dataV0.Then[0].ExecuteCommands[0].OutputFormat,
				Username:       dataV0.Then[0].ExecuteCommands[0].Username,
			}
			if len(dataV0.Then[0].ExecuteCommands[0].Destination) > 0 {
				dataV1.Then.ExecuteCommands.Destination = &eventoptionsPolicyBlockThenBlockDestination{
					Name:          dataV0.Then[0].ExecuteCommands[0].Destination[0].Name,
					RetryCount:    dataV0.Then[0].ExecuteCommands[0].Destination[0].RetryCount,
					RetryInterval: dataV0.Then[0].ExecuteCommands[0].Destination[0].RetryInterval,
					TransferDelay: dataV0.Then[0].ExecuteCommands[0].Destination[0].TransferDelay,
				}
			}
		}
		for _, blockV0 := range dataV0.Then[0].Upload {
			dataV1.Then.Upload = append(dataV1.Then.Upload, eventoptionsPolicyBlockThenBlockUpload{
				Filename:      blockV0.Filename,
				Destination:   blockV0.Destination,
				RetryCount:    blockV0.RetryCount,
				RetryInterval: blockV0.RetryInterval,
				TransferDelay: blockV0.TransferDelay,
				Username:      blockV0.Username,
			})
		}
	}
	for _, blockV0 := range dataV0.AttributesMatch {
		dataV1.AttributesMatch = append(dataV1.AttributesMatch, eventoptionsPolicyBlockAttributesMatch{
			From:    blockV0.From,
			Compare: blockV0.Compare,
			To:      blockV0.To,
		})
	}
	for _, blockV0 := range dataV0.Within {
		dataV1.Within = append(dataV1.Within, eventoptionsPolicyBlockWithin{
			TimeInterval: blockV0.TimeInterval,
			Events:       blockV0.Events,
			NotEvents:    blockV0.NotEvents,
			TriggerCount: blockV0.TriggerCount,
			TriggerWhen:  blockV0.TriggerWhen,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
