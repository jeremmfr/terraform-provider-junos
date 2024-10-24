package providerfwk

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &securityGlobalPolicyUnordered{}
	_ resource.ResourceWithConfigure      = &securityGlobalPolicyUnordered{}
	_ resource.ResourceWithValidateConfig = &securityGlobalPolicyUnordered{}
	_ resource.ResourceWithImportState    = &securityGlobalPolicyUnordered{}
)

type securityGlobalPolicyUnordered struct {
	client *junos.Client
}

func newSecurityGlobalPolicyUnorderedResource() resource.Resource {
	return &securityGlobalPolicyUnordered{}
}

func (rsc *securityGlobalPolicyUnordered) typeName() string {
	return providerName + "_security_global_policy_unordered"
}

func (rsc *securityGlobalPolicyUnordered) junosName() string {
	return "security policies global"
}

func (rsc *securityGlobalPolicyUnordered) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityGlobalPolicyUnordered) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityGlobalPolicyUnordered) Configure(
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

func (rsc *securityGlobalPolicyUnordered) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `security_global_policy`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"policy": schema.SetNestedBlock{
				Description: "For each policy name.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityGlobalPolicyBlockPolicy{}.attributesSchema(),
					Blocks:     securityGlobalPolicyBlockPolicy{}.blocksSchema(),
				},
			},
		},
	}
}

type securityGlobalPolicyUnorderedConfig struct {
	ID     types.String `tfsdk:"id"`
	Policy types.Set    `tfsdk:"policy"`
}

func (rsc *securityGlobalPolicyUnordered) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityGlobalPolicyUnorderedConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Policy.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("policy").AtName("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one policy block must be specified",
		)
	} else if !config.Policy.IsUnknown() {
		var configPolicy []securityGlobalPolicyBlockPolicyConfig
		asDiags := config.Policy.ElementsAs(ctx, &configPolicy, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		policyName := make(map[string]struct{})
		for i, block := range configPolicy {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := policyName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple policy blocks with the same name %q", name),
					)
				}
				policyName[name] = struct{}{}
			}
			if block.MatchApplication.IsNull() && block.MatchDynamicApplication.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("policy").AtListIndex(i).AtName("name"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("at least one of match_application or match_dynamic_application "+
						"must be specified in policy %q", block.Name.ValueString()),
				)
			}
			if block.PermitApplicationServices != nil {
				if block.PermitApplicationServices.isEmpty() {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("permit_application_services"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("permit_application_services block is empty in policy %q", block.Name.ValueString()),
					)
				} else if block.PermitApplicationServices.hasKnownValue() &&
					!block.Then.IsNull() && !block.Then.IsUnknown() && block.Then.ValueString() != junos.PermitW {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("then"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("then is not %q (got %q) and permit_application_services is set in policy %q",
							junos.PermitW, block.Then.ValueString(), block.Name.ValueString()),
					)
				}
				if !block.PermitApplicationServices.RedirectWx.IsNull() &&
					!block.PermitApplicationServices.RedirectWx.IsUnknown() &&
					!block.PermitApplicationServices.ReverseRedirectWx.IsNull() &&
					!block.PermitApplicationServices.ReverseRedirectWx.IsUnknown() {
					resp.Diagnostics.AddAttributeError(
						path.Root("policy").AtListIndex(i).AtName("redirect_wx"),
						tfdiag.ConflictConfigErrSummary,
						fmt.Sprintf("redirect_wx and reverse_redirect_wx enabled both in policy %q", block.Name.ValueString()),
					)
				}
			}
		}
	}
}

func (rsc *securityGlobalPolicyUnordered) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityGlobalPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			if !junSess.CheckCompatibilitySecurity() {
				resp.Diagnostics.AddError(
					tfdiag.CompatibilityErrSummary,
					rsc.junosName()+junSess.SystemInformation.NotCompatibleMsg(),
				)

				return false
			}
			var check securityGlobalPolicyData
			if err := check.read(fnCtx, junSess); err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if len(check.Policy) > 0 {
				resp.Diagnostics.AddError(tfdiag.DuplicateConfigErrSummary, rsc.junosName()+" already exists")

				return false
			}

			return true
		},
		nil,
		&plan,
		resp,
	)
}

func (rsc *securityGlobalPolicyUnordered) Read(
	ctx context.Context, _ resource.ReadRequest, resp *resource.ReadResponse,
) {
	var data securityGlobalPolicyData

	var _ resourceDataReadWithoutArg = &data
	defaultResourceRead(
		ctx,
		rsc,
		nil,
		&data,
		nil,
		resp,
	)
}

func (rsc *securityGlobalPolicyUnordered) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan securityGlobalPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceUpdate(
		ctx,
		rsc,
		&plan,
		&plan,
		resp,
	)
}

func (rsc *securityGlobalPolicyUnordered) Delete(
	ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var empty securityGlobalPolicyData

	defaultResourceDelete(
		ctx,
		rsc,
		&empty,
		resp,
	)
}

func (rsc *securityGlobalPolicyUnordered) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityGlobalPolicyData

	var _ resourceDataReadWithoutArg = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		"",
	)
}
