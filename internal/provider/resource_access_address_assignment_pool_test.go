package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAccessAddressAssignmentPool_basic(t *testing.T) {
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
					ResourceName:      "junos_access_address_assignment_pool.testacc_accessAddAssP4",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_access_address_assignment_pool.testacc_accessAddAssP6_1",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_access_address_assignment_pool.testacc_accessAddAssP6_2",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
