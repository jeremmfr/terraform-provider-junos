package junos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccJunosApplication_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosApplicationConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_application.testacc_app", "protocol", "tcp"),
					resource.TestCheckResourceAttr("junos_application.testacc_app", "destination_port", "22"),
				),
			},
			{
				Config: testAccJunosApplicationConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_application.testacc_app", "protocol", "tcp"),
					resource.TestCheckResourceAttr("junos_application.testacc_app", "destination_port", "22"),
					resource.TestCheckResourceAttr("junos_application.testacc_app", "source_port", "1024-65535"),
				),
			},
			{
				ResourceName:      "junos_application.testacc_app",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccJunosApplicationConfigCreate() string {
	return fmt.Sprintf(`
resource "junos_application" "testacc_app" {
  name = "testacc_app"
  protocol = "tcp"
  destination_port = 22
}
`)
}
func testAccJunosApplicationConfigUpdate() string {
	return fmt.Sprintf(`
resource "junos_application" "testacc_app" {
  name = "testacc_app"
  protocol = "tcp"
  destination_port = "22"
  source_port = "1024-65535"
}
`)
}
