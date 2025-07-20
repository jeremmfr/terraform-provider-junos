package provider

import (
	"context"
	"errors"
	"regexp"
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
	_ resource.Resource                   = &ribGroup{}
	_ resource.ResourceWithConfigure      = &ribGroup{}
	_ resource.ResourceWithValidateConfig = &ribGroup{}
	_ resource.ResourceWithImportState    = &ribGroup{}
)

type ribGroup struct {
	client *junos.Client
}

func newRibGroupResource() resource.Resource {
	return &ribGroup{}
}

func (rsc *ribGroup) typeName() string {
	return providerName + "_rib_group"
}

func (rsc *ribGroup) junosName() string {
	return "routing-options rib-groups"
}

func (rsc *ribGroup) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *ribGroup) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *ribGroup) Configure(
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

func (rsc *ribGroup) Schema(
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
				Description: "The name of rib group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"import_policy": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of policy for import route.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"import_rib": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of import routing table.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
						stringvalidator.RegexMatches(regexp.MustCompile(
							`^(.+\.)?(inet6?\.0)$`),
							"must be equal to or end with `inet.0` or `inet6.0`",
						),
					),
				},
			},
			"export_rib": schema.StringAttribute{
				Optional:    true,
				Description: "Export routing table.",
				Validators: []validator.String{
					tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(.+\.)?(inet6?\.0)$`),
						"must be equal to or end with `inet.0` or `inet6.0`",
					),
				},
			},
		},
	}
}

type ribGroupData struct {
	ID           types.String   `tfsdk:"id"            tfdata:"skip_isempty"`
	Name         types.String   `tfsdk:"name"          tfdata:"skip_isempty"`
	ImportPolicy []types.String `tfsdk:"import_policy"`
	ImportRib    []types.String `tfsdk:"import_rib"`
	ExportRib    types.String   `tfsdk:"export_rib"`
}

func (rscData *ribGroupData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

type ribGroupConfig struct {
	ID           types.String `tfsdk:"id"            tfdata:"skip_isempty"`
	Name         types.String `tfsdk:"name"          tfdata:"skip_isempty"`
	ImportPolicy types.List   `tfsdk:"import_policy"`
	ImportRib    types.List   `tfsdk:"import_rib"`
	ExportRib    types.String `tfsdk:"export_rib"`
}

func (config *ribGroupConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(config)
}

func (rsc *ribGroup) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config ribGroupConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one of arguments need to be set (in addition to `name`)",
		)
	}
}

func (rsc *ribGroup) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan ribGroupData
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
			groupExists, err := checkRibGroupExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if groupExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			groupExists, err := checkRibGroupExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !groupExists {
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

func (rsc *ribGroup) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data ribGroupData
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

func (rsc *ribGroup) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state ribGroupData
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

func (rsc *ribGroup) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state ribGroupData
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

func (rsc *ribGroup) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data ribGroupData

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

func checkRibGroupExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		junos.RoutingOptionsWS + "rib-groups \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *ribGroupData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *ribGroupData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *ribGroupData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	if rscData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name`)")
	}

	configSet := make([]string, 0, 100)
	setPrefix := "set routing-options rib-groups \"" + rscData.Name.ValueString() + "\" "

	for _, v := range rscData.ImportPolicy {
		configSet = append(configSet, setPrefix+"import-policy \""+v.ValueString()+"\"")
	}
	for _, v := range rscData.ImportRib {
		configSet = append(configSet, setPrefix+"import-rib "+v.ValueString())
	}
	if v := rscData.ExportRib.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"export-rib "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *ribGroupData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		junos.RoutingOptionsWS + "rib-groups \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "import-policy "):
				rscData.ImportPolicy = append(rscData.ImportPolicy,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "import-rib "):
				rscData.ImportRib = append(rscData.ImportRib,
					types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "export-rib "):
				rscData.ExportRib = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (rscData *ribGroupData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete routing-options rib-groups \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
