package providersdk

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Provider junos for terraform.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "This is the target for Netconf session (ip or dns name)." +
					" May also be provided via " + junos.EnvHost + " environment variable.",
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "This is the tcp port for ssh connection." +
					" May also be provided via " + junos.EnvPort + " environment variable.",
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description: "This is the username for ssh connection." +
					" May also be provided via " + junos.EnvUsername + " environment variable.",
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "This is a password for ssh connection." +
					" May also be provided via " + junos.EnvPassword + " environment variable.",
			},
			"sshkey_pem": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "This is the ssh key in PEM format for establish ssh connection." +
					" May also be provided via " + junos.EnvKeyPem + " environment variable.",
			},
			"sshkeyfile": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "This is the path to ssh key for establish ssh connection." +
					" May also be provided via " + junos.EnvKeyFile + " environment variable.",
			},
			"keypass": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "This is the passphrase for open `sshkeyfile` or `sshkey_pem`." +
					" May also be provided via " + junos.EnvKeyPass + " environment variable.",
			},
			"group_interface_delete": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "This is the Junos group used to remove configuration on a physical interface." +
					" May also be provided via " + junos.EnvGroupInterfaceDelete + " environment variable.",
			},
			"cmd_sleep_short": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Milliseconds to wait after Terraform  provider executes an action on the Junos device." +
					" May also be provided via " + junos.EnvSleepShort + " environment variable.",
			},
			"cmd_sleep_lock": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Seconds of standby while waiting for Terraform provider " +
					"to lock candidate configuration on a Junos device." +
					" May also be provided via " + junos.EnvSleepLock + " environment variable.",
			},
			"commit_confirmed": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Number of minutes until automatic rollback." +
					" May also be provided via " + junos.EnvCommitConfirmed + " environment variable." +
					" For each resource action with commit, commit with `confirmed` option and" +
					" with the value ot this argument as `confirm-timeout`, " +
					" wait for `<commit_confirmed_wait_percent>`% of the minutes defined in the value of this argument," +
					" and confirm commit to avoid rollback with the `commit check` command.",
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"commit_confirmed_wait_percent": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Percentage of `<commit_confirmed>` minute(s) to wait between" +
					" `commit confirmed` (commit with automatic rollback) and" +
					" `commit check` (confirmation) commands." +
					" No effect if `<commit_confirmed>` is not used." +
					" May also be provided via " + junos.EnvCommitConfirmedWaitPercent + " environment variable." +
					" Defaults to 90.",
				ValidateFunc: validation.IntBetween(0, 99),
			},
			"ssh_sleep_closed": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Seconds to wait after Terraform provider closed a ssh connection." +
					" May also be provided via " + junos.EnvSleepSSHClosed + " environment variable.",
			},
			"ssh_ciphers": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Ciphers used in SSH connection.",
			},
			"ssh_timeout_to_establish": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Seconds to wait for establishing TCP connections when initiating SSH connections." +
					" May also be provided via " + junos.EnvSSHTimeoutToEstablish + " environment variable.",
			},
			"ssh_retry_to_establish": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Description: "Number of retries to establish SSH connections." +
					"The provider waits after each try, with the sleep time increasing by 1 second each time." +
					" May also be provided via " + junos.EnvSSHRetryToEstablish + " environment variable.",
			},
			"file_permission": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The permission to set for the created file (debug, setfile)." +
					" May also be provided via " + junos.EnvFilePermission + " environment variable.",
			},
			"debug_netconf_log_path": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "More detailed log (netconf) in the specified file." +
					" May also be provided via " + junos.EnvLogPath + " environment variable.",
			},
			"fake_create_with_setfile": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The normal process to create resources skipped to generate set lines, " +
					"append them to the specified file, " +
					"and respond with a `fake` successful creation of resources to Terraform." +
					" May also be provided via " + junos.EnvFakecreateSetfile + " environment variable.",
			},
			"fake_update_also": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "The normal process to update resources skipped to generate set/delete lines, " +
					"append them to the same file as `fake_create_with_setfile`, " +
					"and respond with a `fake` successful update of resources to Terraform." +
					" May also be provided via " + junos.EnvFakeupdateAlso + " environment variable.",
			},
			"fake_delete_also": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "The normal process to delete resources skipped to generate delete lines, " +
					"append them to the same file as `fake_create_with_setfile`, " +
					"and respond with a `fake` successful delete of resources to Terraform." +
					" May also be provided via " + junos.EnvFakedeleteAlso + " environment variable.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"junos_access_address_assignment_pool": resourceAccessAddressAssignPool(),

			"junos_chassis_cluster":    resourceChassisCluster(),
			"junos_chassis_redundancy": resourceChassisRedundancy(),

			"junos_forwardingoptions_dhcprelay":             resourceForwardingOptionsDhcpRelay(),
			"junos_forwardingoptions_dhcprelay_group":       resourceForwardingOptionsDhcpRelayGroup(),
			"junos_forwardingoptions_dhcprelay_servergroup": resourceForwardingOptionsDhcpRelayServerGroup(),

			"junos_group_dual_system": resourceGroupDualSystem(),

			"junos_igmp_snooping_vlan": resourceIgmpSnoopingVlan(),

			"junos_layer2_control": resourceLayer2Control(),

			"junos_lldp_interface":    resourceLldpInterface(),
			"junos_lldpmed_interface": resourceLldpMedInterface(),

			"junos_null_commit_file": resourceNullCommitFile(),

			"junos_rib_group": resourceRibGroup(),

			"junos_routing_options": resourceRoutingOptions(),

			"junos_security_dynamic_address_feed_server":                 resourceSecurityDynamicAddressFeedServer(),
			"junos_security_dynamic_address_name":                        resourceSecurityDynamicAddressName(),
			"junos_security_idp_custom_attack":                           resourceSecurityIdpCustomAttack(),
			"junos_security_idp_custom_attack_group":                     resourceSecurityIdpCustomAttackGroup(),
			"junos_security_idp_policy":                                  resourceSecurityIdpPolicy(),
			"junos_security_screen":                                      resourceSecurityScreen(),
			"junos_security_screen_whitelist":                            resourceSecurityScreenWhiteList(),
			"junos_security_utm_custom_url_category":                     resourceSecurityUtmCustomURLCategory(),
			"junos_security_utm_custom_url_pattern":                      resourceSecurityUtmCustomURLPattern(),
			"junos_security_utm_policy":                                  resourceSecurityUtmPolicy(),
			"junos_security_utm_profile_web_filtering_juniper_enhanced":  resourceSecurityUtmProfileWebFilteringEnhanced(),
			"junos_security_utm_profile_web_filtering_juniper_local":     resourceSecurityUtmProfileWebFilteringLocal(),
			"junos_security_utm_profile_web_filtering_websense_redirect": resourceSecurityUtmProfileWebFilteringWebsense(),

			"junos_services": resourceServices(),

			"junos_services_advanced_anti_malware_policy":                resourceServicesAdvancedAntiMalwarePolicy(),
			"junos_services_proxy_profile":                               resourceServicesProxyProfile(),
			"junos_services_rpm_probe":                                   resourceServicesRpmProbe(),
			"junos_services_ssl_initiation_profile":                      resourceServicesSSLInitiationProfile(),
			"junos_services_security_intelligence_policy":                resourceServicesSecurityIntellPolicy(),
			"junos_services_security_intelligence_profile":               resourceServicesSecurityIntellProfile(),
			"junos_services_user_identification_ad_access_domain":        resourceServicesUserIdentAdAccessDomain(),
			"junos_services_user_identification_device_identity_profile": resourceServicesUserIdentDeviceIdentityProfile(),

			"junos_system_login_class":                     resourceSystemLoginClass(),
			"junos_system_login_user":                      resourceSystemLoginUser(),
			"junos_system_ntp_server":                      resourceSystemNtpServer(),
			"junos_system_root_authentication":             resourceSystemRootAuthentication(),
			"junos_system_services_dhcp_localserver_group": resourceSystemServicesDhcpLocalServerGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"junos_system_information": dataSourceSystemInformation(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	hostIP := os.Getenv(junos.EnvHost)
	if v, ok := d.GetOk("ip"); ok {
		hostIP = v.(string)
	}
	if hostIP == "" {
		return nil, diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing Junos IP target",
			Detail: "The provider cannot create the Junos client as there is a missing or empty value for the Junos IP. " +
				"Set the ip value in the configuration or use the " + junos.EnvHost + " environment variable. " +
				"If either is already set, ensure the value is not empty.",
		}}
	}

	client := junos.NewClient(hostIP)

	var diagWarns diag.Diagnostics

	client.WithPort(830) // default value for port
	if v, ok := d.GetOk("port"); ok {
		client.WithPort(v.(int))
	} else if ev := os.Getenv(junos.EnvPort); ev != "" {
		if port, err := strconv.Atoi(ev); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Error to parse " + junos.EnvPort,
				Detail: fmt.Sprintf("Error to parse value in "+junos.EnvPort+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			client.WithPort(port)
		}
	}

	client.WithUserName("netconf") // default value for username
	if v, ok := d.GetOk("username"); ok {
		client.WithUserName(v.(string))
	} else if ev := os.Getenv(junos.EnvUsername); ev != "" {
		client.WithUserName(ev)
	}

	if v, ok := d.GetOk("password"); ok {
		client.WithPassword(v.(string))
	} else if ev := os.Getenv(junos.EnvPassword); ev != "" {
		client.WithPassword(ev)
	}

	if v, ok := d.GetOk("sshkey_pem"); ok {
		client.WithSSHKeyPEM(v.(string))
	} else if ev := os.Getenv(junos.EnvKeyPem); ev != "" {
		client.WithSSHKeyPEM(ev)
	}

	if v, ok := d.GetOk("sshkeyfile"); ok {
		keyFile := v.(string)
		if err := utils.ReplaceTildeToHomeDir(&keyFile); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Bad value in sshkeyfile",
				Detail: fmt.Sprintf("Error to use value in sshkeyfile attribute: %s\n"+
					"So the attribute is not used", err),
			})
		} else {
			client.WithSSHKeyFile(keyFile)
		}
	} else if ev := os.Getenv(junos.EnvKeyFile); ev != "" {
		if err := utils.ReplaceTildeToHomeDir(&ev); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Bad value in " + junos.EnvKeyFile,
				Detail: fmt.Sprintf("Error to use value in "+junos.EnvKeyFile+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			client.WithSSHKeyFile(ev)
		}
	}

	if v, ok := d.GetOk("keypass"); ok {
		client.WithSSHKeyPassphrase(v.(string))
	} else if ev := os.Getenv(junos.EnvKeyPass); ev != "" {
		client.WithSSHKeyPassphrase(ev)
	}

	if v, ok := d.GetOk("group_interface_delete"); ok {
		client.WithGroupInterfaceDelete(v.(string))
	} else if ev := os.Getenv(junos.EnvGroupInterfaceDelete); ev != "" {
		client.WithGroupInterfaceDelete(ev)
	}

	client.WithSleepShort(100) // default value for cmd_sleep_short
	if v, ok := d.GetOk("cmd_sleep_short"); ok {
		client.WithSleepShort(v.(int))
	} else if ev := os.Getenv(junos.EnvSleepShort); ev != "" {
		if ms, err := strconv.Atoi(ev); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Error to parse " + junos.EnvSleepShort,
				Detail: fmt.Sprintf("Error to parse value in "+junos.EnvSleepShort+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			client.WithSleepShort(ms)
		}
	}

	client.WithSleepLock(10) // default value for cmd_sleep_lock
	if v, ok := d.GetOk("cmd_sleep_lock"); ok {
		client.WithSleepLock(v.(int))
	} else if v := os.Getenv(junos.EnvSleepLock); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Error to parse " + junos.EnvSleepLock,
				Detail: fmt.Sprintf("Error to parse value in "+junos.EnvSleepLock+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			client.WithSleepLock(d)
		}
	}

	if v, ok := d.GetOk("commit_confirmed"); ok {
		if _, err := client.WithCommitConfirmed(v.(int)); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Bad value in commit_confirmed",
				Detail: fmt.Sprintf("Error to use value in commit_confirmed attribute: %s\n"+
					"So the attribute is not used", err),
			})
		}
	} else if v := os.Getenv(junos.EnvCommitConfirmed); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Error to parse " + junos.EnvCommitConfirmed,
				Detail: fmt.Sprintf("Error to parse value in "+junos.EnvCommitConfirmed+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			if _, err := client.WithCommitConfirmed(d); err != nil {
				diagWarns = append(diagWarns, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Bad value in " + junos.EnvCommitConfirmed,
					Detail: fmt.Sprintf("Error to use value in "+junos.EnvCommitConfirmed+" environment variable: %s\n"+
						"So the variable is not used", err),
				})
			}
		}
	}

	_, _ = client.WithCommitConfirmedWaitPercent(90) // default value for commit_confirmed_wait_percent
	if v, ok := d.GetOk("commit_confirmed_wait_percent"); ok {
		if _, err := client.WithCommitConfirmedWaitPercent(v.(int)); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Bad value in commit_confirmed_wait_percent",
				Detail: fmt.Sprintf("Error to use value in commit_confirmed_wait_percent attribute: %s\n"+
					"So the attribute is not used", err),
			})
		}
	} else if v := os.Getenv(junos.EnvCommitConfirmedWaitPercent); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Error to parse " + junos.EnvCommitConfirmedWaitPercent,
				Detail: fmt.Sprintf("Error to parse value in "+junos.EnvCommitConfirmedWaitPercent+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			if _, err := client.WithCommitConfirmedWaitPercent(d); err != nil {
				diagWarns = append(diagWarns, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Bad value in " + junos.EnvCommitConfirmedWaitPercent,
					Detail: fmt.Sprintf("Error to use value in "+junos.EnvCommitConfirmedWaitPercent+" environment variable: %s\n"+
						"So the variable is not used", err),
				})
			}
		}
	}

	if v, ok := d.GetOk("ssh_sleep_closed"); ok {
		client.WithSleepSSHClosed(v.(int))
	} else if v := os.Getenv(junos.EnvSleepSSHClosed); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Error to parse " + junos.EnvSleepSSHClosed,
				Detail: fmt.Sprintf("Error to parse value in "+junos.EnvSleepSSHClosed+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			client.WithSleepSSHClosed(d)
		}
	}

	sshCiphers := make([]string, len(d.Get("ssh_ciphers").([]interface{})))
	for i, v := range d.Get("ssh_ciphers").([]interface{}) {
		sshCiphers[i] = v.(string)
	}
	if len(sshCiphers) == 0 {
		client.WithSSHCiphers(junos.DefaultSSHCiphers())
	} else {
		client.WithSSHCiphers(sshCiphers)
	}

	if v, ok := d.GetOk("ssh_timeout_to_establish"); ok {
		client.WithSSHTimeoutToEstablish(v.(int))
	} else if v := os.Getenv(junos.EnvSSHTimeoutToEstablish); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Error to parse " + junos.EnvSSHTimeoutToEstablish,
				Detail: fmt.Sprintf("Error to parse value in "+junos.EnvSSHTimeoutToEstablish+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			client.WithSSHTimeoutToEstablish(d)
		}
	}

	_, _ = client.WithSSHRetryToEstablish(1) // default value for ssh_retry_to_establish
	if v, ok := d.GetOk("ssh_retry_to_establish"); ok {
		if _, err := client.WithSSHRetryToEstablish(v.(int)); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Bad value in ssh_retry_to_establish",
				Detail: fmt.Sprintf("Error to use value in 'ssh_retry_to_establish' attribute: %s\n"+
					"So the attribute has the default value", err),
			})
		}
	} else if v := os.Getenv(junos.EnvSSHRetryToEstablish); v != "" {
		d, err := strconv.Atoi(v)
		if err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Error to parse " + junos.EnvSSHRetryToEstablish,
				Detail: fmt.Sprintf("Error to parse value in "+junos.EnvSSHRetryToEstablish+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			if _, err := client.WithSSHRetryToEstablish(d); err != nil {
				diagWarns = append(diagWarns, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Bad value in " + junos.EnvSSHRetryToEstablish,
					Detail: fmt.Sprintf("Error to use value in "+junos.EnvSSHRetryToEstablish+" environment variable: %s\n"+
						"So the variable is not used", err),
				})
			}
		}
	}

	_, _ = client.WithFilePermission(0o644) // default value for file_permission
	if v, ok := d.GetOk("file_permission"); ok {
		filePerm, err := strconv.ParseInt(v.(string), 8, 64)
		if err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Error to parse file_permission",
				Detail: fmt.Sprintf("Error to parse value in file_permission attribute: %s\n"+
					"So the attribute has the default value", err),
			})
		} else {
			if _, err := client.WithFilePermission(filePerm); err != nil {
				diagWarns = append(diagWarns, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Bad value in file_permission",
					Detail: fmt.Sprintf("Error to use value in file_permission attribute: %s\n"+
						"So the attribute has the default value", err),
				})
			}
		}
	} else if v := os.Getenv(junos.EnvFilePermission); v != "" {
		filePerm, err := strconv.ParseInt(v, 8, 64)
		if err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Error to parse " + junos.EnvFilePermission,
				Detail: fmt.Sprintf("Error to parse value in "+junos.EnvFilePermission+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			if _, err := client.WithFilePermission(filePerm); err != nil {
				diagWarns = append(diagWarns, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Bad value in " + junos.EnvFilePermission,
					Detail: fmt.Sprintf("Error to use value in "+junos.EnvFilePermission+" environment variable: %s\n"+
						"So the variable is not used", err),
				})
			}
		}
	}

	if v, ok := d.GetOk("debug_netconf_log_path"); ok {
		logPath := v.(string)
		if err := utils.ReplaceTildeToHomeDir(&logPath); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Bad value in debug_netconf_log_path",
				Detail: fmt.Sprintf("Error to use value in debug_netconf_log_path attribute: %s\n"+
					"So the attribute is not used", err),
			})
		} else {
			client.WithDebugLogFile(logPath)
		}
	} else if v := os.Getenv(junos.EnvLogPath); v != "" {
		if err := utils.ReplaceTildeToHomeDir(&v); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Bad value in " + junos.EnvLogPath,
				Detail: fmt.Sprintf("Error to use value in "+junos.EnvLogPath+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			client.WithDebugLogFile(v)
		}
	}

	if v, ok := d.GetOk("fake_create_with_setfile"); ok {
		setFile := v.(string)
		if err := utils.ReplaceTildeToHomeDir(&setFile); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Bad value in fake_create_with_setfile",
				Detail: fmt.Sprintf("Error to use value in fake_create_with_setfile attribute: %s\n"+
					"So the attribute is not used", err),
			})
		} else {
			client.WithFakeCreateSetFile(setFile)
		}
	} else if v := os.Getenv(junos.EnvFakecreateSetfile); v != "" {
		if err := utils.ReplaceTildeToHomeDir(&v); err != nil {
			diagWarns = append(diagWarns, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Bad value in " + junos.EnvFakecreateSetfile,
				Detail: fmt.Sprintf("Error to use value in "+junos.EnvFakecreateSetfile+" environment variable: %s\n"+
					"So the variable is not used", err),
			})
		} else {
			client.WithFakeCreateSetFile(v)
		}
	}

	if v, ok := d.GetOk("fake_update_also"); ok {
		if v.(bool) {
			client.WithFakeUpdateAlso()
		}
	} else if v := os.Getenv(junos.EnvFakeupdateAlso); strings.EqualFold(v, "true") || v == "1" {
		client.WithFakeUpdateAlso()
	}

	if v, ok := d.GetOk("fake_delete_also"); ok {
		if v.(bool) {
			client.WithFakeDeleteAlso()
		}
	} else if v := os.Getenv(junos.EnvFakedeleteAlso); strings.EqualFold(v, "true") || v == "1" {
		client.WithFakeDeleteAlso()
	}

	if !client.FakeCreateSetFile() &&
		(client.FakeUpdateAlso() || client.FakeDeleteAlso()) {
		return client, append(diagWarns, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Inconsistency fake attributes",
			Detail:   "'fake_create_with_setfile' need to be set with 'fake_update_also' and 'fake_delete_also'",
		})
	}

	return client, diagWarns
}
