package provider

import (
	"context"
	"errors"
	"maps"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                   = &applicationSet{}
	_ resource.ResourceWithConfigure      = &applicationSet{}
	_ resource.ResourceWithValidateConfig = &applicationSet{}
	_ resource.ResourceWithImportState    = &applicationSet{}
)

type applicationSet struct {
	client *junos.Client
}

func newApplicationSetResource() resource.Resource {
	return &applicationSet{}
}

func (rsc *applicationSet) typeName() string {
	return providerName + "_application_set"
}

func (rsc *applicationSet) junosName() string {
	return "applications application-set"
}

func (rsc *applicationSet) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *applicationSet) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *applicationSet) Configure(
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

func (rsc *applicationSet) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	attributes := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "An identifier for the resource with format `<name>`.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
	maps.Copy(attributes, applicationSetAttrData{}.attributesSchema())
	// add RequiresReplace on name to force replacement of application_set resource
	// when value change (useless in applications resource)
	nameAttribute := attributes["name"].(schema.StringAttribute)
	nameAttribute.PlanModifiers = []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
	}
	attributes["name"] = nameAttribute

	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes:  attributes,
	}
}

type applicationSetData struct {
	applicationSetAttrData

	ID types.String `tfsdk:"id"`
}

type applicationSetAttrData struct {
	Name           types.String   `tfsdk:"name"            tfdata:"identifier,skip_isempty"`
	Applications   []types.String `tfsdk:"applications"`
	ApplicationSet []types.String `tfsdk:"application_set"`
	Description    types.String   `tfsdk:"description"`
}

func (rscData *applicationSetAttrData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

func (applicationSetAttrData) attributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required:    true,
			Description: "Application set name.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 63),
				tfvalidator.StringFormat(tfvalidator.DefaultFormat),
			},
		},
		"applications": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Application to be included in the set.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.NoNullValues(),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				),
			},
		},
		"application_set": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Application-set to be included in the set.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.NoNullValues(),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				),
			},
		},
		"description": schema.StringAttribute{
			Optional:    true,
			Description: "Description for application-set.",
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 900),
				tfvalidator.StringDoubleQuoteExclusion(),
			},
		},
	}
}

type applicationSetConfig struct {
	applicationSetAttrConfig

	ID types.String `tfsdk:"id"`
}

type applicationSetAttrConfig struct {
	Name           types.String `tfsdk:"name"            tfdata:"skip_isempty"`
	Applications   types.List   `tfsdk:"applications"`
	ApplicationSet types.List   `tfsdk:"application_set"`
	Description    types.String `tfsdk:"description"`
}

func (config *applicationSetAttrConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(config)
}

func (config *applicationSetAttrConfig) validateConfig(
	_ context.Context, rootPath *path.Path, blockErrorSuffix string, resp *resource.ValidateConfigResponse,
) {
	if config.isEmpty() {
		errPath := path.Root("name")
		if rootPath != nil {
			errPath = *rootPath
		}
		resp.Diagnostics.AddAttributeError(
			errPath,
			tfdiag.MissingConfigErrSummary,
			"at least one of applications, application_set or description must be specified"+blockErrorSuffix,
		)
	}
}

func (rsc *applicationSet) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config applicationSetConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.applicationSetAttrConfig.validateConfig(ctx, nil, "", resp)
}

func (rsc *applicationSet) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan applicationSetData
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
			applicationSetExists, err := checkApplicationSetExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if applicationSetExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			applicationSetExists, err := checkApplicationSetExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !applicationSetExists {
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

func (rsc *applicationSet) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data applicationSetData
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

func (rsc *applicationSet) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state applicationSetData
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

func (rsc *applicationSet) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state applicationSetData
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

func (rsc *applicationSet) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data applicationSetData

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

func checkApplicationSetExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"applications application-set " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *applicationSetData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *applicationSetData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *applicationSetData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	if rscData.applicationSetAttrData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name`)")
	}

	return path.Empty(), junSess.ConfigSet(rscData.applicationSetAttrData.configSet())
}

func (rscData *applicationSetAttrData) configSet() []string {
	configSet := make([]string, 0, len(rscData.Applications)+len(rscData.ApplicationSet)+1)
	setPrefix := "set applications application-set " + rscData.Name.ValueString() + " "

	for _, v := range rscData.Applications {
		configSet = append(configSet, setPrefix+"application "+v.ValueString())
	}
	for _, v := range rscData.ApplicationSet {
		configSet = append(configSet, setPrefix+"application-set "+v.ValueString())
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}

	return configSet
}

func (rscData *applicationSetData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"applications application-set " + name + junos.PipeDisplaySetRelative)
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
			rscData.applicationSetAttrData.read(itemTrim)
		}
	}

	return nil
}

func (rscData *applicationSetAttrData) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "application "):
		rscData.Applications = append(rscData.Applications, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "application-set "):
		rscData.ApplicationSet = append(rscData.ApplicationSet, types.StringValue(itemTrim))
	case balt.CutPrefixInString(&itemTrim, "description "):
		rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
	}
}

func (rscData *applicationSetData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete applications application-set " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
