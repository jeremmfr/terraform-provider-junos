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
	_ resource.Resource                   = &application{}
	_ resource.ResourceWithConfigure      = &application{}
	_ resource.ResourceWithValidateConfig = &application{}
	_ resource.ResourceWithImportState    = &application{}
)

type application struct {
	client *junos.Client
}

func newApplicationResource() resource.Resource {
	return &application{}
}

func (rsc *application) typeName() string {
	return providerName + "_application"
}

func (rsc *application) junosName() string {
	return "applications application"
}

func (rsc *application) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *application) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *application) Configure(
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

func (rsc *application) Schema(
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
				Description: "Application name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"application_protocol": schema.StringAttribute{
				Optional:    true,
				Description: "Application protocol type.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Text description of application.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 900),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"destination_port": schema.StringAttribute{
				Optional:    true,
				Description: "Match TCP/UDP destination port.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"ether_type": schema.StringAttribute{
				Optional:    true,
				Description: "Match ether type.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^0[xX][0-9a-fA-F]{4}$`),
						"must be in hex (example: 0x8906)",
					),
				},
			},
			"inactivity_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Application-specific inactivity timeout.",
				Validators: []validator.Int64{
					int64validator.Between(4, 86400),
				},
			},
			"inactivity_timeout_never": schema.BoolAttribute{
				Optional:    true,
				Description: "Disables inactivity timeout.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Description: "Match IP protocol type.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"rpc_program_number": schema.StringAttribute{
				Optional:    true,
				Description: "Match range of RPC program numbers.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^\d+(-\d+)?$`),
						"must be an integer or a range of integers"),
				},
			},
			"source_port": schema.StringAttribute{
				Optional:    true,
				Description: "Match TCP/UDP source port.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"uuid": schema.StringAttribute{
				Optional:    true,
				Description: "Match universal unique identifier for DCE RPC objects.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
						"must be of the form xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"term": schema.ListNestedBlock{
				Description: "For each name of term to declare.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Term name.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 63),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"protocol": schema.StringAttribute{
							Required:    true,
							Description: "Match IP protocol type.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"alg": schema.StringAttribute{
							Optional:    true,
							Description: "Application Layer Gateway.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"destination_port": schema.StringAttribute{
							Optional:    true,
							Description: "Match TCP/UDP destination port.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"icmp_code": schema.StringAttribute{
							Optional:    true,
							Description: "Match ICMP message code.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"icmp_type": schema.StringAttribute{
							Optional:    true,
							Description: "Match ICMP message type.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"icmp6_code": schema.StringAttribute{
							Optional:    true,
							Description: "Match ICMP6 message code.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"icmp6_type": schema.StringAttribute{
							Optional:    true,
							Description: "Match ICMP6 message type.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormat),
							},
						},
						"inactivity_timeout": schema.Int64Attribute{
							Optional:    true,
							Description: "Application-specific inactivity timeout.",
							Validators: []validator.Int64{
								int64validator.Between(4, 86400),
							},
						},
						"inactivity_timeout_never": schema.BoolAttribute{
							Optional:    true,
							Description: "Disables inactivity timeout.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"rpc_program_number": schema.StringAttribute{
							Optional:    true,
							Description: "Match range of RPC program numbers.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(
									`^\d+(-\d+)?$`),
									"must be an integer or a range of integers"),
							},
						},
						"source_port": schema.StringAttribute{
							Optional:    true,
							Description: "Match TCP/UDP source port.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"uuid": schema.StringAttribute{
							Optional:    true,
							Description: "Match universal unique identifier for DCE RPC objects.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile(
									`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
									"must be of the form xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
								),
							},
						},
					},
				},
			},
		},
	}
}

type applicationData struct {
	ID                     types.String           `tfsdk:"id"`
	Name                   types.String           `tfsdk:"name"`
	ApplicationProtocol    types.String           `tfsdk:"application_protocol"`
	Description            types.String           `tfsdk:"description"`
	DestinationPort        types.String           `tfsdk:"destination_port"`
	EtherType              types.String           `tfsdk:"ether_type"`
	InactivityTimeout      types.Int64            `tfsdk:"inactivity_timeout"`
	InactivityTimeoutNever types.Bool             `tfsdk:"inactivity_timeout_never"`
	Protocol               types.String           `tfsdk:"protocol"`
	RPCProgramNumber       types.String           `tfsdk:"rpc_program_number"`
	SourcePort             types.String           `tfsdk:"source_port"`
	UUID                   types.String           `tfsdk:"uuid"`
	Term                   []applicationBlockTerm `tfsdk:"term"`
}

type applicationConfig struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ApplicationProtocol    types.String `tfsdk:"application_protocol"`
	Description            types.String `tfsdk:"description"`
	DestinationPort        types.String `tfsdk:"destination_port"`
	EtherType              types.String `tfsdk:"ether_type"`
	InactivityTimeout      types.Int64  `tfsdk:"inactivity_timeout"`
	InactivityTimeoutNever types.Bool   `tfsdk:"inactivity_timeout_never"`
	Protocol               types.String `tfsdk:"protocol"`
	RPCProgramNumber       types.String `tfsdk:"rpc_program_number"`
	SourcePort             types.String `tfsdk:"source_port"`
	UUID                   types.String `tfsdk:"uuid"`
	Term                   types.List   `tfsdk:"term"`
}

type applicationBlockTerm struct {
	Name                   types.String `tfsdk:"name"`
	Protocol               types.String `tfsdk:"protocol"`
	Alg                    types.String `tfsdk:"alg"`
	DestinationPort        types.String `tfsdk:"destination_port"`
	IcmpCode               types.String `tfsdk:"icmp_code"`
	IcmpType               types.String `tfsdk:"icmp_type"`
	Icmp6Code              types.String `tfsdk:"icmp6_code"`
	Icmp6Type              types.String `tfsdk:"icmp6_type"`
	InactivityTimeout      types.Int64  `tfsdk:"inactivity_timeout"`
	InactivityTimeoutNever types.Bool   `tfsdk:"inactivity_timeout_never"`
	RPCRrogramNumber       types.String `tfsdk:"rpc_program_number"`
	SourcePort             types.String `tfsdk:"source_port"`
	UUID                   types.String `tfsdk:"uuid"`
}

func (rsc *application) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config applicationConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.InactivityTimeout.IsNull() && !config.InactivityTimeout.IsUnknown() &&
		!config.InactivityTimeoutNever.IsNull() && !config.InactivityTimeoutNever.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("inactivity_timeout"),
			tfdiag.ConflictConfigErrSummary,
			"inactivity_timeout and inactivity_timeout_never cannot be configured together",
		)
	}

	if !config.Term.IsNull() && !config.Term.IsUnknown() {
		if !config.ApplicationProtocol.IsNull() && !config.ApplicationProtocol.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("application_protocol"),
				tfdiag.ConflictConfigErrSummary,
				"application_protocol and term cannot be configured together",
			)
		}
		if !config.DestinationPort.IsNull() && !config.DestinationPort.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("destination_port"),
				tfdiag.ConflictConfigErrSummary,
				"destination_port and term cannot be configured together",
			)
		}
		if !config.InactivityTimeout.IsNull() && !config.InactivityTimeout.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("inactivity_timeout"),
				tfdiag.ConflictConfigErrSummary,
				"inactivity_timeout and term cannot be configured together",
			)
		}
		if !config.InactivityTimeoutNever.IsNull() && !config.InactivityTimeoutNever.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("inactivity_timeout_never"),
				tfdiag.ConflictConfigErrSummary,
				"inactivity_timeout_never and term cannot be configured together",
			)
		}
		if !config.Protocol.IsNull() && !config.Protocol.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("protocol"),
				tfdiag.ConflictConfigErrSummary,
				"protocol and term cannot be configured together",
			)
		}
		if !config.RPCProgramNumber.IsNull() && !config.RPCProgramNumber.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("rpc_program_number"),
				tfdiag.ConflictConfigErrSummary,
				"rpc_program_number and term cannot be configured together",
			)
		}
		if !config.SourcePort.IsNull() && !config.SourcePort.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("source_port"),
				tfdiag.ConflictConfigErrSummary,
				"source_port and term cannot be configured together",
			)
		}
		if !config.UUID.IsNull() && !config.UUID.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("uuid"),
				tfdiag.ConflictConfigErrSummary,
				"uuid and term cannot be configured together",
			)
		}

		var term []applicationBlockTerm
		asDiags := config.Term.ElementsAs(ctx, &term, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		termName := make(map[string]struct{})
		for i, block := range term {
			if !block.Name.IsUnknown() {
				name := block.Name.ValueString()
				if _, ok := termName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("term").AtListIndex(i).AtName("name"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple term blocks with the same name %q", name),
					)
				}
				termName[name] = struct{}{}
			}

			if !block.InactivityTimeout.IsNull() && !block.InactivityTimeout.IsUnknown() &&
				!block.InactivityTimeoutNever.IsNull() && !block.InactivityTimeoutNever.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root("term").AtListIndex(i).AtName("inactivity_timeout"),
					tfdiag.ConflictConfigErrSummary,
					fmt.Sprintf("inactivity_timeout and inactivity_timeout_never cannot be configured together"+
						" in term block %q", block.Name.ValueString()),
				)
			}
		}
	}
}

func (rsc *application) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan applicationData
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
			applicationExists, err := checkApplicationExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PreCheckErrSummary, err.Error())

				return false
			}
			if applicationExists {
				resp.Diagnostics.AddError(
					tfdiag.DuplicateConfigErrSummary,
					defaultResourceAlreadyExistsMessage(rsc, plan.Name),
				)

				return false
			}

			return true
		},
		func(fnCtx context.Context, junSess *junos.Session) bool {
			applicationExists, err := checkApplicationExists(
				fnCtx,
				plan.Name.ValueString(),
				junSess,
			)
			if err != nil {
				resp.Diagnostics.AddError(tfdiag.PostCheckErrSummary, err.Error())

				return false
			}
			if !applicationExists {
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

func (rsc *application) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data applicationData
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

func (rsc *application) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state applicationData
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

func (rsc *application) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state applicationData
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

func (rsc *application) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data applicationData

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

func checkApplicationExists(
	_ context.Context, name string, junSess *junos.Session,
) (
	bool, error,
) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"applications application " + name + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func (rscData *applicationData) fillID() {
	rscData.ID = types.StringValue(rscData.Name.ValueString())
}

func (rscData *applicationData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *applicationData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set applications application " + rscData.Name.ValueString() + " "

	if v := rscData.ApplicationProtocol.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"application-protocol "+v)
	}
	if v := rscData.Description.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := rscData.DestinationPort.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"destination-port \""+v+"\"")
	}
	if v := rscData.EtherType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"ether-type "+v)
	}
	if !rscData.InactivityTimeout.IsNull() {
		configSet = append(configSet, setPrefix+
			"inactivity-timeout "+utils.ConvI64toa(rscData.InactivityTimeout.ValueInt64()))
	} else if rscData.InactivityTimeoutNever.ValueBool() {
		configSet = append(configSet, setPrefix+"inactivity-timeout never")
	}
	if v := rscData.Protocol.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"protocol "+v)
	}
	if v := rscData.RPCProgramNumber.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"rpc-program-number "+v)
	}
	if v := rscData.SourcePort.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-port \""+v+"\"")
	}
	termName := make(map[string]struct{})
	for i, block := range rscData.Term {
		name := block.Name.ValueString()
		if _, ok := termName[name]; ok {
			return path.Root("term").AtListIndex(i).AtName("name"),
				fmt.Errorf("multiple term blocks with the same name %q", name)
		}
		termName[name] = struct{}{}

		blockSet, pathErr, err := block.configSet(setPrefix, path.Root("term").AtListIndex(i))
		if err != nil {
			return pathErr, err
		}
		configSet = append(configSet, blockSet...)
	}
	if v := rscData.UUID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"uuid "+v)
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *applicationBlockTerm) configSet(
	setPrefix string, pathRoot path.Path,
) (
	[]string, // configSet
	path.Path, // pathErr
	error, // error
) {
	setPrefix += "term " + block.Name.ValueString() + " "

	configSet := []string{
		setPrefix,
		setPrefix + "protocol " + block.Protocol.ValueString(),
	}

	if v := block.Alg.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"alg "+v)
	}
	if v := block.DestinationPort.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"destination-port \""+v+"\"")
	}
	if v := block.IcmpCode.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"icmp-code "+v)
	}
	if v := block.IcmpType.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"icmp-type "+v)
	}
	if v := block.Icmp6Code.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"icmp6-code "+v)
	}
	if v := block.Icmp6Type.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"icmp6-type "+v)
	}
	if !block.InactivityTimeout.IsNull() {
		configSet = append(configSet, setPrefix+
			"inactivity-timeout "+utils.ConvI64toa(block.InactivityTimeout.ValueInt64()))
		if block.InactivityTimeoutNever.ValueBool() {
			return configSet,
				pathRoot.AtName("inactivity_timeout_never"),
				fmt.Errorf("inactivity_timeout and inactivity_timeout_never cannot be configured together"+
					" in term block %q", block.Name.ValueString())
		}
	} else if block.InactivityTimeoutNever.ValueBool() {
		configSet = append(configSet, setPrefix+"inactivity-timeout never")
	}
	if v := block.RPCRrogramNumber.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"rpc-program-number "+v)
	}
	if v := block.SourcePort.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"source-port \""+v+"\"")
	}
	if v := block.UUID.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"uuid "+v)
	}

	return configSet, path.Empty(), nil
}

func (rscData *applicationData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"applications application " + name + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "application-protocol "):
				rscData.ApplicationProtocol = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "description "):
				rscData.Description = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "destination-port "):
				rscData.DestinationPort = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "ether-type "):
				rscData.EtherType = types.StringValue(itemTrim)
			case itemTrim == "inactivity-timeout never":
				rscData.InactivityTimeoutNever = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "inactivity-timeout "):
				rscData.InactivityTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
				if err != nil {
					return err
				}
			case balt.CutPrefixInString(&itemTrim, "protocol "):
				rscData.Protocol = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "rpc-program-number "):
				rscData.RPCProgramNumber = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "source-port "):
				rscData.SourcePort = types.StringValue(strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "term "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var term applicationBlockTerm
				rscData.Term, term = tfdata.ExtractBlockWithTFTypesString(
					rscData.Term, "Name", itemTrimFields[0],
				)
				term.Name = types.StringValue(itemTrimFields[0])
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				if err := term.read(strings.TrimPrefix(itemTrim, itemTrimFields[0]+" ")); err != nil {
					return err
				}
				rscData.Term = append(rscData.Term, term)
			case balt.CutPrefixInString(&itemTrim, "uuid "):
				rscData.UUID = types.StringValue(itemTrim)
			}
		}
	}

	return nil
}

func (block *applicationBlockTerm) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "protocol "):
		block.Protocol = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "alg "):
		block.Alg = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "destination-port "):
		block.DestinationPort = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "icmp-code "):
		block.IcmpCode = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "icmp-type "):
		block.IcmpType = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "icmp6-code "):
		block.Icmp6Code = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "icmp6-type "):
		block.Icmp6Type = types.StringValue(itemTrim)
	case itemTrim == "inactivity-timeout never":
		block.InactivityTimeoutNever = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "inactivity-timeout "):
		block.InactivityTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case balt.CutPrefixInString(&itemTrim, "rpc-program-number "):
		block.RPCRrogramNumber = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "source-port "):
		block.SourcePort = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "uuid "):
		block.UUID = types.StringValue(itemTrim)
	}

	return nil
}

func (rscData *applicationData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete applications application " + rscData.Name.ValueString(),
	}

	return junSess.ConfigSet(configSet)
}
