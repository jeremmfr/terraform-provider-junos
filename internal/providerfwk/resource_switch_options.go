package providerfwk

import (
	"context"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
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
	_ resource.Resource                = &switchOptions{}
	_ resource.ResourceWithConfigure   = &switchOptions{}
	_ resource.ResourceWithImportState = &switchOptions{}
)

type switchOptions struct {
	client *junos.Client
}

func newSwitchOptionsResource() resource.Resource {
	return &switchOptions{}
}

func (rsc *switchOptions) typeName() string {
	return providerName + "_switch_options"
}

func (rsc *switchOptions) junosName() string {
	return "switch-options"
}

func (rsc *switchOptions) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *switchOptions) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *switchOptions) Configure(
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

func (rsc *switchOptions) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `switch_options`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"clean_on_destroy": schema.BoolAttribute{
				Optional:    true,
				Description: "Clean supported lines when destroy this resource.",
			},
			"remote_vtep_list": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Configure static remote VXLAN tunnel endpoints.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringIPAddress().IPv4Only(),
					),
				},
			},
			"remote_vtep_v6_list": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Configure static ipv6 remote VXLAN tunnel endpoints.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						tfvalidator.StringIPAddress().IPv6Only(),
					),
				},
			},
			"service_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Service ID required if multi-chassis AE is part of a bridge-domain.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"vtep_source_interface": schema.StringAttribute{
				Optional:    true,
				Description: "Source layer-3 IFL for VXLAN.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					tfvalidator.String1DotCount(),
				},
			},
		},
	}
}

type switchOptionsData struct {
	ID                  types.String   `tfsdk:"id"`
	CleanOnDestroy      types.Bool     `tfsdk:"clean_on_destroy"`
	RemoteVtepList      []types.String `tfsdk:"remote_vtep_list"`
	RemoteVtepV6List    []types.String `tfsdk:"remote_vtep_v6_list"`
	ServiceID           types.Int64    `tfsdk:"service_id"`
	VTEPSourceInterface types.String   `tfsdk:"vtep_source_interface"`
}

func (rsc *switchOptions) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan switchOptionsData
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

func (rsc *switchOptions) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data switchOptionsData
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

func (rsc *switchOptions) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state switchOptionsData
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

func (rsc *switchOptions) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state switchOptionsData
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

func (rsc *switchOptions) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data switchOptionsData

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

func (rscData *switchOptionsData) fillID() {
	rscData.ID = types.StringValue("switch_options")
}

func (rscData *switchOptionsData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *switchOptionsData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set switch-options "

	for _, v := range rscData.RemoteVtepList {
		configSet = append(configSet, setPrefix+"remote-vtep-list "+v.ValueString())
	}
	for _, v := range rscData.RemoteVtepV6List {
		configSet = append(configSet, setPrefix+"remote-vtep-v6-list "+v.ValueString())
	}
	if !rscData.ServiceID.IsNull() {
		configSet = append(configSet, setPrefix+"service-id "+
			utils.ConvI64toa(rscData.ServiceID.ValueInt64()))
	}
	if v := rscData.VTEPSourceInterface.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"vtep-source-interface "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *switchOptionsData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"switch-options" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "remote-vtep-list "):
				rscData.RemoteVtepList = append(rscData.RemoteVtepList, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "remote-vtep-v6-list "):
				rscData.RemoteVtepV6List = append(rscData.RemoteVtepV6List, types.StringValue(itemTrim))
			case balt.CutPrefixInString(&itemTrim, "service-id "):
				rscData.ServiceID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "vtep-source-interface "):
				rscData.VTEPSourceInterface = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (rscData *switchOptionsData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := "delete switch-options "

	configSet := []string{
		delPrefix + "remote-vtep-list",
		delPrefix + "remote-vtep-v6-list",
		delPrefix + "service-id",
		delPrefix + "vtep-source-interface",
	}

	return junSess.ConfigSet(configSet)
}
