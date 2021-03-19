package junos_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testaccNullCommitFile = "/tmp/testacc/terraform-provider-junos/null-commit-file"

func TestAccJunosNullCommitFile_basic(t *testing.T) {
	var testaccInterface string
	if os.Getenv("TESTACC_INTERFACE") != "" {
		testaccInterface = os.Getenv("TESTACC_INTERFACE")
	} else {
		testaccInterface = defaultInterfaceTestAcc
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
	return `
resource junos_interface_physical testacc_nullcommitfile {
  name         = "` + interFace + `"
  description  = "testacc_null"
  vlan_tagging = true
}
`
}
func testAccJunosNullCommitFileCreate(interFace string) string {
	return `
resource junos_interface_physical testacc_nullcommitfile {
  name         = "` + interFace + `"
  description  = "testacc_null"
  vlan_tagging = true
}
resource "local_file" "hostname" {
  content  = "set interfaces ` + interFace + ` description testacc_nullfile"
  filename = "` + testaccNullCommitFile + `"
}
resource junos_null_commit_file "testacc_nullcommitfile" {
  filename                = local_file.hostname.filename
  clear_file_after_commit = true
}
`
}
func testAccJunosNullCommitFileRead(interFace string) string {
	return `
resource junos_interface_physical testacc_nullcommitfile {
  name         = "` + interFace + `"
  description  = "testacc_null"
  vlan_tagging = true
}
data junos_interface_physical testacc_nullcommitfile {
  config_interface = "` + interFace + `"
}
`
}
