package provider

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
	_ datasource.DataSource              = &policyoptionsASPathDataSource{}
	_ datasource.DataSourceWithConfigure = &policyoptionsASPathDataSource{}
)

type policyoptionsASPathDataSource struct {
	client *junos.Client
}

func (dsc *policyoptionsASPathDataSource) typeName() string {
	return providerName + "_policyoptions_as_path"
}

func (dsc *policyoptionsASPathDataSource) junosName() string {
	return "policy-options as-path"
}

func (dsc *policyoptionsASPathDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newPolicyoptionsASPathDataSource() datasource.DataSource {
	return &policyoptionsASPathDataSource{}
}

func (dsc *policyoptionsASPathDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *policyoptionsASPathDataSource) Configure(
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

func (dsc *policyoptionsASPathDataSource) Schema(
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
				Description: "Name to identify AS path regular expression.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"dynamic_db": schema.BoolAttribute{
				Computed:    true,
				Description: "Object may exist in dynamic database.",
			},
			"path": schema.StringAttribute{
				Computed:    true,
				Description: "AS path regular expression.",
			},
		},
	}
}

type policyoptionsASPathDataSourceData struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	DynamicDB types.Bool   `tfsdk:"dynamic_db"`
	Path      types.String `tfsdk:"path"`
}

func (dsc *policyoptionsASPathDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var name types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data policyoptionsASPathDataSourceData
	var rscData policyoptionsASPathData

	var _ resourceDataReadFrom1String = &rscData
	defaultDataSourceReadFromResource(
		ctx,
		dsc,
		[]string{
			name.ValueString(),
		},
		&data,
		&rscData,
		resp,
		fmt.Sprintf(dsc.junosName()+" %q doesn't exist", name.ValueString()),
	)
}

func (dscData *policyoptionsASPathDataSourceData) copyFromResourceData(data any) {
	rscData := data.(*policyoptionsASPathData)
	dscData.ID = rscData.ID
	dscData.Name = rscData.Name
	dscData.DynamicDB = rscData.DynamicDB
	dscData.Path = rscData.Path
}
