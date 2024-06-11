package providerfwk

import (
	"context"
	"errors"
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
	_ resource.Resource                   = &forwardingoptionsSampling{}
	_ resource.ResourceWithConfigure      = &forwardingoptionsSampling{}
	_ resource.ResourceWithValidateConfig = &forwardingoptionsSampling{}
	_ resource.ResourceWithImportState    = &forwardingoptionsSampling{}
)

type forwardingoptionsSampling struct {
	client *junos.Client
}

func newForwardingoptionsSamplingResource() resource.Resource {
	return &forwardingoptionsSampling{}
}

func (rsc *forwardingoptionsSampling) typeName() string {
	return providerName + "_forwardingoptions_sampling"
}

func (rsc *forwardingoptionsSampling) junosName() string {
	return "forwarding-options sampling"
}

func (rsc *forwardingoptionsSampling) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *forwardingoptionsSampling) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *forwardingoptionsSampling) Configure(
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

func (rsc *forwardingoptionsSampling) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance if not root level.",
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
				Description: "Disable global sampling instance.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"pre_rewrite_tos": schema.BoolAttribute{
				Optional:    true,
				Description: "Sample the packet retaining tos value before normalization.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"sample_once": schema.BoolAttribute{
				Optional:    true,
				Description: "Sample the packet for active-monitoring only once.",
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
					"file": schema.SingleNestedBlock{
						Description: "Configure parameters for dumping sampled packets.",
						Attributes: map[string]schema.Attribute{
							"filename": schema.StringAttribute{
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "Name of file to contain sampled packet dumps.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
									tfvalidator.StringSpaceExclusion(),
									tfvalidator.StringRuneExclusion('/', '%'),
								},
							},
							"disable": schema.BoolAttribute{
								Optional:    true,
								Description: "Disable sampled packet dumps.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"files": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum number of sampled packet dump files (2..10000).",
								Validators: []validator.Int64{
									int64validator.Between(2, 10000),
								},
							},
							"size": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum sample dump file size (1024..104857600).",
								Validators: []validator.Int64{
									int64validator.Between(1024, 104857600),
								},
							},
							"stamp": schema.BoolAttribute{
								Optional:    true,
								Description: "Timestamp every packet in the dump.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_stamp": schema.BoolAttribute{
								Optional:    true,
								Description: "Don't timestamp every packet in the dump.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"world_readable": schema.BoolAttribute{
								Optional:    true,
								Description: "Allow any user to read the sampled dump.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_world_readable": schema.BoolAttribute{
								Optional:    true,
								Description: "Don't allow any user to read the sampled dump.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
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
										int64validator.OneOf(5, 8, 500),
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

func (rsc *forwardingoptionsSampling) schemaInputAttributes() map[string]schema.Attribute {
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

func (rsc *forwardingoptionsSampling) schemaFamilyInetOutputAttributes() map[string]schema.Attribute {
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

func (rsc *forwardingoptionsSampling) schemaOutputBlock() map[string]schema.Block {
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

func (rsc *forwardingoptionsSampling) schemaOutputInterfaceAttributes() map[string]schema.Attribute {
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

type forwardingoptionsSamplingData struct {
	ID                types.String                                     `tfsdk:"id"`
	RoutingInstance   types.String                                     `tfsdk:"routing_instance"`
	Disable           types.Bool                                       `tfsdk:"disable"`
	PreRewriteTos     types.Bool                                       `tfsdk:"pre_rewrite_tos"`
	SampleOnce        types.Bool                                       `tfsdk:"sample_once"`
	FamilyInetInput   *forwardingoptionsSamplingBlockInput             `tfsdk:"family_inet_input"`
	FamilyInetOutput  *forwardingoptionsSamplingBlockFamilyInetOutput  `tfsdk:"family_inet_output"`
	FamilyInet6Input  *forwardingoptionsSamplingBlockInput             `tfsdk:"family_inet6_input"`
	FamilyInet6Output *forwardingoptionsSamplingBlockFamilyInet6Output `tfsdk:"family_inet6_output"`
	FamilyMplsInput   *forwardingoptionsSamplingBlockInput             `tfsdk:"family_mpls_input"`
	FamilyMplsOutput  *forwardingoptionsSamplingBlockFamilyMplsOutput  `tfsdk:"family_mpls_output"`
	Input             *forwardingoptionsSamplingBlockInput             `tfsdk:"input"`
}

type forwardingoptionsSamplingConfig struct {
	ID                types.String                                           `tfsdk:"id"`
	RoutingInstance   types.String                                           `tfsdk:"routing_instance"`
	Disable           types.Bool                                             `tfsdk:"disable"`
	PreRewriteTos     types.Bool                                             `tfsdk:"pre_rewrite_tos"`
	SampleOnce        types.Bool                                             `tfsdk:"sample_once"`
	FamilyInetInput   *forwardingoptionsSamplingBlockInput                   `tfsdk:"family_inet_input"`
	FamilyInetOutput  *forwardingoptionsSamplingBlockFamilyInetOutputConfig  `tfsdk:"family_inet_output"`
	FamilyInet6Input  *forwardingoptionsSamplingBlockInput                   `tfsdk:"family_inet6_input"`
	FamilyInet6Output *forwardingoptionsSamplingBlockFamilyInet6OutputConfig `tfsdk:"family_inet6_output"`
	FamilyMplsInput   *forwardingoptionsSamplingBlockInput                   `tfsdk:"family_mpls_input"`
	FamilyMplsOutput  *forwardingoptionsSamplingBlockFamilyMplsOutputConfig  `tfsdk:"family_mpls_output"`
	Input             *forwardingoptionsSamplingBlockInput                   `tfsdk:"input"`
}

type forwardingoptionsSamplingBlockInput struct {
	MaxPacketsPerSecond types.Int64 `tfsdk:"max_packets_per_second"`
	MaximumPacketLength types.Int64 `tfsdk:"maximum_packet_length"`
	Rate                types.Int64 `tfsdk:"rate"`
	RunLength           types.Int64 `tfsdk:"run_length"`
}

func (block *forwardingoptionsSamplingBlockInput) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (block *forwardingoptionsSamplingBlockInput) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

//nolint:lll
type forwardingoptionsSamplingBlockFamilyInetOutput struct {
	AggregateExportInterval  types.Int64                                                     `tfsdk:"aggregate_export_interval"`
	ExtensionService         []types.String                                                  `tfsdk:"extension_service"`
	File                     *forwardingoptionsSamplingBlockFamilyInetOutputBlockFile        `tfsdk:"file"`
	FlowActiveTimeout        types.Int64                                                     `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout      types.Int64                                                     `tfsdk:"flow_inactive_timeout"`
	InlineJflowExportRate    types.Int64                                                     `tfsdk:"inline_jflow_export_rate"`
	InlineJflowSourceAddress types.String                                                    `tfsdk:"inline_jflow_source_address"`
	FlowServer               []forwardingoptionsSamplingBlockFamilyInetOutputBlockFlowServer `tfsdk:"flow_server"`
	Interface                []forwardingoptionsSamplingBlockOutputBlockInterface            `tfsdk:"interface"`
}

type forwardingoptionsSamplingBlockFamilyInetOutputConfig struct {
	AggregateExportInterval  types.Int64                                              `tfsdk:"aggregate_export_interval"`
	ExtensionService         types.List                                               `tfsdk:"extension_service"`
	File                     *forwardingoptionsSamplingBlockFamilyInetOutputBlockFile `tfsdk:"file"`
	FlowActiveTimeout        types.Int64                                              `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout      types.Int64                                              `tfsdk:"flow_inactive_timeout"`
	InlineJflowExportRate    types.Int64                                              `tfsdk:"inline_jflow_export_rate"`
	InlineJflowSourceAddress types.String                                             `tfsdk:"inline_jflow_source_address"`
	FlowServer               types.Set                                                `tfsdk:"flow_server"`
	Interface                types.List                                               `tfsdk:"interface"`
}

func (block *forwardingoptionsSamplingBlockFamilyInetOutputConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type forwardingoptionsSamplingBlockFamilyInetOutputBlockFile struct {
	Filename        types.String `tfsdk:"filename"`
	Disable         types.Bool   `tfsdk:"disable"`
	Files           types.Int64  `tfsdk:"files"`
	Size            types.Int64  `tfsdk:"size"`
	Stamp           types.Bool   `tfsdk:"stamp"`
	NoStamp         types.Bool   `tfsdk:"no_stamp"`
	WorldReadable   types.Bool   `tfsdk:"world_readable"`
	NoWorldReadable types.Bool   `tfsdk:"no_world_readable"`
}

//nolint:lll
type forwardingoptionsSamplingBlockFamilyInetOutputBlockFlowServer struct {
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
}

type forwardingoptionsSamplingBlockFamilyInet6Output struct {
	AggregateExportInterval  types.Int64                                           `tfsdk:"aggregate_export_interval"`
	ExtensionService         []types.String                                        `tfsdk:"extension_service"`
	FlowActiveTimeout        types.Int64                                           `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout      types.Int64                                           `tfsdk:"flow_inactive_timeout"`
	InlineJflowExportRate    types.Int64                                           `tfsdk:"inline_jflow_export_rate"`
	InlineJflowSourceAddress types.String                                          `tfsdk:"inline_jflow_source_address"`
	FlowServer               []forwardingoptionsSamplingBlockOutputBlockFlowServer `tfsdk:"flow_server"`
	Interface                []forwardingoptionsSamplingBlockOutputBlockInterface  `tfsdk:"interface"`
}

type forwardingoptionsSamplingBlockFamilyInet6OutputConfig struct {
	AggregateExportInterval  types.Int64  `tfsdk:"aggregate_export_interval"`
	ExtensionService         types.List   `tfsdk:"extension_service"`
	FlowActiveTimeout        types.Int64  `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout      types.Int64  `tfsdk:"flow_inactive_timeout"`
	InlineJflowExportRate    types.Int64  `tfsdk:"inline_jflow_export_rate"`
	InlineJflowSourceAddress types.String `tfsdk:"inline_jflow_source_address"`
	FlowServer               types.Set    `tfsdk:"flow_server"`
	Interface                types.List   `tfsdk:"interface"`
}

func (block *forwardingoptionsSamplingBlockFamilyInet6OutputConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type forwardingoptionsSamplingBlockFamilyMplsOutput struct {
	AggregateExportInterval types.Int64                                           `tfsdk:"aggregate_export_interval"`
	FlowActiveTimeout       types.Int64                                           `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout     types.Int64                                           `tfsdk:"flow_inactive_timeout"`
	FlowServer              []forwardingoptionsSamplingBlockOutputBlockFlowServer `tfsdk:"flow_server"`
	Interface               []forwardingoptionsSamplingBlockOutputBlockInterface  `tfsdk:"interface"`
}

type forwardingoptionsSamplingBlockFamilyMplsOutputConfig struct {
	AggregateExportInterval types.Int64 `tfsdk:"aggregate_export_interval"`
	FlowActiveTimeout       types.Int64 `tfsdk:"flow_active_timeout"`
	FlowInactiveTimeout     types.Int64 `tfsdk:"flow_inactive_timeout"`
	FlowServer              types.Set   `tfsdk:"flow_server"`
	Interface               types.List  `tfsdk:"interface"`
}

func (block *forwardingoptionsSamplingBlockFamilyMplsOutputConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

//nolint:lll
type forwardingoptionsSamplingBlockOutputBlockFlowServer struct {
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
}

type forwardingoptionsSamplingBlockOutputBlockInterface struct {
	Name          types.String `tfsdk:"name"`
	EngineID      types.Int64  `tfsdk:"engine_id"`
	EngineType    types.Int64  `tfsdk:"engine_type"`
	SourceAddress types.String `tfsdk:"source_address"`
}

func (rsc *forwardingoptionsSampling) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config forwardingoptionsSamplingConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Input != nil {
		if config.Input.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("input").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"input block is empty",
			)
		} else if config.Input.hasKnownValue() {
			if config.FamilyInetInput != nil && config.FamilyInetInput.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet_input").AtName("*"),
					tfdiag.ConflictConfigErrSummary,
					"cannot set family_inet_input block if input block is used",
				)
			}
			if config.FamilyInet6Input != nil && config.FamilyInet6Input.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet6_input").AtName("*"),
					tfdiag.ConflictConfigErrSummary,
					"cannot set family_inet6_input block if input block is used",
				)
			}
			if config.FamilyMplsInput != nil && config.FamilyMplsInput.hasKnownValue() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_mpls_input").AtName("*"),
					tfdiag.ConflictConfigErrSummary,
					"cannot set family_mpls_input block if input block is used",
				)
			}
		}
	}

	if config.FamilyInetInput != nil {
		if config.FamilyInetInput.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet_input").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"family_inet_input block is empty",
			)
		}
	}
	if config.FamilyInet6Input != nil {
		if config.FamilyInet6Input.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet6_input").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"family_inet6_input block is empty",
			)
		}
	}
	if config.FamilyMplsInput != nil {
		if config.FamilyMplsInput.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_mpls_input").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"family_mpls_input block is empty",
			)
		}
	}

	if config.FamilyInetOutput == nil &&
		config.FamilyInet6Output == nil &&
		config.FamilyMplsOutput == nil {
		resp.Diagnostics.AddError(
			tfdiag.MissingConfigErrSummary,
			"one of family_inet_output, family_inet6_output or family_mpls_output must be specified",
		)
	}

	if config.FamilyInetOutput != nil {
		if config.Input == nil &&
			config.FamilyInetInput == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet_output").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"one of input or family_inet_input must be specified with family_inet_output",
			)
		}
		if config.FamilyInetOutput.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet_output").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"family_inet_output block is empty",
			)
		}
		if config.FamilyInetOutput.File != nil {
			if config.FamilyInetOutput.File.Filename.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet_output").AtName("file"),
					tfdiag.MissingConfigErrSummary,
					"filename must be specified in family_inet_output.file block",
				)
			}
			if !config.FamilyInetOutput.File.NoStamp.IsNull() &&
				!config.FamilyInetOutput.File.NoStamp.IsUnknown() &&
				!config.FamilyInetOutput.File.Stamp.IsNull() &&
				!config.FamilyInetOutput.File.Stamp.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet_output").AtName("file").AtName("stamp"),
					tfdiag.ConflictConfigErrSummary,
					"no_stamp and stamp can't be true in same time in family_inet_output.file block",
				)
			}
			if !config.FamilyInetOutput.File.NoWorldReadable.IsNull() &&
				!config.FamilyInetOutput.File.NoWorldReadable.IsUnknown() &&
				!config.FamilyInetOutput.File.WorldReadable.IsNull() &&
				!config.FamilyInetOutput.File.WorldReadable.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet_output").AtName("file").AtName("world_readable"),
					tfdiag.ConflictConfigErrSummary,
					"no_world_readable and world_readable can't be true in same time in family_inet_output.file block",
				)
			}
		}
		if config.FamilyInetOutput.InlineJflowSourceAddress.IsNull() {
			if !config.FamilyInetOutput.InlineJflowExportRate.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet_output").AtName("inline_jflow_export_rate"),
					tfdiag.MissingConfigErrSummary,
					"inline_jflow_source_address must be specified with inline_jflow_export_rate in family_inet_output block",
				)
			}
		} else if config.FamilyInetOutput.FlowServer.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet_output").AtName("inline_jflow_source_address"),
				tfdiag.MissingConfigErrSummary,
				"flow_server must be specified with inline_jflow_source_address in family_inet_output block",
			)
		}
		if !config.FamilyInetOutput.FlowServer.IsNull() && !config.FamilyInetOutput.FlowServer.IsUnknown() {
			var configFlowServer []forwardingoptionsSamplingBlockFamilyInetOutputBlockFlowServer
			asDiags := config.FamilyInetOutput.FlowServer.ElementsAs(ctx, &configFlowServer, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}
			flowServerHostname := make(map[string]struct{})
			for _, block := range configFlowServer {
				if block.Hostname.IsUnknown() {
					continue
				}
				hostname := block.Hostname.ValueString()
				if _, ok := flowServerHostname[hostname]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet_output").AtName("flow_server"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple flow_server blocks with the same hostname %q"+
							" in family_inet_output block", hostname),
					)
				}
				flowServerHostname[hostname] = struct{}{}
			}
		}
		if !config.FamilyInetOutput.Interface.IsNull() && !config.FamilyInetOutput.Interface.IsUnknown() {
			var configInterface []forwardingoptionsSamplingBlockOutputBlockInterface
			asDiags := config.FamilyInetOutput.Interface.ElementsAs(ctx, &configInterface, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}
			interfaceName := make(map[string]struct{})
			for i, block := range configInterface {
				if block.Name.IsUnknown() {
					continue
				}
				name := block.Name.ValueString()
				if _, ok := interfaceName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet_output").AtName("interface").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple interface blocks with the same name %q"+
							" in family_inet_output block", name),
					)
				}
				interfaceName[name] = struct{}{}
			}
		}
	}
	if config.FamilyInet6Output != nil {
		if config.Input == nil &&
			config.FamilyInet6Input == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet6_output").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"one of input or family_inet6_input must be specified with family_inet6_output",
			)
		}
		if config.FamilyInet6Output.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet6_output").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"family_inet6_output block is empty",
			)
		}
		if config.FamilyInet6Output.InlineJflowSourceAddress.IsNull() {
			if !config.FamilyInet6Output.InlineJflowExportRate.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("family_inet6_output").AtName("inline_jflow_export_rate"),
					tfdiag.MissingConfigErrSummary,
					"inline_jflow_source_address must be specified with inline_jflow_export_rate in family_inet6_output block",
				)
			}
		} else if config.FamilyInet6Output.FlowServer.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_inet6_output").AtName("inline_jflow_source_address"),
				tfdiag.MissingConfigErrSummary,
				"flow_server must be specified with inline_jflow_source_address in family_inet6_output block",
			)
		}
		if !config.FamilyInet6Output.FlowServer.IsNull() && !config.FamilyInet6Output.FlowServer.IsUnknown() {
			var configFlowServer []forwardingoptionsSamplingBlockOutputBlockFlowServer
			asDiags := config.FamilyInet6Output.FlowServer.ElementsAs(ctx, &configFlowServer, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}
			flowServerHostname := make(map[string]struct{})
			for _, block := range configFlowServer {
				if block.Hostname.IsUnknown() {
					continue
				}
				hostname := block.Hostname.ValueString()
				if _, ok := flowServerHostname[hostname]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet6_output").AtName("flow_server"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple flow_server blocks with the same hostname %q"+
							" in family_inet6_output block", hostname),
					)
				}
				flowServerHostname[hostname] = struct{}{}
			}
		}
		if !config.FamilyInet6Output.Interface.IsNull() && !config.FamilyInet6Output.Interface.IsUnknown() {
			var configInterface []forwardingoptionsSamplingBlockOutputBlockInterface
			asDiags := config.FamilyInet6Output.Interface.ElementsAs(ctx, &configInterface, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}
			interfaceName := make(map[string]struct{})
			for i, block := range configInterface {
				if block.Name.IsUnknown() {
					continue
				}
				name := block.Name.ValueString()
				if _, ok := interfaceName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_inet6_output").AtName("interface").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple interface blocks with the same name %q"+
							" in family_inet6_output block", name),
					)
				}
				interfaceName[name] = struct{}{}
			}
		}
	}
	if config.FamilyMplsOutput != nil {
		if config.Input == nil &&
			config.FamilyMplsInput == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_mpls_output").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"one of input or family_mpls_input must be specified with family_mpls_output",
			)
		}
		if config.FamilyMplsOutput.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("family_mpls_output").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"family_mpls_output block is empty",
			)
		}
		if !config.FamilyMplsOutput.FlowServer.IsNull() && !config.FamilyMplsOutput.FlowServer.IsUnknown() {
			var configFlowServer []forwardingoptionsSamplingBlockOutputBlockFlowServer
			asDiags := config.FamilyMplsOutput.FlowServer.ElementsAs(ctx, &configFlowServer, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}
			flowServerHostname := make(map[string]struct{})
			for _, block := range configFlowServer {
				if block.Hostname.IsUnknown() {
					continue
				}
				hostname := block.Hostname.ValueString()
				if _, ok := flowServerHostname[hostname]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_mpls_output").AtName("flow_server"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple flow_server blocks with the same hostname %q"+
							" in family_mpls_output block", hostname),
					)
				}
				flowServerHostname[hostname] = struct{}{}
			}
		}
		if !config.FamilyMplsOutput.Interface.IsNull() && !config.FamilyMplsOutput.Interface.IsUnknown() {
			var configInterface []forwardingoptionsSamplingBlockOutputBlockInterface
			asDiags := config.FamilyMplsOutput.Interface.ElementsAs(ctx, &configInterface, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}
			interfaceName := make(map[string]struct{})
			for i, block := range configInterface {
				if block.Name.IsUnknown() {
					continue
				}
				name := block.Name.ValueString()
				if _, ok := interfaceName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("family_mpls_output").AtName("interface").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple interface blocks with the same name %q"+
							" in family_mpls_output block", name),
					)
				}
				interfaceName[name] = struct{}{}
			}
		}
	}
}

func (rsc *forwardingoptionsSampling) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan forwardingoptionsSamplingData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
				instanceExists, err := checkRoutingInstanceExists(fnCtx, v, junSess)
				if err != nil {
					resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

					return false
				}
				if !instanceExists {
					resp.Diagnostics.AddAttributeError(
						path.Root("routing_instance"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("routing instance %q doesn't exist", v),
					)

					return false
				}
			}
			var check forwardingoptionsSamplingData
			if err := check.read(fnCtx, plan.RoutingInstance.ValueString(), junSess); err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if check.FamilyInetInput != nil ||
				check.FamilyInetOutput != nil ||
				check.FamilyInet6Input != nil ||
				check.FamilyInet6Output != nil ||
				check.FamilyMplsInput != nil ||
				check.FamilyMplsOutput != nil ||
				check.Input != nil ||
				!check.Disable.IsNull() ||
				!check.PreRewriteTos.IsNull() ||
				!check.SampleOnce.IsNull() {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf(rsc.junosName()+" with routing-instance %q already configured", plan.RoutingInstance.ValueString()),
				)

				return false
			}

			return true
		},
		nil,
		&plan,
		resp,
	)
}

func (rsc *forwardingoptionsSampling) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data forwardingoptionsSamplingData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	defer junos.MutexUnlock()

	if v := state.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, v, junSess)
		if err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !instanceExists {
			resp.State.RemoveResource(ctx)

			return
		}
	}
	if err := data.read(ctx, state.RoutingInstance.ValueString(), junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if data.ID.IsNull() {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *forwardingoptionsSampling) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state forwardingoptionsSamplingData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *forwardingoptionsSampling) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state forwardingoptionsSamplingData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceDelete(
		ctx,
		rsc,
		&state,
		resp,
	)
}

func (rsc *forwardingoptionsSampling) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	var data forwardingoptionsSamplingData
	if req.ID != junos.DefaultW {
		instanceExists, err := checkRoutingInstanceExists(ctx, req.ID, junSess)
		if err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

			return
		}
		if !instanceExists {
			resp.Diagnostics.AddError(
				tfdiag.NotFoundErrSummary,
				fmt.Sprintf("routing instance %q doesn't exist", req.ID),
			)

			return
		}
	}
	if err := data.read(ctx, req.ID, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rscData *forwardingoptionsSamplingData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(v)
	} else {
		rscData.ID = types.StringValue(junos.DefaultW)
	}
}

func (rscData *forwardingoptionsSamplingData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "forwarding-options sampling "

	if rscData.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if rscData.PreRewriteTos.ValueBool() {
		configSet = append(configSet, setPrefix+"pre-rewrite-tos")
	}
	if rscData.SampleOnce.ValueBool() {
		configSet = append(configSet, setPrefix+"sample-once")
	}
	if rscData.FamilyInetInput != nil {
		blockSet := rscData.FamilyInetInput.configSet(setPrefix + "family inet input ")
		if len(blockSet) == 0 {
			return path.Root("family_inet_input").AtName("*"), errors.New("family_inet_input block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.FamilyInetOutput != nil {
		blockSet, pathErr, err := rscData.FamilyInetOutput.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		if len(blockSet) == 0 {
			return path.Root("family_inet_output").AtName("*"), errors.New("family_inet_output block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.FamilyInet6Input != nil {
		blockSet := rscData.FamilyInet6Input.configSet(setPrefix + "family inet6 input ")
		if len(blockSet) == 0 {
			return path.Root("family_inet6_input").AtName("*"), errors.New("family_inet6_input block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.FamilyInet6Output != nil {
		blockSet, pathErr, err := rscData.FamilyInet6Output.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		if len(blockSet) == 0 {
			return path.Root("family_inet6_output").AtName("*"), errors.New("family_inet6_output block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.FamilyMplsInput != nil {
		blockSet := rscData.FamilyMplsInput.configSet(setPrefix + "family mpls input ")
		if len(blockSet) == 0 {
			return path.Root("family_mpls_input").AtName("*"), errors.New("family_mpls_input block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.FamilyMplsOutput != nil {
		blockSet, pathErr, err := rscData.FamilyMplsOutput.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		if len(blockSet) == 0 {
			return path.Root("family_mpls_output").AtName("*"), errors.New("family_mpls_output block is empty")
		}
		configSet = append(configSet, blockSet...)
	}
	if rscData.Input != nil {
		blockSet := rscData.Input.configSet(setPrefix + "input ")
		if len(blockSet) == 0 {
			return path.Root("input").AtName("*"), errors.New("input block is empty")
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *forwardingoptionsSamplingBlockInput) configSet(
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

func (block *forwardingoptionsSamplingBlockFamilyInetOutput) configSet(
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
	if block.File != nil {
		configSet = append(configSet, setPrefix+"file filename \""+block.File.Filename.ValueString()+"\"")
		if block.File.Disable.ValueBool() {
			configSet = append(configSet, setPrefix+"file disable")
		}
		if !block.File.Files.IsNull() {
			configSet = append(configSet, setPrefix+"file files "+utils.ConvI64toa(block.File.Files.ValueInt64()))
		}
		if block.File.NoStamp.ValueBool() {
			configSet = append(configSet, setPrefix+"file no-stamp")
		}
		if block.File.NoWorldReadable.ValueBool() {
			configSet = append(configSet, setPrefix+"file no-world-readable")
		}
		if !block.File.Size.IsNull() {
			configSet = append(configSet, setPrefix+"file size "+utils.ConvI64toa(block.File.Size.ValueInt64()))
		}
		if block.File.Stamp.ValueBool() {
			configSet = append(configSet, setPrefix+"file stamp")
		}
		if block.File.WorldReadable.ValueBool() {
			configSet = append(configSet, setPrefix+"file world-readable")
		}
	}
	if !block.FlowActiveTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"flow-active-timeout "+
			utils.ConvI64toa(block.FlowActiveTimeout.ValueInt64()))
	}
	if !block.FlowInactiveTimeout.IsNull() {
		configSet = append(configSet, setPrefix+"flow-inactive-timeout "+
			utils.ConvI64toa(block.FlowInactiveTimeout.ValueInt64()))
	}
	if !block.InlineJflowExportRate.IsNull() {
		configSet = append(configSet, setPrefix+"inline-jflow flow-export-rate "+
			utils.ConvI64toa(block.InlineJflowExportRate.ValueInt64()))
	}
	if v := block.InlineJflowSourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"inline-jflow source-address "+v)
	}
	flowServerHostname := make(map[string]struct{})
	for _, blockFlowServer := range block.FlowServer {
		hostname := blockFlowServer.Hostname.ValueString()
		if _, ok := flowServerHostname[hostname]; ok {
			return configSet,
				path.Root("family_inet_output").AtName("flow_server"),
				fmt.Errorf("multiple flow_server blocks with the same hostname %q in family_inet_output block",
					hostname)
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
				path.Root("family_inet_output").AtName("flow_server"),
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
	}
	interfaceName := make(map[string]struct{})
	for i, blockInterface := range block.Interface {
		name := blockInterface.Name.ValueString()
		if _, ok := interfaceName[name]; ok {
			return configSet,
				path.Root("family_inet_output").AtName("interface").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple interface blocks with the same name %q"+
					" in family_inet_output block", name)
		}
		interfaceName[name] = struct{}{}

		configSet = append(configSet, blockInterface.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsSamplingBlockFamilyInet6Output) configSet(
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
	for _, blockFlowServer := range block.FlowServer {
		hostname := blockFlowServer.Hostname.ValueString()
		if _, ok := flowServerHostname[hostname]; ok {
			return configSet,
				path.Root("family_inet6_output").AtName("flow_server"),
				fmt.Errorf("multiple flow_server blocks with the same hostname %q"+
					" in family_inet6_output block", hostname)
		}
		flowServerHostname[hostname] = struct{}{}

		blockSet, err := blockFlowServer.configSet(setPrefix)
		if err != nil {
			return configSet, path.Root("family_inet6_output").AtName("flow_server"), err
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
		name := blockInterface.Name.ValueString()
		if _, ok := interfaceName[name]; ok {
			return configSet,
				path.Root("family_inet6_output").AtName("interface").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple interface blocks with the same name %q"+
					" in family_inet6_output block", name)
		}
		interfaceName[name] = struct{}{}

		configSet = append(configSet, blockInterface.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsSamplingBlockFamilyMplsOutput) configSet(
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
	for _, blockFlowServer := range block.FlowServer {
		hostname := blockFlowServer.Hostname.ValueString()
		if _, ok := flowServerHostname[hostname]; ok {
			return configSet,
				path.Root("family_mpls_output").AtName("flow_server"),
				fmt.Errorf("multiple flow_server blocks with the same hostname %q in family_mpls_output block",
					hostname)
		}
		flowServerHostname[hostname] = struct{}{}

		blockSet, err := blockFlowServer.configSet(setPrefix)
		if err != nil {
			return configSet, path.Root("family_mpls_output").AtName("flow_server"), err
		}
		configSet = append(configSet, blockSet...)
	}
	interfaceName := make(map[string]struct{})
	for i, blockInterface := range block.Interface {
		name := blockInterface.Name.ValueString()
		if _, ok := interfaceName[name]; ok {
			return configSet,
				path.Root("family_mpls_output").AtName("interface").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple interface blocks with the same name %q"+
					" in family_mpls_output block", name)
		}
		interfaceName[name] = struct{}{}

		configSet = append(configSet, blockInterface.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (block *forwardingoptionsSamplingBlockOutputBlockFlowServer) configSet(
	setPrefix string,
) (
	[]string, error,
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

	return configSet, nil
}

func (block *forwardingoptionsSamplingBlockOutputBlockInterface) configSet(setPrefix string) []string {
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

func (rscData *forwardingoptionsSamplingData) read(
	_ context.Context, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"forwarding-options sampling" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if routingInstance == "" {
		rscData.RoutingInstance = types.StringValue(junos.DefaultW)
	} else {
		rscData.RoutingInstance = types.StringValue(routingInstance)
	}
	rscData.fillID()
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == junos.DisableW:
				rscData.Disable = types.BoolValue(true)
			case itemTrim == "pre-rewrite-tos":
				rscData.PreRewriteTos = types.BoolValue(true)
			case itemTrim == "sample-once":
				rscData.SampleOnce = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "family inet input "):
				if rscData.FamilyInetInput == nil {
					rscData.FamilyInetInput = &forwardingoptionsSamplingBlockInput{}
				}
				if err := rscData.FamilyInetInput.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet6 input "):
				if rscData.FamilyInet6Input == nil {
					rscData.FamilyInet6Input = &forwardingoptionsSamplingBlockInput{}
				}
				if err := rscData.FamilyInet6Input.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "family mpls input "):
				if rscData.FamilyMplsInput == nil {
					rscData.FamilyMplsInput = &forwardingoptionsSamplingBlockInput{}
				}
				if err := rscData.FamilyMplsInput.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "input "):
				if rscData.Input == nil {
					rscData.Input = &forwardingoptionsSamplingBlockInput{}
				}
				if err := rscData.Input.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet output "):
				if rscData.FamilyInetOutput == nil {
					rscData.FamilyInetOutput = &forwardingoptionsSamplingBlockFamilyInetOutput{}
				}
				if err := rscData.FamilyInetOutput.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "family inet6 output "):
				if rscData.FamilyInet6Output == nil {
					rscData.FamilyInet6Output = &forwardingoptionsSamplingBlockFamilyInet6Output{}
				}
				if err := rscData.FamilyInet6Output.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "family mpls output "):
				if rscData.FamilyMplsOutput == nil {
					rscData.FamilyMplsOutput = &forwardingoptionsSamplingBlockFamilyMplsOutput{}
				}
				if err := rscData.FamilyMplsOutput.read(itemTrim); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *forwardingoptionsSamplingBlockInput) read(itemTrim string) (err error) {
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

func (block *forwardingoptionsSamplingBlockFamilyInetOutput) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "aggregate-export-interval "):
		block.AggregateExportInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "extension-service "):
		block.ExtensionService = append(block.ExtensionService, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "file "):
		if block.File == nil {
			block.File = &forwardingoptionsSamplingBlockFamilyInetOutputBlockFile{}
		}
		switch {
		case itemTrim == "disable":
			block.File.Disable = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "filename "):
			block.File.Filename = types.StringValue(strings.Trim(itemTrim, "\""))
		case balt.CutPrefixInString(&itemTrim, "files "):
			block.File.Files, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		case itemTrim == "no-stamp":
			block.File.NoStamp = types.BoolValue(true)
		case itemTrim == "no-world-readable":
			block.File.NoWorldReadable = types.BoolValue(true)
		case balt.CutPrefixInString(&itemTrim, "size "):
			switch {
			case balt.CutSuffixInString(&itemTrim, "k"):
				block.File.Size, err = tfdata.ConvAtoi64Value(itemTrim)
				block.File.Size = types.Int64Value(block.File.Size.ValueInt64() * 1024)
			case balt.CutSuffixInString(&itemTrim, "m"):
				block.File.Size, err = tfdata.ConvAtoi64Value(itemTrim)
				block.File.Size = types.Int64Value(block.File.Size.ValueInt64() * 1024 * 1024)
			case balt.CutSuffixInString(&itemTrim, "g"):
				block.File.Size = types.Int64Value(block.File.Size.ValueInt64() * 1024 * 1024 * 1024)
			default:
				block.File.Size, err = tfdata.ConvAtoi64Value(itemTrim)
			}
			if err != nil {
				return err
			}
		case itemTrim == "stamp":
			block.File.Stamp = types.BoolValue(true)
		case itemTrim == "world-readable":
			block.File.WorldReadable = types.BoolValue(true)
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
	case balt.CutPrefixInString(&itemTrim, "inline-jflow flow-export-rate "):
		block.InlineJflowExportRate, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "inline-jflow source-address "):
		block.InlineJflowSourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "flow-server "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var flowServer forwardingoptionsSamplingBlockFamilyInetOutputBlockFlowServer
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
		}
		block.FlowServer = append(block.FlowServer, flowServer)
	case balt.CutPrefixInString(&itemTrim, "interface "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var iFace forwardingoptionsSamplingBlockOutputBlockInterface
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

func (block *forwardingoptionsSamplingBlockFamilyInet6Output) read(itemTrim string) (err error) {
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
		var flowServer forwardingoptionsSamplingBlockOutputBlockFlowServer
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
		var iFace forwardingoptionsSamplingBlockOutputBlockInterface
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

func (block *forwardingoptionsSamplingBlockFamilyMplsOutput) read(itemTrim string) (err error) {
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
		var flowServer forwardingoptionsSamplingBlockOutputBlockFlowServer
		block.FlowServer, flowServer = tfdata.ExtractBlockWithTFTypesString(
			block.FlowServer, "Hostname", itemTrimFields[0],
		)
		flowServer.Hostname = types.StringValue(itemTrimFields[0])
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		if err := flowServer.read(itemTrim); err != nil {
			return err
		}
		block.FlowServer = append(block.FlowServer, flowServer)
	case balt.CutPrefixInString(&itemTrim, "interface "):
		itemTrimFields := strings.Split(itemTrim, " ")
		var iFace forwardingoptionsSamplingBlockOutputBlockInterface
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

func (block *forwardingoptionsSamplingBlockOutputBlockFlowServer) read(itemTrim string) (err error) {
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
	}

	return nil
}

func (block *forwardingoptionsSamplingBlockOutputBlockInterface) read(itemTrim string) (err error) {
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

func (rscData *forwardingoptionsSamplingData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}
	delPrefix += "forwarding-options sampling "

	configSet := []string{
		delPrefix + junos.DisableW,
		delPrefix + "pre-rewrite-tos",
		delPrefix + "sample-once",
		delPrefix + "family inet input",
		delPrefix + "family inet output",
		delPrefix + "family inet6 input",
		delPrefix + "family inet6 output",
		delPrefix + "family mpls input",
		delPrefix + "family mpls output",
		delPrefix + "input",
	}

	return junSess.ConfigSet(configSet)
}
