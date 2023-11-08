package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceForwardingOptionsDhcpRelayServerGroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" || os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceForwardingOptionsDhcpRelayServerGroupConfigCreate(),
				},
				{
					Config: testAccResourceForwardingOptionsDhcpRelayServerGroupConfigUpdate(),
				},
				{
					ResourceName:      "junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_servergroup_ri",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_servergroup_v6_ri",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccResourceForwardingOptionsDhcpRelayServerGroupConfigUpdate2(),
				},
				{
					ResourceName:      "junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_servergroup",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_forwardingoptions_dhcprelay_servergroup.testacc_dhcprelay_servergroup_v6",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccResourceForwardingOptionsDhcpRelayServerGroupConfigCreate() string {
	return `
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup" {
  name = "testacc_dhcprelay_servergroup"
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_v6" {
  name    = "testacc_dhcprelay_servergroup_v6"
  version = "v6"
}

resource "junos_routing_instance" "testacc_dhcprelay_servergroup" {
  name = "testacc_dhcprelay_servergroup"
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_ri" {
  name             = "testacc_dhcprelay_servergroup_ri"
  routing_instance = junos_routing_instance.testacc_dhcprelay_servergroup.name
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_v6_ri" {
  name             = "testacc_dhcprelay_servergroup_v6_ri"
  routing_instance = junos_routing_instance.testacc_dhcprelay_servergroup.name
  version          = "v6"
}
`
}

func testAccResourceForwardingOptionsDhcpRelayServerGroupConfigUpdate() string {
	return `
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup" {
  name = "testacc_dhcprelay_servergroup"
  ip_address = [
    "192.0.2.8",
  ]
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_v6" {
  name    = "testacc_dhcprelay_servergroup_v6"
  version = "v6"
  ip_address = [
    "fe80::b",
  ]
}

resource "junos_routing_instance" "testacc_dhcprelay_servergroup" {
  name = "testacc_dhcprelay_servergroup"
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_ri" {
  name             = "testacc_dhcprelay_servergroup_ri"
  routing_instance = junos_routing_instance.testacc_dhcprelay_servergroup.name
  ip_address = [
    "192.0.2.88",
  ]
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_v6_ri" {
  name             = "testacc_dhcprelay_servergroup_v6_ri"
  routing_instance = junos_routing_instance.testacc_dhcprelay_servergroup.name
  version          = "v6"
  ip_address = [
    "fe80::bb",
  ]
}
`
}

func testAccResourceForwardingOptionsDhcpRelayServerGroupConfigUpdate2() string {
	return `
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup" {
  name = "testacc_dhcprelay_servergroup"
  ip_address = [
    "fe80::b",
    "192.0.2.8",
  ]
}
resource "junos_forwardingoptions_dhcprelay_servergroup" "testacc_dhcprelay_servergroup_v6" {
  name    = "testacc_dhcprelay_servergroup_v6"
  version = "v6"
  ip_address = [
    "fe80::b",
    "192.0.2.9",
    "fe80::a",
  ]
}
`
}
