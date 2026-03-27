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
	_ datasource.DataSource              = &policyoptionsPrefixListDataSource{}
	_ datasource.DataSourceWithConfigure = &policyoptionsPrefixListDataSource{}
)

type policyoptionsPrefixListDataSource struct {
	client *junos.Client
}

func (dsc *policyoptionsPrefixListDataSource) typeName() string {
	return providerName + "_policyoptions_prefix_list"
}

func (dsc *policyoptionsPrefixListDataSource) junosName() string {
	return "policy-options prefix-list"
}

func (dsc *policyoptionsPrefixListDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newPolicyoptionsPrefixListDataSource() datasource.DataSource {
	return &policyoptionsPrefixListDataSource{}
}

func (dsc *policyoptionsPrefixListDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *policyoptionsPrefixListDataSource) Configure(
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

func (dsc *policyoptionsPrefixListDataSource) Schema(
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
				Description: "Prefix list name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"apply_path": schema.StringAttribute{
				Computed:    true,
				Description: "Apply IP prefixes from a configuration statement.",
			},
			"dynamic_db": schema.BoolAttribute{
				Computed:    true,
				Description: "Object may exist in dynamic database.",
			},
			"prefix": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Address prefixes.",
			},
		},
	}
}

type policyoptionsPrefixListDataSourceData struct {
	ID        types.String   `tfsdk:"id"`
	Name      types.String   `tfsdk:"name"`
	ApplyPath types.String   `tfsdk:"apply_path"`
	DynamicDB types.Bool     `tfsdk:"dynamic_db"`
	Prefix    []types.String `tfsdk:"prefix"`
}

func (dsc *policyoptionsPrefixListDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var name types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data policyoptionsPrefixListDataSourceData
	var rscData policyoptionsPrefixListData

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

func (dscData *policyoptionsPrefixListDataSourceData) copyFromResourceData(data any) {
	rscData := data.(*policyoptionsPrefixListData)
	dscData.ID = rscData.ID
	dscData.Name = rscData.Name
	dscData.ApplyPath = rscData.ApplyPath
	dscData.DynamicDB = rscData.DynamicDB
	dscData.Prefix = rscData.Prefix
}
