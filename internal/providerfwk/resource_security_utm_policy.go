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
	_ resource.Resource                   = &securityUtmPolicy{}
	_ resource.ResourceWithConfigure      = &securityUtmPolicy{}
	_ resource.ResourceWithValidateConfig = &securityUtmPolicy{}
	_ resource.ResourceWithImportState    = &securityUtmPolicy{}
	_ resource.ResourceWithUpgradeState   = &securityUtmPolicy{}
)

type securityUtmPolicy struct {
	client *junos.Client
}

func newSecurityUtmPolicyResource() resource.Resource {
	return &securityUtmPolicy{}
}

func (rsc *securityUtmPolicy) typeName() string {
	return providerName + "_security_utm_policy"
}

func (rsc *securityUtmPolicy) junosName() string {
	return "security utm utm-policy"
}

func (rsc *securityUtmPolicy) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityUtmPolicy) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityUtmPolicy) Configure(
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

func (rsc *securityUtmPolicy) Schema(
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
				Description: "The name of security utm utm-policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 29),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"anti_spam_smtp_profile": schema.StringAttribute{
				Optional:    true,
				Description: "Name of anti-spam profile.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"web_filtering_profile": schema.StringAttribute{
				Optional:    true,
				Description: "Web-filtering HTTP profile (local, enhanced, websense).",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"anti_virus": schema.SingleNestedBlock{
				Description: "Configure for utm anti-virus profile.",
				Attributes: map[string]schema.Attribute{
					"ftp_download_profile": schema.StringAttribute{
						Optional:    true,
						Description: "FTP download anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"ftp_upload_profile": schema.StringAttribute{
						Optional:    true,
						Description: "FTP upload anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"http_profile": schema.StringAttribute{
						Optional:    true,
						Description: "HTTP anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"imap_profile": schema.StringAttribute{
						Optional:    true,
						Description: "IMAP anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"pop3_profile": schema.StringAttribute{
						Optional:    true,
						Description: "POP3 anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"smtp_profile": schema.StringAttribute{
						Optional:    true,
						Description: "SMTP anti-virus profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"content_filtering": schema.SingleNestedBlock{
				Description: "Configure for utm content-filtering profile.",
				Attributes: map[string]schema.Attribute{
					"ftp_download_profile": schema.StringAttribute{
						Optional:    true,
						Description: "FTP download content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"ftp_upload_profile": schema.StringAttribute{
						Optional:    true,
						Description: "FTP upload content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"http_profile": schema.StringAttribute{
						Optional:    true,
						Description: "HTTP content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"imap_profile": schema.StringAttribute{
						Optional:    true,
						Description: "HTTP content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"pop3_profile": schema.StringAttribute{
						Optional:    true,
						Description: "POP3 content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"smtp_profile": schema.StringAttribute{
						Optional:    true,
						Description: "SMTP content-filtering profile.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"traffic_sessions_per_client": schema.SingleNestedBlock{
				Description: " Configure for traffic option session per client.",
				Attributes: map[string]schema.Attribute{
					"limit": schema.Int64Attribute{
						Optional:    true,
						Description: "Sessions limit.",
						Validators: []validator.Int64{
							int64validator.Between(0, 2000),
						},
					},
					"over_limit": schema.StringAttribute{
						Optional:    true,
						Description: "Over limit action.",
						Validators: []validator.String{
							stringvalidator.OneOf("block", "log-and-permit"),
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

type securityUtmPolicyData struct {
	ID                       types.String                                    `tfsdk:"id"`
	Name                     types.String                                    `tfsdk:"name"`
	AntiSpamSMTPProfile      types.String                                    `tfsdk:"anti_spam_smtp_profile"`
	WebFilteringProfile      types.String                                    `tfsdk:"web_filtering_profile"`
	AntiVirus                *securityUtmPolicyBlockProtocolProfile          `tfsdk:"anti_virus"`
	ContentFiltering         *securityUtmPolicyBlockProtocolProfile          `tfsdk:"content_filtering"`
	TrafficSessionsPerClient *securityUtmPolicyBlockTrafficSessionsPerClient `tfsdk:"traffic_sessions_per_client"`
}

func (rscData *securityUtmPolicyData) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(rscData)
}

type securityUtmPolicyBlockProtocolProfile struct {
	FTPDownloadProfile types.String `tfsdk:"ftp_download_profile"`
	FTPUploadProfile   types.String `tfsdk:"ftp_upload_profile"`
	HTTPProfile        types.String `tfsdk:"http_profile"`
	IMAPProfile        types.String `tfsdk:"imap_profile"`
	POP3Profile        types.String `tfsdk:"pop3_profile"`
	SMTPProfile        types.String `tfsdk:"smtp_profile"`
}

func (block *securityUtmPolicyBlockProtocolProfile) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type securityUtmPolicyBlockTrafficSessionsPerClient struct {
	Limit     types.Int64  `tfsdk:"limit"`
	OverLimit types.String `tfsdk:"over_limit"`
}

func (block *securityUtmPolicyBlockTrafficSessionsPerClient) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (rsc *securityUtmPolicy) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityUtmPolicyData
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

	if config.AntiVirus != nil &&
		config.AntiVirus.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("anti_virus").AtName("*"),
			tfdiag.MissingConfigErrSummary,
			"anti_virus block is empty",
		)
	}
	if config.ContentFiltering != nil &&
		config.ContentFiltering.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("content_filtering").AtName("*"),
			tfdiag.MissingConfigErrSummary,
			"content_filtering block is empty",
		)
	}
	if config.TrafficSessionsPerClient != nil &&
		config.TrafficSessionsPerClient.isEmpty() {
		resp.Diagnostics.AddAttributeError(
			path.Root("traffic_sessions_per_client").AtName("*"),
			tfdiag.MissingConfigErrSummary,
			"traffic_sessions_per_client block is empty",
		)
	}
}

func (rsc *securityUtmPolicy) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityUtmPolicyData
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
			policyExists, err := checkSecurityUtmPolicyExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if policyExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			policyExists, err := checkSecurityUtmPolicyExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !policyExists {
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

func (rsc *securityUtmPolicy) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityUtmPolicyData
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

func (rsc *securityUtmPolicy) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityUtmPolicyData
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

func (rsc *securityUtmPolicy) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityUtmPolicyData
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

func (rsc *securityUtmPolicy) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityUtmPolicyData

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

func checkSecurityUtmPolicyExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm utm-policy \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityUtmPolicyData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityUtmPolicyData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityUtmPolicyData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	if rscData.isEmpty() {
		return path.Root("name"),
			errors.New("at least one of arguments need to be set (in addition to `name`)")
	}

	configSet := make([]string, 0, 100)
	setPrefix := "set security utm utm-policy \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.AntiSpamSMTPProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"anti-spam smtp-profile \""+v+"\"")
	}
	if v := rscData.WebFilteringProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"web-filtering http-profile \""+v+"\"")
	}

	if rscData.AntiVirus != nil {
		if rscData.AntiVirus.isEmpty() {
			return path.Root("anti_virus").AtName("*"),
				errors.New("anti_virus block is empty")
		}

		configSet = append(configSet, rscData.AntiVirus.configSet(setPrefix+"anti-virus ")...)
	}
	if rscData.ContentFiltering != nil {
		if rscData.ContentFiltering.isEmpty() {
			return path.Root("content_filtering").AtName("*"),
				errors.New("content_filtering block is empty")
		}

		configSet = append(configSet, rscData.ContentFiltering.configSet(setPrefix+"content-filtering ")...)
	}
	if rscData.TrafficSessionsPerClient != nil {
		if rscData.TrafficSessionsPerClient.isEmpty() {
			return path.Root("traffic_sessions_per_client").AtName("*"),
				errors.New("traffic_sessions_per_client block is empty")
		}

		configSet = append(configSet, rscData.TrafficSessionsPerClient.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityUtmPolicyBlockProtocolProfile) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)

	if v := block.FTPDownloadProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ftp download-profile \""+v+"\"")
	}
	if v := block.FTPUploadProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ftp upload-profile \""+v+"\"")
	}
	if v := block.HTTPProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"http-profile \""+v+"\"")
	}
	if v := block.IMAPProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"imap-profile \""+v+"\"")
	}
	if v := block.POP3Profile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"pop3-profile \""+v+"\"")
	}
	if v := block.SMTPProfile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"smtp-profile \""+v+"\"")
	}

	return configSet
}

func (block *securityUtmPolicyBlockTrafficSessionsPerClient) configSet(setPrefix string) []string {
	configSet := make([]string, 0, 100)
	setPrefix += "traffic-options sessions-per-client "

	if !block.Limit.IsNull() {
		configSet = append(configSet, setPrefix+"limit "+
			utils.ConvI64toa(block.Limit.ValueInt64()))
	}
	if v := block.OverLimit.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"over-limit "+v)
	}

	return configSet
}

func (rscData *securityUtmPolicyData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm utm-policy \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "anti-spam smtp-profile "):
				rscData.AntiSpamSMTPProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "anti-virus "):
				if rscData.AntiVirus == nil {
					rscData.AntiVirus = &securityUtmPolicyBlockProtocolProfile{}
				}

				rscData.AntiVirus.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "content-filtering "):
				if rscData.ContentFiltering == nil {
					rscData.ContentFiltering = &securityUtmPolicyBlockProtocolProfile{}
				}

				rscData.ContentFiltering.read(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "traffic-options sessions-per-client "):
				if rscData.TrafficSessionsPerClient == nil {
					rscData.TrafficSessionsPerClient = &securityUtmPolicyBlockTrafficSessionsPerClient{}
				}

				if err := rscData.TrafficSessionsPerClient.read(itemTrim); err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "web-filtering http-profile "):
				rscData.WebFilteringProfile = types.StringValue(strings.Trim(itemTrim, "\""))
			}
		}
	}

	return nil
}

func (block *securityUtmPolicyBlockProtocolProfile) read(itemTrim string) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "ftp download-profile "):
		block.FTPDownloadProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "ftp upload-profile "):
		block.FTPUploadProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "http-profile "):
		block.HTTPProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "imap-profile "):
		block.IMAPProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "pop3-profile "):
		block.POP3Profile = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "smtp-profile "):
		block.SMTPProfile = types.StringValue(strings.Trim(itemTrim, "\""))
	}
}

func (block *securityUtmPolicyBlockTrafficSessionsPerClient) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "limit "):
		block.Limit, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "over-limit "):
		block.OverLimit = types.StringValue(itemTrim)
	}

	return err
}

func (rscData *securityUtmPolicyData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security utm utm-policy \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
