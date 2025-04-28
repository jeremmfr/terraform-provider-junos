package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                   = &securityDynamicAddressName{}
	_ resource.ResourceWithConfigure      = &securityDynamicAddressName{}
	_ resource.ResourceWithValidateConfig = &securityDynamicAddressName{}
	_ resource.ResourceWithImportState    = &securityDynamicAddressName{}
	_ resource.ResourceWithUpgradeState   = &securityDynamicAddressName{}
)

type securityDynamicAddressName struct {
	client *junos.Client
}

func newSecurityDynamicAddressNameResource() resource.Resource {
	return &securityDynamicAddressName{}
}

func (rsc *securityDynamicAddressName) typeName() string {
	return providerName + "_security_dynamic_address_name"
}

func (rsc *securityDynamicAddressName) junosName() string {
	return "security dynamic-address address-name"
}

func (rsc *securityDynamicAddressName) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityDynamicAddressName) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityDynamicAddressName) Configure(
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

func (rsc *securityDynamicAddressName) Schema(
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
				Description: "Security dynamic address name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 240),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of dynamic address.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"profile_feed_name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of feed in feed-server for this dynamic address.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 240),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"session_scan": schema.BoolAttribute{
				Optional:    true,
				Description: "Trigger session scan when the IP address is added.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"profile_category": schema.SingleNestedBlock{
				Description: "Declare `profile category` configuration to categorize feed data into this dynamic address.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Name of category.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 32),
							tfvalidator.StringFormat(tfvalidator.DefaultFormat),
						},
					},
					"feed": schema.StringAttribute{
						Optional:    true,
						Description: "Name of feed under category.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 32),
							tfvalidator.StringFormat(tfvalidator.DefaultFormat),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"property": schema.ListNestedBlock{
						Description: "For each name of property to match.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "Name of property.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(1, 32),
										tfvalidator.StringDoubleQuoteExclusion(),
									},
								},
								"string": schema.ListAttribute{
									ElementType: types.StringType,
									Required:    true,
									Description: "String value.",
									Validators: []validator.List{
										listvalidator.SizeAtLeast(1),
										listvalidator.NoNullValues(),
										listvalidator.ValueStringsAre(
											stringvalidator.LengthBetween(1, 32),
											tfvalidator.StringDoubleQuoteExclusion(),
										),
									},
								},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeBetween(1, 3),
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

type securityDynamicAddressNameData struct {
	ID              types.String                                    `tfsdk:"id"`
	Name            types.String                                    `tfsdk:"name"`
	Description     types.String                                    `tfsdk:"description"`
	ProfileFeedName types.String                                    `tfsdk:"profile_feed_name"`
	SessionScan     types.Bool                                      `tfsdk:"session_scan"`
	ProfileCategory *securityDynamicAddressNameBlockProfileCategory `tfsdk:"profile_category"`
}

type securityDynamicAddressNameConfig struct {
	ID              types.String                                          `tfsdk:"id"`
	Name            types.String                                          `tfsdk:"name"`
	Description     types.String                                          `tfsdk:"description"`
	ProfileFeedName types.String                                          `tfsdk:"profile_feed_name"`
	SessionScan     types.Bool                                            `tfsdk:"session_scan"`
	ProfileCategory *securityDynamicAddressNameBlockProfileCategoryConfig `tfsdk:"profile_category"`
}

type securityDynamicAddressNameBlockProfileCategory struct {
	Name     types.String                                                  `tfsdk:"name"`
	Feed     types.String                                                  `tfsdk:"feed"`
	Property []securityDynamicAddressNameBlockProfileCategoryBlockProperty `tfsdk:"property"`
}

type securityDynamicAddressNameBlockProfileCategoryConfig struct {
	Name     types.String `tfsdk:"name"`
	Feed     types.String `tfsdk:"feed"`
	Property types.List   `tfsdk:"property"`
}

func (block *securityDynamicAddressNameBlockProfileCategoryConfig) hasKnownValue() bool {
	return tfdata.CheckBlockHasKnownValue(block)
}

type securityDynamicAddressNameBlockProfileCategoryBlockProperty struct {
	Name   types.String   `tfsdk:"name"   tfdata:"identifier"`
	String []types.String `tfsdk:"string"`
}

func (rsc *securityDynamicAddressName) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityDynamicAddressNameConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ProfileFeedName.IsNull() && config.ProfileCategory == nil {
		resp.Diagnostics.AddError(
			tfdiag.MissingConfigErrSummary,
			"one of profile_feed_name or profile_category must be specified",
		)
	}
	if !config.ProfileFeedName.IsNull() && !config.ProfileFeedName.IsUnknown() &&
		config.ProfileCategory != nil && config.ProfileCategory.hasKnownValue() {
		resp.Diagnostics.AddAttributeError(
			path.Root("profile_feed_name"),
			tfdiag.ConflictConfigErrSummary,
			"only one of profile_feed_name or profile_category must be specified",
		)
	}

	if config.ProfileCategory != nil {
		if config.ProfileCategory.Name.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("profile_category").AtName("name"),
				tfdiag.MissingConfigErrSummary,
				"name must be specified in profile_category block",
			)
		}

		if !config.ProfileCategory.Property.IsNull() &&
			!config.ProfileCategory.Property.IsUnknown() {
			var configProperty []securityDynamicAddressNameBlockProfileCategoryBlockProperty
			asDiags := config.ProfileCategory.Property.ElementsAs(ctx, &configProperty, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			propertyName := make(map[string]struct{})
			for i, block := range configProperty {
				if block.Name.IsUnknown() {
					continue
				}

				name := block.Name.ValueString()
				if _, ok := propertyName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("profile_category").AtName("property").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple property blocks with the same name %q"+
							" in profile_category block", name),
					)
				}
				propertyName[name] = struct{}{}
			}
		}
	}
}

func (rsc *securityDynamicAddressName) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityDynamicAddressNameData
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
			addressExists, err := checkSecurityDynamicAddressNameExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if addressExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			addressExists, err := checkSecurityDynamicAddressNameExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !addressExists {
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

func (rsc *securityDynamicAddressName) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityDynamicAddressNameData
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

func (rsc *securityDynamicAddressName) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityDynamicAddressNameData
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

func (rsc *securityDynamicAddressName) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityDynamicAddressNameData
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

func (rsc *securityDynamicAddressName) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityDynamicAddressNameData

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

func checkSecurityDynamicAddressNameExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security dynamic-address address-name " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityDynamicAddressNameData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityDynamicAddressNameData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityDynamicAddressNameData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)
	setPrefix := "set security dynamic-address address-name " + rscData.Name.ValueString() + " "

	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := rscData.ProfileFeedName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"profile feed-name "+v)
	}
	if rscData.SessionScan.ValueBool() {
		configSet = append(configSet, setPrefix+"session-scan")
	}

	if rscData.ProfileCategory != nil {
		blockSet, pathErr, err := rscData.ProfileCategory.configSet(setPrefix)
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityDynamicAddressNameBlockProfileCategory) configSet(
	setPrefix string,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "profile category " + block.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
	}

	if v := block.Feed.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"feed "+v)
	}
	propertyName := make(map[string]struct{})
	for i, subBlock := range block.Property {
		name := subBlock.Name.ValueString()
		if _, ok := propertyName[name]; ok {
			return configSet,
				path.Root("profile_category").AtName("property").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple property blocks with the same name %q"+
					" in profile_category block", name)
		}
		propertyName[name] = struct{}{}

		for _, v := range subBlock.String {
			configSet = append(configSet, setPrefix+"property \""+name+"\""+
				" string \""+v.ValueString()+"\"")
		}
	}

	return configSet, path.Empty(), nil
}

func (rscData *securityDynamicAddressNameData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security dynamic-address address-name " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "profile feed-name "):
				rscData.ProfileFeedName = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "profile category "):
				if rscData.ProfileCategory == nil {
					rscData.ProfileCategory = &securityDynamicAddressNameBlockProfileCategory{}
				}

				profileCategoryName := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.ProfileCategory.Name = types.StringValue(profileCategoryName)
				if balt.CutPrefixInString(&itemTrim, profileCategoryName+" ") {
					if err := rscData.ProfileCategory.read(itemTrim); err != nil {
						return err
					}
				}
			case itemTrim == "session-scan":
				rscData.SessionScan = types.BoolValue(true)
			}
		}
	}

	return nil
}

func (block *securityDynamicAddressNameBlockProfileCategory) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "feed "):
		block.Feed = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "property "):
		name := tfdata.FirstElementOfJunosLine(itemTrim)
		block.Property = tfdata.AppendPotentialNewBlock(block.Property, types.StringValue(strings.Trim(name, "\"")))
		property := &block.Property[len(block.Property)-1]
		balt.CutPrefixInString(&itemTrim, name+" ")

		if balt.CutPrefixInString(&itemTrim, "string ") {
			property.String = append(property.String, types.StringValue(strings.Trim(itemTrim, "\"")))
		}
	}

	return nil
}

func (rscData *securityDynamicAddressNameData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security dynamic-address address-name " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
