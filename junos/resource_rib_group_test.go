package junos

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccJunosRibGroup_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosRibGroupConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_policy.#", "1"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_policy.0", "testacc_ribGroup"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_rib.#", "1"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_rib.0", "testacc_ribGroup1.inet.0"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"export_rib", "testacc_ribGroup1.inet.0"),
					),
				},
				{
					Config: testAccJunosRibGroupConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_rib.#", "2"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"import_rib.1", "testacc_ribGroup2.inet.0"),
						resource.TestCheckResourceAttr("junos_rib_group.testacc_ribGroup",
							"export_rib", "testacc_ribGroup2.inet.0"),
					),
				},
				{
					ResourceName:      "junos_rib_group.testacc_ribGroup",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosRibGroupConfigCreate() string {
	return fmt.Sprintf(`
resource junos_routing_instance "testacc_ribGroup1" {
  name = "testacc_ribGroup1"
}
resource junos_policyoptions_policy_statement "testacc_ribGroup" {
  name = "testacc_ribGroup"
  then {
    action = "accept"
  }
}
resource junos_rib_group testacc_ribGroup {
  name = "testacc_ribGroup-test"
  import_policy = [ junos_policyoptions_policy_statement.testacc_ribGroup.name, ]
  import_rib = [
    "${junos_routing_instance.testacc_ribGroup1.name}.inet.0",
  ]
  export_rib = "${junos_routing_instance.testacc_ribGroup1.name}.inet.0"
}
`)
}
func testAccJunosRibGroupConfigUpdate() string {
	return fmt.Sprintf(`
resource junos_routing_instance "testacc_ribGroup1" {
  name = "testacc_ribGroup1"
}
resource junos_routing_instance "testacc_ribGroup2" {
  name = "testacc_ribGroup2"
}
resource junos_policyoptions_policy_statement "testacc_ribGroup" {
  name = "testacc_ribGroup"
  then {
    action = "accept"
  }
}
resource junos_rib_group testacc_ribGroup {
  name = "testacc_ribGroup-test"
  import_policy = [ junos_policyoptions_policy_statement.testacc_ribGroup.name, ]
  import_rib = [
    "${junos_routing_instance.testacc_ribGroup1.name}.inet.0",
    "${junos_routing_instance.testacc_ribGroup2.name}.inet.0",
  ]
  export_rib = "${junos_routing_instance.testacc_ribGroup2.name}.inet.0"
}
`)
}
