package junos_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testaccNullCommitFile = "/tmp/testacc/terraform-provider-junos/null-commit-file"

// export TESTACC_INTERFACE=<inteface> for choose interface available else it's ge-0/0/3.
func TestAccJunosNullCommitFile_basic(t *testing.T) {
	testaccInterface := defaultInterfaceTestAcc
	if iface := os.Getenv("TESTACC_INTERFACE"); iface != "" {
		testaccInterface = iface
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"local": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccJunosNullCommitFilePreCreate(testaccInterface),
			},
			{
				Config:             testAccJunosNullCommitFileCreate(testaccInterface),
				ExpectNonEmptyPlan: true,
			},
			{
				Config:   testAccJunosNullCommitFileRead(testaccInterface),
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

func testAccJunosNullCommitFilePreCreate(interFace string) string {
	return fmt.Sprintf(`
resource "junos_interface_physical" "testacc_nullcommitfile" {
  name         = "%s"
  description  = "testacc_null"
  vlan_tagging = true
}
`, interFace)
}

func testAccJunosNullCommitFileCreate(interFace string) string {
	return fmt.Sprintf(`
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

func testAccJunosNullCommitFileRead(interFace string) string {
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
