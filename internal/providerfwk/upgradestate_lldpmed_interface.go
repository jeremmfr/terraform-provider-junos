package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *lldpMedInterface) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
				},
				Blocks: map[string]schema.Block{
					"location": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"civic_based_country_code": schema.StringAttribute{
									Optional: true,
								},
								"civic_based_what": schema.Int64Attribute{
									Optional: true,
								},
								"co_ordinate_latitude": schema.Int64Attribute{
									Optional: true,
								},
								"co_ordinate_longitude": schema.Int64Attribute{
									Optional: true,
								},
								"elin": schema.StringAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"civic_based_ca_type": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"ca_type": schema.Int64Attribute{
												Required: true,
											},
											"ca_value": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeLldpMedInterfaceStateV0toV1,
		},
	}
}

func upgradeLldpMedInterfaceStateV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID       types.String `tfsdk:"id"`
		Name     types.String `tfsdk:"name"`
		Disable  types.Bool   `tfsdk:"disable"`
		Enable   types.Bool   `tfsdk:"enable"`
		Location []struct {
			CivicBasedCountryCode types.String `tfsdk:"civic_based_country_code"`
			CivicBasedWhat        types.Int64  `tfsdk:"civic_based_what"`
			CoOrdinateLatitude    types.Int64  `tfsdk:"co_ordinate_latitude"`
			CoOrdinateLongitude   types.Int64  `tfsdk:"co_ordinate_longitude"`
			Elin                  types.String `tfsdk:"elin"`
			CivicBasedCaType      []struct {
				CaType  types.Int64  `tfsdk:"ca_type"`
				CaValue types.String `tfsdk:"ca_value"`
			} `tfsdk:"civic_based_ca_type"`
		} `tfsdk:"location"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 lldpMedInterfaceData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Disable = dataV0.Disable
	dataV1.Enable = dataV0.Enable
	if len(dataV0.Location) > 0 {
		dataV1.Location = &lldpMedInterfaceBlockLocation{
			CivicBasedCountryCode: dataV0.Location[0].CivicBasedCountryCode,
			CivicBasedWhat:        dataV0.Location[0].CivicBasedWhat,
			CoOrdinateLatitude:    dataV0.Location[0].CoOrdinateLatitude,
			CoOrdinateLongitude:   dataV0.Location[0].CoOrdinateLongitude,
			Elin:                  dataV0.Location[0].Elin,
		}

		for _, blockV0 := range dataV0.Location[0].CivicBasedCaType {
			dataV1.Location.CivicBasedCaType = append(dataV1.Location.CivicBasedCaType,
				lldpMedInterfaceBlockLocationBlockCivicBasedCaType{
					CaType:  blockV0.CaType,
					CaValue: blockV0.CaValue,
				})
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
