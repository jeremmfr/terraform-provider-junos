package providerfwk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
			"authentication_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "User's authentication password.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 1024),
					tfvalidator.StringDoubleQuoteExclusion(),
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
			"privacy_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "User's privacy password.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 1024),
					tfvalidator.StringDoubleQuoteExclusion(),
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
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	EngineType             types.String `tfsdk:"engine_type"`
	EngineID               types.String `tfsdk:"engine_id"`
	AuthenticationKey      types.String `tfsdk:"authentication_key"`
	AuthenticationPassword types.String `tfsdk:"authentication_password"`
	AuthenticationType     types.String `tfsdk:"authentication_type"`
	PrivacyKey             types.String `tfsdk:"privacy_key"`
	PrivacyPassword        types.String `tfsdk:"privacy_password"`
	PrivacyType            types.String `tfsdk:"privacy_type"`
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
		return
	}

	if data != nil {
		if err := json.Unmarshal(data, ste); err != nil {
			diags.AddError(tfdiag.GetPrivateStateErrSummary, fmt.Sprintf("json unmarshal: %s", err))
		}
	}

	return
}

func (rsc *snmpV3UsmUser) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config snmpV3UsmUserData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	engineType := "<unknown>"
	if config.EngineType.IsNull() {
		engineType = "local"
	} else if !config.EngineType.IsUnknown() {
		engineType = config.EngineType.ValueString()
	}
	switch engineType {
	case "local":
		if !config.EngineID.IsNull() && !config.EngineID.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("engine_id"),
				tfdiag.ConflictConfigErrSummary,
				"could not create "+rsc.junosName()+" with engine_type = local and engine_id specified",
			)
		}
	case "remote":
		if config.EngineID.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("engine_id"),
				tfdiag.MissingConfigErrSummary,
				"could not create "+rsc.junosName()+" with engine_type = remote and empty engine_id",
			)
		}
	}

	if !config.AuthenticationKey.IsNull() && !config.AuthenticationKey.IsUnknown() &&
		!config.AuthenticationPassword.IsNull() && !config.AuthenticationPassword.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("authentication_key"),
			tfdiag.ConflictConfigErrSummary,
			"authentication_key and authentication_password cannot be configured together",
		)
	}
	if !config.PrivacyKey.IsNull() && !config.PrivacyKey.IsUnknown() &&
		!config.PrivacyPassword.IsNull() && !config.PrivacyPassword.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("privacy_key"),
			tfdiag.ConflictConfigErrSummary,
			"privacy_key and privacy_password cannot be configured together",
		)
	}
	if !config.AuthenticationType.IsNull() && !config.AuthenticationType.IsUnknown() &&
		config.AuthenticationType.ValueString() != "authentication-none" {
		if config.AuthenticationKey.IsNull() && config.AuthenticationPassword.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_type"),
				tfdiag.MissingConfigErrSummary,
				"authentication_key or authentication_password must be specified when authentication_type != authentication-none",
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
		if !config.AuthenticationPassword.IsNull() && !config.AuthenticationPassword.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication_password"),
				tfdiag.ConflictConfigErrSummary,
				"authentication_password not compatible when authentication_type = authentication-none",
			)
		}
	}
	if !config.PrivacyType.IsNull() && !config.PrivacyType.IsUnknown() &&
		config.PrivacyType.ValueString() != "privacy-none" {
		if config.PrivacyKey.IsNull() && config.PrivacyPassword.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_type"),
				tfdiag.MissingConfigErrSummary,
				"privacy_key or privacy_password must be specified when privacy_type != privacy-none",
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
		if !config.PrivacyPassword.IsNull() && !config.PrivacyPassword.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("privacy_password"),
				tfdiag.ConflictConfigErrSummary,
				"privacy_password not compatible when privacy_type = privacy-none",
			)
		}
	}
}

func (rsc *snmpV3UsmUser) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan snmpV3UsmUserData
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
	_ context.Context, name, engineType, engineID string, junSess *junos.Session,
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
	showConfig, err := junSess.Command(showPrefix +
		"user \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
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
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
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
		if rscData.AuthenticationKey.ValueString() == "" && rscData.AuthenticationPassword.ValueString() == "" {
			return path.Root("authentication_type"),
				errors.New("authentication_key or authentication_password must be specified " +
					"when authentication_type != authentication-none")
		}
		if v := rscData.AuthenticationKey.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+authenticationType+" authentication-key \""+v+"\"")
		}
		if v := rscData.AuthenticationPassword.ValueString(); v != "" {
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
		if rscData.AuthenticationPassword.ValueString() != "" {
			return path.Root("authentication_password"),
				errors.New("authentication_password not compatible when authentication_type = authentication-none")
		}
		configSet = append(configSet, setPrefix+"authentication-none")
	}
	if privacyType := rscData.PrivacyType.ValueString(); privacyType != "privacy-none" {
		if rscData.PrivacyKey.ValueString() == "" && rscData.PrivacyPassword.ValueString() == "" {
			return path.Root("privacy_type"),
				errors.New("privacy_key or privacy_password must be specified when privacy_type != privacy-none")
		}
		if v := rscData.PrivacyKey.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+privacyType+" privacy-key \""+v+"\"")
		}
		if v := rscData.PrivacyPassword.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+privacyType+" privacy-password \""+v+"\"")
		}
	} else {
		if rscData.PrivacyKey.ValueString() != "" {
			return path.Root("privacy_key"),
				errors.New("privacy_key not compatible when privacy_type = privacy-none")
		}
		if rscData.PrivacyPassword.ValueString() != "" {
			return path.Root("privacy_password"),
				errors.New("privacy_password not compatible when privacy_type = privacy-none")
		}
		configSet = append(configSet, setPrefix+"privacy-none")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *snmpV3UsmUserData) read(
	_ context.Context, name, engineType, engineID string, junSess *junos.Session,
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
	showConfig, err := junSess.Command(showPrefix +
		"user \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
		for _, item := range strings.Split(showConfig, "\n") {
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
					rscData.AuthenticationKey, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "authentication-key")
					if err != nil {
						return err
					}
				}
			case strings.HasPrefix(itemTrim, "privacy-"):
				itemTrimFields := strings.Split(itemTrim, " ")
				rscData.PrivacyType = types.StringValue(itemTrimFields[0])
				if balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" privacy-key ") {
					rscData.PrivacyKey, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "privacy-key")
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
	showConfig, err := junSess.Command(showPrefix +
		"user \"" + rscData.Name.ValueString() + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	var privateState snmpV3UsmUserPrivateState
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
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
				if balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" authentication-key ") {
					authenticationKey, err := tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "authentication-key")
					if err != nil {
						return err
					}
					privateState.AuthenticationKey = authenticationKey.ValueString()
				}
			case strings.HasPrefix(itemTrim, "privacy-"):
				itemTrimFields := strings.Split(itemTrim, " ")
				if balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" privacy-key ") {
					privacyKey, err := tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "privacy-key")
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
	_ context.Context, junSess *junos.Session,
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

	return junSess.ConfigSet(configSet)
}
