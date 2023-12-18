package providerfwk

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &interfacePhysicalDisable{}
	_ resource.ResourceWithConfigure = &interfacePhysicalDisable{}
)

type interfacePhysicalDisable struct {
	client *junos.Client
}

func newInterfacePhysicalDisableResource() resource.Resource {
	return &interfacePhysicalDisable{}
}

func (rsc *interfacePhysicalDisable) typeName() string {
	return providerName + "_interface_physical_disable"
}

func (rsc *interfacePhysicalDisable) junosName() string {
	return "not configured physical interface"
}

func (rsc *interfacePhysicalDisable) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *interfacePhysicalDisable) Configure(
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

func (rsc *interfacePhysicalDisable) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Disable " + rsc.junosName(),
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
				Description: "Name of physical interface (without dot).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					tfvalidator.StringDotExclusion(),
				},
			},
		},
	}
}

type interfacePhysicalDisableData struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (rsc *interfacePhysicalDisable) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan interfacePhysicalDisableData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := addInterfaceNC(
			ctx,
			plan.Name.ValueString(),
			rsc.client.GroupInterfaceDelete(),
			junSess,
		); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigClearUnlockWarnSummary, junSess.ConfigClear())...)
	}()

	ncInt, emptyInt, err := checkInterfacePhysicalNCEmpty(
		ctx,
		plan.Name.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

		return
	}
	if !ncInt && !emptyInt {
		resp.Diagnostics.AddError(
			"Conflict Error",
			fmt.Sprintf("interface %q is configured", plan.Name.ValueString()),
		)

		return
	}
	if ncInt {
		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}
	if emptyInt {
		if containsUnit, err := checkInterfacePhysicalContainsUnit(ctx, plan.Name.ValueString(), junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

			return
		} else if containsUnit {
			resp.Diagnostics.AddError(
				"Conflict Error",
				fmt.Sprintf("interface %q is used for a logical unit interface", plan.Name.ValueString()),
			)

			return
		}
	}

	if err := addInterfaceNC(
		ctx,
		plan.Name.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())

		return
	}
	warns, err := junSess.CommitConf(ctx, "create resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	ncInt, _, err = checkInterfacePhysicalNCEmpty(
		ctx,
		plan.Name.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

		return
	}
	if !ncInt {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			fmt.Sprintf("interface %q not disable after commit "+
				"=> check your config", plan.Name.ValueString()),
		)

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *interfacePhysicalDisable) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state interfacePhysicalDisableData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	ncInt, _, err := checkInterfacePhysicalNCEmpty(
		ctx,
		state.Name.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if !ncInt {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (rsc *interfacePhysicalDisable) Update(
	_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse,
) {
}

func (rsc *interfacePhysicalDisable) Delete(
	_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse,
) {
}

func (rscData *interfacePhysicalDisableData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}
