package provider

import (
	"context"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &securityAuthenticationKeyChain{}
	_ resource.ResourceWithConfigure      = &securityAuthenticationKeyChain{}
	_ resource.ResourceWithValidateConfig = &securityAuthenticationKeyChain{}
	_ resource.ResourceWithImportState    = &securityAuthenticationKeyChain{}
)

type securityAuthenticationKeyChain struct {
	client *junos.Client
}

func newSecurityAuthenticationKeyChainResource() resource.Resource {
	return &securityAuthenticationKeyChain{}
}

func (rsc *securityAuthenticationKeyChain) typeName() string {
	return providerName + "_security_authentication_key_chain"
}

func (rsc *securityAuthenticationKeyChain) junosName() string {
	return "security authentication-key-chains key-chain"
}

func (rsc *securityAuthenticationKeyChain) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *securityAuthenticationKeyChain) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *securityAuthenticationKeyChain) Configure(
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

func (rsc *securityAuthenticationKeyChain) Schema(
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
				Description: "Name of authentication key chain.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
					tfvalidator.StringSpaceExclusion(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of this authentication-key-chain.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 80),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"key_secret_wo": schema.MapNestedAttribute{
				Optional: true,
				Description: "For each authentication element identifier, " +
					"authentication key not stored in state.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Required:    true,
							Sensitive:   true,
							WriteOnly:   true,
							Description: "Authentication key.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 126),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"version": schema.Int64Attribute{
							Required:    true,
							Description: "Version of `value` to trigger the sending of its value.",
						},
					},
				},
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
					mapvalidator.KeysAre(
						stringvalidator.RegexMatches(regexp.MustCompile(`^\d+$`),
							"must be an authentication element identifier of a key block",
						),
					),
				},
			},
			"tolerance": schema.Int64Attribute{
				Optional:    true,
				Description: "Clock skew tolerance (seconds).",
				Validators: []validator.Int64{
					int64validator.Between(0, 4294967295),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"key": schema.SetNestedBlock{
				Description: "Authentication element configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Required:    true,
							Description: "Authentication element identifier.",
							Validators: []validator.Int64{
								int64validator.Between(0, 63),
							},
						},
						"secret": schema.StringAttribute{
							Optional:    true,
							Sensitive:   true,
							Description: "Authentication key.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 126),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"start_time": schema.StringAttribute{
							Required:    true,
							Description: "Start time for key transmission (YYYY-MM-DD.HH:MM:SS).",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(
									`^\d{4}\-\d\d?\-\d\d?\.\d{2}:\d{2}:\d{2}$`),
									"must be in the format 'YYYY-MM-DD.HH:MM:SS'",
								),
							},
						},
						"algorithm": schema.StringAttribute{
							Optional:    true,
							Description: "Authentication algorithm.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"ao_cryptographic_algorithm": schema.StringAttribute{
							Optional:    true,
							Description: "Cryptographic algorithm for TCP-AO Traffic key and MAC digest generation.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"ao_recv_id": schema.Int64Attribute{
							Optional:    true,
							Description: "Recv id for TCP-AO entry.",
							Validators: []validator.Int64{
								int64validator.Between(0, 255),
							},
						},
						"ao_send_id": schema.Int64Attribute{
							Optional:    true,
							Description: "Send id for TCP-AO entry.",
							Validators: []validator.Int64{
								int64validator.Between(0, 255),
							},
						},
						"ao_tcp_ao_option": schema.StringAttribute{
							Optional:    true,
							Description: "Include TCP-AO option within message header.",
							Validators: []validator.String{
								stringvalidator.OneOf("disabled", "enabled"),
							},
						},
						"key_name": schema.StringAttribute{
							Optional:    true,
							Description: "Key name in hexadecimal format used for macsec.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 64),
								tfvalidator.StringFormat(tfvalidator.HexadecimalFormat),
							},
						},
						"options": schema.StringAttribute{
							Optional:    true,
							Description: "Protocol's transmission encoding format.",
							Validators: []validator.String{
								stringvalidator.OneOf("basic", "isis-enhanced"),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.IsRequired(),
					setvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

type securityAuthenticationKeyChainData struct {
	ID          types.String                                             `tfsdk:"id"`
	Name        types.String                                             `tfsdk:"name"`
	Description types.String                                             `tfsdk:"description"`
	Tolerance   types.Int64                                              `tfsdk:"tolerance"`
	KeySecretWO map[string]securityAuthenticationKeyChainAttrKeySecretWO `tfsdk:"key_secret_wo"`
	Key         []securityAuthenticationKeyChainBlockKey                 `tfsdk:"key"`
}

type securityAuthenticationKeyChainConfig struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Tolerance   types.Int64  `tfsdk:"tolerance"`
	KeySecretWO types.Map    `tfsdk:"key_secret_wo"`
	Key         types.Set    `tfsdk:"key"`
}

type securityAuthenticationKeyChainAttrKeySecretWO struct {
	Value   types.String `tfsdk:"value"`
	Version types.Int64  `tfsdk:"version"`
}

type securityAuthenticationKeyChainBlockKey struct {
	ID                       types.Int64  `tfsdk:"id"                         tfdata:"identifier"`
	Secret                   types.String `tfsdk:"secret"`
	StartTime                types.String `tfsdk:"start_time"`
	Algorithm                types.String `tfsdk:"algorithm"`
	AOCryptographicAlgorithm types.String `tfsdk:"ao_cryptographic_algorithm"`
	AORecvID                 types.Int64  `tfsdk:"ao_recv_id"`
	AOSendID                 types.Int64  `tfsdk:"ao_send_id"`
	AOTcpAOOption            types.String `tfsdk:"ao_tcp_ao_option"`
	KeyName                  types.String `tfsdk:"key_name"`
	Options                  types.String `tfsdk:"options"`
}

func (rsc *securityAuthenticationKeyChain) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config securityAuthenticationKeyChainConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	keySecretWO := make(map[string]struct{})
	if !config.KeySecretWO.IsNull() &&
		!config.KeySecretWO.IsUnknown() {
		for id := range config.KeySecretWO.Elements() {
			keySecretWO[id] = struct{}{}
		}
	}

	if !config.Key.IsNull() &&
		!config.Key.IsUnknown() {
		var configKey []securityAuthenticationKeyChainBlockKey
		asDiags := config.Key.ElementsAs(ctx, &configKey, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		keyID := make(map[int64]struct{})
		unknownKeyID := false
		for _, block := range configKey {
			if block.ID.IsUnknown() {
				unknownKeyID = true

				continue
			}

			id := block.ID.ValueInt64()
			if _, ok := keyID[id]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("key"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple key blocks with the same id %d", id),
				)
			}
			keyID[id] = struct{}{}

			_, withSecretWO := keySecretWO[utils.ConvI64toa(id)]
			if block.Secret.IsNull() && !withSecretWO {
				resp.Diagnostics.AddAttributeError(
					path.Root("key"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("one of secret in key block %d"+
						" or an entry with this id in key_secret_wo must be specified", id),
				)
			}
			if !block.Secret.IsNull() && !block.Secret.IsUnknown() && withSecretWO {
				resp.Diagnostics.AddAttributeError(
					path.Root("key"),
					tfdiag.ConflictConfigErrSummary,
					fmt.Sprintf("only one of secret in key block %d"+
						" or an entry with this id in key_secret_wo can be specified", id),
				)
			}
		}

		// the key blocks are all known, so an entry of key_secret_wo
		// without a key block with the same id can be detected
		if !unknownKeyID {
			for _, id := range slices.Sorted(maps.Keys(keySecretWO)) {
				keyIDNum, err := utils.ConvAtoi64(id)
				if err != nil {
					continue // rejected by the map keys validator
				}
				if _, ok := keyID[keyIDNum]; !ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("key_secret_wo").AtMapKey(id),
						tfdiag.MissingConfigErrSummary,
						fmt.Sprintf("no key block with id %s to associate with this key_secret_wo entry", id),
					)
				}
			}
		}
	}
}

func (rsc *securityAuthenticationKeyChain) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityAuthenticationKeyChainData
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

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			chainExists, err := checkSecurityAuthenticationKeyChainExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if chainExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			chainExists, err := checkSecurityAuthenticationKeyChainExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !chainExists {
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

func (rsc *securityAuthenticationKeyChain) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data securityAuthenticationKeyChainData
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
			data.keepWriteOnly(&state)
		},
		resp,
	)
}

func (rsc *securityAuthenticationKeyChain) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityAuthenticationKeyChainData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(plan.getWriteOnly(ctx, req.Config)...)
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

func (rsc *securityAuthenticationKeyChain) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state securityAuthenticationKeyChainData
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

func (rsc *securityAuthenticationKeyChain) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data securityAuthenticationKeyChainData

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

func checkSecurityAuthenticationKeyChainExists(
	ctx context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"security authentication-key-chains key-chain \""+name+"\""+junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *securityAuthenticationKeyChainData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *securityAuthenticationKeyChainData) nullID() bool {
	return rscData.ID.IsNull()
}

// getWriteOnly read the write-only arguments from the configuration,
// their values aren't present in the plan or the state.
//
// Only the write-only argument of each key_secret_wo entry is read:
// reading the whole entry from the configuration would also read the version,
// which may be unknown at this time whereas it's resolved in the plan.
func (rscData *securityAuthenticationKeyChainData) getWriteOnly(
	ctx context.Context, config tfsdk.Config,
) (diags diag.Diagnostics) {
	for id, attribute := range rscData.KeySecretWO {
		diags.Append(config.GetAttribute(ctx,
			path.Root("key_secret_wo").AtMapKey(id).AtName("value"),
			&attribute.Value)...)
		rscData.KeySecretWO[id] = attribute
	}

	return diags
}

// keepWriteOnly carry over the write-only arguments from the state,
// only the version of each entry is present in it,
// and don't read the secrets in the standard arguments when the write-only ones are used.
//
// The keys of key_secret_wo are the identifiers of the key blocks using
// the write-only argument, the key blocks read on the device aren't in the order
// of the state so they are matched with their identifier.
func (rscData *securityAuthenticationKeyChainData) keepWriteOnly(state *securityAuthenticationKeyChainData) {
	rscData.KeySecretWO = state.KeySecretWO

	for i, block := range rscData.Key {
		if _, ok := state.KeySecretWO[utils.ConvI64toa(block.ID.ValueInt64())]; ok {
			rscData.Key[i].Secret = types.StringNull()
		}
	}
}

func (rscData *securityAuthenticationKeyChainData) set(
	ctx context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)
	setPrefix := "set security authentication-key-chains key-chain \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if !rscData.Tolerance.IsNull() {
		configSet = append(configSet, setPrefix+"tolerance "+
			utils.ConvI64toa(rscData.Tolerance.ValueInt64()))
	}
	// index the write-only secrets with the id of the key block they are associated with,
	// an entry without a key block with the same id is ignored here, it's detected in ValidateConfig
	keySecretWO := make(map[int64]string, len(rscData.KeySecretWO))
	for id, attribute := range rscData.KeySecretWO {
		keyID, err := utils.ConvAtoi64(id)
		if err != nil {
			return path.Root("key_secret_wo").AtMapKey(id),
				fmt.Errorf("invalid key id %q in key_secret_wo", id)
		}
		keySecretWO[keyID] = attribute.Value.ValueString()
	}

	keyID := make(map[int64]struct{})
	for _, block := range rscData.Key {
		id := block.ID.ValueInt64()
		if _, ok := keyID[id]; ok {
			return path.Root("key"),
				fmt.Errorf("multiple key blocks with the same id %d", id)
		}
		keyID[id] = struct{}{}

		blockConfigSet, err := block.configSet(setPrefix, keySecretWO[id])
		if err != nil {
			return path.Root("key"), err
		}
		configSet = append(configSet, blockConfigSet...)
	}

	return path.Empty(), junSess.ConfigSet(ctx, configSet)
}

func (block *securityAuthenticationKeyChainBlockKey) configSet(
	setPrefix, secretWO string,
) (
	[]string, error,
) {
	setPrefix += "key " + utils.ConvI64toa(block.ID.ValueInt64()) + " "

	configSet := make([]string, 2, 100)
	if v := block.Secret.ValueString(); v != "" {
		configSet[0] = setPrefix + "secret \"" + v + "\""
	} else if secretWO != "" {
		configSet[0] = setPrefix + "secret \"" + secretWO + "\""
	} else {
		return nil, fmt.Errorf("one of secret in key block %d"+
			" or an entry with this id in key_secret_wo must be specified", block.ID.ValueInt64())
	}
	configSet[1] = setPrefix + "start-time " + block.StartTime.ValueString()

	if v := block.Algorithm.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"algorithm "+v)
	}
	if v := block.AOCryptographicAlgorithm.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ao-attribute cryptographic-algorithm "+v)
	}
	if !block.AORecvID.IsNull() {
		configSet = append(configSet, setPrefix+"ao-attribute recv-id "+
			utils.ConvI64toa(block.AORecvID.ValueInt64()))
	}
	if !block.AOSendID.IsNull() {
		configSet = append(configSet, setPrefix+"ao-attribute send-id "+
			utils.ConvI64toa(block.AOSendID.ValueInt64()))
	}
	if v := block.AOTcpAOOption.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ao-attribute tcp-ao-option "+v)
	}
	if v := block.KeyName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"key-name "+v)
	}
	if v := block.Options.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"options "+v)
	}

	return configSet, nil
}

func (rscData *securityAuthenticationKeyChainData) read(
	ctx context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"security authentication-key-chains key-chain \""+name+"\""+junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "tolerance "):
				rscData.Tolerance, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "key "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var key securityAuthenticationKeyChainBlockKey
				keyID, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
				if err != nil {
					return err
				}
				rscData.Key, key = tfdata.ExtractBlock(rscData.Key, keyID)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")

				if err := key.read(itemTrim, junSess); err != nil {
					return err
				}
				rscData.Key = append(rscData.Key, key)
			}
		}
	}

	return nil
}

func (block *securityAuthenticationKeyChainBlockKey) read(
	itemTrim string, junSess *junos.Session,
) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "secret "):
		block.Secret, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "secret")
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "start-time "):
		block.StartTime = types.StringValue(strings.Split(strings.Trim(itemTrim, "\""), " ")[0])
	case balt.CutPrefixInString(&itemTrim, "algorithm "):
		block.Algorithm = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ao-attribute cryptographic-algorithm "):
		block.AOCryptographicAlgorithm = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "ao-attribute recv-id "):
		block.AORecvID, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "ao-attribute send-id "):
		block.AOSendID, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "ao-attribute tcp-ao-option "):
		block.AOTcpAOOption = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "key-name "):
		block.KeyName = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "options "):
		block.Options = types.StringValue(itemTrim)
	}

	return nil
}

func (rscData *securityAuthenticationKeyChainData) del(
	ctx context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security authentication-key-chains key-chain \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(ctx, configSet)
}
