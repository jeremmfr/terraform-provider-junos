package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
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
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &policyoptionsASPathGroup{}
	_ resource.ResourceWithConfigure      = &policyoptionsASPathGroup{}
	_ resource.ResourceWithValidateConfig = &policyoptionsASPathGroup{}
	_ resource.ResourceWithImportState    = &policyoptionsASPathGroup{}
)

type policyoptionsASPathGroup struct {
	client *junos.Client
}

func newPolicyoptionsASPathGroupResource() resource.Resource {
	return &policyoptionsASPathGroup{}
}

func (rsc *policyoptionsASPathGroup) typeName() string {
	return providerName + "_policyoptions_as_path_group"
}

func (rsc *policyoptionsASPathGroup) junosName() string {
	return "policy-options as-path-group"
}

func (rsc *policyoptionsASPathGroup) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *policyoptionsASPathGroup) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *policyoptionsASPathGroup) Configure(
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

func (rsc *policyoptionsASPathGroup) Schema(
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
				Description: "Name to identify AS path group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"dynamic_db": schema.BoolAttribute{
				Optional:    true,
				Description: "Object may exist in dynamic database.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"as_path": schema.ListNestedBlock{
				Description: "For each name of as-path to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name to identify AS path regular expression.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 250),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"path": schema.StringAttribute{
							Required:    true,
							Description: "AS path regular expression.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.ASPathRegularExpression),
							},
						},
					},
				},
			},
		},
	}
}

type policyoptionsASPathGroupData struct {
	ID        types.String                          `tfsdk:"id"`
	Name      types.String                          `tfsdk:"name"`
	DynamicDB types.Bool                            `tfsdk:"dynamic_db"`
	ASPath    []policyoptionsASPathGroupBlockASPAth `tfsdk:"as_path"`
}

type policyoptionsASPathGroupConfig struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	DynamicDB types.Bool   `tfsdk:"dynamic_db"`
	ASPath    types.List   `tfsdk:"as_path"`
}

type policyoptionsASPathGroupBlockASPAth struct {
	Name types.String `tfsdk:"name"`
	Path types.String `tfsdk:"path"`
}

func (rsc *policyoptionsASPathGroup) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config policyoptionsASPathGroupConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ASPath.IsNull() &&
		config.DynamicDB.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one of as_path or dynamic_db must be specified",
		)
	}
	if !config.ASPath.IsNull() && !config.ASPath.IsUnknown() {
		var configASPath []policyoptionsASPathGroupBlockASPAth
		asDiags := config.ASPath.ElementsAs(ctx, &configASPath, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		asPathName := make(map[string]struct{})
		for i, asPath := range configASPath {
			if asPath.Name.IsUnknown() {
				continue
			}
			name := asPath.Name.ValueString()
			if _, ok := asPathName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("as_path").AtListIndex(i).AtName("name"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple as_path blocks with the same name %q", name),
				)
			}
			asPathName[name] = struct{}{}
		}
	}
}

func (rsc *policyoptionsASPathGroup) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan policyoptionsASPathGroupData
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
			groupExists, err := checkPolicyoptionsAsPathGroupExists(fnCtx, plan.Name.ValueString(), junSess)
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
			groupExists, err := checkPolicyoptionsAsPathGroupExists(fnCtx, plan.Name.ValueString(), junSess)
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

func (rsc *policyoptionsASPathGroup) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data policyoptionsASPathGroupData
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

func (rsc *policyoptionsASPathGroup) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state policyoptionsASPathGroupData
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

func (rsc *policyoptionsASPathGroup) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state policyoptionsASPathGroupData
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

func (rsc *policyoptionsASPathGroup) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data policyoptionsASPathGroupData

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

func checkPolicyoptionsAsPathGroupExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"policy-options as-path-group \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *policyoptionsASPathGroupData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *policyoptionsASPathGroupData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *policyoptionsASPathGroupData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set policy-options as-path-group \"" + rscData.Name.ValueString() + "\" "

	asPathName := make(map[string]struct{})
	for i, block := range rscData.ASPath {
		name := block.Name.ValueString()
		if _, ok := asPathName[name]; ok {
			return path.Root("as_path").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple as_path blocks with the same name %q", name)
		}
		asPathName[name] = struct{}{}

		configSet = append(configSet, setPrefix+
			"as-path \""+block.Name.ValueString()+"\""+
			" \""+block.Path.ValueString()+"\"")
	}
	if rscData.DynamicDB.ValueBool() {
		configSet = append(configSet, setPrefix+"dynamic-db")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *policyoptionsASPathGroupData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"policy-options as-path-group \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case itemTrim == "dynamic-db":
				rscData.DynamicDB = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "as-path "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.ASPath = append(rscData.ASPath, policyoptionsASPathGroupBlockASPAth{
					Name: types.StringValue(strings.Trim(name, "\"")),
					Path: types.StringValue(strings.Trim(strings.TrimPrefix(itemTrim, name+" "), "\"")),
				})
			}
		}
	}

	return nil
}

func (rscData *policyoptionsASPathGroupData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete policy-options as-path-group \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
