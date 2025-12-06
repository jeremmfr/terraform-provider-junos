package provider

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &configRawDataSource{}
	_ datasource.DataSourceWithConfigure = &configRawDataSource{}
)

type configRawDataSource struct {
	client *junos.Client
}

func (dsc *configRawDataSource) typeName() string {
	return providerName + "_config_raw"
}

func (dsc *configRawDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newConfigRawDataSource() datasource.DataSource {
	return &configRawDataSource{}
}

func (dsc *configRawDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *configRawDataSource) Configure(
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

func (dsc *configRawDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get raw configuration from the Junos device in the specified format.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the data source with format `<format>`.",
			},
			"format": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Configuration format. Defaults to 'text'.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						junos.ConfigFormatJSON,
						junos.ConfigFormatJSONMinified,
						junos.ConfigFormatSet,
						junos.ConfigFormatText,
						junos.ConfigFormatXML,
						junos.ConfigFormatXMLMinified,
					),
				},
			},
			"config": schema.StringAttribute{
				Computed:    true,
				Description: "The raw configuration output in the requested format.",
			},
		},
	}
}

type configRawDataSourceData struct {
	ID     types.String `tfsdk:"id"`
	Format types.String `tfsdk:"format"`
	Config types.String `tfsdk:"config"`
}

func (dsc *configRawDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var data configRawDataSourceData

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _ dataSourceDataReadWithoutArg = &data
	defaultDataSourceRead(
		ctx,
		dsc,
		nil,
		&data,
		resp,
	)
}

func (dscData *configRawDataSourceData) fillID() {
	if v := dscData.Format.ValueString(); v != "" {
		dscData.ID = types.StringValue(dscData.Format.ValueString())
	} else {
		dscData.ID = types.StringValue(junos.ConfigFormatText)
	}
}

func (dscData *configRawDataSourceData) read(
	_ context.Context, junSess *junos.Session,
) error {
	if v := dscData.Format.ValueString(); v == "" {
		dscData.Format = types.StringValue(junos.ConfigFormatText)
	}

	config, err := junSess.ConfigGet(dscData.Format.ValueString())
	if err != nil {
		return fmt.Errorf("getting configuration: %w", err)
	}

	dscData.Config = types.StringValue(config)

	return nil
}
