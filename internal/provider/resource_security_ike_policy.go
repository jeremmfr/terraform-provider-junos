package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &securityIkePolicy{}
	_ resource.ResourceWithConfigure      = &securityIkePolicy{}
	_ resource.ResourceWithValidateConfig = &securityIkePolicy{}
	_ resource.ResourceWithImportState    = &securityIkePolicy{}
)

type securityIkePolicy struct {
	client *junos.Client
}

func newSecurityIkePolicyResource() resource.Resource {
	return &securityIkePolicy{}
}

func (rsc *securityIkePolicy) typeName() string {
	return providerName + "_security_ike_policy"
}

func (rsc *securityIkePolicy) junosName() string {
	return "security ike policy"
}

func (rsc *securityIkePolicy) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityIkePolicy) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityIkePolicy) Configure(
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

func (rsc *securityIkePolicy) Schema(
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
				Description: "The name of IKE policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"proposals": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "IKE proposals list.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 32),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
			"proposal_set": schema.StringAttribute{
				Optional:    true,
				Description: "Types of default IKE proposal-set.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"basic",
						"compatible",
						"prime-128",
						"prime-256",
						"standard",
						"suiteb-gcm-128",
						"suiteb-gcm-256",
					),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of IKE policy.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("main"),
				Description: "IKE mode for Phase 1.",
				Validators: []validator.String{
					stringvalidator.OneOf("main", "aggressive"),
				},
			},
			"pre_shared_key_hexa": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Preshared key with format as hexadecimal.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"pre_shared_key_hexa_wo": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
				Description: "Preshared key with format as hexadecimal, not stored in state.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
					stringvalidator.AlsoRequires(path.MatchRoot("pre_shared_key_hexa_wo_version")),
				},
			},
			"pre_shared_key_hexa_wo_version": schema.Int64Attribute{
				Optional:    true,
				Description: "Version of `pre_shared_key_hexa_wo` to trigger the sending of its value.",
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("pre_shared_key_hexa_wo")),
				},
			},
			"pre_shared_key_text": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Preshared key wit format as text.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"pre_shared_key_text_wo": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
				Description: "Preshared key wit format as text, not stored in state.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
					stringvalidator.AlsoRequires(path.MatchRoot("pre_shared_key_text_wo_version")),
				},
			},
			"pre_shared_key_text_wo_version": schema.Int64Attribute{
				Optional:    true,
				Description: "Version of `pre_shared_key_text_wo` to trigger the sending of its value.",
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("pre_shared_key_text_wo")),
				},
			},
			"reauth_frequency": schema.Int64Attribute{
				Optional:    true,
				Description: "Re-auth Peer after reauth-frequency times hard lifetime. (0-100)",
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
			},
		},
	}
}

type securityIkePolicyData struct {
	ID                        types.String   `tfsdk:"id"`
	Name                      types.String   `tfsdk:"name"`
	Description               types.String   `tfsdk:"description"`
	Mode                      types.String   `tfsdk:"mode"`
	PreSharedKeyHexa          types.String   `tfsdk:"pre_shared_key_hexa"`
	PreSharedKeyHexaWO        types.String   `tfsdk:"pre_shared_key_hexa_wo"`
	PreSharedKeyHexaWOVersion types.Int64    `tfsdk:"pre_shared_key_hexa_wo_version"`
	PreSharedKeyText          types.String   `tfsdk:"pre_shared_key_text"`
	PreSharedKeyTextWO        types.String   `tfsdk:"pre_shared_key_text_wo"`
	PreSharedKeyTextWOVersion types.Int64    `tfsdk:"pre_shared_key_text_wo_version"`
	Proposals                 []types.String `tfsdk:"proposals"`
	ProposalSet               types.String   `tfsdk:"proposal_set"`
	ReauthFrequency           types.Int64    `tfsdk:"reauth_frequency"`
}

type securityIkePolicyConfig struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	Mode                      types.String `tfsdk:"mode"`
	PreSharedKeyHexa          types.String `tfsdk:"pre_shared_key_hexa"`
	PreSharedKeyHexaWO        types.String `tfsdk:"pre_shared_key_hexa_wo"`
	PreSharedKeyHexaWOVersion types.Int64  `tfsdk:"pre_shared_key_hexa_wo_version"`
	PreSharedKeyText          types.String `tfsdk:"pre_shared_key_text"`
	PreSharedKeyTextWO        types.String `tfsdk:"pre_shared_key_text_wo"`
	PreSharedKeyTextWOVersion types.Int64  `tfsdk:"pre_shared_key_text_wo_version"`
	Proposals                 types.List   `tfsdk:"proposals"`
	ProposalSet               types.String `tfsdk:"proposal_set"`
	ReauthFrequency           types.Int64  `tfsdk:"reauth_frequency"`
}

func (rsc *securityIkePolicy) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityIkePolicyConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Proposals.IsNull() && !config.Proposals.IsUnknown() &&
		!config.ProposalSet.IsNull() && !config.ProposalSet.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("proposals"),
			tfdiag.ConflictConfigErrSummary,
			"proposals and proposal_set cannot be configured together",
		)
	}
	// the pre-shared key can only be set once, whatever its format
	// and whether it's a write-only argument or not
	var preSharedKeyConfigured []string
	for _, preSharedKey := range []struct {
		name  string
		value types.String
	}{
		{name: "pre_shared_key_hexa", value: config.PreSharedKeyHexa},
		{name: "pre_shared_key_hexa_wo", value: config.PreSharedKeyHexaWO},
		{name: "pre_shared_key_text", value: config.PreSharedKeyText},
		{name: "pre_shared_key_text_wo", value: config.PreSharedKeyTextWO},
	} {
		if !preSharedKey.value.IsNull() && !preSharedKey.value.IsUnknown() {
			preSharedKeyConfigured = append(preSharedKeyConfigured, preSharedKey.name)
		}
	}
	if len(preSharedKeyConfigured) > 1 {
		resp.Diagnostics.AddAttributeError(
			path.Root(preSharedKeyConfigured[0]),
			tfdiag.ConflictConfigErrSummary,
			"only one of "+strings.Join(preSharedKeyConfigured, ", ")+" can be specified",
		)
	}
}

func (rsc *securityIkePolicy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityIkePolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(plan.getWriteOnly(ctx, req.Config)...)
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
			policyExists, err := checkSecurityIkePolicyExists(fnCtx, plan.Name.ValueString(), junSess)
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
			policyExists, err := checkSecurityIkePolicyExists(fnCtx, plan.Name.ValueString(), junSess)
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

func (rsc *securityIkePolicy) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityIkePolicyData
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
		func() {
			data.keepWriteOnly(&state)
		},
		resp,
	)
}

func (rsc *securityIkePolicy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityIkePolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(plan.getWriteOnly(ctx, req.Config)...)
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

func (rsc *securityIkePolicy) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityIkePolicyData
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

func (rsc *securityIkePolicy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityIkePolicyData

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

func checkSecurityIkePolicyExists(
	ctx context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"security ike policy \""+name+"\""+junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

// getWriteOnly read the write-only arguments from the configuration,
// their values aren't present in the plan or the state.
func (rscData *securityIkePolicyData) getWriteOnly(
	ctx context.Context, config tfsdk.Config,
) (diags diag.Diagnostics) {
	diags.Append(config.GetAttribute(ctx,
		path.Root("pre_shared_key_hexa_wo"), &rscData.PreSharedKeyHexaWO)...)
	diags.Append(config.GetAttribute(ctx,
		path.Root("pre_shared_key_text_wo"), &rscData.PreSharedKeyTextWO)...)

	return diags
}

// keepWriteOnly carry over the version arguments of the write-only arguments from the state,
// and don't read the secrets in the standard arguments when the write-only ones are used.
func (rscData *securityIkePolicyData) keepWriteOnly(state *securityIkePolicyData) {
	rscData.PreSharedKeyHexaWOVersion = state.PreSharedKeyHexaWOVersion
	rscData.PreSharedKeyTextWOVersion = state.PreSharedKeyTextWOVersion
	if !state.PreSharedKeyHexaWOVersion.IsNull() {
		rscData.PreSharedKeyHexa = types.StringNull()
	}
	if !state.PreSharedKeyTextWOVersion.IsNull() {
		rscData.PreSharedKeyText = types.StringNull()
	}
}

func (rscData *securityIkePolicyData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityIkePolicyData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityIkePolicyData) set(
	ctx context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)
	setPrefix := "set security ike policy \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.Mode.ValueString(); v != "" {
		if v != "main" && v != "aggressive" {
			return path.Root("mode"),
				fmt.Errorf("unknown ike mode %q", v)
		}
		configSet = append(configSet, setPrefix+"mode "+v)
	}
	for _, v := range rscData.Proposals {
		configSet = append(configSet, setPrefix+"proposals \""+v.ValueString()+"\"")
	}
	if v := rscData.ProposalSet.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"proposal-set "+v)
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := rscData.PreSharedKeyHexa.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pre-shared-key hexadecimal \""+v+"\"")
	} else if v := rscData.PreSharedKeyHexaWO.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pre-shared-key hexadecimal \""+v+"\"")
	}
	if v := rscData.PreSharedKeyText.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pre-shared-key ascii-text \""+v+"\"")
	} else if v := rscData.PreSharedKeyTextWO.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pre-shared-key ascii-text \""+v+"\"")
	}
	if !rscData.ReauthFrequency.IsNull() {
		configSet = append(configSet, setPrefix+"reauth-frequency "+
			utils.ConvI64toa(rscData.ReauthFrequency.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(ctx, configSet)
}

func (rscData *securityIkePolicyData) read(
	ctx context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"security ike policy \""+name+"\""+junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.fillID()
		for item := range strings.SplitSeq(showConfig, "\n") {
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
			case balt.CutPrefixInString(&itemTrim, "mode "):
				rscData.Mode = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "proposals "):
				rscData.Proposals = append(rscData.Proposals, types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "proposal-set "):
				rscData.ProposalSet = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "pre-shared-key hexadecimal "):
				rscData.PreSharedKeyHexa, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "pre-shared-key hexadecimal")
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "pre-shared-key ascii-text "):
				rscData.PreSharedKeyText, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "pre-shared-key ascii-text")
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "reauth-frequency "):
				rscData.ReauthFrequency, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *securityIkePolicyData) del(
	ctx context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security ike policy \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(ctx, configSet)
}
