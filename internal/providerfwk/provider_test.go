package providerfwk_test

import (
	"context"
	"log"

	"github.com/jeremmfr/terraform-provider-junos/internal/providerfwk"
	"github.com/jeremmfr/terraform-provider-junos/internal/providersdk"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

var testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){ //nolint: gochecknoglobals
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
