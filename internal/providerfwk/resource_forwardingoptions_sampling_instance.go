package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &forwardingoptionsSamplingInstance{}
	_ resource.ResourceWithConfigure      = &forwardingoptionsSamplingInstance{}
	_ resource.ResourceWithValidateConfig = &forwardingoptionsSamplingInstance{}
	_ resource.ResourceWithImportState    = &forwardingoptionsSamplingInstance{}
	_ resource.ResourceWithUpgradeState   = &forwardingoptionsSamplingInstance{}
)

type forwardingoptionsSamplingInstance struct {
	client *junos.Client
}

func newForwardingoptionsSamplingInstance() resource.Resource {
	return &forwardingoptionsSamplingInstance{}
}

func (rsc *forwardingoptionsSamplingInstance) typeName() string {
	return providerName + "_forwardingoptions_sampling_instance"
}

func (rsc *forwardingoptionsSamplingInstance) junosName() string {
	return "forwarding-options sampling instance"
}

func (rsc *forwardingoptionsSamplingInstance) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *forwardingoptionsSamplingInstance) Configure(
	ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*junos.Client)
	if !ok {
		unexpectedResourceConfigureType(ctx, req, resp)

		return
	}
	rsc.client = client
}

func (rsc *forwardingoptionsSamplingInstance) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Provides a " + rsc.junosName() + ".",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<name>_-_<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name for sampling instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for sampling instance if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable sampling instance.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"family_inet_input": schema.SingleNestedBlock{
				Description: "Declare `family inet input` configuration.",
				Attributes:  rsc.schemaInputAttributes(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"family_inet_output": schema.SingleNestedBlock{
				Description: "Declare `family inet output` configuration.",
				Attributes:  rsc.schemaFamilyInetOutputAttributes(),
				Blocks: map[string]schema.Block{
					"flow_server": schema.SetNestedBlock{
						Description: "For each hostname, configure sending traffic aggregates in cflowd format.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"hostname": schema.StringAttribute{
									Required:    true,
									Description: "Name of host collecting cflowd packets.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress(),
									},
								},
								"port": schema.Int64Attribute{
									Required:    true,
									Description: "UDP port number on host collecting cflowd packets (1..65535).",
									Validators: []validator.Int64{
										int64validator.Between(1, 65535),
									},
								},
								"aggregation_autonomous_system": schema.BoolAttribute{
									Optional:    true,
									Description: "Aggregate by autonomous system number.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"aggregation_destination_prefix": schema.BoolAttribute{
									Optional:    true,
									Description: "Aggregate by destination prefix.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"aggregation_protocol_port": schema.BoolAttribute{
									Optional:    true,
									Description: "Aggregate by protocol and port number.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"aggregation_source_destination_prefix": schema.BoolAttribute{
									Optional:    true,
									Description: "Aggregate by source and destination prefix.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"aggregation_source_destination_prefix_caida_compliant": schema.BoolAttribute{
									Optional:    true,
									Description: "Compatible with Caida record format for prefix aggregation (v8).",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"aggregation_source_prefix": schema.BoolAttribute{
									Optional:    true,
									Description: "Aggregate by source prefix.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"autonomous_system_type": schema.StringAttribute{
									Optional:    true,
									Description: "Type of autonomous system number to export.",
									Validators: []validator.String{
										stringvalidator.OneOf("origin", "peer"),
									},
								},
								"dscp": schema.Int64Attribute{
									Optional:    true,
									Description: "Numeric DSCP value in the range 0 to 63 (0..63).",
									Validators: []validator.Int64{
										int64validator.Between(0, 63),
									},
								},
								"forwarding_class": schema.StringAttribute{
									Optional:    true,
									Description: "Forwarding-class for exported jflow packets, applicable only for inline-jflow.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 64),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"local_dump": schema.BoolAttribute{
									Optional:    true,
									Description: "Dump cflowd records to log file before exporting.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"no_local_dump": schema.BoolAttribute{
									Optional:    true,
									Description: "Don't dump cflowd records to log file before exporting.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"routing_instance": schema.StringAttribute{
									Optional:    true,
									Description: "Name of routing instance on which flow collector is reachable.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 63),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
										stringvalidator.NoneOfCaseInsensitive(junos.DefaultW),
									},
								},
								"source_address": schema.StringAttribute{
									Optional:    true,
									Description: "Source IPv4 address for cflowd packets.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv4Only(),
									},
								},
								"version": schema.Int64Attribute{
									Optional:    true,
									Description: "Format of exported cflowd aggregates.",
									Validators: []validator.Int64{
										int64validator.OneOf(5, 8),
									},
								},
								"version9_template": schema.StringAttribute{
									Optional:    true,
									Description: "Template to export data in version 9 format.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"version_ipfix_template": schema.StringAttribute{
									Optional:    true,
									Description: "Template to export data in version ipfix format.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
							},
						},
					},
					"interface": schema.ListNestedBlock{
						Description: "For each name of interface, configure interfaces used to send monitored information.",
						NestedObject: schema.NestedBlockObject{
							Attributes: rsc.schemaOutputInterfaceAttributes(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"family_inet6_input": schema.SingleNestedBlock{
				Description: "Declare `family inet6 input` configuration.",
				Attributes:  rsc.schemaInputAttributes(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"family_inet6_output": schema.SingleNestedBlock{
				Description: "Declare `family inet6 output` configuration.",
				Attributes:  rsc.schemaFamilyInetOutputAttributes(),
				Blocks:      rsc.schemaOutputBlock(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"family_mpls_input": schema.SingleNestedBlock{
				Description: "Declare `family mpls input` configuration.",
				Attributes:  rsc.schemaInputAttributes(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"family_mpls_output": schema.SingleNestedBlock{
				Description: "Declare `family mpls output` configuration.",
				Attributes: map[string]schema.Attribute{
					"aggregate_export_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Interval of exporting aggregate accounting information (90..1800 seconds).",
						Validators: []validator.Int64{
							int64validator.Between(90, 1800),
						},
					},
					"flow_active_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Interval after which an active flow is exported (60..1800 seconds).",
						Validators: []validator.Int64{
							int64validator.Between(60, 1800),
						},
					},
					"flow_inactive_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Interval of inactivity that marks a flow inactive (15..1800 seconds).",
						Validators: []validator.Int64{
							int64validator.Between(15, 1800),
						},
					},
					"inline_jflow_export_rate": schema.Int64Attribute{
						Optional:    true,
						Description: "Inline processing of sampled packets with flow export rate of monitored packets in kpps (1..3200).",
						Validators: []validator.Int64{
							int64validator.Between(1, 3200),
						},
					},
					"inline_jflow_source_address": schema.StringAttribute{
						Optional:    true,
						Description: "Inline processing of sampled packets with address to use for generating monitored packets.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
				},
				Blocks: rsc.schemaOutputBlock(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"input": schema.SingleNestedBlock{
				Description: "Declare `input` configuration.",
				Attributes:  rsc.schemaInputAttributes(),
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
		},
	}
}

func (rsc *forwardingoptionsSamplingInstance) schemaInputAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"max_packets_per_second": schema.Int64Attribute{
			Optional:    true,
			Description: "Threshold of samples per second before dropping.",
			Validators: []validator.Int64{
				int64validator.Between(0, 65535),
			},
		},
		"maximum_packet_length": schema.Int64Attribute{
			Optional:    true,
			Description: "Maximum length of the sampled packet (0..9192 bytes).",
			Validators: []validator.Int64{
				int64validator.Between(0, 9192),
			},
		},
		"rate": schema.Int64Attribute{
			Optional:    true,
			Description: "Ratio of packets to be sampled (1 out of N) (1..16000000).",
			Validators: []validator.Int64{
				int64validator.Between(1, 16000000),
			},
		},
		"run_length": schema.Int64Attribute{
			Optional:    true,
			Description: "Number of samples after initial trigger (0..20).",
			Validators: []validator.Int64{
				int64validator.Between(0, 20),
			},
		},
	}
}

func (rsc *forwardingoptionsSamplingInstance) schemaFamilyInetOutputAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"aggregate_export_interval": schema.Int64Attribute{
			Optional:    true,
			Description: "Interval of exporting aggregate accounting information (90..1800 seconds).",
			Validators: []validator.Int64{
				int64validator.Between(90, 1800),
			},
		},
		"extension_service": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Define the customer specific sampling configuration.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				),
			},
		},
		"flow_active_timeout": schema.Int64Attribute{
			Optional:    true,
			Description: "Interval after which an active flow is exported (60..1800 seconds).",
			Validators: []validator.Int64{
				int64validator.Between(60, 1800),
			},
		},
		"flow_inactive_timeout": schema.Int64Attribute{
			Optional:    true,
			Description: "Interval of inactivity that marks a flow inactive (15..1800 seconds).",
			Validators: []validator.Int64{
				int64validator.Between(15, 1800),
			},
		},
		"inline_jflow_export_rate": schema.Int64Attribute{
			Optional:    true,
			Description: "Inline processing of sampled packets with flow export rate of monitored packets in kpps (1..3200).",
			Validators: []validator.Int64{
				int64validator.Between(1, 3200),
			},
		},
		"inline_jflow_source_address": schema.StringAttribute{
			Optional:    true,
			Description: "Inline processing of sampled packets with address to use for generating monitored packets.",
			Validators: []validator.String{
				tfvalidator.StringIPAddress(),
			},
		},
	}
}

func (rsc *forwardingoptionsSamplingInstance) schemaOutputBlock() map[string]schema.Block {
	return map[string]schema.Block{
		"flow_server": schema.SetNestedBlock{
			Description: "For each hostname, configure sending traffic aggregates in cflowd format.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"hostname": schema.StringAttribute{
						Required:    true,
						Description: "Name of host collecting cflowd packets.",
						Validators: []validator.String{
							tfvalidator.StringIPAddress(),
						},
					},
					"port": schema.Int64Attribute{
						Required:    true,
						Description: "UDP port number on host collecting cflowd packets (1..65535).",
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
					},
					"aggregation_autonomous_system": schema.BoolAttribute{
						Optional:    true,
						Description: "Aggregate by autonomous system number.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"aggregation_destination_prefix": schema.BoolAttribute{
						Optional:    true,
						Description: "Aggregate by destination prefix.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"aggregation_protocol_port": schema.BoolAttribute{
						Optional:    true,
						Description: "Aggregate by protocol and port number.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"aggregation_source_destination_prefix": schema.BoolAttribute{
						Optional:    true,
						Description: "Aggregate by source and destination prefix.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"aggregation_source_destination_prefix_caida_compliant": schema.BoolAttribute{
						Optional:    true,
						Description: "Compatible with Caida record format for prefix aggregation (v8).",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"aggregation_source_prefix": schema.BoolAttribute{
						Optional:    true,
						Description: "Aggregate by source prefix.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"autonomous_system_type": schema.StringAttribute{
						Optional:    true,
						Description: "Type of autonomous system number to export.",
						Validators: []validator.String{
							stringvalidator.OneOf("origin", "peer"),
						},
					},
					"dscp": schema.Int64Attribute{
						Optional:    true,
						Description: "Numeric DSCP value in the range 0 to 63 (0..63).",
						Validators: []validator.Int64{
							int64validator.Between(0, 63),
						},
					},
					"forwarding_class": schema.StringAttribute{
						Optional:    true,
						Description: "Forwarding-class for exported jflow packets, applicable only for inline-jflow.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 64),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"local_dump": schema.BoolAttribute{
						Optional:    true,
						Description: "Dump cflowd records to log file before exporting.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_local_dump": schema.BoolAttribute{
						Optional:    true,
						Description: "Don't dump cflowd records to log file before exporting.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"routing_instance": schema.StringAttribute{
						Optional:    true,
						Description: "Name of routing instance on which flow collector is reachable.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 63),
							tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							stringvalidator.NoneOfCaseInsensitive(junos.DefaultW),
						},
					},
					"source_address": schema.StringAttribute{
						Optional:    true,
						Description: "Source IPv4 address for cflowd packets",
						Validators: []validator.String{
							tfvalidator.StringIPAddress().IPv4Only(),
						},
					},
					"version9_template": schema.StringAttribute{
						Optional:    true,
						Description: "Template to export data in version 9 format.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 250),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"version_ipfix_template": schema.StringAttribute{
						Optional:    true,
						Description: "Template to export data in version ipfix format.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 250),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
			},
		},
		"interface": schema.ListNestedBlock{
			Description: "For each name of interface, configure interfaces used to send monitored information.",
			NestedObject: schema.NestedBlockObject{
				Attributes: rsc.schemaOutputInterfaceAttributes(),
			},
		},
	}
}

func (rsc *forwardingoptionsSamplingInstance) schemaOutputInterfaceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required:    true,
			Description: "Name of interface.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
			},
		},
		"engine_id": schema.Int64Attribute{
			Optional:    true,
			Description: "Identity (number) of this accounting interface (0..255).",
			Validators: []validator.Int64{
				int64validator.Between(0, 255),
			},
		},
		"engine_type": schema.Int64Attribute{
			Optional:    true,
			Description: "Type (number) of this accounting interface (0..255).",
			Validators: []validator.Int64{
				int64validator.Between(0, 255),
			},
		},
		"source_address": schema.StringAttribute{
			Optional:    true,
			Description: "Address to use for generating monitored packets.",
			Validators: []validator.String{
				tfvalidator.StringIPAddress(),
			},
		},
	}
}

type forwardingoptionsSamplingInstanceData struct {
	Disable           types.Bool                                              `tfsdk:"disable"`
	ID                types.String                                            `tfsdk:"id"`
	Name              types.String                                            `tfsdk:"name"`
	RoutingInstance   types.String                                            `tfsdk:"routing_instance"`
	FamilyInetInput   *forwardingoptionsSamplingInstanceInput                 `tfsdk:"family_inet_input"`
	FamilyInetOutput  *forwardingoptionsSamplingInstanceFamilyInetOutputData  `tfsdk:"family_inet_output"`
	FamilyInet6Input  *forwardingoptionsSamplingInstanceInput                 `tfsdk:"family_inet6_input"`
	FamilyInet6Output *forwardingoptionsSamplingInstanceFamilyInet6OutputData `tfsdk:"family_inet6_output"`
	FamilyMplsInput   *forwardingoptionsSamplingInstanceInput                 `tfsdk:"family_mpls_input"`
	FamilyMplsOutput  *forwardingoptionsSamplingInstanceFamilyMplsOutputData  `tfsdk:"family_mpls_output"`
	Input             *forwardingoptionsSamplingInstanceInput                 `tfsdk:"input"`
}

type forwardingoptionsSamplingInstanceConfig struct {
	Disable           types.Bool                                               `tfsdk:"disable"`
	ID                types.String                                             `tfsdk:"id"`
	Name              types.String                                             `tfsdk:"name"`
	RoutingInstance   types.String                                             `tfsdk:"routing_instance"`
	FamilyInetInput   *forwardingoptionsSamplingInstanceInput                  `tfsdk:"family_inet_input"`
	FamilyInetOutput  *forwardingoptionsSamplingInstanceFamilyInetOutputConfig `tfsdk:"family_inet_output"`
	FamilyInet6Input  *forwardingoptionsSamplingInstanceInput                  `tfsdk:"family_inet6_input"`
	FamilyInet6Output *forwardingoptionsSamplingInstanceFamilyInetOutputConfig `tfsdk:"family_inet6_output"`
	FamilyMplsInput   *forwardingoptionsSamplingInstanceInput                  `tfsdk:"family_mpls_input"`
	FamilyMplsOutput  *forwardingoptionsSamplingInstanceFamilyMplsOutputConfig `tfsdk:"family_mpls_output"`
	Input             *forwardingoptionsSamplingInstanceInput                  `tfsdk:"input"`
}

type forwardingoptionsSamplingInstanceInput struct {
	MaxPacketsPerSecond types.Int64 `tfsdk:"max_packets_per_second"`
	MaximumPacketLength types.Int64 `tfsdk:"maximum_packet_length"`
	Rate                types.Int64 `tfsdk:"rate"`
	RunLength           types.Int64 `tfsdk:"run_length"`
}

//nolint:lll
type forwardingoptionsSamplingInstanceFamilyInetOutputData struct {
	AggregateExportInterval  types.Int64                                                   `tfsdk:"aggregate_export_interval"`
	ExtensionService         []types.String                                                `tfsdk:"extension_service"`
	FlowActiveTimeout        types.Int64                                                   `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout      types.Int64                                                   `tfsdk:"flow_inactive_timeout"`
	InlineJflowExportRate    types.Int64                                                   `tfsdk:"inline_jflow_export_rate"`
	InlineJflowSourceAddress types.String                                                  `tfsdk:"inline_jflow_source_address"`
	FlowServer               []forwardingoptionsSamplingInstanceFamilyInetOutputFlowServer `tfsdk:"flow_server"`
	Interface                []forwardingoptionsSamplingInstanceOutputInterface            `tfsdk:"interface"`
}

type forwardingoptionsSamplingInstanceFamilyInetOutputConfig struct {
	AggregateExportInterval  types.Int64  `tfsdk:"aggregate_export_interval"`
	ExtensionService         types.List   `tfsdk:"extension_service"`
	FlowActiveTimeout        types.Int64  `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout      types.Int64  `tfsdk:"flow_inactive_timeout"`
	InlineJflowExportRate    types.Int64  `tfsdk:"inline_jflow_export_rate"`
	InlineJflowSourceAddress types.String `tfsdk:"inline_jflow_source_address"`
	FlowServer               types.Set    `tfsdk:"flow_server"`
	Interface                types.List   `tfsdk:"interface"`
}

//nolint:lll
type forwardingoptionsSamplingInstanceFamilyInetOutputFlowServer struct {
	AggregationAutonomousSystem                      types.Bool   `tfsdk:"aggregation_autonomous_system"`
	AggregationDestinationPrefix                     types.Bool   `tfsdk:"aggregation_destination_prefix"`
	AggregationProtocolPort                          types.Bool   `tfsdk:"aggregation_protocol_port"`
	AggregationSourceDestinationPrefix               types.Bool   `tfsdk:"aggregation_source_destination_prefix"`
	AggregationSourceDestinationPrefixCaidaCompliant types.Bool   `tfsdk:"aggregation_source_destination_prefix_caida_compliant"`
	AggregationSourcePrefix                          types.Bool   `tfsdk:"aggregation_source_prefix"`
	LocalDump                                        types.Bool   `tfsdk:"local_dump"`
	NoLocalDump                                      types.Bool   `tfsdk:"no_local_dump"`
	Hostname                                         types.String `tfsdk:"hostname"`
	Port                                             types.Int64  `tfsdk:"port"`
	AutonomousSystemType                             types.String `tfsdk:"autonomous_system_type"`
	Dscp                                             types.Int64  `tfsdk:"dscp"`
	ForwardingClass                                  types.String `tfsdk:"forwarding_class"`
	RoutingInstance                                  types.String `tfsdk:"routing_instance"`
	SourceAddress                                    types.String `tfsdk:"source_address"`
	Version                                          types.Int64  `tfsdk:"version"`
	Version9Template                                 types.String `tfsdk:"version9_template"`
	VersionIPFixTemplate                             types.String `tfsdk:"version_ipfix_template"`
}

type forwardingoptionsSamplingInstanceFamilyInet6OutputData struct {
	AggregateExportInterval  types.Int64                                         `tfsdk:"aggregate_export_interval"`
	ExtensionService         []types.String                                      `tfsdk:"extension_service"`
	FlowActiveTimeout        types.Int64                                         `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout      types.Int64                                         `tfsdk:"flow_inactive_timeout"`
	InlineJflowExportRate    types.Int64                                         `tfsdk:"inline_jflow_export_rate"`
	InlineJflowSourceAddress types.String                                        `tfsdk:"inline_jflow_source_address"`
	FlowServer               []forwardingoptionsSamplingInstanceOutputFlowServer `tfsdk:"flow_server"`
	Interface                []forwardingoptionsSamplingInstanceOutputInterface  `tfsdk:"interface"`
}

type forwardingoptionsSamplingInstanceFamilyMplsOutputData struct {
	AggregateExportInterval  types.Int64                                         `tfsdk:"aggregate_export_interval"`
	FlowActiveTimeout        types.Int64                                         `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout      types.Int64                                         `tfsdk:"flow_inactive_timeout"`
	InlineJflowExportRate    types.Int64                                         `tfsdk:"inline_jflow_export_rate"`
	InlineJflowSourceAddress types.String                                        `tfsdk:"inline_jflow_source_address"`
	FlowServer               []forwardingoptionsSamplingInstanceOutputFlowServer `tfsdk:"flow_server"`
	Interface                []forwardingoptionsSamplingInstanceOutputInterface  `tfsdk:"interface"`
}

type forwardingoptionsSamplingInstanceFamilyMplsOutputConfig struct {
	AggregateExportInterval  types.Int64  `tfsdk:"aggregate_export_interval"`
	FlowActiveTimeout        types.Int64  `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout      types.Int64  `tfsdk:"flow_inactive_timeout"`
	InlineJflowExportRate    types.Int64  `tfsdk:"inline_jflow_export_rate"`
	InlineJflowSourceAddress types.String `tfsdk:"inline_jflow_source_address"`
	FlowServer               types.Set    `tfsdk:"flow_server"`
	Interface                types.List   `tfsdk:"interface"`
}

//nolint:lll
type forwardingoptionsSamplingInstanceOutputFlowServer struct {
	AggregationAutonomousSystem                      types.Bool   `tfsdk:"aggregation_autonomous_system"`
	AggregationDestinationPrefix                     types.Bool   `tfsdk:"aggregation_destination_prefix"`
	AggregationProtocolPort                          types.Bool   `tfsdk:"aggregation_protocol_port"`
	AggregationSourceDestinationPrefix               types.Bool   `tfsdk:"aggregation_source_destination_prefix"`
	AggregationSourceDestinationPrefixCaidaCompliant types.Bool   `tfsdk:"aggregation_source_destination_prefix_caida_compliant"`
	AggregationSourcePrefix                          types.Bool   `tfsdk:"aggregation_source_prefix"`
	LocalDump                                        types.Bool   `tfsdk:"local_dump"`
	NoLocalDump                                      types.Bool   `tfsdk:"no_local_dump"`
	Hostname                                         types.String `tfsdk:"hostname"`
	Port                                             types.Int64  `tfsdk:"port"`
	AutonomousSystemType                             types.String `tfsdk:"autonomous_system_type"`
	Dscp                                             types.Int64  `tfsdk:"dscp"`
	ForwardingClass                                  types.String `tfsdk:"forwarding_class"`
	RoutingInstance                                  types.String `tfsdk:"routing_instance"`
	SourceAddress                                    types.String `tfsdk:"source_address"`
	Version9Template                                 types.String `tfsdk:"version9_template"`
	VersionIPFixTemplate                             types.String `tfsdk:"version_ipfix_template"`
}

type forwardingoptionsSamplingInstanceOutputInterface struct {
	Name          types.String `tfsdk:"name"`
	EngineID      types.Int64  `tfsdk:"engine_id"`
	EngineType    types.Int64  `tfsdk:"engine_type"`
	SourceAddress types.String `tfsdk:"source_address"`
}

func (block *forwardingoptionsSamplingInstanceInput) IsEmpty() bool {
	switch {
	case !block.MaxPacketsPerSecond.IsNull():
		return false
	case !block.MaximumPacketLength.IsNull():
		return false
	case !block.Rate.IsNull():
		return false
	case !block.RunLength.IsNull():
		return false
	default:
		return true
	}
}

func (rsc *forwardingoptionsSamplingInstance) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config forwardingoptionsSamplingInstanceConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Input != nil {
		if config.Input.IsEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("input").AtName("*"),
				"Missing Configuration Error",
				"input block is empty",
			)
		}
		if config.FamilyInetInput != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet_input").AtName("*"),
				"Conflict Configuration Error",
				"cannot set family_inet_input block if input block is used",
			)
		}
		if config.FamilyInet6Input != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet6_input").AtName("*"),
				"Conflict Configuration Error",
				"cannot set family_inet6_input block if input block is used",
			)
		}
		if config.FamilyMplsInput != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_mpls_input").AtName("*"),
				"Conflict Configuration Error",
				"cannot set family_mpls_input block if input block is used",
			)
		}
	} else if config.FamilyInetInput == nil &&
		config.FamilyInet6Input == nil &&
		config.FamilyMplsInput == nil {
		resp.Diagnostics.AddError(
			"Missing Configuration Error",
			"one of input, family_inet_input, family_inet6_input or family_mpls_input must be specified",
		)
	}
	if config.FamilyInetInput != nil {
		if config.FamilyInetInput.IsEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet_input").AtName("*"),
				"Missing Configuration Error",
				"family_inet_input block is empty",
			)
		}
	}
	if config.FamilyInet6Input != nil {
		if config.FamilyInet6Input.IsEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet6_input").AtName("*"),
				"Missing Configuration Error",
				"family_inet6_input block is empty",
			)
		}
	}
	if config.FamilyMplsInput != nil {
		if config.FamilyMplsInput.IsEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_mpls_input").AtName("*"),
				"Missing Configuration Error",
				"family_mpls_input block is empty",
			)
		}
	}

	if config.FamilyInetOutput == nil &&
		config.FamilyInet6Output == nil &&
		config.FamilyMplsOutput == nil {
		resp.Diagnostics.AddError(
			"Missing Configuration Error",
			"one of family_inet_output, family_inet6_output or family_mpls_output must be specified",
		)
	}

	if config.FamilyInetOutput != nil {
		if config.FamilyInetOutput.AggregateExportInterval.IsNull() &&
			config.FamilyInetOutput.ExtensionService.IsNull() &&
			config.FamilyInetOutput.FlowActiveTimeout.IsNull() &&
			config.FamilyInetOutput.FlowInactiveTimeout.IsNull() &&
			config.FamilyInetOutput.FlowServer.IsNull() &&
			config.FamilyInetOutput.InlineJflowExportRate.IsNull() &&
			config.FamilyInetOutput.InlineJflowSourceAddress.IsNull() &&
			config.FamilyInetOutput.Interface.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet_output").AtName("*"),
				"Missing Configuration Error",
				"family_inet_output block is empty",
			)
		}
		if config.FamilyInetOutput.InlineJflowSourceAddress.IsNull() {
			if !config.FamilyInetOutput.InlineJflowExportRate.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet_output").AtName("inline_jflow_export_rate"),
					"Missing Configuration Error",
					"inline_jflow_source_address must be specified with inline_jflow_export_rate in family_inet_output block",
				)
			}
		} else if config.FamilyInetOutput.FlowServer.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet_output").AtName("inline_jflow_source_address"),
				"Missing Configuration Error",
				"flow_server must be specified with inline_jflow_source_address in family_inet_output block",
			)
		}
	}
	if config.FamilyInet6Output != nil {
		if config.FamilyInet6Output.AggregateExportInterval.IsNull() &&
			config.FamilyInet6Output.ExtensionService.IsNull() &&
			config.FamilyInet6Output.FlowActiveTimeout.IsNull() &&
			config.FamilyInet6Output.FlowInactiveTimeout.IsNull() &&
			config.FamilyInet6Output.FlowServer.IsNull() &&
			config.FamilyInet6Output.InlineJflowExportRate.IsNull() &&
			config.FamilyInet6Output.InlineJflowSourceAddress.IsNull() &&
			config.FamilyInet6Output.Interface.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet6_output").AtName("*"),
				"Missing Configuration Error",
				"family_inet6_output block is empty",
			)
		}
		if config.FamilyInet6Output.InlineJflowSourceAddress.IsNull() {
			if !config.FamilyInet6Output.InlineJflowExportRate.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet6_output").AtName("inline_jflow_export_rate"),
					"Missing Configuration Error",
					"inline_jflow_source_address must be specified with inline_jflow_export_rate in family_inet6_output block",
				)
			}
		} else if config.FamilyInet6Output.FlowServer.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet6_output").AtName("inline_jflow_source_address"),
				"Missing Configuration Error",
				"flow_server must be specified with inline_jflow_source_address in family_inet6_output block",
			)
		}
	}
	if config.FamilyMplsOutput != nil {
		if config.FamilyMplsOutput.AggregateExportInterval.IsNull() &&
			config.FamilyMplsOutput.FlowActiveTimeout.IsNull() &&
			config.FamilyMplsOutput.FlowInactiveTimeout.IsNull() &&
			config.FamilyMplsOutput.FlowServer.IsNull() &&
			config.FamilyMplsOutput.InlineJflowExportRate.IsNull() &&
			config.FamilyMplsOutput.InlineJflowSourceAddress.IsNull() &&
			config.FamilyMplsOutput.Interface.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_mpls_output").AtName("*"),
				"Missing Configuration Error",
				"family_mpls_output block is empty",
			)
		}
		if config.FamilyMplsOutput.InlineJflowSourceAddress.IsNull() {
			if !config.FamilyMplsOutput.InlineJflowExportRate.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_mpls_output").AtName("inline_jflow_export_rate"),
					"Missing Configuration Error",
					"inline_jflow_source_address must be specified with inline_jflow_export_rate in family_mpls_output block",
				)
			}
		} else if config.FamilyMplsOutput.FlowServer.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_mpls_output").AtName("inline_jflow_source_address"),
				"Missing Configuration Error",
				"flow_server must be specified with inline_jflow_source_address in family_mpls_output block",
			)
		}
	}
}

func (rsc *forwardingoptionsSamplingInstance) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan forwardingoptionsSamplingInstanceData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Empty Name",
			"could not create "+rsc.junosName()+" with empty name",
		)

		return
	}

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}
	if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, v, junSess)
		if err != nil {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddError("Pre Check Error", err.Error())

			return
		}
		if !instanceExists {
			resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
			resp.Diagnostics.AddAttributeError(
				path.Root("routing_instance"),
				"Missing Configuration Error",
				fmt.Sprintf("routing instance %q doesn't exist", v),
			)

			return
		}
	}
	instanceExists, err := checkForwardingoptionsSamplingInstanceExists(
		ctx,
		plan.Name.ValueString(),
		plan.RoutingInstance.ValueString(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Pre Check Error", err.Error())

		return
	}
	if instanceExists {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError(
			"Duplicate Configuration Error",
			fmt.Sprintf(rsc.junosName()+" %q already exists", plan.Name.ValueString()),
		)

		return
	}

	if errPath, err := plan.set(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("create resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	instanceExists, err = checkForwardingoptionsSamplingInstanceExists(
		ctx,
		plan.Name.ValueString(),
		plan.RoutingInstance.ValueString(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError("Post Check Error", err.Error())

		return
	}
	if !instanceExists {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf(rsc.junosName()+" %q does not exists after commit "+
				"=> check your config", plan.Name.ValueString()),
		)

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *forwardingoptionsSamplingInstance) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data forwardingoptionsSamplingInstanceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	err = data.read(ctx, state.Name.ValueString(), state.RoutingInstance.ValueString(), junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}
	if data.ID.IsNull() {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *forwardingoptionsSamplingInstance) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state forwardingoptionsSamplingInstanceData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("update resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *forwardingoptionsSamplingInstance) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state forwardingoptionsSamplingInstanceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeDeleteAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	warns, err := junSess.CommitConf("delete resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}
}

func (rsc *forwardingoptionsSamplingInstance) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	var data forwardingoptionsSamplingInstanceData
	idSplit := strings.Split(req.ID, junos.IDSeparator)
	if len(idSplit) > 1 {
		if err := data.read(ctx, idSplit[0], idSplit[1], junSess); err != nil {
			resp.Diagnostics.AddError("Config Read Error", err.Error())

			return
		}
	} else {
		if err := data.read(ctx, idSplit[0], junos.DefaultW, junSess); err != nil {
			resp.Diagnostics.AddError("Config Read Error", err.Error())

			return
		}
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf("don't find "+rsc.junosName()+" with id %q "+
				"(id must be <name> or <name>"+junos.IDSeparator+"<routing_instance>)", req.ID),
		)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkForwardingoptionsSamplingInstanceExists(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) (
	_ bool, err error,
) {
	var showConfig string
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"forwarding-options sampling instance \"" + name + "\"" + junos.PipeDisplaySet)
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"forwarding-options sampling instance \"" + name + "\"" + junos.PipeDisplaySet)
	}
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *forwardingoptionsSamplingInstanceData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *forwardingoptionsSamplingInstanceData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set forwarding-options sampling instance \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix = junos.SetRoutingInstances + v +
			" forwarding-options sampling instance \"" + rscData.Name.ValueString() + "\" "
	}

	if rscData.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if rscData.FamilyInetInput != nil {
		blockSet := rscData.FamilyInetInput.configSet(setPrefix + "family inet input ")
		if len(blockSet) == 0 {
			return path.Root("family_inet_input").AtName("*"), fmt.Errorf("family_inet_input block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.FamilyInetOutput != nil {
		blockSet, pathErr, err := rscData.FamilyInetOutput.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		if len(blockSet) == 0 {
			return path.Root("family_inet_output").AtName("*"), fmt.Errorf("family_inet_output block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.FamilyInet6Input != nil {
		blockSet := rscData.FamilyInet6Input.configSet(setPrefix + "family inet6 input ")
		if len(blockSet) == 0 {
			return path.Root("family_inet6_input").AtName("*"), fmt.Errorf("family_inet6_input block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.FamilyInet6Output != nil {
		blockSet, pathErr, err := rscData.FamilyInet6Output.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		if len(blockSet) == 0 {
			return path.Root("family_inet6_output").AtName("*"), fmt.Errorf("family_inet6_output block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.FamilyMplsInput != nil {
		blockSet := rscData.FamilyMplsInput.configSet(setPrefix + "family mpls input ")
		if len(blockSet) == 0 {
			return path.Root("family_mpls_input").AtName("*"), fmt.Errorf("family_mpls_input block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.FamilyMplsOutput != nil {
		blockSet, pathErr, err := rscData.FamilyMplsOutput.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		if len(blockSet) == 0 {
			return path.Root("family_mpls_output").AtName("*"), fmt.Errorf("family_mpls_output block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.Input != nil {
		blockSet := rscData.Input.configSet(setPrefix + "input ")
		if len(blockSet) == 0 {
			return path.Root("input").AtName("*"), fmt.Errorf("input block is empty")
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *forwardingoptionsSamplingInstanceInput) configSet(
	setPrefix string,
) []string {
	configSet := make([]string, 0)

	if !block.MaxPacketsPerSecond.IsNull() {
		configSet = append(configSet, setPrefix+"max-packets-per-second "+
			utils.ConvI64toa(block.MaxPacketsPerSecond.ValueInt64()))
	}
	if !block.MaximumPacketLength.IsNull() {
		configSet = append(configSet, setPrefix+"maximum-packet-length "+
			utils.ConvI64toa(block.MaximumPacketLength.ValueInt64()))
	}
	if !block.Rate.IsNull() {
		configSet = append(configSet, setPrefix+"rate "+
			utils.ConvI64toa(block.Rate.ValueInt64()))
	}
	if !block.RunLength.IsNull() {
		configSet = append(configSet, setPrefix+"run-length "+
			utils.ConvI64toa(block.RunLength.ValueInt64()))
	}

	return configSet
}

func (block *forwardingoptionsSamplingInstanceFamilyInetOutputData) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "family inet output "

	if !block.AggregateExportInterval.IsNull() {
		configSet = append(configSet, setPrefix+"aggregate-export-interval "+
			utils.ConvI64toa(block.AggregateExportInterval.ValueInt64()))
	}
	for _, v := range block.ExtensionService {
		configSet = append(configSet, setPrefix+"extension-service \""+v.ValueString()+"\"")
	}
	if !block.FlowActiveTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"flow-active-timeout "+
			utils.ConvI64toa(block.FlowActiveTimeout.ValueInt64()))
	}
	if !block.FlowInactiveTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"flow-inactive-timeout "+
			utils.ConvI64toa(block.FlowInactiveTimeout.ValueInt64()))
	}
	flowServerHostname := make(map[string]struct{})
	for i, blockFlowServer := range block.FlowServer {
		hostname := blockFlowServer.Hostname.ValueString()
		if _, ok := flowServerHostname[hostname]; ok {
			return configSet, path.Root("family_inet_output").AtName("flow_server").AtListIndex(i).AtName("hostname"),
				fmt.Errorf("multiple blocks flow_server with the same hostname %q", hostname)
		}
		flowServerHostname[hostname] = struct{}{}
		setPrefixFlowServer := setPrefix + "flow-server " + hostname + " "
		configSet = append(configSet, setPrefixFlowServer+"port "+
			utils.ConvI64toa(blockFlowServer.Port.ValueInt64()))
		if blockFlowServer.AggregationAutonomousSystem.ValueBool() {
			configSet = append(configSet, setPrefixFlowServer+"aggregation autonomous-system")
		}
		if blockFlowServer.AggregationDestinationPrefix.ValueBool() {
			configSet = append(configSet, setPrefixFlowServer+"aggregation destination-prefix")
		}
		if blockFlowServer.AggregationProtocolPort.ValueBool() {
			configSet = append(configSet, setPrefixFlowServer+"aggregation protocol-port")
		}
		if blockFlowServer.AggregationSourceDestinationPrefix.ValueBool() {
			configSet = append(configSet, setPrefixFlowServer+"aggregation source-destination-prefix")
			if blockFlowServer.AggregationSourceDestinationPrefixCaidaCompliant.ValueBool() {
				configSet = append(configSet, setPrefixFlowServer+"aggregation source-destination-prefix caida-compliant")
			}
		} else if blockFlowServer.AggregationSourceDestinationPrefixCaidaCompliant.ValueBool() {
			return configSet,
				path.Root("family_inet_output").AtName("flow_server").AtListIndex(i).
					AtName("aggregation_source_destination_prefix_caida_compliant"),
				fmt.Errorf("aggregation_source_destination_prefix_caida_compliant = true "+
					"without aggregation_source_destination_prefix on flow-server %q", hostname)
		}
		if blockFlowServer.AggregationSourcePrefix.ValueBool() {
			configSet = append(configSet, setPrefixFlowServer+"aggregation source-prefix")
		}
		if v := blockFlowServer.AutonomousSystemType.ValueString(); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"autonomous-system-type "+v)
		}
		if !blockFlowServer.Dscp.IsNull() {
			configSet = append(configSet, setPrefixFlowServer+"dscp "+
				utils.ConvI64toa(blockFlowServer.Dscp.ValueInt64()))
		}
		if v := blockFlowServer.ForwardingClass.ValueString(); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"forwarding-class \""+v+"\"")
		}
		if blockFlowServer.LocalDump.ValueBool() {
			configSet = append(configSet, setPrefixFlowServer+"local-dump")
		}
		if blockFlowServer.NoLocalDump.ValueBool() {
			configSet = append(configSet, setPrefixFlowServer+"no-local-dump")
		}
		if v := blockFlowServer.RoutingInstance.ValueString(); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"routing-instance "+v)
		}
		if v := blockFlowServer.SourceAddress.ValueString(); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"source-address "+v)
		}
		if !blockFlowServer.Version.IsNull() {
			configSet = append(configSet, setPrefixFlowServer+"version "+
				utils.ConvI64toa(blockFlowServer.Version.ValueInt64()))
		}
		if v := blockFlowServer.Version9Template.ValueString(); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"version9 template \""+v+"\"")
		}
		if v := blockFlowServer.VersionIPFixTemplate.ValueString(); v != "" {
			configSet = append(configSet, setPrefixFlowServer+"version-ipfix template \""+v+"\"")
		}
	}
	if !block.InlineJflowExportRate.IsNull() {
		configSet = append(configSet, setPrefix+"inline-jflow flow-export-rate "+
			utils.ConvI64toa(block.InlineJflowExportRate.ValueInt64()))
	}
	if v := block.InlineJflowSourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"inline-jflow source-address "+v)
	}
	interfaceName := make(map[string]struct{})
	for i, blockInterface := range block.Interface {
		if _, ok := interfaceName[blockInterface.Name.ValueString()]; ok {
			return configSet, path.Root("family_inet_output").AtName("interface").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple blocks interface with the same name %q", blockInterface.Name.ValueString())
		}
		interfaceName[blockInterface.Name.ValueString()] = struct{}{}
		configSet = append(configSet, blockInterface.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsSamplingInstanceFamilyInet6OutputData) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "family inet6 output "

	if !block.AggregateExportInterval.IsNull() {
		configSet = append(configSet, setPrefix+"aggregate-export-interval "+
			utils.ConvI64toa(block.AggregateExportInterval.ValueInt64()))
	}
	for _, v := range block.ExtensionService {
		configSet = append(configSet, setPrefix+"extension-service \""+v.ValueString()+"\"")
	}
	if !block.FlowActiveTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"flow-active-timeout "+
			utils.ConvI64toa(block.FlowActiveTimeout.ValueInt64()))
	}
	if !block.FlowInactiveTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"flow-inactive-timeout "+
			utils.ConvI64toa(block.FlowInactiveTimeout.ValueInt64()))
	}
	flowServerHostname := make(map[string]struct{})
	for i, blockFlowServer := range block.FlowServer {
		if _, ok := flowServerHostname[blockFlowServer.Hostname.ValueString()]; ok {
			return configSet, path.Root("family_inet6_output").AtName("flow_server").AtListIndex(i).AtName("hostname"),
				fmt.Errorf("multiple blocks flow_server with the same hostname %q", blockFlowServer.Hostname.ValueString())
		}
		flowServerHostname[blockFlowServer.Hostname.ValueString()] = struct{}{}
		blockSet, pathErr, err := blockFlowServer.configSet(
			setPrefix,
			path.Root("family_inet6_output").AtName("flow_server").AtListIndex(i),
		)
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if !block.InlineJflowExportRate.IsNull() {
		configSet = append(configSet, setPrefix+"inline-jflow flow-export-rate "+
			utils.ConvI64toa(block.InlineJflowExportRate.ValueInt64()))
	}
	if v := block.InlineJflowSourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"inline-jflow source-address "+v)
	}
	interfaceName := make(map[string]struct{})
	for i, blockInterface := range block.Interface {
		if _, ok := interfaceName[blockInterface.Name.ValueString()]; ok {
			return configSet, path.Root("family_inet6_output").AtName("interface").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple blocks interface with the same name %q", blockInterface.Name.ValueString())
		}
		interfaceName[blockInterface.Name.ValueString()] = struct{}{}
		configSet = append(configSet, blockInterface.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsSamplingInstanceFamilyMplsOutputData) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "family mpls output "

	if !block.AggregateExportInterval.IsNull() {
		configSet = append(configSet, setPrefix+"aggregate-export-interval "+
			utils.ConvI64toa(block.AggregateExportInterval.ValueInt64()))
	}
	if !block.FlowActiveTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"flow-active-timeout "+
			utils.ConvI64toa(block.FlowActiveTimeout.ValueInt64()))
	}
	if !block.FlowInactiveTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"flow-inactive-timeout "+
			utils.ConvI64toa(block.FlowInactiveTimeout.ValueInt64()))
	}
	flowServerHostname := make(map[string]struct{})
	for i, blockFlowServer := range block.FlowServer {
		hostname := blockFlowServer.Hostname.ValueString()
		if _, ok := flowServerHostname[hostname]; ok {
			return configSet, path.Root("family_mpls_output").AtName("flow_server").AtListIndex(i).AtName("hostname"),
				fmt.Errorf("multiple blocks flow_server with the same hostname %q", hostname)
		}
		flowServerHostname[hostname] = struct{}{}
		blockSet, pathErr, err := blockFlowServer.configSet(
			setPrefix,
			path.Root("family_mpls_output").AtName("flow_server").AtListIndex(i),
		)
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if !block.InlineJflowExportRate.IsNull() {
		configSet = append(configSet, setPrefix+"inline-jflow flow-export-rate "+
			utils.ConvI64toa(block.InlineJflowExportRate.ValueInt64()))
	}
	if v := block.InlineJflowSourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"inline-jflow source-address "+v)
	}
	interfaceName := make(map[string]struct{})
	for i, blockInterface := range block.Interface {
		if _, ok := interfaceName[blockInterface.Name.ValueString()]; ok {
			return configSet, path.Root("family_mpls_output").AtName("interface").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple blocks interface with the same name %q", blockInterface.Name.ValueString())
		}
		interfaceName[blockInterface.Name.ValueString()] = struct{}{}
		configSet = append(configSet, blockInterface.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsSamplingInstanceOutputFlowServer) configSet(
	setPrefix string,
	pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "flow-server " + block.Hostname.ValueString() + " "

	configSet = append(configSet, setPrefix+"port "+
		utils.ConvI64toa(block.Port.ValueInt64()))
	if block.AggregationAutonomousSystem.ValueBool() {
		configSet = append(configSet, setPrefix+"aggregation autonomous-system")
	}
	if block.AggregationDestinationPrefix.ValueBool() {
		configSet = append(configSet, setPrefix+"aggregation destination-prefix")
	}
	if block.AggregationProtocolPort.ValueBool() {
		configSet = append(configSet, setPrefix+"aggregation protocol-port")
	}
	if block.AggregationSourceDestinationPrefix.ValueBool() {
		configSet = append(configSet, setPrefix+"aggregation source-destination-prefix")
		if block.AggregationSourceDestinationPrefixCaidaCompliant.ValueBool() {
			configSet = append(configSet, setPrefix+"aggregation source-destination-prefix caida-compliant")
		}
	} else if block.AggregationSourceDestinationPrefixCaidaCompliant.ValueBool() {
		return configSet,
			pathRoot.AtName("aggregation_source_destination_prefix_caida_compliant"),
			fmt.Errorf("aggregation_source_destination_prefix_caida_compliant = true "+
				"without aggregation_source_destination_prefix on flow-server %q", block.Hostname.ValueString())
	}
	if block.AggregationSourcePrefix.ValueBool() {
		configSet = append(configSet, setPrefix+"aggregation source-prefix")
	}
	if v := block.AutonomousSystemType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"autonomous-system-type "+v)
	}
	if !block.Dscp.IsNull() {
		configSet = append(configSet, setPrefix+"dscp "+
			utils.ConvI64toa(block.Dscp.ValueInt64()))
	}
	if v := block.ForwardingClass.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"forwarding-class \""+v+"\"")
	}
	if block.LocalDump.ValueBool() {
		configSet = append(configSet, setPrefix+"local-dump")
	}
	if block.NoLocalDump.ValueBool() {
		configSet = append(configSet, setPrefix+"no-local-dump")
	}
	if v := block.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-instance "+v)
	}
	if v := block.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}
	if v := block.Version9Template.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"version9 template \""+v+"\"")
	}
	if v := block.VersionIPFixTemplate.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"version-ipfix template \""+v+"\"")
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsSamplingInstanceOutputInterface) configSet(setPrefix string) []string {
	configSet := make([]string, 0)
	setPrefix += "interface " + block.Name.ValueString() + " "

	configSet = append(configSet, setPrefix)
	if !block.EngineID.IsNull() {
		configSet = append(configSet, setPrefix+"engine-id "+
			utils.ConvI64toa(block.EngineID.ValueInt64()))
	}
	if !block.EngineType.IsNull() {
		configSet = append(configSet, setPrefix+"engine-type "+
			utils.ConvI64toa(block.EngineType.ValueInt64()))
	}
	if v := block.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}

	return configSet
}

func (rscData *forwardingoptionsSamplingInstanceData) read(
	_ context.Context, name, routingInstance string, junSess *junos.Session,
) (
	err error,
) {
	var showConfig string
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showConfig, err = junSess.Command(junos.CmdShowConfig + junos.RoutingInstancesWS + routingInstance + " " +
			"forwarding-options sampling instance \"" + name + "\"" + junos.PipeDisplaySetRelative)
	} else {
		showConfig, err = junSess.Command(junos.CmdShowConfig +
			"forwarding-options sampling instance \"" + name + "\"" + junos.PipeDisplaySetRelative)
	}
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		if routingInstance == "" {
			rscData.RoutingInstance = types.StringValue(junos.DefaultW)
		} else {
			rscData.RoutingInstance = types.StringValue(routingInstance)
		}
		rscData.fillID()
		for _, item := range strings.Split(showConfig, "\n") {
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			switch {
			case itemTrim == junos.DisableW:
				rscData.Disable = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "family inet input "):
				if rscData.FamilyInetInput == nil {
					rscData.FamilyInetInput = &forwardingoptionsSamplingInstanceInput{}
				}
				if err := rscData.FamilyInetInput.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet6 input "):
				if rscData.FamilyInet6Input == nil {
					rscData.FamilyInet6Input = &forwardingoptionsSamplingInstanceInput{}
				}
				if err := rscData.FamilyInet6Input.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "family mpls input "):
				if rscData.FamilyMplsInput == nil {
					rscData.FamilyMplsInput = &forwardingoptionsSamplingInstanceInput{}
				}
				if err := rscData.FamilyMplsInput.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "input "):
				if rscData.Input == nil {
					rscData.Input = &forwardingoptionsSamplingInstanceInput{}
				}
				if err := rscData.Input.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet output "):
				if rscData.FamilyInetOutput == nil {
					rscData.FamilyInetOutput = &forwardingoptionsSamplingInstanceFamilyInetOutputData{}
				}
				if err := rscData.FamilyInetOutput.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet6 output "):
				if rscData.FamilyInet6Output == nil {
					rscData.FamilyInet6Output = &forwardingoptionsSamplingInstanceFamilyInet6OutputData{}
				}
				if err := rscData.FamilyInet6Output.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "family mpls output "):
				if rscData.FamilyMplsOutput == nil {
					rscData.FamilyMplsOutput = &forwardingoptionsSamplingInstanceFamilyMplsOutputData{}
				}
				if err := rscData.FamilyMplsOutput.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *forwardingoptionsSamplingInstanceInput) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "max-packets-per-second "):
		block.MaxPacketsPerSecond, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "maximum-packet-length "):
		block.MaximumPacketLength, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "rate "):
		block.Rate, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "run-length "):
		block.RunLength, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (block *forwardingoptionsSamplingInstanceFamilyInetOutputData) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "aggregate-export-interval "):
		block.AggregateExportInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "extension-service "):
		block.ExtensionService = append(block.ExtensionService, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "flow-active-timeout "):
		block.FlowActiveTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "flow-inactive-timeout "):
		block.FlowInactiveTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "flow-server "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var flowServer forwardingoptionsSamplingInstanceFamilyInetOutputFlowServer
		block.FlowServer, flowServer = tfdata.ExtractBlockWithTFTypesString(
			block.FlowServer, "Hostname", itemTrimFields[0],
		)
		flowServer.Hostname = types.StringValue(itemTrimFields[0])
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		switch {
		case balt.CutPrefixInString(&itemTrim, "port "):
			flowServer.Port, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case itemTrim == "aggregation autonomous-system":
			flowServer.AggregationAutonomousSystem = types.BoolValue(true)
		case itemTrim == "aggregation destination-prefix":
			flowServer.AggregationDestinationPrefix = types.BoolValue(true)
		case itemTrim == "aggregation protocol-port":
			flowServer.AggregationProtocolPort = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "aggregation source-destination-prefix"):
			flowServer.AggregationSourceDestinationPrefix = types.BoolValue(true)
			if itemTrim == " caida-compliant" {
				flowServer.AggregationSourceDestinationPrefixCaidaCompliant = types.BoolValue(true)
			}
		case itemTrim == "aggregation source-prefix":
			flowServer.AggregationSourcePrefix = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "autonomous-system-type "):
			flowServer.AutonomousSystemType = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "dscp "):
			flowServer.Dscp, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "forwarding-class "):
			flowServer.ForwardingClass = types.StringValue(strings.Trim(itemTrim, "\""))
		case itemTrim == "local-dump":
			flowServer.LocalDump = types.BoolValue(true)
		case itemTrim == "no-local-dump":
			flowServer.NoLocalDump = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "routing-instance "):
			flowServer.RoutingInstance = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "source-address "):
			flowServer.SourceAddress = types.StringValue(itemTrim)
		case balt.CutPrefixInString(&itemTrim, "version "):
			flowServer.Version, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case balt.CutPrefixInString(&itemTrim, "version9 template "):
			flowServer.Version9Template = types.StringValue(strings.Trim(itemTrim, "\""))
		case balt.CutPrefixInString(&itemTrim, "version-ipfix template "):
			flowServer.VersionIPFixTemplate = types.StringValue(strings.Trim(itemTrim, "\""))
		}
		block.FlowServer = append(block.FlowServer, flowServer)
	case balt.CutPrefixInString(&itemTrim, "inline-jflow flow-export-rate "):
		block.InlineJflowExportRate, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "inline-jflow source-address "):
		block.InlineJflowSourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "interface "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var iFace forwardingoptionsSamplingInstanceOutputInterface
		block.Interface, iFace = tfdata.ExtractBlockWithTFTypesString(
			block.Interface, "Name", itemTrimFields[0],
		)
		iFace.Name = types.StringValue(itemTrimFields[0])
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		if err := iFace.read(itemTrim); err != nil {
			return err
		}
		block.Interface = append(block.Interface, iFace)
	}

	return nil
}

func (block *forwardingoptionsSamplingInstanceFamilyInet6OutputData) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "aggregate-export-interval "):
		block.AggregateExportInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "extension-service "):
		block.ExtensionService = append(block.ExtensionService, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "flow-active-timeout "):
		block.FlowActiveTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "flow-inactive-timeout "):
		block.FlowInactiveTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "flow-server "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var flowServer forwardingoptionsSamplingInstanceOutputFlowServer
		block.FlowServer, flowServer = tfdata.ExtractBlockWithTFTypesString(
			block.FlowServer, "Hostname", itemTrimFields[0],
		)
		flowServer.Hostname = types.StringValue(itemTrimFields[0])
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		if err := flowServer.read(itemTrim); err != nil {
			return err
		}
		block.FlowServer = append(block.FlowServer, flowServer)
	case balt.CutPrefixInString(&itemTrim, "inline-jflow flow-export-rate "):
		block.InlineJflowExportRate, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "inline-jflow source-address "):
		block.InlineJflowSourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "interface "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var iFace forwardingoptionsSamplingInstanceOutputInterface
		block.Interface, iFace = tfdata.ExtractBlockWithTFTypesString(
			block.Interface, "Name", itemTrimFields[0],
		)
		iFace.Name = types.StringValue(itemTrimFields[0])
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		if err := iFace.read(itemTrim); err != nil {
			return err
		}
		block.Interface = append(block.Interface, iFace)
	}

	return nil
}

func (block *forwardingoptionsSamplingInstanceFamilyMplsOutputData) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "aggregate-export-interval "):
		block.AggregateExportInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "flow-active-timeout "):
		block.FlowActiveTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "flow-inactive-timeout "):
		block.FlowInactiveTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "flow-server "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var flowServer forwardingoptionsSamplingInstanceOutputFlowServer
		block.FlowServer, flowServer = tfdata.ExtractBlockWithTFTypesString(
			block.FlowServer, "Hostname", itemTrimFields[0],
		)
		flowServer.Hostname = types.StringValue(itemTrimFields[0])
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		if err := flowServer.read(itemTrim); err != nil {
			return err
		}
		block.FlowServer = append(block.FlowServer, flowServer)
	case balt.CutPrefixInString(&itemTrim, "inline-jflow flow-export-rate "):
		block.InlineJflowExportRate, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "inline-jflow source-address "):
		block.InlineJflowSourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "interface "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var iFace forwardingoptionsSamplingInstanceOutputInterface
		block.Interface, iFace = tfdata.ExtractBlockWithTFTypesString(
			block.Interface, "Name", itemTrimFields[0],
		)
		iFace.Name = types.StringValue(itemTrimFields[0])
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		if err := iFace.read(itemTrim); err != nil {
			return err
		}
		block.Interface = append(block.Interface, iFace)
	}

	return nil
}

func (block *forwardingoptionsSamplingInstanceOutputFlowServer) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "port "):
		block.Port, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "aggregation autonomous-system":
		block.AggregationAutonomousSystem = types.BoolValue(true)
	case itemTrim == "aggregation destination-prefix":
		block.AggregationDestinationPrefix = types.BoolValue(true)
	case itemTrim == "aggregation protocol-port":
		block.AggregationProtocolPort = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "aggregation source-destination-prefix"):
		block.AggregationSourceDestinationPrefix = types.BoolValue(true)
		if itemTrim == " caida-compliant" {
			block.AggregationSourceDestinationPrefixCaidaCompliant = types.BoolValue(true)
		}
	case itemTrim == "aggregation source-prefix":
		block.AggregationSourcePrefix = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "autonomous-system-type "):
		block.AutonomousSystemType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "dscp "):
		block.Dscp, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "forwarding-class "):
		block.ForwardingClass = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "local-dump":
		block.LocalDump = types.BoolValue(true)
	case itemTrim == "no-local-dump":
		block.NoLocalDump = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "routing-instance "):
		block.RoutingInstance = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		block.SourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "version9 template "):
		block.Version9Template = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "version-ipfix template "):
		block.VersionIPFixTemplate = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return nil
}

func (block *forwardingoptionsSamplingInstanceOutputInterface) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "engine-id "):
		block.EngineID, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "engine-type "):
		block.EngineType, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		block.SourceAddress = types.StringValue(itemTrim)
	}

	return nil
}

func (rscData *forwardingoptionsSamplingInstanceData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := make([]string, 1)
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		configSet[0] = junos.DelRoutingInstances + v +
			" forwarding-options sampling instance \"" + rscData.Name.ValueString() + "\""
	} else {
		configSet[0] = "delete " +
			" forwarding-options sampling instance \"" + rscData.Name.ValueString() + "\""
	}

	return junSess.ConfigSet(configSet)
}
