package provider

import (
	"context"
	"errors"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	_ resource.Resource                   = &servicesSSLInitiationProfile{}
	_ resource.ResourceWithConfigure      = &servicesSSLInitiationProfile{}
	_ resource.ResourceWithValidateConfig = &servicesSSLInitiationProfile{}
	_ resource.ResourceWithImportState    = &servicesSSLInitiationProfile{}
	_ resource.ResourceWithUpgradeState   = &servicesSSLInitiationProfile{}
)

type servicesSSLInitiationProfile struct {
	client *junos.Client
}

func newServicesSSLInitiationProfileResource() resource.Resource {
	return &servicesSSLInitiationProfile{}
}

func (rsc *servicesSSLInitiationProfile) typeName() string {
	return providerName + "_services_ssl_initiation_profile"
}

func (rsc *servicesSSLInitiationProfile) junosName() string {
	return "services ssl initiation profile"
}

func (rsc *servicesSSLInitiationProfile) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *servicesSSLInitiationProfile) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *servicesSSLInitiationProfile) Configure(
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

func (rsc *servicesSSLInitiationProfile) Schema(
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
				Description: "Profile name (Profile identifier).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"client_certificate": schema.StringAttribute{
				Optional:    true,
				Description: "Local certificate identifier.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"custom_ciphers": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Custom cipher list.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.NoNullValues(),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
						tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					),
				},
			},
			"enable_flow_tracing": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable flow tracing for the profile.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"enable_session_cache": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable SSL session cache.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"preferred_ciphers": schema.StringAttribute{
				Optional:    true,
				Description: "Select preferred ciphers.",
				Validators: []validator.String{
					stringvalidator.OneOf("custom", "medium", "strong", "weak"),
				},
			},
			"protocol_version": schema.StringAttribute{
				Optional:    true,
				Description: "Protocol SSL version accepted.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"trusted_ca": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of trusted certificate authority profiles.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.NoNullValues(),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 32),
						tfvalidator.StringDoubleQuoteExclusion(),
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"actions": schema.SingleNestedBlock{
				Description: "Traffic related actions.",
				Attributes: map[string]schema.Attribute{
					"crl_disable": schema.BoolAttribute{
						Optional:    true,
						Description: "Disable CRL validation.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"crl_if_not_present": schema.StringAttribute{
						Optional:    true,
						Description: "Action if CRL information is not present.",
						Validators: []validator.String{
							stringvalidator.OneOf("allow", "drop"),
						},
					},
					"crl_ignore_hold_instruction_code": schema.BoolAttribute{
						Optional:    true,
						Description: "Ignore 'Hold Instruction Code' present in the CRL entry.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ignore_server_auth_failure": schema.BoolAttribute{
						Optional:    true,
						Description: "Ignore server authentication failure.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
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

type servicesSSLInitiationProfileData struct {
	ID                 types.String                              `tfsdk:"id"`
	Name               types.String                              `tfsdk:"name"`
	ClientCertificate  types.String                              `tfsdk:"client_certificate"`
	CustomCiphers      []types.String                            `tfsdk:"custom_ciphers"`
	EnableFlowTracing  types.Bool                                `tfsdk:"enable_flow_tracing"`
	EnableSessionCache types.Bool                                `tfsdk:"enable_session_cache"`
	PreferredCiphers   types.String                              `tfsdk:"preferred_ciphers"`
	ProtocolVersion    types.String                              `tfsdk:"protocol_version"`
	TrustedCA          []types.String                            `tfsdk:"trusted_ca"`
	Actions            *servicesSSLInitiationProfileBlockActions `tfsdk:"actions"`
}

type servicesSSLInitiationProfileConfig struct {
	ID                 types.String                              `tfsdk:"id"`
	Name               types.String                              `tfsdk:"name"`
	ClientCertificate  types.String                              `tfsdk:"client_certificate"`
	CustomCiphers      types.Set                                 `tfsdk:"custom_ciphers"`
	EnableFlowTracing  types.Bool                                `tfsdk:"enable_flow_tracing"`
	EnableSessionCache types.Bool                                `tfsdk:"enable_session_cache"`
	PreferredCiphers   types.String                              `tfsdk:"preferred_ciphers"`
	ProtocolVersion    types.String                              `tfsdk:"protocol_version"`
	TrustedCA          types.Set                                 `tfsdk:"trusted_ca"`
	Actions            *servicesSSLInitiationProfileBlockActions `tfsdk:"actions"`
}

type servicesSSLInitiationProfileBlockActions struct {
	CrlDisable                   types.Bool   `tfsdk:"crl_disable"`
	CrlIfNotPresent              types.String `tfsdk:"crl_if_not_present"`
	CrlIgnoreHoldInstructionCode types.Bool   `tfsdk:"crl_ignore_hold_instruction_code"`
	IgnoreServerAuthFailure      types.Bool   `tfsdk:"ignore_server_auth_failure"`
}

func (block *servicesSSLInitiationProfileBlockActions) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

func (rsc *servicesSSLInitiationProfile) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config servicesSSLInitiationProfileConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Actions != nil {
		if config.Actions.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("actions").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"actions block is empty",
			)
		}
	}
}

func (rsc *servicesSSLInitiationProfile) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan servicesSSLInitiationProfileData
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
			profileExists, err := checkServicesSSLInitiationProfileExists(fnCtx, plan.Name.ValueString(), junSess)
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
			profileExists, err := checkServicesSSLInitiationProfileExists(fnCtx, plan.Name.ValueString(), junSess)
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

func (rsc *servicesSSLInitiationProfile) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data servicesSSLInitiationProfileData
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

func (rsc *servicesSSLInitiationProfile) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state servicesSSLInitiationProfileData
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

func (rsc *servicesSSLInitiationProfile) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state servicesSSLInitiationProfileData
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

func (rsc *servicesSSLInitiationProfile) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data servicesSSLInitiationProfileData

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

func checkServicesSSLInitiationProfileExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services ssl initiation profile \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *servicesSSLInitiationProfileData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *servicesSSLInitiationProfileData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *servicesSSLInitiationProfileData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set services ssl initiation profile \"" + rscData.Name.ValueString() + "\" "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if v := rscData.ClientCertificate.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"client-certificate \""+v+"\"")
	}
	for _, v := range rscData.CustomCiphers {
		configSet = append(configSet, setPrefix+"custom-ciphers "+v.ValueString())
	}
	if rscData.EnableFlowTracing.ValueBool() {
		configSet = append(configSet, setPrefix+"enable-flow-tracing")
	}
	if rscData.EnableSessionCache.ValueBool() {
		configSet = append(configSet, setPrefix+"enable-session-cache")
	}
	if v := rscData.PreferredCiphers.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"preferred-ciphers "+v)
	}
	if v := rscData.ProtocolVersion.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"protocol-version "+v)
	}
	for _, v := range rscData.TrustedCA {
		configSet = append(configSet, setPrefix+"trusted-ca \""+v.ValueString()+"\"")
	}

	if rscData.Actions != nil {
		if rscData.Actions.isEmpty() {
			return path.Root("actions").AtName("*"),
				errors.New("actions block is empty")
		}

		if rscData.Actions.CrlDisable.ValueBool() {
			configSet = append(configSet, setPrefix+"actions crl disable")
		}
		if v := rscData.Actions.CrlIfNotPresent.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"actions crl if-not-present "+v)
		}
		if rscData.Actions.CrlIgnoreHoldInstructionCode.ValueBool() {
			configSet = append(configSet, setPrefix+"actions crl ignore-hold-instruction-code")
		}
		if rscData.Actions.IgnoreServerAuthFailure.ValueBool() {
			configSet = append(configSet, setPrefix+"actions ignore-server-auth-failure")
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *servicesSSLInitiationProfileData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services ssl initiation profile \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.fillID()
		for item := range strings.SplitSeq(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "actions "):
				if rscData.Actions == nil {
					rscData.Actions = &servicesSSLInitiationProfileBlockActions{}
				}

				switch {
				case itemTrim == "crl disable":
					rscData.Actions.CrlDisable = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "crl if-not-present "):
					rscData.Actions.CrlIfNotPresent = types.StringValue(itemTrim)
				case itemTrim == "crl ignore-hold-instruction-code":
					rscData.Actions.CrlIgnoreHoldInstructionCode = types.BoolValue(true)
				case itemTrim == "ignore-server-auth-failure":
					rscData.Actions.IgnoreServerAuthFailure = types.BoolValue(true)
				}
			case balt.CutPrefixInString(&itemTrim, "client-certificate "):
				rscData.ClientCertificate = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "custom-ciphers "):
				rscData.CustomCiphers = append(rscData.CustomCiphers, types.StringValue(itemTrim))
			case itemTrim == "enable-flow-tracing":
				rscData.EnableFlowTracing = types.BoolValue(true)
			case itemTrim == "enable-session-cache":
				rscData.EnableSessionCache = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "preferred-ciphers "):
				rscData.PreferredCiphers = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "protocol-version "):
				rscData.ProtocolVersion = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "trusted-ca "):
				rscData.TrustedCA = append(rscData.TrustedCA, types.StringValue(strings.Trim(itemTrim, "\"")))
			}
		}
	}

	return nil
}

func (rscData *servicesSSLInitiationProfileData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete services ssl initiation profile \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
