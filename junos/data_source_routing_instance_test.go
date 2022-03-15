package junos_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccDataSourceRoutingInstance_basic(t *testing.T) {
	testaccInterface := defaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceRoutingInstanceConfigCreate(testaccInterface),
				},
				{
					Config: testAccDataSourceRoutingInstanceConfigData(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_routing_instance.testacc_dataRoutingInstance",
							"id", "testacc_dataRoutingInstance"),
						resource.TestCheckResourceAttr("data.junos_routing_instance.testacc_dataRoutingInstance",
							"type", "virtual-router"),
						resource.TestCheckResourceAttr("data.junos_routing_instance.testacc_dataRoutingInstance",
							"interface.#", "1"),
					),
				},
				{
					Config:      testAccDataSourceRoutingInstanceConfigDataFailed(),
					ExpectError: regexp.MustCompile("routing instance .* doesn't exist"),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}

func testAccDataSourceRoutingInstanceConfigCreate(interFace string) string {
	return `
resource junos_interface_physical testacc_dataRoutingInstance {
  name         = "` + interFace + `"
  description  = "testacc_dataRoutingInstance"
  vlan_tagging = true
}
resource junos_routing_instance testacc_dataRoutingInstance {
  name = "testacc_dataRoutingInstance"
}
resource junos_interface_logical testacc_dataRoutingInstance {
  name             = "${junos_interface_physical.testacc_dataRoutingInstance.name}.100"
  description      = "testacc_dataRoutingInstance"
  routing_instance = junos_routing_instance.testacc_dataRoutingInstance.name
}
`
}

func testAccDataSourceRoutingInstanceConfigData(interFace string) string {
	return `
resource junos_interface_physical testacc_dataRoutingInstance {
  name         = "` + interFace + `"
  description  = "testacc_dataRoutingInstance"
  vlan_tagging = true
}
resource junos_routing_instance testacc_dataRoutingInstance {
  name = "testacc_dataRoutingInstance"
}
resource junos_interface_logical testacc_dataRoutingInstance {
  name             = "${junos_interface_physical.testacc_dataRoutingInstance.name}.100"
  description      = "testacc_dataRoutingInstance"
  routing_instance = junos_routing_instance.testacc_dataRoutingInstance.name
}

data junos_routing_instance testacc_dataRoutingInstance {
  name = "testacc_dataRoutingInstance"
}
`
}

func testAccDataSourceRoutingInstanceConfigDataFailed() string {
	return `
data junos_routing_instance testacc_dataRoutingInstance {
  name = "testacc"
}
`
}
