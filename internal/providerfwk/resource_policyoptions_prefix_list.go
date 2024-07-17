package providerfwk

import (
	"context"
	"html"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	_ resource.Resource                = &policyoptionsPrefixList{}
	_ resource.ResourceWithConfigure   = &policyoptionsPrefixList{}
	_ resource.ResourceWithImportState = &policyoptionsPrefixList{}
)

type policyoptionsPrefixList struct {
	client *junos.Client
}

func newPolicyoptionsPrefixListResource() resource.Resource {
	return &policyoptionsPrefixList{}
}

func (rsc *policyoptionsPrefixList) typeName() string {
	return providerName + "_policyoptions_prefix_list"
}

func (rsc *policyoptionsPrefixList) junosName() string {
	return "policy-options prefix-list"
}

func (rsc *policyoptionsPrefixList) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *policyoptionsPrefixList) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *policyoptionsPrefixList) Configure(
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

func (rsc *policyoptionsPrefixList) Schema(
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
				Description: "Prefix list name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"apply_path": schema.StringAttribute{
				Optional:    true,
				Description: "Apply IP prefixes from a configuration statement.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
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
			"prefix": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Address prefixes.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringCIDRNetwork(),
					),
				},
			},
		},
	}
}

type policyoptionsPrefixListData struct {
	ID        types.String   `tfsdk:"id"`
	Name      types.String   `tfsdk:"name"`
	ApplyPath types.String   `tfsdk:"apply_path"`
	DynamicDB types.Bool     `tfsdk:"dynamic_db"`
	Prefix    []types.String `tfsdk:"prefix"`
}

func (rsc *policyoptionsPrefixList) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan policyoptionsPrefixListData
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
			listExists, err := checkPolicyoptionsPrefixListExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if listExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			listExists, err := checkPolicyoptionsPrefixListExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !listExists {
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

func (rsc *policyoptionsPrefixList) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data policyoptionsPrefixListData
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

func (rsc *policyoptionsPrefixList) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state policyoptionsPrefixListData
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

func (rsc *policyoptionsPrefixList) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state policyoptionsPrefixListData
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

func (rsc *policyoptionsPrefixList) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data policyoptionsPrefixListData

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

func checkPolicyoptionsPrefixListExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"policy-options prefix-list \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *policyoptionsPrefixListData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *policyoptionsPrefixListData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *policyoptionsPrefixListData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set policy-options prefix-list \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.ApplyPath.ValueString(); v != "" {
		replaceSign := strings.ReplaceAll(strings.ReplaceAll(v, "<", "&lt;"), ">", "&gt;")
		configSet = append(configSet, setPrefix+" apply-path \""+replaceSign+"\"")
	}
	if rscData.DynamicDB.ValueBool() {
		configSet = append(configSet, setPrefix+"dynamic-db")
	}
	for _, v := range rscData.Prefix {
		configSet = append(configSet, setPrefix+v.ValueString())
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *policyoptionsPrefixListData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"policy-options prefix-list \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "apply-path "):
				rscData.ApplyPath = types.StringValue(html.UnescapeString(strings.Trim(itemTrim, "\"")))
			case itemTrim == "dynamic-db":
				rscData.DynamicDB = types.BoolValue(true)
			case strings.Contains(itemTrim, "/"):
				rscData.Prefix = append(rscData.Prefix, types.StringValue(itemTrim))
			}
		}
	}

	return nil
}

func (rscData *policyoptionsPrefixListData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete policy-options prefix-list \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
