package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSecurityLogStream_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config:             testAccJunosSecurityLogStreamConfigPreCreate(),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: testAccJunosSecurityLogStreamConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"category.#", "1"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"category.0", "idp"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"format", "syslog"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"host.#", "1"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"host.0.ip_address", "192.0.2.1"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"host.0.port", "514"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"host.0.routing_instance", "testacclogstream"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"rate_limit", "50"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"severity", "error"),
					),
				},
				{
					Config: testAccJunosSecurityLogStreamConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"file.#", "1"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"file.0.name", "test"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"file.0.allow_duplicates", "true"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"file.0.size", "3"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"file.0.rotation", "3"),
						resource.TestCheckResourceAttr("junos_security_log_stream.testacc_logstream",
							"filter_threat_attack", "true"),
					),
				},
				{
					ResourceName:      "junos_security_log_stream.testacc_logstream",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosSecurityLogStreamConfigPreCreate() string {
	return `
resource junos_security "security" {
  log {
    source_address = "192.0.2.2"
  }
}
`
}

func testAccJunosSecurityLogStreamConfigCreate() string {
	return `
resource junos_routing_instance "testacc_logstream" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacclogstream"
}
resource junos_security_log_stream "testacc_logstream" {
  name     = "testacc_logstream"
  category = ["idp"]
  format   = "syslog"
  host {
    ip_address       = "192.0.2.1"
    port             = 514
    routing_instance = junos_routing_instance.testacc_logstream.name
  }
  rate_limit = 50
  severity   = "error"
}
`
}

func testAccJunosSecurityLogStreamConfigUpdate() string {
	return `
resource junos_security_log_stream "testacc_logstream" {
  name = "testacc_logstream"
  file {
    name             = "test"
    allow_duplicates = true
    size             = 3
    rotation         = 3
  }
  filter_threat_attack = true
}
`
}
