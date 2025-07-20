package provider

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
	_ resource.Resource                   = &layer2Control{}
	_ resource.ResourceWithConfigure      = &layer2Control{}
	_ resource.ResourceWithValidateConfig = &layer2Control{}
	_ resource.ResourceWithImportState    = &layer2Control{}
	_ resource.ResourceWithUpgradeState   = &layer2Control{}
)

type layer2Control struct {
	client *junos.Client
}

func newLayer2ControlResource() resource.Resource {
	return &layer2Control{}
}

func (rsc *layer2Control) typeName() string {
	return providerName + "_layer2_control"
}

func (rsc *layer2Control) junosName() string {
	return "protocols layer2-control"
}

func (rsc *layer2Control) junosClient() *junos.Client {
	return rsc.client
}

func (rsc *layer2Control) Metadata(
	_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = rsc.typeName()
}

func (rsc *layer2Control) Configure(
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

func (rsc *layer2Control) Schema(
	_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Version:     1,
		Description: "Configure static configuration in `" + rsc.junosName() + "` block",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the resource with value `layer2_control`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"nonstop_bridging": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable nonstop operation.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"bpdu_block": schema.SingleNestedBlock{
				Description: "Block BPDU on interface (BPDU Protect).",
				Attributes: map[string]schema.Attribute{
					"disable_timeout": schema.Int64Attribute{
						Optional:    true,
						Description: "Disable timeout for BPDU Protect (seconds).",
						Validators: []validator.Int64{
							int64validator.Between(10, 3600),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"interface": schema.SetNestedBlock{
						Description: "For each interface, to block BPDU on.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "Name of interface.",
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
										tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
										tfvalidator.StringDotExclusion(),
									},
								},
								"disable": schema.BoolAttribute{
									Optional:    true,
									Description: "Disable bpdu-block on a port.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
								"drop": schema.BoolAttribute{
									Optional:    true,
									Description: "Drop xSTP BPDUs.",
									Validators: []validator.Bool{
										tfvalidator.BoolTrue(),
									},
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					tfplanmodifier.BlockRemoveNull(),
				},
			},
			"mac_rewrite_interface": schema.SetNestedBlock{
				Description: "For each interface, Mac rewrite functionality.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of interface.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
								tfvalidator.StringDotExclusion(),
							},
						},
						"enable_all_ifl": schema.BoolAttribute{
							Optional:    true,
							Description: "Enable tunneling for all the IFLs under the interface.",
							Validators: []validator.Bool{
								tfvalidator.BoolTrue(),
							},
						},
						"protocol": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "Protocols for which mac rewrite need to be enabled.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.NoNullValues(),
								setvalidator.ValueStringsAre(
									stringvalidator.LengthAtLeast(1),
									tfvalidator.StringFormat(tfvalidator.DefaultFormat),
								),
							},
						},
					},
				},
			},
		},
	}
}

type layer2ControlData struct {
	ID                  types.String                            `tfsdk:"id"`
	NonstopBridging     types.Bool                              `tfsdk:"nonstop_bridging"`
	BpduBlock           *layer2ControlBlockBpduBlock            `tfsdk:"bpdu_block"`
	MacRewriteInterface []layer2ControlBlockMacRewriteInterface `tfsdk:"mac_rewrite_interface"`
}

type layer2ControlConfig struct {
	ID                  types.String                       `tfsdk:"id"`
	NonstopBridging     types.Bool                         `tfsdk:"nonstop_bridging"`
	BpduBlock           *layer2ControlBlockBpduBlockConfig `tfsdk:"bpdu_block"`
	MacRewriteInterface types.Set                          `tfsdk:"mac_rewrite_interface"`
}

type layer2ControlBlockBpduBlock struct {
	DisableTimeout types.Int64                                 `tfsdk:"disable_timeout"`
	Interface      []layer2ControlBlockBpduBlockBlockInterface `tfsdk:"interface"`
}

type layer2ControlBlockBpduBlockConfig struct {
	DisableTimeout types.Int64 `tfsdk:"disable_timeout"`
	Interface      types.Set   `tfsdk:"interface"`
}

type layer2ControlBlockBpduBlockBlockInterface struct {
	Name    types.String `tfsdk:"name"    tfdata:"identifier"`
	Disable types.Bool   `tfsdk:"disable"`
	Drop    types.Bool   `tfsdk:"drop"`
}

type layer2ControlBlockMacRewriteInterface struct {
	Name         types.String   `tfsdk:"name"           tfdata:"identifier"`
	EnableAllIfl types.Bool     `tfsdk:"enable_all_ifl"`
	Protocol     []types.String `tfsdk:"protocol"`
}

func (rsc *layer2Control) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var config layer2ControlConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.BpduBlock != nil {
		if !config.BpduBlock.Interface.IsNull() && !config.BpduBlock.Interface.IsUnknown() {
			var configInterface []layer2ControlBlockBpduBlockBlockInterface
			asDiags := config.BpduBlock.Interface.ElementsAs(ctx, &configInterface, false)
			if asDiags.HasError() {
				resp.Diagnostics.Append(asDiags...)

				return
			}

			interfaceName := make(map[string]struct{})
			for _, block := range configInterface {
				if block.Name.IsUnknown() {
					continue
				}
				name := block.Name.ValueString()
				if _, ok := interfaceName[name]; ok {
					resp.Diagnostics.AddAttributeError(
						path.Root("bpdu_block").AtName("interface"),
						tfdiag.DuplicateConfigErrSummary,
						fmt.Sprintf("multiple interface blocks with the same name %q"+
							" in bpdu_block block", name),
					)
				}
				interfaceName[name] = struct{}{}
			}
		}
	}
	if !config.MacRewriteInterface.IsNull() && !config.MacRewriteInterface.IsUnknown() {
		var configMacRewriteInterface []layer2ControlBlockMacRewriteInterface
		asDiags := config.MacRewriteInterface.ElementsAs(ctx, &configMacRewriteInterface, false)
		if asDiags.HasError() {
			resp.Diagnostics.Append(asDiags...)

			return
		}

		macRewriteInterfaceName := make(map[string]struct{})
		for _, block := range configMacRewriteInterface {
			if block.Name.IsUnknown() {
				continue
			}
			name := block.Name.ValueString()
			if _, ok := macRewriteInterfaceName[name]; ok {
				resp.Diagnostics.AddAttributeError(
					path.Root("mac_rewrite_interface"),
					tfdiag.DuplicateConfigErrSummary,
					fmt.Sprintf("multiple mac_rewrite_interface blocks with the same name %q", name),
				)
			}
			macRewriteInterfaceName[name] = struct{}{}
		}
	}
}

func (rsc *layer2Control) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan layer2ControlData
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

func (rsc *layer2Control) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state, data layer2ControlData
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

func (rsc *layer2Control) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, state layer2ControlData
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

func (rsc *layer2Control) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state layer2ControlData
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

func (rsc *layer2Control) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	var data layer2ControlData

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

func (rscData *layer2ControlData) fillID() {
	rscData.ID = types.StringValue("layer2_control")
}

func (rscData *layer2ControlData) nullID() bool {
	return rscData.ID.IsNull()
}

func (rscData *layer2ControlData) set(
	_ context.Context, junSess *junos.Session,
) (
	path.Path, error,
) {
	configSet := make([]string, 0, 100)
	setPrefix := "set protocols layer2-control "

	if rscData.NonstopBridging.ValueBool() {
		configSet = append(configSet, setPrefix+"nonstop-bridging")
	}

	if rscData.BpduBlock != nil {
		configSet = append(configSet, setPrefix+"bpdu-block")

		if !rscData.BpduBlock.DisableTimeout.IsNull() {
			configSet = append(configSet, setPrefix+"bpdu-block disable-timeout "+
				utils.ConvI64toa(rscData.BpduBlock.DisableTimeout.ValueInt64()))
		}

		interfaceName := make(map[string]struct{})
		for _, block := range rscData.BpduBlock.Interface {
			name := block.Name.ValueString()
			if _, ok := interfaceName[name]; ok {
				return path.Root("bpdu_block").AtName("interface"),
					fmt.Errorf("multiple interface blocks with the same name %q"+
						" in bpdu_block block", name)
			}
			interfaceName[name] = struct{}{}

			configSet = append(configSet, setPrefix+"bpdu-block interface "+name)
			if block.Disable.ValueBool() {
				configSet = append(configSet, setPrefix+"bpdu-block interface "+name+" disable")
			}
			if block.Drop.ValueBool() {
				configSet = append(configSet, setPrefix+"bpdu-block interface "+name+" drop")
			}
		}
	}
	macRewriteInterfaceName := make(map[string]struct{})
	for _, block := range rscData.MacRewriteInterface {
		name := block.Name.ValueString()
		if _, ok := macRewriteInterfaceName[name]; ok {
			return path.Root("mac_rewrite_interface"),
				fmt.Errorf("multiple mac_rewrite_interface blocks with the same name %q", name)
		}
		macRewriteInterfaceName[name] = struct{}{}

		configSet = append(configSet, setPrefix+"mac-rewrite interface "+name)
		if block.EnableAllIfl.ValueBool() {
			configSet = append(configSet, setPrefix+"mac-rewrite interface "+name+" enable-all-ifl")
		}
		for _, v := range block.Protocol {
			configSet = append(configSet, setPrefix+"mac-rewrite interface "+name+" protocol "+v.ValueString())
		}
	}

	return path.Empty(), junSess.ConfigSet(configSet)
}

func (rscData *layer2ControlData) read(
	_ context.Context, junSess *junos.Session,
) error {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"protocols layer2-control" + junos.PipeDisplaySetRelative)
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
			case balt.CutPrefixInString(&itemTrim, "bpdu-block"):
				if rscData.BpduBlock == nil {
					rscData.BpduBlock = &layer2ControlBlockBpduBlock{}
				}

				switch {
				case balt.CutPrefixInString(&itemTrim, " disable-timeout "):
					rscData.BpduBlock.DisableTimeout, err = tfdata.ConvAtoi64Value(itemTrim)
					if err != nil {
						return err
					}
				case balt.CutPrefixInString(&itemTrim, " interface "):
					name := tfdata.FirstElementOfJunosLine(itemTrim)
					var interFace layer2ControlBlockBpduBlockBlockInterface
					rscData.BpduBlock.Interface, interFace = tfdata.ExtractBlock(rscData.BpduBlock.Interface, types.StringValue(name))

					if balt.CutPrefixInString(&itemTrim, name+" ") {
						switch {
						case itemTrim == "disable":
							interFace.Disable = types.BoolValue(true)
						case itemTrim == "drop":
							interFace.Drop = types.BoolValue(true)
						}
					}
					rscData.BpduBlock.Interface = append(rscData.BpduBlock.Interface, interFace)
				}
			case balt.CutPrefixInString(&itemTrim, "mac-rewrite interface "):
				name := tfdata.FirstElementOfJunosLine(itemTrim)
				var macRewriteInterface layer2ControlBlockMacRewriteInterface
				rscData.MacRewriteInterface, macRewriteInterface = tfdata.ExtractBlock(
					rscData.MacRewriteInterface, types.StringValue(name),
				)

				if balt.CutPrefixInString(&itemTrim, name+" ") {
					switch {
					case itemTrim == "enable-all-ifl":
						macRewriteInterface.EnableAllIfl = types.BoolValue(true)
					case balt.CutPrefixInString(&itemTrim, "protocol "):
						macRewriteInterface.Protocol = append(macRewriteInterface.Protocol, types.StringValue(itemTrim))
					}
				}
				rscData.MacRewriteInterface = append(rscData.MacRewriteInterface, macRewriteInterface)
			case itemTrim == "nonstop-bridging":
				rscData.NonstopBridging = types.BoolValue(true)
			}
		}
	}

	return nil
}

func (rscData *layer2ControlData) del(
	_ context.Context, junSess *junos.Session,
) error {
	configSet := []string{
		"delete protocols layer2-control",
	}

	return junSess.ConfigSet(configSet)
}
