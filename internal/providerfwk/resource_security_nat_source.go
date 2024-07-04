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
	_ resource.Resource                   = &securityNatSource{}
	_ resource.ResourceWithConfigure      = &securityNatSource{}
	_ resource.ResourceWithValidateConfig = &securityNatSource{}
	_ resource.ResourceWithImportState    = &securityNatSource{}
	_ resource.ResourceWithUpgradeState   = &securityNatSource{}
)

type securityNatSource struct {
	client *junos.Client
}

func newSecurityNatSourceResource() resource.Resource {
	return &securityNatSource{}
}

func (rsc *securityNatSource) typeName() string {
	return providerName + "_security_nat_source"
}

func (rsc *securityNatSource) junosName() string {
	return "security nat source rule-set"
}

func (rsc *securityNatSource) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityNatSource) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityNatSource) Configure(
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

func (rsc *securityNatSource) Schema(
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
				Description: "Source nat rule-set name.",
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
				Description: "Text description of rule set.",
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
			"to": schema.SingleNestedBlock{
				Description: "Declare where is the traffic to.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Type of traffic destination.",
						Validators: []validator.String{
							stringvalidator.OneOf("interface", "routing-instance", "zone"),
						},
					},
					"value": schema.SetAttribute{
						ElementType: types.StringType,
						Required:    true,
						Description: "Name of interface, routing-instance or zone for traffic destination.",
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
				Description: "For each name of source nat rule to declare.",
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
					},
					Blocks: map[string]schema.Block{
						"match": schema.SingleNestedBlock{
							Description: "Specify source nat rule match criteria.",
							Attributes: map[string]schema.Attribute{
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
								"destination_address": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "CIDR destination address to match.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											tfvalidator.StringCIDRNetwork(),
										),
									},
								},
								"destination_address_name": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "Destination address from address book to match.",
									Validators: []validator.Set{
										setvalidator.SizeAtLeast(1),
										setvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 63),
											tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
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
							PlanModifiers: []planmodifier.Object{
								tfplanmodifier.BlockRemoveNull(),
							},
						},
						"then": schema.SingleNestedBlock{
							Description: "Declare `then` action.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:    true,
									Description: "Type of source nat.",
									Validators: []validator.String{
										stringvalidator.OneOf("interface", "off", "pool"),
									},
								},
								"pool": schema.StringAttribute{
									Optional:    true,
									Description: "Name of source nat pool when type is pool.",
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

type securityNatSourceData struct {
	ID          types.String                  `tfsdk:"id"`
	Name        types.String                  `tfsdk:"name"`
	Description types.String                  `tfsdk:"description"`
	From        *securityNatSourceBlockFromTo `tfsdk:"from"`
	To          *securityNatSourceBlockFromTo `tfsdk:"to"`
	Rule        []securityNatSourceBlockRule  `tfsdk:"rule"`
}

type securityNatSourceConfig struct {
	ID          types.String                        `tfsdk:"id"`
	Name        types.String                        `tfsdk:"name"`
	Description types.String                        `tfsdk:"description"`
	From        *securityNatSourceBlockFromToConfig `tfsdk:"from"`
	To          *securityNatSourceBlockFromToConfig `tfsdk:"to"`
	Rule        types.List                          `tfsdk:"rule"`
}

type securityNatSourceBlockFromTo struct {
	Type  types.String   `tfsdk:"type"`
	Value []types.String `tfsdk:"value"`
}

type securityNatSourceBlockFromToConfig struct {
	Type  types.String `tfsdk:"type"`
	Value types.Set    `tfsdk:"value"`
}

type securityNatSourceBlockRule struct {
	Name  types.String                          `tfsdk:"name"`
	Match *securityNatSourceBlockRuleBlockMatch `tfsdk:"match"`
	Then  *securityNatSourceBlockRuleBlockThen  `tfsdk:"then"`
}

type securityNatSourceBlockRuleConfig struct {
	Name  types.String                                `tfsdk:"name"`
	Match *securityNatSourceBlockRuleBlockMatchConfig `tfsdk:"match"`
	Then  *securityNatSourceBlockRuleBlockThen        `tfsdk:"then"`
}

type securityNatSourceBlockRuleBlockMatch struct {
	Application            []types.String `tfsdk:"application"`
	DestinationAddress     []types.String `tfsdk:"destination_address"`
	DestinationAddressName []types.String `tfsdk:"destination_address_name"`
	DestinationPort        []types.String `tfsdk:"destination_port"`
	Protocol               []types.String `tfsdk:"protocol"`
	SourceAddress          []types.String `tfsdk:"source_address"`
	SourceAddressName      []types.String `tfsdk:"source_address_name"`
	SourcePort             []types.String `tfsdk:"source_port"`
}

func (block *securityNatSourceBlockRuleBlockMatch) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityNatSourceBlockRuleBlockMatchConfig struct {
	Application            types.Set `tfsdk:"application"`
	DestinationAddress     types.Set `tfsdk:"destination_address"`
	DestinationAddressName types.Set `tfsdk:"destination_address_name"`
	DestinationPort        types.Set `tfsdk:"destination_port"`
	Protocol               types.Set `tfsdk:"protocol"`
	SourceAddress          types.Set `tfsdk:"source_address"`
	SourceAddressName      types.Set `tfsdk:"source_address_name"`
	SourcePort             types.Set `tfsdk:"source_port"`
}

func (block *securityNatSourceBlockRuleBlockMatchConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityNatSourceBlockRuleBlockThen struct {
	Type types.String `tfsdk:"type"`
	Pool types.String `tfsdk:"pool"`
}

func (rsc *securityNatSource) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityNatSourceConfig
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
		var rule []securityNatSourceBlockRuleConfig
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

			if block.Match == nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("rule").AtListIndex(i).AtName("match"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("match block must be specified"+
						" in rule block %q", block.Name.ValueString()),
				)
			} else {
				if block.Match.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("rule").AtListIndex(i).AtName("match"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("match block is empty"+
							" in rule block %q", block.Name.ValueString()),
					)
				} else if block.Match.DestinationAddress.IsNull() &&
					block.Match.DestinationAddressName.IsNull() &&
					block.Match.SourceAddress.IsNull() &&
					block.Match.SourceAddressName.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("rule").AtListIndex(i).AtName("match"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("one of destination_address, destination_address_name, "+
							"source_address or source_address_name must be specified"+
							" in rule block %q", block.Name.ValueString()),
					)
				}
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

func (rsc *securityNatSource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityNatSourceData
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
			natSourceExists, err := checkSecurityNatSourceExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if natSourceExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			natSourceExists, err := checkSecurityNatSourceExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !natSourceExists {
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

func (rsc *securityNatSource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityNatSourceData
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

func (rsc *securityNatSource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityNatSourceData
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

func (rsc *securityNatSource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityNatSourceData
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

func (rsc *securityNatSource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityNatSourceData

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

func checkSecurityNatSourceExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat source rule-set " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityNatSourceData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityNatSourceData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityNatSourceData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security nat source rule-set " + rscData.Name.ValueString() + " "

	regexpPort := regexp.MustCompile(`^\d+( to \d+)?$`)

	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if rscData.From != nil {
		for _, v := range rscData.From.Value {
			configSet = append(configSet, setPrefix+"from "+rscData.From.Type.ValueString()+" "+v.ValueString())
		}
	}
	if rscData.To != nil {
		for _, v := range rscData.To.Value {
			configSet = append(configSet, setPrefix+"to "+rscData.To.Type.ValueString()+" "+v.ValueString())
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

		setPrefixRule := setPrefix + " rule " + name + " "
		if block.Match != nil {
			if block.Match.isEmpty() {
				return path.Root("rule").AtListIndex(i).AtName("match"),
					fmt.Errorf("match block is empty in rule block %q", block.Name.ValueString())
			}
			if len(block.Match.DestinationAddress) == 0 &&
				len(block.Match.DestinationAddressName) == 0 &&
				len(block.Match.SourceAddress) == 0 &&
				len(block.Match.SourceAddressName) == 0 {
				return path.Root("rule").AtListIndex(i).AtName("match"),
					fmt.Errorf("one of destination_address, destination_address_name,"+
						" source_address or source_address_name arguments must be specified"+
						" in rule block %q", block.Name.ValueString())
			}
			for _, v := range block.Match.Application {
				configSet = append(configSet, setPrefixRule+"match application \""+v.ValueString()+"\"")
			}
			for _, v := range block.Match.DestinationAddress {
				configSet = append(configSet, setPrefixRule+"match destination-address "+v.ValueString())
			}
			for _, v := range block.Match.DestinationAddressName {
				configSet = append(configSet, setPrefixRule+"match destination-address-name \""+v.ValueString()+"\"")
			}
			for _, v := range block.Match.DestinationPort {
				if !regexpPort.MatchString(v.ValueString()) {
					return path.Root("rule").AtListIndex(i).AtName("destination_port"),
						fmt.Errorf("destination_port must be use format `x` or `x to y`"+
							" in rule block %q", name)
				}
				configSet = append(configSet, setPrefixRule+"match destination-port "+v.ValueString())
			}
			for _, v := range block.Match.Protocol {
				configSet = append(configSet, setPrefixRule+"match protocol "+v.ValueString())
			}
			for _, v := range block.Match.SourceAddress {
				configSet = append(configSet, setPrefixRule+"match source-address "+v.ValueString())
			}
			for _, v := range block.Match.SourceAddressName {
				configSet = append(configSet, setPrefixRule+"match source-address-name \""+v.ValueString()+"\"")
			}
			for _, v := range block.Match.SourcePort {
				if !regexpPort.MatchString(v.ValueString()) {
					return path.Root("rule").AtListIndex(i).AtName("destination_port"),
						fmt.Errorf("source_port must be use format `x` or `x to y`"+
							" in rule block %q", name)
				}
				configSet = append(configSet, setPrefixRule+"match source-port "+v.ValueString())
			}
		} else {
			return path.Root("rule").AtListIndex(i).AtName("match"),
				fmt.Errorf("match block must be specified"+
					" in rule block %q", block.Name.ValueString())
		}
		if block.Then != nil {
			switch thenType := block.Then.Type.ValueString(); thenType {
			case "interface", "off":
				configSet = append(configSet, setPrefixRule+"then source-nat "+thenType)
			case "pool":
				if block.Then.Pool.ValueString() == "" {
					return path.Root("rule").AtListIndex(i).AtName("then").AtName("pool"),
						fmt.Errorf("missing pool value to source-nat pool"+
							" in rule block %q", name)
				}
				configSet = append(configSet, setPrefixRule+"then source-nat "+thenType+" "+block.Then.Pool.ValueString())
			}
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityNatSourceData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat source rule-set " + name + junos.PipeDisplaySetRelative)
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
					rscData.From = &securityNatSourceBlockFromTo{
						Type: types.StringValue(itemTrimFields[0]),
					}
				}
				rscData.From.Value = append(rscData.From.Value, types.StringValue(itemTrimFields[1]))
			case balt.CutPrefixInString(&itemTrim, "to "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 2 { // <type> <value>
					return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "to", itemTrim)
				}
				if rscData.To == nil {
					rscData.To = &securityNatSourceBlockFromTo{
						Type: types.StringValue(itemTrimFields[0]),
					}
				}
				rscData.To.Value = append(rscData.To.Value, types.StringValue(itemTrimFields[1]))
			case balt.CutPrefixInString(&itemTrim, "rule "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var rule securityNatSourceBlockRule
				rscData.Rule, rule = tfdata.ExtractBlockWithTFTypesString(
					rscData.Rule, "Name", itemTrimFields[0],
				)
				rule.Name = types.StringValue(itemTrimFields[0])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "match "):
					if rule.Match == nil {
						rule.Match = &securityNatSourceBlockRuleBlockMatch{}
					}
					switch {
					case balt.CutPrefixInString(&itemTrim, "application "):
						rule.Match.Application = append(rule.Match.Application,
							types.StringValue(strings.Trim(itemTrim, "\"")),
						)
					case balt.CutPrefixInString(&itemTrim, "destination-address "):
						rule.Match.DestinationAddress = append(rule.Match.DestinationAddress,
							types.StringValue(itemTrim),
						)
					case balt.CutPrefixInString(&itemTrim, "destination-address-name "):
						rule.Match.DestinationAddressName = append(rule.Match.DestinationAddressName,
							types.StringValue(strings.Trim(itemTrim, "\"")),
						)
					case balt.CutPrefixInString(&itemTrim, "destination-port "):
						rule.Match.DestinationPort = append(rule.Match.DestinationPort,
							types.StringValue(itemTrim),
						)
					case balt.CutPrefixInString(&itemTrim, "protocol "):
						rule.Match.Protocol = append(rule.Match.Protocol,
							types.StringValue(itemTrim),
						)
					case balt.CutPrefixInString(&itemTrim, "source-address "):
						rule.Match.SourceAddress = append(rule.Match.SourceAddress,
							types.StringValue(itemTrim),
						)
					case balt.CutPrefixInString(&itemTrim, "source-address-name "):
						rule.Match.SourceAddressName = append(rule.Match.SourceAddressName,
							types.StringValue(strings.Trim(itemTrim, "\"")),
						)
					case balt.CutPrefixInString(&itemTrim, "source-port "):
						rule.Match.SourcePort = append(rule.Match.SourcePort,
							types.StringValue(itemTrim),
						)
					}
				case balt.CutPrefixInString(&itemTrim, "then source-nat "):
					if rule.Then == nil {
						rule.Then = &securityNatSourceBlockRuleBlockThen{}
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

func (rscData *securityNatSourceData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security nat source rule-set " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
