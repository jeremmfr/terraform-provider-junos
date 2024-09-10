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
	_ resource.Resource                   = &securityNatDestination{}
	_ resource.ResourceWithConfigure      = &securityNatDestination{}
	_ resource.ResourceWithValidateConfig = &securityNatDestination{}
	_ resource.ResourceWithImportState    = &securityNatDestination{}
	_ resource.ResourceWithUpgradeState   = &securityNatDestination{}
)

type securityNatDestination struct {
	client *junos.Client
}

func newSecurityNatDestinationResource() resource.Resource {
	return &securityNatDestination{}
}

func (rsc *securityNatDestination) typeName() string {
	return providerName + "_security_nat_destination"
}

func (rsc *securityNatDestination) junosName() string {
	return "security nat destination rule-set"
}

func (rsc *securityNatDestination) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityNatDestination) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityNatDestination) Configure(
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

func (rsc *securityNatDestination) Schema(
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
				Description: "Destination nat rule-set name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of destination nat rule-set.",
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
								stringvalidator.LengthAtLeast(1),
								stringvalidator.Any(
									tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
									tfvalidator.StringFormat(tfvalidator.DefaultFormat),
								),
							),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"rule": schema.ListNestedBlock{
				Description: "For each name of destination nat rule to declare.",
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
						"application": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "Application or application-set name to match.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringFormat(tfvalidator.DefaultFormat),
								),
							},
						},
						"destination_port": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "Destination port to match.",
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
						"protocol": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "IP Protocol to match.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringFormat(tfvalidator.DefaultFormat),
								),
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
					},
					Blocks: map[string]schema.Block{
						"then": schema.SingleNestedBlock{
							Description: "Declare `then` action.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:    true,
									Description: "Type of destination nat.",
									Validators: []validator.String{
										stringvalidator.OneOf("off", "pool"),
									},
								},
								"pool": schema.StringAttribute{
									Optional:    true,
									Description: "Name of destination nat pool when type is pool.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 31),
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

type securityNatDestinationData struct {
	ID          types.String                      `tfsdk:"id"`
	Name        types.String                      `tfsdk:"name"`
	Description types.String                      `tfsdk:"description"`
	From        *securityNatDestinationBlockFrom  `tfsdk:"from"`
	Rule        []securityNatDestinationBlockRule `tfsdk:"rule"`
}

type securityNatDestinationConfig struct {
	ID          types.String                           `tfsdk:"id"`
	Name        types.String                           `tfsdk:"name"`
	Description types.String                           `tfsdk:"description"`
	From        *securityNatDestinationBlockFromConfig `tfsdk:"from"`
	Rule        types.List                             `tfsdk:"rule"`
}

type securityNatDestinationBlockFrom struct {
	Type  types.String   `tfsdk:"type"`
	Value []types.String `tfsdk:"value"`
}

type securityNatDestinationBlockFromConfig struct {
	Type  types.String `tfsdk:"type"`
	Value types.Set    `tfsdk:"value"`
}

type securityNatDestinationBlockRule struct {
	Name                   types.String                              `tfsdk:"name"`
	DestinationAddress     types.String                              `tfsdk:"destination_address"`
	DestinationAddressName types.String                              `tfsdk:"destination_address_name"`
	Application            []types.String                            `tfsdk:"application"`
	DestinationPort        []types.String                            `tfsdk:"destination_port"`
	Protocol               []types.String                            `tfsdk:"protocol"`
	SourceAddress          []types.String                            `tfsdk:"source_address"`
	SourceAddressName      []types.String                            `tfsdk:"source_address_name"`
	Then                   *securityNatDestinationBlockRuleBlockThen `tfsdk:"then"`
}

type securityNatDestinationBlockRuleConfig struct {
	Name                   types.String                              `tfsdk:"name"`
	DestinationAddress     types.String                              `tfsdk:"destination_address"`
	DestinationAddressName types.String                              `tfsdk:"destination_address_name"`
	Application            types.Set                                 `tfsdk:"application"`
	DestinationPort        types.Set                                 `tfsdk:"destination_port"`
	Protocol               types.Set                                 `tfsdk:"protocol"`
	SourceAddress          types.Set                                 `tfsdk:"source_address"`
	SourceAddressName      types.Set                                 `tfsdk:"source_address_name"`
	Then                   *securityNatDestinationBlockRuleBlockThen `tfsdk:"then"`
}

type securityNatDestinationBlockRuleBlockThen struct {
	Type types.String `tfsdk:"type"`
	Pool types.String `tfsdk:"pool"`
}

func (rsc *securityNatDestination) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityNatDestinationConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Rule.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("rule"),
			tfdiag.MissingConfigErrSummary,
			"at least one rule block must be specified",
		)
	} else if !config.Rule.IsUnknown() {
		var rule []securityNatDestinationBlockRuleConfig
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
			if block.Then != nil {
				if !block.Then.Type.IsUnknown() && block.Then.Type.ValueString() == "pool" &&
					block.Then.Pool.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("rule").AtListIndex(i).AtName("then").AtName("type"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("pool must be specified when type is set to %q"+
							" in then block in rule block %q",
							block.Then.Type.ValueString(), block.Name.ValueString()),
					)
				}
			}
		}
	}
}

func (rsc *securityNatDestination) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityNatDestinationData
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
			natDestinationExists, err := checkSecurityNatDestinationExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if natDestinationExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			natDestinationExists, err := checkSecurityNatDestinationExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !natDestinationExists {
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

func (rsc *securityNatDestination) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityNatDestinationData
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

func (rsc *securityNatDestination) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityNatDestinationData
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

func (rsc *securityNatDestination) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityNatDestinationData
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

func (rsc *securityNatDestination) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityNatDestinationData

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

func checkSecurityNatDestinationExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat destination rule-set " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityNatDestinationData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityNatDestinationData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityNatDestinationData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security nat destination rule-set " + rscData.Name.ValueString() + " "

	regexpDestPort := regexp.MustCompile(`^\d+( to \d+)?$`)

	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if rscData.From != nil {
		for _, v := range rscData.From.Value {
			configSet = append(configSet, setPrefix+"from "+rscData.From.Type.ValueString()+" "+v.ValueString())
		}
	}
	ruleName := make(map[string]struct{})
	for i, block := range rscData.Rule {
		name := block.Name.ValueString()
		if _, ok := ruleName[name]; ok {
			return path.Root("rule").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple rule blocks with the same name %q", name)
		}
		ruleName[name] = struct{}{}

		if block.DestinationAddress.IsNull() &&
			block.DestinationAddressName.IsNull() {
			return path.Root("rule").AtListIndex(i).AtName("destination_address"),
				fmt.Errorf("one of destination_address or destination_address_name must be specified"+
					" in rule block %q", name)
		}
		if !block.DestinationAddress.IsNull() &&
			!block.DestinationAddressName.IsNull() {
			return path.Root("rule").AtListIndex(i).AtName("destination_address"),
				fmt.Errorf("destination_address and destination_address_name cannot be configured together"+
					" in rule block %q", name)
		}
		setPrefixRule := setPrefix + "rule " + name + " "
		if v := block.DestinationAddress.ValueString(); v != "" {
			configSet = append(configSet, setPrefixRule+"match destination-address "+v)
		}
		if v := block.DestinationAddressName.ValueString(); v != "" {
			configSet = append(configSet, setPrefixRule+"match destination-address-name \""+v+"\"")
		}
		for _, v := range block.Application {
			configSet = append(configSet, setPrefixRule+"match application \""+v.ValueString()+"\"")
		}
		for _, v := range block.DestinationPort {
			if !regexpDestPort.MatchString(v.ValueString()) {
				return path.Root("rule").AtListIndex(i).AtName("destination_port"),
					fmt.Errorf("destination_port must be use format `x` or `x to y`"+
						" in rule block %q", name)
			}
			configSet = append(configSet, setPrefixRule+"match destination-port "+v.ValueString())
		}
		for _, v := range block.Protocol {
			configSet = append(configSet, setPrefixRule+"match protocol "+v.ValueString())
		}
		for _, v := range block.SourceAddress {
			configSet = append(configSet, setPrefixRule+"match source-address "+v.ValueString())
		}
		for _, v := range block.SourceAddressName {
			configSet = append(configSet, setPrefixRule+"match source-address-name \""+v.ValueString()+"\"")
		}
		if block.Then != nil {
			switch block.Then.Type.ValueString() {
			case "off":
				configSet = append(configSet, setPrefixRule+"then destination-nat off")
			case "pool":
				pool := block.Then.Pool.ValueString()
				if pool == "" {
					return path.Root("rule").AtListIndex(i).AtName("then").AtName("pool"),
						fmt.Errorf("missing pool value to destination-nat pool"+
							" in rule block %q", name)
				}
				configSet = append(configSet, setPrefixRule+"then destination-nat pool "+pool)
			}
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityNatDestinationData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat destination rule-set " + name + junos.PipeDisplaySetRelative)
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
					rscData.From = &securityNatDestinationBlockFrom{
						Type: types.StringValue(itemTrimFields[0]),
					}
				}
				rscData.From.Value = append(rscData.From.Value, types.StringValue(itemTrimFields[1]))
			case balt.CutPrefixInString(&itemTrim, "rule "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var rule securityNatDestinationBlockRule
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
				case balt.CutPrefixInString(&itemTrim, "match application "):
					rule.Application = append(rule.Application,
						types.StringValue(strings.Trim(itemTrim, "\"")),
					)
				case balt.CutPrefixInString(&itemTrim, "match destination-port "):
					rule.DestinationPort = append(rule.DestinationPort,
						types.StringValue(itemTrim),
					)
				case balt.CutPrefixInString(&itemTrim, "match protocol "):
					rule.Protocol = append(rule.Protocol,
						types.StringValue(itemTrim),
					)
				case balt.CutPrefixInString(&itemTrim, "match source-address "):
					rule.SourceAddress = append(rule.SourceAddress,
						types.StringValue(itemTrim),
					)
				case balt.CutPrefixInString(&itemTrim, "match source-address-name "):
					rule.SourceAddressName = append(rule.SourceAddressName,
						types.StringValue(strings.Trim(itemTrim, "\"")),
					)
				case balt.CutPrefixInString(&itemTrim, "then destination-nat "):
					if rule.Then == nil {
						rule.Then = &securityNatDestinationBlockRuleBlockThen{}
					}
					if balt.CutPrefixInString(&itemTrim, "pool ") {
						rule.Then.Type = types.StringValue("pool")
						rule.Then.Pool = types.StringValue(itemTrim)
					} else {
						rule.Then.Type = types.StringValue(itemTrim)
					}
				}
				rscData.Rule = append(rscData.Rule, rule)
			}
		}
	}

	return nil
}

func (rscData *securityNatDestinationData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security nat destination rule-set " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
