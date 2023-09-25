package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceServicesFlowMonitoringV9Template_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ResourceName:      "junos_services_flowmonitoring_v9_template.testacc_flow_v9_template_ipv4",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_services_flowmonitoring_v9_template.testacc_flow_v9_template_ipv6",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_services_flowmonitoring_v9_template.testacc_flow_v9_template_mpls",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
