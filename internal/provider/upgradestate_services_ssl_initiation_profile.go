package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *servicesSSLInitiationProfile) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"client_certificate": schema.StringAttribute{
						Optional: true,
					},
					"custom_ciphers": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"enable_flow_tracing": schema.BoolAttribute{
						Optional: true,
					},
					"enable_session_cache": schema.BoolAttribute{
						Optional: true,
					},
					"preferred_ciphers": schema.StringAttribute{
						Optional: true,
					},
					"protocol_version": schema.StringAttribute{
						Optional: true,
					},
					"trusted_ca": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"actions": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"crl_disable": schema.BoolAttribute{
									Optional: true,
								},
								"crl_if_not_present": schema.StringAttribute{
									Optional: true,
								},
								"crl_ignore_hold_instruction_code": schema.BoolAttribute{
									Optional: true,
								},
								"ignore_server_auth_failure": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeServicesSSLInitiationProfileV0toV1,
		},
	}
}

func upgradeServicesSSLInitiationProfileV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                 types.String   `tfsdk:"id"`
		Name               types.String   `tfsdk:"name"`
		ClientCertificate  types.String   `tfsdk:"client_certificate"`
		CustomCiphers      []types.String `tfsdk:"custom_ciphers"`
		EnableFlowTracing  types.Bool     `tfsdk:"enable_flow_tracing"`
		EnableSessionCache types.Bool     `tfsdk:"enable_session_cache"`
		PreferredCiphers   types.String   `tfsdk:"preferred_ciphers"`
		ProtocolVersion    types.String   `tfsdk:"protocol_version"`
		TrustedCA          []types.String `tfsdk:"trusted_ca"`
		Actions            []struct {
			CrlDisable                   types.Bool   `tfsdk:"crl_disable"`
			CrlIfNotPresent              types.String `tfsdk:"crl_if_not_present"`
			CrlIgnoreHoldInstructionCode types.Bool   `tfsdk:"crl_ignore_hold_instruction_code"`
			IgnoreServerAuthFailure      types.Bool   `tfsdk:"ignore_server_auth_failure"`
		} `tfsdk:"actions"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 servicesSSLInitiationProfileData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.ClientCertificate = dataV0.ClientCertificate
	dataV1.CustomCiphers = dataV0.CustomCiphers
	dataV1.EnableFlowTracing = dataV0.EnableFlowTracing
	dataV1.EnableSessionCache = dataV0.EnableSessionCache
	dataV1.PreferredCiphers = dataV0.PreferredCiphers
	dataV1.ProtocolVersion = dataV0.ProtocolVersion
	dataV1.TrustedCA = dataV0.TrustedCA
	for _, blockV0 := range dataV0.Actions {
		dataV1.Actions = &servicesSSLInitiationProfileBlockActions{
			CrlDisable:                   blockV0.CrlDisable,
			CrlIfNotPresent:              blockV0.CrlIfNotPresent,
			CrlIgnoreHoldInstructionCode: blockV0.CrlIgnoreHoldInstructionCode,
			IgnoreServerAuthFailure:      blockV0.IgnoreServerAuthFailure,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
