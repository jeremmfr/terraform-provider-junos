package junos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosSystemRadiusServer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJunosSystemRadiusServerConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"address", "192.0.2.1"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"secret", "password"),
				),
			},
			{
				Config: testAccJunosSystemRadiusServerConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"preauthentication_secret", "password"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"source_address", "192.0.2.2"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"port", "1645"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"accounting_port", "1646"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"dynamic_request_port", "3799"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"preauthentication_port", "1812"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"timeout", "10"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"accounting_timeout", "5"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"retry", "3"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"accounting_retry", "2"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"max_outstanding_requests", "1000"),
					resource.TestCheckResourceAttr("junos_system_radius_server.testacc_radiusServer",
						"routing_instance", "testacc_radiusServer"),
				),
			},
			{
				ResourceName:      "junos_system_radius_server.testacc_radiusServer",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccJunosSystemRadiusServerConfigCreate() string {
	return `
resource junos_system_radius_server testacc_radiusServer {
  address = "192.0.2.1"
  secret  = "password"
}
`
}
func testAccJunosSystemRadiusServerConfigUpdate() string {
	return `
resource junos_routing_instance testacc_radiusServer {
  name = "testacc_radiusServer"
}
resource junos_system_radius_server testacc_radiusServer {
  address                  = "192.0.2.1"
  secret                   = "password"
  preauthentication_secret = "password"
  source_address           = "192.0.2.2"
  port                     = 1645
  accounting_port          = 1646
  dynamic_request_port     = 3799
  preauthentication_port   = 1812
  timeout                  = 10
  accounting_timeout       = 5
  retry                    = 3
  accounting_retry         = 2
  max_outstanding_requests = 1000
  routing_instance         = junos_routing_instance.testacc_radiusServer.name
}
`
}
