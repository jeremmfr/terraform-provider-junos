package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceForwardingoptionsSamplingInstance_basic(t *testing.T) {
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
					ResourceName:      "junos_forwardingoptions_sampling_instance.testacc_sampInstance",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_forwardingoptions_sampling_instance.testacc_sampInstance2",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					ResourceName:      "junos_forwardingoptions_sampling_instance.testacc_sampInstance3",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}
