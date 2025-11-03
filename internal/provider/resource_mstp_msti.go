package provider

import (
	"context"
	"fmt"
	"regexp"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &mstpMsti{}
	_ resource.ResourceWithConfigure      = &mstpMsti{}
	_ resource.ResourceWithValidateConfig = &mstpMsti{}
	_ resource.ResourceWithImportState    = &mstpMsti{}
)

type mstpMsti struct {
	client *junos.Client
}

func newMstpMstiResource() resource.Resource {
	return &mstpMsti{}
}

func (rsc *mstpMsti) typeName() string {
	return providerName + "_mstp_msti"
}

func (rsc *mstpMsti) junosName() string {
	return "mstp msti"
}

func (rsc *mstpMsti) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *mstpMsti) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *mstpMsti) Configure(
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

func (rsc *mstpMsti) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<msti_id>" + junos.IDSeparator + "<routing_instance>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"msti_id": schema.Int64Attribute{
				Required:    true,
				Description: "MSTI identifier.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 4094),
				},
			},
			"routing_instance": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(junos.DefaultW),
				Description: "Routing instance for mstp protocol if not root level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"vlan": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "VLAN ID or VLAN ID range.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.NoNullValues(),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(
							`^(409[0-4]|(40[0-8]|[1-3]\d\d|[1-9]\d|[1-9])\d|[1-9])`+
								`(-(409[0-4]|(40[0-8]|[1-3]\d\d|[1-9]\d|[1-9])\d|[1-9]))?$`),
							"must be a VLAN id (1..4094) or a range of VLAN id (1..4094)-(1..4094)"),
					),
				},
			},
			"backup_bridge_priority": schema.StringAttribute{
				Optional:    true,
				Description: "Priority of the bridge.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d\d?k$`),
						"must be a number with increments of 4k - 4k,8k,..60k",
					),
				},
			},
			"bridge_priority": schema.StringAttribute{
				Optional:    true,
				Description: "Priority of the bridge.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^(0|\d\d?k)$`),
						"must be a number with increments of 4k - 0,4k,8k,..60k",
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"interface": schema.SetNestedBlock{
				Description: "Interface options.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Interface name or `all`.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
								tfvalidator.StringDotExclusion(),
							},
						},
						"cost": schema.Int64Attribute{
							Optional:    true,
							Description: "Cost of the interface.",
							Validators: []validator.Int64{
								int64validator.Between(1, 200000000),
							},
						},
						"priority": schema.Int64Attribute{
							Optional:    true,
							Description: "Interface priority (in increments of 16).",
							Validators: []validator.Int64{
								int64validator.Between(0, 240),
							},
						},
					},
				},
			},
		},
	}
}

type mstpMstiData struct {
	ID                   types.String             `tfsdk:"id"`
	MstiID               types.Int64              `tfsdk:"msti_id"`
	RoutingInstance      types.String             `tfsdk:"routing_instance"`
	Vlan                 []types.String           `tfsdk:"vlan"`
	BackupBridgePriority types.String             `tfsdk:"backup_bridge_priority"`
	BridgePriority       types.String             `tfsdk:"bridge_priority"`
	Interface            []mstpMstiBlockInterface `tfsdk:"interface"`
}

type mstpMstiConfig struct {
	ID                   types.String `tfsdk:"id"`
	MstiID               types.Int64  `tfsdk:"msti_id"`
	RoutingInstance      types.String `tfsdk:"routing_instance"`
	Vlan                 types.Set    `tfsdk:"vlan"`
	BackupBridgePriority types.String `tfsdk:"backup_bridge_priority"`
	BridgePriority       types.String `tfsdk:"bridge_priority"`
	Interface            types.Set    `tfsdk:"interface"`
}

type mstpMstiBlockInterface struct {
	Name     types.String `tfsdk:"name"     tfdata:"identifier,skip_isempty"`
	Cost     types.Int64  `tfsdk:"cost"`
	Priority types.Int64  `tfsdk:"priority"`
}

func (block *mstpMstiBlockInterface) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (rsc *mstpMsti) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config mstpMstiConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Interface.IsNull() && !config.Interface.IsUnknown() {
		var interFace []mstpMstiBlockInterface
		asDiags := config.Interface.ElementsAs(ctx, &interFace, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		interfaceName := make(map[string]struct{})
		for _, block := range interFace {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := interfaceName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple interface blocks with the same name %q", name),
					)
				}
				interfaceName[name] = struct{}{}
			}
			if block.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("interface"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("cost or priority must be specified"+
						" in interface block %q", block.Name.ValueString()),
				)
			}
			if !block.Priority.IsNull() && !block.Priority.IsUnknown() {
				if block.Priority.ValueInt64()%16 != 0 {
					resp.Diagnostics.AddAttributeError(
						path.Root("interface"),
						"Bad Value Error",
						fmt.Sprintf("priority must be a multiple of 16"+
							" in interface block %q", block.Name.ValueString()),
					)
				}
			}
		}
	}
}

func (rsc *mstpMsti) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan mstpMstiData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.MstiID.ValueInt64() == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("msti_id"),
			"Zero msti-id",
			"could not create "+rsc.junosName()+" with msti_id at zero",
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
				instanceExists, err := checkRoutingInstanceExists(fnCtx, v, junSess)
				if err != nil {
					resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

					return false
				}
				if !instanceExists {
					resp.Diagnostics.AddAttributeError(
						path.Root("routing_instance"),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("routing instance %q doesn't exist", v),
					)

					return false
				}
			}
			mstiInstanceExists, err := checkMstpMstiExists(
				fnCtx,
				plan.MstiID.ValueInt64(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if mstiInstanceExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf(rsc.junosName()+" %d already exists in routing-instance %q",
							plan.MstiID.ValueInt64(), v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf(rsc.junosName()+" %d already exists", plan.MstiID.ValueInt64()),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			mstiInstanceExists, err := checkMstpMstiExists(
				fnCtx,
				plan.MstiID.ValueInt64(),
				plan.RoutingInstance.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !mstiInstanceExists {
				if v := plan.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						fmt.Sprintf(rsc.junosName()+" %d does not exists in routing-instance %q after commit "+
							"=> check your config", plan.MstiID.ValueInt64(), v),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						fmt.Sprintf(rsc.junosName()+" %d does not exists after commit "+
							"=> check your config", plan.MstiID.ValueInt64()),
					)
				}

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *mstpMsti) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data mstpMstiData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1Int1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.MstiID.ValueInt64(),
			state.RoutingInstance.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *mstpMsti) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state mstpMstiData
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

func (rsc *mstpMsti) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state mstpMstiData
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

func (rsc *mstpMsti) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data mstpMstiData

	var _ resourceDataReadFrom1Int1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <msti_id>"+junos.IDSeparator+"<routing_instance>)",
	)
}

func checkMstpMstiExists(
	_ context.Context, mstiID int64, routingInstance string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols mstp msti " + utils.ConvI64toa(mstiID) + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *mstpMstiData) fillID() {
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		rscData.ID = types.StringValue(utils.ConvI64toa(rscData.MstiID.ValueInt64()) + junos.IDSeparator + v)
	} else {
		rscData.ID = types.StringValue(utils.ConvI64toa(rscData.MstiID.ValueInt64()) + junos.IDSeparator + junos.DefaultW)
	}
}

func (rscData *mstpMstiData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *mstpMstiData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)
	setPrefix := junos.SetLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		setPrefix += junos.RoutingInstancesWS + v + " "
	}
	setPrefix += "protocols mstp msti " + utils.ConvI64toa(rscData.MstiID.ValueInt64()) + " "

	for _, v := range rscData.Vlan {
		configSet = append(configSet, setPrefix+"vlan "+v.ValueString())
	}

	if v := rscData.BackupBridgePriority.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"backup-bridge-priority "+v)
	}
	if v := rscData.BridgePriority.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"bridge-priority "+v)
	}

	interfaceName := make(map[string]struct{})
	for _, block := range rscData.Interface {
		name := block.Name.ValueString()
		if _, ok := interfaceName[name]; ok {
			return path.Root("system_id"),
				fmt.Errorf("multiple interface blocks with the same name %q", name)
		}
		interfaceName[name] = struct{}{}

		if block.isEmpty() {
			return path.Root("interface"),
				fmt.Errorf("cost or priority must be specified"+
					" in interface block %q", name)
		}

		if !block.Cost.IsNull() {
			configSet = append(configSet, setPrefix+"interface "+name+" cost "+
				utils.ConvI64toa(block.Cost.ValueInt64()))
		}
		if !block.Priority.IsNull() {
			configSet = append(configSet, setPrefix+"interface "+name+" priority "+
				utils.ConvI64toa(block.Priority.ValueInt64()))
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *mstpMstiData) read(
	_ context.Context, mstiID int64, routingInstance string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig
	if routingInstance != "" && routingInstance != junos.DefaultW {
		showPrefix += junos.RoutingInstancesWS + routingInstance + " "
	}
	showConfig, err := junSess.Command(showPrefix +
		"protocols mstp msti " + utils.ConvI64toa(mstiID) + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.MstiID = types.Int64Value(mstiID)
		if routingInstance == "" {
			rscData.RoutingInstance = types.StringValue(junos.DefaultW)
		} else {
			rscData.RoutingInstance = types.StringValue(routingInstance)
		}
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
			case balt.CutPrefixInString(&itemTrim, "vlan "):
				rscData.Vlan = append(rscData.Vlan, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "backup-bridge-priority "):
				rscData.BackupBridgePriority = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "bridge-priority "):
				rscData.BridgePriority = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "interface "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				var interFace mstpMstiBlockInterface
				rscData.Interface, interFace = tfdata.ExtractBlock(rscData.Interface, types.StringValue(name))
				balt.CutPrefixInString(&itemTrim, name+" ")

				switch {
				case balt.CutPrefixInString(&itemTrim, "cost "):
					interFace.Cost, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, "priority "):
					interFace.Priority, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				}
				rscData.Interface = append(rscData.Interface, interFace)
			}
		}
	}

	return nil
}

func (rscData *mstpMstiData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS
	if v := rscData.RoutingInstance.ValueString(); v != "" && v != junos.DefaultW {
		delPrefix += junos.RoutingInstancesWS + v + " "
	}

	configSet := []string{
		delPrefix + "protocols mstp msti " + utils.ConvI64toa(rscData.MstiID.ValueInt64()),
	}

	return junSess.ConfigSet(configSet)
}
