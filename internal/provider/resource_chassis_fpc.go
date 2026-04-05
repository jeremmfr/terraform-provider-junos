package provider

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &chassisFpc{}
	_ resource.ResourceWithConfigure      = &chassisFpc{}
	_ resource.ResourceWithValidateConfig = &chassisFpc{}
	_ resource.ResourceWithImportState    = &chassisFpc{}
)

type chassisFpc struct {
	client *junos.Client
}

func newChassisFpcResource() resource.Resource {
	return &chassisFpc{}
}

func (rsc *chassisFpc) typeName() string {
	return providerName + "_chassis_fpc"
}

func (rsc *chassisFpc) junosName() string {
	return "chassis fpc"
}

func (rsc *chassisFpc) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *chassisFpc) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *chassisFpc) Configure(
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

func (rsc *chassisFpc) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<slot_number>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"slot_number": schema.Int64Attribute{
				Required:    true,
				Description: "FPC number.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"cfp_to_et": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable ET interface and remove CFP client",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"sampling_instance": schema.StringAttribute{
				Optional:    true,
				Description: "Name for sampling instance.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"error": schema.SingleNestedBlock{
				Description: "Error level configuration for FPC.",
				Attributes: map[string]schema.Attribute{
					"fatal_action": schema.StringAttribute{
						Optional:    true,
						Description: "Configure the action for fatal level.",
						Validators: []validator.String{
							stringvalidator.OneOf("alarm", "disable-pfe", "get-state", "log", "offline", "reset", "trap"),
						},
					},
					"fatal_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Error count at which to take the action.",
						Validators: []validator.Int64{
							int64validator.Between(1, 1024),
						},
					},
					"major_action": schema.StringAttribute{
						Optional:    true,
						Description: "Configure the action for major level.",
						Validators: []validator.String{
							stringvalidator.OneOf("alarm", "disable-pfe", "get-state", "log", "offline", "reset", "trap"),
						},
					},
					"major_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Error count at which to take the action.",
						Validators: []validator.Int64{
							int64validator.Between(1, 1024),
						},
					},
					"minor_action": schema.StringAttribute{
						Optional:    true,
						Description: "Configure the action for minor level.",
						Validators: []validator.String{
							stringvalidator.OneOf("alarm", "disable-pfe", "get-state", "log", "offline", "reset", "trap"),
						},
					},
					"minor_threshold": schema.Int64Attribute{
						Optional:    true,
						Description: "Error count at which to take the action.",
						Validators: []validator.Int64{
							int64validator.Between(0, 1024),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
		},
	}
}

type chassisFpcData struct {
	ID               types.String          `tfsdk:"id"                tfdata:"skip_isempty"`
	SlotNumber       types.Int64           `tfsdk:"slot_number"       tfdata:"skip_isempty"`
	CfpToEt          types.Bool            `tfsdk:"cfp_to_et"`
	Error            *chassisFpcBlockError `tfsdk:"error"`
	SamplingInstance types.String          `tfsdk:"sampling_instance"`
}

func (rscData *chassisFpcData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

type chassisFpcBlockError struct {
	FatalAction    types.String `tfsdk:"fatal_action"`
	FatalThreshold types.Int64  `tfsdk:"fatal_threshold"`
	MajorAction    types.String `tfsdk:"major_action"`
	MajorThreshold types.Int64  `tfsdk:"major_threshold"`
	MinorAction    types.String `tfsdk:"minor_action"`
	MinorThreshold types.Int64  `tfsdk:"minor_threshold"`
}

func (block *chassisFpcBlockError) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (rsc *chassisFpc) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config chassisFpcData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("slot_number"),
			tfdiag.MissingConfigErrSummary,
			"at least one of arguments need to be set (in addition to `slot_number`)",
		)
	}

	if config.Error != nil {
		if config.Error.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("error").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"error block is empty",
			)
		}
	}
}

func (rsc *chassisFpc) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan chassisFpcData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			slotNumberExists, err := checkChassisFpcExists(fnCtx, plan.SlotNumber.ValueInt64(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if slotNumberExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf(rsc.junosName()+" %d already exists", plan.SlotNumber.ValueInt64()),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			slotNumberExists, err := checkChassisFpcExists(fnCtx, plan.SlotNumber.ValueInt64(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !slotNumberExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					fmt.Sprintf(rsc.junosName()+" %d does not exists after commit "+
						"=> check your config", plan.SlotNumber.ValueInt64()),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *chassisFpc) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data chassisFpcData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1Int64 = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.SlotNumber.ValueInt64(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *chassisFpc) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state chassisFpcData
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

func (rsc *chassisFpc) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state chassisFpcData
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

func (rsc *chassisFpc) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data chassisFpcData

	var _ resourceDataReadFrom1Int64 = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		"",
	)
}

func checkChassisFpcExists(
	ctx context.Context, slotNumber int64, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"chassis fpc "+utils.ConvI64toa(slotNumber)+junos.PipeDisplaySetRelative)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}
	if slices.ContainsFunc(strings.Split(showConfig, "\n"), func(item string) bool {
		itemTrim := strings.TrimPrefix(item, junos.SetLS)
		switch {
		case itemTrim == "cfp-to-et":
			return true
		case strings.HasPrefix(itemTrim, "error fatal "):
			return true
		case strings.HasPrefix(itemTrim, "error major "):
			return true
		case strings.HasPrefix(itemTrim, "error minor "):
			return true
		case strings.HasPrefix(itemTrim, "sampling-instance "):
			return true
		default:
			return false
		}
	}) {
		return true, nil
	}

	return false, nil
}

func (rscData *chassisFpcData) fillID() {
	rscData.ID = types.StringValue(utils.ConvI64toa(rscData.SlotNumber.ValueInt64()))
}

func (rscData *chassisFpcData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *chassisFpcData) set(
	ctx context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	if rscData.isEmpty() {
		return path.Root("slot_number"),
			errors.New("at least one of arguments need to be set (in addition to `slot_number`)")
	}

	configSet := make([]string, 0, 100)
	setPrefix := "set chassis fpc " + utils.ConvI64toa(rscData.SlotNumber.ValueInt64()) + " "

	if rscData.CfpToEt.ValueBool() {
		configSet = append(configSet, setPrefix+"cfp-to-et")
	}
	if v := rscData.SamplingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"sampling-instance \""+v+"\"")
	}

	if rscData.Error != nil {
		blockSet, pathErr, err := rscData.Error.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(ctx, configSet)
}

func (block *chassisFpcBlockError) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	if block.isEmpty() {
		return nil, path.Root("error").AtName("*"),
			errors.New("error block is empty")
	}

	configSet := make([]string, 0, 100)
	setPrefix += "error "

	if v := block.FatalAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"fatal action "+v)
	}
	if !block.FatalThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"fatal threshold "+
			utils.ConvI64toa(block.FatalThreshold.ValueInt64()))
	}
	if v := block.MajorAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"major action "+v)
	}
	if !block.MajorThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"major threshold "+
			utils.ConvI64toa(block.MajorThreshold.ValueInt64()))
	}
	if v := block.MinorAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"minor action "+v)
	}
	if !block.MinorThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"minor threshold "+
			utils.ConvI64toa(block.MinorThreshold.ValueInt64()))
	}

	return configSet, path.Empty(), nil
}

func (rscData *chassisFpcData) read(
	ctx context.Context, slotNumber int64, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"chassis fpc "+utils.ConvI64toa(slotNumber)+junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		for item := range strings.SplitSeq(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case itemTrim == "cfp-to-et":
				rscData.CfpToEt = types.BoolValue(true)
			case strings.HasPrefix(itemTrim, "error fatal "),
				strings.HasPrefix(itemTrim, "error major "),
				strings.HasPrefix(itemTrim, "error minor "):
				if rscData.Error == nil {
					rscData.Error = &chassisFpcBlockError{}
				}
				if err := rscData.Error.read(strings.TrimPrefix(itemTrim, "error ")); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "sampling-instance "):
				rscData.SamplingInstance = types.StringValue(strings.Trim(itemTrim, "\""))
			}
		}
	}
	if !rscData.CfpToEt.IsNull() ||
		rscData.Error != nil ||
		!rscData.SamplingInstance.IsNull() {
		rscData.SlotNumber = types.Int64Value(slotNumber)
		rscData.fillID()
	}

	return nil
}

func (block *chassisFpcBlockError) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "fatal action "):
		block.FatalAction = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "fatal threshold "):
		block.FatalThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "major action "):
		block.MajorAction = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "major threshold "):
		block.MajorThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "minor action "):
		block.MinorAction = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "minor threshold "):
		block.MinorThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (rscData *chassisFpcData) del(
	ctx context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS + "chassis fpc " + utils.ConvI64toa(rscData.SlotNumber.ValueInt64()) + " "

	configSet := []string{
		delPrefix + "cfp-to-et",
		delPrefix + "error fatal",
		delPrefix + "error major",
		delPrefix + "error minor",
		delPrefix + "sampling-instance \"" + rscData.SamplingInstance.ValueString() + "\"",
	}

	return junSess.ConfigSet(ctx, configSet)
}
