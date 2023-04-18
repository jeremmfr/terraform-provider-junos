package providerfwk

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
	return "protocol oam gre-tunnel interface"
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
		Description: "Provides a " + rsc.junosName() + ".",
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
					stringvalidator.RegexMatches(regexp.MustCompile(`^gr-`),
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
			"could not create "+rsc.junosName()+" with empty name",
		)

		return
	}

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}
	defer func() { resp.Diagnostics.Append(tfdiag.Warns("Config Clear/Unlock Warning", junSess.ConfigClear())...) }()

	interfaceExists, err := checkOamGretunnelInterfaceExists(ctx, plan.Name.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.AddError("Pre Check Error", err.Error())

		return
	}
	if interfaceExists {
		resp.Diagnostics.AddError(
			"Duplicate Configuration Error",
			fmt.Sprintf(rsc.junosName()+" %q already exists", plan.Name.ValueString()),
		)

		return
	}

	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("create resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	interfaceExists, err = checkOamGretunnelInterfaceExists(ctx, plan.Name.ValueString(), junSess)
	if err != nil {
		resp.Diagnostics.AddError("Post Check Error", err.Error())

		return
	}
	if !interfaceExists {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf(rsc.junosName()+" %q does not exists after commit "+
				"=> check your config", plan.Name.ValueString()),
		)

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *oamGretunnelInterface) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data oamGretunnelInterfaceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	err = data.read(ctx, state.Name.ValueString(), junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}
	if data.ID.IsNull() {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

	if rsc.client.FakeUpdateAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
			} else {
				resp.Diagnostics.AddError("Config Set Error", err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}
	defer func() { resp.Diagnostics.Append(tfdiag.Warns("Config Clear/Unlock Warning", junSess.ConfigClear())...) }()

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, "Config Set Error", err.Error())
		} else {
			resp.Diagnostics.AddError("Config Set Error", err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf("update resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (rsc *oamGretunnelInterface) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state oamGretunnelInterfaceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if rsc.client.FakeDeleteAlso() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Config Del Error", err.Error())

			return
		}

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError("Config Lock Error", err.Error())

		return
	}
	defer func() { resp.Diagnostics.Append(tfdiag.Warns("Config Clear/Unlock Warning", junSess.ConfigClear())...) }()

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.AddError("Config Del Error", err.Error())

		return
	}
	warns, err := junSess.CommitConf("delete resource " + rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns("Config Commit Warning", warns)...)
	if err != nil {
		resp.Diagnostics.AddError("Config Commit Error", err.Error())

		return
	}
}

func (rsc *oamGretunnelInterface) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	var data oamGretunnelInterfaceData
	if err := data.read(ctx, req.ID, junSess); err != nil {
		resp.Diagnostics.AddError("Config Read Error", err.Error())

		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf("don't find "+rsc.junosName()+" with id %q "+
				"(id must be <name>)", req.ID),
		)

		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkOamGretunnelInterfaceExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	_ bool, err error,
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
) (
	err error,
) {
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
