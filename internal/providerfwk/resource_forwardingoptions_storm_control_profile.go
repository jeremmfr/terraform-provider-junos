package providerfwk

import (
	"context"
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
	_ resource.Resource                   = &forwardingoptionsStormControlProfile{}
	_ resource.ResourceWithConfigure      = &forwardingoptionsStormControlProfile{}
	_ resource.ResourceWithValidateConfig = &forwardingoptionsStormControlProfile{}
	_ resource.ResourceWithImportState    = &forwardingoptionsStormControlProfile{}
)

type forwardingoptionsStormControlProfile struct {
	client *junos.Client
}

func newForwardingoptionsStormControlProfileResource() resource.Resource {
	return &forwardingoptionsStormControlProfile{}
}

func (rsc *forwardingoptionsStormControlProfile) typeName() string {
	return providerName + "_forwardingoptions_storm_control_profile"
}

func (rsc *forwardingoptionsStormControlProfile) junosName() string {
	return "forwarding-options storm-control-profile"
}

func (rsc *forwardingoptionsStormControlProfile) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *forwardingoptionsStormControlProfile) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *forwardingoptionsStormControlProfile) Configure(
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

func (rsc *forwardingoptionsStormControlProfile) Schema(
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
				Description: "Storm control profile name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 127),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"action_shutdown": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable port for excessive storm control errors.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"all": schema.SingleNestedBlock{
				Description: "For all BUM traffic.",
				Attributes: map[string]schema.Attribute{
					"bandwidth_level": schema.Int64Attribute{
						Optional:    true,
						Description: "Link bandwidth (kbps)",
						Validators: []validator.Int64{
							int64validator.Between(100, 100000000),
						},
					},
					"bandwidth_percentage": schema.Int64Attribute{
						Optional:    true,
						Description: "Percentage of link bandwidth.",
						Validators: []validator.Int64{
							int64validator.Between(1, 100),
						},
					},
					"burst_size": schema.Int64Attribute{
						Optional:    true,
						Description: "Burst size (bytes).",
						Validators: []validator.Int64{
							int64validator.Between(1500, 100000000),
						},
					},
					"no_broadcast": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable broadcast storm control.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_multicast": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable multicast storm control.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_registered_multicast": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable registered multicast storm control.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_unknown_unicast": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable unknown unicast storm control.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"no_unregistered_multicast": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable unregistered multicast storm control.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
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

type forwardingoptionsStormControlProfileData struct {
	ID             types.String                                  `tfsdk:"id"`
	Name           types.String                                  `tfsdk:"name"`
	ActionShutdown types.Bool                                    `tfsdk:"action_shutdown"`
	All            *forwardingoptionsStormControlProfileBlockAll `tfsdk:"all"`
}

type forwardingoptionsStormControlProfileBlockAll struct {
	BandwidthLevel          types.Int64 `tfsdk:"bandwidth_level"`
	BandwidthPercentage     types.Int64 `tfsdk:"bandwidth_percentage"`
	BurstSize               types.Int64 `tfsdk:"burst_size"`
	NoBroadcast             types.Bool  `tfsdk:"no_broadcast"`
	NoMulticast             types.Bool  `tfsdk:"no_multicast"`
	NoRegisteredMulticast   types.Bool  `tfsdk:"no_registered_multicast"`
	NoUnknownUnicast        types.Bool  `tfsdk:"no_unknown_unicast"`
	NoUnregisteredMulticast types.Bool  `tfsdk:"no_unregistered_multicast"`
}

func (rsc *forwardingoptionsStormControlProfile) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config forwardingoptionsStormControlProfileData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.All == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("all"),
			tfdiag.MissingConfigErrSummary,
			"all block must be specified",
		)
	} else {
		if !config.All.BandwidthLevel.IsNull() && !config.All.BandwidthLevel.IsUnknown() &&
			!config.All.BandwidthPercentage.IsNull() && !config.All.BandwidthPercentage.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("all").AtName("bandwidth_level"),
				tfdiag.ConflictConfigErrSummary,
				"bandwidth_level and bandwidth_percentage cannot be configured together"+
					" in all block",
			)
		}
		if !config.All.NoMulticast.IsNull() && !config.All.NoMulticast.IsUnknown() {
			if !config.All.NoRegisteredMulticast.IsNull() && !config.All.NoRegisteredMulticast.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("all").AtName("no_multicast"),
					tfdiag.ConflictConfigErrSummary,
					"no_multicast and no_registered_multicast cannot be configured together"+
						" in all block",
				)
			}
			if !config.All.NoUnregisteredMulticast.IsNull() && !config.All.NoUnregisteredMulticast.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("all").AtName("no_multicast"),
					tfdiag.ConflictConfigErrSummary,
					"no_multicast and no_unregistered_multicast cannot be configured together"+
						" in all block",
				)
			}
		}
		if !config.All.NoRegisteredMulticast.IsNull() && !config.All.NoRegisteredMulticast.IsUnknown() &&
			!config.All.NoUnregisteredMulticast.IsNull() && !config.All.NoUnregisteredMulticast.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("all").AtName("no_registered_multicast"),
				tfdiag.ConflictConfigErrSummary,
				"no_registered_multicast and no_unregistered_multicast cannot be configured together"+
					" in all block",
			)
		}
	}
}

func (rsc *forwardingoptionsStormControlProfile) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan forwardingoptionsStormControlProfileData
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
			profileExists, err := checkForwardingoptionsStormControlProfileExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if profileExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			profileExists, err := checkForwardingoptionsStormControlProfileExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !profileExists {
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

func (rsc *forwardingoptionsStormControlProfile) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data forwardingoptionsStormControlProfileData
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

func (rsc *forwardingoptionsStormControlProfile) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state forwardingoptionsStormControlProfileData
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

func (rsc *forwardingoptionsStormControlProfile) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state forwardingoptionsStormControlProfileData
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

func (rsc *forwardingoptionsStormControlProfile) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data forwardingoptionsStormControlProfileData

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

func checkForwardingoptionsStormControlProfileExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"forwarding-options storm-control-profiles \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *forwardingoptionsStormControlProfileData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *forwardingoptionsStormControlProfileData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *forwardingoptionsStormControlProfileData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set forwarding-options storm-control-profiles \"" + rscData.Name.ValueString() + "\" "

	if rscData.ActionShutdown.ValueBool() {
		configSet = append(configSet, setPrefix+"action-shutdown")
	}
	if rscData.All != nil {
		configSet = append(configSet, setPrefix+"all")

		if !rscData.All.BandwidthLevel.IsNull() {
			configSet = append(configSet, setPrefix+"all bandwidth-level "+
				utils.ConvI64toa(rscData.All.BandwidthLevel.ValueInt64()))
		}
		if !rscData.All.BandwidthPercentage.IsNull() {
			configSet = append(configSet, setPrefix+"all bandwidth-percentage "+
				utils.ConvI64toa(rscData.All.BandwidthPercentage.ValueInt64()))
		}
		if !rscData.All.BurstSize.IsNull() {
			configSet = append(configSet, setPrefix+"all burst-size "+
				utils.ConvI64toa(rscData.All.BurstSize.ValueInt64()))
		}
		if rscData.All.NoBroadcast.ValueBool() {
			configSet = append(configSet, setPrefix+"all no-broadcast")
		}
		if rscData.All.NoMulticast.ValueBool() {
			configSet = append(configSet, setPrefix+"all no-multicast")
		}
		if rscData.All.NoRegisteredMulticast.ValueBool() {
			configSet = append(configSet, setPrefix+"all no-registered-multicast")
		}
		if rscData.All.NoUnknownUnicast.ValueBool() {
			configSet = append(configSet, setPrefix+"all no-unknown-unicast")
		}
		if rscData.All.NoUnregisteredMulticast.ValueBool() {
			configSet = append(configSet, setPrefix+"all no-unregistered-multicast")
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *forwardingoptionsStormControlProfileData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"forwarding-options storm-control-profiles \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case itemTrim == "action-shutdown":
				rscData.ActionShutdown = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "all"):
				if rscData.All == nil {
					rscData.All = &forwardingoptionsStormControlProfileBlockAll{}
				}
				var err error
				switch {
				case balt.CutPrefixInString(&itemTrim, " bandwidth-level "):
					rscData.All.BandwidthLevel, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " bandwidth-percentage "):
					rscData.All.BandwidthPercentage, err = tfdata.ConvAtoi64Value(itemTrim)
				case balt.CutPrefixInString(&itemTrim, " burst-size "):
					rscData.All.BurstSize, err = tfdata.ConvAtoi64Value(itemTrim)
				case itemTrim == " no-broadcast":
					rscData.All.NoBroadcast = types.BoolValue(true)
				case itemTrim == " no-multicast":
					rscData.All.NoMulticast = types.BoolValue(true)
				case itemTrim == " no-registered-multicast":
					rscData.All.NoRegisteredMulticast = types.BoolValue(true)
				case itemTrim == " no-unknown-unicast":
					rscData.All.NoUnknownUnicast = types.BoolValue(true)
				case itemTrim == " no-unregistered-multicast":
					rscData.All.NoUnregisteredMulticast = types.BoolValue(true)
				}
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *forwardingoptionsStormControlProfileData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete forwarding-options storm-control-profiles \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
