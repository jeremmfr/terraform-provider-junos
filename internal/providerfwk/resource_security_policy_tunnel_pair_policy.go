package providerfwk

import (
	"context"
	"fmt"
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
	_ resource.Resource                = &securityPolicyTunnelPairPolicy{}
	_ resource.ResourceWithConfigure   = &securityPolicyTunnelPairPolicy{}
	_ resource.ResourceWithImportState = &securityPolicyTunnelPairPolicy{}
)

type securityPolicyTunnelPairPolicy struct {
	client *junos.Client
}

func newSecurityPolicyTunnelPairPolicyResource() resource.Resource {
	return &securityPolicyTunnelPairPolicy{}
}

func (rsc *securityPolicyTunnelPairPolicy) typeName() string {
	return providerName + "_security_policy_tunnel_pair_policy"
}

func (rsc *securityPolicyTunnelPairPolicy) junosName() string {
	return "security policy tunnel pair policy"
}

func (rsc *securityPolicyTunnelPairPolicy) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityPolicyTunnelPairPolicy) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityPolicyTunnelPairPolicy) Configure(
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

func (rsc *securityPolicyTunnelPairPolicy) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Provides a tunnel pair policy resource options in each policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<zone_a>" + junos.IDSeparator + "<policy_a_to_b>" +
					junos.IDSeparator +
					"<zone_b>" + junos.IDSeparator + "<policy_b_to_a>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_a": schema.StringAttribute{
				Required:    true,
				Description: "The name of first zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"zone_b": schema.StringAttribute{
				Required:    true,
				Description: "The name of second zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"policy_a_to_b": schema.StringAttribute{
				Required:    true,
				Description: "The name of policy when from zone zone_a to zone zone_b.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"policy_b_to_a": schema.StringAttribute{
				Required:    true,
				Description: "The name of policy when from zone zone_b to zone zone_a.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
		},
	}
}

type securityPolicyTunnelPairPolicyData struct {
	ID         types.String `tfsdk:"id"`
	ZoneA      types.String `tfsdk:"zone_a"`
	ZoneB      types.String `tfsdk:"zone_b"`
	PolicyAtoB types.String `tfsdk:"policy_a_to_b"`
	PolicyBtoA types.String `tfsdk:"policy_b_to_a"`
}

func (rsc *securityPolicyTunnelPairPolicy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityPolicyTunnelPairPolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.ZoneA.ValueString() == "" ||
		plan.ZoneB.ValueString() == "" ||
		plan.PolicyAtoB.ValueString() == "" ||
		plan.PolicyBtoA.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Empty Zone",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "zone_a, zone_b, policy_a_to_b or policy_b_to_a"),
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
				plan.ZoneA.ValueString(),
				plan.ZoneB.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !policyExists {
				resp.Diagnostics.AddError(
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("security policy from %q to %q doesn't exist",
						plan.ZoneA.ValueString(), plan.ZoneB.ValueString()),
				)

				return false
			}
			policyExists, err = checkSecurityPolicyExists(
				fnCtx,
				plan.ZoneB.ValueString(),
				plan.ZoneA.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !policyExists {
				resp.Diagnostics.AddError(
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("security policy from %q to %q doesn't exist",
						plan.ZoneB.ValueString(), plan.ZoneA.ValueString()),
				)

				return false
			}
			pairPolicyExists, err := checkSecurityPolicyTunnelPairPolicyExists(
				fnCtx,
				plan.ZoneA.ValueString(),
				plan.PolicyAtoB.ValueString(),
				plan.ZoneB.ValueString(),
				plan.PolicyBtoA.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if pairPolicyExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf(rsc.junosName()+" %q(%q) / %q(%q) already exists",
						plan.ZoneA.ValueString(),
						plan.PolicyAtoB.ValueString(),
						plan.ZoneB.ValueString(),
						plan.PolicyBtoA.ValueString()),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			pairPolicyExists, err := checkSecurityPolicyTunnelPairPolicyExists(
				fnCtx,
				plan.ZoneA.ValueString(),
				plan.PolicyAtoB.ValueString(),
				plan.ZoneB.ValueString(),
				plan.PolicyBtoA.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !pairPolicyExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					rsc.junosName()+" does not exists after commit "+
						"=> check your config",
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *securityPolicyTunnelPairPolicy) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityPolicyTunnelPairPolicyData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom4String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.ZoneA.ValueString(),
			state.PolicyAtoB.ValueString(),
			state.ZoneB.ValueString(),
			state.PolicyBtoA.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *securityPolicyTunnelPairPolicy) Update(
	_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse,
) {
}

func (rsc *securityPolicyTunnelPairPolicy) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityPolicyTunnelPairPolicyData
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

func (rsc *securityPolicyTunnelPairPolicy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityPolicyTunnelPairPolicyData

	var _ resourceDataReadFrom4String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be "+
			"<zone_a>"+junos.IDSeparator+
			"<policy_a_to_b>"+junos.IDSeparator+
			"<zone_b>"+junos.IDSeparator+
			"<policy_b_to_a>)",
	)
}

func checkSecurityPolicyTunnelPairPolicyExists(
	_ context.Context, zoneA, policyAtoB, zoneB, policyBtoA string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfigPairAtoB, err := junSess.Command(junos.CmdShowConfig +
		"security policies from-zone " + zoneA + " to-zone " + zoneB + " policy " + policyAtoB +
		" then permit tunnel pair-policy" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	showConfigPairBtoA, err := junSess.Command(junos.CmdShowConfig +
		"security policies from-zone " + zoneB + " to-zone " + zoneA + " policy " + policyBtoA +
		" then permit tunnel pair-policy" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfigPairAtoB == junos.EmptyW && showConfigPairBtoA == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityPolicyTunnelPairPolicyData) fillID() {
	rscData.ID = types.StringValue(
		rscData.ZoneA.ValueString() + junos.IDSeparator +
			rscData.PolicyAtoB.ValueString() + junos.IDSeparator +
			rscData.ZoneB.ValueString() + junos.IDSeparator +
			rscData.PolicyBtoA.ValueString(),
	)
}

func (rscData *securityPolicyTunnelPairPolicyData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityPolicyTunnelPairPolicyData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 2)

	configSet = append(configSet, "set security policies from-zone "+
		rscData.ZoneA.ValueString()+" to-zone "+rscData.ZoneB.ValueString()+
		" policy "+rscData.PolicyAtoB.ValueString()+
		" then permit tunnel pair-policy "+rscData.PolicyBtoA.ValueString())
	configSet = append(configSet, "set security policies from-zone "+
		rscData.ZoneB.ValueString()+" to-zone "+rscData.ZoneA.ValueString()+
		" policy "+rscData.PolicyBtoA.ValueString()+
		" then permit tunnel pair-policy "+rscData.PolicyAtoB.ValueString())

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityPolicyTunnelPairPolicyData) read(
	_ context.Context, zoneA, policyAtoB, zoneB, policyBtoA string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security policies from-zone " + zoneA + " to-zone " + zoneB + " policy " + policyAtoB +
		" then permit tunnel pair-policy" + junos.PipeDisplaySet)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.ZoneA = types.StringValue(zoneA)
		rscData.ZoneB = types.StringValue(zoneB)
		rscData.fillID()
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			if strings.Contains(item, " tunnel pair-policy ") {
				rscData.PolicyBtoA = types.StringValue(strings.TrimPrefix(item,
					"set security policies from-zone "+zoneA+" to-zone "+zoneB+
						" policy "+policyAtoB+" then permit tunnel pair-policy "))
			}
		}
	}
	showConfig, err = junSess.Command(junos.CmdShowConfig +
		"security policies from-zone " + zoneB + " to-zone " + zoneA + " policy " + policyBtoA +
		" then permit tunnel pair-policy" + junos.PipeDisplaySet)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.ZoneA = types.StringValue(zoneA)
		rscData.ZoneB = types.StringValue(zoneB)
		rscData.fillID()
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			if strings.Contains(item, " tunnel pair-policy ") {
				rscData.PolicyAtoB = types.StringValue(strings.TrimPrefix(item,
					"set security policies from-zone "+zoneB+" to-zone "+zoneA+
						" policy "+policyBtoA+" then permit tunnel pair-policy "))
			}
		}
	}

	return nil
}

func (rscData *securityPolicyTunnelPairPolicyData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security policies" +
			" from-zone " + rscData.ZoneA.ValueString() + " to-zone " + rscData.ZoneB.ValueString() +
			" policy " + rscData.PolicyAtoB.ValueString() +
			" then permit tunnel pair-policy " + rscData.PolicyBtoA.ValueString(),
		"delete security policies" +
			" from-zone " + rscData.ZoneB.ValueString() + " to-zone " + rscData.ZoneA.ValueString() +
			" policy " + rscData.PolicyBtoA.ValueString() +
			" then permit tunnel pair-policy " + rscData.PolicyAtoB.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
