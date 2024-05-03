package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *systemSyslogHost) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"host": schema.StringAttribute{
						Required: true,
					},
					"allow_duplicates": schema.BoolAttribute{
						Optional: true,
					},
					"exclude_hostname": schema.BoolAttribute{
						Optional: true,
					},
					"explicit_priority": schema.BoolAttribute{
						Optional: true,
					},
					"facility_override": schema.StringAttribute{
						Optional: true,
					},
					"log_prefix": schema.StringAttribute{
						Optional: true,
					},
					"match": schema.StringAttribute{
						Optional: true,
					},
					"match_strings": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"port": schema.Int64Attribute{
						Optional: true,
					},
					"source_address": schema.StringAttribute{
						Optional: true,
					},
					"any_severity": schema.StringAttribute{
						Optional: true,
					},
					"authorization_severity": schema.StringAttribute{
						Optional: true,
					},
					"changelog_severity": schema.StringAttribute{
						Optional: true,
					},
					"conflictlog_severity": schema.StringAttribute{
						Optional: true,
					},
					"daemon_severity": schema.StringAttribute{
						Optional: true,
					},
					"dfc_severity": schema.StringAttribute{
						Optional: true,
					},
					"external_severity": schema.StringAttribute{
						Optional: true,
					},
					"firewall_severity": schema.StringAttribute{
						Optional: true,
					},
					"ftp_severity": schema.StringAttribute{
						Optional: true,
					},
					"interactivecommands_severity": schema.StringAttribute{
						Optional: true,
					},
					"kernel_severity": schema.StringAttribute{
						Optional: true,
					},
					"ntp_severity": schema.StringAttribute{
						Optional: true,
					},
					"pfe_severity": schema.StringAttribute{
						Optional: true,
					},
					"security_severity": schema.StringAttribute{
						Optional: true,
					},
					"user_severity": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"structured_data": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"brief": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: upgradeSystemSyslogHostV0toV1,
		},
	}
}

func upgradeSystemSyslogHostV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                          types.String   `tfsdk:"id"`
		Host                        types.String   `tfsdk:"host"`
		AllowDuplicates             types.Bool     `tfsdk:"allow_duplicates"`
		ExcludeHostname             types.Bool     `tfsdk:"exclude_hostname"`
		ExplicitPriority            types.Bool     `tfsdk:"explicit_priority"`
		FacilityOverride            types.String   `tfsdk:"facility_override"`
		LogPrefix                   types.String   `tfsdk:"log_prefix"`
		Match                       types.String   `tfsdk:"match"`
		MatchStrings                []types.String `tfsdk:"match_strings"`
		Port                        types.Int64    `tfsdk:"port"`
		SourceAddress               types.String   `tfsdk:"source_address"`
		AnySeverity                 types.String   `tfsdk:"any_severity"`
		AuthorizationSeverity       types.String   `tfsdk:"authorization_severity"`
		ChangelogSeverity           types.String   `tfsdk:"changelog_severity"`
		ConflictlogSeverity         types.String   `tfsdk:"conflictlog_severity"`
		DaemonSeverity              types.String   `tfsdk:"daemon_severity"`
		DfcSeverity                 types.String   `tfsdk:"dfc_severity"`
		ExternalSeverity            types.String   `tfsdk:"external_severity"`
		FirewallSeverity            types.String   `tfsdk:"firewall_severity"`
		FtpSeverity                 types.String   `tfsdk:"ftp_severity"`
		InteractivecommandsSeverity types.String   `tfsdk:"interactivecommands_severity"`
		KernelSeverity              types.String   `tfsdk:"kernel_severity"`
		NtpSeverity                 types.String   `tfsdk:"ntp_severity"`
		PfeSeverity                 types.String   `tfsdk:"pfe_severity"`
		SecuritySeverity            types.String   `tfsdk:"security_severity"`
		UserSeverity                types.String   `tfsdk:"user_severity"`
		StructuredData              []struct {
			Brief types.Bool `tfsdk:"brief"`
		} `tfsdk:"structured_data"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 systemSyslogHostData
	dataV1.ID = dataV0.ID
	dataV1.Host = dataV0.Host
	dataV1.AllowDuplicates = dataV0.AllowDuplicates
	dataV1.ExcludeHostname = dataV0.ExcludeHostname
	dataV1.ExplicitPriority = dataV0.ExplicitPriority
	dataV1.FacilityOverride = dataV0.FacilityOverride
	dataV1.LogPrefix = dataV0.LogPrefix
	dataV1.Match = dataV0.Match
	dataV1.MatchStrings = dataV0.MatchStrings
	dataV1.Port = dataV0.Port
	dataV1.SourceAddress = dataV0.SourceAddress
	dataV1.AnySeverity = dataV0.AnySeverity
	dataV1.AuthorizationSeverity = dataV0.AuthorizationSeverity
	dataV1.ChangelogSeverity = dataV0.ChangelogSeverity
	dataV1.ConflictlogSeverity = dataV0.ConflictlogSeverity
	dataV1.DaemonSeverity = dataV0.DaemonSeverity
	dataV1.DfcSeverity = dataV0.DfcSeverity
	dataV1.ExternalSeverity = dataV0.ExternalSeverity
	dataV1.FirewallSeverity = dataV0.FirewallSeverity
	dataV1.FtpSeverity = dataV0.FtpSeverity
	dataV1.InteractivecommandsSeverity = dataV0.InteractivecommandsSeverity
	dataV1.KernelSeverity = dataV0.KernelSeverity
	dataV1.NtpSeverity = dataV0.NtpSeverity
	dataV1.PfeSeverity = dataV0.PfeSeverity
	dataV1.SecuritySeverity = dataV0.SecuritySeverity
	dataV1.UserSeverity = dataV0.UserSeverity
	if len(dataV0.StructuredData) > 0 {
		dataV1.StructuredData = &systemSyslogHostBlockStructuredData{
			Brief: dataV0.StructuredData[0].Brief,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
