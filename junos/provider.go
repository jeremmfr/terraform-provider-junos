package junos

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var mutex = &sync.Mutex{} // nolint: gochecknoglobals

// Provider junos for terraform.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_HOST", nil),
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_PORT", 830),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_USERNAME", "netconf"),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_PASSWORD", nil),
			},
			"sshkey_pem": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_KEYPEM", nil),
			},
			"sshkeyfile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_KEYFILE", nil),
			},
			"keypass": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_KEYPASS", nil),
			},
			"group_interface_delete": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_GROUP_INTERFACE_DELETE", nil),
			},
			"cmd_sleep_short": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_SLEEP_SHORT", 100),
			},
			"cmd_sleep_lock": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_SLEEP_LOCK", 10),
			},
			"ssh_sleep_closed": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_SLEEP_SSH_CLOSED", 0),
			},
			"ssh_ciphers": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				DefaultFunc: defaultSSHCiphers(),
			},
			"file_permission": {
				Type:             schema.TypeString,
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc("JUNOS_FILE_PERMISSION", "644"),
				ValidateDiagFunc: validateFilePermission(),
			},
			"debug_netconf_log_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_LOG_PATH", ""),
			},
			"fake_create_with_setfile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_FAKECREATE_SETFILE", ""),
			},
			"fake_update_also": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"fake_create_with_setfile"},
				DefaultFunc:  EnvDefaultBooleanFunc("JUNOS_FAKEUPDATE_ALSO"),
			},
			"fake_delete_also": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"fake_create_with_setfile"},
				DefaultFunc:  EnvDefaultBooleanFunc("JUNOS_FAKEDELETE_ALSO"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"junos_access_address_assignment_pool":                       resourceAccessAddressAssignPool(),
			"junos_aggregate_route":                                      resourceAggregateRoute(),
			"junos_application":                                          resourceApplication(),
			"junos_application_set":                                      resourceApplicationSet(),
			"junos_bgp_group":                                            resourceBgpGroup(),
			"junos_bgp_neighbor":                                         resourceBgpNeighbor(),
			"junos_bridge_domain":                                        resourceBridgeDomain(),
			"junos_chassis_cluster":                                      resourceChassisCluster(),
			"junos_eventoptions_destination":                             resourceEventoptionsDestination(),
			"junos_eventoptions_generate_event":                          resourceEventoptionsGenerateEvent(),
			"junos_eventoptions_policy":                                  resourceEventoptionsPolicy(),
			"junos_evpn":                                                 resourceEvpn(),
			"junos_firewall_filter":                                      resourceFirewallFilter(),
			"junos_firewall_policer":                                     resourceFirewallPolicer(),
			"junos_forwardingoptions_sampling_instance":                  resourceForwardingoptionsSamplingInstance(),
			"junos_generate_route":                                       resourceGenerateRoute(),
			"junos_group_dual_system":                                    resourceGroupDualSystem(),
			"junos_interface":                                            resourceInterface(),
			"junos_interface_logical":                                    resourceInterfaceLogical(),
			"junos_interface_physical":                                   resourceInterfacePhysical(),
			"junos_interface_physical_disable":                           resourceInterfacePhysicalDisable(),
			"junos_interface_st0_unit":                                   resourceInterfaceSt0Unit(),
			"junos_null_commit_file":                                     resourceNullCommitFile(),
			"junos_ospf":                                                 resourceOspf(),
			"junos_ospf_area":                                            resourceOspfArea(),
			"junos_policyoptions_as_path":                                resourcePolicyoptionsAsPath(),
			"junos_policyoptions_as_path_group":                          resourcePolicyoptionsAsPathGroup(),
			"junos_policyoptions_community":                              resourcePolicyoptionsCommunity(),
			"junos_policyoptions_policy_statement":                       resourcePolicyoptionsPolicyStatement(),
			"junos_policyoptions_prefix_list":                            resourcePolicyoptionsPrefixList(),
			"junos_rib_group":                                            resourceRibGroup(),
			"junos_routing_instance":                                     resourceRoutingInstance(),
			"junos_routing_options":                                      resourceRoutingOptions(),
			"junos_security":                                             resourceSecurity(),
			"junos_security_address_book":                                resourceSecurityAddressBook(),
			"junos_security_dynamic_address_feed_server":                 resourceSecurityDynamicAddressFeedServer(),
			"junos_security_dynamic_address_name":                        resourceSecurityDynamicAddressName(),
			"junos_security_global_policy":                               resourceSecurityGlobalPolicy(),
			"junos_security_idp_custom_attack":                           resourceSecurityIdpCustomAttack(),
			"junos_security_idp_custom_attack_group":                     resourceSecurityIdpCustomAttackGroup(),
			"junos_security_idp_policy":                                  resourceSecurityIdpPolicy(),
			"junos_security_ike_gateway":                                 resourceIkeGateway(),
			"junos_security_ike_policy":                                  resourceIkePolicy(),
			"junos_security_ike_proposal":                                resourceIkeProposal(),
			"junos_security_ipsec_policy":                                resourceIpsecPolicy(),
			"junos_security_ipsec_proposal":                              resourceIpsecProposal(),
			"junos_security_ipsec_vpn":                                   resourceIpsecVpn(),
			"junos_security_log_stream":                                  resourceSecurityLogStream(),
			"junos_security_nat_destination":                             resourceSecurityNatDestination(),
			"junos_security_nat_destination_pool":                        resourceSecurityNatDestinationPool(),
			"junos_security_nat_source":                                  resourceSecurityNatSource(),
			"junos_security_nat_source_pool":                             resourceSecurityNatSourcePool(),
			"junos_security_nat_static":                                  resourceSecurityNatStatic(),
			"junos_security_nat_static_rule":                             resourceSecurityNatStaticRule(),
			"junos_security_policy":                                      resourceSecurityPolicy(),
			"junos_security_policy_tunnel_pair_policy":                   resourceSecurityPolicyTunnelPairPolicy(),
			"junos_security_screen":                                      resourceSecurityScreen(),
			"junos_security_screen_whitelist":                            resourceSecurityScreenWhiteList(),
			"junos_security_utm_custom_url_category":                     resourceSecurityUtmCustomURLCategory(),
			"junos_security_utm_custom_url_pattern":                      resourceSecurityUtmCustomURLPattern(),
			"junos_security_utm_policy":                                  resourceSecurityUtmPolicy(),
			"junos_security_utm_profile_web_filtering_juniper_enhanced":  resourceSecurityUtmProfileWebFilteringEnhanced(),
			"junos_security_utm_profile_web_filtering_juniper_local":     resourceSecurityUtmProfileWebFilteringLocal(),
			"junos_security_utm_profile_web_filtering_websense_redirect": resourceSecurityUtmProfileWebFilteringWebsense(),
			"junos_security_zone":                                        resourceSecurityZone(),
			"junos_security_zone_book_address":                           resourceSecurityZoneBookAddress(),
			"junos_security_zone_book_address_set":                       resourceSecurityZoneBookAddressSet(),
			"junos_services":                                             resourceServices(),
			"junos_services_advanced_anti_malware_policy":                resourceServicesAdvancedAntiMalwarePolicy(),
			"junos_services_flowmonitoring_vipfix_template":              resourceServicesFlowMonitoringVIPFixTemplate(),
			"junos_services_proxy_profile":                               resourceServicesProxyProfile(),
			"junos_services_rpm_probe":                                   resourceServicesRpmProbe(),
			"junos_services_ssl_initiation_profile":                      resourceServicesSSLInitiationProfile(),
			"junos_services_security_intelligence_policy":                resourceServicesSecurityIntellPolicy(),
			"junos_services_security_intelligence_profile":               resourceServicesSecurityIntellProfile(),
			"junos_services_user_identification_ad_access_domain":        resourceServicesUserIdentAdAccessDomain(),
			"junos_services_user_identification_device_identity_profile": resourceServicesUserIdentDeviceIdentityProfile(),
			"junos_snmp":                                                 resourceSnmp(),
			"junos_snmp_clientlist":                                      resourceSnmpClientlist(),
			"junos_snmp_community":                                       resourceSnmpCommunity(),
			"junos_snmp_v3_usm_user":                                     resourceSnmpV3UsmUser(),
			"junos_snmp_v3_vacm_accessgroup":                             resourceSnmpV3VacmAccessGroup(),
			"junos_snmp_v3_vacm_securitytogroup":                         resourceSnmpV3VacmSecurityToGroup(),
			"junos_snmp_view":                                            resourceSnmpView(),
			"junos_static_route":                                         resourceStaticRoute(),
			"junos_switch_options":                                       resourceSwitchOptions(),
			"junos_system":                                               resourceSystem(),
			"junos_system_login_class":                                   resourceSystemLoginClass(),
			"junos_system_login_user":                                    resourceSystemLoginUser(),
			"junos_system_ntp_server":                                    resourceSystemNtpServer(),
			"junos_system_radius_server":                                 resourceSystemRadiusServer(),
			"junos_system_root_authentication":                           resourceSystemRootAuthentication(),
			"junos_system_services_dhcp_localserver_group":               resourceSystemServicesDhcpLocalServerGroup(),
			"junos_system_syslog_file":                                   resourceSystemSyslogFile(),
			"junos_system_syslog_host":                                   resourceSystemSyslogHost(),
			"junos_vlan":                                                 resourceVlan(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"junos_interface":                   dataSourceInterface(),
			"junos_interface_logical":           dataSourceInterfaceLogical(),
			"junos_interface_physical":          dataSourceInterfacePhysical(),
			"junos_interfaces_physical_present": dataSourceInterfacesPhysicalPresent(),
			"junos_system_information":          dataSourceSystemInformation(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	if d.Get("fake_update_also").(bool) || d.Get("fake_delete_also").(bool) {
		if d.Get("fake_create_with_setfile").(string) == "" {
			return nil, diag.FromErr(fmt.Errorf(
				"'fake_create_with_setfile' need to be set with 'fake_update_also' and 'fake_delete_also'"))
		}
	}
	c := configProvider{
		junosIP:                  d.Get("ip").(string),
		junosPort:                d.Get("port").(int),
		junosUserName:            d.Get("username").(string),
		junosPassword:            d.Get("password").(string),
		junosSSHKeyPEM:           d.Get("sshkey_pem").(string),
		junosSSHKeyFile:          d.Get("sshkeyfile").(string),
		junosKeyPass:             d.Get("keypass").(string),
		junosGroupIntDel:         d.Get("group_interface_delete").(string),
		junosCmdSleepShort:       d.Get("cmd_sleep_short").(int),
		junosCmdSleepLock:        d.Get("cmd_sleep_lock").(int),
		junosSSHSleepClosed:      d.Get("ssh_sleep_closed").(int),
		junosFilePermission:      d.Get("file_permission").(string),
		junosDebugNetconfLogPath: d.Get("debug_netconf_log_path").(string),
		junosFakeCreateSetFile:   d.Get("fake_create_with_setfile").(string),
		junosFakeUpdateAlso:      d.Get("fake_update_also").(bool),
		junosFakeDeleteAlso:      d.Get("fake_delete_also").(bool),
	}
	for _, v := range d.Get("ssh_ciphers").([]interface{}) {
		c.junosSSHCiphers = append(c.junosSSHCiphers, v.(string))
	}

	return c.prepareSession()
}

func EnvDefaultBooleanFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); strings.ToLower(v) == "true" {
			return true, nil
		}

		return false, nil
	}
}
