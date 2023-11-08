package providersdk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSystemNtpServer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSystemNtpServerConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_ntp_server.testacc_ntpServer",
						"address", "192.0.2.1"),
					resource.TestCheckResourceAttr("junos_system_ntp_server.testacc_ntpServer",
						"prefer", "true"),
					resource.TestCheckResourceAttr("junos_system_ntp_server.testacc_ntpServer",
						"version", "4"),
					resource.TestCheckResourceAttr("junos_system_ntp_server.testacc_ntpServer",
						"key", "1"),
				),
			},
			{
				Config: testAccResourceSystemNtpServerConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_ntp_server.testacc_ntpServer",
						"routing_instance", "testacc_ntpServer"),
				),
			},
			{
				ResourceName:      "junos_system_ntp_server.testacc_ntpServer",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceSystemNtpServerConfigCreate() string {
	return `
resource "junos_system_ntp_server" "testacc_ntpServer" {
  address = "192.0.2.1"
  prefer  = true
  version = 4
  key     = 1
}
`
}

func testAccResourceSystemNtpServerConfigUpdate() string {
	return `
resource "junos_routing_instance" "testacc_ntpServer" {
  name = "testacc_ntpServer"
}
resource "junos_system_ntp_server" "testacc_ntpServer" {
  address          = "192.0.2.1"
  prefer           = true
  version          = 4
  routing_instance = junos_routing_instance.testacc_ntpServer.name
}
`
}
