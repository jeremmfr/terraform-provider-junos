package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
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
	_ resource.Resource                   = &securityZone{}
	_ resource.ResourceWithConfigure      = &securityZone{}
	_ resource.ResourceWithValidateConfig = &securityZone{}
	_ resource.ResourceWithImportState    = &securityZone{}
)

type securityZone struct {
	client *junos.Client
}

func newSecurityZoneResource() resource.Resource {
	return &securityZone{}
}

func (rsc *securityZone) typeName() string {
	return providerName + "_security_zone"
}

func (rsc *securityZone) junosName() string {
	return "security zone"
}

func (rsc *securityZone) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityZone) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityZone) Configure(
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

func (rsc *securityZone) Schema(
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
				Required:    true,
				Description: "The name of security zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"address_book_configure_singly": schema.BoolAttribute{
				Optional: true,
				Description: "Disable management of address-book in this resource " +
					"to be able to manage them with specific resources.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"advance_policy_based_routing_profile": schema.StringAttribute{
				Optional:    true,
				Description: "Enable Advance Policy Based Routing on this zone with a profile.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"application_tracking": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable Application tracking support for this zone.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of zone.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"inbound_protocols": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The inbound protocols allowed.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 32),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"inbound_services": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The inbound services allowed.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 32),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"reverse_reroute": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable Reverse route lookup when there is change in ingress interface.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"screen": schema.StringAttribute{
				Optional:    true,
				Description: "Name of ids option object (screen) applied to the zone.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"source_identity_log": schema.BoolAttribute{
				Optional:    true,
				Description: "Show user and group info in session log for this zone.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"tcp_rst": schema.BoolAttribute{
				Optional:    true,
				Description: "Send RST for NON-SYN packet not matching TCP session.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"address_book": schema.SetNestedBlock{
				Description: "For each name of network address to declare.",
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
						"network": schema.StringAttribute{
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
			"address_book_dns": schema.SetNestedBlock{
				Description: "For each name of dns-name address to declare.",
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
						"fqdn": schema.StringAttribute{
							Required:    true,
							Description: "Fully qualified domain name.",
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
			"address_book_range": schema.SetNestedBlock{
				Description: "For each name of range-address to declare.",
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
							Description: "Lower limit of address range.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress(),
							},
						},
						"to": schema.StringAttribute{
							Required:    true,
							Description: "Upper limit of address range.",
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
			"address_book_set": schema.SetNestedBlock{
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
			"address_book_wildcard": schema.SetNestedBlock{
				Description: "For each name of wildcard-address to declare.",
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
						"network": schema.StringAttribute{
							Required:    true,
							Description: "Numeric IPv4 wildcard address with in the form of a.d.d.r/netmask.",
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
		},
	}
}

type securityZoneData struct {
	ID                               types.String                           `tfsdk:"id"`
	Name                             types.String                           `tfsdk:"name"`
	AddressBookConfigureSingly       types.Bool                             `tfsdk:"address_book_configure_singly"`
	AdvancePolicyBasedRoutingProfile types.String                           `tfsdk:"advance_policy_based_routing_profile"`
	ApplicationTracking              types.Bool                             `tfsdk:"application_tracking"`
	Description                      types.String                           `tfsdk:"description"`
	InboundProtocols                 []types.String                         `tfsdk:"inbound_protocols"`
	InboundServices                  []types.String                         `tfsdk:"inbound_services"`
	ReverseReroute                   types.Bool                             `tfsdk:"reverse_reroute"`
	Screen                           types.String                           `tfsdk:"screen"`
	SourceIdentityLog                types.Bool                             `tfsdk:"source_identity_log"`
	TCPRst                           types.Bool                             `tfsdk:"tcp_rst"`
	AddressBook                      []securityZoneBlockAddressBook         `tfsdk:"address_book"`
	AddressBookDNS                   []securityZoneBlockAddressBookDNS      `tfsdk:"address_book_dns"`
	AddressBookRange                 []securityZoneBlockAddressBookRange    `tfsdk:"address_book_range"`
	AddressBookSet                   []securityZoneBlockAddressBookSet      `tfsdk:"address_book_set"`
	AddressBookWildcard              []securityZoneBlockAddressBookWildcard `tfsdk:"address_book_wildcard"`
	Interface                        []securityZoneDataSourceBlockInterface `tfsdk:"-"` // to data source
}

type securityZoneConfig struct {
	ID                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	AddressBookConfigureSingly       types.Bool   `tfsdk:"address_book_configure_singly"`
	AdvancePolicyBasedRoutingProfile types.String `tfsdk:"advance_policy_based_routing_profile"`
	ApplicationTracking              types.Bool   `tfsdk:"application_tracking"`
	Description                      types.String `tfsdk:"description"`
	InboundProtocols                 types.Set    `tfsdk:"inbound_protocols"`
	InboundServices                  types.Set    `tfsdk:"inbound_services"`
	ReverseReroute                   types.Bool   `tfsdk:"reverse_reroute"`
	Screen                           types.String `tfsdk:"screen"`
	SourceIdentityLog                types.Bool   `tfsdk:"source_identity_log"`
	TCPRst                           types.Bool   `tfsdk:"tcp_rst"`
	AddressBook                      types.Set    `tfsdk:"address_book"`
	AddressBookDNS                   types.Set    `tfsdk:"address_book_dns"`
	AddressBookRange                 types.Set    `tfsdk:"address_book_range"`
	AddressBookSet                   types.Set    `tfsdk:"address_book_set"`
	AddressBookWildcard              types.Set    `tfsdk:"address_book_wildcard"`
}

type securityZoneBlockAddressBook struct {
	Name        types.String `tfsdk:"name"`
	Network     types.String `tfsdk:"network"`
	Description types.String `tfsdk:"description"`
}

type securityZoneBlockAddressBookDNS struct {
	Name        types.String `tfsdk:"name"`
	FQDN        types.String `tfsdk:"fqdn"`
	Description types.String `tfsdk:"description"`
	IPv4Only    types.Bool   `tfsdk:"ipv4_only"`
	IPv6Only    types.Bool   `tfsdk:"ipv6_only"`
}

type securityZoneBlockAddressBookRange struct {
	Name        types.String `tfsdk:"name"`
	From        types.String `tfsdk:"from"`
	To          types.String `tfsdk:"to"`
	Description types.String `tfsdk:"description"`
}

type securityZoneBlockAddressBookSet struct {
	Name        types.String   `tfsdk:"name"`
	Address     []types.String `tfsdk:"address"`
	AddressSet  []types.String `tfsdk:"address_set"`
	Description types.String   `tfsdk:"description"`
}

type securityZoneBlockAddressBookSetConfig struct {
	Name        types.String `tfsdk:"name"`
	Address     types.Set    `tfsdk:"address"`
	AddressSet  types.Set    `tfsdk:"address_set"`
	Description types.String `tfsdk:"description"`
}

type securityZoneBlockAddressBookWildcard struct {
	Name        types.String `tfsdk:"name"`
	Network     types.String `tfsdk:"network"`
	Description types.String `tfsdk:"description"`
}

func (rsc *securityZone) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityZoneConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.AddressBookConfigureSingly.ValueBool() &&
		((!config.AddressBook.IsNull() && !config.AddressBook.IsUnknown()) ||
			(!config.AddressBookDNS.IsNull() && !config.AddressBookDNS.IsUnknown()) ||
			(!config.AddressBookRange.IsNull() && !config.AddressBookRange.IsUnknown()) ||
			(!config.AddressBookSet.IsNull() && !config.AddressBookSet.IsUnknown()) ||
			(!config.AddressBookWildcard.IsNull() && !config.AddressBookWildcard.IsUnknown())) {
		resp.Diagnostics.AddAttributeError(
			path.Root("address_book_configure_singly"),
			tfdiag.ConflictConfigErrSummary,
			"cannot have address_book_configure_singly and want to configure address book at the same time",
		)
	}

	addressName := make(map[string]struct{})
	if !config.AddressBook.IsNull() && !config.AddressBook.IsUnknown() {
		var configAddressBook []securityZoneBlockAddressBook
		asDiags := config.AddressBook.ElementsAs(ctx, &configAddressBook, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, block := range configAddressBook {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}
	if !config.AddressBookDNS.IsNull() && !config.AddressBookDNS.IsUnknown() {
		var configAddressBookDNS []securityZoneBlockAddressBookDNS
		asDiags := config.AddressBookDNS.ElementsAs(ctx, &configAddressBookDNS, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, block := range configAddressBookDNS {
			if block.IPv4Only.ValueBool() && block.IPv6Only.ValueBool() {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book_dns"),
					tfdiag.ConflictConfigErrSummary,
					fmt.Sprintf("ipv4_only and ipv6_only cannot be configured together in address_book_dns %q",
						block.Name.ValueString()),
				)
			}
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book_dns"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}
	if !config.AddressBookRange.IsNull() && !config.AddressBookRange.IsUnknown() {
		var configAddressBookRange []securityZoneBlockAddressBookRange
		asDiags := config.AddressBookRange.ElementsAs(ctx, &configAddressBookRange, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, block := range configAddressBookRange {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book_range"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}
	if !config.AddressBookSet.IsNull() && !config.AddressBookSet.IsUnknown() {
		var configAddressBookSet []securityZoneBlockAddressBookSetConfig
		asDiags := config.AddressBookSet.ElementsAs(ctx, &configAddressBookSet, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, block := range configAddressBookSet {
			if block.Address.IsNull() && block.AddressSet.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book_set"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("at least one of address or address_set must be specified in address_book_set %q",
						block.Name.ValueString()),
				)
			}
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book_set"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses or address-sets with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}
	if !config.AddressBookWildcard.IsNull() && !config.AddressBookWildcard.IsUnknown() {
		var configAddressBookWildcard []securityZoneBlockAddressBookWildcard
		asDiags := config.AddressBookWildcard.ElementsAs(ctx, &configAddressBookWildcard, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, block := range configAddressBookWildcard {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book_wildcard"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}
}

func (rsc *securityZone) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityZoneData
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
			zoneExists, err := checkSecurityZonesExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if zoneExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			zoneExists, err := checkSecurityZonesExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !zoneExists {
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

func (rsc *securityZone) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityZoneData
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
			data.AddressBookConfigureSingly = state.AddressBookConfigureSingly
			if data.AddressBookConfigureSingly.ValueBool() {
				data.AddressBook = nil
				data.AddressBookDNS = nil
				data.AddressBookRange = nil
				data.AddressBookSet = nil
				data.AddressBookWildcard = nil
			}
		},
		resp,
	)
}

func (rsc *securityZone) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityZoneData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	addressBookConfiguredSingly := plan.AddressBookConfigureSingly.ValueBool()
	if !plan.AddressBookConfigureSingly.Equal(state.AddressBookConfigureSingly) {
		if state.AddressBookConfigureSingly.ValueBool() {
			addressBookConfiguredSingly = state.AddressBookConfigureSingly.ValueBool()
			resp.Diagnostics.AddAttributeWarning(
				path.Root("address_book_configure_singly"),
				"Disable address_book_configure_singly on resource already created",
				"It's doesn't delete addresses and address-sets already configured. "+
					"So refresh resource after apply to detect address-book entries that need to be deleted",
			)
		} else {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("address_book_configure_singly"),
				"Enable address_book_configure_singly on resource already created",
				"It's doesn't delete addresses and address-sets already configured. "+
					"So import address-book entries in dedicated resource(s) to be able to manage them",
			)
		}
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.delOpts(ctx, addressBookConfiguredSingly, junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

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

	if err := state.delOpts(ctx, addressBookConfiguredSingly, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

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

func (rsc *securityZone) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityZoneData
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

func (rsc *securityZone) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityZoneData

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

func checkSecurityZonesExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityZoneData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityZoneData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityZoneData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security zones security-zone " + rscData.Name.ValueString() + " "

	configSet = append(configSet, setPrefix)
	if !rscData.AddressBookConfigureSingly.ValueBool() {
		addressName := make(map[string]struct{})
		for _, block := range rscData.AddressBook {
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				return path.Root("address_book"),
					fmt.Errorf("multiple addresses with the same name %q", name)
			}
			addressName[name] = struct{}{}

			configSet = append(configSet, setPrefix+"address-book address "+name+" "+block.Network.ValueString())
			if v := block.Description.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"address-book address "+name+" description \""+v+"\"")
			}
		}
		for _, block := range rscData.AddressBookDNS {
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				return path.Root("address_book_dns"),
					fmt.Errorf("multiple addresses with the same name %q", name)
			}
			addressName[name] = struct{}{}

			setLine := setPrefix + "address-book address " + name + " dns-name " + block.FQDN.ValueString()
			configSet = append(configSet, setLine)
			if block.IPv4Only.ValueBool() {
				configSet = append(configSet, setLine+" ipv4-only")
			}
			if block.IPv6Only.ValueBool() {
				configSet = append(configSet, setLine+" ipv6-only")
			}
			if v := block.Description.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"address-book address "+name+" description \""+v+"\"")
			}
		}
		for _, block := range rscData.AddressBookRange {
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				return path.Root("address_book_range"),
					fmt.Errorf("multiple addresses with the same name %q", name)
			}
			addressName[name] = struct{}{}

			configSet = append(configSet, setPrefix+"address-book address "+
				name+" range-address "+block.From.ValueString()+
				" to "+block.To.ValueString())
			if v := block.Description.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"address-book address "+name+" description \""+v+"\"")
			}
		}
		for _, block := range rscData.AddressBookSet {
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				return path.Root("address_book_set"),
					fmt.Errorf("multiple addresses or address-sets with the same name %q", name)
			}
			addressName[name] = struct{}{}

			if len(block.Address) == 0 &&
				len(block.AddressSet) == 0 {
				return path.Root("address_book_set"),
					fmt.Errorf("at least one of address or address_set must be specified in address_set %q", name)
			}
			for _, v := range block.Address {
				configSet = append(configSet, setPrefix+"address-book address-set "+name+" address "+v.ValueString())
			}
			for _, v := range block.AddressSet {
				configSet = append(configSet, setPrefix+"address-book address-set "+name+" address-set "+v.ValueString())
			}
			if v := block.Description.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"address-book address-set "+name+" description \""+v+"\"")
			}
		}
		for _, block := range rscData.AddressBookWildcard {
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				return path.Root("address_book_wildcard"),
					fmt.Errorf("multiple addresses with the same name %q", name)
			}
			addressName[name] = struct{}{}

			configSet = append(configSet, setPrefix+"address-book address "+
				name+" wildcard-address "+block.Network.ValueString())
			if v := block.Description.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"address-book address "+name+" description \""+v+"\"")
			}
		}
	}
	if v := rscData.AdvancePolicyBasedRoutingProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"advance-policy-based-routing-profile \""+v+"\"")
	}
	if rscData.ApplicationTracking.ValueBool() {
		configSet = append(configSet, setPrefix+"application-tracking")
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	for _, v := range rscData.InboundProtocols {
		configSet = append(configSet, setPrefix+"host-inbound-traffic protocols \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.InboundServices {
		configSet = append(configSet, setPrefix+"host-inbound-traffic system-services \""+v.ValueString()+"\"")
	}
	if rscData.ReverseReroute.ValueBool() {
		configSet = append(configSet, setPrefix+"enable-reverse-reroute")
	}
	if v := rscData.Screen.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"screen \""+v+"\"")
	}
	if rscData.SourceIdentityLog.ValueBool() {
		configSet = append(configSet, setPrefix+"source-identity-log")
	}
	if rscData.TCPRst.ValueBool() {
		configSet = append(configSet, setPrefix+"tcp-rst")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityZoneData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	descAddressBookMap := make(map[string]string)
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
			case balt.CutPrefixInString(&itemTrim, "address-book address "):
				itemTrimFields := strings.Split(itemTrim, " ")
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "description "):
					descAddressBookMap[itemTrimFields[0]] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "dns-name "):
					switch {
					case balt.CutSuffixInString(&itemTrim, " ipv4-only"):
						rscData.AddressBookDNS = append(rscData.AddressBookDNS, securityZoneBlockAddressBookDNS{
							Name:     types.StringValue(itemTrimFields[0]),
							FQDN:     types.StringValue(itemTrim),
							IPv4Only: types.BoolValue(true),
						})
					case balt.CutSuffixInString(&itemTrim, " ipv6-only"):
						rscData.AddressBookDNS = append(rscData.AddressBookDNS, securityZoneBlockAddressBookDNS{
							Name:     types.StringValue(itemTrimFields[0]),
							FQDN:     types.StringValue(itemTrim),
							IPv6Only: types.BoolValue(true),
						})
					default:
						rscData.AddressBookDNS = append(rscData.AddressBookDNS, securityZoneBlockAddressBookDNS{
							Name: types.StringValue(itemTrimFields[0]),
							FQDN: types.StringValue(itemTrim),
						})
					}
				case balt.CutPrefixInString(&itemTrim, "range-address "):
					rangeAddressFields := strings.Split(itemTrim, " ")
					if len(rangeAddressFields) < 3 { // <from> to <to>
						return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "range-address", itemTrim)
					}
					rscData.AddressBookRange = append(rscData.AddressBookRange, securityZoneBlockAddressBookRange{
						Name: types.StringValue(itemTrimFields[0]),
						From: types.StringValue(rangeAddressFields[0]),
						To:   types.StringValue(rangeAddressFields[2]),
					})
				case balt.CutPrefixInString(&itemTrim, "wildcard-address "):
					rscData.AddressBookWildcard = append(rscData.AddressBookWildcard, securityZoneBlockAddressBookWildcard{
						Name:    types.StringValue(itemTrimFields[0]),
						Network: types.StringValue(itemTrim),
					})
				default:
					rscData.AddressBook = append(rscData.AddressBook, securityZoneBlockAddressBook{
						Name:    types.StringValue(itemTrimFields[0]),
						Network: types.StringValue(itemTrim),
					})
				}
			case balt.CutPrefixInString(&itemTrim, "address-book address-set "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var adSet securityZoneBlockAddressBookSet
				rscData.AddressBookSet, adSet = tfdata.ExtractBlockWithTFTypesString(
					rscData.AddressBookSet, "Name", itemTrimFields[0])
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
				rscData.AddressBookSet = append(rscData.AddressBookSet, adSet)
			case balt.CutPrefixInString(&itemTrim, "advance-policy-based-routing-profile "):
				rscData.AdvancePolicyBasedRoutingProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "application-tracking":
				rscData.ApplicationTracking = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "host-inbound-traffic protocols "):
				rscData.InboundProtocols = append(rscData.InboundProtocols, types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "host-inbound-traffic system-services "):
				rscData.InboundServices = append(rscData.InboundServices, types.StringValue(strings.Trim(itemTrim, "\"")))
			case itemTrim == "enable-reverse-reroute":
				rscData.ReverseReroute = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "screen "):
				rscData.Screen = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "source-identity-log":
				rscData.SourceIdentityLog = types.BoolValue(true)
			case itemTrim == "tcp-rst":
				rscData.TCPRst = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "interfaces "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var interFace securityZoneDataSourceBlockInterface
				rscData.Interface, interFace = tfdata.ExtractBlockWithTFTypesString(rscData.Interface, "Name", itemTrimFields[0])
				interFace.Name = types.StringValue(itemTrimFields[0])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "host-inbound-traffic protocols "):
					interFace.InboundProtocols = append(interFace.InboundProtocols, types.StringValue(itemTrim))
				case balt.CutPrefixInString(&itemTrim, "host-inbound-traffic system-services "):
					interFace.InboundServices = append(interFace.InboundServices, types.StringValue(itemTrim))
				}
				rscData.Interface = append(rscData.Interface, interFace)
			}
		}
	}
	// copy description to struct
	for i, b := range rscData.AddressBook {
		if v, ok := descAddressBookMap[b.Name.ValueString()]; ok {
			rscData.AddressBook[i].Description = types.StringValue(v)
		}
	}
	for i, b := range rscData.AddressBookDNS {
		if v, ok := descAddressBookMap[b.Name.ValueString()]; ok {
			rscData.AddressBookDNS[i].Description = types.StringValue(v)
		}
	}
	for i, b := range rscData.AddressBookRange {
		if v, ok := descAddressBookMap[b.Name.ValueString()]; ok {
			rscData.AddressBookRange[i].Description = types.StringValue(v)
		}
	}
	for i, b := range rscData.AddressBookWildcard {
		if v, ok := descAddressBookMap[b.Name.ValueString()]; ok {
			rscData.AddressBookWildcard[i].Description = types.StringValue(v)
		}
	}

	return nil
}

func (rscData *securityZoneData) delOpts(
	_ context.Context, addressBookSingly bool, junSess *junos.Session,
) error {
	delPrefix := "delete security zones security-zone " + rscData.Name.ValueString() + " "

	configSet := []string{
		delPrefix + "advance-policy-based-routing-profile",
		delPrefix + "description",
		delPrefix + "application-tracking",
		delPrefix + "host-inbound-traffic",
		delPrefix + "enable-reverse-reroute",
		delPrefix + "screen",
		delPrefix + "source-identity-log",
		delPrefix + "tcp-rst",
	}
	if !addressBookSingly {
		configSet = append(configSet, delPrefix+"address-book")
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *securityZoneData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security zones security-zone " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
