package providersdk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceChassisRedundancy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceChassisRedundancyConfigCreate(),
			},
			{
				Config: testAccResourceChassisRedundancyConfigUpdate(),
			},
			{
				ResourceName:      "junos_chassis_redundancy.testacc_cred",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceChassisRedundancyConfigCreate() string {
	return `
resource "junos_chassis_redundancy" "testacc_cred" {
  graceful_switchover = true
}
`
}

func testAccResourceChassisRedundancyConfigUpdate() string {
	return `
resource "junos_chassis_redundancy" "testacc_cred" {
  failover_disk_read_threshold      = 2000
  failover_disk_write_threshold     = 3000
  failover_not_on_disk_underperform = true
  failover_on_disk_failure          = true
  failover_on_loss_of_keepalives    = true
  keepalive_time                    = 300
  routing_engine {
    slot = 0
    role = "master"
  }
}
`
}
