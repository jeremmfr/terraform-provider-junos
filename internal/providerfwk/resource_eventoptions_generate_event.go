package providerfwk

import (
	"context"
	"regexp"
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
	_ resource.Resource                   = &eventoptionsGenerateEvent{}
	_ resource.ResourceWithConfigure      = &eventoptionsGenerateEvent{}
	_ resource.ResourceWithValidateConfig = &eventoptionsGenerateEvent{}
	_ resource.ResourceWithImportState    = &eventoptionsGenerateEvent{}
)

type eventoptionsGenerateEvent struct {
	client *junos.Client
}

func newEventoptionsGenerateEventResource() resource.Resource {
	return &eventoptionsGenerateEvent{}
}

func (rsc *eventoptionsGenerateEvent) typeName() string {
	return providerName + "_eventoptions_generate_event"
}

func (rsc *eventoptionsGenerateEvent) junosName() string {
	return "event-options generate-event"
}

func (rsc *eventoptionsGenerateEvent) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *eventoptionsGenerateEvent) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *eventoptionsGenerateEvent) Configure(
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

func (rsc *eventoptionsGenerateEvent) Schema(
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
				Description: "Name of the event to be generated.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"no_drift": schema.BoolAttribute{
				Optional:    true,
				Description: "Avoid event generation delay propagating to next event.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"start_time": schema.StringAttribute{
				Optional:    true,
				Description: "Start-time to generate event.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d{4}\-\d\d?\-\d\d?\.\d{2}:\d{2}:\d{2}$`),
						"must be in the format 'YYYY-MM-DD.HH:MM:SS'",
					),
				},
			},
			"time_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Frequency for generating the event.",
				Validators: []validator.Int64{
					int64validator.Between(60, 2592000),
				},
			},
			"time_of_day": schema.StringAttribute{
				Optional:    true,
				Description: "Time of day at which to generate event.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d{2}:\d{2}:\d{2}$`),
						"must be in the format 'HH:MM:SS'",
					),
				},
			},
		},
	}
}

type eventoptionsGenerateEventData struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	NoDrift      types.Bool   `tfsdk:"no_drift"`
	StartTime    types.String `tfsdk:"start_time"`
	TimeInterval types.Int64  `tfsdk:"time_interval"`
	TimeOfDay    types.String `tfsdk:"time_of_day"`
}

func (rsc *eventoptionsGenerateEvent) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config eventoptionsGenerateEventData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.TimeInterval.IsNull() &&
		config.TimeOfDay.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"one of time_interval or time_of_day must be specified",
		)
	}
	if !config.TimeInterval.IsNull() && !config.TimeInterval.IsUnknown() &&
		!config.TimeOfDay.IsNull() && !config.TimeOfDay.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("time_of_day"),
			tfdiag.ConflictConfigErrSummary,
			"only one of time_interval or time_of_day can be specified",
		)
	}
	if !config.StartTime.IsNull() && !config.StartTime.IsUnknown() {
		if config.TimeInterval.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("start_time"),
				tfdiag.MissingConfigErrSummary,
				"time_interval must be specified with start_time",
			)
		}
	}
}

func (rsc *eventoptionsGenerateEvent) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan eventoptionsGenerateEventData
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
			eventExists, err := checkEventoptionsGenerateEventExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if eventExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			eventExists, err := checkEventoptionsGenerateEventExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !eventExists {
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

func (rsc *eventoptionsGenerateEvent) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data eventoptionsGenerateEventData
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

func (rsc *eventoptionsGenerateEvent) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state eventoptionsGenerateEventData
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

func (rsc *eventoptionsGenerateEvent) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state eventoptionsGenerateEventData
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

func (rsc *eventoptionsGenerateEvent) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data eventoptionsGenerateEventData

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

func checkEventoptionsGenerateEventExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"event-options generate-event \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *eventoptionsGenerateEventData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *eventoptionsGenerateEventData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *eventoptionsGenerateEventData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set event-options generate-event \"" + rscData.Name.ValueString() + "\" "

	if rscData.NoDrift.ValueBool() {
		configSet = append(configSet, setPrefix+"no-drift")
	}
	if v := rscData.StartTime.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"start-time "+v)
	}
	if !rscData.TimeInterval.IsNull() {
		configSet = append(configSet, setPrefix+"time-interval "+
			utils.ConvI64toa(rscData.TimeInterval.ValueInt64()))
	}
	if v := rscData.TimeOfDay.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"time-of-day "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *eventoptionsGenerateEventData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"event-options generate-event \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case itemTrim == "no-drift":
				rscData.NoDrift = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "start-time "):
				rscData.StartTime = types.StringValue(strings.Split(strings.Trim(itemTrim, "\""), " ")[0])
			case balt.CutPrefixInString(&itemTrim, "time-interval "):
				rscData.TimeInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "time-of-day "):
				rscData.TimeOfDay = types.StringValue(strings.Split(strings.Trim(itemTrim, "\""), " ")[0])
			}
		}
	}

	return nil
}

func (rscData *eventoptionsGenerateEventData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete event-options generate-event \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
