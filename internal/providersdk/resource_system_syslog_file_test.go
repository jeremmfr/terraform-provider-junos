package providersdk_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSystemSyslogFile_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosSystemSyslogFileConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"filename", "testacc"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"structured_data.#", "0"),
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
				Config: testAccJunosSystemSyslogFileConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"structured_data.#", "1"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"structured_data.0.brief", "false"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.#", "1"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.sites.#", "0"),
				),
			},
			{
				Config: testAccJunosSystemSyslogFileConfigUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"structured_data.#", "1"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"structured_data.0.brief", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.#", "1"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.binary_data", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.world_readable", "true"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.size", "1073741823"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.files", "5"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.transfer_interval", "5"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.sites.#", "2"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.sites.0.url", "example.fr"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.sites.1.url", "example.com"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.sites.1.password", "password"),
					resource.TestCheckResourceAttr("junos_system_syslog_file.testacc_syslogFile",
						"archive.0.sites.1.routing_instance", "testacc_syslogFile"),
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

func testAccJunosSystemSyslogFileConfigCreate() string {
	return `
resource "junos_system_syslog_file" "testacc_syslogFile" {
  filename                     = "testacc"
  allow_duplicates             = true
  explicit_priority            = true
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
}
`
}

func testAccJunosSystemSyslogFileConfigUpdate() string {
	return `
resource "junos_system_syslog_file" "testacc_syslogFile" {
  filename                     = "testacc"
  allow_duplicates             = true
  match                        = "match testacc"
  any_severity                 = "emergency"
  changelog_severity           = "critical"
  conflictlog_severity         = "error"
  daemon_severity              = "warning"
  dfc_severity                 = "alert"
  external_severity            = "any"
  firewall_severity            = "info"
  ftp_severity                 = "none"
  interactivecommands_severity = "notice"
  kernel_severity              = "error"
  ntp_severity                 = "error"
  pfe_severity                 = "error"
  security_severity            = "error"
  user_severity                = "error"
  structured_data {}
  archive {}
}
resource "junos_system_syslog_file" "testacc_syslogFile2" {
  filename          = "testacc2"
  explicit_priority = true
}
`
}

func testAccJunosSystemSyslogFileConfigUpdate2() string {
	return `
resource "junos_routing_instance" "testacc_syslogFile" {
  name = "testacc_syslogFile"
}
resource "junos_system_syslog_file" "testacc_syslogFile" {
  filename                     = "testacc"
  allow_duplicates             = true
  match                        = "match testacc"
  any_severity                 = "emergency"
  changelog_severity           = "critical"
  conflictlog_severity         = "error"
  daemon_severity              = "warning"
  dfc_severity                 = "alert"
  external_severity            = "any"
  firewall_severity            = "info"
  ftp_severity                 = "none"
  interactivecommands_severity = "notice"
  kernel_severity              = "error"
  ntp_severity                 = "error"
  pfe_severity                 = "error"
  security_severity            = "error"
  user_severity                = "error"
  structured_data {
    brief = true
  }
  archive {
    binary_data       = true
    world_readable    = true
    size              = 1073741823
    files             = 5
    transfer_interval = 5
    sites {
      url = "example.fr"
    }
    sites {
      url              = "example.com"
      password         = "password"
      routing_instance = junos_routing_instance.testacc_syslogFile.name
    }
  }
}
`
}
