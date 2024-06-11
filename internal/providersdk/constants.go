package providersdk

import "github.com/jeremmfr/terraform-provider-junos/internal/junos"

const failedConvAtoiError = "failed to convert value from '%s' to integer: %w"

const (
	setRoutingInstances = junos.SetLS + junos.RoutingInstancesWS
	delRoutingInstances = junos.DeleteLS + junos.RoutingInstancesWS
)
