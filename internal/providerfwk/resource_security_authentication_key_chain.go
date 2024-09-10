package providerfwk

import (
	"context"
	"fmt"
	"regexp"
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
							Required:    true,
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
			},
		},
	}
}

type securityAuthenticationKeyChainData struct {
	ID          types.String                             `tfsdk:"id"`
	Name        types.String                             `tfsdk:"name"`
	Description types.String                             `tfsdk:"description"`
	Tolerance   types.Int64                              `tfsdk:"tolerance"`
	Key         []securityAuthenticationKeyChainBlockKey `tfsdk:"key"`
}

type securityAuthenticationKeyChainConfig struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Tolerance   types.Int64  `tfsdk:"tolerance"`
	Key         types.Set    `tfsdk:"key"`
}

type securityAuthenticationKeyChainBlockKey struct {
	ID                       types.Int64  `tfsdk:"id"`
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

	if config.Key.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			tfdiag.MissingConfigErrSummary,
			"key block must be specified",
		)
	} else if !config.Key.IsUnknown() {
		var configKey []securityAuthenticationKeyChainBlockKey
		asDiags := config.Key.ElementsAs(ctx, &configKey, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		keyID := make(map[int64]struct{})
		for _, block := range configKey {
			if block.ID.IsUnknown() {
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
		}
	}
}

func (rsc *securityAuthenticationKeyChain) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan securityAuthenticationKeyChainData
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
		nil,
		resp,
	)
}

func (rsc *securityAuthenticationKeyChain) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state securityAuthenticationKeyChainData
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
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security authentication-key-chains key-chain \"" + name + "\"" + junos.PipeDisplaySet)
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

func (rscData *securityAuthenticationKeyChainData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set security authentication-key-chains key-chain \"" + rscData.Name.ValueString() + "\" "

	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if !rscData.Tolerance.IsNull() {
		configSet = append(configSet, setPrefix+"tolerance "+
			utils.ConvI64toa(rscData.Tolerance.ValueInt64()))
	}
	keyID := make(map[int64]struct{})
	for _, block := range rscData.Key {
		id := block.ID.ValueInt64()
		if _, ok := keyID[id]; ok {
			return path.Root("key"),
				fmt.Errorf("multiple key blocks with the same id %d", id)
		}
		keyID[id] = struct{}{}

		configSet = append(configSet, block.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *securityAuthenticationKeyChainBlockKey) configSet(setPrefix string) []string {
	setPrefix += "key " + utils.ConvI64toa(block.ID.ValueInt64()) + " "

	configSet := []string{
		setPrefix + "secret \"" + block.Secret.ValueString() + "\"",
		setPrefix + "start-time " + block.StartTime.ValueString(),
	}

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

	return configSet
}

func (rscData *securityAuthenticationKeyChainData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security authentication-key-chains key-chain \"" + name + "\"" + junos.PipeDisplaySetRelative)
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
				rscData.Key, key = tfdata.ExtractBlockWithTFTypesInt64(
					rscData.Key, "ID", keyID.ValueInt64(),
				)
				key.ID = keyID
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
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete security authentication-key-chains key-chain \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
