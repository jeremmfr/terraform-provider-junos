package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceRoutingInstance_router(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"type", "virtual-router"),
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"as", "65000"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"as", "65001"),
					),
				},
				{
					ResourceName:      "junos_routing_instance.testacc_routingInst",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_routing_instance.testacc_routingInst2",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func TestAccResourceRoutingInstance_srx(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"type", "virtual-router"),
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"as", "65000"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_routing_instance.testacc_routingInst",
							"as", "65001"),
					),
				},
				{
					ResourceName:      "junos_routing_instance.testacc_routingInst",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
