package provider

import (
	"context"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource                   = &policyoptionsDamping{}
	_ resource.ResourceWithConfigure      = &policyoptionsDamping{}
	_ resource.ResourceWithValidateConfig = &policyoptionsDamping{}
	_ resource.ResourceWithImportState    = &policyoptionsDamping{}
)

type policyoptionsDamping struct {
	client *junos.Client
}

func newPolicyoptionsDampingResource() resource.Resource {
	return &policyoptionsDamping{}
}

func (rsc *policyoptionsDamping) typeName() string {
	return providerName + "_policyoptions_damping"
}

func (rsc *policyoptionsDamping) junosName() string {
	return "policy-options damping"
}

func (rsc *policyoptionsDamping) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *policyoptionsDamping) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *policyoptionsDamping) Configure(
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

func (rsc *policyoptionsDamping) Schema(
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
				Description: "Name to identify route flap damping parameters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable damping.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"half_life": schema.Int64Attribute{
				Optional:    true,
				Description: "Decay half-life (minutes).",
				Validators: []validator.Int64{
					int64validator.Between(1, 45),
				},
			},
			"max_suppress": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum hold-down time (minutes).",
				Validators: []validator.Int64{
					int64validator.Between(1, 720),
				},
			},
			"reuse": schema.Int64Attribute{
				Optional:    true,
				Description: "Reuse threshold.",
				Validators: []validator.Int64{
					int64validator.Between(1, 20000),
				},
			},
			"suppress": schema.Int64Attribute{
				Optional:    true,
				Description: "Cutoff threshold.",
				Validators: []validator.Int64{
					int64validator.Between(1, 20000),
				},
			},
		},
	}
}

type policyoptionsDampingData struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Disable     types.Bool   `tfsdk:"disable"`
	HalfLife    types.Int64  `tfsdk:"half_life"`
	MaxSuppress types.Int64  `tfsdk:"max_suppress"`
	Reuse       types.Int64  `tfsdk:"reuse"`
	Suppress    types.Int64  `tfsdk:"suppress"`
}

func (rsc *policyoptionsDamping) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config policyoptionsDampingData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.HalfLife.IsNull() && !config.HalfLife.IsUnknown() &&
		!config.MaxSuppress.IsNull() && !config.MaxSuppress.IsUnknown() {
		if config.MaxSuppress.ValueInt64() <= config.HalfLife.ValueInt64() {
			resp.Diagnostics.AddAttributeError(
				path.Root("max_suppress"),
				tfdiag.ConflictConfigErrSummary,
				"max_suppress must be greater than half_life",
			)
		}
	}
}

func (rsc *policyoptionsDamping) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan policyoptionsDampingData
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
			dampingExists, err := checkPolicyoptionsDampingExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if dampingExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			dampingExists, err := checkPolicyoptionsDampingExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !dampingExists {
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

func (rsc *policyoptionsDamping) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data policyoptionsDampingData
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

func (rsc *policyoptionsDamping) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state policyoptionsDampingData
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

func (rsc *policyoptionsDamping) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state policyoptionsDampingData
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

func (rsc *policyoptionsDamping) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data policyoptionsDampingData

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

func checkPolicyoptionsDampingExists(
	ctx context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"policy-options damping \""+name+"\""+junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *policyoptionsDampingData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *policyoptionsDampingData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *policyoptionsDampingData) set(
	ctx context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set policy-options damping \"" + rscData.Name.ValueString() + "\" "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if rscData.Disable.ValueBool() {
		configSet = append(configSet, setPrefix+"disable")
	}
	if !rscData.HalfLife.IsNull() {
		configSet = append(configSet, setPrefix+"half-life "+
			utils.ConvI64toa(rscData.HalfLife.ValueInt64()))
	}
	if !rscData.MaxSuppress.IsNull() {
		configSet = append(configSet, setPrefix+"max-suppress "+
			utils.ConvI64toa(rscData.MaxSuppress.ValueInt64()))
	}
	if !rscData.Reuse.IsNull() {
		configSet = append(configSet, setPrefix+"reuse "+
			utils.ConvI64toa(rscData.Reuse.ValueInt64()))
	}
	if !rscData.Suppress.IsNull() {
		configSet = append(configSet, setPrefix+"suppress "+
			utils.ConvI64toa(rscData.Suppress.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(ctx, configSet)
}

func (rscData *policyoptionsDampingData) read(
	ctx context.Context, name string, junSess *junos.Session,
) (
	err error,
) {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"policy-options damping \""+name+"\""+junos.PipeDisplaySetRelative)
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
			case itemTrim == "disable":
				rscData.Disable = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "half-life "):
				rscData.HalfLife, err = tfdata.ConvAtoi64Value(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "max-suppress "):
				rscData.MaxSuppress, err = tfdata.ConvAtoi64Value(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "reuse "):
				rscData.Reuse, err = tfdata.ConvAtoi64Value(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "suppress "):
				rscData.Suppress, err = tfdata.ConvAtoi64Value(itemTrim)
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (rscData *policyoptionsDampingData) del(
	ctx context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete policy-options damping \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(ctx, configSet)
}
