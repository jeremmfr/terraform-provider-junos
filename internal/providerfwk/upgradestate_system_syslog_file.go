package providerfwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (rsc *systemSyslogFile) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"filename": schema.StringAttribute{
						Required: true,
					},
					"allow_duplicates": schema.BoolAttribute{
						Optional: true,
					},
					"explicit_priority": schema.BoolAttribute{
						Optional: true,
					},
					"match": schema.StringAttribute{
						Optional: true,
					},
					"match_strings": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
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
					"archive": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"binary_data": schema.BoolAttribute{
									Optional: true,
								},
								"no_binary_data": schema.BoolAttribute{
									Optional: true,
								},
								"files": schema.Int64Attribute{
									Optional: true,
								},
								"size": schema.Int64Attribute{
									Optional: true,
								},
								"start_time": schema.StringAttribute{
									Optional: true,
								},
								"transfer_interval": schema.Int64Attribute{
									Optional: true,
								},
								"world_readable": schema.BoolAttribute{
									Optional: true,
								},
								"no_world_readable": schema.BoolAttribute{
									Optional: true,
								},
							},
							Blocks: map[string]schema.Block{
								"sites": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"url": schema.StringAttribute{
												Required: true,
											},
											"password": schema.StringAttribute{
												Optional:  true,
												Sensitive: true,
											},
											"routing_instance": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
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
			StateUpgrader: upgradeSystemSyslogFileV0toV1,
		},
	}
}

func upgradeSystemSyslogFileV0toV1(
	ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse,
) {
	type modelV0 struct {
		ID                          types.String   `tfsdk:"id"`
		Filename                    types.String   `tfsdk:"filename"`
		AllowDuplicates             types.Bool     `tfsdk:"allow_duplicates"`
		ExplicitPriority            types.Bool     `tfsdk:"explicit_priority"`
		Match                       types.String   `tfsdk:"match"`
		MatchStrings                []types.String `tfsdk:"match_strings"`
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
		Archive                     []struct {
			BinaryData       types.Bool   `tfsdk:"binary_data"`
			NoBinaryData     types.Bool   `tfsdk:"no_binary_data"`
			Files            types.Int64  `tfsdk:"files"`
			Size             types.Int64  `tfsdk:"size"`
			StartTime        types.String `tfsdk:"start_time"`
			TransferInterval types.Int64  `tfsdk:"transfer_interval"`
			WorldReadable    types.Bool   `tfsdk:"world_readable"`
			NoWorldReadable  types.Bool   `tfsdk:"no_world_readable"`
			Sites            []struct {
				URL             types.String `tfsdk:"url"`
				Password        types.String `tfsdk:"password"`
				RoutingInstance types.String `tfsdk:"routing_instance"`
			} `tfsdk:"sites"`
		} `tfsdk:"archive"`
		StructuredData []struct {
			Brief types.Bool `tfsdk:"brief"`
		} `tfsdk:"structured_data"`
	}

	var dataV0 modelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &dataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataV1 systemSyslogFileData
	dataV1.ID = dataV0.ID
	dataV1.Filename = dataV0.Filename
	dataV1.AllowDuplicates = dataV0.AllowDuplicates
	dataV1.ExplicitPriority = dataV0.ExplicitPriority
	dataV1.Match = dataV0.Match
	dataV1.MatchStrings = dataV0.MatchStrings
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
	if len(dataV0.Archive) > 0 {
		dataV1.Archive = &systemSyslogFileBlockArchive{
			BinaryData:       dataV0.Archive[0].BinaryData,
			NoBinaryData:     dataV0.Archive[0].NoBinaryData,
			WorldReadable:    dataV0.Archive[0].WorldReadable,
			NoWorldReadable:  dataV0.Archive[0].NoWorldReadable,
			Files:            dataV0.Archive[0].Files,
			Size:             dataV0.Archive[0].Size,
			StartTime:        dataV0.Archive[0].StartTime,
			TransferInterval: dataV0.Archive[0].TransferInterval,
		}
		for _, block := range dataV0.Archive[0].Sites {
			dataV1.Archive.Sites = append(dataV1.Archive.Sites,
				systemSyslogFileBlockArchiveBlockSites{
					URL:             block.URL,
					Password:        block.Password,
					RoutingInstance: block.RoutingInstance,
				},
			)
		}
	}
	if len(dataV0.StructuredData) > 0 {
		dataV1.StructuredData = &systemSyslogFileBlockStructuredData{
			Brief: dataV0.StructuredData[0].Brief,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, dataV1)...)
}
