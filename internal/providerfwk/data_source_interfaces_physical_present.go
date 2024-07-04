package providerfwk

import (
	"context"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &interfacesPhysicalPresentDataSource{}
	_ datasource.DataSourceWithConfigure = &interfacesPhysicalPresentDataSource{}
)

type interfacesPhysicalPresentDataSource struct {
	client *junos.Client
}

func (dsc *interfacesPhysicalPresentDataSource) typeName() string {
	return providerName + "_interfaces_physical_present"
}

func (dsc *interfacesPhysicalPresentDataSource) junosName() string {
	return "physical interfaces present on device"
}

func (dsc *interfacesPhysicalPresentDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newInterfacesPhysicalPresentDataSource() datasource.DataSource {
	return &interfacesPhysicalPresentDataSource{}
}

func (dsc *interfacesPhysicalPresentDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *interfacesPhysicalPresentDataSource) Configure(
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

func (dsc *interfacesPhysicalPresentDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get list of all of filtered " + dsc.junosName() + ".",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the data source.",
			},
			"match_name": schema.StringAttribute{
				Optional:    true,
				Description: "A regexp to apply filter on name.",
				Validators: []validator.String{
					tfvalidator.StringRegex(),
				},
			},
			"match_admin_up": schema.BoolAttribute{
				Optional:    true,
				Description: "Filter on interfaces that have admin status `up`.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"match_oper_up": schema.BoolAttribute{
				Optional:    true,
				Description: "Filter on interfaces that have operational status `up`.",
				Validators: []validator.Bool{
					tfvalidator.BoolTrue(),
				},
			},
			"interface_names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Found interface names.",
			},
			"interfaces": schema.MapAttribute{
				Computed:    true,
				Description: "Dictionary of found interfaces with interface name as key.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":                    types.StringType,
						"admin_status":            types.StringType,
						"oper_status":             types.StringType,
						"logical_interface_names": types.ListType{}.WithElementType(types.StringType),
					},
				},
			},
			"interface_statuses": schema.ListAttribute{
				Computed:           true,
				DeprecationMessage: "Use the \"interfaces\" attribute instead.",
				Description:        "For each found interface name, its status.",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":         types.StringType,
						"admin_status": types.StringType,
						"oper_status":  types.StringType,
					},
				},
			},
		},
	}
}

type interfacesPhysicalPresentDataSourceData struct {
	ID                types.String                                                  `tfsdk:"id"`
	MatchName         types.String                                                  `tfsdk:"match_name"`
	MatchAdminUp      types.Bool                                                    `tfsdk:"match_admin_up"`
	MatchOperUp       types.Bool                                                    `tfsdk:"match_oper_up"`
	InterfaceNames    []types.String                                                `tfsdk:"interface_names"`
	Interfaces        map[string]interfacesPhysicalPresentDataSourceBlockInterfaces `tfsdk:"interfaces"`
	InterfaceStatuses []interfacesPhysicalPresentDataSourceBlockInterfaceStatuses   `tfsdk:"interface_statuses"`
}

type interfacesPhysicalPresentDataSourceConfig struct {
	ID                types.String `tfsdk:"id"`
	MatchName         types.String `tfsdk:"match_name"`
	MatchAdminUp      types.Bool   `tfsdk:"match_admin_up"`
	MatchOperUp       types.Bool   `tfsdk:"match_oper_up"`
	InterfaceNames    types.List   `tfsdk:"interface_names"`
	Interfaces        types.Map    `tfsdk:"interfaces"`
	InterfaceStatuses types.List   `tfsdk:"interface_statuses"`
}

type interfacesPhysicalPresentDataSourceBlockInterfaces struct {
	Name                  types.String   `tfsdk:"name"`
	AdminStatus           types.String   `tfsdk:"admin_status"`
	OperStatus            types.String   `tfsdk:"oper_status"`
	LogicalInterfaceNames []types.String `tfsdk:"logical_interface_names"`
}

type interfacesPhysicalPresentDataSourceBlockInterfaceStatuses struct {
	Name        types.String `tfsdk:"name"`
	AdminStatus types.String `tfsdk:"admin_status"`
	OperStatus  types.String `tfsdk:"oper_status"`
}

func (dsc *interfacesPhysicalPresentDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var config interfacesPhysicalPresentDataSourceConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data interfacesPhysicalPresentDataSourceData
	data.MatchName = config.MatchName
	data.MatchAdminUp = config.MatchAdminUp
	data.MatchOperUp = config.MatchOperUp

	var _ dataSourceDataReadWith1String2Bool = &data
	defaultDataSourceRead(
		ctx,
		dsc,
		[]any{
			config.MatchName.ValueString(),
			config.MatchAdminUp.ValueBool(),
			config.MatchOperUp.ValueBool(),
		},
		&data,
		resp,
	)
}

func (dscData *interfacesPhysicalPresentDataSourceData) fillID() {
	idString := "match=" + dscData.MatchName.ValueString()
	if dscData.MatchAdminUp.ValueBool() {
		idString += junos.IDSeparator + "admin_up=true"
	}
	if dscData.MatchOperUp.ValueBool() {
		idString += junos.IDSeparator + "oper_up=true"
	}
	dscData.ID = types.StringValue(idString)
}

func (dscData *interfacesPhysicalPresentDataSourceData) read(
	_ context.Context,
	matchName string,
	matchAdminUp bool,
	matchOperUp bool,
	junSess *junos.Session,
) error {
	replyData, err := junSess.CommandXML(junos.RPCGetInterfacesInformationTerse)
	if err != nil {
		return err
	}
	var reply junos.RPCGetPhysicalInterfaceTerseReply
	err = xml.Unmarshal([]byte(replyData), &reply)
	if err != nil {
		return fmt.Errorf("unmarshaling xml reply %q: %w", replyData, err)
	}

	dscData.InterfaceNames = make([]types.String, 0, len(reply.PhysicalInterface))
	dscData.Interfaces = make(map[string]interfacesPhysicalPresentDataSourceBlockInterfaces)
	dscData.InterfaceStatuses = make(
		[]interfacesPhysicalPresentDataSourceBlockInterfaceStatuses, 0, len(reply.PhysicalInterface),
	)
	for _, iFace := range reply.PhysicalInterface {
		if matchName != "" {
			matched, err := regexp.MatchString(matchName, strings.TrimSpace(iFace.Name))
			if err != nil {
				return fmt.Errorf("matching with regexp %q: %w", matchName, err)
			}
			if !matched {
				continue
			}
		}
		if matchAdminUp && strings.TrimSpace(iFace.AdminStatus) != "up" {
			continue
		}
		if matchOperUp && strings.TrimSpace(iFace.OperStatus) != "up" {
			continue
		}
		dscData.InterfaceNames = append(dscData.InterfaceNames, types.StringValue(strings.TrimSpace(iFace.Name)))
		logicalInterfaceNames := make([]types.String, len(iFace.LogicalInterface))
		for i, loIface := range iFace.LogicalInterface {
			logicalInterfaceNames[i] = types.StringValue(strings.TrimSpace(loIface.Name))
		}
		dscData.Interfaces[strings.TrimSpace(iFace.Name)] = interfacesPhysicalPresentDataSourceBlockInterfaces{
			Name:                  types.StringValue(strings.TrimSpace(iFace.Name)),
			AdminStatus:           types.StringValue(strings.TrimSpace(iFace.AdminStatus)),
			OperStatus:            types.StringValue(strings.TrimSpace(iFace.OperStatus)),
			LogicalInterfaceNames: logicalInterfaceNames,
		}
		dscData.InterfaceStatuses = append(dscData.InterfaceStatuses,
			interfacesPhysicalPresentDataSourceBlockInterfaceStatuses{
				Name:        types.StringValue(strings.TrimSpace(iFace.Name)),
				AdminStatus: types.StringValue(strings.TrimSpace(iFace.AdminStatus)),
				OperStatus:  types.StringValue(strings.TrimSpace(iFace.OperStatus)),
			})
	}

	return nil
}
