package providersdk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEventoptionsGenerateEvent_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEventoptionsGenerateEventConfigCreate(),
			},
			{
				Config: testAccResourceEventoptionsGenerateEventConfigUpdate(),
			},
			{
				ResourceName:      "junos_eventoptions_generate_event.testacc_evtopts_genevent",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceEventoptionsGenerateEventConfigCreate() string {
	return `
resource "junos_eventoptions_generate_event" "testacc_evtopts_genevent" {
  name          = "testacc_evtopts_genevent#1"
  time_interval = 3600
  no_drift      = true
}
`
}

func testAccResourceEventoptionsGenerateEventConfigUpdate() string {
	return `
resource "junos_eventoptions_generate_event" "testacc_evtopts_genevent" {
  name        = "testacc_evtopts_genevent#1"
  time_of_day = "01:02:03"
}
`
}
