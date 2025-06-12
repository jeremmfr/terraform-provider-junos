package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                   = &servicesUserIdentificationADAccessDomain{}
	_ resource.ResourceWithConfigure      = &servicesUserIdentificationADAccessDomain{}
	_ resource.ResourceWithValidateConfig = &servicesUserIdentificationADAccessDomain{}
	_ resource.ResourceWithImportState    = &servicesUserIdentificationADAccessDomain{}
	_ resource.ResourceWithUpgradeState   = &servicesUserIdentificationADAccessDomain{}
)

type servicesUserIdentificationADAccessDomain struct {
	client *junos.Client
}

func newServicesUserIdentificationADAccessDomainResource() resource.Resource {
	return &servicesUserIdentificationADAccessDomain{}
}

func (rsc *servicesUserIdentificationADAccessDomain) typeName() string {
	return providerName + "_services_user_identification_ad_access_domain"
}

func (rsc *servicesUserIdentificationADAccessDomain) junosName() string {
	return "services user-identification active-directory-access domain"
}

func (rsc *servicesUserIdentificationADAccessDomain) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *servicesUserIdentificationADAccessDomain) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *servicesUserIdentificationADAccessDomain) Configure(
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

func (rsc *servicesUserIdentificationADAccessDomain) Schema(
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
				Description: "Domain name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
				},
			},
			"user_name": schema.StringAttribute{
				Required:    true,
				Description: "User name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
				},
			},
			"user_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Password string.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"domain_controller": schema.ListNestedBlock{
				Description: "For each name of domain-controller, configure address.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Domain controller name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 64),
								tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
							},
						},
						"address": schema.StringAttribute{
							Required:    true,
							Description: "Address of domain controller.",
							Validators: []validator.String{
								tfvalidator.StringIPAddress(),
							},
						},
					},
				},
			},
			"ip_user_mapping_discovery_wmi": schema.SingleNestedBlock{
				Description: "Enable `ip-user-mapping discovery-method wmi`.",
				Attributes: map[string]schema.Attribute{
					"event_log_scanning_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "Interval of event log scanning (seconds).",
						Validators: []validator.Int64{
							int64validator.Between(5, 60),
						},
					},
					"initial_event_log_timespan": schema.Int64Attribute{
						Optional:    true,
						Description: "Event log scanning timespan (hours).",
						Validators: []validator.Int64{
							int64validator.Between(1, 168),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"user_group_mapping_ldap": schema.SingleNestedBlock{
				Description: "User-group-mapping LDAP.",
				Attributes: map[string]schema.Attribute{
					"base": schema.StringAttribute{
						Required:    false, // true when SingleNestedBlock is specified
						Optional:    true,
						Description: "Base distinguished name.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 128),
							tfvalidator.StringDoubleQuoteExclusion(),
						},
					},
					"address": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Address of LDAP server.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.NoNullValues(),
							listvalidator.ValueStringsAre(
								tfvalidator.StringIPAddress(),
							),
						},
					},
					"auth_algo_simple": schema.BoolAttribute{
						Optional:    true,
						Description: "Authentication-algorithm simple.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"ssl": schema.BoolAttribute{
						Optional:    true,
						Description: "SSL.",
						Validators: []validator.Bool{
							tfvalidator.BoolTrue(),
						},
					},
					"user_name": schema.StringAttribute{
						Optional:    true,
						Description: "User name.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 64),
							tfvalidator.StringFormat(tfvalidator.DNSNameFormat),
						},
					},
					"user_password": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Password string.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 128),
							tfvalidator.StringDoubleQuoteExclusion(),
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

//nolint:lll
type servicesUserIdentificationADAccessDomainData struct {
	ID                        types.String                                                            `tfsdk:"id"`
	Name                      types.String                                                            `tfsdk:"name"`
	UserName                  types.String                                                            `tfsdk:"user_name"`
	UserPassword              types.String                                                            `tfsdk:"user_password"`
	DomainController          []servicesUserIdentificationADAccessDomainBlockDomainController         `tfsdk:"domain_controller"`
	IPUserMappingDiscoveryWmi *servicesUserIdentificationADAccessDomainBlockIPUserMappingDiscoveryWmi `tfsdk:"ip_user_mapping_discovery_wmi"`
	UserGroupMappingLdap      *servicesUserIdentificationADAccessDomainBlockUserGroupMappingLdap      `tfsdk:"user_group_mapping_ldap"`
}

//nolint:lll
type servicesUserIdentificationADAccessDomainConfig struct {
	ID                        types.String                                                             `tfsdk:"id"`
	Name                      types.String                                                             `tfsdk:"name"`
	UserName                  types.String                                                             `tfsdk:"user_name"`
	UserPassword              types.String                                                             `tfsdk:"user_password"`
	DomainController          types.List                                                               `tfsdk:"domain_controller"`
	IPUserMappingDiscoveryWmi *servicesUserIdentificationADAccessDomainBlockIPUserMappingDiscoveryWmi  `tfsdk:"ip_user_mapping_discovery_wmi"`
	UserGroupMappingLdap      *servicesUserIdentificationADAccessDomainBlockUserGroupMappingLdapConfig `tfsdk:"user_group_mapping_ldap"`
}

type servicesUserIdentificationADAccessDomainBlockDomainController struct {
	Name    types.String `tfsdk:"name"    tfdata:"identifier"`
	Address types.String `tfsdk:"address"`
}

type servicesUserIdentificationADAccessDomainBlockIPUserMappingDiscoveryWmi struct {
	EventLogScanningInterval types.Int64 `tfsdk:"event_log_scanning_interval"`
	InitialEventLogTimespan  types.Int64 `tfsdk:"initial_event_log_timespan"`
}

type servicesUserIdentificationADAccessDomainBlockUserGroupMappingLdap struct {
	Base           types.String   `tfsdk:"base"`
	Address        []types.String `tfsdk:"address"`
	AuthAlgoSimple types.Bool     `tfsdk:"auth_algo_simple"`
	Ssl            types.Bool     `tfsdk:"ssl"`
	UserName       types.String   `tfsdk:"user_name"`
	UserPassword   types.String   `tfsdk:"user_password"`
}

type servicesUserIdentificationADAccessDomainBlockUserGroupMappingLdapConfig struct {
	Base           types.String `tfsdk:"base"`
	Address        types.List   `tfsdk:"address"`
	AuthAlgoSimple types.Bool   `tfsdk:"auth_algo_simple"`
	Ssl            types.Bool   `tfsdk:"ssl"`
	UserName       types.String `tfsdk:"user_name"`
	UserPassword   types.String `tfsdk:"user_password"`
}

func (rsc *servicesUserIdentificationADAccessDomain) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config servicesUserIdentificationADAccessDomainConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.DomainController.IsNull() &&
		!config.DomainController.IsUnknown() {
		var configDomainController []servicesUserIdentificationADAccessDomainBlockDomainController
		asDiags := config.DomainController.ElementsAs(ctx, &configDomainController, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}
		domainControllerNamme := make(map[string]struct{})
		for i, block := range configDomainController {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := domainControllerNamme[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("domain_controller").AtListIndex(i).AtName("name"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple domain_controller blocks with the same name %q", name),
				)
			}
			domainControllerNamme[name] = struct{}{}
		}
	}
	if config.UserGroupMappingLdap != nil {
		if config.UserGroupMappingLdap.Base.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("user_group_mapping_ldap").AtName("base"),
				tfdiag.MissingConfigErrSummary,
				"base must be specified in user_group_mapping_ldap block",
			)
		}
	}
}

func (rsc *servicesUserIdentificationADAccessDomain) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan servicesUserIdentificationADAccessDomainData
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
			domainExists, err := checkServicesUserIdentADAccessDomainExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if domainExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			domainExists, err := checkServicesUserIdentADAccessDomainExists(fnCtx, plan.Name.ValueString(), junSess)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !domainExists {
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

func (rsc *servicesUserIdentificationADAccessDomain) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data servicesUserIdentificationADAccessDomainData
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

func (rsc *servicesUserIdentificationADAccessDomain) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state servicesUserIdentificationADAccessDomainData
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

func (rsc *servicesUserIdentificationADAccessDomain) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state servicesUserIdentificationADAccessDomainData
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

func (rsc *servicesUserIdentificationADAccessDomain) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data servicesUserIdentificationADAccessDomainData

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

func checkServicesUserIdentADAccessDomainExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services user-identification active-directory-access domain " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *servicesUserIdentificationADAccessDomainData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *servicesUserIdentificationADAccessDomainData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *servicesUserIdentificationADAccessDomainData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	setPrefix := "set services user-identification active-directory-access domain " + rscData.Name.ValueString() + " "

	configSet := make([]string, 2, 100)
	configSet[0] = setPrefix + "user " + rscData.UserName.ValueString()
	configSet[1] = setPrefix + "user password \"" + rscData.UserPassword.ValueString() + "\""

	domainControllerName := make(map[string]struct{})
	for i, block := range rscData.DomainController {
		name := block.Name.ValueString()
		if _, ok := domainControllerName[name]; ok {
			return path.Root("domain_controller").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple domain_controller blocks with the same name %q", name)
		}
		domainControllerName[name] = struct{}{}

		configSet = append(configSet, setPrefix+"domain-controller "+name+
			" address "+block.Address.ValueString())
	}
	if rscData.IPUserMappingDiscoveryWmi != nil {
		configSet = append(configSet, rscData.IPUserMappingDiscoveryWmi.configSet(setPrefix)...)
	}
	if rscData.UserGroupMappingLdap != nil {
		configSet = append(configSet, rscData.UserGroupMappingLdap.configSet(setPrefix)...)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *servicesUserIdentificationADAccessDomainBlockIPUserMappingDiscoveryWmi) configSet(
	setPrefix string,
) []string {
	setPrefix += "ip-user-mapping discovery-method wmi "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix

	if !block.EventLogScanningInterval.IsNull() {
		configSet = append(configSet, setPrefix+"event-log-scanning-interval "+
			utils.ConvI64toa(block.EventLogScanningInterval.ValueInt64()))
	}
	if !block.InitialEventLogTimespan.IsNull() {
		configSet = append(configSet, setPrefix+"initial-event-log-timespan "+
			utils.ConvI64toa(block.InitialEventLogTimespan.ValueInt64()))
	}

	return configSet
}

func (block *servicesUserIdentificationADAccessDomainBlockUserGroupMappingLdap) configSet(setPrefix string) []string {
	setPrefix += "user-group-mapping ldap "

	configSet := make([]string, 1, 100)
	configSet[0] = setPrefix + "base \"" + block.Base.ValueString() + "\""

	for _, v := range block.Address {
		configSet = append(configSet, setPrefix+"address "+v.ValueString())
	}
	if block.AuthAlgoSimple.ValueBool() {
		configSet = append(configSet, setPrefix+"authentication-algorithm simple")
	}
	if block.Ssl.ValueBool() {
		configSet = append(configSet, setPrefix+"ssl")
	}
	if v := block.UserName.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"user "+v)
	}
	if v := block.UserPassword.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"user password \""+v+"\"")
	}

	return configSet
}

func (rscData *servicesUserIdentificationADAccessDomainData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"services user-identification active-directory-access domain " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "user password "):
				rscData.UserPassword, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "user password")
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "user "):
				rscData.UserName = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "domain-controller "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				rscData.DomainController = tfdata.AppendPotentialNewBlock(rscData.DomainController, types.StringValue(name))
				domainController := &rscData.DomainController[len(rscData.DomainController)-1]
				balt.CutPrefixInString(&itemTrim, name+" ")

				if balt.CutPrefixInString(&itemTrim, "address ") {
					domainController.Address = types.StringValue(itemTrim)
				}
			case balt.CutPrefixInString(&itemTrim, "ip-user-mapping discovery-method wmi"):
				if rscData.IPUserMappingDiscoveryWmi == nil {
					rscData.IPUserMappingDiscoveryWmi = &servicesUserIdentificationADAccessDomainBlockIPUserMappingDiscoveryWmi{}
				}

				if balt.CutPrefixInString(&itemTrim, " ") {
					if err := rscData.IPUserMappingDiscoveryWmi.read(itemTrim); err != nil {
						return err
					}
				}
			case balt.CutPrefixInString(&itemTrim, "user-group-mapping ldap "):
				if rscData.UserGroupMappingLdap == nil {
					rscData.UserGroupMappingLdap = &servicesUserIdentificationADAccessDomainBlockUserGroupMappingLdap{}
				}

				if err := rscData.UserGroupMappingLdap.read(itemTrim, junSess); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (block *servicesUserIdentificationADAccessDomainBlockIPUserMappingDiscoveryWmi) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "event-log-scanning-interval "):
		block.EventLogScanningInterval, err = tfdata.ConvAtoi64Value(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "initial-event-log-timespan "):
		block.InitialEventLogTimespan, err = tfdata.ConvAtoi64Value(itemTrim)
	}

	return err
}

func (block *servicesUserIdentificationADAccessDomainBlockUserGroupMappingLdap) read(
	itemTrim string, junSess *junos.Session,
) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "base "):
		block.Base = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "address "):
		block.Address = append(block.Address, types.StringValue(itemTrim))
	case itemTrim == "authentication-algorithm simple":
		block.AuthAlgoSimple = types.BoolValue(true)
	case itemTrim == "ssl":
		block.Ssl = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "user password "):
		block.UserPassword, err = junSess.JunosDecode(strings.Trim(itemTrim, "\""), "user password")
	case balt.CutPrefixInString(&itemTrim, "user "):
		block.UserName = types.StringValue(itemTrim)
	}

	return err
}

func (rscData *servicesUserIdentificationADAccessDomainData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete services user-identification active-directory-access domain " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
