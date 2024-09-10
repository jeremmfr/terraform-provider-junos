package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
	_ resource.Resource                   = &securityZoneBookAddressSet{}
	_ resource.ResourceWithConfigure      = &securityZoneBookAddressSet{}
	_ resource.ResourceWithValidateConfig = &securityZoneBookAddressSet{}
	_ resource.ResourceWithImportState    = &securityZoneBookAddressSet{}
)

type securityZoneBookAddressSet struct {
	client *junos.Client
}

func newSecurityZoneBookAddressSetResource() resource.Resource {
	return &securityZoneBookAddressSet{}
}

func (rsc *securityZoneBookAddressSet) typeName() string {
	return providerName + "_security_zone_book_address_set"
}

func (rsc *securityZoneBookAddressSet) junosName() string {
	return "security zone address-book address-set"
}

func (rsc *securityZoneBookAddressSet) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityZoneBookAddressSet) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityZoneBookAddressSet) Configure(
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

func (rsc *securityZoneBookAddressSet) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Provides an address-set resource in address-book of security zone.",
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
				Description: "The name of address-set.",
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
	}
}

type securityZoneBookAddressSetData struct {
	ID          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Zone        types.String   `tfsdk:"zone"`
	Description types.String   `tfsdk:"description"`
	Address     []types.String `tfsdk:"address"`
	AddressSet  []types.String `tfsdk:"address_set"`
}

type securityZoneBookAddressSetConfig struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Zone        types.String `tfsdk:"zone"`
	Description types.String `tfsdk:"description"`
	Address     types.Set    `tfsdk:"address"`
	AddressSet  types.Set    `tfsdk:"address_set"`
}

func (rsc *securityZoneBookAddressSet) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityZoneBookAddressSetConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.Address.IsNull() && config.AddressSet.IsNull() {
		resp.Diagnostics.AddError(
			tfdiag.MissingConfigErrSummary,
			"at least one of address or address_set must be specified",
		)
	}
}

func (rsc *securityZoneBookAddressSet) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityZoneBookAddressSetData
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
			setExists, err := checkSecurityZoneBookAddressSetExists(
				fnCtx,
				plan.Zone.ValueString(),
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if setExists {
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
			setExists, err := checkSecurityZoneBookAddressSetExists(
				fnCtx,
				plan.Zone.ValueString(),
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !setExists {
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

func (rsc *securityZoneBookAddressSet) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityZoneBookAddressSetData
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

func (rsc *securityZoneBookAddressSet) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityZoneBookAddressSetData
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

func (rsc *securityZoneBookAddressSet) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityZoneBookAddressSetData
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

func (rsc *securityZoneBookAddressSet) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityZoneBookAddressSetData

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

func checkSecurityZoneBookAddressSetExists(
	_ context.Context, zone, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + zone + " address-book address-set " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityZoneBookAddressSetData) fillID() {
	rscData.ID = types.StringValue(rscData.Zone.ValueString() + junos.IDSeparator + rscData.Name.ValueString())
}

func (rscData *securityZoneBookAddressSetData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityZoneBookAddressSetData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security zones security-zone " +
		rscData.Zone.ValueString() + " address-book address-set " + rscData.Name.ValueString() + " "

	if len(rscData.Address) == 0 && len(rscData.AddressSet) == 0 {
		return path.Empty(), errors.New("at least one element of address or address_set must be specified")
	}
	for _, v := range rscData.Address {
		configSet = append(configSet, setPrefix+"address "+v.ValueString())
	}
	for _, v := range rscData.AddressSet {
		configSet = append(configSet, setPrefix+"address-set "+v.ValueString())
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityZoneBookAddressSetData) read(
	_ context.Context, zone, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + zone + " address-book address-set " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "address "):
				rscData.Address = append(rscData.Address, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "address-set "):
				rscData.AddressSet = append(rscData.AddressSet, types.StringValue(itemTrim))
			}
		}
	}

	return nil
}

func (rscData *securityZoneBookAddressSetData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security zones security-zone " + rscData.Zone.ValueString() +
			" address-book address-set " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
