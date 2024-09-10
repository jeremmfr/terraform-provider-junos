package providerfwk

import (
	"context"
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
	_ resource.Resource                   = &securityNatStatic{}
	_ resource.ResourceWithConfigure      = &securityNatStatic{}
	_ resource.ResourceWithValidateConfig = &securityNatStatic{}
	_ resource.ResourceWithImportState    = &securityNatStatic{}
	_ resource.ResourceWithUpgradeState   = &securityNatStatic{}
)

type securityNatStatic struct {
	client *junos.Client
}

func newSecurityNatStaticResource() resource.Resource {
	return &securityNatStatic{}
}

func (rsc *securityNatStatic) typeName() string {
	return providerName + "_security_nat_static"
}

func (rsc *securityNatStatic) junosName() string {
	return "security nat static rule-set"
}

func (rsc *securityNatStatic) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityNatStatic) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityNatStatic) Configure(
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

func (rsc *securityNatStatic) Schema(
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
				Description: "Static nat rule-set name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"configure_rules_singly": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable management of rules.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of static nat rule-set.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"from": schema.SingleNestedBlock{
				Description: "Declare where is the traffic from.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Type of traffice source.",
						Validators: []validator.String{
							stringvalidator.OneOf("interface", "routing-instance", "zone"),
						},
					},
					"value": schema.SetAttribute{
						ElementType: types.StringType,
						Required:    true,
						Description: "Name of interface, routing-instance or zone for traffic source.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(
								stringvalidator.Any(
									tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
									tfvalidator.StringFormat(tfvalidator.DefaultFormat),
								),
								stringvalidator.LengthAtLeast(1),
							),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"rule": schema.ListNestedBlock{
				Description: "For each name of static nat rule to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Rule name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 31),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"destination_address": schema.StringAttribute{
							Optional:    true,
							Description: "CIDR destination address to match.",
							Validators: []validator.String{
								tfvalidator.StringCIDRNetwork(),
							},
						},
						"destination_address_name": schema.StringAttribute{
							Optional:    true,
							Description: "Destination address from address book to match.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
							},
						},
						"destination_port": schema.Int64Attribute{
							Optional:    true,
							Description: "Destination port or lower limit of port range to match.",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"destination_port_to": schema.Int64Attribute{
							Optional:    true,
							Description: "Port range upper limit to match.",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"source_address": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "CIDR source address to match.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									tfvalidator.StringCIDRNetwork(),
								),
							},
						},
						"source_address_name": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "Source address from address book to match.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
								),
							},
						},
						"source_port": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "Source port to match.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(regexp.MustCompile(
										`^\d+( to \d+)?$`),
										"must be use format `x` or `x to y`",
									),
								),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"then": schema.SingleNestedBlock{
							Description: "Declare `then` action.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:    true,
									Description: "Type of static nat.",
									Validators: []validator.String{
										stringvalidator.OneOf("inet", "prefix", "prefix-name"),
									},
								},
								"mapped_port": schema.Int64Attribute{
									Optional:    true,
									Description: "Port or lower limit of port range to mapped port.",
									Validators: []validator.Int64{
										int64validator.Between(1, 65535),
									},
								},
								"mapped_port_to": schema.Int64Attribute{
									Optional:    true,
									Description: "Port range upper limit to mapped port.",
									Validators: []validator.Int64{
										int64validator.Between(1, 65535),
									},
								},
								"prefix": schema.StringAttribute{
									Optional:    true,
									Description: "CIDR or address from address book to prefix static nat.",
									Validators: []validator.String{
										stringvalidator.Any(
											stringvalidator.All(
												stringvalidator.LengthBetween(1, 63),
												tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
											),
											tfvalidator.StringCIDRNetwork(),
										),
									},
								},
								"routing_instance": schema.StringAttribute{
									Optional:    true,
									Description: "Name of routing instance to switch instance with nat.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 63),
										tfvalidator.StringFormat(tfvalidator.DefaultFormat),
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

type securityNatStaticData struct {
	ID                   types.String                 `tfsdk:"id"`
	Name                 types.String                 `tfsdk:"name"`
	ConfigureRulesSingly types.Bool                   `tfsdk:"configure_rules_singly"`
	Description          types.String                 `tfsdk:"description"`
	From                 *securityNatStaticBlockFrom  `tfsdk:"from"`
	Rule                 []securityNatStaticBlockRule `tfsdk:"rule"`
}

type securityNatStaticConfig struct {
	ID                   types.String                      `tfsdk:"id"`
	Name                 types.String                      `tfsdk:"name"`
	ConfigureRulesSingly types.Bool                        `tfsdk:"configure_rules_singly"`
	Description          types.String                      `tfsdk:"description"`
	From                 *securityNatStaticBlockFromConfig `tfsdk:"from"`
	Rule                 types.List                        `tfsdk:"rule"`
}

type securityNatStaticBlockFrom struct {
	Type  types.String   `tfsdk:"type"`
	Value []types.String `tfsdk:"value"`
}

type securityNatStaticBlockFromConfig struct {
	Type  types.String `tfsdk:"type"`
	Value types.Set    `tfsdk:"value"`
}

type securityNatStaticBlockRule struct {
	Name                   types.String                    `tfsdk:"name"`
	DestinationAddress     types.String                    `tfsdk:"destination_address"`
	DestinationAddressName types.String                    `tfsdk:"destination_address_name"`
	DestinationPort        types.Int64                     `tfsdk:"destination_port"`
	DestiantionPortTo      types.Int64                     `tfsdk:"destination_port_to"`
	SourceAddress          []types.String                  `tfsdk:"source_address"`
	SourceAddressName      []types.String                  `tfsdk:"source_address_name"`
	SourcePort             []types.String                  `tfsdk:"source_port"`
	Then                   *securityNatStaticRuleBlockThen `tfsdk:"then"`
}

type securityNatStaticBlockRuleConfig struct {
	Name                   types.String                    `tfsdk:"name"`
	DestinationAddress     types.String                    `tfsdk:"destination_address"`
	DestinationAddressName types.String                    `tfsdk:"destination_address_name"`
	DestinationPort        types.Int64                     `tfsdk:"destination_port"`
	DestiantionPortTo      types.Int64                     `tfsdk:"destination_port_to"`
	SourceAddress          types.Set                       `tfsdk:"source_address"`
	SourceAddressName      types.Set                       `tfsdk:"source_address_name"`
	SourcePort             types.Set                       `tfsdk:"source_port"`
	Then                   *securityNatStaticRuleBlockThen `tfsdk:"then"`
}

func (rsc *securityNatStatic) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityNatStaticConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ConfigureRulesSingly.IsNull() &&
		config.Rule.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"one of configure_rules_singly or rule must be specified",
		)
	}
	if config.ConfigureRulesSingly.ValueBool() &&
		!config.Rule.IsNull() && !config.Rule.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("configure_rules_singly"),
			tfdiag.ConflictConfigErrSummary,
			"only one of configure_rules_singly or rule must be specified",
		)
	}

	if !config.Rule.IsNull() && !config.Rule.IsUnknown() {
		var rule []securityNatStaticBlockRuleConfig
		asDiags := config.Rule.ElementsAs(ctx, &rule, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		ruleName := make(map[string]struct{})
		for i, block := range rule {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := ruleName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("rule").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple rule blocks with the same name %q", name),
					)
				}
				ruleName[name] = struct{}{}
			}

			if block.DestinationAddress.IsNull() &&
				block.DestinationAddressName.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("rule").AtListIndex(i).AtName("destination_address"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("one of destination_address or destination_address_name must be specified"+
						" in rule block %q", block.Name.ValueString()),
				)
			}
			if !block.DestinationAddress.IsNull() && !block.DestinationAddress.IsUnknown() &&
				!block.DestinationAddressName.IsNull() && !block.DestinationAddressName.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("rule").AtListIndex(i).AtName("destination_address"),
					tfdiag.ConflictConfigErrSummary,
					fmt.Sprintf("destination_address and destination_address_name cannot be configured together"+
						" in rule block %q", block.Name.ValueString()),
				)
			}
			if !block.DestiantionPortTo.IsNull() &&
				block.DestinationPort.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("rule").AtListIndex(i).AtName("destination_port_to"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("cannot have destination_port_to without destination_port"+
						" in rule block %q", block.Name.ValueString()),
				)
			}
			if block.Then != nil {
				if !block.Then.Type.IsUnknown() {
					switch block.Then.Type.ValueString() {
					case junos.InetW:
						if !block.Then.Prefix.IsNull() && !block.Then.Prefix.IsUnknown() {
							resp.Diagnostics.AddAttributeError(
								path.Root("rule").AtListIndex(i).AtName("then").AtName("prefix"),
								tfdiag.ConflictConfigErrSummary,
								fmt.Sprintf("only routing_instance can be set when type = inet"+
									" in then block in rule block %q", block.Name.ValueString()),
							)
						}
						if !block.Then.MappedPort.IsNull() && !block.Then.MappedPort.IsUnknown() {
							resp.Diagnostics.AddAttributeError(
								path.Root("rule").AtListIndex(i).AtName("then").AtName("mapped_port"),
								tfdiag.ConflictConfigErrSummary,
								fmt.Sprintf("only routing_instance can be set when type = inet"+
									" in then block in rule block %q", block.Name.ValueString()),
							)
						}
						if !block.Then.MappedPortTo.IsNull() && !block.Then.MappedPortTo.IsUnknown() {
							resp.Diagnostics.AddAttributeError(
								path.Root("rule").AtListIndex(i).AtName("then").AtName("mapped_port_to"),
								tfdiag.ConflictConfigErrSummary,
								fmt.Sprintf("only routing_instance can be set when type = inet"+
									" in then block in rule block %q", block.Name.ValueString()),
							)
						}
					case "prefix":
						if block.Then.Prefix.IsNull() {
							resp.Diagnostics.AddAttributeError(
								path.Root("rule").AtListIndex(i).AtName("then").AtName("type"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("prefix must be specified when type = prefix"+
									" in then block in rule block %q", block.Name.ValueString()),
							)
						} else if !block.Then.Prefix.IsUnknown() {
							if err := tfvalidator.StringCIDRNetworkValidateAttribute(ctx, block.Then.Prefix); err != nil {
								resp.Diagnostics.AddAttributeError(
									path.Root("rule").AtListIndex(i).AtName("then").AtName("prefix"),
									"Invalid CIDR Network",
									err.Error(),
								)
							}
						}
					case "prefix-name":
						if block.Then.Prefix.IsNull() {
							resp.Diagnostics.AddAttributeError(
								path.Root("rule").AtListIndex(i).AtName("then").AtName("type"),
								tfdiag.MissingConfigErrSummary,
								fmt.Sprintf("prefix must be specified when type = prefix-name"+
									" in then block in rule block %q", block.Name.ValueString()),
							)
						}
					}
				}
				if !block.Then.MappedPortTo.IsNull() &&
					block.Then.MappedPort.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("rule").AtListIndex(i).AtName("then").AtName("mapped_port_to"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("cannot have mapped_port_to without mapped_port"+
							" in then block in rule block %q", block.Name.ValueString()),
					)
				}
			}
		}
	}
}

func (rsc *securityNatStatic) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityNatStaticData
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
			if !junSess.CheckCompatibilitySecurity() {
				resp.Diagnostics.AddError(
					tfdiag.CompatibilityErrSummary,
					rsc.junosName()+junSess.SystemInformation.NotCompatibleMsg(),
				)

				return false
			}
			natStaticExists, err := checkSecurityNatStaticExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if natStaticExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			natStaticExists, err := checkSecurityNatStaticExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !natStaticExists {
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

func (rsc *securityNatStatic) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityNatStaticData
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
		func() {
			data.ConfigureRulesSingly = state.ConfigureRulesSingly
			if data.ConfigureRulesSingly.ValueBool() {
				data.Rule = nil
			}
		},
		resp,
	)
}

func (rsc *securityNatStatic) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityNatStaticData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	configureRulesSingly := plan.ConfigureRulesSingly.ValueBool()
	if !plan.ConfigureRulesSingly.Equal(state.ConfigureRulesSingly) {
		if state.ConfigureRulesSingly.ValueBool() {
			configureRulesSingly = state.ConfigureRulesSingly.ValueBool()
			resp.Diagnostics.AddAttributeWarning(
				path.Root("configure_rules_singly"),
				"Disable configure_rules_singly on resource already created",
				"It's doesn't delete rule(s) already configured. "+
					"So refresh resource after apply to detect rule(s) that need to be deleted",
			)
		} else {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("configure_rules_singly"),
				"Enable configure_rules_singly on resource already created",
				"It's doesn't delete rule(s) already configured. "+
					"So import rule(s) in dedicated resource(s) to be able to manage them",
			)
		}
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		var delErr error
		if configureRulesSingly {
			delErr = state.delOptsWithoutRules(ctx, junSess)
		} else {
			delErr = state.del(ctx, junSess)
		}
		if delErr != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, delErr.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
			} else {
				resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigClearUnlockWarnSummary, junSess.ConfigClear())...)
	}()

	var delErr error
	if configureRulesSingly {
		delErr = state.delOptsWithoutRules(ctx, junSess)
	} else {
		delErr = state.del(ctx, junSess)
	}
	if delErr != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, delErr.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf(ctx, "update resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityNatStatic) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityNatStaticData
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

func (rsc *securityNatStatic) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	var data securityNatStaticData
	idList := strings.Split(req.ID, junos.IDSeparator)
	if err := data.read(ctx, idList[0], junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be <name> or <name>"+junos.IDSeparator+"no_rules)",
		)
	}
	if len(idList) > 1 && idList[1] == "no_rules" {
		data.ConfigureRulesSingly = types.BoolValue(true)
		data.Rule = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkSecurityNatStaticExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat static rule-set " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityNatStaticData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityNatStaticData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityNatStaticData) set(
	ctx context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security nat static rule-set " + rscData.Name.ValueString() + " "

	regexpSourcePort := regexp.MustCompile(`^\d+( to \d+)?$`)

	if rscData.From != nil {
		for _, value := range rscData.From.Value {
			configSet = append(configSet, setPrefix+"from "+rscData.From.Type.ValueString()+" "+value.ValueString())
		}
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+" description \""+v+"\"")
	}

	if !rscData.ConfigureRulesSingly.ValueBool() {
		ruleName := make(map[string]struct{})
		for i, block := range rscData.Rule {
			name := block.Name.ValueString()
			if _, ok := ruleName[name]; ok {
				return path.Root("rule").AtListIndex(i).AtName("name"),
					fmt.Errorf("multiple rule blocks with the same name %q", name)
			}
			ruleName[name] = struct{}{}

			setPrefixRule := setPrefix + "rule " + name + " "
			if block.DestinationAddress.IsNull() && block.DestinationAddressName.IsNull() {
				return path.Root("rule").AtListIndex(i).AtName("destination_address"),
					fmt.Errorf("destination_address or destination_address_name must be specified"+
						" in rule block %q", name)
			}
			if !block.DestinationAddress.IsNull() && !block.DestinationAddressName.IsNull() {
				return path.Root("rule").AtListIndex(i).AtName("destination_address"),
					fmt.Errorf("destination_address and destination_address_name cannot be configured together"+
						" in rule block %q", name)
			}
			if v := block.DestinationAddress.ValueString(); v != "" {
				configSet = append(configSet, setPrefixRule+"match destination-address "+v)
			}
			if v := block.DestinationAddressName.ValueString(); v != "" {
				configSet = append(configSet, setPrefixRule+"match destination-address-name \""+v+"\"")
			}
			if !block.DestinationPort.IsNull() {
				configSet = append(configSet, setPrefixRule+"match destination-port "+
					utils.ConvI64toa(block.DestinationPort.ValueInt64()))
				if !block.DestiantionPortTo.IsNull() {
					configSet = append(configSet, setPrefixRule+"match destination-port to "+
						utils.ConvI64toa(block.DestiantionPortTo.ValueInt64()))
				}
			} else if !block.DestiantionPortTo.IsNull() {
				return path.Root("rule").AtListIndex(i).AtName("destination_port_to"),
					fmt.Errorf("cannot have destination_port_to without destination_port"+
						" in rule block %q", name)
			}
			for _, v := range block.SourceAddress {
				configSet = append(configSet, setPrefixRule+"match source-address "+v.ValueString())
			}
			for _, v := range block.SourceAddressName {
				configSet = append(configSet, setPrefixRule+"match source-address-name \""+v.ValueString()+"\"")
			}
			for _, v := range block.SourcePort {
				if !regexpSourcePort.MatchString(v.ValueString()) {
					return path.Root("rule").AtListIndex(i).AtName("source_port"),
						fmt.Errorf("source_port need to have format `x` or `x to y`"+
							" in rule block %q", name)
				}
				configSet = append(configSet, setPrefixRule+"match source-port "+v.ValueString())
			}
			if block.Then != nil {
				setPrefixRuleThenStaticNat := setPrefixRule + "then static-nat "
				switch thenType := block.Then.Type.ValueString(); thenType {
				case junos.InetW:
					if !block.Then.Prefix.IsNull() {
						return path.Root("rule").AtListIndex(i).AtName("then").AtName("prefix"),
							fmt.Errorf("only routing_instance can be set when type = inet"+
								" in then block in rule block %q", name)
					}
					if !block.Then.MappedPort.IsNull() {
						return path.Root("rule").AtListIndex(i).AtName("then").AtName("mapped_port"),
							fmt.Errorf("only routing_instance can be set when type = inet"+
								" in then block in rule block %q", name)
					}
					if !block.Then.MappedPortTo.IsNull() {
						return path.Root("rule").AtListIndex(i).AtName("then").AtName("mapped_port_to"),
							fmt.Errorf("only routing_instance can be set when type = inet"+
								" in then block in rule block %q", name)
					}
					configSet = append(configSet, setPrefixRuleThenStaticNat+"inet")
					if v := block.Then.RoutingInstance.ValueString(); v != "" {
						configSet = append(configSet, setPrefixRuleThenStaticNat+"inet routing-instance "+v)
					}
				case "prefix", "prefix-name":
					if block.Then.Prefix.ValueString() == "" {
						return path.Root("rule").AtListIndex(i).AtName("then").AtName("type"),
							fmt.Errorf("type = %s and prefix without value"+
								" in then block in rule block %q", thenType, name)
					}
					switch thenType {
					case "prefix":
						setPrefixRuleThenStaticNat += "prefix "
						if err := tfvalidator.StringCIDRNetworkValidateAttribute(ctx, block.Then.Prefix); err != nil {
							return path.Root("rule").AtListIndex(i).AtName("then").AtName("prefix"), err
						}
					case "prefix-name":
						setPrefixRuleThenStaticNat += "prefix-name "
					}
					configSet = append(configSet, setPrefixRuleThenStaticNat+"\""+block.Then.Prefix.ValueString()+"\"")

					if !block.Then.MappedPort.IsNull() {
						configSet = append(configSet, setPrefixRuleThenStaticNat+"mapped-port "+
							utils.ConvI64toa(block.Then.MappedPort.ValueInt64()))
						if !block.Then.MappedPortTo.IsNull() {
							configSet = append(configSet, setPrefixRuleThenStaticNat+"mapped-port to "+
								utils.ConvI64toa(block.Then.MappedPortTo.ValueInt64()))
						}
					} else if !block.Then.MappedPortTo.IsNull() {
						return path.Root("rule").AtListIndex(i).AtName("then").AtName("mapped_port_to"),
							fmt.Errorf("cannot have mapped_port_to without mapped_port"+
								" in then block in rule block %q", name)
					}
					if v := block.Then.RoutingInstance.ValueString(); v != "" {
						configSet = append(configSet, setPrefixRuleThenStaticNat+"routing-instance "+v)
					}
				}
			}
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityNatStaticData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat static rule-set " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "from "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 2 { // <type> <value>
					return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "from", itemTrim)
				}
				if rscData.From == nil {
					rscData.From = &securityNatStaticBlockFrom{
						Type: types.StringValue(itemTrimFields[0]),
					}
				}
				rscData.From.Value = append(rscData.From.Value, types.StringValue(itemTrimFields[1]))
			case balt.CutPrefixInString(&itemTrim, "rule "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var rule securityNatStaticBlockRule
				rscData.Rule, rule = tfdata.ExtractBlockWithTFTypesString(
					rscData.Rule, "Name", itemTrimFields[0],
				)
				rule.Name = types.StringValue(itemTrimFields[0])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "match destination-address "):
					rule.DestinationAddress = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "match destination-address-name "):
					rule.DestinationAddressName = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "match destination-port to "):
					rule.DestiantionPortTo, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "match destination-port "):
					rule.DestinationPort, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "match source-address "):
					rule.SourceAddress = append(rule.SourceAddress, types.StringValue(itemTrim))
				case balt.CutPrefixInString(&itemTrim, "match source-address-name "):
					rule.SourceAddressName = append(rule.SourceAddressName, types.StringValue(strings.Trim(itemTrim, "\"")))
				case balt.CutPrefixInString(&itemTrim, "match source-port "):
					rule.SourcePort = append(rule.SourcePort, types.StringValue(itemTrim))
				case balt.CutPrefixInString(&itemTrim, "then static-nat "):
					if rule.Then == nil {
						rule.Then = &securityNatStaticRuleBlockThen{}
					}
					switch {
					case balt.CutPrefixInString(&itemTrim, "prefix"):
						rule.Then.Type = types.StringValue("prefix")
						if balt.CutPrefixInString(&itemTrim, "-name") {
							rule.Then.Type = types.StringValue("prefix-name")
						}
						balt.CutPrefixInString(&itemTrim, " ")
						switch {
						case balt.CutPrefixInString(&itemTrim, "routing-instance "):
							rule.Then.RoutingInstance = types.StringValue(itemTrim)
						case balt.CutPrefixInString(&itemTrim, "mapped-port to "):
							rule.Then.MappedPortTo, err = tfdata.ConvAtoi64Value(itemTrim)
							if err != nil {
								return err
							}
						case balt.CutPrefixInString(&itemTrim, "mapped-port "):
							rule.Then.MappedPort, err = tfdata.ConvAtoi64Value(itemTrim)
							if err != nil {
								return err
							}
						default:
							rule.Then.Prefix = types.StringValue(strings.Trim(itemTrim, "\""))
						}
					case balt.CutPrefixInString(&itemTrim, junos.InetW):
						rule.Then.Type = types.StringValue(junos.InetW)
						if balt.CutPrefixInString(&itemTrim, " routing-instance ") {
							rule.Then.RoutingInstance = types.StringValue(itemTrim)
						}
					}
				}
				rscData.Rule = append(rscData.Rule, rule)
			}
		}
	}

	return nil
}

func (rscData *securityNatStaticData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security nat static rule-set " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *securityNatStaticData) delOptsWithoutRules(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security nat static rule-set " + rscData.Name.ValueString() + " description",
		"delete security nat static rule-set " + rscData.Name.ValueString() + " from",
	}

	return junSess.ConfigSet(configSet)
}
