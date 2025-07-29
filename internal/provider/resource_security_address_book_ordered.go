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
	_ resource.Resource                   = &securityAddressBookOrdered{}
	_ resource.ResourceWithConfigure      = &securityAddressBookOrdered{}
	_ resource.ResourceWithValidateConfig = &securityAddressBookOrdered{}
	_ resource.ResourceWithImportState    = &securityAddressBookOrdered{}
)

type securityAddressBookOrdered struct {
	client *junos.Client
}

func newSecurityAddressBookOrderedResource() resource.Resource {
	return &securityAddressBookOrdered{}
}

func (rsc *securityAddressBookOrdered) typeName() string {
	return providerName + "_security_address_book_ordered"
}

func (rsc *securityAddressBookOrdered) junosName() string {
	return "security address book"
}

func (rsc *securityAddressBookOrdered) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityAddressBookOrdered) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityAddressBookOrdered) Configure(
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

func (rsc *securityAddressBookOrdered) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes:  securityAddressBookData{}.attributesSchema(),
		Blocks: map[string]schema.Block{
			"network_address": schema.ListNestedBlock{
				Description: "For each name of network address.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityAddressBookBlockNetworkAddress{}.attributesSchema(),
				},
			},
			"dns_name": schema.ListNestedBlock{
				Description: "For each name of dns name address.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityAddressBookBlockDNSName{}.attributesSchema(),
				},
			},
			"range_address": schema.ListNestedBlock{
				Description: "For each name of range address.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityAddressBookBlockRangeAddress{}.attributesSchema(),
				},
			},
			"wildcard_address": schema.ListNestedBlock{
				Description: "For each name of wildcard address.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityAddressBookBlockWildcardAddress{}.attributesSchema(),
				},
			},
			"address_set": schema.ListNestedBlock{
				Description: "For each name of address-set to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityAddressBookBlockAddressSet{}.attributesSchema(),
				},
			},
		},
	}
}

type securityAddressBookOrderedConfig struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	AttachZone      types.List   `tfsdk:"attach_zone"`
	NetworkAddress  types.List   `tfsdk:"network_address"`
	DNSName         types.List   `tfsdk:"dns_name"`
	RangeAddress    types.List   `tfsdk:"range_address"`
	WildcardAddress types.List   `tfsdk:"wildcard_address"`
	AddressSet      types.List   `tfsdk:"address_set"`
}

func (rsc *securityAddressBookOrdered) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityAddressBookOrderedConfig
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
		for i, block := range configNetworkAddress {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("network_address").AtListIndex(i).AtName("name"),
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
		for i, block := range configDNSName {
			if !block.IPv4Only.IsNull() && !block.IPv4Only.IsUnknown() &&
				!block.IPv6Only.IsNull() && !block.IPv6Only.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("dns_name").AtListIndex(i).AtName("ipv4_only"),
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
					path.Root("dns_name").AtListIndex(i).AtName("name"),
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
		for i, block := range configRangeAddress {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[block.Name.ValueString()]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("range_address").AtListIndex(i).AtName("name"),
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
		for i, block := range configWildcardAddress {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := addressName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("wildcard_address").AtListIndex(i).AtName("name"),
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
		for i, block := range configAddressSet {
			if block.Address.IsNull() && block.AddressSet.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("address_set").AtListIndex(i).AtName("name"),
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
					path.Root("address_set").AtListIndex(i).AtName("name"),
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

func (rsc *securityAddressBookOrdered) Create(
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

func (rsc *securityAddressBookOrdered) Read(
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

func (rsc *securityAddressBookOrdered) Update(
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

func (rsc *securityAddressBookOrdered) Delete(
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

func (rsc *securityAddressBookOrdered) ImportState(
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
