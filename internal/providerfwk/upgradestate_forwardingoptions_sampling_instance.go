package providerfwk

import (
	"context"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

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
							Attributes: map[string]schema.Attribute{
								"max_packets_per_second": schema.Int64Attribute{
									Optional: true,
								},
								"maximum_packet_length": schema.Int64Attribute{
									Optional: true,
								},
								"rate": schema.Int64Attribute{
									Optional: true,
								},
								"run_length": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"family_inet_output": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"aggregate_export_interval": schema.Int64Attribute{
									Optional: true,
								},
								"extension_service": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
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
							Attributes: map[string]schema.Attribute{
								"max_packets_per_second": schema.Int64Attribute{
									Optional: true,
								},
								"maximum_packet_length": schema.Int64Attribute{
									Optional: true,
								},
								"rate": schema.Int64Attribute{
									Optional: true,
								},
								"run_length": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
					"family_inet6_output": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"aggregate_export_interval": schema.Int64Attribute{
									Optional: true,
								},
								"extension_service": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
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
							Blocks: map[string]schema.Block{
								"flow_server": schema.SetNestedBlock{
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
					"family_mpls_input": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"max_packets_per_second": schema.Int64Attribute{
									Optional: true,
								},
								"maximum_packet_length": schema.Int64Attribute{
									Optional: true,
								},
								"rate": schema.Int64Attribute{
									Optional: true,
								},
								"run_length": schema.Int64Attribute{
									Optional: true,
								},
							},
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
							Blocks: map[string]schema.Block{
								"flow_server": schema.SetNestedBlock{
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
					"input": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"max_packets_per_second": schema.Int64Attribute{
									Optional: true,
								},
								"maximum_packet_length": schema.Int64Attribute{
									Optional: true,
								},
								"rate": schema.Int64Attribute{
									Optional: true,
								},
								"run_length": schema.Int64Attribute{
									Optional: true,
								},
							},
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
	//nolint:lll
	type modelV0 struct {
		ID              types.String `tfsdk:"id"`
		Name            types.String `tfsdk:"name"`
		Disable         types.Bool   `tfsdk:"disable"`
		FamilyInetInput []struct {
			MaxPacketsPerSecond types.Int64 `tfsdk:"max_packets_per_second"`
			MaximumPacketLength types.Int64 `tfsdk:"maximum_packet_length"`
			Rate                types.Int64 `tfsdk:"rate"`
			RunLength           types.Int64 `tfsdk:"run_length"`
		} `tfsdk:"family_inet_input"`
		FamilyInetOutput []struct {
			AggregateExportInterval  types.Int64    `tfsdk:"aggregate_export_interval"`
			ExtensionService         []types.String `tfsdk:"extension_service"`
			FlowActiveTimeout        types.Int64    `tfsdk:"flow_active_timeout"`
			FlowInactiveTimeout      types.Int64    `tfsdk:"flow_inactive_timeout"`
			InlineJflowExportRate    types.Int64    `tfsdk:"inline_jflow_export_rate"`
			InlineJflowSourceAddress types.String   `tfsdk:"inline_jflow_source_address"`
			FlowServer               []struct {
				Hostname                                         types.String `tfsdk:"hostname"`
				Port                                             types.Int64  `tfsdk:"port"`
				AggregationAutonomousSystem                      types.Bool   `tfsdk:"aggregation_autonomous_system"`
				AggregationDestinationPrefix                     types.Bool   `tfsdk:"aggregation_destination_prefix"`
				AggregationProtocolPort                          types.Bool   `tfsdk:"aggregation_protocol_port"`
				AggregationSourceDestinationPrefix               types.Bool   `tfsdk:"aggregation_source_destination_prefix"`
				AggregationSourceDestinationPrefixCaidaCompliant types.Bool   `tfsdk:"aggregation_source_destination_prefix_caida_compliant"`
				AggregationSourcePrefix                          types.Bool   `tfsdk:"aggregation_source_prefix"`
				AutonomousSystemType                             types.String `tfsdk:"autonomous_system_type"`
				Dscp                                             types.Int64  `tfsdk:"dscp"`
				ForwardingClass                                  types.String `tfsdk:"forwarding_class"`
				LocalDump                                        types.Bool   `tfsdk:"local_dump"`
				NoLocalDump                                      types.Bool   `tfsdk:"no_local_dump"`
				RoutingInstance                                  types.String `tfsdk:"routing_instance"`
				SourceAddress                                    types.String `tfsdk:"source_address"`
				Version                                          types.Int64  `tfsdk:"version"`
				Version9Template                                 types.String `tfsdk:"version9_template"`
				VersionIPFixTemplate                             types.String `tfsdk:"version_ipfix_template"`
			} `tfsdk:"flow_server"`
			Interface []forwardingoptionsSamplingInstanceBlockOutputBlockInterface `tfsdk:"interface"`
		} `tfsdk:"family_inet_output"`
		FamilyInet6Input []struct {
			MaxPacketsPerSecond types.Int64 `tfsdk:"max_packets_per_second"`
			MaximumPacketLength types.Int64 `tfsdk:"maximum_packet_length"`
			Rate                types.Int64 `tfsdk:"rate"`
			RunLength           types.Int64 `tfsdk:"run_length"`
		} `tfsdk:"family_inet6_input"`
		FamilyInet6Output []struct {
			AggregateExportInterval  types.Int64    `tfsdk:"aggregate_export_interval"`
			ExtensionService         []types.String `tfsdk:"extension_service"`
			FlowActiveTimeout        types.Int64    `tfsdk:"flow_active_timeout"`
			FlowInactiveTimeout      types.Int64    `tfsdk:"flow_inactive_timeout"`
			InlineJflowExportRate    types.Int64    `tfsdk:"inline_jflow_export_rate"`
			InlineJflowSourceAddress types.String   `tfsdk:"inline_jflow_source_address"`
			FlowServer               []struct {
				Hostname                                         types.String `tfsdk:"hostname"`
				Port                                             types.Int64  `tfsdk:"port"`
				AggregationAutonomousSystem                      types.Bool   `tfsdk:"aggregation_autonomous_system"`
				AggregationDestinationPrefix                     types.Bool   `tfsdk:"aggregation_destination_prefix"`
				AggregationProtocolPort                          types.Bool   `tfsdk:"aggregation_protocol_port"`
				AggregationSourceDestinationPrefix               types.Bool   `tfsdk:"aggregation_source_destination_prefix"`
				AggregationSourceDestinationPrefixCaidaCompliant types.Bool   `tfsdk:"aggregation_source_destination_prefix_caida_compliant"`
				AggregationSourcePrefix                          types.Bool   `tfsdk:"aggregation_source_prefix"`
				AutonomousSystemType                             types.String `tfsdk:"autonomous_system_type"`
				Dscp                                             types.Int64  `tfsdk:"dscp"`
				ForwardingClass                                  types.String `tfsdk:"forwarding_class"`
				LocalDump                                        types.Bool   `tfsdk:"local_dump"`
				NoLocalDump                                      types.Bool   `tfsdk:"no_local_dump"`
				RoutingInstance                                  types.String `tfsdk:"routing_instance"`
				SourceAddress                                    types.String `tfsdk:"source_address"`
				Version9Template                                 types.String `tfsdk:"version9_template"`
				VersionIPFixTemplate                             types.String `tfsdk:"version_ipfix_template"`
			} `tfsdk:"flow_server"`
			Interface []forwardingoptionsSamplingInstanceBlockOutputBlockInterface `tfsdk:"interface"`
		} `tfsdk:"family_inet6_output"`
		FamilyMplsInput []struct {
			MaxPacketsPerSecond types.Int64 `tfsdk:"max_packets_per_second"`
			MaximumPacketLength types.Int64 `tfsdk:"maximum_packet_length"`
			Rate                types.Int64 `tfsdk:"rate"`
			RunLength           types.Int64 `tfsdk:"run_length"`
		} `tfsdk:"family_mpls_input"`
		FamilyMplsOutput []struct {
			AggregateExportInterval  types.Int64  `tfsdk:"aggregate_export_interval"`
			FlowActiveTimeout        types.Int64  `tfsdk:"flow_active_timeout"`
			FlowInactiveTimeout      types.Int64  `tfsdk:"flow_inactive_timeout"`
			InlineJflowExportRate    types.Int64  `tfsdk:"inline_jflow_export_rate"`
			InlineJflowSourceAddress types.String `tfsdk:"inline_jflow_source_address"`
			FlowServer               []struct {
				Hostname                                         types.String `tfsdk:"hostname"`
				Port                                             types.Int64  `tfsdk:"port"`
				AggregationAutonomousSystem                      types.Bool   `tfsdk:"aggregation_autonomous_system"`
				AggregationDestinationPrefix                     types.Bool   `tfsdk:"aggregation_destination_prefix"`
				AggregationProtocolPort                          types.Bool   `tfsdk:"aggregation_protocol_port"`
				AggregationSourceDestinationPrefix               types.Bool   `tfsdk:"aggregation_source_destination_prefix"`
				AggregationSourceDestinationPrefixCaidaCompliant types.Bool   `tfsdk:"aggregation_source_destination_prefix_caida_compliant"`
				AggregationSourcePrefix                          types.Bool   `tfsdk:"aggregation_source_prefix"`
				AutonomousSystemType                             types.String `tfsdk:"autonomous_system_type"`
				Dscp                                             types.Int64  `tfsdk:"dscp"`
				ForwardingClass                                  types.String `tfsdk:"forwarding_class"`
				LocalDump                                        types.Bool   `tfsdk:"local_dump"`
				NoLocalDump                                      types.Bool   `tfsdk:"no_local_dump"`
				RoutingInstance                                  types.String `tfsdk:"routing_instance"`
				SourceAddress                                    types.String `tfsdk:"source_address"`
				Version9Template                                 types.String `tfsdk:"version9_template"`
				VersionIPFixTemplate                             types.String `tfsdk:"version_ipfix_template"`
			} `tfsdk:"flow_server"`
			Interface []forwardingoptionsSamplingInstanceBlockOutputBlockInterface `tfsdk:"interface"`
		} `tfsdk:"family_mpls_output"`
		Input []struct {
			MaxPacketsPerSecond types.Int64 `tfsdk:"max_packets_per_second"`
			MaximumPacketLength types.Int64 `tfsdk:"maximum_packet_length"`
			Rate                types.Int64 `tfsdk:"rate"`
			RunLength           types.Int64 `tfsdk:"run_length"`
		} `tfsdk:"input"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 forwardingoptionsSamplingInstanceData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.RoutingInstance = types.StringValue(junos.DefaultW)
	dataV1.Disable = dataV0.Disable
	if len(dataV0.FamilyInetInput) > 0 {
		dataV1.FamilyInetInput = &forwardingoptionsSamplingInstanceBlockInput{
			MaxPacketsPerSecond: dataV0.FamilyInetInput[0].MaxPacketsPerSecond,
			MaximumPacketLength: dataV0.FamilyInetInput[0].MaximumPacketLength,
			Rate:                dataV0.FamilyInetInput[0].Rate,
			RunLength:           dataV0.FamilyInetInput[0].RunLength,
		}
	}
	if len(dataV0.FamilyInetOutput) > 0 {
		dataV1.FamilyInetOutput = &forwardingoptionsSamplingInstanceBlockFamilyInetOutput{
			AggregateExportInterval:  dataV0.FamilyInetOutput[0].AggregateExportInterval,
			ExtensionService:         dataV0.FamilyInetOutput[0].ExtensionService,
			FlowActiveTimeout:        dataV0.FamilyInetOutput[0].FlowActiveTimeout,
			FlowInactiveTimeout:      dataV0.FamilyInetOutput[0].FlowInactiveTimeout,
			InlineJflowExportRate:    dataV0.FamilyInetOutput[0].InlineJflowExportRate,
			InlineJflowSourceAddress: dataV0.FamilyInetOutput[0].InlineJflowSourceAddress,
			Interface:                dataV0.FamilyInetOutput[0].Interface,
		}
		for _, blockV0 := range dataV0.FamilyInetOutput[0].FlowServer {
			dataV1.FamilyInetOutput.FlowServer = append(dataV1.FamilyInetOutput.FlowServer,
				forwardingoptionsSamplingInstanceBlockFamilyInetOutputBlockFlowServer{
					AggregationAutonomousSystem:                      blockV0.AggregationAutonomousSystem,
					AggregationDestinationPrefix:                     blockV0.AggregationDestinationPrefix,
					AggregationProtocolPort:                          blockV0.AggregationProtocolPort,
					AggregationSourceDestinationPrefix:               blockV0.AggregationSourceDestinationPrefix,
					AggregationSourceDestinationPrefixCaidaCompliant: blockV0.AggregationSourceDestinationPrefixCaidaCompliant,
					AggregationSourcePrefix:                          blockV0.AggregationSourcePrefix,
					LocalDump:                                        blockV0.LocalDump,
					NoLocalDump:                                      blockV0.NoLocalDump,
					Hostname:                                         blockV0.Hostname,
					Port:                                             blockV0.Port,
					AutonomousSystemType:                             blockV0.AutonomousSystemType,
					Dscp:                                             blockV0.Dscp,
					ForwardingClass:                                  blockV0.ForwardingClass,
					RoutingInstance:                                  blockV0.RoutingInstance,
					SourceAddress:                                    blockV0.SourceAddress,
					Version:                                          blockV0.Version,
					Version9Template:                                 blockV0.Version9Template,
					VersionIPFixTemplate:                             blockV0.VersionIPFixTemplate,
				},
			)
		}
	}
	if len(dataV0.FamilyInet6Input) > 0 {
		dataV1.FamilyInet6Input = &forwardingoptionsSamplingInstanceBlockInput{
			MaxPacketsPerSecond: dataV0.FamilyInet6Input[0].MaxPacketsPerSecond,
			MaximumPacketLength: dataV0.FamilyInet6Input[0].MaximumPacketLength,
			Rate:                dataV0.FamilyInet6Input[0].Rate,
			RunLength:           dataV0.FamilyInet6Input[0].RunLength,
		}
	}
	if len(dataV0.FamilyInet6Output) > 0 {
		dataV1.FamilyInet6Output = &forwardingoptionsSamplingInstanceBlockFamilyInet6Output{
			AggregateExportInterval:  dataV0.FamilyInet6Output[0].AggregateExportInterval,
			ExtensionService:         dataV0.FamilyInet6Output[0].ExtensionService,
			FlowActiveTimeout:        dataV0.FamilyInet6Output[0].FlowActiveTimeout,
			FlowInactiveTimeout:      dataV0.FamilyInet6Output[0].FlowInactiveTimeout,
			InlineJflowExportRate:    dataV0.FamilyInet6Output[0].InlineJflowExportRate,
			InlineJflowSourceAddress: dataV0.FamilyInet6Output[0].InlineJflowSourceAddress,
			Interface:                dataV0.FamilyInet6Output[0].Interface,
		}
		for _, blockV0 := range dataV0.FamilyInet6Output[0].FlowServer {
			dataV1.FamilyInet6Output.FlowServer = append(dataV1.FamilyInet6Output.FlowServer,
				forwardingoptionsSamplingInstanceBlockOutputBlockFlowServer{
					AggregationAutonomousSystem:                      blockV0.AggregationAutonomousSystem,
					AggregationDestinationPrefix:                     blockV0.AggregationDestinationPrefix,
					AggregationProtocolPort:                          blockV0.AggregationProtocolPort,
					AggregationSourceDestinationPrefix:               blockV0.AggregationSourceDestinationPrefix,
					AggregationSourceDestinationPrefixCaidaCompliant: blockV0.AggregationSourceDestinationPrefixCaidaCompliant,
					AggregationSourcePrefix:                          blockV0.AggregationSourcePrefix,
					LocalDump:                                        blockV0.LocalDump,
					NoLocalDump:                                      blockV0.NoLocalDump,
					Hostname:                                         blockV0.Hostname,
					Port:                                             blockV0.Port,
					AutonomousSystemType:                             blockV0.AutonomousSystemType,
					Dscp:                                             blockV0.Dscp,
					ForwardingClass:                                  blockV0.ForwardingClass,
					RoutingInstance:                                  blockV0.RoutingInstance,
					SourceAddress:                                    blockV0.SourceAddress,
					Version9Template:                                 blockV0.Version9Template,
					VersionIPFixTemplate:                             blockV0.VersionIPFixTemplate,
				},
			)
		}
	}
	if len(dataV0.FamilyMplsInput) > 0 {
		dataV1.FamilyMplsInput = &forwardingoptionsSamplingInstanceBlockInput{
			MaxPacketsPerSecond: dataV0.FamilyMplsInput[0].MaxPacketsPerSecond,
			MaximumPacketLength: dataV0.FamilyMplsInput[0].MaximumPacketLength,
			Rate:                dataV0.FamilyMplsInput[0].Rate,
			RunLength:           dataV0.FamilyMplsInput[0].RunLength,
		}
	}
	if len(dataV0.FamilyMplsOutput) > 0 {
		dataV1.FamilyMplsOutput = &forwardingoptionsSamplingInstanceBlockFamilyMplsOutput{
			AggregateExportInterval:  dataV0.FamilyMplsOutput[0].AggregateExportInterval,
			FlowActiveTimeout:        dataV0.FamilyMplsOutput[0].FlowActiveTimeout,
			FlowInactiveTimeout:      dataV0.FamilyMplsOutput[0].FlowInactiveTimeout,
			InlineJflowExportRate:    dataV0.FamilyMplsOutput[0].InlineJflowExportRate,
			InlineJflowSourceAddress: dataV0.FamilyMplsOutput[0].InlineJflowSourceAddress,
			Interface:                dataV0.FamilyMplsOutput[0].Interface,
		}
		for _, blockV0 := range dataV0.FamilyMplsOutput[0].FlowServer {
			dataV1.FamilyMplsOutput.FlowServer = append(dataV1.FamilyMplsOutput.FlowServer,
				forwardingoptionsSamplingInstanceBlockOutputBlockFlowServer{
					AggregationAutonomousSystem:                      blockV0.AggregationAutonomousSystem,
					AggregationDestinationPrefix:                     blockV0.AggregationDestinationPrefix,
					AggregationProtocolPort:                          blockV0.AggregationProtocolPort,
					AggregationSourceDestinationPrefix:               blockV0.AggregationSourceDestinationPrefix,
					AggregationSourceDestinationPrefixCaidaCompliant: blockV0.AggregationSourceDestinationPrefixCaidaCompliant,
					AggregationSourcePrefix:                          blockV0.AggregationSourcePrefix,
					LocalDump:                                        blockV0.LocalDump,
					NoLocalDump:                                      blockV0.NoLocalDump,
					Hostname:                                         blockV0.Hostname,
					Port:                                             blockV0.Port,
					AutonomousSystemType:                             blockV0.AutonomousSystemType,
					Dscp:                                             blockV0.Dscp,
					ForwardingClass:                                  blockV0.ForwardingClass,
					RoutingInstance:                                  blockV0.RoutingInstance,
					SourceAddress:                                    blockV0.SourceAddress,
					Version9Template:                                 blockV0.Version9Template,
					VersionIPFixTemplate:                             blockV0.VersionIPFixTemplate,
				},
			)
		}
	}
	if len(dataV0.Input) > 0 {
		dataV1.Input = &forwardingoptionsSamplingInstanceBlockInput{
			MaxPacketsPerSecond: dataV0.Input[0].MaxPacketsPerSecond,
			MaximumPacketLength: dataV0.Input[0].MaximumPacketLength,
			Rate:                dataV0.Input[0].Rate,
			RunLength:           dataV0.Input[0].RunLength,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
