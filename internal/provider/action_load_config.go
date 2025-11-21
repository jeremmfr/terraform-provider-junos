package provider

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ action.Action                   = &loadConfigAction{}
	_ action.ActionWithConfigure      = &loadConfigAction{}
	_ action.ActionWithValidateConfig = &loadConfigAction{}
)

type loadConfigAction struct {
	client *junos.Client
}

func newLoadConfigAction() action.Action {
	return &loadConfigAction{}
}

func (act *loadConfigAction) typeName() string {
	return providerName + "_load_config"
}

func (act *loadConfigAction) junosClient() *junos.Client {
	return act.client
}

func (act *loadConfigAction) Metadata(
	_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_load_config"
}

func (act *loadConfigAction) Configure(
	ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*junos.Client)
	if !ok {
		unexpectedActionConfigureType(ctx, req, resp)

		return
	}
	act.client = client
}

func (act *loadConfigAction) Schema(
	_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Load an arbitrary configuration and commit it.",
		Attributes: map[string]schema.Attribute{
			"config": schema.StringAttribute{
				Required:    true,
				Description: "The configuration to load and apply.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"action": schema.StringAttribute{
				Optional:    true,
				Description: "Specify how to load the configuration data. Defaults to 'merge'.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						junos.LoadConfigActionMerge,
						junos.LoadConfigActionOverride,
						junos.LoadConfigActionReplace,
						junos.LoadConfigActionSet,
						junos.LoadConfigActionUpdate,
					),
				},
			},
			"format": schema.StringAttribute{
				Optional:    true,
				Description: "The format used for the configuration data. Defaults to 'text'.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						junos.LoadConfigFormatText,
						junos.LoadConfigFormatJSON,
						junos.LoadConfigFormatXML,
					),
				},
			},
		},
	}
}

type loadConfigActionData struct {
	Config types.String `tfsdk:"config"`
	Action types.String `tfsdk:"action"`
	Format types.String `tfsdk:"format"`
}

func (act *loadConfigAction) ValidateConfig(
	ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse,
) {
	var config loadConfigActionData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Format.IsNull() && !config.Format.IsUnknown() {
		format := config.Format.ValueString()
		actionValue := config.Action.ValueString()
		if actionValue == "" {
			actionValue = junos.LoadConfigActionMerge
		}
		if actionValue == junos.LoadConfigActionSet && format != junos.LoadConfigFormatText {
			resp.Diagnostics.AddAttributeError(
				path.Root("format"),
				tfdiag.ConflictConfigErrSummary,
				fmt.Sprintf("format cannot be %q when action = %q, must be 'text'", format, actionValue),
			)
		}
	}
}

func (act *loadConfigAction) Invoke(
	ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse,
) {
	var config loadConfigActionData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	format := config.Format.ValueString()
	if format == "" {
		format = junos.LoadConfigFormatText
	}
	actionValue := config.Action.ValueString()
	if actionValue == "" {
		actionValue = junos.LoadConfigActionMerge
	}

	if format != junos.LoadConfigFormatText && actionValue == junos.LoadConfigActionSet {
		resp.Diagnostics.AddAttributeError(
			path.Root("format"),
			tfdiag.ConflictConfigErrSummary,
			fmt.Sprintf("format cannot be %q when action = %q, must be %q", format, actionValue, junos.LoadConfigFormatText),
		)

		return
	}

	clt := act.junosClient()
	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Starting session to device",
	})
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Locking candidate configuration",
	})
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigUnlockWarnSummary, junSess.ConfigUnlock())...)
	}()

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Loading configuration",
	})
	if err := junSess.ConfigLoad(actionValue, format, config.Config.ValueString()); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())

		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Committing configuration",
	})
	warns, err := junSess.CommitConf(ctx, "load a config with action "+act.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: "Configuration loaded and committed",
	})
}
