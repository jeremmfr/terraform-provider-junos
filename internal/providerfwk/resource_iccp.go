package providerfwk

import (
	"context"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
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
	_ resource.Resource                = &iccp{}
	_ resource.ResourceWithConfigure   = &iccp{}
	_ resource.ResourceWithImportState = &iccp{}
)

type iccp struct {
	client *junos.Client
}

func newIccpResource() resource.Resource {
	return &iccp{}
}

func (rsc *iccp) typeName() string {
	return providerName + "_iccp"
}

func (rsc *iccp) junosName() string {
	return "protocols iccp"
}

func (rsc *iccp) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *iccp) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *iccp) Configure(
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

func (rsc *iccp) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with value " +
					"`iccp`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"local_ip_addr": schema.StringAttribute{
				Required:    true,
				Description: "Local IP address to use by default for all peers.",
				Validators: []validator.String{
					tfvalidator.StringIPAddress().IPv4Only(),
				},
			},
			"authentication_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "MD5 authentication key for all peers.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 126),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"session_establishment_hold_time": schema.Int64Attribute{
				Optional:    true,
				Description: "Time within which connection must succeed with peers (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(45, 600),
				},
			},
		},
	}
}

type iccpData struct {
	ID                           types.String `tfsdk:"id"`
	LocalIPAddr                  types.String `tfsdk:"local_ip_addr"`
	AuthenticationKey            types.String `tfsdk:"authentication_key"`
	SessionEstablishmentHoldTime types.Int64  `tfsdk:"session_establishment_hold_time"`
}

func (rsc *iccp) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan iccpData
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

func (rsc *iccp) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data iccpData
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

func (rsc *iccp) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state iccpData
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

func (rsc *iccp) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state iccpData
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

func (rsc *iccp) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data iccpData

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

func (rscData *iccpData) fillID() {
	rscData.ID = types.StringValue("iccp")
}

func (rscData *iccpData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *iccpData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set protocols iccp "

	configSet := []string{
		setPrefix + "local-ip-addr " + rscData.LocalIPAddr.ValueString(),
	}

	if v := rscData.AuthenticationKey.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"authentication-key \""+v+"\"")
	}
	if !rscData.SessionEstablishmentHoldTime.IsNull() {
		configSet = append(configSet, setPrefix+"session-establishment-hold-time "+
			utils.ConvI64toa(rscData.SessionEstablishmentHoldTime.ValueInt64()))
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *iccpData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols iccp" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
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
			case balt.CutPrefixInString(&itemTrim, "local-ip-addr "):
				rscData.LocalIPAddr = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "authentication-key "):
				rscData.AuthenticationKey, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "authentication-key")
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "session-establishment-hold-time "):
				rscData.SessionEstablishmentHoldTime, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *iccpData) del(
	_ context.Context, junSess *junos.Session,
) error {
	delPrefix := "delete protocols iccp "

	configSet := []string{
		delPrefix + "local-ip-addr",
		delPrefix + "authentication-key",
		delPrefix + "session-establishment-hold-time",
	}

	return junSess.ConfigSet(configSet)
}
