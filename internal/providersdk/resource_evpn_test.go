package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosEvpn_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosEvpnConfigCreate(),
				},
				{
					ResourceName:      "junos_evpn.testacc_evpn_default",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_evpn.testacc_evpn_ri1",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosEvpnConfigUpdate(),
				},
			},
		})
	}
}

func testAccJunosEvpnConfigCreate() string {
	return `
resource "junos_interface_logical" "testacc_evpn" {
  depends_on = [
    junos_routing_options.testacc_evpn,
  ]
  name        = "lo0.0"
  description = "testacc_evpn"
  family_inet {
    address {
      cidr_ip = "192.0.2.18/32"
    }
  }
}
resource "junos_routing_options" "testacc_evpn" {
  clean_on_destroy = true
  router_id        = "192.0.2.18"
}
resource "junos_switch_options" "testacc_evpn" {
  clean_on_destroy      = true
  vtep_source_interface = junos_interface_logical.testacc_evpn.name
}
resource "junos_policyoptions_community" "testacc_evpn" {
  lifecycle {
    create_before_destroy = true
  }
  name    = "testacc_evpn"
  members = ["target:65000:100"]
}
resource "junos_policyoptions_policy_statement" "testacc_evpn" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_evpn1"
  from {
    bgp_community = [junos_policyoptions_community.testacc_evpn.name]
  }
  then {
    action = "accept"
  }
}
resource "junos_evpn" "testacc_evpn_default" {
  depends_on = [
    junos_switch_options.testacc_evpn,
  ]
  encapsulation = "vxlan"
  switch_or_ri_options {
    route_distinguisher = "20:1"
    vrf_target          = "target:20:2"
    vrf_import          = [junos_policyoptions_policy_statement.testacc_evpn.name]
    vrf_export          = [junos_policyoptions_policy_statement.testacc_evpn.name]
  }
}
resource "junos_routing_instance" "testacc_evpn_ri1" {
  name                  = "testacc_evpn_ri1"
  type                  = "virtual-switch"
  route_distinguisher   = "1:1"
  vrf_target            = "target:1:2"
  vrf_import            = [junos_policyoptions_policy_statement.testacc_evpn.name]
  vrf_export            = [junos_policyoptions_policy_statement.testacc_evpn.name]
  vtep_source_interface = junos_interface_logical.testacc_evpn.name
}
resource "junos_evpn" "testacc_evpn_ri1" {
  routing_instance = junos_routing_instance.testacc_evpn_ri1.name
  encapsulation    = "vxlan"
}
resource "junos_routing_instance" "testacc_evpn_ri2" {
  name                        = "testacc_evpn_ri2"
  type                        = "virtual-switch"
  configure_rd_vrfopts_singly = true
  vtep_source_interface       = junos_interface_logical.testacc_evpn.name
}
resource "junos_evpn" "testacc_evpn_ri2" {
  routing_instance = junos_routing_instance.testacc_evpn_ri2.name
  encapsulation    = "vxlan"
  default_gateway  = "advertise"
  switch_or_ri_options {
    route_distinguisher = "10:1"
    vrf_import          = [junos_policyoptions_policy_statement.testacc_evpn.name]
    vrf_export          = [junos_policyoptions_policy_statement.testacc_evpn.name]
    vrf_target          = "target:10:2"
  }
}
`
}

func testAccJunosEvpnConfigUpdate() string {
	return `
resource "junos_interface_logical" "testacc_evpn" {
  depends_on = [
    junos_routing_options.testacc_evpn,
  ]
  name        = "lo0.0"
  description = "testacc_evpn"
  family_inet {
    address {
      cidr_ip = "192.0.2.18/32"
    }
  }
}
resource "junos_routing_options" "testacc_evpn" {
  clean_on_destroy = true
  router_id        = "192.0.2.18"
}
resource "junos_switch_options" "testacc_evpn" {
  clean_on_destroy      = true
  vtep_source_interface = junos_interface_logical.testacc_evpn.name
}
resource "junos_evpn" "testacc_evpn_default" {
  depends_on = [
    junos_switch_options.testacc_evpn,
  ]
  encapsulation = "vxlan"
  switch_or_ri_options {
    route_distinguisher = "201:1"
    vrf_target          = "target:201:2"
  }
}
resource "junos_routing_instance" "testacc_evpn_ri1" {
  name                  = "testacc_evpn_ri1"
  type                  = "virtual-switch"
  route_distinguisher   = "11:1"
  vrf_target            = "target:11:2"
  vtep_source_interface = junos_interface_logical.testacc_evpn.name
}
resource "junos_evpn" "testacc_evpn_ri1" {
  routing_instance = junos_routing_instance.testacc_evpn_ri1.name
  encapsulation    = "vxlan"
}
resource "junos_routing_instance" "testacc_evpn_ri2" {
  name                        = "testacc_evpn_ri2"
  type                        = "virtual-switch"
  configure_rd_vrfopts_singly = true
  vtep_source_interface       = junos_interface_logical.testacc_evpn.name
}
resource "junos_evpn" "testacc_evpn_ri2" {
  routing_instance = junos_routing_instance.testacc_evpn_ri2.name
  encapsulation    = "vxlan"
  switch_or_ri_options {
    route_distinguisher = "101:1"
    vrf_target          = "target:101:2"
  }
}
`
}
