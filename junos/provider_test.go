package junos_test

import (
	"context"
	"os"
	"terraform-provider-junos/junos"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testAccProviders = map[string]*schema.Provider{
		"junos": testAccProvider,
	}
	testAccProvider = junos.Provider()
)

const defaultInterfaceTestAcc = "ge-0/0/3"

func TestProvider(t *testing.T) {
	if err := junos.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = junos.Provider()
}

// export TESTACC_SWITCH not empty for test switch options (interface mode trunk, vlan native/members)
// with switch Junos device, else it's test for all others parameters
// (interface inet, 802.3ad, routing instance, security zone/nat/ike/ipsec, etc  ).
// Few resources and parameters works on both devices, but most tested without TESTACC_SWITCH

func testAccPreCheck(t *testing.T) {
	if os.Getenv("JUNOS_HOST") == "" && os.Getenv("JUNOS_KEYFILE") == "" {
		t.Fatal("JUNOS_HOST must be set for acceptance tests")
	}
	if os.Getenv("JUNOS_KEYFILE") == "" && os.Getenv("JUNOS_PASSWORD") == "" {
		t.Fatal("JUNOS_KEYFILE or JUNOS_PASSWORD must be set for acceptance tests")
	}

	err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}
