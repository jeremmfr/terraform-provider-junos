package providerfwk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccDataSourceInterfaceLogical_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceInterfaceLogicalConfigCreate(testaccInterface),
				},
				{
					Config: testAccDataSourceInterfaceLogicalConfigData(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_interface_logical.testacc_datainterfaceL",
							"id", testaccInterface+".100"),
						resource.TestCheckResourceAttr("data.junos_interface_logical.testacc_datainterfaceL",
							"name", testaccInterface+".100"),
						resource.TestCheckResourceAttr("data.junos_interface_logical.testacc_datainterfaceL",
							"family_inet.address.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interface_logical.testacc_datainterfaceL",
							"family_inet.address.0.cidr_ip", "192.0.2.1/25"),
						resource.TestCheckResourceAttr("data.junos_interface_logical.testacc_datainterfaceL2",
							"id", testaccInterface+".100"),
					),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}

func testAccDataSourceInterfaceLogicalConfigCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_datainterfaceP" {
  name         = "%s"
  description  = "testacc_datainterfaceP"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_datainterfaceL" {
  name        = "${junos_interface_physical.testacc_datainterfaceP.name}.100"
  description = "testacc_datainterfaceL"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
    }
  }
  family_inet6 {
    address {
      cidr_ip = "2001:db8::1/64"
    }
  }
}
`, interFace)
}

func testAccDataSourceInterfaceLogicalConfigData(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_datainterfaceP" {
  name         = "%s"
  description  = "testacc_datainterfaceP"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_datainterfaceL" {
  name        = "${junos_interface_physical.testacc_datainterfaceP.name}.100"
  description = "testacc_datainterfaceL"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
    }
  }
}

data "junos_interface_logical" "testacc_datainterfaceL" {
  config_interface = "%s"
  match            = "192.0.2.1/"
}

data "junos_interface_logical" "testacc_datainterfaceL2" {
  match = "192.0.2.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
}
`, interFace, interFace)
}
