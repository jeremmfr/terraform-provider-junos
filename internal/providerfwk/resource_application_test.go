package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceApplication_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_application.testacc_app", "protocol", "tcp"),
						resource.TestCheckResourceAttr("junos_application.testacc_app", "destination_port", "22"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_application.testacc_app", "protocol", "tcp"),
						resource.TestCheckResourceAttr("junos_application.testacc_app", "destination_port", "22"),
						resource.TestCheckResourceAttr("junos_application.testacc_app", "source_port", "1024-65535"),
					),
				},
				{
					ResourceName:      "junos_application.testacc_app",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_application.testacc_app2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_application.testacc_app3",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_application.testacc_app4",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_application.testacc_app5",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
