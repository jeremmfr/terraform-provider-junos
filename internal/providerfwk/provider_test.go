package providerfwk_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/providerfwk"
	"github.com/jeremmfr/terraform-provider-junos/internal/providersdk"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

var testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){ //nolint:gochecknoglobals
	"junos": testAccNewProtoV5MuxProviderServer(),
}

func testAccNewProtoV5MuxProviderServer() func() (tfprotov5.ProviderServer, error) {
	providers := []func() tfprotov5.ProviderServer{
		providerserver.NewProtocol5(providerfwk.New()),
		providersdk.Provider().GRPCProvider,
	}

	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatal(err)
	}

	return func() (tfprotov5.ProviderServer, error) { return muxServer.ProviderServer(), nil }
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
