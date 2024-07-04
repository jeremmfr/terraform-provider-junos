package providerfwk

import (
	"context"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
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
	_ resource.Resource                = &snmpV3Community{}
	_ resource.ResourceWithConfigure   = &snmpV3Community{}
	_ resource.ResourceWithImportState = &snmpV3Community{}
)

type snmpV3Community struct {
	client *junos.Client
}

func newSnmpV3CommunityResource() resource.Resource {
	return &snmpV3Community{}
}

func (rsc *snmpV3Community) typeName() string {
	return providerName + "_snmp_v3_community"
}

func (rsc *snmpV3Community) junosName() string {
	return "snmp v3 snmp-community"
}

func (rsc *snmpV3Community) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *snmpV3Community) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *snmpV3Community) Configure(
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

func (rsc *snmpV3Community) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with format `<community_index>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"community_index": schema.StringAttribute{
				Required:    true,
				Description: "Unique index value in this community table entry.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"security_name": schema.StringAttribute{
				Required:    true,
				Description: "Security name used when performing access control.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"community_name": schema.StringAttribute{
				Optional:    true,
				Description: "SNMPv1/v2c community name (default is same as community-index).",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"context": schema.StringAttribute{
				Optional:    true,
				Description: "Context used when performing access control.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"tag": schema.StringAttribute{
				Optional:    true,
				Description: "Tag identifier for set of targets allowed to use this community string.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
	}
}

type snmpV3CommunityData struct {
	ID             types.String `tfsdk:"id"`
	CommunityIndex types.String `tfsdk:"community_index"`
	SecurityName   types.String `tfsdk:"security_name"`
	CommunityName  types.String `tfsdk:"community_name"`
	Context        types.String `tfsdk:"context"`
	Tag            types.String `tfsdk:"tag"`
}

func (rsc *snmpV3Community) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan snmpV3CommunityData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.CommunityIndex.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("community_index"),
			"Empty Community Index",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "community_index"),
		)

		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		func(fnCtx context.Context, junSess *junos.Session) bool {
			communityExists, err := checkSnmpV3CommunityExists(
				fnCtx,
				plan.CommunityIndex.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if communityExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.CommunityIndex),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			communityExists, err := checkSnmpV3CommunityExists(
				fnCtx,
				plan.CommunityIndex.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !communityExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					defaultResourceDoesNotExistsAfterCommitMessage(rsc, plan.CommunityIndex),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *snmpV3Community) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data snmpV3CommunityData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom1String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.CommunityIndex.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *snmpV3Community) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state snmpV3CommunityData
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

func (rsc *snmpV3Community) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state snmpV3CommunityData
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

func (rsc *snmpV3Community) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data snmpV3CommunityData

	var _ resourceDataReadFrom1String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindIDStrMessage(rsc, req.ID, "community_index"),
	)
}

func checkSnmpV3CommunityExists(
	_ context.Context, communityIndex string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp v3 snmp-community \"" + communityIndex + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *snmpV3CommunityData) fillID() {
	rscData.ID = types.StringValue(rscData.CommunityIndex.ValueString())
}

func (rscData *snmpV3CommunityData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *snmpV3CommunityData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set snmp v3 snmp-community \"" + rscData.CommunityIndex.ValueString() + "\" "

	configSet := []string{
		setPrefix + "security-name \"" + rscData.SecurityName.ValueString() + "\"",
	}

	if v := rscData.CommunityName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"community-name \""+v+"\"")
	}
	if v := rscData.Context.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"context \""+v+"\"")
	}
	if v := rscData.Tag.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"tag \""+v+"\"")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *snmpV3CommunityData) read(
	_ context.Context, communityIndex string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp v3 snmp-community \"" + communityIndex + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.CommunityIndex = types.StringValue(communityIndex)
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
			case balt.CutPrefixInString(&itemTrim, "security-name "):
				rscData.SecurityName = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "community-name "):
				rscData.CommunityName, err = tfdata.JunosDecode(strings.Trim(itemTrim, "\""), "community-name")
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "context "):
				rscData.Context = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "tag "):
				rscData.Tag = types.StringValue(strings.Trim(itemTrim, "\""))
			}
		}
	}

	return nil
}

func (rscData *snmpV3CommunityData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete snmp v3 snmp-community \"" + rscData.CommunityIndex.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
