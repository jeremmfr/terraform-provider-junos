package provider

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &securityPolicyUnordered{}
	_ resource.ResourceWithConfigure      = &securityPolicyUnordered{}
	_ resource.ResourceWithValidateConfig = &securityPolicyUnordered{}
	_ resource.ResourceWithImportState    = &securityPolicyUnordered{}
)

type securityPolicyUnordered struct {
	client *junos.Client
}

func newSecurityPolicyUnorderedResource() resource.Resource {
	return &securityPolicyUnordered{}
}

func (rsc *securityPolicyUnordered) typeName() string {
	return providerName + "_security_policy_unordered"
}

func (rsc *securityPolicyUnordered) junosName() string {
	return "security policy"
}

func (rsc *securityPolicyUnordered) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityPolicyUnordered) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityPolicyUnordered) Configure(
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

func (rsc *securityPolicyUnordered) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes:  securityPolicyData{}.attributesSchema(),
		Blocks: map[string]schema.Block{
			"policy": schema.SetNestedBlock{
				Description: "For each name of policy.",
				NestedObject: schema.NestedBlockObject{
					Attributes: securityPolicyBlockPolicy{}.attributesSchema(),
					Blocks:     securityPolicyBlockPolicy{}.blocksSchema(),
				},
			},
		},
	}
}

type securityPolicyUnorderedConfig struct {
	ID       types.String `tfsdk:"id"`
	FromZone types.String `tfsdk:"from_zone"`
	ToZone   types.String `tfsdk:"to_zone"`
	Policy   types.Set    `tfsdk:"policy"`
}

func (rsc *securityPolicyUnordered) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityPolicyUnorderedConfig
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
		var configPolicy []securityPolicyBlockPolicyConfig
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
			if !block.PermitTunnelIpsecVpn.IsNull() && !block.PermitTunnelIpsecVpn.IsUnknown() &&
				!block.Then.IsNull() && !block.Then.IsUnknown() && block.Then.ValueString() != junos.PermitW {
				resp.Diagnostics.AddAttributeError(
					path.Root("policy").AtListIndex(i).AtName("then"),
					tfdiag.ConflictConfigErrSummary,
					fmt.Sprintf("then is not %q (got %q) and permit_tunnel_ipsec_vpn is set in policy %q",
						junos.PermitW, block.Then.ValueString(), block.Name.ValueString()),
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

func (rsc *securityPolicyUnordered) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.FromZone.ValueString() == "" || plan.ToZone.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Empty Zone",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "from_zone or to_zone"),
		)

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
			policyExists, err := checkSecurityPolicyExists(
				fnCtx,
				plan.FromZone.ValueString(),
				plan.ToZone.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if policyExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf(rsc.junosName()+" from %q to %q already exists",
						plan.FromZone.ValueString(), plan.ToZone.ValueString()),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			policyExists, err := checkSecurityPolicyExists(
				fnCtx,
				plan.FromZone.ValueString(),
				plan.ToZone.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !policyExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					fmt.Sprintf(rsc.junosName()+" from %q to %q not exists after commit "+
						"=> check your config", plan.FromZone.ValueString(), plan.ToZone.ValueString()),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *securityPolicyUnordered) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityPolicyData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.FromZone.ValueString(),
			state.ToZone.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *securityPolicyUnordered) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

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

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
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

	listLinesToPairPolicy, err := readSecurityPolicyTunnelPairPolicyLines(
		ctx,
		state.FromZone.ValueString(),
		state.ToZone.ValueString(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
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
	if err := junSess.ConfigSet(listLinesToPairPolicy); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())

		return
	}
	warns, err := junSess.CommitConf(ctx, "update resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityPolicyUnordered) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityPolicyData
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

func (rsc *securityPolicyUnordered) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityPolicyData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <zone>"+junos.IDSeparator+"<name>)",
	)
}
