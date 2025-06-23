package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *securityUtmPolicy) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
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
					"anti_spam_smtp_profile": schema.StringAttribute{
						Optional: true,
					},
					"web_filtering_profile": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"anti_virus": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ftp_download_profile": schema.StringAttribute{
									Optional: true,
								},
								"ftp_upload_profile": schema.StringAttribute{
									Optional: true,
								},
								"http_profile": schema.StringAttribute{
									Optional: true,
								},
								"imap_profile": schema.StringAttribute{
									Optional: true,
								},
								"pop3_profile": schema.StringAttribute{
									Optional: true,
								},
								"smtp_profile": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"content_filtering": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ftp_download_profile": schema.StringAttribute{
									Optional: true,
								},
								"ftp_upload_profile": schema.StringAttribute{
									Optional: true,
								},
								"http_profile": schema.StringAttribute{
									Optional: true,
								},
								"imap_profile": schema.StringAttribute{
									Optional: true,
								},
								"pop3_profile": schema.StringAttribute{
									Optional: true,
								},
								"smtp_profile": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
					"traffic_sessions_per_client": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"limit": schema.Int64Attribute{
									Optional: true,
								},
								"over_limit": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSecurityUtmPolicyV0toV1,
		},
	}
}

func upgradeSecurityUtmPolicyV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                  types.String `tfsdk:"id"`
		Name                types.String `tfsdk:"name"`
		AntiSpamSMTPProfile types.String `tfsdk:"anti_spam_smtp_profile"`
		WebFilteringProfile types.String `tfsdk:"web_filtering_profile"`
		AntiVirus           []struct {
			FTPDownloadProfile types.String `tfsdk:"ftp_download_profile"`
			FTPUploadProfile   types.String `tfsdk:"ftp_upload_profile"`
			HTTPProfile        types.String `tfsdk:"http_profile"`
			IMAPProfile        types.String `tfsdk:"imap_profile"`
			POP3Profile        types.String `tfsdk:"pop3_profile"`
			SMTPProfile        types.String `tfsdk:"smtp_profile"`
		} `tfsdk:"anti_virus"`
		ContentFiltering []struct {
			FTPDownloadProfile types.String `tfsdk:"ftp_download_profile"`
			FTPUploadProfile   types.String `tfsdk:"ftp_upload_profile"`
			HTTPProfile        types.String `tfsdk:"http_profile"`
			IMAPProfile        types.String `tfsdk:"imap_profile"`
			POP3Profile        types.String `tfsdk:"pop3_profile"`
			SMTPProfile        types.String `tfsdk:"smtp_profile"`
		} `tfsdk:"content_filtering"`
		TrafficSessionsPerClient []struct {
			Limit     types.Int64  `tfsdk:"limit"`
			OverLimit types.String `tfsdk:"over_limit"`
		} `tfsdk:"traffic_sessions_per_client"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 securityUtmPolicyData
	dataV1.ID = dataV0.ID
	dataV1.Name = dataV0.Name
	dataV1.AntiSpamSMTPProfile = dataV0.AntiSpamSMTPProfile
	dataV1.WebFilteringProfile = dataV0.WebFilteringProfile
	if len(dataV0.AntiVirus) > 0 {
		dataV1.AntiVirus = &securityUtmPolicyBlockProtocolProfile{
			FTPDownloadProfile: dataV0.AntiVirus[0].FTPDownloadProfile,
			FTPUploadProfile:   dataV0.AntiVirus[0].FTPUploadProfile,
			HTTPProfile:        dataV0.AntiVirus[0].HTTPProfile,
			IMAPProfile:        dataV0.AntiVirus[0].IMAPProfile,
			POP3Profile:        dataV0.AntiVirus[0].POP3Profile,
			SMTPProfile:        dataV0.AntiVirus[0].SMTPProfile,
		}
	}
	if len(dataV0.ContentFiltering) > 0 {
		dataV1.ContentFiltering = &securityUtmPolicyBlockProtocolProfile{
			FTPDownloadProfile: dataV0.ContentFiltering[0].FTPDownloadProfile,
			FTPUploadProfile:   dataV0.ContentFiltering[0].FTPUploadProfile,
			HTTPProfile:        dataV0.ContentFiltering[0].HTTPProfile,
			IMAPProfile:        dataV0.ContentFiltering[0].IMAPProfile,
			POP3Profile:        dataV0.ContentFiltering[0].POP3Profile,
			SMTPProfile:        dataV0.ContentFiltering[0].SMTPProfile,
		}
	}
	if len(dataV0.TrafficSessionsPerClient) > 0 {
		dataV1.TrafficSessionsPerClient = &securityUtmPolicyBlockTrafficSessionsPerClient{
			Limit:     dataV0.TrafficSessionsPerClient[0].Limit,
			OverLimit: dataV0.TrafficSessionsPerClient[0].OverLimit,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
