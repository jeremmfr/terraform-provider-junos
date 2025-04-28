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
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"golang.org/x/crypto/ssh"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &systemLoginUser{}
	_ resource.ResourceWithConfigure      = &systemLoginUser{}
	_ resource.ResourceWithValidateConfig = &systemLoginUser{}
	_ resource.ResourceWithImportState    = &systemLoginUser{}
	_ resource.ResourceWithUpgradeState   = &systemLoginUser{}
)

type systemLoginUser struct {
	client *junos.Client
}

func newSystemLoginUserResource() resource.Resource {
	return &systemLoginUser{}
}

func (rsc *systemLoginUser) typeName() string {
	return providerName + "_system_login_user"
}

func (rsc *systemLoginUser) junosName() string {
	return "system login user"
}

func (rsc *systemLoginUser) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *systemLoginUser) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *systemLoginUser) Configure(
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

func (rsc *systemLoginUser) Schema(
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
				Description: "The name of system login user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
				},
			},
			"class": schema.StringAttribute{
				Required:    true,
				Description: "Login class.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"uid": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "User identifier (uid).",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Between(100, 64000),
				},
			},
			"cli_prompt": schema.StringAttribute{
				Optional:    true,
				Description: "Cli prompt name for this user.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"full_name": schema.StringAttribute{
				Optional:    true,
				Description: "Full name.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"authentication": schema.SingleNestedBlock{
				Description: "Authentication method.",
				Attributes: map[string]schema.Attribute{
					"encrypted_password": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Encrypted password string.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 128),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"no_public_keys": schema.BoolAttribute{
						Optional:    true,
						Description: "Disables ssh public key based authentication.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"plain_text_password": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Plain text password.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 128),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"ssh_public_keys": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Secure shell (ssh) public key string.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.NoNullValues(),
							setvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							),
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

type systemLoginUserData struct {
	ID             types.String                        `tfsdk:"id"`
	Name           types.String                        `tfsdk:"name"`
	Class          types.String                        `tfsdk:"class"`
	UID            types.Int64                         `tfsdk:"uid"`
	CliPrompt      types.String                        `tfsdk:"cli_prompt"`
	FullName       types.String                        `tfsdk:"full_name"`
	Authentication *systemLoginUserBlockAuthentication `tfsdk:"authentication"`
}

type systemLoginUserConfig struct {
	ID             types.String                              `tfsdk:"id"`
	Name           types.String                              `tfsdk:"name"`
	Class          types.String                              `tfsdk:"class"`
	UID            types.Int64                               `tfsdk:"uid"`
	CliPrompt      types.String                              `tfsdk:"cli_prompt"`
	FullName       types.String                              `tfsdk:"full_name"`
	Authentication *systemLoginUserBlockAuthenticationConfig `tfsdk:"authentication"`
}

type systemLoginUserBlockAuthentication struct {
	EncryptedPassword types.String   `tfsdk:"encrypted_password"`
	NoPublicKeys      types.Bool     `tfsdk:"no_public_keys"`
	PlainTextPassword types.String   `tfsdk:"plain_text_password"`
	SSHPublicKeys     []types.String `tfsdk:"ssh_public_keys"`
}

func (block *systemLoginUserBlockAuthentication) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemLoginUserBlockAuthenticationConfig struct {
	EncryptedPassword types.String `tfsdk:"encrypted_password"`
	NoPublicKeys      types.Bool   `tfsdk:"no_public_keys"`
	PlainTextPassword types.String `tfsdk:"plain_text_password"`
	SSHPublicKeys     types.Set    `tfsdk:"ssh_public_keys"`
}

func (block *systemLoginUserBlockAuthenticationConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type systemLoginUserPrivateState struct {
	AuthenticationEncryptedPassword string `json:"authentication_encrypted_password"`
}

func (ste *systemLoginUserPrivateState) key() string {
	return "v0"
}

func (ste *systemLoginUserPrivateState) get(
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

func (rsc *systemLoginUser) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config systemLoginUserConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Authentication != nil {
		if config.Authentication.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"authentication block is empty",
			)
		}
		if !config.Authentication.EncryptedPassword.IsNull() &&
			!config.Authentication.EncryptedPassword.IsUnknown() &&
			!config.Authentication.PlainTextPassword.IsNull() &&
			!config.Authentication.PlainTextPassword.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication").AtName("encrypted_password"),
				tfdiag.ConflictConfigErrSummary,
				"encrypted_password and plain_text_password cannot be configured together"+
					" in authentication block",
			)
		}
		if !config.Authentication.NoPublicKeys.IsNull() &&
			!config.Authentication.NoPublicKeys.IsUnknown() &&
			!config.Authentication.SSHPublicKeys.IsNull() &&
			!config.Authentication.SSHPublicKeys.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("authentication").AtName("no_public_keys"),
				tfdiag.ConflictConfigErrSummary,
				"no_public_keys and ssh_public_keys cannot be configured together"+
					" in authentication block",
			)
		}
		if !config.Authentication.SSHPublicKeys.IsNull() &&
			!config.Authentication.SSHPublicKeys.IsUnknown() {
			var configAuthenticationSSHPublicKeys []types.String
			asDiags := config.Authentication.SSHPublicKeys.ElementsAs(ctx, &configAuthenticationSSHPublicKeys, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}
			for _, v := range configAuthenticationSSHPublicKeys {
				if v.IsUnknown() {
					continue
				}
				key := v.ValueString()
				switch {
				case strings.HasPrefix(key, ssh.KeyAlgoDSA):
					continue
				case strings.HasPrefix(key, ssh.KeyAlgoRSA):
					continue
				case strings.HasPrefix(key, ssh.KeyAlgoECDSA256),
					strings.HasPrefix(key, ssh.KeyAlgoECDSA384),
					strings.HasPrefix(key, ssh.KeyAlgoECDSA521):
					continue
				case strings.HasPrefix(key, ssh.KeyAlgoED25519):
					continue
				default:
					resp.Diagnostics.AddAttributeError(
						path.Root("authentication").AtName("ssh_public_keys"),
						tfdiag.CompatibilityErrSummary,
						fmt.Sprintf("format in key %q not supported in ssh_public_keys"+
							" in authentication block", key),
					)
				}
			}
		}
	}
}

func (rsc *systemLoginUser) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemLoginUserData
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
	if plan.Class.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("class"),
			"Empty Class",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "class"),
		)

		return
	}

	var _ resourceDataReadComputed = &plan
	var _ resourceDataReadPrivateToState = &plan
	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			userExists, err := checkSystemLoginUserExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if userExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			userExists, err := checkSystemLoginUserExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !userExists {
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

func (rsc *systemLoginUser) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemLoginUserData
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
		func() {
			var privateState systemLoginUserPrivateState
			resp.Diagnostics.Append(privateState.get(ctx, req.Private)...)
			if resp.Diagnostics.HasError() {
				return
			}

			if data.Authentication != nil &&
				data.Authentication.EncryptedPassword.ValueString() != "" &&
				state.Authentication != nil &&
				state.Authentication.PlainTextPassword.ValueString() != "" {
				if privateState.AuthenticationEncryptedPassword != "" {
					if privateState.AuthenticationEncryptedPassword == data.Authentication.EncryptedPassword.ValueString() {
						data.Authentication.PlainTextPassword = state.Authentication.PlainTextPassword
					} else {
						data.Authentication.PlainTextPassword = types.StringValue(`?`)
					}
					data.Authentication.EncryptedPassword = types.StringNull()
				} else {
					data.Authentication.PlainTextPassword = state.Authentication.PlainTextPassword
					data.Authentication.EncryptedPassword = types.StringNull()
				}
			}
		},
		resp,
	)
}

func (rsc *systemLoginUser) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemLoginUserData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadComputed = &plan
	var _ resourceDataReadPrivateToState = &plan
	defaultResourceUpdate(
		ctx,
		rsc,
		&state,
		&plan,
		resp,
	)
}

func (rsc *systemLoginUser) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state systemLoginUserData
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

func (rsc *systemLoginUser) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data systemLoginUserData

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

func checkSystemLoginUserExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system login user " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *systemLoginUserData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *systemLoginUserData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemLoginUserData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set system login user " + rscData.Name.ValueString() + " "

	configSet := []string{
		setPrefix + "class " + rscData.Class.ValueString(),
	}

	if !rscData.UID.IsNull() {
		configSet = append(configSet, setPrefix+"uid "+
			utils.ConvI64toa(rscData.UID.ValueInt64()))
	}
	if v := rscData.CliPrompt.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"cli prompt \""+v+"\"")
	}
	if v := rscData.FullName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"full-name \""+v+"\"")
	}
	if rscData.Authentication != nil {
		if rscData.Authentication.isEmpty() {
			return path.Root("authentication").AtName("*"),
				errors.New("authentication block is empty")
		}

		if v := rscData.Authentication.PlainTextPassword.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"authentication plain-text-password-value \""+v+"\"")
		} else if v := rscData.Authentication.EncryptedPassword.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"authentication encrypted-password \""+v+"\"")
		}
		if rscData.Authentication.NoPublicKeys.ValueBool() {
			configSet = append(configSet, setPrefix+"authentication no-public-keys")
		}
		for _, v := range rscData.Authentication.SSHPublicKeys {
			key := v.ValueString()
			switch {
			case strings.HasPrefix(key, ssh.KeyAlgoDSA):
				configSet = append(configSet, setPrefix+"authentication ssh-dsa \""+key+"\"")
			case strings.HasPrefix(key, ssh.KeyAlgoRSA):
				configSet = append(configSet, setPrefix+"authentication ssh-rsa \""+key+"\"")
			case strings.HasPrefix(key, ssh.KeyAlgoECDSA256),
				strings.HasPrefix(key, ssh.KeyAlgoECDSA384),
				strings.HasPrefix(key, ssh.KeyAlgoECDSA521):
				configSet = append(configSet, setPrefix+"authentication ssh-ecdsa \""+key+"\"")
			case strings.HasPrefix(key, ssh.KeyAlgoED25519):
				configSet = append(configSet, setPrefix+"authentication ssh-ed25519 \""+key+"\"")
			default:
				return path.Root("authentication").AtName("ssh_public_keys"),
					fmt.Errorf("format in key %q not supported in ssh_public_keys"+
						" in authentication block", key)
			}
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *systemLoginUserData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system login user " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "class "):
				rscData.Class = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "uid "):
				rscData.UID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "cli prompt "):
				rscData.CliPrompt = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "full-name "):
				rscData.FullName = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "authentication "):
				if rscData.Authentication == nil {
					rscData.Authentication = &systemLoginUserBlockAuthentication{}
				}

				switch {
				case balt.CutPrefixInString(&itemTrim, "encrypted-password "):
					rscData.Authentication.EncryptedPassword = types.StringValue(strings.Trim(itemTrim, "\""))
				case itemTrim == "no-public-keys":
					rscData.Authentication.NoPublicKeys = types.BoolValue(true)
				case balt.CutPrefixInString(&itemTrim, "ssh-dsa "),
					balt.CutPrefixInString(&itemTrim, "ssh-ecdsa "),
					balt.CutPrefixInString(&itemTrim, "ssh-ed25519 "),
					balt.CutPrefixInString(&itemTrim, "ssh-rsa "):
					rscData.Authentication.SSHPublicKeys = append(rscData.Authentication.SSHPublicKeys,
						types.StringValue(strings.Trim(itemTrim, "\"")),
					)
				}
			}
		}
	}

	return nil
}

func (rscData *systemLoginUserData) readComputed(
	_ context.Context, junSess *junos.Session,
) error {
	defer func() {
		// set unknown to null if still unknown after reading config
		if rscData.UID.IsUnknown() {
			rscData.UID = types.Int64Null()
		}
	}()

	if !junSess.HasNetconf() {
		return nil
	}

	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system login user " + rscData.Name.ValueString() + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "uid ") && rscData.UID.IsUnknown() {
				rscData.UID, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (rscData *systemLoginUserData) readPrivateToState(
	ctx context.Context, junSess *junos.Session, private privateStateSetter,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system login user " + rscData.Name.ValueString() + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	var privateState systemLoginUserPrivateState
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "authentication encrypted-password ") {
				privateState.AuthenticationEncryptedPassword = strings.Trim(itemTrim, "\"")
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

func (rscData *systemLoginUserData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete system login user " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
