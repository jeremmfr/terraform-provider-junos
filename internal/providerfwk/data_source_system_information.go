package providerfwk

import (
	"context"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &systemInformationDataSource{}
	_ datasource.DataSourceWithConfigure = &systemInformationDataSource{}
)

type systemInformationDataSource struct {
	client *junos.Client
}

func (dsc *systemInformationDataSource) typeName() string {
	return providerName + "_system_information"
}

func (dsc *systemInformationDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newSystemInformationDataSource() datasource.DataSource {
	return &systemInformationDataSource{}
}

func (dsc *systemInformationDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *systemInformationDataSource) Configure(
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

func (dsc *systemInformationDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get information of the Junos device system information.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Hostname of the Junos device or `Null-Hostname` if not set.",
			},
			"hardware_model": schema.StringAttribute{
				Computed:    true,
				Description: "Type of hardware/software of Junos device.",
			},
			"os_name": schema.StringAttribute{
				Computed:    true,
				Description: "Operating system name of Junos.",
			},
			"os_version": schema.StringAttribute{
				Computed:    true,
				Description: "Software version of Junos.",
			},
			"serial_number": schema.StringAttribute{
				Computed:    true,
				Description: "Serial number of the device.",
			},
			"cluster_node": schema.BoolAttribute{
				Computed:    true,
				Description: "Boolean flag that indicates if device is part of a cluster or not.",
			},
		},
	}
}

type systemInformationDataSourceData struct {
	ID            types.String `tfsdk:"id"`
	HardwareModel types.String `tfsdk:"hardware_model"`
	OSName        types.String `tfsdk:"os_name"`
	OSVersion     types.String `tfsdk:"os_version"`
	SerialNumber  types.String `tfsdk:"serial_number"`
	ClusterNode   types.Bool   `tfsdk:"cluster_node"`

	hostName string `tfsdk:"-"`
}

func (dsc *systemInformationDataSource) Read(
	ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var data systemInformationDataSourceData

	var _ dataSourceDataReadWithoutArg = &data
	defaultDataSourceRead(
		ctx,
		dsc,
		nil,
		&data,
		resp,
	)
}

func (dscData *systemInformationDataSourceData) fillID() {
	if v := dscData.hostName; v != "" {
		dscData.ID = types.StringValue(v)
	} else {
		dscData.ID = types.StringValue("Null-Hostname")
	}
}

func (dscData *systemInformationDataSourceData) read(
	_ context.Context, junSess *junos.Session,
) error {
	dscData.hostName = junSess.SystemInformation.HostName
	dscData.HardwareModel = types.StringValue(junSess.SystemInformation.HardwareModel)
	dscData.OSName = types.StringValue(junSess.SystemInformation.OSName)
	dscData.OSVersion = types.StringValue(junSess.SystemInformation.OSVersion)
	dscData.SerialNumber = types.StringValue(junSess.SystemInformation.SerialNumber)

	// Pointer will be nil if the tag does not exist
	if junSess.SystemInformation.ClusterNode != nil {
		dscData.ClusterNode = types.BoolValue(true)
	} else {
		dscData.ClusterNode = types.BoolValue(false)
	}

	return nil
}
