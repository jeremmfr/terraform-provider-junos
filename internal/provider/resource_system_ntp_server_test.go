package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSystemNtpServer_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
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
					ConfigDirectory: config.TestStepDirectory(),
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
}
