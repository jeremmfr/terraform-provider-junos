package providerfwk

import (
	"context"
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
	_ resource.Resource                = &securityIdpCustomAttackGroup{}
	_ resource.ResourceWithConfigure   = &securityIdpCustomAttackGroup{}
	_ resource.ResourceWithImportState = &securityIdpCustomAttackGroup{}
)

type securityIdpCustomAttackGroup struct {
	client *junos.Client
}

func newSecurityIdpCustomAttackGroupResource() resource.Resource {
	return &securityIdpCustomAttackGroup{}
}

func (rsc *securityIdpCustomAttackGroup) typeName() string {
	return providerName + "_security_idp_custom_attack_group"
}

func (rsc *securityIdpCustomAttackGroup) junosName() string {
	return "security idp custom-attack-group"
}

func (rsc *securityIdpCustomAttackGroup) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityIdpCustomAttackGroup) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityIdpCustomAttackGroup) Configure(
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

func (rsc *securityIdpCustomAttackGroup) Schema(
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
				Description: "Custom attack group name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"member": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "List of attacks/attack groups belonging to this group.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.NoNullValues(),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 250),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
		},
	}
}

type securityIdpCustomAttackGroupData struct {
	ID     types.String   `tfsdk:"id"`
	Name   types.String   `tfsdk:"name"`
	Member []types.String `tfsdk:"member"`
}

func (rsc *securityIdpCustomAttackGroup) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityIdpCustomAttackGroupData
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
			if !junSess.CheckCompatibilitySecurity() {
				resp.Diagnostics.AddError(
					tfdiag.CompatibilityErrSummary,
					rsc.junosName()+junSess.SystemInformation.NotCompatibleMsg(),
				)

				return false
			}
			groupExists, err := checkSecurityIdpCustomAttackGroupExists(fnCtx, plan.Name.ValueString(), junSess)
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
			groupExists, err := checkSecurityIdpCustomAttackGroupExists(fnCtx, plan.Name.ValueString(), junSess)
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

func (rsc *securityIdpCustomAttackGroup) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityIdpCustomAttackGroupData
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

func (rsc *securityIdpCustomAttackGroup) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityIdpCustomAttackGroupData
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

func (rsc *securityIdpCustomAttackGroup) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityIdpCustomAttackGroupData
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

func (rsc *securityIdpCustomAttackGroup) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityIdpCustomAttackGroupData

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

func checkSecurityIdpCustomAttackGroupExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security idp custom-attack-group \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityIdpCustomAttackGroupData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityIdpCustomAttackGroupData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityIdpCustomAttackGroupData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set security idp custom-attack-group \"" + rscData.Name.ValueString() + "\" "

	configSet := []string{
		setPrefix,
	}

	for _, v := range rscData.Member {
		configSet = append(configSet, setPrefix+"group-members \""+v.ValueString()+"\"")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityIdpCustomAttackGroupData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security idp custom-attack-group \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			if balt.CutPrefixInString(&itemTrim, "group-members ") {
				rscData.Member = append(rscData.Member, types.StringValue(strings.Trim(itemTrim, "\"")))
			}
		}
	}

	return nil
}

func (rscData *securityIdpCustomAttackGroupData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security idp custom-attack-group \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
