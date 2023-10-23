package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceEvpn_basic(t *testing.T) {
	if os.Getenv("TESTACC_ROUTER") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ResourceName:      "junos_evpn.testacc_evpn_default",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_evpn.testacc_evpn_ri1",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}
