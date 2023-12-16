package providerfwk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSystemSyslogUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"allow_duplicates", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"match", "match testacc"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"any_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"changelog_severity", "critical"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"conflictlog_severity", "error"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"daemon_severity", "warning"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"dfc_severity", "alert"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"external_severity", "any"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"firewall_severity", "info"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"ftp_severity", "none"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"interactivecommands_severity", "notice"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"kernel_severity", "error"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"ntp_severity", "error"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"pfe_severity", "error"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"security_severity", "error"),
					resource.TestCheckResourceAttr("junos_system_syslog_user.testacc_syslogUser",
						"user_severity", "error"),
				),
			},
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ConfigDirectory:          config.TestStepDirectory(),
			},
			{
				ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
				ResourceName:             "junos_system_syslog_user.testacc_syslogUser",
				ImportState:              true,
				ImportStateVerify:        true,
			},
		},
	})
}
