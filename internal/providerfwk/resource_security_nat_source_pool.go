package providerfwk

import (
	"context"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
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
	_ resource.Resource                   = &securityNatSourcePool{}
	_ resource.ResourceWithConfigure      = &securityNatSourcePool{}
	_ resource.ResourceWithValidateConfig = &securityNatSourcePool{}
	_ resource.ResourceWithImportState    = &securityNatSourcePool{}
)

type securityNatSourcePool struct {
	client *junos.Client
}

func newSecurityNatSourcePoolResource() resource.Resource {
	return &securityNatSourcePool{}
}

func (rsc *securityNatSourcePool) typeName() string {
	return providerName + "_security_nat_source_pool"
}

func (rsc *securityNatSourcePool) junosName() string {
	return "security nat source pool"
}

func (rsc *securityNatSourcePool) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityNatSourcePool) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityNatSourcePool) Configure(
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

func (rsc *securityNatSourcePool) Schema(
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
				Description: "Pool name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"address": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "CIDR address to source nat pool.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						tfvalidator.StringCIDR(),
					),
				},
			},
			"address_pooling": schema.StringAttribute{
				Optional:    true,
				Description: "Type of address pooling.",
				Validators: []validator.String{
					stringvalidator.OneOf("no-paired", "paired"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of pool.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"pool_utilization_alarm_clear_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "Lower threshold at which an SNMP trap is triggered.",
				Validators: []validator.Int64{
					int64validator.Between(40, 100),
				},
			},
			"pool_utilization_alarm_raise_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "Upper threshold at which an SNMP trap is triggered.",
				Validators: []validator.Int64{
					int64validator.Between(50, 100),
				},
			},
			"port_no_translation": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not perform port translation.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"port_overloading_factor": schema.Int64Attribute{
				Optional:    true,
				Description: "Port overloading factor for each IP.",
				Validators: []validator.Int64{
					int64validator.Between(2, 32),
				},
			},
			"port_range": schema.StringAttribute{
				Optional:    true,
				Description: "Range of port to source nat.",
				Validators: []validator.String{
					tfvalidator.StringNumberRange(1024, 65535).WithNameInError("Nat Source Port"),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Description: "Name of routing instance to switch instance with nat.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
		},
	}
}

type securityNatSourcePoolData struct {
	ID                                 types.String   `tfsdk:"id"`
	Name                               types.String   `tfsdk:"name"`
	Address                            []types.String `tfsdk:"address"`
	AddressPooling                     types.String   `tfsdk:"address_pooling"`
	Description                        types.String   `tfsdk:"description"`
	PoolUtilizationAlarmClearThreshold types.Int64    `tfsdk:"pool_utilization_alarm_clear_threshold"`
	PoolUtilizationAlarmRaiseThreshold types.Int64    `tfsdk:"pool_utilization_alarm_raise_threshold"`
	PortNoTranslation                  types.Bool     `tfsdk:"port_no_translation"`
	PortOverloadingFactor              types.Int64    `tfsdk:"port_overloading_factor"`
	PortRange                          types.String   `tfsdk:"port_range"`
	RoutingInstance                    types.String   `tfsdk:"routing_instance"`
}

type securityNatSourcePoolConfig struct {
	ID                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	Address                            types.List   `tfsdk:"address"`
	AddressPooling                     types.String `tfsdk:"address_pooling"`
	Description                        types.String `tfsdk:"description"`
	PoolUtilizationAlarmClearThreshold types.Int64  `tfsdk:"pool_utilization_alarm_clear_threshold"`
	PoolUtilizationAlarmRaiseThreshold types.Int64  `tfsdk:"pool_utilization_alarm_raise_threshold"`
	PortNoTranslation                  types.Bool   `tfsdk:"port_no_translation"`
	PortOverloadingFactor              types.Int64  `tfsdk:"port_overloading_factor"`
	PortRange                          types.String `tfsdk:"port_range"`
	RoutingInstance                    types.String `tfsdk:"routing_instance"`
}

func (rsc *securityNatSourcePool) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityNatSourcePoolConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.PoolUtilizationAlarmClearThreshold.IsNull() &&
		config.PoolUtilizationAlarmRaiseThreshold.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("pool_utilization_alarm_clear_threshold"),
			tfdiag.MissingConfigErrSummary,
			"pool_utilization_alarm_raise_threshold must be specified with pool_utilization_alarm_clear_threshold",
		)
	}
	if !config.PortNoTranslation.IsNull() && !config.PortNoTranslation.IsUnknown() &&
		!config.PortOverloadingFactor.IsNull() && !config.PortOverloadingFactor.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("port_no_translation"),
			tfdiag.ConflictConfigErrSummary,
			"port_no_translation and port_overloading_factor cannot be configured together",
		)
	}
	if !config.PortNoTranslation.IsNull() && !config.PortNoTranslation.IsUnknown() &&
		!config.PortRange.IsNull() && !config.PortRange.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("port_no_translation"),
			tfdiag.ConflictConfigErrSummary,
			"port_no_translation and port_range cannot be configured together",
		)
	}
	if !config.PortOverloadingFactor.IsNull() && !config.PortOverloadingFactor.IsUnknown() &&
		!config.PortRange.IsNull() && !config.PortRange.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("port_overloading_factor"),
			tfdiag.ConflictConfigErrSummary,
			"port_overloading_factor and port_range cannot be configured together",
		)
	}
}

func (rsc *securityNatSourcePool) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityNatSourcePoolData
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
			poolExists, err := checkSecurityNatSourcePoolExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if poolExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			poolExists, err := checkSecurityNatSourcePoolExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if !poolExists {
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

func (rsc *securityNatSourcePool) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityNatSourcePoolData
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

func (rsc *securityNatSourcePool) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityNatSourcePoolData
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

func (rsc *securityNatSourcePool) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityNatSourcePoolData
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

func (rsc *securityNatSourcePool) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityNatSourcePoolData

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

func checkSecurityNatSourcePoolExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat source pool " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityNatSourcePoolData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityNatSourcePoolData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityNatSourcePoolData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security nat source pool " + rscData.Name.ValueString() + " "

	for _, v := range rscData.Address {
		configSet = append(configSet, setPrefix+"address "+v.ValueString())
	}
	if v := rscData.AddressPooling.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"address-pooling "+v)
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if !rscData.PoolUtilizationAlarmClearThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"pool-utilization-alarm clear-threshold "+
			utils.ConvI64toa(rscData.PoolUtilizationAlarmClearThreshold.ValueInt64()))
	}
	if !rscData.PoolUtilizationAlarmRaiseThreshold.IsNull() {
		configSet = append(configSet, setPrefix+"pool-utilization-alarm raise-threshold "+
			utils.ConvI64toa(rscData.PoolUtilizationAlarmRaiseThreshold.ValueInt64()))
	}
	if rscData.PortNoTranslation.ValueBool() {
		configSet = append(configSet, setPrefix+"port no-translation")
	}
	if !rscData.PortOverloadingFactor.IsNull() {
		configSet = append(configSet, setPrefix+"port port-overloading-factor "+
			utils.ConvI64toa(rscData.PortOverloadingFactor.ValueInt64()))
	}
	if v := rscData.PortRange.ValueString(); v != "" {
		vSplit := strings.Split(v, "-")
		if len(vSplit) > 1 {
			configSet = append(configSet, setPrefix+"port range "+vSplit[0]+" to "+vSplit[1])
		} else {
			configSet = append(configSet, setPrefix+"port range "+vSplit[0])
		}
	}
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-instance "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityNatSourcePoolData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat source pool " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "address "):
				rscData.Address = append(rscData.Address, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "address-pooling "):
				rscData.AddressPooling = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "pool-utilization-alarm clear-threshold "):
				rscData.PoolUtilizationAlarmClearThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "pool-utilization-alarm raise-threshold "):
				rscData.PoolUtilizationAlarmRaiseThreshold, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "port no-translation":
				rscData.PortNoTranslation = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "port port-overloading-factor "):
				rscData.PortOverloadingFactor, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "port range to "):
				rscData.PortRange = types.StringValue(rscData.PortRange.ValueString() + "-" + itemTrim)
			case balt.CutPrefixInString(&itemTrim, "port range "):
				rscData.PortRange = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "routing-instance "):
				rscData.RoutingInstance = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (rscData *securityNatSourcePoolData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security nat source pool " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
