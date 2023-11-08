package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceFirewallPolicer_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic2",
							"filter_specific", "true"),
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
							"if_exceeding.bandwidth_percent", "80"),
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
							"if_exceeding.burst_size_limit", "50k"),
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
							"then.discard", "true"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
							"if_exceeding.bandwidth_limit", "32k"),
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
							"then.forwarding_class", "best-effort"),
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
							"then.loss_priority", "high"),
					),
				},
				{
					ResourceName:      "junos_firewall_policer.testacc_fwPolic",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func TestAccResourceFirewallPolicer_srx(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
