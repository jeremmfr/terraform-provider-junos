package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3 or xe-0/0/3.
func TestAccDataSourceInterfaceLogicalInfo_basic(t *testing.T) {
	testaccInterface := defaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceInterfaceLogicalInfoPre(testaccInterface),
				},
				{
					Config: testAccDataSourceInterfaceLogicalInfoConfig(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_interface_logical_info.testacc_dataIfaceLogInfo",
							"admin_status", "up"),
						resource.TestCheckResourceAttr("data.junos_interface_logical_info.testacc_dataIfaceLogInfo",
							"family_inet.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interface_logical_info.testacc_dataIfaceLogInfo",
							"family_inet.0.address_cidr.#", "2"),
						resource.TestCheckResourceAttr("data.junos_interface_logical_info.testacc_dataIfaceLogInfo",
							"family_inet6.#", "1"),
						resource.TestCheckTypeSetElemAttr("data.junos_interface_logical_info.testacc_dataIfaceLogInfo",
							"family_inet6.0.address_cidr.*", "2001:db8::1/64"),
					),
				},
			},
		})
	}
}

func testAccDataSourceInterfaceLogicalInfoPre(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_dataIfaceLogInfo" {
  name         = "%s"
  description  = "testacc_dataIfaceLogInfo"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_dataIfaceLogInfo" {
  name        = "${junos_interface_physical.testacc_dataIfaceLogInfo.name}.10"
  description = "testacc_dataIfaceLogInfo"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
    }
    address {
      cidr_ip = "192.0.2.2/25"
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

func testAccDataSourceInterfaceLogicalInfoConfig(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_dataIfaceLogInfo" {
  name         = "%s"
  description  = "testacc_dataIfaceLogInfo"
  vlan_tagging = true
}
resource "junos_interface_logical" "testacc_dataIfaceLogInfo" {
  name        = "${junos_interface_physical.testacc_dataIfaceLogInfo.name}.10"
  description = "testacc_dataIfaceLogInfo"
  family_inet {
    address {
      cidr_ip = "192.0.2.1/25"
    }
    address {
      cidr_ip = "192.0.2.2/25"
    }
  }
  family_inet6 {
    address {
      cidr_ip = "2001:db8::1/64"
    }
  }
}
data "junos_interface_logical_info" "testacc_dataIfaceLogInfo" {
  name = junos_interface_logical.testacc_dataIfaceLogInfo.name
}
`, interFace)
}
