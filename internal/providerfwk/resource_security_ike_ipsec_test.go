package providerfwk_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccResourceSecurityIkeIpsec_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectSensitiveValue("junos_security_ike_policy.testacc_ikepol",
								tfjsonpath.New("pre_shared_key_text")),
						},
					},
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
							"external_interface", testaccInterface+".0"),
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
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
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
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_security_ike_proposal.testacc_ikeprop",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_security_ike_policy.testacc_ikepol",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_security_ike_gateway.testacc_ikegateway",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_security_ipsec_proposal.testacc_ipsecprop",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_security_ipsec_policy.testacc_ipsecpol",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ResourceName:      "junos_security_ipsec_vpn.testacc_ipsecvpn2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectSensitiveValue("junos_security_ike_gateway.testacc_ikegateway",
								tfjsonpath.New("aaa").AtMapKey("client_password")),
						},
					},
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
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
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
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectSensitiveValue("junos_security_ipsec_vpn.testacc_ipsecvpn",
								tfjsonpath.New("manual").AtMapKey("authentication_key_text")),
							plancheck.ExpectSensitiveValue("junos_security_ipsec_vpn.testacc_ipsecvpn",
								tfjsonpath.New("manual").AtMapKey("encryption_key_text")),
						},
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.inet", "192.168.0.4"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectSensitiveValue("junos_security_ipsec_vpn.testacc_ipsecvpn2",
								tfjsonpath.New("manual").AtMapKey("authentication_key_hexa")),
							plancheck.ExpectSensitiveValue("junos_security_ipsec_vpn.testacc_ipsecvpn2",
								tfjsonpath.New("manual").AtMapKey("encryption_key_hexa")),
						},
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.inet6", "2001:db8::1"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
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
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_ike_gateway.testacc_ikegateway",
							"dynamic_remote.distinguished_name.wildcard", "*.com"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					ConfigVariables: map[string]config.Variable{
						"interface": config.StringVariable(testaccInterface),
					},
				},
			},
		})
	}
}
