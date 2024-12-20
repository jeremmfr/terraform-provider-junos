package providerfwk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"golang.org/x/crypto/ssh"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &systemRootAuthentication{}
	_ resource.ResourceWithConfigure      = &systemRootAuthentication{}
	_ resource.ResourceWithValidateConfig = &systemRootAuthentication{}
	_ resource.ResourceWithImportState    = &systemRootAuthentication{}
)

type systemRootAuthentication struct {
	client *junos.Client
}

func newSystemRootAuthenticationResource() resource.Resource {
	return &systemRootAuthentication{}
}

func (rsc *systemRootAuthentication) typeName() string {
	return providerName + "_system_root_authentication"
}

func (rsc *systemRootAuthentication) junosName() string {
	return "system root-authentication"
}

func (rsc *systemRootAuthentication) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *systemRootAuthentication) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *systemRootAuthentication) Configure(
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

func (rsc *systemRootAuthentication) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `system_root_authentication`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"encrypted_password": schema.StringAttribute{
				Optional:    true,
				Description: "Encrypted password string.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					tfvalidator.StringDoubleQuoteExclusion(),
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
			"no_public_keys": schema.BoolAttribute{
				Optional:    true,
				Description: "Disables ssh public key based authentication.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
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
	}
}

type systemRootAuthenticationData struct {
	ID                types.String   `tfsdk:"id"`
	EncryptedPassword types.String   `tfsdk:"encrypted_password"`
	PlainTextPassword types.String   `tfsdk:"plain_text_password"`
	NoPublicKeys      types.Bool     `tfsdk:"no_public_keys"`
	SSHPublicKeys     []types.String `tfsdk:"ssh_public_keys"`
}

type systemRootAuthenticationConfig struct {
	ID                types.String `tfsdk:"id"`
	EncryptedPassword types.String `tfsdk:"encrypted_password"`
	PlainTextPassword types.String `tfsdk:"plain_text_password"`
	NoPublicKeys      types.Bool   `tfsdk:"no_public_keys"`
	SSHPublicKeys     types.Set    `tfsdk:"ssh_public_keys"`
}

type systemRootAuthenticationPrivateState struct {
	EncryptedPassword string `json:"encrypted_password"`
}

func (ste *systemRootAuthenticationPrivateState) key() string {
	return "v0"
}

func (ste *systemRootAuthenticationPrivateState) get(
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

func (rsc *systemRootAuthentication) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config systemRootAuthenticationConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.EncryptedPassword.IsNull() &&
		config.PlainTextPassword.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("plain_text_password"),
			tfdiag.ConflictConfigErrSummary,
			"encrypted_password or plain_text_password must be specified",
		)
	}
	if !config.EncryptedPassword.IsNull() &&
		!config.EncryptedPassword.IsUnknown() &&
		!config.PlainTextPassword.IsNull() &&
		!config.PlainTextPassword.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("encrypted_password"),
			tfdiag.ConflictConfigErrSummary,
			"encrypted_password and plain_text_password cannot be configured together",
		)
	}
	if !config.NoPublicKeys.IsNull() &&
		!config.NoPublicKeys.IsUnknown() &&
		!config.SSHPublicKeys.IsNull() &&
		!config.SSHPublicKeys.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("no_public_keys"),
			tfdiag.ConflictConfigErrSummary,
			"no_public_keys and ssh_public_keys cannot be configured together",
		)
	}
	if !config.SSHPublicKeys.IsNull() &&
		!config.SSHPublicKeys.IsUnknown() {
		var configSSHPublicKeys []types.String
		asDiags := config.SSHPublicKeys.ElementsAs(ctx, &configSSHPublicKeys, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		for _, v := range configSSHPublicKeys {
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
					path.Root("ssh_public_keys"),
					tfdiag.CompatibilityErrSummary,
					fmt.Sprintf("format in key %q not supported in ssh_public_keys", key),
				)
			}
		}
	}
}

func (rsc *systemRootAuthentication) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan systemRootAuthenticationData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.EncryptedPassword.ValueString() == "" &&
		plan.PlainTextPassword.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("plain_text_password"),
			"Empty Password",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "encrypted_password and plain_text_password"),
		)

		return
	}

	if rsc.client.FakeCreateSetFile() {
		junSess := rsc.client.NewSessionWithoutNetconf(ctx)

		if plan.PlainTextPassword.ValueString() != "" {
			// To be able detect a plain text password not accepted by system
			if err := plan.delPassword(ctx, junSess); err != nil {
				resp.Diagnostics.AddError("Pre Config Set Error", err.Error())

				return
			}
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
			} else {
				resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
			}

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigUnlockWarnSummary, junSess.ConfigUnlock())...)
	}()

	if plan.PlainTextPassword.ValueString() != "" {
		// To be able detect a plain text password not accepted by system
		if err := plan.delPassword(ctx, junSess); err != nil {
			resp.Diagnostics.AddError("Pre Config Set Error", err.Error())

			return
		}
	}

	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf(ctx, "create resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := plan.readPrivateToState(ctx, junSess, resp.Private); err != nil {
		resp.Diagnostics.AddError(tfdiag.ReadPrivateToStateErrSummary, err.Error())
	}
}

func (rsc *systemRootAuthentication) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data systemRootAuthenticationData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadWithoutArg = &data
	defaultResourceRead(
		ctx,
		rsc,
		nil,
		&data,
		func() {
			var privateState systemRootAuthenticationPrivateState
			resp.Diagnostics.Append(privateState.get(ctx, req.Private)...)
			if resp.Diagnostics.HasError() {
				return
			}

			if data.EncryptedPassword.ValueString() != "" &&
				state.PlainTextPassword.ValueString() != "" {
				if privateState.EncryptedPassword != "" {
					if privateState.EncryptedPassword == data.EncryptedPassword.ValueString() {
						data.PlainTextPassword = state.PlainTextPassword
					} else {
						data.PlainTextPassword = types.StringValue(`?`)
					}
					data.EncryptedPassword = types.StringNull()
				} else {
					data.PlainTextPassword = state.PlainTextPassword
					data.EncryptedPassword = types.StringNull()
				}
			}
		},
		resp,
	)
}

func (rsc *systemRootAuthentication) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state systemRootAuthenticationData
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

func (rsc *systemRootAuthentication) Delete(
	_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse,
) {
	// no-op
}

func (rsc *systemRootAuthentication) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data systemRootAuthenticationData

	var _ resourceDataReadWithoutArg = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		"don't find "+rsc.junosName(),
	)
}

func (rscData *systemRootAuthenticationData) fillID() {
	rscData.ID = types.StringValue("system_root_authentication")
}

func (rscData *systemRootAuthenticationData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *systemRootAuthenticationData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set system root-authentication "

	if v := rscData.PlainTextPassword.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"plain-text-password-value \""+v+"\"")
	} else {
		configSet = append(configSet, setPrefix+"encrypted-password \""+rscData.EncryptedPassword.ValueString()+"\"")
	}
	if rscData.NoPublicKeys.ValueBool() {
		configSet = append(configSet, setPrefix+"no-public-keys")
	}
	for _, v := range rscData.SSHPublicKeys {
		key := v.ValueString()
		switch {
		case strings.HasPrefix(key, ssh.KeyAlgoDSA):
			configSet = append(configSet, setPrefix+"ssh-dsa \""+key+"\"")
		case strings.HasPrefix(key, ssh.KeyAlgoRSA):
			configSet = append(configSet, setPrefix+"ssh-rsa \""+key+"\"")
		case strings.HasPrefix(key, ssh.KeyAlgoECDSA256),
			strings.HasPrefix(key, ssh.KeyAlgoECDSA384),
			strings.HasPrefix(key, ssh.KeyAlgoECDSA521):
			configSet = append(configSet, setPrefix+"ssh-ecdsa \""+key+"\"")
		case strings.HasPrefix(key, ssh.KeyAlgoED25519):
			configSet = append(configSet, setPrefix+"ssh-ed25519 \""+key+"\"")
		default:
			return path.Root("ssh_public_keys"),
				fmt.Errorf("format in key %q not supported in ssh_public_keys", key)
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *systemRootAuthenticationData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system root-authentication" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
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
			case balt.CutPrefixInString(&itemTrim, "encrypted-password "):
				rscData.EncryptedPassword = types.StringValue(strings.Trim(itemTrim, "\""))
			case itemTrim == "no-public-keys":
				rscData.NoPublicKeys = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "ssh-dsa "),
				balt.CutPrefixInString(&itemTrim, "ssh-ecdsa "),
				balt.CutPrefixInString(&itemTrim, "ssh-ed25519 "),
				balt.CutPrefixInString(&itemTrim, "ssh-rsa "):
				rscData.SSHPublicKeys = append(rscData.SSHPublicKeys,
					types.StringValue(strings.Trim(itemTrim, "\"")))
			}
		}
	}

	return nil
}

func (rscData *systemRootAuthenticationData) readPrivateToState(
	ctx context.Context, junSess *junos.Session, private privateStateSetter,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"system root-authentication" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	var privateState systemRootAuthenticationPrivateState
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			if balt.CutPrefixInString(&itemTrim, "encrypted-password ") {
				privateState.EncryptedPassword = strings.Trim(itemTrim, "\"")
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

func (rscData *systemRootAuthenticationData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete system root-authentication",
	}

	return junSess.ConfigSet(configSet)
}

func (rscData *systemRootAuthenticationData) delPassword(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete system root-authentication encrypted-password",
	}

	return junSess.ConfigSet(configSet)
}
