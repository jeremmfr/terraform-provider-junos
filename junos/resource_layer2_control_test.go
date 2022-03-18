package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's xe-0/0/3.
// export TESTACC_INTERFACE2=<interface> for choose 2nd interface available else it's xe-0/0/4.
func TestAccJunosLayer2Control_basic(t *testing.T) {
	testaccInterface := defaultInterfaceTestAcc
	testaccInterface2 := defaultInterfaceTestAcc2
	if os.Getenv("TESTACC_SWITCH") != "" {
		testaccInterface = defaultInterfaceSwitchTestAcc
		testaccInterface2 = defaultInterfaceSwitchTestAcc2
	}
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if iface := os.Getenv("TESTACC_INTERFACE2"); iface != "" {
		testaccInterface2 = iface
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosLayer2ControlConfigCreate(),
			},
			{
				Config: testAccJunosLayer2ControlConfigUpdate(testaccInterface),
			},
			{
				ResourceName:      "junos_layer2_control.l2c",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccJunosLayer2ControlConfigUpdate2(testaccInterface, testaccInterface2),
			},
		},
	})
}

func testAccJunosLayer2ControlConfigCreate() string {
	return `
resource "junos_layer2_control" "l2c" {
  bpdu_block {}
}
`
}

func testAccJunosLayer2ControlConfigUpdate(interFace string) string {
	return `
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
      name    = "` + interFace + `"
      disable = true
      drop    = true
    }
  }
  mac_rewrite_interface {
    name           = "` + interFace + `"
    enable_all_ifl = true
    protocol       = ["cdp", "stp"]
  }
  nonstop_bridging = true
}
`
}

func testAccJunosLayer2ControlConfigUpdate2(interFace, interFace2 string) string {
	return `
resource "junos_layer2_control" "l2c" {
  bpdu_block {
    interface {
      name = "` + interFace2 + `"
    }
    interface {
      name    = "` + interFace + `"
      disable = true
      drop    = true
    }
  }
  mac_rewrite_interface {
    name           = "` + interFace2 + `"
    enable_all_ifl = true
    protocol       = ["cdp", "stp"]
  }
  mac_rewrite_interface {
    name = "` + interFace + `"
  }
}
`
}
