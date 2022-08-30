package main

import (
	"log"

	"github.com/jeremmfr/terraform-provider-junos/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	// Remove any date and time prefix in log package
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: junos.Provider,
	})
}
