package providerfwk

import (
	"context"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource                = &multichassis{}
	_ resource.ResourceWithConfigure   = &multichassis{}
	_ resource.ResourceWithModifyPlan  = &multichassis{}
	_ resource.ResourceWithImportState = &multichassis{}
)

type multichassis struct {
	client *junos.Client
}

func newMultichassisResource() resource.Resource {
	return &multichassis{}
}

func (rsc *multichassis) typeName() string {
	return providerName + "_multichassis"
}

func (rsc *multichassis) junosName() string {
	return "multi-chassis"
}

func (rsc *multichassis) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *multichassis) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *multichassis) Configure(
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

func (rsc *multichassis) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with value " +
					"`multichassis`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"clean_on_destroy": schema.BoolAttribute{
				Optional:    true,
				Description: "Clean entirely `" + rsc.junosName() + "` block when destroy this resource.",
			},
			"mc_lag_consistency_check": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Consistency Check.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"mc_lag_consistency_check_comparison_delay_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Time after which local and remote config are compared (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(5, 600),
				},
			},
		},
	}
}

type multichassisData struct {
	ID                                        types.String `tfsdk:"id"`
	CleanOnDestroy                            types.Bool   `tfsdk:"clean_on_destroy"`
	MCLagConsistencyCheck                     types.Bool   `tfsdk:"mc_lag_consistency_check"`
	MCLagConsistencyCheckComparaisonDelayTime types.Int64  `tfsdk:"mc_lag_consistency_check_comparison_delay_time"`
}

func (rsc *multichassis) ModifyPlan(
	ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var config, plan multichassisData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.MCLagConsistencyCheck.IsNull() {
		if config.MCLagConsistencyCheckComparaisonDelayTime.IsNull() {
			plan.MCLagConsistencyCheck = types.BoolNull()
		} else if !plan.MCLagConsistencyCheckComparaisonDelayTime.IsNull() &&
			!plan.MCLagConsistencyCheckComparaisonDelayTime.IsUnknown() {
			plan.MCLagConsistencyCheck = types.BoolValue(true)
		}
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func (rsc *multichassis) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan multichassisData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.MCLagConsistencyCheck.IsUnknown() {
		plan.MCLagConsistencyCheck = types.BoolNull()
		if !plan.MCLagConsistencyCheckComparaisonDelayTime.IsNull() {
			plan.MCLagConsistencyCheck = types.BoolValue(true)
		}
	}

	defaultResourceCreate(
		ctx,
		rsc,
		nil,
		nil,
		&plan,
		resp,
	)
}

func (rsc *multichassis) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data multichassisData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadWithoutArg = &data
	defaultResourceRead(
		ctx,
		rsc,
		nil,
		&data,
		func() {
			data.CleanOnDestroy = state.CleanOnDestroy
		},
		resp,
	)
}

func (rsc *multichassis) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state multichassisData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.MCLagConsistencyCheck.IsUnknown() {
		plan.MCLagConsistencyCheck = types.BoolNull()
		if !plan.MCLagConsistencyCheckComparaisonDelayTime.IsNull() {
			plan.MCLagConsistencyCheck = types.BoolValue(true)
		}
	}

	var _ resourceDataDelWithOpts = &state
	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *multichassis) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state multichassisData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.CleanOnDestroy.ValueBool() {
		defaultResourceDelete(
			ctx,
			rsc,
			&state,
			resp,
		)
	}
}

func (rsc *multichassis) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data multichassisData

	var _ resourceDataReadWithoutArg = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		"the `"+rsc.junosName()+"` block is not configured on device",
	)
}

func (rscData *multichassisData) fillID() {
	rscData.ID = types.StringValue("multichassis")
}

func (rscData *multichassisData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *multichassisData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set multi-chassis "

	configSet := []string{
		setPrefix,
	}

	if rscData.MCLagConsistencyCheck.ValueBool() {
		configSet = append(configSet, setPrefix+"mc-lag consistency-check")
	}
	if !rscData.MCLagConsistencyCheckComparaisonDelayTime.IsNull() {
		configSet = append(configSet, setPrefix+"mc-lag consistency-check comparison-delay-time "+
			utils.ConvI64toa(rscData.MCLagConsistencyCheckComparaisonDelayTime.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *multichassisData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"multi-chassis" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.fillID()
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "mc-lag consistency-check") {
				rscData.MCLagConsistencyCheck = types.BoolValue(true)
				if balt.CutPrefixInString(&itemTrim, " comparison-delay-time ") {
					rscData.MCLagConsistencyCheckComparaisonDelayTime, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (rscData *multichassisData) delOpts(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := "delete multi-chassis "

	configSet := []string{
		delPrefix + "mc-lag",
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *multichassisData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete multi-chassis",
	}

	return junSess.ConfigSet(configSet)
}
