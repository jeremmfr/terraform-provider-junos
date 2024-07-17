package providerfwk

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

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
	_ datasource.DataSource              = &interfaceLogicalInfoDataSource{}
	_ datasource.DataSourceWithConfigure = &interfaceLogicalInfoDataSource{}
)

type interfaceLogicalInfoDataSource struct {
	client *junos.Client
}

func (dsc *interfaceLogicalInfoDataSource) typeName() string {
	return providerName + "_interface_logical_info"
}

func (dsc *interfaceLogicalInfoDataSource) junosName() string {
	return "summary information about a logical interface"
}

func (dsc *interfaceLogicalInfoDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newInterfaceLogicalInfoDataSource() datasource.DataSource {
	return &interfaceLogicalInfoDataSource{}
}

func (dsc *interfaceLogicalInfoDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *interfaceLogicalInfoDataSource) Configure(
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

func (dsc *interfaceLogicalInfoDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get " + dsc.junosName() + ".",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The name of interface read.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of logical interface (with dot).",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					tfvalidator.StringFormat(tfvalidator.InterfaceFormat),
					tfvalidator.String1DotCount(),
				},
			},
			"admin_status": schema.StringAttribute{
				Computed:    true,
				Description: "Admin status.",
			},
			"oper_status": schema.StringAttribute{
				Computed:    true,
				Description: "Operational status.",
			},
			"family_inet": schema.ObjectAttribute{
				Computed:    true,
				Description: "Family inet enabled.",
				AttributeTypes: map[string]attr.Type{
					"address_cidr": types.ListType{}.WithElementType(types.StringType),
				},
			},
			"family_inet6": schema.ObjectAttribute{
				Computed:    true,
				Description: "Family inet6 enabled.",
				AttributeTypes: map[string]attr.Type{
					"address_cidr": types.ListType{}.WithElementType(types.StringType),
				},
			},
		},
	}
}

type interfaceLogicalInfoDataSourceeData struct {
	ID          types.String                                    `tfsdk:"id"`
	Name        types.String                                    `tfsdk:"name"`
	AdminStatus types.String                                    `tfsdk:"admin_status"`
	OperStatus  types.String                                    `tfsdk:"oper_status"`
	FamilyInet  *interfaceLogicalInfoDataSourceeBlockFamilyInet `tfsdk:"family_inet"`
	FamilyInet6 *interfaceLogicalInfoDataSourceeBlockFamilyInet `tfsdk:"family_inet6"`
}

type interfaceLogicalInfoDataSourceeBlockFamilyInet struct {
	AddressCIDR []types.String `tfsdk:"address_cidr"`
}

func (dsc *interfaceLogicalInfoDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var name types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data interfaceLogicalInfoDataSourceeData

	var _ dataSourceDataReadWith1String = &data
	defaultDataSourceRead(
		ctx,
		dsc,
		[]any{
			name.ValueString(),
		},
		&data,
		resp,
	)
}

func (dscData *interfaceLogicalInfoDataSourceeData) fillID() {
	dscData.ID = types.StringValue(dscData.Name.ValueString())
}

func (dscData *interfaceLogicalInfoDataSourceeData) read(
	_ context.Context, name string, junSess *junos.Session,
) error {
	replyData, err := junSess.CommandXML(fmt.Sprintf(junos.RPCGetInterfaceInformationTerse, name))
	if err != nil {
		return err
	}
	var reply junos.RPCGetLogicalInterfaceTerseReply
	err = xml.Unmarshal([]byte(replyData), &reply)
	if err != nil {
		return fmt.Errorf("unmarshaling xml reply %q: %w", replyData, err)
	}
	if len(reply.LogicalInterface) == 0 {
		return fmt.Errorf("logical-interface not found in xml: %v", replyData)
	}
	ifaceInfo := reply.LogicalInterface[0]
	dscData.Name = types.StringValue(strings.TrimSpace(ifaceInfo.Name))
	dscData.AdminStatus = types.StringValue(strings.TrimSpace(ifaceInfo.AdminStatus))
	dscData.OperStatus = types.StringValue(strings.TrimSpace(ifaceInfo.OperStatus))
	for _, family := range ifaceInfo.AddressFamily {
		switch strings.TrimSpace(family.Name) {
		case junos.InetW:
			if dscData.FamilyInet == nil {
				dscData.FamilyInet = &interfaceLogicalInfoDataSourceeBlockFamilyInet{}
			}
			for _, address := range family.Address {
				dscData.FamilyInet.AddressCIDR = append(
					dscData.FamilyInet.AddressCIDR,
					types.StringValue(strings.TrimSpace(address.Local)),
				)
			}
		case junos.Inet6W:
			if dscData.FamilyInet6 == nil {
				dscData.FamilyInet6 = &interfaceLogicalInfoDataSourceeBlockFamilyInet{}
			}
			for _, address := range family.Address {
				dscData.FamilyInet6.AddressCIDR = append(
					dscData.FamilyInet6.AddressCIDR,
					types.StringValue(strings.TrimSpace(address.Local)),
				)
			}
		}
	}

	return nil
}
