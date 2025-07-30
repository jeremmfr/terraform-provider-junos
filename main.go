package main

import (
	"context"
	"flag"
	"log"

	"github.com/jeremmfr/terraform-provider-junos/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	ctx := context.Background()

	// Remove any date and time prefix in log package
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(ctx, provider.New, providerserver.ServeOpts{
		Address:         "registry.terraform.io/jeremmfr/junos",
		Debug:           debug,
		ProtocolVersion: 5,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
