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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &chassisInventoryDataSource{}
	_ datasource.DataSourceWithConfigure = &chassisInventoryDataSource{}
)

type chassisInventoryDataSource struct {
	client *junos.Client
}

func (dsc *chassisInventoryDataSource) typeName() string {
	return providerName + "_chassis_inventory"
}

func (dsc *chassisInventoryDataSource) junosName() string {
	return "chassis hardware"
}

func (dsc *chassisInventoryDataSource) junosClient() *junos.Client {
	return dsc.client
}

func newChassisInventoryDataSource() datasource.DataSource {
	return &chassisInventoryDataSource{}
}

func (dsc *chassisInventoryDataSource) Metadata(
	_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = dsc.typeName()
}

func (dsc *chassisInventoryDataSource) Configure(
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

func (dsc *chassisInventoryDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Get chassis inventory (" + dsc.junosName() + ")",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "An identifier for the data source with value `chassis_inventory`.",
			},
			"chassis": schema.ListAttribute{
				Computed:    true,
				Description: "Chassis inventory for each routing engine.",
				ElementType: types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
					"re_name":       types.StringType,
					"clei_code":     types.StringType,
					"description":   types.StringType,
					"model_number":  types.StringType,
					"name":          types.StringType,
					"part_number":   types.StringType,
					"serial_number": types.StringType,
					"version":       types.StringType,
					"module": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
						"clei_code":     types.StringType,
						"description":   types.StringType,
						"model_number":  types.StringType,
						"name":          types.StringType,
						"part_number":   types.StringType,
						"serial_number": types.StringType,
						"version":       types.StringType,
						"sub_module": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
							"clei_code":     types.StringType,
							"description":   types.StringType,
							"model_number":  types.StringType,
							"name":          types.StringType,
							"part_number":   types.StringType,
							"serial_number": types.StringType,
							"version":       types.StringType,
							"sub_sub_module": types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
								"clei_code":     types.StringType,
								"description":   types.StringType,
								"model_number":  types.StringType,
								"name":          types.StringType,
								"part_number":   types.StringType,
								"serial_number": types.StringType,
								"version":       types.StringType,
							})),
						})),
					})),
				}),
			},
		},
	}
}

type chassisInventoryDataSourceData struct {
	ID      types.String                             `tfsdk:"id"`
	Chassis []chassisInventoryDataSourceBlockChassis `tfsdk:"chassis"`
}

type chassisInventoryDataSourceBlockChassis struct {
	ReName       types.String                                        `tfsdk:"re_name"`
	CleiCode     types.String                                        `tfsdk:"clei_code"`
	Description  types.String                                        `tfsdk:"description"`
	ModelNumber  types.String                                        `tfsdk:"model_number"`
	Name         types.String                                        `tfsdk:"name"`
	PartNumber   types.String                                        `tfsdk:"part_number"`
	SerialNumber types.String                                        `tfsdk:"serial_number"`
	Version      types.String                                        `tfsdk:"version"`
	Module       []chassisInventoryDataSourceBlockChassisBlockModule `tfsdk:"module"`
}

type chassisInventoryDataSourceBlockChassisBlockModule struct {
	CleiCode     types.String                                                      `tfsdk:"clei_code"`
	Description  types.String                                                      `tfsdk:"description"`
	ModelNumber  types.String                                                      `tfsdk:"model_number"`
	Name         types.String                                                      `tfsdk:"name"`
	PartNumber   types.String                                                      `tfsdk:"part_number"`
	SerialNumber types.String                                                      `tfsdk:"serial_number"`
	Version      types.String                                                      `tfsdk:"version"`
	SubModule    []chassisInventoryDataSourceBlockChassisBlockModuleBlockSubModule `tfsdk:"sub_module"`
}

//nolint:lll
type chassisInventoryDataSourceBlockChassisBlockModuleBlockSubModule struct {
	CleiCode     types.String                                                                       `tfsdk:"clei_code"`
	Description  types.String                                                                       `tfsdk:"description"`
	ModelNumber  types.String                                                                       `tfsdk:"model_number"`
	Name         types.String                                                                       `tfsdk:"name"`
	PartNumber   types.String                                                                       `tfsdk:"part_number"`
	SerialNumber types.String                                                                       `tfsdk:"serial_number"`
	Version      types.String                                                                       `tfsdk:"version"`
	SubSubModule []chassisInventoryDataSourceBlockChassisBlockModuleBlockSubModuleBlockSubSubModule `tfsdk:"sub_sub_module"`
}

type chassisInventoryDataSourceBlockChassisBlockModuleBlockSubModuleBlockSubSubModule struct {
	CleiCode     types.String `tfsdk:"clei_code"`
	Description  types.String `tfsdk:"description"`
	ModelNumber  types.String `tfsdk:"model_number"`
	Name         types.String `tfsdk:"name"`
	PartNumber   types.String `tfsdk:"part_number"`
	SerialNumber types.String `tfsdk:"serial_number"`
	Version      types.String `tfsdk:"version"`
}

func (dsc *chassisInventoryDataSource) Read(
	ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var data chassisInventoryDataSourceData

	var _ dataSourceDataReadWithoutArg = &data
	defaultDataSourceRead(
		ctx,
		dsc,
		nil,
		&data,
		resp,
	)
}

func (dscData *chassisInventoryDataSourceData) fillID() {
	dscData.ID = types.StringValue("chassis_inventory")
}

func (dscData *chassisInventoryDataSourceData) read(
	_ context.Context, junSess *junos.Session,
) error {
	replyData, err := junSess.CommandXML(junos.RPCGetChassisInventory)
	if err != nil {
		return err
	}
	if strings.Contains(replyData, "<multi-routing-engine-results") {
		type multiChassisInventory struct {
			XMLName                xml.Name `xml:"multi-routing-engine-results"`
			MultiRoutingEngineItem []struct {
				ChassisInventory junos.RPCGetChassisInventoryReply `xml:"chassis-inventory"`
				ReName           string                            `xml:"re-name"`
			} `xml:"multi-routing-engine-item"`
		}

		var reply multiChassisInventory
		if err := xml.Unmarshal([]byte(replyData), &reply); err != nil {
			return fmt.Errorf("unmarshaling xml reply '%s': %w", replyData, err)
		}

		for i, chassis := range reply.MultiRoutingEngineItem {
			dscData.fillChassisItem(&chassis.ChassisInventory)

			dscData.Chassis[i].ReName = types.StringValue(chassis.ReName)
		}
	} else {
		var reply junos.RPCGetChassisInventoryReply
		err = xml.Unmarshal([]byte(replyData), &reply)
		if err != nil {
			return fmt.Errorf("unmarshaling xml reply '%s': %w", replyData, err)
		}

		dscData.fillChassisItem(&reply)
	}

	return nil
}

func (dscData *chassisInventoryDataSourceData) fillChassisItem(
	reply *junos.RPCGetChassisInventoryReply,
) {
	chassis := chassisInventoryDataSourceBlockChassis{}

	if reply.Chassis.CleiCode != nil {
		chassis.CleiCode = types.StringValue(*reply.Chassis.CleiCode)
	}
	if reply.Chassis.Description != nil {
		chassis.Description = types.StringValue(*reply.Chassis.Description)
	}
	if reply.Chassis.ModelNumber != nil {
		chassis.ModelNumber = types.StringValue(*reply.Chassis.ModelNumber)
	}
	if reply.Chassis.Name != nil {
		chassis.Name = types.StringValue(*reply.Chassis.Name)
	}
	if reply.Chassis.PartNumber != nil {
		chassis.PartNumber = types.StringValue(*reply.Chassis.PartNumber)
	}
	if reply.Chassis.SerialNumber != nil {
		chassis.SerialNumber = types.StringValue(*reply.Chassis.SerialNumber)
	}
	if reply.Chassis.Version != nil {
		chassis.Version = types.StringValue(*reply.Chassis.Version)
	}

	for _, moduleInfo := range reply.Chassis.Module {
		module := chassisInventoryDataSourceBlockChassisBlockModule{}
		if moduleInfo.CleiCode != nil {
			module.CleiCode = types.StringValue(*moduleInfo.CleiCode)
		}
		if moduleInfo.Description != nil {
			module.Description = types.StringValue(*moduleInfo.Description)
		}
		if moduleInfo.ModelNumber != nil {
			module.ModelNumber = types.StringValue(*moduleInfo.ModelNumber)
		}
		if moduleInfo.Name != nil {
			module.Name = types.StringValue(*moduleInfo.Name)
		}
		if moduleInfo.PartNumber != nil {
			module.PartNumber = types.StringValue(*moduleInfo.PartNumber)
		}
		if moduleInfo.SerialNumber != nil {
			module.SerialNumber = types.StringValue(*moduleInfo.SerialNumber)
		}
		if moduleInfo.Version != nil {
			module.Version = types.StringValue(*moduleInfo.Version)
		}

		for _, subModuleInfo := range moduleInfo.SubModule {
			subModule := chassisInventoryDataSourceBlockChassisBlockModuleBlockSubModule{}
			if subModuleInfo.CleiCode != nil {
				subModule.CleiCode = types.StringValue(*subModuleInfo.CleiCode)
			}
			if subModuleInfo.Description != nil {
				subModule.Description = types.StringValue(*subModuleInfo.Description)
			}
			if subModuleInfo.ModelNumber != nil {
				subModule.ModelNumber = types.StringValue(*subModuleInfo.ModelNumber)
			}
			if subModuleInfo.Name != nil {
				subModule.Name = types.StringValue(*subModuleInfo.Name)
			}
			if subModuleInfo.PartNumber != nil {
				subModule.PartNumber = types.StringValue(*subModuleInfo.PartNumber)
			}
			if subModuleInfo.SerialNumber != nil {
				subModule.SerialNumber = types.StringValue(*subModuleInfo.SerialNumber)
			}
			if subModuleInfo.Version != nil {
				subModule.Version = types.StringValue(*subModuleInfo.Version)
			}

			for _, subSubModuleInfo := range subModuleInfo.SubSubModule {
				subSubModule := chassisInventoryDataSourceBlockChassisBlockModuleBlockSubModuleBlockSubSubModule{}
				if subSubModuleInfo.CleiCode != nil {
					subSubModule.CleiCode = types.StringValue(*subSubModuleInfo.CleiCode)
				}
				if subSubModuleInfo.Description != nil {
					subSubModule.Description = types.StringValue(*subSubModuleInfo.Description)
				}
				if subSubModuleInfo.ModelNumber != nil {
					subSubModule.ModelNumber = types.StringValue(*subSubModuleInfo.ModelNumber)
				}
				if subSubModuleInfo.Name != nil {
					subSubModule.Name = types.StringValue(*subSubModuleInfo.Name)
				}
				if subSubModuleInfo.PartNumber != nil {
					subSubModule.PartNumber = types.StringValue(*subSubModuleInfo.PartNumber)
				}
				if subSubModuleInfo.SerialNumber != nil {
					subSubModule.SerialNumber = types.StringValue(*subSubModuleInfo.SerialNumber)
				}
				if subSubModuleInfo.Version != nil {
					subSubModule.Version = types.StringValue(*subSubModuleInfo.Version)
				}

				subModule.SubSubModule = append(subModule.SubSubModule, subSubModule)
			}

			module.SubModule = append(module.SubModule, subModule)
		}

		chassis.Module = append(chassis.Module, module)
	}

	dscData.Chassis = append(dscData.Chassis, chassis)
}
