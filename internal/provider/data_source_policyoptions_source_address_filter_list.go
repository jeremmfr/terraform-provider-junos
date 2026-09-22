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
	_ datasource.DataSource              = &policyoptionsSourceAddressFilterListDataSource{}
	_ datasource.DataSourceWithConfigure = &policyoptionsSourceAddressFilterListDataSource{}
)

type policyoptionsSourceAddressFilterListDataSource struct {
	client *junos.Client
}

func (dsc *policyoptionsSourceAddressFilterListDataSource) typeName() string {
	return providerName + "_policyoptions_source_address_filter_list"
}

func (dsc *policyoptionsSourceAddressFilterListDataSource) junosName() string {
	return "policy-options source-address-filter-list"
}

func (dsc *policyoptionsSourceAddressFilterListDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newPolicyoptionsSourceAddressFilterListDataSource() datasource.DataSource {
	return &policyoptionsSourceAddressFilterListDataSource{}
}

func (dsc *policyoptionsSourceAddressFilterListDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *policyoptionsSourceAddressFilterListDataSource) Configure(
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

func (dsc *policyoptionsSourceAddressFilterListDataSource) Schema(
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
				Description: "Source address filter list name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"dynamic_db": schema.BoolAttribute{
				Computed:    true,
				Description: "Object may exist in dynamic database.",
			},
		},
		Blocks: map[string]schema.Block{
			"address": schema.SetNestedBlock{
				Description: "List of source addresses.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Computed:    true,
							Description: "IP address.",
						},
						"option": schema.StringAttribute{
							Computed:    true,
							Description: "Mask option.",
						},
						"option_value": schema.StringAttribute{
							Computed:    true,
							Description: "For options that need an argument.",
						},
					},
				},
			},
		},
	}
}

type policyoptionsSourceAddressFilterListDataSourceData struct {
	ID        types.String                                       `tfsdk:"id"`
	Name      types.String                                       `tfsdk:"name"`
	DynamicDB types.Bool                                         `tfsdk:"dynamic_db"`
	Address   []policyoptionsSourceAddressFilterListBlockAddress `tfsdk:"address"`
}

func (dsc *policyoptionsSourceAddressFilterListDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var name types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data policyoptionsSourceAddressFilterListDataSourceData
	var rscData policyoptionsSourceAddressFilterListData

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

func (dscData *policyoptionsSourceAddressFilterListDataSourceData) copyFromResourceData(data any) {
	rscData := data.(*policyoptionsSourceAddressFilterListData)
	dscData.ID = rscData.ID
	dscData.Name = rscData.Name
	dscData.DynamicDB = rscData.DynamicDB
	dscData.Address = rscData.Address
}
