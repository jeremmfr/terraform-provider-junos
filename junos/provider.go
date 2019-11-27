package junos

import (
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	idSeparator    = "_-_"
	defaultWord    = "default"
	inetWord       = "inet"
	inet6Word      = "inet6"
	emptyWord      = "empty"
	matchWord      = "match"
	thenWord       = "then"
	prefixWord     = "prefix"
	actionNoneWord = "none"
	setLineStart   = "set "
	st0Word        = "st0"
	opsfV2         = "ospf"
	ospfV3         = "ospf3"
)

var (
	mutex = &sync.Mutex{}
)

// Provider junos for terraform
func Provider() terraform.ResourceProvider {
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
			"sshkeyfile": {
				Type:        schema.TypeString,
				Required:    true,
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"junos_application_set":                    resourceApplicationSet(),
			"junos_application":                        resourceApplication(),
			"junos_bgp_group":                          resourceBgpGroup(),
			"junos_bgp_neighbor":                       resourceBgpNeighbor(),
			"junos_firewall_filter":                    resourceFirewallFilter(),
			"junos_firewall_policer":                   resourceFirewallPolicer(),
			"junos_interface":                          resourceInterface(),
			"junos_ospf_area":                          resourceOspfArea(),
			"junos_policyoptions_as_path_group":        resourcePolicyoptionsAsPathGroup(),
			"junos_policyoptions_as_path":              resourcePolicyoptionsAsPath(),
			"junos_policyoptions_community":            resourcePolicyoptionsCommunity(),
			"junos_policyoptions_policy_statement":     resourcePolicyoptionsPolicyStatement(),
			"junos_policyoptions_prefix_list":          resourcePolicyoptionsPrefixList(),
			"junos_rib_group":                          resourceRibGroup(),
			"junos_routing_instance":                   resourceRoutingInstance(),
			"junos_security_ike_gateway":               resourceIkeGateway(),
			"junos_security_ike_policy":                resourceIkePolicy(),
			"junos_security_ike_proposal":              resourceIkeProposal(),
			"junos_security_ipsec_policy":              resourceIpsecPolicy(),
			"junos_security_ipsec_proposal":            resourceIpsecProposal(),
			"junos_security_ipsec_vpn":                 resourceIpsecVpn(),
			"junos_security_nat_destination_pool":      resourceSecurityNatDestinationPool(),
			"junos_security_nat_destination":           resourceSecurityNatDestination(),
			"junos_security_nat_source_pool":           resourceSecurityNatSourcePool(),
			"junos_security_nat_source":                resourceSecurityNatSource(),
			"junos_security_nat_static":                resourceSecurityNatStatic(),
			"junos_security_policy_tunnel_pair_policy": resourceSecurityPolicyTunnelPairPolicy(),
			"junos_security_policy":                    resourceSecurityPolicy(),
			"junos_security_zone":                      resourceSecurityZone(),
			"junos_static_route":                       resourceStaticRoute(),
			"junos_vlan":                               resourceVlan(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		junosIP:          d.Get("ip").(string),
		junosPort:        d.Get("port").(int),
		junosUserName:    d.Get("username").(string),
		junosSSHKeyFile:  d.Get("sshkeyfile").(string),
		junosKeyPass:     d.Get("keypass").(string),
		junosGroupIntDel: d.Get("group_interface_delete").(string),
	}
	return config.Session()
}
