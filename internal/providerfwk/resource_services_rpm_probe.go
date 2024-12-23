package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &servicesRpmProbe{}
	_ resource.ResourceWithConfigure      = &servicesRpmProbe{}
	_ resource.ResourceWithValidateConfig = &servicesRpmProbe{}
	_ resource.ResourceWithImportState    = &servicesRpmProbe{}
	_ resource.ResourceWithUpgradeState   = &servicesRpmProbe{}
)

type servicesRpmProbe struct {
	client *junos.Client
}

func newServicesRpmProbeResource() resource.Resource {
	return &servicesRpmProbe{}
}

func (rsc *servicesRpmProbe) typeName() string {
	return providerName + "_services_rpm_probe"
}

func (rsc *servicesRpmProbe) junosName() string {
	return "services rpm probe"
}

func (rsc *servicesRpmProbe) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *servicesRpmProbe) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *servicesRpmProbe) Configure(
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

func (rsc *servicesRpmProbe) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of owner.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"delegate_probes": schema.BoolAttribute{
				Optional:    true,
				Description: "Offload real-time performance monitoring probes to MS-MIC/MS-MPC card.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"test": schema.ListNestedBlock{
				Description: "For each name of test, configure a test.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of test.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"data_fill": schema.StringAttribute{
							Optional:    true,
							Description: "Define contents of the data portion of the probes.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[0-9a-fA-F]+$`),
									"must be hexadecimal digits (0-9, a-f, A-F)",
								),
							},
						},
						"data_size": schema.Int64Attribute{
							Optional:    true,
							Description: "Size of the data portion of the probes.",
							Validators: []validator.Int64{
								int64validator.Between(0, 65400),
							},
						},
						"destination_interface": schema.StringAttribute{
							Optional:    true,
							Description: "Name of output interface for probes.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
								tfvalidator.String1DotCount(),
							},
						},
						"destination_port": schema.Int64Attribute{
							Optional:    true,
							Description: "TCP/UDP port number.",
							Validators: []validator.Int64{
								int64validator.Between(7, 65535),
							},
						},
						"dscp_code_points": schema.StringAttribute{
							Optional:    true,
							Description: "Differentiated Services code point bits or alias.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 64),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"hardware_timestamp": schema.BoolAttribute{
							Optional:    true,
							Description: "Packet Forwarding Engine updates timestamps.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"history_size": schema.Int64Attribute{
							Optional:    true,
							Description: "Number of stored history entries.",
							Validators: []validator.Int64{
								int64validator.Between(0, 512),
							},
						},
						"inet6_source_address": schema.StringAttribute{
							Optional:    true,
							Description: "Inet6 Source Address of the probe.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress().IPv6Only(),
							},
						},
						"moving_average_size": schema.Int64Attribute{
							Optional:    true,
							Description: "Number of samples used for moving average.",
							Validators: []validator.Int64{
								int64validator.Between(0, 1024),
							},
						},
						"one_way_hardware_timestamp": schema.BoolAttribute{
							Optional:    true,
							Description: "Enable hardware timestamps for one-way measurements.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"probe_count": schema.Int64Attribute{
							Optional:    true,
							Description: "Total number of probes per test.",
							Validators: []validator.Int64{
								int64validator.Between(1, 15),
							},
						},
						"probe_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Delay between probes (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 255),
							},
						},
						"probe_type": schema.StringAttribute{
							Optional:    true,
							Description: "Probe request type.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"http-get", "http-metadata-get",
									"icmp-ping", "icmp-ping-timestamp", "icmp6-ping",
									"tcp-ping",
									"udp-ping", "udp-ping-timestamp",
								),
							},
						},
						"routing_instance": schema.StringAttribute{
							Optional:    true,
							Description: "Routing instance used by probes.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 64),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"source_address": schema.StringAttribute{
							Optional:    true,
							Description: "Source address for probe.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress().IPv4Only(),
							},
						},
						"target_type": schema.StringAttribute{
							Optional:    true,
							Description: "Type of target destination for probe.",
							Validators: []validator.String{
								stringvalidator.OneOf("address", "inet6-address", "inet6-url", "url"),
							},
						},
						"target_value": schema.StringAttribute{
							Optional:    true,
							Description: "Target destination for probe.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"test_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Delay between tests (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(0, 86400),
							},
						},
						"traps": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "Trap to send if threshold is met or exceeded.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringFormat(tfvalidator.DefaultFormat),
								),
							},
						},
						"ttl": schema.Int64Attribute{
							Optional:    true,
							Description: "Time to Live (hop-limit) value for an RPM IPv4(IPv6) packet.",
							Validators: []validator.Int64{
								int64validator.Between(1, 254),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"rpm_scale": schema.SingleNestedBlock{
							Description: "Configure real-time performance monitoring scale tests.",
							Attributes: map[string]schema.Attribute{
								"tests_count": schema.Int64Attribute{
									Required:    false, // true when SingleNestedBlock is specified
									Optional:    true,
									Description: "Number of probe-tests generated using scale config.",
									Validators: []validator.Int64{
										int64validator.Between(1, 500000),
									},
								},
								"destination_interface": schema.StringAttribute{
									Optional:    true,
									Description: "Base destination interface for scale test.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
										tfvalidator.String1DotCount(),
									},
								},
								"destination_subunit_cnt": schema.Int64Attribute{
									Optional:    true,
									Description: "Subunit count for destination interface for scale test.",
									Validators: []validator.Int64{
										int64validator.Between(1, 500000),
									},
								},
								"source_address_base": schema.StringAttribute{
									Optional:    true,
									Description: "Source base address of host.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv4Only(),
									},
								},
								"source_count": schema.Int64Attribute{
									Optional:    true,
									Description: "Source-address count.",
									Validators: []validator.Int64{
										int64validator.Between(1, 500000),
									},
								},
								"source_step": schema.StringAttribute{
									Optional:    true,
									Description: "Steps to increment src address.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv4Only(),
									},
								},
								"source_inet6_address_base": schema.StringAttribute{
									Optional:    true,
									Description: "Source base inet6 address of host.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv6Only(),
									},
								},
								"source_inet6_count": schema.Int64Attribute{
									Optional:    true,
									Description: "Source-inet6-address count.",
									Validators: []validator.Int64{
										int64validator.Between(1, 500000),
									},
								},
								"source_inet6_step": schema.StringAttribute{
									Optional:    true,
									Description: "Steps to increment src inet6 address.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv6Only(),
									},
								},
								"target_address_base": schema.StringAttribute{
									Optional:    true,
									Description: "Base address of target host.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv4Only(),
									},
								},
								"target_count": schema.Int64Attribute{
									Optional:    true,
									Description: "Target address count.",
									Validators: []validator.Int64{
										int64validator.Between(1, 500000),
									},
								},
								"target_step": schema.StringAttribute{
									Optional:    true,
									Description: "Steps to increment target address.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv4Only(),
									},
								},
								"target_inet6_address_base": schema.StringAttribute{
									Optional:    true,
									Description: "Base inet6 address of target host.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv6Only(),
									},
								},
								"target_inet6_count": schema.Int64Attribute{
									Optional:    true,
									Description: "Target inet6 address count.",
									Validators: []validator.Int64{
										int64validator.Between(1, 500000),
									},
								},
								"target_inet6_step": schema.StringAttribute{
									Optional:    true,
									Description: "Steps to increment target inet6 address.",
									Validators: []validator.String{
										tfvalidator.StringIPAddress().IPv6Only(),
									},
								},
							},
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
						"thresholds": schema.SingleNestedBlock{
							Description: "Declare `thresholds` configuration.",
							Attributes: map[string]schema.Attribute{
								"egress_time": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum source to destination time per probe (microseconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 60000000),
									},
								},
								"ingress_time": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum destination to source time per probe (microseconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 60000000),
									},
								},
								"jitter_egress": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum source to destination jitter per test (microseconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 60000000),
									},
								},
								"jitter_ingress": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum destination to source jitter per test (microseconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 60000000),
									},
								},
								"jitter_rtt": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum jitter per test (microseconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 60000000),
									},
								},
								"rtt": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum round trip time per probe (microseconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 60000000),
									},
								},
								"std_dev_egress": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum source to destination standard deviation per test (microseconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 60000000),
									},
								},
								"std_dev_ingress": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum destination to source standard deviation per test (microseconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 60000000),
									},
								},
								"std_dev_rtt": schema.Int64Attribute{
									Optional:    true,
									Description: "Maximum standard deviation per test (microseconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 60000000),
									},
								},
								"successive_loss": schema.Int64Attribute{
									Optional:    true,
									Description: "Successive probe loss count indicating probe failure.",
									Validators: []validator.Int64{
										int64validator.Between(0, 15),
									},
								},
								"total_loss": schema.Int64Attribute{
									Optional:    true,
									Description: "Total probe loss count indicating test failure.",
									Validators: []validator.Int64{
										int64validator.Between(0, 15),
									},
								},
							},
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
					},
				},
			},
		},
	}
}

type servicesRpmProbeData struct {
	ID             types.String                `tfsdk:"id"`
	Name           types.String                `tfsdk:"name"`
	DelegateProbes types.Bool                  `tfsdk:"delegate_probes"`
	Test           []servicesRpmProbeBlockTest `tfsdk:"test"`
}

type servicesRpmProbeConfig struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	DelegateProbes types.Bool   `tfsdk:"delegate_probes"`
	Test           types.List   `tfsdk:"test"`
}

//nolint:lll
type servicesRpmProbeBlockTest struct {
	Name                    types.String                              `tfsdk:"name"                       tfdata:"identifier"`
	DataFill                types.String                              `tfsdk:"data_fill"`
	DataSize                types.Int64                               `tfsdk:"data_size"`
	DestinationInterface    types.String                              `tfsdk:"destination_interface"`
	DestinationPort         types.Int64                               `tfsdk:"destination_port"`
	DscpCodePoints          types.String                              `tfsdk:"dscp_code_points"`
	HardwareTimestamp       types.Bool                                `tfsdk:"hardware_timestamp"`
	HistorySize             types.Int64                               `tfsdk:"history_size"`
	Inet6SourceAddress      types.String                              `tfsdk:"inet6_source_address"`
	MovingAverageSize       types.Int64                               `tfsdk:"moving_average_size"`
	OneWayHardwareTimestamp types.Bool                                `tfsdk:"one_way_hardware_timestamp"`
	ProbeCount              types.Int64                               `tfsdk:"probe_count"`
	ProbeInterval           types.Int64                               `tfsdk:"probe_interval"`
	ProbeType               types.String                              `tfsdk:"probe_type"`
	RoutingInstance         types.String                              `tfsdk:"routing_instance"`
	SourceAddress           types.String                              `tfsdk:"source_address"`
	TargetType              types.String                              `tfsdk:"target_type"`
	TargetValue             types.String                              `tfsdk:"target_value"`
	TestInterval            types.Int64                               `tfsdk:"test_interval"`
	Traps                   []types.String                            `tfsdk:"traps"`
	TTL                     types.Int64                               `tfsdk:"ttl"`
	RpmScale                *servicesRpmProbeBlockTestBlockRpmScale   `tfsdk:"rpm_scale"`
	Thresholds              *servicesRpmProbeBlockTestBlockThresholds `tfsdk:"thresholds"`
}

type servicesRpmProbeBlockTestConfig struct {
	Name                    types.String                              `tfsdk:"name"`
	DataFill                types.String                              `tfsdk:"data_fill"`
	DataSize                types.Int64                               `tfsdk:"data_size"`
	DestinationInterface    types.String                              `tfsdk:"destination_interface"`
	DestinationPort         types.Int64                               `tfsdk:"destination_port"`
	DscpCodePoints          types.String                              `tfsdk:"dscp_code_points"`
	HardwareTimestamp       types.Bool                                `tfsdk:"hardware_timestamp"`
	HistorySize             types.Int64                               `tfsdk:"history_size"`
	Inet6SourceAddress      types.String                              `tfsdk:"inet6_source_address"`
	MovingAverageSize       types.Int64                               `tfsdk:"moving_average_size"`
	OneWayHardwareTimestamp types.Bool                                `tfsdk:"one_way_hardware_timestamp"`
	ProbeCount              types.Int64                               `tfsdk:"probe_count"`
	ProbeInterval           types.Int64                               `tfsdk:"probe_interval"`
	ProbeType               types.String                              `tfsdk:"probe_type"`
	RoutingInstance         types.String                              `tfsdk:"routing_instance"`
	SourceAddress           types.String                              `tfsdk:"source_address"`
	TargetType              types.String                              `tfsdk:"target_type"`
	TargetValue             types.String                              `tfsdk:"target_value"`
	TestInterval            types.Int64                               `tfsdk:"test_interval"`
	Traps                   types.Set                                 `tfsdk:"traps"`
	TTL                     types.Int64                               `tfsdk:"ttl"`
	RpmScale                *servicesRpmProbeBlockTestBlockRpmScale   `tfsdk:"rpm_scale"`
	Thresholds              *servicesRpmProbeBlockTestBlockThresholds `tfsdk:"thresholds"`
}

type servicesRpmProbeBlockTestBlockRpmScale struct {
	TestsCount             types.Int64  `tfsdk:"tests_count"`
	DestinationInterface   types.String `tfsdk:"destination_interface"`
	DestinationSubunitCnt  types.Int64  `tfsdk:"destination_subunit_cnt"`
	SourceAddressBase      types.String `tfsdk:"source_address_base"`
	SourceCount            types.Int64  `tfsdk:"source_count"`
	SourceStep             types.String `tfsdk:"source_step"`
	SourceInet6AddressBase types.String `tfsdk:"source_inet6_address_base"`
	SourceInet6Count       types.Int64  `tfsdk:"source_inet6_count"`
	SourceInet6Step        types.String `tfsdk:"source_inet6_step"`
	TargetAddressBase      types.String `tfsdk:"target_address_base"`
	TargetCount            types.Int64  `tfsdk:"target_count"`
	TargetStep             types.String `tfsdk:"target_step"`
	TargetInet6AddressBase types.String `tfsdk:"target_inet6_address_base"`
	TargetInet6Count       types.Int64  `tfsdk:"target_inet6_count"`
	TargetInet6Step        types.String `tfsdk:"target_inet6_step"`
}

type servicesRpmProbeBlockTestBlockThresholds struct {
	EgressTime     types.Int64 `tfsdk:"egress_time"`
	IngressTime    types.Int64 `tfsdk:"ingress_time"`
	JitterEgress   types.Int64 `tfsdk:"jitter_egress"`
	JitterIngress  types.Int64 `tfsdk:"jitter_ingress"`
	JitterRtt      types.Int64 `tfsdk:"jitter_rtt"`
	Rtt            types.Int64 `tfsdk:"rtt"`
	StdDevEgress   types.Int64 `tfsdk:"std_dev_egress"`
	StdDevIngress  types.Int64 `tfsdk:"std_dev_ingress"`
	StdDevRtt      types.Int64 `tfsdk:"std_dev_rtt"`
	SuccessiveLoss types.Int64 `tfsdk:"successive_loss"`
	TotalLoss      types.Int64 `tfsdk:"total_loss"`
}

func (rsc *servicesRpmProbe) ValidateConfig( //nolint:gocognit
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config servicesRpmProbeConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Test.IsNull() && !config.Test.IsUnknown() {
		var configTest []servicesRpmProbeBlockTestConfig
		asDiags := config.Test.ElementsAs(ctx, &configTest, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		testName := make(map[string]struct{})
		for i, block := range configTest {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := testName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("test").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple test blocks with the same name %q", name),
					)
				}
				testName[name] = struct{}{}
			}

			if !block.TargetType.IsNull() && !block.TargetType.IsUnknown() &&
				block.TargetValue.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("test").AtListIndex(i).AtName("target_type"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("target_value must be specified with target_type"+
						" in test block %q", block.Name.ValueString()),
				)
			}
			if !block.TargetValue.IsNull() && !block.TargetValue.IsUnknown() &&
				block.TargetType.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("test").AtListIndex(i).AtName("target_value"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("target_type must be specified with target_value"+
						" in test block %q", block.Name.ValueString()),
				)
			}
			if block.RpmScale != nil {
				if block.RpmScale.TestsCount.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("tests_count"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("tests_count must be specified"+
							" in rpm_scale block in test block %q", block.Name.ValueString()),
					)
				}
				if !block.RpmScale.DestinationInterface.IsNull() &&
					!block.RpmScale.DestinationInterface.IsUnknown() &&
					block.RpmScale.DestinationSubunitCnt.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("destination_interface"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("destination_subunit_cnt must be specified with destination_interface"+
							" in rpm_scale block in test block %q", block.Name.ValueString()),
					)
				}
				if !block.RpmScale.DestinationSubunitCnt.IsNull() &&
					!block.RpmScale.DestinationSubunitCnt.IsUnknown() &&
					block.RpmScale.DestinationInterface.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("destination_subunit_cnt"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("destination_interface must be specified with destination_subunit_cnt"+
							" in rpm_scale block in test block %q", block.Name.ValueString()),
					)
				}
				if !block.RpmScale.SourceAddressBase.IsNull() &&
					!block.RpmScale.SourceAddressBase.IsUnknown() {
					if block.RpmScale.SourceCount.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_address_base"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_count must be specified with source_address_base"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.SourceStep.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_address_base"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_step must be specified with source_address_base"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.SourceStep.IsNull() &&
					!block.RpmScale.SourceStep.IsUnknown() {
					if block.RpmScale.SourceAddressBase.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_step"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_address_base must be specified with source_step"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.SourceCount.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_step"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_count must be specified with source_step"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.SourceCount.IsNull() &&
					!block.RpmScale.SourceCount.IsUnknown() {
					if block.RpmScale.SourceAddressBase.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_count"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_address_base must be specified with source_count"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.SourceStep.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_count"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_step must be specified with source_count"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.SourceInet6AddressBase.IsNull() &&
					!block.RpmScale.SourceInet6AddressBase.IsUnknown() {
					if block.RpmScale.SourceInet6Count.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_inet6_address_base"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_inet6_count must be specified with source_inet6_address_base"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.SourceInet6Step.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_inet6_address_base"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_inet6_step must be specified with source_inet6_address_base"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.SourceInet6Step.IsNull() &&
					!block.RpmScale.SourceInet6Step.IsUnknown() {
					if block.RpmScale.SourceInet6AddressBase.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_inet6_step"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_inet6_address_base must be specified with source_inet6_step"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.SourceInet6Count.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_inet6_step"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_inet6_count must be specified with source_inet6_step"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.SourceInet6Count.IsNull() &&
					!block.RpmScale.SourceInet6Count.IsUnknown() {
					if block.RpmScale.SourceInet6AddressBase.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_inet6_count"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_inet6_address_base must be specified with source_inet6_count"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.SourceInet6Step.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("source_inet6_count"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("source_inet6_step must be specified with source_inet6_count"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.TargetAddressBase.IsNull() &&
					!block.RpmScale.TargetAddressBase.IsUnknown() {
					if block.RpmScale.TargetCount.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_address_base"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_count must be specified with target_address_base"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.TargetStep.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_address_base"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_step must be specified with target_address_base"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.TargetStep.IsNull() &&
					!block.RpmScale.TargetStep.IsUnknown() {
					if block.RpmScale.TargetAddressBase.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_step"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_address_base must be specified with target_step"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.TargetCount.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_step"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_count must be specified with target_step"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.TargetCount.IsNull() &&
					!block.RpmScale.TargetCount.IsUnknown() {
					if block.RpmScale.TargetAddressBase.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_count"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_address_base must be specified with target_count"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.TargetStep.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_count"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_step must be specified with target_count"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.TargetInet6AddressBase.IsNull() &&
					!block.RpmScale.TargetInet6AddressBase.IsUnknown() {
					if block.RpmScale.TargetInet6Count.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_inet6_address_base"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_inet6_count must be specified with target_inet6_address_base"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.TargetInet6Step.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_inet6_address_base"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_inet6_step must be specified with target_inet6_address_base"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.TargetInet6Step.IsNull() &&
					!block.RpmScale.TargetInet6Step.IsUnknown() {
					if block.RpmScale.TargetInet6AddressBase.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_inet6_step"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_inet6_address_base must be specified with target_inet6_step"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.TargetInet6Count.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_inet6_step"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_inet6_count must be specified with target_inet6_step"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
				if !block.RpmScale.TargetInet6Count.IsNull() &&
					!block.RpmScale.TargetInet6Count.IsUnknown() {
					if block.RpmScale.TargetInet6AddressBase.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_inet6_count"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_inet6_address_base must be specified with target_inet6_count"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
					if block.RpmScale.TargetInet6Step.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("test").AtListIndex(i).AtName("rpm_scale").AtName("target_inet6_count"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf("target_inet6_step must be specified with target_inet6_count"+
								" in rpm_scale block in test block %q", block.Name.ValueString()),
						)
					}
				}
			}
		}
	}
}

func (rsc *servicesRpmProbe) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan servicesRpmProbeData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Empty Name",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "name"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			probeExists, err := checkServicesRpmProbeExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if probeExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			probeExists, err := checkServicesRpmProbeExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !probeExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *servicesRpmProbe) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data servicesRpmProbeData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *servicesRpmProbe) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state servicesRpmProbeData
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

func (rsc *servicesRpmProbe) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state servicesRpmProbeData
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

func (rsc *servicesRpmProbe) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data servicesRpmProbeData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "name"),
	)
}

func checkServicesRpmProbeExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services rpm probe \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *servicesRpmProbeData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *servicesRpmProbeData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *servicesRpmProbeData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set services rpm probe \"" + rscData.Name.ValueString() + "\" "

	configSet := []string{
		setPrefix,
	}

	if rscData.DelegateProbes.ValueBool() {
		configSet = append(configSet, setPrefix+"delegate-probes")
	}

	testName := make(map[string]struct{})
	for i, block := range rscData.Test {
		name := block.Name.ValueString()
		if _, ok := testName[name]; ok {
			return path.Root("test").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple test blocks with the same name %q", name)
		}
		testName[name] = struct{}{}

		blockSet, pathErr, err := block.configSet(setPrefix, path.Root("test").AtListIndex(i))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *servicesRpmProbeBlockTest) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "test \"" + block.Name.ValueString() + "\" "

	configSet := []string{
		setPrefix,
	}

	if v := block.DataFill.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"data-fill "+v)
	}
	if !block.DataSize.IsNull() {
		configSet = append(configSet, setPrefix+"data-size "+
			utils.ConvI64toa(block.DataSize.ValueInt64()))
	}
	if v := block.DestinationInterface.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"destination-interface "+v)
	}
	if !block.DestinationPort.IsNull() {
		configSet = append(configSet, setPrefix+"destination-port "+
			utils.ConvI64toa(block.DestinationPort.ValueInt64()))
	}
	if v := block.DscpCodePoints.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"dscp-code-points \""+v+"\"")
	}
	if block.HardwareTimestamp.ValueBool() {
		configSet = append(configSet, setPrefix+"hardware-timestamp")
	}
	if !block.HistorySize.IsNull() {
		configSet = append(configSet, setPrefix+"history-size "+
			utils.ConvI64toa(block.HistorySize.ValueInt64()))
	}
	if v := block.Inet6SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"inet6-options source-address "+v)
	}
	if !block.MovingAverageSize.IsNull() {
		configSet = append(configSet, setPrefix+"moving-average-size "+
			utils.ConvI64toa(block.MovingAverageSize.ValueInt64()))
	}
	if block.OneWayHardwareTimestamp.ValueBool() {
		configSet = append(configSet, setPrefix+"one-way-hardware-timestamp")
	}
	if !block.ProbeCount.IsNull() {
		configSet = append(configSet, setPrefix+"probe-count "+
			utils.ConvI64toa(block.ProbeCount.ValueInt64()))
	}
	if !block.ProbeInterval.IsNull() {
		configSet = append(configSet, setPrefix+"probe-interval "+
			utils.ConvI64toa(block.ProbeInterval.ValueInt64()))
	}
	if v := block.ProbeType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"probe-type "+v)
	}
	if v := block.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-instance \""+v+"\"")
	}
	if v := block.SourceAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-address "+v)
	}
	if v := block.TargetType.ValueString(); v != "" {
		if block.TargetValue.IsNull() {
			return configSet,
				pathRoot.AtName("target_type"),
				errors.New("target_value must be specified with target_type")
		}

		configSet = append(configSet, setPrefix+"target "+v+" \""+block.TargetValue.ValueString()+"\"")
	} else if !block.TargetValue.IsNull() {
		return configSet,
			pathRoot.AtName("target_value"),
			errors.New("target_type must be specified with target_value")
	}
	if !block.TestInterval.IsNull() {
		configSet = append(configSet, setPrefix+"test-interval "+
			utils.ConvI64toa(block.TestInterval.ValueInt64()))
	}
	for _, v := range block.Traps {
		configSet = append(configSet, setPrefix+"traps "+v.ValueString())
	}
	if !block.TTL.IsNull() {
		configSet = append(configSet, setPrefix+"ttl "+
			utils.ConvI64toa(block.TTL.ValueInt64()))
	}

	if block.RpmScale != nil {
		blockSet, pathErr, err := block.RpmScale.configSet(setPrefix, pathRoot.AtName("rpm_scale"))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if block.Thresholds != nil {
		configSet = append(configSet, block.Thresholds.configSet(setPrefix)...)
	}

	return configSet, path.Empty(), nil
}

func (block *servicesRpmProbeBlockTestBlockRpmScale) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "rpm-scale "

	configSet := []string{
		setPrefix + "tests-count " + utils.ConvI64toa(block.TestsCount.ValueInt64()),
	}

	if v := block.DestinationInterface.ValueString(); v != "" {
		if block.DestinationSubunitCnt.IsNull() {
			return configSet,
				pathRoot.AtName("destination_interface"),
				errors.New("destination_subunit_cnt must be specified with destination_interface")
		}
		configSet = append(configSet, setPrefix+"destination interface "+v)
		configSet = append(configSet, setPrefix+"destination subunit-cnt "+
			utils.ConvI64toa(block.DestinationSubunitCnt.ValueInt64()))
	} else if !block.DestinationSubunitCnt.IsNull() {
		return configSet,
			pathRoot.AtName("destination_subunit_cnt"),
			errors.New("destination_interface must be specified with destination_subunit_cnt")
	}
	if v := block.SourceAddressBase.ValueString(); v != "" {
		if block.SourceCount.IsNull() {
			return configSet,
				pathRoot.AtName("source_address_base"),
				errors.New("source_count must be specified with source_address_base")
		}
		if block.SourceStep.IsNull() {
			return configSet,
				pathRoot.AtName("source_address_base"),
				errors.New("source_step must be specified with source_address_base")
		}

		configSet = append(configSet, setPrefix+"source address-base "+v)
		configSet = append(configSet, setPrefix+"source count "+
			utils.ConvI64toa(block.SourceCount.ValueInt64()))
		configSet = append(configSet, setPrefix+"source step "+block.SourceStep.ValueString())
	} else {
		if !block.SourceCount.IsNull() {
			return configSet,
				pathRoot.AtName("source_count"),
				errors.New("source_address_base must be specified with source_count")
		}
		if !block.SourceStep.IsNull() {
			return configSet,
				pathRoot.AtName("source_step"),
				errors.New("source_address_base must be specified with source_step")
		}
	}
	if v := block.SourceInet6AddressBase.ValueString(); v != "" {
		if block.SourceInet6Count.IsNull() {
			return configSet,
				pathRoot.AtName("source_inet6_address_base"),
				errors.New("source_inet6_count must be specified with source_inet6_address_base")
		}
		if block.SourceInet6Step.IsNull() {
			return configSet,
				pathRoot.AtName("source_inet6_address_base"),
				errors.New("source_inet6_step must be specified with source_inet6_address_base")
		}

		configSet = append(configSet, setPrefix+"source-inet6 address-base "+v)
		configSet = append(configSet, setPrefix+"source-inet6 count "+
			utils.ConvI64toa(block.SourceInet6Count.ValueInt64()))
		configSet = append(configSet, setPrefix+"source-inet6 step "+block.SourceInet6Step.ValueString())
	} else {
		if !block.SourceInet6Count.IsNull() {
			return configSet,
				pathRoot.AtName("source_inet6_count"),
				errors.New("source_inet6_address_base must be specified with source_inet6_count")
		}
		if !block.SourceInet6Step.IsNull() {
			return configSet,
				pathRoot.AtName("source_inet6_step"),
				errors.New("source_inet6_address_base must be specified with source_inet6_step")
		}
	}
	if v := block.TargetAddressBase.ValueString(); v != "" {
		if block.TargetCount.IsNull() {
			return configSet,
				pathRoot.AtName("target_address_base"),
				errors.New("target_count must be specified with target_address_base")
		}
		if block.TargetStep.IsNull() {
			return configSet,
				pathRoot.AtName("target_address_base"),
				errors.New("target_step must be specified with target_address_base")
		}

		configSet = append(configSet, setPrefix+"target address-base "+v)
		configSet = append(configSet, setPrefix+"target count "+
			utils.ConvI64toa(block.TargetCount.ValueInt64()))
		configSet = append(configSet, setPrefix+"target step "+block.TargetStep.ValueString())
	} else {
		if !block.TargetCount.IsNull() {
			return configSet,
				pathRoot.AtName("target_count"),
				errors.New("target_address_base must be specified with target_count")
		}
		if !block.TargetStep.IsNull() {
			return configSet,
				pathRoot.AtName("target_step"),
				errors.New("target_address_base must be specified with target_step")
		}
	}
	if v := block.TargetInet6AddressBase.ValueString(); v != "" {
		if block.TargetInet6Count.IsNull() {
			return configSet,
				pathRoot.AtName("target_inet6_address_base"),
				errors.New("target_inet6_count must be specified with target_inet6_address_base")
		}
		if block.TargetInet6Step.IsNull() {
			return configSet,
				pathRoot.AtName("target_inet6_address_base"),
				errors.New("target_inet6_step must be specified with target_inet6_address_base")
		}

		configSet = append(configSet, setPrefix+"target-inet6 address-base "+v)
		configSet = append(configSet, setPrefix+"target-inet6 count "+
			utils.ConvI64toa(block.TargetInet6Count.ValueInt64()))
		configSet = append(configSet, setPrefix+"target-inet6 step "+block.TargetInet6Step.ValueString())
	} else {
		if !block.TargetInet6Count.IsNull() {
			return configSet,
				pathRoot.AtName("target_inet6_count"),
				errors.New("target_inet6_address_base must be specified with target_inet6_count")
		}
		if !block.TargetInet6Step.IsNull() {
			return configSet,
				pathRoot.AtName("target_inet6_step"),
				errors.New("target_inet6_address_base must be specified with target_inet6_step")
		}
	}

	return configSet, path.Empty(), nil
}

func (block *servicesRpmProbeBlockTestBlockThresholds) configSet(setPrefix string) []string {
	setPrefix += "thresholds "

	configSet := []string{
		setPrefix,
	}

	if !block.EgressTime.IsNull() {
		configSet = append(configSet, setPrefix+"egress-time "+
			utils.ConvI64toa(block.EgressTime.ValueInt64()))
	}
	if !block.IngressTime.IsNull() {
		configSet = append(configSet, setPrefix+"ingress-time "+
			utils.ConvI64toa(block.IngressTime.ValueInt64()))
	}
	if !block.JitterEgress.IsNull() {
		configSet = append(configSet, setPrefix+"jitter-egress "+
			utils.ConvI64toa(block.JitterEgress.ValueInt64()))
	}
	if !block.JitterIngress.IsNull() {
		configSet = append(configSet, setPrefix+"jitter-ingress "+
			utils.ConvI64toa(block.JitterIngress.ValueInt64()))
	}
	if !block.JitterRtt.IsNull() {
		configSet = append(configSet, setPrefix+"jitter-rtt "+
			utils.ConvI64toa(block.JitterRtt.ValueInt64()))
	}
	if !block.Rtt.IsNull() {
		configSet = append(configSet, setPrefix+"rtt "+
			utils.ConvI64toa(block.Rtt.ValueInt64()))
	}
	if !block.StdDevEgress.IsNull() {
		configSet = append(configSet, setPrefix+"std-dev-egress "+
			utils.ConvI64toa(block.StdDevEgress.ValueInt64()))
	}
	if !block.StdDevIngress.IsNull() {
		configSet = append(configSet, setPrefix+"std-dev-ingress "+
			utils.ConvI64toa(block.StdDevIngress.ValueInt64()))
	}
	if !block.StdDevRtt.IsNull() {
		configSet = append(configSet, setPrefix+"std-dev-rtt "+
			utils.ConvI64toa(block.StdDevRtt.ValueInt64()))
	}
	if !block.SuccessiveLoss.IsNull() {
		configSet = append(configSet, setPrefix+"successive-loss "+
			utils.ConvI64toa(block.SuccessiveLoss.ValueInt64()))
	}
	if !block.TotalLoss.IsNull() {
		configSet = append(configSet, setPrefix+"total-loss "+
			utils.ConvI64toa(block.TotalLoss.ValueInt64()))
	}

	return configSet
}

func (rscData *servicesRpmProbeData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services rpm probe \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.fillID()
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "delegate-probes":
				rscData.DelegateProbes = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "test "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				var test servicesRpmProbeBlockTest
				rscData.Test, test = tfdata.ExtractBlock(rscData.Test, types.StringValue(strings.Trim(name, "\"")))
				balt.CutPrefixInString(&itemTrim, name+" ")

				if err := test.read(itemTrim); err != nil {
					return err
				}
				rscData.Test = append(rscData.Test, test)
			}
		}
	}

	return nil
}

func (block *servicesRpmProbeBlockTest) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "data-fill "):
		block.DataFill = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "data-size "):
		block.DataSize, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination-interface "):
		block.DestinationInterface = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination-port "):
		block.DestinationPort, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "dscp-code-points "):
		block.DscpCodePoints = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "hardware-timestamp":
		block.HardwareTimestamp = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "history-size "):
		block.HistorySize, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "inet6-options source-address "):
		block.Inet6SourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "moving-average-size "):
		block.MovingAverageSize, err = tfdata.ConvAtoi64Value(itemTrim)
	case itemTrim == "one-way-hardware-timestamp":
		block.OneWayHardwareTimestamp = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "probe-count "):
		block.ProbeCount, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "probe-interval "):
		block.ProbeInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "probe-type "):
		block.ProbeType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "routing-instance "):
		block.RoutingInstance = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-address "):
		block.SourceAddress = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "target "):
		if len(strings.Split(itemTrim, " ")) < 2 { // <type> <value>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "target", itemTrim)
		}

		targetType := tfdata.FirstElementOfJunosLine(itemTrim)
		block.TargetType = types.StringValue(targetType)
		block.TargetValue = types.StringValue(strings.Trim(strings.TrimPrefix(itemTrim, targetType+" "), "\""))
	case balt.CutPrefixInString(&itemTrim, "traps "):
		block.Traps = append(block.Traps, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "ttl "):
		block.TTL, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "test-interval "):
		block.TestInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "rpm-scale "):
		if block.RpmScale == nil {
			block.RpmScale = &servicesRpmProbeBlockTestBlockRpmScale{}
		}

		err = block.RpmScale.read(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "thresholds"):
		if block.Thresholds == nil {
			block.Thresholds = &servicesRpmProbeBlockTestBlockThresholds{}
		}

		if balt.CutPrefixInString(&itemTrim, " ") {
			err = block.Thresholds.read(itemTrim)
		}
	}

	return err
}

func (block *servicesRpmProbeBlockTestBlockRpmScale) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "tests-count "):
		block.TestsCount, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination interface "):
		block.DestinationInterface = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination subunit-cnt "):
		block.DestinationSubunitCnt, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source address-base "):
		block.SourceAddressBase = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source count "):
		block.SourceCount, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source step "):
		block.SourceStep = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-inet6 address-base "):
		block.SourceInet6AddressBase = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-inet6 count "):
		block.SourceInet6Count, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-inet6 step "):
		block.SourceInet6Step = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "target address-base "):
		block.TargetAddressBase = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "target count "):
		block.TargetCount, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "target step "):
		block.TargetStep = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "target-inet6 address-base "):
		block.TargetInet6AddressBase = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "target-inet6 count "):
		block.TargetInet6Count, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "target-inet6 step "):
		block.TargetInet6Step = types.StringValue(itemTrim)
	}

	return err
}

func (block *servicesRpmProbeBlockTestBlockThresholds) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "egress-time "):
		block.EgressTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ingress-time "):
		block.IngressTime, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "jitter-egress "):
		block.JitterEgress, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "jitter-ingress "):
		block.JitterIngress, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "jitter-rtt "):
		block.JitterRtt, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "rtt "):
		block.Rtt, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "std-dev-egress "):
		block.StdDevEgress, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "std-dev-ingress "):
		block.StdDevIngress, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "std-dev-rtt "):
		block.StdDevRtt, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "successive-loss "):
		block.SuccessiveLoss, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "total-loss "):
		block.TotalLoss, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (rscData *servicesRpmProbeData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete services rpm probe \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
