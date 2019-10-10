package junos

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccJunosSecurityIkeIpsec_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityIkeIpsecConfigCreate(),
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
							"external_interface", "lo0"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"general_ike_id", "true"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"no_nat_traversal", "true"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dead_peer_detection.#", "1"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dead_peer_detection.0.interval", "10"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dead_peer_detection.0.threshold", "3"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"local_address", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"local_identity.#", "1"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"local_identity.0.type", "hostname"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"local_identity.0.value", "testacc"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"remote_identity.#", "1"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"remote_identity.0.type", "hostname"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"remote_identity.0.value", "testacc_remote"),
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
							"ike.#", "1"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"ike.0.gateway", "testacc_ikegateway"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"ike.0.policy", "testacc_ipsecpol"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"ike.0.identity_local", "192.0.2.64/26"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"ike.0.identity_remote", "192.0.2.128/26"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"ike.0.identity_service", "any"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"vpn_monitor.#", "1"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"vpn_monitor.0.destination_ip", "192.0.2.129"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"vpn_monitor.0.optimized", "true"),
						resource.TestMatchResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"vpn_monitor.0.source_interface", regexp.MustCompile(`^st0\.\d+$`)),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"establish_tunnels", "on-traffic"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"df_bit", "clear"),
					),
				},
				{
					Config: testAccJunosSecurityIkeIpsecConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_proposal.testacc_ikeprop",
							"dh_group", "group1"),
						resource.TestCheckResourceAttr("junos_security_ike_policy.testacc_ikepol",
							"pre_shared_key_text", "mysecret"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"address.#", "1"),
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"address.0", "192.0.2.4"),
						resource.TestCheckResourceAttr("junos_security_ipsec_proposal.testacc_ipsecprop",
							"encryption_algorithm", "aes-256-cbc"),
						resource.TestCheckResourceAttr("junos_security_ipsec_policy.testacc_ipsecpol",
							"pfs_keys", "group1"),
						resource.TestCheckResourceAttr("junos_security_ipsec_vpn.testacc_ipsecvpn",
							"bind_interface", ""),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_policyIpsecRemToLoc",
							"policy.#", "1"),
						resource.TestCheckResourceAttr("junos_security_policy.testacc_policyIpsecRemToLoc",
							"policy.0.permit_tunnel_ipsec_vpn", "testacc_ipsecvpn"),
						resource.TestCheckResourceAttr("junos_security_policy_pair_policy.testacc_vpn-in-out",
							"policy_a_to_b", "testacc_vpn-out"),
						resource.TestCheckResourceAttr("junos_security_policy_pair_policy.testacc_vpn-in-out",
							"policy_b_to_a", "testacc_vpn-in"),
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
					ResourceName:      "junos_security_ipsec_vpn.testacc_ipsecvpn",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityIkeIpsecConfigCreate() string {
	return fmt.Sprintf(`
resource junos_security_ike_proposal "testacc_ikeprop" {
  name = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm = "aes-256-cbc"
  dh_group = "group2"
  lifetime_seconds = 3600
}
resource junos_security_ike_policy "testacc_ikepol" {
  name = "testacc_ikepol"
  proposals = [ junos_security_ike_proposal.testacc_ikeprop.name ]
  mode = "main"
  pre_shared_key_text = "thePassWord"
}
resource junos_security_ike_gateway "testacc_ikegateway" {
  name = "testacc_ikegateway"
  address = [ "192.0.2.3" ]
  policy = junos_security_ike_policy.testacc_ikepol.name
  external_interface = "lo0"
  general_ike_id = true
  no_nat_traversal = true
  dead_peer_detection {
    interval = 10
    threshold = 3
  }
  local_address = "192.0.2.4"
  local_identity {
    type = "hostname"
    value = "testacc"
  }
  remote_identity {
    type = "hostname"
    value = "testacc_remote"
  }
  version = "v2-only"
}

resource junos_security_ipsec_proposal "testacc_ipsecprop" {
  name = "testacc_ipsecprop"
  authentication_algorithm = "hmac-sha1-96"
  protocol = "esp"
  encryption_algorithm = "aes-128-cbc"
}
resource junos_security_ipsec_policy "testacc_ipsecpol" {
  name = "testacc_ipsecpol"
  proposals = [ junos_security_ipsec_proposal.testacc_ipsecprop.name ]
  pfs_keys = "group2"
}
resource junos_security_ipsec_vpn "testacc_ipsecvpn" {
  name = "testacc_ipsecvpn"
  bind_interface_auto = true
  ike {
    gateway = junos_security_ike_gateway.testacc_ikegateway.name
    policy = junos_security_ipsec_policy.testacc_ipsecpol.name
    identity_local = "192.0.2.64/26"
    identity_remote = "192.0.2.128/26"
    identity_service = "any"
  }
  vpn_monitor {
    destination_ip = "192.0.2.129"
    optimized = true
    source_interface_auto = true
  }
  establish_tunnels = "on-traffic"
  df_bit = "clear"
}
`)
}
func testAccJunosSecurityIkeIpsecConfigUpdate() string {
	return fmt.Sprintf(`
resource junos_security_ike_proposal "testacc_ikeprop" {
  name = "testacc_ikeprop"
  authentication_algorithm = "sha1"
  encryption_algorithm = "aes-256-cbc"
  dh_group = "group1"
  lifetime_seconds = 3600
}
resource junos_security_ike_policy "testacc_ikepol" {
  name = "testacc_ikepol"
  proposals = [ junos_security_ike_proposal.testacc_ikeprop.name ]
  mode = "main"
  pre_shared_key_text = "mysecret"
}
resource junos_security_ike_gateway "testacc_ikegateway" {
  name = "testacc_ikegateway"
  address = [ "192.0.2.4" ]
  policy = junos_security_ike_policy.testacc_ikepol.name
  external_interface = "lo0"
  general_ike_id = true
  no_nat_traversal = true
  dead_peer_detection {
    interval = 10
    threshold = 3
  }
  local_address = "192.0.2.4"
  local_identity {
    type = "hostname"
    value = "testacc"
  }
  remote_identity {
    type = "hostname"
    value = "testacc_remote"
  }
  version = "v2-only"
}

resource junos_security_ipsec_proposal "testacc_ipsecprop" {
  name = "testacc_ipsecprop"
  authentication_algorithm = "hmac-sha1-96"
  protocol = "esp"
  encryption_algorithm = "aes-256-cbc"
}
resource junos_security_ipsec_policy "testacc_ipsecpol" {
  name = "testacc_ipsecpol"
  proposals = [ junos_security_ipsec_proposal.testacc_ipsecprop.name ]
  pfs_keys = "group1"
}
resource junos_security_ipsec_vpn "testacc_ipsecvpn" {
  name = "testacc_ipsecvpn"
  bind_interface_auto = false
  ike {
    gateway = junos_security_ike_gateway.testacc_ikegateway.name
    policy = junos_security_ipsec_policy.testacc_ipsecpol.name
    identity_local = "192.0.2.64/26"
    identity_remote = "192.0.2.128/26"
    identity_service = "any"
  }
  establish_tunnels = "immediately"
  df_bit = "clear"
}
resource junos_security_zone testacc_secIkeIpsec_local {
  name = "testacc_secIkeIPsec_local"
  address_book {
    name = "testacc_vpnlocal"
    network = "192.0.2.64/26"
  }
}
resource junos_security_zone testacc_secIkeIpsec_remote {
  name = "testacc_secIkeIPsec_remote"
  address_book {
    name = "testacc_vpnremote"
    network = "192.0.2.128/26"
  }
}
resource junos_security_policy testacc_policyIpsecLocToRem {
  from_zone = junos_security_zone.testacc_secIkeIpsec_local.name
  to_zone = junos_security_zone.testacc_secIkeIpsec_remote.name
  policy {
      name = "testacc_vpn-out"
      match_source_address = [ junos_security_zone.testacc_secIkeIpsec_local.address_book[0].name ]
      match_destination_address = [ junos_security_zone.testacc_secIkeIpsec_remote.address_book[0].name ]
      match_application = [ "any" ]
      permit_tunnel_ipsec_vpn = junos_security_ipsec_vpn.testacc_ipsecvpn.name
  }
}

resource junos_security_policy testacc_policyIpsecRemToLoc {
  from_zone = junos_security_zone.testacc_secIkeIpsec_remote.name
  to_zone = junos_security_zone.testacc_secIkeIpsec_local.name
  policy {
    name = "testacc_vpn-in"
    match_source_address = [ junos_security_zone.testacc_secIkeIpsec_remote.address_book[0].name ]
    match_destination_address = [ junos_security_zone.testacc_secIkeIpsec_local.address_book[0].name ]
    match_application = [ "any" ]
    permit_tunnel_ipsec_vpn = junos_security_ipsec_vpn.testacc_ipsecvpn.name
  }
}

resource junos_security_policy_tunnel_pair_policy testacc_vpn-in-out {
  zone_a = junos_security_zone.testacc_secIkeIpsec_local.name
  zone_b =  junos_security_zone.testacc_secIkeIpsec_remote.name
  policy_a_to_b = junos_security_policy.testacc_policyIpsecLocToRem.policy[0].name
  policy_b_to_a = junos_security_policy.testacc_policyIpsecRemToLoc.policy[0].name
}
`)
}
