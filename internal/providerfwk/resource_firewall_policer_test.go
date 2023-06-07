package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosFirewallPolicer_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosFirewallPolicerConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
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
					Config: testAccJunosFirewallPolicerConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
							"if_exceeding.bandwidth_limit", "32k"),
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
							"then.forwarding_class", "best-effort"),
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
							"then.loss_priority", "high"),
						resource.TestCheckResourceAttr("junos_firewall_policer.testacc_fwPolic",
							"then.out_of_profile", "true"),
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

func testAccJunosFirewallPolicerConfigCreate() string {
	return `
resource "junos_firewall_policer" "testacc_fwPolic" {
  name            = "testacc_fwPolic"
  filter_specific = true
  if_exceeding {
    bandwidth_percent = 80
    burst_size_limit  = "50k"
  }
  then {
    discard = true
  }
}
`
}

func testAccJunosFirewallPolicerConfigUpdate() string {
	return `
resource "junos_firewall_policer" "testacc_fwPolic" {
  name = "testacc_fwPolic"
  if_exceeding {
    bandwidth_limit  = "32k"
    burst_size_limit = "50k"
  }
  then {
    forwarding_class = "best-effort"
    loss_priority    = "high"
    out_of_profile   = true
  }
}
`
}
