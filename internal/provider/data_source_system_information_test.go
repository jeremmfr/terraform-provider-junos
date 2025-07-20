package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSystemInformation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.junos_system_information.test", "id"),
					resource.TestCheckResourceAttrSet("data.junos_system_information.test", "hardware_model"),
					resource.TestCheckResourceAttrSet("data.junos_system_information.test", "os_name"),
					resource.TestCheckResourceAttrSet("data.junos_system_information.test", "os_version"),
					resource.TestCheckResourceAttrSet("data.junos_system_information.test", "serial_number"),
					resource.TestCheckResourceAttrSet("data.junos_system_information.test", "cluster_node"),
				),
			},
		},
	})
}
