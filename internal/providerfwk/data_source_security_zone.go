package providerfwk

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
	_ datasource.DataSource              = &securityZoneDataSource{}
	_ datasource.DataSourceWithConfigure = &securityZoneDataSource{}
)

type securityZoneDataSource struct {
	client *junos.Client
}

func (dsc *securityZoneDataSource) typeName() string {
	return providerName + "_security_zone"
}

func (dsc *securityZoneDataSource) junosName() string {
	return "security zone"
}

func (dsc *securityZoneDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newSecurityZoneDataSource() datasource.DataSource {
	return &securityZoneDataSource{}
}

func (dsc *securityZoneDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *securityZoneDataSource) Configure(
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

func (dsc *securityZoneDataSource) Schema(
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
				Description: "The name of security zone.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 63),
					tfvalidator.StringFormat(tfvalidator.DefaultFormat),
				},
			},
			"advance_policy_based_routing_profile": schema.StringAttribute{
				Computed:    true,
				Description: "Enable Advance Policy Based Routing on this zone with a profile.",
			},
			"application_tracking": schema.BoolAttribute{
				Computed:    true,
				Description: "Enable Application tracking support for this zone.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Text description of zone.",
			},
			"inbound_protocols": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The inbound protocols allowed.",
			},
			"inbound_services": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The inbound services allowed.",
			},
			"reverse_reroute": schema.BoolAttribute{
				Computed:    true,
				Description: "Enable Reverse route lookup when there is change in ingress interface.",
			},
			"screen": schema.StringAttribute{
				Computed:    true,
				Description: "Name of ids option object (screen) applied to the zone.",
			},
			"source_identity_log": schema.BoolAttribute{
				Computed:    true,
				Description: "Show user and group info in session log for this zone.",
			},
			"tcp_rst": schema.BoolAttribute{
				Computed:    true,
				Description: "Send RST for NON-SYN packet not matching TCP session.",
			},
			"address_book": schema.SetAttribute{
				Computed:    true,
				Description: "For each name of address.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"network":     types.StringType,
						"description": types.StringType,
					},
				},
			},
			"address_book_dns": schema.SetAttribute{
				Computed:    true,
				Description: "For each name of dns-name address.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"fqdn":        types.StringType,
						"description": types.StringType,
						"ipv4_only":   types.BoolType,
						"ipv6_only":   types.BoolType,
					},
				},
			},
			"address_book_range": schema.SetAttribute{
				Computed:    true,
				Description: "For each name of range-address.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"from":        types.StringType,
						"to":          types.StringType,
						"description": types.StringType,
					},
				},
			},
			"address_book_set": schema.SetAttribute{
				Computed:    true,
				Description: "For each name of address-set.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"address":     types.SetType{}.WithElementType(types.StringType),
						"address_set": types.SetType{}.WithElementType(types.StringType),
						"description": types.StringType,
					},
				},
			},
			"address_book_wildcard": schema.SetAttribute{
				Computed:    true,
				Description: "For each name of wildcard-address.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"network":     types.StringType,
						"description": types.StringType,
					},
				},
			},
			"interface": schema.SetAttribute{
				Computed:    true,
				Description: "List of interfaces in security-zone.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":              types.StringType,
						"inbound_protocols": types.SetType{}.WithElementType(types.StringType),
						"inbound_services":  types.SetType{}.WithElementType(types.StringType),
					},
				},
			},
		},
	}
}

type securityZoneDataSourceData struct {
	ID                               types.String                           `tfsdk:"id"`
	Name                             types.String                           `tfsdk:"name"`
	AdvancePolicyBasedRoutingProfile types.String                           `tfsdk:"advance_policy_based_routing_profile"`
	ApplicationTracking              types.Bool                             `tfsdk:"application_tracking"`
	Description                      types.String                           `tfsdk:"description"`
	Screen                           types.String                           `tfsdk:"screen"`
	InboundProtocols                 []types.String                         `tfsdk:"inbound_protocols"`
	InboundServices                  []types.String                         `tfsdk:"inbound_services"`
	ReverseReroute                   types.Bool                             `tfsdk:"reverse_reroute"`
	SourceIdentityLog                types.Bool                             `tfsdk:"source_identity_log"`
	TCPRst                           types.Bool                             `tfsdk:"tcp_rst"`
	AddressBook                      []securityZoneBlockAddressBook         `tfsdk:"address_book"`
	AddressBookDNS                   []securityZoneBlockAddressBookDNS      `tfsdk:"address_book_dns"`
	AddressBookRange                 []securityZoneBlockAddressBookRange    `tfsdk:"address_book_range"`
	AddressBookSet                   []securityZoneBlockAddressBookSet      `tfsdk:"address_book_set"`
	AddressBookWildcard              []securityZoneBlockAddressBookWildcard `tfsdk:"address_book_wildcard"`
	Interface                        []securityZoneDataSourceBlockInterface `tfsdk:"interface"`
}

type securityZoneDataSourceBlockInterface struct {
	Name             types.String   `tfsdk:"name"`
	InboundProtocols []types.String `tfsdk:"inbound_protocols"`
	InboundServices  []types.String `tfsdk:"inbound_services"`
}

func (dsc *securityZoneDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var name types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data securityZoneDataSourceData
	var rscData securityZoneData

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

func (dscData *securityZoneDataSourceData) copyFromResourceData(data any) {
	rscData := data.(*securityZoneData)
	dscData.ID = rscData.ID
	dscData.Name = rscData.Name
	dscData.AdvancePolicyBasedRoutingProfile = rscData.AdvancePolicyBasedRoutingProfile
	dscData.ApplicationTracking = rscData.ApplicationTracking
	dscData.Description = rscData.Description
	dscData.InboundProtocols = rscData.InboundProtocols
	dscData.InboundServices = rscData.InboundServices
	dscData.ReverseReroute = rscData.ReverseReroute
	dscData.Screen = rscData.Screen
	dscData.SourceIdentityLog = rscData.SourceIdentityLog
	dscData.TCPRst = rscData.TCPRst
	dscData.AddressBook = rscData.AddressBook
	dscData.AddressBookDNS = rscData.AddressBookDNS
	dscData.AddressBookRange = rscData.AddressBookRange
	dscData.AddressBookSet = rscData.AddressBookSet
	dscData.AddressBookWildcard = rscData.AddressBookWildcard
	dscData.Interface = rscData.Interface
}
