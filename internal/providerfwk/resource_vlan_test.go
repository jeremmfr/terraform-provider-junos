package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceVlan_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"description", "testacc_vlansw"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"vlan_id", "1000"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"service_id", "1000"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"l3_interface", "irb.1000"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"forward_filter_input", "testacc_vlansw"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"forward_filter_output", "testacc_vlansw"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"forward_flood_input", "testacc_vlansw"),
					),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"vlan_id_list.#", "1"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"vlan_id_list.0", "1001-1002"),
						resource.TestCheckResourceAttr("junos_vlan.testacc_vlansw",
							"private_vlan", "community"),
					),
				},
				{
					ResourceName:      "junos_vlan.testacc_vlansw",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
				{
					ConfigDirectory: config.TestStepDirectory(),
				},
			},
		})
	}
}

func TestAccResourceVlan_router(t *testing.T) {
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
			},
		})
	}
}
