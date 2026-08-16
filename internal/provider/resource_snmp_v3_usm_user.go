package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &snmpV3UsmUser{}
	_ resource.ResourceWithConfigure      = &snmpV3UsmUser{}
	_ resource.ResourceWithValidateConfig = &snmpV3UsmUser{}
	_ resource.ResourceWithImportState    = &snmpV3UsmUser{}
)

type snmpV3UsmUser struct {
	client *junos.Client
}

func newSnmpV3UsmUserResource() resource.Resource {
	return &snmpV3UsmUser{}
}

func (rsc *snmpV3UsmUser) typeName() string {
	return providerName + "_snmp_v3_usm_user"
}

func (rsc *snmpV3UsmUser) junosName() string {
	return "snmp v3 usm user"
}

func (rsc *snmpV3UsmUser) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *snmpV3UsmUser) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *snmpV3UsmUser) Configure(
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

func (rsc *snmpV3UsmUser) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format" +
					" `local_<name>` or `remote_<engine_id>_<name>` (according to <engine_type>).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of snmp v3 USM user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"engine_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("local"),
				Description: "Local or remote engine user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("local", "remote"),
				},
			},
			"engine_id": schema.StringAttribute{
				Optional:    true,
				Description: "Remote engine id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(5, 64),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"authentication_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Encrypted key used for user authentication.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"authentication_key_wo": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
				Description: "Encrypted key used for user authentication, not stored in state.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
					stringvalidator.AlsoRequires(path.MatchRoot("authentication_key_wo_version")),
				},
			},
			"authentication_key_wo_version": schema.Int64Attribute{
				Optional:    true,
				Description: "Version of `authentication_key_wo` to trigger the sending of its value.",
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("authentication_key_wo")),
				},
			},
			"authentication_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "User's authentication password.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 1024),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"authentication_password_wo": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
				Description: "User's authentication password, not stored in state.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 1024),
					tfvalidator.StringDoubleQuoteExclusion(),
					stringvalidator.AlsoRequires(path.MatchRoot("authentication_password_wo_version")),
				},
			},
			"authentication_password_wo_version": schema.Int64Attribute{
				Optional:    true,
				Description: "Version of `authentication_password_wo` to trigger the sending of its value.",
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("authentication_password_wo")),
				},
			},
			"authentication_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("authentication-none"),
				Description: "Define authentication type.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"authentication-md5",
						"authentication-none",
						"authentication-sha",
						"authentication-sha224",
						"authentication-sha256",
						"authentication-sha384",
						"authentication-sha512",
					),
				},
			},
			"privacy_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Encrypted key used for user privacy.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"privacy_key_wo": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
				Description: "Encrypted key used for user privacy, not stored in state.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
					stringvalidator.AlsoRequires(path.MatchRoot("privacy_key_wo_version")),
				},
			},
			"privacy_key_wo_version": schema.Int64Attribute{
				Optional:    true,
				Description: "Version of `privacy_key_wo` to trigger the sending of its value.",
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("privacy_key_wo")),
				},
			},
			"privacy_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "User's privacy password.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 1024),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"privacy_password_wo": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				WriteOnly:   true,
				Description: "User's privacy password, not stored in state.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 1024),
					tfvalidator.StringDoubleQuoteExclusion(),
					stringvalidator.AlsoRequires(path.MatchRoot("privacy_password_wo_version")),
				},
			},
			"privacy_password_wo_version": schema.Int64Attribute{
				Optional:    true,
				Description: "Version of `privacy_password_wo` to trigger the sending of its value.",
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("privacy_password_wo")),
				},
			},
			"privacy_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("privacy-none"),
				Description: "Define privacy type.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"privacy-3des",
						"privacy-aes128",
						"privacy-des",
						"privacy-none",
					),
				},
			},
		},
	}
}

type snmpV3UsmUserData struct {
	ID                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	EngineType                      types.String `tfsdk:"engine_type"`
	EngineID                        types.String `tfsdk:"engine_id"`
	AuthenticationKey               types.String `tfsdk:"authentication_key"`
	AuthenticationKeyWO             types.String `tfsdk:"authentication_key_wo"`
	AuthenticationKeyWOVersion      types.Int64  `tfsdk:"authentication_key_wo_version"`
	AuthenticationPassword          types.String `tfsdk:"authentication_password"`
	AuthenticationPasswordWO        types.String `tfsdk:"authentication_password_wo"`
	AuthenticationPasswordWOVersion types.Int64  `tfsdk:"authentication_password_wo_version"`
	AuthenticationType              types.String `tfsdk:"authentication_type"`
	PrivacyKey                      types.String `tfsdk:"privacy_key"`
	PrivacyKeyWO                    types.String `tfsdk:"privacy_key_wo"`
	PrivacyKeyWOVersion             types.Int64  `tfsdk:"privacy_key_wo_version"`
	PrivacyPassword                 types.String `tfsdk:"privacy_password"`
	PrivacyPasswordWO               types.String `tfsdk:"privacy_password_wo"`
	PrivacyPasswordWOVersion        types.Int64  `tfsdk:"privacy_password_wo_version"`
	PrivacyType                     types.String `tfsdk:"privacy_type"`
}

type snmpV3UsmUserPrivateState struct {
	AuthenticationKey string `json:"authentication_key"`
	PrivacyKey        string `json:"privacy_key"`
}

func (ste *snmpV3UsmUserPrivateState) key() string {
	return "v0"
}

func (ste *snmpV3UsmUserPrivateState) get(
	ctx context.Context, private privateStateGetter,
) (diags diag.Diagnostics) {
	data, getDiags := private.GetKey(ctx, ste.key())
	diags.Append(getDiags...)
	if diags.HasError() {
		return diags
	}

	if data != nil {
		if err := json.Unmarshal(data, ste); err != nil {
			diags.AddError(tfdiag.GetPrivateStateErrSummary, fmt.Sprintf("json unmarshal: %s", err))
		}
	}

	return diags
}

func (rsc *snmpV3UsmUser) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config snmpV3UsmUserData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	engineType := config.EngineType.ValueString()
	switch {
	case config.EngineType.IsUnknown():
	case engineType == "local" || config.EngineType.IsNull():
		if !config.EngineID.IsNull() && !config.EngineID.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("engine_id"),
				tfdiag.ConflictConfigErrSummary,
				"could not create "+rsc.junosName()+" with engine_type = local and engine_id specified",
			)
		}
	case engineType == "remote":
		if config.EngineID.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("engine_id"),
				tfdiag.MissingConfigErrSummary,
				"could not create "+rsc.junosName()+" with engine_type = remote and empty engine_id",
			)
		}
	}

	if !config.AuthenticationKey.IsNull() && !config.AuthenticationKey.IsUnknown() {
		if !config.AuthenticationPassword.IsNull() && !config.AuthenticationPassword.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_key"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_key and authentication_password cannot be configured together",
			)
		}
		if !config.AuthenticationKeyWO.IsNull() && !config.AuthenticationKeyWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_key"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_key and authentication_key_wo cannot be configured together",
			)
		}
		if !config.AuthenticationPasswordWO.IsNull() && !config.AuthenticationPasswordWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_key"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_key and authentication_password_wo cannot be configured together",
			)
		}
	}
	if !config.AuthenticationPassword.IsNull() && !config.AuthenticationPassword.IsUnknown() {
		if !config.AuthenticationKeyWO.IsNull() && !config.AuthenticationKeyWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_password"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_password and authentication_key_wo cannot be configured together",
			)
		}
		if !config.AuthenticationPasswordWO.IsNull() && !config.AuthenticationPasswordWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_password"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_password and authentication_password_wo cannot be configured together",
			)
		}
	}
	if !config.AuthenticationKeyWO.IsNull() && !config.AuthenticationKeyWO.IsUnknown() &&
		!config.AuthenticationPasswordWO.IsNull() && !config.AuthenticationPasswordWO.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("authentication_key_wo"),
			tfdiag.ConflictConfigErrSummary,
			"authentication_key_wo and authentication_password_wo cannot be configured together",
		)
	}
	if !config.PrivacyKey.IsNull() && !config.PrivacyKey.IsUnknown() {
		if !config.PrivacyPassword.IsNull() && !config.PrivacyPassword.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_key"),
				tfdiag.ConflictConfigErrSummary,
				"privacy_key and privacy_password cannot be configured together",
			)
		}
		if !config.PrivacyKeyWO.IsNull() && !config.PrivacyKeyWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_key"),
				tfdiag.ConflictConfigErrSummary,
				"privacy_key and privacy_key_wo cannot be configured together",
			)
		}
		if !config.PrivacyPasswordWO.IsNull() && !config.PrivacyPasswordWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_key"),
				tfdiag.ConflictConfigErrSummary,
				"privacy_key and privacy_password_wo cannot be configured together",
			)
		}
	}
	if !config.PrivacyPassword.IsNull() && !config.PrivacyPassword.IsUnknown() {
		if !config.PrivacyKeyWO.IsNull() && !config.PrivacyKeyWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_password"),
				tfdiag.ConflictConfigErrSummary,
				"privacy_password and privacy_key_wo cannot be configured together",
			)
		}
		if !config.PrivacyPasswordWO.IsNull() && !config.PrivacyPasswordWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_password"),
				tfdiag.ConflictConfigErrSummary,
				"privacy_password and privacy_password_wo cannot be configured together",
			)
		}
	}
	if !config.PrivacyKeyWO.IsNull() && !config.PrivacyKeyWO.IsUnknown() &&
		!config.PrivacyPasswordWO.IsNull() && !config.PrivacyPasswordWO.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("privacy_key_wo"),
			tfdiag.ConflictConfigErrSummary,
			"privacy_key_wo and privacy_password_wo cannot be configured together",
		)
	}

	if !config.AuthenticationType.IsNull() && !config.AuthenticationType.IsUnknown() &&
		config.AuthenticationType.ValueString() != "authentication-none" {
		if config.AuthenticationKey.IsNull() && config.AuthenticationPassword.IsNull() &&
			config.AuthenticationKeyWO.IsNull() && config.AuthenticationPasswordWO.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_type"),
				tfdiag.MissingConfigErrSummary,
				"one of authentication_key, authentication_password, authentication_key_wo or authentication_password_wo"+
					" must be specified when authentication_type != authentication-none",
			)
		}
	} else if config.AuthenticationType.IsNull() || config.AuthenticationType.ValueString() == "authentication-none" {
		if !config.PrivacyType.IsNull() && !config.PrivacyType.IsUnknown() &&
			config.PrivacyType.ValueString() != "privacy-none" {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_type"),
				tfdiag.MissingConfigErrSummary,
				"authentication should be configured before configuring the privacy",
			)
		}
		if !config.AuthenticationKey.IsNull() && !config.AuthenticationKey.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_key"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_key not compatible when authentication_type = authentication-none",
			)
		}
		if !config.AuthenticationKeyWO.IsNull() && !config.AuthenticationKeyWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_key_wo"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_key_wo not compatible when authentication_type = authentication-none",
			)
		}
		if !config.AuthenticationPassword.IsNull() && !config.AuthenticationPassword.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_password"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_password not compatible when authentication_type = authentication-none",
			)
		}
		if !config.AuthenticationPasswordWO.IsNull() && !config.AuthenticationPasswordWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_password_wo"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_password_wo not compatible when authentication_type = authentication-none",
			)
		}
	}
	if !config.PrivacyType.IsNull() && !config.PrivacyType.IsUnknown() &&
		config.PrivacyType.ValueString() != "privacy-none" {
		if config.PrivacyKey.IsNull() && config.PrivacyPassword.IsNull() &&
			config.PrivacyKeyWO.IsNull() && config.PrivacyPasswordWO.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_type"),
				tfdiag.MissingConfigErrSummary,
				"one of privacy_key, privacy_password, privacy_key_wo or privacy_password_wo"+
					" must be specified when privacy_type != privacy-none",
			)
		}
	} else if config.PrivacyType.IsNull() || config.PrivacyType.ValueString() == "privacy-none" {
		if !config.PrivacyKey.IsNull() && !config.PrivacyKey.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_key"),
				tfdiag.ConflictConfigErrSummary,
				"privacy_key not compatible when privacy_type = privacy-none",
			)
		}
		if !config.PrivacyKeyWO.IsNull() && !config.PrivacyKeyWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_key_wo"),
				tfdiag.ConflictConfigErrSummary,
				"privacy_key_wo not compatible when privacy_type = privacy-none",
			)
		}
		if !config.PrivacyPassword.IsNull() && !config.PrivacyPassword.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_password"),
				tfdiag.ConflictConfigErrSummary,
				"privacy_password not compatible when privacy_type = privacy-none",
			)
		}
		if !config.PrivacyPasswordWO.IsNull() && !config.PrivacyPasswordWO.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_password_wo"),
				tfdiag.ConflictConfigErrSummary,
				"privacy_password_wo not compatible when privacy_type = privacy-none",
			)
		}
	}
}

func (rsc *snmpV3UsmUser) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan snmpV3UsmUserData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(plan.getWriteOnly(ctx, req.Config)...)
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
	if plan.EngineType.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("engine_type"),
			"Empty Engine Type",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "engine_type"),
		)

		return
	}
	switch v := plan.EngineType.ValueString(); v {
	case "local":
		if plan.EngineID.ValueString() != "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("engine_id"),
				tfdiag.ConflictConfigErrSummary,
				"could not create "+rsc.junosName()+" with engine_type = local and engine_id specified",
			)

			return
		}
	case "remote":
		if plan.EngineID.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("engine_id"),
				tfdiag.MissingConfigErrSummary,
				"could not create "+rsc.junosName()+" with engine_type = remote and empty engine_id",
			)

			return
		}
	default:
		resp.Diagnostics.AddAttributeError(
			path.Root("engine_type"),
			"Bad Engine Type",
			fmt.Sprintf("could not create "+rsc.junosName()+" with engine_type %q", v),
		)

		return
	}

	var _ resourceDataReadPrivateToState = &plan
	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			userExists, err := checkSnmpV3UsmUserExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.EngineType.ValueString(),
				plan.EngineID.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if userExists {
				if plan.EngineType.ValueString() == "remote" {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf(
							rsc.junosName()+" %q in remote-engine %q already exists",
							plan.Name.ValueString(), plan.EngineID.ValueString(),
						),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf(
							rsc.junosName()+" %q in local-engine already exists",
							plan.Name.ValueString(),
						),
					)
				}

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			userExists, err := checkSnmpV3UsmUserExists(
				fnCtx,
				plan.Name.ValueString(),
				plan.EngineType.ValueString(),
				plan.EngineID.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !userExists {
				if plan.EngineType.ValueString() == "remote" {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						fmt.Sprintf(
							rsc.junosName()+" %q in remote-engine %q does not exists after commit "+
								"=> check your config",
							plan.Name.ValueString(), plan.EngineID.ValueString(),
						),
					)
				} else {
					resp.Diagnostics.AddError(
						tfdiag.NotFoundErrSummary,
						fmt.Sprintf(
							rsc.junosName()+" %q in local-engine does not exists after commit "+
								"=> check your config",
							plan.Name.ValueString(),
						),
					)
				}

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *snmpV3UsmUser) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data snmpV3UsmUserData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom3String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Name.ValueString(),
			state.EngineType.ValueString(),
			state.EngineID.ValueString(),
		},
		&data,
		func() {
			data.keepWriteOnly(&state)

			var privateState snmpV3UsmUserPrivateState
			resp.Diagnostics.Append(privateState.get(ctx, req.Private)...)
			if resp.Diagnostics.HasError() {
				return
			}

			if data.AuthenticationType.ValueString() != "authentication-none" &&
				data.AuthenticationType.ValueString() == state.AuthenticationType.ValueString() &&
				data.AuthenticationKey.ValueString() != "" &&
				state.AuthenticationPassword.ValueString() != "" {
				if privateState.AuthenticationKey != "" {
					if privateState.AuthenticationKey == data.AuthenticationKey.ValueString() {
						data.AuthenticationPassword = state.AuthenticationPassword
					} else {
						data.AuthenticationPassword = types.StringValue(`?`)
					}
					data.AuthenticationKey = types.StringNull()
				} else {
					data.AuthenticationPassword = state.AuthenticationPassword
					data.AuthenticationKey = types.StringNull()
				}
			}
			if data.PrivacyType.ValueString() != "privacy-none" &&
				data.PrivacyType.ValueString() == state.PrivacyType.ValueString() &&
				data.PrivacyKey.ValueString() != "" &&
				state.PrivacyPassword.ValueString() != "" {
				if privateState.PrivacyKey != "" {
					if privateState.PrivacyKey == data.PrivacyKey.ValueString() {
						data.PrivacyPassword = state.PrivacyPassword
					} else {
						data.PrivacyPassword = types.StringValue(`?`)
					}
					data.PrivacyKey = types.StringNull()
				} else {
					data.PrivacyPassword = state.PrivacyPassword
					data.PrivacyKey = types.StringNull()
				}
			}
		},
		resp,
	)
}

func (rsc *snmpV3UsmUser) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state snmpV3UsmUserData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(plan.getWriteOnly(ctx, req.Config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadPrivateToState = &plan
	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *snmpV3UsmUser) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state snmpV3UsmUserData
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

func (rsc *snmpV3UsmUser) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	idList := strings.Split(req.ID, junos.IDSeparator)
	var name, engineType, engineID string
	switch {
	case len(idList) == 2 && idList[0] == "local":
		engineType = idList[0]
		name = idList[1]
	case len(idList) == 3 && idList[0] == "remote":
		engineType = idList[0]
		engineID = idList[1]
		name = idList[2]
	default:
		resp.Diagnostics.AddError(
			"Bad ID Format",
			fmt.Sprintf(
				"can't find snmp v3 usm user with id '%v' (id must be "+
					"local"+junos.IDSeparator+"<name> or "+
					"remote"+junos.IDSeparator+"<engine_id>"+junos.IDSeparator+"<name>)",
				req.ID,
			))

		return
	}
	var data snmpV3UsmUserData
	if err := data.read(ctx, name, engineType, engineID, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}
	if data.nullID() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			defaultResourceImportDontFindMessage(rsc, req.ID)+
				" (id must be local"+junos.IDSeparator+"<name> or "+
				"remote"+junos.IDSeparator+"<engine_id>"+junos.IDSeparator+"<name>)",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func checkSnmpV3UsmUserExists(
	ctx context.Context, name, engineType, engineID string, junSess *junos.Session,
) (
	bool, error,
) {
	showPrefix := junos.CmdShowConfig + "snmp v3 usm "
	switch engineType {
	case "local":
		showPrefix += "local-engine "
	case "remote":
		showPrefix += "remote-engine \"" + engineID + "\" "
	default:
		return false, fmt.Errorf("can't check config with engine_type %q", engineType)
	}
	showConfig, err := junSess.Command(ctx, showPrefix+
		"user \""+name+"\""+junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

// getWriteOnly read the write-only arguments from the configuration,
// their values aren't present in the plan or the state.
func (rscData *snmpV3UsmUserData) getWriteOnly(
	ctx context.Context, config tfsdk.Config,
) (diags diag.Diagnostics) {
	diags.Append(config.GetAttribute(ctx,
		path.Root("authentication_key_wo"), &rscData.AuthenticationKeyWO)...)
	diags.Append(config.GetAttribute(ctx,
		path.Root("authentication_password_wo"), &rscData.AuthenticationPasswordWO)...)
	diags.Append(config.GetAttribute(ctx,
		path.Root("privacy_key_wo"), &rscData.PrivacyKeyWO)...)
	diags.Append(config.GetAttribute(ctx,
		path.Root("privacy_password_wo"), &rscData.PrivacyPasswordWO)...)

	return diags
}

// keepWriteOnly carry over the version arguments of the write-only arguments from the state,
// and don't read the secrets in the standard arguments when the write-only ones are used.
func (rscData *snmpV3UsmUserData) keepWriteOnly(state *snmpV3UsmUserData) {
	rscData.AuthenticationKeyWOVersion = state.AuthenticationKeyWOVersion
	rscData.AuthenticationPasswordWOVersion = state.AuthenticationPasswordWOVersion
	if !state.AuthenticationKeyWOVersion.IsNull() ||
		!state.AuthenticationPasswordWOVersion.IsNull() {
		rscData.AuthenticationKey = types.StringNull()
		rscData.AuthenticationPassword = types.StringNull()
	}
	rscData.PrivacyKeyWOVersion = state.PrivacyKeyWOVersion
	rscData.PrivacyPasswordWOVersion = state.PrivacyPasswordWOVersion
	if !state.PrivacyKeyWOVersion.IsNull() ||
		!state.PrivacyPasswordWOVersion.IsNull() {
		rscData.PrivacyKey = types.StringNull()
		rscData.PrivacyPassword = types.StringNull()
	}
}

func (rscData *snmpV3UsmUserData) fillID() {
	switch v := rscData.EngineType.ValueString(); v {
	case "local":
		rscData.ID = types.StringValue(
			v + junos.IDSeparator + rscData.Name.ValueString(),
		)
	case "remote":
		rscData.ID = types.StringValue(
			v + junos.IDSeparator + rscData.EngineID.ValueString() + junos.IDSeparator + rscData.Name.ValueString(),
		)
	}
}

func (rscData *snmpV3UsmUserData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *snmpV3UsmUserData) set(
	ctx context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)
	setPrefix := "set snmp v3 usm "
	switch v := rscData.EngineType.ValueString(); v {
	case "local":
		setPrefix += "local-engine "
	case "remote":
		setPrefix += "remote-engine \"" + rscData.EngineID.ValueString() + "\" "
	default:
		return path.Root("engine_type"), fmt.Errorf("can't set config with engine_type %q", v)
	}
	setPrefix += "user \"" + rscData.Name.ValueString() + "\" "

	if authenticationType := rscData.AuthenticationType.ValueString(); authenticationType != "authentication-none" {
		if rscData.AuthenticationKey.ValueString() == "" && rscData.AuthenticationPassword.ValueString() == "" &&
			rscData.AuthenticationKeyWO.ValueString() == "" && rscData.AuthenticationPasswordWO.ValueString() == "" {
			return path.Root("authentication_type"),
				errors.New("one of authentication_key, authentication_password, authentication_key_wo or " +
					"authentication_password_wo must be specified when authentication_type != authentication-none")
		}
		if v := rscData.AuthenticationKey.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+authenticationType+" authentication-key \""+v+"\"")
		} else if v := rscData.AuthenticationKeyWO.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+authenticationType+" authentication-key \""+v+"\"")
		}
		if v := rscData.AuthenticationPassword.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+authenticationType+" authentication-password \""+v+"\"")
		} else if v := rscData.AuthenticationPasswordWO.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+authenticationType+" authentication-password \""+v+"\"")
		}
	} else {
		if rscData.PrivacyType.ValueString() != "privacy-none" {
			return path.Root("privacy_type"),
				errors.New("authentication should be configured before configuring the privacy")
		}
		if rscData.AuthenticationKey.ValueString() != "" {
			return path.Root("authentication_key"),
				errors.New("authentication_key not compatible when authentication_type = authentication-none")
		}
		if rscData.AuthenticationKeyWO.ValueString() != "" {
			return path.Root("authentication_key_wo"),
				errors.New("authentication_key_wo not compatible when authentication_type = authentication-none")
		}
		if rscData.AuthenticationPassword.ValueString() != "" {
			return path.Root("authentication_password"),
				errors.New("authentication_password not compatible when authentication_type = authentication-none")
		}
		if rscData.AuthenticationPasswordWO.ValueString() != "" {
			return path.Root("authentication_password_wo"),
				errors.New("authentication_password_wo not compatible when authentication_type = authentication-none")
		}
		configSet = append(configSet, setPrefix+"authentication-none")
	}
	if privacyType := rscData.PrivacyType.ValueString(); privacyType != "privacy-none" {
		if rscData.PrivacyKey.ValueString() == "" && rscData.PrivacyPassword.ValueString() == "" &&
			rscData.PrivacyKeyWO.ValueString() == "" && rscData.PrivacyPasswordWO.ValueString() == "" {
			return path.Root("privacy_type"),
				errors.New("one of privacy_key, privacy_password, privacy_key_wo or privacy_password_wo " +
					"must be specified when privacy_type != privacy-none")
		}
		if v := rscData.PrivacyKey.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+privacyType+" privacy-key \""+v+"\"")
		} else if v := rscData.PrivacyKeyWO.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+privacyType+" privacy-key \""+v+"\"")
		}
		if v := rscData.PrivacyPassword.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+privacyType+" privacy-password \""+v+"\"")
		} else if v := rscData.PrivacyPasswordWO.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+privacyType+" privacy-password \""+v+"\"")
		}
	} else {
		if rscData.PrivacyKey.ValueString() != "" {
			return path.Root("privacy_key"),
				errors.New("privacy_key not compatible when privacy_type = privacy-none")
		}
		if rscData.PrivacyKeyWO.ValueString() != "" {
			return path.Root("privacy_key_wo"),
				errors.New("privacy_key_wo not compatible when privacy_type = privacy-none")
		}
		if rscData.PrivacyPassword.ValueString() != "" {
			return path.Root("privacy_password"),
				errors.New("privacy_password not compatible when privacy_type = privacy-none")
		}
		if rscData.PrivacyPasswordWO.ValueString() != "" {
			return path.Root("privacy_password_wo"),
				errors.New("privacy_password_wo not compatible when privacy_type = privacy-none")
		}
		configSet = append(configSet, setPrefix+"privacy-none")
	}

	return path.Empty(), junSess.ConfigSet(ctx, configSet)
}

func (rscData *snmpV3UsmUserData) read(
	ctx context.Context, name, engineType, engineID string, junSess *junos.Session,
) error {
	showPrefix := junos.CmdShowConfig + "snmp v3 usm "
	switch engineType {
	case "remote":
		showPrefix += "remote-engine \"" + engineID + "\" "
	default:
		if engineType != "local" {
			engineType = "local"
		}
		showPrefix += "local-engine "
	}
	showConfig, err := junSess.Command(ctx, showPrefix+
		"user \""+name+"\""+junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Name = types.StringValue(name)
		rscData.EngineType = types.StringValue(engineType)
		if engineType == "remote" {
			rscData.EngineID = types.StringValue(engineID)
		}
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
			case strings.HasPrefix(itemTrim, "authentication-"):
				itemTrimFields := strings.Split(itemTrim, " ")
				rscData.AuthenticationType = types.StringValue(itemTrimFields[0])
				if balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" authentication-key ") {
					rscData.AuthenticationKey, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "authentication-key")
					if err != nil {
						return err
					}
				}
			case strings.HasPrefix(itemTrim, "privacy-"):
				itemTrimFields := strings.Split(itemTrim, " ")
				rscData.PrivacyType = types.StringValue(itemTrimFields[0])
				if balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" privacy-key ") {
					rscData.PrivacyKey, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "privacy-key")
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (rscData *snmpV3UsmUserData) readPrivateToState(
	ctx context.Context, junSess *junos.Session, private privateStateSetter,
) error {
	showPrefix := junos.CmdShowConfig + "snmp v3 usm "
	switch engineType := rscData.EngineType.ValueString(); engineType {
	case "remote":
		showPrefix += "remote-engine \"" + rscData.EngineID.ValueString() + "\" "
	default:
		showPrefix += "local-engine "
	}
	showConfig, err := junSess.Command(ctx, showPrefix+
		"user \""+rscData.Name.ValueString()+"\""+junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	// the private state is stored in the Terraform state, so don't keep the keys read on the device
	// when the write-only arguments are used: they are only compared with the key generated by
	// authentication_password or privacy_password, which cannot be set in this case
	authenticationWriteOnly := !rscData.AuthenticationKeyWOVersion.IsNull() ||
		!rscData.AuthenticationPasswordWOVersion.IsNull()
	privacyWriteOnly := !rscData.PrivacyKeyWOVersion.IsNull() ||
		!rscData.PrivacyPasswordWOVersion.IsNull()

	var privateState snmpV3UsmUserPrivateState
	if showConfig != junos.EmptyW {
		for item := range strings.SplitSeq(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case !authenticationWriteOnly && strings.HasPrefix(itemTrim, "authentication-"):
				itemTrimFields := strings.Split(itemTrim, " ")
				if balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" authentication-key ") {
					authenticationKey, err := junSess.JunosDecode(strings.Trim(itemTrim, "\""), "authentication-key")
					if err != nil {
						return err
					}
					privateState.AuthenticationKey = authenticationKey.ValueString()
				}
			case !privacyWriteOnly && strings.HasPrefix(itemTrim, "privacy-"):
				itemTrimFields := strings.Split(itemTrim, " ")
				if balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" privacy-key ") {
					privacyKey, err := junSess.JunosDecode(strings.Trim(itemTrim, "\""), "privacy-key")
					if err != nil {
						return err
					}
					privateState.PrivacyKey = privacyKey.ValueString()
				}
			}
		}
	}

	privateStateJSON, err := json.Marshal(privateState)
	if err != nil {
		return fmt.Errorf("internal error: json marshal private state: %w", err)
	}
	private.SetKey(ctx, privateState.key(), privateStateJSON)

	return nil
}

func (rscData *snmpV3UsmUserData) del(
	ctx context.Context, junSess *junos.Session,
) error {
	delPrefix := junos.DeleteLS + "snmp v3 usm "
	switch v := rscData.EngineType.ValueString(); v {
	case "local":
		delPrefix += "local-engine "
	case "remote":
		delPrefix += "remote-engine \"" + rscData.EngineID.ValueString() + "\" "
	default:
		return fmt.Errorf("can't del config with engine_type %q", v)
	}

	configSet := []string{
		delPrefix + "user \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(ctx, configSet)
}
