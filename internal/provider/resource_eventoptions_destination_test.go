package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccResourceEventoptionsDestination_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
			},
			{
				ResourceName:      "junos_eventoptions_destination.testacc_evtopts_dest",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceEventoptionsDestination_writeOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				// 1
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("junos_eventoptions_destination.testacc_evtopts_dest_wo",
						"archive_site.1.password"),
					resource.TestCheckNoResourceAttr("junos_eventoptions_destination.testacc_evtopts_dest_wo",
						"archive_site.1.password_wo"),
					resource.TestCheckResourceAttr("junos_eventoptions_destination.testacc_evtopts_dest_wo",
						"archive_site.1.password_wo_version", "1"),
				),
			},
			{
				// 2 check that the write-only password has really been sent to the device
				ConfigDirectory: config.TestStepDirectory(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.junos_config_raw.testacc_evtopts_dest_wo",
						tfjsonpath.New("config"),
						knownvalue.StringRegexp(regexp.MustCompile(
							`destinations "?testacc_evtopts_dest_wo"? archive-sites "?https://example\.fr"? password `)),
					),
				},
			},
			{
				// 3
				ConfigDirectory: config.TestStepDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_eventoptions_destination.testacc_evtopts_dest_wo",
						"archive_site.1.password_wo_version", "2"),
				),
			},
		},
	})
}
