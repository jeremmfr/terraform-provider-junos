package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

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
	_ resource.Resource                   = &securityZoneBookAddress{}
	_ resource.ResourceWithConfigure      = &securityZoneBookAddress{}
	_ resource.ResourceWithValidateConfig = &securityZoneBookAddress{}
	_ resource.ResourceWithImportState    = &securityZoneBookAddress{}
)

type securityZoneBookAddress struct {
	client *junos.Client
}

func newSecurityZoneBookAddressResource() resource.Resource {
	return &securityZoneBookAddress{}
}

func (rsc *securityZoneBookAddress) typeName() string {
	return providerName + "_security_zone_book_address"
}

func (rsc *securityZoneBookAddress) junosName() string {
	return "security zone address-book address"
}

func (rsc *securityZoneBookAddress) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityZoneBookAddress) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityZoneBookAddress) Configure(
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

func (rsc *securityZoneBookAddress) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Provides an address resource in address-book of security zone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<zone>" + junos.IDSeparator + "<name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.AddressNameFormat),
				},
			},
			"zone": schema.StringAttribute{
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
			"cidr": schema.StringAttribute{
				Optional:    true,
				Description: "CIDR of address.",
				Validators: []validator.String{
					tfvalidator.StringCIDRNetwork(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of address.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"dns_ipv4_only": schema.BoolAttribute{
				Optional:    true,
				Description: "IPv4 dns address.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"dns_ipv6_only": schema.BoolAttribute{
				Optional:    true,
				Description: "IPv6 dns address.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"dns_name": schema.StringAttribute{
				Optional:    true,
				Description: "DNS address name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 253),
					tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
				},
			},
			"range_from": schema.StringAttribute{
				Optional:    true,
				Description: "Lower limit of address range.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
			"range_to": schema.StringAttribute{
				Optional:    true,
				Description: "Upper limit of address range.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress(),
				},
			},
			"wildcard": schema.StringAttribute{
				Optional:    true,
				Description: "Numeric IPv4 wildcard address in the form of `a.d.d.r/netmask`.",
				Validators: []validator.String{
					tfvalidator.StringWildcardNetwork(),
				},
			},
		},
	}
}

type securityZoneBookAddressData struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Zone        types.String `tfsdk:"zone"`
	CIDR        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
	DNSIPv4Only types.Bool   `tfsdk:"dns_ipv4_only"`
	DNSIPv6Only types.Bool   `tfsdk:"dns_ipv6_only"`
	DNSName     types.String `tfsdk:"dns_name"`
	RangeFrom   types.String `tfsdk:"range_from"`
	RangeTo     types.String `tfsdk:"range_to"`
	Wildcard    types.String `tfsdk:"wildcard"`
}

func (rsc *securityZoneBookAddress) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityZoneBookAddressData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.DNSIPv4Only.IsNull() && config.DNSName.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dns_ipv4_only"),
			tfdiag.MissingConfigErrSummary,
			"cannot have dns_ipv4_only without dns_name",
		)
	}
	if !config.DNSIPv6Only.IsNull() && config.DNSName.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dns_ipv6_only"),
			tfdiag.MissingConfigErrSummary,
			"cannot have dns_ipv6_only without dns_name",
		)
	}
	if !config.DNSIPv4Only.IsNull() && !config.DNSIPv4Only.IsUnknown() &&
		!config.DNSIPv6Only.IsNull() && !config.DNSIPv6Only.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dns_ipv4_only"),
			tfdiag.ConflictConfigErrSummary,
			"only one of dns_ipv4_only or dns_ipv6_only can be specified",
		)
	}
	if !config.RangeTo.IsNull() && config.RangeFrom.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("range_to"),
			tfdiag.MissingConfigErrSummary,
			"cannot have range_to without range_from",
		)
	}
	if !config.RangeFrom.IsNull() && config.RangeTo.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("range_from"),
			tfdiag.MissingConfigErrSummary,
			"cannot have range_from without range_to",
		)
	}
	switch {
	case !config.CIDR.IsNull() && !config.CIDR.IsUnknown():
		if (!config.DNSName.IsNull() && !config.DNSName.IsUnknown()) ||
			(!config.RangeFrom.IsNull() && !config.RangeFrom.IsUnknown()) ||
			(!config.Wildcard.IsNull() && !config.Wildcard.IsUnknown()) {
			resp.Diagnostics.AddAttributeError(
				path.Root("cidr"),
				tfdiag.ConflictConfigErrSummary,
				"only one of cidr, dns_name, range_from or wildcard must be specified",
			)
		}
	case !config.DNSName.IsNull() && !config.DNSName.IsUnknown():
		if (!config.CIDR.IsNull() && !config.CIDR.IsUnknown()) ||
			(!config.RangeFrom.IsNull() && !config.RangeFrom.IsUnknown()) ||
			(!config.Wildcard.IsNull() && !config.Wildcard.IsUnknown()) {
			resp.Diagnostics.AddAttributeError(
				path.Root("dns_name"),
				tfdiag.ConflictConfigErrSummary,
				"only one of cidr, dns_name, range_from or wildcard must be specified",
			)
		}
	case !config.RangeFrom.IsNull() && !config.RangeFrom.IsUnknown():
		if (!config.CIDR.IsNull() && !config.CIDR.IsUnknown()) ||
			(!config.DNSName.IsNull() && !config.DNSName.IsUnknown()) ||
			(!config.Wildcard.IsNull() && !config.Wildcard.IsUnknown()) {
			resp.Diagnostics.AddAttributeError(
				path.Root("range_from"),
				tfdiag.ConflictConfigErrSummary,
				"only one of cidr, dns_name, range_from or wildcard must be specified",
			)
		}
	case !config.Wildcard.IsNull() && !config.Wildcard.IsUnknown():
		if (!config.CIDR.IsNull() && !config.CIDR.IsUnknown()) ||
			(!config.DNSName.IsNull() && !config.DNSName.IsUnknown()) ||
			(!config.RangeFrom.IsNull() && !config.RangeFrom.IsUnknown()) {
			resp.Diagnostics.AddAttributeError(
				path.Root("wildcard"),
				tfdiag.ConflictConfigErrSummary,
				"only one of cidr, dns_name, range_from or wildcard must be specified",
			)
		}
	case config.CIDR.IsNull() && config.DNSName.IsNull() && config.RangeFrom.IsNull() && config.Wildcard.IsNull():
		resp.Diagnostics.AddError(
			tfdiag.MissingConfigErrSummary,
			"one of cidr, dns_name, range_from or wildcard must be specified",
		)
	}
}

func (rsc *securityZoneBookAddress) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityZoneBookAddressData
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
	if plan.Zone.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("zone"),
			"Empty Zone",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "zone"),
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
			zonesExists, err := checkSecurityZonesExists(fnCtx, plan.Zone.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !zonesExists {
				resp.Diagnostics.AddAttributeError(
					path.Root("zone"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("security zone %q doesn't exist", plan.Zone.ValueString()),
				)

				return false
			}
			addressExists, err := checkSecurityZoneBookAddressExists(
				fnCtx,
				plan.Zone.ValueString(),
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if addressExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf(rsc.junosName()+" %q already exists in zone %q",
						plan.Name.ValueString(), plan.Zone.ValueString()),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			addressExists, err := checkSecurityZoneBookAddressExists(
				fnCtx,
				plan.Zone.ValueString(),
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !addressExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					fmt.Sprintf(rsc.junosName()+" %q does not exists in zone %q after commit "+
						"=> check your config", plan.Name.ValueString(), plan.Zone.ValueString()),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *securityZoneBookAddress) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityZoneBookAddressData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Zone.ValueString(),
			state.Name.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *securityZoneBookAddress) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityZoneBookAddressData
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

func (rsc *securityZoneBookAddress) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityZoneBookAddressData
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

func (rsc *securityZoneBookAddress) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityZoneBookAddressData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <zone>"+junos.IDSeparator+"<name>)",
	)
}

func checkSecurityZoneBookAddressExists(
	_ context.Context, zone, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + zone + " address-book address " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityZoneBookAddressData) fillID() {
	rscData.ID = types.StringValue(rscData.Zone.ValueString() + junos.IDSeparator + rscData.Name.ValueString())
}

func (rscData *securityZoneBookAddressData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityZoneBookAddressData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security zones security-zone " +
		rscData.Zone.ValueString() + " address-book address " + rscData.Name.ValueString() + " "

	if v := rscData.CIDR.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+v)
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := rscData.DNSName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"dns-name "+v)
		if rscData.DNSIPv4Only.ValueBool() {
			configSet = append(configSet, setPrefix+"dns-name "+v+" ipv4-only")
		}
		if rscData.DNSIPv6Only.ValueBool() {
			configSet = append(configSet, setPrefix+"dns-name "+v+" ipv6-only")
		}
	}
	if v := rscData.RangeFrom.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"range-address "+v+" to "+rscData.RangeTo.ValueString())
	}
	if v := rscData.Wildcard.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"wildcard-address "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityZoneBookAddressData) read(
	_ context.Context, zone, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + zone + " address-book address " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.Zone = types.StringValue(zone)
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
			case balt.CutPrefixInString(&itemTrim, "dns-name "):
				switch {
				case balt.CutSuffixInString(&itemTrim, " ipv4-only"):
					rscData.DNSIPv4Only = types.BoolValue(true)
					rscData.DNSName = types.StringValue(itemTrim)
				case balt.CutSuffixInString(&itemTrim, " ipv6-only"):
					rscData.DNSIPv6Only = types.BoolValue(true)
					rscData.DNSName = types.StringValue(itemTrim)
				default:
					rscData.DNSName = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "range-address "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 3 { // <from> to <to>
					return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "range-address", itemTrim)
				}
				rscData.RangeFrom = types.StringValue(itemTrimFields[0])
				rscData.RangeTo = types.StringValue(itemTrimFields[2])
			case balt.CutPrefixInString(&itemTrim, "wildcard-address "):
				rscData.Wildcard = types.StringValue(itemTrim)
			case strings.Contains(itemTrim, "/"):
				rscData.CIDR = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (rscData *securityZoneBookAddressData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security zones security-zone " + rscData.Zone.ValueString() +
			" address-book address " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
