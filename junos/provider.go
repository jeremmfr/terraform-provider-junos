package junos

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	idSeparator        = "_-_"
	defaultWord        = "default"
	inetWord           = "inet"
	inet6Word          = "inet6"
	emptyWord          = "empty"
	matchWord          = "match"
	permitWord         = "permit"
	thenWord           = "then"
	prefixWord         = "prefix"
	actionNoneWord     = "none"
	addWord            = "add"
	deleteWord         = "delete"
	setWord            = "set"
	setLineStart       = setWord + " "
	st0Word            = "st0"
	opsfV2             = "ospf"
	ospfV3             = "ospf3"
	activeW            = "active"
	passiveW           = "passive"
	discardW           = "discard"
	disableW           = "disable"
	dynamicDB          = "dynamic-db"
	preemptWord        = "preempt"
	flowControlWords   = "flow-control"
	noFlowControlWords = "no-flow-control"
	loopbackWord       = "loopback"
	noLoopbackWord     = "no-loopback"
)

var (
	mutex = &sync.Mutex{}
)

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
			"debug_netconf_log_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUNOS_LOG_PATH", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"junos_aggregate_route":                                      resourceAggregateRoute(),
			"junos_application":                                          resourceApplication(),
			"junos_application_set":                                      resourceApplicationSet(),
			"junos_bgp_group":                                            resourceBgpGroup(),
			"junos_bgp_neighbor":                                         resourceBgpNeighbor(),
			"junos_chassis_cluster":                                      resourceChassisCluster(),
			"junos_firewall_filter":                                      resourceFirewallFilter(),
			"junos_firewall_policer":                                     resourceFirewallPolicer(),
			"junos_interface":                                            resourceInterface(),
			"junos_interface_logical":                                    resourceInterfaceLogical(),
			"junos_interface_physical":                                   resourceInterfacePhysical(),
			"junos_interface_st0_unit":                                   resourceInterfaceSt0Unit(),
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
			"junos_static_route":                                         resourceStaticRoute(),
			"junos_system":                                               resourceSystem(),
			"junos_system_login_class":                                   resourceSystemLoginClass(),
			"junos_system_login_user":                                    resourceSystemLoginUser(),
			"junos_system_ntp_server":                                    resourceSystemNtpServer(),
			"junos_system_radius_server":                                 resourceSystemRadiusServer(),
			"junos_system_root_authentication":                           resourceSystemRootAuthentication(),
			"junos_system_syslog_file":                                   resourceSystemSyslogFile(),
			"junos_system_syslog_host":                                   resourceSystemSyslogHost(),
			"junos_vlan":                                                 resourceVlan(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"junos_interface":          dataSourceInterface(),
			"junos_interface_logical":  dataSourceInterfaceLogical(),
			"junos_interface_physical": dataSourceInterfacePhysical(),
			"junos_system_information": dataSourceSystemInformation(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
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
		junosDebugNetconfLogPath: d.Get("debug_netconf_log_path").(string),
	}

	return config.Session()
}
