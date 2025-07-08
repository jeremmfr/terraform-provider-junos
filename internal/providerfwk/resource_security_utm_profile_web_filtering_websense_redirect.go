package providerfwk

import (
	"context"
	"errors"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
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
	_ resource.Resource                   = &securityUtmProfileWebFilteringWebsenseRedirect{}
	_ resource.ResourceWithConfigure      = &securityUtmProfileWebFilteringWebsenseRedirect{}
	_ resource.ResourceWithValidateConfig = &securityUtmProfileWebFilteringWebsenseRedirect{}
	_ resource.ResourceWithImportState    = &securityUtmProfileWebFilteringWebsenseRedirect{}
	_ resource.ResourceWithUpgradeState   = &securityUtmProfileWebFilteringWebsenseRedirect{}
)

type securityUtmProfileWebFilteringWebsenseRedirect struct {
	client *junos.Client
}

func newSecurityUtmProfileWebFilteringWebsenseRedirectResource() resource.Resource {
	return &securityUtmProfileWebFilteringWebsenseRedirect{}
}

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) typeName() string {
	return providerName + "_security_utm_profile_web_filtering_websense_redirect"
}

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) junosName() string {
	return "security utm feature-profile web-filtering websense-redirect profile"
}

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) Configure(
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

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
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
				Description: "The name of security utm feature-profile web-filtering websense-redirect profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 29),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"account": schema.StringAttribute{
				Optional:    true,
				Description: "Set websense redirect account.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 28),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"custom_block_message": schema.StringAttribute{
				Optional:    true,
				Description: "Custom block message sent to HTTP client.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"no_safe_search": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not perform safe-search for Juniper local protocol.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"sockets": schema.Int64Attribute{
				Optional:    true,
				Description: "Set sockets number.",
				Validators: []validator.Int64{
					int64validator.Between(1, 32),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Set timeout.",
				Validators: []validator.Int64{
					int64validator.Between(1, 1800),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"fallback_settings": securityUtmProfileWebFilteringBlockFallbackSettings{}.schema(),
			"server": schema.SingleNestedBlock{
				Description: "Configure server settings.",
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						Optional:    true,
						Description: "Server host IP address or string host name.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"port": schema.Int64Attribute{
						Optional:    true,
						Description: "Server port.",
						Validators: []validator.Int64{
							int64validator.Between(1024, 65535),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
		},
	}
}

type securityUtmProfileWebFilteringWebsenseRedirectData struct {
	ID                 types.String                                               `tfsdk:"id"`
	Name               types.String                                               `tfsdk:"name"`
	Account            types.String                                               `tfsdk:"account"`
	CustomBlockMessage types.String                                               `tfsdk:"custom_block_message"`
	NoSafeSearch       types.Bool                                                 `tfsdk:"no_safe_search"`
	Sockets            types.Int64                                                `tfsdk:"sockets"`
	Timeout            types.Int64                                                `tfsdk:"timeout"`
	FallbackSettings   *securityUtmProfileWebFilteringBlockFallbackSettings       `tfsdk:"fallback_settings"`
	Server             *securityUtmProfileWebFilteringWebsenseRedirectBlockServer `tfsdk:"server"`
}

func (rscData *securityUtmProfileWebFilteringWebsenseRedirectData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

type securityUtmProfileWebFilteringWebsenseRedirectBlockServer struct {
	Host types.String `tfsdk:"host"`
	Port types.Int64  `tfsdk:"port"`
}

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityUtmProfileWebFilteringWebsenseRedirectData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one of arguments need to be set (in addition to `name`)",
		)
	}
}

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityUtmProfileWebFilteringWebsenseRedirectData
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
			profileExists, err := checkSecurityUtmProfileWebFilteringWebsenseRedirectExists(
				fnCtx, plan.Name.ValueString(), junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if profileExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			profileExists, err := checkSecurityUtmProfileWebFilteringWebsenseRedirectExists(
				fnCtx, plan.Name.ValueString(), junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !profileExists {
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

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityUtmProfileWebFilteringWebsenseRedirectData
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

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityUtmProfileWebFilteringWebsenseRedirectData
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

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityUtmProfileWebFilteringWebsenseRedirectData
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

func (rsc *securityUtmProfileWebFilteringWebsenseRedirect) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityUtmProfileWebFilteringWebsenseRedirectData

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

func checkSecurityUtmProfileWebFilteringWebsenseRedirectExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm feature-profile web-filtering websense-redirect profile \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityUtmProfileWebFilteringWebsenseRedirectData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityUtmProfileWebFilteringWebsenseRedirectData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityUtmProfileWebFilteringWebsenseRedirectData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	if rscData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name`)")
	}

	configSet := make([]string, 0, 100)
	setPrefix := "set security utm feature-profile web-filtering websense-redirect " +
		"profile \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.Account.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"account \""+v+"\"")
	}
	if v := rscData.CustomBlockMessage.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"custom-block-message \""+v+"\"")
	}
	if rscData.NoSafeSearch.ValueBool() {
		configSet = append(configSet, setPrefix+"no-safe-search")
	}
	if !rscData.Sockets.IsNull() {
		configSet = append(configSet, setPrefix+"sockets "+
			utils.ConvI64toa(rscData.Sockets.ValueInt64()))
	}
	if !rscData.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(rscData.Timeout.ValueInt64()))
	}

	if rscData.FallbackSettings != nil {
		configSet = append(configSet, rscData.FallbackSettings.configSet(setPrefix)...)
	}
	if rscData.Server != nil {
		configSet = append(configSet, rscData.Server.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityUtmProfileWebFilteringWebsenseRedirectBlockServer) configSet(setPrefix string) []string {
	setPrefix += "server "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if v := block.Host.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"host \""+v+"\"")
	}
	if !block.Port.IsNull() {
		configSet = append(configSet, setPrefix+"port "+
			utils.ConvI64toa(block.Port.ValueInt64()))
	}

	return configSet
}

func (rscData *securityUtmProfileWebFilteringWebsenseRedirectData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm feature-profile web-filtering websense-redirect profile \"" + name + "\"" +
		junos.PipeDisplaySetRelative,
	)
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
			case balt.CutPrefixInString(&itemTrim, "account "):
				rscData.Account = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "custom-block-message "):
				rscData.CustomBlockMessage = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "fallback-settings"):
				if rscData.FallbackSettings == nil {
					rscData.FallbackSettings = &securityUtmProfileWebFilteringBlockFallbackSettings{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					rscData.FallbackSettings.read(itemTrim)
				}
			case itemTrim == "no-safe-search":
				rscData.NoSafeSearch = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "server"):
				if rscData.Server == nil {
					rscData.Server = &securityUtmProfileWebFilteringWebsenseRedirectBlockServer{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.Server.read(itemTrim); err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "sockets "):
				rscData.Sockets, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "timeout "):
				rscData.Timeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *securityUtmProfileWebFilteringWebsenseRedirectBlockServer) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "host "):
		block.Host = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "port "):
		block.Port, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (rscData *securityUtmProfileWebFilteringWebsenseRedirectData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security utm feature-profile web-filtering websense-redirect profile \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
