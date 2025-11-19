package provider

import (
	"context"
	"errors"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &groupRaw{}
	_ resource.ResourceWithConfigure      = &groupRaw{}
	_ resource.ResourceWithValidateConfig = &groupRaw{}
	_ resource.ResourceWithImportState    = &groupRaw{}
)

type groupRaw struct {
	client *junos.Client
}

func newGroupRawResource() resource.Resource {
	return &groupRaw{}
}

func (rsc *groupRaw) typeName() string {
	return providerName + "_group_raw"
}

func (rsc *groupRaw) junosName() string {
	return "groups"
}

func (rsc *groupRaw) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *groupRaw) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *groupRaw) Configure(
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

func (rsc *groupRaw) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<name>_-_<format>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 254),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"config": schema.StringAttribute{
				Required:    true,
				Description: "The raw configuration to load.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"format": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("text"),
				Description: "The format used for the configuration data.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("text", "set"),
				},
			},
		},
	}
}

type groupRawData struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Config types.String `tfsdk:"config"`
	Format types.String `tfsdk:"format"`
}

func (rsc *groupRaw) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config groupRawData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Format.ValueString() == "set" {
		for item := range strings.SplitSeq(config.Config.ValueString(), "\n") {
			if strings.TrimSpace(item) == "" {
				continue
			}
			if !strings.HasPrefix(item, "set ") {
				resp.Diagnostics.AddAttributeError(
					path.Root("config"),
					"Bad Value Error",
					"all line of config must be start with 'set ' when format = 'set'",
				)
			}
		}
	}
}

func (rsc *groupRaw) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan groupRawData
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

	junSess, err := rsc.junosClient().StartNewSession(ctx)
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

	groupExists, err := checkGroupRawExists(ctx, plan.Name.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

		return
	}
	if groupExists {
		resp.Diagnostics.AddError(
			tfdiag.DuplicateConfigErrSummary,
			defaultResourceAlreadyExistsMessage(rsc, plan.Name),
		)

		return
	}

	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf(ctx, "create resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	groupExists, err = checkGroupRawExists(ctx, plan.Name.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

		return
	}
	if !groupExists {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.Name),
		)

		return
	}

	plan.fillID()

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rsc *groupRaw) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data groupRawData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Format = state.Format

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

func (rsc *groupRaw) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state groupRawData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	junSess, err := rsc.junosClient().StartNewSession(ctx)
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

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf(ctx, "update resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rsc *groupRaw) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state groupRawData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	junSess, err := rsc.junosClient().StartNewSession(ctx)
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

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

		return
	}
	warns, err := junSess.CommitConf(ctx, "delete resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}
}

func (rsc *groupRaw) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	var data groupRawData

	idList := strings.Split(req.ID, junos.IDSeparator)
	if len(idList) > 1 {
		data.Format = types.StringValue(idList[1])
	}

	if err := data.read(ctx, idList[0], junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if data.nullID() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be <name> or <name>_-_<format>)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkGroupRawExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"groups \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *groupRawData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString() + junos.IDSeparator + rscData.Format.ValueString())
}

func (rscData *groupRawData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *groupRawData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	switch rscData.Format.ValueString() {
	case "set":
		rawConfig := strings.Builder{}
		for item := range strings.SplitSeq(rscData.Config.ValueString(), "\n") {
			if strings.TrimSpace(item) == "" {
				continue
			}
			if !strings.HasPrefix(item, "set ") {
				return path.Root("config"),
					errors.New("all line of config must be start with 'set ' when format = 'set'")
			}
			_, _ = rawConfig.WriteString(
				strings.Replace(item, "set ", "set groups \""+rscData.Name.ValueString()+"\" ", 1) +
					"\n",
			)
		}

		return path.Empty(), junSess.ConfigLoad("set", "text", rawConfig.String())
	case "text":
		fallthrough
	default:
		rawConfig := "groups {\n\"" + rscData.Name.ValueString() + "\" {\n" + rscData.Config.ValueString() + "}\n}\n"

		// merge action as there is a delete of group before set when update
		return path.Empty(), junSess.ConfigLoad("merge", "text", rawConfig)
	}
}

func (rscData *groupRawData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	switch rscData.Format.ValueString() {
	case "set":
		showConfig, err := junSess.Command(junos.CmdShowConfig +
			"groups \"" + name + "\"" + junos.PipeDisplaySetRelative)
		if err != nil {
			return err
		}
		if showConfig != junos.EmptyW {
			rscData.Name = types.StringValue(name)
			rscData.fillID()

			configFiltered := strings.Builder{}
			for item := range strings.SplitSeq(showConfig, "\n") {
				if strings.Contains(item, junos.XMLStartTagConfigOut) {
					continue
				}
				if item == "" {
					continue
				}
				if strings.Contains(item, junos.XMLEndTagConfigOut) {
					break
				}
				_, _ = configFiltered.WriteString(item + "\n")
			}

			rscData.Config = types.StringValue(configFiltered.String())
		}
	case "text":
		fallthrough
	default:
		showConfig, err := junSess.Command(junos.CmdShowConfig +
			"groups \"" + name + "\"")
		if err != nil {
			return err
		}
		if showConfig != junos.EmptyW {
			rscData.Name = types.StringValue(name)
			rscData.Format = types.StringValue("text")
			rscData.fillID()

			configFiltered := strings.Builder{}
			for item := range strings.SplitSeq(showConfig, "\n") {
				if strings.Contains(item, junos.XMLStartTagConfigOut) {
					continue
				}
				if item == "" {
					continue
				}
				if strings.Contains(item, junos.XMLEndTagConfigOut) {
					break
				}
				_, _ = configFiltered.WriteString(item + "\n")
			}
			rscData.Config = types.StringValue(configFiltered.String())
		}
	}

	return nil
}

func (rscData *groupRawData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete groups \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
