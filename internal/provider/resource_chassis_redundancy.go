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
	_ resource.Resource                   = &chassisRedundancy{}
	_ resource.ResourceWithConfigure      = &chassisRedundancy{}
	_ resource.ResourceWithValidateConfig = &chassisRedundancy{}
	_ resource.ResourceWithImportState    = &chassisRedundancy{}
)

type chassisRedundancy struct {
	client *junos.Client
}

func newChassisRedundancyResource() resource.Resource {
	return &chassisRedundancy{}
}

func (rsc *chassisRedundancy) typeName() string {
	return providerName + "_chassis_redundancy"
}

func (rsc *chassisRedundancy) junosName() string {
	return "chassis redundancy"
}

func (rsc *chassisRedundancy) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *chassisRedundancy) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *chassisRedundancy) Configure(
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

func (rsc *chassisRedundancy) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Configure static configuration in `" + rsc.junosName() + "` block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with value " +
					"`redundancy`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"failover_disk_read_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "To failover, read threshold (ms) on disk underperform monitoring.",
				Validators: []validator.Int64{
					int64validator.Between(1000, 10000),
				},
			},
			"failover_disk_write_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "To failover, write threshold (ms) on disk underperform monitoring.",
				Validators: []validator.Int64{
					int64validator.Between(1000, 10000),
				},
			},
			"failover_not_on_disk_underperform": schema.BoolAttribute{
				Optional:    true,
				Description: "Prevent gstatd from initiating failovers in response to slow disks.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"failover_on_disk_failure": schema.BoolAttribute{
				Optional:    true,
				Description: "Failover on disk failure.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"failover_on_loss_of_keepalives": schema.BoolAttribute{
				Optional:    true,
				Description: "Failover on loss of keepalives.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"graceful_switchover": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable graceful switchover on supported hardware.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"keepalive_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Time before Routing Engine failover (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(2, 10000),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"routing_engine": schema.SetNestedBlock{
				Description: "For each slot, redundancy options.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"slot": schema.Int64Attribute{
							Required:    true,
							Description: "Routing Engine slot number.",
							Validators: []validator.Int64{
								int64validator.Between(0, 1),
							},
						},
						"role": schema.StringAttribute{
							Required:    true,
							Description: "Define role.",
							Validators: []validator.String{
								stringvalidator.OneOf("backup", "disabled", "master"),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(2),
				},
			},
		},
	}
}

type chassisRedundancyData struct {
	ID                            types.String                          `tfsdk:"id"`
	FailoverDiskReadThreshold     types.Int64                           `tfsdk:"failover_disk_read_threshold"`
	FailoverDiskWriteThreshold    types.Int64                           `tfsdk:"failover_disk_write_threshold"`
	FailoverNotOnDiskUnderperform types.Bool                            `tfsdk:"failover_not_on_disk_underperform"`
	FailoverOnDiskFailure         types.Bool                            `tfsdk:"failover_on_disk_failure"`
	FailoverOnLossOfKeepalives    types.Bool                            `tfsdk:"failover_on_loss_of_keepalives"`
	GracefulSwitchover            types.Bool                            `tfsdk:"graceful_switchover"`
	KeepaliveTime                 types.Int64                           `tfsdk:"keepalive_time"`
	RoutingEngine                 []chassisRedundancyBlockRoutingEngine `tfsdk:"routing_engine"`
}

type chassisRedundancyConfig struct {
	ID                            types.String `tfsdk:"id"`
	FailoverDiskReadThreshold     types.Int64  `tfsdk:"failover_disk_read_threshold"`
	FailoverDiskWriteThreshold    types.Int64  `tfsdk:"failover_disk_write_threshold"`
	FailoverNotOnDiskUnderperform types.Bool   `tfsdk:"failover_not_on_disk_underperform"`
	FailoverOnDiskFailure         types.Bool   `tfsdk:"failover_on_disk_failure"`
	FailoverOnLossOfKeepalives    types.Bool   `tfsdk:"failover_on_loss_of_keepalives"`
	GracefulSwitchover            types.Bool   `tfsdk:"graceful_switchover"`
	KeepaliveTime                 types.Int64  `tfsdk:"keepalive_time"`
	RoutingEngine                 types.Set    `tfsdk:"routing_engine"`
}

type chassisRedundancyBlockRoutingEngine struct {
	Slot types.Int64  `tfsdk:"slot" tfdata:"identifier"`
	Role types.String `tfsdk:"role"`
}

func (rsc *chassisRedundancy) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config chassisRedundancyConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.RoutingEngine.IsNull() &&
		!config.RoutingEngine.IsUnknown() {
		var configRoutingEngine []chassisRedundancyBlockRoutingEngine
		asDiags := config.RoutingEngine.ElementsAs(ctx, &configRoutingEngine, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		routingEngineSlot := make(map[int64]struct{})
		for _, block := range configRoutingEngine {
			if block.Slot.IsUnknown() {
				continue
			}

			slot := block.Slot.ValueInt64()
			if _, ok := routingEngineSlot[slot]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("routing_engine"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple routing_engine blocks with the same slot %d", slot),
				)
			}
			routingEngineSlot[slot] = struct{}{}
		}
	}
}

func (rsc *chassisRedundancy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan chassisRedundancyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
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

func (rsc *chassisRedundancy) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data chassisRedundancyData
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
		nil,
		resp,
	)
}

func (rsc *chassisRedundancy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state chassisRedundancyData
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

func (rsc *chassisRedundancy) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state chassisRedundancyData
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

func (rsc *chassisRedundancy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data chassisRedundancyData

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

func (rscData *chassisRedundancyData) fillID() {
	rscData.ID = types.StringValue("redundancy")
}

func (rscData *chassisRedundancyData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *chassisRedundancyData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set chassis redundancy "
	configSet := make([]string, 0, 100)

	if !rscData.FailoverDiskReadThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"failover disk-read-threshold "+
			utils.ConvI64toa(rscData.FailoverDiskReadThreshold.ValueInt64()))
	}
	if !rscData.FailoverDiskWriteThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"failover disk-write-threshold "+
			utils.ConvI64toa(rscData.FailoverDiskWriteThreshold.ValueInt64()))
	}
	if rscData.FailoverNotOnDiskUnderperform.ValueBool() {
		configSet = append(configSet, setPrefix+"failover not-on-disk-underperform")
	}
	if rscData.FailoverOnDiskFailure.ValueBool() {
		configSet = append(configSet, setPrefix+"failover on-disk-failure")
	}
	if rscData.FailoverOnLossOfKeepalives.ValueBool() {
		configSet = append(configSet, setPrefix+"failover on-loss-of-keepalives")
	}
	if rscData.GracefulSwitchover.ValueBool() {
		configSet = append(configSet, setPrefix+"graceful-switchover")
	}
	if !rscData.KeepaliveTime.IsNull() {
		configSet = append(configSet, setPrefix+"keepalive-time "+
			utils.ConvI64toa(rscData.KeepaliveTime.ValueInt64()))
	}

	routingEngineSlot := make(map[int64]struct{})
	for _, v := range rscData.RoutingEngine {
		slot := v.Slot.ValueInt64()
		if _, ok := routingEngineSlot[slot]; ok {
			return path.Root("routing_engine"),
				fmt.Errorf("multiple routing_engine blocks with the same slot %d", slot)
		}
		routingEngineSlot[slot] = struct{}{}

		configSet = append(configSet,
			setPrefix+"routing-engine "+utils.ConvI64toa(slot)+" "+v.Role.ValueString())
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *chassisRedundancyData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"chassis redundancy" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	rscData.fillID()
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "failover disk-read-threshold "):
				rscData.FailoverDiskReadThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "failover disk-write-threshold "):
				rscData.FailoverDiskWriteThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "failover not-on-disk-underperform":
				rscData.FailoverNotOnDiskUnderperform = types.BoolValue(true)
			case itemTrim == "failover on-disk-failure":
				rscData.FailoverOnDiskFailure = types.BoolValue(true)
			case itemTrim == "failover on-loss-of-keepalives":
				rscData.FailoverOnLossOfKeepalives = types.BoolValue(true)
			case itemTrim == "graceful-switchover":
				rscData.GracefulSwitchover = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "keepalive-time "):
				rscData.KeepaliveTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "routing-engine "):
				slotStr := tfdata.FirstElementOfJunosLine(itemTrim)
				slot, err := tfdata.ConvAtoi64Value(slotStr)
				if err != nil {
					return err
				}
				var routingEngine chassisRedundancyBlockRoutingEngine
				rscData.RoutingEngine, routingEngine = tfdata.ExtractBlock(rscData.RoutingEngine, slot)

				if balt.CutPrefixInString(&itemTrim, slotStr+" ") {
					routingEngine.Role = types.StringValue(itemTrim)
				}
				rscData.RoutingEngine = append(rscData.RoutingEngine, routingEngine)
			}
		}
	}

	return nil
}

func (rscData *chassisRedundancyData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete chassis redundancy",
	}

	return junSess.ConfigSet(configSet)
}
