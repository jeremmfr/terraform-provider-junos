package provider

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &securityZoneOrdered{}
	_ resource.ResourceWithConfigure      = &securityZoneOrdered{}
	_ resource.ResourceWithValidateConfig = &securityZoneOrdered{}
	_ resource.ResourceWithImportState    = &securityZoneOrdered{}
)

type securityZoneOrdered struct {
	client *junos.Client
}

func newSecurityZoneOrderedResource() resource.Resource {
	return &securityZoneOrdered{}
}

func (rsc *securityZoneOrdered) typeName() string {
	return providerName + "_security_zone_ordered"
}

func (rsc *securityZoneOrdered) junosName() string {
	return "security zone"
}

func (rsc *securityZoneOrdered) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityZoneOrdered) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityZoneOrdered) Configure(
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

func (rsc *securityZoneOrdered) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes:  securityZoneData{}.attributesSchema(),
		Blocks: map[string]schema.Block{
			"address_book": schema.ListNestedBlock{
				Description: "For each name of network address to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityZoneBlockAddressBook{}.attributesSchema(),
				},
			},
			"address_book_dns": schema.ListNestedBlock{
				Description: "For each name of dns-name address to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityZoneBlockAddressBookDNS{}.attributesSchema(),
				},
			},
			"address_book_range": schema.ListNestedBlock{
				Description: "For each name of range-address to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityZoneBlockAddressBookRange{}.attributesSchema(),
				},
			},
			"address_book_set": schema.ListNestedBlock{
				Description: "For each name of address-set to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityZoneBlockAddressBookSet{}.attributesSchema(),
				},
			},
			"address_book_wildcard": schema.ListNestedBlock{
				Description: "For each name of wildcard-address to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityZoneBlockAddressBookWildcard{}.attributesSchema(),
				},
			},
		},
	}
}

type securityZoneOrderedConfig struct {
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
	AddressBook                      types.List   `tfsdk:"address_book"`
	AddressBookDNS                   types.List   `tfsdk:"address_book_dns"`
	AddressBookRange                 types.List   `tfsdk:"address_book_range"`
	AddressBookSet                   types.List   `tfsdk:"address_book_set"`
	AddressBookWildcard              types.List   `tfsdk:"address_book_wildcard"`
}

func (rsc *securityZoneOrdered) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityZoneOrderedConfig
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
		for i, block := range configAddressBook {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book").AtListIndex(i).AtName("name"),
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
		for i, block := range configAddressBookDNS {
			if block.IPv4Only.ValueBool() && block.IPv6Only.ValueBool() {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book_dns").AtListIndex(i).AtName("ipv4_only"),
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
					path.Root("address_book_dns").AtListIndex(i).AtName("name"),
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
		for i, block := range configAddressBookRange {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book_range").AtListIndex(i).AtName("name"),
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
		for i, block := range configAddressBookSet {
			if block.Address.IsNull() && block.AddressSet.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book_set").AtListIndex(i).AtName("name"),
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
					path.Root("address_book_set").AtListIndex(i).AtName("name"),
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
		for i, block := range configAddressBookWildcard {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_book_wildcard").AtListIndex(i).AtName("name"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple addresses with the same name %q", name),
				)
			}
			addressName[name] = struct{}{}
		}
	}
}

func (rsc *securityZoneOrdered) Create(
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

func (rsc *securityZoneOrdered) Read(
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

func (rsc *securityZoneOrdered) Update(
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
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigUnlockWarnSummary, junSess.ConfigUnlock())...)
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

func (rsc *securityZoneOrdered) Delete(
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

func (rsc *securityZoneOrdered) ImportState(
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
