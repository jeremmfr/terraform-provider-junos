package providerfwk_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJunosInterfaceSt0Unit_basic(t *testing.T) {
	regexpSt0 := regexp.MustCompile("st0.")
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosInterfaceSt0UnitConfig(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestMatchResourceAttr("junos_interface_st0_unit.testacc", "id", regexpSt0),
					),
				},
			},
		})
	}
}

func testAccJunosInterfaceSt0UnitConfig() string {
	return `
resource "junos_interface_st0_unit" "testacc" {}
`
}
