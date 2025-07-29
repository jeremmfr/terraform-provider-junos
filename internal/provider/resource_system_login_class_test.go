package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSystemLoginClass_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
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
				ConfigDirectory: config.TestStepDirectory(),
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
