package junos

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccJunosOspfArea_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosOspfAreaConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"area_id", "0.0.0.0"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"version", "v2"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"routing_instance", "default"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.#", "1"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.name", "all"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.disable", "true"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.passive", "true"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.metric", "100"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.retransmit_interval", "10"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.hello_interval", "10"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.0.dead_interval", "10"),
					),
				},
				{
					Config: testAccJunosOspfAreaConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"routing_instance", "testacc_ospfarea"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"version", "v3"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.#", "2"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.1.name", "vlan.100"),
						resource.TestCheckResourceAttr("junos_ospf_area.testacc_ospfarea",
							"interface.1.disable", "true"),
					),
				},
				{
					ResourceName:      "junos_ospf_area.testacc_ospfarea",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosOspfAreaConfigCreate() string {
	return `
resource junos_ospf_area "testacc_ospfarea" {
  area_id = "0.0.0.0"
  interface {
    name = "all"
    disable = true
    passive = true
    metric = 100
    retransmit_interval = 10
    hello_interval = 10
    dead_interval = 10
  }
}
`
}
func testAccJunosOspfAreaConfigUpdate() string {
	return `
resource junos_interface "testacc_ospfarea" {
  name = "vlan.100"
  routing_instance = junos_routing_instance.testacc_ospfarea.name
}
resource junos_routing_instance "testacc_ospfarea" {
  name = "testacc_ospfarea"
}
resource junos_ospf_area "testacc_ospfarea" {
  area_id = "0.0.0.0"
  version = "v3"
  routing_instance = junos_routing_instance.testacc_ospfarea.name
  interface {
    name = "all"
    passive = true
    metric = 100
    retransmit_interval = 10
    hello_interval = 10
    dead_interval = 10
  }
  interface {
    name = junos_interface.testacc_ospfarea.name
    disable = true
  }
}
`
}
