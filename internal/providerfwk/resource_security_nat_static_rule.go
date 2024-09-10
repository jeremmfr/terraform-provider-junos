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
	_ resource.Resource                   = &securityNatStaticRule{}
	_ resource.ResourceWithConfigure      = &securityNatStaticRule{}
	_ resource.ResourceWithValidateConfig = &securityNatStaticRule{}
	_ resource.ResourceWithImportState    = &securityNatStaticRule{}
	_ resource.ResourceWithUpgradeState   = &securityNatStaticRule{}
)

type securityNatStaticRule struct {
	client *junos.Client
}

func newSecurityNatStaticRuleResource() resource.Resource {
	return &securityNatStaticRule{}
}

func (rsc *securityNatStaticRule) typeName() string {
	return providerName + "_security_nat_static_rule"
}

func (rsc *securityNatStaticRule) junosName() string {
	return "security nat static rule in rule-set"
}

func (rsc *securityNatStaticRule) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityNatStaticRule) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityNatStaticRule) Configure(
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

func (rsc *securityNatStaticRule) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<rule_set>" + junos.IDSeparator + "<name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Static Rule name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"rule_set": schema.StringAttribute{
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
	}
}

type securityNatStaticRuleData struct {
	ID                     types.String                    `tfsdk:"id"`
	Name                   types.String                    `tfsdk:"name"`
	RuleSet                types.String                    `tfsdk:"rule_set"`
	DestinationAddress     types.String                    `tfsdk:"destination_address"`
	DestinationAddressName types.String                    `tfsdk:"destination_address_name"`
	DestinationPort        types.Int64                     `tfsdk:"destination_port"`
	DestiantionPortTo      types.Int64                     `tfsdk:"destination_port_to"`
	SourceAddress          []types.String                  `tfsdk:"source_address"`
	SourceAddressName      []types.String                  `tfsdk:"source_address_name"`
	SourcePort             []types.String                  `tfsdk:"source_port"`
	Then                   *securityNatStaticRuleBlockThen `tfsdk:"then"`
}

type securityNatStaticRuleConfig struct {
	ID                     types.String                    `tfsdk:"id"`
	Name                   types.String                    `tfsdk:"name"`
	RuleSet                types.String                    `tfsdk:"rule_set"`
	DestinationAddress     types.String                    `tfsdk:"destination_address"`
	DestinationAddressName types.String                    `tfsdk:"destination_address_name"`
	DestinationPort        types.Int64                     `tfsdk:"destination_port"`
	DestiantionPortTo      types.Int64                     `tfsdk:"destination_port_to"`
	SourceAddress          types.Set                       `tfsdk:"source_address"`
	SourceAddressName      types.Set                       `tfsdk:"source_address_name"`
	SourcePort             types.Set                       `tfsdk:"source_port"`
	Then                   *securityNatStaticRuleBlockThen `tfsdk:"then"`
}

type securityNatStaticRuleBlockThen struct {
	Type            types.String `tfsdk:"type"`
	MappedPort      types.Int64  `tfsdk:"mapped_port"`
	MappedPortTo    types.Int64  `tfsdk:"mapped_port_to"`
	Prefix          types.String `tfsdk:"prefix"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
}

func (rsc *securityNatStaticRule) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityNatStaticRuleConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.DestinationAddress.IsNull() &&
		config.DestinationAddressName.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"one of destination_address or destination_address_name must be specified",
		)
	}
	if !config.DestinationAddress.IsNull() && !config.DestinationAddress.IsUnknown() &&
		!config.DestinationAddressName.IsNull() && !config.DestinationAddressName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("destination_address"),
			tfdiag.ConflictConfigErrSummary,
			"only one of destination_address or destination_address_name must be specified",
		)
	}
	if !config.DestiantionPortTo.IsNull() &&
		config.DestinationPort.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("destination_port_to"),
			tfdiag.MissingConfigErrSummary,
			"cannot have destination_port_to without destination_port",
		)
	}
	if config.Then != nil {
		if !config.Then.Type.IsUnknown() {
			switch config.Then.Type.ValueString() {
			case junos.InetW:
				if !config.Then.Prefix.IsNull() && !config.Then.Prefix.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("prefix"),
						tfdiag.ConflictConfigErrSummary,
						"only routing_instance can be set when type = inet"+
							" in then block",
					)
				}
				if !config.Then.MappedPort.IsNull() && !config.Then.MappedPort.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("mapped_port"),
						tfdiag.ConflictConfigErrSummary,
						"only routing_instance can be set when type = inet"+
							" in then block",
					)
				}
				if !config.Then.MappedPortTo.IsNull() && !config.Then.MappedPortTo.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("mapped_port_to"),
						tfdiag.ConflictConfigErrSummary,
						"only routing_instance can be set when type = inet"+
							" in then block",
					)
				}
			case "prefix":
				if config.Then.Prefix.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("type"),
						tfdiag.MissingConfigErrSummary,
						"prefix must be specified when type = prefix"+
							" in then block",
					)
				} else if !config.Then.Prefix.IsUnknown() {
					if err := tfvalidator.StringCIDRNetworkValidateAttribute(ctx, config.Then.Prefix); err != nil {
						resp.Diagnostics.AddAttributeError(
							path.Root("then").AtName("prefix"),
							"Invalid CIDR Network",
							err.Error(),
						)
					}
				}
			case "prefix-name":
				if config.Then.Prefix.IsNull() {
					resp.Diagnostics.AddAttributeError(
						path.Root("then").AtName("type"),
						tfdiag.MissingConfigErrSummary,
						"prefix must be specified when type = prefix-name"+
							" in then block",
					)
				}
			}
		}
		if !config.Then.MappedPortTo.IsNull() &&
			config.Then.MappedPort.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("then").AtName("mapped_port_to"),
				tfdiag.MissingConfigErrSummary,
				"cannot have mapped_port_to without mapped_port"+
					" in then block",
			)
		}
	}
}

func (rsc *securityNatStaticRule) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityNatStaticRuleData
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
	if plan.RuleSet.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("rule_set"),
			"Empty rule-set",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "rule-set"),
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
			natStaticExists, err := checkSecurityNatStaticExists(fnCtx, plan.RuleSet.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !natStaticExists {
				resp.Diagnostics.AddError(
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("security nat static rule-set %q does not exists", plan.RuleSet.ValueString()),
				)

				return false
			}
			natStaticRuleExists, err := checkSecurityNatStaticRuleExists(
				fnCtx,
				plan.RuleSet.ValueString(),
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if natStaticRuleExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf(
						rsc.junosName()+" %q already exists in rule-set %q",
						plan.Name.ValueString(), plan.RuleSet.ValueString(),
					),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			natStaticRuleExists, err := checkSecurityNatStaticRuleExists(
				fnCtx,
				plan.RuleSet.ValueString(),
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !natStaticRuleExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					fmt.Sprintf(rsc.junosName()+" %q does not exists in rule-set %q after commit "+
						"=> check your config", plan.Name.ValueString(), plan.RuleSet.ValueString()),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *securityNatStaticRule) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityNatStaticRuleData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.RuleSet.ValueString(),
			state.Name.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *securityNatStaticRule) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityNatStaticRuleData
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

func (rsc *securityNatStaticRule) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityNatStaticRuleData
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

func (rsc *securityNatStaticRule) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityNatStaticRuleData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <rule_set>"+junos.IDSeparator+"<name>)",
	)
}

func checkSecurityNatStaticRuleExists(
	_ context.Context, ruleSet, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat static rule-set " + ruleSet + " rule " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityNatStaticRuleData) fillID() {
	rscData.ID = types.StringValue(rscData.RuleSet.ValueString() + junos.IDSeparator + rscData.Name.ValueString())
}

func (rscData *securityNatStaticRuleData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityNatStaticRuleData) set(
	ctx context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set security nat static " +
		"rule-set " + rscData.RuleSet.ValueString() +
		" rule " + rscData.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	regexpSourcePort := regexp.MustCompile(`^\d+( to \d+)?$`)

	if v := rscData.DestinationAddress.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"match destination-address "+v)
	}
	if v := rscData.DestinationAddressName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"match destination-address-name \""+v+"\"")
	}
	if !rscData.DestinationPort.IsNull() {
		configSet = append(configSet, setPrefix+"match destination-port "+
			utils.ConvI64toa(rscData.DestinationPort.ValueInt64()))
		if !rscData.DestiantionPortTo.IsNull() {
			configSet = append(configSet, setPrefix+"match destination-port to "+
				utils.ConvI64toa(rscData.DestiantionPortTo.ValueInt64()))
		}
	} else if !rscData.DestiantionPortTo.IsNull() {
		return path.Root("destination_port_to"),
			errors.New("cannot have destination_port_to without destination_port")
	}
	for _, v := range rscData.SourceAddress {
		configSet = append(configSet, setPrefix+"match source-address "+v.ValueString())
	}
	for _, v := range rscData.SourceAddressName {
		configSet = append(configSet, setPrefix+"match source-address-name \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.SourcePort {
		if !regexpSourcePort.MatchString(v.ValueString()) {
			return path.Root("source_port"),
				errors.New("source_port must be use format `x` or `x to y`")
		}
		configSet = append(configSet, setPrefix+"match source-port "+v.ValueString())
	}
	if rscData.Then != nil {
		setPrefixRuleThenStaticNat := setPrefix + "then static-nat "
		switch thenType := rscData.Then.Type.ValueString(); thenType {
		case junos.InetW:
			if !rscData.Then.Prefix.IsNull() {
				return path.Root("then").AtName("prefix"),
					errors.New("only routing_instance can be set when type = inet in then block")
			}
			if !rscData.Then.MappedPort.IsNull() {
				return path.Root("then").AtName("mapped_port"),
					errors.New("only routing_instance can be set when type = inet in then block")
			}
			if !rscData.Then.MappedPortTo.IsNull() {
				return path.Root("then").AtName("mapped_port_to"),
					errors.New("only routing_instance can be set when type = inet in then block")
			}
			configSet = append(configSet, setPrefixRuleThenStaticNat+"inet")
			if v := rscData.Then.RoutingInstance.ValueString(); v != "" {
				configSet = append(configSet, setPrefixRuleThenStaticNat+"inet routing-instance "+v)
			}
		case "prefix", "prefix-name":
			if rscData.Then.Prefix.ValueString() == "" {
				return path.Root("then").AtName("type"),
					fmt.Errorf("type = %s and prefix without value in then block", thenType)
			}
			switch thenType {
			case "prefix":
				setPrefixRuleThenStaticNat += "prefix "
				if err := tfvalidator.StringCIDRNetworkValidateAttribute(ctx, rscData.Then.Prefix); err != nil {
					return path.Root("then").AtName("prefix"), err
				}
			case "prefix-name":
				setPrefixRuleThenStaticNat += "prefix-name "
			}
			configSet = append(configSet, setPrefixRuleThenStaticNat+"\""+rscData.Then.Prefix.ValueString()+"\"")

			if !rscData.Then.MappedPort.IsNull() {
				configSet = append(configSet, setPrefixRuleThenStaticNat+"mapped-port "+
					utils.ConvI64toa(rscData.Then.MappedPort.ValueInt64()))
				if !rscData.Then.MappedPortTo.IsNull() {
					configSet = append(configSet, setPrefixRuleThenStaticNat+"mapped-port to "+
						utils.ConvI64toa(rscData.Then.MappedPortTo.ValueInt64()))
				}
			} else if !rscData.Then.MappedPortTo.IsNull() {
				return path.Root("then").AtName("mapped_port_to"),
					errors.New("cannot have mapped_port_to without mapped_port" +
						" in then block")
			}
			if v := rscData.Then.RoutingInstance.ValueString(); v != "" {
				configSet = append(configSet, setPrefixRuleThenStaticNat+"routing-instance "+v)
			}
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityNatStaticRuleData) read(
	_ context.Context, ruleSet, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat static rule-set " + ruleSet + " rule " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.RuleSet = types.StringValue(ruleSet)
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
			case balt.CutPrefixInString(&itemTrim, "match destination-address "):
				rscData.DestinationAddress = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "match destination-address-name "):
				rscData.DestinationAddressName = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "match destination-port to "):
				rscData.DestiantionPortTo, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "match destination-port "):
				rscData.DestinationPort, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "match source-address "):
				rscData.SourceAddress = append(rscData.SourceAddress,
					types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "match source-address-name "):
				rscData.SourceAddressName = append(rscData.SourceAddressName,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "match source-port "):
				rscData.SourcePort = append(rscData.SourcePort,
					types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "then static-nat "):
				if rscData.Then == nil {
					rscData.Then = &securityNatStaticRuleBlockThen{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "prefix"):
					rscData.Then.Type = types.StringValue("prefix")
					if balt.CutPrefixInString(&itemTrim, "-name") {
						rscData.Then.Type = types.StringValue("prefix-name")
					}
					balt.CutPrefixInString(&itemTrim, " ")
					switch {
					case balt.CutPrefixInString(&itemTrim, "routing-instance "):
						rscData.Then.RoutingInstance = types.StringValue(itemTrim)
					case balt.CutPrefixInString(&itemTrim, "mapped-port to "):
						rscData.Then.MappedPortTo, err = tfdata.ConvAtoi64Value(itemTrim)
						if err != nil {
							return err
						}
					case balt.CutPrefixInString(&itemTrim, "mapped-port "):
						rscData.Then.MappedPort, err = tfdata.ConvAtoi64Value(itemTrim)
						if err != nil {
							return err
						}
					default:
						rscData.Then.Prefix = types.StringValue(strings.Trim(itemTrim, "\""))
					}
				case balt.CutPrefixInString(&itemTrim, junos.InetW):
					rscData.Then.Type = types.StringValue(junos.InetW)
					if balt.CutPrefixInString(&itemTrim, " routing-instance ") {
						rscData.Then.RoutingInstance = types.StringValue(itemTrim)
					}
				}
			}
		}
	}

	return nil
}

func (rscData *securityNatStaticRuleData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security nat static rule-set " + rscData.RuleSet.ValueString() + " rule " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
