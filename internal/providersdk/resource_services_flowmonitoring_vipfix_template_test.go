package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosServicesFlowMonitoringVIPFixTemplate_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosServicesFlowMonitoringVIPFixTemplateConfigCreate(),
				},
				{
					Config: testAccJunosServicesFlowMonitoringVIPFixTemplateConfigUpdate(),
				},
				{
					ResourceName:      "junos_services_flowmonitoring_vipfix_template.testacc_flow_vipfix_template_ipv4",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_services_flowmonitoring_vipfix_template.testacc_flow_vipfix_template_ipv6",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_services_flowmonitoring_vipfix_template.testacc_flow_vipfix_template_mpls",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testAccJunosServicesFlowMonitoringVIPFixTemplateConfigCreate() string {
	return `
resource "junos_services_flowmonitoring_vipfix_template" "testacc_flow_vipfix_template_ipv4" {
  name = "testacc_template@1"
  type = "ipv4-template"
}
resource "junos_services_flowmonitoring_vipfix_template" "testacc_flow_vipfix_template_ipv6" {
  name = "testacc_template@3"
  type = "ipv6-template"
}
resource "junos_services_flowmonitoring_vipfix_template" "testacc_flow_vipfix_template_mpls" {
  name = "testacc_template@2"
  type = "mpls-template"
}
`
}

func testAccJunosServicesFlowMonitoringVIPFixTemplateConfigUpdate() string {
	return `
resource "junos_services_flowmonitoring_vipfix_template" "testacc_flow_vipfix_template_ipv4" {
  name                         = "testacc_template@1"
  type                         = "ipv4-template"
  ip_template_export_extension = ["app-id", "flow-dir"]
  flow_active_timeout          = 60
  flow_inactive_timeout        = 30
  flow_key_flow_direction      = true
  flow_key_vlan_id             = true
  nexthop_learning_enable      = true
  observation_domain_id        = 10
  option_refresh_rate {}
  option_template_id = 2000
  template_id        = 2001
}
resource "junos_services_flowmonitoring_vipfix_template" "testacc_flow_vipfix_template_ipv6" {
  name                         = "testacc_template@3"
  type                         = "ipv6-template"
  ip_template_export_extension = ["app-id", "flow-dir"]
  flow_active_timeout          = 60
  flow_inactive_timeout        = 30
}
resource "junos_services_flowmonitoring_vipfix_template" "testacc_flow_vipfix_template_mpls" {
  name                    = "testacc_template@2"
  type                    = "mpls-template"
  flow_active_timeout     = 60
  flow_inactive_timeout   = 30
  flow_key_flow_direction = true
  flow_key_vlan_id        = true
  nexthop_learning_enable = true
  observation_domain_id   = 10
  option_refresh_rate {
    packets = 100
    seconds = 60
  }
  option_template_id      = 2002
  template_id             = 2003
  tunnel_observation_ipv4 = true
  tunnel_observation_ipv6 = true
}
`
}
