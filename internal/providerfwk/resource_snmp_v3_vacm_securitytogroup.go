package providerfwk

import (
	"context"
	"errors"
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
	balt "github.com/jeremmfr/go-utils/basicalter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &snmpV3VacmSecuritytogroup{}
	_ resource.ResourceWithConfigure   = &snmpV3VacmSecuritytogroup{}
	_ resource.ResourceWithImportState = &snmpV3VacmSecuritytogroup{}
)

type snmpV3VacmSecuritytogroup struct {
	client *junos.Client
}

func newSnmpV3VacmSecuritytogroupResource() resource.Resource {
	return &snmpV3VacmSecuritytogroup{}
}

func (rsc *snmpV3VacmSecuritytogroup) typeName() string {
	return providerName + "_snmp_v3_vacm_securitytogroup"
}

func (rsc *snmpV3VacmSecuritytogroup) junosName() string {
	return "snmp v3 vacm security-to-group"
}

func (rsc *snmpV3VacmSecuritytogroup) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *snmpV3VacmSecuritytogroup) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *snmpV3VacmSecuritytogroup) Configure(
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

func (rsc *snmpV3VacmSecuritytogroup) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: defaultResourceSchemaDescription(rsc),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "An identifier for the resource with format " +
					"`<model>" + junos.IDSeparator + "<name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"model": schema.StringAttribute{
				Required:    true,
				Description: "Security model context for group assignment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("usm", "v1", "v2c"),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Security name to assign to group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"group": schema.StringAttribute{
				Required:    true,
				Description: "Group to which to assign security name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
	}
}

type snmpV3VacmSecuritytogroupData struct {
	ID    types.String `tfsdk:"id"`
	Model types.String `tfsdk:"model"`
	Name  types.String `tfsdk:"name"`
	Group types.String `tfsdk:"group"`
}

func (rsc *snmpV3VacmSecuritytogroup) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan snmpV3VacmSecuritytogroupData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Model.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("model"),
			"Empty Model",
			defaultResourceCouldNotCreateWithEmptyMessage(rsc, "model"),
		)

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
			securityToGroupExists, err := checkSnmpV3VacmSecuritytogroupExists(
				fnCtx,
				plan.Model.ValueString(),
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if securityToGroupExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf(
						rsc.junosName()+" security-model %q security-name %q already exists",
						plan.Model.ValueString(), plan.Name.ValueString(),
					),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			securityToGroupExists, err := checkSnmpV3VacmSecuritytogroupExists(
				fnCtx,
				plan.Model.ValueString(),
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !securityToGroupExists {
				resp.Diagnostics.AddError(
					tfdiag.NotFoundErrSummary,
					fmt.Sprintf(
						rsc.junosName()+" security-model %q security-name %q does not exists after commit "+
							"=> check your config",
						plan.Model.ValueString(), plan.Name.ValueString(),
					),
				)

				return false
			}

			return true
		},
		&plan,
		resp,
	)
}

func (rsc *snmpV3VacmSecuritytogroup) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data snmpV3VacmSecuritytogroupData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ resourceDataReadFrom2String = &data
	defaultResourceRead(
		ctx,
		rsc,
		[]any{
			state.Model.ValueString(),
			state.Name.ValueString(),
		},
		&data,
		nil,
		resp,
	)
}

func (rsc *snmpV3VacmSecuritytogroup) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state snmpV3VacmSecuritytogroupData
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

func (rsc *snmpV3VacmSecuritytogroup) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state snmpV3VacmSecuritytogroupData
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

func (rsc *snmpV3VacmSecuritytogroup) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data snmpV3VacmSecuritytogroupData

	var _ resourceDataReadFrom2String = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		defaultResourceImportDontFindMessage(rsc, req.ID)+
			" (id must be <model>"+junos.IDSeparator+"<name>)",
	)
}

func checkSnmpV3VacmSecuritytogroupExists(
	_ context.Context, model, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp v3 vacm security-to-group security-model " + model + " security-name \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *snmpV3VacmSecuritytogroupData) fillID() {
	rscData.ID = types.StringValue(rscData.Model.ValueString() + junos.IDSeparator + rscData.Name.ValueString())
}

func (rscData *snmpV3VacmSecuritytogroupData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *snmpV3VacmSecuritytogroupData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 1)
	if v := rscData.Group.ValueString(); v != "" {
		configSet[0] = "set snmp v3 vacm security-to-group" +
			" security-model " + rscData.Model.ValueString() +
			" security-name \"" + rscData.Name.ValueString() + "\"" +
			" group \"" + v + "\""
	} else {
		return path.Root("group"), errors.New("group must be specified")
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *snmpV3VacmSecuritytogroupData) read(
	_ context.Context, model, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"snmp v3 vacm security-to-group security-model " + model + " security-name \"" + name + "\"" +
		junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	if showConfig != junos.EmptyW {
		rscData.Model = types.StringValue(model)
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
			if balt.CutPrefixInString(&itemTrim, "group ") {
				rscData.Group = types.StringValue(strings.Trim(itemTrim, "\""))
			}
		}
	}

	return nil
}

func (rscData *snmpV3VacmSecuritytogroupData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete snmp v3 vacm security-to-group " +
			"security-model " + rscData.Model.ValueString() + " security-name \"" + rscData.Name.ValueString() + "\"",
	}

	return junSess.ConfigSet(configSet)
}
