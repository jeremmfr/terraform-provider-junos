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
				Description: "An identifier for the resource with format `<zone>_-_<name>`.",
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
			"Missing Configuration Error",
			"cannot have dns_ipv4_only without dns_name",
		)
	}
	if !config.DNSIPv6Only.IsNull() && config.DNSName.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dns_ipv6_only"),
			"Missing Configuration Error",
			"cannot have dns_ipv6_only without dns_name",
		)
	}
	if !config.DNSIPv4Only.IsNull() && !config.DNSIPv6Only.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("dns_ipv4_only"),
			"Conflict Configuration Error",
			"only one of dns_ipv4_only or dns_ipv6_only can be specified",
		)
	}
	if !config.RangeTo.IsNull() && config.RangeFrom.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("range_to"),
			"Missing Configuration Error",
			"cannot have range_to without range_from",
		)
	}
	if !config.RangeFrom.IsNull() && config.RangeTo.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("range_from"),
			"Missing Configuration Error",
			"cannot have range_from without range_to",
		)
	}
	switch {
	case !config.CIDR.IsNull():
		if !config.DNSName.IsNull() ||
			!config.RangeFrom.IsNull() ||
			!config.Wildcard.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("cidr"),
				"Conflict Configuration Error",
				"only one of cidr, dns_name, range_from or wildcard must be specified",
			)
		}
	case !config.DNSName.IsNull():
		if !config.CIDR.IsNull() ||
			!config.RangeFrom.IsNull() ||
			!config.Wildcard.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("dns_name"),
				"Conflict Configuration Error",
				"only one of cidr, dns_name, range_from or wildcard must be specified",
			)
		}
	case !config.RangeFrom.IsNull():
		if !config.CIDR.IsNull() ||
			!config.DNSName.IsNull() ||
			!config.Wildcard.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("range_from"),
				"Conflict Configuration Error",
				"only one of cidr, dns_name, range_from or wildcard must be specified",
			)
		}
	case !config.Wildcard.IsNull():
		if !config.CIDR.IsNull() ||
			!config.DNSName.IsNull() ||
			!config.RangeFrom.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("range_from"),
				"Conflict Configuration Error",
				"only one of cidr, dns_name, range_from or wildcard must be specified",
			)
		}
	default:
		resp.Diagnostics.AddError(
			"Missing Configuration Error",
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
			"could not create "+rsc.junosName()+" address with empty name",
		)

		return
	}
	if plan.Zone.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("zone"),
			"Empty Zone",
			"could not create "+rsc.junosName()+" address with empty zone",
		)

		return
	}

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		resp.Diagnostics.AddError(
			"Compatibility Error",
			fmt.Sprintf(rsc.junosName()+" not compatible "+
				"with Junos device %q", junSess.SystemInformation.HardwareModel),
		)

		return
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}
	defer func() { resp.Diagnostics.Append(tfdiag.Warns("Config Clear/Unlock Warning", junSess.ConfigClear())...) }()

	zonesExists, err := checkSecurityZonesExists(ctx, plan.Zone.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.AddError("Pre Check Error", err.Error())

		return
	}
	if !zonesExists {
		resp.Diagnostics.AddAttributeError(
			path.Root("zone"),
			"Missing Configuration Error",
			fmt.Sprintf("security zone %q doesn't exist", plan.Zone.ValueString()),
		)

		return
	}
	addressExists, err := checkSecurityZoneBookAddressExists(
		ctx,
		plan.Zone.ValueString(),
		plan.Name.ValueString(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError("Pre Check Error", err.Error())

		return
	}
	if addressExists {
		resp.Diagnostics.AddError(
			"Duplicate Configuration Error",
			fmt.Sprintf(rsc.junosName()+" %q already exists in zone %q",
				plan.Name.ValueString(), plan.Zone.ValueString()),
		)

		return
	}

	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("create resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	addressExists, err = checkSecurityZoneBookAddressExists(
		ctx,
		plan.Zone.ValueString(),
		plan.Name.ValueString(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError("Post Check Error", err.Error())

		return
	}
	if !addressExists {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf(rsc.junosName()+" %q does not exists in zone %q after commit "+
				"=> check your config", plan.Name.ValueString(), plan.Zone.ValueString()),
		)

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityZoneBookAddress) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityZoneBookAddressData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	err = data.read(
		ctx,
		state.Zone.ValueString(),
		state.Name.ValueString(),
		junSess,
	)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}
	if data.ID.IsNull() {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}
	defer func() { resp.Diagnostics.Append(tfdiag.Warns("Config Clear/Unlock Warning", junSess.ConfigClear())...) }()

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("update resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityZoneBookAddress) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityZoneBookAddressData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeDeleteAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}
	defer func() { resp.Diagnostics.Append(tfdiag.Warns("Config Clear/Unlock Warning", junSess.ConfigClear())...) }()

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	warns, err := junSess.CommitConf("delete resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}
}

func (rsc *securityZoneBookAddress) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	idList := strings.Split(req.ID, junos.IDSeparator)
	if len(idList) < 2 {
		resp.Diagnostics.AddError(
			"Bad ID Format",
			fmt.Sprintf("missing element(s) in id with separator %q", junos.IDSeparator),
		)

		return
	}

	var data securityZoneBookAddressData
	if err := data.read(ctx, idList[0], idList[1], junSess); err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf("don't find "+rsc.junosName()+" with id %q "+
				"(id must be <zone>"+junos.IDSeparator+"<name>)", req.ID),
		)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
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
) (
	err error,
) {
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
