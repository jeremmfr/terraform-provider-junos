package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *lldpInterface) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"disable": schema.BoolAttribute{
						Optional: true,
					},
					"enable": schema.BoolAttribute{
						Optional: true,
					},
					"trap_notification_disable": schema.BoolAttribute{
						Optional: true,
					},
					"trap_notification_enable": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"power_negotiation": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"disable": schema.BoolAttribute{
									Optional: true,
								},
								"enable": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeLldpInterfaceStateV0toV1,
		},
	}
}

func upgradeLldpInterfaceStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                      types.String `tfsdk:"id"`
		Name                    types.String `tfsdk:"name"`
		Disable                 types.Bool   `tfsdk:"disable"`
		Enable                  types.Bool   `tfsdk:"enable"`
		TrapNotificationDisable types.Bool   `tfsdk:"trap_notification_disable"`
		TrapNotificationEnable  types.Bool   `tfsdk:"trap_notification_enable"`
		PowerNegotiation        []struct {
			Disable types.Bool `tfsdk:"disable"`
			Enable  types.Bool `tfsdk:"enable"`
		} `tfsdk:"power_negotiation"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 lldpInterfaceData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Disable = dataV0.Disable
	dataV1.Enable = dataV0.Enable
	dataV1.TrapNotificationDisable = dataV0.TrapNotificationDisable
	dataV1.TrapNotificationEnable = dataV0.TrapNotificationEnable
	if len(dataV0.PowerNegotiation) > 0 {
		dataV1.PowerNegotiation = &lldpInterfaceBlockPowerNegotiation{
			Enable:  dataV0.PowerNegotiation[0].Enable,
			Disable: dataV0.PowerNegotiation[0].Disable,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
