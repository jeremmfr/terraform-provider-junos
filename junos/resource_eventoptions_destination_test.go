package junos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosEventoptionsDestination_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosEventoptionsDestinationConfigCreate(),
			},
			{
				Config: testAccJunosEventoptionsDestinationConfigUpdate(),
			},
			{
				ResourceName:      "junos_eventoptions_destination.testacc_evtopts_dest",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccJunosEventoptionsDestinationConfigCreate() string {
	return `
resource "junos_eventoptions_destination" "testacc_evtopts_dest" {
  name = "testacc_evtopts_dest#1"
  archive_site {
    url = "https://example.com"
  }
}
`
}

func testAccJunosEventoptionsDestinationConfigUpdate() string {
	return `
resource "junos_eventoptions_destination" "testacc_evtopts_dest" {
  name = "testacc_evtopts_dest#1"
  archive_site {
    url = "https://example.com"
  }
  archive_site {
    url      = "https://example.fr"
    password = "thePassword"
  }
  transfer_delay = 10
}
`
}
