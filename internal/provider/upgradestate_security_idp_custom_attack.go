package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityIdpCustomAttack) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"recommended_action": schema.StringAttribute{
						Required: true,
					},
					"severity": schema.StringAttribute{
						Required: true,
					},
					"time_binding_count": schema.Int64Attribute{
						Optional: true,
					},
					"time_binding_scope": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"attack_type_anomaly": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"direction": schema.StringAttribute{
									Required: true,
								},
								"test": schema.StringAttribute{
									Required: true,
								},
								"service": schema.StringAttribute{
									Required: true,
								},
								"shellcode": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"attack_type_chain": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"expression": schema.StringAttribute{
									Optional: true,
								},
								"order": schema.BoolAttribute{
									Optional: true,
								},
								"protocol_binding": schema.StringAttribute{
									Optional: true,
								},
								"reset": schema.BoolAttribute{
									Optional: true,
								},
								"scope": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"member": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Required: true,
											},
										},
										Blocks: map[string]schema.Block{
											"attack_type_anomaly": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"direction": schema.StringAttribute{
															Required: true,
														},
														"test": schema.StringAttribute{
															Required: true,
														},
														"shellcode": schema.StringAttribute{
															Optional: true,
														},
													},
												},
											},
											"attack_type_signature": schema.ListNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"context": schema.StringAttribute{
															Required: true,
														},
														"direction": schema.StringAttribute{
															Required: true,
														},
														"negate": schema.BoolAttribute{
															Optional: true,
														},
														"pattern": schema.StringAttribute{
															Optional: true,
														},
														"pattern_pcre": schema.StringAttribute{
															Optional: true,
														},
														"regexp": schema.StringAttribute{
															Optional: true,
														},
														"shellcode": schema.StringAttribute{
															Optional: true,
														},
													},
													Blocks: map[string]schema.Block{
														"protocol_icmp": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"checksum_validate_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"checksum_validate_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"code_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"code_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"data_length_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"data_length_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"identification_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"identification_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"sequence_number_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"sequence_number_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"type_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"type_value": schema.Int64Attribute{
																		Optional: true,
																	},
																},
															},
														},
														"protocol_icmpv6": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"checksum_validate_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"checksum_validate_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"code_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"code_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"data_length_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"data_length_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"identification_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"identification_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"sequence_number_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"sequence_number_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"type_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"type_value": schema.Int64Attribute{
																		Optional: true,
																	},
																},
															},
														},
														"protocol_ipv4": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"checksum_validate_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"checksum_validate_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"destination_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"destination_value": schema.StringAttribute{
																		Optional: true,
																	},
																	"identification_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"identification_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"ihl_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"ihl_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"ip_flags": schema.SetAttribute{
																		ElementType: types.StringType,
																		Optional:    true,
																	},
																	"protocol_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"protocol_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"source_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"source_value": schema.StringAttribute{
																		Optional: true,
																	},
																	"tos_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"tos_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"total_length_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"total_length_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"ttl_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"ttl_value": schema.Int64Attribute{
																		Optional: true,
																	},
																},
															},
														},
														"protocol_ipv6": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"destination_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"destination_value": schema.StringAttribute{
																		Optional: true,
																	},
																	"extension_header_destination_option_home_address_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"extension_header_destination_option_home_address_value": schema.StringAttribute{
																		Optional: true,
																	},
																	"extension_header_destination_option_type_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"extension_header_destination_option_type_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"extension_header_routing_header_type_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"extension_header_routing_header_type_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"flow_label_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"flow_label_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"hop_limit_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"hop_limit_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"next_header_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"next_header_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"payload_length_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"payload_length_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"source_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"source_value": schema.StringAttribute{
																		Optional: true,
																	},
																	"traffic_class_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"traffic_class_value": schema.Int64Attribute{
																		Optional: true,
																	},
																},
															},
														},
														"protocol_tcp": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"ack_number_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"ack_number_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"checksum_validate_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"checksum_validate_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"data_length_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"data_length_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"destination_port_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"destination_port_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"header_length_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"header_length_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"mss_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"mss_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"option_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"option_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"reserved_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"reserved_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"sequence_number_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"sequence_number_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"source_port_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"source_port_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"tcp_flags": schema.SetAttribute{
																		ElementType: types.StringType,
																		Optional:    true,
																	},
																	"urgent_pointer_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"urgent_pointer_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"window_scale_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"window_scale_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"window_size_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"window_size_value": schema.Int64Attribute{
																		Optional: true,
																	},
																},
															},
														},
														"protocol_udp": schema.ListNestedBlock{
															NestedObject: schema.NestedBlockObject{
																Attributes: map[string]schema.Attribute{
																	"checksum_validate_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"checksum_validate_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"data_length_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"data_length_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"destination_port_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"destination_port_value": schema.Int64Attribute{
																		Optional: true,
																	},
																	"source_port_match": schema.StringAttribute{
																		Optional: true,
																	},
																	"source_port_value": schema.Int64Attribute{
																		Optional: true,
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
					},
					"attack_type_signature": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"context": schema.StringAttribute{
									Required: true,
								},
								"direction": schema.StringAttribute{
									Required: true,
								},
								"negate": schema.BoolAttribute{
									Optional: true,
								},
								"pattern": schema.StringAttribute{
									Optional: true,
								},
								"pattern_pcre": schema.StringAttribute{
									Optional: true,
								},
								"protocol_binding": schema.StringAttribute{
									Optional: true,
								},
								"regexp": schema.StringAttribute{
									Optional: true,
								},
								"shellcode": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"protocol_icmp": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"checksum_validate_match": schema.StringAttribute{
												Optional: true,
											},
											"checksum_validate_value": schema.Int64Attribute{
												Optional: true,
											},
											"code_match": schema.StringAttribute{
												Optional: true,
											},
											"code_value": schema.Int64Attribute{
												Optional: true,
											},
											"data_length_match": schema.StringAttribute{
												Optional: true,
											},
											"data_length_value": schema.Int64Attribute{
												Optional: true,
											},
											"identification_match": schema.StringAttribute{
												Optional: true,
											},
											"identification_value": schema.Int64Attribute{
												Optional: true,
											},
											"sequence_number_match": schema.StringAttribute{
												Optional: true,
											},
											"sequence_number_value": schema.Int64Attribute{
												Optional: true,
											},
											"type_match": schema.StringAttribute{
												Optional: true,
											},
											"type_value": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"protocol_icmpv6": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"checksum_validate_match": schema.StringAttribute{
												Optional: true,
											},
											"checksum_validate_value": schema.Int64Attribute{
												Optional: true,
											},
											"code_match": schema.StringAttribute{
												Optional: true,
											},
											"code_value": schema.Int64Attribute{
												Optional: true,
											},
											"data_length_match": schema.StringAttribute{
												Optional: true,
											},
											"data_length_value": schema.Int64Attribute{
												Optional: true,
											},
											"identification_match": schema.StringAttribute{
												Optional: true,
											},
											"identification_value": schema.Int64Attribute{
												Optional: true,
											},
											"sequence_number_match": schema.StringAttribute{
												Optional: true,
											},
											"sequence_number_value": schema.Int64Attribute{
												Optional: true,
											},
											"type_match": schema.StringAttribute{
												Optional: true,
											},
											"type_value": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"protocol_ipv4": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"checksum_validate_match": schema.StringAttribute{
												Optional: true,
											},
											"checksum_validate_value": schema.Int64Attribute{
												Optional: true,
											},
											"destination_match": schema.StringAttribute{
												Optional: true,
											},
											"destination_value": schema.StringAttribute{
												Optional: true,
											},
											"identification_match": schema.StringAttribute{
												Optional: true,
											},
											"identification_value": schema.Int64Attribute{
												Optional: true,
											},
											"ihl_match": schema.StringAttribute{
												Optional: true,
											},
											"ihl_value": schema.Int64Attribute{
												Optional: true,
											},
											"ip_flags": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"protocol_match": schema.StringAttribute{
												Optional: true,
											},
											"protocol_value": schema.Int64Attribute{
												Optional: true,
											},
											"source_match": schema.StringAttribute{
												Optional: true,
											},
											"source_value": schema.StringAttribute{
												Optional: true,
											},
											"tos_match": schema.StringAttribute{
												Optional: true,
											},
											"tos_value": schema.Int64Attribute{
												Optional: true,
											},
											"total_length_match": schema.StringAttribute{
												Optional: true,
											},
											"total_length_value": schema.Int64Attribute{
												Optional: true,
											},
											"ttl_match": schema.StringAttribute{
												Optional: true,
											},
											"ttl_value": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"protocol_ipv6": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"destination_match": schema.StringAttribute{
												Optional: true,
											},
											"destination_value": schema.StringAttribute{
												Optional: true,
											},
											"extension_header_destination_option_home_address_match": schema.StringAttribute{
												Optional: true,
											},
											"extension_header_destination_option_home_address_value": schema.StringAttribute{
												Optional: true,
											},
											"extension_header_destination_option_type_match": schema.StringAttribute{
												Optional: true,
											},
											"extension_header_destination_option_type_value": schema.Int64Attribute{
												Optional: true,
											},
											"extension_header_routing_header_type_match": schema.StringAttribute{
												Optional: true,
											},
											"extension_header_routing_header_type_value": schema.Int64Attribute{
												Optional: true,
											},
											"flow_label_match": schema.StringAttribute{
												Optional: true,
											},
											"flow_label_value": schema.Int64Attribute{
												Optional: true,
											},
											"hop_limit_match": schema.StringAttribute{
												Optional: true,
											},
											"hop_limit_value": schema.Int64Attribute{
												Optional: true,
											},
											"next_header_match": schema.StringAttribute{
												Optional: true,
											},
											"next_header_value": schema.Int64Attribute{
												Optional: true,
											},
											"payload_length_match": schema.StringAttribute{
												Optional: true,
											},
											"payload_length_value": schema.Int64Attribute{
												Optional: true,
											},
											"source_match": schema.StringAttribute{
												Optional: true,
											},
											"source_value": schema.StringAttribute{
												Optional: true,
											},
											"traffic_class_match": schema.StringAttribute{
												Optional: true,
											},
											"traffic_class_value": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"protocol_tcp": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"ack_number_match": schema.StringAttribute{
												Optional: true,
											},
											"ack_number_value": schema.Int64Attribute{
												Optional: true,
											},
											"checksum_validate_match": schema.StringAttribute{
												Optional: true,
											},
											"checksum_validate_value": schema.Int64Attribute{
												Optional: true,
											},
											"data_length_match": schema.StringAttribute{
												Optional: true,
											},
											"data_length_value": schema.Int64Attribute{
												Optional: true,
											},
											"destination_port_match": schema.StringAttribute{
												Optional: true,
											},
											"destination_port_value": schema.Int64Attribute{
												Optional: true,
											},
											"header_length_match": schema.StringAttribute{
												Optional: true,
											},
											"header_length_value": schema.Int64Attribute{
												Optional: true,
											},
											"mss_match": schema.StringAttribute{
												Optional: true,
											},
											"mss_value": schema.Int64Attribute{
												Optional: true,
											},
											"option_match": schema.StringAttribute{
												Optional: true,
											},
											"option_value": schema.Int64Attribute{
												Optional: true,
											},
											"reserved_match": schema.StringAttribute{
												Optional: true,
											},
											"reserved_value": schema.Int64Attribute{
												Optional: true,
											},
											"sequence_number_match": schema.StringAttribute{
												Optional: true,
											},
											"sequence_number_value": schema.Int64Attribute{
												Optional: true,
											},
											"source_port_match": schema.StringAttribute{
												Optional: true,
											},
											"source_port_value": schema.Int64Attribute{
												Optional: true,
											},
											"tcp_flags": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
											},
											"urgent_pointer_match": schema.StringAttribute{
												Optional: true,
											},
											"urgent_pointer_value": schema.Int64Attribute{
												Optional: true,
											},
											"window_scale_match": schema.StringAttribute{
												Optional: true,
											},
											"window_scale_value": schema.Int64Attribute{
												Optional: true,
											},
											"window_size_match": schema.StringAttribute{
												Optional: true,
											},
											"window_size_value": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
								"protocol_udp": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"checksum_validate_match": schema.StringAttribute{
												Optional: true,
											},
											"checksum_validate_value": schema.Int64Attribute{
												Optional: true,
											},
											"data_length_match": schema.StringAttribute{
												Optional: true,
											},
											"data_length_value": schema.Int64Attribute{
												Optional: true,
											},
											"destination_port_match": schema.StringAttribute{
												Optional: true,
											},
											"destination_port_value": schema.Int64Attribute{
												Optional: true,
											},
											"source_port_match": schema.StringAttribute{
												Optional: true,
											},
											"source_port_value": schema.Int64Attribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSecurityIdpCustomAttackV0toV1,
		},
	}
}

func upgradeSecurityIdpCustomAttackV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	//nolint:lll
	type modelV0 struct {
		ID                types.String `tfsdk:"id"`
		Name              types.String `tfsdk:"name"`
		RecommendedAction types.String `tfsdk:"recommended_action"`
		Severity          types.String `tfsdk:"severity"`
		TimeBindingCount  types.Int64  `tfsdk:"time_binding_count"`
		TimeBindingScope  types.String `tfsdk:"time_binding_scope"`
		AttackTypeAnomaly []struct {
			Direction types.String `tfsdk:"direction"`
			Test      types.String `tfsdk:"test"`
			Service   types.String `tfsdk:"service"`
			Shellcode types.String `tfsdk:"shellcode"`
		} `tfsdk:"attack_type_anomaly"`
		AttackTypeChain []struct {
			Expression      types.String `tfsdk:"expression"`
			Order           types.Bool   `tfsdk:"order"`
			ProtocolBinding types.String `tfsdk:"protocol_binding"`
			Reset           types.Bool   `tfsdk:"reset"`
			Scope           types.String `tfsdk:"scope"`
			Member          []struct {
				Name              types.String `tfsdk:"name"`
				AttackTypeAnomaly []struct {
					Direction types.String `tfsdk:"direction"`
					Test      types.String `tfsdk:"test"`
					Shellcode types.String `tfsdk:"shellcode"`
				} `tfsdk:"attack_type_anomaly"`
				AttackTypeSignature []struct {
					Context      types.String `tfsdk:"context"`
					Direction    types.String `tfsdk:"direction"`
					Negate       types.Bool   `tfsdk:"negate"`
					Pattern      types.String `tfsdk:"pattern"`
					PatternPcre  types.String `tfsdk:"pattern_pcre"`
					Regexp       types.String `tfsdk:"regexp"`
					Shellcode    types.String `tfsdk:"shellcode"`
					ProtocolIcmp []struct {
						ChecksumValidateMatch types.String `tfsdk:"checksum_validate_match"`
						ChecksumValidateValue types.Int64  `tfsdk:"checksum_validate_value"`
						CodeMatch             types.String `tfsdk:"code_match"`
						CodeValue             types.Int64  `tfsdk:"code_value"`
						DataLengthMatch       types.String `tfsdk:"data_length_match"`
						DataLengthValue       types.Int64  `tfsdk:"data_length_value"`
						IdentificationMatch   types.String `tfsdk:"identification_match"`
						IdentificationValue   types.Int64  `tfsdk:"identification_value"`
						SequenceNumberMatch   types.String `tfsdk:"sequence_number_match"`
						SequenceNumberValue   types.Int64  `tfsdk:"sequence_number_value"`
						TypeMatch             types.String `tfsdk:"type_match"`
						TypeValue             types.Int64  `tfsdk:"type_value"`
					} `tfsdk:"protocol_icmp"`
					ProtocolIcmpv6 []struct {
						ChecksumValidateMatch types.String `tfsdk:"checksum_validate_match"`
						ChecksumValidateValue types.Int64  `tfsdk:"checksum_validate_value"`
						CodeMatch             types.String `tfsdk:"code_match"`
						CodeValue             types.Int64  `tfsdk:"code_value"`
						DataLengthMatch       types.String `tfsdk:"data_length_match"`
						DataLengthValue       types.Int64  `tfsdk:"data_length_value"`
						IdentificationMatch   types.String `tfsdk:"identification_match"`
						IdentificationValue   types.Int64  `tfsdk:"identification_value"`
						SequenceNumberMatch   types.String `tfsdk:"sequence_number_match"`
						SequenceNumberValue   types.Int64  `tfsdk:"sequence_number_value"`
						TypeMatch             types.String `tfsdk:"type_match"`
						TypeValue             types.Int64  `tfsdk:"type_value"`
					} `tfsdk:"protocol_icmpv6"`
					ProtocolIPv4 []struct {
						ChecksumValidateMatch types.String   `tfsdk:"checksum_validate_match"`
						ChecksumValidateValue types.Int64    `tfsdk:"checksum_validate_value"`
						DestinationMatch      types.String   `tfsdk:"destination_match"`
						DestinationValue      types.String   `tfsdk:"destination_value"`
						IdentificationMatch   types.String   `tfsdk:"identification_match"`
						IdentificationValue   types.Int64    `tfsdk:"identification_value"`
						IhlMatch              types.String   `tfsdk:"ihl_match"`
						IhlValue              types.Int64    `tfsdk:"ihl_value"`
						IPFlags               []types.String `tfsdk:"ip_flags"`
						ProtocolMatch         types.String   `tfsdk:"protocol_match"`
						ProtocolValue         types.Int64    `tfsdk:"protocol_value"`
						SourceMatch           types.String   `tfsdk:"source_match"`
						SourceValue           types.String   `tfsdk:"source_value"`
						TosMatch              types.String   `tfsdk:"tos_match"`
						TosValue              types.Int64    `tfsdk:"tos_value"`
						TotalLengthMatch      types.String   `tfsdk:"total_length_match"`
						TotalLengthValue      types.Int64    `tfsdk:"total_length_value"`
						TTLMatch              types.String   `tfsdk:"ttl_match"`
						TTLValue              types.Int64    `tfsdk:"ttl_value"`
					} `tfsdk:"protocol_ipv4"`
					ProtocolIPv6 []struct {
						DestinationMatch                                 types.String `tfsdk:"destination_match"`
						DestinationValue                                 types.String `tfsdk:"destination_value"`
						ExtensionHeaderDestinationOptionHomeAddressMatch types.String `tfsdk:"extension_header_destination_option_home_address_match"`
						ExtensionHeaderDestinationOptionHomeAddressValue types.String `tfsdk:"extension_header_destination_option_home_address_value"`
						ExtensionHeaderDestinationOptionTypeMatch        types.String `tfsdk:"extension_header_destination_option_type_match"`
						ExtensionHeaderDestinationOptionTypeValue        types.Int64  `tfsdk:"extension_header_destination_option_type_value"`
						ExtensionHeaderRoutingHeaderTypeMatch            types.String `tfsdk:"extension_header_routing_header_type_match"`
						ExtensionHeaderRoutingHeaderTypeValue            types.Int64  `tfsdk:"extension_header_routing_header_type_value"`
						FlowLabelMatch                                   types.String `tfsdk:"flow_label_match"`
						FlowLabelValue                                   types.Int64  `tfsdk:"flow_label_value"`
						HopLimitMatch                                    types.String `tfsdk:"hop_limit_match"`
						HopLimitValue                                    types.Int64  `tfsdk:"hop_limit_value"`
						NextHeaderMatch                                  types.String `tfsdk:"next_header_match"`
						NextHeaderValue                                  types.Int64  `tfsdk:"next_header_value"`
						PayloadLengthMatch                               types.String `tfsdk:"payload_length_match"`
						PayloadLengthValue                               types.Int64  `tfsdk:"payload_length_value"`
						SourceMatch                                      types.String `tfsdk:"source_match"`
						SourceValue                                      types.String `tfsdk:"source_value"`
						TrafficClassMatch                                types.String `tfsdk:"traffic_class_match"`
						TrafficClassValue                                types.Int64  `tfsdk:"traffic_class_value"`
					} `tfsdk:"protocol_ipv6"`
					ProtocolTCP []struct {
						AckNumberMatch        types.String   `tfsdk:"ack_number_match"`
						AckNumberValue        types.Int64    `tfsdk:"ack_number_value"`
						ChecksumValidateMatch types.String   `tfsdk:"checksum_validate_match"`
						ChecksumValidateValue types.Int64    `tfsdk:"checksum_validate_value"`
						DataLengthMatch       types.String   `tfsdk:"data_length_match"`
						DataLengthValue       types.Int64    `tfsdk:"data_length_value"`
						DestinationPortMatch  types.String   `tfsdk:"destination_port_match"`
						DestinationPortValue  types.Int64    `tfsdk:"destination_port_value"`
						HeaderLengthMatch     types.String   `tfsdk:"header_length_match"`
						HeaderLengthValue     types.Int64    `tfsdk:"header_length_value"`
						MssMatch              types.String   `tfsdk:"mss_match"`
						MssValue              types.Int64    `tfsdk:"mss_value"`
						OptionMatch           types.String   `tfsdk:"option_match"`
						OptionValue           types.Int64    `tfsdk:"option_value"`
						ReservedMatch         types.String   `tfsdk:"reserved_match"`
						ReservedValue         types.Int64    `tfsdk:"reserved_value"`
						SequenceNumberMatch   types.String   `tfsdk:"sequence_number_match"`
						SequenceNumberValue   types.Int64    `tfsdk:"sequence_number_value"`
						SourcePortMatch       types.String   `tfsdk:"source_port_match"`
						SourcePortValue       types.Int64    `tfsdk:"source_port_value"`
						TCPFlags              []types.String `tfsdk:"tcp_flags"`
						UrgentPointerMatch    types.String   `tfsdk:"urgent_pointer_match"`
						UrgentPointerValue    types.Int64    `tfsdk:"urgent_pointer_value"`
						WindowScaleMatch      types.String   `tfsdk:"window_scale_match"`
						WindowScaleValue      types.Int64    `tfsdk:"window_scale_value"`
						WindowSizeMatch       types.String   `tfsdk:"window_size_match"`
						WindowSizeValue       types.Int64    `tfsdk:"window_size_value"`
					} `tfsdk:"protocol_tcp"`
					ProtocolUDP []struct {
						ChecksumValidateMatch types.String `tfsdk:"checksum_validate_match"`
						ChecksumValidateValue types.Int64  `tfsdk:"checksum_validate_value"`
						DataLengthMatch       types.String `tfsdk:"data_length_match"`
						DataLengthValue       types.Int64  `tfsdk:"data_length_value"`
						DestinationPortMatch  types.String `tfsdk:"destination_port_match"`
						DestinationPortValue  types.Int64  `tfsdk:"destination_port_value"`
						SourcePortMatch       types.String `tfsdk:"source_port_match"`
						SourcePortValue       types.Int64  `tfsdk:"source_port_value"`
					} `tfsdk:"protocol_udp"`
				} `tfsdk:"attack_type_signature"`
			} `tfsdk:"member"`
		} `tfsdk:"attack_type_chain"`
		AttackTypeSignature []struct {
			Context         types.String `tfsdk:"context"`
			Direction       types.String `tfsdk:"direction"`
			Negate          types.Bool   `tfsdk:"negate"`
			Pattern         types.String `tfsdk:"pattern"`
			PatternPcre     types.String `tfsdk:"pattern_pcre"`
			ProtocolBinding types.String `tfsdk:"protocol_binding"`
			Regexp          types.String `tfsdk:"regexp"`
			Shellcode       types.String `tfsdk:"shellcode"`
			ProtocolIcmp    []struct {
				ChecksumValidateMatch types.String `tfsdk:"checksum_validate_match"`
				ChecksumValidateValue types.Int64  `tfsdk:"checksum_validate_value"`
				CodeMatch             types.String `tfsdk:"code_match"`
				CodeValue             types.Int64  `tfsdk:"code_value"`
				DataLengthMatch       types.String `tfsdk:"data_length_match"`
				DataLengthValue       types.Int64  `tfsdk:"data_length_value"`
				IdentificationMatch   types.String `tfsdk:"identification_match"`
				IdentificationValue   types.Int64  `tfsdk:"identification_value"`
				SequenceNumberMatch   types.String `tfsdk:"sequence_number_match"`
				SequenceNumberValue   types.Int64  `tfsdk:"sequence_number_value"`
				TypeMatch             types.String `tfsdk:"type_match"`
				TypeValue             types.Int64  `tfsdk:"type_value"`
			} `tfsdk:"protocol_icmp"`
			ProtocolIcmpv6 []struct {
				ChecksumValidateMatch types.String `tfsdk:"checksum_validate_match"`
				ChecksumValidateValue types.Int64  `tfsdk:"checksum_validate_value"`
				CodeMatch             types.String `tfsdk:"code_match"`
				CodeValue             types.Int64  `tfsdk:"code_value"`
				DataLengthMatch       types.String `tfsdk:"data_length_match"`
				DataLengthValue       types.Int64  `tfsdk:"data_length_value"`
				IdentificationMatch   types.String `tfsdk:"identification_match"`
				IdentificationValue   types.Int64  `tfsdk:"identification_value"`
				SequenceNumberMatch   types.String `tfsdk:"sequence_number_match"`
				SequenceNumberValue   types.Int64  `tfsdk:"sequence_number_value"`
				TypeMatch             types.String `tfsdk:"type_match"`
				TypeValue             types.Int64  `tfsdk:"type_value"`
			} `tfsdk:"protocol_icmpv6"`
			ProtocolIPv4 []struct {
				ChecksumValidateMatch types.String   `tfsdk:"checksum_validate_match"`
				ChecksumValidateValue types.Int64    `tfsdk:"checksum_validate_value"`
				DestinationMatch      types.String   `tfsdk:"destination_match"`
				DestinationValue      types.String   `tfsdk:"destination_value"`
				IdentificationMatch   types.String   `tfsdk:"identification_match"`
				IdentificationValue   types.Int64    `tfsdk:"identification_value"`
				IhlMatch              types.String   `tfsdk:"ihl_match"`
				IhlValue              types.Int64    `tfsdk:"ihl_value"`
				IPFlags               []types.String `tfsdk:"ip_flags"`
				ProtocolMatch         types.String   `tfsdk:"protocol_match"`
				ProtocolValue         types.Int64    `tfsdk:"protocol_value"`
				SourceMatch           types.String   `tfsdk:"source_match"`
				SourceValue           types.String   `tfsdk:"source_value"`
				TosMatch              types.String   `tfsdk:"tos_match"`
				TosValue              types.Int64    `tfsdk:"tos_value"`
				TotalLengthMatch      types.String   `tfsdk:"total_length_match"`
				TotalLengthValue      types.Int64    `tfsdk:"total_length_value"`
				TTLMatch              types.String   `tfsdk:"ttl_match"`
				TTLValue              types.Int64    `tfsdk:"ttl_value"`
			} `tfsdk:"protocol_ipv4"`
			ProtocolIPv6 []struct {
				DestinationMatch                                 types.String `tfsdk:"destination_match"`
				DestinationValue                                 types.String `tfsdk:"destination_value"`
				ExtensionHeaderDestinationOptionHomeAddressMatch types.String `tfsdk:"extension_header_destination_option_home_address_match"`
				ExtensionHeaderDestinationOptionHomeAddressValue types.String `tfsdk:"extension_header_destination_option_home_address_value"`
				ExtensionHeaderDestinationOptionTypeMatch        types.String `tfsdk:"extension_header_destination_option_type_match"`
				ExtensionHeaderDestinationOptionTypeValue        types.Int64  `tfsdk:"extension_header_destination_option_type_value"`
				ExtensionHeaderRoutingHeaderTypeMatch            types.String `tfsdk:"extension_header_routing_header_type_match"`
				ExtensionHeaderRoutingHeaderTypeValue            types.Int64  `tfsdk:"extension_header_routing_header_type_value"`
				FlowLabelMatch                                   types.String `tfsdk:"flow_label_match"`
				FlowLabelValue                                   types.Int64  `tfsdk:"flow_label_value"`
				HopLimitMatch                                    types.String `tfsdk:"hop_limit_match"`
				HopLimitValue                                    types.Int64  `tfsdk:"hop_limit_value"`
				NextHeaderMatch                                  types.String `tfsdk:"next_header_match"`
				NextHeaderValue                                  types.Int64  `tfsdk:"next_header_value"`
				PayloadLengthMatch                               types.String `tfsdk:"payload_length_match"`
				PayloadLengthValue                               types.Int64  `tfsdk:"payload_length_value"`
				SourceMatch                                      types.String `tfsdk:"source_match"`
				SourceValue                                      types.String `tfsdk:"source_value"`
				TrafficClassMatch                                types.String `tfsdk:"traffic_class_match"`
				TrafficClassValue                                types.Int64  `tfsdk:"traffic_class_value"`
			} `tfsdk:"protocol_ipv6"`
			ProtocolTCP []struct {
				AckNumberMatch        types.String   `tfsdk:"ack_number_match"`
				AckNumberValue        types.Int64    `tfsdk:"ack_number_value"`
				ChecksumValidateMatch types.String   `tfsdk:"checksum_validate_match"`
				ChecksumValidateValue types.Int64    `tfsdk:"checksum_validate_value"`
				DataLengthMatch       types.String   `tfsdk:"data_length_match"`
				DataLengthValue       types.Int64    `tfsdk:"data_length_value"`
				DestinationPortMatch  types.String   `tfsdk:"destination_port_match"`
				DestinationPortValue  types.Int64    `tfsdk:"destination_port_value"`
				HeaderLengthMatch     types.String   `tfsdk:"header_length_match"`
				HeaderLengthValue     types.Int64    `tfsdk:"header_length_value"`
				MssMatch              types.String   `tfsdk:"mss_match"`
				MssValue              types.Int64    `tfsdk:"mss_value"`
				OptionMatch           types.String   `tfsdk:"option_match"`
				OptionValue           types.Int64    `tfsdk:"option_value"`
				ReservedMatch         types.String   `tfsdk:"reserved_match"`
				ReservedValue         types.Int64    `tfsdk:"reserved_value"`
				SequenceNumberMatch   types.String   `tfsdk:"sequence_number_match"`
				SequenceNumberValue   types.Int64    `tfsdk:"sequence_number_value"`
				SourcePortMatch       types.String   `tfsdk:"source_port_match"`
				SourcePortValue       types.Int64    `tfsdk:"source_port_value"`
				TCPFlags              []types.String `tfsdk:"tcp_flags"`
				UrgentPointerMatch    types.String   `tfsdk:"urgent_pointer_match"`
				UrgentPointerValue    types.Int64    `tfsdk:"urgent_pointer_value"`
				WindowScaleMatch      types.String   `tfsdk:"window_scale_match"`
				WindowScaleValue      types.Int64    `tfsdk:"window_scale_value"`
				WindowSizeMatch       types.String   `tfsdk:"window_size_match"`
				WindowSizeValue       types.Int64    `tfsdk:"window_size_value"`
			} `tfsdk:"protocol_tcp"`
			ProtocolUDP []struct {
				ChecksumValidateMatch types.String `tfsdk:"checksum_validate_match"`
				ChecksumValidateValue types.Int64  `tfsdk:"checksum_validate_value"`
				DataLengthMatch       types.String `tfsdk:"data_length_match"`
				DataLengthValue       types.Int64  `tfsdk:"data_length_value"`
				DestinationPortMatch  types.String `tfsdk:"destination_port_match"`
				DestinationPortValue  types.Int64  `tfsdk:"destination_port_value"`
				SourcePortMatch       types.String `tfsdk:"source_port_match"`
				SourcePortValue       types.Int64  `tfsdk:"source_port_value"`
			} `tfsdk:"protocol_udp"`
		} `tfsdk:"attack_type_signature"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityIdpCustomAttackData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.RecommendedAction = dataV0.RecommendedAction
	dataV1.Severity = dataV0.Severity
	dataV1.TimeBindingCount = dataV0.TimeBindingCount
	dataV1.TimeBindingScope = dataV0.TimeBindingScope
	if len(dataV0.AttackTypeAnomaly) > 0 {
		dataV1.AttackTypeAnomaly = &securityIdpCustomAttackBlockAttackTypeAnomaly{
			Service: dataV0.AttackTypeAnomaly[0].Service,
		}
		dataV1.AttackTypeAnomaly.Direction = dataV0.AttackTypeAnomaly[0].Direction
		dataV1.AttackTypeAnomaly.Test = dataV0.AttackTypeAnomaly[0].Test
		dataV1.AttackTypeAnomaly.Shellcode = dataV0.AttackTypeAnomaly[0].Shellcode
	}
	//nolint:lll
	if len(dataV0.AttackTypeChain) > 0 {
		dataV1.AttackTypeChain = &securityIdpCustomAttackBlockAttackTypeChain{
			Expression:      dataV0.AttackTypeChain[0].Expression,
			Order:           dataV0.AttackTypeChain[0].Order,
			ProtocolBinding: dataV0.AttackTypeChain[0].ProtocolBinding,
			Reset:           dataV0.AttackTypeChain[0].Reset,
			Scope:           dataV0.AttackTypeChain[0].Scope,
		}
		for _, blockV0 := range dataV0.AttackTypeChain[0].Member {
			blockV1 := securityIdpCustomAttackBlockAttackTypeChainBlockMember{
				Name: blockV0.Name,
			}
			if len(blockV0.AttackTypeAnomaly) > 0 {
				blockV1.AttackTypeAnomaly = &securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeAnomaly{
					Direction: blockV0.AttackTypeAnomaly[0].Direction,
					Test:      blockV0.AttackTypeAnomaly[0].Test,
					Shellcode: blockV0.AttackTypeAnomaly[0].Shellcode,
				}
			}
			if len(blockV0.AttackTypeSignature) > 0 {
				blockV1.AttackTypeSignature = &securityIdpCustomAttackBlockAttackTypeChainBlockMemberBlockAttackTypeSignature{
					Context:     blockV0.AttackTypeSignature[0].Context,
					Direction:   blockV0.AttackTypeSignature[0].Direction,
					Negate:      blockV0.AttackTypeSignature[0].Negate,
					Pattern:     blockV0.AttackTypeSignature[0].Pattern,
					PatternPcre: blockV0.AttackTypeSignature[0].PatternPcre,
					Regexp:      blockV0.AttackTypeSignature[0].Regexp,
					Shellcode:   blockV0.AttackTypeSignature[0].Shellcode,
				}
				if len(blockV0.AttackTypeSignature[0].ProtocolIcmp) > 0 {
					blockV1.AttackTypeSignature.ProtocolIcmp = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp{
						ChecksumValidateMatch: blockV0.AttackTypeSignature[0].ProtocolIcmp[0].ChecksumValidateMatch,
						ChecksumValidateValue: blockV0.AttackTypeSignature[0].ProtocolIcmp[0].ChecksumValidateValue,
						CodeMatch:             blockV0.AttackTypeSignature[0].ProtocolIcmp[0].CodeMatch,
						CodeValue:             blockV0.AttackTypeSignature[0].ProtocolIcmp[0].CodeValue,
						DataLengthMatch:       blockV0.AttackTypeSignature[0].ProtocolIcmp[0].DataLengthMatch,
						DataLengthValue:       blockV0.AttackTypeSignature[0].ProtocolIcmp[0].DataLengthValue,
						IdentificationMatch:   blockV0.AttackTypeSignature[0].ProtocolIcmp[0].IdentificationMatch,
						IdentificationValue:   blockV0.AttackTypeSignature[0].ProtocolIcmp[0].IdentificationValue,
						SequenceNumberMatch:   blockV0.AttackTypeSignature[0].ProtocolIcmp[0].SequenceNumberMatch,
						SequenceNumberValue:   blockV0.AttackTypeSignature[0].ProtocolIcmp[0].SequenceNumberValue,
						TypeMatch:             blockV0.AttackTypeSignature[0].ProtocolIcmp[0].TypeMatch,
						TypeValue:             blockV0.AttackTypeSignature[0].ProtocolIcmp[0].TypeValue,
					}
				}
				if len(blockV0.AttackTypeSignature[0].ProtocolIcmpv6) > 0 {
					blockV1.AttackTypeSignature.ProtocolIcmpv6 = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp{
						ChecksumValidateMatch: blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].ChecksumValidateMatch,
						ChecksumValidateValue: blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].ChecksumValidateValue,
						CodeMatch:             blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].CodeMatch,
						CodeValue:             blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].CodeValue,
						DataLengthMatch:       blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].DataLengthMatch,
						DataLengthValue:       blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].DataLengthValue,
						IdentificationMatch:   blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].IdentificationMatch,
						IdentificationValue:   blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].IdentificationValue,
						SequenceNumberMatch:   blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].SequenceNumberMatch,
						SequenceNumberValue:   blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].SequenceNumberValue,
						TypeMatch:             blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].TypeMatch,
						TypeValue:             blockV0.AttackTypeSignature[0].ProtocolIcmpv6[0].TypeValue,
					}
				}
				if len(blockV0.AttackTypeSignature[0].ProtocolIPv4) > 0 {
					blockV1.AttackTypeSignature.ProtocolIPv4 = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4{
						ChecksumValidateMatch: blockV0.AttackTypeSignature[0].ProtocolIPv4[0].ChecksumValidateMatch,
						ChecksumValidateValue: blockV0.AttackTypeSignature[0].ProtocolIPv4[0].ChecksumValidateValue,
						DestinationMatch:      blockV0.AttackTypeSignature[0].ProtocolIPv4[0].DestinationMatch,
						DestinationValue:      blockV0.AttackTypeSignature[0].ProtocolIPv4[0].DestinationValue,
						IdentificationMatch:   blockV0.AttackTypeSignature[0].ProtocolIPv4[0].IdentificationMatch,
						IdentificationValue:   blockV0.AttackTypeSignature[0].ProtocolIPv4[0].IdentificationValue,
						IhlMatch:              blockV0.AttackTypeSignature[0].ProtocolIPv4[0].IhlMatch,
						IhlValue:              blockV0.AttackTypeSignature[0].ProtocolIPv4[0].IhlValue,
						IPFlags:               blockV0.AttackTypeSignature[0].ProtocolIPv4[0].IPFlags,
						ProtocolMatch:         blockV0.AttackTypeSignature[0].ProtocolIPv4[0].ProtocolMatch,
						ProtocolValue:         blockV0.AttackTypeSignature[0].ProtocolIPv4[0].ProtocolValue,
						SourceMatch:           blockV0.AttackTypeSignature[0].ProtocolIPv4[0].SourceMatch,
						SourceValue:           blockV0.AttackTypeSignature[0].ProtocolIPv4[0].SourceValue,
						TosMatch:              blockV0.AttackTypeSignature[0].ProtocolIPv4[0].TosMatch,
						TosValue:              blockV0.AttackTypeSignature[0].ProtocolIPv4[0].TosValue,
						TotalLengthMatch:      blockV0.AttackTypeSignature[0].ProtocolIPv4[0].TotalLengthMatch,
						TotalLengthValue:      blockV0.AttackTypeSignature[0].ProtocolIPv4[0].TotalLengthValue,
						TTLMatch:              blockV0.AttackTypeSignature[0].ProtocolIPv4[0].TTLMatch,
						TTLValue:              blockV0.AttackTypeSignature[0].ProtocolIPv4[0].TTLValue,
					}
				}
				if len(blockV0.AttackTypeSignature[0].ProtocolIPv6) > 0 {
					blockV1.AttackTypeSignature.ProtocolIPv6 = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv6{
						DestinationMatch: blockV0.AttackTypeSignature[0].ProtocolIPv6[0].DestinationMatch,
						DestinationValue: blockV0.AttackTypeSignature[0].ProtocolIPv6[0].DestinationValue,
						ExtensionHeaderDestinationOptionHomeAddressMatch: blockV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderDestinationOptionHomeAddressMatch,
						ExtensionHeaderDestinationOptionHomeAddressValue: blockV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderDestinationOptionHomeAddressValue,
						ExtensionHeaderDestinationOptionTypeMatch:        blockV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderDestinationOptionTypeMatch,
						ExtensionHeaderDestinationOptionTypeValue:        blockV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderDestinationOptionTypeValue,
						ExtensionHeaderRoutingHeaderTypeMatch:            blockV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderRoutingHeaderTypeMatch,
						ExtensionHeaderRoutingHeaderTypeValue:            blockV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderRoutingHeaderTypeValue,
						FlowLabelMatch:                                   blockV0.AttackTypeSignature[0].ProtocolIPv6[0].FlowLabelMatch,
						FlowLabelValue:                                   blockV0.AttackTypeSignature[0].ProtocolIPv6[0].FlowLabelValue,
						HopLimitMatch:                                    blockV0.AttackTypeSignature[0].ProtocolIPv6[0].HopLimitMatch,
						HopLimitValue:                                    blockV0.AttackTypeSignature[0].ProtocolIPv6[0].HopLimitValue,
						NextHeaderMatch:                                  blockV0.AttackTypeSignature[0].ProtocolIPv6[0].NextHeaderMatch,
						NextHeaderValue:                                  blockV0.AttackTypeSignature[0].ProtocolIPv6[0].NextHeaderValue,
						PayloadLengthMatch:                               blockV0.AttackTypeSignature[0].ProtocolIPv6[0].PayloadLengthMatch,
						PayloadLengthValue:                               blockV0.AttackTypeSignature[0].ProtocolIPv6[0].PayloadLengthValue,
						SourceMatch:                                      blockV0.AttackTypeSignature[0].ProtocolIPv6[0].SourceMatch,
						SourceValue:                                      blockV0.AttackTypeSignature[0].ProtocolIPv6[0].SourceValue,
						TrafficClassMatch:                                blockV0.AttackTypeSignature[0].ProtocolIPv6[0].TrafficClassMatch,
						TrafficClassValue:                                blockV0.AttackTypeSignature[0].ProtocolIPv6[0].TrafficClassValue,
					}
				}
				if len(blockV0.AttackTypeSignature[0].ProtocolTCP) > 0 {
					blockV1.AttackTypeSignature.ProtocolTCP = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCP{
						AckNumberMatch:        blockV0.AttackTypeSignature[0].ProtocolTCP[0].AckNumberMatch,
						AckNumberValue:        blockV0.AttackTypeSignature[0].ProtocolTCP[0].AckNumberValue,
						ChecksumValidateMatch: blockV0.AttackTypeSignature[0].ProtocolTCP[0].ChecksumValidateMatch,
						ChecksumValidateValue: blockV0.AttackTypeSignature[0].ProtocolTCP[0].ChecksumValidateValue,
						DataLengthMatch:       blockV0.AttackTypeSignature[0].ProtocolTCP[0].DataLengthMatch,
						DataLengthValue:       blockV0.AttackTypeSignature[0].ProtocolTCP[0].DataLengthValue,
						DestinationPortMatch:  blockV0.AttackTypeSignature[0].ProtocolTCP[0].DestinationPortMatch,
						DestinationPortValue:  blockV0.AttackTypeSignature[0].ProtocolTCP[0].DestinationPortValue,
						HeaderLengthMatch:     blockV0.AttackTypeSignature[0].ProtocolTCP[0].HeaderLengthMatch,
						HeaderLengthValue:     blockV0.AttackTypeSignature[0].ProtocolTCP[0].HeaderLengthValue,
						MssMatch:              blockV0.AttackTypeSignature[0].ProtocolTCP[0].MssMatch,
						MssValue:              blockV0.AttackTypeSignature[0].ProtocolTCP[0].MssValue,
						OptionMatch:           blockV0.AttackTypeSignature[0].ProtocolTCP[0].OptionMatch,
						OptionValue:           blockV0.AttackTypeSignature[0].ProtocolTCP[0].OptionValue,
						ReservedMatch:         blockV0.AttackTypeSignature[0].ProtocolTCP[0].ReservedMatch,
						ReservedValue:         blockV0.AttackTypeSignature[0].ProtocolTCP[0].ReservedValue,
						SequenceNumberMatch:   blockV0.AttackTypeSignature[0].ProtocolTCP[0].SequenceNumberMatch,
						SequenceNumberValue:   blockV0.AttackTypeSignature[0].ProtocolTCP[0].SequenceNumberValue,
						SourcePortMatch:       blockV0.AttackTypeSignature[0].ProtocolTCP[0].SourcePortMatch,
						SourcePortValue:       blockV0.AttackTypeSignature[0].ProtocolTCP[0].SourcePortValue,
						TCPFlags:              blockV0.AttackTypeSignature[0].ProtocolTCP[0].TCPFlags,
						UrgentPointerMatch:    blockV0.AttackTypeSignature[0].ProtocolTCP[0].UrgentPointerMatch,
						UrgentPointerValue:    blockV0.AttackTypeSignature[0].ProtocolTCP[0].UrgentPointerValue,
						WindowScaleMatch:      blockV0.AttackTypeSignature[0].ProtocolTCP[0].WindowScaleMatch,
						WindowScaleValue:      blockV0.AttackTypeSignature[0].ProtocolTCP[0].WindowScaleValue,
						WindowSizeMatch:       blockV0.AttackTypeSignature[0].ProtocolTCP[0].WindowSizeMatch,
						WindowSizeValue:       blockV0.AttackTypeSignature[0].ProtocolTCP[0].WindowSizeValue,
					}
				}
				if len(blockV0.AttackTypeSignature[0].ProtocolUDP) > 0 {
					blockV1.AttackTypeSignature.ProtocolUDP = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP{
						ChecksumValidateMatch: blockV0.AttackTypeSignature[0].ProtocolUDP[0].ChecksumValidateMatch,
						ChecksumValidateValue: blockV0.AttackTypeSignature[0].ProtocolUDP[0].ChecksumValidateValue,
						DataLengthMatch:       blockV0.AttackTypeSignature[0].ProtocolUDP[0].DataLengthMatch,
						DataLengthValue:       blockV0.AttackTypeSignature[0].ProtocolUDP[0].DataLengthValue,
						DestinationPortMatch:  blockV0.AttackTypeSignature[0].ProtocolUDP[0].DestinationPortMatch,
						DestinationPortValue:  blockV0.AttackTypeSignature[0].ProtocolUDP[0].DestinationPortValue,
						SourcePortMatch:       blockV0.AttackTypeSignature[0].ProtocolUDP[0].SourcePortMatch,
						SourcePortValue:       blockV0.AttackTypeSignature[0].ProtocolUDP[0].SourcePortValue,
					}
				}
			}
			dataV1.AttackTypeChain.Member = append(dataV1.AttackTypeChain.Member, blockV1)
		}
	}
	//nolint:lll
	if len(dataV0.AttackTypeSignature) > 0 {
		dataV1.AttackTypeSignature = &securityIdpCustomAttackBlockAttackTypeSignature{
			ProtocolBinding: dataV0.AttackTypeSignature[0].ProtocolBinding,
		}
		dataV1.AttackTypeSignature.Context = dataV0.AttackTypeSignature[0].Context
		dataV1.AttackTypeSignature.Direction = dataV0.AttackTypeSignature[0].Direction
		dataV1.AttackTypeSignature.Negate = dataV0.AttackTypeSignature[0].Negate
		dataV1.AttackTypeSignature.Pattern = dataV0.AttackTypeSignature[0].Pattern
		dataV1.AttackTypeSignature.PatternPcre = dataV0.AttackTypeSignature[0].PatternPcre
		dataV1.AttackTypeSignature.Regexp = dataV0.AttackTypeSignature[0].Regexp
		dataV1.AttackTypeSignature.Shellcode = dataV0.AttackTypeSignature[0].Shellcode
		if len(dataV0.AttackTypeSignature[0].ProtocolIcmp) > 0 {
			dataV1.AttackTypeSignature.ProtocolIcmp = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp{
				ChecksumValidateMatch: dataV0.AttackTypeSignature[0].ProtocolIcmp[0].ChecksumValidateMatch,
				ChecksumValidateValue: dataV0.AttackTypeSignature[0].ProtocolIcmp[0].ChecksumValidateValue,
				CodeMatch:             dataV0.AttackTypeSignature[0].ProtocolIcmp[0].CodeMatch,
				CodeValue:             dataV0.AttackTypeSignature[0].ProtocolIcmp[0].CodeValue,
				DataLengthMatch:       dataV0.AttackTypeSignature[0].ProtocolIcmp[0].DataLengthMatch,
				DataLengthValue:       dataV0.AttackTypeSignature[0].ProtocolIcmp[0].DataLengthValue,
				IdentificationMatch:   dataV0.AttackTypeSignature[0].ProtocolIcmp[0].IdentificationMatch,
				IdentificationValue:   dataV0.AttackTypeSignature[0].ProtocolIcmp[0].IdentificationValue,
				SequenceNumberMatch:   dataV0.AttackTypeSignature[0].ProtocolIcmp[0].SequenceNumberMatch,
				SequenceNumberValue:   dataV0.AttackTypeSignature[0].ProtocolIcmp[0].SequenceNumberValue,
				TypeMatch:             dataV0.AttackTypeSignature[0].ProtocolIcmp[0].TypeMatch,
				TypeValue:             dataV0.AttackTypeSignature[0].ProtocolIcmp[0].TypeValue,
			}
		}
		if len(dataV0.AttackTypeSignature[0].ProtocolIcmpv6) > 0 {
			dataV1.AttackTypeSignature.ProtocolIcmpv6 = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIcmp{
				ChecksumValidateMatch: dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].ChecksumValidateMatch,
				ChecksumValidateValue: dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].ChecksumValidateValue,
				CodeMatch:             dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].CodeMatch,
				CodeValue:             dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].CodeValue,
				DataLengthMatch:       dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].DataLengthMatch,
				DataLengthValue:       dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].DataLengthValue,
				IdentificationMatch:   dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].IdentificationMatch,
				IdentificationValue:   dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].IdentificationValue,
				SequenceNumberMatch:   dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].SequenceNumberMatch,
				SequenceNumberValue:   dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].SequenceNumberValue,
				TypeMatch:             dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].TypeMatch,
				TypeValue:             dataV0.AttackTypeSignature[0].ProtocolIcmpv6[0].TypeValue,
			}
		}
		if len(dataV0.AttackTypeSignature[0].ProtocolIPv4) > 0 {
			dataV1.AttackTypeSignature.ProtocolIPv4 = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv4{
				ChecksumValidateMatch: dataV0.AttackTypeSignature[0].ProtocolIPv4[0].ChecksumValidateMatch,
				ChecksumValidateValue: dataV0.AttackTypeSignature[0].ProtocolIPv4[0].ChecksumValidateValue,
				DestinationMatch:      dataV0.AttackTypeSignature[0].ProtocolIPv4[0].DestinationMatch,
				DestinationValue:      dataV0.AttackTypeSignature[0].ProtocolIPv4[0].DestinationValue,
				IdentificationMatch:   dataV0.AttackTypeSignature[0].ProtocolIPv4[0].IdentificationMatch,
				IdentificationValue:   dataV0.AttackTypeSignature[0].ProtocolIPv4[0].IdentificationValue,
				IhlMatch:              dataV0.AttackTypeSignature[0].ProtocolIPv4[0].IhlMatch,
				IhlValue:              dataV0.AttackTypeSignature[0].ProtocolIPv4[0].IhlValue,
				IPFlags:               dataV0.AttackTypeSignature[0].ProtocolIPv4[0].IPFlags,
				ProtocolMatch:         dataV0.AttackTypeSignature[0].ProtocolIPv4[0].ProtocolMatch,
				ProtocolValue:         dataV0.AttackTypeSignature[0].ProtocolIPv4[0].ProtocolValue,
				SourceMatch:           dataV0.AttackTypeSignature[0].ProtocolIPv4[0].SourceMatch,
				SourceValue:           dataV0.AttackTypeSignature[0].ProtocolIPv4[0].SourceValue,
				TosMatch:              dataV0.AttackTypeSignature[0].ProtocolIPv4[0].TosMatch,
				TosValue:              dataV0.AttackTypeSignature[0].ProtocolIPv4[0].TosValue,
				TotalLengthMatch:      dataV0.AttackTypeSignature[0].ProtocolIPv4[0].TotalLengthMatch,
				TotalLengthValue:      dataV0.AttackTypeSignature[0].ProtocolIPv4[0].TotalLengthValue,
				TTLMatch:              dataV0.AttackTypeSignature[0].ProtocolIPv4[0].TTLMatch,
				TTLValue:              dataV0.AttackTypeSignature[0].ProtocolIPv4[0].TTLValue,
			}
		}
		if len(dataV0.AttackTypeSignature[0].ProtocolIPv6) > 0 {
			dataV1.AttackTypeSignature.ProtocolIPv6 = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolIPv6{
				DestinationMatch: dataV0.AttackTypeSignature[0].ProtocolIPv6[0].DestinationMatch,
				DestinationValue: dataV0.AttackTypeSignature[0].ProtocolIPv6[0].DestinationValue,
				ExtensionHeaderDestinationOptionHomeAddressMatch: dataV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderDestinationOptionHomeAddressMatch,
				ExtensionHeaderDestinationOptionHomeAddressValue: dataV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderDestinationOptionHomeAddressValue,
				ExtensionHeaderDestinationOptionTypeMatch:        dataV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderDestinationOptionTypeMatch,
				ExtensionHeaderDestinationOptionTypeValue:        dataV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderDestinationOptionTypeValue,
				ExtensionHeaderRoutingHeaderTypeMatch:            dataV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderRoutingHeaderTypeMatch,
				ExtensionHeaderRoutingHeaderTypeValue:            dataV0.AttackTypeSignature[0].ProtocolIPv6[0].ExtensionHeaderRoutingHeaderTypeValue,
				FlowLabelMatch:                                   dataV0.AttackTypeSignature[0].ProtocolIPv6[0].FlowLabelMatch,
				FlowLabelValue:                                   dataV0.AttackTypeSignature[0].ProtocolIPv6[0].FlowLabelValue,
				HopLimitMatch:                                    dataV0.AttackTypeSignature[0].ProtocolIPv6[0].HopLimitMatch,
				HopLimitValue:                                    dataV0.AttackTypeSignature[0].ProtocolIPv6[0].HopLimitValue,
				NextHeaderMatch:                                  dataV0.AttackTypeSignature[0].ProtocolIPv6[0].NextHeaderMatch,
				NextHeaderValue:                                  dataV0.AttackTypeSignature[0].ProtocolIPv6[0].NextHeaderValue,
				PayloadLengthMatch:                               dataV0.AttackTypeSignature[0].ProtocolIPv6[0].PayloadLengthMatch,
				PayloadLengthValue:                               dataV0.AttackTypeSignature[0].ProtocolIPv6[0].PayloadLengthValue,
				SourceMatch:                                      dataV0.AttackTypeSignature[0].ProtocolIPv6[0].SourceMatch,
				SourceValue:                                      dataV0.AttackTypeSignature[0].ProtocolIPv6[0].SourceValue,
				TrafficClassMatch:                                dataV0.AttackTypeSignature[0].ProtocolIPv6[0].TrafficClassMatch,
				TrafficClassValue:                                dataV0.AttackTypeSignature[0].ProtocolIPv6[0].TrafficClassValue,
			}
		}
		if len(dataV0.AttackTypeSignature[0].ProtocolTCP) > 0 {
			dataV1.AttackTypeSignature.ProtocolTCP = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolTCP{
				AckNumberMatch:        dataV0.AttackTypeSignature[0].ProtocolTCP[0].AckNumberMatch,
				AckNumberValue:        dataV0.AttackTypeSignature[0].ProtocolTCP[0].AckNumberValue,
				ChecksumValidateMatch: dataV0.AttackTypeSignature[0].ProtocolTCP[0].ChecksumValidateMatch,
				ChecksumValidateValue: dataV0.AttackTypeSignature[0].ProtocolTCP[0].ChecksumValidateValue,
				DataLengthMatch:       dataV0.AttackTypeSignature[0].ProtocolTCP[0].DataLengthMatch,
				DataLengthValue:       dataV0.AttackTypeSignature[0].ProtocolTCP[0].DataLengthValue,
				DestinationPortMatch:  dataV0.AttackTypeSignature[0].ProtocolTCP[0].DestinationPortMatch,
				DestinationPortValue:  dataV0.AttackTypeSignature[0].ProtocolTCP[0].DestinationPortValue,
				HeaderLengthMatch:     dataV0.AttackTypeSignature[0].ProtocolTCP[0].HeaderLengthMatch,
				HeaderLengthValue:     dataV0.AttackTypeSignature[0].ProtocolTCP[0].HeaderLengthValue,
				MssMatch:              dataV0.AttackTypeSignature[0].ProtocolTCP[0].MssMatch,
				MssValue:              dataV0.AttackTypeSignature[0].ProtocolTCP[0].MssValue,
				OptionMatch:           dataV0.AttackTypeSignature[0].ProtocolTCP[0].OptionMatch,
				OptionValue:           dataV0.AttackTypeSignature[0].ProtocolTCP[0].OptionValue,
				ReservedMatch:         dataV0.AttackTypeSignature[0].ProtocolTCP[0].ReservedMatch,
				ReservedValue:         dataV0.AttackTypeSignature[0].ProtocolTCP[0].ReservedValue,
				SequenceNumberMatch:   dataV0.AttackTypeSignature[0].ProtocolTCP[0].SequenceNumberMatch,
				SequenceNumberValue:   dataV0.AttackTypeSignature[0].ProtocolTCP[0].SequenceNumberValue,
				SourcePortMatch:       dataV0.AttackTypeSignature[0].ProtocolTCP[0].SourcePortMatch,
				SourcePortValue:       dataV0.AttackTypeSignature[0].ProtocolTCP[0].SourcePortValue,
				TCPFlags:              dataV0.AttackTypeSignature[0].ProtocolTCP[0].TCPFlags,
				UrgentPointerMatch:    dataV0.AttackTypeSignature[0].ProtocolTCP[0].UrgentPointerMatch,
				UrgentPointerValue:    dataV0.AttackTypeSignature[0].ProtocolTCP[0].UrgentPointerValue,
				WindowScaleMatch:      dataV0.AttackTypeSignature[0].ProtocolTCP[0].WindowScaleMatch,
				WindowScaleValue:      dataV0.AttackTypeSignature[0].ProtocolTCP[0].WindowScaleValue,
				WindowSizeMatch:       dataV0.AttackTypeSignature[0].ProtocolTCP[0].WindowSizeMatch,
				WindowSizeValue:       dataV0.AttackTypeSignature[0].ProtocolTCP[0].WindowSizeValue,
			}
		}
		if len(dataV0.AttackTypeSignature[0].ProtocolUDP) > 0 {
			dataV1.AttackTypeSignature.ProtocolUDP = &securityIdpCustomAttackBlockAttackTypeSignatureBlockProtocolUDP{
				ChecksumValidateMatch: dataV0.AttackTypeSignature[0].ProtocolUDP[0].ChecksumValidateMatch,
				ChecksumValidateValue: dataV0.AttackTypeSignature[0].ProtocolUDP[0].ChecksumValidateValue,
				DataLengthMatch:       dataV0.AttackTypeSignature[0].ProtocolUDP[0].DataLengthMatch,
				DataLengthValue:       dataV0.AttackTypeSignature[0].ProtocolUDP[0].DataLengthValue,
				DestinationPortMatch:  dataV0.AttackTypeSignature[0].ProtocolUDP[0].DestinationPortMatch,
				DestinationPortValue:  dataV0.AttackTypeSignature[0].ProtocolUDP[0].DestinationPortValue,
				SourcePortMatch:       dataV0.AttackTypeSignature[0].ProtocolUDP[0].SourcePortMatch,
				SourcePortValue:       dataV0.AttackTypeSignature[0].ProtocolUDP[0].SourcePortValue,
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
