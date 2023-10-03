package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSecurityDynamicAddressFeedServer_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSecurityDynamicAddressFeedServerConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_dynamic_address_feed_server.testacc_dyn_add_feed_srv",
							"feed_name.#", "2"),
					),
				},
				{
					Config: testAccResourceSecurityDynamicAddressFeedServerConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_dynamic_address_feed_server.testacc_dyn_add_feed_srv",
							"feed_name.#", "3"),
					),
				},
				{
					ResourceName:      "junos_security_dynamic_address_feed_server.testacc_dyn_add_feed_srv",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccResourceSecurityDynamicAddressFeedServerConfigCreate() string {
	return `
resource "junos_security_dynamic_address_feed_server" "testacc_dyn_add_feed_srv" {
  name     = "tfacc_dafeedsrv"
  hostname = "example.com"
  feed_name {
    name = "feed_b"
    path = "/srx/"
  }
  feed_name {
    name = "feed_a"
    path = "/srx/"
  }
}
resource "junos_security_dynamic_address_feed_server" "testacc_dyn_add_feed_srv2" {
  name = "tfacc_dafeedsr2"
  url  = "https://example.com/?test=#2"
}
`
}

func testAccResourceSecurityDynamicAddressFeedServerConfigUpdate() string {
	return `
resource "junos_security_dynamic_address_feed_server" "testacc_dyn_add_feed_srv" {
  name        = "tfacc_dafeedsrv"
  hostname    = "example.com/?test=#1"
  description = "testacc junos_security_dynamic_address_feed_server"
  feed_name {
    name            = "feed_b"
    path            = "/srx/"
    description     = "testacc junos_security_dynamic_address_feed_server feed_b"
    hold_interval   = 1110
    update_interval = 120
  }
  feed_name {
    name          = "feed_a"
    path          = "/srx/"
    hold_interval = 0
  }
  feed_name {
    name            = "feed_0"
    path            = "/srx/"
    description     = "testacc junos_security_dynamic_address_feed_server feed_0"
    hold_interval   = 1130
    update_interval = 140
  }
  hold_interval   = 1150
  update_interval = 160
}
`
}
