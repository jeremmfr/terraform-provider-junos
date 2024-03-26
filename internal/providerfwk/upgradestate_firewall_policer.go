package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *firewallPolicer) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"filter_specific": schema.BoolAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"if_exceeding": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"burst_size_limit": schema.StringAttribute{
									Required: true,
								},
								"bandwidth_percent": schema.Int64Attribute{
									Optional: true,
								},
								"bandwidth_limit": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"then": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"discard": schema.BoolAttribute{
									Optional: true,
								},
								"forwarding_class": schema.StringAttribute{
									Optional: true,
								},
								"loss_priority": schema.StringAttribute{
									Optional: true,
								},
								"out_of_profile": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeFirewallPolicerStateV0toV1,
		},
	}
}

func upgradeFirewallPolicerStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID             types.String `tfsdk:"id"`
		Name           types.String `tfsdk:"name"`
		FilterSpecific types.Bool   `tfsdk:"filter_specific"`
		IfExceeding    []struct {
			BurstSizeLimit   types.String `tfsdk:"burst_size_limit"`
			BandwidthPercent types.Int64  `tfsdk:"bandwidth_percent"`
			BandwidthLimit   types.String `tfsdk:"bandwidth_limit"`
		} `tfsdk:"if_exceeding"`
		Then []struct {
			Discard         types.Bool   `tfsdk:"discard"`
			ForwardingClass types.String `tfsdk:"forwarding_class"`
			LossPriority    types.String `tfsdk:"loss_priority"`
			OutOfProfile    types.Bool   `tfsdk:"out_of_profile"`
		} `tfsdk:"then"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 firewallPolicerData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.FilterSpecific = dataV0.FilterSpecific
	if len(dataV0.IfExceeding) > 0 {
		dataV1.IfExceeding = &firewallPolicerBlockIfExceeding{
			BurstSizeLimit:   dataV0.IfExceeding[0].BurstSizeLimit,
			BandwidthPercent: dataV0.IfExceeding[0].BandwidthPercent,
			BandwidthLimit:   dataV0.IfExceeding[0].BandwidthLimit,
		}
	}
	if len(dataV0.Then) > 0 {
		dataV1.Then = &firewallPolicerBlockThen{
			Discard:         dataV0.Then[0].Discard,
			OutOfProfile:    dataV0.Then[0].OutOfProfile,
			ForwardingClass: dataV0.Then[0].ForwardingClass,
			LossPriority:    dataV0.Then[0].LossPriority,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
