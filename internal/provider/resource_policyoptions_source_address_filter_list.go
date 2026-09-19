package provider

import (
	"context"
	"fmt"
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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &policyoptionsSourceAddressFilterList{}
	_ resource.ResourceWithConfigure      = &policyoptionsSourceAddressFilterList{}
	_ resource.ResourceWithValidateConfig = &policyoptionsSourceAddressFilterList{}
	_ resource.ResourceWithImportState    = &policyoptionsSourceAddressFilterList{}
)

type policyoptionsSourceAddressFilterList struct {
	client *junos.Client
}

func newPolicyoptionsSourceAddressFilterListResource() resource.Resource {
	return &policyoptionsSourceAddressFilterList{}
}

func (rsc *policyoptionsSourceAddressFilterList) typeName() string {
	return providerName + "_policyoptions_source_address_filter_list"
}

func (rsc *policyoptionsSourceAddressFilterList) junosName() string {
	return "policy-options source-address-filter-list"
}

func (rsc *policyoptionsSourceAddressFilterList) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *policyoptionsSourceAddressFilterList) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *policyoptionsSourceAddressFilterList) Configure(
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

func (rsc *policyoptionsSourceAddressFilterList) Schema(
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
				Description: "Source address filter list name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"dynamic_db": schema.BoolAttribute{
				Optional:    true,
				Description: "Object may exist in dynamic database.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"address": schema.SetNestedBlock{
				Description: "List of source addresses.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Required:    true,
							Description: "IP address.",
							Validators: []validator.String{
								tfvalidator.StringCIDRNetwork(),
							},
						},
						"option": schema.StringAttribute{
							Required:    true,
							Description: "Mask option.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"exact", "longer", "orlonger", "prefix-length-range", "through", "upto",
								),
							},
						},
						"option_value": schema.StringAttribute{
							Optional:    true,
							Description: "For options that need an argument.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
		},
	}
}

type policyoptionsSourceAddressFilterListData struct {
	ID        types.String                                       `tfsdk:"id"`
	Name      types.String                                       `tfsdk:"name"`
	DynamicDB types.Bool                                         `tfsdk:"dynamic_db"`
	Address   []policyoptionsSourceAddressFilterListBlockAddress `tfsdk:"address"`
}

type policyoptionsSourceAddressFilterListConfig struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	DynamicDB types.Bool   `tfsdk:"dynamic_db"`
	Address   types.Set    `tfsdk:"address"`
}

type policyoptionsSourceAddressFilterListBlockAddress struct {
	Address     types.String `tfsdk:"address"`
	Option      types.String `tfsdk:"option"`
	OptionValue types.String `tfsdk:"option_value"`
}

func (rsc *policyoptionsSourceAddressFilterList) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config policyoptionsSourceAddressFilterListConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Address.IsNull() && !config.Address.IsUnknown() &&
		!config.DynamicDB.IsNull() && !config.DynamicDB.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			tfdiag.ConflictConfigErrSummary,
			"only one of address or dynamic_db must be specified",
		)
	}
}

func (rsc *policyoptionsSourceAddressFilterList) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan policyoptionsSourceAddressFilterListData
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
			listExists, err := checkPolicyoptionsSourceAddressFilterListExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if listExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			listExists, err := checkPolicyoptionsSourceAddressFilterListExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !listExists {
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

func (rsc *policyoptionsSourceAddressFilterList) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data policyoptionsSourceAddressFilterListData
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

func (rsc *policyoptionsSourceAddressFilterList) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state policyoptionsSourceAddressFilterListData
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

func (rsc *policyoptionsSourceAddressFilterList) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state policyoptionsSourceAddressFilterListData
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

func (rsc *policyoptionsSourceAddressFilterList) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data policyoptionsSourceAddressFilterListData

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

func checkPolicyoptionsSourceAddressFilterListExists(
	ctx context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"policy-options source-address-filter-list \""+name+"\""+junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *policyoptionsSourceAddressFilterListData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *policyoptionsSourceAddressFilterListData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *policyoptionsSourceAddressFilterListData) set(
	ctx context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set policy-options source-address-filter-list \"" + rscData.Name.ValueString() + "\" "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	for _, block := range rscData.Address {
		setAddress := setPrefix + block.Address.ValueString() + " " + block.Option.ValueString()
		if v := block.OptionValue.ValueString(); v != "" {
			setAddress += " " + v
		}
		configSet = append(configSet, setAddress)
	}
	if rscData.DynamicDB.ValueBool() {
		configSet = append(configSet, setPrefix+"dynamic-db")
	}

	return path.Empty(), junSess.ConfigSet(ctx, configSet)
}

func (rscData *policyoptionsSourceAddressFilterListData) read(
	ctx context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(ctx, junos.CmdShowConfig+
		"policy-options source-address-filter-list \""+name+"\""+junos.PipeDisplaySetRelative)
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
			case itemTrim == "dynamic-db":
				rscData.DynamicDB = types.BoolValue(true)
			case strings.Contains(itemTrim, "/"):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 2 { // <address> <option> <option_value>?
					return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "source-address-filter-list", itemTrim)
				}
				address := policyoptionsSourceAddressFilterListBlockAddress{
					Address: types.StringValue(itemTrimFields[0]),
					Option:  types.StringValue(itemTrimFields[1]),
				}
				if len(itemTrimFields) > 2 {
					address.OptionValue = types.StringValue(itemTrimFields[2])
				}
				rscData.Address = append(rscData.Address, address)
			}
		}
	}

	return nil
}

func (rscData *policyoptionsSourceAddressFilterListData) del(
	ctx context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete policy-options source-address-filter-list \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(ctx, configSet)
}
