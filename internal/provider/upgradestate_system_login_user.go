package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *systemLoginUser) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"class": schema.StringAttribute{
						Required: true,
					},
					"uid": schema.Int64Attribute{
						Optional: true,
						Computed: true,
					},
					"cli_prompt": schema.StringAttribute{
						Optional: true,
					},
					"full_name": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"authentication": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"encrypted_password": schema.StringAttribute{
									Optional: true,
								},
								"no_public_keys": schema.BoolAttribute{
									Optional: true,
								},
								"plain_text_password": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
								"ssh_public_keys": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSystemLoginUserV0toV1,
		},
	}
}

func upgradeSystemLoginUserV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID             types.String `tfsdk:"id"`
		Name           types.String `tfsdk:"name"`
		Class          types.String `tfsdk:"class"`
		UID            types.Int64  `tfsdk:"uid"`
		CliPrompt      types.String `tfsdk:"cli_prompt"`
		FullName       types.String `tfsdk:"full_name"`
		Authentication []struct {
			EncryptedPassword types.String   `tfsdk:"encrypted_password"`
			NoPublicKeys      types.Bool     `tfsdk:"no_public_keys"`
			PlainTextPassword types.String   `tfsdk:"plain_text_password"`
			SSHPublicKeys     []types.String `tfsdk:"ssh_public_keys"`
		} `tfsdk:"authentication"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 systemLoginUserData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.Class = dataV0.Class
	if dataV0.UID.ValueInt64() != 0 {
		dataV1.UID = dataV0.UID
	}
	dataV1.CliPrompt = dataV0.CliPrompt
	dataV1.FullName = dataV0.FullName
	if len(dataV0.Authentication) > 0 {
		dataV1.Authentication = &systemLoginUserBlockAuthentication{
			EncryptedPassword: dataV0.Authentication[0].EncryptedPassword,
			NoPublicKeys:      dataV0.Authentication[0].NoPublicKeys,
			PlainTextPassword: dataV0.Authentication[0].PlainTextPassword,
			SSHPublicKeys:     dataV0.Authentication[0].SSHPublicKeys,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
