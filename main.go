package main

import (
	"context"
	"log"

	"github.com/jeremmfr/terraform-provider-junos/internal/providerfwk"
	"github.com/jeremmfr/terraform-provider-junos/internal/providersdk"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func main() {
	ctx := context.Background()

	// Remove any date and time prefix in log package
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	providers := []func() tfprotov5.ProviderServer{
		providerserver.NewProtocol5(providerfwk.New()),
		providersdk.Provider().GRPCProvider,
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	err = tf5server.Serve(
		"registry.terraform.io/jeremmfr/junos",
		muxServer.ProviderServer,
	)
	if err != nil {
		log.Fatal(err)
	}
}
