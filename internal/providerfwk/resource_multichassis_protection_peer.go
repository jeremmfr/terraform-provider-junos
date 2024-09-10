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
	_ resource.Resource                = &multichassisProtectionPeer{}
	_ resource.ResourceWithConfigure   = &multichassisProtectionPeer{}
	_ resource.ResourceWithImportState = &multichassisProtectionPeer{}
)

type multichassisProtectionPeer struct {
	client *junos.Client
}

func newMultichassisProtectionPeerResource() resource.Resource {
	return &multichassisProtectionPeer{}
}

func (rsc *multichassisProtectionPeer) typeName() string {
	return providerName + "_multichassis_protection_peer"
}

func (rsc *multichassisProtectionPeer) junosName() string {
	return "multi-chassis multi-chassis-protection"
}

func (rsc *multichassisProtectionPeer) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *multichassisProtectionPeer) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *multichassisProtectionPeer) Configure(
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

func (rsc *multichassisProtectionPeer) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<ip_address>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip_address": schema.StringAttribute{
				Required:    true,
				Description: "IP address for this peer.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					tfvalidator.StringIPAddress().IPv4Only(),
				},
			},
			"interface": schema.StringAttribute{
				Required:    true,
				Description: "Inter-Chassis protection link.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
				},
			},
			"icl_down_delay": schema.Int64Attribute{
				Optional:    true,
				Description: "Time in seconds between ICL down and MCAEs moving to standby (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(1, 6000),
				},
			},
		},
	}
}

type multichassisProtectionPeerData struct {
	ID           types.String `tfsdk:"id"`
	IPAddress    types.String `tfsdk:"ip_address"`
	Interface    types.String `tfsdk:"interface"`
	IclDownDelay types.Int64  `tfsdk:"icl_down_delay"`
}

func (rsc *multichassisProtectionPeer) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan multichassisProtectionPeerData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.IPAddress.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("ip_address"),
			"Empty IP Address",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "ip_address"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			peerExists, err := checkMultichassisProtectionPeerExists(fnCtx, plan.IPAddress.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if peerExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.IPAddress),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			peerExists, err := checkMultichassisProtectionPeerExists(fnCtx, plan.IPAddress.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !peerExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.IPAddress),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *multichassisProtectionPeer) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data multichassisProtectionPeerData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.IPAddress.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *multichassisProtectionPeer) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state multichassisProtectionPeerData
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

func (rsc *multichassisProtectionPeer) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state multichassisProtectionPeerData
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

func (rsc *multichassisProtectionPeer) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data multichassisProtectionPeerData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "ip_address"),
	)
}

func checkMultichassisProtectionPeerExists(
	_ context.Context, ipAddress string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"multi-chassis multi-chassis-protection " + ipAddress + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *multichassisProtectionPeerData) fillID() {
	rscData.ID = types.StringValue(rscData.IPAddress.ValueString())
}

func (rscData *multichassisProtectionPeerData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *multichassisProtectionPeerData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set multi-chassis multi-chassis-protection " + rscData.IPAddress.ValueString() + " "

	configSet := []string{
		setPrefix,
		setPrefix + "interface " + rscData.Interface.ValueString(),
	}

	if !rscData.IclDownDelay.IsNull() {
		configSet = append(configSet, setPrefix+"icl-down-delay "+
			utils.ConvI64toa(rscData.IclDownDelay.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *multichassisProtectionPeerData) read(
	_ context.Context, ipAddress string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"multi-chassis multi-chassis-protection " + ipAddress + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.IPAddress = types.StringValue(ipAddress)
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
			case balt.CutPrefixInString(&itemTrim, "interface "):
				rscData.Interface = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "icl-down-delay "):
				rscData.IclDownDelay, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *multichassisProtectionPeerData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete multi-chassis multi-chassis-protection " + rscData.IPAddress.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
