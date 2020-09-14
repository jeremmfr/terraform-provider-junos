package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccDataSourceInterface_basic(t *testing.T) {
	var testaccInterface string
	if os.Getenv("TESTACC_INTERFACE") != "" {
		testaccInterface = os.Getenv("TESTACC_INTERFACE")
	} else {
		testaccInterface = defaultInterfaceTestAcc
	}
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceInterfaceConfigCreate(testaccInterface),
				},
				{
					Config: testAccDataSourceInterfaceConfigData(testaccInterface),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_interface.testacc_datainterface",
							"id", testaccInterface+".100"),
						resource.TestCheckResourceAttr("data.junos_interface.testacc_datainterface",
							"name", testaccInterface+".100"),
						resource.TestCheckResourceAttr("data.junos_interface.testacc_datainterface",
							"inet_address.#", "1"),
						resource.TestCheckResourceAttr("data.junos_interface.testacc_datainterface",
							"inet_address.0.address", "192.0.2.1/25"),
						resource.TestCheckResourceAttr("data.junos_interface.testacc_datainterface2",
							"id", testaccInterface+".100"),
					),
				},
			},
			PreventPostDestroyRefresh: true,
		})
	}
}

func testAccDataSourceInterfaceConfigCreate(interFace string) string {
	return `
resource junos_interface testacc_datainterfaceP {
  name         = "` + interFace + `"
  description  = "testacc_datainterfaceP"
  vlan_tagging = true
}
resource junos_interface testacc_datainterface {
  name        = "${junos_interface.testacc_datainterfaceP.name}.100"
  description = "testacc_datainterface"
  inet_address {
    address = "192.0.2.1/25"
  }
}
`
}

func testAccDataSourceInterfaceConfigData(interFace string) string {
	return `
resource junos_interface testacc_datainterfaceP {
  name         = "` + interFace + `"
  description  = "testacc_datainterfaceP"
  vlan_tagging = true
}
resource junos_interface testacc_datainterface {
  name        = "${junos_interface.testacc_datainterfaceP.name}.100"
  description = "testacc_datainterface"
  inet_address {
    address = "192.0.2.1/25"
  }
}

data junos_interface testacc_datainterface {
  config_interface = "` + interFace + `"
  match            = "192.0.2.1/"
}

data junos_interface testacc_datainterface2 {
  match      = "192.0.2.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
}
`
}
