package junos

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var (
	testAccProviders = map[string]terraform.ResourceProvider{
		"junos": testAccProvider,
	}
	testAccProvider = Provider().(*schema.Provider)
)

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

// export TESTACC_SWITCH not empty for test switch options (interface mode trunk, vlan native/members)
// with switch Junos device, else it's test for all others parameters
// (interface inet, 802.3ad, routing instance, security zone/nat/ike/ipsec, etc  ).
// Some resources and parameters works on both devices, but most tested without TESTACC_SWITCH

func testAccPreCheck(t *testing.T) {
	if os.Getenv("JUNOS_HOST") == "" && os.Getenv("JUNOS_KEYFILE") == "" {
		t.Fatal("JUNOS_HOST or JUNOS_KEYFILE must be set for acceptance tests")
	}

	err := testAccProvider.Configure(terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}
