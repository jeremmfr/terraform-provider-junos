package providersdk_test

import (
	"context"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/providersdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testAccProviders = map[string]*schema.Provider{ //nolint: gochecknoglobals
		"junos": testAccProvider,
	}
	testAccProvider = providersdk.Provider() //nolint: gochecknoglobals
)

const (
	defaultInterfaceTestAcc        = "ge-0/0/3"
	defaultInterfaceTestAcc2       = "ge-0/0/4"
	defaultInterfaceSwitchTestAcc  = "xe-0/0/3"
	defaultInterfaceSwitchTestAcc2 = "xe-0/0/4"
)

func TestProvider(t *testing.T) {
	if err := providersdk.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = providersdk.Provider()
}

// export TESTACC_SWITCH not empty to test specific switch options
// export TESTACC_ROUTER not empty to test specific router options
// export TESTACC_SRX not empty to test specific SRX options
// export TESTACC_DEPRECATED not empty to launch testacc on deprecated resources

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv(junos.EnvHost) == "" {
		t.Fatal(junos.EnvHost + " must be set for acceptance tests")
	}
	if os.Getenv(junos.EnvKeyFile) == "" && os.Getenv(junos.EnvPassword) == "" && os.Getenv("SSH_AUTH_SOCK") == "" &&
		os.Getenv(junos.EnvKeyPem) == "" {
		t.Fatal(junos.EnvKeyPem + ", " + junos.EnvKeyFile + ", SSH_AUTH_SOCK or " + junos.EnvPassword +
			" must be set for acceptance tests")
	}
	if os.Getenv(junos.EnvFakecreateSetfile) != "" {
		t.Fatal("can't run testacc with " + junos.EnvFakecreateSetfile)
	}

	if err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil)); err != nil {
		t.Fatal(err)
	}
}
