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
	_ datasource.DataSource              = &policyoptionsDampingDataSource{}
	_ datasource.DataSourceWithConfigure = &policyoptionsDampingDataSource{}
)

type policyoptionsDampingDataSource struct {
	client *junos.Client
}

func (dsc *policyoptionsDampingDataSource) typeName() string {
	return providerName + "_policyoptions_damping"
}

func (dsc *policyoptionsDampingDataSource) junosName() string {
	return "policy-options damping"
}

func (dsc *policyoptionsDampingDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newPolicyoptionsDampingDataSource() datasource.DataSource {
	return &policyoptionsDampingDataSource{}
}

func (dsc *policyoptionsDampingDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *policyoptionsDampingDataSource) Configure(
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

func (dsc *policyoptionsDampingDataSource) Schema(
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
				Description: "Name to identify route flap damping parameters.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 250),
					tfvalidator.StringDoubleQuoteExclusion(),
				},
			},
			"disable": schema.BoolAttribute{
				Computed:    true,
				Description: "Disable damping.",
			},
			"half_life": schema.Int64Attribute{
				Computed:    true,
				Description: "Decay half-life (minutes).",
			},
			"max_suppress": schema.Int64Attribute{
				Computed:    true,
				Description: "Maximum hold-down time (minutes).",
			},
			"reuse": schema.Int64Attribute{
				Computed:    true,
				Description: "Reuse threshold.",
			},
			"suppress": schema.Int64Attribute{
				Computed:    true,
				Description: "Cutoff threshold.",
			},
		},
	}
}

type policyoptionsDampingDataSourceData struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Disable     types.Bool   `tfsdk:"disable"`
	HalfLife    types.Int64  `tfsdk:"half_life"`
	MaxSuppress types.Int64  `tfsdk:"max_suppress"`
	Reuse       types.Int64  `tfsdk:"reuse"`
	Suppress    types.Int64  `tfsdk:"suppress"`
}

func (dsc *policyoptionsDampingDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var name types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data policyoptionsDampingDataSourceData
	var rscData policyoptionsDampingData

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

func (dscData *policyoptionsDampingDataSourceData) copyFromResourceData(data any) {
	rscData := data.(*policyoptionsDampingData)
	dscData.ID = rscData.ID
	dscData.Name = rscData.Name
	dscData.Disable = rscData.Disable
	dscData.HalfLife = rscData.HalfLife
	dscData.MaxSuppress = rscData.MaxSuppress
	dscData.Reuse = rscData.Reuse
	dscData.Suppress = rscData.Suppress
}
