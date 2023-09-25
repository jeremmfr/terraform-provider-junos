package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSystemLoginUser_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSystemLoginUserConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_system_login_user.testacc",
							"name", "testacc"),
						resource.TestCheckResourceAttrSet("junos_system_login_user.testacc",
							"uid"),
						resource.TestCheckResourceAttr("junos_system_login_user.testacc",
							"authentication.#", "1"),
					),
				},
				{
					Config: testAccResourceSystemLoginUserConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_system_login_user.testacc",
							"name", "testacc"),
						resource.TestCheckResourceAttrSet("junos_system_login_user.testacc",
							"uid"),
						resource.TestCheckResourceAttr("junos_system_login_user.testacc",
							"authentication.#", "1"),
						resource.TestCheckResourceAttr("junos_system_login_user.testacc",
							"authentication.0.ssh_public_keys.#", "1"),
					),
				},
				{
					ResourceName:      "junos_system_login_user.testacc",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccResourceSystemLoginUserConfigCreate() string {
	return `
resource "junos_system_login_user" "testacc" {
  name       = "testacc"
  class      = "unauthorized"
  cli_prompt = "test cli"
  full_name  = "test name"
  authentication {
    encrypted_password = "test"
    no_public_keys     = true
  }
}
resource "junos_system_login_user" "testacc2" {
  name  = "test.acc2"
  class = "unauthorized"
}
resource "junos_system_login_user" "testacc3" {
  name  = "test.acc3"
  class = "unauthorized"
  authentication {
    plain_text_password = "test1234"
  }
}
`
}

func testAccResourceSystemLoginUserConfigUpdate() string {
	return `
resource "junos_system_login_user" "testacc" {
  name  = "testacc"
  class = "unauthorized"
  authentication {
    encrypted_password = "test"
    ssh_public_keys    = ["ssh-rsa testkey"]
  }
}
`
}
