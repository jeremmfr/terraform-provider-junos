package provider

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
	_ resource.Resource                   = &servicesSecurityIntelligencePolicy{}
	_ resource.ResourceWithConfigure      = &servicesSecurityIntelligencePolicy{}
	_ resource.ResourceWithValidateConfig = &servicesSecurityIntelligencePolicy{}
	_ resource.ResourceWithImportState    = &servicesSecurityIntelligencePolicy{}
)

type servicesSecurityIntelligencePolicy struct {
	client *junos.Client
}

func newServicesSecurityIntelligencePolicyResource() resource.Resource {
	return &servicesSecurityIntelligencePolicy{}
}

func (rsc *servicesSecurityIntelligencePolicy) typeName() string {
	return providerName + "_services_security_intelligence_policy"
}

func (rsc *servicesSecurityIntelligencePolicy) junosName() string {
	return "services security-intelligence policy"
}

func (rsc *servicesSecurityIntelligencePolicy) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *servicesSecurityIntelligencePolicy) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *servicesSecurityIntelligencePolicy) Configure(
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

func (rsc *servicesSecurityIntelligencePolicy) Schema(
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
				Description: "Security intelligence policy name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of policy.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"category": schema.ListNestedBlock{
				Description: "For each name of security intelligence category, configure a profile.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of security intelligence category.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"profile_name": schema.StringAttribute{
							Required:    true,
							Description: "Name of profile.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
				},
			},
		},
	}
}

type servicesSecurityIntelligencePolicyData struct {
	ID          types.String                                      `tfsdk:"id"`
	Name        types.String                                      `tfsdk:"name"`
	Description types.String                                      `tfsdk:"description"`
	Category    []servicesSecurityIntelligencePolicyBlockCategory `tfsdk:"category"`
}

type servicesSecurityIntelligencePolicyConfig struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Category    types.List   `tfsdk:"category"`
}

type servicesSecurityIntelligencePolicyBlockCategory struct {
	Name        types.String `tfsdk:"name"         tfdata:"identifier"`
	ProfileName types.String `tfsdk:"profile_name"`
}

func (rsc *servicesSecurityIntelligencePolicy) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config servicesSecurityIntelligencePolicyConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Category.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("category"),
			tfdiag.MissingConfigErrSummary,
			"category block must be specified",
		)
	} else if !config.Category.IsUnknown() {
		var configCategory []servicesSecurityIntelligencePolicyBlockCategory
		asDiags := config.Category.ElementsAs(ctx, &configCategory, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		categoryName := make(map[string]struct{})
		for i, block := range configCategory {
			if block.Name.IsUnknown() {
				continue
			}

			name := block.Name.ValueString()
			if _, ok := categoryName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("category").AtListIndex(i).AtName("name"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple category blocks with the same name %q", name),
				)
			}
			categoryName[name] = struct{}{}
		}
	}
}

func (rsc *servicesSecurityIntelligencePolicy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan servicesSecurityIntelligencePolicyData
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
			policyExists, err := checkServicesSecurityIntelligencePolicyExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if policyExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			policyExists, err := checkServicesSecurityIntelligencePolicyExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !policyExists {
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

func (rsc *servicesSecurityIntelligencePolicy) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data servicesSecurityIntelligencePolicyData
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

func (rsc *servicesSecurityIntelligencePolicy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state servicesSecurityIntelligencePolicyData
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

func (rsc *servicesSecurityIntelligencePolicy) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state servicesSecurityIntelligencePolicyData
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

func (rsc *servicesSecurityIntelligencePolicy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data servicesSecurityIntelligencePolicyData

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

func checkServicesSecurityIntelligencePolicyExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services security-intelligence policy \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *servicesSecurityIntelligencePolicyData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *servicesSecurityIntelligencePolicyData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *servicesSecurityIntelligencePolicyData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)
	setPrefix := "set services security-intelligence policy \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}

	categoryName := make(map[string]struct{})
	for i, block := range rscData.Category {
		name := block.Name.ValueString()
		if _, ok := categoryName[name]; ok {
			return path.Root("category").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple category blocks with the same name %q", name)
		}
		categoryName[name] = struct{}{}

		configSet = append(configSet, setPrefix+name+" \""+block.ProfileName.ValueString()+"\"")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *servicesSecurityIntelligencePolicyData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services security-intelligence policy \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case len(strings.Split(itemTrim, " ")) >= 2:
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.Category = tfdata.AppendPotentialNewBlock(rscData.Category, types.StringValue(name))
				category := &rscData.Category[len(rscData.Category)-1]
				balt.CutPrefixInString(&itemTrim, name+" ")

				category.ProfileName = types.StringValue(strings.Trim(itemTrim, "\""))
			}
		}
	}

	return nil
}

func (rscData *servicesSecurityIntelligencePolicyData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete services security-intelligence policy \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
