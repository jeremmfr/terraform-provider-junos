package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSwitchOptions_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSwitchOptionsConfigCreate(),
				},
				{
					ResourceName:      "junos_switch_options.testacc_switchOpts",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosSwitchOptionsConfigUpdate(),
				},
			},
		})
	}
}

func testAccJunosSwitchOptionsConfigCreate() string {
	return `
resource junos_interface_logical "testacc_switchOpts" {
  lifecycle {
    create_before_destroy = true
  }
  name = "lo0.0"
  family_inet {
    address {
      cidr_ip = "192.0.2.16/32"
    }
  }
}
resource "junos_switch_options" "testacc_switchOpts" {
  vtep_source_interface = junos_interface_logical.testacc_switchOpts.name
}
`
}

func testAccJunosSwitchOptionsConfigUpdate() string {
	return `
resource "junos_switch_options" "testacc_switchOpts" {
}
`
}
