package providerfwk

import (
	"context"
	"errors"
	"fmt"
	"regexp"
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
	_ resource.Resource                   = &virtualChassis{}
	_ resource.ResourceWithConfigure      = &virtualChassis{}
	_ resource.ResourceWithValidateConfig = &virtualChassis{}
	_ resource.ResourceWithImportState    = &virtualChassis{}
)

type virtualChassis struct {
	client *junos.Client
}

func newVirtualChassisResource() resource.Resource {
	return &virtualChassis{}
}

func (rsc *virtualChassis) typeName() string {
	return providerName + "_virtual_chassis"
}

func (rsc *virtualChassis) junosName() string {
	return "virtual-chassis"
}

func (rsc *virtualChassis) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *virtualChassis) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *virtualChassis) Configure(
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

func (rsc *virtualChassis) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Configure `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `virtual-chassis`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_sw_update": schema.BoolAttribute{
				Optional:    true,
				Description: "Auto software update.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"auto_sw_update_package_name": schema.StringAttribute{
				Optional:    true,
				Description: "URL or pathname of software package to auto software update.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"graceful_restart_disable": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable graceful restart.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"identifier": schema.StringAttribute{
				Optional:    true,
				Description: "Virtual chassis identifier, of type ISO system-id.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^[0-9a-f]{4}\.[0-9a-f]{4}\.[0-9a-f]{4}$`),
						"must be of type ISO system-id",
					),
				},
			},
			"mac_persistence_timer": schema.StringAttribute{
				Optional:    true,
				Description: "MAC persistence time (minutes) or disable.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(
						`^([1-9]|[1-5][0-9]|60|disable)$`),
						"must be a number between 1 to 60 or disable",
					),
				},
			},
			"no_split_detection": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable split detection.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"preprovisioned": schema.BoolAttribute{
				Optional:    true,
				Description: "Only accept preprovisioned members.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"vcp_no_hold_time": schema.BoolAttribute{
				Optional:    true,
				Description: "Set no hold time for vcp interfaces.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"alias": schema.SetNestedBlock{
				Description: "Alias for serial number.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"serial_number": schema.StringAttribute{
							Required:    true,
							Description: "Member's serial number.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 12),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"alias_name": schema.StringAttribute{
							Required:    true,
							Description: "Alias name for this serial-number.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 12),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
				},
			},
			"member": schema.ListNestedBlock{
				Description: "Member of virtual chassis configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Required:    true,
							Description: "Member identifier.",
							Validators: []validator.Int64{
								int64validator.Between(0, 9),
							},
						},
						"location": schema.StringAttribute{
							Optional:    true,
							Description: "Member's location.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 254),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
						"mastership_priority": schema.Int64Attribute{
							Optional:    true,
							Description: "Member's mastership priority.",
							Validators: []validator.Int64{
								int64validator.Between(0, 255),
							},
						},
						"no_management_vlan": schema.BoolAttribute{
							Optional:    true,
							Description: "Disable management VLAN.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"role": schema.StringAttribute{
							Optional:    true,
							Description: "Member's role.",
							Validators: []validator.String{
								stringvalidator.OneOf("line-card", "routing-engine"),
							},
						},
						"serial_number": schema.StringAttribute{
							Optional:    true,
							Description: "Member's serial number.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 12),
								tfvalidator.StringDoubleQuoteExclusion(),
							},
						},
					},
				},
			},
			"traceoptions": schema.SingleNestedBlock{
				Description: "Trace options for virtual chassis.",
				Attributes: map[string]schema.Attribute{
					"flag": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "Tracing parameters.",
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.DefaultFormatAndSpace),
							),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"file": schema.SingleNestedBlock{
						Description: "Declare `file` configuration.",
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Required:    false, // true when SingleNestedBlock is specified
								Optional:    true,
								Description: "Name of file in which to write trace information.",
								Validators: []validator.String{
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringDoubleQuoteExclusion(),
									tfvalidator.StringSpaceExclusion(),
									tfvalidator.StringRuneExclusion('/', '%'),
								},
							},
							"files": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum number of trace files.",
								Validators: []validator.Int64{
									int64validator.Between(2, 1000),
								},
							},
							"no_stamp": schema.BoolAttribute{
								Optional:    true,
								Description: "Do not timestamp trace file.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"replace": schema.BoolAttribute{
								Optional:    true,
								Description: "Replace trace file rather than appending to it.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"size": schema.Int64Attribute{
								Optional:    true,
								Description: "Maximum trace file size.",
								Validators: []validator.Int64{
									int64validator.Between(10240, 4294967295),
								},
							},
							"world_readable": schema.BoolAttribute{
								Optional:    true,
								Description: "Allow any user to read the log file.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
							"no_world_readable": schema.BoolAttribute{
								Optional:    true,
								Description: "Don't allow any user to read the log file.",
								Validators: []validator.Bool{
									tfvalidator.BoolTrue(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							tfplanmodifier.BlockRemoveNull(),
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

type virtualChassisData struct {
	ID                      types.String                     `tfsdk:"id"`
	AutoSWUpdate            types.Bool                       `tfsdk:"auto_sw_update"`
	AutoSWUpdatePackageName types.String                     `tfsdk:"auto_sw_update_package_name"`
	GracefulRestartDisable  types.Bool                       `tfsdk:"graceful_restart_disable"`
	Identifier              types.String                     `tfsdk:"identifier"`
	MacPersistenceTimer     types.String                     `tfsdk:"mac_persistence_timer"`
	NoSplitDetection        types.Bool                       `tfsdk:"no_split_detection"`
	Preprovisioned          types.Bool                       `tfsdk:"preprovisioned"`
	VcpNoHoldTime           types.Bool                       `tfsdk:"vcp_no_hold_time"`
	Alias                   []virtualChassisBlockAlias       `tfsdk:"alias"`
	Member                  []virtualChassisBlockMember      `tfsdk:"member"`
	Traceoptions            *virtualChassisBlockTraceoptions `tfsdk:"traceoptions"`
}

type virtualChassisConfig struct {
	ID                      types.String                           `tfsdk:"id"`
	AutoSWUpdate            types.Bool                             `tfsdk:"auto_sw_update"`
	AutoSWUpdatePackageName types.String                           `tfsdk:"auto_sw_update_package_name"`
	GracefulRestartDisable  types.Bool                             `tfsdk:"graceful_restart_disable"`
	Identifier              types.String                           `tfsdk:"identifier"`
	MacPersistenceTimer     types.String                           `tfsdk:"mac_persistence_timer"`
	NoSplitDetection        types.Bool                             `tfsdk:"no_split_detection"`
	Preprovisioned          types.Bool                             `tfsdk:"preprovisioned"`
	VcpNoHoldTime           types.Bool                             `tfsdk:"vcp_no_hold_time"`
	Alias                   types.Set                              `tfsdk:"alias"`
	Member                  types.List                             `tfsdk:"member"`
	Traceoptions            *virtualChassisBlockTraceoptionsConfig `tfsdk:"traceoptions"`
}

type virtualChassisBlockAlias struct {
	SerialNumber types.String `tfsdk:"serial_number"`
	AliasName    types.String `tfsdk:"alias_name"`
}

type virtualChassisBlockMember struct {
	ID                 types.Int64  `tfsdk:"id"`
	Location           types.String `tfsdk:"location"`
	MastershipPriority types.Int64  `tfsdk:"mastership_priority"`
	NoManagementVlan   types.Bool   `tfsdk:"no_management_vlan"`
	Role               types.String `tfsdk:"role"`
	SerialNumber       types.String `tfsdk:"serial_number"`
}

func (block *virtualChassisBlockMember) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block, "ID")
}

type virtualChassisBlockTraceoptions struct {
	Flag []types.String                            `tfsdk:"flag"`
	File *virtualChassisBlockTraceoptionsBlockFile `tfsdk:"file"`
}

func (block *virtualChassisBlockTraceoptions) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type virtualChassisBlockTraceoptionsConfig struct {
	Flag types.Set                                 `tfsdk:"flag"`
	File *virtualChassisBlockTraceoptionsBlockFile `tfsdk:"file"`
}

func (block *virtualChassisBlockTraceoptionsConfig) isEmpty() bool {
	return tfdata.CheckBlockIsEmpty(block)
}

type virtualChassisBlockTraceoptionsBlockFile struct {
	Name            types.String `tfsdk:"name"`
	Files           types.Int64  `tfsdk:"files"`
	NoStamp         types.Bool   `tfsdk:"no_stamp"`
	Replace         types.Bool   `tfsdk:"replace"`
	Size            types.Int64  `tfsdk:"size"`
	WorldReadable   types.Bool   `tfsdk:"world_readable"`
	NoWorldReadable types.Bool   `tfsdk:"no_world_readable"`
}

func (rsc *virtualChassis) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config virtualChassisConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.AutoSWUpdatePackageName.IsNull() && !config.AutoSWUpdatePackageName.IsUnknown() {
		if config.AutoSWUpdate.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("auto_sw_update_package_name"),
				tfdiag.MissingConfigErrSummary,
				"auto_sw_update must be specified with auto_sw_update_package_name",
			)
		}
	}

	if !config.Alias.IsNull() && !config.Alias.IsUnknown() {
		var configAlias []virtualChassisBlockAlias
		asDiags := config.Alias.ElementsAs(ctx, &configAlias, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		aliasSerialNumber := make(map[string]struct{})
		aliasAliasName := make(map[string]struct{})
		for _, block := range configAlias {
			if block.SerialNumber.IsUnknown() {
				continue
			}

			serialNumber := block.SerialNumber.ValueString()
			if _, ok := aliasSerialNumber[serialNumber]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("alias"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple alias blocks with the same serial_number %q", serialNumber),
				)
			}
			aliasSerialNumber[serialNumber] = struct{}{}

			aliasName := block.AliasName.ValueString()
			if _, ok := aliasAliasName[aliasName]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("alias"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple alias blocks with the same alias_name %q", aliasName),
				)
			}
			aliasAliasName[aliasName] = struct{}{}
		}
	}
	if !config.Member.IsNull() && !config.Member.IsUnknown() {
		var configMember []virtualChassisBlockMember
		asDiags := config.Member.ElementsAs(ctx, &configMember, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		memberID := make(map[int64]struct{})
		for i, block := range configMember {
			if !block.ID.IsUnknown() {
				id := block.ID.ValueInt64()
				if _, ok := memberID[id]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("member").AtListIndex(i).AtName("id"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple member blocks with the same id %d", id),
					)
				}
				memberID[id] = struct{}{}
			}
			if block.isEmpty() {
				resp.Diagnostics.AddAttributeError(
					path.Root("member").AtListIndex(i).AtName("id"),
					tfdiag.MissingConfigErrSummary,
					fmt.Sprintf("member block %d is empty", block.ID.ValueInt64()),
				)
			}
		}
	}
	if config.Traceoptions != nil {
		if config.Traceoptions.isEmpty() {
			resp.Diagnostics.AddAttributeError(
				path.Root("traceoptions").AtName("*"),
				tfdiag.MissingConfigErrSummary,
				"traceoptions block is empty",
			)
		}

		if config.Traceoptions.File != nil {
			if config.Traceoptions.File.Name.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("traceoptions").AtName("file").AtName("name"),
					tfdiag.MissingConfigErrSummary,
					"name must be specified in file block in traceoptions block",
				)
			}
		}
	}
}

func (rsc *virtualChassis) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan virtualChassisData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultResourceCreate(
		ctx,
		rsc,
		nil,
		nil,
		&plan,
		resp,
	)
}

func (rsc *virtualChassis) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data virtualChassisData
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
		nil,
		resp,
	)
}

func (rsc *virtualChassis) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state virtualChassisData
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

func (rsc *virtualChassis) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state virtualChassisData
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

func (rsc *virtualChassis) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data virtualChassisData

	var _ resourceDataReadWithoutArg = &data
	defaultResourceImportState(
		ctx,
		rsc,
		&data,
		req,
		resp,
		"",
	)
}

func (rscData *virtualChassisData) fillID() {
	rscData.ID = types.StringValue("virtual-chassis")
}

func (rscData *virtualChassisData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *virtualChassisData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0)
	setPrefix := "set virtual-chassis "

	if rscData.AutoSWUpdate.ValueBool() {
		configSet = append(configSet, setPrefix+"auto-sw-update")
		if v := rscData.AutoSWUpdatePackageName.ValueString(); v != "" {
			configSet = append(configSet, setPrefix+"auto-sw-update package-name \""+v+"\"")
		}
	} else if rscData.AutoSWUpdatePackageName.ValueString() != "" {
		return path.Root("auto_sw_update_package_name"),
			errors.New("auto_sw_update must be specified with auto_sw_update_package_name")
	}
	if rscData.GracefulRestartDisable.ValueBool() {
		configSet = append(configSet, setPrefix+"graceful-restart disable")
	}
	if v := rscData.Identifier.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"id "+v)
	}
	if v := rscData.MacPersistenceTimer.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"mac-persistence-timer "+v)
	}
	if rscData.NoSplitDetection.ValueBool() {
		configSet = append(configSet, setPrefix+"no-split-detection")
	}
	if rscData.Preprovisioned.ValueBool() {
		configSet = append(configSet, setPrefix+"preprovisioned")
	}
	if rscData.VcpNoHoldTime.ValueBool() {
		configSet = append(configSet, setPrefix+"vcp-no-hold-time")
	}

	aliasSerialNumber := make(map[string]struct{})
	aliasAliasName := make(map[string]struct{})
	for _, block := range rscData.Alias {
		serialNumber := block.SerialNumber.ValueString()
		if _, ok := aliasSerialNumber[serialNumber]; ok {
			return path.Root("alias"),
				fmt.Errorf("multiple alias blocks with the same serial_number %q", serialNumber)
		}
		aliasSerialNumber[serialNumber] = struct{}{}

		aliasName := block.AliasName.ValueString()
		if _, ok := aliasAliasName[aliasName]; ok {
			return path.Root("alias"),
				fmt.Errorf("multiple alias blocks with the same alias_name %q", aliasName)
		}
		aliasAliasName[aliasName] = struct{}{}

		configSet = append(configSet, setPrefix+"aliases"+
			" serial-number \""+serialNumber+"\""+
			" alias-name \""+aliasName+"\"")
	}
	memberID := make(map[int64]struct{})
	for i, block := range rscData.Member {
		id := block.ID.ValueInt64()
		if _, ok := memberID[id]; ok {
			return path.Root("member").AtListIndex(i).AtName("id"),
				fmt.Errorf("multiple member blocks with the same id %d", id)
		}
		memberID[id] = struct{}{}
		if block.isEmpty() {
			return path.Root("member").AtListIndex(i).AtName("id"),
				fmt.Errorf("member block %d is empty", id)
		}

		configSet = append(configSet, block.configSet()...)
	}
	if rscData.Traceoptions != nil {
		if rscData.Traceoptions.isEmpty() {
			return path.Root("traceoptions").AtName("*"),
				errors.New("traceoptions block is empty")
		}

		for _, v := range rscData.Traceoptions.Flag {
			configSet = append(configSet, setPrefix+"traceoptions flag "+v.ValueString())
		}
		if rscData.Traceoptions.File != nil {
			configSet = append(configSet, rscData.Traceoptions.File.configSet()...)
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (block *virtualChassisBlockMember) configSet() []string {
	configSet := make([]string, 0, 1)
	setPrefix := "set virtual-chassis member " + utils.ConvI64toa(block.ID.ValueInt64()) + " "

	if v := block.Location.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"location \""+v+"\"")
	}
	if !block.MastershipPriority.IsNull() {
		configSet = append(configSet, setPrefix+"mastership-priority "+
			utils.ConvI64toa(block.MastershipPriority.ValueInt64()))
	}
	if block.NoManagementVlan.ValueBool() {
		configSet = append(configSet, setPrefix+"no-management-vlan")
	}
	if v := block.Role.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"role "+v)
	}
	if v := block.SerialNumber.ValueString(); v != "" {
		configSet = append(configSet, setPrefix+"serial-number \""+v+"\"")
	}

	return configSet
}

func (block *virtualChassisBlockTraceoptionsBlockFile) configSet() []string {
	setPrefix := "set virtual-chassis traceoptions file "

	configSet := []string{
		setPrefix + "\"" + block.Name.ValueString() + "\"",
	}

	if !block.Files.IsNull() {
		configSet = append(configSet, setPrefix+"files "+
			utils.ConvI64toa(block.Files.ValueInt64()))
	}
	if block.NoStamp.ValueBool() {
		configSet = append(configSet, setPrefix+"no-stamp")
	}
	if block.Replace.ValueBool() {
		configSet = append(configSet, setPrefix+"replace")
	}
	if !block.Size.IsNull() {
		configSet = append(configSet, setPrefix+"size "+
			utils.ConvI64toa(block.Size.ValueInt64()))
	}
	if block.WorldReadable.ValueBool() {
		configSet = append(configSet, setPrefix+"world-readable")
	}
	if block.NoWorldReadable.ValueBool() {
		configSet = append(configSet, setPrefix+"no-world-readable")
	}

	return configSet
}

func (rscData *virtualChassisData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"virtual-chassis" + junos.PipeDisplaySetRelative)
	if err != nil {
		return err
	}
	rscData.fillID()
	if showConfig != junos.EmptyW {
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "auto-sw-update"):
				rscData.AutoSWUpdate = types.BoolValue(true)
				if balt.CutPrefixInString(&itemTrim, " package-name ") {
					rscData.AutoSWUpdatePackageName = types.StringValue(strings.Trim(itemTrim, "\""))
				}
			case itemTrim == "graceful-restart disable":
				rscData.GracefulRestartDisable = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "id "):
				rscData.Identifier = types.StringValue(itemTrim)
			case balt.CutPrefixInString(&itemTrim, "mac-persistence-timer "):
				rscData.MacPersistenceTimer = types.StringValue(itemTrim)
			case itemTrim == "no-split-detection":
				rscData.NoSplitDetection = types.BoolValue(true)
			case itemTrim == "preprovisioned":
				rscData.Preprovisioned = types.BoolValue(true)
			case itemTrim == "vcp-no-hold-time":
				rscData.VcpNoHoldTime = types.BoolValue(true)
			case balt.CutPrefixInString(&itemTrim, "aliases serial-number "):
				serialNumber := tfdata.FirstElementOfJunosLine(itemTrim)
				alias := virtualChassisBlockAlias{
					SerialNumber: types.StringValue(strings.Trim(serialNumber, "\"")),
				}
				if balt.CutPrefixInString(&itemTrim, serialNumber+" alias-name ") {
					alias.AliasName = types.StringValue(strings.Trim(itemTrim, "\""))
					rscData.Alias = append(rscData.Alias, alias)
				}
			case balt.CutPrefixInString(&itemTrim, "member "):
				itemTrimFields := strings.Split(itemTrim, " ")
				var member virtualChassisBlockMember
				memberID, err := tfdata.ConvAtoi64Value(itemTrimFields[0])
				if err != nil {
					return err
				}
				rscData.Member, member = tfdata.ExtractBlockWithTFTypesInt64(
					rscData.Member, "ID", memberID.ValueInt64(),
				)
				member.ID = memberID
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")

				if err := member.read(itemTrim); err != nil {
					return err
				}
				rscData.Member = append(rscData.Member, member)
			case balt.CutPrefixInString(&itemTrim, "traceoptions "):
				if rscData.Traceoptions == nil {
					rscData.Traceoptions = &virtualChassisBlockTraceoptions{}
				}
				switch {
				case balt.CutPrefixInString(&itemTrim, "flag "):
					rscData.Traceoptions.Flag = append(rscData.Traceoptions.Flag, types.StringValue(itemTrim))
				case balt.CutPrefixInString(&itemTrim, "file "):
					if rscData.Traceoptions.File == nil {
						rscData.Traceoptions.File = &virtualChassisBlockTraceoptionsBlockFile{}
					}

					if err := rscData.Traceoptions.File.read(itemTrim); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (block *virtualChassisBlockMember) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "location "):
		block.Location = types.StringValue(strings.Trim(itemTrim, "\""))
	case balt.CutPrefixInString(&itemTrim, "mastership-priority "):
		block.MastershipPriority, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "no-management-vlan":
		block.NoManagementVlan = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "role "):
		block.Role = types.StringValue(itemTrim)
	case balt.CutPrefixInString(&itemTrim, "serial-number "):
		block.SerialNumber = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return nil
}

func (block *virtualChassisBlockTraceoptionsBlockFile) read(itemTrim string) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "files "):
		block.Files, err = tfdata.ConvAtoi64Value(itemTrim)
		if err != nil {
			return err
		}
	case itemTrim == "no-stamp":
		block.NoStamp = types.BoolValue(true)
	case itemTrim == "replace":
		block.Replace = types.BoolValue(true)
	case balt.CutPrefixInString(&itemTrim, "size "):
		var size types.Int64
		switch {
		case balt.CutSuffixInString(&itemTrim, "k"):
			size, err = tfdata.ConvAtoi64Value(itemTrim)
			size = types.Int64Value(size.ValueInt64() * 1024)
		case balt.CutSuffixInString(&itemTrim, "m"):
			size, err = tfdata.ConvAtoi64Value(itemTrim)
			size = types.Int64Value(size.ValueInt64() * 1024 * 1024)
		case balt.CutSuffixInString(&itemTrim, "g"):
			size, err = tfdata.ConvAtoi64Value(itemTrim)
			size = types.Int64Value(size.ValueInt64() * 1024 * 1024 * 1024)
		default:
			size, err = tfdata.ConvAtoi64Value(itemTrim)
		}
		if err != nil {
			return err
		}
		block.Size = size
	case itemTrim == "world-readable":
		block.WorldReadable = types.BoolValue(true)
	case itemTrim == "no-world-readable":
		block.NoWorldReadable = types.BoolValue(true)
	default:
		block.Name = types.StringValue(strings.Trim(itemTrim, "\""))
	}

	return nil
}

func (rscData *virtualChassisData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete virtual-chassis",
	}

	return junSess.ConfigSet(configSet)
}
