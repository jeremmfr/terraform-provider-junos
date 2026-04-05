package provider

import (
	"context"
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &policyoptionsASPathGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &policyoptionsASPathGroupDataSource{}
)

type policyoptionsASPathGroupDataSource struct {
	client *junos.Client
}

func (dsc *policyoptionsASPathGroupDataSource) typeName() string {
	return providerName + "_policyoptions_as_path_group"
}

func (dsc *policyoptionsASPathGroupDataSource) junosName() string {
	return "policy-options as-path-group"
}

func (dsc *policyoptionsASPathGroupDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newPolicyoptionsASPathGroupDataSource() datasource.DataSource {
	return &policyoptionsASPathGroupDataSource{}
}

func (dsc *policyoptionsASPathGroupDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *policyoptionsASPathGroupDataSource) Configure(
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

func (dsc *policyoptionsASPathGroupDataSource) Schema(
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
				Description: "Name to identify AS path group.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"dynamic_db": schema.BoolAttribute{
				Computed:    true,
				Description: "Object may exist in dynamic database.",
			},
			"as_path": schema.ListAttribute{
				Computed:    true,
				Description: "List of as-path entries in this group.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name": types.StringType,
						"path": types.StringType,
					},
				},
			},
		},
	}
}

type policyoptionsASPathGroupDataSourceData struct {
	ID        types.String                          `tfsdk:"id"`
	Name      types.String                          `tfsdk:"name"`
	DynamicDB types.Bool                            `tfsdk:"dynamic_db"`
	ASPath    []policyoptionsASPathGroupBlockASPAth `tfsdk:"as_path"`
}

func (dsc *policyoptionsASPathGroupDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var name types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data policyoptionsASPathGroupDataSourceData
	var rscData policyoptionsASPathGroupData

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

func (dscData *policyoptionsASPathGroupDataSourceData) copyFromResourceData(data any) {
	rscData := data.(*policyoptionsASPathGroupData)
	dscData.ID = rscData.ID
	dscData.Name = rscData.Name
	dscData.DynamicDB = rscData.DynamicDB
	dscData.ASPath = rscData.ASPath
}
