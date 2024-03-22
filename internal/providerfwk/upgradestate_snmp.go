package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *snmp) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"clean_on_destroy": schema.BoolAttribute{
						Optional: true,
					},
					"arp": schema.BoolAttribute{
						Optional: true,
					},
					"arp_host_name_resolution": schema.BoolAttribute{
						Optional: true,
					},
					"contact": schema.StringAttribute{
						Optional: true,
					},
					"description": schema.StringAttribute{
						Optional: true,
					},
					"engine_id": schema.StringAttribute{
						Optional: true,
					},
					"filter_duplicates": schema.BoolAttribute{
						Optional: true,
					},
					"filter_interfaces": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"filter_internal_interfaces": schema.BoolAttribute{
						Optional: true,
					},
					"if_count_with_filter_interfaces": schema.BoolAttribute{
						Optional: true,
					},
					"interface": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"location": schema.StringAttribute{
						Optional: true,
					},
					"routing_instance_access": schema.BoolAttribute{
						Optional: true,
					},
					"routing_instance_access_list": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"health_monitor": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"falling_threshold": schema.Int64Attribute{
									Optional: true,
								},
								"idp": schema.BoolAttribute{
									Optional: true,
								},
								"idp_falling_threshold": schema.Int64Attribute{
									Optional: true,
								},
								"idp_interval": schema.Int64Attribute{
									Optional: true,
								},
								"idp_rising_threshold": schema.Int64Attribute{
									Optional: true,
								},
								"interval": schema.Int64Attribute{
									Optional: true,
								},
								"rising_threshold": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSnmpV0toV1,
		},
	}
}

func upgradeSnmpV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                          types.String   `tfsdk:"id"`
		CleanOnDestroy              types.Bool     `tfsdk:"clean_on_destroy"`
		ARP                         types.Bool     `tfsdk:"arp"`
		ARPHostNameResolution       types.Bool     `tfsdk:"arp_host_name_resolution"`
		Contact                     types.String   `tfsdk:"contact"`
		Description                 types.String   `tfsdk:"description"`
		EngineID                    types.String   `tfsdk:"engine_id"`
		FilterDuplicates            types.Bool     `tfsdk:"filter_duplicates"`
		FilterInterfaces            []types.String `tfsdk:"filter_interfaces"`
		FilterInternalInterfaces    types.Bool     `tfsdk:"filter_internal_interfaces"`
		IfCountWithFilterInterfaces types.Bool     `tfsdk:"if_count_with_filter_interfaces"`
		Interface                   []types.String `tfsdk:"interface"`
		Location                    types.String   `tfsdk:"location"`
		RoutingInstanceAccess       types.Bool     `tfsdk:"routing_instance_access"`
		RoutingInstanceAccessList   []types.String `tfsdk:"routing_instance_access_list"`
		HealthMonitor               []struct {
			FallingThreshold    types.Int64 `tfsdk:"falling_threshold"`
			Idp                 types.Bool  `tfsdk:"idp"`
			IdpFallingThreshold types.Int64 `tfsdk:"idp_falling_threshold"`
			IdpInterval         types.Int64 `tfsdk:"idp_interval"`
			IdpRisingThreshold  types.Int64 `tfsdk:"idp_rising_threshold"`
			Interval            types.Int64 `tfsdk:"interval"`
			RisingThreshold     types.Int64 `tfsdk:"rising_threshold"`
		} `tfsdk:"health_monitor"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 snmpData
	dataV1.ID = dataV0.ID
	dataV1.CleanOnDestroy = dataV0.CleanOnDestroy
	if !dataV1.CleanOnDestroy.IsNull() && !dataV1.CleanOnDestroy.ValueBool() {
		dataV1.CleanOnDestroy = types.BoolNull()
	}
	dataV1.ARP = dataV0.ARP
	dataV1.ARPHostNameResolution = dataV0.ARPHostNameResolution
	dataV1.FilterDuplicates = dataV0.FilterDuplicates
	dataV1.FilterInternalInterfaces = dataV0.FilterInternalInterfaces
	dataV1.IfCountWithFilterInterfaces = dataV0.IfCountWithFilterInterfaces
	dataV1.RoutingInstanceAccess = dataV0.RoutingInstanceAccess
	dataV1.Contact = dataV0.Contact
	dataV1.Description = dataV0.Description
	dataV1.EngineID = dataV0.EngineID
	dataV1.FilterInterfaces = dataV0.FilterInterfaces
	dataV1.Interface = dataV0.Interface
	dataV1.Location = dataV0.Location
	dataV1.RoutingInstanceAccessList = dataV0.RoutingInstanceAccessList
	if len(dataV0.HealthMonitor) > 0 {
		dataV1.HealthMonitor = &snmpBlockHealthMonitor{
			Idp:                 dataV0.HealthMonitor[0].Idp,
			FallingThreshold:    dataV0.HealthMonitor[0].FallingThreshold,
			IdpFallingThreshold: dataV0.HealthMonitor[0].IdpFallingThreshold,
			IdpInterval:         dataV0.HealthMonitor[0].IdpInterval,
			IdpRisingThreshold:  dataV0.HealthMonitor[0].IdpRisingThreshold,
			Interval:            dataV0.HealthMonitor[0].Interval,
			RisingThreshold:     dataV0.HealthMonitor[0].RisingThreshold,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
