package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceForwardingoptionsDhcprelayServergroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" || os.Getenv("TESTACC_ROUTER") != "" {
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
					ConfigDirectory: config.TestStepDirectory(),
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
