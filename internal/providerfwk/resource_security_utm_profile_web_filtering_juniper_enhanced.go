package providerfwk

import (
	"context"
	"errors"
	"fmt"
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
	_ resource.Resource                   = &securityUtmProfileWebFilteringJuniperEnhanced{}
	_ resource.ResourceWithConfigure      = &securityUtmProfileWebFilteringJuniperEnhanced{}
	_ resource.ResourceWithValidateConfig = &securityUtmProfileWebFilteringJuniperEnhanced{}
	_ resource.ResourceWithImportState    = &securityUtmProfileWebFilteringJuniperEnhanced{}
	_ resource.ResourceWithUpgradeState   = &securityUtmProfileWebFilteringJuniperEnhanced{}
)

type securityUtmProfileWebFilteringJuniperEnhanced struct {
	client *junos.Client
}

func newSecurityUtmProfileWebFilteringJuniperEnhancedResource() resource.Resource {
	return &securityUtmProfileWebFilteringJuniperEnhanced{}
}

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) typeName() string {
	return providerName + "_security_utm_profile_web_filtering_juniper_enhanced"
}

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) junosName() string {
	return "security utm feature-profile web-filtering juniper-enhanced profile"
}

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) Configure(
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

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) Schema(
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
				Description: "The name of security utm feature-profile web-filtering juniper-enhanced profile.",
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
			"default_action": schema.StringAttribute{
				Optional:    true,
				Description: "Default action.",
				Validators: []validator.String{
					stringvalidator.OneOf("block", "log-and-permit", "permit", "quarantine"),
				},
			},
			"no_safe_search": schema.BoolAttribute{
				Optional:    true,
				Description: "Do not perform safe-search for Juniper enhanced protocol.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"quarantine_custom_message": schema.StringAttribute{
				Optional:    true,
				Description: "Quarantine custom message.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 512),
					tfvalidator.StringDoubleQuoteExclusion(),
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
			"block_message": schema.SingleNestedBlock{
				Description: "Configure block message.",
				Attributes: map[string]schema.Attribute{
					"type_custom_redirect_url": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable Custom redirect URL server type.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"url": schema.StringAttribute{
						Optional:    true,
						Description: "URL of block message.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 1024),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"category": schema.ListNestedBlock{
				Description: "For each name of category, configure enhanced category actions.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of category.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"action": schema.StringAttribute{
							Required:    true,
							Description: "Action when web traffic matches category.",
							Validators: []validator.String{
								stringvalidator.OneOf("block", "log-and-permit", "permit", "quarantine"),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"reputation_action": schema.ListNestedBlock{
							Description: "For each site_reputation, configure action.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"site_reputation": schema.StringAttribute{
										Required:    true,
										Description: "Level of reputation.",
										Validators: []validator.String{
											stringvalidator.OneOf("fairly-safe", "harmful", "moderately-safe", "suspicious", "very-safe"),
										},
									},
									"action": schema.StringAttribute{
										Required:    true,
										Description: "Action for site-reputation.",
										Validators: []validator.String{
											stringvalidator.OneOf("block", "log-and-permit", "permit", "quarantine"),
										},
									},
								},
							},
						},
					},
				},
			},
			"fallback_settings": securityUtmProfileWebFilteringBlockFallbackSettings{}.schema(),
			"quarantine_message": schema.SingleNestedBlock{
				Description: "Configure quarantine message.",
				Attributes: map[string]schema.Attribute{
					"type_custom_redirect_url": schema.BoolAttribute{
						Optional:    true,
						Description: "Enable Custom redirect URL server type.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"url": schema.StringAttribute{
						Optional:    true,
						Description: "URL of quarantine message.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 1024),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"site_reputation_action": schema.ListNestedBlock{
				Description: "For each site_reputation, configure action.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"site_reputation": schema.StringAttribute{
							Required:    true,
							Description: "Level of reputation.",
							Validators: []validator.String{
								stringvalidator.OneOf("fairly-safe", "harmful", "moderately-safe", "suspicious", "very-safe"),
							},
						},
						"action": schema.StringAttribute{
							Required:    true,
							Description: "Action for site-reputation.",
							Validators: []validator.String{
								stringvalidator.OneOf("block", "log-and-permit", "permit", "quarantine"),
							},
						},
					},
				},
			},
		},
	}
}

//nolint:lll
type securityUtmProfileWebFilteringJuniperEnhancedData struct {
	ID                      types.String                                                         `tfsdk:"id"`
	Name                    types.String                                                         `tfsdk:"name"`
	CustomBlockMessage      types.String                                                         `tfsdk:"custom_block_message"`
	DefaultAction           types.String                                                         `tfsdk:"default_action"`
	NoSafeSearch            types.Bool                                                           `tfsdk:"no_safe_search"`
	QuarantineCustomMessage types.String                                                         `tfsdk:"quarantine_custom_message"`
	Timeout                 types.Int64                                                          `tfsdk:"timeout"`
	BlockMessage            *securityUtmProfileWebFilteringJuniperEnhancedBlockMessage           `tfsdk:"block_message"`
	Category                []securityUtmProfileWebFilteringJuniperEnhancedBlockCategory         `tfsdk:"category"`
	FallbackSettings        *securityUtmProfileWebFilteringBlockFallbackSettings                 `tfsdk:"fallback_settings"`
	QuarantineMessage       *securityUtmProfileWebFilteringJuniperEnhancedBlockMessage           `tfsdk:"quarantine_message"`
	SiteReputationAction    []securityUtmProfileWebFilteringJuniperEnhancedBlockReputationAction `tfsdk:"site_reputation_action"`
}

func (rscData *securityUtmProfileWebFilteringJuniperEnhancedData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

type securityUtmProfileWebFilteringJuniperEnhancedConfig struct {
	ID                      types.String                                               `tfsdk:"id"`
	Name                    types.String                                               `tfsdk:"name"`
	CustomBlockMessage      types.String                                               `tfsdk:"custom_block_message"`
	DefaultAction           types.String                                               `tfsdk:"default_action"`
	NoSafeSearch            types.Bool                                                 `tfsdk:"no_safe_search"`
	QuarantineCustomMessage types.String                                               `tfsdk:"quarantine_custom_message"`
	Timeout                 types.Int64                                                `tfsdk:"timeout"`
	BlockMessage            *securityUtmProfileWebFilteringJuniperEnhancedBlockMessage `tfsdk:"block_message"`
	Category                types.List                                                 `tfsdk:"category"`
	FallbackSettings        *securityUtmProfileWebFilteringBlockFallbackSettings       `tfsdk:"fallback_settings"`
	QuarantineMessage       *securityUtmProfileWebFilteringJuniperEnhancedBlockMessage `tfsdk:"quarantine_message"`
	SiteReputationAction    types.List                                                 `tfsdk:"site_reputation_action"`
}

func (config *securityUtmProfileWebFilteringJuniperEnhancedConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(config)
}

type securityUtmProfileWebFilteringJuniperEnhancedBlockMessage struct {
	TypeCustomRedirectURL types.Bool   `tfsdk:"type_custom_redirect_url"`
	URL                   types.String `tfsdk:"url"`
}

//nolint:lll
type securityUtmProfileWebFilteringJuniperEnhancedBlockCategory struct {
	Name             types.String                                                         `tfsdk:"name"              tfdata:"identifier"`
	Action           types.String                                                         `tfsdk:"action"`
	ReputationAction []securityUtmProfileWebFilteringJuniperEnhancedBlockReputationAction `tfsdk:"reputation_action"`
}

type securityUtmProfileWebFilteringJuniperEnhancedBlockCategoryConfig struct {
	Name             types.String `tfsdk:"name"`
	Action           types.String `tfsdk:"action"`
	ReputationAction types.List   `tfsdk:"reputation_action"`
}

type securityUtmProfileWebFilteringJuniperEnhancedBlockReputationAction struct {
	SiteReputation types.String `tfsdk:"site_reputation" tfdata:"identifier"`
	Action         types.String `tfsdk:"action"`
}

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityUtmProfileWebFilteringJuniperEnhancedConfig
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
		var configCategory []securityUtmProfileWebFilteringJuniperEnhancedBlockCategoryConfig
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

			if !block.ReputationAction.IsNull() &&
				!block.ReputationAction.IsUnknown() {
				var configReputationAction []securityUtmProfileWebFilteringJuniperEnhancedBlockReputationAction
				asDiags := block.ReputationAction.ElementsAs(ctx, &configReputationAction, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				reputationActionSiteReputation := make(map[string]struct{})
				for ii, subBlock := range configReputationAction {
					if subBlock.SiteReputation.IsUnknown() {
						continue
					}
					siteReputation := subBlock.SiteReputation.ValueString()
					if _, ok := reputationActionSiteReputation[siteReputation]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("category").AtListIndex(i).AtName("reputation_action").AtListIndex(ii).AtName("site_reputation"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple reputation_action blocks with the same site_reputation %q"+
								" in category blocks %q", siteReputation, block.Name.ValueString()),
						)
					}
					reputationActionSiteReputation[siteReputation] = struct{}{}
				}
			}
		}
	}
	if !config.SiteReputationAction.IsNull() &&
		!config.SiteReputationAction.IsUnknown() {
		var configSiteReputationAction []securityUtmProfileWebFilteringJuniperEnhancedBlockReputationAction
		asDiags := config.SiteReputationAction.ElementsAs(ctx, &configSiteReputationAction, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		siteReputationActionSiteReputation := make(map[string]struct{})
		for i, block := range configSiteReputationAction {
			if block.SiteReputation.IsUnknown() {
				continue
			}
			siteReputation := block.SiteReputation.ValueString()
			if _, ok := siteReputationActionSiteReputation[siteReputation]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("site_reputation_action").AtListIndex(i).AtName("site_reputation"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple site_reputation_action blocks with the same site_reputation %q", siteReputation),
				)
			}
			siteReputationActionSiteReputation[siteReputation] = struct{}{}
		}
	}
}

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityUtmProfileWebFilteringJuniperEnhancedData
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
			profileExists, err := checkSecurityUtmProfileWebFilteringJuniperEnhancedExists(
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
			profileExists, err := checkSecurityUtmProfileWebFilteringJuniperEnhancedExists(
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

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityUtmProfileWebFilteringJuniperEnhancedData
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

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityUtmProfileWebFilteringJuniperEnhancedData
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

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityUtmProfileWebFilteringJuniperEnhancedData
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

func (rsc *securityUtmProfileWebFilteringJuniperEnhanced) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityUtmProfileWebFilteringJuniperEnhancedData

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

func checkSecurityUtmProfileWebFilteringJuniperEnhancedExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm feature-profile web-filtering juniper-enhanced profile \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityUtmProfileWebFilteringJuniperEnhancedData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityUtmProfileWebFilteringJuniperEnhancedData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityUtmProfileWebFilteringJuniperEnhancedData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	if rscData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name`)")
	}

	configSet := make([]string, 0, 100)
	setPrefix := "set security utm feature-profile web-filtering juniper-enhanced" +
		" profile \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.CustomBlockMessage.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"custom-block-message \""+v+"\"")
	}
	if v := rscData.DefaultAction.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"default "+v)
	}
	if rscData.NoSafeSearch.ValueBool() {
		configSet = append(configSet, setPrefix+"no-safe-search")
	}
	if v := rscData.QuarantineCustomMessage.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"quarantine-custom-message \""+v+"\"")
	}
	if !rscData.Timeout.IsNull() {
		configSet = append(configSet, setPrefix+"timeout "+
			utils.ConvI64toa(rscData.Timeout.ValueInt64()))
	}

	if rscData.BlockMessage != nil {
		configSet = append(configSet, rscData.BlockMessage.configSet(setPrefix+"block-message ")...)
	}
	categoryName := make(map[string]struct{})
	for i, block := range rscData.Category {
		name := block.Name.ValueString()
		if _, ok := categoryName[name]; ok {
			return path.Root("category").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple category blocks with the same name %q", name)
		}
		categoryName[name] = struct{}{}

		configSet = append(configSet, setPrefix+"category \""+name+"\" action "+block.Action.ValueString())

		reputationActionSiteReputation := make(map[string]struct{})
		for ii, subBlock := range block.ReputationAction {
			siteReputation := subBlock.SiteReputation.ValueString()
			if _, ok := reputationActionSiteReputation[siteReputation]; ok {
				return path.Root("category").AtListIndex(i).AtName("reputation_action").AtListIndex(ii).AtName("site_reputation"),
					fmt.Errorf("multiple reputation_action blocks with the same site_reputation %q"+
						" in category block %q", siteReputation, name)
			}
			reputationActionSiteReputation[siteReputation] = struct{}{}

			configSet = append(configSet, setPrefix+"category \""+name+"\""+
				" reputation-action "+siteReputation+" "+subBlock.Action.ValueString())
		}
	}
	if rscData.FallbackSettings != nil {
		configSet = append(configSet, rscData.FallbackSettings.configSet(setPrefix)...)
	}
	if rscData.QuarantineMessage != nil {
		configSet = append(configSet, rscData.QuarantineMessage.configSet(setPrefix+"quarantine-message ")...)
	}
	siteReputationActionSiteReputation := make(map[string]struct{})
	for i, block := range rscData.SiteReputationAction {
		siteReputation := block.SiteReputation.ValueString()
		if _, ok := siteReputationActionSiteReputation[siteReputation]; ok {
			return path.Root("site_reputation_action").AtListIndex(i).AtName("site_reputation"),
				fmt.Errorf("multiple site_reputation_action blocks with the same site_reputation %q", siteReputation)
		}
		siteReputationActionSiteReputation[siteReputation] = struct{}{}

		configSet = append(configSet, setPrefix+"site-reputation-action "+siteReputation+" "+block.Action.ValueString())
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityUtmProfileWebFilteringJuniperEnhancedBlockMessage) configSet(setPrefix string) []string {
	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if block.TypeCustomRedirectURL.ValueBool() {
		configSet = append(configSet, setPrefix+"type custom-redirect-url")
	}
	if v := block.URL.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"url \""+v+"\"")
	}

	return configSet
}

func (rscData *securityUtmProfileWebFilteringJuniperEnhancedData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm feature-profile web-filtering juniper-enhanced profile \"" + name + "\"" +
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
			case balt.CutPrefixInString(&itemTrim, "block-message"):
				if rscData.BlockMessage == nil {
					rscData.BlockMessage = &securityUtmProfileWebFilteringJuniperEnhancedBlockMessage{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					rscData.BlockMessage.read(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "category "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.Category = tfdata.AppendPotentialNewBlock(rscData.Category, types.StringValue(strings.Trim(name, "\"")))
				category := &rscData.Category[len(rscData.Category)-1]
				balt.CutPrefixInString(&itemTrim, name+" ")

				switch {
				case balt.CutPrefixInString(&itemTrim, "action "):
					category.Action = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "reputation-action "):
					siteReputation := tfdata.FirstElementOfJunosLine(itemTrim)
					category.ReputationAction = tfdata.AppendPotentialNewBlock(
						category.ReputationAction, types.StringValue(siteReputation),
					)

					if balt.CutPrefixInString(&itemTrim, siteReputation+" ") {
						category.ReputationAction[len(category.ReputationAction)-1].Action = types.StringValue(itemTrim)
					}
				}
			case balt.CutPrefixInString(&itemTrim, "custom-block-message "):
				rscData.CustomBlockMessage = types.StringValue(strings.Trim(itemTrim, "\""))
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
			case balt.CutPrefixInString(&itemTrim, "quarantine-custom-message "):
				rscData.QuarantineCustomMessage = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "quarantine-message"):
				if rscData.QuarantineMessage == nil {
					rscData.QuarantineMessage = &securityUtmProfileWebFilteringJuniperEnhancedBlockMessage{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					rscData.QuarantineMessage.read(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "site-reputation-action "):
				siteReputation := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.SiteReputationAction = tfdata.AppendPotentialNewBlock(
					rscData.SiteReputationAction, types.StringValue(siteReputation),
				)

				if balt.CutPrefixInString(&itemTrim, siteReputation+" ") {
					rscData.SiteReputationAction[len(rscData.SiteReputationAction)-1].Action = types.StringValue(itemTrim)
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

func (block *securityUtmProfileWebFilteringJuniperEnhancedBlockMessage) read(itemTrim string) {
	switch {
	case itemTrim == "type custom-redirect-url":
		block.TypeCustomRedirectURL = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "url "):
		block.URL = types.StringValue(strings.Trim(itemTrim, "\""))
	}
}

func (rscData *securityUtmProfileWebFilteringJuniperEnhancedData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security utm feature-profile web-filtering juniper-enhanced profile \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
