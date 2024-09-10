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
	_ resource.Resource                   = &eventoptionsPolicy{}
	_ resource.ResourceWithConfigure      = &eventoptionsPolicy{}
	_ resource.ResourceWithValidateConfig = &eventoptionsPolicy{}
	_ resource.ResourceWithImportState    = &eventoptionsPolicy{}
	_ resource.ResourceWithUpgradeState   = &eventoptionsPolicy{}
)

type eventoptionsPolicy struct {
	client *junos.Client
}

func newEventoptionsPolicyResource() resource.Resource {
	return &eventoptionsPolicy{}
}

func (rsc *eventoptionsPolicy) typeName() string {
	return providerName + "_eventoptions_policy"
}

func (rsc *eventoptionsPolicy) junosName() string {
	return "event-options policy"
}

func (rsc *eventoptionsPolicy) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *eventoptionsPolicy) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *eventoptionsPolicy) Configure(
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

func (rsc *eventoptionsPolicy) Schema(
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
				Description: "Name of policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"events": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "List of events that trigger this policy.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"then": schema.SingleNestedBlock{
				Description: "Declare `then` configuration.",
				Attributes: map[string]schema.Attribute{
					"ignore": schema.BoolAttribute{
						Optional:    true,
						Description: "Do not log event or perform any other action.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"priority_override_facility": schema.StringAttribute{
						Optional:    true,
						Description: "Change syslog priority facility value.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"authorization",
								"change-log",
								"conflict-log",
								"daemon",
								"dfc",
								"external",
								"firewall",
								"ftp",
								"interactive-commands",
								"kernel",
								"ntp",
								"pfe",
								"security",
								"user",
							),
						},
					},
					"priority_override_severity": schema.StringAttribute{
						Optional:    true,
						Description: "Change syslog priority severity value.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"alert",
								"critical",
								"emergency",
								"error",
								"info",
								"notice",
								"warning",
							),
						},
					},
					"raise_trap": schema.BoolAttribute{
						Optional:    true,
						Description: "Raise SNMP trap.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"change_configuration": schema.SingleNestedBlock{
						Description: "Declare `change-configuration` configuration.",
						Attributes: map[string]schema.Attribute{
							"commands": schema.ListAttribute{
								ElementType: types.StringType,
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "List of configuration commands.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									),
								},
							},
							"commit_options_check": schema.BoolAttribute{
								Optional:    true,
								Description: "Check correctness of syntax; do not apply changes.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"commit_options_check_synchronize": schema.BoolAttribute{
								Optional:    true,
								Description: "Synchronize commit check on both Routing Engines.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"commit_options_force": schema.BoolAttribute{
								Optional:    true,
								Description: "Force commit on other Routing Engine (ignore warnings).",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"commit_options_log": schema.StringAttribute{
								Optional:    true,
								Description: "Message to write to commit log.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
							"commit_options_synchronize": schema.BoolAttribute{
								Optional:    true,
								Description: "Synchronize commit on both Routing Engines.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"retry_count": schema.Int64Attribute{
								Optional:    true,
								Description: "Change configuration retry attempt count.",
								Validators: []validator.Int64{
									int64validator.Between(0, 10),
								},
							},
							"retry_interval": schema.Int64Attribute{
								Optional:    true,
								Description: "Time interval between each retry (seconds).",
								Validators: []validator.Int64{
									int64validator.Between(0, 4294967295),
								},
							},
							"user_name": schema.StringAttribute{
								Optional:    true,
								Description: "User under whose privileges configuration should be changed.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"event_script": schema.ListNestedBlock{
						Description: "For each filename, invoke event scripts.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"filename": schema.StringAttribute{
									Required:    true,
									Description: "Local filename of the script file.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
										tfvalidator.StringSpaceExclusion(),
										tfvalidator.StringRuneExclusion('/', '%'),
									},
								},
								"output_filename": schema.StringAttribute{
									Optional:    true,
									Description: "Name of file in which to write event script output.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
										tfvalidator.StringSpaceExclusion(),
										tfvalidator.StringRuneExclusion('/', '%'),
									},
								},
								"output_format": schema.StringAttribute{
									Optional:    true,
									Description: "Format of output from event-script.",
									Validators: []validator.String{
										stringvalidator.OneOf("text", "xml"),
									},
								},
								"user_name": schema.StringAttribute{
									Optional:    true,
									Description: "User under whose privileges event script will execute.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
							},
							Blocks: map[string]schema.Block{
								"arguments": schema.ListNestedBlock{
									Description: "For each name of arguments, command line argument to the script.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Required:    true,
												Description: "Name of the argument.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
													tfvalidator.StringDoubleQuoteExclusion(),
												},
											},
											"value": schema.StringAttribute{
												Required:    true,
												Description: "Value of the argument.",
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
													tfvalidator.StringDoubleQuoteExclusion(),
												},
											},
										},
									},
								},
								"destination": schema.SingleNestedBlock{
									Description: "Location to which to upload event script output.",
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Required:    false, // true when SingleNestedBlock is specified
											Optional:    true,
											Description: "Destination name.",
											Validators: []validator.String{
												stringvalidator.LengthBetween(1, 250),
												tfvalidator.StringDoubleQuoteExclusion(),
											},
										},
										"retry_count": schema.Int64Attribute{
											Optional:    true,
											Description: "Upload output-filename retry attempt count.",
											Validators: []validator.Int64{
												int64validator.Between(0, 10),
											},
										},
										"retry_interval": schema.Int64Attribute{
											Optional:    true,
											Description: "Time interval between each retry (seconds).",
											Validators: []validator.Int64{
												int64validator.Between(0, 4294967295),
											},
										},
										"transfer_delay": schema.Int64Attribute{
											Optional:    true,
											Description: "Delay before uploading files (seconds).",
											Validators: []validator.Int64{
												int64validator.Between(0, 4294967295),
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
					"execute_commands": schema.SingleNestedBlock{
						Description: "Issue one or more CLI commands.",
						Attributes: map[string]schema.Attribute{
							"commands": schema.ListAttribute{
								ElementType: types.StringType,
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "List of CLI commands to issue.",
								Validators: []validator.List{
									listvalidator.SizeAtLeast(1),
									listvalidator.ValueStringsAre(
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									),
								},
							},
							"output_filename": schema.StringAttribute{
								Optional:    true,
								Description: "Name of file in which to write command output.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 250),
									tfvalidator.StringDoubleQuoteExclusion(),
									tfvalidator.StringSpaceExclusion(),
									tfvalidator.StringRuneExclusion('/', '%'),
								},
							},
							"output_format": schema.StringAttribute{
								Optional:    true,
								Description: "Format of output from CLI commands.",
								Validators: []validator.String{
									stringvalidator.OneOf("text", "xml"),
								},
							},
							"user_name": schema.StringAttribute{
								Optional:    true,
								Description: "User under whose privileges command will execute.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								},
							},
						},
						Blocks: map[string]schema.Block{
							"destination": schema.SingleNestedBlock{
								Description: "Location to which to upload command output.",
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Required:    false, // true when SingleNestedBlock is specified
										Optional:    true,
										Description: "Destination name.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 250),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"retry_count": schema.Int64Attribute{
										Optional:    true,
										Description: "Upload output-filename retry attempt count.",
										Validators: []validator.Int64{
											int64validator.Between(0, 10),
										},
									},
									"retry_interval": schema.Int64Attribute{
										Optional:    true,
										Description: "Time interval between each retry (seconds).",
										Validators: []validator.Int64{
											int64validator.Between(0, 4294967295),
										},
									},
									"transfer_delay": schema.Int64Attribute{
										Optional:    true,
										Description: "Delay before uploading file to the destination (seconds).",
										Validators: []validator.Int64{
											int64validator.Between(0, 4294967295),
										},
									},
								},
								PlanModifiers: []planmodifier.Object{
									tfplanmodifier.BlockRemoveNull(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
						},
					},
					"upload": schema.ListNestedBlock{
						Description: "For each combination of `filename` and `destination` arguments," +
							" upload file to specified destination.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"filename": schema.StringAttribute{
									Required:    true,
									Description: "Name of file to upload.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
										tfvalidator.StringSpaceExclusion(),
										tfvalidator.StringRuneExclusion('/', '%'),
									},
								},
								"destination": schema.StringAttribute{
									Required:    true,
									Description: "Location to which to output file.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 250),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"retry_count": schema.Int64Attribute{
									Optional:    true,
									Description: "Upload output-filename retry attempt count.",
									Validators: []validator.Int64{
										int64validator.Between(0, 10),
									},
								},
								"retry_interval": schema.Int64Attribute{
									Optional:    true,
									Description: "Time interval between each retry (seconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 4294967295),
									},
								},
								"transfer_delay": schema.Int64Attribute{
									Optional:    true,
									Description: "Delay before uploading file to the destination (seconds).",
									Validators: []validator.Int64{
										int64validator.Between(0, 4294967295),
									},
								},
								"user_name": schema.StringAttribute{
									Optional:    true,
									Description: "User under whose privileges upload action will execute.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"attributes_match": schema.ListNestedBlock{
				Description: "For each combination of block arguments, attributes to compare for two events.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"from": schema.StringAttribute{
							Required:    true,
							Description: "First attribute to compare.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"compare": schema.StringAttribute{
							Required:    true,
							Description: "Type to compare.",
							Validators: []validator.String{
								stringvalidator.OneOf("equals", "matches", "starts-with"),
							},
						},
						"to": schema.StringAttribute{
							Required:    true,
							Description: "Second attribute or value to compare.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
				},
			},
			"within": schema.ListNestedBlock{
				Description: "For each time interval, list of events correlated with triggering events.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"time_interval": schema.Int64Attribute{
							Required:    true,
							Description: "Time within which correlated events must occur (or not) (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(1, 604800),
							},
						},
						"events": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of events that must occur within time interval.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								),
							},
						},
						"not_events": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of events must not occur within time interval.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								),
							},
						},
						"trigger_count": schema.Int64Attribute{
							Optional:    true,
							Description: " Number of occurrences of triggering event.",
							Validators: []validator.Int64{
								int64validator.Between(0, 4294967295),
							},
						},
						"trigger_when": schema.StringAttribute{
							Optional:    true,
							Description: "To compare with `trigger_count`.",
							Validators: []validator.String{
								stringvalidator.OneOf("after", "on", "until"),
							},
						},
					},
				},
			},
		},
	}
}

type eventoptionsPolicyData struct {
	ID              types.String                             `tfsdk:"id"`
	Name            types.String                             `tfsdk:"name"`
	Events          []types.String                           `tfsdk:"events"`
	Then            *eventoptionsPolicyBlockThen             `tfsdk:"then"`
	AttributesMatch []eventoptionsPolicyBlockAttributesMatch `tfsdk:"attributes_match"`
	Within          []eventoptionsPolicyBlockWithin          `tfsdk:"within"`
}

type eventoptionsPolicyConfig struct {
	ID              types.String                       `tfsdk:"id"`
	Name            types.String                       `tfsdk:"name"`
	Events          types.Set                          `tfsdk:"events"`
	Then            *eventoptionsPolicyBlockThenConfig `tfsdk:"then"`
	AttributesMatch types.List                         `tfsdk:"attributes_match"`
	Within          types.List                         `tfsdk:"within"`
}

type eventoptionsPolicyBlockThen struct {
	Ignore                   types.Bool                                          `tfsdk:"ignore"`
	PriorityOverrideFacility types.String                                        `tfsdk:"priority_override_facility"`
	PriorityOverrideSeverity types.String                                        `tfsdk:"priority_override_severity"`
	RaiseTrap                types.Bool                                          `tfsdk:"raise_trap"`
	ChangeConfiguration      *eventoptionsPolicyBlockThenBlockChangeConfigurtion `tfsdk:"change_configuration"`
	EventScript              []eventoptionsPolicyBlockThenBlockEventScript       `tfsdk:"event_script"`
	ExecuteCommands          *eventoptionsPolicyBlockThenBlockExecuteCommands    `tfsdk:"execute_commands"`
	Upload                   []eventoptionsPolicyBlockThenBlockUpload            `tfsdk:"upload"`
}

func (block *eventoptionsPolicyBlockThen) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type eventoptionsPolicyBlockThenConfig struct {
	Ignore                   types.Bool                                                `tfsdk:"ignore"`
	PriorityOverrideFacility types.String                                              `tfsdk:"priority_override_facility"`
	PriorityOverrideSeverity types.String                                              `tfsdk:"priority_override_severity"`
	RaiseTrap                types.Bool                                                `tfsdk:"raise_trap"`
	ChangeConfiguration      *eventoptionsPolicyBlockThenBlockChangeConfigurtionConfig `tfsdk:"change_configuration"`
	EventScript              types.List                                                `tfsdk:"event_script"`
	ExecuteCommands          *eventoptionsPolicyBlockThenBlockExecuteCommandsConfig    `tfsdk:"execute_commands"`
	Upload                   types.List                                                `tfsdk:"upload"`
}

func (block *eventoptionsPolicyBlockThenConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type eventoptionsPolicyBlockThenBlockChangeConfigurtion struct {
	Commands                      []types.String `tfsdk:"commands"`
	CommitOptionsCheck            types.Bool     `tfsdk:"commit_options_check"`
	CommitOptionsCheckSynchronize types.Bool     `tfsdk:"commit_options_check_synchronize"`
	CommitOptionsForce            types.Bool     `tfsdk:"commit_options_force"`
	CommitOptionsLog              types.String   `tfsdk:"commit_options_log"`
	CommitOptionsSynchronize      types.Bool     `tfsdk:"commit_options_synchronize"`
	RetryCount                    types.Int64    `tfsdk:"retry_count"`
	RetryInterval                 types.Int64    `tfsdk:"retry_interval"`
	Username                      types.String   `tfsdk:"user_name"`
}

type eventoptionsPolicyBlockThenBlockChangeConfigurtionConfig struct {
	Commands                      types.List   `tfsdk:"commands"`
	CommitOptionsCheck            types.Bool   `tfsdk:"commit_options_check"`
	CommitOptionsCheckSynchronize types.Bool   `tfsdk:"commit_options_check_synchronize"`
	CommitOptionsForce            types.Bool   `tfsdk:"commit_options_force"`
	CommitOptionsLog              types.String `tfsdk:"commit_options_log"`
	CommitOptionsSynchronize      types.Bool   `tfsdk:"commit_options_synchronize"`
	RetryCount                    types.Int64  `tfsdk:"retry_count"`
	RetryInterval                 types.Int64  `tfsdk:"retry_interval"`
	Username                      types.String `tfsdk:"user_name"`
}

type eventoptionsPolicyBlockThenBlockEventScript struct {
	Filename       types.String                                                `tfsdk:"filename"`
	OutputFilename types.String                                                `tfsdk:"output_filename"`
	OutputFormat   types.String                                                `tfsdk:"output_format"`
	Username       types.String                                                `tfsdk:"user_name"`
	Arguments      []eventoptionsPolicyBlockThenBlockEventScriptBlockArguments `tfsdk:"arguments"`
	Destination    *eventoptionsPolicyBlockThenBlockDestination                `tfsdk:"destination"`
}

type eventoptionsPolicyBlockThenBlockEventScriptConfig struct {
	Filename       types.String                                 `tfsdk:"filename"`
	OutputFilename types.String                                 `tfsdk:"output_filename"`
	OutputFormat   types.String                                 `tfsdk:"output_format"`
	Username       types.String                                 `tfsdk:"user_name"`
	Arguments      types.List                                   `tfsdk:"arguments"`
	Destination    *eventoptionsPolicyBlockThenBlockDestination `tfsdk:"destination"`
}

type eventoptionsPolicyBlockThenBlockEventScriptBlockArguments struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type eventoptionsPolicyBlockThenBlockExecuteCommands struct {
	Commands       []types.String                               `tfsdk:"commands"`
	OutputFilename types.String                                 `tfsdk:"output_filename"`
	OutputFormat   types.String                                 `tfsdk:"output_format"`
	Username       types.String                                 `tfsdk:"user_name"`
	Destination    *eventoptionsPolicyBlockThenBlockDestination `tfsdk:"destination"`
}

type eventoptionsPolicyBlockThenBlockExecuteCommandsConfig struct {
	Commands       types.List                                   `tfsdk:"commands"`
	OutputFilename types.String                                 `tfsdk:"output_filename"`
	OutputFormat   types.String                                 `tfsdk:"output_format"`
	Username       types.String                                 `tfsdk:"user_name"`
	Destination    *eventoptionsPolicyBlockThenBlockDestination `tfsdk:"destination"`
}

type eventoptionsPolicyBlockThenBlockDestination struct {
	Name          types.String `tfsdk:"name"`
	RetryCount    types.Int64  `tfsdk:"retry_count"`
	RetryInterval types.Int64  `tfsdk:"retry_interval"`
	TransferDelay types.Int64  `tfsdk:"transfer_delay"`
}

func (block *eventoptionsPolicyBlockThenBlockDestination) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type eventoptionsPolicyBlockThenBlockUpload struct {
	Filename      types.String `tfsdk:"filename"`
	Destination   types.String `tfsdk:"destination"`
	RetryCount    types.Int64  `tfsdk:"retry_count"`
	RetryInterval types.Int64  `tfsdk:"retry_interval"`
	TransferDelay types.Int64  `tfsdk:"transfer_delay"`
	Username      types.String `tfsdk:"user_name"`
}

type eventoptionsPolicyBlockAttributesMatch struct {
	From    types.String `tfsdk:"from"`
	Compare types.String `tfsdk:"compare"`
	To      types.String `tfsdk:"to"`
}

type eventoptionsPolicyBlockWithin struct {
	TimeInterval types.Int64    `tfsdk:"time_interval"`
	Events       []types.String `tfsdk:"events"`
	NotEvents    []types.String `tfsdk:"not_events"`
	TriggerCount types.Int64    `tfsdk:"trigger_count"`
	TriggerWhen  types.String   `tfsdk:"trigger_when"`
}

type eventoptionsPolicyBlockWithinConfig struct {
	TimeInterval types.Int64  `tfsdk:"time_interval"`
	Events       types.Set    `tfsdk:"events"`
	NotEvents    types.Set    `tfsdk:"not_events"`
	TriggerCount types.Int64  `tfsdk:"trigger_count"`
	TriggerWhen  types.String `tfsdk:"trigger_when"`
}

func (block *eventoptionsPolicyBlockWithin) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block, "TimeInterval")
}

//nolint:gocognit
func (rsc *eventoptionsPolicy) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config eventoptionsPolicyConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Then != nil {
		if config.Then.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("then"),
				tfdiag.MissingConfigErrSummary,
				"then block is empty",
			)
		} else {
			if !config.Then.Ignore.IsNull() && !config.Then.Ignore.IsUnknown() {
				if tfdata.CheckBlockHasKnownValue(&config.Then, "Ignore") {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("ignore"),
						tfdiag.ConflictConfigErrSummary,
						"ignore must be specified alone in then block",
					)
				}
			}
			if config.Then.ChangeConfiguration != nil {
				if config.Then.ChangeConfiguration.Commands.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("change_configuration").AtName("commands"),
						tfdiag.MissingConfigErrSummary,
						"commands must be specified in change_configuration block in then block",
					)
				}
				if !config.Then.ChangeConfiguration.CommitOptionsCheckSynchronize.IsNull() &&
					!config.Then.ChangeConfiguration.CommitOptionsCheckSynchronize.IsUnknown() &&
					config.Then.ChangeConfiguration.CommitOptionsCheck.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("change_configuration").AtName("commit_options_check_synchronize"),
						tfdiag.MissingConfigErrSummary,
						"commit_options_check must be specified with commit_options_check_synchronize"+
							" in change_configuration block in then block",
					)
				}
				if !config.Then.ChangeConfiguration.CommitOptionsForce.IsNull() &&
					!config.Then.ChangeConfiguration.CommitOptionsForce.IsUnknown() &&
					!config.Then.ChangeConfiguration.CommitOptionsCheck.IsNull() &&
					!config.Then.ChangeConfiguration.CommitOptionsCheck.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("change_configuration").AtName("commit_options_force"),
						tfdiag.ConflictConfigErrSummary,
						"commit_options_force and commit_options_check cannot be configured together",
					)
				}
				if !config.Then.ChangeConfiguration.CommitOptionsLog.IsNull() &&
					!config.Then.ChangeConfiguration.CommitOptionsLog.IsUnknown() &&
					!config.Then.ChangeConfiguration.CommitOptionsCheck.IsNull() &&
					!config.Then.ChangeConfiguration.CommitOptionsCheck.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("change_configuration").AtName("commit_options_log"),
						tfdiag.ConflictConfigErrSummary,
						"commit_options_log and commit_options_check cannot be configured together",
					)
				}
				if !config.Then.ChangeConfiguration.CommitOptionsSynchronize.IsNull() &&
					!config.Then.ChangeConfiguration.CommitOptionsSynchronize.IsUnknown() &&
					!config.Then.ChangeConfiguration.CommitOptionsCheck.IsNull() &&
					!config.Then.ChangeConfiguration.CommitOptionsCheck.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("change_configuration").AtName("commit_options_synchronize"),
						tfdiag.ConflictConfigErrSummary,
						"commit_options_synchronize and commit_options_check cannot be configured together",
					)
				}
				if !config.Then.ChangeConfiguration.RetryCount.IsNull() &&
					!config.Then.ChangeConfiguration.RetryCount.IsUnknown() &&
					config.Then.ChangeConfiguration.RetryInterval.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("change_configuration").AtName("retry_count"),
						tfdiag.MissingConfigErrSummary,
						"retry_interval must be specified with retry_count"+
							" in change_configuration block in then block",
					)
				}
				if !config.Then.ChangeConfiguration.RetryInterval.IsNull() &&
					!config.Then.ChangeConfiguration.RetryInterval.IsUnknown() &&
					config.Then.ChangeConfiguration.RetryCount.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("change_configuration").AtName("retry_interval"),
						tfdiag.MissingConfigErrSummary,
						"retry_count must be specified with retry_interval"+
							" in change_configuration block in then block",
					)
				}
			}
			if !config.Then.EventScript.IsNull() && !config.Then.EventScript.IsUnknown() {
				var configEventScript []eventoptionsPolicyBlockThenBlockEventScriptConfig
				asDiags := config.Then.EventScript.ElementsAs(ctx, &configEventScript, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				eventScriptFilename := make(map[string]struct{})
				for i, block := range configEventScript {
					if !block.Filename.IsUnknown() {
						filename := block.Filename.ValueString()
						if _, ok := eventScriptFilename[filename]; ok {
							resp.Diagnostics.AddAttributeError(
								path.Root("then").AtName("event_script").AtListIndex(i).AtName("filename"),
								tfdiag.DuplicateConfigErrSummary,
								fmt.Sprintf("multiple event_script blocks with the same filename %q", filename),
							)
						}
						eventScriptFilename[filename] = struct{}{}
					}
					if block.Destination != nil {
						if block.Destination.Name.IsNull() {
							resp.Diagnostics.AddAttributeError(
								path.Root("then").AtName("event_script").AtListIndex(i).AtName("destination").AtName("name"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("name must be specified"+
									" in destination block in event_script block %q in then block", block.Filename.ValueString()),
							)
						}
						if !block.Destination.RetryCount.IsNull() &&
							!block.Destination.RetryCount.IsUnknown() &&
							block.Destination.RetryInterval.IsNull() {
							resp.Diagnostics.AddAttributeError(
								path.Root("then").AtName("event_script").AtListIndex(i).AtName("destination").AtName("retry_count"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("retry_interval must be specified with retry_count"+
									" in destination block in event_script block %q in then block", block.Filename.ValueString()),
							)
						}
						if !block.Destination.RetryInterval.IsNull() &&
							!block.Destination.RetryInterval.IsUnknown() &&
							block.Destination.RetryCount.IsNull() {
							resp.Diagnostics.AddAttributeError(
								path.Root("then").AtName("execute_commands").AtName("destination").AtName("retry_interval"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("retry_count must be specified with retry_interval"+
									" in destination block in event_script block %q in then block", block.Filename.ValueString()),
							)
						}
					}
					if !block.Arguments.IsNull() && !block.Arguments.IsUnknown() {
						var configArguments []eventoptionsPolicyBlockThenBlockEventScriptBlockArguments
						asDiags := block.Arguments.ElementsAs(ctx, &configArguments, false)
						if asDiags.HasError() {
							resp.Diagnostics.Append(asDiags...)

							return
						}

						argumentsName := make(map[string]struct{})
						for ii, block2 := range configArguments {
							if block2.Name.IsUnknown() {
								continue
							}
							name := block2.Name.ValueString()
							if _, ok := argumentsName[name]; ok {
								resp.Diagnostics.AddAttributeError(
									path.Root("then").AtName("event_script").AtListIndex(i).AtName("arguments").AtListIndex(ii).AtName("name"),
									tfdiag.DuplicateConfigErrSummary,
									fmt.Sprintf("multiple arguments blocks with the same name %q"+
										" in event_script block %q in then block", name, block.Filename.ValueString()),
								)
							}
							argumentsName[name] = struct{}{}
						}
					}
				}
			}
			if config.Then.ExecuteCommands != nil {
				if config.Then.ExecuteCommands.Commands.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("execute_commands").AtName("commands"),
						tfdiag.MissingConfigErrSummary,
						"commands must be specified in execute_commands block in then block",
					)
				}
				if !config.Then.ExecuteCommands.OutputFilename.IsNull() &&
					!config.Then.ExecuteCommands.OutputFilename.IsUnknown() &&
					config.Then.ExecuteCommands.Destination == nil {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("execute_commands").AtName("output_filename"),
						tfdiag.MissingConfigErrSummary,
						"destination must be specified with output_filename"+
							" in execute_commands block in then block",
					)
				}

				if config.Then.ExecuteCommands.Destination != nil {
					if config.Then.ExecuteCommands.Destination.hasKnownValue() &&
						config.Then.ExecuteCommands.OutputFilename.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("then").AtName("execute_commands").AtName("destination"),
							tfdiag.MissingConfigErrSummary,
							"output_filename must be specified with destination"+
								" in execute_commands block in then block",
						)
					}
					if config.Then.ExecuteCommands.Destination.Name.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("then").AtName("execute_commands").AtName("destination").AtName("name"),
							tfdiag.MissingConfigErrSummary,
							"name must be specified"+
								" in destination block in execute_commands block in then block",
						)
					}
					if !config.Then.ExecuteCommands.Destination.RetryCount.IsNull() &&
						!config.Then.ExecuteCommands.Destination.RetryCount.IsUnknown() &&
						config.Then.ExecuteCommands.Destination.RetryInterval.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("then").AtName("execute_commands").AtName("destination").AtName("retry_count"),
							tfdiag.MissingConfigErrSummary,
							"retry_interval must be specified with retry_count"+
								" in destination block in execute_commands block in then block",
						)
					}
					if !config.Then.ExecuteCommands.Destination.RetryInterval.IsNull() &&
						!config.Then.ExecuteCommands.Destination.RetryInterval.IsUnknown() &&
						config.Then.ExecuteCommands.Destination.RetryCount.IsNull() {
						resp.Diagnostics.AddAttributeError(
							path.Root("then").AtName("execute_commands").AtName("destination").AtName("retry_interval"),
							tfdiag.MissingConfigErrSummary,
							"retry_count must be specified with retry_interval"+
								" in destination block in execute_commands block in then block",
						)
					}
				}
			}
			if !config.Then.Upload.IsNull() && !config.Then.Upload.IsUnknown() {
				var configUpload []eventoptionsPolicyBlockThenBlockUpload
				asDiags := config.Then.Upload.ElementsAs(ctx, &configUpload, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				uploadFilenameDestination := make(map[string]struct{})
				for i, block := range configUpload {
					if block.Filename.IsUnknown() || block.Destination.IsUnknown() {
						continue
					}
					filename := block.Filename.ValueString()
					destination := block.Destination.ValueString()
					if _, ok := uploadFilenameDestination[filename+junos.IDSeparator+destination]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("then").AtName("upload").AtListIndex(i).AtName("filename"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple upload blocks with the same filename %q and destination %q", filename, destination),
						)
					}
					uploadFilenameDestination[filename+junos.IDSeparator+destination] = struct{}{}
				}
			}
		}
	} else {
		resp.Diagnostics.AddAttributeError(
			path.Root("then"),
			tfdiag.MissingConfigErrSummary,
			"then block must be specified",
		)
	}
	if !config.AttributesMatch.IsNull() && !config.AttributesMatch.IsUnknown() {
		var configAttributesMatch []eventoptionsPolicyBlockAttributesMatch
		asDiags := config.AttributesMatch.ElementsAs(ctx, &configAttributesMatch, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		attributesMatchArgs := make(map[string]struct{})
		for i, block := range configAttributesMatch {
			if block.From.IsUnknown() || block.Compare.IsUnknown() || block.To.IsUnknown() {
				continue
			}
			args := "\"" + block.From.ValueString() + "\" " + block.Compare.ValueString() + " \"" + block.To.ValueString() + "\""
			if _, ok := attributesMatchArgs[args]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("attributes_match").AtListIndex(i).AtName("from"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple attributes_match blocks with the same from %q, compare %q and to %q",
						block.From.ValueString(), block.Compare.ValueString(), block.To.ValueString()),
				)
			}
			attributesMatchArgs[args] = struct{}{}
		}
	}
	if !config.Within.IsNull() && !config.Within.IsUnknown() {
		var configWithin []eventoptionsPolicyBlockWithinConfig
		asDiags := config.Within.ElementsAs(ctx, &configWithin, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		withinTimeInterval := make(map[int64]struct{})
		for i, block := range configWithin {
			if !block.TimeInterval.IsNull() {
				timeInterval := block.TimeInterval.ValueInt64()
				if _, ok := withinTimeInterval[timeInterval]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("within").AtListIndex(i).AtName("time_interval"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple within blocks with the same time_interval %d", timeInterval),
					)
				}
				withinTimeInterval[timeInterval] = struct{}{}
			}
			if !block.TriggerWhen.IsNull() &&
				!block.TriggerWhen.IsUnknown() &&
				block.TriggerCount.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("within").AtListIndex(i).AtName("trigger_when"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("trigger_count must be specified with trigger_when"+
						" in within block %d", block.TimeInterval.ValueInt64()),
				)
			}
			if !block.TriggerCount.IsNull() &&
				!block.TriggerCount.IsUnknown() &&
				block.TriggerWhen.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("within").AtListIndex(i).AtName("trigger_count"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("trigger_when must be specified with trigger_count"+
						" in within block %d", block.TimeInterval.ValueInt64()),
				)
			}
		}
	}
}

func (rsc *eventoptionsPolicy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan eventoptionsPolicyData
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
			policyExists, err := checkEventoptionsPolicyExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if policyExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			policyExists, err := checkEventoptionsPolicyExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !policyExists {
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

func (rsc *eventoptionsPolicy) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data eventoptionsPolicyData
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

func (rsc *eventoptionsPolicy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state eventoptionsPolicyData
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

func (rsc *eventoptionsPolicy) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state eventoptionsPolicyData
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

func (rsc *eventoptionsPolicy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data eventoptionsPolicyData

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

func checkEventoptionsPolicyExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"event-options policy \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *eventoptionsPolicyData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *eventoptionsPolicyData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *eventoptionsPolicyData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set event-options policy \"" + rscData.Name.ValueString() + "\" "

	for _, v := range rscData.Events {
		configSet = append(configSet, setPrefix+"events \""+v.ValueString()+"\"")
	}
	if rscData.Then != nil {
		if rscData.Then.isEmpty() {
			return path.Root("then").AtName("*"),
				errors.New("then block is empty")
		}
		blockSet, pathErr, err := rscData.Then.configSet(setPrefix, path.Root("then"))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	attributesMatchArgs := make(map[string]struct{})
	for i, block := range rscData.AttributesMatch {
		setLine := "attributes-match" +
			" \"" + block.From.ValueString() + "\"" +
			" " + block.Compare.ValueString() +
			" \"" + block.To.ValueString() + "\""
		if _, ok := attributesMatchArgs[setLine]; ok {
			return path.Root("attributes_match").AtListIndex(i).AtName("from"),
				fmt.Errorf("multiple attributes_match blocks with the same from %q, compare %q and to %q",
					block.From.ValueString(), block.Compare.ValueString(), block.To.ValueString(),
				)
		}
		attributesMatchArgs[setLine] = struct{}{}

		configSet = append(configSet, setPrefix+setLine)
	}
	withinTimeInterval := make(map[int64]struct{})
	for i, block := range rscData.Within {
		timeInterval := block.TimeInterval.ValueInt64()
		if _, ok := withinTimeInterval[timeInterval]; ok {
			return path.Root("within").AtListIndex(i).AtName("time_interval"),
				fmt.Errorf("multiple within blocks with the same time_interval %d", timeInterval)
		}
		withinTimeInterval[timeInterval] = struct{}{}
		if block.isEmpty() {
			return path.Root("within").AtListIndex(i).AtName("time_interval"),
				fmt.Errorf("within block %q is empty", timeInterval)
		}

		blockSet, pathErr, err := block.configSet(setPrefix, path.Root("within").AtListIndex(i))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *eventoptionsPolicyBlockThen) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "then "

	if block.Ignore.ValueBool() {
		configSet = append(configSet, setPrefix+"ignore")
	}
	if v := block.PriorityOverrideFacility.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"priority-override facility "+v)
	}
	if v := block.PriorityOverrideSeverity.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"priority-override severity "+v)
	}
	if block.RaiseTrap.ValueBool() {
		configSet = append(configSet, setPrefix+"raise-trap")
	}
	if block.ChangeConfiguration != nil {
		blockSet, pathErr, err := block.ChangeConfiguration.configSet(setPrefix, pathRoot.AtName("change_configuration"))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	eventScriptFilename := make(map[string]struct{})
	for i, eventScript := range block.EventScript {
		fileName := eventScript.Filename.ValueString()
		if _, ok := eventScriptFilename[fileName]; ok {
			return configSet,
				pathRoot.AtName("event_script").AtListIndex(i).AtName("filename"),
				fmt.Errorf("multiple event_script blocks with the same filename %q", fileName)
		}
		eventScriptFilename[fileName] = struct{}{}

		blockSet, pathErr, err := eventScript.configSet(setPrefix, pathRoot.AtName("event_script").AtListIndex(i))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if block.ExecuteCommands != nil {
		blockSet, pathErr, err := block.ExecuteCommands.configSet(setPrefix, pathRoot.AtName("execute_commands"))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	uploadFilenameDestination := make(map[string]struct{})
	for i, upload := range block.Upload {
		filename := upload.Filename.ValueString()
		destination := upload.Destination.ValueString()
		if _, ok := uploadFilenameDestination[filename+junos.IDSeparator+destination]; ok {
			return configSet,
				pathRoot.AtName("upload").AtListIndex(i).AtName("filename"),
				fmt.Errorf("multiple upload blocks with the same filename %q and destination %q", filename, destination)
		}
		uploadFilenameDestination[filename+junos.IDSeparator+destination] = struct{}{}

		blockSet, pathErr, err := upload.configSet(setPrefix, pathRoot.AtName("upload").AtListIndex(i))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return configSet, path.Empty(), nil
}

func (block *eventoptionsPolicyBlockThenBlockChangeConfigurtion) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, len(block.Commands))
	setPrefix += "change-configuration "

	for _, v := range block.Commands {
		configSet = append(configSet, setPrefix+"commands \""+v.ValueString()+"\"")
	}
	if block.CommitOptionsCheck.ValueBool() {
		configSet = append(configSet, setPrefix+"commit-options check")
		if block.CommitOptionsCheckSynchronize.ValueBool() {
			configSet = append(configSet, setPrefix+"commit-options check synchronize")
		}
	} else if block.CommitOptionsCheckSynchronize.ValueBool() {
		return configSet,
			pathRoot.AtName("commit_options_check_synchronize"),
			errors.New("commit_options_check must be specified with commit_options_check_synchronize")
	}
	if block.CommitOptionsForce.ValueBool() {
		configSet = append(configSet, setPrefix+"commit-options force")
	}
	if v := block.CommitOptionsLog.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"commit-options log \""+v+"\"")
	}
	if block.CommitOptionsSynchronize.ValueBool() {
		configSet = append(configSet, setPrefix+"commit-options synchronize")
	}
	if !block.RetryCount.IsNull() {
		if block.RetryInterval.IsNull() {
			return configSet,
				pathRoot.AtName("retry_count"),
				errors.New("retry_interval must be specified with retry_count")
		}
		configSet = append(configSet, setPrefix+"retry count "+
			utils.ConvI64toa(block.RetryCount.ValueInt64())+
			" interval "+utils.ConvI64toa(block.RetryInterval.ValueInt64()))
	} else if !block.RetryInterval.IsNull() {
		return configSet,
			pathRoot.AtName("retry_interval"),
			errors.New("retry_count must be specified with retry_interval")
	}
	if v := block.Username.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"user-name \""+v+"\"")
	}

	return configSet, path.Empty(), nil
}

func (block *eventoptionsPolicyBlockThenBlockEventScript) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "event-script \"" + block.Filename.ValueString() + "\" "

	configSet := []string{
		setPrefix,
	}

	if v := block.OutputFilename.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"output-filename \""+v+"\"")
	}
	if v := block.OutputFormat.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"output-format "+v)
	}
	if v := block.Username.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"user-name \""+v+"\"")
	}
	argumentsName := make(map[string]struct{})
	for i, arguments := range block.Arguments {
		name := arguments.Name.ValueString()
		if _, ok := argumentsName[name]; ok {
			return configSet,
				pathRoot.AtName("arguments").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple arguments blocks with the same name %q", name)
		}
		argumentsName[name] = struct{}{}

		configSet = append(configSet, setPrefix+"arguments \""+name+"\" \""+arguments.Value.ValueString()+"\"")
	}
	if block.Destination != nil {
		blockSet, pathErr, err := block.Destination.configSet(setPrefix, pathRoot.AtName("destination"))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return configSet, path.Empty(), nil
}

func (block *eventoptionsPolicyBlockThenBlockExecuteCommands) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0, len(block.Commands))
	setPrefix += "execute-commands "

	for _, v := range block.Commands {
		configSet = append(configSet, setPrefix+"commands \""+v.ValueString()+"\"")
	}
	if v := block.OutputFilename.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"output-filename \""+v+"\"")
	}
	if v := block.OutputFormat.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"output-format "+v)
	}
	if v := block.Username.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"user-name \""+v+"\"")
	}
	if block.Destination != nil {
		blockSet, pathErr, err := block.Destination.configSet(setPrefix, pathRoot.AtName("destination"))
		if err != nil {
			return configSet, pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return configSet, path.Empty(), nil
}

func (block *eventoptionsPolicyBlockThenBlockUpload) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "upload filename \"" + block.Filename.ValueString() + "\"" +
		" destination \"" + block.Destination.ValueString() + "\" "

	configSet := []string{
		setPrefix,
	}

	if !block.RetryCount.IsNull() {
		if block.RetryInterval.IsNull() {
			return configSet,
				pathRoot.AtName("retry_count"),
				errors.New("retry_interval must be specified with retry_count")
		}
		configSet = append(configSet, setPrefix+"retry-count "+
			utils.ConvI64toa(block.RetryCount.ValueInt64())+
			" retry-interval "+utils.ConvI64toa(block.RetryInterval.ValueInt64()))
	} else if !block.RetryInterval.IsNull() {
		return configSet,
			pathRoot.AtName("retry_interval"),
			errors.New("retry_count must be specified with retry_interval")
	}
	if !block.TransferDelay.IsNull() {
		configSet = append(configSet, setPrefix+"transfer-delay "+
			utils.ConvI64toa(block.TransferDelay.ValueInt64()))
	}
	if v := block.Username.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"user-name \""+v+"\"")
	}

	return configSet, path.Empty(), nil
}

func (block *eventoptionsPolicyBlockThenBlockDestination) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "destination \"" + block.Name.ValueString() + "\" "

	configSet := []string{
		setPrefix,
	}

	if !block.RetryCount.IsNull() {
		if block.RetryInterval.IsNull() {
			return configSet,
				pathRoot.AtName("retry_count"),
				errors.New("retry_interval must be specified with retry_count")
		}
		configSet = append(configSet, setPrefix+"retry-count "+
			utils.ConvI64toa(block.RetryCount.ValueInt64())+
			" retry-interval "+utils.ConvI64toa(block.RetryInterval.ValueInt64()))
	} else if !block.RetryInterval.IsNull() {
		return configSet,
			pathRoot.AtName("retry_interval"),
			errors.New("retry_count must be specified with retry_interval")
	}
	if !block.TransferDelay.IsNull() {
		configSet = append(configSet, setPrefix+"transfer-delay "+
			utils.ConvI64toa(block.TransferDelay.ValueInt64()))
	}

	return configSet, path.Empty(), nil
}

func (block *eventoptionsPolicyBlockWithin) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	configSet := make([]string, 0)
	setPrefix += "within " + utils.ConvI64toa(block.TimeInterval.ValueInt64()) + " "

	for _, v := range block.Events {
		configSet = append(configSet, setPrefix+"events \""+v.ValueString()+"\"")
	}
	for _, v := range block.NotEvents {
		configSet = append(configSet, setPrefix+"not events \""+v.ValueString()+"\"")
	}
	if v := block.TriggerWhen.ValueString(); v != "" {
		if block.TriggerCount.IsNull() {
			return configSet,
				pathRoot.AtName("trigger_when"),
				errors.New("trigger_count must be specified with trigger_when")
		}
		configSet = append(configSet, setPrefix+"trigger "+v+" "+
			utils.ConvI64toa(block.TriggerCount.ValueInt64()))
	} else if !block.TriggerCount.IsNull() {
		return configSet,
			pathRoot.AtName("trigger_count"),
			errors.New("trigger_when must be specified with trigger_count")
	}

	return configSet, path.Empty(), nil
}

func (rscData *eventoptionsPolicyData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"event-options policy \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "events "):
				rscData.Events = append(rscData.Events, types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "then "):
				if rscData.Then == nil {
					rscData.Then = &eventoptionsPolicyBlockThen{}
				}
				if err := rscData.Then.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "attributes-match "):
				from := tfdata.FirstElementOfJunosLine(itemTrim)
				attributesMatch := eventoptionsPolicyBlockAttributesMatch{
					From: types.StringValue(strings.Trim(from, "\"")),
				}
				balt.CutPrefixInString(&itemTrim, from+" ")

				compare := tfdata.FirstElementOfJunosLine(itemTrim)
				attributesMatch.Compare = types.StringValue(compare)
				balt.CutPrefixInString(&itemTrim, compare+" ")

				attributesMatch.To = types.StringValue(strings.Trim(itemTrim, "\""))

				rscData.AttributesMatch = append(rscData.AttributesMatch, attributesMatch)
			case balt.CutPrefixInString(&itemTrim, "within "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var within eventoptionsPolicyBlockWithin
				withinTimeInterval, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
				if err != nil {
					return err
				}
				rscData.Within, within = tfdata.ExtractBlockWithTFTypesInt64(
					rscData.Within, "TimeInterval", withinTimeInterval.ValueInt64(),
				)
				within.TimeInterval = withinTimeInterval
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")

				if err := within.read(itemTrim); err != nil {
					return err
				}
				rscData.Within = append(rscData.Within, within)
			}
		}
	}

	return nil
}

func (block *eventoptionsPolicyBlockThen) read(itemTrim string) (err error) {
	switch {
	case itemTrim == "ignore":
		block.Ignore = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "priority-override facility "):
		block.PriorityOverrideFacility = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "priority-override severity "):
		block.PriorityOverrideSeverity = types.StringValue(itemTrim)
	case itemTrim == "raise-trap":
		block.RaiseTrap = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "change-configuration "):
		if block.ChangeConfiguration == nil {
			block.ChangeConfiguration = &eventoptionsPolicyBlockThenBlockChangeConfigurtion{}
		}
		if err := block.ChangeConfiguration.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "event-script "):
		filename := tfdata.FirstElementOfJunosLine(itemTrim)
		var eventScript eventoptionsPolicyBlockThenBlockEventScript
		block.EventScript, eventScript = tfdata.ExtractBlockWithTFTypesString(
			block.EventScript, "Filename", strings.Trim(filename, "\""),
		)
		eventScript.Filename = types.StringValue(strings.Trim(filename, "\""))
		balt.CutPrefixInString(&itemTrim, filename+" ")

		if err := eventScript.read(itemTrim); err != nil {
			return err
		}
		block.EventScript = append(block.EventScript, eventScript)
	case balt.CutPrefixInString(&itemTrim, "execute-commands "):
		if block.ExecuteCommands == nil {
			block.ExecuteCommands = &eventoptionsPolicyBlockThenBlockExecuteCommands{}
		}
		if err := block.ExecuteCommands.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "upload filename "):
		filename := tfdata.FirstElementOfJunosLine(itemTrim)
		var destination string
		if balt.CutPrefixInString(&itemTrim, filename+" destination ") {
			destination = tfdata.FirstElementOfJunosLine(itemTrim)
			balt.CutPrefixInString(&itemTrim, destination+" ")
		} else {
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "upload filename destination", itemTrim)
		}
		var upload eventoptionsPolicyBlockThenBlockUpload
		block.Upload, upload = tfdata.ExtractBlockWith2TFTypesString(
			block.Upload, "Filename", strings.Trim(filename, "\""), "Destination", strings.Trim(destination, "\""),
		)
		upload.Filename = types.StringValue(strings.Trim(filename, "\""))
		upload.Destination = types.StringValue(strings.Trim(destination, "\""))
		if err := upload.read(itemTrim); err != nil {
			return err
		}
		block.Upload = append(block.Upload, upload)
	}

	return nil
}

func (block *eventoptionsPolicyBlockThenBlockChangeConfigurtion) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "commands "):
		block.Commands = append(block.Commands, types.StringValue(strings.Trim(itemTrim, "\"")))
	case itemTrim == "commit-options check":
		block.CommitOptionsCheck = types.BoolValue(true)
	case itemTrim == "commit-options check synchronize":
		block.CommitOptionsCheck = types.BoolValue(true)
		block.CommitOptionsCheckSynchronize = types.BoolValue(true)
	case itemTrim == "commit-options force":
		block.CommitOptionsForce = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "commit-options log "):
		block.CommitOptionsLog = types.StringValue(strings.Trim(itemTrim, "\""))
	case itemTrim == "commit-options synchronize":
		block.CommitOptionsSynchronize = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "retry count "):
		block.RetryCount, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "retry interval "):
		block.RetryInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "user-name "):
		block.Username = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return nil
}

func (block *eventoptionsPolicyBlockThenBlockEventScript) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "arguments "):
		name := tfdata.FirstElementOfJunosLine(itemTrim)
		var arguments eventoptionsPolicyBlockThenBlockEventScriptBlockArguments
		block.Arguments, arguments = tfdata.ExtractBlockWithTFTypesString(
			block.Arguments, "Name", strings.Trim(name, "\""),
		)
		arguments.Name = types.StringValue(strings.Trim(name, "\""))
		balt.CutPrefixInString(&itemTrim, name+" ")
		arguments.Value = types.StringValue(strings.Trim(strings.TrimPrefix(itemTrim, name+" "), "\""))
		block.Arguments = append(block.Arguments, arguments)
	case balt.CutPrefixInString(&itemTrim, "destination "):
		name := tfdata.FirstElementOfJunosLine(itemTrim)
		if block.Destination == nil {
			block.Destination = &eventoptionsPolicyBlockThenBlockDestination{
				Name: types.StringValue(strings.Trim(name, "\"")),
			}
		}
		balt.CutPrefixInString(&itemTrim, name+" ")
		if err := block.Destination.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "output-filename "):
		block.OutputFilename = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "output-format "):
		block.OutputFormat = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "user-name "):
		block.Username = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return nil
}

func (block *eventoptionsPolicyBlockThenBlockDestination) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "retry-count retry-interval "):
		block.RetryInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "retry-count "):
		block.RetryCount, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "transfer-delay "):
		block.TransferDelay, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	}

	return nil
}

func (block *eventoptionsPolicyBlockThenBlockExecuteCommands) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "commands "):
		block.Commands = append(block.Commands, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "destination "):
		name := tfdata.FirstElementOfJunosLine(itemTrim)
		if block.Destination == nil {
			block.Destination = &eventoptionsPolicyBlockThenBlockDestination{
				Name: types.StringValue(strings.Trim(name, "\"")),
			}
		}
		balt.CutPrefixInString(&itemTrim, name+" ")
		if err := block.Destination.read(itemTrim); err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "output-filename "):
		block.OutputFilename = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "output-format "):
		block.OutputFormat = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "user-name "):
		block.Username = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return nil
}

func (block *eventoptionsPolicyBlockThenBlockUpload) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "retry-count retry-interval "):
		block.RetryInterval, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "retry-count "):
		block.RetryCount, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "transfer-delay "):
		block.TransferDelay, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "user-name "):
		block.Username = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return nil
}

func (block *eventoptionsPolicyBlockWithin) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "events "):
		block.Events = append(block.Events, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "not events "):
		block.NotEvents = append(block.NotEvents, types.StringValue(strings.Trim(itemTrim, "\"")))
	case balt.CutPrefixInString(&itemTrim, "trigger "):
		switch itemTrim {
		case "after", "on", "until":
			block.TriggerWhen = types.StringValue(itemTrim)
		default:
			block.TriggerCount, err = tfdata.ConvAtoi64Value(itemTrim)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (rscData *eventoptionsPolicyData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete event-options policy \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
