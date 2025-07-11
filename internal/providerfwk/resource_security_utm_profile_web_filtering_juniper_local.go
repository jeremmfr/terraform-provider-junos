package providerfwk

import (
	"context"
	"errors"
	"fmt"
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
	_ resource.Resource                   = &securityUtmProfileWebFilteringJuniperLocal{}
	_ resource.ResourceWithConfigure      = &securityUtmProfileWebFilteringJuniperLocal{}
	_ resource.ResourceWithValidateConfig = &securityUtmProfileWebFilteringJuniperLocal{}
	_ resource.ResourceWithImportState    = &securityUtmProfileWebFilteringJuniperLocal{}
	_ resource.ResourceWithUpgradeState   = &securityUtmProfileWebFilteringJuniperLocal{}
)

type securityUtmProfileWebFilteringJuniperLocal struct {
	client *junos.Client
}

func newSecurityUtmProfileWebFilteringJuniperLocalResource() resource.Resource {
	return &securityUtmProfileWebFilteringJuniperLocal{}
}

func (rsc *securityUtmProfileWebFilteringJuniperLocal) typeName() string {
	return providerName + "_security_utm_profile_web_filtering_juniper_local"
}

func (rsc *securityUtmProfileWebFilteringJuniperLocal) junosName() string {
	return "security utm feature-profile web-filtering juniper-local profile"
}

func (rsc *securityUtmProfileWebFilteringJuniperLocal) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityUtmProfileWebFilteringJuniperLocal) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityUtmProfileWebFilteringJuniperLocal) Configure(
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

func (rsc *securityUtmProfileWebFilteringJuniperLocal) Schema(
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
				Description: "The name of security utm feature-profile web-filtering juniper-local profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 29),
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
			"custom_message": schema.StringAttribute{
				Optional:    true,
				Description: "Custom message.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"default_action": schema.StringAttribute{
				Optional:    true,
				Description: "Default action.",
				Validators: []validator.String{
					stringvalidator.OneOf("block", "log-and-permit", "permit"),
				},
			},
			"no_safe_search": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not perform safe-search for Juniper local protocol.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
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
			"category":          securityUtmProfileWebFilteringBlockCategoryCustom{}.schema(),
			"fallback_settings": securityUtmProfileWebFilteringBlockFallbackSettings{}.schema(),
		},
	}
}

//nolint:lll
type securityUtmProfileWebFilteringJuniperLocalData struct {
	ID                 types.String                                         `tfsdk:"id"                   tfdata:"skip_isempty"`
	Name               types.String                                         `tfsdk:"name"                 tfdata:"skip_isempty"`
	CustomBlockMessage types.String                                         `tfsdk:"custom_block_message"`
	CustomMessage      types.String                                         `tfsdk:"custom_message"`
	DefaultAction      types.String                                         `tfsdk:"default_action"`
	NoSafeSearch       types.Bool                                           `tfsdk:"no_safe_search"`
	Timeout            types.Int64                                          `tfsdk:"timeout"`
	Category           []securityUtmProfileWebFilteringBlockCategoryCustom  `tfsdk:"category"`
	FallbackSettings   *securityUtmProfileWebFilteringBlockFallbackSettings `tfsdk:"fallback_settings"`
}

func (rscData *securityUtmProfileWebFilteringJuniperLocalData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

//nolint:lll
type securityUtmProfileWebFilteringJuniperLocalConfig struct {
	ID                 types.String                                         `tfsdk:"id"                   tfdata:"skip_isempty"`
	Name               types.String                                         `tfsdk:"name"                 tfdata:"skip_isempty"`
	CustomBlockMessage types.String                                         `tfsdk:"custom_block_message"`
	CustomMessage      types.String                                         `tfsdk:"custom_message"`
	DefaultAction      types.String                                         `tfsdk:"default_action"`
	NoSafeSearch       types.Bool                                           `tfsdk:"no_safe_search"`
	Timeout            types.Int64                                          `tfsdk:"timeout"`
	Category           types.List                                           `tfsdk:"category"`
	FallbackSettings   *securityUtmProfileWebFilteringBlockFallbackSettings `tfsdk:"fallback_settings"`
}

func (config *securityUtmProfileWebFilteringJuniperLocalConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(config)
}

func (rsc *securityUtmProfileWebFilteringJuniperLocal) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityUtmProfileWebFilteringJuniperLocalConfig
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

	if !config.Category.IsNull() &&
		!config.Category.IsUnknown() {
		var configCategory []securityUtmProfileWebFilteringBlockCategoryCustom
		asDiags := config.Category.ElementsAs(ctx, &configCategory, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		categoryName := make(map[string]struct{})
		for i, block := range configCategory {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := categoryName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("category").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple category blocks with the same name %q", name),
					)
				}
				categoryName[name] = struct{}{}
			}
		}
	}
}

func (rsc *securityUtmProfileWebFilteringJuniperLocal) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityUtmProfileWebFilteringJuniperLocalData
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
			profileExists, err := checkSecurityUtmProfileWebFilteringJuniperLocalExists(fnCtx, plan.Name.ValueString(), junSess)
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
			profileExists, err := checkSecurityUtmProfileWebFilteringJuniperLocalExists(fnCtx, plan.Name.ValueString(), junSess)
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

func (rsc *securityUtmProfileWebFilteringJuniperLocal) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityUtmProfileWebFilteringJuniperLocalData
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

func (rsc *securityUtmProfileWebFilteringJuniperLocal) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityUtmProfileWebFilteringJuniperLocalData
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

func (rsc *securityUtmProfileWebFilteringJuniperLocal) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityUtmProfileWebFilteringJuniperLocalData
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

func (rsc *securityUtmProfileWebFilteringJuniperLocal) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityUtmProfileWebFilteringJuniperLocalData

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

func checkSecurityUtmProfileWebFilteringJuniperLocalExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm feature-profile web-filtering juniper-local profile \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityUtmProfileWebFilteringJuniperLocalData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityUtmProfileWebFilteringJuniperLocalData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityUtmProfileWebFilteringJuniperLocalData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	if rscData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name`)")
	}

	configSet := make([]string, 0, 100)
	setPrefix := "set security utm feature-profile web-filtering juniper-local " +
		"profile \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.CustomBlockMessage.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"custom-block-message \""+v+"\"")
	}
	if v := rscData.CustomMessage.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"custom-message \""+v+"\"")
	}
	if v := rscData.DefaultAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"default "+v)
	}
	if rscData.NoSafeSearch.ValueBool() {
		configSet = append(configSet, setPrefix+"no-safe-search")
	}
	if !rscData.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(rscData.Timeout.ValueInt64()))
	}

	categoryName := make(map[string]struct{})
	for i, block := range rscData.Category {
		name := block.Name.ValueString()
		if _, ok := categoryName[name]; ok {
			return path.Root("category").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple category blocks with the same name %q", name)
		}
		categoryName[name] = struct{}{}

		configSet = append(configSet, block.configSet(setPrefix)...)
	}
	if rscData.FallbackSettings != nil {
		configSet = append(configSet, rscData.FallbackSettings.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityUtmProfileWebFilteringJuniperLocalData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm feature-profile web-filtering juniper-local profile \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "category "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.Category = tfdata.AppendPotentialNewBlock(rscData.Category, types.StringValue(strings.Trim(name, "\"")))
				category := &rscData.Category[len(rscData.Category)-1]
				balt.CutPrefixInString(&itemTrim, name+" ")

				category.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "custom-block-message "):
				rscData.CustomBlockMessage = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "custom-message "):
				rscData.CustomMessage = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "default "):
				rscData.DefaultAction = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "fallback-settings"):
				if rscData.FallbackSettings == nil {
					rscData.FallbackSettings = &securityUtmProfileWebFilteringBlockFallbackSettings{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					rscData.FallbackSettings.read(itemTrim)
				}
			case itemTrim == "no-safe-search":
				rscData.NoSafeSearch = types.BoolValue(true)
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

func (rscData *securityUtmProfileWebFilteringJuniperLocalData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security utm feature-profile web-filtering juniper-local profile \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
