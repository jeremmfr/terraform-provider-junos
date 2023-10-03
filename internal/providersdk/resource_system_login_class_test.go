package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSystemLoginClass_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSystemLoginClassConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"name", "testacc"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"allowed_days.#", "2"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"confirm_commands.#", "1"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"permissions.#", "2"),
					),
				},
				{
					Config: testAccResourceSystemLoginClassConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"name", "testacc"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"allow_commands_regexps.#", "1"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"allow_configuration_regexps.#", "1"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"no_hidden_commands_except.#", "1"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"deny_commands_regexps.#", "1"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"deny_configuration_regexps.#", "1"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"allow_configuration_regexps.#", "1"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"allow_configuration_regexps.#", "1"),
						resource.TestCheckResourceAttr("junos_system_login_class.testacc",
							"permissions.#", "1"),
					),
				},
				{
					ResourceName:      "junos_system_login_class.testacc",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccResourceSystemLoginClassConfigCreate() string {
	return `
resource "junos_system_login_class" "testacc" {
  name                      = "testacc"
  access_start              = "08:00:00"
  access_end                = "18:00:00"
  allow_commands            = ".*"
  allow_configuration       = ".*"
  allow_hidden_commands     = true
  allowed_days              = ["sunday", "monday"]
  cli_prompt                = "prompt cli"
  configuration_breadcrumbs = true
  confirm_commands          = ["confirm commands"]
  deny_commands             = "request"
  deny_configuration        = "system"
  idle_timeout              = 120
  login_alarms              = true
  login_tip                 = true
  permissions               = ["view", "floppy"]
  security_role             = "security-administrator"
}
`
}

func testAccResourceSystemLoginClassConfigUpdate() string {
	return `
resource "junos_system_login_class" "testacc" {
  name                        = "testacc"
  access_start                = "08:00:00"
  access_end                  = "18:00:00"
  allow_commands_regexps      = [".*"]
  allow_configuration_regexps = [".*"]
  no_hidden_commands_except   = [".*"]
  deny_commands_regexps       = ["request"]
  deny_configuration_regexps  = ["system"]
  idle_timeout                = 120
  login_alarms                = true
  login_tip                   = true
  permissions                 = ["view"]
  security_role               = "security-administrator"
}
`
}
