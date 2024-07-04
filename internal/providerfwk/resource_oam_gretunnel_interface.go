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
	_ resource.Resource                   = &oamGretunnelInterface{}
	_ resource.ResourceWithConfigure      = &oamGretunnelInterface{}
	_ resource.ResourceWithValidateConfig = &oamGretunnelInterface{}
	_ resource.ResourceWithImportState    = &oamGretunnelInterface{}
)

type oamGretunnelInterface struct {
	client *junos.Client
}

func newOamGretunnelInterfaceResource() resource.Resource {
	return &oamGretunnelInterface{}
}

func (rsc *oamGretunnelInterface) typeName() string {
	return providerName + "_oam_gretunnel_interface"
}

func (rsc *oamGretunnelInterface) junosName() string {
	return "oam gre-tunnel interface"
}

func (rsc *oamGretunnelInterface) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *oamGretunnelInterface) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *oamGretunnelInterface) Configure(
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

func (rsc *oamGretunnelInterface) Schema(
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
				Description: "Name of interface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^gr-`),
						"must be a gr interface"),
				},
			},
			"hold_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Hold time (5..250 seconds).",
				Validators: []validator.Int64{
					int64validator.Between(5, 250),
				},
			},
			"keepalive_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Keepalive time (1..50 seconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 50),
				},
			},
		},
	}
}

type oamGretunnelInterfaceData struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	HoldTime      types.Int64  `tfsdk:"hold_time"`
	KeepaliveTime types.Int64  `tfsdk:"keepalive_time"`
}

func (rsc *oamGretunnelInterface) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config oamGretunnelInterfaceData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.HoldTime.IsNull() && !config.HoldTime.IsUnknown() &&
		!config.KeepaliveTime.IsNull() && !config.KeepaliveTime.IsUnknown() {
		if config.KeepaliveTime.ValueInt64()*2 > config.HoldTime.ValueInt64() {
			resp.Diagnostics.AddAttributeError(
				path.Root("hold_time"),
				"Bad Value Error",
				"hold_time has to be at least twice the keepalive_time",
			)
		}
	}
}

func (rsc *oamGretunnelInterface) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan oamGretunnelInterfaceData
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
			interfaceExists, err := checkOamGretunnelInterfaceExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if interfaceExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			interfaceExists, err := checkOamGretunnelInterfaceExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !interfaceExists {
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

func (rsc *oamGretunnelInterface) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data oamGretunnelInterfaceData
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

func (rsc *oamGretunnelInterface) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state oamGretunnelInterfaceData
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

func (rsc *oamGretunnelInterface) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state oamGretunnelInterfaceData
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

func (rsc *oamGretunnelInterface) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data oamGretunnelInterfaceData

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

func checkOamGretunnelInterfaceExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols oam gre-tunnel interface " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *oamGretunnelInterfaceData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *oamGretunnelInterfaceData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *oamGretunnelInterfaceData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set protocols oam gre-tunnel interface " + rscData.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if !rscData.HoldTime.IsNull() {
		configSet = append(configSet, setPrefix+"hold-time "+
			utils.ConvI64toa(rscData.HoldTime.ValueInt64()))
	}
	if !rscData.KeepaliveTime.IsNull() {
		configSet = append(configSet, setPrefix+"keepalive-time "+
			utils.ConvI64toa(rscData.KeepaliveTime.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *oamGretunnelInterfaceData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols oam gre-tunnel interface " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "hold-time "):
				rscData.HoldTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "keepalive-time "):
				rscData.KeepaliveTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *oamGretunnelInterfaceData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete protocols oam gre-tunnel interface " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
