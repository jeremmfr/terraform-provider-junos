package provider

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &nullLoadConfig{}
	_ resource.ResourceWithConfigure      = &nullLoadConfig{}
	_ resource.ResourceWithValidateConfig = &nullLoadConfig{}
)

type nullLoadConfig struct {
	client *junos.Client
}

func newNullLoadConfigResource() resource.Resource {
	return &nullLoadConfig{}
}

func (rsc *nullLoadConfig) typeName() string {
	return providerName + "_null_load_config"
}

func (rsc *nullLoadConfig) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *nullLoadConfig) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *nullLoadConfig) Configure(
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

func (rsc *nullLoadConfig) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Load an arbitrary configuration and commit it.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with value " +
					"`null_load_config`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"config": schema.StringAttribute{
				Required:    true,
				Description: "The configuration to load and apply.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"action": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.LoadConfigActionMerge),
				Description: "Specify how to load the configuration data.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
				Computed:    true,
				Default:     stringdefault.StaticString(junos.LoadConfigFormatText),
				Description: "The format used for the configuration data.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						junos.LoadConfigFormatText,
						junos.LoadConfigFormatJSON,
						junos.LoadConfigFormatXML,
					),
				},
			},
			"triggers": schema.DynamicAttribute{
				Optional:    true,
				Description: "Any value that, when changed, will force the resource to be replaced.",
				PlanModifiers: []planmodifier.Dynamic{
					dynamicplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

type nullLoadConfigData struct {
	ID       types.String  `tfsdk:"id"`
	Config   types.String  `tfsdk:"config"`
	Action   types.String  `tfsdk:"action"`
	Format   types.String  `tfsdk:"format"`
	Triggers types.Dynamic `tfsdk:"triggers"`
}

func (rsc *nullLoadConfig) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config nullLoadConfigData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Format.IsNull() && !config.Format.IsUnknown() {
		format := config.Format.ValueString()
		if action := config.Action.ValueString(); action == junos.LoadConfigActionSet &&
			format != junos.LoadConfigFormatText {
			resp.Diagnostics.AddAttributeError(
				path.Root("format"),
				tfdiag.ConflictConfigErrSummary,
				fmt.Sprintf("format cannot be %q when action = %q, must be 'text'", format, action),
			)
		}
	}
}

func (rsc *nullLoadConfig) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan nullLoadConfigData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clt := rsc.junosClient()
	junSess, err := clt.StartNewSession(ctx)
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

	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}

	warns, err := junSess.CommitConf(ctx, "load a config with resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rsc *nullLoadConfig) Read(
	_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse,
) {
	// no-op
}

func (rsc *nullLoadConfig) Update(
	_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse,
) {
	// no-op
}

func (rsc *nullLoadConfig) Delete(
	_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse,
) {
	// no-op
}

func (rscData *nullLoadConfigData) fillID() {
	rscData.ID = types.StringValue("null_load_config")
}

func (rscData *nullLoadConfigData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	format := rscData.Format.ValueString()
	if format == "" {
		format = junos.LoadConfigFormatText
	}
	action := rscData.Action.ValueString()
	if action == "" {
		action = junos.LoadConfigActionMerge
	}

	if format != junos.LoadConfigFormatText && action == junos.LoadConfigActionSet {
		return path.Root("format"),
			fmt.Errorf("format cannot be %q when action = %q, must be %q", format, action, junos.LoadConfigFormatText)
	}

	return path.Empty(), junSess.ConfigLoad(action, format, rscData.Config.ValueString())
}
