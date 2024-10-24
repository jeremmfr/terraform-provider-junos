package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceForwardingoptionsDhcprelay_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					// 1
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					// 2
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					// 3
					ResourceName:      "junos_forwardingoptions_dhcprelay.testacc_dhcprelay_v4_default",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					// 4
					ResourceName:      "junos_forwardingoptions_dhcprelay.testacc_dhcprelay_v6_default",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					// 5
					ResourceName:      "junos_forwardingoptions_dhcprelay.testacc_dhcprelay_v4_ri",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					// 6
					ResourceName:      "junos_forwardingoptions_dhcprelay.testacc_dhcprelay_v6_ri",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					// 7
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
