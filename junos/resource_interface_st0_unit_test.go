package junos_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosInterfaceSt0Unit_basic(t *testing.T) {
	regexpSt0 := regexp.MustCompile("st0.")
	if os.Getenv("TESTACC_SWITCH") == "" && os.Getenv("TESTACC_ROUTER") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
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
resource junos_interface_st0_unit testacc {}
`
}
