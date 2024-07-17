package providerfwk

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &routesDataSource{}
	_ datasource.DataSourceWithConfigure = &routesDataSource{}
)

type routesDataSource struct {
	client *junos.Client
}

func (dsc *routesDataSource) typeName() string {
	return providerName + "_routes"
}

func (dsc *routesDataSource) junosName() string {
	return "present routes"
}

func (dsc *routesDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newRoutesDataSource() datasource.DataSource {
	return &routesDataSource{}
}

func (dsc *routesDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *routesDataSource) Configure(
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

func (dsc *routesDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get " + dsc.junosName() + " in table(s).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the data source.",
			},
			"table_name": schema.StringAttribute{
				Optional:    true,
				Description: "Get routes only on a specific routing table with the name.",
			},
			"table": schema.ListAttribute{
				Computed:    true,
				Description: "For each routing table.",
				ElementType: types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
					"name": types.StringType,
					"route": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"destination": types.StringType,
						"entry": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
							"as_path":          types.StringType,
							"current_active":   types.BoolType,
							"local_preference": types.Int64Type,
							"metric":           types.Int64Type,
							"next_hop": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
								"local_interface":   types.StringType,
								"selected_next_hop": types.BoolType,
								"to":                types.StringType,
								"via":               types.StringType,
							})),
							"next_hop_type": types.StringType,
							"preference":    types.Int64Type,
							"protocol":      types.StringType,
						})),
					})),
				}),
			},
		},
	}
}

type routesDataSourceData struct {
	ID        types.String                 `tfsdk:"id"`
	TableName types.String                 `tfsdk:"table_name"`
	Table     []routesDataSourceBlockTable `tfsdk:"table"`
}

type routesDataSourceBlockTable struct {
	Name  types.String                           `tfsdk:"name"`
	Route []routesDataSourceBlockTableBlockRoute `tfsdk:"route"`
}

type routesDataSourceBlockTableBlockRoute struct {
	Destination types.String                                     `tfsdk:"destination"`
	Entry       []routesDataSourceBlockTableBlockRouteBlockEntry `tfsdk:"entry"`
}

type routesDataSourceBlockTableBlockRouteBlockEntry struct {
	ASPath          types.String                                                 `tfsdk:"as_path"`
	CurrentActive   types.Bool                                                   `tfsdk:"current_active"`
	LocalPreference types.Int64                                                  `tfsdk:"local_preference"`
	Metric          types.Int64                                                  `tfsdk:"metric"`
	NextHop         []routesDataSourceBlockTableBlockRouteBlockEntryBlockNextHop `tfsdk:"next_hop"`
	NextHopType     types.String                                                 `tfsdk:"next_hop_type"`
	Preference      types.Int64                                                  `tfsdk:"preference"`
	Protocol        types.String                                                 `tfsdk:"protocol"`
}

type routesDataSourceBlockTableBlockRouteBlockEntryBlockNextHop struct {
	LocalInterface  types.String `tfsdk:"local_interface"`
	SelectedNextHop types.Bool   `tfsdk:"selected_next_hop"`
	To              types.String `tfsdk:"to"`
	Via             types.String `tfsdk:"via"`
}

func (dsc *routesDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var tableName types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("table_name"), &tableName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data routesDataSourceData
	data.TableName = tableName

	var _ dataSourceDataReadWith1String = &data
	defaultDataSourceRead(
		ctx,
		dsc,
		[]any{
			tableName.ValueString(),
		},
		&data,
		resp,
	)
}

func (dscData *routesDataSourceData) fillID() {
	if v := dscData.TableName.ValueString(); v != "" {
		dscData.ID = types.StringValue(v)
	} else {
		dscData.ID = types.StringValue("all")
	}
}

func (dscData *routesDataSourceData) read(
	_ context.Context, tableName string, junSess *junos.Session,
) error {
	rpcReq := junos.RPCGetRouteAllInformation
	if tableName != "" {
		rpcReq = fmt.Sprintf(junos.RPCGetRouteAllTableInformation, tableName)
	}
	replyData, err := junSess.CommandXML(rpcReq)
	if err != nil {
		return err
	}
	var reply junos.RPCGetRouteInformationReply
	err = xml.Unmarshal([]byte(replyData), &reply)
	if err != nil {
		return fmt.Errorf("unmarshaling xml reply '%s': %w", replyData, err)
	}

	for _, tableInfo := range reply.RouteTable {
		table := routesDataSourceBlockTable{
			Name: types.StringValue(tableInfo.TableName),
		}
		for _, routeInfo := range tableInfo.Route {
			route := routesDataSourceBlockTableBlockRoute{
				Destination: types.StringValue(routeInfo.Destination),
			}
			for _, entryInfo := range routeInfo.Entry {
				entry := routesDataSourceBlockTableBlockRouteBlockEntry{
					CurrentActive: types.BoolValue(entryInfo.CurrentActive != nil),
				}
				if entryInfo.ASPath != nil {
					entry.ASPath = types.StringValue(strings.Trim(*entryInfo.ASPath, "\n"))
				}
				if entryInfo.LocalPreference != nil {
					entry.LocalPreference = types.Int64Value(int64(*entryInfo.LocalPreference))
				}
				if entryInfo.Metric != nil {
					entry.Metric = types.Int64Value(int64(*entryInfo.Metric))
				}
				if entryInfo.NextHopType != nil {
					entry.NextHopType = types.StringValue(*entryInfo.NextHopType)
				}
				if entryInfo.Preference != nil {
					entry.Preference = types.Int64Value(int64(*entryInfo.Preference))
				}
				if entryInfo.Protocol != nil {
					entry.Protocol = types.StringValue(*entryInfo.Protocol)
				}
				for _, nextHopInfo := range entryInfo.NextHop {
					nextHop := routesDataSourceBlockTableBlockRouteBlockEntryBlockNextHop{
						SelectedNextHop: types.BoolValue(nextHopInfo.SelectedNextHop != nil),
					}
					if nextHopInfo.LocalInterface != nil {
						nextHop.LocalInterface = types.StringValue(*nextHopInfo.LocalInterface)
					}
					if nextHopInfo.To != nil {
						nextHop.To = types.StringValue(*nextHopInfo.To)
					}
					if nextHopInfo.Via != nil {
						nextHop.Via = types.StringValue(*nextHopInfo.Via)
					}
					entry.NextHop = append(entry.NextHop, nextHop)
				}
				route.Entry = append(route.Entry, entry)
			}
			table.Route = append(table.Route, route)
		}
		dscData.Table = append(dscData.Table, table)
	}

	return nil
}
