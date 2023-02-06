package providersdk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSystemSyslogHost_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosSystemSyslogHostConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"host", "192.0.2.1"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"port", "514"),
				),
			},
			{
				Config: testAccJunosSystemSyslogHostConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"structured_data.#", "1"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"structured_data.0.brief", "true"),
				),
			},
			{
				Config: testAccJunosSystemSyslogHostConfigUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"structured_data.#", "1"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"structured_data.0.brief", "false"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"source_address", "192.0.2.2"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"allow_duplicates", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"exclude_hostname", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"explicit_priority", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"facility_override", "local3"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"log_prefix", "prefix"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"match", "match testacc"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"match_strings.#", "1"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"match_strings.0", "match testacc"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"any_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"changelog_severity", "critical"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"conflictlog_severity", "error"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"daemon_severity", "warning"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"dfc_severity", "alert"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"external_severity", "any"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"firewall_severity", "info"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"ftp_severity", "none"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"interactivecommands_severity", "notice"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"kernel_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"ntp_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"pfe_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"security_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"user_severity", "emergency"),
					resource.TestCheckResourceAttr("junos_system_syslog_host.testacc_syslogHost",
						"source_address", "192.0.2.2"),
				),
			},
			{
				ResourceName:      "junos_system_syslog_host.testacc_syslogHost",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccJunosSystemSyslogHostConfigCreate() string {
	return `
resource "junos_system_syslog_host" "testacc_syslogHost" {
  host = "192.0.2.1"
  port = 514
}
`
}

func testAccJunosSystemSyslogHostConfigUpdate() string {
	return `
resource "junos_system_syslog_host" "testacc_syslogHost" {
  host = "192.0.2.1"
  structured_data {
    brief = true
  }
}
`
}

func testAccJunosSystemSyslogHostConfigUpdate2() string {
	return `
resource "junos_system_syslog_host" "testacc_syslogHost" {
  host                         = "192.0.2.1"
  port                         = 514
  allow_duplicates             = true
  exclude_hostname             = true
  explicit_priority            = true
  facility_override            = "local3"
  log_prefix                   = "prefix"
  match                        = "match testacc"
  match_strings                = ["match testacc"]
  any_severity                 = "emergency"
  changelog_severity           = "critical"
  conflictlog_severity         = "error"
  daemon_severity              = "warning"
  dfc_severity                 = "alert"
  external_severity            = "any"
  firewall_severity            = "info"
  ftp_severity                 = "none"
  interactivecommands_severity = "notice"
  kernel_severity              = "emergency"
  ntp_severity                 = "emergency"
  pfe_severity                 = "emergency"
  security_severity            = "emergency"
  user_severity                = "emergency"
  structured_data {}
  source_address = "192.0.2.2"
}
`
}
