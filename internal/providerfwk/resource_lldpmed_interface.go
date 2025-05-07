package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource                   = &lldpMedInterface{}
	_ resource.ResourceWithConfigure      = &lldpMedInterface{}
	_ resource.ResourceWithValidateConfig = &lldpMedInterface{}
	_ resource.ResourceWithImportState    = &lldpMedInterface{}
	_ resource.ResourceWithUpgradeState   = &lldpMedInterface{}
)

type lldpMedInterface struct {
	client *junos.Client
}

func newLldpMedInterfaceResource() resource.Resource {
	return &lldpMedInterface{}
}

func (rsc *lldpMedInterface) typeName() string {
	return providerName + "_lldpmed_interface"
}

func (rsc *lldpMedInterface) junosName() string {
	return "lldp-med interface"
}

func (rsc *lldpMedInterface) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *lldpMedInterface) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *lldpMedInterface) Configure(
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

func (rsc *lldpMedInterface) Schema(
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
				Description: "Interface name or `all`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					tfvalidator.StringDotExclusion(),
				},
			},
			"disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable LLDP-MED.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"enable": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable LLDP-MED.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"location": schema.SingleNestedBlock{
				Description: "Define location.",
				Attributes: map[string]schema.Attribute{
					"civic_based_country_code": schema.StringAttribute{
						Optional:    true,
						Description: "Two-letter country code.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(2, 2),
						},
					},
					"civic_based_what": schema.Int64Attribute{
						Optional:    true,
						Description: "Type of address.",
						Validators: []validator.Int64{
							int64validator.Between(0, 2),
						},
					},
					"co_ordinate_latitude": schema.Int64Attribute{
						Optional:    true,
						Description: "Latitude value to address based on longitude and latitude coordinates.",
						Validators: []validator.Int64{
							int64validator.Between(0, 360),
						},
					},
					"co_ordinate_longitude": schema.Int64Attribute{
						Optional:    true,
						Description: "Longitude value to address based on longitude and latitude coordinates.",
						Validators: []validator.Int64{
							int64validator.Between(0, 360),
						},
					},
					"elin": schema.StringAttribute{
						Optional:    true,
						Description: "Emergency line identification (ELIN) string.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"civic_based_ca_type": schema.ListNestedBlock{
						Description: "For each ca-type, configure civic-based ca-type.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ca_type": schema.Int64Attribute{
									Required:    true,
									Description: "Address element type.",
									Validators: []validator.Int64{
										int64validator.Between(0, 255),
									},
								},
								"ca_value": schema.StringAttribute{
									Optional:    true,
									Description: "Address element value.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
							},
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

type lldpMedInterfaceData struct {
	ID       types.String                   `tfsdk:"id"`
	Name     types.String                   `tfsdk:"name"`
	Disable  types.Bool                     `tfsdk:"disable"`
	Enable   types.Bool                     `tfsdk:"enable"`
	Location *lldpMedInterfaceBlockLocation `tfsdk:"location"`
}

type lldpMedInterfaceConfig struct {
	ID       types.String                         `tfsdk:"id"`
	Name     types.String                         `tfsdk:"name"`
	Disable  types.Bool                           `tfsdk:"disable"`
	Enable   types.Bool                           `tfsdk:"enable"`
	Location *lldpMedInterfaceBlockLocationConfig `tfsdk:"location"`
}

type lldpMedInterfaceBlockLocation struct {
	CivicBasedCountryCode types.String                                         `tfsdk:"civic_based_country_code"`
	CivicBasedWhat        types.Int64                                          `tfsdk:"civic_based_what"`
	CoOrdinateLatitude    types.Int64                                          `tfsdk:"co_ordinate_latitude"`
	CoOrdinateLongitude   types.Int64                                          `tfsdk:"co_ordinate_longitude"`
	Elin                  types.String                                         `tfsdk:"elin"`
	CivicBasedCaType      []lldpMedInterfaceBlockLocationBlockCivicBasedCaType `tfsdk:"civic_based_ca_type"`
}

type lldpMedInterfaceBlockLocationConfig struct {
	CivicBasedCountryCode types.String `tfsdk:"civic_based_country_code"`
	CivicBasedWhat        types.Int64  `tfsdk:"civic_based_what"`
	CoOrdinateLatitude    types.Int64  `tfsdk:"co_ordinate_latitude"`
	CoOrdinateLongitude   types.Int64  `tfsdk:"co_ordinate_longitude"`
	Elin                  types.String `tfsdk:"elin"`
	CivicBasedCaType      types.List   `tfsdk:"civic_based_ca_type"`
}

type lldpMedInterfaceBlockLocationBlockCivicBasedCaType struct {
	CaType  types.Int64  `tfsdk:"ca_type"  tfdata:"identifier"`
	CaValue types.String `tfsdk:"ca_value"`
}

func (rsc *lldpMedInterface) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config lldpMedInterfaceConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Disable.IsNull() && !config.Disable.IsUnknown() &&
		!config.Enable.IsNull() && !config.Enable.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("enable"),
			tfdiag.ConflictConfigErrSummary,
			"enable and disable cannot be configured together",
		)
	}

	if config.Location != nil {
		if !config.Location.CivicBasedCountryCode.IsNull() &&
			!config.Location.CivicBasedCountryCode.IsUnknown() {
			if !config.Location.CoOrdinateLatitude.IsNull() &&
				!config.Location.CoOrdinateLatitude.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("civic_based_country_code"),
					tfdiag.ConflictConfigErrSummary,
					"civic_based_country_code and co_ordinate_latitude cannot be configured together"+
						" in location block",
				)
			}
			if !config.Location.CoOrdinateLongitude.IsNull() &&
				!config.Location.CoOrdinateLongitude.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("civic_based_country_code"),
					tfdiag.ConflictConfigErrSummary,
					"civic_based_country_code and co_ordinate_longitude cannot be configured together"+
						" in location block",
				)
			}
			if !config.Location.Elin.IsNull() &&
				!config.Location.Elin.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("civic_based_country_code"),
					tfdiag.ConflictConfigErrSummary,
					"civic_based_country_code and elin cannot be configured together"+
						" in location block",
				)
			}
		}
		if !config.Location.CivicBasedWhat.IsNull() &&
			!config.Location.CivicBasedWhat.IsUnknown() &&
			config.Location.CivicBasedCountryCode.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("location").AtName("civic_based_what"),
				tfdiag.MissingConfigErrSummary,
				"civic_based_country_code must be specified with civic_based_what"+
					" in location block",
			)
		}
		if !config.Location.Elin.IsNull() &&
			!config.Location.Elin.IsUnknown() {
			if !config.Location.CoOrdinateLatitude.IsNull() &&
				!config.Location.CoOrdinateLatitude.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("elin"),
					tfdiag.ConflictConfigErrSummary,
					"co_ordinate_latitude and elin cannot be configured together"+
						" in location block",
				)
			}
			if !config.Location.CoOrdinateLongitude.IsNull() &&
				!config.Location.CoOrdinateLongitude.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("elin"),
					tfdiag.ConflictConfigErrSummary,
					"co_ordinate_longitude and elin cannot be configured together"+
						" in location block",
				)
			}
		}
		if !config.Location.CivicBasedCaType.IsNull() &&
			!config.Location.CivicBasedCaType.IsUnknown() {
			if config.Location.CivicBasedCountryCode.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("location").AtName("civic_based_ca_type"),
					tfdiag.MissingConfigErrSummary,
					"civic_based_country_code must be specified with civic_based_ca_type"+
						" in location block",
				)
			}

			var configCivicBasedCaType []lldpMedInterfaceBlockLocationBlockCivicBasedCaType
			asDiags := config.Location.CivicBasedCaType.ElementsAs(ctx, &configCivicBasedCaType, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			civicBasedCaTypeCaType := make(map[int64]struct{})
			for i, block := range configCivicBasedCaType {
				if block.CaType.IsUnknown() {
					continue
				}

				caType := block.CaType.ValueInt64()
				if _, ok := civicBasedCaTypeCaType[caType]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("location").AtName("civic_based_ca_type").AtListIndex(i).AtName("ca_type"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple civic_based_ca_type blocks with the same ca_type %d"+
							" in location block", caType),
					)
				}
				civicBasedCaTypeCaType[caType] = struct{}{}
			}
		}
	}
}

func (rsc *lldpMedInterface) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan lldpMedInterfaceData
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
			interfaceExists, err := checkLldpMedInterfaceExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if interfaceExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			interfaceExists, err := checkLldpMedInterfaceExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !interfaceExists {
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

func (rsc *lldpMedInterface) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data lldpMedInterfaceData
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

func (rsc *lldpMedInterface) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state lldpMedInterfaceData
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

func (rsc *lldpMedInterface) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state lldpMedInterfaceData
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

func (rsc *lldpMedInterface) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data lldpMedInterfaceData

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

func checkLldpMedInterfaceExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols lldp-med interface " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *lldpMedInterfaceData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *lldpMedInterfaceData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *lldpMedInterfaceData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set protocols lldp-med interface " + rscData.Name.ValueString() + " "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if rscData.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if rscData.Enable.ValueBool() {
		configSet = append(configSet, setPrefix+"enable")
	}
	if rscData.Location != nil {
		blockSet, pathErr, err := rscData.Location.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *lldpMedInterfaceBlockLocation) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "location "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if v := block.CivicBasedCountryCode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"civic-based country-code "+v)

		if !block.CivicBasedWhat.IsNull() {
			configSet = append(configSet, setPrefix+"civic-based what "+
				utils.ConvI64toa(block.CivicBasedWhat.ValueInt64()))
		}

		civicBasedCaTypeCaType := make(map[int64]struct{})
		for i, subBlock := range block.CivicBasedCaType {
			caType := subBlock.CaType.ValueInt64()
			if _, ok := civicBasedCaTypeCaType[caType]; ok {
				return configSet,
					path.Root("location").AtName("civic_based_ca_type").AtListIndex(i).AtName("ca_type"),
					fmt.Errorf("multiple civic_based_ca_type blocks with the same ca_type %d"+
						" in location block", caType)
			}
			civicBasedCaTypeCaType[caType] = struct{}{}

			configSet = append(configSet, setPrefix+"civic-based ca-type "+utils.ConvI64toa(caType))
			if v := subBlock.CaValue.ValueString(); v != "" {
				configSet = append(configSet, setPrefix+"civic-based ca-type "+utils.ConvI64toa(caType)+" ca-value \""+v+"\"")
			}
		}
	} else {
		if len(block.CivicBasedCaType) > 0 {
			return configSet,
				path.Root("location").AtName("civic_based_ca_type"),
				errors.New("civic_based_country_code must be specified with " +
					"civic_based_ca_type in location block")
		}
		if !block.CivicBasedWhat.IsNull() {
			return configSet,
				path.Root("location").AtName("civic_based_what"),
				errors.New("civic_based_country_code must be specified with " +
					"civic_based_what in location block")
		}
	}

	if !block.CoOrdinateLatitude.IsNull() {
		configSet = append(configSet, setPrefix+"co-ordinate lattitude "+ //nolint:misspell
			utils.ConvI64toa(block.CoOrdinateLatitude.ValueInt64()))
		configSet = append(configSet, setPrefix+"co-ordinate latitude "+
			utils.ConvI64toa(block.CoOrdinateLatitude.ValueInt64()))
	}
	if !block.CoOrdinateLongitude.IsNull() {
		configSet = append(configSet, setPrefix+"co-ordinate longitude "+
			utils.ConvI64toa(block.CoOrdinateLongitude.ValueInt64()))
	}
	if v := block.Elin.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"elin \""+v+"\"")
	}

	return configSet, path.Empty(), nil
}

func (rscData *lldpMedInterfaceData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols lldp-med interface " + name + junos.PipeDisplaySetRelative)
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
			case itemTrim == junos.DisableW:
				rscData.Disable = types.BoolValue(true)
			case itemTrim == "enable":
				rscData.Enable = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "location"):
				if rscData.Location == nil {
					rscData.Location = &lldpMedInterfaceBlockLocation{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.Location.read(itemTrim); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (block *lldpMedInterfaceBlockLocation) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "civic-based ca-type "):
		caTypeStr := tfdata.FirstElementOfJunosLine(itemTrim)
		caType, err := tfdata.ConvAtoi64Value(caTypeStr)
		if err != nil {
			return err
		}
		block.CivicBasedCaType = tfdata.AppendPotentialNewBlock(block.CivicBasedCaType, caType)

		if balt.CutPrefixInString(&itemTrim, caTypeStr+" ca-value ") {
			block.CivicBasedCaType[len(block.CivicBasedCaType)-1].
				CaValue = types.StringValue(strings.Trim(itemTrim, "\""))
		}
	case balt.CutPrefixInString(&itemTrim, "civic-based country-code "):
		block.CivicBasedCountryCode = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "civic-based what "):
		block.CivicBasedWhat, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "co-ordinate lattitude "): //nolint:misspell
		block.CoOrdinateLatitude, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "co-ordinate latitude "):
		block.CoOrdinateLatitude, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "co-ordinate longitude "):
		block.CoOrdinateLongitude, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "elin "):
		block.Elin = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return err
}

func (rscData *lldpMedInterfaceData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete protocols lldp-med interface " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
