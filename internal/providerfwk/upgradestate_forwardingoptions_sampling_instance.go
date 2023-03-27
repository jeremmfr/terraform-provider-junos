package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *forwardingoptionsSamplingInstance) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"disable": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"family_inet_input": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: rsc.schemaInputAttributes(),
						},
					},
					"family_inet_output": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: rsc.schemaFamilyInetOutputAttributes(),
							Blocks: map[string]schema.Block{
								"flow_server": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"hostname": schema.StringAttribute{
												Required: true,
											},
											"port": schema.Int64Attribute{
												Required: true,
											},
											"aggregation_autonomous_system": schema.BoolAttribute{
												Optional: true,
											},
											"aggregation_destination_prefix": schema.BoolAttribute{
												Optional: true,
											},
											"aggregation_protocol_port": schema.BoolAttribute{
												Optional: true,
											},
											"aggregation_source_destination_prefix": schema.BoolAttribute{
												Optional: true,
											},
											"aggregation_source_destination_prefix_caida_compliant": schema.BoolAttribute{
												Optional: true,
											},
											"aggregation_source_prefix": schema.BoolAttribute{
												Optional: true,
											},
											"autonomous_system_type": schema.StringAttribute{
												Optional: true,
											},
											"dscp": schema.Int64Attribute{
												Optional: true,
											},
											"forwarding_class": schema.StringAttribute{
												Optional: true,
											},
											"local_dump": schema.BoolAttribute{
												Optional: true,
											},
											"no_local_dump": schema.BoolAttribute{
												Optional: true,
											},
											"routing_instance": schema.StringAttribute{
												Optional: true,
											},
											"source_address": schema.StringAttribute{
												Optional: true,
											},
											"version": schema.Int64Attribute{
												Optional: true,
											},
											"version9_template": schema.StringAttribute{
												Optional: true,
											},
											"version_ipfix_template": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
								"interface": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: rsc.schemaOutputInterfaceAttributes(),
									},
								},
							},
						},
					},
					"family_inet6_input": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: rsc.schemaInputAttributes(),
						},
					},
					"family_inet6_output": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: rsc.schemaFamilyInetOutputAttributes(),
							Blocks:     rsc.schemaOutputBlock(),
						},
					},
					"family_mpls_input": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: rsc.schemaInputAttributes(),
						},
					},
					"family_mpls_output": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"aggregate_export_interval": schema.Int64Attribute{
									Optional: true,
								},
								"flow_active_timeout": schema.Int64Attribute{
									Optional: true,
								},
								"flow_inactive_timeout": schema.Int64Attribute{
									Optional: true,
								},
								"inline_jflow_export_rate": schema.Int64Attribute{
									Optional: true,
								},
								"inline_jflow_source_address": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: rsc.schemaOutputBlock(),
						},
					},
					"input": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: rsc.schemaInputAttributes(),
						},
					},
				},
			},
			StateUpgrader: upgradeForwardingoptionsSamplingInstanceStateV0toV1,
		},
	}
}

func upgradeForwardingoptionsSamplingInstanceStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		Disable           types.Bool                                                `tfsdk:"disable"`
		ID                types.String                                              `tfsdk:"id"`
		Name              types.String                                              `tfsdk:"name"`
		FamilyInetInput   []forwardingoptionsSamplingInstanceBlockInput             `tfsdk:"family_inet_input"`
		FamilyInetOutput  []forwardingoptionsSamplingInstanceBlockFamilyInetOutput  `tfsdk:"family_inet_output"`
		FamilyInet6Input  []forwardingoptionsSamplingInstanceBlockInput             `tfsdk:"family_inet6_input"`
		FamilyInet6Output []forwardingoptionsSamplingInstanceBlockFamilyInet6Output `tfsdk:"family_inet6_output"`
		FamilyMplsInput   []forwardingoptionsSamplingInstanceBlockInput             `tfsdk:"family_mpls_input"`
		FamilyMplsOutput  []forwardingoptionsSamplingInstanceBlockFamilyMplsOutput  `tfsdk:"family_mpls_output"`
		Input             []forwardingoptionsSamplingInstanceBlockInput             `tfsdk:"input"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 forwardingoptionsSamplingInstanceData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Disable = dataV0.Disable
	if len(dataV0.FamilyInetInput) > 0 {
		dataV1.FamilyInetInput = &dataV0.FamilyInetInput[0]
	}
	if len(dataV0.FamilyInetOutput) > 0 {
		dataV1.FamilyInetOutput = &dataV0.FamilyInetOutput[0]
	}
	if len(dataV0.FamilyInet6Input) > 0 {
		dataV1.FamilyInet6Input = &dataV0.FamilyInet6Input[0]
	}
	if len(dataV0.FamilyInet6Output) > 0 {
		dataV1.FamilyInet6Output = &dataV0.FamilyInet6Output[0]
	}
	if len(dataV0.FamilyMplsInput) > 0 {
		dataV1.FamilyMplsInput = &dataV0.FamilyMplsInput[0]
	}
	if len(dataV0.FamilyMplsOutput) > 0 {
		dataV1.FamilyMplsOutput = &dataV0.FamilyMplsOutput[0]
	}
	if len(dataV0.Input) > 0 {
		dataV1.Input = &dataV0.Input[0]
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
