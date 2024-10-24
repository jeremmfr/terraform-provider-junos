package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceForwardingoptionsDhcprelayGroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ResourceName:      "junos_forwardingoptions_dhcprelay_group.testacc_dhcprelaygroup_v4_default",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_forwardingoptions_dhcprelay_group.testacc_dhcprelaygroup_v6_default",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_forwardingoptions_dhcprelay_group.testacc_dhcprelaygroup_v4_ri",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_forwardingoptions_dhcprelay_group.testacc_dhcprelaygroup_v6_ri",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
