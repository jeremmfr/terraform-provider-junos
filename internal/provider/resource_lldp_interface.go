package provider

import (
	"context"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
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
	_ resource.Resource                   = &lldpInterface{}
	_ resource.ResourceWithConfigure      = &lldpInterface{}
	_ resource.ResourceWithValidateConfig = &lldpInterface{}
	_ resource.ResourceWithImportState    = &lldpInterface{}
	_ resource.ResourceWithUpgradeState   = &lldpInterface{}
)

type lldpInterface struct {
	client *junos.Client
}

func newLldpInterfaceResource() resource.Resource {
	return &lldpInterface{}
}

func (rsc *lldpInterface) typeName() string {
	return providerName + "_lldp_interface"
}

func (rsc *lldpInterface) junosName() string {
	return "lldp interface"
}

func (rsc *lldpInterface) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *lldpInterface) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *lldpInterface) Configure(
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

func (rsc *lldpInterface) Schema(
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
				Description: "Disable LLDP.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"enable": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable LLDP.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"trap_notification_disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable lldp-trap notification.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"trap_notification_enable": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable lldp-trap notification.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"power_negotiation": schema.SingleNestedBlock{
				Description: "LLDP power negotiation.",
				Attributes: map[string]schema.Attribute{
					"disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable power negotiation.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"enable": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable power negotiation.",
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

type lldpInterfaceData struct {
	ID                      types.String                        `tfsdk:"id"`
	Name                    types.String                        `tfsdk:"name"`
	Disable                 types.Bool                          `tfsdk:"disable"`
	Enable                  types.Bool                          `tfsdk:"enable"`
	TrapNotificationDisable types.Bool                          `tfsdk:"trap_notification_disable"`
	TrapNotificationEnable  types.Bool                          `tfsdk:"trap_notification_enable"`
	PowerNegotiation        *lldpInterfaceBlockPowerNegotiation `tfsdk:"power_negotiation"`
}

type lldpInterfaceBlockPowerNegotiation struct {
	Disable types.Bool `tfsdk:"disable"`
	Enable  types.Bool `tfsdk:"enable"`
}

func (rsc *lldpInterface) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config lldpInterfaceData
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
	if !config.TrapNotificationDisable.IsNull() && !config.TrapNotificationDisable.IsUnknown() &&
		!config.TrapNotificationEnable.IsNull() && !config.TrapNotificationEnable.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("trap_notification_enable"),
			tfdiag.ConflictConfigErrSummary,
			"trap_notification_enable and trap_notification_disable cannot be configured together",
		)
	}
	if config.PowerNegotiation != nil {
		if !config.PowerNegotiation.Disable.IsNull() && !config.PowerNegotiation.Disable.IsUnknown() &&
			!config.PowerNegotiation.Enable.IsNull() && !config.PowerNegotiation.Enable.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("power_negotiation").AtName("enable"),
				tfdiag.ConflictConfigErrSummary,
				"enable and disable cannot be configured together"+
					" in power_negotiation block",
			)
		}
	}
}

func (rsc *lldpInterface) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan lldpInterfaceData
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
			interfaceExists, err := checkLldpInterfaceExists(
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
			interfaceExists, err := checkLldpInterfaceExists(
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

func (rsc *lldpInterface) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data lldpInterfaceData
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

func (rsc *lldpInterface) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state lldpInterfaceData
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

func (rsc *lldpInterface) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state lldpInterfaceData
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

func (rsc *lldpInterface) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data lldpInterfaceData

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

func checkLldpInterfaceExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols lldp interface " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *lldpInterfaceData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *lldpInterfaceData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *lldpInterfaceData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set protocols lldp interface " + rscData.Name.ValueString() + " "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if rscData.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if rscData.Enable.ValueBool() {
		configSet = append(configSet, setPrefix+"enable")
	}
	if rscData.TrapNotificationDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"trap-notification disable")
	}
	if rscData.TrapNotificationEnable.ValueBool() {
		configSet = append(configSet, setPrefix+"trap-notification enable")
	}

	if rscData.PowerNegotiation != nil {
		configSet = append(configSet, setPrefix+"power-negotiation")

		if rscData.PowerNegotiation.Disable.ValueBool() {
			configSet = append(configSet, setPrefix+"power-negotiation disable")
		}
		if rscData.PowerNegotiation.Enable.ValueBool() {
			configSet = append(configSet, setPrefix+"power-negotiation enable")
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *lldpInterfaceData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols lldp interface " + name + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.fillID()
		for item := range strings.SplitSeq(showConfig, "\n") {
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
			case balt.CutPrefixInString(&itemTrim, "power-negotiation"):
				if rscData.PowerNegotiation == nil {
					rscData.PowerNegotiation = &lldpInterfaceBlockPowerNegotiation{}
				}

				switch {
				case itemTrim == " disable":
					rscData.PowerNegotiation.Disable = types.BoolValue(true)
				case itemTrim == " enable":
					rscData.PowerNegotiation.Enable = types.BoolValue(true)
				}
			case itemTrim == "trap-notification disable":
				rscData.TrapNotificationDisable = types.BoolValue(true)
			case itemTrim == "trap-notification enable":
				rscData.TrapNotificationEnable = types.BoolValue(true)
			}
		}
	}

	return nil
}

func (rscData *lldpInterfaceData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete protocols lldp interface " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
