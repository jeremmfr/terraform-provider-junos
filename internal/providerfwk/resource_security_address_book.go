package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	_ resource.Resource                   = &securityAddressBook{}
	_ resource.ResourceWithConfigure      = &securityAddressBook{}
	_ resource.ResourceWithValidateConfig = &securityAddressBook{}
	_ resource.ResourceWithImportState    = &securityAddressBook{}
)

type securityAddressBook struct {
	client *junos.Client
}

func newSecurityAddressBookResource() resource.Resource {
	return &securityAddressBook{}
}

func (rsc *securityAddressBook) typeName() string {
	return providerName + "_security_address_book"
}

func (rsc *securityAddressBook) junosName() string {
	return "security address book"
}

func (rsc *securityAddressBook) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityAddressBook) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityAddressBook) Configure(
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

func (rsc *securityAddressBook) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
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
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("global"),
				Description: "The name of address book.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the address book.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"attach_zone": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of zones to attach address book to.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 63),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"network_address": schema.SetNestedBlock{
				Description: "For each name of network address.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of network address.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
							},
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "CIDR value of network address (`192.0.0.0/24`).",
							Validators: []validator.String{
								tfvalidator.StringCIDRNetwork(),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Description of network address.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 900),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
				},
			},
			"dns_name": schema.SetNestedBlock{
				Description: "For each name of dns name address.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of dns name address.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
							},
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "DNS name string value.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 253),
								tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Description of dns name address.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 900),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"ipv4_only": schema.BoolAttribute{
							Optional:    true,
							Description: "IPv4 dns address.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"ipv6_only": schema.BoolAttribute{
							Optional:    true,
							Description: "IPv6 dns address.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
					},
				},
			},
			"range_address": schema.SetNestedBlock{
				Description: "For each name of range address.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of range address.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
							},
						},
						"from": schema.StringAttribute{
							Required:    true,
							Description: "IP address of start of range.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress(),
							},
						},
						"to": schema.StringAttribute{
							Required:    true,
							Description: "IP address of end of range.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress(),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Description of range address.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 900),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
				},
			},
			"wildcard_address": schema.SetNestedBlock{
				Description: "For each name of wildcard address.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of wildcard address.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
							},
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "Network and mask of wildcard address (`192.0.0.0/255.255.0.255`).",
							Validators: []validator.String{
								tfvalidator.StringWildcardNetwork(),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Description of wildcard address.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 900),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
				},
			},
			"address_set": schema.SetNestedBlock{
				Description: "For each name of address-set to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of address-set.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
							},
						},
						"address": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of address names.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
								),
							},
						},
						"address_set": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of address-set names.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthBetween(1, 63),
									tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
								),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Description of address-set.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 900),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
				},
			},
		},
	}
}

type securityAddressBookData struct {
	ID              types.String                              `tfsdk:"id"`
	Name            types.String                              `tfsdk:"name"`
	Description     types.String                              `tfsdk:"description"`
	AttachZone      []types.String                            `tfsdk:"attach_zone"`
	NetworkAddress  []securityAddressBookBlockNetworkAddress  `tfsdk:"network_address"`
	DNSName         []securityAddressBookBlockDNSName         `tfsdk:"dns_name"`
	RangeAddress    []securityAddressBookBlockRangeAddress    `tfsdk:"range_address"`
	WildcardAddress []securityAddressBookBlockWildcardAddress `tfsdk:"wildcard_address"`
	AddressSet      []securityAddressBookBlockAddressSet      `tfsdk:"address_set"`
}

type securityAddressBookConfig struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	AttachZone      types.List   `tfsdk:"attach_zone"`
	NetworkAddress  types.Set    `tfsdk:"network_address"`
	DNSName         types.Set    `tfsdk:"dns_name"`
	RangeAddress    types.Set    `tfsdk:"range_address"`
	WildcardAddress types.Set    `tfsdk:"wildcard_address"`
	AddressSet      types.Set    `tfsdk:"address_set"`
}

type securityAddressBookBlockNetworkAddress struct {
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
}

type securityAddressBookBlockDNSName struct {
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
	IPv4Only    types.Bool   `tfsdk:"ipv4_only"`
	IPv6Only    types.Bool   `tfsdk:"ipv6_only"`
}

type securityAddressBookBlockRangeAddress struct {
	Name        types.String `tfsdk:"name"`
	From        types.String `tfsdk:"from"`
	To          types.String `tfsdk:"to"`
	Description types.String `tfsdk:"description"`
}

type securityAddressBookBlockWildcardAddress struct {
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
}

type securityAddressBookBlockAddressSet struct {
	Name        types.String   `tfsdk:"name"`
	Address     []types.String `tfsdk:"address"`
	AddressSet  []types.String `tfsdk:"address_set"`
	Description types.String   `tfsdk:"description"`
}

type securityAddressBookBlockAddressSetConfig struct {
	Name        types.String `tfsdk:"name"`
	Address     types.Set    `tfsdk:"address"`
	AddressSet  types.Set    `tfsdk:"address_set"`
	Description types.String `tfsdk:"description"`
}

func (rsc *securityAddressBook) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityAddressBookConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if (config.Name.IsNull() || config.Name.ValueString() == "global") &&
		!config.AttachZone.IsNull() && !config.AttachZone.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("attach_zone"),
			tfdiag.ConflictConfigErrSummary,
			"cannot attach global address book to a zone",
		)
	}

	addressName := make(map[string]struct{})
	if !config.NetworkAddress.IsNull() && !config.NetworkAddress.IsUnknown() {
		var configNetworkAddress []securityAddressBookBlockNetworkAddress
		asDiags := config.NetworkAddress.ElementsAs(ctx, &configNetworkAddress, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, block := range configNetworkAddress {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("network_address"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}
	if !config.DNSName.IsNull() && !config.DNSName.IsUnknown() {
		var configDNSName []securityAddressBookBlockDNSName
		asDiags := config.DNSName.ElementsAs(ctx, &configDNSName, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, block := range configDNSName {
			if !block.IPv4Only.IsNull() && !block.IPv4Only.IsUnknown() &&
				!block.IPv6Only.IsNull() && !block.IPv6Only.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("dns_name"),
					tfdiag.ConflictConfigErrSummary,
					fmt.Sprintf("ipv4_only and ipv6_only cannot be configured together in dns_name %q", block.Name.ValueString()),
				)
			}
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("dns_name"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}
	if !config.RangeAddress.IsNull() && !config.RangeAddress.IsUnknown() {
		var configRangeAddress []securityAddressBookBlockRangeAddress
		asDiags := config.RangeAddress.ElementsAs(ctx, &configRangeAddress, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, block := range configRangeAddress {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[block.Name.ValueString()]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("range_address"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}
	if !config.WildcardAddress.IsNull() && !config.WildcardAddress.IsUnknown() {
		var configWildcardAddress []securityAddressBookBlockWildcardAddress
		asDiags := config.WildcardAddress.ElementsAs(ctx, &configWildcardAddress, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, block := range configWildcardAddress {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("wildcard_address"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}
	if !config.AddressSet.IsNull() && !config.AddressSet.IsUnknown() {
		var configAddressSet []securityAddressBookBlockAddressSetConfig
		asDiags := config.AddressSet.ElementsAs(ctx, &configAddressSet, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, block := range configAddressSet {
			if block.Address.IsNull() && block.AddressSet.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_set"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("at least one of address or address_set must be specified in address_set %q",
						block.Name.ValueString()),
				)
			}
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_set"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses or address-sets with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}

	if config.Description.IsNull() &&
		config.AttachZone.IsNull() &&
		config.NetworkAddress.IsNull() &&
		config.DNSName.IsNull() &&
		config.RangeAddress.IsNull() &&
		config.WildcardAddress.IsNull() &&
		config.AddressSet.IsNull() {
		if config.Name.IsNull() {
			resp.Diagnostics.AddError(
				tfdiag.MissingConfigErrSummary,
				"resource without argument is not supported",
			)
		} else {
			resp.Diagnostics.AddError(
				tfdiag.MissingConfigErrSummary,
				"resource with only the name argument is not supported",
			)
		}
	}
}

func (rsc *securityAddressBook) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityAddressBookData
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
			bookExists, err := checkSecurityAddressBookExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if bookExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			bookExists, err := checkSecurityAddressBookExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !bookExists {
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

func (rsc *securityAddressBook) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityAddressBookData
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

func (rsc *securityAddressBook) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityAddressBookData
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

func (rsc *securityAddressBook) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityAddressBookData
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

func (rsc *securityAddressBook) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityAddressBookData

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

func checkSecurityAddressBookExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security address-book \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityAddressBookData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityAddressBookData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityAddressBookData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security address-book \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	for _, v := range rscData.AttachZone {
		if rscData.Name.ValueString() == "global" {
			return path.Root("attach_zone"),
				errors.New("cannot attach global address book to a zone")
		}
		configSet = append(configSet, setPrefix+"attach zone "+v.ValueString())
	}
	addressName := make(map[string]struct{})
	for _, block := range rscData.NetworkAddress {
		name := block.Name.ValueString()
		if _, ok := addressName[name]; ok {
			return path.Root("network_address"),
				fmt.Errorf("multiple addresses with the same name %q", name)
		}
		addressName[name] = struct{}{}

		setPrefixAddr := setPrefix + "address " + name + " "
		configSet = append(configSet, setPrefixAddr+block.Value.ValueString())
		if v := block.Description.ValueString(); v != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+v+"\"")
		}
	}
	for _, block := range rscData.DNSName {
		name := block.Name.ValueString()
		if _, ok := addressName[name]; ok {
			return path.Root("dns_name"),
				fmt.Errorf("multiple addresses with the same name %q", name)
		}
		addressName[name] = struct{}{}

		setPrefixAddr := setPrefix + "address " + name + " "
		configSet = append(configSet, setPrefixAddr+"dns-name "+block.Value.ValueString())
		if block.IPv4Only.ValueBool() {
			configSet = append(configSet, setPrefixAddr+"dns-name "+block.Value.ValueString()+" ipv4-only")
		}
		if block.IPv6Only.ValueBool() {
			configSet = append(configSet, setPrefixAddr+"dns-name "+block.Value.ValueString()+" ipv6-only")
		}
		if v := block.Description.ValueString(); v != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+v+"\"")
		}
	}
	for _, block := range rscData.RangeAddress {
		name := block.Name.ValueString()
		if _, ok := addressName[name]; ok {
			return path.Root("range_address"),
				fmt.Errorf("multiple addresses with the same name %q", name)
		}
		addressName[name] = struct{}{}

		setPrefixAddr := setPrefix + "address " + name + " "
		configSet = append(configSet, setPrefixAddr+"range-address "+block.From.ValueString()+" to "+block.To.ValueString())
		if v := block.Description.ValueString(); v != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+v+"\"")
		}
	}
	for _, block := range rscData.WildcardAddress {
		name := block.Name.ValueString()
		if _, ok := addressName[name]; ok {
			return path.Root("wildcard_address"),
				fmt.Errorf("multiple addresses with the same name %q", name)
		}
		addressName[name] = struct{}{}

		setPrefixAddr := setPrefix + "address " + name + " "
		configSet = append(configSet, setPrefixAddr+"wildcard-address "+block.Value.ValueString())
		if v := block.Description.ValueString(); v != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+v+"\"")
		}
	}
	for _, block := range rscData.AddressSet {
		name := block.Name.ValueString()
		if _, ok := addressName[name]; ok {
			return path.Root("address_set"),
				fmt.Errorf("multiple addresses or address-sets with the same name %q", name)
		}
		addressName[name] = struct{}{}

		setPrefixAddrSet := setPrefix + "address-set " + name + " "
		if len(block.Address) == 0 && len(block.AddressSet) == 0 {
			return path.Root("address_set"),
				fmt.Errorf("at least one of address or address_set must be specified in address_set %q", name)
		}
		for _, v := range block.Address {
			configSet = append(configSet, setPrefixAddrSet+"address "+v.ValueString())
		}
		for _, v := range block.AddressSet {
			configSet = append(configSet, setPrefixAddrSet+"address-set "+v.ValueString())
		}
		if v := block.Description.ValueString(); v != "" {
			configSet = append(configSet, setPrefixAddrSet+"description \""+v+"\"")
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityAddressBookData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security address-book \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	descMap := make(map[string]string)
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
			case balt.CutPrefixInString(&itemTrim, "address "):
				itemTrimFields := strings.Split(itemTrim, " ")
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "description "):
					descMap[itemTrimFields[0]] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "dns-name "):
					switch {
					case balt.CutSuffixInString(&itemTrim, " ipv4-only"):
						rscData.DNSName = append(rscData.DNSName, securityAddressBookBlockDNSName{
							Name:     types.StringValue(itemTrimFields[0]),
							Value:    types.StringValue(itemTrim),
							IPv4Only: types.BoolValue(true),
						})
					case balt.CutSuffixInString(&itemTrim, " ipv6-only"):
						rscData.DNSName = append(rscData.DNSName, securityAddressBookBlockDNSName{
							Name:     types.StringValue(itemTrimFields[0]),
							Value:    types.StringValue(itemTrim),
							IPv6Only: types.BoolValue(true),
						})
					default:
						rscData.DNSName = append(rscData.DNSName, securityAddressBookBlockDNSName{
							Name:  types.StringValue(itemTrimFields[0]),
							Value: types.StringValue(itemTrim),
						})
					}
				case balt.CutPrefixInString(&itemTrim, "range-address "):
					rangeAddressFields := strings.Split(itemTrim, " ")
					if len(rangeAddressFields) < 3 { // <from> to <to>
						return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "range-address", itemTrim)
					}
					rscData.RangeAddress = append(rscData.RangeAddress, securityAddressBookBlockRangeAddress{
						Name: types.StringValue(itemTrimFields[0]),
						From: types.StringValue(rangeAddressFields[0]),
						To:   types.StringValue(rangeAddressFields[2]),
					})
				case balt.CutPrefixInString(&itemTrim, "wildcard-address "):
					rscData.WildcardAddress = append(rscData.WildcardAddress, securityAddressBookBlockWildcardAddress{
						Name:  types.StringValue(itemTrimFields[0]),
						Value: types.StringValue(itemTrim),
					})
				default:
					rscData.NetworkAddress = append(rscData.NetworkAddress, securityAddressBookBlockNetworkAddress{
						Name:  types.StringValue(itemTrimFields[0]),
						Value: types.StringValue(itemTrim),
					})
				}
			case balt.CutPrefixInString(&itemTrim, "address-set "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var adSet securityAddressBookBlockAddressSet
				rscData.AddressSet, adSet = tfdata.ExtractBlockWithTFTypesString(rscData.AddressSet, "Name", itemTrimFields[0])
				adSet.Name = types.StringValue(itemTrimFields[0])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "description "):
					adSet.Description = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "address "):
					adSet.Address = append(adSet.Address, types.StringValue(itemTrim))
				case balt.CutPrefixInString(&itemTrim, "address-set "):
					adSet.AddressSet = append(adSet.AddressSet, types.StringValue(itemTrim))
				}
				rscData.AddressSet = append(rscData.AddressSet, adSet)
			case balt.CutPrefixInString(&itemTrim, "attach zone "):
				rscData.AttachZone = append(rscData.AttachZone, types.StringValue(itemTrim))
			}
		}
	}
	// copy description to struct
	for i, b := range rscData.NetworkAddress {
		if v, ok := descMap[b.Name.ValueString()]; ok {
			rscData.NetworkAddress[i].Description = types.StringValue(v)
		}
	}
	for i, b := range rscData.DNSName {
		if v, ok := descMap[b.Name.ValueString()]; ok {
			rscData.DNSName[i].Description = types.StringValue(v)
		}
	}
	for i, b := range rscData.RangeAddress {
		if v, ok := descMap[b.Name.ValueString()]; ok {
			rscData.RangeAddress[i].Description = types.StringValue(v)
		}
	}
	for i, b := range rscData.WildcardAddress {
		if v, ok := descMap[b.Name.ValueString()]; ok {
			rscData.WildcardAddress[i].Description = types.StringValue(v)
		}
	}

	return nil
}

func (rscData *securityAddressBookData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security address-book \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
