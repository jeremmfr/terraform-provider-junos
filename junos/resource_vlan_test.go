package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosVlan_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosVlanSwConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"description", "testacc_vlansw"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"vlan_id", "1000"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"service_id", "1000"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"l3_interface", "irb.1000"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"forward_filter_input", "testacc_vlansw"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"forward_filter_output", "testacc_vlansw"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"forward_flood_input", "testacc_vlansw"),
					),
				},
				{
					Config: testAccJunosVlanSwConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"vlan_id", "0"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"vlan_id_list.#", "1"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"vlan_id_list.0", "1001-1002"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"private_vlan", "community"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"l3_interface", ""),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"forward_filter_input", ""),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"forward_filter_output", ""),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"forward_flood_input", ""),
					),
				},
				{
					ResourceName:      "junos_vlan.testacc_vlansw",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosVlanSwConfigCreate() string {
	return `
resource "junos_firewall_filter" "testacc_vlansw" {
  lifecycle {
    create_before_destroy = true
  }
  name   = "testacc_vlansw"
  family = "ethernet-switching"
  term {
    name = "testacc_vlansw_term1"
    then {
      action = "accept"
    }
  }
}
resource "junos_interface_logical" "testacc_vlansw" {
  lifecycle {
    create_before_destroy = true
  }
  name = "irb.1000"
}
resource "junos_vlan" "testacc_vlansw" {
  name                  = "testacc_vlansw"
  description           = "testacc_vlansw"
  vlan_id               = 1000
  service_id            = 1000
  l3_interface          = junos_interface_logical.testacc_vlansw.name
  forward_filter_input  = junos_firewall_filter.testacc_vlansw.name
  forward_filter_output = junos_firewall_filter.testacc_vlansw.name
  forward_flood_input   = junos_firewall_filter.testacc_vlansw.name
}
`
}

func testAccJunosVlanSwConfigUpdate() string {
	return `
resource "junos_vlan" "testacc_vlansw" {
  name         = "testacc_vlansw"
  description  = "testacc_vlansw"
  vlan_id_list = ["1001-1002"]
  private_vlan = "community"
}
`
}
