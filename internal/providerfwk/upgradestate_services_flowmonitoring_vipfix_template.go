package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *servicesFlowMonitoringVIPFixTemplate) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
					},
					"flow_active_timeout": schema.Int64Attribute{
						Optional: true,
					},
					"flow_inactive_timeout": schema.Int64Attribute{
						Optional: true,
					},
					"flow_key_flow_direction": schema.BoolAttribute{
						Optional: true,
					},
					"flow_key_vlan_id": schema.BoolAttribute{
						Optional: true,
					},
					"ip_template_export_extension": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"nexthop_learning_enable": schema.BoolAttribute{
						Optional: true,
					},
					"nexthop_learning_disable": schema.BoolAttribute{
						Optional: true,
					},
					"observation_domain_id": schema.Int64Attribute{
						Optional: true,
					},
					"option_template_id": schema.Int64Attribute{
						Optional: true,
					},
					"template_id": schema.Int64Attribute{
						Optional: true,
					},
					"tunnel_observation_ipv4": schema.BoolAttribute{
						Optional: true,
					},
					"tunnel_observation_ipv6": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"option_refresh_rate": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"packets": schema.Int64Attribute{
									Optional: true,
								},
								"seconds": schema.Int64Attribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeServicesFlowMonitoringVIPFixTemplateStateV0toV1,
		},
	}
}

func upgradeServicesFlowMonitoringVIPFixTemplateStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                        types.String   `tfsdk:"id"`
		Name                      types.String   `tfsdk:"name"`
		Type                      types.String   `tfsdk:"type"`
		FlowActiveTimeout         types.Int64    `tfsdk:"flow_active_timeout"`
		FlowInactiveTimeout       types.Int64    `tfsdk:"flow_inactive_timeout"`
		FlowKeyFlowDirection      types.Bool     `tfsdk:"flow_key_flow_direction"`
		FlowKeyVlanID             types.Bool     `tfsdk:"flow_key_vlan_id"`
		IPTemplateExportExtension []types.String `tfsdk:"ip_template_export_extension"`
		NexthopLearningEnable     types.Bool     `tfsdk:"nexthop_learning_enable"`
		NexthopLearningDisable    types.Bool     `tfsdk:"nexthop_learning_disable"`
		ObservationDomainID       types.Int64    `tfsdk:"observation_domain_id"`
		OptionTemplateID          types.Int64    `tfsdk:"option_template_id"`
		TemplateID                types.Int64    `tfsdk:"template_id"`
		TunnelObservationIPv4     types.Bool     `tfsdk:"tunnel_observation_ipv4"`
		TunnelObservationIPv6     types.Bool     `tfsdk:"tunnel_observation_ipv6"`
		OptionRefreshRate         []struct {
			Packets types.Int64 `tfsdk:"packets"`
			Seconds types.Int64 `tfsdk:"seconds"`
		} `tfsdk:"option_refresh_rate"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 servicesFlowMonitoringVIPFixTemplateData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Type = dataV0.Type
	dataV1.FlowActiveTimeout = dataV0.FlowActiveTimeout
	dataV1.FlowInactiveTimeout = dataV0.FlowInactiveTimeout
	dataV1.FlowKeyFlowDirection = dataV0.FlowKeyFlowDirection
	dataV1.FlowKeyVlanID = dataV0.FlowKeyVlanID
	dataV1.IPTemplateExportExtension = dataV0.IPTemplateExportExtension
	dataV1.NexthopLearningEnable = dataV0.NexthopLearningEnable
	dataV1.NexthopLearningDisable = dataV0.NexthopLearningDisable
	dataV1.ObservationDomainID = dataV0.ObservationDomainID
	dataV1.OptionTemplateID = dataV0.OptionTemplateID
	dataV1.TemplateID = dataV0.TemplateID
	dataV1.TunnelObservationIPv4 = dataV0.TunnelObservationIPv4
	dataV1.TunnelObservationIPv6 = dataV0.TunnelObservationIPv6
	if len(dataV0.OptionRefreshRate) > 0 {
		dataV1.OptionRefreshRate = &servicesFlowMonitoringVIPFixTemplateBlockRefreshRate{
			Packets: dataV0.OptionRefreshRate[0].Packets,
			Seconds: dataV0.OptionRefreshRate[0].Seconds,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
