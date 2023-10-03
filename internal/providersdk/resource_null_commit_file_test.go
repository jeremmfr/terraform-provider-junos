package providersdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testaccNullCommitFile = "/tmp/testacc/terraform-provider-junos/null-commit-file"

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccResourceNullCommitFile_basic(t *testing.T) {
	testaccInterface := junos.DefaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"local": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNullCommitFileFakeCreate(testaccInterface),
			},
			{
				Config: testAccResourceNullCommitFileFakeUpdate(testaccInterface),
			},
			{
				Config: testAccResourceNullCommitFilePreCustom(testaccInterface),
			},
			{
				Config:             testAccResourceNullCommitFileCustom(testaccInterface),
				ExpectNonEmptyPlan: true,
			},
			{
				Config:   testAccResourceNullCommitFileRead(testaccInterface),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("junos_interface_physical.testacc_nullcommitfile",
						"description", "testacc_nullfile"),
					resource.TestCheckResourceAttr("data.junos_interface_physical.testacc_nullcommitfile",
						"description", "testacc_nullfile"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResourceNullCommitFileFakeCreate(interFace string) string {
	return fmt.Sprintf(`
provider "junos" {
  alias                    = "fake"
  fake_create_with_setfile = "%s"
}
resource "junos_interface_physical" "testacc_nullcommitfile" {
  provider     = junos.fake
  name         = "%s"
  description  = "testacc_fakecreate"
  vlan_tagging = true
}
resource "junos_null_commit_file" "setfile" {
  provider = junos.fake
  depends_on = [
    junos_interface_physical.testacc_nullcommitfile
  ]
  filename                = "%s"
  clear_file_after_commit = true
}
`, testaccNullCommitFile, interFace, testaccNullCommitFile)
}

func testAccResourceNullCommitFileFakeUpdate(interFace string) string {
	return fmt.Sprintf(`
provider "junos" {
  alias                    = "fake"
  fake_create_with_setfile = "%s"
  fake_update_also         = true
}
resource "junos_interface_physical" "testacc_nullcommitfile" {
  provider     = junos.fake
  name         = "%s"
  description  = "testacc_fakeupdate"
  vlan_tagging = true
}
resource "junos_null_commit_file" "setfile2" {
  provider = junos.fake
  depends_on = [
    junos_interface_physical.testacc_nullcommitfile
  ]
  filename                = "%s"
  clear_file_after_commit = true
}
`, testaccNullCommitFile, interFace, testaccNullCommitFile)
}

func testAccResourceNullCommitFilePreCustom(interFace string) string {
	return fmt.Sprintf(`
provider "junos" {
  alias = "fake"
}
resource "junos_interface_physical" "testacc_nullcommitfile" {
  provider     = junos.fake
  name         = "%s"
  description  = "testacc_null"
  vlan_tagging = true
}
`, interFace)
}

func testAccResourceNullCommitFileCustom(interFace string) string {
	return fmt.Sprintf(`
provider "junos" {
  alias = "fake"
}
resource "junos_interface_physical" "testacc_nullcommitfile" {
  name         = "%s"
  description  = "testacc_null"
  vlan_tagging = true
}
resource "local_file" "hostname" {
  content  = "set interfaces %s description testacc_nullfile"
  filename = "%s"
}
resource "junos_null_commit_file" "testacc_nullcommitfile" {
  filename                = local_file.hostname.filename
  clear_file_after_commit = true
}
`, interFace, interFace, testaccNullCommitFile)
}

func testAccResourceNullCommitFileRead(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_nullcommitfile" {
  name         = "%s"
  description  = "testacc_null"
  vlan_tagging = true
}
data "junos_interface_physical" "testacc_nullcommitfile" {
  config_interface = "%s"
}
`, interFace, interFace)
}
