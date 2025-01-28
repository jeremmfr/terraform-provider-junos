package providerfwk

import (
	"context"
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
	_ resource.Resource                   = &securityDynamicAddressFeedServer{}
	_ resource.ResourceWithConfigure      = &securityDynamicAddressFeedServer{}
	_ resource.ResourceWithValidateConfig = &securityDynamicAddressFeedServer{}
	_ resource.ResourceWithImportState    = &securityDynamicAddressFeedServer{}
)

type securityDynamicAddressFeedServer struct {
	client *junos.Client
}

func newSecurityDynamicAddressFeedServerResource() resource.Resource {
	return &securityDynamicAddressFeedServer{}
}

func (rsc *securityDynamicAddressFeedServer) typeName() string {
	return providerName + "_security_dynamic_address_feed_server"
}

func (rsc *securityDynamicAddressFeedServer) junosName() string {
	return "security dynamic-address feed-server"
}

func (rsc *securityDynamicAddressFeedServer) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityDynamicAddressFeedServer) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityDynamicAddressFeedServer) Configure(
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

func (rsc *securityDynamicAddressFeedServer) Schema(
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
				Description: "Security dynamic address feed-server name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 16),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Description: "Hostname or IP address of feed-server.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"url": schema.StringAttribute{
				Optional:    true,
				Description: "URL.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of feed-server.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"hold_interval": schema.Int64Attribute{
				Optional:    true,
				Description: "Time to keep IP entry when update failed (seconds)",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
			"tls_profile": schema.StringAttribute{
				Optional:    true,
				Description: "TLS profile.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"update_interval": schema.Int64Attribute{
				Optional:    true,
				Description: " Interval to retrieve update (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(30, 4294967295),
				},
			},
			"validate_certificate_attributes_subject_or_san": schema.BoolAttribute{
				Optional: true,
				Description: "Validate certificate attributes subject or subject alternative names" +
					" must match feed server host name.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"feed_name": schema.ListNestedBlock{
				Description: "Feed name in feed-server.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Security dynamic address feed name in feed-server.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 240),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"path": schema.StringAttribute{
							Required:    true,
							Description: "Path of feed, appended to feed-server to form a complete URL.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 1024),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Text description of feed in feed-server.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 900),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"hold_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Time to keep IP entry when update failed (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(0, 4294967295),
							},
						},
						"update_interval": schema.Int64Attribute{
							Optional:    true,
							Description: "Interval to retrieve update (seconds).",
							Validators: []validator.Int64{
								int64validator.Between(30, 4294967295),
							},
						},
					},
				},
			},
		},
	}
}

//nolint:lll
type securityDynamicAddressFeedServerData struct {
	ID                                        types.String                                    `tfsdk:"id"`
	Name                                      types.String                                    `tfsdk:"name"`
	Hostname                                  types.String                                    `tfsdk:"hostname"`
	URL                                       types.String                                    `tfsdk:"url"`
	Description                               types.String                                    `tfsdk:"description"`
	HoldInterval                              types.Int64                                     `tfsdk:"hold_interval"`
	TLSProfile                                types.String                                    `tfsdk:"tls_profile"`
	UpdateInterval                            types.Int64                                     `tfsdk:"update_interval"`
	ValidateCertificateAttributesSubjectOrSan types.Bool                                      `tfsdk:"validate_certificate_attributes_subject_or_san"`
	FeedName                                  []securityDynamicAddressFeedServerBlockFeedName `tfsdk:"feed_name"`
}

type securityDynamicAddressFeedServerConfig struct {
	ID                                        types.String `tfsdk:"id"`
	Name                                      types.String `tfsdk:"name"`
	Hostname                                  types.String `tfsdk:"hostname"`
	URL                                       types.String `tfsdk:"url"`
	Description                               types.String `tfsdk:"description"`
	HoldInterval                              types.Int64  `tfsdk:"hold_interval"`
	TLSProfile                                types.String `tfsdk:"tls_profile"`
	UpdateInterval                            types.Int64  `tfsdk:"update_interval"`
	ValidateCertificateAttributesSubjectOrSan types.Bool   `tfsdk:"validate_certificate_attributes_subject_or_san"`
	FeedName                                  types.List   `tfsdk:"feed_name"`
}

type securityDynamicAddressFeedServerBlockFeedName struct {
	Name           types.String `tfsdk:"name"            tfdata:"identifier"`
	Path           types.String `tfsdk:"path"`
	Description    types.String `tfsdk:"description"`
	HoldInterval   types.Int64  `tfsdk:"hold_interval"`
	UpdateInterval types.Int64  `tfsdk:"update_interval"`
}

func (rsc *securityDynamicAddressFeedServer) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityDynamicAddressFeedServerConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Hostname.IsNull() &&
		config.URL.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hostname"),
			tfdiag.MissingConfigErrSummary,
			"one of hostname or url must be specified",
		)
	}
	if !config.Hostname.IsNull() && !config.Hostname.IsUnknown() &&
		!config.URL.IsNull() && !config.URL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hostname"),
			tfdiag.ConflictConfigErrSummary,
			"only one of hostname or url must be specified",
		)
	}
	if !config.ValidateCertificateAttributesSubjectOrSan.IsNull() &&
		!config.ValidateCertificateAttributesSubjectOrSan.IsUnknown() &&
		config.TLSProfile.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("validate_certificate_attributes_subject_or_san"),
			tfdiag.MissingConfigErrSummary,
			"tls_profile must be specified with validate_certificate_attributes_subject_or_san",
		)
	}

	if !config.FeedName.IsNull() &&
		!config.FeedName.IsUnknown() {
		var configFeedName []securityDynamicAddressFeedServerBlockFeedName
		asDiags := config.FeedName.ElementsAs(ctx, &configFeedName, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		feedNameName := make(map[string]struct{})
		for i, block := range configFeedName {
			if block.Name.IsUnknown() {
				continue
			}

			name := block.Name.ValueString()
			if _, ok := feedNameName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("feed_name").AtListIndex(i).AtName("name"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple feed_name blocks with the same name %q", name),
				)
			}
			feedNameName[name] = struct{}{}
		}
	}
}

func (rsc *securityDynamicAddressFeedServer) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityDynamicAddressFeedServerData
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
			serverExists, err := checkSecurityDynamicAddressFeedServerExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if serverExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			serverExists, err := checkSecurityDynamicAddressFeedServerExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !serverExists {
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

func (rsc *securityDynamicAddressFeedServer) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityDynamicAddressFeedServerData
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

func (rsc *securityDynamicAddressFeedServer) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityDynamicAddressFeedServerData
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

func (rsc *securityDynamicAddressFeedServer) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityDynamicAddressFeedServerData
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

func (rsc *securityDynamicAddressFeedServer) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityDynamicAddressFeedServerData

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

func checkSecurityDynamicAddressFeedServerExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security dynamic-address feed-server " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityDynamicAddressFeedServerData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityDynamicAddressFeedServerData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityDynamicAddressFeedServerData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security dynamic-address feed-server " + rscData.Name.ValueString() + " "

	if v := rscData.Hostname.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"hostname \""+v+"\"")
	}
	if v := rscData.URL.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"url \""+v+"\"")
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if !rscData.HoldInterval.IsNull() {
		configSet = append(configSet, setPrefix+"hold-interval "+
			utils.ConvI64toa(rscData.HoldInterval.ValueInt64()))
	}
	if v := rscData.TLSProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"tls-profile \""+v+"\"")
	}
	if !rscData.UpdateInterval.IsNull() {
		configSet = append(configSet, setPrefix+"update-interval "+
			utils.ConvI64toa(rscData.UpdateInterval.ValueInt64()))
	}
	if rscData.ValidateCertificateAttributesSubjectOrSan.ValueBool() {
		configSet = append(configSet, setPrefix+"validate-certificate-attributes subject-or-subject-alternative-names")
	}

	feedNameName := make(map[string]struct{})
	for i, block := range rscData.FeedName {
		name := block.Name.ValueString()
		if _, ok := feedNameName[name]; ok {
			return path.Root("feed_name").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple feed_name blocks with the same name %q", name)
		}
		feedNameName[name] = struct{}{}

		configSet = append(configSet, block.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityDynamicAddressFeedServerBlockFeedName) configSet(setPrefix string) []string {
	setPrefix += "feed-name " + block.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
		setPrefix + "path \"" + block.Path.ValueString() + "\"",
	}

	if v := block.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if !block.HoldInterval.IsNull() {
		configSet = append(configSet, setPrefix+"hold-interval "+
			utils.ConvI64toa(block.HoldInterval.ValueInt64()))
	}
	if !block.UpdateInterval.IsNull() {
		configSet = append(configSet, setPrefix+"update-interval "+
			utils.ConvI64toa(block.UpdateInterval.ValueInt64()))
	}

	return configSet
}

func (rscData *securityDynamicAddressFeedServerData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security dynamic-address feed-server " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "hostname "):
				rscData.Hostname = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "url "):
				rscData.URL = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "feed-name "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.FeedName = tfdata.AppendPotentialNewBlock(rscData.FeedName, types.StringValue(name))
				feedName := &rscData.FeedName[len(rscData.FeedName)-1]
				balt.CutPrefixInString(&itemTrim, name+" ")

				if err := feedName.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "hold-interval "):
				rscData.HoldInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "tls-profile "):
				rscData.TLSProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "update-interval "):
				rscData.UpdateInterval, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case itemTrim == "validate-certificate-attributes subject-or-subject-alternative-names":
				rscData.ValidateCertificateAttributesSubjectOrSan = types.BoolValue(true)
			}
		}
	}

	return nil
}

func (block *securityDynamicAddressFeedServerBlockFeedName) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "path "):
		block.Path = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "description "):
		block.Description = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "hold-interval "):
		block.HoldInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "update-interval "):
		block.UpdateInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (rscData *securityDynamicAddressFeedServerData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security dynamic-address feed-server " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
