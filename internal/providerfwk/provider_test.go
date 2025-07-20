package providerfwk_test

import (
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/providerfwk"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

var testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){ //nolint:gochecknoglobals
	"junos": providerserver.NewProtocol5WithError(providerfwk.New()),
}

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
}

func testAccUpgradeStatePrecheck(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_CLI_CONFIG_FILE") != "" {
		t.Fatal("can't test state upgrade with TF_CLI_CONFIG_FILE env variable")
	}
}
