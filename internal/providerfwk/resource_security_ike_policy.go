package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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

func (rsc *securityIkePolicy) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
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
		Description: "Provides a " + rsc.junosName() + ".",
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
				Description: "IKE mode for Phase 1.",
				PlanModifiers: []planmodifier.String{
					tfplanmodifier.StringDefault("main"),
				},
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
			"pre_shared_key_text": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Preshared key wit format as text.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
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
	ID               types.String   `tfsdk:"id"`
	Name             types.String   `tfsdk:"name"`
	Description      types.String   `tfsdk:"description"`
	Mode             types.String   `tfsdk:"mode"`
	PreSharedKeyHexa types.String   `tfsdk:"pre_shared_key_hexa"`
	PreSharedKeyText types.String   `tfsdk:"pre_shared_key_text"`
	Proposals        []types.String `tfsdk:"proposals"`
	ProposalSet      types.String   `tfsdk:"proposal_set"`
	ReauthFrequency  types.Int64    `tfsdk:"reauth_frequency"`
}

type securityIkePolicyConfig struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Mode             types.String `tfsdk:"mode"`
	PreSharedKeyHexa types.String `tfsdk:"pre_shared_key_hexa"`
	PreSharedKeyText types.String `tfsdk:"pre_shared_key_text"`
	Proposals        types.List   `tfsdk:"proposals"`
	ProposalSet      types.String `tfsdk:"proposal_set"`
	ReauthFrequency  types.Int64  `tfsdk:"reauth_frequency"`
}

func (rsc *securityIkePolicy) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityIkePolicyConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Proposals.IsNull() && !config.ProposalSet.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("proposals"),
			"Conflict Configuration Error",
			"only one of proposals or proposal_set must be specified",
		)
	}
	if !config.PreSharedKeyText.IsNull() && !config.PreSharedKeyHexa.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("pre_shared_key_text"),
			"Conflict Configuration Error",
			"only one of pre_shared_key_text or pre_shared_key_hexa can be specified",
		)
	}
}

func (rsc *securityIkePolicy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityIkePolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Empty Name",
			"could not create "+rsc.junosName()+" with empty name",
		)

		return
	}

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		resp.Diagnostics.AddError(
			"Compatibility Error",
			fmt.Sprintf(rsc.junosName()+" not compatible "+
				"with Junos device %q", junSess.SystemInformation.HardwareModel),
		)

		return
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}
	policyExists, err := checkSecurityIkePolicyExists(ctx, plan.Name.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Pre Check Error", err.Error())

		return
	}
	if policyExists {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError(
			"Duplicate Configuration Error",
			fmt.Sprintf(rsc.junosName()+" %q already exists", plan.Name.ValueString()),
		)

		return
	}

	if errPath, err := plan.set(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("create resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	policyExists, err = checkSecurityIkePolicyExists(ctx, plan.Name.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Post Check Error", err.Error())

		return
	}
	if !policyExists {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf(rsc.junosName()+" %q not exists after commit "+
				"=> check your config", plan.Name.ValueString()),
		)

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityIkePolicy) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityIkePolicyData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	err = data.read(ctx, state.Name.ValueString(), junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}
	if data.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (rsc *securityIkePolicy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityIkePolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("update resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *securityIkePolicy) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityIkePolicyData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeDeleteAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	warns, err := junSess.CommitConf("delete resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.Append(tfdiag.Warns("Config Clear Warning", junSess.ConfigClear())...)
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}
}

func (rsc *securityIkePolicy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	var data securityIkePolicyData
	if err := data.read(ctx, req.ID, junSess); err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf("don't find "+rsc.junosName()+" with id %q "+
				"(id must be <name>)", req.ID),
		)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkSecurityIkePolicyExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig + "security ike policy \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityIkePolicyData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityIkePolicyData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)

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
	}
	if v := rscData.PreSharedKeyText.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pre-shared-key ascii-text \""+v+"\"")
	}
	if !rscData.ReauthFrequency.IsNull() {
		configSet = append(configSet, setPrefix+"reauth-frequency "+utils.ConvI64toa(rscData.ReauthFrequency.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityIkePolicyData) read(
	_ context.Context, name string, junSess *junos.Session,
) (
	err error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security ike policy \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "mode "):
				rscData.Mode = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "proposals "):
				rscData.Proposals = append(rscData.Proposals, types.StringValue(strings.Trim(itemTrim, "\"")))
			case balt.CutPrefixInString(&itemTrim, "proposal-set "):
				rscData.ProposalSet = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "pre-shared-key hexadecimal "):
				rscData.PreSharedKeyHexa, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "pre-shared-key hexadecimal")
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "pre-shared-key ascii-text "):
				rscData.PreSharedKeyText, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "pre-shared-key ascii-text")
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
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security ike policy \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
