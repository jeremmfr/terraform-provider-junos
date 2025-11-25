package provider

import (
	"context"
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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &applyGroupExcept{}
	_ resource.ResourceWithConfigure      = &applyGroupExcept{}
	_ resource.ResourceWithValidateConfig = &applyGroupExcept{}
	_ resource.ResourceWithImportState    = &applyGroupExcept{}
)

type applyGroupExcept struct {
	client *junos.Client
}

func newApplyGroupExceptResource() resource.Resource {
	return &applyGroupExcept{}
}

func (rsc *applyGroupExcept) typeName() string {
	return providerName + "_apply_group_except"
}

func (rsc *applyGroupExcept) junosName() string {
	return "apply-groups-except"
}

func (rsc *applyGroupExcept) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *applyGroupExcept) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *applyGroupExcept) Configure(
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

func (rsc *applyGroupExcept) Schema(
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
				Description: "Name of group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 254),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"prefix": schema.StringAttribute{
				Required:    true,
				Description: "Prefix path to define where apply-groups-except must be set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

type applyGroupExceptData struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Prefix types.String `tfsdk:"prefix"`
}

func (rsc *applyGroupExcept) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config applyGroupExceptData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Prefix.IsNull() && !config.Prefix.IsUnknown() {
		if !strings.HasSuffix(config.Prefix.ValueString(), " ") {
			resp.Diagnostics.AddAttributeError(
				path.Root("Prefix"),
				"Bad Value Error",
				"prefix must be end with a space",
			)
		}
		if strings.TrimSpace(config.Prefix.ValueString()) == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("prefix"),
				"Empty Prefix",
				"prefix is only fill with space(s)",
			)
		}
	}
}

func (rsc *applyGroupExcept) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan applyGroupExceptData
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
	if strings.TrimSpace(plan.Prefix.ValueString()) == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("prefix"),
			"Empty Prefix",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "prefix")+" or with only space(s)",
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		nil,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			var result applyGroupExceptData
			err := result.read(fnCtx, plan.Name.ValueString(), plan.Prefix.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if result.nullID() {
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

func (rsc *applyGroupExcept) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data applyGroupExceptData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.Prefix.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *applyGroupExcept) Update(
	_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse,
) {
	// no-op
}

func (rsc *applyGroupExcept) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state applyGroupExceptData
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

func (rsc *applyGroupExcept) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data applyGroupExceptData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <name>"+junos.IDSeparator+"<prefix>)",
	)
}

func (rscData *applyGroupExceptData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + rscData.Prefix.ValueString())
}

func (rscData *applyGroupExceptData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *applyGroupExceptData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	prefix := rscData.Prefix.ValueString()
	if !strings.HasSuffix(prefix, " ") {
		prefix += " "
	}

	configSet := []string{
		junos.SetLS + prefix + "apply-groups-except \"" + rscData.Name.ValueString() + "\"",
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *applyGroupExceptData) read(
	_ context.Context, name, prefix string, junSess *junos.Session,
) error {
	if !strings.HasSuffix(prefix, " ") {
		prefix += " "
	}
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		prefix + "apply-groups-except" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		for item := range strings.SplitSeq(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if strings.Trim(strings.TrimSpace(itemTrim), "\"") == name {
				rscData.Name = types.StringValue(name)
				rscData.Prefix = types.StringValue(prefix)
				rscData.fillID()

				break
			}
		}
	}

	return nil
}

func (rscData *applyGroupExceptData) del(
	_ context.Context, junSess *junos.Session,
) error {
	prefix := rscData.Prefix.ValueString()
	if !strings.HasSuffix(prefix, " ") {
		prefix += " "
	}

	configSet := []string{
		junos.DeleteLS + prefix + "apply-groups-except \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
