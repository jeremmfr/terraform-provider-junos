package providerfwk_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccJunosSecurityIkeIpsec_basic(t *testing.T) {
	testaccIkeIpsec := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccIkeIpsec = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityIkeIpsecConfigCreate(testaccIkeIpsec),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_proposal.testacc_ikeprop",
							"authentication_algorithm", "sha1"),
						resource.TestCheckResourceAttr("junos_security_ike_proposal.testacc_ikeprop",
							"encryption_algorithm", "aes-256-cbc"),
						resource.TestCheckResourceAttr("junos_security_ike_proposal.testacc_ikeprop",
							"authentication_method", "pre-shared-keys"),
						resource.TestCheckResourceAttr("junos_security_ike_proposal.testacc_ikeprop",
							"dh_group", "group2"),
						resource.TestCheckResourceAttr("junos_security_ike_proposal.testacc_ikeprop",
							"lifetime_seconds", "3600"),
						resource.TestCheckResourceAttr("junos_security_ike_policy.testacc_ikepol",
							"proposals.#", "1"),
						resource.TestCheckResourceAttr("junos_security_ike_policy.testacc_ikepol",
							"proposals.0", "testacc_ikeprop"),
						resource.TestCheckResourceAttr("junos_security_ike_policy.testacc_ikepol",
							"mode", "main"),
						resource.TestCheckResourceAttr("junos_security_ike_policy.testacc_ikepol",
							"pre_shared_key_text", "thePassWord"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"address.0", "192.0.2.3"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"policy", "testacc_ikepol"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"external_interface", testaccIkeIpsec+".0"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"general_ike_id", "true"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"no_nat_traversal", "true"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dead_peer_detection.interval", "10"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dead_peer_detection.threshold", "3"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dead_peer_detection.send_mode", "always-send"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"local_address", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"local_identity.type", "hostname"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"local_identity.value", "testacc"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"version", "v2-only"),
						resource.TestCheckResourceAttr("junos_security_ipsec_proposal.testacc_ipsecprop",
							"authentication_algorithm", "hmac-sha1-96"),
						resource.TestCheckResourceAttr("junos_security_ipsec_proposal.testacc_ipsecprop",
							"protocol", "esp"),
						resource.TestCheckResourceAttr("junos_security_ipsec_proposal.testacc_ipsecprop",
							"encryption_algorithm", "aes-128-cbc"),
						resource.TestCheckResourceAttr("junos_security_ipsec_policy.testacc_ipsecpol",
							"proposals.#", "1"),
						resource.TestCheckResourceAttr("junos_security_ipsec_policy.testacc_ipsecpol",
							"proposals.0", "testacc_ipsecprop"),
						resource.TestCheckResourceAttr("junos_security_ipsec_policy.testacc_ipsecpol",
							"pfs_keys", "group2"),
						resource.TestMatchResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"bind_interface", regexp.MustCompile(`^st0\.\d+$`)),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"ike.gateway", "testacc_ikegateway"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"ike.policy", "testacc_ipsecpol"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"ike.identity_local", "192.0.2.64/26"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"ike.identity_remote", "192.0.2.128/26"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"ike.identity_service", "any"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"vpn_monitor.destination_ip", "192.0.2.129"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"vpn_monitor.optimized", "true"),
						resource.TestMatchResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"vpn_monitor.source_interface", regexp.MustCompile(`^st0\.\d+$`)),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"establish_tunnels", "on-traffic"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"df_bit", "clear"),
					),
				},
				{
					Config: testAccJunosSecurityIkeIpsecConfigUpdate(testaccIkeIpsec),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_proposal.testacc_ikeprop",
							"dh_group", "group1"),
						resource.TestCheckResourceAttr("junos_security_ike_policy.testacc_ikepol",
							"pre_shared_key_text", "mysecret"),
						resource.TestCheckResourceAttr("junos_security_ike_policy.testacc_ikepol",
							"proposal_set", "standard"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"address.0", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_security_ipsec_proposal.testacc_ipsecprop",
							"encryption_algorithm", "aes-256-cbc"),
						resource.TestCheckResourceAttr("junos_security_ipsec_policy.testacc_ipsecpol",
							"pfs_keys", "group1"),
						resource.TestCheckResourceAttr("junos_security_ipsec_policy.testacc_ipsecpol",
							"proposal_set", "standard"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_policyIpsecRemToLoc",
							"policy.#", "1"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_policyIpsecRemToLoc",
							"policy.0.permit_tunnel_ipsec_vpn", "testacc_ipsecvpn2"),
						resource.TestCheckResourceAttr("junos_security_policy_tunnel_pair_policy.testacc_vpn-in-out",
							"policy_a_to_b", "testacc_vpn-out"),
						resource.TestCheckResourceAttr("junos_security_policy_tunnel_pair_policy.testacc_vpn-in-out",
							"policy_b_to_a", "testacc_vpn-in"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"remote_identity.type", "hostname"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"remote_identity.value", "testacc_remote"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dead_peer_detection.send_mode", "optimized"),
					),
				},
				{
					ResourceName:      "junos_security_ike_proposal.testacc_ikeprop",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_ike_policy.testacc_ikepol",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_ike_gateway.testacc_ikegateway",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_ipsec_proposal.testacc_ipsecprop",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_ipsec_policy.testacc_ipsecpol",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_security_ipsec_vpn.testacc_ipsecvpn2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosSecurityIkeIpsecConfigUpdate2(testaccIkeIpsec),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dead_peer_detection.send_mode", "probe-idle-tunnel"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.distinguished_name.container", "dc=example,dc=com"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.connections_limit", "10"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"aaa.client_username", "user"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"aaa.client_password", "password"),
					),
				},
				{
					Config: testAccJunosSecurityIkeIpsecConfigUpdate3(testaccIkeIpsec),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.hostname", "host1.example.com"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"traffic_selector.#", "2"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"traffic_selector.0.name", "ts-1"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"traffic_selector.0.local_ip", "192.0.2.0/26"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"traffic_selector.0.remote_ip", "192.0.3.64/26"),
					),
				},
				{
					Config: testAccJunosSecurityIkeIpsecConfigUpdate4(testaccIkeIpsec),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.inet", "192.168.0.4"),
					),
				},
				{
					Config: testAccJunosSecurityIkeIpsecConfigUpdate5(testaccIkeIpsec),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.inet6", "2001:db8::1"),
					),
				},
				{
					Config: testAccJunosSecurityIkeIpsecConfigUpdate6(testaccIkeIpsec),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.ike_user_type", "group-ike-id"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.reject_duplicate_connection", "true"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.user_at_hostname", "user@example.com"),
					),
				},
				{
					Config: testAccJunosSecurityIkeIpsecConfigUpdate7(testaccIkeIpsec),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.distinguished_name.wildcard", "*.com"),
					),
				},
				{
					Config: testAccJunosSecurityIkeIpsecConfigUpdate8(testaccIkeIpsec),
				},
				{
					Config: testAccJunosSecurityIkeIpsecConfigUpdate9(testaccIkeIpsec),
				},
			},
		})
	}
}

func testAccJunosSecurityIkeIpsecConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group2"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "main"
  pre_shared_key_text = "thePassWord"
  reauth_frequency    = 50
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name               = "testacc_ikegateway"
  address            = ["192.0.2.3"]
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
  general_ike_id     = true
  no_nat_traversal   = true
  dead_peer_detection {
    interval  = 10
    threshold = 3
    send_mode = "always-send"
  }
  local_address = "192.0.2.4"
  local_identity {
    type  = "hostname"
    value = "testacc"
  }
  version = "v2-only"
}

resource "junos_security_ipsec_proposal" "testacc_ipsecprop" {
  name                     = "testacc_ipsecprop"
  authentication_algorithm = "hmac-sha1-96"
  protocol                 = "esp"
  encryption_algorithm     = "aes-128-cbc"
}
resource "junos_security_ipsec_policy" "testacc_ipsecpol" {
  name      = "testacc_ipsecpol"
  proposals = [junos_security_ipsec_proposal.testacc_ipsecprop.name]
  pfs_keys  = "group2"
}
resource "junos_interface_st0_unit" "testacc_ipsecvpn" {}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn" {
  name           = "testacc_ipsecvpn"
  bind_interface = junos_interface_st0_unit.testacc_ipsecvpn.id
  ike {
    gateway          = junos_security_ike_gateway.testacc_ikegateway.name
    policy           = junos_security_ipsec_policy.testacc_ipsecpol.name
    identity_local   = "192.0.2.64/26"
    identity_remote  = "192.0.2.128/26"
    identity_service = "any"
  }
  vpn_monitor {
    destination_ip        = "192.0.2.129"
    optimized             = true
    source_interface_auto = true
  }
  establish_tunnels = "on-traffic"
  df_bit            = "clear"
}
`, interFace)
}

func testAccJunosSecurityIkeIpsecConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group1"
  lifetime_seconds         = 3600
  description              = "testacc ikeprop"
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposal_set        = "standard"
  mode                = "main"
  pre_shared_key_text = "mysecret"
  description         = "testacc ikepol"
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name               = "testacc_ikegateway"
  address            = ["192.0.2.4"]
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
  no_nat_traversal   = true
  dead_peer_detection {
    interval  = 10
    threshold = 3
    send_mode = "optimized"
  }
  local_address = "192.0.2.4"
  local_identity {
    type  = "hostname"
    value = "testacc"
  }
  remote_identity {
    type  = "hostname"
    value = "testacc_remote"
  }
}

resource "junos_security_ipsec_proposal" "testacc_ipsecprop" {
  name                     = "testacc_ipsecprop"
  authentication_algorithm = "hmac-sha1-96"
  protocol                 = "esp"
  encryption_algorithm     = "aes-256-cbc"
  description              = "testacc ipsecprop"
}
resource "junos_security_ipsec_policy" "testacc_ipsecpol" {
  name         = "testacc_ipsecpol"
  proposal_set = "standard"
  pfs_keys     = "group1"
  description  = "testacc ipsecpol"
}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn2" {
  name = "testacc_ipsecvpn2"
  ike {
    gateway          = junos_security_ike_gateway.testacc_ikegateway.name
    policy           = junos_security_ipsec_policy.testacc_ipsecpol.name
    identity_local   = "192.0.2.64/26"
    identity_remote  = "192.0.2.128/26"
    identity_service = "any"
  }
  establish_tunnels = "immediately"
  df_bit            = "clear"
}
resource "junos_security_zone" "testacc_secIkeIpsec_local" {
  name = "testacc_secIkeIPsec_local"
  address_book {
    name    = "testacc_vpnlocal"
    network = "192.0.2.64/26"
  }
}
resource "junos_security_zone" "testacc_secIkeIpsec_remote" {
  name = "testacc_secIkeIPsec_remote"
  address_book {
    name    = "testacc_vpnremote"
    network = "192.0.2.128/26"
  }
}
resource "junos_security_policy" "testacc_policyIpsecLocToRem" {
  depends_on = [
    junos_security_zone.testacc_secIkeIpsec_local,
    junos_security_zone.testacc_secIkeIpsec_remote
  ]
  from_zone = junos_security_zone.testacc_secIkeIpsec_local.name
  to_zone   = junos_security_zone.testacc_secIkeIpsec_remote.name
  policy {
    name                      = "testacc_vpn-out"
    match_source_address      = ["testacc_vpnlocal"]
    match_destination_address = ["testacc_vpnremote"]
    match_application         = ["any"]
    permit_tunnel_ipsec_vpn   = junos_security_ipsec_vpn.testacc_ipsecvpn2.name
  }
}

resource "junos_security_policy" "testacc_policyIpsecRemToLoc" {
  depends_on = [
    junos_security_zone.testacc_secIkeIpsec_local,
    junos_security_zone.testacc_secIkeIpsec_remote
  ]
  from_zone = junos_security_zone.testacc_secIkeIpsec_remote.name
  to_zone   = junos_security_zone.testacc_secIkeIpsec_local.name
  policy {
    name                      = "testacc_vpn-in"
    match_source_address      = ["testacc_vpnremote"]
    match_destination_address = ["testacc_vpnlocal"]
    match_application         = ["any"]
    permit_tunnel_ipsec_vpn   = junos_security_ipsec_vpn.testacc_ipsecvpn2.name
  }
}

resource "junos_security_policy_tunnel_pair_policy" "testacc_vpn-in-out" {
  zone_a        = junos_security_zone.testacc_secIkeIpsec_local.name
  zone_b        = junos_security_zone.testacc_secIkeIpsec_remote.name
  policy_a_to_b = junos_security_policy.testacc_policyIpsecLocToRem.policy[0].name
  policy_b_to_a = junos_security_policy.testacc_policyIpsecRemToLoc.policy[0].name
}
`, interFace)
}

func testAccJunosSecurityIkeIpsecConfigUpdate2(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group1"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "aggressive"
  pre_shared_key_text = "mysecret"
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name = "testacc_ikegateway"
  dynamic_remote {
    distinguished_name {
      container = "dc=example,dc=com"
    }
    connections_limit = 10
  }
  aaa {
    client_username = "user"
    client_password = "password"
  }
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
  no_nat_traversal   = true
  dead_peer_detection {
    interval  = 10
    threshold = 3
    send_mode = "probe-idle-tunnel"
  }
  local_address = "192.0.2.4"
}

resource "junos_security_ipsec_proposal" "testacc_ipsecprop" {
  name                     = "testacc_ipsecprop"
  authentication_algorithm = "hmac-sha1-96"
  protocol                 = "esp"
  encryption_algorithm     = "aes-256-cbc"
}
resource "junos_security_ipsec_policy" "testacc_ipsecpol" {
  name      = "testacc_ipsecpol"
  proposals = [junos_security_ipsec_proposal.testacc_ipsecprop.name]
  pfs_keys  = "group1"
}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn2" {
  name = "testacc_ipsecvpn2"
  ike {
    gateway          = junos_security_ike_gateway.testacc_ikegateway.name
    policy           = junos_security_ipsec_policy.testacc_ipsecpol.name
    identity_local   = "192.0.2.64/26"
    identity_remote  = "192.0.2.128/26"
    identity_service = "any"
  }
  establish_tunnels = "immediately"
  df_bit            = "clear"
}
resource "junos_security_zone" "testacc_secIkeIpsec_local" {
  name = "testacc_secIkeIPsec_local"
  address_book {
    name    = "testacc_vpnlocal"
    network = "192.0.2.64/26"
  }
}
resource "junos_security_zone" "testacc_secIkeIpsec_remote" {
  name = "testacc_secIkeIPsec_remote"
  address_book {
    name    = "testacc_vpnremote"
    network = "192.0.2.128/26"
  }
}
resource "junos_security_policy" "testacc_policyIpsecLocToRem" {
  depends_on = [
    junos_security_zone.testacc_secIkeIpsec_local,
    junos_security_zone.testacc_secIkeIpsec_remote
  ]
  from_zone = junos_security_zone.testacc_secIkeIpsec_local.name
  to_zone   = junos_security_zone.testacc_secIkeIpsec_remote.name
  policy {
    name                      = "testacc_vpn-out"
    match_source_address      = ["testacc_vpnlocal"]
    match_destination_address = ["testacc_vpnremote"]
    match_application         = ["any"]
    permit_tunnel_ipsec_vpn   = junos_security_ipsec_vpn.testacc_ipsecvpn2.name
  }
}

resource "junos_security_policy" "testacc_policyIpsecRemToLoc" {
  depends_on = [
    junos_security_zone.testacc_secIkeIpsec_local,
    junos_security_zone.testacc_secIkeIpsec_remote
  ]
  from_zone = junos_security_zone.testacc_secIkeIpsec_remote.name
  to_zone   = junos_security_zone.testacc_secIkeIpsec_local.name
  policy {
    name                      = "testacc_vpn-in"
    match_source_address      = ["testacc_vpnremote"]
    match_destination_address = ["testacc_vpnlocal"]
    match_application         = ["any"]
    permit_tunnel_ipsec_vpn   = junos_security_ipsec_vpn.testacc_ipsecvpn2.name
  }
}

resource "junos_security_policy_tunnel_pair_policy" "testacc_vpn-in-out" {
  zone_a        = junos_security_zone.testacc_secIkeIpsec_local.name
  zone_b        = junos_security_zone.testacc_secIkeIpsec_remote.name
  policy_a_to_b = junos_security_policy.testacc_policyIpsecLocToRem.policy[0].name
  policy_b_to_a = junos_security_policy.testacc_policyIpsecRemToLoc.policy[0].name
}
`, interFace)
}

func testAccJunosSecurityIkeIpsecConfigUpdate3(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group1"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "aggressive"
  pre_shared_key_text = "mysecret"
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name = "testacc_ikegateway"
  dynamic_remote {
    hostname = "host1.example.com"
  }
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
}
resource "junos_security_ipsec_proposal" "testacc_ipsecprop" {
  name                     = "testacc_ipsecprop"
  authentication_algorithm = "hmac-sha1-96"
  protocol                 = "esp"
  encryption_algorithm     = "aes-128-cbc"
}
resource "junos_security_ipsec_policy" "testacc_ipsecpol" {
  name      = "testacc_ipsecpol"
  proposals = [junos_security_ipsec_proposal.testacc_ipsecprop.name]
  pfs_keys  = "group2"
}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn" {
  lifecycle {
    create_before_destroy = true
  }

  name           = "testacc_ipsecvpn"
  bind_interface = junos_interface_logical.testacc_ipsecvpn_bind.name
  ike {
    gateway = junos_security_ike_gateway.testacc_ikegateway.name
    policy  = junos_security_ipsec_policy.testacc_ipsecpol.name
  }
  establish_tunnels = "on-traffic"
  traffic_selector {
    name      = "ts-1"
    local_ip  = "192.0.2.0/26"
    remote_ip = "192.0.3.64/26"
  }
  traffic_selector {
    name      = "ts-2"
    local_ip  = "192.0.2.128/26"
    remote_ip = "192.0.3.192/26"
  }
  udp_encapsulate {
    dest_port = "1025"
  }
}
resource "junos_interface_logical" "testacc_ipsecvpn_bind" {
  name = junos_interface_st0_unit.testacc_ipsec_vpn.id
  family_inet {}
}
resource "junos_interface_st0_unit" "testacc_ipsec_vpn" {}
`, interFace)
}

func testAccJunosSecurityIkeIpsecConfigUpdate4(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group1"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "aggressive"
  pre_shared_key_text = "mysecret"
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name = "testacc_ikegateway"
  dynamic_remote {
    inet = "192.168.0.4"
  }
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn" {
  name = "testacc_ipsecvpn"
  manual {
    external_interface       = junos_interface_logical.testacc_ikegateway.name
    protocol                 = "esp"
    spi                      = 256
    authentication_algorithm = "hmac-sha-256-128"
    authentication_key_text  = "AuthenticationKey123456789012345"
    encryption_algorithm     = "aes-256-gcm"
    encryption_key_text      = "Encryp"
  }
}
`, interFace)
}

func testAccJunosSecurityIkeIpsecConfigUpdate5(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group1"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "aggressive"
  pre_shared_key_text = "mysecret"
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name = "testacc_ikegateway"
  dynamic_remote {
    inet6 = "2001:db8::1"
  }
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
}
resource "junos_security_ipsec_vpn" "testacc_ipsecvpn2" {
  name            = "testacc_ipsecvpn2"
  copy_outer_dscp = true
  manual {
    external_interface       = junos_interface_logical.testacc_ikegateway.name
    protocol                 = "esp"
    spi                      = 500
    authentication_algorithm = "hmac-sha-256-128"
    authentication_key_hexa  = "00112233445566778899AABBCCDDEEFFaabbccddeeff00112233445566778899"
    encryption_algorithm     = "aes-256-gcm"
    encryption_key_hexa      = "00112233445566778899AABBCCDDEEFFaabbccddeeff00112233445566778899"
    gateway                  = "192.0.2.128"
  }
  multi_sa_forwarding_class = ["network-control", "best-effort"]
  df_bit                    = "clear"
  udp_encapsulate {}
}
`, interFace)
}

func testAccJunosSecurityIkeIpsecConfigUpdate6(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group1"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "aggressive"
  pre_shared_key_text = "mysecret"
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name = "testacc_ikegateway"
  dynamic_remote {
    ike_user_type               = "group-ike-id"
    reject_duplicate_connection = true
    user_at_hostname            = "user@example.com"
  }
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
}
`, interFace)
}

func testAccJunosSecurityIkeIpsecConfigUpdate7(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group1"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "aggressive"
  pre_shared_key_text = "mysecret"
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name = "testacc_ikegateway"
  dynamic_remote {
    distinguished_name {
      wildcard = "*.com"
    }
  }
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
}
`, interFace)
}

func testAccJunosSecurityIkeIpsecConfigUpdate8(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group1"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "aggressive"
  pre_shared_key_text = "mysecret"
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name               = "testacc_ikegateway"
  address            = ["192.0.2.4"]
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
  local_identity {
    type = "distinguished-name"
  }
  remote_identity {
    type = "distinguished-name"
  }
}
`, interFace)
}

func testAccJunosSecurityIkeIpsecConfigUpdate9(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_logical" "testacc_ikegateway" {
  name = "%s.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.4/25"
    }
  }
}
resource "junos_security_ike_proposal" "testacc_ikeprop" {
  name                     = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm     = "aes-256-cbc"
  dh_group                 = "group1"
  lifetime_seconds         = 3600
}
resource "junos_security_ike_policy" "testacc_ikepol" {
  name                = "testacc_ikepol"
  proposals           = [junos_security_ike_proposal.testacc_ikeprop.name]
  mode                = "aggressive"
  pre_shared_key_text = "mysecret"
}
resource "junos_security_ike_gateway" "testacc_ikegateway" {
  name               = "testacc_ikegateway"
  address            = ["192.0.2.4"]
  policy             = junos_security_ike_policy.testacc_ikepol.name
  external_interface = junos_interface_logical.testacc_ikegateway.name
  remote_identity {
    type                         = "distinguished-name"
    distinguished_name_container = "testacc1"
    distinguished_name_wildcard  = "testacc2"
  }
}
`, interFace)
}
