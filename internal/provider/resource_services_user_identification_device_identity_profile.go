package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                   = &servicesUserIdentificationDeviceIdentityProfile{}
	_ resource.ResourceWithConfigure      = &servicesUserIdentificationDeviceIdentityProfile{}
	_ resource.ResourceWithValidateConfig = &servicesUserIdentificationDeviceIdentityProfile{}
	_ resource.ResourceWithImportState    = &servicesUserIdentificationDeviceIdentityProfile{}
)

type servicesUserIdentificationDeviceIdentityProfile struct {
	client *junos.Client
}

func newServicesUserIdentificationDeviceIdentityProfileResource() resource.Resource {
	return &servicesUserIdentificationDeviceIdentityProfile{}
}

func (rsc *servicesUserIdentificationDeviceIdentityProfile) typeName() string {
	return providerName + "_services_user_identification_device_identity_profile"
}

func (rsc *servicesUserIdentificationDeviceIdentityProfile) junosName() string {
	return "services user-identification device-information end-user-profile profile-name"
}

func (rsc *servicesUserIdentificationDeviceIdentityProfile) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *servicesUserIdentificationDeviceIdentityProfile) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *servicesUserIdentificationDeviceIdentityProfile) Configure(
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

func (rsc *servicesUserIdentificationDeviceIdentityProfile) Schema(
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
				Description: "End-user-profile profile-name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
				},
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "Domain name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"attribute": schema.ListNestedBlock{
				Description: "For each name of attribute, configure list of values.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Attribute name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 64),
								tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
							},
						},
						"value": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "A list of values.",
							Validators: []validator.Set{
								setvalidator.SizeBetween(1, 20),
								setvalidator.NoNullValues(),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
								),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.IsRequired(),
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

type servicesUserIdentificationDeviceIdentityProfileData struct {
	ID        types.String                                                    `tfsdk:"id"`
	Name      types.String                                                    `tfsdk:"name"`
	Domain    types.String                                                    `tfsdk:"domain"`
	Attribute []servicesUserIdentificationDeviceIdentityProfileBlockAttribute `tfsdk:"attribute"`
}

type servicesUserIdentificationDeviceIdentityProfileConfig struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Domain    types.String `tfsdk:"domain"`
	Attribute types.List   `tfsdk:"attribute"`
}

type servicesUserIdentificationDeviceIdentityProfileBlockAttribute struct {
	Name  types.String   `tfsdk:"name"  tfdata:"identifier"`
	Value []types.String `tfsdk:"value"`
}

type servicesUserIdentificationDeviceIdentityProfileBlockAttributeConfig struct {
	Name  types.String `tfsdk:"name"`
	Value types.Set    `tfsdk:"value"`
}

func (rsc *servicesUserIdentificationDeviceIdentityProfile) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config servicesUserIdentificationDeviceIdentityProfileConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Attribute.IsNull() &&
		!config.Attribute.IsUnknown() {
		var configAttribute []servicesUserIdentificationDeviceIdentityProfileBlockAttributeConfig
		asDiags := config.Attribute.ElementsAs(ctx, &configAttribute, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		attributeNamme := make(map[string]struct{})
		for i, block := range configAttribute {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := attributeNamme[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("attribute").AtListIndex(i).AtName("name"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple attribute blocks with the same name %q", name),
				)
			}
			attributeNamme[name] = struct{}{}
		}
	}
}

func (rsc *servicesUserIdentificationDeviceIdentityProfile) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan servicesUserIdentificationDeviceIdentityProfileData
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
			profileExists, err := checkServicesUserIdentDeviceIdentityProfileExists(fnCtx, plan.Name.ValueString(), junSess)
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
			profileExists, err := checkServicesUserIdentDeviceIdentityProfileExists(fnCtx, plan.Name.ValueString(), junSess)
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

func (rsc *servicesUserIdentificationDeviceIdentityProfile) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data servicesUserIdentificationDeviceIdentityProfileData
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

func (rsc *servicesUserIdentificationDeviceIdentityProfile) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state servicesUserIdentificationDeviceIdentityProfileData
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

func (rsc *servicesUserIdentificationDeviceIdentityProfile) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state servicesUserIdentificationDeviceIdentityProfileData
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

func (rsc *servicesUserIdentificationDeviceIdentityProfile) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data servicesUserIdentificationDeviceIdentityProfileData

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

func checkServicesUserIdentDeviceIdentityProfileExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services user-identification device-information end-user-profile profile-name " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *servicesUserIdentificationDeviceIdentityProfileData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *servicesUserIdentificationDeviceIdentityProfileData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *servicesUserIdentificationDeviceIdentityProfileData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set services user-identification device-information end-user-profile profile-name " +
		rscData.Name.ValueString() + " "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix + "domain-name " + rscData.Domain.ValueString()

	attributeNamme := make(map[string]struct{})
	for i, block := range rscData.Attribute {
		name := block.Name.ValueString()
		if _, ok := attributeNamme[name]; ok {
			return path.Root("attribute").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple attribute blocks with the same name %q", name)
		}
		attributeNamme[name] = struct{}{}

		for _, v := range block.Value {
			configSet = append(configSet, setPrefix+"attribute "+name+
				" string \""+v.ValueString()+"\"")
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *servicesUserIdentificationDeviceIdentityProfileData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services user-identification device-information end-user-profile profile-name " + name +
		junos.PipeDisplaySetRelative,
	)
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
			case balt.CutPrefixInString(&itemTrim, "domain-name "):
				rscData.Domain = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "attribute "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.Attribute = tfdata.AppendPotentialNewBlock(rscData.Attribute, types.StringValue(name))
				attribute := &rscData.Attribute[len(rscData.Attribute)-1]
				balt.CutPrefixInString(&itemTrim, name+" ")

				if balt.CutPrefixInString(&itemTrim, "string ") {
					attribute.Value = append(attribute.Value, types.StringValue(strings.Trim(itemTrim, "\"")))
				}
			}
		}
	}

	return nil
}

func (rscData *servicesUserIdentificationDeviceIdentityProfileData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete services user-identification device-information end-user-profile profile-name " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
