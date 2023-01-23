package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRoutes_basic(t *testing.T) {
	testaccInterface := defaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SRX") != "" || os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceRoutesPre(testaccInterface),
				},
				{
					Config: testAccDataSourceRoutesConfig(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckTypeSetElemNestedAttrs("data.junos_routes.all",
							"table.*", map[string]string{"name": "inet.0"}),
						resource.TestCheckTypeSetElemNestedAttrs("data.junos_routes.all",
							"table.*", map[string]string{"name": "testacc_data_routes.inet.0"}),
						resource.TestCheckTypeSetElemNestedAttrs("data.junos_routes.all",
							"table.*.route.*", map[string]string{"destination": "0.0.0.0/0"}),
						resource.TestCheckTypeSetElemNestedAttrs("data.junos_routes.all",
							"table.*.route.*.entry.*", map[string]string{
								"current_active": "true",
								"protocol":       "Local",
								"next_hop_type":  "Local",
							}),
						resource.TestCheckResourceAttr("data.junos_routes.default",
							"table.#", "1"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.#", "1"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.0.route.#", "1"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.0.route.0.destination", "192.0.2.1/32"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.0.route.0.entry.#", "1"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.0.route.0.entry.0.current_active", "true"),
						resource.TestCheckResourceAttr("data.junos_routes.testacc",
							"table.0.route.0.entry.0.protocol", "Local"),
					),
				},
			},
		})
	}
}

func testAccDataSourceRoutesPre(interFace string) string {
	return fmt.Sprintf(`
resource "junos_routing_instance" "testacc_data_routes" {
  name = "testacc_data_routes"
}
resource "junos_interface_physical" "testacc_data_routes" {
  name         = "%s"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  name             = "${junos_interface_physical.testacc_data_routes.name}.100"
  routing_instance = junos_routing_instance.testacc_data_routes.name
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
    }
  }
}
`, interFace)
}

func testAccDataSourceRoutesConfig(interFace string) string {
	return fmt.Sprintf(`
resource "junos_routing_instance" "testacc_data_routes" {
  name = "testacc_data_routes"
}
resource "junos_interface_physical" "testacc_data_routes" {
  name         = "%s"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_interface_logical" {
  name             = "${junos_interface_physical.testacc_data_routes.name}.100"
  routing_instance = junos_routing_instance.testacc_data_routes.name
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
    }
  }
}

data "junos_routes" "all" {}
data "junos_routes" "default" {
  table_name = "inet.0"
}
data "junos_routes" "testacc" {
  table_name = "${junos_routing_instance.testacc_data_routes.name}.inet.0"
}
`, interFace)
}
