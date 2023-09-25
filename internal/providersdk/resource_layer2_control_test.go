package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's xe-0/0/3.
// export TESTACC_INTERFACE2=<interface> for choose 2nd interface available else it's xe-0/0/4.
func TestAccResourceLayer2Control_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	testaccInterface2 := junos.DefaultInterfaceTestAcc2
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface = junos.DefaultInterfaceSwitchTestAcc
		testaccInterface2 = junos.DefaultInterfaceSwitchTestAcc2
	}
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
		testaccInterface2 = iface
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceLayer2ControlConfigCreate(),
			},
			{
				Config: testAccResourceLayer2ControlConfigUpdate(testaccInterface),
			},
			{
				ResourceName:      "junos_layer2_control.l2c",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceLayer2ControlConfigUpdate2(testaccInterface, testaccInterface2),
			},
		},
	})
}

func testAccResourceLayer2ControlConfigCreate() string {
	return `
resource "junos_layer2_control" "l2c" {
  bpdu_block {}
}
`
}

func testAccResourceLayer2ControlConfigUpdate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_chassis_redundancy" "l2c" {
  graceful_switchover = true
}
resource "junos_layer2_control" "l2c" {
  depends_on = [
    junos_chassis_redundancy.l2c
  ]
  lifecycle {
    create_before_destroy = true
  }
  bpdu_block {
    disable_timeout = 300
    interface {
      name    = "%s"
      disable = true
      drop    = true
    }
  }
  mac_rewrite_interface {
    name           = "%s"
    enable_all_ifl = true
    protocol       = ["cdp", "stp"]
  }
  nonstop_bridging = true
}
`, interFace, interFace)
}

func testAccResourceLayer2ControlConfigUpdate2(interFace, interFace2 string) string {
	return fmt.Sprintf(`
resource "junos_layer2_control" "l2c" {
  bpdu_block {
    interface {
      name = "%s"
    }
    interface {
      name    = "%s"
      disable = true
      drop    = true
    }
  }
  mac_rewrite_interface {
    name           = "%s"
    enable_all_ifl = true
    protocol       = ["cdp", "stp"]
  }
  mac_rewrite_interface {
    name = "%s"
  }
}
`, interFace2, interFace, interFace2, interFace)
}
