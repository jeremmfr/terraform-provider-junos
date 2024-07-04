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
	_ resource.Resource                   = &securityNatDestinationPool{}
	_ resource.ResourceWithConfigure      = &securityNatDestinationPool{}
	_ resource.ResourceWithValidateConfig = &securityNatDestinationPool{}
	_ resource.ResourceWithImportState    = &securityNatDestinationPool{}
)

type securityNatDestinationPool struct {
	client *junos.Client
}

func newSecurityNatDestinationPoolResource() resource.Resource {
	return &securityNatDestinationPool{}
}

func (rsc *securityNatDestinationPool) typeName() string {
	return providerName + "_security_nat_destination_pool"
}

func (rsc *securityNatDestinationPool) junosName() string {
	return "security nat destination pool"
}

func (rsc *securityNatDestinationPool) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityNatDestinationPool) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityNatDestinationPool) Configure(
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

func (rsc *securityNatDestinationPool) Schema(
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
			"address": schema.StringAttribute{
				Required:    true,
				Description: "CIDR address to destination nat pool.",
				Validators: []validator.String{
					tfvalidator.StringCIDR(),
				},
			},
			"address_port": schema.Int64Attribute{
				Optional:    true,
				Description: "Port change too with destination nat.",
				Validators: []validator.Int64{
					int64validator.Between(0, 65535),
				},
			},
			"address_to": schema.StringAttribute{
				Optional:    true,
				Description: "CIDR to define range of address to destination nat pool.",
				Validators: []validator.String{
					tfvalidator.StringCIDR(),
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

type securityNatDestinationPoolData struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Address         types.String `tfsdk:"address"`
	AddressPort     types.Int64  `tfsdk:"address_port"`
	AddressTo       types.String `tfsdk:"address_to"`
	Description     types.String `tfsdk:"description"`
	RoutingInstance types.String `tfsdk:"routing_instance"`
}

func (rsc *securityNatDestinationPool) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityNatDestinationPoolData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.AddressPort.IsNull() && !config.AddressPort.IsUnknown() &&
		!config.AddressTo.IsNull() && !config.AddressTo.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("address_port"),
			tfdiag.ConflictConfigErrSummary,
			"address_port and address_to cannot be configured together",
		)
	}
}

func (rsc *securityNatDestinationPool) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityNatDestinationPoolData
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
			poolExists, err := checkSecurityNatDestinationPoolExists(fnCtx, plan.Name.ValueString(), junSess)
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
			poolExists, err := checkSecurityNatDestinationPoolExists(fnCtx, plan.Name.ValueString(), junSess)
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

func (rsc *securityNatDestinationPool) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityNatDestinationPoolData
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

func (rsc *securityNatDestinationPool) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityNatDestinationPoolData
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

func (rsc *securityNatDestinationPool) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityNatDestinationPoolData
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

func (rsc *securityNatDestinationPool) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityNatDestinationPoolData

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

func checkSecurityNatDestinationPoolExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat destination pool " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityNatDestinationPoolData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityNatDestinationPoolData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityNatDestinationPoolData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set security nat destination pool " + rscData.Name.ValueString() + " "

	configSet := []string{
		setPrefix + "address " + rscData.Address.ValueString(),
	}
	if !rscData.AddressPort.IsNull() {
		configSet = append(configSet, setPrefix+
			"address port "+utils.ConvI64toa(rscData.AddressPort.ValueInt64()))
	}
	if v := rscData.AddressTo.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"address to "+v)
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := rscData.RoutingInstance.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"routing-instance "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityNatDestinationPoolData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security nat destination pool " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "address port "):
				var err error
				rscData.AddressPort, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "address to "):
				rscData.AddressTo = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "address "):
				rscData.Address = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "routing-instance "):
				rscData.RoutingInstance = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (rscData *securityNatDestinationPoolData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security nat destination pool " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
