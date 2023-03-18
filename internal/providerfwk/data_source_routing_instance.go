package providerfwk

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &routingInstanceDataSource{}
	_ datasource.DataSourceWithConfigure = &routingInstanceDataSource{}
)

type routingInstanceDataSource struct {
	client *junos.Client
}

func (dsc *routingInstanceDataSource) typeName() string {
	return providerName + "_routing_instance"
}

func (dsc *routingInstanceDataSource) junosName() string {
	return "routing instance"
}

func newRoutingInstanceDataSource() datasource.DataSource {
	return &routingInstanceDataSource{}
}

func (dsc *routingInstanceDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *routingInstanceDataSource) Configure(
	ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*junos.Client)
	if !ok {
		unexpectedDataSourceConfigureType(ctx, req, resp)

		return
	}
	dsc.client = client
}

func (dsc *routingInstanceDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get configuration from a " + dsc.junosName() + ".",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the data source with format `<name>`.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of routing instance.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
					stringvalidator.NoneOfCaseInsensitive(junos.DefaultW),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of routing instance.",
			},
			"as": schema.StringAttribute{
				Computed:    true,
				Description: "Autonomous system number in plain number or `higher 16bits`.`lower 16 bits` (asdot notation) format.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Text description of routing instance.",
			},
			"instance_export": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Export policy for instance RIBs.",
			},
			"instance_import": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Import policy for instance RIBs.",
			},
			"interface": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of interfaces in routing instance.",
			},
			"route_distinguisher": schema.StringAttribute{
				Computed:    true,
				Description: "Route distinguisher for this instance.",
			},
			"router_id": schema.StringAttribute{
				Computed:    true,
				Description: "Router identifier.",
			},
			"vrf_export": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Export policy for VRF instance RIBs.",
			},
			"vrf_import": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Import policy for VRF instance RIBs.",
			},
			"vrf_target": schema.StringAttribute{
				Computed:    true,
				Description: "Target community to use in import and export.",
			},
			"vrf_target_auto": schema.BoolAttribute{
				Computed:    true,
				Description: "Auto derive import and export target community from BGP AS & L2.",
			},
			"vrf_target_export": schema.StringAttribute{
				Computed:    true,
				Description: "Target community to use when marking routes on export.",
			},
			"vrf_target_import": schema.StringAttribute{
				Computed:    true,
				Description: "Target community to use when filtering on import.",
			},
			"vtep_source_interface": schema.StringAttribute{
				Computed:    true,
				Description: "Source layer-3 IFL for VXLAN.",
			},
		},
	}
}

type routingInstanceDataSourceData struct {
	ID                  types.String   `tfsdk:"id"`
	Name                types.String   `tfsdk:"name"`
	Type                types.String   `tfsdk:"type"`
	AS                  types.String   `tfsdk:"as"`
	Description         types.String   `tfsdk:"description"`
	InstanceExport      []types.String `tfsdk:"instance_export"`
	InstanceImport      []types.String `tfsdk:"instance_import"`
	Interface           []types.String `tfsdk:"interface"`
	RouteDistinguisher  types.String   `tfsdk:"route_distinguisher"`
	RouterID            types.String   `tfsdk:"router_id"`
	VRFExport           []types.String `tfsdk:"vrf_export"`
	VRFImport           []types.String `tfsdk:"vrf_import"`
	VRFTarget           types.String   `tfsdk:"vrf_target"`
	VRFTargetAuto       types.Bool     `tfsdk:"vrf_target_auto"`
	VRFTargetExport     types.String   `tfsdk:"vrf_target_export"`
	VRFTargetImport     types.String   `tfsdk:"vrf_target_import"`
	VTEPSourceInterface types.String   `tfsdk:"vtep_source_interface"`
}

func (dsc *routingInstanceDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var name types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}
	junSess, err := dsc.client.StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Start Session Error", err.Error())

		return
	}
	defer junSess.Close()

	var rscData routingInstanceData
	junos.MutexLock()
	err = rscData.read(ctx, name.ValueString(), junSess)
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())

		return
	}
	if rscData.ID.IsNull() {
		resp.Diagnostics.AddError(
			"Not Found Error",
			fmt.Sprintf(dsc.junosName()+" %q doesn't exist", name.ValueString()),
		)

		return
	}

	var data routingInstanceDataSourceData
	data.CopyFromResourceData(rscData)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (dscData *routingInstanceDataSourceData) CopyFromResourceData(rscData routingInstanceData) {
	dscData.ID = rscData.ID
	dscData.Name = rscData.Name
	dscData.Type = rscData.Type
	dscData.AS = rscData.AS
	dscData.Description = rscData.Description
	dscData.InstanceExport = rscData.InstanceExport
	dscData.InstanceImport = rscData.InstanceImport
	dscData.Interface = rscData.Interface
	dscData.RouteDistinguisher = rscData.RouteDistinguisher
	dscData.RouterID = rscData.RouterID
	dscData.VRFExport = rscData.VRFExport
	dscData.VRFImport = rscData.VRFImport
	dscData.VRFTarget = rscData.VRFTarget
	dscData.VRFTargetAuto = rscData.VRFTargetAuto
	dscData.VRFTargetExport = rscData.VRFTargetExport
	dscData.VRFTargetImport = rscData.VRFTargetImport
	dscData.VTEPSourceInterface = rscData.VTEPSourceInterface
}
