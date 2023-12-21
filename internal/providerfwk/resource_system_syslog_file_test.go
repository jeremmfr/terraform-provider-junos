package providerfwk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSystemSyslogFile_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"filename", "testacc"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"allow_duplicates", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"explicit_priority", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"match", "match testacc"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"match_strings.#", "1"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"match_strings.0", "match testacc"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"any_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"changelog_severity", "critical"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"conflictlog_severity", "error"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"daemon_severity", "warning"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"dfc_severity", "alert"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"external_severity", "any"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"firewall_severity", "info"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"ftp_severity", "none"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"interactivecommands_severity", "notice"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"kernel_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"ntp_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"pfe_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"security_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"user_severity", "emergency"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("junos_system_syslog_file.testacc_syslogFile",
						"structured_data.%"),
					resource.TestCheckResourceAttrSet("junos_system_syslog_file.testacc_syslogFile",
						"archive.%"),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"structured_data.brief", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.binary_data", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.world_readable", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.size", "1073741823"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.files", "5"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.transfer_interval", "5"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.sites.#", "1"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.sites.0.url", "192.0.2.1"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.sites.0.password", "password"),
				),
			},
			{
				ResourceName:      "junos_system_syslog_file.testacc_syslogFile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
