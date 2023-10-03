package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRstp_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceRstpSWConfigCreate(),
				},
				{
					Config: testAccResourceRstpSWConfigUpdate(),
				},
				{
					ResourceName:      "junos_rstp.testacc_ri_rstp",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	} else {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceRstpConfigCreate(),
				},
				{
					Config: testAccResourceRstpConfigUpdate(),
				},
				{
					Config: testAccResourceRstpConfigUpdate2(),
				},
				{
					ResourceName:      "junos_rstp.testacc_rstp",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccResourceRstpSWConfigCreate() string {
	return `
resource "junos_rstp" "testacc_rstp" {
  bpdu_block_on_edge = true
}
resource "junos_routing_instance" "testacc_rstp" {
  name = "testacc_rstp"
  type = "virtual-switch"
}
resource "junos_rstp" "testacc_ri_rstp" {
  routing_instance = junos_routing_instance.testacc_rstp.name
  bridge_priority  = 0
  system_id {
    id = "00:11:22:33:44:56"
  }
}
`
}

func testAccResourceRstpSWConfigUpdate() string {
	return `
resource "junos_rstp" "testacc_rstp" {
  bpdu_block_on_edge     = true
  backup_bridge_priority = "8k"
  bridge_priority        = 0
}
resource "junos_routing_instance" "testacc_rstp" {
  name = "testacc_rstp"
  type = "virtual-switch"
}
resource "junos_rstp" "testacc_ri_rstp" {
  routing_instance                                   = junos_routing_instance.testacc_rstp.name
  backup_bridge_priority                             = "60k"
  bridge_priority                                    = "4k"
  bpdu_destination_mac_address_provider_bridge_group = true
  extended_system_id                                 = 0
  force_version_stp                                  = true
  forward_delay                                      = 20
  hello_time                                         = 5
  max_age                                            = 22
  priority_hold_time                                 = 100
  system_id {
    id = "00:11:22:33:44:55"
  }
  system_id {
    id         = "00:22:33:44:55:aa"
    ip_address = "192.0.2.4/31"
  }
  system_identifier             = "66:55:44:33:22:11"
  vpls_flush_on_topology_change = true
}
`
}

func testAccResourceRstpConfigCreate() string {
	return `
resource "junos_rstp" "testacc_rstp" {
  disable = true
}
`
}

func testAccResourceRstpConfigUpdate() string {
	return `
resource "junos_rstp" "testacc_rstp" {
  backup_bridge_priority = "32k"
  bridge_priority        = "16k"
  system_id {
    id = "00:22:33:44:55:aa"
  }
}
`
}

func testAccResourceRstpConfigUpdate2() string {
	return `
resource "junos_rstp" "testacc_rstp" {
  backup_bridge_priority                             = "60k"
  bridge_priority                                    = "4k"
  bpdu_destination_mac_address_provider_bridge_group = true
  extended_system_id                                 = 0
  force_version_stp                                  = true
  forward_delay                                      = 20
  hello_time                                         = 5
  max_age                                            = 22
  priority_hold_time                                 = 100
  system_id {
    id = "00:11:22:33:44:55"
  }
  system_id {
    id         = "00:22:33:44:55:66"
    ip_address = "192.0.2.4/24"
  }
  system_identifier             = "66:55:44:33:22:11"
  vpls_flush_on_topology_change = true
}
`
}
