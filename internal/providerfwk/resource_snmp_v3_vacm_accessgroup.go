package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

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
	_ resource.Resource                   = &snmpV3VacmAccessgroup{}
	_ resource.ResourceWithConfigure      = &snmpV3VacmAccessgroup{}
	_ resource.ResourceWithValidateConfig = &snmpV3VacmAccessgroup{}
	_ resource.ResourceWithImportState    = &snmpV3VacmAccessgroup{}
)

type snmpV3VacmAccessgroup struct {
	client *junos.Client
}

func newSnmpV3VacmAccessgroupResource() resource.Resource {
	return &snmpV3VacmAccessgroup{}
}

func (rsc *snmpV3VacmAccessgroup) typeName() string {
	return providerName + "_snmp_v3_vacm_accessgroup"
}

func (rsc *snmpV3VacmAccessgroup) junosName() string {
	return "snmp v3 vacm access group"
}

func (rsc *snmpV3VacmAccessgroup) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *snmpV3VacmAccessgroup) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *snmpV3VacmAccessgroup) Configure(
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

func (rsc *snmpV3VacmAccessgroup) Schema(
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
				Description: "SNMPv3 VACM group name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"context_prefix": schema.ListNestedBlock{
				Description: "For each prefix of context-prefix access configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"prefix": schema.StringAttribute{
							Required:    true,
							Description: "SNMPv3 VACM context prefix.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"access_config": schema.SetNestedBlock{
							Description: "For each combination of `model` and `level`, define context-prefix access configuration.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"model": schema.StringAttribute{
										Required:    true,
										Description: "Security model access configuration.",
										Validators: []validator.String{
											stringvalidator.OneOf("any", "usm", "v1", "v2c"),
										},
									},
									"level": schema.StringAttribute{
										Required:    true,
										Description: "Security level access configuration.",
										Validators: []validator.String{
											stringvalidator.OneOf("authentication", "none", "privacy"),
										},
									},
									"context_match": schema.StringAttribute{
										Optional:    true,
										Description: "Type of match to perform on context-prefix.",
										Validators: []validator.String{
											stringvalidator.OneOf("exact", "prefix"),
										},
									},
									"notify_view": schema.StringAttribute{
										Optional:    true,
										Description: "View used to notifications.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 32),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"read_view": schema.StringAttribute{
										Optional:    true,
										Description: "View used for read access.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 32),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
									"write_view": schema.StringAttribute{
										Optional:    true,
										Description: "View used for write access.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 32),
											tfvalidator.StringDoubleQuoteExclusion(),
										},
									},
								},
							},
						},
					},
				},
			},
			"default_context_prefix": schema.SetNestedBlock{
				Description: "For each combination of `model` and `level`, define default context-prefix access configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"model": schema.StringAttribute{
							Required:    true,
							Description: "Security model access configuration.",
							Validators: []validator.String{
								stringvalidator.OneOf("any", "usm", "v1", "v2c"),
							},
						},
						"level": schema.StringAttribute{
							Required:    true,
							Description: "Security level access configuration.",
							Validators: []validator.String{
								stringvalidator.OneOf("authentication", "none", "privacy"),
							},
						},
						"context_match": schema.StringAttribute{
							Optional:    true,
							Description: "Type of match to perform on context-prefix.",
							Validators: []validator.String{
								stringvalidator.OneOf("exact", "prefix"),
							},
						},
						"notify_view": schema.StringAttribute{
							Optional:    true,
							Description: "View used to notifications.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"read_view": schema.StringAttribute{
							Optional:    true,
							Description: "View used for read access.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"write_view": schema.StringAttribute{
							Optional:    true,
							Description: "View used for write access.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
				},
			},
		},
	}
}

type snmpV3VacmAccessgroupData struct {
	ID                   types.String                                               `tfsdk:"id"`
	Name                 types.String                                               `tfsdk:"name"`
	ContextPrefix        []snmpV3VacmAccessgroupBlockContextPrefix                  `tfsdk:"context_prefix"`
	DefaultContextPrefix []snmpV3VacmAccessgroupBlockContextPrefixBlockAccessConfig `tfsdk:"default_context_prefix"`
}

type snmpV3VacmAccessgroupConfig struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	ContextPrefix        types.List   `tfsdk:"context_prefix"`
	DefaultContextPrefix types.Set    `tfsdk:"default_context_prefix"`
}

type snmpV3VacmAccessgroupBlockContextPrefix struct {
	Prefix       types.String                                               `tfsdk:"prefix"`
	AccessConfig []snmpV3VacmAccessgroupBlockContextPrefixBlockAccessConfig `tfsdk:"access_config"`
}

type snmpV3VacmAccessgroupBlockContextPrefixConfig struct {
	Prefix       types.String `tfsdk:"prefix"`
	AccessConfig types.Set    `tfsdk:"access_config"`
}

type snmpV3VacmAccessgroupBlockContextPrefixBlockAccessConfig struct {
	Model        types.String `tfsdk:"model"`
	Level        types.String `tfsdk:"level"`
	ContextMatch types.String `tfsdk:"context_match"`
	NotifyView   types.String `tfsdk:"notify_view"`
	ReadView     types.String `tfsdk:"read_view"`
	WriteView    types.String `tfsdk:"write_view"`
}

func (block *snmpV3VacmAccessgroupBlockContextPrefixBlockAccessConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block, "Model", "Level")
}

func (rsc *snmpV3VacmAccessgroup) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config snmpV3VacmAccessgroupConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ContextPrefix.IsNull() &&
		config.DefaultContextPrefix.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.MissingConfigErrSummary,
			"at least one of context_prefix or default_context_prefix must be specified",
		)
	}
	if !config.ContextPrefix.IsNull() && !config.ContextPrefix.IsUnknown() {
		var contextPrefix []snmpV3VacmAccessgroupBlockContextPrefixConfig
		asDiags := config.ContextPrefix.ElementsAs(ctx, &contextPrefix, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		contextPrefixPrefix := make(map[string]struct{})
		for i, block := range contextPrefix {
			if !block.Prefix.IsUnknown() {
				prefix := block.Prefix.ValueString()
				if _, ok := contextPrefixPrefix[prefix]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("context_prefix").AtListIndex(i).AtName("prefix"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple context_prefix blocks with the same prefix %q", prefix),
					)
				}
				contextPrefixPrefix[prefix] = struct{}{}
			}

			if block.AccessConfig.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("context_prefix").AtListIndex(i).AtName("prefix"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("access_config block must be specified"+
						" in context_prefix block %q", block.Prefix.ValueString()),
				)
			} else if !block.AccessConfig.IsUnknown() {
				var accessConfig []snmpV3VacmAccessgroupBlockContextPrefixBlockAccessConfig
				asDiags := block.AccessConfig.ElementsAs(ctx, &accessConfig, false)
				if asDiags.HasError() {
					resp.Diagnostics.Append(asDiags...)

					return
				}

				accessConfigModelLevel := make(map[string]struct{})
				for _, blockAccessConfig := range accessConfig {
					if blockAccessConfig.isEmpty() {
						resp.Diagnostics.AddAttributeError(
							path.Root("context_prefix").AtListIndex(i).AtName("prefix"),
							tfdiag.MissingConfigErrSummary,
							fmt.Sprintf(
								"access_config block model:%q level:%q is empty"+
									" in context_prefix block %q",
								blockAccessConfig.Model.ValueString(), blockAccessConfig.Level.ValueString(), block.Prefix.ValueString(),
							),
						)
					}
					if blockAccessConfig.Model.IsUnknown() {
						continue
					}
					if blockAccessConfig.Level.IsUnknown() {
						continue
					}
					model := blockAccessConfig.Model.ValueString()
					level := blockAccessConfig.Level.ValueString()
					if _, ok := accessConfigModelLevel[model+junos.IDSeparator+level]; ok {
						resp.Diagnostics.AddAttributeError(
							path.Root("context_prefix").AtListIndex(i).AtName("prefix"),
							tfdiag.DuplicateConfigErrSummary,
							fmt.Sprintf("multiple access_config blocks with the same model %q and level %q"+
								" in context_prefix block %q", model, level, block.Prefix.ValueString()),
						)
					}
					accessConfigModelLevel[model+junos.IDSeparator+level] = struct{}{}
				}
			}
		}
	}
	if !config.DefaultContextPrefix.IsNull() && !config.DefaultContextPrefix.IsUnknown() {
		var defaultContextPrefix []snmpV3VacmAccessgroupBlockContextPrefixBlockAccessConfig
		asDiags := config.DefaultContextPrefix.ElementsAs(ctx, &defaultContextPrefix, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		defaultContextPrefixModelLevel := make(map[string]struct{})
		for _, block := range defaultContextPrefix {
			if block.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("default_context_prefix"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf(
						"default_context_prefix block model:%q level:%q is empty",
						block.Model.ValueString(), block.Level.ValueString(),
					),
				)
			}
			if block.Model.IsUnknown() {
				continue
			}
			if block.Level.IsUnknown() {
				continue
			}
			model := block.Model.ValueString()
			level := block.Level.ValueString()
			if _, ok := defaultContextPrefixModelLevel[model+junos.IDSeparator+level]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("default_context_prefix"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple default_context_prefix blocks with the same model %q and level %q", model, level),
				)
			}
			defaultContextPrefixModelLevel[model+junos.IDSeparator+level] = struct{}{}
		}
	}
}

func (rsc *snmpV3VacmAccessgroup) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan snmpV3VacmAccessgroupData
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
			groupExists, err := checkSnmpV3VacmAccessgroupExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if groupExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			groupExists, err := checkSnmpV3VacmAccessgroupExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !groupExists {
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

func (rsc *snmpV3VacmAccessgroup) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data snmpV3VacmAccessgroupData
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

func (rsc *snmpV3VacmAccessgroup) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state snmpV3VacmAccessgroupData
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

func (rsc *snmpV3VacmAccessgroup) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state snmpV3VacmAccessgroupData
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

func (rsc *snmpV3VacmAccessgroup) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data snmpV3VacmAccessgroupData

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

func checkSnmpV3VacmAccessgroupExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp v3 vacm access group \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *snmpV3VacmAccessgroupData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *snmpV3VacmAccessgroupData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *snmpV3VacmAccessgroupData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set snmp v3 vacm access group \"" + rscData.Name.ValueString() + "\" "

	defaultContextPrefixModelLevel := make(map[string]struct{})
	for _, block := range rscData.DefaultContextPrefix {
		model := block.Model.ValueString()
		level := block.Level.ValueString()
		if _, ok := defaultContextPrefixModelLevel[model+junos.IDSeparator+level]; ok {
			return path.Root("default_context_prefix"),
				fmt.Errorf("multiple default_context_prefix blocks with the same model %q and level %q", model, level)
		}
		defaultContextPrefixModelLevel[model+junos.IDSeparator+level] = struct{}{}
		if block.isEmpty() {
			return path.Root("default_context_prefix"),
				fmt.Errorf("default_context_prefix block model:%q level:%q is empty", model, level)
		}

		setPrefixBlock := setPrefix + "default-context-prefix security-model " + model + " security-level " + level + " "
		if v := block.ContextMatch.ValueString(); v != "" {
			configSet = append(configSet, setPrefixBlock+"context-match "+v)
		}
		if v := block.NotifyView.ValueString(); v != "" {
			configSet = append(configSet, setPrefixBlock+"notify-view \""+v+"\"")
		}
		if v := block.ReadView.ValueString(); v != "" {
			configSet = append(configSet, setPrefixBlock+"read-view \""+v+"\"")
		}
		if v := block.WriteView.ValueString(); v != "" {
			configSet = append(configSet, setPrefixBlock+"write-view \""+v+"\"")
		}
	}
	contextPrefixPrefix := make(map[string]struct{})
	for i, block := range rscData.ContextPrefix {
		prefix := block.Prefix.ValueString()
		if _, ok := contextPrefixPrefix[prefix]; ok {
			return path.Root("context_prefix").AtListIndex(i).AtName("prefix"),
				fmt.Errorf("multiple context_prefix blocks with the same prefix %q", prefix)
		}
		contextPrefixPrefix[prefix] = struct{}{}

		accessConfigModelLevel := make(map[string]struct{})
		for _, blockAccessConfig := range block.AccessConfig {
			model := blockAccessConfig.Model.ValueString()
			level := blockAccessConfig.Level.ValueString()
			if _, ok := accessConfigModelLevel[model+junos.IDSeparator+level]; ok {
				return path.Root("context_prefix").AtListIndex(i).AtName("access_config"),
					fmt.Errorf("multiple access_config blocks with the same model %q and level %q"+
						" in context_prefix block %q", model, level, prefix)
			}
			accessConfigModelLevel[model+junos.IDSeparator+level] = struct{}{}
			if blockAccessConfig.isEmpty() {
				return path.Root("context_prefix").AtListIndex(i).AtName("access_config"),
					fmt.Errorf("access_config block model:%q level:%q is empty"+
						" in context_prefix block %q", model, level, prefix)
			}

			setPrefixBlock := setPrefix + " context-prefix \"" + prefix + "\" " +
				"security-model " + model + " security-level " + level + " "
			if v := blockAccessConfig.ContextMatch.ValueString(); v != "" {
				configSet = append(configSet, setPrefixBlock+"context-match "+v)
			}
			if v := blockAccessConfig.NotifyView.ValueString(); v != "" {
				configSet = append(configSet, setPrefixBlock+"notify-view \""+v+"\"")
			}
			if v := blockAccessConfig.ReadView.ValueString(); v != "" {
				configSet = append(configSet, setPrefixBlock+"read-view \""+v+"\"")
			}
			if v := blockAccessConfig.WriteView.ValueString(); v != "" {
				configSet = append(configSet, setPrefixBlock+"write-view \""+v+"\"")
			}
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *snmpV3VacmAccessgroupData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp v3 vacm access group \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "default-context-prefix security-model "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 3 { // <model> security-level <level>
					return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "default-context-prefix security-model", itemTrim)
				}
				var defaultContextPrefix snmpV3VacmAccessgroupBlockContextPrefixBlockAccessConfig
				rscData.DefaultContextPrefix, defaultContextPrefix = tfdata.ExtractBlockWith2TFTypesString(
					rscData.DefaultContextPrefix, "Model", itemTrimFields[0], "Level", itemTrimFields[2])
				defaultContextPrefix.Model = types.StringValue(itemTrimFields[0])
				defaultContextPrefix.Level = types.StringValue(itemTrimFields[2])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" security-level "+itemTrimFields[2]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "context-match "):
					defaultContextPrefix.ContextMatch = types.StringValue(itemTrim)
				case balt.CutPrefixInString(&itemTrim, "notify-view "):
					defaultContextPrefix.NotifyView = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "read-view "):
					defaultContextPrefix.ReadView = types.StringValue(strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "write-view "):
					defaultContextPrefix.WriteView = types.StringValue(strings.Trim(itemTrim, "\""))
				}
				rscData.DefaultContextPrefix = append(rscData.DefaultContextPrefix, defaultContextPrefix)
			case balt.CutPrefixInString(&itemTrim, "context-prefix "):
				prefix := tfdata.FirstElementOfJunosLine(itemTrim)
				var contextPrefix snmpV3VacmAccessgroupBlockContextPrefix
				rscData.ContextPrefix, contextPrefix = tfdata.ExtractBlockWithTFTypesString(
					rscData.ContextPrefix, "Prefix", strings.Trim(prefix, "\""))
				contextPrefix.Prefix = types.StringValue(strings.Trim(prefix, "\""))
				balt.CutPrefixInString(&itemTrim, prefix+" ")
				if balt.CutPrefixInString(&itemTrim, "security-model ") {
					itemTrimFields := strings.Split(itemTrim, " ")
					if len(itemTrimFields) < 3 { // <model> security-level <level>
						return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "security-model", itemTrim)
					}
					var accessConfig snmpV3VacmAccessgroupBlockContextPrefixBlockAccessConfig
					contextPrefix.AccessConfig, accessConfig = tfdata.ExtractBlockWith2TFTypesString(
						contextPrefix.AccessConfig, "Model", itemTrimFields[0], "Level", itemTrimFields[2])
					accessConfig.Model = types.StringValue(itemTrimFields[0])
					accessConfig.Level = types.StringValue(itemTrimFields[2])
					balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" security-level "+itemTrimFields[2]+" ")
					switch {
					case balt.CutPrefixInString(&itemTrim, "context-match "):
						accessConfig.ContextMatch = types.StringValue(itemTrim)
					case balt.CutPrefixInString(&itemTrim, "notify-view "):
						accessConfig.NotifyView = types.StringValue(strings.Trim(itemTrim, "\""))
					case balt.CutPrefixInString(&itemTrim, "read-view "):
						accessConfig.ReadView = types.StringValue(strings.Trim(itemTrim, "\""))
					case balt.CutPrefixInString(&itemTrim, "write-view "):
						accessConfig.WriteView = types.StringValue(strings.Trim(itemTrim, "\""))
					}
					contextPrefix.AccessConfig = append(contextPrefix.AccessConfig, accessConfig)
				}
				rscData.ContextPrefix = append(rscData.ContextPrefix, contextPrefix)
			}
		}
	}

	return nil
}

func (rscData *snmpV3VacmAccessgroupData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete snmp v3 vacm access group \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
