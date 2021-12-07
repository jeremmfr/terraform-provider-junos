package junos_test

import (
	"context"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testAccProviders = map[string]*schema.Provider{ // nolint: gochecknoglobals
		"junos": testAccProvider,
	}
	testAccProvider = junos.Provider() // nolint: gochecknoglobals
)

const (
	defaultInterfaceTestAcc       = "ge-0/0/3"
	defaultInterfaceTestAcc2      = "ge-0/0/4"
	defaultInterfaceSwitchTestAcc = "xe-0/0/3"
)

func TestProvider(t *testing.T) {
	if err := junos.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = junos.Provider()
}

// export TESTACC_SWITCH not empty to test specific switch options
// export TESTACC_ROUTER not empty to test specific router options
// export TESTACC_SRX not empty to test specific SRX options
// export TESTACC_DEPRECATED not empty to launch testacc on deprecated resources

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("JUNOS_HOST") == "" {
		t.Fatal("JUNOS_HOST must be set for acceptance tests")
	}
	if os.Getenv("JUNOS_KEYFILE") == "" && os.Getenv("JUNOS_PASSWORD") == "" && os.Getenv("SSH_AUTH_SOCK") == "" &&
		os.Getenv("JUNOS_KEYPEM") == "" {
		t.Fatal("JUNOS_KEYPEM, JUNOS_KEYFILE, SSH_AUTH_SOCK or JUNOS_PASSWORD must be set for acceptance tests")
	}
	if os.Getenv("JUNOS_FAKECREATE_SETFILE") != "" {
		t.Fatal("can't run testacc with JUNOS_FAKECREATE_SETFILE")
	}

	if err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil)); err != nil {
		t.Fatal(err)
	}
}
