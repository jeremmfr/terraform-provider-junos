package providerfwk

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"
	"github.com/jeremmfr/terraform-provider-junos/internal/version"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func New() provider.Provider {
	return &junosProvider{}
}

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &junosProvider{}

type junosProvider struct{}

type junosProviderModel struct {
	IP                         types.String `tfsdk:"ip"`
	Port                       types.Int64  `tfsdk:"port"`
	Username                   types.String `tfsdk:"username"`
	Password                   types.String `tfsdk:"password"`
	SSHKeyPem                  types.String `tfsdk:"sshkey_pem"`
	SSHKeyFile                 types.String `tfsdk:"sshkeyfile"`
	SSHKeyPass                 types.String `tfsdk:"keypass"`
	GroupIntDel                types.String `tfsdk:"group_interface_delete"`
	CmdSleepShort              types.Int64  `tfsdk:"cmd_sleep_short"`
	CmdSleepLock               types.Int64  `tfsdk:"cmd_sleep_lock"`
	CommitConfirmed            types.Int64  `tfsdk:"commit_confirmed"`
	CommitConfirmedWaitPercent types.Int64  `tfsdk:"commit_confirmed_wait_percent"`
	SleepSSHClosed             types.Int64  `tfsdk:"ssh_sleep_closed"`
	SSHCiphers                 types.List   `tfsdk:"ssh_ciphers"`
	SSHTimeoutToEstab          types.Int64  `tfsdk:"ssh_timeout_to_establish"`
	SSHRetryToEstab            types.Int64  `tfsdk:"ssh_retry_to_establish"`
	FilePermission             types.String `tfsdk:"file_permission"`
	DebugNetconfLogPath        types.String `tfsdk:"debug_netconf_log_path"`
	FakeCreateSetFile          types.String `tfsdk:"fake_create_with_setfile"`
	FakeUpdateAlso             types.Bool   `tfsdk:"fake_update_also"`
	FakeDeleteAlso             types.Bool   `tfsdk:"fake_delete_also"`
}

const (
	providerName = "junos"
)

// Metadata returns the provider type name.
func (p *junosProvider) Metadata(
	_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse,
) {
	resp.TypeName = providerName
	resp.Version = version.Get()
}

func (p *junosProvider) Schema(
	_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Optional: true,
				Description: "This is the target for Netconf session (ip or dns name)." +
					" May also be provided via " + junos.EnvHost + " environment variable.",
			},
			"port": schema.Int64Attribute{
				Optional: true,
				Description: "This is the tcp port for ssh connection." +
					" May also be provided via " + junos.EnvPort + " environment variable.",
			},
			"username": schema.StringAttribute{
				Optional: true,
				Description: "This is the username for ssh connection." +
					" May also be provided via " + junos.EnvUsername + " environment variable.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password": schema.StringAttribute{
				Optional: true,
				Description: "This is a password for ssh connection." +
					" May also be provided via " + junos.EnvPassword + " environment variable.",
			},
			"sshkey_pem": schema.StringAttribute{
				Optional: true,
				Description: "This is the ssh key in PEM format for establish ssh connection." +
					" May also be provided via " + junos.EnvKeyPem + " environment variable.",
			},
			"sshkeyfile": schema.StringAttribute{
				Optional: true,
				Description: "This is the path to ssh key for establish ssh connection." +
					" May also be provided via " + junos.EnvKeyFile + " environment variable.",
			},
			"keypass": schema.StringAttribute{
				Optional: true,
				Description: "This is the passphrase for open `sshkeyfile` or `sshkey_pem`." +
					" May also be provided via " + junos.EnvKeyPass + " environment variable.",
			},
			"group_interface_delete": schema.StringAttribute{
				Optional: true,
				Description: "This is the Junos group used to remove configuration on a physical interface." +
					" May also be provided via " + junos.EnvGroupInterfaceDelete + " environment variable.",
			},
			"cmd_sleep_short": schema.Int64Attribute{
				Optional: true,
				Description: "Milliseconds to wait after Terraform  provider executes an action on the Junos device." +
					" May also be provided via " + junos.EnvSleepShort + " environment variable.",
			},
			"cmd_sleep_lock": schema.Int64Attribute{
				Optional: true,
				Description: "Seconds of standby while waiting for Terraform provider " +
					"to lock candidate configuration on a Junos device." +
					" May also be provided via " + junos.EnvSleepLock + " environment variable.",
			},
			"commit_confirmed": schema.Int64Attribute{
				Optional: true,
				Description: "Number of minutes until automatic rollback." +
					" May also be provided via " + junos.EnvCommitConfirmed + " environment variable." +
					" For each resource action with commit, commit with `confirmed` option and" +
					" with the value ot this argument as `confirm-timeout`, " +
					" wait for `<commit_confirmed_wait_percent>`% of the minutes defined in the value of this argument," +
					" and confirm commit to avoid rollback with the `commit check` command.",
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"commit_confirmed_wait_percent": schema.Int64Attribute{
				Optional: true,
				Description: "Percentage of `<commit_confirmed>` minute(s) to wait between" +
					" `commit confirmed` (commit with automatic rollback) and" +
					" `commit check` (confirmation) commands." +
					" No effect if `<commit_confirmed>` is not used." +
					" May also be provided via " + junos.EnvCommitConfirmedWaitPercent + " environment variable." +
					" Defaults to 90.",
				Validators: []validator.Int64{
					int64validator.Between(0, 99),
				},
			},
			"ssh_sleep_closed": schema.Int64Attribute{
				Optional: true,
				Description: "Seconds to wait after Terraform provider closed a ssh connection." +
					" May also be provided via " + junos.EnvSleepSSHClosed + " environment variable.",
			},
			"ssh_ciphers": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Ciphers used in SSH connection.",
			},
			"ssh_timeout_to_establish": schema.Int64Attribute{
				Optional: true,
				Description: "Seconds to wait for establishing TCP connections when initiating SSH connections." +
					" May also be provided via " + junos.EnvSSHTimeoutToEstablish + " environment variable.",
			},
			"ssh_retry_to_establish": schema.Int64Attribute{
				Optional: true,
				Description: "Number of retries to establish SSH connections." +
					"The provider waits after each try, with the sleep time increasing by 1 second each time." +
					" May also be provided via " + junos.EnvSSHRetryToEstablish + " environment variable.",
				Validators: []validator.Int64{
					int64validator.Between(1, 10),
				},
			},
			"file_permission": schema.StringAttribute{
				Optional: true,
				Description: "The permission to set for the created file (debug, setfile)." +
					" May also be provided via " + junos.EnvFilePermission + " environment variable.",
			},
			"debug_netconf_log_path": schema.StringAttribute{
				Optional: true,
				Description: "More detailed log (netconf) in the specified file." +
					" May also be provided via " + junos.EnvLogPath + " environment variable.",
			},
			"fake_create_with_setfile": schema.StringAttribute{
				Optional: true,
				Description: "The normal process to create resources skipped to generate set lines, " +
					"append them to the specified file, " +
					"and respond with a `fake` successful creation of resources to Terraform." +
					" May also be provided via " + junos.EnvFakecreateSetfile + " environment variable.",
			},
			"fake_update_also": schema.BoolAttribute{
				Optional: true,
				Description: "The normal process to update resources skipped to generate set/delete lines, " +
					"append them to the same file as `fake_create_with_setfile`, " +
					"and respond with a `fake` successful update of resources to Terraform." +
					" May also be provided via " + junos.EnvFakeupdateAlso + " environment variable.",
			},
			"fake_delete_also": schema.BoolAttribute{
				Optional: true,
				Description: "The normal process to delete resources skipped to generate delete lines, " +
					"append them to the same file as `fake_create_with_setfile`, " +
					"and respond with a `fake` successful delete of resources to Terraform." +
					" May also be provided via " + junos.EnvFakedeleteAlso + " environment variable.",
			},
		},
	}
}

func (p *junosProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newApplicationSetsDataSource,
		newApplicationsDataSource,
		newChassisInventoryDataSource,
		newInterfaceLogicalDataSource,
		newInterfaceLogicalInfoDataSource,
		newInterfacePhysicalDataSource,
		newInterfacesPhysicalPresentDataSource,
		newRoutesDataSource,
		newRoutingInstanceDataSource,
		newSecurityZoneDataSource,
	}
}

func (p *junosProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newAggregateRouteResource,
		newApplicationSetResource,
		newApplicationResource,
		newBgpGroupResource,
		newBgpNeighborResource,
		newBridgeDomainResource,
		newEventoptionsDestinationResource,
		newEventoptionsGenerateEventResource,
		newEventoptionsPolicyResource,
		newEvpnResource,
		newFirewallFilterResource,
		newFirewallPolicerResource,
		newForwardingoptionsSamplingResource,
		newForwardingoptionsSamplingInstanceResource,
		newForwardingoptionsStormControlProfileResource,
		newGenerateRouteResource,
		newIccpResource,
		newIccpPeerResource,
		newInterfaceLogicalResource,
		newInterfacePhysicalDisableResource,
		newInterfacePhysicalResource,
		newInterfaceSt0UnitResource,
		newMultichassisResource,
		newMultichassisProtectionPeerResource,
		newOamGretunnelInterfaceResource,
		newPolicyoptionsASPathResource,
		newPolicyoptionsASPathGroupResource,
		newPolicyoptionsCommunityResource,
		newPolicyoptionsPolicyStatementResource,
		newPolicyoptionsPrefixListResource,
		newRoutingInstanceResource,
		newSecurityResource,
		newSecurityAddressBookResource,
		newSecurityGlobalPolicyResource,
		newSecurityIkeGatewayResource,
		newSecurityIkePolicyResource,
		newSecurityIkeProposalResource,
		newSecurityIpsecPolicyResource,
		newSecurityIpsecProposalResource,
		newSecurityIpsecVpnResource,
		newSecurityNatDestinationResource,
		newSecurityNatDestinationPoolResource,
		newSecurityNatSourceResource,
		newSecurityNatSourcePoolResource,
		newSecurityNatStaticResource,
		newSecurityNatStaticRuleResource,
		newSecurityPolicyResource,
		newSecurityPolicyTunnelPairPolicyResource,
		newSecurityZoneResource,
		newSecurityZoneBookAddressResource,
		newSecurityZoneBookAddressSetResource,
		newServicesFlowMonitoringV9TemplateResource,
		newServicesFlowMonitoringVIPFixTemplateResource,
		newSnmpResource,
		newSnmpClientlistResource,
		newSnmpCommunityResource,
		newSnmpV3CommunityResource,
		newSnmpV3UsmUserResource,
		newSnmpV3VacmAccessgroupResource,
		newSnmpV3VacmSecuritytogroupResource,
		newSnmpViewResource,
		newStaticRouteResource,
		newSwitchOptionsResource,
		newSystemResource,
		newSystemRadiusServerResource,
		newSystemSyslogFileResource,
		newSystemSyslogHostResource,
		newSystemSyslogUserResource,
		newSystemTacplusServerResource,
		newVirtualChassisResource,
	}
}

func (p *junosProvider) Configure(
	ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse,
) {
	var config junosProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	unknownValueErrorMessage := "The provider cannot create the Junos client as there is an unknown configuration value "
	instructionUnknownMessage := " Either target apply the source of the value first, " +
		"set the value statically in the configuration, or use the %s environment variable."
	if config.IP.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ip"),
			"Unknown Junos IP target",
			unknownValueErrorMessage+"for the Junos IP."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvHost),
		)
	}
	if config.Port.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("port"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'port' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvPort),
		)
	}
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'username' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvUsername),
		)
	}
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'password' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvPassword),
		)
	}
	if config.SSHKeyPem.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("sshkey_pem"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'sshkey_pem' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvKeyPem),
		)
	}
	if config.SSHKeyFile.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("sshkeyfile"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'sshkeyfile' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvKeyFile),
		)
	}
	if config.SSHKeyPass.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("keypass"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'keypass' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvKeyPass),
		)
	}
	if config.GroupIntDel.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("group_interface_delete"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'group_interface_delete' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvGroupInterfaceDelete),
		)
	}
	if config.CmdSleepShort.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("cmd_sleep_short"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'cmd_sleep_short' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvSleepShort),
		)
	}
	if config.CmdSleepLock.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("cmd_sleep_lock"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'cmd_sleep_lock' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvSleepLock),
		)
	}
	if config.CommitConfirmed.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("commit_confirmed"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'commit_confirmed' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvCommitConfirmed),
		)
	}
	if config.CommitConfirmedWaitPercent.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("commit_confirmed_wait_percent"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'commit_confirmed_wait_percent' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvCommitConfirmedWaitPercent),
		)
	}
	if config.SleepSSHClosed.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ssh_sleep_closed"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'ssh_sleep_closed' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvSleepSSHClosed),
		)
	}
	if config.SSHCiphers.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ssh_ciphers"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'ssh_ciphers' attribute."+
				" Either target apply the source of the value first or set the value statically in the configuration.",
		)
	}
	for _, v := range config.SSHCiphers.Elements() {
		if v.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ssh_ciphers"),
				tfdiag.UnknownJunosAttrErrSummary,
				unknownValueErrorMessage+"for 'ssh_ciphers' attribute."+
					" Either target apply the source of the value first or set the value statically in the configuration.",
			)
		}
	}
	if config.SSHTimeoutToEstab.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ssh_timeout_to_establish"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'ssh_timeout_to_establish' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvSSHTimeoutToEstablish),
		)
	}
	if config.SSHRetryToEstab.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ssh_retry_to_establish"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'ssh_retry_to_establish' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvSSHRetryToEstablish),
		)
	}
	if config.FilePermission.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("file_permission"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'file_permission' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvFilePermission),
		)
	}
	if config.DebugNetconfLogPath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("debug_netconf_log_path"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'debug_netconf_log_path' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvLogPath),
		)
	}
	if config.FakeCreateSetFile.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("fake_create_with_setfile"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'fake_create_with_setfile' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvFakecreateSetfile),
		)
	}
	if config.FakeUpdateAlso.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("fake_update_also"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'fake_update_also' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvFakeupdateAlso),
		)
	}
	if config.FakeDeleteAlso.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("fake_delete_also"),
			tfdiag.UnknownJunosAttrErrSummary,
			unknownValueErrorMessage+"for 'fake_delete_also' attribute."+
				fmt.Sprintf(instructionUnknownMessage, junos.EnvFakedeleteAlso),
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	hostIP := os.Getenv(junos.EnvHost)
	if !config.IP.IsNull() {
		hostIP = config.IP.ValueString()
	}
	if hostIP == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("ip"),
			"Missing Junos IP target",
			"The provider cannot create the Junos client as there is a missing or empty value for the Junos IP."+
				" Set the ip value in the configuration or use the "+junos.EnvHost+" environment variable."+
				" If either is already set, ensure the value is not empty.",
		)

		return
	}

	client := junos.NewClient(hostIP)

	client.WithPort(830) // default value for port
	if !config.Port.IsNull() {
		client.WithPort(int(config.Port.ValueInt64()))
	} else if v := os.Getenv(junos.EnvPort); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("port"),
				"Error to parse "+junos.EnvPort,
				fmt.Sprintf("Error to parse value in "+junos.EnvPort+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			client.WithPort(d)
		}
	}

	client.WithUserName("netconf") // default value for username
	if !config.Username.IsNull() {
		client.WithUserName(config.Username.ValueString())
	} else if v := os.Getenv(junos.EnvUsername); v != "" {
		client.WithUserName(v)
	}

	if !config.Password.IsNull() {
		client.WithPassword(config.Password.ValueString())
	} else if v := os.Getenv(junos.EnvPassword); v != "" {
		client.WithPassword(v)
	}

	if !config.SSHKeyPem.IsNull() {
		client.WithSSHKeyPEM(config.SSHKeyPem.ValueString())
	} else if v := os.Getenv(junos.EnvKeyPem); v != "" {
		client.WithSSHKeyPEM(v)
	}

	if !config.SSHKeyFile.IsNull() {
		keyFile := config.SSHKeyFile.ValueString()
		if err := utils.ReplaceTildeToHomeDir(&keyFile); err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("sshkeyfile"),
				"Bad value in sshkeyfile",
				fmt.Sprintf("Error to use value in sshkeyfile attribute: %s\n"+
					"So the attribute is not used", err),
			)
		} else {
			client.WithSSHKeyFile(keyFile)
		}
	} else if v := os.Getenv(junos.EnvKeyFile); v != "" {
		if err := utils.ReplaceTildeToHomeDir(&v); err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("sshkeyfile"),
				"Bad value in "+junos.EnvKeyFile,
				fmt.Sprintf("Error to use value in "+junos.EnvKeyFile+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			client.WithSSHKeyFile(v)
		}
	}

	if !config.SSHKeyPass.IsNull() {
		client.WithSSHKeyPassphrase(config.SSHKeyPass.ValueString())
	} else if v := os.Getenv(junos.EnvKeyPass); v != "" {
		client.WithSSHKeyPassphrase(v)
	}

	if !config.GroupIntDel.IsNull() {
		client.WithGroupInterfaceDelete(config.GroupIntDel.ValueString())
	} else if v := os.Getenv(junos.EnvGroupInterfaceDelete); v != "" {
		client.WithGroupInterfaceDelete(v)
	}

	client.WithSleepShort(100) // default value for cmd_sleep_short
	if !config.CmdSleepShort.IsNull() {
		client.WithSleepShort(int(config.CmdSleepShort.ValueInt64()))
	} else if v := os.Getenv(junos.EnvSleepShort); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("cmd_sleep_short"),
				"Error to parse "+junos.EnvSleepShort,
				fmt.Sprintf("Error to parse value in "+junos.EnvSleepShort+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			client.WithSleepShort(d)
		}
	}

	client.WithSleepLock(10) // default value for cmd_sleep_lock
	if !config.CmdSleepLock.IsNull() {
		client.WithSleepLock(int(config.CmdSleepLock.ValueInt64()))
	} else if v := os.Getenv(junos.EnvSleepLock); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("cmd_sleep_lock"),
				"Error to parse "+junos.EnvSleepLock,
				fmt.Sprintf("Error to parse value in "+junos.EnvSleepLock+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			client.WithSleepLock(d)
		}
	}

	if !config.CommitConfirmed.IsNull() {
		if _, err := client.WithCommitConfirmed(int(config.CommitConfirmed.ValueInt64())); err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("commit_confirmed"),
				"Bad value in commit_confirmed",
				fmt.Sprintf("Error to use value in commit_confirmed attribute: %s\n"+
					"So the attribute is not used", err),
			)
		}
	} else if v := os.Getenv(junos.EnvCommitConfirmed); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("commit_confirmed"),
				"Error to parse "+junos.EnvCommitConfirmed,
				fmt.Sprintf("Error to parse value in "+junos.EnvCommitConfirmed+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			if _, err := client.WithCommitConfirmed(d); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("commit_confirmed"),
					"Bad value in "+junos.EnvCommitConfirmed,
					fmt.Sprintf("Error to use value in "+junos.EnvCommitConfirmed+" environment variable: %s\n"+
						"So the variable is not used", err),
				)
			}
		}
	}

	_, _ = client.WithCommitConfirmedWaitPercent(90) // default value for commit_confirmed_wait_percent
	if !config.CommitConfirmedWaitPercent.IsNull() {
		if _, err := client.WithCommitConfirmedWaitPercent(int(config.CommitConfirmedWaitPercent.ValueInt64())); err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("commit_confirmed_wait_percent"),
				"Bad value in commit_confirmed_wait_percent",
				fmt.Sprintf("Error to use value in commit_confirmed_wait_percent attribute: %s\n"+
					"So the attribute is not used", err),
			)
		}
	} else if v := os.Getenv(junos.EnvCommitConfirmedWaitPercent); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("commit_confirmed_wait_percent"),
				"Error to parse "+junos.EnvCommitConfirmedWaitPercent,
				fmt.Sprintf("Error to parse value in "+junos.EnvCommitConfirmedWaitPercent+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			if _, err := client.WithCommitConfirmedWaitPercent(d); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("commit_confirmed_wait_percent"),
					"Bad value in "+junos.EnvCommitConfirmedWaitPercent,
					fmt.Sprintf("Error to use value in "+junos.EnvCommitConfirmedWaitPercent+" environment variable: %s\n"+
						"So the variable is not used", err),
				)
			}
		}
	}

	if !config.SleepSSHClosed.IsNull() {
		client.WithSleepSSHClosed(int(config.SleepSSHClosed.ValueInt64()))
	} else if v := os.Getenv(junos.EnvSleepSSHClosed); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("ssh_sleep_closed"),
				"Error to parse "+junos.EnvSleepSSHClosed,
				fmt.Sprintf("Error to parse value in "+junos.EnvSleepSSHClosed+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			client.WithSleepSSHClosed(d)
		}
	}

	client.WithSSHCiphers(junos.DefaultSSHCiphers())
	if !config.SSHCiphers.IsNull() && len(config.SSHCiphers.Elements()) > 0 {
		sshCiphers := make([]string, len(config.SSHCiphers.Elements()))
		for i, v := range config.SSHCiphers.Elements() {
			sshCiphers[i] = v.(types.String).ValueString()
		}
		client.WithSSHCiphers(sshCiphers)
	}

	if !config.SSHTimeoutToEstab.IsNull() {
		client.WithSSHTimeoutToEstablish(int(config.SSHTimeoutToEstab.ValueInt64()))
	} else if v := os.Getenv(junos.EnvSSHTimeoutToEstablish); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("ssh_timeout_to_establish"),
				"Error to parse "+junos.EnvSSHTimeoutToEstablish,
				fmt.Sprintf("Error to parse value in "+junos.EnvSSHTimeoutToEstablish+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			client.WithSSHTimeoutToEstablish(d)
		}
	}

	_, _ = client.WithSSHRetryToEstablish(1) // default value for ssh_retry_to_establish
	if !config.SSHRetryToEstab.IsNull() {
		if _, err := client.WithSSHRetryToEstablish(int(config.SSHRetryToEstab.ValueInt64())); err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("ssh_retry_to_establish"),
				"Bad value in ssh_retry_to_establish",
				fmt.Sprintf("Error to use value in 'ssh_retry_to_establish' attribute: %s\n"+
					"So the attribute has the default value", err),
			)
		}
	} else if v := os.Getenv(junos.EnvSSHRetryToEstablish); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("ssh_retry_to_establish"),
				"Error to parse "+junos.EnvSSHRetryToEstablish,
				fmt.Sprintf("Error to parse value in "+junos.EnvSSHRetryToEstablish+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			if _, err := client.WithSSHRetryToEstablish(d); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("ssh_retry_to_establish"),
					"Bad value in "+junos.EnvSSHRetryToEstablish,
					fmt.Sprintf("Error to use value in "+junos.EnvSSHRetryToEstablish+" environment variable: %s\n"+
						"So the variable is not used", err),
				)
			}
		}
	}

	_, _ = client.WithFilePermission(0o644) // default value for file_permission
	if !config.FilePermission.IsNull() {
		filePerm, err := strconv.ParseInt(config.FilePermission.ValueString(), 8, 64)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("file_permission"),
				"Error to parse file_permission",
				fmt.Sprintf("Error to parse value in file_permission attribute: %s\n"+
					"So the attribute has the default value", err),
			)
		} else {
			if _, err := client.WithFilePermission(filePerm); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("file_permission"),
					"Bad value in file_permission",
					fmt.Sprintf("Error to use value in file_permission attribute: %s\n"+
						"So the attribute has the default value", err),
				)
			}
		}
	} else if v := os.Getenv(junos.EnvFilePermission); v != "" {
		filePerm, err := strconv.ParseInt(v, 8, 64)
		if err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("file_permission"),
				"Error to parse "+junos.EnvFilePermission,
				fmt.Sprintf("Error to parse value in "+junos.EnvFilePermission+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			if _, err := client.WithFilePermission(filePerm); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("file_permission"),
					"Bad value in "+junos.EnvFilePermission,
					fmt.Sprintf("Error to use value in "+junos.EnvFilePermission+" environment variable: %s\n"+
						"So the variable is not used", err),
				)
			}
		}
	}

	if !config.DebugNetconfLogPath.IsNull() {
		logPath := config.DebugNetconfLogPath.ValueString()
		if err := utils.ReplaceTildeToHomeDir(&logPath); err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("debug_netconf_log_path"),
				"Bad value in debug_netconf_log_path",
				fmt.Sprintf("Error to use value in debug_netconf_log_path attribute: %s\n"+
					"So the attribute is not used", err),
			)
		} else {
			client.WithDebugLogFile(logPath)
		}
	} else if v := os.Getenv(junos.EnvLogPath); v != "" {
		if err := utils.ReplaceTildeToHomeDir(&v); err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("debug_netconf_log_path"),
				"Bad value in "+junos.EnvLogPath,
				fmt.Sprintf("Error to use value in "+junos.EnvLogPath+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			client.WithDebugLogFile(v)
		}
	}

	if !config.FakeCreateSetFile.IsNull() {
		setFile := config.FakeCreateSetFile.ValueString()
		if err := utils.ReplaceTildeToHomeDir(&setFile); err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("fake_create_with_setfile"),
				"Bad value in fake_create_with_setfile",
				fmt.Sprintf("Error to use value in fake_create_with_setfile attribute: %s\n"+
					"So the attribute is not used", err),
			)
		} else {
			client.WithFakeCreateSetFile(setFile)
		}
	} else if v := os.Getenv(junos.EnvFakecreateSetfile); v != "" {
		if err := utils.ReplaceTildeToHomeDir(&v); err != nil {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("fake_create_with_setfile"),
				"Bad value in "+junos.EnvFakecreateSetfile,
				fmt.Sprintf("Error to use value in "+junos.EnvFakecreateSetfile+" environment variable: %s\n"+
					"So the variable is not used", err),
			)
		} else {
			client.WithFakeCreateSetFile(v)
		}
	}

	if !config.FakeUpdateAlso.IsNull() {
		if config.FakeUpdateAlso.ValueBool() {
			client.WithFakeUpdateAlso()
		}
	} else if v := os.Getenv(junos.EnvFakeupdateAlso); strings.EqualFold(v, "true") || v == "1" {
		client.WithFakeUpdateAlso()
	}

	if !config.FakeDeleteAlso.IsNull() {
		if config.FakeDeleteAlso.ValueBool() {
			client.WithFakeDeleteAlso()
		}
	} else if v := os.Getenv(junos.EnvFakedeleteAlso); strings.EqualFold(v, "true") || v == "1" {
		client.WithFakeDeleteAlso()
	}

	if !client.FakeCreateSetFile() &&
		(client.FakeUpdateAlso() || client.FakeDeleteAlso()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("fake_create_with_setfile"),
			"Inconsistency fake attributes",
			"'fake_create_with_setfile' need to be set with 'fake_update_also' and 'fake_delete_also'",
		)

		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func unexpectedResourceConfigureType(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
	resp.Diagnostics.AddError(
		"Unexpected Resource Configure Type",
		fmt.Sprintf(
			"Expected *junos.Client, got: %T. Please report this issue to the provider developers.",
			req.ProviderData,
		),
	)
}

func unexpectedDataSourceConfigureType(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	resp.Diagnostics.AddError(
		"Unexpected Data Source Configure Type",
		fmt.Sprintf(
			"Expected *junos.Client, got: %T. Please report this issue to the provider developers.",
			req.ProviderData,
		),
	)
}
