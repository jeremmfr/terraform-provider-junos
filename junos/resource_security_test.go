package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurity_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosSecurityConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.name", "ike.log"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.files", "5"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.match", "test"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.size", "100000"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.world_readable", "true"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.flag.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.flag.0", "all"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.rate_limit", "100"),
						/*resource.TestCheckResourceAttr("junos_security.testacc_security",
						"ike_traceoptions.0.no_remote_trace", "true"),*/
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"utm.#", "1"),
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"utm.0.feature_profile_web_filtering_type", "juniper-enhanced"),
					),
				},
				{
					Config: testAccJunosSecurityConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security.testacc_security",
							"ike_traceoptions.0.file.0.no_world_readable", "true"),
					),
				},
				{
					ResourceName:      "junos_security.testacc_security",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityConfigCreate() string {
	return `
resource junos_security "testacc_security" {
  ike_traceoptions {
    file {
      name           = "ike.log"
      files          = 5
      match          = "test"
      size           = 100000
      world_readable = true
    }
    flag       = ["all"]
    rate_limit = 100
    # no_remote_trace = true
  }
  utm {
    feature_profile_web_filtering_type = "juniper-enhanced"
  }
}
`
}
func testAccJunosSecurityConfigUpdate() string {
	return `
resource junos_security "testacc_security" {
  ike_traceoptions {
    file {
      name           = "ike.log"
      files          = 5
      match          = "test"
      size           = 100000
      no_world_readable = true
    }
    flag       = ["all"]
    rate_limit = 100
    # no_remote_trace = true
  }
}
`
}