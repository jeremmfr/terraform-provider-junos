package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &interfaceSt0Unit{}
	_ resource.ResourceWithConfigure   = &interfaceSt0Unit{}
	_ resource.ResourceWithImportState = &interfaceSt0Unit{}
)

type interfaceSt0Unit struct {
	client *junos.Client
}

func newInterfaceSt0UnitResource() resource.Resource {
	return &interfaceSt0Unit{}
}

func (rsc *interfaceSt0Unit) typeName() string {
	return providerName + "_interface_st0_unit"
}

func (rsc *interfaceSt0Unit) junosName() string {
	return "st0 logical interface"
}

func (rsc *interfaceSt0Unit) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *interfaceSt0Unit) Configure(
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

func (rsc *interfaceSt0Unit) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Find an available " + rsc.junosName() + " and create it.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with the name of interface found.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type interfaceSt0UnitData struct {
	ID types.String `tfsdk:"id"`
}

func (rsc *interfaceSt0Unit) Create(
	ctx context.Context, _ resource.CreateRequest, resp *resource.CreateResponse,
) {
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

	newSt0, err := rsc.searchNewAvailable(junSess)
	if err != nil {
		resp.Diagnostics.AddError("Search Error", err.Error())

		return
	}
	if err := junSess.ConfigSet([]string{
		"set interfaces " + newSt0,
	}); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())

		return
	}
	warns, err := junSess.CommitConf(ctx, "create resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(
		ctx,
		newSt0,
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

		return
	}
	if ncInt {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			fmt.Sprintf(rsc.junosName()+" %q always disable (NC) after commit "+
				"=> check your config", newSt0),
		)

		return
	}
	if emptyInt && !setInt {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			fmt.Sprintf("create new "+rsc.junosName()+" %q doesn't works, "+
				"can't find the new interface after commit", newSt0),
		)

		return
	}

	data := interfaceSt0UnitData{
		ID: types.StringValue(newSt0),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rsc *interfaceSt0Unit) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state interfaceSt0UnitData
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
	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(
		ctx,
		state.ID.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if ncInt || (emptyInt && !setInt) {
		resp.State.RemoveResource(ctx)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (rsc *interfaceSt0Unit) Update(
	_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse,
) {
}

func (rsc *interfaceSt0Unit) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state interfaceSt0UnitData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeDeleteAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := junSess.ConfigSet([]string{
			"delete interfaces " + state.ID.ValueString(),
		}); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

			return
		}

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

	ncInt, emptyInt, _, err := checkInterfaceLogicalNCEmpty(
		ctx,
		state.ID.ValueString(),
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

		return
	}
	if !ncInt && !emptyInt {
		resp.Diagnostics.AddError(
			tfdiag.PreCheckErrSummary,
			fmt.Sprintf("interface %q not empty or disable", state.ID.ValueString()),
		)

		return
	}
	if err := junSess.ConfigSet([]string{
		"delete interfaces " + state.ID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

		return
	}
	warns, err := junSess.CommitConf(ctx, "delete resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}
}

func (rsc *interfaceSt0Unit) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	if !strings.HasPrefix(req.ID, "st0.") {
		resp.Diagnostics.AddError(
			tfdiag.PreCheckErrSummary,
			fmt.Sprintf("name of interface need to state with 'st0.', got %q", req.ID),
		)

		return
	}
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	ncInt, emptyInt, setInt, err := checkInterfaceLogicalNCEmpty(
		ctx,
		req.ID,
		rsc.client.GroupInterfaceDelete(),
		junSess,
	)
	if err != nil {
		resp.Diagnostics.AddError("Interface Read Error", err.Error())

		return
	}
	if ncInt {
		resp.Diagnostics.AddError(
			"Disable Error",
			fmt.Sprintf("interface %q is disabled (NC), import is not possible", req.ID),
		)

		return
	}
	if emptyInt && !setInt {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be the name of st0 unit interface <st0.?>)",
		)

		return
	}

	data := interfaceSt0UnitData{
		ID: types.StringValue(req.ID),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (rsc *interfaceSt0Unit) searchNewAvailable(junSess *junos.Session) (string, error) {
	st0, err := junSess.Command("show interfaces st0 terse")
	if err != nil {
		return "", err
	}
	st0Line := strings.Split(st0, "\n")
	st0int := make([]string, 0)
	for _, line := range st0Line {
		if strings.HasPrefix(line, "st0.") {
			lineFields := strings.Split(line, " ")
			st0int = append(st0int, lineFields[0])
		}
	}
	for i := 0; i <= 1073741823; i++ {
		if !slices.Contains(st0int, "st0."+strconv.Itoa(i)) {
			return "st0." + strconv.Itoa(i), nil
		}
	}

	return "", errors.New("error for find st0 unit to create")
}
