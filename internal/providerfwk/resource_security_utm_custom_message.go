package providerfwk

import (
	"context"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
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
	_ resource.Resource                   = &securityUtmCustomMessage{}
	_ resource.ResourceWithConfigure      = &securityUtmCustomMessage{}
	_ resource.ResourceWithValidateConfig = &securityUtmCustomMessage{}
	_ resource.ResourceWithImportState    = &securityUtmCustomMessage{}
)

type securityUtmCustomMessage struct {
	client *junos.Client
}

func newSecurityUtmCustomMessageResource() resource.Resource {
	return &securityUtmCustomMessage{}
}

func (rsc *securityUtmCustomMessage) typeName() string {
	return providerName + "_security_utm_custom_message"
}

func (rsc *securityUtmCustomMessage) junosName() string {
	return "security utm custom-objects custom-message"
}

func (rsc *securityUtmCustomMessage) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityUtmCustomMessage) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityUtmCustomMessage) Configure(
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

func (rsc *securityUtmCustomMessage) Schema(
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
				Description: "The name of security utm custom-object custom-message.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 59),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Type of custom message.",
				Validators: []validator.String{
					stringvalidator.OneOf("custom-page", "redirect-url", "user-message"),
				},
			},
			"content": schema.StringAttribute{
				Optional:    true,
				Description: "Content of custom message.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1024),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"custom_page_file": schema.StringAttribute{
				Optional:    true,
				Description: "Name of custom page file.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringDoubleQuoteExclusion(),
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^.+\.html$`),
						"must be end with '.html'"),
				},
			},
		},
	}
}

type securityUtmCustomMessageData struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	Content        types.String `tfsdk:"content"`
	CustomPageFile types.String `tfsdk:"custom_page_file"`
}

func (rsc *securityUtmCustomMessage) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityUtmCustomMessageData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Content.IsNull() &&
		!config.Content.IsUnknown() &&
		!config.CustomPageFile.IsNull() &&
		!config.CustomPageFile.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("content"),
			tfdiag.ConflictConfigErrSummary,
			"content and custom_page_file cannot be configured together",
		)
	}

	switch {
	case config.Type.IsNull():
	case config.Type.IsUnknown():
	case config.Type.ValueString() == "custom-page":
		if !config.Content.IsNull() &&
			!config.Content.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("content"),
				tfdiag.ConflictConfigErrSummary,
				"content cannot be configured when type = "+config.Type.ValueString(),
			)
		}
		if config.CustomPageFile.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("type"),
				tfdiag.MissingConfigErrSummary,
				"custom_page_file must be specified when type = "+config.Type.ValueString(),
			)
		}
	case config.Type.ValueString() == "redirect-url":
		if !config.CustomPageFile.IsNull() &&
			!config.CustomPageFile.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("custom_page_file"),
				tfdiag.ConflictConfigErrSummary,
				"custom_page_file cannot be configured when type = "+config.Type.ValueString(),
			)
		}
		if config.Content.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("type"),
				tfdiag.MissingConfigErrSummary,
				"content must be specified when type = "+config.Type.ValueString(),
			)
		} else if !config.Content.IsUnknown() {
			if !strings.HasPrefix(config.Content.ValueString(), "http://") &&
				!strings.HasPrefix(config.Content.ValueString(), "https://") {
				resp.Diagnostics.AddAttributeError(
					path.Root("type"),
					tfdiag.MissingConfigErrSummary,
					"content must be specified with 'http(s)://' prefix when type = "+config.Type.ValueString(),
				)
			}
		}
	case config.Type.ValueString() == "user-message":
		if !config.CustomPageFile.IsNull() &&
			!config.CustomPageFile.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("custom_page_file"),
				tfdiag.ConflictConfigErrSummary,
				"custom_page_file cannot be configured when type = "+config.Type.ValueString(),
			)
		}
		if config.Content.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("type"),
				tfdiag.MissingConfigErrSummary,
				"content must be specified when type = "+config.Type.ValueString(),
			)
		}
	}
}

func (rsc *securityUtmCustomMessage) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityUtmCustomMessageData
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
			messageExists, err := checkSecurityUtmCustomMessageExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if messageExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			messageExists, err := checkSecurityUtmCustomMessageExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !messageExists {
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

func (rsc *securityUtmCustomMessage) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityUtmCustomMessageData
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

func (rsc *securityUtmCustomMessage) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityUtmCustomMessageData
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

func (rsc *securityUtmCustomMessage) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityUtmCustomMessageData
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

func (rsc *securityUtmCustomMessage) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityUtmCustomMessageData

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

func checkSecurityUtmCustomMessageExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm custom-objects custom-message " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityUtmCustomMessageData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityUtmCustomMessageData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *securityUtmCustomMessageData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set security utm custom-objects custom-message " + rscData.Name.ValueString() + " "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix + "type " + rscData.Type.ValueString()

	if v := rscData.Content.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"content \""+v+"\"")
	}
	if v := rscData.CustomPageFile.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"custom-page-file \""+v+"\"")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *securityUtmCustomMessageData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security utm custom-objects custom-message " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "content "):
				rscData.Content = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "custom-page-file "):
				rscData.CustomPageFile = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "type "):
				rscData.Type = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (rscData *securityUtmCustomMessageData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security utm custom-objects custom-message " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
